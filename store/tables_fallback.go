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

// Auto generated via tableToStruct

import (
	"sort"
	"time"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/storage/money"
)

var (
	_ = (*sort.IntSlice)(nil)
	_ = (*time.Time)(nil)
	_ = (*money.Money)(nil)
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
	TableCollection = csdb.NewTableManager(
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

// ErrIDNotFoundTableStoreSlice gets returned when an ID cannot be found in a TableStoreSlice slice.
// Generated via tableToStruct.
var ErrIDNotFoundTableStoreSlice = csdb.NewError("ID not found in TableStoreSlice")

// FindByStoreID searches the primary keys and returns a
// *TableStore if found or an error.
// Generated via tableToStruct.
func (s TableStoreSlice) FindByStoreID(
	store_id int64,
) (*TableStore, error) {
	for _, u := range s {
		if u != nil && u.StoreID == store_id {
			return u, nil
		}
	}
	return nil, ErrIDNotFoundTableStoreSlice
}

// FindByCode searches through this unique key and returns
// a *TableStore if found or an error.
// Generated via tableToStruct.
func (s TableStoreSlice) FindByCode(code string) (*TableStore, error) {
	for _, u := range s {
		if u != nil && u.Code.String == code {
			return u, nil
		}
	}
	return nil, ErrIDNotFoundTableStoreSlice
}

var _ sort.Interface = (*TableStoreSlice)(nil)

// Len returns the length and  will satisfy the sort.Interface.
// Generated via tableToStruct.
func (s TableStoreSlice) Len() int { return len(s) }

// Less will satisfy the sort.Interface and compares via
// the primary key.
// Generated via tableToStruct.
func (s TableStoreSlice) Less(i, j int) bool {
	return s[i].StoreID < s[j].StoreID && 1 == 1
}

// Swap will satisfy the sort.Interface.
// Generated via tableToStruct.
func (s TableStoreSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Sort will sort TableStoreSlice.
// Generated via tableToStruct.
func (s TableStoreSlice) Sort() { sort.Sort(s) }

// FilterThis filters the current slice by predicate f without memory allocation.
// Generated via tableToStruct.
func (s TableStoreSlice) FilterThis(f func(*TableStore) bool) TableStoreSlice {
	b := s[:0]
	for _, x := range s {
		if f(x) {
			b = append(b, x)
		}
	}
	return b
}

// Filter returns a new slice filtered by predicate f.
// Generated via tableToStruct.
func (s TableStoreSlice) Filter(f func(*TableStore) bool) TableStoreSlice {
	sl := make(TableStoreSlice, 0, len(s))
	for _, w := range s {
		if f(w) {
			sl = append(sl, w)
		}
	}
	return sl
}

// FilterNot will return a new TableStoreSlice that does not match
// by calling the function f
// Generated via tableToStruct.
func (s TableStoreSlice) FilterNot(f func(*TableStore) bool) TableStoreSlice {
	sl := make(TableStoreSlice, 0, len(s))
	for _, v := range s {
		if f(v) == false {
			sl = append(sl, v)
		}
	}
	return sl
}

// parentMap will run function f on all items in TableStoreSlice.
// Generated via tableToStruct.
func (s TableStoreSlice) Map(f func(*TableStore)) TableStoreSlice {
	for i := range s {
		f(s[i])
	}
	return s
}

// Cut will remove items i through j-1.
// Generated via tableToStruct.
func (s *TableStoreSlice) Cut(i, j int) {
	z := *s // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this should avoid the memory leak
	}
	z = z[:len(z)-j+i]
	*s = z
}

// Delete will remove an item from the slice.
// Generated via tableToStruct.
func (s *TableStoreSlice) Delete(i int) {
	z := *s // copy the slice header
	end := len(z) - 1
	s.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	*s = z
}

// Insert will place a new item at position i.
// Generated via tableToStruct.
func (s *TableStoreSlice) Insert(n *TableStore, i int) {
	z := *s // copy the slice header
	z = append(z, &TableStore{})
	copy(z[i+1:], z[i:])
	z[i] = n
	*s = z
}

// Append will add a new item at the end of TableStoreSlice.
// Generated via tableToStruct.
func (s *TableStoreSlice) Append(n ...*TableStore) {
	*s = append(*s, n...)
}

// Prepend will add a new item at the beginning of TableStoreSlice.
// Generated via tableToStruct.
func (s *TableStoreSlice) Prepend(n *TableStore) {
	s.Insert(n, 0)
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
				if v != nil {
					ext = append(ext, v.StoreID)
				}
			}
			return ext
		},
		Code: func() []string {
			ext := make([]string, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.Code.String)
				}
			}
			return ext
		},
		WebsiteID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.WebsiteID)
				}
			}
			return ext
		},
		GroupID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.GroupID)
				}
			}
			return ext
		},
		Name: func() []string {
			ext := make([]string, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.Name)
				}
			}
			return ext
		},
		SortOrder: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.SortOrder)
				}
			}
			return ext
		},
		IsActive: func() []bool {
			ext := make([]bool, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.IsActive)
				}
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

// ErrIDNotFoundTableGroupSlice gets returned when an ID cannot be found in a TableGroupSlice slice.
// Generated via tableToStruct.
var ErrIDNotFoundTableGroupSlice = csdb.NewError("ID not found in TableGroupSlice")

// FindByGroupID searches the primary keys and returns a
// *TableGroup if found or an error.
// Generated via tableToStruct.
func (s TableGroupSlice) FindByGroupID(
	group_id int64,
) (*TableGroup, error) {
	for _, u := range s {
		if u != nil && u.GroupID == group_id {
			return u, nil
		}
	}
	return nil, ErrIDNotFoundTableGroupSlice
}

var _ sort.Interface = (*TableGroupSlice)(nil)

// Len returns the length and  will satisfy the sort.Interface.
// Generated via tableToStruct.
func (s TableGroupSlice) Len() int { return len(s) }

// Less will satisfy the sort.Interface and compares via
// the primary key.
// Generated via tableToStruct.
func (s TableGroupSlice) Less(i, j int) bool {
	return s[i].GroupID < s[j].GroupID && 1 == 1
}

// Swap will satisfy the sort.Interface.
// Generated via tableToStruct.
func (s TableGroupSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Sort will sort TableGroupSlice.
// Generated via tableToStruct.
func (s TableGroupSlice) Sort() { sort.Sort(s) }

// FilterThis filters the current slice by predicate f without memory allocation.
// Generated via tableToStruct.
func (s TableGroupSlice) FilterThis(f func(*TableGroup) bool) TableGroupSlice {
	b := s[:0]
	for _, x := range s {
		if f(x) {
			b = append(b, x)
		}
	}
	return b
}

// Filter returns a new slice filtered by predicate f.
// Generated via tableToStruct.
func (s TableGroupSlice) Filter(f func(*TableGroup) bool) TableGroupSlice {
	sl := make(TableGroupSlice, 0, len(s))
	for _, w := range s {
		if f(w) {
			sl = append(sl, w)
		}
	}
	return sl
}

// FilterNot will return a new TableGroupSlice that does not match
// by calling the function f
// Generated via tableToStruct.
func (s TableGroupSlice) FilterNot(f func(*TableGroup) bool) TableGroupSlice {
	sl := make(TableGroupSlice, 0, len(s))
	for _, v := range s {
		if f(v) == false {
			sl = append(sl, v)
		}
	}
	return sl
}

// parentMap will run function f on all items in TableGroupSlice.
// Generated via tableToStruct.
func (s TableGroupSlice) Map(f func(*TableGroup)) TableGroupSlice {
	for i := range s {
		f(s[i])
	}
	return s
}

// Cut will remove items i through j-1.
// Generated via tableToStruct.
func (s *TableGroupSlice) Cut(i, j int) {
	z := *s // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this should avoid the memory leak
	}
	z = z[:len(z)-j+i]
	*s = z
}

// Delete will remove an item from the slice.
// Generated via tableToStruct.
func (s *TableGroupSlice) Delete(i int) {
	z := *s // copy the slice header
	end := len(z) - 1
	s.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	*s = z
}

// Insert will place a new item at position i.
// Generated via tableToStruct.
func (s *TableGroupSlice) Insert(n *TableGroup, i int) {
	z := *s // copy the slice header
	z = append(z, &TableGroup{})
	copy(z[i+1:], z[i:])
	z[i] = n
	*s = z
}

// Append will add a new item at the end of TableGroupSlice.
// Generated via tableToStruct.
func (s *TableGroupSlice) Append(n ...*TableGroup) {
	*s = append(*s, n...)
}

// Prepend will add a new item at the beginning of TableGroupSlice.
// Generated via tableToStruct.
func (s *TableGroupSlice) Prepend(n *TableGroup) {
	s.Insert(n, 0)
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
				if v != nil {
					ext = append(ext, v.GroupID)
				}
			}
			return ext
		},
		WebsiteID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.WebsiteID)
				}
			}
			return ext
		},
		Name: func() []string {
			ext := make([]string, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.Name)
				}
			}
			return ext
		},
		RootCategoryID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.RootCategoryID)
				}
			}
			return ext
		},
		DefaultStoreID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.DefaultStoreID)
				}
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

// ErrIDNotFoundTableWebsiteSlice gets returned when an ID cannot be found in a TableWebsiteSlice slice.
// Generated via tableToStruct.
var ErrIDNotFoundTableWebsiteSlice = csdb.NewError("ID not found in TableWebsiteSlice")

// FindByWebsiteID searches the primary keys and returns a
// *TableWebsite if found or an error.
// Generated via tableToStruct.
func (s TableWebsiteSlice) FindByWebsiteID(
	website_id int64,
) (*TableWebsite, error) {
	for _, u := range s {
		if u != nil && u.WebsiteID == website_id {
			return u, nil
		}
	}
	return nil, ErrIDNotFoundTableWebsiteSlice
}

// FindByCode searches through this unique key and returns
// a *TableWebsite if found or an error.
// Generated via tableToStruct.
func (s TableWebsiteSlice) FindByCode(code string) (*TableWebsite, error) {
	for _, u := range s {
		if u != nil && u.Code.String == code {
			return u, nil
		}
	}
	return nil, ErrIDNotFoundTableWebsiteSlice
}

var _ sort.Interface = (*TableWebsiteSlice)(nil)

// Len returns the length and  will satisfy the sort.Interface.
// Generated via tableToStruct.
func (s TableWebsiteSlice) Len() int { return len(s) }

// Less will satisfy the sort.Interface and compares via
// the primary key.
// Generated via tableToStruct.
func (s TableWebsiteSlice) Less(i, j int) bool {
	return s[i].WebsiteID < s[j].WebsiteID && 1 == 1
}

// Swap will satisfy the sort.Interface.
// Generated via tableToStruct.
func (s TableWebsiteSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Sort will sort TableWebsiteSlice.
// Generated via tableToStruct.
func (s TableWebsiteSlice) Sort() { sort.Sort(s) }

// FilterThis filters the current slice by predicate f without memory allocation.
// Generated via tableToStruct.
func (s TableWebsiteSlice) FilterThis(f func(*TableWebsite) bool) TableWebsiteSlice {
	b := s[:0]
	for _, x := range s {
		if f(x) {
			b = append(b, x)
		}
	}
	return b
}

// Filter returns a new slice filtered by predicate f.
// Generated via tableToStruct.
func (s TableWebsiteSlice) Filter(f func(*TableWebsite) bool) TableWebsiteSlice {
	sl := make(TableWebsiteSlice, 0, len(s))
	for _, w := range s {
		if f(w) {
			sl = append(sl, w)
		}
	}
	return sl
}

// FilterNot will return a new TableWebsiteSlice that does not match
// by calling the function f
// Generated via tableToStruct.
func (s TableWebsiteSlice) FilterNot(f func(*TableWebsite) bool) TableWebsiteSlice {
	sl := make(TableWebsiteSlice, 0, len(s))
	for _, v := range s {
		if f(v) == false {
			sl = append(sl, v)
		}
	}
	return sl
}

// parentMap will run function f on all items in TableWebsiteSlice.
// Generated via tableToStruct.
func (s TableWebsiteSlice) Map(f func(*TableWebsite)) TableWebsiteSlice {
	for i := range s {
		f(s[i])
	}
	return s
}

// Cut will remove items i through j-1.
// Generated via tableToStruct.
func (s *TableWebsiteSlice) Cut(i, j int) {
	z := *s // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this should avoid the memory leak
	}
	z = z[:len(z)-j+i]
	*s = z
}

// Delete will remove an item from the slice.
// Generated via tableToStruct.
func (s *TableWebsiteSlice) Delete(i int) {
	z := *s // copy the slice header
	end := len(z) - 1
	s.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	*s = z
}

// Insert will place a new item at position i.
// Generated via tableToStruct.
func (s *TableWebsiteSlice) Insert(n *TableWebsite, i int) {
	z := *s // copy the slice header
	z = append(z, &TableWebsite{})
	copy(z[i+1:], z[i:])
	z[i] = n
	*s = z
}

// Append will add a new item at the end of TableWebsiteSlice.
// Generated via tableToStruct.
func (s *TableWebsiteSlice) Append(n ...*TableWebsite) {
	*s = append(*s, n...)
}

// Prepend will add a new item at the beginning of TableWebsiteSlice.
// Generated via tableToStruct.
func (s *TableWebsiteSlice) Prepend(n *TableWebsite) {
	s.Insert(n, 0)
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
				if v != nil {
					ext = append(ext, v.WebsiteID)
				}
			}
			return ext
		},
		Code: func() []string {
			ext := make([]string, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.Code.String)
				}
			}
			return ext
		},
		Name: func() []string {
			ext := make([]string, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.Name.String)
				}
			}
			return ext
		},
		SortOrder: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.SortOrder)
				}
			}
			return ext
		},
		DefaultGroupID: func() []int64 {
			ext := make([]int64, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.DefaultGroupID)
				}
			}
			return ext
		},
		IsDefault: func() []bool {
			ext := make([]bool, 0, len(s))
			for _, v := range s {
				if v != nil {
					ext = append(ext, v.IsDefault.Bool)
				}
			}
			return ext
		},
	}
}
