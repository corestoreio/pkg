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

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/stretchr/testify/assert"
)

func TestWebsiteSlice_Map_Each(t *testing.T) {
	ws := store.WebsiteSlice{
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			store.TableGroupSlice{
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			},
			store.TableStoreSlice{
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
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
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), SortOrder: 4, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("uk"), SortOrder: 3, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 3, Code: dbr.NewNullString("ch"), SortOrder: 5, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
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
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), SortOrder: 4, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("uk"), SortOrder: 3, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 3, Code: dbr.NewNullString("ch"), SortOrder: 5, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
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
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), SortOrder: 4, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("uk"), SortOrder: 3, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			nil,
			nil,
		),
		store.MustNewWebsite(
			cfgmock.NewService(),
			&store.TableWebsite{WebsiteID: 3, Code: dbr.NewNullString("ch"), SortOrder: 5, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			nil,
			nil,
		),
	}
	assert.Exactly(t, []int64{1, 2, 3}, ws.IDs())
	assert.Nil(t, store.WebsiteSlice{}.IDs())
}
