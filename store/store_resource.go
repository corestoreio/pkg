// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"github.com/corestoreio/pkg/storage/csdb"
	"github.com/corestoreio/pkg/storage/dbr"
)

func init() {
	TableCollection = csdb.MustInitTables(TableCollection,
		csdb.WithTableDMLListeners(TableIndexStore,
			dbr.MustNewListenerBucket(
				dbr.Listen{
					Name:      "admin store on top",
					EventType: dbr.OnBeforeToSQL,
					SelectFunc: func(sb *dbr.Select) {
						sb.OrderBy("CASE WHEN main_table.store_id = 0 THEN 0 ELSE 1 END ASC")
						sb.OrderBy("main_table.sort_order ASC")
						sb.OrderBy("main_table.name ASC")
					},
				},
				dbr.Listen{
					EventType: dbr.OnBeforeToSQL,
					InsertFunc: func(ib *dbr.Insert) {
						// todo ... ?
					},
				},
			),
		),
	)
}

// IsDefault returns true if the current store is the default store.
func (s TableStore) IsDefault() bool {
	return s.StoreID == DefaultStoreID
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
