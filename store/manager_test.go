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
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/stretchr/testify/assert"
)

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
		haveErr := test.haveManager.Init(test.haveID, config.ScopeStore)
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
	}).Init(store.ID(1), config.ScopeGroup)
	assert.EqualError(t, store.ErrGroupDefaultStoreNotFound, err.Error(), "Incorrect DefaultStore for a Group")

	err = getTestManager().Init(store.ID(21), config.ScopeGroup)
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
	err = tm3.Init(store.ID(1), config.ScopeGroup)
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
	}).Init(store.Code("euro"), config.ScopeWebsite)
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

	err = managerWebsite.Init(store.Code("euro"), config.ScopeWebsite)
	assert.NoError(t, err)

	w2, err := managerWebsite.Website()
	assert.NoError(t, err)
	assert.EqualValues(t, "euro", w2.Data().Code.String)

	err3 := getTestManager(func(ms *mockStorage) {}).Init(store.Code("euro"), config.ScopeWebsite)
	assert.EqualError(t, store.ErrWebsiteNotFound, err3.Error())
}

func TestNewManagerError(t *testing.T) {
	err := getTestManager().Init(store.Code("euro"), config.ScopeDefault)
	assert.EqualError(t, err, store.ErrUnsupportedScopeID.Error())
}

var storeManagerRequestStore = store.NewManager(
	store.NewStorage(
		store.TableWebsiteSlice{
			&store.TableWebsite{WebsiteID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Admin", Valid: true}}, SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		},
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
			&store.TableStore{IsActive: false, StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30},
		},
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
			assert.Nil(t, haveStore, "%#v", test)
			assert.EqualError(t, test.wantErr, haveErr.Error(), "%#v", test)
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
	testScope := config.ScopeStore

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
		{store.Code("ch"), "", store.ErrStoreNotFound},

		{store.Code("nz"), "nz", nil},
		{store.Code("de"), "de", nil},
		{store.ID(2), "at", nil},

		{store.ID(2), "at", nil},
		{store.Code("au"), "au", nil},
		{store.Code("ch"), "", store.ErrStoreNotFound},
	}
	runNewManagerGetRequestStore(t, testScope, tests)
}

func TestNewManagerGetRequestStore_ScopeGroup(t *testing.T) {
	testCode := store.ID(1)
	testScope := config.ScopeGroup

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
		{store.Code("ch"), "", store.ErrStoreNotFound},

		{store.Code("de"), "de", nil},
		{store.ID(2), "at", nil},

		{store.ID(2), "at", nil},
		{store.Code("au"), "au", store.ErrStoreChangeNotAllowed},
		{store.Code("ch"), "", store.ErrStoreNotFound},
	}
	runNewManagerGetRequestStore(t, testScope, tests)
}

func TestNewManagerGetRequestStore_ScopeWebsite(t *testing.T) {
	testCode := store.ID(1)
	testScope := config.ScopeWebsite

	if haveStore, haveErr := storeManagerRequestStore.GetRequestStore(store.ID(1), testScope); haveErr == nil {
		t.Error("appStore should not be set!")
		t.Fail()
	} else {
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
	}

	assert.EqualError(t, store.ErrUnsupportedScopeID, storeManagerRequestStore.Init(store.ID(123), config.ScopeDefault).Error())
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
		{store.Code("ch"), "", store.ErrStoreNotFound},

		{store.Code("de"), "de", nil},
		{store.ID(2), "at", nil},

		{store.ID(2), "at", nil},
		{store.Code("au"), "au", store.ErrStoreChangeNotAllowed},
		{store.Code("ch"), "", store.ErrStoreNotFound},
	}
	runNewManagerGetRequestStore(t, testScope, tests)
}

func TestInitByRequest_Group(t *testing.T) {
	testCode := store.ID(1)
	testScope := config.ScopeGroup
	if err := storeManagerRequestStore.Init(testCode, testScope); err != nil {
		t.Error(err)
		t.Fail()
	}

	tests := []struct {
		res           http.ResponseWriter
		req           *http.Request
		wantStoreCode string
		wantErr       error
	}{}

	for _, test := range tests {
		haveStore, haveErr := storeManagerRequestStore.InitByRequest(test.res, test.req, testScope)
	}
}

/*
	MOCKS
*/

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
