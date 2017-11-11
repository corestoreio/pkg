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

package backendgeoip_test

import (
	"path/filepath"

	"github.com/corestoreio/cspkg/net/geoip/backendgeoip"
	"github.com/corestoreio/cspkg/net/geoip/maxmindfile"
)

// backend overall backend models for all tests
var backend *backendgeoip.Configuration

var filePathGeoIP string

// this would belong into the test suit setup
func init() {

	filePathGeoIP = filepath.Join("..", "testdata", "GeoIP2-Country-Test.mmdb")

	cfgStruct, err := backendgeoip.NewConfigStructure()
	if err != nil {
		panic(err)
	}
	backend = backendgeoip.New(cfgStruct)

	backend.Register(
		maxmindfile.NewOptionFactory(backend.MaxmindLocalFile),
	)
}
