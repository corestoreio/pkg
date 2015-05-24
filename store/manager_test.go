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
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/dgrijalva/jwt-go"
	"github.com/juju/errgo"
	"github.com/stretchr/testify/assert"
)

func init() {
	store.SetConfigReader(newMockScopeReader(nil, nil))
}

func getTestManager(opts ...func(ms *mockStorage)) *store.Manager {
	ms := &mockStorage{}
	for _, opt := range opts {
		opt(ms)
	}
	return store.NewManager(ms)
}

var managerStoreSimpleTest = getTestManager(func(ms *mockStorage) {
	ms.s = func() (*store.Store, error) {
		return store.NewStore(
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		), nil
	}
})

func TestNewManagerStore(t *testing.T) {
	assert.True(t, managerStoreSimpleTest.IsCacheEmpty())
	for j := 0; j < 3; j++ {
		s, err := managerStoreSimpleTest.Store(store.Code("notNil"))
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.EqualValues(t, "de", s.Data().Code.String)
	}
	assert.False(t, managerStoreSimpleTest.IsCacheEmpty())
	managerStoreSimpleTest.ClearCache()
	assert.True(t, managerStoreSimpleTest.IsCacheEmpty())

	tests := []struct {
		have    store.Retriever
		wantErr error
	}{
		{store.Code("nilSlices"), store.ErrStoreNotFound},
		{store.ID(2), store.ErrStoreNotFound},
		{nil, store.ErrAppStoreNotSet},
	}

	managerEmpty := getTestManager()
	for _, test := range tests {
		s, err := managerEmpty.Store(test.have)
		assert.Nil(t, s)
		assert.EqualError(t, test.wantErr, err.Error())
	}
	assert.True(t, managerStoreSimpleTest.IsCacheEmpty())
}

func TestNewManagerDefaultStoreView(t *testing.T) {
	managerDefaultStore := getTestManager(func(ms *mockStorage) {
		ms.dsv = func() (*store.Store, error) {
			return store.NewStore(
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			), nil
		}
	})

	// call it twice to test internal caching
	s, err := managerDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.NoError(t, err)
	assert.NotEmpty(t, s.Data().Code.String)

	s, err = managerDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.NoError(t, err)
	assert.NotEmpty(t, s.Data().Code.String)
	assert.False(t, managerDefaultStore.IsCacheEmpty())
	managerDefaultStore.ClearCache()
	assert.True(t, managerDefaultStore.IsCacheEmpty())
}

func TestNewManagerStoreInit(t *testing.T) {

	tms := getTestManager(func(ms *mockStorage) {
		ms.s = func() (*store.Store, error) {
			return store.NewStore(
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			), nil
		}
	})
	tests := []struct {
		haveManager *store.Manager
		haveID      store.Retriever
		wantErr     error
	}{
		{tms, store.ID(1), nil},
		{tms, store.ID(1), store.ErrAppStoreSet},
		{tms, nil, store.ErrAppStoreSet},
		{tms, nil, store.ErrAppStoreSet},
	}

	for _, test := range tests {
		haveErr := test.haveManager.Init(test.haveID, config.IDScopeStore)
		if test.wantErr != nil {
			assert.Error(t, haveErr)
			assert.EqualError(t, test.wantErr, haveErr.Error())
		} else {
			assert.NoError(t, haveErr)
		}
		s, err := test.haveManager.Store()
		assert.NotNil(t, s)
		assert.NoError(t, err)
	}
}

var benchmarkManagerStore *store.Store

// BenchmarkManagerGetStore	 5000000	       355 ns/op	      24 B/op	       2 allocs/op
func BenchmarkManagerGetStore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkManagerStore, err = managerStoreSimpleTest.Store(store.Code("de"))
		if err != nil {
			b.Error(err)
		}
		if benchmarkManagerStore == nil {
			b.Error("benchmarkManagerStore is nil")
		}
	}
}

func TestNewManagerStores(t *testing.T) {
	managerStores := getTestManager(func(ms *mockStorage) {
		ms.ss = func() (store.StoreSlice, error) {
			return store.StoreSlice{
				store.NewStore(
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
					&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				),
				store.NewStore(
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
					&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
				),
				store.NewStore(
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
					&store.TableStore{StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
				),
			}, nil
		}
	})

	// call it twice to test internal caching
	ss, err := managerStores.Stores()
	assert.NotNil(t, ss)
	assert.NoError(t, err)
	assert.Equal(t, "at", ss[1].Data().Code.String)

	ss, err = managerStores.Stores()
	assert.NotNil(t, ss)
	assert.NoError(t, err)
	assert.NotEmpty(t, ss[2].Data().Code.String)

	assert.False(t, managerStores.IsCacheEmpty())
	managerStores.ClearCache()
	assert.True(t, managerStores.IsCacheEmpty())

	ss, err = getTestManager(func(ms *mockStorage) {
		ms.ss = func() (store.StoreSlice, error) {
			return nil, nil
		}
	}).Stores()
	assert.Nil(t, ss)
	assert.NoError(t, err)
}

func TestNewManagerGroup(t *testing.T) {
	var managerGroupSimpleTest = getTestManager(func(ms *mockStorage) {
		ms.g = func() (*store.Group, error) {
			return store.NewGroup(
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			), nil
		}
		ms.s = func() (*store.Store, error) {
			return store.NewStore(
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			), nil
		}
	})

	tests := []struct {
		m               *store.Manager
		have            store.Retriever
		wantErr         error
		wantGroupName   string
		wantWebsiteCode string
	}{
		{managerGroupSimpleTest, nil, store.ErrAppStoreNotSet, "", ""},
		{getTestManager(), store.ID(20), store.ErrGroupNotFound, "", ""},
		{managerGroupSimpleTest, store.ID(1), nil, "DACH Group", "euro"},
		{managerGroupSimpleTest, store.ID(1), nil, "DACH Group", "euro"},
	}

	for _, test := range tests {
		g, err := test.m.Group(test.have)
		if test.wantErr != nil {
			assert.Nil(t, g)
			assert.EqualError(t, test.wantErr, err.Error(), "test %#v", test)
		} else {
			assert.NotNil(t, g, "test %#v", test)
			assert.NoError(t, err, "test %#v", test)
			assert.Equal(t, test.wantGroupName, g.Data().Name)
			assert.Equal(t, test.wantWebsiteCode, g.Website().Data().Code.String)
		}
	}
	assert.False(t, managerGroupSimpleTest.IsCacheEmpty())
	managerGroupSimpleTest.ClearCache()
	assert.True(t, managerGroupSimpleTest.IsCacheEmpty())
}

func TestNewManagerGroupInit(t *testing.T) {

	err := getTestManager(func(ms *mockStorage) {
		ms.g = func() (*store.Group, error) {
			return store.NewGroup(
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			), nil
		}
	}).Init(store.ID(1), config.IDScopeGroup)
	assert.EqualError(t, store.ErrGroupDefaultStoreNotFound, err.Error(), "Incorrect DefaultStore for a Group")

	err = getTestManager().Init(store.ID(21), config.IDScopeGroup)
	assert.EqualError(t, store.ErrGroupNotFound, err.Error())

	tm3 := getTestManager(func(ms *mockStorage) {
		ms.g = func() (*store.Group, error) {
			return store.NewGroup(
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			).SetStores(store.TableStoreSlice{
				&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			}, nil), nil
		}
	})
	err = tm3.Init(store.ID(1), config.IDScopeGroup)
	assert.NoError(t, err)
	g, err := tm3.Group()
	assert.NoError(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, 2, g.Data().DefaultStoreID)
}

func TestNewManagerGroups(t *testing.T) {
	managerGroups := getTestManager(func(ms *mockStorage) {
		ms.gs = func() (store.GroupSlice, error) {
			return store.GroupSlice{}, nil
		}
	})

	// call it twice to test internal caching
	ss, err := managerGroups.Groups()
	assert.NotNil(t, ss)
	assert.NoError(t, err)
	assert.Len(t, ss, 0)

	ss, err = managerGroups.Groups()
	assert.NotNil(t, ss)
	assert.NoError(t, err)
	assert.Len(t, ss, 0)

	assert.False(t, managerGroups.IsCacheEmpty())
	managerGroups.ClearCache()
	assert.True(t, managerGroups.IsCacheEmpty())
}

func TestNewManagerWebsite(t *testing.T) {

	var managerWebsite = getTestManager(func(ms *mockStorage) {
		ms.w = func() (*store.Website, error) {
			return store.NewWebsite(
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			), nil
		}
	})

	tests := []struct {
		m               *store.Manager
		have            store.Retriever
		wantErr         error
		wantWebsiteCode string
	}{
		{managerWebsite, nil, store.ErrAppStoreNotSet, ""},
		{getTestManager(), store.ID(20), store.ErrGroupNotFound, ""},
		{managerWebsite, store.ID(1), nil, "euro"},
		{managerWebsite, store.ID(1), nil, "euro"},
		{managerWebsite, store.Code("notImportant"), nil, "euro"},
		{managerWebsite, store.Code("notImportant"), nil, "euro"},
	}

	for _, test := range tests {
		haveW, haveErr := test.m.Website(test.have)
		if test.wantErr != nil {
			assert.Error(t, haveErr, "%#v", test)
			assert.Nil(t, haveW, "%#v", test)
		} else {
			assert.NoError(t, haveErr, "%#v", test)
			assert.NotNil(t, haveW, "%#v", test)
			assert.Equal(t, test.wantWebsiteCode, haveW.Data().Code.String)
		}
	}
	assert.False(t, managerWebsite.IsCacheEmpty())
	managerWebsite.ClearCache()
	assert.True(t, managerWebsite.IsCacheEmpty())

}

func TestNewManagerWebsites(t *testing.T) {
	managerWebsites := getTestManager(func(ms *mockStorage) {
		ms.ws = func() (store.WebsiteSlice, error) {
			return store.WebsiteSlice{}, nil
		}
	})

	tests := []struct {
		m       *store.Manager
		wantErr error
		wantNil bool
	}{
		{managerWebsites, nil, false},
		{managerWebsites, nil, false},
		{getTestManager(func(ms *mockStorage) {
			ms.ws = func() (store.WebsiteSlice, error) {
				return nil, nil
			}
		}), nil, true},
	}

	for _, test := range tests {
		haveWS, haveErr := test.m.Websites()
		if test.wantErr != nil {
			assert.Error(t, haveErr, "%#v", test)
			assert.Nil(t, haveWS, "%#v", test)
		} else {
			assert.NoError(t, haveErr, "%#v", test)
			if test.wantNil {
				assert.Nil(t, haveWS, "%#v", test)
			} else {
				assert.NotNil(t, haveWS, "%#v", test)
			}
		}
	}

	assert.False(t, managerWebsites.IsCacheEmpty())
	managerWebsites.ClearCache()
	assert.True(t, managerWebsites.IsCacheEmpty())
}

func TestNewManagerWebsiteInit(t *testing.T) {

	err := getTestManager(func(ms *mockStorage) {
		ms.w = func() (*store.Website, error) {
			return store.NewWebsite(
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			), nil
		}
	}).Init(store.Code("euro"), config.IDScopeWebsite)
	assert.EqualError(t, store.ErrWebsiteDefaultGroupNotFound, err.Error())

	managerWebsite := getTestManager(func(ms *mockStorage) {
		ms.w = func() (*store.Website, error) {
			return store.NewWebsite(
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			).SetGroupsStores(
				store.TableGroupSlice{
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				},
				store.TableStoreSlice{
					&store.TableStore{StoreID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
					&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
					&store.TableStore{StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
				},
			), nil
		}
	})
	w1, err := managerWebsite.Website()
	assert.EqualError(t, store.ErrAppStoreNotSet, err.Error())
	assert.Nil(t, w1)

	err = managerWebsite.Init(store.Code("euro"), config.IDScopeWebsite)
	assert.NoError(t, err)

	w2, err := managerWebsite.Website()
	assert.NoError(t, err)
	assert.EqualValues(t, "euro", w2.Data().Code.String)

	err3 := getTestManager(func(ms *mockStorage) {}).Init(store.Code("euronen"), config.IDScopeWebsite)
	assert.Error(t, err3, "store.Code(euro), config.ScopeWebsite: %#v => %s", err3, err3)
	assert.EqualError(t, store.ErrWebsiteNotFound, err3.Error())
}

func TestNewManagerError(t *testing.T) {
	err := getTestManager().Init(store.Code("euro"), config.IDScopeDefault)
	assert.EqualError(t, err, store.ErrUnsupportedScopeID.Error())
}

var storeManagerRequestStore = store.NewManager(
	store.NewStorage(
		store.StorageTableWebsites(
			&store.TableWebsite{WebsiteID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Admin", Valid: true}}, SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		),
		store.StorageTableGroups(
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
		),
		store.StorageTableStores(
			&store.TableStore{StoreID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
			&store.TableStore{StoreID: 5, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "uk", Valid: true}}, WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NullString{NullString: sql.NullString{String: "nz", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
			&store.TableStore{IsActive: false, StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30},
		),
	),
)

type testNewManagerGetRequestStore struct {
	haveR         store.Retriever
	wantStoreCode string
	wantErr       error
}

func runNewManagerGetRequestStore(t *testing.T, testScope config.ScopeID, tests []testNewManagerGetRequestStore) {
	for _, test := range tests {
		haveStore, haveErr := storeManagerRequestStore.GetRequestStore(test.haveR, testScope)
		if test.wantErr != nil {
			assert.Nil(t, haveStore, "testScope %d: %#v", testScope, test)
			assert.EqualError(t, test.wantErr, haveErr.Error(), "testScope %d: %#v", testScope, test)
		} else {
			assert.NotNil(t, haveStore)
			assert.NoError(t, haveErr, "%#v", test)
			assert.EqualValues(t, test.wantStoreCode, haveStore.Data().Code.String)
		}
	}
	storeManagerRequestStore.ClearCache(true)
}
func TestNewManagerGetRequestStore_ScopeStore(t *testing.T) {

	testCode := store.Code("de")
	testScope := config.IDScopeStore

	if haveStore, haveErr := storeManagerRequestStore.GetRequestStore(store.ID(1), testScope); haveErr == nil {
		t.Error("appStore should not be set!")
		t.Fail()
	} else {
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
	}

	// init with scope store
	if err := storeManagerRequestStore.Init(testCode, testScope); err != nil {
		t.Error(err)
		t.Fail()
	}
	assert.EqualError(t, store.ErrAppStoreSet, storeManagerRequestStore.Init(testCode, testScope).Error())

	if s, err := storeManagerRequestStore.Store(); err == nil {
		assert.EqualValues(t, "de", s.Data().Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	tests := []testNewManagerGetRequestStore{
		{store.ID(232), "", store.ErrStoreNotFound},
		{nil, "", store.ErrStoreNotFound},
		{store.Code("\U0001f631"), "", store.ErrStoreNotFound},

		{store.ID(6), "nz", nil},
		{store.Code("ch"), "", store.ErrStoreNotActive},

		{store.Code("nz"), "nz", nil},
		{store.Code("de"), "de", nil},
		{store.ID(2), "at", nil},

		{store.ID(2), "at", nil},
		{store.Code("au"), "au", nil},
		{store.Code("ch"), "", store.ErrStoreNotActive},
	}
	runNewManagerGetRequestStore(t, testScope, tests)
}

func TestNewManagerGetRequestStore_ScopeGroup(t *testing.T) {
	testCode := store.ID(1)
	testScope := config.IDScopeGroup

	if haveStore, haveErr := storeManagerRequestStore.GetRequestStore(store.ID(1), testScope); haveErr == nil {
		t.Error("appStore should not be set!")
		t.Fail()
	} else {
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
	}

	assert.EqualError(t, store.ErrGroupNotFound, storeManagerRequestStore.Init(store.ID(123), testScope).Error())
	if err := storeManagerRequestStore.Init(testCode, testScope); err != nil {
		t.Error(err)
		t.Fail()
	}
	assert.EqualError(t, store.ErrAppStoreSet, storeManagerRequestStore.Init(testCode, testScope).Error())

	if s, err := storeManagerRequestStore.Store(); err == nil {
		assert.EqualValues(t, "at", s.Data().Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	if g, err := storeManagerRequestStore.Group(); err == nil {
		assert.EqualValues(t, 1, g.Data().GroupID)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	tests := []testNewManagerGetRequestStore{
		{store.ID(232), "", store.ErrStoreNotFound},
		{nil, "", store.ErrStoreNotFound},
		{store.Code("\U0001f631"), "", store.ErrStoreNotFound},

		{store.ID(6), "nz", store.ErrStoreChangeNotAllowed},
		{store.Code("ch"), "", store.ErrStoreNotActive},

		{store.Code("de"), "de", nil},
		{store.ID(2), "at", nil},

		{store.ID(2), "at", nil},
		{store.Code("au"), "au", store.ErrStoreChangeNotAllowed},
		{store.Code("ch"), "", store.ErrStoreNotActive},
	}
	runNewManagerGetRequestStore(t, testScope, tests)
}

func TestNewManagerGetRequestStore_ScopeWebsite(t *testing.T) {
	testCode := store.ID(1)
	testScope := config.IDScopeWebsite

	if haveStore, haveErr := storeManagerRequestStore.GetRequestStore(store.ID(1), testScope); haveErr == nil {
		t.Error("appStore should not be set!")
		t.Fail()
	} else {
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
	}

	assert.EqualError(t, store.ErrUnsupportedScopeID, storeManagerRequestStore.Init(store.ID(123), config.IDScopeDefault).Error())
	assert.EqualError(t, store.ErrWebsiteNotFound, storeManagerRequestStore.Init(store.ID(123), testScope).Error())
	if err := storeManagerRequestStore.Init(testCode, testScope); err != nil {
		t.Error(err)
		t.Fail()
	}
	assert.EqualError(t, store.ErrAppStoreSet, storeManagerRequestStore.Init(testCode, testScope).Error())

	if s, err := storeManagerRequestStore.Store(); err == nil {
		assert.EqualValues(t, "at", s.Data().Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	if w, err := storeManagerRequestStore.Website(); err == nil {
		assert.EqualValues(t, "euro", w.Data().Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	tests := []testNewManagerGetRequestStore{
		{store.ID(232), "", store.ErrStoreNotFound},
		{nil, "", store.ErrStoreNotFound},
		{store.Code("\U0001f631"), "", store.ErrStoreNotFound},

		{store.ID(6), "nz", store.ErrStoreChangeNotAllowed},
		{store.Code("ch"), "", store.ErrStoreNotActive},

		{store.Code("de"), "de", nil},
		{store.ID(2), "at", nil},

		{store.ID(2), "at", nil},
		{store.Code("au"), "au", store.ErrStoreChangeNotAllowed},
		{store.Code("ch"), "", store.ErrStoreNotActive},
	}
	runNewManagerGetRequestStore(t, testScope, tests)
}

func getTestRequest(t *testing.T, m, u string, c *http.Cookie) *http.Request {
	req, err := http.NewRequest(m, u, nil)
	if err != nil {
		t.Fatal(err)
	}
	if c != nil {
		req.AddCookie(c)
	}
	return req
}

// cyclomatic complexity 12 of function TestInitByRequest() is high (> 10) (gocyclo)
func TestInitByRequest(t *testing.T) {
	store.SetConfigReader(newMockScopeReader(func(path string, scope config.ScopeID, r ...config.Retriever) string {
		switch path {
		case store.PathSecureBaseURL:
			return store.PlaceholderBaseURL
		case store.PathUnsecureBaseURL:
			return store.PlaceholderBaseURL
		case config.PathCSBaseURL:
			return "http://cs.io/"
		}
		return ""
	}, nil))

	tests := []struct {
		req                  *http.Request
		haveR                store.Retriever
		haveScopeType        config.ScopeID
		wantStoreCode        string // this is the default store in a scope, lookup in storeManagerRequestStore
		wantRequestStoreCode store.CodeRetriever
		wantErr              error
		wantCookie           string
	}{
		{
			getTestRequest(t, "GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "uk"}),
			store.ID(1), config.IDScopeStore, "de", store.Code("uk"), nil, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
			store.ID(1), config.IDScopeStore, "de", store.Code("uk"), nil, store.CookieName + "=uk;", // generates a new 1y valid cookie
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=%20uk", nil),
			store.ID(1), config.IDScopeStore, "de", store.Code("uk"), store.ErrStoreNotFound, "",
		},

		{
			getTestRequest(t, "GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "de"}),
			store.ID(1), config.IDScopeGroup, "at", store.Code("de"), nil, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io", nil),
			store.ID(1), config.IDScopeGroup, "at", nil, nil, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=de", nil),
			store.ID(1), config.IDScopeGroup, "at", store.Code("de"), nil, store.CookieName + "=de;", // generates a new 1y valid cookie
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=at", nil),
			store.ID(1), config.IDScopeGroup, "at", store.Code("at"), nil, store.CookieName + "=;", // generates a delete cookie
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=cz", nil),
			store.ID(1), config.IDScopeGroup, "at", nil, store.ErrStoreNotFound, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
			store.ID(1), config.IDScopeGroup, "at", nil, store.ErrStoreChangeNotAllowed, "",
		},

		{
			getTestRequest(t, "GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "nz"}),
			store.ID(2), config.IDScopeWebsite, "au", store.Code("nz"), nil, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "n'z"}),
			store.ID(2), config.IDScopeWebsite, "au", nil, nil, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
			store.ID(2), config.IDScopeWebsite, "au", nil, store.ErrStoreChangeNotAllowed, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=nz", nil),
			store.ID(2), config.IDScopeWebsite, "au", store.Code("nz"), nil, store.CookieName + "=nz;",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=ch", nil),
			store.ID(1), config.IDScopeWebsite, "at", nil, store.ErrStoreNotActive, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=nz", nil),
			store.ID(1), config.IDScopeDefault, "at", store.Code("nz"), nil, "",
		},
	}

	for _, test := range tests {
		if _, haveErr := storeManagerRequestStore.InitByRequest(nil, nil, test.haveScopeType); haveErr != nil {
			assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
		} else {
			t.Fatal("InitByRequest should return an error if used without running Init() first.")
		}

		if err := storeManagerRequestStore.Init(test.haveR, test.haveScopeType); err != nil {
			assert.EqualError(t, store.ErrUnsupportedScopeID, err.Error())
			t.Log("continuing for loop because of expected store.ErrUnsupportedScopeID")
			storeManagerRequestStore.ClearCache(true)
			continue
		}

		if s, err := storeManagerRequestStore.Store(); err == nil {
			assert.EqualValues(t, test.wantStoreCode, s.Data().Code.String)
		} else {
			assert.EqualError(t, err, store.ErrStoreNotFound.Error())
			t.Log("continuing for loop because of expected store.ErrStoreNotFound")
			storeManagerRequestStore.ClearCache(true)
			continue
		}

		resRec := httptest.NewRecorder()

		haveStore, haveErr := storeManagerRequestStore.InitByRequest(resRec, test.req, test.haveScopeType)
		if test.wantErr != nil {
			assert.Nil(t, haveStore)
			assert.EqualError(t, test.wantErr, haveErr.Error())
		} else {
			if msg, ok := haveErr.(errgo.Locationer); ok {
				t.Logf("\nLocation: %s => %s\n", haveErr, msg.Location())
			}
			assert.NoError(t, haveErr, "%#v", test)
			if test.wantRequestStoreCode != nil {
				assert.NotNil(t, haveStore, "%#v", test.req.URL.Query())
				assert.EqualValues(t, test.wantRequestStoreCode.Code(), haveStore.Data().Code.String)

				newKeks := resRec.HeaderMap.Get("Set-Cookie")
				if test.wantCookie != "" {
					assert.Contains(t, newKeks, test.wantCookie, "%#v", test)
					//					t.Logf(
					//						"\nwantRequestStoreCode: %s\nCookie Str: %#v\n",
					//						test.wantRequestStoreCode.Code(),
					//						newKeks,
					//					)
				} else {
					assert.Empty(t, newKeks, "%#v", test)
				}

			} else {
				assert.Nil(t, haveStore, "%#v", haveStore)
			}
		}
		storeManagerRequestStore.ClearCache(true)
	}
}

func TestInitByToken(t *testing.T) {

	getToken := func(code string) *jwt.Token {
		t := jwt.New(jwt.SigningMethodHS256)
		t.Claims[store.CookieName] = code
		return t
	}

	tests := []struct {
		haveR              store.Retriever
		haveCodeToken      string
		haveScopeType      config.ScopeID
		wantStoreCode      string // this is the default store in a scope, lookup in storeManagerRequestStore
		wantTokenStoreCode store.CodeRetriever
		wantErr            error
	}{
		{store.Code("de"), "de", config.IDScopeStore, "de", store.Code("de"), nil},
		{store.Code("de"), "at", config.IDScopeStore, "de", store.Code("at"), nil},
		{store.Code("de"), "a$t", config.IDScopeStore, "de", nil, nil},
		{store.Code("at"), "ch", config.IDScopeStore, "at", nil, store.ErrStoreNotActive},
		{store.Code("at"), "", config.IDScopeStore, "at", nil, nil},

		{store.ID(1), "de", config.IDScopeGroup, "at", store.Code("de"), nil},
		{store.ID(1), "ch", config.IDScopeGroup, "at", nil, store.ErrStoreNotActive},
		{store.ID(1), " ch", config.IDScopeGroup, "at", nil, nil},
		{store.ID(1), "uk", config.IDScopeGroup, "at", nil, store.ErrStoreChangeNotAllowed},

		{store.ID(2), "uk", config.IDScopeWebsite, "au", nil, store.ErrStoreChangeNotAllowed},
		{store.ID(2), "nz", config.IDScopeWebsite, "au", store.Code("nz"), nil},
		{store.ID(2), "n z", config.IDScopeWebsite, "au", nil, nil},
		{store.ID(2), "", config.IDScopeWebsite, "au", nil, nil},
	}
	for _, test := range tests {

		haveStore, haveErr := storeManagerRequestStore.InitByToken(nil, test.haveScopeType)
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())

		if err := storeManagerRequestStore.Init(test.haveR, test.haveScopeType); err != nil {
			t.Fatal(err)
		}

		if s, err := storeManagerRequestStore.Store(); err == nil {
			assert.EqualValues(t, test.wantStoreCode, s.Data().Code.String)
		} else {
			assert.EqualError(t, err, store.ErrStoreNotFound.Error())
			t.Fail()
		}

		haveStore, haveErr = storeManagerRequestStore.InitByToken(getToken(test.haveCodeToken), test.haveScopeType)
		if test.wantErr != nil {
			assert.Nil(t, haveStore, "%#v", test)
			assert.Error(t, haveErr, "%#v", test)
			assert.EqualError(t, test.wantErr, haveErr.Error())
		} else {
			if test.wantTokenStoreCode != nil {
				assert.NotNil(t, haveStore, "%#v", test)
				assert.NoError(t, haveErr)
				assert.Equal(t, test.wantTokenStoreCode.Code(), haveStore.Data().Code.String)
			} else {
				assert.Nil(t, haveStore, "%#v", test)
				assert.NoError(t, haveErr, "%#v", test)
			}

		}
		storeManagerRequestStore.ClearCache(true)
	}
}

func TestNewManagerReInit(t *testing.T) {
	numCPU := runtime.NumCPU()
	prevCPU := runtime.GOMAXPROCS(numCPU)
	t.Logf("GOMAXPROCS was: %d now: %d", prevCPU, numCPU)
	defer runtime.GOMAXPROCS(prevCPU)

	// quick implement, use mock of dbr.SessionRunner and remove connection
	db := csdb.MustConnectTest()
	defer db.Close()
	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)

	storeManager := store.NewManager(store.NewStorage(nil /* trick it*/))
	if err := storeManager.ReInit(dbrSess); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		have    store.Retriever
		wantErr error
	}{
		{store.Code("de"), nil},
		{store.Code("cz"), store.ErrStoreNotFound},
		{store.Code("de"), nil},
		{store.ID(1), nil},
		{store.ID(100), store.ErrStoreNotFound},
		{mockIDCode{1, "de"}, nil},
		{mockIDCode{2, "cz"}, store.ErrStoreNotFound},
		{mockIDCode{2, ""}, nil},
		{nil, store.ErrAppStoreNotSet}, // if set returns default store
	}

	for _, test := range tests {
		s, err := storeManager.Store(test.have)
		if test.wantErr == nil {
			assert.NoError(t, err, "For test: %#v", test)
			assert.NotNil(t, s)
			//			assert.NotEmpty(t, s.Data().Code.String, "%#v", s.Data())
		} else {
			assert.Error(t, err, "For test: %#v", test)
			assert.EqualError(t, test.wantErr, err.Error(), "For test: %#v", test)
			assert.Nil(t, s)
		}
	}
	assert.False(t, storeManager.IsCacheEmpty())
	storeManager.ClearCache()
	assert.True(t, storeManager.IsCacheEmpty())
}

/*
	MOCKS
*/

var _ config.Reader = (*mockScopeReader)(nil)

type mockScopeReader struct {
	s func(...config.OptionFunc) string
	b func(...config.OptionFunc) bool
}

func newMockScopeReader(
	s func(...config.OptionFunc) string,
	b func(...config.OptionFunc) bool,
) *mockScopeReader {

	return &mockScopeReader{
		s: s,
		b: b,
	}
}

func (sr mockScopeReader) GetString(opts ...config.OptionFunc) string {
	if sr.s == nil {
		return ""
	}
	return sr.s(opts...)
}

func (sr mockScopeReader) GetBool(opts ...config.OptionFunc) bool {
	if sr.b == nil {
		return false
	}
	return sr.b(opts...)
}

type mockIDCode struct {
	id   int64
	code string
}

func (ic mockIDCode) ID() int64 {
	return ic.id
}
func (ic mockIDCode) Code() string {
	return ic.code
}

type mockStorage struct {
	w   func() (*store.Website, error)
	ws  func() (store.WebsiteSlice, error)
	g   func() (*store.Group, error)
	gs  func() (store.GroupSlice, error)
	s   func() (*store.Store, error)
	dsv func() (*store.Store, error)
	ss  func() (store.StoreSlice, error)
}

var _ store.Storager = (*mockStorage)(nil)

func (ms *mockStorage) Website(_ store.Retriever) (*store.Website, error) {
	if ms.w == nil {
		return nil, store.ErrWebsiteNotFound
	}
	return ms.w()
}
func (ms *mockStorage) Websites() (store.WebsiteSlice, error) {
	if ms.ws == nil {
		return nil, nil
	}
	return ms.ws()
}
func (ms *mockStorage) Group(_ store.Retriever) (*store.Group, error) {
	if ms.g == nil {
		return nil, store.ErrGroupNotFound
	}
	return ms.g()
}
func (ms *mockStorage) Groups() (store.GroupSlice, error) {
	if ms.gs == nil {
		return nil, nil
	}
	return ms.gs()
}
func (ms *mockStorage) Store(_ store.Retriever) (*store.Store, error) {
	if ms.s == nil {
		return nil, store.ErrStoreNotFound
	}
	return ms.s()
}

func (ms *mockStorage) Stores() (store.StoreSlice, error) {
	if ms.ss == nil {
		return nil, nil
	}
	return ms.ss()
}
func (ms *mockStorage) DefaultStoreView() (*store.Store, error) {
	if ms.dsv == nil {
		return nil, store.ErrStoreNotFound
	}
	return ms.dsv()
}
func (ms *mockStorage) ReInit(dbr.SessionRunner, ...csdb.DbrSelectCb) error {
	return nil
}
