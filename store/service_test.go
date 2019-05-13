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
	"fmt"
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
	store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1, Code: "dach"}),
	store.WithStores(&store.Store{StoreID: 1, Code: "de", WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true}),
)

func TestNewServiceStore_QueryInvalidStore(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.ErrorIsKind(t, errors.NotValid, err)
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
	t.Parallel()
	tests := []struct {
		name    string
		opts    []store.Option
		wantErr errors.Kind
	}{
		{"All valid",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true, Code: "dach"}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1, Code: "dach"}),
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
		{"Website.Code is empty",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1, Code: "dach"}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 1, GroupID: 1, Code: "de-de"}),
			},
			errors.NotValid},
		{"Group.Code is empty",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true, Code: "dach"}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 1, GroupID: 1, Code: "de-de"}),
			},
			errors.NotValid},
		{"Store.Code is empty",
			[]store.Option{
				store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, DefaultGroupID: 1, IsDefault: true, Code: "dach"}),
				store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, DefaultStoreID: 1, Code: "dach"}),
				store.WithStores(&store.Store{StoreID: 1, WebsiteID: 1, GroupID: 1}),
			},
			errors.NotValid},
		// TODO add Website.Validate and Group.Validate
	}

	s := store.MustNewService()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := s.Options(test.opts...); test.wantErr > 0 {
				assert.ErrorIsKind(t, test.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
			s.ClearCache()
		})
	}
}

func TestNewService_DefaultStoreView(t *testing.T) {
	t.Parallel()
	srv := store.MustNewService(
		store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: "euro", Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: true}),
		store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, Code: "dach", Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1}),
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
	t.Parallel()
	stores := store.Stores{
		Data: []*store.Store{
			{StoreID: 1, Code: "de", WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			{StoreID: 2, Code: "at", WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			{StoreID: 3, Code: "ch", WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
		},
	}
	groups := store.StoreGroups{
		Data: []*store.StoreGroup{
			{GroupID: 1, WebsiteID: 1, Name: "DACH Group", Code: "dach", RootCategoryID: 2, DefaultStoreID: 2},
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
	t.Parallel()
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
		{scope.MakeTypeID(scope.Store, 0), 3, false, "ch"}, // #23
		{scope.MakeTypeID(scope.Store, 1), 4, true, "uk"},
		{scope.MakeTypeID(scope.Store, 9999), 4, true, "uk"},
		{scope.MakeTypeID(124, 1), 4, false, ""},
		{scope.MakeTypeID(124, 0), 4, false, ""},
	}
	eoSrv := storemock.NewServiceEuroOZ()
	for i, test := range tests {
		t.Run(fmt.Sprintf("Index_%02d", i), func(t *testing.T) {
			haveIsAllowed, haveCode, haveErr := eoSrv.IsAllowedStoreID(test.runMode, test.storeID)
			assert.NoError(t, haveErr)
			assert.Exactly(t, test.wantIsAllowed, haveIsAllowed)
			assert.Exactly(t, test.wantCode, haveCode)
		})
	}
}

func TestService_DefaultStoreID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		runMode       scope.TypeID
		wantStoreID   uint32
		wantWebsiteID uint32
		wantErrBhf    errors.Kind
	}{
		{0, 2, 1, errors.NoKind},                                  // fall back to default website -> default group -> default store
		{scope.MakeTypeID(scope.Website, 0), 0, 0, errors.NoKind}, // admin scope
		{scope.MakeTypeID(scope.Website, 1), 2, 1, errors.NoKind}, // euro scope, not included ch, because not active, and UK, different group
		{scope.MakeTypeID(scope.Website, 2), 5, 2, errors.NoKind}, // oz scope
		{scope.MakeTypeID(scope.Website, 9999), 0, 0, errors.NotFound},
		{scope.MakeTypeID(scope.Group, 0), 0, 0, errors.NoKind}, // admin scope
		{scope.MakeTypeID(scope.Group, 1), 2, 1, errors.NoKind}, // dach scope
		{scope.MakeTypeID(scope.Group, 2), 4, 1, errors.NoKind}, // uk scope
		{scope.MakeTypeID(scope.Group, 3), 5, 2, errors.NoKind}, // au scope
		{scope.MakeTypeID(scope.Group, 9999), 0, 0, errors.NotFound},
		{scope.MakeTypeID(scope.Store, 0), 0, 0, errors.NoKind},
		{scope.MakeTypeID(scope.Store, 1), 1, 1, errors.NoKind},
		{scope.MakeTypeID(scope.Store, 9999), 0, 0, errors.NotFound},
		{scope.MakeTypeID(scope.Store, 3), 0, 0, errors.NotValid}, // ch store is not active
	}
	eurSrv := storemock.NewServiceEuroOZ()
	for i, test := range tests {
		t.Run(fmt.Sprintf("Index_%02d", i), func(t *testing.T) {
			haveStoreID, haveWebsiteID, haveErr := eurSrv.DefaultStoreID(test.runMode)
			if test.wantErrBhf > 0 {
				assert.ErrorIsKind(t, test.wantErrBhf, haveErr)
				assert.Exactly(t, test.wantStoreID, haveStoreID)
				return
			}
			assert.NoError(t, haveErr)
			assert.Exactly(t,
				fmt.Sprintf("StoreID %d WebsiteID %d", test.wantStoreID, test.wantWebsiteID),
				fmt.Sprintf("StoreID %d WebsiteID %d", haveStoreID, haveWebsiteID),
			)
		})
	}

	t.Run("default store not active", func(t *testing.T) {
		srv := store.MustNewService( // default store not active
			store.WithWebsites(&store.StoreWebsite{WebsiteID: 1, Code: "euro", Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: true}),
			store.WithGroups(&store.StoreGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 1, Code: "dach"}),
			store.WithStores(&store.Store{StoreID: 1, Code: "de", WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: false}),
		)
		haveStoreID, haveWebsiteID, haveErr := srv.DefaultStoreID(scope.MakeTypeID(scope.Website, 1))
		assert.ErrorIsKind(t, errors.NotFound, haveErr)
		assert.Exactly(t,
			fmt.Sprintf("StoreID %d WebsiteID %d", 0, 0),
			fmt.Sprintf("StoreID %d WebsiteID %d", haveStoreID, haveWebsiteID),
		)
	})
}

func TestService_StoreIDbyCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		runMode       scope.TypeID
		code          string
		wantStoreID   uint32
		wantWebsiteID uint32
		wantErrBhf    errors.Kind
	}{
		{0, "", 2, 1, errors.NoKind},
		{scope.DefaultTypeID, "x", 0, 0, errors.NotFound},
		{scope.DefaultTypeID, "uk", 0, 0, errors.NotFound},
		{scope.Website.WithID(0), "admin", 0, 0, errors.NoKind},
		{scope.Website.WithID(1), "de", 1, 1, errors.NoKind},
		{scope.Website.WithID(2), "nz", 6, 2, errors.NoKind},
		{scope.Website.WithID(3), "uk", 0, 0, errors.NotFound},
		{scope.Absent.WithID(0), "uk", 0, 0, errors.NotFound},
		{scope.Absent.WithID(0), "at", 2, 1, errors.NoKind},
		{scope.Group.WithID(2), "uk", 4, 1, errors.NoKind},
		{scope.Group.WithID(99), "uk", 0, 0, errors.NotFound},
		{scope.Store.WithID(0), "admin", 0, 0, errors.NoKind},
		{scope.Store.WithID(0), "au", 5, 2, errors.NoKind},
		{scope.Store.WithID(0), "xx", 0, 0, errors.NotFound},
	}
	eurSrv := storemock.NewServiceEuroOZ()
	for i, test := range tests {
		t.Run(fmt.Sprintf("Index_%02d", i), func(t *testing.T) {
			haveStoreID, haveWebsiteID, haveErr := eurSrv.StoreIDbyCode(test.runMode, test.code)
			if test.wantErrBhf > 0 {
				assert.ErrorIsKind(t, test.wantErrBhf, haveErr)
				assert.Exactly(t, test.wantStoreID, haveStoreID)
				assert.Exactly(t, test.wantWebsiteID, haveWebsiteID)
				return
			}
			assert.NoError(t, haveErr)
			assert.Exactly(t, test.wantStoreID, haveStoreID)
			assert.Exactly(t, test.wantWebsiteID, haveWebsiteID)
		})
	}
}

func TestService_AllowedStores_Implementation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		runMode      scope.TypeID
		wantStoreIDs []uint32
		wantErrBhf   errors.Kind
	}{
		{0, nil, errors.NotImplemented},
		{scope.DefaultTypeID, []uint32{0}, errors.NoKind},
		{scope.Website.WithID(0), []uint32{0}, errors.NoKind},
		{scope.Website.WithID(1), []uint32{1, 2, 4}, errors.NoKind},
		{scope.Website.WithID(2), []uint32{5, 6}, errors.NoKind},
		{scope.Website.WithID(3), []uint32{}, errors.NoKind},
		{scope.Absent.WithID(0), nil, errors.NotImplemented},
		{scope.Group.WithID(2), []uint32{4}, errors.NoKind},
		{scope.Group.WithID(99), []uint32{}, errors.NoKind},
		{scope.Store.WithID(0), []uint32{0, 1, 2, 4, 5, 6}, errors.NoKind},
		{87987, nil, errors.NotImplemented},
		{2789987, nil, errors.NotImplemented},
	}
	eurSrv := storemock.NewServiceEuroOZ()
	for i, test := range tests {
		t.Run(fmt.Sprintf("Index_%02d_%s", i, test.runMode.String()), func(t *testing.T) {
			haveStores, haveErr := eurSrv.AllowedStores(test.runMode)
			if test.wantErrBhf > 0 {
				assert.ErrorIsKind(t, test.wantErrBhf, haveErr)
				assert.Exactly(t, test.wantStoreIDs, haveStores.StoreIDs())
				return
			}
			assert.NoError(t, haveErr)
			assert.Exactly(t, test.wantStoreIDs, haveStores.StoreIDs())

			// todo check if store pointer not returned still exists.
		})
	}
}

func TestService_AllowedStores_PointerCheck(t *testing.T) {
	eurSrv := storemock.NewServiceEuroOZ()
	sts, err := eurSrv.AllowedStores(scope.Website.WithID(1))
	assert.NoError(t, err)

	assert.Exactly(t, []uint32{1, 2, 4}, sts.StoreIDs())

	st, err := eurSrv.Store(5)
	assert.NoError(t, err)
	assert.Exactly(t, uint32(5), st.StoreID)

	st, err = eurSrv.Store(6)
	assert.NoError(t, err)
	assert.Exactly(t, uint32(5), st.StoreID)
}

func TestService_LoadFromDB_OK(t *testing.T) {

	t.Skip("TODO")

	// dbrCon, dbMock := cstesting.MockDB(t)
	//
	// dbMock.ExpectQuery("SELECT (.+) FROM `store`(.+) ORDER BY CASE WHEN(.+)").WillReturnRows(
	// 	cstesting.MustMockRows(cstesting.WithFile("testdata", "m1_core_store_view.csv")),
	// )
	// dbMock.ExpectQuery("SELECT (.+) FROM `store_website`(.+) ORDER BY(.+)").WillReturnRows(
	// 	cstesting.MustMockRows(cstesting.WithFile("testdata", "m1_core_website_view.csv")),
	// )
	// dbMock.ExpectQuery("SELECT (.+) FROM `store_group`(.+) ORDER BY main_table(.+)").WillReturnRows(
	// 	cstesting.MustMockRows(cstesting.WithFile("testdata", "m1_core_store_group_view.csv")),
	// )
	// dbMock.MatchExpectationsInOrder(false) // we're using goroutines!
	//
	// srv := MustNewService(cfgmock.NewService())
	// if err := srv.LoadFromResource(dbrCon.NewSession()); err != nil {
	// 	t.Fatalf("%+v", err)
	// }
	//
	// if err := dbMock.ExpectationsWereMet(); err != nil {
	// 	t.Fatalf("%+v", err)
	// }
	// assert.Len(t, srv.Websites(), 9)
	// assert.Len(t, srv.Groups(), 9)
	// assert.Len(t, srv.Stores(), 16)
	//
	// tree, err := json.Marshal(srv.Websites().Tree())
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// assert.Exactly(t,
	// 	`{"scope":"Default","id":0,"scopes":[{"scope":"Website","id":0,"scopes":[{"scope":"Group","id":0,"scopes":[{"scope":"Store","id":0}]}]},{"scope":"Website","id":2,"scopes":[{"scope":"Group","id":2,"scopes":[{"scope":"Store","id":2},{"scope":"Store","id":5}]}]},{"scope":"Website","id":3,"scopes":[{"scope":"Group","id":3,"scopes":[{"scope":"Store","id":6},{"scope":"Store","id":7},{"scope":"Store","id":8},{"scope":"Store","id":9}]}]},{"scope":"Website","id":4,"scopes":[{"scope":"Group","id":4,"scopes":[{"scope":"Store","id":10},{"scope":"Store","id":11}]}]},{"scope":"Website","id":5,"scopes":[{"scope":"Group","id":5,"scopes":[{"scope":"Store","id":12}]}]},{"scope":"Website","id":6,"scopes":[{"scope":"Group","id":6,"scopes":[{"scope":"Store","id":13},{"scope":"Store","id":14}]}]},{"scope":"Website","id":7,"scopes":[{"scope":"Group","id":7,"scopes":[{"scope":"Store","id":15},{"scope":"Store","id":16}]}]},{"scope":"Website","id":8,"scopes":[{"scope":"Group","id":8,"scopes":[{"scope":"Store","id":17}]}]},{"scope":"Website","id":9,"scopes":[{"scope":"Group","id":9,"scopes":[{"scope":"Store","id":18}]}]}]}`,
	// 	string(tree))
}
