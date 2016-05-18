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

	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func TestMmws_Country_Failure_Response(t *testing.T) {

	c := newMMWS("gopher", "passw0rd", time.Second)
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

	c := newMMWS("a", "b", time.Second)
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

	c := newMMWS("gopher", "passw0rd", time.Second)
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
