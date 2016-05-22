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
	"net"
	"testing"
	"time"

	"io/ioutil"

	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ CountryRetriever = (*mmws)(nil)

type mockCacher struct {
	cache  map[string]interface{}
	setErr error
	getErr error
}

func newMockCacher() *mockCacher {
	return &mockCacher{
		cache: make(map[string]interface{}),
	}
}

func (mc *mockCacher) Set(key []byte, src interface{}) error {
	if mc.setErr != nil {
		return mc.setErr
	}
	mc.cache[string(key)] = src
	return nil
}

func (mc *mockCacher) Get(key []byte, dst interface{}) error {
	if mc.getErr != nil {
		return mc.getErr
	}
	if val, ok := mc.cache[string(key)]; ok {
		dst = val
	}
	return errors.NewNotFoundf("[mockCacher] Key %q not found", string(key))

}

func TestMmws_Country_Failure_Response(t *testing.T) {

	c := newMMWS(newMockCacher(), "gopher", "passw0rd", time.Second)
	trip := cstesting.NewHttpTrip(400, `{"error":"Invalid user_id or license_key provided","code":"AUTHORIZATION_INVALID"}`, nil)
	c.client.Transport = trip
	ret, err := c.Country(net.ParseIP("123.123.123.123"))
	assert.Nil(t, ret)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)

	u, p, ok := trip.Req.BasicAuth()
	assert.True(t, ok)
	assert.Exactly(t, "gopher", u)
	assert.Exactly(t, "passw0rd", p)
}

func TestMmws_Country_Failure_JSON(t *testing.T) {

	c := newMMWS(newMockCacher(), "a", "b", time.Second)
	trip := cstesting.NewHttpTrip(200, `"error":"Invalid user_id or license_key provided","code":"AUTHORIZATION_INVALID"}`, nil)
	c.client.Transport = trip
	ret, err := c.Country(net.ParseIP("123.123.123.123"))
	assert.NotNil(t, ret)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
}

func TestMmws_Country_Success(t *testing.T) {
	td, err := ioutil.ReadFile("testdata/response.json")
	if err != nil {
		t.Fatal(err)
	}

	c := newMMWS(newMockCacher(), "gopher", "passw0rd", time.Second)
	trip := cstesting.NewHttpTrip(200, string(td), nil)
	c.client.Transport = trip
	ret, err := c.Country(net.ParseIP("123.123.123.123"))
	assert.NotNil(t, ret)
	assert.NoError(t, err)
	assert.Exactly(t, "US", ret.Country.IsoCode)
	u, p, ok := trip.Req.BasicAuth()
	assert.True(t, ok)
	assert.Exactly(t, "gopher", u)
	assert.Exactly(t, "passw0rd", p)
}
