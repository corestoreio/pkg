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

/*
package main generates Go structs and slices from SQL tables.

Example

    type (
        // TableStoreSlice contains pointers to TableStore types
        TableStoreSlice []*TableStore
        // TableStore a type for the MySQL table core_store
        TableStore struct {
            StoreID   int64          `db:"store_id"`   // store_id smallint(5) unsigned NOT NULL PRI  auto_increment
            Code      dbr.NullString `db:"code"`       // code varchar(32) NULL UNI
            WebsiteID int64          `db:"website_id"` // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
            GroupID   int64          `db:"group_id"`   // group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
            Name      string         `db:"name"`       // name varchar(255) NOT NULL
            SortOrder int64          `db:"sort_order"` // sort_order smallint(5) unsigned NOT NULL  DEFAULT '0'
            IsActive  bool           `db:"is_active"`  // is_active smallint(5) unsigned NOT NULL MUL DEFAULT '0'
        }
    )

and table structure collections:

	tableMap = csdb.TableStructureSlice{
		TableIndexStore: csdb.NewTableStructure(
			"core_store",
			[]string{
				"store_id",
			},
			[]string{

				"code",
				"website_id",
				"group_id",
				"name",
				"sort_order",
				"is_active",
			},
		),
    ...

*/
package main
