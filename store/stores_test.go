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

func TestStoreSlice_Map_Each(t *testing.T) {
	ss := store.StoreSlice{
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
		),
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 2, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Swiss", SortOrder: 20, IsActive: true},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
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
