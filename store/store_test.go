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
	"encoding/json"
	"testing"

	"github.com/corestoreio/csfw/backend"
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util"
	"github.com/stretchr/testify/assert"
)

const TODO_Better_Test_Data = "@todo implement better test data which is equal for each Magento version"

func TestNewStore(t *testing.T) {
	defer debugLogBuf.Reset()
	tests := []struct {
		w *store.TableWebsite
		g *store.TableGroup
		s *store.TableStore
	}{
		{
			w: &store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
			g: &store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			s: &store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		},
		{
			w: &store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
			g: &store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
			s: &store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
		},
	}
	for _, test := range tests {
		s, err := store.NewStore(test.s, test.w, test.g)
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.EqualValues(t, test.w.WebsiteID, s.Website.Data.WebsiteID)
		assert.EqualValues(t, test.g.GroupID, s.Group.Data.GroupID)
		assert.EqualValues(t, test.s.Code, s.Data.Code)
		assert.NotNil(t, s.Group.Website)
		assert.NotEmpty(t, s.Group.Website.WebsiteID())
		assert.Nil(t, s.Group.Stores)
		assert.EqualValues(t, test.s.StoreID, s.StoreID())
		assert.EqualValues(t, test.s.GroupID, s.GroupID())
		assert.EqualValues(t, test.s.WebsiteID, s.WebsiteID())
	}
}

func TestNewStoreErrorArgsNil(t *testing.T) {
	s, err := store.NewStore(nil, nil, nil)
	assert.Nil(t, s)
	assert.EqualError(t, store.ErrArgumentCannotBeNil, err.Error())
}

func TestNewStoreErrorIncorrectGroup(t *testing.T) {
	s, err := store.NewStore(
		&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
		&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
	)
	assert.Nil(t, s)
	assert.EqualError(t, store.ErrStoreIncorrectGroup, err.Error())
}

func TestNewStoreErrorIncorrectWebsite(t *testing.T) {
	s, err := store.NewStore(
		&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
	)
	assert.Nil(t, s)
	assert.EqualError(t, store.ErrStoreIncorrectWebsite, err.Error())
}

func TestStoreSlice(t *testing.T) {
	defer debugLogBuf.Reset()
	storeSlice := store.StoreSlice{
		store.MustNewStore(
			&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
		),
		nil,
		store.MustNewStore(
			&store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		),
	}
	assert.True(t, storeSlice.Len() == 3)
	assert.EqualValues(t, util.Int64Slice{1, 5}, storeSlice.IDs())
	assert.EqualValues(t, util.StringSlice{"de", "au"}, storeSlice.Codes())
	assert.EqualValues(t, 5, storeSlice.LastItem().Data.StoreID)
	assert.Nil(t, (store.StoreSlice{}).LastItem())

	storeSlice2 := storeSlice.Filter(func(s *store.Store) bool {
		return s.Website.Data.WebsiteID == 2
	})
	assert.True(t, storeSlice2.Len() == 1)
	assert.Equal(t, "au", storeSlice2[0].Data.Code.String)
	assert.EqualValues(t, util.Int64Slice{5}, storeSlice2.IDs())
	assert.EqualValues(t, util.StringSlice{"au"}, storeSlice2.Codes())

	assert.Nil(t, (store.StoreSlice{}).IDs())
	assert.Nil(t, (store.StoreSlice{}).Codes())
}

var testStores = store.TableStoreSlice{
	&store.TableStore{StoreID: 0, Code: dbr.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
	nil,
	&store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
	&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
	&store.TableStore{StoreID: 4, Code: dbr.NewNullString("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
	&store.TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
	&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
	&store.TableStore{StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
	nil,
}

func TestTableStoreSliceLoad(t *testing.T) {
	defer debugLogBuf.Reset()
	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()
	dbrSess := dbc.NewSession()

	var stores store.TableStoreSlice
	_, err := stores.SQLSelect(dbrSess)
	assert.NoError(t, err)
	assert.True(t, stores.Len() >= 2) // @todo proper test data in database
	for _, s := range stores {
		assert.True(t, len(s.Code.String) > 1)
	}
}

func TestTableStoreSliceFindByID(t *testing.T) {
	eLen := 9
	assert.True(t, testStores.Len() == eLen, "Length of TableStoreSlice is not %d", eLen)

	s1, err := testStores.FindByStoreID(999)
	assert.Nil(t, s1)
	assert.EqualError(t, store.ErrIDNotFoundTableStoreSlice, err.Error())

	s2, err := testStores.FindByStoreID(6)
	assert.NotNil(t, s2)
	assert.NoError(t, err)
	assert.Equal(t, int64(6), s2.StoreID)
}

func TestTableStoreSliceFindByCode(t *testing.T) {

	s1, err := testStores.FindByCode("corestore")
	assert.Nil(t, s1)
	assert.EqualError(t, store.ErrIDNotFoundTableStoreSlice, err.Error())

	s2, err := testStores.FindByCode("ch")
	assert.NotNil(t, s2)
	assert.NoError(t, err)
	assert.Equal(t, "ch", s2.Code.String)
}

func TestTableStoreSliceFilterByGroupID(t *testing.T) {
	gStores := testStores.FilterByGroupID(3)
	assert.NotNil(t, gStores)
	assert.Len(t, gStores, 2)
	gStores2 := testStores.FilterByGroupID(32)
	assert.NotNil(t, gStores2)
	assert.Len(t, gStores2, 0)
}

func TestTableStoreSliceFilterByWebsiteID(t *testing.T) {
	gStores := testStores.FilterByWebsiteID(0)
	assert.NotNil(t, gStores)
	assert.Len(t, gStores, 1)
	gStores2 := testStores.FilterByWebsiteID(32)
	assert.NotNil(t, gStores2)
	assert.Len(t, gStores2, 0)

	var ts = store.TableStoreSlice{}
	tsRes := ts.FilterByGroupID(2)
	assert.NotNil(t, tsRes)
	assert.Len(t, tsRes, 0)
}

func TestTableStoreSliceCodes(t *testing.T) {

	t.Skip(TODO_Better_Test_Data)

	codes := testStores.Extract().Code()
	assert.NotNil(t, codes)
	assert.Equal(t, util.StringSlice{"admin", "au", "de", "uk", "at", "nz", "ch"}, codes)

	var ts = store.TableStoreSlice{}
	assert.Nil(t, ts.Extract().Code())
}

func TestTableStoreSliceIDs(t *testing.T) {

	t.Skip(TODO_Better_Test_Data)

	ids := testStores.Extract().StoreID()
	assert.NotNil(t, ids)
	assert.Equal(t, util.Int64Slice{0, 5, 1, 4, 2, 6, 3}, ids)

	var ts = store.TableStoreSlice{}
	assert.Nil(t, ts.Extract().StoreID())
}

func TestStoreBaseURLandPath(t *testing.T) {
	defer debugLogBuf.Reset()

	s, err := store.NewStore(
		&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 1},
	)
	assert.NoError(t, err)
	if s == nil {
		t.Fail()
	}

	tests := []struct {
		haveR        config.Getter
		haveUT       config.URLType
		haveIsSecure bool
		wantBaseUrl  string
		wantPath     string
	}{
		{
			config.NewMockGetter(config.WithMockString(
				func(path string) (string, error) {

					switch path {
					// scope is here store but config.ScopedGetter must fall back to default
					case backend.Backend.WebSecureBaseURL.FQPathInt64(scope.StrDefault, 0):
						return "https://corestore.io", nil
					case backend.Backend.WebUnsecureBaseURL.FQPathInt64(scope.StrDefault, 0):
						return "http://corestore.io", nil
					}
					return "", config.ErrKeyNotFound
				},
			)),
			config.URLTypeWeb, true, "https://corestore.io/", "/",
		},
		{
			config.NewMockGetter(config.WithMockString(
				func(path string) (string, error) {
					switch path {
					case backend.Backend.WebSecureBaseURL.FQPathInt64(scope.StrDefault, 0):
						return "https://myplatform.io/customer1", nil
					case backend.Backend.WebUnsecureBaseURL.FQPathInt64(scope.StrDefault, 0):
						return "http://myplatform.io/customer1", nil
					}
					return "", config.ErrKeyNotFound
				},
			)),
			config.URLTypeWeb, false, "http://myplatform.io/customer1/", "/customer1/",
		},
		{
			config.NewMockGetter(config.WithMockString(
				func(path string) (string, error) {
					switch path {
					case backend.Backend.WebSecureBaseURL.FQPathInt64(scope.StrDefault, 0):
						return model.PlaceholderBaseURL, nil
					case backend.Backend.WebUnsecureBaseURL.FQPathInt64(scope.StrDefault, 0):
						return model.PlaceholderBaseURL, nil
					case scope.StrDefault.FQPathInt64(0, config.PathCSBaseURL):
						return config.CSBaseURL, nil
					}
					return "", config.ErrKeyNotFound
				},
			)),
			config.URLTypeWeb, false, config.CSBaseURL, "/",
		},
	}

	for i, test := range tests {
		s.ApplyOptions(store.WithStoreConfig(test.haveR))
		assert.NotNil(t, s.Config, "Index %d", i)
		baseURL, err := s.BaseURL(test.haveUT, test.haveIsSecure)
		assert.NoError(t, err)
		assert.EqualValues(t, test.wantBaseUrl, baseURL.String())
		assert.EqualValues(t, test.wantPath, s.Path())

		_, err = s.BaseURL(config.URLTypeAbsent, false)
		assert.EqualError(t, err, config.ErrURLCacheCleared.Error())
	}
	t.Log(debugLogBuf.String()) // bug in config ScopedGetter
}

func TestMarshalJSON(t *testing.T) {
	s := store.MustNewStore(
		&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
	)

	jdata, err := json.Marshal(s)
	assert.NoError(t, err)
	have := []byte(`{"StoreID":1,"Code":"de","WebsiteID":1,"GroupID":1,"Name":"Germany","SortOrder":10,"IsActive":true}`)
	assert.Equal(t, have, jdata, "Have: %s\nWant: %s", have, jdata)
}
