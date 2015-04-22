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
		w *TableWebsite

		// groups contains a slice to all groups associated to one website. This slice can be nil.
		groups []*GroupBucket
		// stores contains a slice to all stores associated to one website. This slice can be nil.
		stores []*StoreBucket
	}
	// WebsiteGetter methods to retrieve a store pointer
//	WebsiteGetter interface {
//		// Get first arg website id or 2nd arg website code. Multiple 2nd args will be ignored.
//		Get(int64, ...string) (*TableWebsite, error)
//		Collection() TableWebsiteSlice
//	}
)

var (
	ErrWebsiteNotFound             = errors.New("Website not found")
	ErrWebsiteDefaultGroupNotFound = errors.New("Website Default Group not found")
	ErrWebsiteGroupsNotAvailable   = errors.New("Website Groups not available")
	ErrWebsiteStoresNotAvailable   = errors.New("Website Stores not available")
)

//var _ WebsiteGetter = (*WebsiteBucket)(nil)

// NewWebsiteBucket returns a new pointer to a WebsiteBucket.
func NewWebsiteBucket(w *TableWebsite) *WebsiteBucket {
	return &WebsiteBucket{
		w: w,
	}
}

// Data returns the data from the database
func (wb *WebsiteBucket) Data() *TableWebsite { return wb.w }

// DefaultGroup returns the default GroupBucket or an error if not found
func (wb *WebsiteBucket) DefaultGroup() (*GroupBucket, error) {
	for _, g := range wb.groups {
		if wb.w.DefaultGroupID == g.Data().GroupID {
			return g, nil
		}
	}
	return nil, ErrWebsiteDefaultGroupNotFound
}

// Stores
func (wb *WebsiteBucket) Stores() ([]*StoreBucket, error) {
	if len(wb.stores) > 0 {
		return wb.stores, nil
	}
	return nil, ErrWebsiteStoresNotAvailable
}

// Groups
func (wb *WebsiteBucket) Groups() ([]*GroupBucket, error) {
	if len(wb.groups) > 0 {
		return wb.groups, nil
	}
	return nil, ErrWebsiteGroupsNotAvailable
}

func (wb *WebsiteBucket) SetGroupsStores(tgs TableGroupSlice, tss TableStoreSlice) *WebsiteBucket {
	wb.groups = make([]TableGroupSlice, len(wb.s), len(wb.s))
	for i, website := range wb.s {
		if website == nil {
			continue
		}
		wb.groups[i] = gg.Collection().FilterByWebsiteID(website.WebsiteID)
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
