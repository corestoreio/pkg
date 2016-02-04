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
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
)

func ExampleService() {

	// default storage engine with build-in in-memory map.
	srv := config.NewService( /*options*/ )

	// now add some configuration values with different scopes.
	// normally these config values will be loaded from the core_config_data table
	// via function ApplyCoreConfigData()
	// The scope is default -> website (ID 3) -> group (ID 1) -> store (ID 2).
	// The IDs are picked randomly here. Group config values are officially not
	// supported, but we do ;-)

	// We focus here on type String() other primitive types are of course also available.

	var pathScopeTestString = path.MustNewByParts("scope/test/string") // panics on incorrect argument.

	// scope default:
	if err := srv.Write(pathScopeTestString, "DefaultGopher"); err != nil {
		fmt.Printf("Write Error: %s", err)
		return
	}

	// scope website. The number 3 is made up and comes usually from DB table
	// (M1) core_website or (M2) store_website.
	if err := srv.Write(pathScopeTestString.Bind(scope.WebsiteID, 3), "WebsiteGopher"); err != nil {
		fmt.Printf("Write Error: %s", err)
		return
	}

	// scope store. The number 2 is made up and comes usually from DB table
	// (M1) core_store or (M2) store.
	if err := srv.Write(pathScopeTestString.Bind(scope.StoreID, 2), "StoreGopher"); err != nil {
		fmt.Printf("Write Error: %s", err)
		return
	}

	val, err := srv.String(pathScopeTestString)
	if err != nil {
		fmt.Printf("srvString Error: %s", err)
		return
	}
	fmt.Println("Scope1:", val)

	//todo more getters with different scope

	// Output:
	// Scope1: DefaultGopher
}
