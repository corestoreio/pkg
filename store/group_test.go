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
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	dbc := csdb.MustConnectTest()
	defer func() {
		if err := dbc.Close(); err != nil {
			panic(err)
		}
	}()
	if err := store.TableCollection.Init(dbc.NewSession()); err != nil {
		panic(err)
	}
}

func TestNewGroup(t *testing.T) {
	g, err := store.NewGroup(
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		nil,
	)
	assert.NoError(t, err)
	assert.EqualValues(t, "DACH Group", g.Data.Name)
	assert.Nil(t, g.Stores)

	gStores2, err := g.DefaultStore()
	assert.Nil(t, gStores2)
	assert.EqualError(t, store.ErrGroupDefaultStoreNotFound, err.Error())
}

func TestNewGroupErrorArgument(t *testing.T) {
	ng, err := store.NewGroup(nil, nil)
	assert.Nil(t, ng)
	assert.EqualError(t, store.ErrArgumentCannotBeNil, err.Error())
}

func TestNewGroupPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), store.ErrArgumentCannotBeNil.Error())
		}
	}()
	_ = store.MustNewGroup(nil, nil)
}

func TestNewGroupPanicWebsiteIncorrect(t *testing.T) {
	ng, err := store.NewGroup(
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		store.SetGroupWebsite(&store.TableWebsite{WebsiteID: 2, Code: csdb.NewNullString("oz"), Name: csdb.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: csdb.NewNullBool(false)}),
	)
	assert.Nil(t, ng)
	assert.EqualError(t, store.ErrGroupWebsiteNotFound, err.Error())
}

func TestNewGroupSetStoresErrorWebsiteIsNil(t *testing.T) {
	g, err := store.NewGroup(
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		nil,
		store.SetGroupStores(
			store.TableStoreSlice{
				&store.TableStore{StoreID: 0, Code: csdb.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
			},
			nil,
		),
	)
	assert.Nil(t, g)
	assert.EqualError(t, store.ErrGroupWebsiteNotFound, err.Error())
}

func TestNewGroupSetStoresErrorWebsiteIncorrect(t *testing.T) {
	g, err := store.NewGroup(
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		store.SetGroupStores(
			store.TableStoreSlice{
				&store.TableStore{StoreID: 0, Code: csdb.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
			},
			&store.TableWebsite{WebsiteID: 2, Code: csdb.NewNullString("oz"), Name: csdb.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: csdb.NewNullBool(false)},
		),
	)
	assert.Nil(t, g)
	assert.EqualError(t, store.ErrGroupWebsiteIntegrityFailed, err.Error())
}

func TestNewGroupSetStores(t *testing.T) {

	g, err := store.NewGroup(
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		nil, // this is just for to confuse the NewGroup ApplyOption function
		store.SetGroupStores(
			store.TableStoreSlice{
				&store.TableStore{StoreID: 0, Code: csdb.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
				&store.TableStore{StoreID: 5, Code: csdb.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 1, Code: csdb.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 4, Code: csdb.NewNullString("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 2, Code: csdb.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
				&store.TableStore{StoreID: 6, Code: csdb.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
				&store.TableStore{StoreID: 3, Code: csdb.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
			},
			&store.TableWebsite{WebsiteID: 1, Code: csdb.NewNullString("euro"), Name: csdb.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: csdb.NewNullBool(true)},
		),
	)
	assert.NoError(t, err)

	assert.NotNil(t, g.Stores)
	assert.EqualValues(t, utils.StringSlice{"de", "at", "ch"}, g.Stores.Codes())

	gDefaultStore, err := g.DefaultStore()
	assert.NoError(t, err)
	assert.EqualValues(t, "euro", gDefaultStore.Website.Data.Code.String)
	assert.EqualValues(t, "DACH Group", gDefaultStore.Group.Data.Name)
	assert.EqualValues(t, "at", gDefaultStore.Data.Code.String)
}

var testGroups = store.TableGroupSlice{
	&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
	&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
	&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
	&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
}

func TestTableGroupSliceLoad(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()
	dbrSess := dbc.NewSession()

	var groups store.TableGroupSlice
	rows, err := groups.SQLSelect(dbrSess)
	assert.NoError(t, err)
	assert.True(t, rows > 0)

	assert.True(t, groups.Len() >= 2)
	for _, s := range groups {
		assert.True(t, len(s.Name) > 1)
	}
}

func TestTableGroupSliceIDs(t *testing.T) {
	assert.EqualValues(t, []int64{3, 1, 0, 2}, testGroups.Extract().GroupID())
	assert.True(t, testGroups.Len() == 4)
}

func TestTableGroupSliceFindByID(t *testing.T) {
	g1, err := testGroups.FindByGroupID(999)
	assert.EqualError(t, err, "ID not found in TableGroupSlice")
	assert.Nil(t, g1)

	g2, err := testGroups.FindByGroupID(3)
	assert.NoError(t, err)
	assert.EqualValues(t, "Australia", g2.Name)
}
