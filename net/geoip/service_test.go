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

package geoip_test

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func mustGetTestService(opts ...geoip.Option) *geoip.Service {
	maxMindDB := filepath.Join("testdata", "GeoIP2-Country-Test.mmdb")
	s, err := geoip.New(append(opts, geoip.WithGeoIP2Reader(maxMindDB))...)
	if err != nil {
		panic(err)
	}
	return s
}

func finalHandlerFinland(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipc, err := geoip.FromContextCountry(r.Context())
		assert.NotNil(t, ipc)
		assert.NoError(t, err)
		assert.Exactly(t, "2a02:d200::", ipc.IP.String())
		assert.Exactly(t, "FI", ipc.Country.IsoCode)
	}
}

func mustGetRequestFinland() *http.Request {
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Forwarded-For", "2a02:d200::")
	return req
}

func TestNewServiceErrorWithoutOptions(t *testing.T) {
	s, err := geoip.New()
	assert.Nil(t, s)
	assert.EqualError(t, err, "Please provide a GeoIP Reader.")
}

func TestNewServiceErrorWithAlternativeHandler(t *testing.T) {
	s, err := geoip.New(geoip.WithAlternativeHandler(scope.Absent, 314152, nil))
	assert.Nil(t, s)
	assert.True(t, errors.IsNotSupported(err), "Error: %s", err)
}

func TestNewServiceErrorWithGeoIP2Reader(t *testing.T) {
	s, err := geoip.New(geoip.WithGeoIP2Reader("Walhalla/GeoIP2-Country-Test.mmdb"))
	assert.Nil(t, s)
	assert.EqualError(t, err, "File Walhalla/GeoIP2-Country-Test.mmdb not found")
}

func deferClose(t *testing.T, c io.Closer) {
	assert.NoError(t, c.Close())
}

func TestNewServiceWithGeoIP2Reader(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)
	ip, _, err := net.ParseCIDR("2a02:d200::/29") // IP range for Finland

	assert.NoError(t, err)
	haveCty, err := s.GeoIP.Country(ip)
	assert.NoError(t, err)
	assert.Exactly(t, "FI", haveCty.Country.IsoCode)
}

type geoReaderMock struct{}

func (geoReaderMock) Country(ipAddress net.IP) (*geoip.Country, error) {
	return nil, errors.New("Failed to read country from MMDB")
}
func (geoReaderMock) Close() error { return nil }

func TestWithCountryByIPErrorGetCountryByIP(t *testing.T) {
	s := mustGetTestService()
	s.GeoIP = geoReaderMock{}
	defer deferClose(t, s.GeoIP)

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipc, err := geoip.FromContextCountry(r.Context())
		assert.Nil(t, ipc)
		assert.EqualError(t, err, "Failed to read country from MMDB")
	})

	countryHandler := s.WithCountryByIP()(finalHandler)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "2a02:d200::")
	countryHandler.ServeHTTP(rec, req)
}

func TestWithCountryByIPSuccess(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	countryHandler := s.WithCountryByIP()(finalHandlerFinland(t))
	rec := httptest.NewRecorder()

	countryHandler.ServeHTTP(rec, mustGetRequestFinland())
}

func TestWithIsCountryAllowedByIPErrorStoreManager(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipc, err := geoip.FromContextCountry(r.Context())
		assert.Nil(t, ipc)
		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	})

	countryHandler := s.WithIsCountryAllowedByIP()(finalHandler)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	countryHandler.ServeHTTP(rec, req)
}

func ipErrorFinalHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipc, err := geoip.FromContextCountry(r.Context())
		assert.Nil(t, ipc)
		assert.True(t, errors.IsFatal(err), "Error: %s", err)
	}
}

func TestWithCountryByIPErrorRemoteAddr(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	countryHandler := s.WithCountryByIP()(ipErrorFinalHandler(t))
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "2324.2334.432.534")
	countryHandler.ServeHTTP(rec, req)
}

func TestWithIsCountryAllowedByIPErrorWithContextCountryByIP(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	countryHandler := s.WithIsCountryAllowedByIP()(ipErrorFinalHandler(t))
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "2R02:d2'0.:")

	countryHandler.ServeHTTP(rec, req) // managerStoreSimpleTest,  ?
}

func TestWithIsCountryAllowedByIPErrorAllowedCountries(t *testing.T) {
	t.Skip("@todo once store package has been refactored")
}
