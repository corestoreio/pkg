// Copyright 2015 CoreStore Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"errors"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
)

const (
	DefaultWebsiteId int64 = 0
)

type (
	WebsiteIndexCodeMap map[string]IDX
	WebsiteIndexIDMap   map[int64]IDX
	// WebsiteBucket contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface WebsiteGetter.
	WebsiteBucket struct {
		// store collection
		s TableWebsiteSlice
		// c map bei code
		c WebsiteIndexCodeMap
		// i map by store_id
		i WebsiteIndexIDMap

		// groups contains a slice to all groups associated to one website.
		// Slice index is the iota value of a website constant.
		groups []TableGroupSlice
		// stores contains a slice to all stores associated to one website.
		// Slice index is the iota value of a website constant.
		stores []TableStoreSlice
	}
	// WebsiteGetter methods to retrieve a store pointer
	WebsiteGetter interface {
		ByID(id int64) (*TableWebsite, error)
		ByCode(code string) (*TableWebsite, error)
		Collection() TableWebsiteSlice
	}
)

var (
	ErrWebsiteNotFound = errors.New("Website not found")
)
var _ WebsiteGetter = (*WebsiteBucket)(nil)

// NewWebsiteBucket returns a new pointer to a WebsiteBucket.
func NewWebsiteBucket(s TableWebsiteSlice, i WebsiteIndexIDMap, c WebsiteIndexCodeMap) *WebsiteBucket {
	// @todo idea if i and c is nil generate them from s.
	return &WebsiteBucket{
		i: i,
		c: c,
		s: s,
	}
}

// ByID uses the database store id to return a TableWebsite struct.
func (s *WebsiteBucket) ByID(id int64) (*TableWebsite, error) {
	if i, ok := s.i[id]; ok {
		return s.s[i], nil
	}
	return nil, ErrWebsiteNotFound
}

// ByCode uses the database store code to return a TableWebsite struct.
func (s *WebsiteBucket) ByCode(code string) (*TableWebsite, error) {
	if i, ok := s.c[code]; ok {
		return s.s[i], nil
	}
	return nil, ErrWebsiteNotFound
}

// Collection returns the TableWebsiteSlice
func (s *WebsiteBucket) Collection() TableWebsiteSlice { return s.s }

// GroupByID @todo
func (s *WebsiteBucket) GroupByID(id int64) *GroupBucket { return nil }

// GroupByCode @todo
func (s *WebsiteBucket) GroupByCode(code string) *GroupBucket { return nil }

// SetGroups uses the full group collection to extract the groups which are
// assigned to a website.
func (wb *WebsiteBucket) SetGroups(gg GroupGetter) *WebsiteBucket {
	wb.groups = make([]TableGroupSlice, len(wb.s), len(wb.s))
	for i, website := range wb.s {
		if website == nil {
			continue
		}
		wb.groups[i] = gg.Collection().FilterByWebsiteID(website.WebsiteID)
	}
	return wb
}

// SetStores uses the full store collection to extract the stores which are
// assigned to a website.
func (wb *WebsiteBucket) SetStores(sg StoreGetter) *WebsiteBucket {
	wb.stores = make([]TableStoreSlice, len(wb.s), len(wb.s))
	for i, website := range wb.s {
		if website == nil {
			continue
		}
		wb.stores[i] = sg.Collection().FilterByWebsiteID(website.WebsiteID)
	}
	return wb
}

// Load uses a dbr session to load all data from the core_website table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Website/Collection.php::Load()
func (s *TableWebsiteSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return loadSlice(dbrSess, TableIndexWebsite, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.sort_order ASC").OrderBy("main_table.name ASC")
	})...)
}

// Len returns the length
func (s TableWebsiteSlice) Len() int { return len(s) }

// Filter returns a new slice filtered by predicate f
func (s TableWebsiteSlice) Filter(f func(*TableWebsite) bool) TableWebsiteSlice {
	var tws TableWebsiteSlice
	for _, v := range s {
		if v != nil && f(v) {
			tws = append(tws, v)
		}
	}
	return tws
}

// @todo review Magento code because of column is_default
//func (s Website) IsDefault() bool {
//	return s.WebsiteID == DefaultWebsiteId
//}

/*
	@todo implement Magento\Store\Model\Website
*/
