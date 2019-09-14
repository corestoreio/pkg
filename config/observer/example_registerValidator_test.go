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

// +build csall json

package observer_test

import (
	"bytes"
	"fmt"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/observer"
	"github.com/corestoreio/pkg/config/storage"
)

// ExampleRegisterValidator register a custom validator with the name
// `check_xyzUser`. It checks that the configuration value for route
// `payment/serviceX/username` must be equal to string `xyzUser`. (boring
// example) The validator `check_xyzUser` gets activated for route
// `payment/serviceX/username` with event after_get.
func ExampleRegisterValidator() {
	observer.RegisterValidator("check_xyzUser", func(s string) bool {
		return s == "xyzUser"
	})

	cfgSrv := config.MustNewService(storage.NewMap(
		"stores/2/payment/serviceX/username", "xyzUser",
		"stores/3/payment/serviceX/username", "abcUser",
	), config.Options{})

	err := observer.RegisterWithJSON(cfgSrv, bytes.NewBufferString(`{"Collection":[ { 
		  "event":"after_get", "route":"payment/serviceX/username", "type":"validator",
		  "condition":{"funcs":["utf8","check_xyzUser"]}}
		]}`))
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	ps2 := config.MustNewPath("payment/serviceX/username").BindStore(2)
	val := cfgSrv.Get(ps2)
	fmt.Printf("%s\n", val.String())

	val = cfgSrv.Get(ps2.BindStore(3))
	fmt.Printf("%s\n", val.Error())

	// Output:
	// "xyzUser"
	// [config/observer] The value "<redacted>" can't be validated against ["utf8" "check_xyzUser"]
}
