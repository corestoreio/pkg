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
	DefaultStoreId int64 = 0
)

type (
	// Store contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface StoreGetter.
	Store struct {
		Website *Website
		Group   *Group
		s       *TableStore
	}
)

var (
	ErrStoreNotFound = errors.New("Store not found")
)

// NewStore returns a new pointer to a Store. Panics if one of the arguments is nil.
func NewStore(w *TableWebsite, g *TableGroup, s *TableStore) *Store {
	if w == nil || g == nil || s == nil {
		panic("An argument cannot be nil")
	}
	wb := NewWebsite(w)
	return &Store{
		Website: wb,
		Group:   NewGroup(g),
		s:       s,
	}
}

// Data returns the real store data from the database
func (s *Store) Data() *TableStore {
	return s.s
}

// Load uses a dbr session to load all data from the core_store table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Store/Collection.php::Load() for sort order
func (s *TableStoreSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return loadSlice(dbrSess, TableIndexStore, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		sb.OrderBy("CASE WHEN main_table.store_id = 0 THEN 0 ELSE 1 END ASC")
		sb.OrderBy("main_table.sort_order ASC")
		return sb.OrderBy("main_table.name ASC")
	})...)
}

// Len returns the length
func (s TableStoreSlice) Len() int { return len(s) }

// FindByID returns a TableStore if found by id or an error
func (s TableStoreSlice) FindByID(id int64) (*TableStore, error) {
	for _, store := range s {
		if store.StoreID == id {
			return store, nil
		}
	}
	return nil, ErrStoreNotFound
}

// ByGroupID returns a new slice with all stores belonging to a group id
func (s TableStoreSlice) FilterByGroupID(id int64) TableStoreSlice {
	return s.Filter(func(store *TableStore) bool {
		return store.GroupID == id
	})
}

// FilterByWebsiteID returns a new slice with all stores belonging to a website id
func (s TableStoreSlice) FilterByWebsiteID(id int64) TableStoreSlice {
	return s.Filter(func(store *TableStore) bool {
		return store.WebsiteID == id
	})
}

// Filter returns a new slice filtered by predicate f
func (s TableStoreSlice) Filter(f func(*TableStore) bool) TableStoreSlice {
	var tss TableStoreSlice
	for _, v := range s {
		if v != nil && f(v) {
			tss = append(tss, v)
		}
	}
	return tss
}

// Codes returns a StringSlice with all store codes
func (s TableStoreSlice) Codes() utils.StringSlice {
	c := make(utils.StringSlice, len(s))
	for i, store := range s {
		c[i] = store.Code.String
	}
	return c
}

// IDs returns an Int64Slice with all store ids
func (s TableStoreSlice) IDs() utils.Int64Slice {
	id := make(utils.Int64Slice, len(s))
	for i, store := range s {
		id[i] = store.StoreID
	}
	return id
}

func (s TableStore) IsDefault() bool {
	return s.StoreID == DefaultStoreId
}

/*
	@todo implement Magento\Store\Model\Store
*/
