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

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util"
	"github.com/stretchr/testify/assert"
)

var _ scope.WebsiteIDer = (*store.Website)(nil)
var _ scope.StoreIDer = (*store.Website)(nil)
var _ scope.GroupIDer = (*store.Website)(nil)
var _ scope.WebsiteCoder = (*store.Website)(nil)

func TestNewWebsite(t *testing.T) {
	w, err := store.NewWebsite(
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
	)
	assert.NoError(t, err)
	assert.Equal(t, "euro", w.Data.Code.String)

	dg, err := w.DefaultGroup()
	assert.Nil(t, dg)
	assert.EqualError(t, store.errWebsiteDefaultGroupNotFound, err.Error())

	ds, err := w.DefaultStore()
	assert.Nil(t, ds)
	assert.EqualError(t, store.errWebsiteDefaultGroupNotFound, err.Error())
	assert.Nil(t, w.Stores)
	assert.Nil(t, w.Groups)
}

func TestMustNewWebsite(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), store.errArgumentCannotBeNil.Error())
		}
	}()
	_ = store.MustNewWebsite(nil, nil)
}

func TestNewWebsiteSetGroupsStores(t *testing.T) {
	w, err := store.NewWebsite(
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
		store.SetWebsiteGroupsStores(
			store.TableGroupSlice{
				&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
				&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
			},
			store.TableStoreSlice{
				&store.TableStore{StoreID: 0, Code: dbr.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
				&store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 4, Code: dbr.NewNullString("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
				&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
				&store.TableStore{StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
			},
		),
	)
	assert.NoError(t, err)

	dg, err := w.DefaultGroup()
	assert.NotNil(t, dg)
	assert.EqualValues(t, "DACH Group", dg.Data.Name, "get default group: %#v", dg)
	assert.NoError(t, err)

	ds, err := w.DefaultStore()
	assert.NotNil(t, ds)
	assert.EqualValues(t, "at", ds.Data.Code.String, "get default store: %#v", ds)
	assert.NoError(t, err)

	assert.NotNil(t, dg.Stores)
	assert.EqualValues(t, util.StringSlice{"de", "at", "ch"}, dg.Stores.Codes())

	for _, st := range dg.Stores {
		assert.EqualValues(t, "DACH Group", st.Group.Data.Name)
		assert.EqualValues(t, "Europe", st.Website.Data.Name.String)
	}

	assert.NotNil(t, w.Stores)
	assert.EqualValues(t, util.StringSlice{"de", "uk", "at", "ch"}, w.Stores.Codes())

	assert.NotNil(t, w.Groups)
	assert.EqualValues(t, util.Int64Slice{1, 2}, w.Groups.IDs())

	assert.Exactly(t, int64(2), w.StoreID())
	assert.Exactly(t, int64(1), w.GroupID())
	assert.Equal(t, "euro", w.WebsiteCode())
}

func TestNewWebsiteStoreIDError(t *testing.T) {
	t.Parallel()
	w, err := store.NewWebsite(
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
	)
	assert.NoError(t, err)
	assert.Exactly(t, scope.UnavailableStoreID, w.StoreID())
}

func TestNewWebsiteSetGroupsStoresError1(t *testing.T) {
	t.Parallel()
	w, err := store.NewWebsite(
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
		store.SetWebsiteGroupsStores(
			store.TableGroupSlice{
				0: &store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			},
			store.TableStoreSlice{
				&store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 4, Code: dbr.NewNullString("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
				&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
				&store.TableStore{StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
			},
		),
	)
	assert.Nil(t, w)
	assert.Contains(t, err.Error(), "Integrity error")
}

// TODO
//func getWebsiteBaseCurrency(priceScope int, curGlobal, curWebsite string) (*store.Website, error) {
//	return store.NewWebsite(
//		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
//		store.SetWebsiteGroupsStores(
//			store.TableGroupSlice{
//				0: &store.TableGroup{GroupID: 0, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 1},
//			},
//			store.TableStoreSlice{
//				0: &store.TableStore{StoreID: 0, Code: dbr.NewNullString("Admin"), WebsiteID: 1, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
//				1: &store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 0, Name: "Germany", SortOrder: 10, IsActive: true},
//			},
//		),
//		store.SetWebsiteConfig(
//			cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
//				catconfig.Backend.CatalogPriceScope.FQPathInt64(scope.StrDefault, 0):    priceScope,
//				directory.Backend.CurrencyOptionsBase.FQPathInt64(scope.StrDefault, 0):  curGlobal,
//				directory.Backend.CurrencyOptionsBase.FQPathInt64(scope.StrWebsites, 1): curWebsite,
//			})),
//		),
//	)
//}
//
//func TestWebsiteBaseCurrency(t *testing.T) {
//	t.Parallel()
//	tests := []struct {
//		priceScope int
//		curGlobal  string
//		curWebsite string
//		curWant    string
//		wantErr    error
//	}{
//		{catconfig.PriceScopeGlobal, "USD", "EUR", "USD", nil},
//		{catconfig.PriceScopeGlobal, "ZZ", "EUR", "XXX", errors.New("currency: tag is not well-formed")},
//		{catconfig.PriceScopeWebsite, "USD", "EUR", "EUR", nil},
//		{catconfig.PriceScopeWebsite, "USD", "YYY", "XXX", errors.New("currency: tag is not a recognized currency")},
//	}
//
//	for _, test := range tests {
//		w, err := getWebsiteBaseCurrency(test.priceScope, test.curGlobal, test.curWebsite)
//		assert.NoError(t, err)
//		if false == assert.NotNil(t, w) {
//			t.Fatal("website is nil")
//		}
//
//		haveCur, haveErr := w.BaseCurrency()
//
//		if test.wantErr != nil {
//			assert.EqualError(t, haveErr, test.wantErr.Error())
//			assert.Exactly(t, test.curWant, haveCur.Unit.String())
//			continue
//		}
//
//		assert.NoError(t, haveErr)
//
//		wantCur, err := directory.NewCurrencyISO(test.curWant)
//		assert.NoError(t, err)
//		assert.Exactly(t, wantCur, haveCur)
//	}
//}

func TestTableWebsiteSlice(t *testing.T) {
	t.Parallel()
	websites := store.TableWebsiteSlice{
		0: &store.TableWebsite{WebsiteID: 0, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
		1: &store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
		2: nil,
		3: &store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
	}
	assert.True(t, websites.Len() == 4)

	w1, err := websites.FindByWebsiteID(999)
	assert.Nil(t, w1)
	assert.EqualError(t, store.ErrIDNotFoundTableWebsiteSlice, err.Error())

	w2, err := websites.FindByWebsiteID(2)
	assert.NotNil(t, w2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), w2.WebsiteID)

	w3, err := websites.FindByCode("euro")
	assert.NotNil(t, w3)
	assert.NoError(t, err)
	assert.Equal(t, "euro", w3.Code.String)

	w4, err := websites.FindByCode("corestore")
	assert.Nil(t, w4)
	assert.EqualError(t, store.ErrIDNotFoundTableWebsiteSlice, err.Error())

	wf1 := websites.Filter(func(w *store.TableWebsite) bool {
		return w != nil && w.WebsiteID == 1
	})
	assert.EqualValues(t, "Europe", wf1[0].Name.String)
}

func TestTableWebsiteSliceLoad(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()
	dbrSess := dbc.NewSession()

	var websites store.TableWebsiteSlice
	_, err := websites.SQLSelect(dbrSess)
	assert.NoError(t, err)

	assert.True(t, websites.Len() >= 2)
	for _, s := range websites {
		assert.True(t, len(s.Code.String) > 1)
	}
}
