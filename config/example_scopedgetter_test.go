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

package config_test

import (
	"fmt"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// Default storage engine with build-in in-memory map.
// the config.NewService or config.MustNewService gets only instantiated once
// during app start up.
var configSrv2 = config.MustNewService( /*options*/ )

// The store.MustNewService gets only instantiated once during app start up.
var storeSrv = store.MustNewService(
	scope.Option{
		// bound to website ID 1 = Europe
		// This gets set during app start up and a HTTP/RPC request cannot changed the bound scope.
		Website: scope.MockID(1),
	},
	store.MustNewStorage(
		// Storage gets usually loaded from the database tables containing
		// website, group and store. For the sake of this example the storage
		// is hard coded.
		store.SetStorageWebsites(
			&store.TableWebsite{WebsiteID: 0, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
		),
		store.SetStorageGroups(
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
		),
		store.SetStorageStores(
			&store.TableStore{StoreID: 0, Code: dbr.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
			&store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 4, Code: dbr.NewNullString("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
			&store.TableStore{IsActive: false, StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30},
		),
		store.WithStorageConfig(configSrv2),
	),
)

// We focus here on type Int() other primitive types are of course also available.
// The int numbers here are converted floats. Can you spot the origin?
var pathInt = cfgpath.MustNewByParts("scope/test/integer") // panics on incorrect argument. Use only during app start up.

var defaultsInt = []struct {
	key cfgpath.Path
	val int
}{
	{pathInt, 314159},                // Default
	{pathInt.BindWebsite(1), 271828}, // Scope 1 = Website euro
	{pathInt.BindStore(2), 141421},   // Scope 2 = Store at
}

func ExampleScopedGetter() {

	// now add some configuration values with different scopes.
	// normally these config values will be loaded from the core_config_data table
	// via function ApplyCoreConfigData()

	for _, vi := range defaultsInt {
		if err := configSrv2.Write(vi.key, vi.val); err != nil {
			fmt.Printf("Write Error: %s", err)
			return
		}
	}

	atStore, err := storeSrv.Store(scope.MockID(2))
	if err != nil {
		fmt.Printf("testStoreService.Store Error: %s", err)
		return
	}

	// deStore.Config contains our ScopedGetter interface which has been bound
	// already to the appropriate scopes.

	// Scope1 use store config and hence store value
	val, h, err := atStore.Config.Int(pathInt.Route)
	if err != nil {
		fmt.Printf("srvString Error: %s", err)
		return
	}
	fmt.Println("Scope Value for Store ID 2:", val, " | ", h.String())

	// Scope2 use website config and hence website value
	val, h, err = atStore.Website.Config.Int(pathInt.Route)
	if err != nil {
		fmt.Printf("srvString Error: %s", err)
		return
	}
	fmt.Println("Scope Value for Website ID 1:", val, " | ", h.String())

	// Scope3 force default value
	val, h, err = atStore.Config.Int(pathInt.Route, scope.Default)
	if err != nil {
		fmt.Printf("srvString Error: %s", err)
		return
	}
	fmt.Println("Scope Value for Default:", val, " | ", h.String())

	// Scope4 route not found
	_, _, err = atStore.Config.Int(cfgpath.MustNewByParts("xx/yy/zz").Route)
	if err != nil {
		fmt.Printf("Scope4: srvString Error: %s\n", err)
	}
	fmt.Printf("Route IsNotFound %t\n", errors.IsNotFound(err))

	// Output:
	// Scope Value for Store ID 2: 141421  |  Scope(Store) ID(2)
	// Scope Value for Website ID 1: 271828  |  Scope(Website) ID(1)
	// Scope Value for Default: 314159  |  Scope(Default) ID(0)
	// Scope4: srvString Error: [config] Storage.Int.get: [storage] Key not found
	// Route IsNotFound true
}
