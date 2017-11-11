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
	"testing"

	"github.com/corestoreio/cspkg/config/cfgmock"
	"github.com/corestoreio/cspkg/store"
	"github.com/corestoreio/cspkg/util/null"
	"github.com/stretchr/testify/assert"
)

func TestStoreSlice_Map_Each(t *testing.T) {
	ss := store.StoreSlice{
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("admin"), Name: null.StringFrom("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: null.BoolFrom(false)},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
		),
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 2, Code: null.StringFrom("ch"), WebsiteID: 1, GroupID: 1, Name: "Swiss", SortOrder: 20, IsActive: true},
			&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("admin"), Name: null.StringFrom("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: null.BoolFrom(false)},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
		),
	}

	ss.
		Map(func(s *store.Store) {
			s.Data.StoreID = 4
			s.Website.Data.WebsiteID = 2
		}).
		Each(func(s store.Store) {
			assert.Exactly(t, int64(2), s.Website.ID())
		})

	assert.Exactly(t, []int64{4, 4}, ss.IDs())
}

func TestStoreSlice_ActiveCodes(t *testing.T) {
	ss := store.StoreSlice{
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			nil,
			nil,
		),
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 2, Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: false},
			nil,
			nil,
		),
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 3, Code: null.StringFrom("ch"), WebsiteID: 1, GroupID: 1, Name: "Swiss", SortOrder: 20, IsActive: true},
			nil,
			nil,
		),
	}
	assert.Exactly(t, []string{"de", "ch"}, ss.ActiveCodes())
	assert.Nil(t, store.StoreSlice{}.ActiveCodes())
	fs, ok := ss.FindOne(func(s store.Store) bool {
		return s.Code() == "at"
	})
	assert.True(t, ok)
	assert.Exactly(t, "at", fs.Code())

	fs, ok = ss.FindOne(func(s store.Store) bool {
		return s.Code() == "xx"
	})
	assert.False(t, ok)
	assert.Exactly(t, "", fs.Code())
}

func TestStoreSlice_ActiveIDs(t *testing.T) {
	ss := store.StoreSlice{
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			nil,
			nil,
		),
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 2, Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: false},
			nil,
			nil,
		),
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 3, Code: null.StringFrom("ch"), WebsiteID: 1, GroupID: 1, Name: "Swiss", SortOrder: 20, IsActive: true},
			nil,
			nil,
		),
	}
	assert.Exactly(t, []int64{1, 3}, ss.ActiveIDs())
	assert.Nil(t, store.StoreSlice{}.ActiveIDs())
}

func TestStoreSlice_Sort(t *testing.T) {
	ss := store.StoreSlice{
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 2, IsActive: true},
			nil,
			nil,
		),
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 2, Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 1, IsActive: false},
			nil,
			nil,
		),
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 3, Code: null.StringFrom("ch"), WebsiteID: 1, GroupID: 1, Name: "Swiss", SortOrder: 3, IsActive: true},
			nil,
			nil,
		),
	}
	ss.Sort()
	assert.Exactly(t, []int64{2, 1, 3}, ss.IDs())
}

var benchmarkStoreSlice_Filter store.StoreSlice

func BenchmarkStoreSlice_Filter(b *testing.B) {
	const count = 1000
	cfg := cfgmock.NewService()
	stores := make(store.StoreSlice, count)
	for i := 0; i < count; i++ {
		stores[i] = store.MustNewStore(cfg,
			&store.TableStore{StoreID: int64(i), Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 1, IsActive: (i % 2) == 0},
			nil, nil)
	}
	f := func(s store.Store) bool {
		return s.IsActive()
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkStoreSlice_Filter = stores.Filter(f)
	}
	if have, want := len(benchmarkStoreSlice_Filter), count/2; have != want {
		b.Errorf("Have: %v Want: %v", have, want)
	}

}
