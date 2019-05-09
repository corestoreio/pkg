// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/store"
	storemock "github.com/corestoreio/pkg/store/mock"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
)

var _ store.Finder = (*store.Service)(nil)

var serviceStoreSimpleTest = store.MustNewService(
	store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: "euro", Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: true}),
	store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
	store.WithStores(&store.Store{StoreID: 1, Code: "de", WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
)

func TestNewServiceStore_QueryInvalidStore(t *testing.T) {

	assert.True(t, serviceStoreSimpleTest.IsCacheEmpty())

	s, err := serviceStoreSimpleTest.Store(10000)
	assert.ErrorIsKind(t, errors.NotFound, err)
	err = s.Validate()
	assert.ErrorIsKind(t, errors.NotValid, err)

	assert.True(t, serviceStoreSimpleTest.IsCacheEmpty())
	serviceStoreSimpleTest.ClearCache()
	assert.True(t, serviceStoreSimpleTest.IsCacheEmpty())

}

func TestMustNewService_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.ErrorIsKind(t, errors.NotFound, err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	_ = store.MustNewService(
		store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: "euro", Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: true}),
		store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 0, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
	)
}

func TestNewService_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    []store.Option
		wantErr errors.Kind
	}{
		{"All valid",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true, Code: "dach"}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 1, GroupID: 1, Code: "de-de"}),
			},
			errors.NoKind},

		{"Website DefaultGroupID not found",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 0, DefaultStoreID: 2}),
			},
			errors.NotValid},
		{"No Default Website found",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: false}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 0, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 1, GroupID: 1}),
			},
			errors.NotValid},
		{"too many Default Websites",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true}),
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 2, DefaultGroupID: 1, IsDefault: true}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 0, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 1, GroupID: 1}),
			},
			errors.NotValid},
		{"Group WebsiteID not found",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 10000, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 1, GroupID: 1}),
			},
			errors.NotValid},
		{"Group DefaultStoreID not found",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 100000}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 1, GroupID: 1}),
			},
			errors.NotValid},
		{"Store WebsiteID not found",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 10000, GroupID: 1}),
			},
			errors.NotValid},
		{"Store GroupID not found",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 1, GroupID: 10000}),
			},
			errors.NotValid},
		{"Missing Store Code",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true, Code: "dach"}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 1, GroupID: 1}),
			},
			errors.NotValid},
		{"Store.StoreWebsite.WebsiteID not found in Store",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true, Code: "dach"}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 2220, GroupID: 1, Code: "de-de",
					StoreWebsite: &store.StoreWebsite{
						WebsiteID: 2221, DefaultGroupID: 1, IsDefault: true, Code: "dach",
					},
				}),
			},
			errors.NotValid},
		{"Store.StoreGroup.WebsiteID not found in Store",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true, Code: "dach"}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 2220, GroupID: 1, Code: "de-de",
					StoreGroup: &store.StoreGroup{GroupID: 1, WebsiteID: 2221, DefaultStoreID: 1},
				}),
			},
			errors.NotValid},
		{"Store.StoreGroup.GroupID not found in Store",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true, Code: "dach"}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 2220, GroupID: 1, Code: "de-de",
					StoreGroup: &store.StoreGroup{GroupID: 1, WebsiteID: 2221, DefaultStoreID: 1},
				}),
			},
			errors.NotValid},
		// TODO add Website.Validate and Group.Validate
	}

	s := store.MustNewService()
	for _, test := range tests {
		if err := s.Options(test.opts...); test.wantErr > 0 {
			assert.ErrorIsKind(t, test.wantErr, err, test.name)
		} else {
			assert.NoError(t, err)
		}
		s.ClearCache()
	}
}

func TestNewService_DefaultStoreView(t *testing.T) {

	srv := store.MustNewService(
		store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: "euro", Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: true}),
		store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
		store.WithStores(&store.Store{StoreID: 1, Code: "de", WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
	)

	bgwork.Wait(20, func(index int) {
		s, err := srv.DefaultStoreView()
		if err != nil {
			t.Fatal(err)
		}
		if s.Code != "de" {
			t.Fatalf("Expecting store code `de`, got %q", s.Code)
		}
	})

	srv.ClearCache()
	assert.True(t, srv.IsCacheEmpty())
}

func TestNewService_WebsiteGroupStore(t *testing.T) {
	stores := store.Stores{
		Data: []*store.Store{
			{StoreID: 1, Code: "de", WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			{StoreID: 2, Code: "at", WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			{StoreID: 3, Code: "ch", WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
		},
	}
	groups := store.StoreGroups{
		Data: []*store.StoreGroup{
			{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		},
	}
	websites := store.StoreWebsites{
		Data: []*store.StoreWebsite{
			{WebsiteID: 1, Code: "euro", Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: true},
		},
	}

	srv := store.MustNewService(
		store.WithWebsites(websites.Data...),
		store.WithGroups(groups.Data...),
		store.WithStores(stores.Data...),
	)

	t.Run("Store", func(t *testing.T) {
		assert.Exactly(t, stores, srv.Stores())

		bgwork.Wait(20, func(index int) {
			st, err := srv.Store(3)
			assert.NoError(t, err)
			assert.Exactly(t, stores.Data[2], st)

			st, err = srv.Store(4)
			assert.ErrorIsKind(t, errors.NotFound, err)
			assert.Nil(t, st)
		})

	})
	t.Run("Group", func(t *testing.T) {
		assert.Exactly(t, groups, srv.Groups())
		bgwork.Wait(20, func(index int) {
			st, err := srv.Group(1)
			assert.NoError(t, err)
			assert.Exactly(t, groups.Data[0], st)

			st, err = srv.Group(4)
			assert.ErrorIsKind(t, errors.NotFound, err)
			assert.Nil(t, st)
		})
	})
	t.Run("Website", func(t *testing.T) {
		assert.Exactly(t, websites, srv.Websites())
		bgwork.Wait(20, func(index int) {
			st, err := srv.Website(1)
			assert.NoError(t, err)
			assert.Exactly(t, websites.Data[0], st)

			st, err = srv.Website(4)
			assert.ErrorIsKind(t, errors.NotFound, err)
			assert.Nil(t, st)
		})
	})
}

func TestService_IsAllowedStoreID(t *testing.T) {

	tests := []struct {
		runMode       scope.TypeID
		storeID       uint32
		wantIsAllowed bool
		wantCode      string
	}{
		{0, 1, true, "de"}, // fall back to default website -> default group -> default store
		{0, 2, true, "at"}, // fall back to default website -> default group -> default store
		{0, 5, false, ""},  // fall back to default website -> default group -> default store Australia not allowed
		{0, 0, false, ""},  // fall back to default website -> default group -> default store admin not allowed
		{scope.MakeTypeID(scope.Website, 0), 0, true, "admin"}, // admin scope or single website scope
		{scope.MakeTypeID(scope.Website, 0), 2, false, ""},     // admin scope or single website scope
		{scope.MakeTypeID(scope.Website, 1), 1, true, "de"},    // euro scope, not included ch, because not active, and UK, different group
		{scope.MakeTypeID(scope.Website, 1), 2, true, "at"},    // euro scope, not included ch, because not active, and UK, different group
		{scope.MakeTypeID(scope.Website, 1), 3, false, ""},     // euro scope, not included ch
		{scope.MakeTypeID(scope.Website, 1), 4, true, "uk"},    // euro scope, uk allowed
		{scope.MakeTypeID(scope.Website, 2), 5, true, "au"},    // oz scope
		{scope.MakeTypeID(scope.Website, 2), 6, true, "nz"},    // oz scope
		{scope.MakeTypeID(scope.Website, 2), 1, false, ""},     // oz scope
		{scope.MakeTypeID(scope.Website, 9999), 1, false, ""},
		{scope.MakeTypeID(scope.Website, 1), 9999, false, ""},
		{scope.MakeTypeID(scope.Group, 0), 0, true, "admin"}, // admin scope
		{scope.MakeTypeID(scope.Group, 1), 1, true, "de"},    // dach scope
		{scope.MakeTypeID(scope.Group, 1), 2, true, "at"},    // dach scope
		{scope.MakeTypeID(scope.Group, 2), 4, true, "uk"},    // uk scope
		{scope.MakeTypeID(scope.Group, 2), 5, false, ""},     // uk scope
		{scope.MakeTypeID(scope.Group, 9999), 4, false, ""},  // uk scope
		{scope.MakeTypeID(scope.Store, 0), 5, true, "au"},
		{scope.MakeTypeID(scope.Store, 0), 1, true, "de"},
		{scope.MakeTypeID(scope.Store, 0), 3, false, ""},
		{scope.MakeTypeID(scope.Store, 1), 4, true, "uk"},
		{scope.MakeTypeID(scope.Store, 9999), 4, true, "uk"},
		{scope.MakeTypeID(124, 1), 4, false, ""},
		{scope.MakeTypeID(124, 0), 4, false, ""},
	}
	eoSrv := storemock.NewServiceEuroOZ()
	for i, test := range tests {
		haveIsAllowed, haveCode, haveErr := eoSrv.IsAllowedStoreID(test.runMode, test.storeID)
		assert.NoError(t, haveErr, "(Index %d)", i)
		assert.Exactly(t, test.wantIsAllowed, haveIsAllowed, "Index %d", i)
		assert.Exactly(t, test.wantCode, haveCode, "Index %d", i)
	}
}

// func TestService_DefaultStoreID(t *testing.T) {
// 	eurSrv := storemock.NewEurozzyService(cfgmock.NewService())
// 	tests := []struct {
// 		srv           *Service
// 		runMode       scope.TypeID
// 		wantStoreID   int64
// 		wantWebsiteID int64
// 		wantErrBhf    errors.BehaviourFunc
// 	}{
// 		{eurSrv, 0, 2, 1, nil},                                  // fall back to default website -> default group -> default store
// 		{eurSrv, scope.MakeTypeID(scope.Website, 0), 0, 0, nil}, // admin scope
// 		{eurSrv, scope.MakeTypeID(scope.Website, 1), 2, 1, nil}, // euro scope, not included ch, because not active, and UK, different group
// 		{eurSrv, scope.MakeTypeID(scope.Website, 2), 5, 2, nil}, // oz scope
// 		{eurSrv, scope.MakeTypeID(scope.Website, 9999), 0, 0, errors.IsNotFound},
// 		{MustNewService(cfgmock.NewService(), // default store not active
// 			store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: null.MakeString("euro"), Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)}),
// 			store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
// 			store.WithStores(&store.Store{StoreID: 1, Code: null.MakeString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: false}),
// 		), scope.MakeTypeID(scope.Website, 1), 0, 0, errors.IsNotValid},
//
// 		{eurSrv, scope.MakeTypeID(scope.Group, 0), 0, 0, nil}, // admin scope
// 		{eurSrv, scope.MakeTypeID(scope.Group, 1), 2, 1, nil}, // dach scope
// 		{eurSrv, scope.MakeTypeID(scope.Group, 2), 4, 1, nil}, // uk scope
// 		{eurSrv, scope.MakeTypeID(scope.Group, 3), 5, 2, nil}, // au scope
// 		{eurSrv, scope.MakeTypeID(scope.Group, 9999), 0, 0, errors.IsNotFound},
// 		{MustNewService(cfgmock.NewService(), // default store not active
// 			store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: null.MakeString("euro"), Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
// 			store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
// 			store.WithStores(&store.Store{StoreID: 1, Code: null.MakeString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: false}),
// 		), scope.MakeTypeID(scope.Group, 1), 0, 0, errors.IsNotValid},
// 		{MustNewService(cfgmock.NewService(), // default store not found
// 			store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: null.MakeString("euro"), Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
// 			store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
// 		), scope.MakeTypeID(scope.Group, 1), 0, 0, errors.IsNotFound},
//
// 		{eurSrv, scope.MakeTypeID(scope.Store, 0), 0, 0, nil},
// 		{eurSrv, scope.MakeTypeID(scope.Store, 1), 1, 1, nil},
// 		{eurSrv, scope.MakeTypeID(scope.Store, 9999), 0, 0, errors.IsNotFound},
// 		{eurSrv, scope.MakeTypeID(scope.Store, 3), 0, 0, errors.IsNotValid}, // ch store is not active
//
// 		{MustNewService(cfgmock.NewService(),
// 			store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: null.MakeString("euro"), Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
// 			store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2}),
// 			store.WithStores(&store.Store{StoreID: 1, Code: null.MakeString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
// 		), 0, 0, 0, errors.IsNotFound},
// 	}
// 	for i, test := range tests {
// 		haveStoreID, haveWebsiteID, haveErr := test.srv.DefaultStoreID(test.runMode)
// 		if test.wantErrBhf != nil {
// 			assert.True(t, test.wantErrBhf(haveErr), "(%d) %+v", i, haveErr)
// 			assert.Exactly(t, test.wantStoreID, haveStoreID, "Index %d", i)
// 			continue
// 		}
// 		assert.NoError(t, haveErr, "(%d) %+v", i, haveErr)
// 		assert.Exactly(t, test.wantStoreID, haveStoreID, "Index %d", i)
// 		assert.Exactly(t, test.wantWebsiteID, haveWebsiteID, "Index %d", i)
// 	}
// }
//
// func TestService_StoreIDbyCode(t *testing.T) {
// 	eurSrv := storemock.NewEurozzyService(cfgmock.NewService())
// 	tests := []struct {
// 		srv           *Service
// 		runMode       scope.TypeID
// 		code          string
// 		wantStoreID   int64
// 		wantWebsiteID int64
// 		wantErrBhf    errors.BehaviourFunc
// 	}{
// 		{eurSrv, 0, "", 2, 1, nil},
// 		{eurSrv, scope.DefaultTypeID, "x", 0, 0, errors.IsNotFound},
// 		{eurSrv, scope.DefaultTypeID, "uk", 0, 0, errors.IsNotFound},
// 		{eurSrv, scope.Website.WithID(0), "admin", 0, 0, nil},
// 		{eurSrv, scope.Website.WithID(1), "de", 1, 1, nil},
// 		{eurSrv, scope.Website.WithID(2), "nz", 6, 2, nil},
// 		{eurSrv, scope.Website.WithID(3), "uk", 0, 0, errors.IsNotFound},
// 		{eurSrv, scope.Absent.WithID(0), "uk", 0, 0, errors.IsNotFound},
// 		{eurSrv, scope.Absent.WithID(0), "at", 2, 1, nil},
// 		{eurSrv, scope.Group.WithID(2), "uk", 4, 1, nil},
// 		{eurSrv, scope.Group.WithID(99), "uk", 0, 0, errors.IsNotFound},
// 		{eurSrv, scope.Store.WithID(0), "admin", 0, 0, nil},
// 		{eurSrv, scope.Store.WithID(0), "au", 5, 2, nil},
// 		{eurSrv, scope.Store.WithID(0), "xx", 0, 0, errors.IsNotFound},
// 	}
// 	for i, test := range tests {
// 		haveStoreID, haveWebsiteID, haveErr := test.srv.StoreIDbyCode(test.runMode, test.code)
// 		if test.wantErrBhf != nil {
// 			assert.True(t, test.wantErrBhf(haveErr), "(%d) %+v", i, haveErr)
// 			assert.Exactly(t, test.wantStoreID, haveStoreID, "Index %d", i)
// 			assert.Exactly(t, test.wantWebsiteID, haveWebsiteID, "Index %d", i)
// 			continue
// 		}
// 		assert.NoError(t, haveErr, "(%d) %+v", i, haveErr)
// 		assert.Exactly(t, test.wantStoreID, haveStoreID, "Index %d", i)
// 		assert.Exactly(t, test.wantWebsiteID, haveWebsiteID, "Index %d", i)
// 	}
// }
//
// func TestService_AllowedStores(t *testing.T) {
// 	eurSrv := storemock.NewEurozzyService(cfgmock.NewService())
// 	tests := []struct {
// 		srv          *Service
// 		runMode      scope.TypeID
// 		wantStoreIDs []int64
// 		wantErrBhf   errors.BehaviourFunc
// 	}{
// 		{eurSrv, 0, []int64{1, 2}, nil},
// 		{eurSrv, scope.DefaultTypeID, []int64{1, 2}, nil},
// 		{eurSrv, scope.Website.WithID(0), []int64{0}, nil},
// 		{eurSrv, scope.Website.WithID(1), []int64{1, 4, 2}, nil},
// 		{eurSrv, scope.Website.WithID(2), []int64{5, 6}, nil},
// 		{eurSrv, scope.Website.WithID(3), nil, nil},
// 		{eurSrv, scope.Absent.WithID(0), []int64{1, 2}, nil},
// 		{eurSrv, scope.Group.WithID(2), []int64{4}, nil},
// 		{eurSrv, scope.Group.WithID(99), nil, nil},
// 		{eurSrv, scope.Store.WithID(0), []int64{0, 5, 1, 4, 2, 6}, nil},
// 	}
// 	for i, test := range tests {
// 		haveStores, haveErr := test.srv.AllowedStores(test.runMode)
// 		if test.wantErrBhf != nil {
// 			assert.True(t, test.wantErrBhf(haveErr), "(%d) %+v", i, haveErr)
// 			assert.Exactly(t, test.wantStoreIDs, haveStores.IDs(), "Index %d", i)
// 			continue
// 		}
// 		assert.NoError(t, haveErr, "(%d) %+v", i, haveErr)
// 		assert.Exactly(t, test.wantStoreIDs, haveStores.IDs(), "Index %d", i)
// 	}
// }
//
// func TestService_HasSingleStore(t *testing.T) {
// 	s := MustNewService(cfgmock.NewService(),
// 		store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: null.MakeString("euro"), Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
// 	)
// 	s1 := MustNewService(cfgmock.NewService(),
// 		store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: null.MakeString("euro"), Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
// 	)
// 	s1.SingleStoreModeEnabled = false
//
// 	s2 := storemock.NewEurozzyService(cfgmock.NewService())
//
// 	const iterations = 10
// 	var wg sync.WaitGroup
// 	wg.Add(iterations)
// 	for i := 0; i < iterations; i++ {
// 		go func(wg *sync.WaitGroup) {
// 			defer wg.Done()
// 			assert.True(t, s.HasSingleStore())   // no stores so true
// 			assert.False(t, s1.HasSingleStore()) // no stores but globally disabled so false
// 			assert.False(t, s2.HasSingleStore()) // lots of stores so false
// 		}(&wg)
// 	}
// 	wg.Wait()
// }
//
// func TestService_IsSingleStoreMode(t *testing.T) {
//
// 	const xPath = `general/single_store_mode/enabled`
//
// 	s := MustNewService(cfgmock.NewService(),
// 		store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: null.MakeString("euro"), Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 12, IsDefault: null.BoolFrom(true)}),
// 	)
//
// 	// no stores and backend not set so true
// 	sCfg := cfgmock.NewService().NewScoped(0, 0)
// 	b, err := s.IsSingleStoreMode(sCfg)
// 	assert.NoError(t, err, "%+v", err)
// 	assert.True(t, b)
//
// 	// no stores and backend set but configured with false
// 	s.ClearCache()
// 	sCfg = cfgmock.NewService(cfgmock.PathValue{
// 		cfgpath.MustMakeByString(xPath).BindStore(2).String(): 0,
// 	}).NewScoped(1, 2)
// 	s.BackendSingleStore = cfgmodel.NewBool(xPath, cfgmodel.WithScopeStore())
// 	b, err = s.IsSingleStoreMode(sCfg)
// 	assert.NoError(t, err, "%+v", err)
// 	assert.False(t, b)
//
// 	// no stores and backend set but returns an error
// 	s.ClearCache()
// 	tErr := errors.NewNotImplementedf("Ups")
// 	s.BackendSingleStore = cfgmodel.NewBool(xPath)
// 	s.BackendSingleStore.LastError = tErr
// 	b, err = s.IsSingleStoreMode(config.Scoped{})
// 	assert.True(t, errors.IsNotImplemented(tErr), "%+v", tErr)
// 	assert.False(t, b)
//
// 	s2 := storemock.NewEurozzyService(cfgmock.NewService())
// 	s2.BackendSingleStore = cfgmodel.NewBool(xPath) // returns false always no error
// 	assert.False(t, s2.HasSingleStore())
//
// 	b, err = s2.IsSingleStoreMode(sCfg)
// 	assert.NoError(t, err, "%+v", err)
// 	assert.False(t, b)
//
// 	s2.ClearCache()
// 	s2.BackendSingleStore = cfgmodel.NewBool(xPath, cfgmodel.WithField(&element.Field{ID: cfgpath.MakeRoute(`enabled`), Default: `1`})) // returns true
// 	b, err = s2.IsSingleStoreMode(sCfg)
// 	assert.NoError(t, err, "%+v", err)
// 	assert.True(t, b)
//
// 	// call it twice to test cache
// 	b, err = s2.IsSingleStoreMode(sCfg)
// 	assert.NoError(t, err, "%+v", err)
// 	assert.True(t, b)
// }

func TestService_LoadFromDB_OK(t *testing.T) {

	t.Skip("TODO")

	//dbrCon, dbMock := cstesting.MockDB(t)
	//
	//dbMock.ExpectQuery("SELECT (.+) FROM `store`(.+) ORDER BY CASE WHEN(.+)").WillReturnRows(
	//	cstesting.MustMockRows(cstesting.WithFile("testdata", "m1_core_store_view.csv")),
	//)
	//dbMock.ExpectQuery("SELECT (.+) FROM `store_website`(.+) ORDER BY(.+)").WillReturnRows(
	//	cstesting.MustMockRows(cstesting.WithFile("testdata", "m1_core_website_view.csv")),
	//)
	//dbMock.ExpectQuery("SELECT (.+) FROM `store_group`(.+) ORDER BY main_table(.+)").WillReturnRows(
	//	cstesting.MustMockRows(cstesting.WithFile("testdata", "m1_core_store_group_view.csv")),
	//)
	//dbMock.MatchExpectationsInOrder(false) // we're using goroutines!
	//
	//srv := MustNewService(cfgmock.NewService())
	//if err := srv.LoadFromResource(dbrCon.NewSession()); err != nil {
	//	t.Fatalf("%+v", err)
	//}
	//
	//if err := dbMock.ExpectationsWereMet(); err != nil {
	//	t.Fatalf("%+v", err)
	//}
	//assert.Len(t, srv.Websites(), 9)
	//assert.Len(t, srv.Groups(), 9)
	//assert.Len(t, srv.Stores(), 16)
	//
	//tree, err := json.Marshal(srv.Websites().Tree())
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//assert.Exactly(t,
	//	`{"scope":"Default","id":0,"scopes":[{"scope":"Website","id":0,"scopes":[{"scope":"Group","id":0,"scopes":[{"scope":"Store","id":0}]}]},{"scope":"Website","id":2,"scopes":[{"scope":"Group","id":2,"scopes":[{"scope":"Store","id":2},{"scope":"Store","id":5}]}]},{"scope":"Website","id":3,"scopes":[{"scope":"Group","id":3,"scopes":[{"scope":"Store","id":6},{"scope":"Store","id":7},{"scope":"Store","id":8},{"scope":"Store","id":9}]}]},{"scope":"Website","id":4,"scopes":[{"scope":"Group","id":4,"scopes":[{"scope":"Store","id":10},{"scope":"Store","id":11}]}]},{"scope":"Website","id":5,"scopes":[{"scope":"Group","id":5,"scopes":[{"scope":"Store","id":12}]}]},{"scope":"Website","id":6,"scopes":[{"scope":"Group","id":6,"scopes":[{"scope":"Store","id":13},{"scope":"Store","id":14}]}]},{"scope":"Website","id":7,"scopes":[{"scope":"Group","id":7,"scopes":[{"scope":"Store","id":15},{"scope":"Store","id":16}]}]},{"scope":"Website","id":8,"scopes":[{"scope":"Group","id":8,"scopes":[{"scope":"Store","id":17}]}]},{"scope":"Website","id":9,"scopes":[{"scope":"Group","id":9,"scopes":[{"scope":"Store","id":18}]}]}]}`,
	//	string(tree))
}

func TestService_LoadFromDB_NOK_Store(t *testing.T) {
	t.Skip("TODO")

	//dbrCon, dbMock := cstesting.MockDB(t)
	//
	//wsErr := errors.NewAlreadyClosedf("DB Already closed")
	//dbMock.ExpectQuery("SELECT (.+) FROM `store`(.+) ORDER BY CASE WHEN(.+)").WillReturnError(wsErr)
	//dbMock.ExpectQuery("SELECT (.+) FROM `store_website`(.+) ORDER BY(.+)").WillReturnRows(
	//	cstesting.MustMockRows(cstesting.WithFile("testdata", "m1_core_website_view.csv")),
	//)
	//dbMock.ExpectQuery("SELECT (.+) FROM `store_group`(.+) ORDER BY main_table(.+)").WillReturnRows(
	//	cstesting.MustMockRows(cstesting.WithFile("testdata", "m1_core_store_group_view.csv")),
	//)
	//dbMock.MatchExpectationsInOrder(false) // we're using goroutines!
	//
	//srv := MustNewService(cfgmock.NewService())
	//err := srv.LoadFromResource(dbrCon.NewSession())
	//assert.True(t, errors.IsAlreadyClosed(err))
	//
	//if err := dbMock.ExpectationsWereMet(); err != nil {
	//	t.Fatalf("%+v", err)
	//}
	//assert.Len(t, srv.Websites(), 0)
	//assert.Len(t, srv.Groups(), 0)
	//assert.Len(t, srv.Stores(), 0)

}

func TestService_LoadFromDB_NOK_All(t *testing.T) {

	t.Skip("TODO")

	//dbrCon, dbMock := cstesting.MockDB(t)
	//
	//wsErr1 := errors.NewAlreadyClosedf("DB Already closed")
	//wsErr2 := errors.NewNotImplementedf("DB is NoSQL")
	//wsErr3 := errors.NewEmptyf("DB empty")
	//dbMock.ExpectQuery("SELECT (.+) FROM `store`(.+) ORDER BY CASE WHEN(.+)").WillReturnError(wsErr1)
	//dbMock.ExpectQuery("SELECT (.+) FROM `store_website`(.+) ORDER BY(.+)").WillReturnError(wsErr2)
	//dbMock.ExpectQuery("SELECT (.+) FROM `store_group`(.+) ORDER BY main_table(.+)").WillReturnError(wsErr3)
	//dbMock.MatchExpectationsInOrder(false) // we're using goroutines!
	//
	//srv := MustNewService(cfgmock.NewService())
	//err := srv.LoadFromResource(dbrCon.NewSession())
	//assert.True(t, errors.IsAlreadyClosed(err) || errors.IsNotImplemented(err) || errors.IsEmpty(err), "%+v", err)
	//
	//if err := dbMock.ExpectationsWereMet(); err != nil {
	//	t.Fatalf("%+v", err)
	//}
	//assert.Len(t, srv.Websites(), 0)
	//assert.Len(t, srv.Groups(), 0)
	//assert.Len(t, srv.Stores(), 0)

}
