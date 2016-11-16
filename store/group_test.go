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
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/null"
	"github.com/stretchr/testify/assert"
)

func TestNewGroup(t *testing.T) {

	g, err := store.NewGroup(
		cfgmock.NewService(),
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		nil,
		nil,
	)
	assert.NoError(t, err)
	assert.EqualValues(t, "DACH Group", g.Data.Name)
	assert.Nil(t, g.Stores)

	gStores2, err := g.DefaultStore()
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)

	err = gStores2.Validate()
	assert.True(t, errors.IsNotValid(err), "%+v", err)
}

func TestNewGroupErrorWebsiteIncorrect(t *testing.T) {

	ng, err := store.NewGroup(
		cfgmock.NewService(),
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		&store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(false)},
		nil,
	)
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	assert.NoError(t, ng.Validate())
}

func TestNewGroupSetStores_WebsiteIsNil(t *testing.T) {
	g, err := store.NewGroup(
		cfgmock.NewService(),
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		nil,
		store.TableStoreSlice{
			&store.TableStore{StoreID: 0, Code: null.StringFrom("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
		},
	)
	assert.False(t, errors.IsNotValid(err), "Error: %s", err)
	assert.NoError(t, g.Validate())
	assert.Exactly(t, int64(1), g.ID())
	assert.Exactly(t, []int64(nil), g.Stores.IDs())
	assert.Exactly(t, int64(-1), g.Website.ID())
}

func TestNewGroupSetStoresErrorWebsiteIncorrect(t *testing.T) {
	g, err := store.NewGroup(
		cfgmock.NewService(),
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		&store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(false)},
		store.TableStoreSlice{
			&store.TableStore{StoreID: 0, Code: null.StringFrom("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
		},
	)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
	assert.NoError(t, g.Validate())
}

func TestNewGroupSetStores(t *testing.T) {
	g := store.MustNewGroup(
		cfgmock.NewService(),
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
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

	assert.NotNil(t, g.Stores)
	assert.Exactly(t, []string{"de", "at", "ch"}, g.Stores.Codes())

	gDefaultStore, err := g.DefaultStore()
	assert.NoError(t, err)
	assert.Exactly(t, g.WebsiteID(), gDefaultStore.WebsiteID())
	assert.Exactly(t, g.ID(), gDefaultStore.GroupID())
	assert.Exactly(t, "at", gDefaultStore.Code())
}

var testGroups = store.TableGroupSlice{
	&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
	&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
	&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
	&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
}

func TestTableGroupSliceLoad(t *testing.T) {

	dbrCon, dbMock := cstesting.MockDB(t)
	dbMock.ExpectQuery("SELECT (.+) FROM `store_group`(.+) ORDER BY(.+)").WillReturnRows(
		cstesting.MustMockRows(cstesting.WithFile("testdata", "core_store_group_view.csv")),
	)

	// store.TableCollection already initialized

	var groups store.TableGroupSlice
	rows, err := groups.SQLSelect(dbrCon.NewSession())
	assert.NoError(t, err)

	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("%+v", err)
	}

	assert.Exactly(t, 9, rows)

	assert.Len(t, groups, 9)
	for _, s := range groups {
		assert.True(t, len(s.Name) > 1)
	}
}

func TestTableGroupSliceIDs(t *testing.T) {

	assert.EqualValues(t, []int64{3, 1, 0, 2}, testGroups.Extract().GroupID())
	assert.True(t, testGroups.Len() == 4)
}

func TestTableGroupSliceFindByID(t *testing.T) {

	g1, found := testGroups.FindByGroupID(999)
	assert.False(t, found, "ID not found in TableGroupSlice")
	assert.Nil(t, g1)

	g2, found := testGroups.FindByGroupID(3)
	assert.True(t, found)
	assert.EqualValues(t, "Australia", g2.Name)
}
