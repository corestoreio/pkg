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
	"testing"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	storemock "github.com/corestoreio/csfw/store/mock"
	"github.com/corestoreio/csfw/utils/log"
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
// Reminder to myself:
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

var serviceStoreSimpleTest = storemock.MustNewService(scope.Option{}, func(ms *storemock.Storage) {
	ms.MockStore = func() (*store.Store, error) {
		return store.NewStore(
			&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true, true)},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		)
	}
})

func TestNewServiceStore(t *testing.T) {
	assert.True(t, serviceStoreSimpleTest.IsCacheEmpty())
	for j := 0; j < 3; j++ {
		s, err := serviceStoreSimpleTest.Store(scope.MockCode("notNil"))
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.EqualValues(t, "de", s.Data.Code.String)
	}
	assert.False(t, serviceStoreSimpleTest.IsCacheEmpty())
	serviceStoreSimpleTest.ClearCache()
	assert.True(t, serviceStoreSimpleTest.IsCacheEmpty())

}

func TestMustNewService(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), store.ErrStoreNotFound.Error())
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	tests := []struct {
		have    scope.StoreIDer
		wantErr error
	}{
		{scope.MockCode("nilSlices"), store.ErrStoreNotFound},
		{scope.MockID(2), store.ErrStoreNotFound},
		{nil, store.ErrStoreNotFound},
	}
	serviceEmpty := storemock.MustNewService(scope.Option{})
	for _, test := range tests {
		s, err := serviceEmpty.Store(test.have)
		assert.Nil(t, s)
		assert.EqualError(t, test.wantErr, err.Error())
	}
	assert.True(t, serviceStoreSimpleTest.IsCacheEmpty())

}

func TestNewServiceDefaultStoreView(t *testing.T) {
	serviceDefaultStore := storemock.MustNewService(scope.Option{}, func(ms *storemock.Storage) {
		ms.MockDefaultStore = func() (*store.Store, error) {
			return store.NewStore(
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true, true)},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			)
		}
	})

	// call it twice to test internal caching
	s, err := serviceDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.NoError(t, err)
	assert.NotEmpty(t, s.Data.Code.String)

	s, err = serviceDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.NoError(t, err)
	assert.NotEmpty(t, s.Data.Code.String)
	assert.False(t, serviceDefaultStore.IsCacheEmpty())
	serviceDefaultStore.ClearCache()
	assert.True(t, serviceDefaultStore.IsCacheEmpty())
}

var benchmarkServiceStore *store.Store

// BenchmarkServiceGetStore-4              	 5000000	       256 ns/op	      16 B/op	       1 allocs/op
func BenchmarkServiceGetStore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkServiceStore, err = serviceStoreSimpleTest.Store(scope.MockCode("de"))
		if err != nil {
			b.Error(err)
		}
		if benchmarkServiceStore == nil {
			b.Error("benchmarkServiceStore is nil")
		}
	}
}

func TestNewServiceStores(t *testing.T) {
	serviceStores := storemock.MustNewService(scope.Option{}, func(ms *storemock.Storage) {

		ms.MockDefaultStore = func() (*store.Store, error) {
			return store.MustNewStore(
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			), nil
		}

		ms.MockStoreSlice = func() (store.StoreSlice, error) {
			return store.StoreSlice{
				store.MustNewStore(
					&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				),
				store.MustNewStore(
					&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				),
				store.MustNewStore(
					&store.TableStore{StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				),
			}, nil
		}
	})

	// call it twice to test internal caching
	ss, err := serviceStores.Stores()
	assert.NotNil(t, ss)
	assert.NoError(t, err)
	assert.Equal(t, "at", ss[1].Data.Code.String)

	ss, err = serviceStores.Stores()
	assert.NotNil(t, ss)
	assert.NoError(t, err)
	assert.NotEmpty(t, ss[2].Data.Code.String)

	assert.False(t, serviceStores.IsCacheEmpty())
	serviceStores.ClearCache()
	assert.True(t, serviceStores.IsCacheEmpty())
}

func TestMustNewServiceStores(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), store.ErrStoreNotFound.Error())
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	ss, err := storemock.MustNewService(scope.Option{}, func(ms *storemock.Storage) {
		ms.MockStoreSlice = func() (store.StoreSlice, error) {
			return nil, nil
		}
	}).Stores()
	assert.Nil(t, ss)
	assert.NoError(t, err)
}

func TestNewServiceGroup(t *testing.T) {
	var serviceGroupSimpleTest = storemock.MustNewService(scope.Option{}, func(ms *storemock.Storage) {
		ms.MockGroup = func() (*store.Group, error) {
			return store.NewGroup(
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				store.SetGroupWebsite(&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}}),
			)
		}
		ms.MockStore = func() (*store.Store, error) {
			return store.NewStore(
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			)
		}
	})

	tests := []struct {
		m               *store.Service
		have            scope.GroupIDer
		wantErr         error
		wantGroupName   string
		wantWebsiteCode string
	}{
		{serviceGroupSimpleTest, scope.MockID(20), nil, "DACH Group", "euro"},
		{serviceGroupSimpleTest, scope.MockID(1), nil, "DACH Group", "euro"},
		{serviceGroupSimpleTest, scope.MockID(1), nil, "DACH Group", "euro"},
	}

	for i, test := range tests {
		g, err := test.m.Group(test.have)
		if test.wantErr != nil {
			assert.Nil(t, g, "Index %d", i)
			assert.EqualError(t, test.wantErr, err.Error(), "test %#v", test)
		} else {
			assert.NotNil(t, g, "test %#v", test)
			assert.NoError(t, err, "test %#v", test)
			assert.Equal(t, test.wantGroupName, g.Data.Name)
			assert.Equal(t, test.wantWebsiteCode, g.Website.Data.Code.String)
		}
	}
	assert.False(t, serviceGroupSimpleTest.IsCacheEmpty())
	serviceGroupSimpleTest.ClearCache()
	assert.True(t, serviceGroupSimpleTest.IsCacheEmpty())
}

func TestNewServiceGroups(t *testing.T) {
	serviceGroups := storemock.MustNewService(scope.Option{}, func(ms *storemock.Storage) {
		ms.MockGroupSlice = func() (store.GroupSlice, error) {
			return store.GroupSlice{}, nil
		}
		ms.MockStore = func() (*store.Store, error) {
			return store.NewStore(
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			)
		}
	})

	// call it twice to test internal caching
	ss, err := serviceGroups.Groups()
	assert.NotNil(t, ss)
	assert.NoError(t, err)
	assert.Len(t, ss, 0)

	ss, err = serviceGroups.Groups()
	assert.NotNil(t, ss)
	assert.NoError(t, err)
	assert.Len(t, ss, 0)

	assert.False(t, serviceGroups.IsCacheEmpty())
	serviceGroups.ClearCache()
	assert.True(t, serviceGroups.IsCacheEmpty())
}

func TestNewServiceWebsite(t *testing.T) {

	var serviceWebsite = storemock.MustNewService(scope.Option{}, func(ms *storemock.Storage) {
		ms.MockWebsite = func() (*store.Website, error) {
			return store.NewWebsite(
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			)
		}
		ms.MockStore = func() (*store.Store, error) {
			return store.NewStore(
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			)
		}
	})

	tests := []struct {
		m               *store.Service
		have            scope.WebsiteIDer
		wantErr         error
		wantWebsiteCode string
	}{
		{serviceWebsite, scope.MockID(1), nil, "euro"},
		{serviceWebsite, scope.MockID(1), nil, "euro"},
		{serviceWebsite, scope.MockCode("notImportant"), nil, "euro"},
		{serviceWebsite, scope.MockCode("notImportant"), nil, "euro"},
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
	assert.False(t, serviceWebsite.IsCacheEmpty())
	serviceWebsite.ClearCache()
	assert.True(t, serviceWebsite.IsCacheEmpty())

}

func TestNewServiceWebsites(t *testing.T) {
	serviceWebsites := storemock.MustNewService(scope.Option{}, func(ms *storemock.Storage) {
		ms.MockWebsiteSlice = func() (store.WebsiteSlice, error) {
			return store.WebsiteSlice{}, nil
		}
		ms.MockStore = func() (*store.Store, error) {
			return store.NewStore(
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			)
		}
	})

	tests := []struct {
		m       *store.Service
		wantErr error
		wantNil bool
	}{
		{serviceWebsites, nil, false},
		{serviceWebsites, nil, false},
		{storemock.MustNewService(scope.Option{}, func(ms *storemock.Storage) {
			ms.MockWebsiteSlice = func() (store.WebsiteSlice, error) {
				return nil, nil
			}
			ms.MockStore = func() (*store.Store, error) {
				return store.NewStore(
					&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				)
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

	assert.False(t, serviceWebsites.IsCacheEmpty())
	serviceWebsites.ClearCache()
	assert.True(t, serviceWebsites.IsCacheEmpty())
}

func getInitializedStoreService(so scope.Option) *store.Service {
	return store.MustNewService(so,
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
}

type testNewServiceGetRequestStore struct {
	haveSO        scope.Option
	wantStoreCode string
	wantErr       error
}

func runTestsGetRequestedStore(t *testing.T, sm *store.Service, tests []testNewServiceGetRequestStore) {
	for i, test := range tests {
		haveStore, haveErr := sm.GetRequestedStore(test.haveSO)
		if test.wantErr != nil {
			assert.Nil(t, haveStore, "Index: %d: %#v", i, test)
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index: %d: %#v", i, test)
		} else {
			assert.NotNil(t, haveStore)
			assert.NoError(t, haveErr, "%#v", test)
			assert.EqualValues(t, test.wantStoreCode, haveStore.Data.Code.String)
		}
	}
	sm.ClearCache(true)
}

func TestNewServiceGetRequestStore_ScopeStore(t *testing.T) {

	initScope := scope.Option{Store: scope.MockID(1)}
	sm := getInitializedStoreService(initScope)

	if haveStore, haveErr := sm.GetRequestedStore(initScope); haveErr != nil {
		t.Fatal(haveErr)
	} else {
		assert.NoError(t, haveErr)
		assert.Exactly(t, int64(1), haveStore.StoreID())
	}

	if s, err := sm.Store(); err == nil {
		assert.EqualValues(t, "de", s.Data.Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	if s, err := sm.Store(scope.MockID(123)); err == nil {
		assert.Error(t, err)
		t.Fail()
	} else {
		assert.Nil(t, s)
		assert.EqualError(t, err, store.ErrIDNotFoundTableStoreSlice.Error())
	}

	tests := []testNewServiceGetRequestStore{
		{scope.Option{Store: scope.MockID(232)}, "", store.ErrIDNotFoundTableStoreSlice},
		{scope.Option{}, "at", nil},
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
	runTestsGetRequestedStore(t, sm, tests)
}

func TestNewServiceGetRequestStore_ScopeGroup(t *testing.T) {
	initScope := scope.Option{Group: scope.MockID(1)}

	sm := getInitializedStoreService(initScope)
	if haveStore, haveErr := sm.GetRequestedStore(initScope); haveErr != nil {
		t.Fatal(haveErr)
	} else {
		assert.NoError(t, haveErr)
		assert.Exactly(t, int64(2), haveStore.StoreID())
	}

	if s, err := sm.Store(); err == nil {
		assert.EqualValues(t, "at", s.Data.Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	if g, err := sm.Group(); err == nil {
		assert.EqualValues(t, 1, g.Data.GroupID)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	//	// we're testing here against Group ID = 1
	tests := []testNewServiceGetRequestStore{
		{scope.Option{Group: scope.MockID(232)}, "", store.ErrIDNotFoundTableGroupSlice},
		{scope.Option{Store: scope.MockID(232)}, "", store.ErrIDNotFoundTableStoreSlice},
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
	runTestsGetRequestedStore(t, sm, tests)
}

func TestNewServiceGetRequestStore_ScopeWebsite(t *testing.T) {
	initScope := scope.Option{Website: scope.MockID(1)}

	sm := getInitializedStoreService(initScope)

	if haveStore, haveErr := sm.GetRequestedStore(initScope); haveErr != nil {
		t.Fatal(haveErr)
	} else {
		assert.NoError(t, haveErr)
		assert.Exactly(t, int64(2), haveStore.StoreID())
	}

	if s, err := sm.Store(); err == nil {
		assert.EqualValues(t, "at", s.Data.Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	if w, err := sm.Website(); err == nil {
		assert.EqualValues(t, "euro", w.Data.Code.String)
	} else {
		assert.EqualError(t, err, store.ErrStoreNotFound.Error())
		t.Fail()
	}

	// test against website euro
	tests := []testNewServiceGetRequestStore{
		{scope.Option{Website: scope.MockID(232)}, "", store.ErrIDNotFoundTableWebsiteSlice},
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
	runTestsGetRequestedStore(t, sm, tests)
}

func TestNewServiceReInit(t *testing.T) {

	t.Skip(TODO_Better_Test_Data)

	// quick implement, use mock of dbr.SessionRunner and remove connection
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	dbrSess := dbc.NewSession()

	storeService := store.MustNewService(scope.Option{}, store.NewStorage(nil /* trick it*/))
	if err := storeService.ReInit(dbrSess); err != nil {
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
	}

	for _, test := range tests {
		s, err := storeService.Store(test.have)
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
	assert.False(t, storeService.IsCacheEmpty())
	storeService.ClearCache()
	assert.True(t, storeService.IsCacheEmpty())
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
