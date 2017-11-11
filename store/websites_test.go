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

package store_test

import (
	"encoding/json"
	"testing"

	"github.com/corestoreio/cspkg/config/cfgmock"
	"github.com/corestoreio/cspkg/store"
	"github.com/corestoreio/cspkg/util/null"
	"github.com/stretchr/testify/assert"
)

func TestWebsiteSlice_Map_Each(t *testing.T) {
	ws := store.WebsiteSlice{
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			store.TableGroupSlice{
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			},
			store.TableStoreSlice{
				&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 2, Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			},
		),
	}

	ws.
		Map(func(w *store.Website) {
			w.Data.WebsiteID = 4
			w.Groups.Map(func(g *store.Group) {
				g.Data.Name = "Gopher"
			})
		}).
		Each(func(w store.Website) {
			w.Groups.Each(func(g store.Group) {
				assert.Exactly(t, "Gopher", g.Name())
			})

		})
	assert.Exactly(t, []int64{4}, ws.IDs())
}

func TestWebsiteSlice_Sort(t *testing.T) {
	ws := store.WebsiteSlice{
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), SortOrder: 4, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("uk"), SortOrder: 3, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 3, Code: null.StringFrom("ch"), SortOrder: 5, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			nil,
			nil,
		),
	}
	ws.Sort()
	assert.Exactly(t, []int64{2, 1, 3}, ws.IDs())
}

func TestWebsiteSlice_Codes(t *testing.T) {
	ws := store.WebsiteSlice{
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), SortOrder: 4, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("uk"), SortOrder: 3, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 3, Code: null.StringFrom("ch"), SortOrder: 5, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			nil,
			nil,
		),
	}
	assert.Exactly(t, []string{"euro", "uk", "ch"}, ws.Codes())
	assert.Nil(t, store.WebsiteSlice{}.Codes())
}

func TestWebsiteSlice_IDs(t *testing.T) {
	ws := store.WebsiteSlice{
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), SortOrder: 4, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("uk"), SortOrder: 3, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 3, Code: null.StringFrom("ch"), SortOrder: 5, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			nil,
			nil,
		),
	}
	assert.Exactly(t, []int64{1, 2, 3}, ws.IDs())
	assert.Nil(t, store.WebsiteSlice{}.IDs())
}

var treeStoreSrv = store.MustNewService(
	cfgmock.NewService(),
	store.WithTableWebsites(
		&store.TableWebsite{WebsiteID: 0, Code: null.StringFrom("admin"), Name: null.StringFrom("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: null.BoolFrom(false)},
		&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
		&store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(false)},
	),
	store.WithTableGroups(
		&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
		&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
	),
	store.WithTableStores(
		&store.TableStore{StoreID: 0, Code: null.StringFrom("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
		&store.TableStore{StoreID: 5, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
		&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableStore{StoreID: 4, Code: null.StringFrom("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
		&store.TableStore{StoreID: 2, Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
		&store.TableStore{StoreID: 6, Code: null.StringFrom("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		&store.TableStore{StoreID: 3, Code: null.StringFrom("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
	),
)

func TestScopeTree(t *testing.T) {

	tr := treeStoreSrv.Websites().Tree()

	j, err := json.Marshal(tr)
	if err != nil {
		t.Fatal(err)
	}

	const want = `{"scope":"Default","id":0,"scopes":[{"scope":"Website","id":0,"scopes":[{"scope":"Group","id":0,"scopes":[{"scope":"Store","id":0}]}]},{"scope":"Website","id":1,"scopes":[{"scope":"Group","id":1,"scopes":[{"scope":"Store","id":1},{"scope":"Store","id":2},{"scope":"Store","id":3}]},{"scope":"Group","id":2,"scopes":[{"scope":"Store","id":4}]}]},{"scope":"Website","id":2,"scopes":[{"scope":"Group","id":3,"scopes":[{"scope":"Store","id":5},{"scope":"Store","id":6}]}]}]}`
	assert.Exactly(t, want, string(j))
}

var benchmarkScopeTree store.Tree

// 1000000	      1406 ns/op	     608 B/op	       8 allocs/op
func BenchmarkScopeTree(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkScopeTree = treeStoreSrv.Websites().Tree()
	}
}
