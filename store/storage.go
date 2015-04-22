package store

import (
	"sync"

	"github.com/corestoreio/csfw/config"
)

type (
	Storage struct {
		Websites  TableWebsiteSlice
		websiteIM *indexMap
		Groups    TableGroupSlice
		groupIM   *indexMap
		Stores    TableStoreSlice
		storeIM   *indexMap
	}
	IDRetriever interface {
		ID() int64
	}
	CodeRetriever interface {
		Code() string
	}
	// ID is convenience helper to satisfy the interface IDRetriever
	ID int64
	// Code is convenience helper to satisfy the interface CodeRetriever
	Code string
)

// ID is convenience helper to satisfy the interface IDRetriever
func (i ID) ID() int64 { return int64(i) }

// Code is convenience helper to satisfy the interface CodeRetriever
func (c Code) Code() string { return string(c) }

func NewStorage(tws TableWebsiteSlice, tgs TableGroupSlice, tss TableStoreSlice) *Storage {
	// maybe we can totally remove the maps and just rely on the for loop to find an entity.
	return &Storage{
		Websites:  tws,
		websiteIM: newIndexMap(tws),
		Groups:    tgs,
		groupIM:   newIndexMap(tgs),
		Stores:    tss,
		storeIM:   newIndexMap(tss),
	}
}

func (s *Storage) News(scopeCode string, scopeType config.ScopeID) (sb *Store, gb *Group, wb *Website) {
	return nil, nil, nil
}

// website returns a TableWebsite
func (s *Storage) website(id IDRetriever, c CodeRetriever) (*TableWebsite, error) {
	var idx int = -1
	switch {
	case id != nil:
		if i, ok := s.websiteIM.id[id.ID()]; ok {
			idx = i
		}
		break
	case c != nil:
		if i, ok := s.websiteIM.code[c.Code()]; ok {
			idx = i
		}
		break
	default:
		return nil, ErrWebsiteNotFound
	}
	if idx < 0 {
		return nil, ErrWebsiteNotFound
	}
	return s.Websites[idx], nil
}

// Website creates a new Website which contains the current website, all its groups and
// all its related stores. Groups and stores can be nil. It panics when the integrity is incorrect.
func (s *Storage) Website(id IDRetriever, c CodeRetriever) (*Website, error) {
	w, err := s.website(id, c)
	if err != nil {
		return nil, err
	}
	return NewWebsite(w).SetGroupsStores(s.Groups, s.Stores), nil
}

// group returns a TableGroup
func (s *Storage) group(id IDRetriever) (*TableGroup, error) {
	var idx int = -1
	switch {
	case id != nil:
		if i, ok := s.groupIM.id[id.ID()]; ok {
			idx = i
		}
		break
	default:
		return nil, ErrGroupNotFound
	}
	if idx < 0 {
		return nil, ErrGroupNotFound
	}
	return s.Groups[idx], nil
}

// Group creates a new Group which contains all related stores and its website
func (s *Storage) Group(id IDRetriever) (*Group, error) {
	g, err := s.group(id)
	if err != nil {
		return nil, err
	}

	w, err := s.website(ID(g.WebsiteID), nil)
	if err != nil {
		return nil, err
	}
	return NewGroup(g).SetStores(s.Stores, w), nil
}

// store returns a TableStore
func (s *Storage) store(id IDRetriever, c CodeRetriever) (*TableStore, error) {
	var idx int = -1
	switch {
	case id != nil:
		if i, ok := s.storeIM.id[id.ID()]; ok {
			idx = i
		}
		break
	case c != nil:
		if i, ok := s.storeIM.code[c.Code()]; ok {
			idx = i
		}
		break
	default:
		return nil, ErrStoreNotFound
	}
	if idx < 0 {
		return nil, ErrStoreNotFound
	}
	return s.Stores[idx], nil
}

// Store creates a new Store which contains the current store, its Group and Website
func (s *Storage) Store(id IDRetriever, c CodeRetriever) (*Store, error) {
	store, err := s.store(id, c)
	if err != nil {
		return nil, err
	}
	website, err := s.website(ID(store.WebsiteID), nil)
	if err != nil {
		return nil, err
	}
	group, err := s.group(ID(store.GroupID))
	if err != nil {
		return nil, err
	}
	return NewStore(website, group, store), nil
}

// DefaultStoreView traverses through the websites to find the default website and gets
// the group which has the default store id assigned to. Only one website can be the default one.
func (s *Storage) DefaultStoreView() (*Store, error) {
	for _, website := range s.Websites {
		if website.IsDefault.Bool && website.IsDefault.Valid {
			group, err := s.group(ID(website.DefaultGroupID))
			if err != nil {
				return nil, err
			}
			return s.Store(ID(group.DefaultStoreID), nil)
		}
	}
	return nil, ErrStoreNotFound
}

/*
	INTERNAL
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
	im.code = make(map[string]int)
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
	for i := 0; i < len(s); i++ {
		im.id[s[i].GroupID] = i
	}
	return im
}

// populateStore fills the map (itself) with the store ids and codes and the index of the slice. Thread safe.
func (im *indexMap) populateStore(s TableStoreSlice) *indexMap {
	im.Lock()
	defer im.Unlock()
	im.code = make(map[string]int)
	for i := 0; i < len(s); i++ {
		im.id[s[i].StoreID] = i
		im.code[s[i].Code.String] = i
	}
	return im
}
