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
	"encoding/json"
	"fmt"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/observer"
)

func NewAllowedObserver(data []byte) (config.Observer, error) {
	a := &allowed{}
	if err := json.Unmarshal(data, a); err != nil {
		return nil, err
	}
	return a, nil
}

type allowed struct {
	Paths []string
}

func (a *allowed) Observe(p config.Path, rawData []byte, found bool) (newRawData []byte, err error) {
	for _, ap := range a.Paths {
		if _, r := p.ScopeRoute(); r == ap {
			return rawData, nil
		}
	}
	return nil, errors.NotAllowed.Newf("Access not allowed to path %q", p.String())
}

// ExampleRegisterFactory shows how to create a custom observer based on the
// JSON input data. Usually the JSON data gets send via HTTP or protobuf.
func ExampleRegisterFactory() {
	observer.RegisterFactory("path_allowed", NewAllowedObserver)

	cfgSrv := config.MustNewService(nil, config.Options{})

	err := observer.RegisterWithJSON(cfgSrv, bytes.NewBufferString(`{"Collection":[ { 
		  "event":"before_set", "route":"aa/gg", "type":"path_allowed",
		  "condition":{"Paths":["aa/gg/hh","aa/gg/jj"]}}
		]}`))
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	ps3 := config.MustNewPath("aa/gg/kk").BindStore(3)
	err = cfgSrv.Set(ps3, []byte(`GopherCon in San Diego`))
	fmt.Printf("%s\n", err)

	// Output:
	// Access not allowed to path "stores/3/aa/gg/kk"
}
