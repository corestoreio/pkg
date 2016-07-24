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

	"sync"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ store.CodeToIDMapper = (*store.Service)(nil)
var _ store.AvailabilityChecker = (*store.Service)(nil)

var serviceStoreSimpleTest = store.MustNewService(
	cfgmock.NewService(),
	store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
	store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
	store.WithTableStores(&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
)

func TestNewServiceStore_QueryInvalidStore(t *testing.T) {

	assert.False(t, serviceStoreSimpleTest.IsCacheEmpty())

	s, err := serviceStoreSimpleTest.Store(-1)
	assert.True(t, errors.IsNotFound(err), "%+v", err)
	err = s.Validate()
	assert.True(t, errors.IsNotValid(err), "%+v", err)
	assert.EqualValues(t, "", s.Code())

	assert.False(t, serviceStoreSimpleTest.IsCacheEmpty())
	serviceStoreSimpleTest.ClearCache()
	assert.True(t, serviceStoreSimpleTest.IsCacheEmpty())
}

func TestMustNewService_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotFound(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	_ = store.MustNewService(cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 0, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
	)
}

func TestMustNewService_NoPanic(t *testing.T) {
	tests := []struct {
		have       int64
		wantErrBhf errors.BehaviourFunc
	}{
		{-1, errors.IsNotFound},
		{4444, errors.IsNotFound},
		{0, errors.IsNotFound},
	}
	serviceEmpty := store.MustNewService(cfgmock.NewService())
	for i, test := range tests {
		s, err := serviceEmpty.Store(test.have)
		assert.Error(t, s.Validate(), "Index %d", i)
		assert.True(t, test.wantErrBhf(err), "Index %d => %s", i, err)
	}
}

func TestNewService_DefaultStoreView_OK(t *testing.T) {

	serviceDefaultStore := store.MustNewService(
		cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
		store.WithTableStores(&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
	)

	// call it twice to test internal caching
	s, err := serviceDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.NoError(t, err)
	assert.Exactly(t, "de", s.Code())

	s, err = serviceDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.NoError(t, err)
	assert.Exactly(t, "de", s.Code())
	assert.False(t, serviceDefaultStore.IsCacheEmpty())
	serviceDefaultStore.ClearCache()
	assert.True(t, serviceDefaultStore.IsCacheEmpty())
}

func TestNewService_DefaultStoreView_NOK(t *testing.T) {

	serviceDefaultStore := store.MustNewService(
		cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
		store.WithTableStores(&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
	)

	// call it twice to test internal caching
	s, err := serviceDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.True(t, errors.IsNotFound(err), "%+v", err)
	assert.Exactly(t, "", s.Code())

}
func TestNewServiceStores(t *testing.T) {

	serviceStores := store.MustNewService(
		cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
		store.WithTableStores(
			&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			&store.TableStore{StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
		),
	)

	// call it twice to test internal caching
	ss := serviceStores.Stores()
	assert.NotNil(t, ss)
	assert.Equal(t, "at", ss[1].Data.Code.String)

	ss = serviceStores.Stores()
	assert.NotNil(t, ss)
	assert.NotEmpty(t, ss[2].Data.Code.String)

	assert.False(t, serviceStores.IsCacheEmpty())
	serviceStores.ClearCache()
	assert.True(t, serviceStores.IsCacheEmpty())
}

func TestMustNewServiceStores(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotFound(err), "Error: %+v", err)
			//t.Logf("%+v", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	_ = store.MustNewService(cfgmock.NewService(),
		store.WithTableStores(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 10, WebsiteID: 21, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
	)
}

func TestNewServiceGroup(t *testing.T) {

	serviceGroupSimpleTest := store.MustNewService(cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
		store.WithTableStores(&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
	)

	tests := []struct {
		m               *store.Service
		have            int64
		wantErr         errors.BehaviourFunc
		wantGroupName   string
		wantWebsiteCode string
	}{
		{serviceGroupSimpleTest, 20, errors.IsNotFound, "DACH Group", "euro"},
		{serviceGroupSimpleTest, 1, nil, "DACH Group", "euro"},
	}

	for i, test := range tests {
		g, err := test.m.Group(test.have)
		if test.wantErr != nil {
			assert.NoError(t, g.Validate(), "Index %d", i)
			assert.True(t, test.wantErr(err), "Index %d\n%+v", i, err)
			continue
		}
		assert.NotNil(t, g, "Index %d", i)
		assert.NoError(t, err, "Index %d\n%#v", i, err)
		assert.Exactly(t, test.wantGroupName, g.Name(), "Index %d", i)
		assert.Exactly(t, test.wantWebsiteCode, g.Website.Code(), "Index %d", i)

	}
	assert.False(t, serviceGroupSimpleTest.IsCacheEmpty())
	serviceGroupSimpleTest.ClearCache()
	assert.True(t, serviceGroupSimpleTest.IsCacheEmpty())
}

func TestNewServiceGroups(t *testing.T) {

	serviceGroups := store.MustNewService(cfgmock.NewService(),
		store.WithTableStores(&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
	)
	const iterations = 10
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			ss := serviceGroups.Groups()
			assert.NotNil(t, ss)
			assert.Len(t, ss, 1)
		}(&wg)
	}
	wg.Wait()

	assert.False(t, serviceGroups.IsCacheEmpty())
	serviceGroups.ClearCache()
	assert.True(t, serviceGroups.IsCacheEmpty())
}

func TestNewServiceWebsite(t *testing.T) {

	serviceWebsite := store.MustNewService(cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
		store.WithTableStores(&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
	)

	tests := []struct {
		m               *store.Service
		have            int64
		wantErr         errors.BehaviourFunc
		wantWebsiteCode string
	}{
		{serviceWebsite, 1, nil, "euro"},
		{serviceWebsite, 0, errors.IsNotFound, ""},
		{serviceWebsite, 0, errors.IsNotFound, ""},
	}

	for i, test := range tests {
		haveW, haveErr := test.m.Website(test.have)
		if test.wantErr != nil {
			assert.True(t, test.wantErr(haveErr), "Index %d\n%+v", i, haveErr)
			assert.True(t, errors.IsNotValid(haveW.Validate()), "Index %d", i)
		} else {
			assert.NoError(t, haveErr, "Index %d\n%+v", i, haveErr)
			assert.NotNil(t, haveW, "Index %d", i)
			assert.Exactly(t, test.wantWebsiteCode, haveW.Code(), "Index %d", i)
		}
	}
	assert.False(t, serviceWebsite.IsCacheEmpty())
	serviceWebsite.ClearCache()
	assert.True(t, serviceWebsite.IsCacheEmpty())

}

//func TestNewServiceWebsites(t *testing.T) {
//
//	serviceWebsites := store.MustNewService(cfgmock.NewService(),
//		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)}),
//		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
//		store.WithTableStores(&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
//	)
//
//	tests := []struct {
//		m       *store.Service
//		wantErr error
//		wantNil bool
//	}{
//		{serviceWebsites, nil, false},
//		{serviceWebsites, nil, false},
//		//{storemock.MustNewService(0, func(ms *storemock.Storage) {
//		//	ms.MockWebsiteSlice = func() (store.WebsiteSlice, error) {
//		//		return nil, nil
//		//	}
//		//	ms.MockStore = func() (*store.Store, error) {
//		//		return store.NewStore(
//		//			cfgmock.NewService(),
//		//			&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
//		//			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
//		//			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
//		//		)
//		//	}
//		//}), nil, true},
//	}
//
//	for _, test := range tests {
//		haveWS, haveErr := test.m.Websites()
//		if test.wantErr != nil {
//			assert.Error(t, haveErr, "%#v", test)
//			assert.Nil(t, haveWS, "%#v", test)
//		} else {
//			assert.NoError(t, haveErr, "%#v", test)
//			if test.wantNil {
//				assert.Nil(t, haveWS, "%#v", test)
//			} else {
//				assert.NotNil(t, haveWS, "%#v", test)
//			}
//		}
//	}
//
//	assert.False(t, serviceWebsites.IsCacheEmpty())
//	serviceWebsites.ClearCache()
//	assert.True(t, serviceWebsites.IsCacheEmpty())
//}

//type testNewServiceRequestedStore struct {
//	haveStoreID   int64
//	wantStoreCode string
//	wantErrBhf    errors.BehaviourFunc
//}
//
//func runTestsRequestedStore(t *testing.T, sm *store.Service, tests []testNewServiceRequestedStore) {
//	for i, test := range tests {
//		haveStore, haveErr := sm.RequestedStore(test.haveStoreID)
//		if test.wantErrBhf != nil {
//			assert.Nil(t, haveStore, "Index: %d: %#v", i, test)
//			assert.True(t, test.wantErrBhf(haveErr), "Index %d Error: %s", i, haveErr)
//		} else {
//			assert.NotNil(t, haveStore, "Index %d", i)
//			assert.NoError(t, haveErr, "Index %d => %#v", i, test)
//			assert.EqualValues(t, test.wantStoreCode, haveStore.Data.Code.String, "Index %d", i)
//		}
//	}
//	sm.ClearCache(true)
//}

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
