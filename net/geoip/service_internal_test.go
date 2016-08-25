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
	"bytes"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ io.Closer = (*Service)(nil)

func deferClose(t *testing.T, c io.Closer) {
	assert.NoError(t, c.Close())
}

func mustGetTestService(opts ...Option) *Service {
	maxMindDB := filepath.Join("testdata", "GeoIP2-Country-Test.mmdb")
	return MustNew(append(opts, WithGeoIP2File(maxMindDB))...)
}

func TestMustNew(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.IsNotFound(err), "Error: %s", err)
			} else {
				t.Fatal("Expecting an error")
			}
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	s := MustNew(WithGeoIP2File("not found"))
	assert.Nil(t, s)
}

func TestNewServiceErrorWithoutOptions(t *testing.T) {
	s, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Nil(t, s.geoIP)
}

func TestNewService_WithGeoIP2File_Atomic(t *testing.T) {
	logBuf := &bytes.Buffer{}
	s, err := New(
		WithLogger(logw.NewLog(logw.WithWriter(logBuf), logw.WithLevel(logw.LevelDebug))),
		WithGeoIP2File(filepath.Join("testdata", "GeoIP2-Country-Test.mmdb")),
	)
	defer deferClose(t, s)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, s.geoIP)
	for i := 0; i < 3; i++ {
		if err := s.Options(WithGeoIP2File(filepath.Join("testdata", "GeoIP2-Country-Test.mmdb"))); err != nil {
			t.Fatal(err)
		}
	}
	assert.True(t, 3 == strings.Count(logBuf.String(), `geoip.WithGeoIP.geoIPDone done: 1`), logBuf.String())
}

func TestNewService_WithGeoIP2Webservice_Atomic(t *testing.T) {
	logBuf := &bytes.Buffer{}
	s, err := New(
		WithLogger(logw.NewLog(logw.WithWriter(logBuf), logw.WithLevel(logw.LevelDebug))),
		WithGeoIP2Webservice(nil, "a", "b", 1),
	)
	defer deferClose(t, s)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, s.geoIP)
	for i := 0; i < 3; i++ {
		assert.NoError(t, s.Options(WithGeoIP2Webservice(nil, "d", "e", 1)))
	}
	assert.True(t, 3 == strings.Count(logBuf.String(), `WithGeoIP.geoIPDone done: 1`), logBuf.String())
}

func TestNewServiceErrorWithGeoIP2Reader(t *testing.T) {
	s, err := New(WithGeoIP2File("Walhalla/GeoIP2-Country-Test.mmdb"))
	assert.Nil(t, s)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
}

func TestNewServiceWithGeoIP2Reader(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s)
	ip, _, err := net.ParseCIDR("2a02:d200::/29") // IP range for Finland

	assert.NoError(t, err)
	haveCty, err := s.geoIP.Country(ip)
	assert.NoError(t, err)
	assert.Exactly(t, "FI", haveCty.Country.IsoCode)
}

func TestNewServiceWithCheckAllow(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s)

	req, _ := http.NewRequest("GET", "http://corestore.io", nil)
	req.Header.Set("Forwarded-For", "2a02:d200::") // IP Range Finland

	t.Run("Scope_Default", func(t *testing.T) {

		if err := s.Options(WithAllowedCountryCodes(scope.Default, 0, "US")); err != nil {
			t.Fatal(err)
		}

		scpCfg := s.ConfigByScopeHash(scope.DefaultHash, 0)
		if err := scpCfg.IsValid(); err != nil {
			t.Fatal(err)
		}

		c, err := s.CountryByIP(req)
		if err != nil {
			t.Fatal(err)
		}
		haveErr := scpCfg.checkAllow(0, c)
		assert.True(t, errors.IsUnauthorized(haveErr), "Error: %s", haveErr)
	})

	t.Run("Scope_Store", func(t *testing.T) {
		if err := s.Options(WithCheckAllow(scope.Store, 331122, func(s scope.Hash, c *Country, allowedCountries []string) error {
			assert.Exactly(t, scope.Hash(0), s, "scope.Hash")
			assert.Exactly(t, "FI", c.Country.IsoCode)
			assert.Exactly(t, []string{"ABC"}, allowedCountries)
			return errors.NewNotImplementedf("You're not allowed")
		}), WithAllowedCountryCodes(scope.Store, 331122, "ABC")); err != nil {
			t.Fatal(err)
		}

		scpCfg := s.ConfigByScopeHash(scope.NewHash(scope.Store, 331122), 0)
		if err := scpCfg.IsValid(); err != nil {
			t.Fatal(err)
		}

		c, err := s.CountryByIP(req)
		if err != nil {
			t.Fatal(err)
		}
		haveErr := scpCfg.checkAllow(0, c)
		assert.True(t, errors.IsNotImplemented(haveErr), "Error: %s", haveErr)
	})
}
