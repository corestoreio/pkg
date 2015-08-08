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

//func init() {
//	// regarding SetConfigReader: https://twitter.com/davecheney/status/602633849374429185
//	store.SetConfigReader(config.NewMockReader(func(path string) string {
//		switch path {
//		case store.PathSecureBaseURL:
//			return store.PlaceholderBaseURL
//		case store.PathUnsecureBaseURL:
//			return store.PlaceholderBaseURL
//		case config.PathCSBaseURL:
//			return "http://cs.io/"
//		}
//		return ""
//	}, nil))
//}

func getTestManager(opts ...func(ms *mockStorage)) *store.Manager {
	ms := &mockStorage{}
	for _, opt := range opts {
		opt(ms)
	}
	return store.NewManager(store.SetManagerStorage(ms))
}

var managerStoreSimpleTest = getTestManager(func(ms *mockStorage) {
	ms.s = func() (*store.Store, error) {
		return store.NewStore(
			&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		), nil
	}
})

func TestNewManagerStore(t *testing.T) {
	assert.True(t, managerStoreSimpleTest.IsCacheEmpty())
	for j := 0; j < 3; j++ {
		s, err := managerStoreSimpleTest.Store(config.ScopeCode("notNil"))
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.EqualValues(t, "de", s.Data().Code.String)
	}
	assert.False(t, managerStoreSimpleTest.IsCacheEmpty())
	managerStoreSimpleTest.ClearCache()
	assert.True(t, managerStoreSimpleTest.IsCacheEmpty())

	tests := []struct {
		have    config.ScopeIDer
		wantErr error
	}{
		{config.ScopeCode("nilSlices"), store.ErrStoreNotFound},
		{config.ScopeID(2), store.ErrStoreNotFound},
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
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
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
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			), nil
		}
	})
	tests := []struct {
		haveManager *store.Manager
		haveID      config.ScopeIDer
		wantErr     error
	}{
		{tms, config.ScopeID(1), nil},
		{tms, config.ScopeID(1), store.ErrAppStoreSet},
		{tms, nil, store.ErrAppStoreSet},
		{tms, nil, store.ErrAppStoreSet},
	}

	for _, test := range tests {
		haveErr := test.haveManager.Init(test.haveID, config.ScopeStoreID)
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
		benchmarkManagerStore, err = managerStoreSimpleTest.Store(config.ScopeCode("de"))
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
					&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				),
				store.NewStore(
					&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				),
				store.NewStore(
					&store.TableStore{StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
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
				store.SetGroupWebsite(&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}}),
			), nil
		}
		ms.s = func() (*store.Store, error) {
			return store.NewStore(
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			), nil
		}
	})

	tests := []struct {
		m               *store.Manager
		have            config.ScopeIDer
		wantErr         error
		wantGroupName   string
		wantWebsiteCode string
	}{
		{managerGroupSimpleTest, nil, store.ErrAppStoreNotSet, "", ""},
		{getTestManager(), config.ScopeID(20), store.ErrGroupNotFound, "", ""},
		{managerGroupSimpleTest, config.ScopeID(1), nil, "DACH Group", "euro"},
		{managerGroupSimpleTest, config.ScopeID(1), nil, "DACH Group", "euro"},
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
				store.SetGroupWebsite(&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}}),
			), nil
		}
	}).Init(config.ScopeID(1), config.ScopeGroupID)
	assert.EqualError(t, store.ErrGroupDefaultStoreNotFound, err.Error(), "Incorrect DefaultStore for a Group")

	err = getTestManager().Init(config.ScopeID(21), config.ScopeGroupID)
	assert.EqualError(t, store.ErrGroupNotFound, err.Error())

	tm3 := getTestManager(func(ms *mockStorage) {
		ms.g = func() (*store.Group, error) {
			return store.NewGroup(
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				store.SetGroupWebsite(&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}}),
			).SetStores(store.TableStoreSlice{
				&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			}, nil), nil
		}
	})
	err = tm3.Init(config.ScopeID(1), config.ScopeGroupID)
	assert.NoError(t, err)
	g, err := tm3.Group()
	assert.NoError(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, int64(2), g.Data().DefaultStoreID)
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
		have            config.ScopeIDer
		wantErr         error
		wantWebsiteCode string
	}{
		{managerWebsite, nil, store.ErrAppStoreNotSet, ""},
		{getTestManager(), config.ScopeID(20), store.ErrGroupNotFound, ""},
		{managerWebsite, config.ScopeID(1), nil, "euro"},
		{managerWebsite, config.ScopeID(1), nil, "euro"},
		{managerWebsite, config.ScopeCode("notImportant"), nil, "euro"},
		{managerWebsite, config.ScopeCode("notImportant"), nil, "euro"},
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
	}).Init(config.ScopeCode("euro"), config.ScopeWebsiteID)
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

	err = managerWebsite.Init(config.ScopeCode("euro"), config.ScopeWebsiteID)
	assert.NoError(t, err)

	w2, err := managerWebsite.Website()
	assert.NoError(t, err)
	assert.EqualValues(t, "euro", w2.Data().Code.String)

	err3 := getTestManager(func(ms *mockStorage) {}).Init(config.ScopeCode("euronen"), config.ScopeWebsiteID)
	assert.Error(t, err3, "config.ScopeCode(euro), config.ScopeWebsite: %#v => %s", err3, err3)
	assert.EqualError(t, store.ErrWebsiteNotFound, err3.Error())
}

func TestNewManagerError(t *testing.T) {
	err := getTestManager().Init(config.ScopeCode("euro"), config.ScopeDefaultID)
	assert.EqualError(t, err, store.ErrUnsupportedScopeGroup.Error())
}

var storeManagerRequestStore = store.NewManager(
	store.NewStorageOption(
		store.SetStorageWebsites(
			&store.TableWebsite{WebsiteID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Admin", Valid: true}}, SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		),
		store.SetStorageGroups(
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
		),
		store.SetStorageStores(
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
	haveR         config.ScopeIDer
	wantStoreCode string
	wantErr       error
}

func runNewManagerGetRequestStore(t *testing.T, testScope config.ScopeGroup, tests []testNewManagerGetRequestStore) {
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

	testCode := config.ScopeCode("de")
	testScope := config.ScopeStoreID

	if haveStore, haveErr := storeManagerRequestStore.GetRequestStore(config.ScopeID(1), testScope); haveErr == nil {
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
		{config.ScopeID(232), "", store.ErrStoreNotFound},
		{nil, "", store.ErrStoreNotFound},
		{config.ScopeCode("\U0001f631"), "", store.ErrStoreNotFound},

		{config.ScopeID(6), "nz", nil},
		{config.ScopeCode("ch"), "", store.ErrStoreNotActive},

		{config.ScopeCode("nz"), "nz", nil},
		{config.ScopeCode("de"), "de", nil},
		{config.ScopeID(2), "at", nil},

		{config.ScopeID(2), "at", nil},
		{config.ScopeCode("au"), "au", nil},
		{config.ScopeCode("ch"), "", store.ErrStoreNotActive},
	}
	runNewManagerGetRequestStore(t, testScope, tests)
}

func TestNewManagerGetRequestStore_ScopeGroup(t *testing.T) {
	testCode := config.ScopeID(1)
	testScope := config.ScopeGroupID

	if haveStore, haveErr := storeManagerRequestStore.GetRequestStore(config.ScopeID(1), testScope); haveErr == nil {
		t.Error("appStore should not be set!")
		t.Fail()
	} else {
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
	}

	assert.EqualError(t, store.ErrGroupNotFound, storeManagerRequestStore.Init(config.ScopeID(123), testScope).Error())
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
		{config.ScopeID(232), "", store.ErrStoreNotFound},
		{nil, "", store.ErrStoreNotFound},
		{config.ScopeCode("\U0001f631"), "", store.ErrStoreNotFound},

		{config.ScopeID(6), "nz", store.ErrStoreChangeNotAllowed},
		{config.ScopeCode("ch"), "", store.ErrStoreNotActive},

		{config.ScopeCode("de"), "de", nil},
		{config.ScopeID(2), "at", nil},

		{config.ScopeID(2), "at", nil},
		{config.ScopeCode("au"), "au", store.ErrStoreChangeNotAllowed},
		{config.ScopeCode("ch"), "", store.ErrStoreNotActive},
	}
	runNewManagerGetRequestStore(t, testScope, tests)
}

func TestNewManagerGetRequestStore_ScopeWebsite(t *testing.T) {
	testCode := config.ScopeID(1)
	testScope := config.ScopeWebsiteID

	if haveStore, haveErr := storeManagerRequestStore.GetRequestStore(config.ScopeID(1), testScope); haveErr == nil {
		t.Error("appStore should not be set!")
		t.Fail()
	} else {
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
	}

	assert.EqualError(t, store.ErrUnsupportedScopeGroup, storeManagerRequestStore.Init(config.ScopeID(123), config.ScopeDefaultID).Error())
	assert.EqualError(t, store.ErrWebsiteNotFound, storeManagerRequestStore.Init(config.ScopeID(123), testScope).Error())
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
		{config.ScopeID(232), "", store.ErrStoreNotFound},
		{nil, "", store.ErrStoreNotFound},
		{config.ScopeCode("\U0001f631"), "", store.ErrStoreNotFound},

		{config.ScopeID(6), "nz", store.ErrStoreChangeNotAllowed},
		{config.ScopeCode("ch"), "", store.ErrStoreNotActive},

		{config.ScopeCode("de"), "de", nil},
		{config.ScopeID(2), "at", nil},

		{config.ScopeID(2), "at", nil},
		{config.ScopeCode("au"), "au", store.ErrStoreChangeNotAllowed},
		{config.ScopeCode("ch"), "", store.ErrStoreNotActive},
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
	tests := []struct {
		req                  *http.Request
		haveR                config.ScopeIDer
		haveScopeType        config.ScopeGroup
		wantStoreCode        string // this is the default store in a scope, lookup in storeManagerRequestStore
		wantRequestStoreCode config.ScopeCoder
		wantErr              error
		wantCookie           string
	}{
		{
			getTestRequest(t, "GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "uk"}),
			config.ScopeID(1), config.ScopeStoreID, "de", config.ScopeCode("uk"), nil, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
			config.ScopeID(1), config.ScopeStoreID, "de", config.ScopeCode("uk"), nil, store.CookieName + "=uk;", // generates a new 1y valid cookie
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=%20uk", nil),
			config.ScopeID(1), config.ScopeStoreID, "de", config.ScopeCode("uk"), store.ErrStoreNotFound, "",
		},

		{
			getTestRequest(t, "GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "de"}),
			config.ScopeID(1), config.ScopeGroupID, "at", config.ScopeCode("de"), nil, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io", nil),
			config.ScopeID(1), config.ScopeGroupID, "at", nil, nil, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=de", nil),
			config.ScopeID(1), config.ScopeGroupID, "at", config.ScopeCode("de"), nil, store.CookieName + "=de;", // generates a new 1y valid cookie
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=at", nil),
			config.ScopeID(1), config.ScopeGroupID, "at", config.ScopeCode("at"), nil, store.CookieName + "=;", // generates a delete cookie
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=cz", nil),
			config.ScopeID(1), config.ScopeGroupID, "at", nil, store.ErrStoreNotFound, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
			config.ScopeID(1), config.ScopeGroupID, "at", nil, store.ErrStoreChangeNotAllowed, "",
		},

		{
			getTestRequest(t, "GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "nz"}),
			config.ScopeID(2), config.ScopeWebsiteID, "au", config.ScopeCode("nz"), nil, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "n'z"}),
			config.ScopeID(2), config.ScopeWebsiteID, "au", nil, nil, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
			config.ScopeID(2), config.ScopeWebsiteID, "au", nil, store.ErrStoreChangeNotAllowed, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=nz", nil),
			config.ScopeID(2), config.ScopeWebsiteID, "au", config.ScopeCode("nz"), nil, store.CookieName + "=nz;",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=ch", nil),
			config.ScopeID(1), config.ScopeWebsiteID, "at", nil, store.ErrStoreNotActive, "",
		},
		{
			getTestRequest(t, "GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=nz", nil),
			config.ScopeID(1), config.ScopeDefaultID, "at", config.ScopeCode("nz"), nil, "",
		},
	}

	for _, test := range tests {
		if _, haveErr := storeManagerRequestStore.InitByRequest(nil, nil, test.haveScopeType); haveErr != nil {
			assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
		} else {
			t.Fatal("InitByRequest should return an error if used without running Init() first.")
		}

		if err := storeManagerRequestStore.Init(test.haveR, test.haveScopeType); err != nil {
			assert.EqualError(t, store.ErrUnsupportedScopeGroup, err.Error())
			t.Log("continuing for loop because of expected store.ErrUnsupportedScopeGroup")
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
				assert.EqualValues(t, test.wantRequestStoreCode.ScopeCode(), haveStore.Data().Code.String)

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
		haveR              config.ScopeIDer
		haveCodeToken      string
		haveScopeType      config.ScopeGroup
		wantStoreCode      string // this is the default store in a scope, lookup in storeManagerRequestStore
		wantTokenStoreCode config.ScopeCoder
		wantErr            error
	}{
		{config.ScopeCode("de"), "de", config.ScopeStoreID, "de", config.ScopeCode("de"), nil},
		{config.ScopeCode("de"), "at", config.ScopeStoreID, "de", config.ScopeCode("at"), nil},
		{config.ScopeCode("de"), "a$t", config.ScopeStoreID, "de", nil, nil},
		{config.ScopeCode("at"), "ch", config.ScopeStoreID, "at", nil, store.ErrStoreNotActive},
		{config.ScopeCode("at"), "", config.ScopeStoreID, "at", nil, nil},

		{config.ScopeID(1), "de", config.ScopeGroupID, "at", config.ScopeCode("de"), nil},
		{config.ScopeID(1), "ch", config.ScopeGroupID, "at", nil, store.ErrStoreNotActive},
		{config.ScopeID(1), " ch", config.ScopeGroupID, "at", nil, nil},
		{config.ScopeID(1), "uk", config.ScopeGroupID, "at", nil, store.ErrStoreChangeNotAllowed},

		{config.ScopeID(2), "uk", config.ScopeWebsiteID, "au", nil, store.ErrStoreChangeNotAllowed},
		{config.ScopeID(2), "nz", config.ScopeWebsiteID, "au", config.ScopeCode("nz"), nil},
		{config.ScopeID(2), "n z", config.ScopeWebsiteID, "au", nil, nil},
		{config.ScopeID(2), "", config.ScopeWebsiteID, "au", nil, nil},
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
				assert.Equal(t, test.wantTokenStoreCode.ScopeCode(), haveStore.Data().Code.String)
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
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	dbrSess := dbc.NewSession()

	storeManager := store.NewManager(store.NewStorageOption(nil /* trick it*/))
	if err := storeManager.ReInit(dbrSess); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		have    config.ScopeIDer
		wantErr error
	}{
		{config.ScopeCode("dede"), nil},
		{config.ScopeCode("czcz"), store.ErrStoreNotFound},
		{config.ScopeID(1), nil},
		{config.ScopeID(100), store.ErrStoreNotFound},
		{mockIDCode{1, "dede"}, nil},
		{mockIDCode{2, "czfr"}, store.ErrStoreNotFound},
		{mockIDCode{2, ""}, nil},
		{nil, store.ErrAppStoreNotSet}, // if set returns default store
	}

	for _, test := range tests {
		s, err := storeManager.Store(test.have)
		if test.wantErr == nil {
			assert.NoError(t, err, "No Err; for test: %#v", test)
			assert.NotNil(t, s)
			//			assert.NotEmpty(t, s.Data().Code.String, "%#v", s.Data())
		} else {
			assert.Error(t, err, "Err for test: %#v", test)
			assert.EqualError(t, test.wantErr, err.Error(), "EqualErr for test: %#v", test)
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

type mockIDCode struct {
	id   int64
	code string
}

func (ic mockIDCode) ScopeID() int64 {
	return ic.id
}
func (ic mockIDCode) ScopeCode() string {
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

func (ms *mockStorage) Website(_ config.ScopeIDer) (*store.Website, error) {
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
func (ms *mockStorage) Group(_ config.ScopeIDer) (*store.Group, error) {
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
func (ms *mockStorage) Store(_ config.ScopeIDer) (*store.Store, error) {
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
