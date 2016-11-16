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
	"sync"
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/null"
	"github.com/stretchr/testify/assert"
)

var _ store.Finder = (*store.Service)(nil)

var serviceStoreSimpleTest = store.MustNewService(
	cfgmock.NewService(),
	store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
	store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
	store.WithTableStores(&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
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
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
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
	serviceEmpty := store.MustNewService(cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
	)
	for i, test := range tests {
		s, err := serviceEmpty.Store(test.have)
		assert.Error(t, s.Validate(), "Index %d", i)
		assert.True(t, test.wantErrBhf(err), "Index %d => %s", i, err)
	}
}

func TestMustNewService_DefaultWebsiteCheck(t *testing.T) {

	s, err := store.NewService(cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), IsDefault: null.BoolFrom(true)}),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 12, Code: null.StringFrom("euro2"), IsDefault: null.BoolFrom(true)}),
	)
	assert.Nil(t, s)
	assert.True(t, errors.IsNotValid(err), "%+v", err)
}

func TestNewService_DefaultStoreView_OK(t *testing.T) {

	serviceDefaultStore := store.MustNewService(
		cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
		store.WithTableStores(&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
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
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
		store.WithTableStores(&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
	)

	// call it twice to test internal caching
	s, err := serviceDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.True(t, errors.IsNotFound(err), "%+v", err)
	assert.Exactly(t, "", s.Code())

}
func TestNewService_Stores(t *testing.T) {

	serviceStores := store.MustNewService(
		cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
		store.WithTableStores(
			&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 2, Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			&store.TableStore{StoreID: 3, Code: null.StringFrom("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
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

func TestMustNewService_Stores_Panic(t *testing.T) {
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
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 10, WebsiteID: 21, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
	)
}

func TestNewService_Group(t *testing.T) {

	serviceGroupSimpleTest := store.MustNewService(cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
		store.WithTableStores(&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
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

func TestNewService_Groups(t *testing.T) {

	serviceGroups := store.MustNewService(cfgmock.NewService(),
		store.WithTableStores(&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
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

func TestNewService_Website(t *testing.T) {

	serviceWebsite := store.MustNewService(cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
		store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
		store.WithTableStores(&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
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

func TestNewService_Websites(t *testing.T) {
	srv := store.MustNewService(cfgmock.NewService(),
		store.WithTableWebsites(
			&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("European Union"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
			&store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("uk"), Name: null.StringFrom("Britain (without Scotland)"), SortOrder: 0, DefaultGroupID: 2},
		),
	)
	assert.Exactly(t, []int64{1, 2}, srv.Websites().IDs())
	assert.Exactly(t, []string{"euro", "uk"}, srv.Websites().Codes())
}

func TestService_IsAllowedStoreID(t *testing.T) {
	eurSrv := storemock.NewEurozzyService(cfgmock.NewService())
	tests := []struct {
		srv           *store.Service
		runMode       scope.TypeID
		storeID       int64
		wantIsAllowed bool
		wantCode      string
		wantErrBhf    errors.BehaviourFunc
	}{
		{eurSrv, 0, 1, true, "de", nil},                                     // fall back to default website -> default group -> default store
		{eurSrv, 0, 2, true, "at", nil},                                     // fall back to default website -> default group -> default store
		{eurSrv, 0, 5, false, "", nil},                                      // fall back to default website -> default group -> default store Australia not allowed
		{eurSrv, 0, 0, false, "", nil},                                      // fall back to default website -> default group -> default store admin not allowed
		{eurSrv, scope.MakeTypeID(scope.Website, 0), 0, true, "admin", nil}, // admin scope or single website scope
		{eurSrv, scope.MakeTypeID(scope.Website, 0), 2, false, "", nil},     // admin scope or single website scope
		{eurSrv, scope.MakeTypeID(scope.Website, 1), 1, true, "de", nil},    // euro scope, not included ch, because not active, and UK, different group
		{eurSrv, scope.MakeTypeID(scope.Website, 1), 2, true, "at", nil},    // euro scope, not included ch, because not active, and UK, different group
		{eurSrv, scope.MakeTypeID(scope.Website, 1), 3, false, "", nil},     // euro scope, not included ch
		{eurSrv, scope.MakeTypeID(scope.Website, 1), 4, true, "uk", nil},    // euro scope, uk allowed
		{eurSrv, scope.MakeTypeID(scope.Website, 2), 5, true, "au", nil},    // oz scope
		{eurSrv, scope.MakeTypeID(scope.Website, 2), 6, true, "nz", nil},    // oz scope
		{eurSrv, scope.MakeTypeID(scope.Website, 2), 1, false, "", nil},     // oz scope
		{eurSrv, scope.MakeTypeID(scope.Website, 9999), 1, false, "", nil},
		{eurSrv, scope.MakeTypeID(scope.Website, 1), 9999, false, "", nil},
		{eurSrv, scope.MakeTypeID(scope.Group, 0), 0, true, "admin", nil}, // admin scope
		{eurSrv, scope.MakeTypeID(scope.Group, 1), 1, true, "de", nil},    // dach scope
		{eurSrv, scope.MakeTypeID(scope.Group, 1), 2, true, "at", nil},    // dach scope
		{eurSrv, scope.MakeTypeID(scope.Group, 2), 4, true, "uk", nil},    // uk scope
		{eurSrv, scope.MakeTypeID(scope.Group, 2), 5, false, "", nil},     // uk scope
		{eurSrv, scope.MakeTypeID(scope.Group, 9999), 4, false, "", nil},  // uk scope
		{eurSrv, scope.MakeTypeID(scope.Store, 0), 5, true, "au", nil},
		{eurSrv, scope.MakeTypeID(scope.Store, 0), 1, true, "de", nil},
		{eurSrv, scope.MakeTypeID(scope.Store, 0), 3, false, "", nil},
		{eurSrv, scope.MakeTypeID(scope.Store, 1), 4, true, "uk", nil},
		{eurSrv, scope.MakeTypeID(scope.Store, 9999), 4, true, "uk", nil},
		{eurSrv, scope.MakeTypeID(124, 1), 4, false, "", nil},
		{eurSrv, scope.MakeTypeID(124, 0), 4, false, "", nil},
		{store.MustNewService(cfgmock.NewService(),
			store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
			store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
			store.WithTableStores(&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
		), 0, 2, false, "", errors.IsNotFound},
	}
	for i, test := range tests {
		haveIsAllowed, haveCode, haveErr := test.srv.IsAllowedStoreID(test.runMode, test.storeID)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "(%d) %+v", i, haveErr)
			assert.False(t, haveIsAllowed, "Index %d", i)
			assert.Empty(t, haveCode, "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "(%d) %+v", i, haveErr)
		assert.Exactly(t, test.wantIsAllowed, haveIsAllowed, "Index %d", i)
		assert.Exactly(t, test.wantCode, haveCode, "Index %d", i)
	}
}

func TestService_DefaultStoreID(t *testing.T) {
	eurSrv := storemock.NewEurozzyService(cfgmock.NewService())
	tests := []struct {
		srv           *store.Service
		runMode       scope.TypeID
		wantStoreID   int64
		wantWebsiteID int64
		wantErrBhf    errors.BehaviourFunc
	}{
		{eurSrv, 0, 2, 1, nil},                                  // fall back to default website -> default group -> default store
		{eurSrv, scope.MakeTypeID(scope.Website, 0), 0, 0, nil}, // admin scope
		{eurSrv, scope.MakeTypeID(scope.Website, 1), 2, 1, nil}, // euro scope, not included ch, because not active, and UK, different group
		{eurSrv, scope.MakeTypeID(scope.Website, 2), 5, 2, nil}, // oz scope
		{eurSrv, scope.MakeTypeID(scope.Website, 9999), 0, 0, errors.IsNotFound},
		{store.MustNewService(cfgmock.NewService(), // default store not active
			store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
			store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
			store.WithTableStores(&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: false}),
		), scope.MakeTypeID(scope.Website, 1), 0, 0, errors.IsNotValid},

		{eurSrv, scope.MakeTypeID(scope.Group, 0), 0, 0, nil}, // admin scope
		{eurSrv, scope.MakeTypeID(scope.Group, 1), 2, 1, nil}, // dach scope
		{eurSrv, scope.MakeTypeID(scope.Group, 2), 4, 1, nil}, // uk scope
		{eurSrv, scope.MakeTypeID(scope.Group, 3), 5, 2, nil}, // au scope
		{eurSrv, scope.MakeTypeID(scope.Group, 9999), 0, 0, errors.IsNotFound},
		{store.MustNewService(cfgmock.NewService(), // default store not active
			store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
			store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
			store.WithTableStores(&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: false}),
		), scope.MakeTypeID(scope.Group, 1), 0, 0, errors.IsNotValid},
		{store.MustNewService(cfgmock.NewService(), // default store not found
			store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
			store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
		), scope.MakeTypeID(scope.Group, 1), 0, 0, errors.IsNotFound},

		{eurSrv, scope.MakeTypeID(scope.Store, 0), 0, 0, nil},
		{eurSrv, scope.MakeTypeID(scope.Store, 1), 1, 1, nil},
		{eurSrv, scope.MakeTypeID(scope.Store, 9999), 0, 0, errors.IsNotFound},
		{eurSrv, scope.MakeTypeID(scope.Store, 3), 0, 0, errors.IsNotValid}, // ch store is not active

		{store.MustNewService(cfgmock.NewService(),
			store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
			store.WithTableGroups(&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
			store.WithTableStores(&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
		), 0, 0, 0, errors.IsNotFound},
	}
	for i, test := range tests {
		haveStoreID, haveWebsiteID, haveErr := test.srv.DefaultStoreID(test.runMode)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "(%d) %+v", i, haveErr)
			assert.Exactly(t, test.wantStoreID, haveStoreID, "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "(%d) %+v", i, haveErr)
		assert.Exactly(t, test.wantStoreID, haveStoreID, "Index %d", i)
		assert.Exactly(t, test.wantWebsiteID, haveWebsiteID, "Index %d", i)
	}
}

func TestService_StoreIDbyCode(t *testing.T) {
	eurSrv := storemock.NewEurozzyService(cfgmock.NewService())
	tests := []struct {
		srv           *store.Service
		runMode       scope.TypeID
		code          string
		wantStoreID   int64
		wantWebsiteID int64
		wantErrBhf    errors.BehaviourFunc
	}{
		{eurSrv, 0, "", 2, 1, nil},
		{eurSrv, scope.DefaultTypeID, "x", 0, 0, errors.IsNotFound},
		{eurSrv, scope.DefaultTypeID, "uk", 0, 0, errors.IsNotFound},
		{eurSrv, scope.Website.Pack(0), "admin", 0, 0, nil},
		{eurSrv, scope.Website.Pack(1), "de", 1, 1, nil},
		{eurSrv, scope.Website.Pack(2), "nz", 6, 2, nil},
		{eurSrv, scope.Website.Pack(3), "uk", 0, 0, errors.IsNotFound},
		{eurSrv, scope.Absent.Pack(0), "uk", 0, 0, errors.IsNotFound},
		{eurSrv, scope.Absent.Pack(0), "at", 2, 1, nil},
		{eurSrv, scope.Group.Pack(2), "uk", 4, 1, nil},
		{eurSrv, scope.Group.Pack(99), "uk", 0, 0, errors.IsNotFound},
		{eurSrv, scope.Store.Pack(0), "admin", 0, 0, nil},
		{eurSrv, scope.Store.Pack(0), "au", 5, 2, nil},
		{eurSrv, scope.Store.Pack(0), "xx", 0, 0, errors.IsNotFound},
	}
	for i, test := range tests {
		haveStoreID, haveWebsiteID, haveErr := test.srv.StoreIDbyCode(test.runMode, test.code)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "(%d) %+v", i, haveErr)
			assert.Exactly(t, test.wantStoreID, haveStoreID, "Index %d", i)
			assert.Exactly(t, test.wantWebsiteID, haveWebsiteID, "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "(%d) %+v", i, haveErr)
		assert.Exactly(t, test.wantStoreID, haveStoreID, "Index %d", i)
		assert.Exactly(t, test.wantWebsiteID, haveWebsiteID, "Index %d", i)
	}
}

func TestService_AllowedStores(t *testing.T) {
	eurSrv := storemock.NewEurozzyService(cfgmock.NewService())
	tests := []struct {
		srv          *store.Service
		runMode      scope.TypeID
		wantStoreIDs []int64
		wantErrBhf   errors.BehaviourFunc
	}{
		{eurSrv, 0, []int64{1, 2}, nil},
		{eurSrv, scope.DefaultTypeID, []int64{1, 2}, nil},
		{eurSrv, scope.Website.Pack(0), []int64{0}, nil},
		{eurSrv, scope.Website.Pack(1), []int64{1, 4, 2}, nil},
		{eurSrv, scope.Website.Pack(2), []int64{5, 6}, nil},
		{eurSrv, scope.Website.Pack(3), nil, nil},
		{eurSrv, scope.Absent.Pack(0), []int64{1, 2}, nil},
		{eurSrv, scope.Group.Pack(2), []int64{4}, nil},
		{eurSrv, scope.Group.Pack(99), nil, nil},
		{eurSrv, scope.Store.Pack(0), []int64{0, 5, 1, 4, 2, 6}, nil},
	}
	for i, test := range tests {
		haveStores, haveErr := test.srv.AllowedStores(test.runMode)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "(%d) %+v", i, haveErr)
			assert.Exactly(t, test.wantStoreIDs, haveStores.IDs(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "(%d) %+v", i, haveErr)
		assert.Exactly(t, test.wantStoreIDs, haveStores.IDs(), "Index %d", i)
	}
}

func TestService_HasSingleStore(t *testing.T) {
	s := store.MustNewService(cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
	)
	s1 := store.MustNewService(cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
	)
	s1.SingleStoreModeEnabled = false

	s2 := storemock.NewEurozzyService(cfgmock.NewService())

	const iterations = 10
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			assert.True(t, s.HasSingleStore())   // no stores so true
			assert.False(t, s1.HasSingleStore()) // no stores but globally disabled so false
			assert.False(t, s2.HasSingleStore()) // lots of stores so false
		}(&wg)
	}
	wg.Wait()
}

func TestService_IsSingleStoreMode(t *testing.T) {

	const xPath = `general/single_store_mode/enabled`

	s := store.MustNewService(cfgmock.NewService(),
		store.WithTableWebsites(&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
	)

	// no stores and backend not set so true
	sCfg := cfgmock.NewService().NewScoped(0, 0)
	b, err := s.IsSingleStoreMode(sCfg)
	assert.NoError(t, err, "%+v", err)
	assert.True(t, b)

	// no stores and backend set but configured with false
	s.ClearCache()
	sCfg = cfgmock.NewService(cfgmock.PathValue{
		cfgpath.MustNewByParts(xPath).BindStore(2).String(): 0,
	}).NewScoped(1, 2)
	s.BackendSingleStore = cfgmodel.NewBool(xPath, cfgmodel.WithScopeStore())
	b, err = s.IsSingleStoreMode(sCfg)
	assert.NoError(t, err, "%+v", err)
	assert.False(t, b)

	// no stores and backend set but returns an error
	s.ClearCache()
	tErr := errors.NewNotImplementedf("Ups")
	s.BackendSingleStore = cfgmodel.NewBool(xPath)
	s.BackendSingleStore.LastError = tErr
	b, err = s.IsSingleStoreMode(config.Scoped{})
	assert.True(t, errors.IsNotImplemented(tErr), "%+v", tErr)
	assert.False(t, b)

	s2 := storemock.NewEurozzyService(cfgmock.NewService())
	s2.BackendSingleStore = cfgmodel.NewBool(xPath) // returns false always no error
	assert.False(t, s2.HasSingleStore())

	b, err = s2.IsSingleStoreMode(sCfg)
	assert.NoError(t, err, "%+v", err)
	assert.False(t, b)

	s2.ClearCache()
	s2.BackendSingleStore = cfgmodel.NewBool(xPath, cfgmodel.WithField(&element.Field{ID: cfgpath.NewRoute(`enabled`), Default: `1`})) // returns true
	b, err = s2.IsSingleStoreMode(sCfg)
	assert.NoError(t, err, "%+v", err)
	assert.True(t, b)

	// call it twice to test cache
	b, err = s2.IsSingleStoreMode(sCfg)
	assert.NoError(t, err, "%+v", err)
	assert.True(t, b)
}

func TestService_LoadFromDB_OK(t *testing.T) {

	dbrCon, dbMock := cstesting.MockDB(t)

	dbMock.ExpectQuery("SELECT (.+) FROM `store`(.+) ORDER BY CASE WHEN(.+)").WillReturnRows(
		cstesting.MustMockRows(cstesting.WithFile("testdata", "core_store_view.csv")),
	)
	dbMock.ExpectQuery("SELECT (.+) FROM `store_website`(.+) ORDER BY(.+)").WillReturnRows(
		cstesting.MustMockRows(cstesting.WithFile("testdata", "core_website_view.csv")),
	)
	dbMock.ExpectQuery("SELECT (.+) FROM `store_group`(.+) ORDER BY main_table(.+)").WillReturnRows(
		cstesting.MustMockRows(cstesting.WithFile("testdata", "core_store_group_view.csv")),
	)
	dbMock.MatchExpectationsInOrder(false) // we're using goroutines!

	srv := store.MustNewService(cfgmock.NewService())
	if err := srv.LoadFromDB(dbrCon.NewSession()); err != nil {
		t.Fatalf("%+v", err)
	}

	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Len(t, srv.Websites(), 9)
	assert.Len(t, srv.Groups(), 9)
	assert.Len(t, srv.Stores(), 16)

	tree, err := json.Marshal(srv.Websites().Tree())
	if err != nil {
		t.Fatal(err)
	}

	assert.Exactly(t,
		`{"scope":"Default","id":0,"scopes":[{"scope":"Website","id":0,"scopes":[{"scope":"Group","id":0,"scopes":[{"scope":"Store","id":0}]}]},{"scope":"Website","id":2,"scopes":[{"scope":"Group","id":2,"scopes":[{"scope":"Store","id":2},{"scope":"Store","id":5}]}]},{"scope":"Website","id":3,"scopes":[{"scope":"Group","id":3,"scopes":[{"scope":"Store","id":6},{"scope":"Store","id":7},{"scope":"Store","id":8},{"scope":"Store","id":9}]}]},{"scope":"Website","id":4,"scopes":[{"scope":"Group","id":4,"scopes":[{"scope":"Store","id":10},{"scope":"Store","id":11}]}]},{"scope":"Website","id":5,"scopes":[{"scope":"Group","id":5,"scopes":[{"scope":"Store","id":12}]}]},{"scope":"Website","id":6,"scopes":[{"scope":"Group","id":6,"scopes":[{"scope":"Store","id":13},{"scope":"Store","id":14}]}]},{"scope":"Website","id":7,"scopes":[{"scope":"Group","id":7,"scopes":[{"scope":"Store","id":15},{"scope":"Store","id":16}]}]},{"scope":"Website","id":8,"scopes":[{"scope":"Group","id":8,"scopes":[{"scope":"Store","id":17}]}]},{"scope":"Website","id":9,"scopes":[{"scope":"Group","id":9,"scopes":[{"scope":"Store","id":18}]}]}]}`,
		string(tree))
}

func TestService_LoadFromDB_NOK_Store(t *testing.T) {

	dbrCon, dbMock := cstesting.MockDB(t)

	wsErr := errors.NewAlreadyClosedf("DB Already closed")
	dbMock.ExpectQuery("SELECT (.+) FROM `store`(.+) ORDER BY CASE WHEN(.+)").WillReturnError(wsErr)
	dbMock.ExpectQuery("SELECT (.+) FROM `store_website`(.+) ORDER BY(.+)").WillReturnRows(
		cstesting.MustMockRows(cstesting.WithFile("testdata", "core_website_view.csv")),
	)
	dbMock.ExpectQuery("SELECT (.+) FROM `store_group`(.+) ORDER BY main_table(.+)").WillReturnRows(
		cstesting.MustMockRows(cstesting.WithFile("testdata", "core_store_group_view.csv")),
	)
	dbMock.MatchExpectationsInOrder(false) // we're using goroutines!

	srv := store.MustNewService(cfgmock.NewService())
	err := srv.LoadFromDB(dbrCon.NewSession())
	assert.True(t, errors.IsAlreadyClosed(err))

	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Len(t, srv.Websites(), 0)
	assert.Len(t, srv.Groups(), 0)
	assert.Len(t, srv.Stores(), 0)

}

func TestService_LoadFromDB_NOK_All(t *testing.T) {

	dbrCon, dbMock := cstesting.MockDB(t)

	wsErr1 := errors.NewAlreadyClosedf("DB Already closed")
	wsErr2 := errors.NewNotImplementedf("DB is NoSQL")
	wsErr3 := errors.NewEmptyf("DB empty")
	dbMock.ExpectQuery("SELECT (.+) FROM `store`(.+) ORDER BY CASE WHEN(.+)").WillReturnError(wsErr1)
	dbMock.ExpectQuery("SELECT (.+) FROM `store_website`(.+) ORDER BY(.+)").WillReturnError(wsErr2)
	dbMock.ExpectQuery("SELECT (.+) FROM `store_group`(.+) ORDER BY main_table(.+)").WillReturnError(wsErr3)
	dbMock.MatchExpectationsInOrder(false) // we're using goroutines!

	srv := store.MustNewService(cfgmock.NewService())
	err := srv.LoadFromDB(dbrCon.NewSession())
	assert.True(t, errors.IsAlreadyClosed(err) || errors.IsNotImplemented(err) || errors.IsEmpty(err), "%+v", err)

	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Len(t, srv.Websites(), 0)
	assert.Len(t, srv.Groups(), 0)
	assert.Len(t, srv.Stores(), 0)

}
