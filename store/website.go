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

// Package store implements the handling of websites, groups and stores
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

// SetGroups @todo
func (wb *WebsiteBucket) SetGroups(gg GroupGetter) *WebsiteBucket {
	return wb
}

// SetStores @todo
func (wb *WebsiteBucket) SetStores(gg StoreGetter) *WebsiteBucket {
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

// @todo review Magento code because of column is_default
//func (s Website) IsDefault() bool {
//	return s.WebsiteID == DefaultWebsiteId
//}

/*
	@todo implement Magento\Store\Model\Website
*/
