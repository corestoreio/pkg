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
	// DefaultStoreID is always 0.
	DefaultStoreID int64 = 0
	// HTTPRequestParamStore name of the GET parameter to set a new store in a current website/group context
	HTTPRequestParamStore = `___store`
	// CookieName important when the user selects a different store within the current website/group context.
	// This cookie permanently saves the new selected store code for one year.
	// The cookie must be removed when the default store of the current website if equal to the current store.
	CookieName = `store`
)

type (
	// Store contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface StoreGetter.
	Store struct {
		// Contains the current website for this store. No integrity checks
		w *Website
		g *Group
		// underlaying raw data
		s *TableStore
	}
	// StoreSlice a collection of pointers to the Store structs. StoreSlice has some nifty method receviers.
	StoreSlice []*Store
)

var (
	ErrStoreNotFound         = errors.New("Store not found")
	ErrStoreNewArgNil        = errors.New("An argument cannot be nil")
	ErrStoreIncorrectGroup   = errors.New("Incorrect group")
	ErrStoreIncorrectWebsite = errors.New("Incorrect website")
)

// NewStore returns a new pointer to a Store. Panics if one of the arguments is nil.
// The integrity checks are done by the database.
func NewStore(w *TableWebsite, g *TableGroup, s *TableStore) *Store {
	if w == nil || g == nil || s == nil {
		panic(ErrStoreNewArgNil)
	}

	if s.GroupID != g.GroupID {
		panic(ErrStoreIncorrectGroup)
	}

	if s.WebsiteID != w.WebsiteID {
		panic(ErrStoreIncorrectWebsite)
	}

	return &Store{
		w: NewWebsite(w),
		g: NewGroup(g, nil),
		s: s,
	}
}

/*
	@todo implement Magento\Store\Model\Store
*/

// Website returns the website associated to this store
func (s *Store) Website() *Website {
	return s.w
}

// Group returns the group associated to this store
func (s *Store) Group() *Group {
	return s.g
}

// Data returns the real store data from the database
func (s *Store) Data() *TableStore {
	return s.s
}

/*
	StoreSlice method receivers
*/

// Len returns the length
func (s StoreSlice) Len() int { return len(s) }

// Filter returns a new slice filtered by predicate f
func (s StoreSlice) Filter(f func(*Store) bool) StoreSlice {
	var stores StoreSlice
	for _, v := range s {
		if v != nil && f(v) {
			stores = append(stores, v)
		}
	}
	return stores
}

// Codes returns a StringSlice with all store codes
func (s StoreSlice) Codes() utils.StringSlice {
	if len(s) == 0 {
		return nil
	}
	var c utils.StringSlice
	for _, st := range s {
		if st != nil {
			c.Append(st.Data().Code.String)
		}
	}
	return c
}

// IDs returns an Int64Slice with all store ids
func (s StoreSlice) IDs() utils.Int64Slice {
	if len(s) == 0 {
		return nil
	}
	var ids utils.Int64Slice
	for _, st := range s {
		if st != nil {
			ids.Append(st.Data().StoreID)
		}
	}
	return ids
}

/*
	TableStore and TableStoreSlice method receivers
*/

// IsDefault returns true if the current store is the default store.
func (s TableStore) IsDefault() bool {
	return s.StoreID == DefaultStoreID
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
	for _, st := range s {
		if st != nil && st.StoreID == id {
			return st, nil
		}
	}
	return nil, ErrStoreNotFound
}

// FilterByGroupID returns a new slice with all TableStores belonging to a group id
func (s TableStoreSlice) FilterByGroupID(id int64) TableStoreSlice {
	return s.Filter(func(ts *TableStore) bool {
		return ts.GroupID == id
	})
}

// FilterByWebsiteID returns a new slice with all TableStores belonging to a website id
func (s TableStoreSlice) FilterByWebsiteID(id int64) TableStoreSlice {
	return s.Filter(func(ts *TableStore) bool {
		return ts.WebsiteID == id
	})
}

// Filter returns a new slice containing TableStores filtered by predicate f
func (s TableStoreSlice) Filter(f func(*TableStore) bool) TableStoreSlice {
	if len(s) == 0 {
		return nil
	}
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
	if len(s) == 0 {
		return nil
	}
	var c utils.StringSlice
	for _, store := range s {
		if store != nil {
			c.Append(store.Code.String)
		}
	}
	return c
}

// IDs returns an Int64Slice with all store ids
func (s TableStoreSlice) IDs() utils.Int64Slice {
	if len(s) == 0 {
		return nil
	}
	var ids utils.Int64Slice
	for _, store := range s {
		if store != nil {
			ids.Append(store.StoreID)
		}
	}
	return ids
}
