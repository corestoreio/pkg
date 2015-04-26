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
		Website(id IDRetriever, c CodeRetriever) (*Website, error)
		Websites() (WebsiteSlice, error)
		Group(id IDRetriever) (*Group, error)
		Groups() (GroupSlice, error)
		Store(id IDRetriever, c CodeRetriever) (*Store, error)
		Stores() (StoreSlice, error)
		DefaultStoreView() (*Store, error)
	}

	StorageMutator interface {
		ReInit(dbr.SessionRunner) error
		Persists(dbr.SessionRunner) error
	}

	// Storage private type which holds the slices and maps
	Storage struct {
		mu        sync.RWMutex
		websites  TableWebsiteSlice
		websiteIM *indexMap
		groups    TableGroupSlice
		groupIM   *indexMap
		stores    TableStoreSlice
		storeIM   *indexMap
	}
	// IDRetriever implements how to get the objects ID
	IDRetriever interface {
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

// ID is convenience helper to satisfy the interface IDRetriever
func (i ID) ID() int64 { return int64(i) }

// Code is convenience helper to satisfy the interface CodeRetriever
func (c Code) Code() string { return string(c) }

// NewStorage creates a new storage object from three slice types
func NewStorage(tws TableWebsiteSlice, tgs TableGroupSlice, tss TableStoreSlice) *Storage {
	// maybe we can totally remove the maps and just rely on the for loop to find an entity.
	return &Storage{
		mu:        sync.RWMutex{},
		websites:  tws,
		websiteIM: newIndexMap(tws),
		groups:    tgs,
		groupIM:   newIndexMap(tgs),
		stores:    tss,
		storeIM:   newIndexMap(tss),
	}
}

// website returns a TableWebsite by using either id or code to find it.
func (st *Storage) website(id IDRetriever, c CodeRetriever) (*TableWebsite, error) {
	var idx = -1
	switch {
	case id != nil:
		if i, ok := st.websiteIM.id[id.ID()]; ok {
			idx = i
		}
		break
	case c != nil:
		if i, ok := st.websiteIM.code[c.Code()]; ok {
			idx = i
		}
		break
	default:
		return nil, ErrWebsiteNotFound
	}
	if idx < 0 {
		return nil, ErrWebsiteNotFound
	}
	return st.websites[idx], nil
}

// Website creates a new Website from an ID or code including all its groups and
// all related stores. It panics when the integrity is incorrect. If both arguments
// are set then the first one will take effect.
func (st *Storage) Website(id IDRetriever, c CodeRetriever) (*Website, error) {
	w, err := st.website(id, c)
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
func (st *Storage) group(id IDRetriever) (*TableGroup, error) {
	var idx = -1
	switch {
	case id != nil:
		if i, ok := st.groupIM.id[id.ID()]; ok {
			idx = i
		}
		break
	default:
		return nil, ErrGroupNotFound
	}
	if idx < 0 {
		return nil, ErrGroupNotFound
	}
	return st.groups[idx], nil
}

// Group creates a new Group which contains all related stores and its website
func (st *Storage) Group(id IDRetriever) (*Group, error) {
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
func (st *Storage) store(id IDRetriever, c CodeRetriever) (*TableStore, error) {
	var idx = -1
	switch {
	case id != nil:
		if i, ok := st.storeIM.id[id.ID()]; ok {
			idx = i
		}
		break
	case c != nil:
		if i, ok := st.storeIM.code[c.Code()]; ok {
			idx = i
		}
		break
	default:
		return nil, ErrStoreNotFound
	}
	if idx < 0 {
		return nil, ErrStoreNotFound
	}
	return st.stores[idx], nil
}

// Store creates a new Store which contains the current store, its group and website.
// One of the arguments can be nil.
func (st *Storage) Store(id IDRetriever, c CodeRetriever) (*Store, error) {
	s, err := st.store(id, c)
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
	// fetch from DB, clear internal maps, pointers, etc, check for mem leak ;-) ...
	defer st.mu.Unlock()
	return errors.New("@todo")
}

// Persists writes all websites, groups and stores to the database @todo
func (st *Storage) Persists(dbrSess dbr.SessionRunner) error {
	st.mu.RLock()
	// save to DB in a transaction
	defer st.mu.RUnlock()
	return errors.New("@todo")
}

/*
	INTERNAL @todo investigate if we maybe can remove the maps and just rely on the for range loops
				   to find a website, group or store.
*/

// indexMap for faster access to the website, store group, store structs instead of
// iterating over the slices.
type indexMap struct {
	sync.RWMutex
	id   map[int64]int  // always initialized
	code map[string]int // lazy initialization
}

func newIndexMap(s interface{}) *indexMap {
	im := &indexMap{
		id: make(map[int64]int),
	}
	switch s.(type) {
	case TableWebsiteSlice:
		im.populateWebsite(s.(TableWebsiteSlice))
		break
	case TableGroupSlice:
		im.populateGroup(s.(TableGroupSlice))
		break
	case TableStoreSlice:
		im.populateStore(s.(TableStoreSlice))
		break
	default:
		panic("Unsupported slice: Either TableStoreSlice, TableGroupSlice or TableWebsiteSlice supported")
	}
	return im
}

// populateWebsite fills the map (itself) with the website ids and codes and the index of the slice. Thread safe.
func (im *indexMap) populateWebsite(s TableWebsiteSlice) *indexMap {
	im.Lock()
	defer im.Unlock()
	im.initCode()
	im.clearID()
	for i := 0; i < len(s); i++ {
		im.id[s[i].WebsiteID] = i
		im.code[s[i].Code.String] = i
	}
	return im
}

// populateGroup fills the map (itself) with the group ids and the index of the slice. Thread safe.
func (im *indexMap) populateGroup(s TableGroupSlice) *indexMap {
	im.Lock()
	defer im.Unlock()
	im.clearID()

	for i := 0; i < len(s); i++ {
		im.id[s[i].GroupID] = i
	}
	return im
}

// populateStore fills the map (itself) with the store ids and codes and the index of the slice. Thread safe.
func (im *indexMap) populateStore(s TableStoreSlice) *indexMap {
	im.Lock()
	defer im.Unlock()
	im.initCode()
	im.clearID()
	for i := 0; i < len(s); i++ {
		im.id[s[i].StoreID] = i
		im.code[s[i].Code.String] = i
	}
	return im
}

func (im *indexMap) initCode() {
	if len(im.code) > 0 {
		for k := range im.code {
			delete(im.code, k)
		}
	} else {
		im.code = make(map[string]int)
	}
}

func (im *indexMap) clearID() {
	if len(im.id) > 0 {
		for k := range im.id {
			delete(im.id, k)
		}
	}
}
