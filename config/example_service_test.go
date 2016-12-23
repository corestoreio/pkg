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
	"github.com/corestoreio/errors"
)

// We focus here on type String() other primitive types are of course also available.
var pathString = cfgpath.MustNewByParts("scope/test/string") // panics on incorrect argument.

// Default storage engine with build-in in-memory map.
// the NewService gets only instantiated once during app start up.
var configSrv = config.MustNewService(config.NewInMemoryStore() /*options*/)

func ExampleService() {

	// now add some configuration values with different scopes.
	// normally these config values will be loaded from the core_config_data table
	// via function ApplyCoreConfigData()
	// The scope is default -> website (ID 3) -> group (ID 1) -> store (ID 2).
	// The IDs are picked randomly here. Group config values are officially not
	// supported, but we do ;-)

	// scope default:
	if err := configSrv.Write(pathString, "DefaultGopher"); err != nil {
		fmt.Printf("Write Error: %s", err)
		return
	}

	// scope website. The number 3 is made up and comes usually from DB table
	// (M1) core_website or (M2) store_website.
	if err := configSrv.Write(pathString.BindWebsite(3), "WebsiteGopher"); err != nil {
		fmt.Printf("Write Error: %s", err)
		return
	}

	// scope store. The number 2 is made up and comes usually from DB table
	// (M1) core_store or (M2) store.
	if err := configSrv.Write(pathString.BindStore(2), "StoreGopher"); err != nil {
		fmt.Printf("Write Error: %s", err)
		return
	}

	// Scope1
	val, err := configSrv.String(pathString)
	if err != nil {
		fmt.Printf("srvString Error: %s", err)
		return
	}
	fmt.Println("Scope1:", val)

	// Scope2
	val, err = configSrv.String(pathString.BindWebsite(3))
	if err != nil {
		fmt.Printf("srvString Error: %s", err)
		return
	}
	fmt.Println("Scope2:", val)

	// Scope3
	val, err = configSrv.String(pathString.BindStore(2))
	if err != nil {
		fmt.Printf("srvString Error: %s", err)
		return
	}
	fmt.Println("Scope3:", val)

	// Scope4
	_, err = configSrv.String(pathString.BindStore(3)) // different scope ID
	if err != nil {
		fmt.Printf("Scope4a: srvString Error: %s\n", err)
		fmt.Printf("Scope4b: srvString Error: %v\n", err) // Use %+v to show the full path! :-)
	}
	fmt.Printf("Scope4: Is KeyNotFound %t\n", errors.IsNotFound(err))

	// Output:
	// Scope1: DefaultGopher
	// Scope2: WebsiteGopher
	// Scope3: StoreGopher
	// Scope4a: srvString Error: [config] Storage.String.get: [config] KVMap Unknown Key: stores/3/scope/test/string
	// Scope4b: srvString Error: [config] Storage.String.get: [config] KVMap Unknown Key: stores/3/scope/test/string
	// Scope4: Is KeyNotFound true
}
