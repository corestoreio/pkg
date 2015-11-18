// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

import "github.com/corestoreio/csfw/storage/dbr"

/*
	TableStore and TableStoreSlice method receivers
*/

// IsDefault returns true if the current store is the default store.
func (s TableStore) IsDefault() bool {
	return s.StoreID == DefaultStoreID
}

// SQLSelect uses a dbr session to load all data from the core_store table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see https://github.com/magento/magento2/blob/0.74.0-beta7/app%2Fcode%2FMagento%2FStore%2FModel%2FResource%2FStore%2FCollection.php#L147
// regarding the sort order.
func (s *TableStoreSlice) SQLSelect(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) (int, error) {
	return s.parentSQLSelect(dbrSess, append(append([]dbr.SelectCb{nil}, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		sb.OrderBy("CASE WHEN main_table.store_id = 0 THEN 0 ELSE 1 END ASC")
		sb.OrderBy("main_table.sort_order ASC")
		return sb.OrderBy("main_table.name ASC")
	}), cbs...)...)
}

// FilterByGroupID returns a new slice with all TableStores belonging to a group id
func (s TableStoreSlice) FilterByGroupID(id int64) TableStoreSlice {
	return s.Filter(func(ts *TableStore) bool {
		return ts != nil && ts.GroupID == id
	})
}

// FilterByWebsiteID returns a new slice with all TableStores belonging to a website id
func (s TableStoreSlice) FilterByWebsiteID(id int64) TableStoreSlice {
	return s.Filter(func(ts *TableStore) bool {
		return ts != nil && ts.WebsiteID == id
	})
}
