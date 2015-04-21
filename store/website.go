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
	// WebsiteBucket contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface WebsiteGetter.
	WebsiteBucket struct {
		// store collection
		//		s TableWebsiteSlice
		//		// im index map by website id and code
		//		im *indexMap
		//
		//		// groups contains a slice to all groups associated to one website.
		//		// Slice index is the iota value of a website constant.
		//		groups []TableGroupSlice
		//		// stores contains a slice to all stores associated to one website.
		//		// Slice index is the iota value of a website constant.
		//		stores []TableStoreSlice
	}
	// WebsiteGetter methods to retrieve a store pointer
	WebsiteGetter interface {
		// Get first arg website id or 2nd arg website code. Multiple 2nd args will be ignored.
		Get(int64, ...string) (*TableWebsite, error)
		Collection() TableWebsiteSlice
	}
)

var (
	ErrWebsiteNotFound       = errors.New("Website not found")
	ErrWebsiteGroupNotFound  = errors.New("Website Group not found")
	ErrWebsiteStoresNotFound = errors.New("Website Stores not found")
)
var _ WebsiteGetter = (*WebsiteBucket)(nil)

// NewWebsiteBucket returns a new pointer to a WebsiteBucket.
func NewWebsiteBucket(s TableWebsiteSlice) *WebsiteBucket {
	return &WebsiteBucket{
		im: (&indexMap{}).populateWebsite(s),
		s:  s,
	}
}

// Get accepts one or two arguments to return a TableWebsite struct. The 2nd argument
// can be the website code. If the 2nd arguments is present then website id as 1st argument
// will be ignored.
func (s *WebsiteBucket) Get(wID int64, wc ...string) (*TableWebsite, error) {

}

// Collection returns the TableWebsiteSlice
func (s *WebsiteBucket) Collection() TableWebsiteSlice { return s.s }

// Group accepts one or two arguments to return a GroupBucket. The 2nd argument
// can be the website code. If the 2nd arguments is present then website id as 1st argument
// will be ignored.
func (s *WebsiteBucket) Group(wID int64, wc ...string) (*GroupBucket, error) {
	i, oki := s.im.id[wID]
	if !oki || s.s[i] == nil {
		return nil, ErrWebsiteNotFound
	}
	if len(s.groups) < int(i) {
		return nil, ErrWebsiteGroupNotFound
	}
	g := s.groups[i]
	if g == nil {
		return nil, ErrWebsiteGroupNotFound
	}
	st := s.stores[i]
	if st == nil {
		return nil, ErrWebsiteStoresNotFound
	}

	gb := NewGroupBucket(g)
	sb := NewStoreBucket(st)
	gb.SetStores(sb).SetWebSite(s)
	return gb, nil
}

// Stores @todo
func (s *WebsiteBucket) Stores(wID int64, wc ...string) (*StoreBucket, error) { return nil, nil }

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
