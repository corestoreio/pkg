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

// package store implements the handling of websites, groups and stores
package store

import (
	"errors"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
)

const (
	DefaultStoreId int64 = 0
)

// StoreIndex used for iota and for not mixing up indexes
type StoreIndex int

var (
	ErrStoreNotFound = errors.New("Store not found")
	storeCollection  StoreSlice
)

// GetStore uses a StoreIndex to return a store or an error.
// One should not modify the store object.
func GetStore(i StoreIndex) (*Store, error) {
	if int(i) < len(storeCollection) {
		return storeCollection[i], nil
	}
	return nil, ErrStoreNotFound
}

// GetStores returns a copy of the main slice of stores.
// One should not modify the slice and its content.
func GetStores() StoreSlice {
	return storeCollection
}

// Load uses a dbr session to load all data from the core_store table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Store/Collection.php::Load() for sort order
func (s *StoreSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return loadSlice(dbrSess, TableStore, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		sb.OrderBy("CASE WHEN main_table.store_id = 0 THEN 0 ELSE 1 END ASC")
		sb.OrderBy("main_table.sort_order ASC")
		return sb.OrderBy("main_table.name ASC")
	})...)
}

func (s Store) IsDefault() bool {
	return s.StoreID == DefaultStoreId
}

/*
	@todo implement Magento\Store\Model\Store
*/
