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
	DefaultStoreId int64 = 0
)

type (
	StoreIndexCodeMap map[string]IDX
	StoreIndexIDMap   map[int64]IDX
	// StoreBucket contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface StoreGetter.
	StoreBucket struct {
		// store collection
		s TableStoreSlice
		// c map bei code
		c StoreIndexCodeMap
		// i map by store_id
		i StoreIndexIDMap
	}
	// StoreGetter methods to retrieve a store pointer
	StoreGetter interface {
		ByID(id int64) (*TableStore, error)
		ByCode(code string) (*TableStore, error)
		Collection() TableStoreSlice
	}
)

var (
	ErrStoreNotFound             = errors.New("Store not found")
	_                StoreGetter = (*StoreBucket)(nil)
)

// NewStoreBucket returns a new pointer to a StoreBucket.
func NewStoreBucket(s TableStoreSlice, i StoreIndexIDMap, c StoreIndexCodeMap) *StoreBucket {
	// @todo idea if i and c is nil generate them from s.
	return &StoreBucket{
		i: i,
		c: c,
		s: s,
	}
}

// ByID uses the database store id to return a TableStore struct.
func (s *StoreBucket) ByID(id int64) (*TableStore, error) {
	if i, ok := s.i[id]; ok {
		return s.s[i], nil
	}
	return nil, ErrStoreNotFound
}

// ByCode uses the database store code to return a TableStore struct.
func (s *StoreBucket) ByCode(code string) (*TableStore, error) {
	if i, ok := s.c[code]; ok {
		return s.s[i], nil
	}
	return nil, ErrStoreNotFound
}

// Collection returns the TableStoreSlice
func (s *StoreBucket) Collection() TableStoreSlice { return s.s }

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

// ByGroupID returns a new slice with all stores belonging to a group id
func (s TableStoreSlice) FilterByGroupID(id int64) TableStoreSlice {
	return s.Filter(func(store *TableStore) bool {
		return store.GroupID == id
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

func (s TableStore) IsDefault() bool {
	return s.StoreID == DefaultStoreId
}

/*
	@todo implement Magento\Store\Model\Store
*/
