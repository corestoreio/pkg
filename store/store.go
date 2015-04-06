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
	"github.com/juju/errgo"
)

const (
	DefaultStoreId int64 = 0
)

type (
	// StoreIndex used for iota and for not mixing up indexes
	StoreIndex int
	// StoreGetter contains generated code from the database to provide easy and fast methods to
	// retrieve the stores
	StoreGetter interface {
		// ByID returns a StoreIndex using the StoreID.  This StoreIndex identifies a store within a StoreSlice.
		ByID(id int64) (StoreIndex, error)
		// ByCode returns a StoreIndex using the code.  This StoreIndex identifies a store within a StoreSlice.
		ByCode(code string) (StoreIndex, error)
	}
)

var (
	ErrStoreNotFound     = errors.New("Store not found")
	ErrStoreGetterNotSet = errors.New("StoreGetter not set")
	storeCollection      TableStoreSlice
	storeGetter          StoreGetter
)

func SetStoreCollection(sc TableStoreSlice) {
	if len(sc) == 0 {
		panic("StoreSlice is empty")
	}
	storeCollection = sc
}

func SetStoreGetter(g StoreGetter) {
	if g == nil {
		panic("StoreGetter cannot be nil")
	}
	storeGetter = g
}

// GetStore uses a StoreIndex to return a store or an error.
// One should not modify the store object.
func GetStore(i StoreIndex) (*TableStore, error) {
	if int(i) < len(storeCollection) {
		return storeCollection[i], nil
	}
	return nil, ErrStoreNotFound
}

func GetStoreByID(id int64) (*TableStore, error) {
	return storeCollection.ByID(id)
}

func GetStoreByCode(code string) (*TableStore, error) {
	return storeCollection.ByCode(code)
}

// GetStores returns a copy of the main slice of stores. There can be nils within the slice.
// One should not modify the slice and its content.
func GetStores() TableStoreSlice {
	return storeCollection
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

func (s TableStoreSlice) ByID(id int64) (*TableStore, error) {
	if storeGetter == nil {
		return nil, ErrStoreGetterNotSet
	}
	i, err := storeGetter.ByID(id)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

func (s TableStoreSlice) ByCode(code string) (*TableStore, error) {
	if storeGetter == nil {
		return nil, ErrStoreGetterNotSet
	}
	i, err := storeGetter.ByCode(code)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

func (s TableStore) IsDefault() bool {
	return s.StoreID == DefaultStoreId
}

/*
	@todo implement Magento\Store\Model\Store
*/
