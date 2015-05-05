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
		// Website creates a new Website pointer from an ID or code including all of its
		// groups and all related stores. It panics when the integrity is incorrect.
		// If ID and code are available then the non-empty code has precedence.
		Website(Retriever) (*Website, error)
		// Websites creates a slice containing all pointers to Websites with its associated
		// groups and stores. It panics when the integrity is incorrect.
		Websites() (WebsiteSlice, error)
		// Group creates a new Group which contains all related stores and its website.
		// Only the argument ID can be used to get a specific Group.
		Group(Retriever) (*Group, error)
		// Groups creates a slice containing all pointers to Groups with its associated
		// stores and websites. It panics when the integrity is incorrect.
		Groups() (GroupSlice, error)
		// Store creates a new Store containing its group and its website.
		// If ID and code are available then the non-empty code has precedence.
		Store(Retriever) (*Store, error)
		// Stores creates a new store slice. Can return an error when the website or
		// the group cannot be found.
		Stores() (StoreSlice, error)
		// DefaultStoreView traverses through the websites to find the default website and gets
		// the default group which has the default store id assigned to. Only one website can be the default one.
		DefaultStoreView() (*Store, error)
		// ActiveStore returns a new Store with all its Websites and Groups but only if the Store
		// is marked as active. Argument can be an ID or a Code. Returns nil if Store not found or inactive.
		// No need here to return an error.
		ActiveStore(Retriever) (*Store, error)
	}

	// StorageMutator allows changes to the internal stored slices.
	StorageMutator interface {
		// ReInit reloads the websites, groups and stores from the database.
		ReInit(dbr.SessionRunner) error
	}

	// Storage contains a mutex and the raw slices from the database.
	Storage struct {
		mu       sync.RWMutex
		websites TableWebsiteSlice
		groups   TableGroupSlice
		stores   TableStoreSlice
	}

	// Retriever implements how to get the ID. If Retriever implements CodeRetriever
	// then CodeRetriever has precedence. ID can be any of the website, group or store IDs.
	Retriever interface {
		ID() int64
	}
	// CodeRetriever implements how to get an object by Code which can be website or store code.
	// Groups doesn't have codes.
	CodeRetriever interface {
		Code() string
	}
	// ID is convenience helper to satisfy the interface IDRetriever.
	ID int64
	// Code is convenience helper to satisfy the interface CodeRetriever and IDRetriever.
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
	return &Storage{
		mu:       sync.RWMutex{},
		websites: tws,
		groups:   tgs,
		stores:   tss,
	}
}

// website returns a TableWebsite by using either id or code to find it. If id and code are
// available then the non-empty code has precedence.
func (st *Storage) website(r Retriever) (*TableWebsite, error) {
	if r == nil {
		return nil, ErrWebsiteNotFound
	}
	if c, ok := r.(CodeRetriever); ok && c.Code() != "" {
		return st.websites.FindByCode(c.Code())
	}
	return st.websites.FindByID(r.ID())
}

// Website creates a new Website according to the interface definition.
func (st *Storage) Website(r Retriever) (*Website, error) {
	w, err := st.website(r)
	if err != nil {
		return nil, err
	}
	return NewWebsite(w).SetGroupsStores(st.groups, st.stores), nil
}

// Websites creates a slice of Website pointers according to the interface definition.
func (st *Storage) Websites() (WebsiteSlice, error) {
	websites := make(WebsiteSlice, len(st.websites), len(st.websites))
	for i, w := range st.websites {
		websites[i] = NewWebsite(w).SetGroupsStores(st.groups, st.stores)
	}
	return websites, nil
}

// group returns a TableGroup by using a group id as argument. If no argument or more than
// one has been supplied it returns an error.
func (st *Storage) group(r Retriever) (*TableGroup, error) {
	if r == nil {
		return nil, ErrGroupNotFound
	}
	return st.groups.FindByID(r.ID())
}

// Group creates a new Group which contains all related stores and its website according to the
// interface definition.
func (st *Storage) Group(id Retriever) (*Group, error) {
	g, err := st.group(id)
	if err != nil {
		return nil, err
	}

	w, err := st.website(ID(g.WebsiteID))
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
		w, err := st.website(ID(g.WebsiteID))
		if err != nil {
			return nil, errgo.Mask(err)
		}
		groups[i] = NewGroup(g, w).SetStores(st.stores, nil)
	}
	return groups, nil
}

// store returns a TableStore by an id or code.
// The non-empty code has precedence if available.
func (st *Storage) store(r Retriever) (*TableStore, error) {
	if r == nil {
		return nil, ErrStoreNotFound
	}
	if c, ok := r.(CodeRetriever); ok && c.Code() != "" {
		return st.stores.FindByCode(c.Code())
	}
	return st.stores.FindByID(r.ID())
}

// Store creates a new Store which contains the the store, its group and website
// according to the interface definition.
func (st *Storage) Store(r Retriever) (*Store, error) {
	s, err := st.store(r)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	w, err := st.website(ID(s.WebsiteID))
	if err != nil {
		return nil, errgo.Mask(err)
	}
	g, err := st.group(ID(s.GroupID))
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return NewStore(w, g, s), nil
}

// ActiveStore returns a new Store with all its Websites and Groups but only if the Store
// is marked as active. Argument can be an ID or a Code. Returns nil if Store not found or inactive.
// No need here to return an error.
func (st *Storage) ActiveStore(r Retriever) (*Store, error) {
	s, err := st.Store(r)
	if err != nil {
		return nil, err
	}
	if s.Data().IsActive {
		s.Website().SetGroupsStores(st.groups, st.stores)
		return s, nil
	}
	return nil, ErrStoreNotActive
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
			return st.Store(ID(g.DefaultStoreID))
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
