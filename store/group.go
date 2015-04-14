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
	"github.com/juju/errgo"
)

const (
	DefaultGroupId int64 = 0
)

type (
	GroupIndexIDMap map[int64]IDX
	// GroupBucket contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface GroupGetter.
	GroupBucket struct {
		// store collection
		s TableGroupSlice
		// i map by store_id
		i GroupIndexIDMap
		// stores contains a slice to all stores associated to one group.
		// Slice index is the iota value of a group constant.
		stores []TableStoreSlice
		// websites is a slice to TableWebsite
		// Slice index is the iota value of a group constant.
		websites []*TableWebsite
	}
	// GroupGetter methods to retrieve a store pointer
	GroupGetter interface {
		ByID(id int64) (*TableGroup, error)
		Collection() TableGroupSlice
	}
)

var (
	ErrGroupNotFound             = errors.New("Store Group not found")
	ErrGroupStoresNotFound       = errors.New("Store Group stores not found")
	ErrGroupDefaultStoreNotFound = errors.New("Group default store not found")
	ErrGroupWebsiteNotFound      = errors.New("Store Group website not found")
)

var _ GroupGetter = (*GroupBucket)(nil)

// NewGroupBucket returns a new pointer to a GroupBucket.
func NewGroupBucket(s TableGroupSlice, i GroupIndexIDMap) *GroupBucket {
	// @todo idea if i and c is nil generate them from s.
	return &GroupBucket{
		i: i,
		s: s,
	}
}

// ByID uses the database store id to return a TableGroup struct.
func (gb *GroupBucket) ByID(id int64) (*TableGroup, error) {
	if i, ok := gb.i[id]; ok {
		return gb.s[i], nil
	}
	return nil, ErrGroupNotFound
}

// DefaultStore returns the default TableStore for a group id
func (gb *GroupBucket) DefaultStore(id int64) (*TableStore, error) {
	group, err := gb.ByID(id)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	stores, err := gb.Stores(id)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	for _, store := range stores {
		if store.StoreID == group.DefaultStoreID {
			return store, nil
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

// Collection returns the TableGroupSlice
func (gb *GroupBucket) Collection() TableGroupSlice { return gb.s }

func (gb *GroupBucket) Stores(id int64) (TableStoreSlice, error) {
	if i, ok := gb.i[id]; ok {
		return gb.stores[i], nil
	}
	return nil, ErrGroupStoresNotFound
}

func (gb *GroupBucket) Website(id int64) (*TableWebsite, error) {
	if i, ok := gb.i[id]; ok {
		return gb.websites[i], nil
	}
	return nil, ErrGroupWebsiteNotFound
}

// SetStores uses the full store collection to extract the stores which are
// assigned to a group.
func (gb *GroupBucket) SetStores(sg StoreGetter) *GroupBucket {
	gb.stores = make([]TableStoreSlice, len(gb.s), len(gb.s))
	for i, group := range gb.s {
		if group == nil {
			continue
		}
		gb.stores[i] = sg.Collection().FilterByGroupID(group.GroupID)
	}
	return gb
}

// SetWebSite assigns a website to a group
func (gb *GroupBucket) SetWebSite(wb WebsiteGetter) *GroupBucket {
	gb.websites = make([]*TableWebsite, len(gb.s), len(gb.s))
	for i, group := range gb.s {
		if group == nil {
			continue
		}
		var err error
		gb.websites[i], err = wb.ByID(group.WebsiteID)
		if err != nil {
			panic(errgo.Mask(err)) // @todo not nice this one. so fix it
		}
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
