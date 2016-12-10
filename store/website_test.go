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
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/null"
	"github.com/corestoreio/csfw/util/slices"
	"github.com/stretchr/testify/assert"
)

func TestNewWebsite(t *testing.T) {

	w, err := store.NewWebsite(
		cfgmock.NewService(),
		&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
		nil,
		nil,
	)
	assert.NoError(t, err)
	assert.Equal(t, "euro", w.Data.Code.String)

	dg, err := w.DefaultGroup()
	assert.Nil(t, dg.Validate())
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)

	ds, err := w.DefaultStore()
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	err = ds.Validate()
	assert.True(t, errors.IsNotValid(err))
	assert.Nil(t, w.Stores)
	assert.Nil(t, w.Groups)
}

func TestNewWebsiteSetGroupsStores(t *testing.T) {

	w := store.MustNewWebsite(
		cfgmock.NewService(),
		&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
		store.TableGroupSlice{
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
		},
		store.TableStoreSlice{
			&store.TableStore{StoreID: 0, Code: null.StringFrom("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
			&store.TableStore{StoreID: 5, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 4, Code: null.StringFrom("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 2, Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			&store.TableStore{StoreID: 6, Code: null.StringFrom("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
			&store.TableStore{StoreID: 3, Code: null.StringFrom("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
		},
	)

	ds, err := w.DefaultStore()
	assert.NotNil(t, ds)
	assert.Exactly(t, "at", ds.Code(), "get default store: %#v", ds)
	assert.NoError(t, err)

	dg, err := w.DefaultGroup()
	assert.NotNil(t, dg)
	assert.Exactly(t, "DACH Group", dg.Name(), "get default group: %#v", dg)
	assert.NoError(t, err)

	assert.NotNil(t, dg.Stores)
	assert.Exactly(t, []string{"de", "at", "ch"}, dg.Stores.Codes())

	for _, st := range dg.Stores {
		assert.Empty(t, st.Group.Name())
		assert.Empty(t, st.Website.Name())
	}

	assert.NotNil(t, w.Stores)
	assert.EqualValues(t, slices.String{"de", "uk", "at", "ch"}, w.Stores.Codes())

	assert.NotNil(t, w.Groups)
	assert.EqualValues(t, slices.Int64{1, 2}, w.Groups.IDs())

	dsi, err := w.DefaultStoreID()
	assert.NoError(t, err)
	assert.Exactly(t, int64(2), dsi)
	assert.Exactly(t, int64(1), w.DefaultGroupID())
	assert.Equal(t, "euro", w.Code())
}

func TestNewWebsiteStoreIDError(t *testing.T) {
	w, err := store.NewWebsite(
		cfgmock.NewService(),
		&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
		nil,
		nil,
	)
	assert.NoError(t, err)
	dsi, err := w.DefaultStoreID()
	assert.True(t, errors.IsNotFound(err), "%+v", err)
	assert.Empty(t, dsi)
}

func TestNewWebsiteSetGroupsStores_Filter_Invalid(t *testing.T) {

	w, err := store.NewWebsite(
		cfgmock.NewService(),
		&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
		store.TableGroupSlice{
			&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
		},
		store.TableStoreSlice{
			&store.TableStore{StoreID: 5, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 4, Code: null.StringFrom("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 2, Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			&store.TableStore{StoreID: 6, Code: null.StringFrom("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
			&store.TableStore{StoreID: 3, Code: null.StringFrom("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
		},
	)
	assert.NotNil(t, w)
	assert.NoError(t, err, "%+v", err)
	err = w.Validate()
	assert.NoError(t, err, "%+v", err)
	assert.Exactly(t, []int64(nil), w.Groups.IDs())
	assert.Exactly(t, []int64{1, 4, 2, 3}, w.Stores.IDs())
}

func TestTableWebsiteSlice(t *testing.T) {

	websites := store.TableWebsiteSlice{
		0: &store.TableWebsite{WebsiteID: 0, Code: null.StringFrom("admin"), Name: null.StringFrom("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: null.BoolFrom(false)},
		1: &store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
		2: nil,
		3: &store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(false)},
	}
	assert.True(t, websites.Len() == 4)

	w1, found := websites.FindByWebsiteID(999)
	assert.Nil(t, w1)
	assert.False(t, found)

	w2, found := websites.FindByWebsiteID(2)
	assert.NotNil(t, w2)
	assert.True(t, found)
	assert.Equal(t, int64(2), w2.WebsiteID)

	w3, found := websites.FindByCode("euro")
	assert.NotNil(t, w3)
	assert.True(t, found)
	assert.Equal(t, "euro", w3.Code.String)

	w4, found := websites.FindByCode("corestore")
	assert.Nil(t, w4)
	assert.False(t, found)

	wf1 := websites.Filter(func(w *store.TableWebsite) bool {
		return w != nil && w.WebsiteID == 1
	})
	assert.EqualValues(t, "Europe", wf1[0].Name.String)
}

func TestTableWebsiteSliceLoad(t *testing.T) {
	t.Skip("TODO")

	//dbrCon, dbMock := cstesting.MockDB(t)
	//dbMock.ExpectQuery("SELECT (.+) FROM `store_website`(.+) ORDER BY(.+)").WillReturnRows(
	//	cstesting.MustMockRows(cstesting.WithFile("testdata", "core_website_view.csv")),
	//)
	//
	//// store.TableCollection already initialized
	//
	//var websites store.TableWebsiteSlice
	//rows, err := websites.SQLSelect(dbrCon.NewSession())
	//assert.NoError(t, err)
	//
	//if err := dbMock.ExpectationsWereMet(); err != nil {
	//	t.Fatalf("%+v", err)
	//}
	//assert.Exactly(t, 9, rows)
	//assert.Len(t, websites, 9)
	//for _, s := range websites {
	//	assert.True(t, len(s.Name.String) > 1)
	//}
}
