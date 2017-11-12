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

package directory_test

import "github.com/corestoreio/pkg/directory"

// backend overall backend models for all tests
var backend *directory.PkgBackend

// this would belong into the test suit setup
func init() {
	cfgStruct, err := directory.NewConfigStructure()
	if err != nil {
		panic(err)
	}

	backend = directory.NewBackend(cfgStruct)

	src, err := backend.InitSources(nil) // TODO(cs) add DB
	if err != nil {
		panic(err)
	}
	if err := backend.InitCountry(nil, src); err != nil { // TODO(cs) add DB
		panic(err)
	}
	if err := backend.InitCurrency(nil, src); err != nil { // TODO(cs) add DB
		panic(err)
	}
}
