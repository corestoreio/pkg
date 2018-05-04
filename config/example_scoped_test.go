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

package config_test

//import (
//	"fmt"
//
//	"github.com/corestoreio/errors"
//	"github.com/corestoreio/pkg/config"
//	"github.com/corestoreio/pkg/storage/null"
//	"github.com/corestoreio/pkg/store"
//	"github.com/corestoreio/pkg/store/scope"
//)
//
//// Config Service, the Default storage engine with build-in in-memory map. The
//// config.NewService or config.MustNewService gets only instantiated once during
//// app start up.
//var configService = config.MustNewService(config.NewInMemoryStore() /*options*/)
//
//// The store.MustNewService gets only instantiated once during app start up.
//var storeSrv = store.MustNewService(
//	configService,
//
//	// Storage gets usually loaded from the database tables containing
//	// website, group and store. For the sake of this example the storage
//	// is hard coded.
//	store.WithTableWebsites(
//		&store.TableWebsite{WebsiteID: 0, Code: null.MakeString("admin"), Name: null.MakeString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: null.MakeBool(false)},
//		&store.TableWebsite{WebsiteID: 1, Code: null.MakeString("euro"), Name: null.MakeString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.MakeBool(true)},
//		&store.TableWebsite{WebsiteID: 2, Code: null.MakeString("oz"), Name: null.MakeString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.MakeBool(false)},
//	),
//	store.WithTableGroups(
//		&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
//		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
//		&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
//		&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
//	),
//	store.WithTableStores(
//		&store.TableStore{StoreID: 0, Code: null.MakeString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
//		&store.TableStore{StoreID: 5, Code: null.MakeString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
//		&store.TableStore{StoreID: 1, Code: null.MakeString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
//		&store.TableStore{StoreID: 4, Code: null.MakeString("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
//		&store.TableStore{StoreID: 2, Code: null.MakeString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
//		&store.TableStore{StoreID: 6, Code: null.MakeString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
//		&store.TableStore{IsActive: false, StoreID: 3, Code: null.MakeString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30},
//	),
//)
//
//const exampleRoute = "scope/test/demo_path"
//
//// We focus here on type Int() other primitive types are of course also available.
//// The int numbers here are converted floats. Can you spot the origin?
//var myConfigPath = config.MustNewPath(exampleRoute) // panics on incorrect argument. Use only during app start up.
//
//var defaultsInt = []struct {
//	key config.Path
//	val []byte
//}{
//	{myConfigPath, []byte(`ScopeDefaultValue`)},                       // Default
//	{myConfigPath.BindWebsite(1), []byte(`ScopeWebsiteOneEuroValue`)}, // Scope 1 = Website euro
//	{myConfigPath.BindStore(2), []byte(`ScopeStoreTwoATValue`)},       // Scope 2 = Store at
//}
//
//func ExampleScopedGetter() {
//
//	// now add some configuration values with different scopes.
//	// normally these config values will be loaded from the core_config_data table
//	// via function ApplyCoreConfigData()
//
//	for _, vi := range defaultsInt {
//		if err := configService.Write(vi.key, vi.val); err != nil {
//			fmt.Printf("Write Error: %s", err)
//			return
//		}
//	}
//
//	atStore, err := storeSrv.Store(2)
//	if err != nil {
//		fmt.Printf("testStoreService.Store Error: %s", err)
//		return
//	}
//
//	// deStore.Config contains our ScopedGetter interface which has been bound
//	// already to the appropriate scopes.
//
//	// Scope1 use store config and hence store value
//	val := atStore.Config.Value(scope.Store, exampleRoute) //.Int(myConfigPath.Route)
//	valStr, ok, err := val.Str()
//	if err != nil {
//		fmt.Printf("srvString1 Error: %+v", err)
//		return
//	}
//	if !ok {
//		fmt.Printf("srvString1 NULL: %+v", err)
//		return
//	}
//	fmt.Println("Scope Value for Store ID 2:", valStr)
//
//	// Scope2 use website config and hence website value
//	val, err = atStore.Website.Config.Int(myConfigPath.Route)
//	if err != nil {
//		fmt.Printf("srvString2 Error: %+v", err)
//		return
//	}
//	fmt.Println("Scope Value for Website ID 1:", val)
//
//	// Scope3 force default value
//	val, err = atStore.Config.Int(myConfigPath.Route, scope.Default)
//	if err != nil {
//		fmt.Printf("srvString3 Error: %+v", err)
//		return
//	}
//	fmt.Println("Scope Value for Default:", val)
//
//	// Scope4 route not found
//	_, err = atStore.Config.Int(cfgpath.MustMakeByString("xx/yy/zz").Route)
//	if err != nil {
//		fmt.Printf("Scope4: srvString Error: %s\n", err)
//	}
//	fmt.Printf("Route IsNotFound %t\n", errors.NotFound.Match(err))
//
//	// Output:
//	// Scope Value for Store ID 2: 141421
//	// Scope Value for Website ID 1: 271828
//	// Scope Value for Default: 314159
//	// Scope4: srvString Error: [config] Storage.Int.get: [config] KVMap Unknown Key: default/0/xx/yy/zz
//	// Route IsNotFound true
//}
