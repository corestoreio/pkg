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

package store_test

import (
	"database/sql"
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewStore(t *testing.T) {

	tests := []struct {
		w *store.TableWebsite
		g *store.TableGroup
		s *store.TableStore
	}{
		{
			w: &store.TableWebsite{WebsiteID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Admin", Valid: true}}, SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
			g: &store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			s: &store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		},
		{
			w: &store.TableWebsite{WebsiteID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
			g: &store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
			s: &store.TableStore{StoreID: 5, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
		},
	}
	for _, test := range tests {
		s := store.NewStore(test.w, test.g, test.s)
		assert.NotNil(t, s)
		assert.EqualValues(t, test.w.WebsiteID, s.Website.Data().WebsiteID)
		assert.EqualValues(t, test.g.GroupID, s.Group.Data().GroupID)
		assert.EqualValues(t, test.s.Code, s.Data().Code)
		assert.Nil(t, s.Group.Website)
		gStores, gErr := s.Group.Stores()
		assert.Nil(t, gStores)
		assert.EqualError(t, store.ErrGroupStoresNotAvailable, gErr.Error())
	}
}

func TestNewStorePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.EqualError(t, store.ErrStoreNewArgNil, err.Error())
		}
	}()
	_ = store.NewStore(nil, nil, nil)
}

var testStores = store.TableStoreSlice{
	&store.TableStore{StoreID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
	nil,
	&store.TableStore{StoreID: 5, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
	&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
	&store.TableStore{StoreID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "uk", Valid: true}}, WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
	&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
	&store.TableStore{StoreID: 6, Code: dbr.NullString{NullString: sql.NullString{String: "nz", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
	&store.TableStore{StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
	nil,
}

func TestTableStoreSliceLoad(t *testing.T) {
	db := csdb.MustConnectTest()
	defer db.Close()
	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)
	var stores store.TableStoreSlice
	stores.Load(dbrSess)
	assert.True(t, stores.Len() > 2)
	for _, s := range stores {
		assert.True(t, len(s.Code.String) > 1)
	}
}

func TestTableStoreSliceFindByID(t *testing.T) {
	eLen := 9
	assert.True(t, testStores.Len() == eLen, "Length of TableStoreSlice is not %d", eLen)

	s1, err := testStores.FindByID(999)
	assert.Nil(t, s1)
	assert.EqualError(t, store.ErrStoreNotFound, err.Error())

	s2, err := testStores.FindByID(6)
	assert.NotNil(t, s2)
	assert.NoError(t, err)
	assert.Equal(t, 6, s2.StoreID)
}

func TestTableStoreSliceFilterByGroupID(t *testing.T) {
	gStores := testStores.FilterByGroupID(3)
	assert.NotNil(t, gStores)
	assert.Len(t, gStores, 2)
	gStores2 := testStores.FilterByGroupID(32)
	assert.Nil(t, gStores2)
	assert.Len(t, gStores2, 0)
}

func TestTableStoreSliceFilterByWebsiteID(t *testing.T) {
	gStores := testStores.FilterByWebsiteID(0)
	assert.NotNil(t, gStores)
	assert.Len(t, gStores, 1)
	gStores2 := testStores.FilterByWebsiteID(32)
	assert.Nil(t, gStores2)
	assert.Len(t, gStores2, 0)

	var ts = store.TableStoreSlice{}
	assert.Nil(t, ts.FilterByGroupID(2))
}

func TestTableStoreSliceCodes(t *testing.T) {
	codes := testStores.Codes()
	assert.NotNil(t, codes)
	assert.Equal(t, utils.StringSlice{"admin", "au", "de", "uk", "at", "nz", "ch"}, codes)

	var ts = store.TableStoreSlice{}
	assert.Nil(t, ts.Codes())
}

func TestTableStoreSliceIDs(t *testing.T) {
	ids := testStores.IDs()
	assert.NotNil(t, ids)
	assert.Equal(t, utils.Int64Slice{0, 5, 1, 4, 2, 6, 3}, ids)

	var ts = store.TableStoreSlice{}
	assert.Nil(t, ts.IDs())
}
