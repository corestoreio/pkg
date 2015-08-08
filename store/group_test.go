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

func TestNewGroup(t *testing.T) {
	g := store.NewGroup(
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		nil,
	)
	assert.EqualValues(t, "DACH Group", g.Data().Name)

	gStores1, err := g.Stores()
	assert.Nil(t, gStores1)
	assert.EqualError(t, store.ErrGroupStoresNotAvailable, err.Error())

	gStores2, err := g.DefaultStore()
	assert.Nil(t, gStores2)
	assert.EqualError(t, store.ErrGroupDefaultStoreNotFound, err.Error())
}

func TestNewGroupPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.EqualError(t, store.ErrStoreNewArgNil, err.Error())
			} else {
				t.Errorf("Failed to convert to type error: %#v", err)
			}
		} else {
			t.Error("Cannot find panic")
		}
	}()
	_ = store.NewGroup(nil, nil)
}

func TestNewGroupPanicWebsiteIncorrect(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.EqualError(t, store.ErrGroupWebsiteNotFound, err.Error())
			} else {
				t.Errorf("Failed to convert to type error: %#v", err)
			}
		} else {
			t.Error("Cannot find panic")
		}
	}()

	_ = store.NewGroup(
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		store.SetGroupWebsite(&store.TableWebsite{WebsiteID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}}),
	)

}

func TestNewGroupSetStoresPanicWebsiteIsNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.EqualError(t, store.ErrGroupWebsiteNotFound, err.Error())
			} else {
				t.Errorf("Failed to convert to type error: %#v", err)
			}
		} else {
			t.Error("Cannot find panic")
		}
	}()

	g := store.NewGroup(
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		nil,
	)
	g.SetStores(
		store.TableStoreSlice{
			&store.TableStore{StoreID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
		},
		nil,
	)
}

func TestNewGroupSetStoresPanicWebsiteIncorrect(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.EqualError(t, store.ErrGroupWebsiteNotFound, err.Error())
			} else {
				t.Errorf("Failed to convert to type error: %#v", err)
			}
		} else {
			t.Error("Cannot find panic")
		}
	}()

	g := store.NewGroup(
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		nil,
	)

	g.SetStores(
		store.TableStoreSlice{
			&store.TableStore{StoreID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
		},
		&store.TableWebsite{WebsiteID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
	)
}

func TestNewGroupSetStores(t *testing.T) {

	g := store.NewGroup(
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		nil,
	)
	g.SetStores(
		store.TableStoreSlice{
			&store.TableStore{StoreID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
			&store.TableStore{StoreID: 5, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "uk", Valid: true}}, WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NullString{NullString: sql.NullString{String: "nz", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
			&store.TableStore{StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
		},
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
	)
	gStores, err := g.Stores()
	assert.NoError(t, err)
	assert.EqualValues(t, utils.StringSlice{"de", "at", "ch"}, gStores.Codes())

	gDefaultStore, err := g.DefaultStore()
	assert.NoError(t, err)
	assert.EqualValues(t, "euro", gDefaultStore.Website().Data().Code.String)
	assert.EqualValues(t, "DACH Group", gDefaultStore.Group().Data().Name)
	assert.EqualValues(t, "at", gDefaultStore.Data().Code.String)
}

var testGroups = store.TableGroupSlice{
	&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
	&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
	&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
	&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
}

func TestTableGroupSliceLoad(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	dbrSess := dbc.NewSession()

	var groups store.TableGroupSlice
	rows, err := groups.Load(dbrSess)
	assert.NoError(t, err)
	assert.True(t, rows > 0)

	assert.True(t, groups.Len() > 2)
	for _, s := range groups {
		assert.True(t, len(s.Name) > 1)
	}
}

func TestTableGroupSliceIDs(t *testing.T) {
	assert.EqualValues(t, utils.Int64Slice{3, 1, 0, 2}, testGroups.IDs())
	assert.True(t, testGroups.Len() == 4)
}

func TestTableGroupSliceFindByID(t *testing.T) {
	g1, err := testGroups.FindByID(999)
	assert.EqualError(t, store.ErrGroupNotFound, err.Error())
	assert.Nil(t, g1)

	g2, err := testGroups.FindByID(3)
	assert.NoError(t, err)
	assert.EqualValues(t, "Australia", g2.Name)
}
