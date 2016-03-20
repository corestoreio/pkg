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
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"io"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/store/storenet"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var _ error = (*geoip.Service)(nil)

func mustGetTestService(opts ...geoip.Option) *geoip.Service {
	maxMindDB := filepath.Join(cstesting.RootPath, "net", "geoip", "GeoIP2-Country-Test.mmdb")
	s, err := geoip.NewService(append(opts, geoip.WithGeoIP2Reader(maxMindDB))...)
	if err != nil {
		panic(err)
	}
	return s
}

func finalHandlerFinland(t *testing.T) ctxhttp.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ipc, err, ok := geoip.FromContextCountry(ctx)
		assert.NotNil(t, ipc)
		assert.True(t, ok)
		assert.NoError(t, err)
		assert.Exactly(t, "2a02:d200::", ipc.IP.String())
		assert.Exactly(t, "FI", ipc.Country.IsoCode)
		return nil
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
	s, err := geoip.NewService()
	assert.Nil(t, s)
	assert.EqualError(t, err, "Please provide a GeoIP Reader.")
}

func TestNewServiceErrorWithAlternativeHandler(t *testing.T) {
	s, err := geoip.NewService(geoip.WithAlternativeHandler(scope.AbsentID, 314152, nil))
	assert.Nil(t, s)
	assert.EqualError(t, err, scope.ErrUnsupportedScopeID.Error())
}

func TestNewServiceErrorWithGeoIP2Reader(t *testing.T) {
	s, err := geoip.NewService(geoip.WithGeoIP2Reader("Walhalla/GeoIP2-Country-Test.mmdb"))
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

	finalHandler := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ipc, err, ok := geoip.FromContextCountry(ctx)
		assert.Nil(t, ipc)
		assert.True(t, ok)
		assert.EqualError(t, err, "Failed to read country from MMDB")
		return nil
	})

	countryHandler := s.WithCountryByIP()(finalHandler)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "2a02:d200::")
	assert.NoError(t, countryHandler.ServeHTTPContext(context.Background(), rec, req))
}

func TestWithCountryByIPSuccess(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	countryHandler := s.WithCountryByIP()(finalHandlerFinland(t))
	rec := httptest.NewRecorder()

	assert.NoError(t, countryHandler.ServeHTTPContext(context.Background(), rec, mustGetRequestFinland()))
}

func TestWithIsCountryAllowedByIPErrorStoreManager(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	finalHandler := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ipc, err, ok := geoip.FromContextCountry(ctx)
		assert.Nil(t, ipc)
		assert.False(t, ok)
		assert.NoError(t, err)
		return nil
	})

	countryHandler := s.WithIsCountryAllowedByIP()(finalHandler)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	assert.EqualError(t, countryHandler.ServeHTTPContext(context.Background(), rec, req), storenet.ErrContextProviderNotFound.Error())
}

var managerStoreSimpleTest = storemock.WithContextMustService(scope.Option{}, func(ms *storemock.Storage) {
	ms.MockStore = func() (*store.Store, error) {
		return store.NewStore(
			&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		)
	}
})

func ipErrorFinalHandler(t *testing.T) ctxhttp.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ipc, err, ok := geoip.FromContextCountry(ctx)
		assert.Nil(t, ipc)
		assert.True(t, ok)
		assert.EqualError(t, err, geoip.ErrCannotGetRemoteAddr.Error())
		return nil
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
	assert.NoError(t, countryHandler.ServeHTTPContext(context.Background(), rec, req))
}

func TestWithIsCountryAllowedByIPErrorWithContextCountryByIP(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	countryHandler := s.WithIsCountryAllowedByIP()(ipErrorFinalHandler(t))
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "2R02:d2'0.:")

	assert.NoError(t, countryHandler.ServeHTTPContext(managerStoreSimpleTest, rec, req))
}

func TestWithIsCountryAllowedByIPErrorAllowedCountries(t *testing.T) {
	t.Skip("@todo once store package has been refactored")
}
