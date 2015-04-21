package store

import "github.com/corestoreio/csfw/config"

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

// Website returns a table website either by id or code, one of them can be nil but not both.
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
	website := s.Websites[idx]

	return nil, nil
}

// Group returns a table group either by id or code, one of them can be nil but not both.
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
	group := s.Groups[idx]

	return nil, nil
}

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
	store := s.Stores[idx]

	return nil, nil
}
