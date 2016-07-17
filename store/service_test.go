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
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ store.Requester = (*store.Service)(nil)
var _ store.CodeToIDMapper = (*store.Service)(nil)

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

var serviceStoreSimpleTest = storemock.MustNewService(0, func(ms *storemock.Storage) {
	ms.MockStore = func() (*store.Store, error) {
		return store.NewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		)
	}
})

func TestNewServiceStore(t *testing.T) {

	assert.False(t, serviceStoreSimpleTest.IsCacheEmpty())

	s, err := serviceStoreSimpleTest.Store(-1)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.EqualValues(t, "de", s.Data.Code.String)

	assert.False(t, serviceStoreSimpleTest.IsCacheEmpty())
	serviceStoreSimpleTest.ClearCache()
	assert.True(t, serviceStoreSimpleTest.IsCacheEmpty())

}

func TestMustNewService(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotFound(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	tests := []struct {
		have       int64
		wantErrBhf errors.BehaviourFunc
	}{
		{-1, errors.IsNotFound},
		{4444, errors.IsNotFound},
		{0, errors.IsNotFound},
	}
	serviceEmpty := storemock.MustNewService(0)
	for i, test := range tests {
		s, err := serviceEmpty.Store(test.have)
		assert.Nil(t, s, "Index %d")
		assert.True(t, test.wantErrBhf(err), "Index %d => %s", i, err)
	}
	assert.True(t, serviceStoreSimpleTest.IsCacheEmpty())
}

func TestNewServiceDefaultStoreView(t *testing.T) {

	serviceDefaultStore := storemock.MustNewService(0, func(ms *storemock.Storage) {
		ms.MockStore = func() (*store.Store, error) {
			return store.NewStore(
				cfgmock.NewService(),
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
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
		benchmarkServiceStore, err = serviceStoreSimpleTest.Store(1)
		if err != nil {
			b.Error(err)
		}
		if benchmarkServiceStore == nil {
			b.Error("benchmarkServiceStore is nil")
		}
	}
}

func TestNewServiceStores(t *testing.T) {

	serviceStores := storemock.MustNewService(0, func(ms *storemock.Storage) {

		ms.MockStore = func() (*store.Store, error) {
			return store.MustNewStore(
				cfgmock.NewService(),
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			), nil
		}

		ms.MockStoreSlice = func() (store.StoreSlice, error) {
			cfg := cfgmock.NewService()
			return store.StoreSlice{
				store.MustNewStore(
					cfg,
					&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				),
				store.MustNewStore(
					cfg,
					&store.TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				),
				store.MustNewStore(
					cfg,
					&store.TableStore{StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
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
			err := r.(error)
			assert.True(t, errors.IsNotFound(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	ss, err := storemock.MustNewService(0, func(ms *storemock.Storage) {
		ms.MockStoreSlice = func() (store.StoreSlice, error) {
			return nil, nil
		}
	}).Stores()
	assert.Nil(t, ss)
	assert.NoError(t, err)
}

func TestNewServiceGroup(t *testing.T) {

	var serviceGroupSimpleTest = storemock.MustNewService(0, func(ms *storemock.Storage) {
		ms.MockGroup = func() (*store.Group, error) {
			return store.NewGroup(
				cfgmock.NewService(),
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				store.SetGroupWebsite(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
			)
		}
		ms.MockStore = func() (*store.Store, error) {
			return store.NewStore(
				cfgmock.NewService(),
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			)
		}
	})

	tests := []struct {
		m               *store.Service
		have            int64
		wantErr         error
		wantGroupName   string
		wantWebsiteCode string
	}{
		{serviceGroupSimpleTest, 20, nil, "DACH Group", "euro"},
		{serviceGroupSimpleTest, 1, nil, "DACH Group", "euro"},
		{serviceGroupSimpleTest, 1, nil, "DACH Group", "euro"},
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

	serviceGroups := storemock.MustNewService(0, func(ms *storemock.Storage) {
		ms.MockGroupSlice = func() (store.GroupSlice, error) {
			return store.GroupSlice{}, nil
		}
		ms.MockStore = func() (*store.Store, error) {
			return store.NewStore(
				cfgmock.NewService(),
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
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

	var serviceWebsite = storemock.MustNewService(0, func(ms *storemock.Storage) {
		ms.MockWebsite = func() (*store.Website, error) {
			return store.NewWebsite(
				cfgmock.NewService(),
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			)
		}
		ms.MockStore = func() (*store.Store, error) {
			return store.NewStore(
				cfgmock.NewService(),
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			)
		}
	})

	tests := []struct {
		m               *store.Service
		have            int64
		wantErr         error
		wantWebsiteCode string
	}{
		{serviceWebsite, 1, nil, "euro"},
		{serviceWebsite, 0, nil, "euro"},
		{serviceWebsite, 0, nil, "euro"},
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

	serviceWebsites := storemock.MustNewService(0, func(ms *storemock.Storage) {
		ms.MockWebsiteSlice = func() (store.WebsiteSlice, error) {
			return store.WebsiteSlice{}, nil
		}
		ms.MockStore = func() (*store.Store, error) {
			return store.NewStore(
				cfgmock.NewService(),
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
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
		{storemock.MustNewService(0, func(ms *storemock.Storage) {
			ms.MockWebsiteSlice = func() (store.WebsiteSlice, error) {
				return nil, nil
			}
			ms.MockStore = func() (*store.Store, error) {
				return store.NewStore(
					cfgmock.NewService(),
					&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
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

type testNewServiceRequestedStore struct {
	haveStoreID   int64
	wantStoreCode string
	wantErrBhf    errors.BehaviourFunc
}

func runTestsRequestedStore(t *testing.T, sm *store.Service, tests []testNewServiceRequestedStore) {
	for i, test := range tests {
		haveStore, haveErr := sm.RequestedStore(test.haveStoreID)
		if test.wantErrBhf != nil {
			assert.Nil(t, haveStore, "Index: %d: %#v", i, test)
			assert.True(t, test.wantErrBhf(haveErr), "Index %d Error: %s", i, haveErr)
		} else {
			assert.NotNil(t, haveStore, "Index %d", i)
			assert.NoError(t, haveErr, "Index %d => %#v", i, test)
			assert.EqualValues(t, test.wantStoreCode, haveStore.Data.Code.String, "Index %d", i)
		}
	}
	sm.ClearCache(true)
}

//func TestNewServiceRequestedStore_ScopeStore(t *testing.T) {
//
//	sm := storemock.NewEurozzyService(scope.NewHash(scope.Store, 1), cfgmock.NewService())
//
//	if haveStore, haveErr := sm.RequestedStore(1); haveErr != nil {
//		t.Fatal(haveErr)
//	} else {
//		assert.NoError(t, haveErr)
//		assert.Exactly(t, int64(1), haveStore.StoreID())
//	}
//
//	if s, err := sm.Store(); err == nil {
//		assert.EqualValues(t, "de", s.Data.Code.String)
//	} else {
//		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
//		t.Fail()
//	}
//
//	if s, err := sm.Store(123); err == nil {
//		assert.Error(t, err)
//		t.Fail()
//	} else {
//		assert.Nil(t, s)
//		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
//	}
//
//	tests := []testNewServiceRequestedStore{
//		{232 , "", errors.IsNotFound},
//		{0, "de", errors.IsNotSupported},
//
//
//		{6 , "nz", nil},
//		{3, "", errors.IsUnauthorized},
//
//		{1, "de", nil},
//		{2 , "at", nil},
//
//		{2 , "at", nil},
//		{5, "au", nil},
//		{3, "", errors.IsUnauthorized},
//	}
//	runTestsRequestedStore(t, sm, tests)
//}
//
//func TestNewServiceRequestedStore_ScopeGroup(t *testing.T) {
//
//	sm := storemock.NewEurozzyService(scope.NewHash(scope.Group,1),cfgmock.NewService())
//	if haveStore, haveErr := sm.RequestedStore(1); haveErr != nil {
//		t.Fatal(haveErr)
//	} else {
//		assert.NoError(t, haveErr)
//		assert.Exactly(t, int64(2), haveStore.StoreID())
//	}
//
//	if s, err := sm.Store(); err == nil {
//		assert.EqualValues(t, "at", s.Data.Code.String)
//	} else {
//		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
//		t.Fail()
//	}
//
//	if g, err := sm.Group(); err == nil {
//		assert.EqualValues(t, 1, g.Data.GroupID)
//	} else {
//		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
//		t.Fail()
//	}
//
//	//	// we're testing here against Group ID = 1
//	tests := []testNewServiceRequestedStore{
//		{scope.Option{Group: scope.MockID(232)}, "", errors.IsNotFound},
//
//		{232)}, "", errors.IsNotFound},
//		{scope.Option{Store: scope.MockCode("\U0001f631")}, "", errors.IsNotFound},
//
//		{6)}, "nz", errors.IsUnauthorized},
//		{scope.Option{Store: scope.MockCode("ch")}, "", errors.IsUnauthorized},
//
//		{scope.Option{Store: scope.MockCode("de")}, "de", nil},
//		{2)}, "at", nil},
//
//		{2)}, "at", nil},
//		{scope.Option{Store: scope.MockCode("au")}, "au", errors.IsUnauthorized},
//		{scope.Option{Store: scope.MockCode("ch")}, "", errors.IsUnauthorized},
//
//		{scope.Option{Group: scope.MockCode("ch")}, "", errors.IsNotFound},
//		{scope.Option{Group: scope.MockID(2)}, "", errors.IsUnauthorized},
//		{scope.Option{Group: scope.MockID(1)}, "at", nil},
//
//		{scope.Option{Website: scope.MockCode("xxxx")}, "", errors.IsNotFound},
//		{scope.Option{Website: scope.MockID(2)}, "", errors.IsUnauthorized},
//		{scope.Option{Website: scope.MockID(1)}, "at", nil},
//	}
//	runTestsRequestedStore(t, sm, tests)
//}
//
//func TestNewServiceRequestedStore_ScopeWebsite(t *testing.T) {
//
//	sm := storemock.NewEurozzyService(scope.NewHash(scope.Website,1),cfgmock.NewService())
//
//	if haveStore, haveErr := sm.RequestedStore(1); haveErr != nil {
//		t.Fatal(haveErr)
//	} else {
//		assert.NoError(t, haveErr)
//		assert.Exactly(t, int64(2), haveStore.StoreID())
//	}
//
//	if s, err := sm.Store(); err == nil {
//		assert.EqualValues(t, "at", s.Data.Code.String)
//	} else {
//		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
//		t.Fail()
//	}
//
//	if w, err := sm.Website(); err == nil {
//		assert.EqualValues(t, "euro", w.Data.Code.String)
//	} else {
//		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
//		t.Fail()
//	}
//
//	// test against website euro
//	tests := []testNewServiceRequestedStore{
//		{scope.Option{Website: scope.MockID(232)}, "", errors.IsNotFound},
//		{scope.Option{Website: scope.MockCode("\U0001f631")}, "", errors.IsNotFound},
//		{scope.Option{Store: scope.MockCode("\U0001f631")}, "", errors.IsNotFound},
//
//		{6)}, "", errors.IsUnauthorized},
//		{scope.Option{Website: scope.MockCode("oz")}, "", errors.IsUnauthorized},
//		{scope.Option{Store: scope.MockCode("ch")}, "", errors.IsUnauthorized},
//
//		{scope.Option{Store: scope.MockCode("de")}, "de", nil},
//		{2)}, "at", nil},
//
//		{2)}, "at", nil},
//		{scope.Option{Store: scope.MockCode("au")}, "au", errors.IsUnauthorized},
//		{scope.Option{Store: scope.MockCode("ch")}, "", errors.IsUnauthorized},
//
//		{scope.Option{Group: scope.MockID(3)}, "", errors.IsUnauthorized},
//	}
//	runTestsRequestedStore(t, sm, tests)
//}

func TestNewServiceReInit(t *testing.T) {

	t.Skip(TODO_Better_Test_Data)

	//// quick implement, use mock of dbr.SessionRunner and remove connection
	//dbc := csdb.MustConnectTest()
	//defer func() { assert.NoError(t, dbc.Close()) }()
	//dbrSess := dbc.NewSession()
	//
	//storeService := store.MustNewService(0, store.MustNewStorage(nil /* trick it*/))
	//if err := storeService.ReInit(dbrSess); err != nil {
	//	t.Fatal(err)
	//}
	//
	//tests := []struct {
	//	have       scope.StoreIDer
	//	wantErrBhf errors.BehaviourFunc
	//}{
	//	{scope.MockCode("dede"), nil},
	//	{scope.MockCode("czcz"), errors.IsNotFound},
	//	{scope.MockID(1), nil},
	//	{scope.MockID(100), errors.IsNotFound},
	//	{mockIDCode{1, "dede"}, nil},
	//	{mockIDCode{2, "czfr"}, errors.IsNotFound},
	//	{mockIDCode{2, ""}, nil},
	//}
	//
	//for i, test := range tests {
	//	s, err := storeService.Store(test.have)
	//	if test.wantErrBhf == nil {
	//		assert.NoError(t, err, "Index %d", i)
	//		assert.NotNil(t, s, "Index %d", i)
	//		//			assert.NotEmpty(t, s.Data.Code.String, "%#v", s.Data)
	//	} else {
	//		assert.True(t, test.wantErrBhf(err), "Index %d Error: %s", i, err)
	//		assert.Nil(t, s, "Index %d", i)
	//	}
	//}
	//assert.False(t, storeService.IsCacheEmpty())
	//storeService.ClearCache()
	//assert.True(t, storeService.IsCacheEmpty())
}

/*
	MOCKS
*/

//type mockIDCode struct {
//	id   int64
//	code string
//}
//
//func (ic mockIDCode) StoreID() int64 {
//	return ic.id
//}
//func (ic mockIDCode) StoreCode() string {
//	return ic.code
//}
//func (ic mockIDCode) WebsiteID() int64 {
//	return ic.id
//}
//func (ic mockIDCode) WebsiteCode() string {
//	return ic.code
//}
