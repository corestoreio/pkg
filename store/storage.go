package store

import (
	"errors"
	"sync"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

type (
	// Storager implements the requirements to get new websites, groups and store views.
	// This interface is used in the StoreManager
	Storager interface {
		Website(...Retriever) (*Website, error)
		Websites() (WebsiteSlice, error)
		Group(Retriever) (*Group, error)
		Groups() (GroupSlice, error)
		Store(...Retriever) (*Store, error)
		Stores() (StoreSlice, error)
		DefaultStoreView() (*Store, error)
	}

	StorageMutator interface {
		ReInit(dbr.SessionRunner) error
	}

	// Storage private type which holds the slices and maps
	Storage struct {
		mu       sync.RWMutex
		websites TableWebsiteSlice
		groups   TableGroupSlice
		stores   TableStoreSlice
	}

	// Retriever implements how to get the objects ID. If Retriever implements CodeRetriever
	// then CodeRetriever has precedence.
	Retriever interface {
		ID() int64
	}
	// CodeRetriever implements how to get the objects Code which can be website or store code.
	CodeRetriever interface {
		Code() string
	}
	// ID is convenience helper to satisfy the interface IDRetriever
	ID int64
	// Code is convenience helper to satisfy the interface CodeRetriever
	Code string
)

// check if interface has been implemented
var _ Storager = (*Storage)(nil)
var _ StorageMutator = (*Storage)(nil)

// ID is convenience helper to satisfy the interface Retriever
func (i ID) ID() int64 { return int64(i) }

// ID is a noop method receiver to satisfy the interface Retriever
func (c Code) ID() int64 { return int64(0) }

// Code is convenience helper to satisfy the interface CodeRetriever
func (c Code) Code() string { return string(c) }

// NewStorage creates a new storage object from three slice types
func NewStorage(tws TableWebsiteSlice, tgs TableGroupSlice, tss TableStoreSlice) *Storage {
	// maybe we can totally remove the maps and just rely on the for loop to find an entity.
	return &Storage{
		mu:       sync.RWMutex{},
		websites: tws,
		groups:   tgs,
		stores:   tss,
	}
}

// website returns a TableWebsite by using either id or code to find it.
func (st *Storage) website(r ...Retriever) (*TableWebsite, error) {
	if r == nil || len(r) > 1 {
		return nil, ErrWebsiteNotFound
	}
	if c, ok := r[0].(CodeRetriever); ok && c.Code() != "" {
		return st.websites.FindByCode(c.Code())
	}
	return st.websites.FindByID(r[0].ID())
}

// Website creates a new Website from an ID or code including all its groups and
// all related stores. It panics when the integrity is incorrect. If both arguments
// are set then the first one will take effect.
func (st *Storage) Website(r ...Retriever) (*Website, error) {
	w, err := st.website(r...)
	if err != nil {
		return nil, err
	}
	return NewWebsite(w).SetGroupsStores(st.groups, st.stores), nil
}

// Websites creates a slice contains all websites with its associated groups and stores.
// It panics when the integrity is incorrect.
func (st *Storage) Websites() (WebsiteSlice, error) {
	websites := make(WebsiteSlice, len(st.websites), len(st.websites))
	for i, w := range st.websites {
		websites[i] = NewWebsite(w).SetGroupsStores(st.groups, st.stores)
	}
	return websites, nil
}

// group returns a TableGroup by using a group id as argument.
func (st *Storage) group(id Retriever) (*TableGroup, error) {
	switch {
	case id != nil:
		return st.groups.FindByID(id.ID())
	default:
		return nil, ErrGroupNotFound
	}
}

// Group creates a new Group which contains all related stores and its website
func (st *Storage) Group(id Retriever) (*Group, error) {
	g, err := st.group(id)
	if err != nil {
		return nil, err
	}

	w, err := st.website(ID(g.WebsiteID), nil)
	if err != nil {
		return nil, err
	}
	return NewGroup(g, w).SetStores(st.stores, nil), nil
}

// Groups creates a new group slice containing its website all related stores.
// May panic when a website pointer is nil.
func (st *Storage) Groups() (GroupSlice, error) {
	groups := make(GroupSlice, len(st.groups), len(st.groups))
	for i, g := range st.groups {
		w, err := st.website(ID(g.WebsiteID), nil)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		groups[i] = NewGroup(g, w).SetStores(st.stores, nil)
	}
	return groups, nil
}

// store returns a TableStore by an id or code. Only one of the args can be nil.
func (st *Storage) store(r ...Retriever) (*TableStore, error) {
	switch {
	case id != nil:
		return st.stores.FindByID(id.ID())
	case c != nil:
		return st.stores.FindByCode(c.Code())
	default:
		return nil, ErrStoreNotFound
	}
}

// Store creates a new Store which contains the current store, its group and website.
// One of the arguments can be nil.
func (st *Storage) Store(r ...Retriever) (*Store, error) {
	s, err := st.store(r)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	w, err := st.website(ID(s.WebsiteID), nil)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	g, err := st.group(ID(s.GroupID))
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return NewStore(w, g, s), nil
}

// Stores creates a new store slice. Can return an error when the website or
// the group cannot be found.
func (st *Storage) Stores() (StoreSlice, error) {
	stores := make(StoreSlice, len(st.stores), len(st.stores))
	for i, s := range st.stores {
		w, err := st.websites.FindByID(s.WebsiteID)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		g, err := st.groups.FindByID(s.GroupID)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		stores[i] = NewStore(w, g, s)
	}
	return stores, nil
}

// DefaultStoreView traverses through the websites to find the default website and gets
// the default group which has the default store id assigned to. Only one website can be the default one.
func (st *Storage) DefaultStoreView() (*Store, error) {
	for _, website := range st.websites {
		if website.IsDefault.Bool && website.IsDefault.Valid {
			g, err := st.group(ID(website.DefaultGroupID))
			if err != nil {
				return nil, err
			}
			return st.Store(ID(g.DefaultStoreID), nil)
		}
	}
	return nil, ErrStoreNotFound
}

// ReInit reloads all websites, groups and stores from the database @todo
func (st *Storage) ReInit(dbrSess dbr.SessionRunner) error {
	st.mu.Lock()
	// fetch from DB, clear slices, pointers, etc, check for mem leak ;-) ...
	defer st.mu.Unlock()
	return errors.New("@todo")
}
