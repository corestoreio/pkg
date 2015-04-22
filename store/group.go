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
	"github.com/corestoreio/csfw/utils"
)

const (
	DefaultGroupId int64 = 0
)

type (

	// GroupBucket contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface GroupGetter.
	GroupBucket struct {
		// g group data
		g *TableGroup
		// stores contains a slice to all stores associated to this group.
		// This slice can be nil
		stores []*StoreBucket
		// website which belongs to this group
		Website *WebsiteBucket
	}
	// GroupGetter methods to retrieve a store pointer
//	GroupGetter interface {
//		Get(int64) (*TableGroup, error)
//	}
)

var (
	ErrGroupNotFound             = errors.New("Store Group not found")
	ErrGroupStoresNotAvailable   = errors.New("Store Group stores not available")
	ErrGroupDefaultStoreNotFound = errors.New("Group default store not found")

//	ErrGroupWebsiteNotFound      = errors.New("Store Group website not found")
)

//var _ GroupGetter = (*GroupBucket)(nil)

// NewGroupBucket returns a new pointer to a GroupBucket. Second argument can be nil.
func NewGroupBucket(g *TableGroup) *GroupBucket {
	if g == nil {
		panic("First argument TableGroup cannot be nil")
	}

	gb := &GroupBucket{
		g: g,
	}
	return gb
}

// Data returns the data from the database
func (gb *GroupBucket) Data() *TableGroup {
	return gb.g
}

// DefaultStore returns the default StoreBucket or an error
func (gb *GroupBucket) DefaultStore(id int64) (*StoreBucket, error) {
	for _, sb := range gb.stores {
		if sb.Data().StoreID == gb.g.DefaultStoreID {
			return sb, nil
		}
	}
	return nil, ErrGroupDefaultStoreNotFound
}

// DefaultStoreByLocale returns the default store using a group ip and a locale
// @todo magento2/app/code/Magento/Store/Model/Group.php::getDefaultStoreByLocale()
// Based on some config values
func (gb *GroupBucket) DefaultStoreByLocale(id int64, locale string) (*TableStore, error) {
	return nil, ErrGroupDefaultStoreNotFound
}

func (gb *GroupBucket) Stores() ([]*StoreBucket, error) {
	if len(gb.stores) > 0 {
		return gb.stores, nil
	}
	return ErrGroupStoresNotAvailable, nil
}

// SetStores uses the full store collection to extract the stores which are
// assigned to a group.
func (gb *GroupBucket) SetStores(tss TableStoreSlice) *GroupBucket {
	if tss == nil {
		gb.stores = nil
		return gb
	}
	@todo
	gb.stores = make([]TableStoreSlice, len(tss), len(tss))
	for i, store := range tss {
		if store == nil {
			continue
		}
		gb.stores[i] = tss.FilterByGroupID(gb.g.GroupID)
	}
	return gb
}

// Load uses a dbr session to load all data from the core_store_group table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Group/Collection.php::_beforeLoad()
func (s *TableGroupSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return loadSlice(dbrSess, TableIndexGroup, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.name ASC")
	})...)
}

// Len returns the length
func (s TableGroupSlice) Len() int { return len(s) }

// FilterByWebsiteID returns a new slice with all groups belonging to a website id
func (s TableGroupSlice) FilterByWebsiteID(id int64) TableGroupSlice {
	return s.Filter(func(w *TableGroup) bool {
		return w.WebsiteID == id
	})
}

// Filter returns a new slice filtered by predicate f
func (s TableGroupSlice) Filter(f func(*TableGroup) bool) TableGroupSlice {
	var tgs TableGroupSlice
	for _, v := range s {
		if v != nil && f(v) {
			tgs = append(tgs, v)
		}
	}
	return tgs
}

// IDs returns an Int64Slice with all group ids
func (s TableGroupSlice) IDs() utils.Int64Slice {
	id := make(utils.Int64Slice, len(s))
	for i, store := range s {
		id[i] = store.GroupID
	}
	return id
}

/*
	@todo implement Magento\Store\Model\Group
*/
