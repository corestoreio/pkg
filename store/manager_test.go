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
	"bytes"
	"database/sql"
	std "log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// Within a test function use defer errLogBuf.Reset() to clean the logger
var errLogBuf bytes.Buffer

func init() {
	log.Set(log.NewStdLogger(
		log.SetStdError(&errLogBuf, "testErr: ", std.Lshortfile),
	))
	log.SetLevel(log.StdLevelError)
}

//func init() {
// TODO(cs): check this
//	// regarding SetConfigReader: https://twitter.com/davecheney/status/602633849374429185
// 		@ianthomasrose @francesc package variables are a smell, modifying them for tests is a stink.
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
	return store.NewManager(ms)
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
		s, err := managerStoreSimpleTest.Store(scope.MockCode("notNil"))
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.EqualValues(t, "de", s.Data.Code.String)
	}
	assert.False(t, managerStoreSimpleTest.IsCacheEmpty())
	managerStoreSimpleTest.ClearCache()
	assert.True(t, managerStoreSimpleTest.IsCacheEmpty())

	tests := []struct {
		have    scope.StoreIDer
		wantErr error
	}{
		{scope.MockCode("nilSlices"), store.ErrStoreNotFound},
		{scope.MockID(2), store.ErrStoreNotFound},
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
	assert.NotEmpty(t, s.Data.Code.String)

	s, err = managerDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.NoError(t, err)
	assert.NotEmpty(t, s.Data.Code.String)
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
		haveID      scope.StoreIDer
		wantErr     error
	}{
		{tms, scope.MockID(1), nil},
		{tms, scope.MockID(1), store.ErrAppStoreSet},
		{tms, nil, store.ErrAppStoreSet},
		{tms, nil, store.ErrAppStoreSet},
	}

	for _, test := range tests {
		haveErr := test.haveManager.Init(scope.Option{Store: test.haveID})
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
		benchmarkManagerStore, err = managerStoreSimpleTest.Store(scope.MockCode("de"))
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
	assert.Equal(t, "at", ss[1].Data.Code.String)

	ss, err = managerStores.Stores()
	assert.NotNil(t, ss)
	assert.NoError(t, err)
	assert.NotEmpty(t, ss[2].Data.Code.String)

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
		have            scope.GroupIDer
		wantErr         error
		wantGroupName   string
		wantWebsiteCode string
	}{
		{managerGroupSimpleTest, nil, store.ErrAppStoreNotSet, "", ""},
		{getTestManager(), scope.MockID(20), store.ErrGroupNotFound, "", ""},
		{managerGroupSimpleTest, scope.MockID(1), nil, "DACH Group", "euro"},
		{managerGroupSimpleTest, scope.MockID(1), nil, "DACH Group", "euro"},
	}

	for _, test := range tests {
		g, err := test.m.Group(test.have)
		if test.wantErr != nil {
			assert.Nil(t, g)
			assert.EqualError(t, test.wantErr, err.Error(), "test %#v", test)
		} else {
			assert.NotNil(t, g, "test %#v", test)
			assert.NoError(t, err, "test %#v", test)
			assert.Equal(t, test.wantGroupName, g.Data.Name)
			assert.Equal(t, test.wantWebsiteCode, g.Website.Data.Code.String)
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
	}).Init(scope.Option{Group: scope.MockID(1)})
	assert.EqualError(t, store.ErrGroupDefaultStoreNotFound, err.Error(), "Incorrect DefaultStore for a Group")

	err = getTestManager().Init(scope.Option{Group: scope.MockID(21)})
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
	err = tm3.Init(scope.Option{Group: scope.MockID(1)})
	assert.NoError(t, err)
	g, err := tm3.Group()
	assert.NoError(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, int64(2), g.Data.DefaultStoreID)
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
			)
		}
	})

	tests := []struct {
		m               *store.Manager
		have            scope.WebsiteIDer
		wantErr         error
		wantWebsiteCode string
	}{
		{managerWebsite, nil, store.ErrAppStoreNotSet, ""},
		{getTestManager(), scope.MockID(20), store.ErrGroupNotFound, ""},
		{managerWebsite, scope.MockID(1), nil, "euro"},
		{managerWebsite, scope.MockID(1), nil, "euro"},
		{managerWebsite, scope.MockCode("notImportant"), nil, "euro"},
		{managerWebsite, scope.MockCode("notImportant"), nil, "euro"},
	}

	for _, test := range tests {
		haveW, haveErr := test.m.Website(test.have)
		if test.wantErr != nil {
			assert.Error(t, haveErr, "%#v", test)
			assert.Nil(t, haveW, "%#v", test)
		} else {
			assert.NoError(t, haveErr, "%#v", test)
			assert.NotNil(t, haveW, "%#v", test)
			assert.Equal(t, test.wantWebsiteCode, haveW.Data.Code.String)
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
			)
		}
	}).Init(scope.Option{Website: scope.MockCode("euro")})
	assert.EqualError(t, store.ErrWebsiteDefaultGroupNotFound, err.Error())

	managerWebsite := getTestManager(func(ms *mockStorage) {
		ms.w = func() (*store.Website, error) {
			nw, err := store.NewWebsite(
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			)
			if err != nil {
				return nil, err
			}
			return nw.SetGroupsStores(
				store.TableGroupSlice{
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				},
				store.TableStoreSlice{
					&store.TableStore{StoreID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
					&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
					&store.TableStore{StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
				},
			)
		}
	})
	w1, err := managerWebsite.Website()
	assert.EqualError(t, store.ErrAppStoreNotSet, err.Error())
	assert.Nil(t, w1)

	err = managerWebsite.Init(scope.Option{Website: scope.MockCode("euro")})
	assert.NoError(t, err)

	w2, err := managerWebsite.Website()
	assert.NoError(t, err)
	assert.EqualValues(t, "euro", w2.Data.Code.String)

	err3 := getTestManager(func(ms *mockStorage) {}).Init(scope.Option{Website: scope.MockCode("euronen")})
	assert.Error(t, err3, "scope.MockCode(euro), config.ScopeWebsite: %#v => %s", err3, err3)
	assert.EqualError(t, store.ErrWebsiteNotFound, err3.Error())
}

func TestNewManagerError(t *testing.T) {
	err := getTestManager().Init(scope.Option{})
	assert.EqualError(t, err, store.ErrUnsupportedScope.Error())
}

var storeManagerRequestStore = store.NewManager(
	store.NewStorage(
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
	haveSO        scope.Option
	wantStoreCode string
	wantErr       error
}

func runNewManagerGetRequestStore(t *testing.T, testScope scope.Scope, tests []testNewManagerGetRequestStore) {
	for i, test := range tests {
		haveStore, haveErr := storeManagerRequestStore.GetRequestStore(test.haveSO, testScope)
		if test.wantErr != nil {
			assert.Nil(t, haveStore, "%d: testScope %d: %#v", i, testScope, test)
			assert.EqualError(t, haveErr, test.wantErr.Error(), "%d: testScope %d: %#v", i, testScope, test)
		} else {
			assert.NotNil(t, haveStore)
			assert.NoError(t, haveErr, "%#v", test)
			assert.EqualValues(t, test.wantStoreCode, haveStore.Data.Code.String)
		}
	}
	storeManagerRequestStore.ClearCache(true)
}

func TestNewManagerGetRequestStore_ScopeStore(t *testing.T) {

	testCode := scope.MockCode("de")
	testScope := scope.StoreID

	if haveStore, haveErr := storeManagerRequestStore.GetRequestStore(scope.Option{Store: scope.MockID(1)}, testScope); haveErr == nil {
		t.Error("appStore should not be set!")
		t.Fail()
	} else {
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
	}

	// init with scope store
	if err := storeManagerRequestStore.Init(scope.Option{Store: testCode}); err != nil {
		t.Error(err)
		t.Fail()
	}
	assert.EqualError(t, store.ErrAppStoreSet, storeManagerRequestStore.Init(scope.Option{Store: testCode}).Error())

	if s, err := storeManagerRequestStore.Store(); err == nil {
		assert.EqualValues(t, "de", s.Data.Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	tests := []testNewManagerGetRequestStore{
		{scope.Option{Store: scope.MockID(232)}, "", store.ErrIDNotFoundTableStoreSlice},
		{scope.Option{}, "", store.ErrUnsupportedScope},
		{scope.Option{Store: scope.MockCode("\U0001f631")}, "", store.ErrIDNotFoundTableStoreSlice},

		{scope.Option{Store: scope.MockID(6)}, "nz", nil},
		{scope.Option{Store: scope.MockCode("ch")}, "", store.ErrStoreNotActive},

		{scope.Option{Store: scope.MockCode("nz")}, "nz", nil},
		{scope.Option{Store: scope.MockCode("de")}, "de", nil},
		{scope.Option{Store: scope.MockID(2)}, "at", nil},

		{scope.Option{Store: scope.MockID(2)}, "at", nil},
		{scope.Option{Store: scope.MockCode("au")}, "au", nil},
		{scope.Option{Store: scope.MockCode("ch")}, "", store.ErrStoreNotActive},
	}
	runNewManagerGetRequestStore(t, testScope, tests)
}

func TestNewManagerGetRequestStore_ScopeGroup(t *testing.T) {
	testOption := scope.Option{Group: scope.MockID(1)}
	testScope := scope.GroupID

	if haveStore, haveErr := storeManagerRequestStore.GetRequestStore(testOption, testScope); haveErr == nil {
		t.Error("appStore should not be set!")
		t.Fail()
	} else {
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
	}

	assert.EqualError(t, store.ErrIDNotFoundTableGroupSlice, storeManagerRequestStore.Init(scope.Option{Group: scope.MockID(123)}).Error())
	if err := storeManagerRequestStore.Init(testOption); err != nil {
		t.Error(err)
		t.Fail()
	}
	assert.EqualError(t, store.ErrAppStoreSet, storeManagerRequestStore.Init(testOption).Error())

	if s, err := storeManagerRequestStore.Store(); err == nil {
		assert.EqualValues(t, "at", s.Data.Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	if g, err := storeManagerRequestStore.Group(); err == nil {
		assert.EqualValues(t, 1, g.Data.GroupID)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	// we're testing here against Group ID = 1
	tests := []testNewManagerGetRequestStore{
		{scope.Option{Group: scope.MockID(232)}, "", store.ErrIDNotFoundTableGroupSlice},
		{scope.Option{Store: scope.MockID(232)}, "", store.ErrIDNotFoundTableStoreSlice},
		{scope.Option{}, "", store.ErrUnsupportedScope},
		{scope.Option{Store: scope.MockCode("\U0001f631")}, "", store.ErrIDNotFoundTableStoreSlice},

		{scope.Option{Store: scope.MockID(6)}, "nz", store.ErrStoreChangeNotAllowed},
		{scope.Option{Store: scope.MockCode("ch")}, "", store.ErrStoreNotActive},

		{scope.Option{Store: scope.MockCode("de")}, "de", nil},
		{scope.Option{Store: scope.MockID(2)}, "at", nil},

		{scope.Option{Store: scope.MockID(2)}, "at", nil},
		{scope.Option{Store: scope.MockCode("au")}, "au", store.ErrStoreChangeNotAllowed},
		{scope.Option{Store: scope.MockCode("ch")}, "", store.ErrStoreNotActive},

		{scope.Option{Group: scope.MockCode("ch")}, "", store.ErrIDNotFoundTableGroupSlice},
		{scope.Option{Group: scope.MockID(2)}, "", store.ErrStoreChangeNotAllowed},
		{scope.Option{Group: scope.MockID(1)}, "at", nil},

		{scope.Option{Website: scope.MockCode("xxxx")}, "", store.ErrIDNotFoundTableWebsiteSlice},
		{scope.Option{Website: scope.MockID(2)}, "", store.ErrStoreChangeNotAllowed},
		{scope.Option{Website: scope.MockID(1)}, "at", nil},
	}
	runNewManagerGetRequestStore(t, scope.GroupID, tests)
}

func TestNewManagerGetRequestStore_ScopeWebsite(t *testing.T) {
	testCode := scope.Option{Website: scope.MockID(1)}
	testScope := scope.WebsiteID

	if haveStore, haveErr := storeManagerRequestStore.GetRequestStore(testCode, testScope); haveErr == nil {
		t.Error("appStore should not be set!")
		t.Fail()
	} else {
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
	}

	assert.EqualError(t, store.ErrUnsupportedScope, storeManagerRequestStore.Init(scope.Option{}).Error())
	assert.EqualError(t, store.ErrIDNotFoundTableWebsiteSlice, storeManagerRequestStore.Init(scope.Option{Website: scope.MockID(123)}).Error())
	if err := storeManagerRequestStore.Init(testCode); err != nil {
		t.Error(err)
		t.Fail()
	}
	assert.EqualError(t, store.ErrAppStoreSet, storeManagerRequestStore.Init(testCode).Error())

	if s, err := storeManagerRequestStore.Store(); err == nil {
		assert.EqualValues(t, "at", s.Data.Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	if w, err := storeManagerRequestStore.Website(); err == nil {
		assert.EqualValues(t, "euro", w.Data.Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	// test against website euro
	tests := []testNewManagerGetRequestStore{
		{scope.Option{Website: scope.MockID(232)}, "", store.ErrIDNotFoundTableWebsiteSlice},
		{scope.Option{}, "", store.ErrUnsupportedScope},
		{scope.Option{Website: scope.MockCode("\U0001f631")}, "", store.ErrIDNotFoundTableWebsiteSlice},
		{scope.Option{Store: scope.MockCode("\U0001f631")}, "", store.ErrIDNotFoundTableStoreSlice},

		{scope.Option{Store: scope.MockID(6)}, "", store.ErrStoreChangeNotAllowed},
		{scope.Option{Website: scope.MockCode("oz")}, "", store.ErrStoreChangeNotAllowed},
		{scope.Option{Store: scope.MockCode("ch")}, "", store.ErrStoreNotActive},

		{scope.Option{Store: scope.MockCode("de")}, "de", nil},
		{scope.Option{Store: scope.MockID(2)}, "at", nil},

		{scope.Option{Store: scope.MockID(2)}, "at", nil},
		{scope.Option{Store: scope.MockCode("au")}, "au", store.ErrStoreChangeNotAllowed},
		{scope.Option{Store: scope.MockCode("ch")}, "", store.ErrStoreNotActive},

		{scope.Option{Group: scope.MockID(3)}, "", store.ErrStoreChangeNotAllowed},
	}
	runNewManagerGetRequestStore(t, scope.WebsiteID, tests)
}

func getTestRequest(m, u string, c *http.Cookie) *http.Request {
	req, err := http.NewRequest(m, u, nil)
	if err != nil {
		panic(err)
	}
	if c != nil {
		req.AddCookie(c)
	}
	return req
}

var testsInitByRequest = []struct {
	req                  *http.Request
	haveSO               scope.Option
	haveScopeType        scope.Scope
	wantStoreCode        string           // this is the default store in a scope, lookup in storeManagerRequestStore
	wantRequestStoreCode scope.StoreCoder // can be nil in tests
	wantErr              error
	wantCookie           string
}{
	{
		getTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "uk"}),
		scope.Option{Store: scope.MockID(1)}, scope.StoreID, "de", scope.MockCode("uk"), nil, store.CookieName + "=uk;",
	},
	{
		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
		scope.Option{Store: scope.MockID(1)}, scope.StoreID, "de", scope.MockCode("uk"), nil, store.CookieName + "=uk;", // generates a new 1y valid cookie
	},
	{
		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=%20uk", nil),
		scope.Option{Store: scope.MockID(1)}, scope.StoreID, "de", scope.MockCode("uk"), store.ErrStoreCodeInvalid, "",
	},
	{
		getTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "de"}),
		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", scope.MockCode("de"), nil, store.CookieName + "=de;",
	},
	{
		getTestRequest("GET", "http://cs.io", nil),
		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", nil, store.ErrUnsupportedScope, "",
	},
	{
		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=de", nil),
		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", scope.MockCode("de"), nil, store.CookieName + "=de;", // generates a new 1y valid cookie
	},
	{
		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=at", nil),
		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", scope.MockCode("at"), nil, store.CookieName + "=;", // generates a delete cookie
	},
	{
		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=cz", nil),
		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", nil, store.ErrIDNotFoundTableStoreSlice, "",
	},
	{
		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", nil, store.ErrStoreChangeNotAllowed, "",
	},

	{
		getTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "nz"}),
		scope.Option{Website: scope.MockID(2)}, scope.WebsiteID, "au", scope.MockCode("nz"), nil, store.CookieName + "=nz;",
	},
	{
		getTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "n'z"}),
		scope.Option{Website: scope.MockID(2)}, scope.WebsiteID, "au", nil, store.ErrStoreCodeInvalid, "",
	},
	{
		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
		scope.Option{Website: scope.MockID(2)}, scope.WebsiteID, "au", nil, store.ErrStoreChangeNotAllowed, "",
	},
	{
		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=nz", nil),
		scope.Option{Website: scope.MockID(2)}, scope.WebsiteID, "au", scope.MockCode("nz"), nil, store.CookieName + "=nz;",
	},
	{
		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=ch", nil),
		scope.Option{Website: scope.MockID(1)}, scope.WebsiteID, "at", nil, store.ErrStoreNotActive, "",
	},
	{
		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=nz", nil),
		scope.Option{Website: scope.MockID(1)}, scope.DefaultID, "at", scope.MockCode("nz"), store.ErrStoreChangeNotAllowed, "",
	},
}

func TestInitByRequestGeneral(t *testing.T) {
	errLogBuf.Reset()
	defer errLogBuf.Reset()

	for _, test := range testsInitByRequest {
		if _, haveErr := storeManagerRequestStore.InitByRequest(nil, nil, test.haveScopeType); haveErr != nil {
			assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
		} else {
			t.Fatal("InitByRequest should return an error if used without running Init() first.")
		}

		if err := storeManagerRequestStore.Init(test.haveSO); err != nil {
			assert.EqualError(t, store.ErrUnsupportedScope, err.Error())
			t.Log("continuing for loop because of expected store.ErrUnsupportedScopeGroup")
			storeManagerRequestStore.ClearCache(true)
			continue
		}

		if s, err := storeManagerRequestStore.Store(); err == nil {
			assert.EqualValues(t, test.wantStoreCode, s.Data.Code.String)
		} else {
			assert.EqualError(t, err, store.ErrStoreNotFound.Error())
			t.Log("continuing for loop because of expected store.ErrStoreNotFound")
			storeManagerRequestStore.ClearCache(true)
			continue
		}
		storeManagerRequestStore.ClearCache(true)
	}
}

func TestInitByRequestInDepth(t *testing.T) {
	errLogBuf.Reset()
	defer errLogBuf.Reset()

	for i, test := range testsInitByRequest {
		if err := storeManagerRequestStore.Init(test.haveSO); err != nil {
			assert.EqualError(t, store.ErrUnsupportedScope, err.Error())
			t.Log("continuing for loop because of expected store.ErrUnsupportedScopeGroup")
			storeManagerRequestStore.ClearCache(true)
			continue
		}

		resRec := httptest.NewRecorder()

		haveStore, haveErr := storeManagerRequestStore.InitByRequest(resRec, test.req, test.haveScopeType)
		if test.wantErr != nil {
			assert.Nil(t, haveStore)
			assert.Error(t, haveErr, "Index %d", i)
			assert.EqualError(t, haveErr, test.wantErr.Error(), "\nIndex: %d\nError: %s", i, errLogBuf.String())
		} else {

			assert.NoError(t, haveErr, "Test: %#v\n\n%s\n\n", test, errLogBuf.String())

			if test.wantRequestStoreCode != nil {
				assert.NotNil(t, haveStore, "URL Query: %#v\nCookies %#v", test.req.URL.Query(), test.req.Cookies())
				assert.EqualValues(t, test.wantRequestStoreCode.StoreCode(), haveStore.Data.Code.String)

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
		haveSO             scope.Option
		haveCodeToken      string
		haveScopeType      scope.Scope
		wantStoreCode      string           // this is the default store in a scope, lookup in storeManagerRequestStore
		wantTokenStoreCode scope.StoreCoder // can be nil
		wantErr            error
	}{
		{scope.Option{Store: scope.MockCode("de")}, "de", scope.StoreID, "de", scope.MockCode("de"), nil},
		{scope.Option{Store: scope.MockCode("de")}, "at", scope.StoreID, "de", scope.MockCode("at"), nil},
		{scope.Option{Store: scope.MockCode("de")}, "a$t", scope.StoreID, "de", nil, nil},
		{scope.Option{Store: scope.MockCode("at")}, "ch", scope.StoreID, "at", nil, store.ErrStoreNotActive},
		{scope.Option{Store: scope.MockCode("at")}, "", scope.StoreID, "at", nil, nil},

		{scope.Option{Group: scope.MockID(1)}, "de", scope.GroupID, "at", scope.MockCode("de"), nil},
		{scope.Option{Group: scope.MockID(1)}, "ch", scope.GroupID, "at", nil, store.ErrStoreNotActive},
		{scope.Option{Group: scope.MockID(1)}, " ch", scope.GroupID, "at", nil, nil},
		{scope.Option{Group: scope.MockID(1)}, "uk", scope.GroupID, "at", nil, store.ErrStoreChangeNotAllowed},

		{scope.Option{Website: scope.MockID(2)}, "uk", scope.WebsiteID, "au", nil, store.ErrStoreChangeNotAllowed},
		{scope.Option{Website: scope.MockID(2)}, "nz", scope.WebsiteID, "au", scope.MockCode("nz"), nil},
		{scope.Option{Website: scope.MockID(2)}, "n z", scope.WebsiteID, "au", nil, nil},
		{scope.Option{Website: scope.MockID(2)}, "", scope.WebsiteID, "au", nil, nil},
	}
	for _, test := range tests {

		haveStore, haveErr := storeManagerRequestStore.InitByToken(nil, test.haveScopeType)
		assert.Nil(t, haveStore)
		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())

		if err := storeManagerRequestStore.Init(test.haveSO); err != nil {
			t.Fatal(err)
		}

		if s, err := storeManagerRequestStore.Store(); err == nil {
			assert.EqualValues(t, test.wantStoreCode, s.Data.Code.String)
		} else {
			assert.EqualError(t, err, store.ErrStoreNotFound.Error())
			t.Fail()
		}

		haveStore, haveErr = storeManagerRequestStore.InitByToken(getToken(test.haveCodeToken).Claims, test.haveScopeType)
		if test.wantErr != nil {
			assert.Nil(t, haveStore, "%#v", test)
			assert.Error(t, haveErr, "%#v", test)
			assert.EqualError(t, test.wantErr, haveErr.Error())
		} else {
			if test.wantTokenStoreCode != nil {
				assert.NotNil(t, haveStore, "%#v", test)
				assert.NoError(t, haveErr)
				assert.Equal(t, test.wantTokenStoreCode.StoreCode(), haveStore.Data.Code.String)
			} else {
				assert.Nil(t, haveStore, "%#v", test)
				assert.NoError(t, haveErr, "%#v", test)
			}

		}
		storeManagerRequestStore.ClearCache(true)
	}
}

func TestNewManagerReInit(t *testing.T) {

	t.Skip(TODO_Better_Test_Data)

	// quick implement, use mock of dbr.SessionRunner and remove connection
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	dbrSess := dbc.NewSession()

	storeManager := store.NewManager(store.NewStorage(nil /* trick it*/))
	if err := storeManager.ReInit(dbrSess); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		have    scope.StoreIDer
		wantErr error
	}{
		{scope.MockCode("dede"), nil},
		{scope.MockCode("czcz"), store.ErrIDNotFoundTableStoreSlice},
		{scope.MockID(1), nil},
		{scope.MockID(100), store.ErrStoreNotFound},
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
			//			assert.NotEmpty(t, s.Data.Code.String, "%#v", s.Data)
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

func (ic mockIDCode) StoreID() int64 {
	return ic.id
}
func (ic mockIDCode) StoreCode() string {
	return ic.code
}
func (ic mockIDCode) WebsiteID() int64 {
	return ic.id
}
func (ic mockIDCode) WebsiteCode() string {
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

func (ms *mockStorage) Website(_ scope.WebsiteIDer) (*store.Website, error) {
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
func (ms *mockStorage) Group(_ scope.GroupIDer) (*store.Group, error) {
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
func (ms *mockStorage) Store(_ scope.StoreIDer) (*store.Store, error) {
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

type mockManager struct{}

func (m *mockManager) IsSingleStoreMode() bool { return false }
func (m *mockManager) HasSingleStore() bool    { return false }
func (m *mockManager) Website(r ...scope.WebsiteIDer) (*store.Website, error) {
	return nil, store.ErrWebsiteNotFound
}
func (m *mockManager) Websites() (store.WebsiteSlice, error) { return nil, store.ErrWebsiteNotFound }
func (m *mockManager) Group(r ...scope.GroupIDer) (*store.Group, error) {
	return nil, store.ErrGroupNotFound
}
func (m *mockManager) Groups() (store.GroupSlice, error) { return nil, store.ErrGroupNotFound }
func (m *mockManager) Store(r ...scope.StoreIDer) (*store.Store, error) {
	return nil, store.ErrStoreNotFound
}
func (m *mockManager) Stores() (store.StoreSlice, error)       { return nil, store.ErrStoreNotFound }
func (m *mockManager) DefaultStoreView() (*store.Store, error) { return nil, store.ErrStoreNotFound }

var _ store.ManagerReader = (*mockManager)(nil)
