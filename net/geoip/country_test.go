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

package geoip

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"path/filepath"
	"testing"

	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ CountryRetriever = (*mmdb)(nil)

func TestCountry_JSON(t *testing.T) {
	td, err := ioutil.ReadFile("testdata/response.json")
	if err != nil {
		t.Fatal(err)
	}
	var c Country
	if err = json.Unmarshal(td, &c); err != nil {
		t.Fatal(err)
	}

	haveTD, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, string(td), string(haveTD)+"\n")
}

func TestMmdb_Country(t *testing.T) {
	maxMindDB := filepath.Join("testdata", "GeoIP2-Country-Test.mmdb")
	r, err := newMMDBByFile(maxMindDB)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	c, err := r.Country(nil)
	assert.Nil(t, c)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)

	ip, _, err := net.ParseCIDR("2a02:d200::/29") // IP range for Finland
	if err != nil {
		t.Fatal(err)
	}
	c, err = r.Country(ip)
	assert.NoError(t, err)
	assert.Exactly(t, "FI", c.Country.IsoCode)
}
