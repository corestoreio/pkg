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
	"bytes"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corestoreio/pkg/net/geoip"
	"github.com/corestoreio/pkg/net/geoip/maxmindfile"
	"github.com/corestoreio/pkg/net/geoip/maxmindwebservice"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log/logw"
	"github.com/corestoreio/pkg/util/assert"
)

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
	s := geoip.MustNew(maxmindfile.WithCountryFinder("not found"))
	assert.Nil(t, s)
}

func TestNewServiceErrorWithoutOptions(t *testing.T) {
	s, err := geoip.New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Nil(t, s.Finder)
}

func TestNewService_WithGeoIP2File_Atomic(t *testing.T) {
	logBuf := &bytes.Buffer{}
	s, err := geoip.New(
		geoip.WithLogger(logw.NewLog(logw.WithWriter(logBuf), logw.WithLevel(logw.LevelDebug))),
		maxmindfile.WithCountryFinder(filepath.Join("testdata", "GeoIP2-Country-Test.mmdb")),
	)
	defer func() {
		if err := s.Close(); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
	}()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, s.Finder)
	for i := 0; i < 3; i++ {
		if err := s.Options(maxmindfile.WithCountryFinder(filepath.Join("testdata", "GeoIP2-Country-Test.mmdb"))); err != nil {
			t.Fatal(err)
		}
	}
	assert.True(t, 3 == strings.Count(logBuf.String(), `geoip.WithGeoIP.geoIPDone done: 1`), logBuf.String())
}

func TestNewService_WithGeoIP2Webservice_Atomic(t *testing.T) {
	logBuf := &bytes.Buffer{}
	s, err := geoip.New(
		geoip.WithDebugLog(logBuf),
		maxmindwebservice.WithCountryFinder(nil, "a", "b", 1),
	)
	defer func() {
		if err := s.Close(); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
	}()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, s.Finder)
	for i := 0; i < 3; i++ {
		assert.NoError(t, s.Options(maxmindwebservice.WithCountryFinder(nil, "d", "e", 1)))
	}
	assert.True(t, 3 == strings.Count(logBuf.String(), `WithGeoIP.geoIPDone done: 1`), logBuf.String())
}

func TestNewServiceErrorWithGeoIP2Reader(t *testing.T) {
	s, err := geoip.New(maxmindfile.WithCountryFinder("Walhalla/GeoIP2-Country-Test.mmdb"))
	assert.Nil(t, s)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
}

func TestNewServiceWithGeoIP2Reader(t *testing.T) {
	s, closeFn := mustGetTestService()
	defer closeFn()

	ip, _, err := net.ParseCIDR("2a02:d200::/29") // IP range for Finland

	assert.NoError(t, err)
	haveCty, err := s.Finder.FindCountry(ip)
	assert.NoError(t, err)
	assert.Exactly(t, "FI", haveCty.Country.IsoCode)
}

func TestNewServiceWithCheckAllow(t *testing.T) {
	s, closeFn := mustGetTestService()
	defer closeFn()

	req, _ := http.NewRequest("GET", "http://corestore.io", nil)
	req.Header.Set("Forwarded-For", "2a02:d200::") // IP Range Finland

	t.Run("Scope_Default", func(t *testing.T) {

		if err := s.Options(geoip.WithAllowedCountryCodes([]string{"US"})); err != nil {
			t.Fatal(err)
		}

		scpCfg, err := s.ConfigByScopeID(scope.DefaultTypeID, 0)
		if err != nil {
			t.Fatal(err)
		}

		c, err := s.CountryByIP(req)
		if err != nil {
			t.Fatal(err)
		}
		haveErr := scpCfg.IsAllowed(c)
		assert.True(t, errors.IsUnauthorized(haveErr), "Error: %s", haveErr)
	})

	t.Run("Scope_Store", func(t *testing.T) {
		var scopeID = scope.Store.WithID(331122)

		isAllowed := func(s scope.TypeID, c *geoip.Country, allowedCountries []string) error {
			assert.Exactly(t, scopeID, s, "Scope_Store @ scope.Hash")
			assert.Exactly(t, "FI", c.Country.IsoCode)
			assert.Exactly(t, []string{"ABC"}, allowedCountries)
			return errors.NewNotImplementedf("You're not allowed")
		}

		if err := s.Options(
			geoip.WithCheckAllow(isAllowed, scopeID),
			geoip.WithAllowedCountryCodes([]string{"ABC"}, scopeID),
		); err != nil {
			t.Fatalf("%+v", err)
		}

		scpCfg, err := s.ConfigByScopeID(scopeID, 0)
		if err != nil {
			t.Fatal(err)
		}

		c, err := s.CountryByIP(req)
		if err != nil {
			t.Fatal(err)
		}
		haveErr := scpCfg.IsAllowed(c)
		assert.True(t, errors.IsNotImplemented(haveErr), "Error: %s", haveErr)
	})
}
