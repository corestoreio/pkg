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

// +build !mage1,!mage2

// Only include this file IF no specific build tag for mage has been set

package store

/*
	Initial idea to use https://github.com/go-goast/goast
	- so far so good but some generics are hard to implement like the Extract* structs
	  need to invest more time
	- Comments are not ported to the newly generated file
	- This code here works just run `$ go generate .`
*/

// Auto generated via tableToStruct
//go:generate goast write impl --prefix=zgen_ github.com/corestoreio/csfw/codegen/tableToStruct/_goast_generic

import (
	"sort"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/errors"
)

// TableIndex... is the index to a table. These constants are guaranteed
// to stay the same for all Magento versions. Please access a table via this
// constant instead of the raw table name. TableIndex iotas must start with 0.
const (
	TableIndexStore   csdb.Index = iota // Table: store
	TableIndexGroup                     // Table: store_group
	TableIndexWebsite                   // Table: store_website
	TableIndexZZZ                       // the maximum index, which is not available.
)

func init() {
	TableCollection = csdb.MustNewTableService(
		csdb.WithTable(TableIndexStore, "store"),
		csdb.WithTable(TableIndexGroup, "store_group"),
		csdb.WithTable(TableIndexWebsite, "store_website"),
	)
	// Don't forget to call TableCollection.ReInit(...) in your code to load the column definitions.
}

// TableStoreSlice represents a collection type for DB table store
// Generated via tableToStruct.
type TableStoreSlice []*TableStore

// TableStore represents a type for DB table store
// Generated via tableToStruct.
type TableStore struct {
	StoreID   int64          `db:"store_id" json:",omitempty"`   // store_id smallint(5) unsigned NOT NULL PRI  auto_increment
	Code      dbr.NullString `db:"code" json:",omitempty"`       // code varchar(32) NULL UNI
	WebsiteID int64          `db:"website_id" json:",omitempty"` // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	GroupID   int64          `db:"group_id" json:",omitempty"`   // group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	Name      string         `db:"name" json:",omitempty"`       // name varchar(255) NOT NULL
	SortOrder int64          `db:"sort_order" json:",omitempty"` // sort_order smallint(5) unsigned NOT NULL  DEFAULT '0'
	IsActive  bool           `db:"is_active" json:",omitempty"`  // is_active smallint(5) unsigned NOT NULL MUL DEFAULT '0'
}

// parentSQLSelect fills this slice with data from the database.
// Generated via tableToStruct.
func (s *TableStoreSlice) parentSQLSelect(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) (int, error) {
	return csdb.LoadSlice(dbrSess, TableCollection, TableIndexStore, &(*s), cbs...)
}

// SQLInsert inserts all records into the database @todo.
// Generated via tableToStruct.
func (s *TableStoreSlice) SQLInsert(dbrSess dbr.SessionRunner, cbs ...dbr.InsertCb) (int, error) {
	return 0, nil
}

// SQLUpdate updates all record in the database @todo.
// Generated via tableToStruct.
func (s *TableStoreSlice) SQLUpdate(dbrSess dbr.SessionRunner, cbs ...dbr.UpdateCb) (int, error) {
	return 0, nil
}

// SQLDelete deletes all record from the database @todo.
// Generated via tableToStruct.
func (s *TableStoreSlice) SQLDelete(dbrSess dbr.SessionRunner, cbs ...dbr.DeleteCb) (int, error) {
	return 0, nil
}

// ExtractStore functions for extracting fields from Store
// slice. Generated via tableToStruct.
type ExtractStore struct {
	StoreID   func() []int64
	Code      func() []string
	WebsiteID func() []int64
	GroupID   func() []int64
	Name      func() []string
	SortOrder func() []int64
	IsActive  func() []bool
}

// Extract extracts from a specified field all values into a slice.
// Generated via tableToStruct.
func (s TableStoreSlice) Extract() ExtractStore {
	return ExtractStore{
		StoreID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.StoreID)
			}
			return ext
		},
		Code: func() []string {
			ext := make([]string, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.Code.String)
			}
			return ext
		},
		WebsiteID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.WebsiteID)
			}
			return ext
		},
		GroupID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.GroupID)
			}
			return ext
		},
		Name: func() []string {
			ext := make([]string, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.Name)
			}
			return ext
		},
		SortOrder: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.SortOrder)
			}
			return ext
		},
		IsActive: func() []bool {
			ext := make([]bool, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.IsActive)
			}
			return ext
		},
	}
}

// TableGroupSlice represents a collection type for DB table store_group
// Generated via tableToStruct.
type TableGroupSlice []*TableGroup

// TableGroup represents a type for DB table store_group
// Generated via tableToStruct.
type TableGroup struct {
	GroupID        int64  `db:"group_id" json:",omitempty"`         // group_id smallint(5) unsigned NOT NULL PRI  auto_increment
	WebsiteID      int64  `db:"website_id" json:",omitempty"`       // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	Name           string `db:"name" json:",omitempty"`             // name varchar(255) NOT NULL
	RootCategoryID int64  `db:"root_category_id" json:",omitempty"` // root_category_id int(10) unsigned NOT NULL  DEFAULT '0'
	DefaultStoreID int64  `db:"default_store_id" json:",omitempty"` // default_store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
}

// parentSQLSelect fills this slice with data from the database.
// Generated via tableToStruct.
func (s *TableGroupSlice) parentSQLSelect(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) (int, error) {
	return csdb.LoadSlice(dbrSess, TableCollection, TableIndexGroup, &(*s), cbs...)
}

// SQLInsert inserts all records into the database @todo.
// Generated via tableToStruct.
func (s *TableGroupSlice) SQLInsert(dbrSess dbr.SessionRunner, cbs ...dbr.InsertCb) (int, error) {
	return 0, nil
}

// SQLUpdate updates all record in the database @todo.
// Generated via tableToStruct.
func (s *TableGroupSlice) SQLUpdate(dbrSess dbr.SessionRunner, cbs ...dbr.UpdateCb) (int, error) {
	return 0, nil
}

// SQLDelete deletes all record from the database @todo.
// Generated via tableToStruct.
func (s *TableGroupSlice) SQLDelete(dbrSess dbr.SessionRunner, cbs ...dbr.DeleteCb) (int, error) {
	return 0, nil
}

// ExtractGroup functions for extracting fields from Group
// slice. Generated via tableToStruct.
type ExtractGroup struct {
	GroupID        func() []int64
	WebsiteID      func() []int64
	Name           func() []string
	RootCategoryID func() []int64
	DefaultStoreID func() []int64
}

// Extract extracts from a specified field all values into a slice.
// Generated via tableToStruct.
func (s TableGroupSlice) Extract() ExtractGroup {
	return ExtractGroup{
		GroupID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.GroupID)
			}
			return ext
		},
		WebsiteID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.WebsiteID)
			}
			return ext
		},
		Name: func() []string {
			ext := make([]string, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.Name)
			}
			return ext
		},
		RootCategoryID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.RootCategoryID)
			}
			return ext
		},
		DefaultStoreID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.DefaultStoreID)
			}
			return ext
		},
	}
}

// TableWebsiteSlice represents a collection type for DB table store_website
// Generated via tableToStruct.
type TableWebsiteSlice []*TableWebsite

// TableWebsite represents a type for DB table store_website
// Generated via tableToStruct.
type TableWebsite struct {
	WebsiteID      int64          `db:"website_id" json:",omitempty"`       // website_id smallint(5) unsigned NOT NULL PRI  auto_increment
	Code           dbr.NullString `db:"code" json:",omitempty"`             // code varchar(32) NULL UNI
	Name           dbr.NullString `db:"name" json:",omitempty"`             // name varchar(64) NULL
	SortOrder      int64          `db:"sort_order" json:",omitempty"`       // sort_order smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	DefaultGroupID int64          `db:"default_group_id" json:",omitempty"` // default_group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	IsDefault      dbr.NullBool   `db:"is_default" json:",omitempty"`       // is_default smallint(5) unsigned NULL  DEFAULT '0'
}

// parentSQLSelect fills this slice with data from the database.
// Generated via tableToStruct.
func (s *TableWebsiteSlice) parentSQLSelect(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) (int, error) {
	return csdb.LoadSlice(dbrSess, TableCollection, TableIndexWebsite, &(*s), cbs...)
}

// SQLInsert inserts all records into the database @todo.
// Generated via tableToStruct.
func (s *TableWebsiteSlice) SQLInsert(dbrSess dbr.SessionRunner, cbs ...dbr.InsertCb) (int, error) {
	return 0, nil
}

// SQLUpdate updates all record in the database @todo.
// Generated via tableToStruct.
func (s *TableWebsiteSlice) SQLUpdate(dbrSess dbr.SessionRunner, cbs ...dbr.UpdateCb) (int, error) {
	return 0, nil
}

// SQLDelete deletes all record from the database @todo.
// Generated via tableToStruct.
func (s *TableWebsiteSlice) SQLDelete(dbrSess dbr.SessionRunner, cbs ...dbr.DeleteCb) (int, error) {
	return 0, nil
}

// ExtractWebsite functions for extracting fields from Website
// slice. Generated via tableToStruct.
type ExtractWebsite struct {
	WebsiteID      func() []int64
	Code           func() []string
	Name           func() []string
	SortOrder      func() []int64
	DefaultGroupID func() []int64
	IsDefault      func() []bool
}

// Extract extracts from a specified field all values into a slice.
// Generated via tableToStruct.
func (s TableWebsiteSlice) Extract() ExtractWebsite {
	return ExtractWebsite{
		WebsiteID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.WebsiteID)
			}
			return ext
		},
		Code: func() []string {
			ext := make([]string, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.Code.String)
			}
			return ext
		},
		Name: func() []string {
			ext := make([]string, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.Name.String)
			}
			return ext
		},
		SortOrder: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.SortOrder)
			}
			return ext
		},
		DefaultGroupID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.DefaultGroupID)
			}
			return ext
		},
		IsDefault: func() []bool {
			ext := make([]bool, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.IsDefault.Bool)
			}
			return ext
		},
	}
}
