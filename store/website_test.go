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

func TestNewWebsite(t *testing.T) {
	w := store.NewWebsite(
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
	)
	assert.Equal(t, "euro", w.Data().Code.String)

	dg, err := w.DefaultGroup()
	assert.Nil(t, dg)
	assert.EqualError(t, store.ErrWebsiteDefaultGroupNotFound, err.Error())

	stores, err := w.Stores()
	assert.Nil(t, stores)
	assert.EqualError(t, store.ErrWebsiteStoresNotAvailable, err.Error())

	groups, err := w.Groups()
	assert.Nil(t, groups)
	assert.EqualError(t, store.ErrWebsiteGroupsNotAvailable, err.Error())
}

func TestNewWebsiteSetGroupsStores(t *testing.T) {
	w := store.NewWebsite(
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
	)
	w.SetGroupsStores(
		store.TableGroupSlice{
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
		},
		store.TableStoreSlice{
			&store.TableStore{StoreID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
			&store.TableStore{StoreID: 5, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "uk", Valid: true}}, WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NullString{NullString: sql.NullString{String: "nz", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
			&store.TableStore{StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
		},
	)

	dg, err := w.DefaultGroup()
	assert.NotNil(t, dg)
	assert.EqualValues(t, "DACH Group", dg.Data().Name, "get default group: %#v", dg)
	assert.NoError(t, err)

	dgStores, err := dg.Stores()
	assert.NoError(t, err)
	assert.EqualValues(t, utils.StringSlice{"de", "at", "ch"}, dgStores.Codes())

	for _, st := range dgStores {
		assert.EqualValues(t, "DACH Group", st.Group().Data().Name)
		assert.EqualValues(t, "Europe", st.Website().Data().Name.String)
	}

	stores, err := w.Stores()
	assert.NoError(t, err)
	assert.NotNil(t, stores)
	assert.EqualValues(t, utils.StringSlice{"de", "uk", "at", "ch"}, stores.Codes())

	groups, err := w.Groups()
	assert.NoError(t, err)
	assert.NotNil(t, groups)
	assert.EqualValues(t, utils.Int64Slice{1, 2}, groups.IDs())

}

func TestNewWebsiteSetGroupsStoresPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if str, ok := r.(string); ok {
				assert.Contains(t, str, "Integrity error")
			} else {
				t.Errorf("Failed to convert to type error: %#v", str)
			}
		} else {
			t.Error("Cannot find panic")
		}
	}()
	w := store.NewWebsite(
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
	)
	w.SetGroupsStores(
		store.TableGroupSlice{
			&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
		},
		store.TableStoreSlice{
			&store.TableStore{StoreID: 5, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "uk", Valid: true}}, WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NullString{NullString: sql.NullString{String: "nz", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
			&store.TableStore{StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
		},
	)
}

func TestTableWebsiteSlice(t *testing.T) {
	websites := store.TableWebsiteSlice{
		&store.TableWebsite{WebsiteID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Admin", Valid: true}}, SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
		nil,
		&store.TableWebsite{WebsiteID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
	}
	assert.True(t, websites.Len() == 4)

	w1, err := websites.FindByID(999)
	assert.Nil(t, w1)
	assert.EqualError(t, store.ErrWebsiteNotFound, err.Error())

	w2, err := websites.FindByID(2)
	assert.NotNil(t, w2)
	assert.NoError(t, err)
	assert.Equal(t, 2, w2.WebsiteID)

	w3, err := websites.FindByCode("euro")
	assert.NotNil(t, w3)
	assert.NoError(t, err)
	assert.Equal(t, "euro", w3.Code.String)

	w4, err := websites.FindByCode("corestore")
	assert.Nil(t, w4)
	assert.EqualError(t, store.ErrWebsiteNotFound, err.Error())

	wf1 := websites.Filter(func(w *store.TableWebsite) bool {
		return w.WebsiteID == 1
	})
	assert.EqualValues(t, "Europe", wf1[0].Name.String)
}

func TestTableWebsiteSliceLoad(t *testing.T) {
	db := csdb.MustConnectTest()
	defer db.Close()
	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)
	var websites store.TableWebsiteSlice
	websites.Load(dbrSess)
	assert.True(t, websites.Len() > 2)
	for _, s := range websites {
		assert.True(t, len(s.Code.String) > 1)
	}
}
