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
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"testing"

	"github.com/corestoreio/csfw/storage/transcache"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ CountryRetriever = (*mmws)(nil)

func TestMmws_Country_Failure_Response(t *testing.T) {

	ws := newMMWS(transcache.NewMock(), "gopher", "passw0rd", http.DefaultClient)
	trip := cstesting.NewHttpTrip(400, `{"error":"Invalid user_id or license_key provided","code":"AUTHORIZATION_INVALID"}`, nil)
	ws.client.Transport = trip
	c, err := ws.Country(net.ParseIP("123.123.123.123"))
	assert.Nil(t, c)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)

	trip.RequestsMatchAll(t, func(r *http.Request) bool {
		u, p, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Exactly(t, "gopher", u)
		assert.Exactly(t, "passw0rd", p)
		return true
	})
}

func TestMmws_Country_Failure_JSON(t *testing.T) {

	ws := newMMWS(transcache.NewMock(), "a", "b", http.DefaultClient)
	trip := cstesting.NewHttpTrip(200, `"error":"Invalid user_id or license_key provided","code":"AUTHORIZATION_INVALID"}`, nil)
	ws.client.Transport = trip
	c, err := ws.Country(net.ParseIP("123.123.123.123"))
	assert.Nil(t, c)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
}

func TestMmws_Country_Cache_GetError(t *testing.T) {
	tcmock := transcache.NewMock()
	tcmock.GetErr = errors.NewAlreadyClosedf("cache already closed ;-)")

	ws := newMMWS(tcmock, "a", "b", http.DefaultClient)
	trip := cstesting.NewHttpTripFromFile(200, "testdata/response.json")
	ws.client.Transport = trip
	c, err := ws.Country(net.ParseIP("123.123.123.123"))
	assert.Nil(t, c)
	assert.True(t, errors.IsAlreadyClosed(err), "Error: %s", err)
}

func TestMmws_Country_Cache_SetError(t *testing.T) {
	tcmock := transcache.NewMock()
	tcmock.SetErr = errors.NewAlreadyClosedf("cache already closed ;-(")

	ws := newMMWS(tcmock, "a", "b", http.DefaultClient)
	trip := cstesting.NewHttpTripFromFile(200, "testdata/response.json")
	ws.client.Transport = trip
	c, err := ws.Country(net.ParseIP("123.123.123.123"))
	assert.Nil(t, c)
	assert.True(t, errors.IsAlreadyClosed(err), "Error: %s", err)
}

func TestMmws_Country_Success(t *testing.T) {
	td, err := ioutil.ReadFile("testdata/response.json")
	if err != nil {
		t.Fatal(err)
	}

	tcmock := transcache.NewMock()
	ws := newMMWS(tcmock, "gopher", "passw0rd", http.DefaultClient)
	trip := cstesting.NewHttpTrip(200, string(td), nil)
	ws.client.Transport = trip

	const iterations = 10
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			c, err := ws.Country(net.ParseIP("123.123.123.123"))
			assert.NotNil(t, c)
			assert.NoError(t, err)
			assert.Exactly(t, "US", c.Country.IsoCode)
		}(&wg)
	}
	wg.Wait()

	assert.Exactly(t, 1, tcmock.SetCount(), "SetCount")
	assert.Exactly(t, iterations-1, tcmock.GetCount(), "GetCount")

	trip.RequestsMatchAll(t, func(r *http.Request) bool {
		u, p, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Exactly(t, "gopher", u)
		assert.Exactly(t, "passw0rd", p)
		return true
	})

}

var maxMindWebServiceClient string

// BenchmarkMaxMindWebServiceClient/Serial-4         	   50000	     25525 ns/op	    5612 B/op	     108 allocs/op
// BenchmarkMaxMindWebServiceClient/Parallel-4       	  100000	     18447 ns/op	    5652 B/op	     108 allocs/op
func BenchmarkMaxMindWebServiceClient(b *testing.B) {

	// transcache.NewMock has gob encoding

	wsc := newMMWS(transcache.NewMock(), "gopher", "passw0rd", &http.Client{
		Transport: cstesting.NewHttpTrip(200, `{ "continent": { "code": "EU", "geoname_id": 6255148, "names": { "de": "Europa", "en": "Europe", "ru": "Европа", "zh-CN": "欧洲" } }, "country": { "geoname_id": 2921044, "iso_code": "DE", "names": { "de": "Deutschland", "en": "Germany", "es": "Alemania", "fr": "Allemagne", "ja": "ドイツ連邦共和国", "pt-BR": "Alemanha", "ru": "Германия", "zh-CN": "德国" } }, "maxmind": { "queries_remaining": 54321 } }`, nil),
	})

	var checkCountry = func(b *testing.B, ip net.IP) {
		ret, err := wsc.Country(ip)
		if err != nil {
			b.Fatal(err)
		}
		var want string
		if maxMindWebServiceClient, want = ret.Country.IsoCode, "DE"; maxMindWebServiceClient != want {
			b.Fatalf("Have: %v Want: %v", maxMindWebServiceClient, want)
		}
	}
	ip := net.ParseIP("123.123.123.123")

	checkCountry(b, ip)

	b.Run("Serial", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			checkCountry(b, ip)
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				checkCountry(b, ip)
			}
		})
	})
}
