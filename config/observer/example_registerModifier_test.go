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

// ExampleRegisterModifier register a custom modifier with the name
// `append_test` which appends the word TEST. The modifier `append_test` gets
// activated for pre-route `payment/serviceX` with event after_get.
func ExampleRegisterModifier() {
	observer.RegisterModifier("append_test", func(_ *config.Path, data []byte) ([]byte, error) {
		return append(data, []byte(` - TEST`)...), nil
	})

	cfgSrv := config.MustNewService(storage.NewMap("stores/2/payment/serviceX/username", "xyzUser"), config.Options{})

	err := observer.RegisterWithJSON(cfgSrv, bytes.NewBufferString(`{"Collection":[ { 
		  "event":"after_get", "route":"payment/serviceX", "type":"modifier",
		  "condition":{"funcs":["append_test"]}}
		]}`))
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	ps2 := config.MustNewPath("payment/serviceX/username").BindStore(2)
	val := cfgSrv.Get(ps2)
	fmt.Printf("%s\n", val.String())

	// Output:
	// "xyzUser - TEST"
}
