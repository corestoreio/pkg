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
)

func NewStorage(tws TableWebsiteSlice, tgs TableGroupSlice, tss TableStoreSlice) *Storage {
	return &Storage{
		Websites:  tws,
		websiteIM: newIndexMap(tws),
		Groups:    tgs,
		groupIM:   newIndexMap(tgs),
		Stores:    tss,
		storeIM:   newIndexMap(tss),
	}
}

func (s *Storage) NewBuckets(scopeCode string, scopeType config.ScopeID) (sb *StoreBucket, gb *GroupBucket, wb *WebsiteBucket) {
	return nil, nil, nil
}

// Website creates a new WebsiteBucket which contains the current website, all its groups and
// all its related stores. Groups and stores can be nil.
func (s *Storage) Website(id IDRetriever, c CodeRetriever) (*WebsiteBucket, error) {
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
	//	website := s.Websites[idx]

	return nil, nil
}

// Group creates a new GroupBucket which contains all related stores and its website
func (s *Storage) Group(id IDRetriever) (*GroupBucket, error) {
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

	return NewGroupBucket(s.Groups[idx]).SetStores(s.Stores), nil
}

// Store creates a new StoreBucket which contains the current store, its GroupBucket and WebsiteBucket
func (s *Storage) Store(id IDRetriever, c CodeRetriever) (*StoreBucket, error) {
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
	//	store := s.Stores[idx]

	return nil, nil
}

// DefaultStoreView traverses through the websites to find the default website and gets
// the group which has the default store id assigned to. Only one website can be the default one.
func (s *Storage) DefaultStoreView() (*StoreBucket, error) {
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
		im.populateWebsite(s)
		break
	case TableGroupSlice:
		im.populateGroup(s)
		break
	case TableStoreSlice:
		im.populateStore(s)
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
