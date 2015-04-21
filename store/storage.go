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
		websiteIM: (&indexMap{}).populateWebsite(tws),
		Groups:    tgs,
		groupIM:   (&indexMap{}).populateGroup(tgs),
		Stores:    tss,
		storeIM:   (&indexMap{}).populateStore(tss),
	}
}

func (s *Storage) New(scopeCode string, scopeType config.ScopeID) (s *StoreBucket, g *GroupBucket, w *WebsiteBucket) {

}

// Website returns a table website either by id or code, one of them can be nil but not both.
func (s *Storage) Website(id IDRetriever, c CodeRetriever) (*WebsiteBucket, error) {
	var idx int = -1
	switch {
	case id != nil:
		if i, ok := s.websiteIM.id[id.ID()]; ok {
			idx = s.Websites[i]
		}
		break
	case c != nil:
		if i, ok := s.websiteIM.code[c.Code()]; ok {
			idx = s.Websites[i]
		}
		break
	default:
		return nil, ErrWebsiteNotFound
	}
	if idx < 0 {
		return nil, ErrWebsiteNotFound
	}

	return nil, nil
}

// Group returns a table group either by id or code, one of them can be nil but not both.
func (s *Storage) Group(id IDRetriever) (*GroupBucket, error) {
	var idx int = -1
	switch {
	case id != nil:
		if i, ok := s.groupIM.id[id.ID()]; ok {
			return s.Groups[i], nil
		}
		break
	}
	return nil, ErrGroupNotFound
}

func (s *Storage) Store(id IDRetriever, c CodeRetriever) (*StoreBucket, error) {
	switch {
	case id != nil:
		if i, ok := s.storeIM.id[id.ID()]; ok {
			return s.Stores[i], nil
		}
		break
	case c != nil:
		if i, ok := s.storeIM.code[c.Code()]; ok {
			return s.Stores[i], nil
		}
		break
	}
	return nil, ErrStoreNotFound
}
