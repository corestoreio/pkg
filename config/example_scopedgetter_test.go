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

//import (
//	"fmt"
//
//	"github.com/corestoreio/csfw/config"
//	"github.com/corestoreio/csfw/config/path"
//	"github.com/corestoreio/csfw/storage/dbr"
//	"github.com/corestoreio/csfw/store"
//	"github.com/corestoreio/csfw/store/scope"
//)
//
//// Default storage engine with build-in in-memory map.
//// the config.NewService or config.MustNewService gets only instantiated once
//// during app start up.
//var configSrv = config.MustNewService( /*options*/ )
//
//// The store.MustNewService gets only instantiated once during app start up.
//var storeSrv = store.MustNewService(
//	scope.Option{
//		// bound to website ID 1 = Europe
//		// This gets set during app start up and a HTTP/RPC request cannot changed the bound scope.
//		Website: scope.MockID(1),
//	},
//	store.MustNewStorage(
//		// Storage gets usually loaded from the database tables containing
//		// website, group and store. For the sake of this example the storage
//		// is hard coded.
//		store.SetStorageWebsites(
//			&store.TableWebsite{WebsiteID: 0, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
//			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
//			&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
//		),
//		store.SetStorageGroups(
//			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
//			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
//			&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
//			&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
//		),
//		store.SetStorageStores(
//			&store.TableStore{StoreID: 0, Code: dbr.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
//			&store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
//			&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
//			&store.TableStore{StoreID: 4, Code: dbr.NewNullString("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
//			&store.TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
//			&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
//			&store.TableStore{IsActive: false, StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30},
//		),
//	),
//	store.WithServiceConfigReader(configSrv),
//)
//
//// We focus here on type Int() other primitive types are of course also available.
//// The int numbers here are converted floats. Can you spot the origin?
//var pathInt = path.MustNewByParts("scope/test/integer") // panics on incorrect argument. Use only during app start up.
//
//var defaultsInt = struct {
//	key path.Path
//	val int
//}{
//	{pathInt, 314159},                          // Default
//	{pathInt.Bind(scope.WebsiteID, 1), 271828}, // Scope 1 = Website euro
//	{pathInt.Bind(scope.StoreID, 2), 141421},   // Scope 2 = Store de
//}
//
//func ExampleScopedGetter() {
//
//	// now add some configuration values with different scopes.
//	// normally these config values will be loaded from the core_config_data table
//	// via function ApplyCoreConfigData()
//
//	for k, v := range defaultsInt {
//		if err := configSrv.Write(k, v); err != nil {
//			fmt.Printf("Write Error: %s", err)
//			return
//		}
//	}
//
//	deStore, err := storeSrv.Store(scope.MockID(1))
//	if err != nil {
//		fmt.Printf("testStoreService.Store Error: %s", err)
//		return
//	}
//
//	// deStore.Config contains our ScopedGetter interface which has been bound
//	// already to the appropriate scopes.
//
//	// Scope1
//	val, err := deStore.Config.Int(pathInt)
//	if err != nil {
//		fmt.Printf("srvString Error: %s", err)
//		return
//	}
//	fmt.Println("Scope1:", val)
//
//	// Scope2
//	val, err = deStore.Config.Int(pathInt.Bind(scope.WebsiteID, 3))
//	if err != nil {
//		fmt.Printf("srvString Error: %s", err)
//		return
//	}
//	fmt.Println("Scope2:", val)
//
//	// Scope3
//	val, err = deStore.Config.Int(pathInt.Bind(scope.StoreID, 2))
//	if err != nil {
//		fmt.Printf("srvString Error: %s", err)
//		return
//	}
//	fmt.Println("Scope3:", val)
//
//	// Scope4
//	val, err = deStore.Config.Int(pathInt.Bind(scope.StoreID, 3)) // different scope ID
//	if err != nil {
//		fmt.Printf("Scope4: srvString Error: %s\n", err)
//	}
//	fmt.Printf("Scope4: Is KeyNotFound %t\n", !config.NotKeyNotFoundError(err))
//
//	// Output:
//	// Scope1: todo
//}
