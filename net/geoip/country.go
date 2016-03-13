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
	"net/http"

	"github.com/corestoreio/csfw/store"
)

// The Country structure corresponds to the data in the GeoIP2/GeoLite2
// Country databases.
type Country struct {
	// IP contains the request IP address even if we run behind a proxy
	IP        net.IP
	Continent struct {
		Code      string
		GeoNameID uint
		Names     map[string]string
	}
	Country struct {
		GeoNameID uint
		IsoCode   string
		Names     map[string]string
	}
	RegisteredCountry struct {
		GeoNameID uint
		IsoCode   string
		Names     map[string]string
	}
	RepresentedCountry struct {
		GeoNameID uint
		IsoCode   string
		Names     map[string]string
		Type      string
	}
	Traits struct {
		IsAnonymousProxy    bool
		IsSatelliteProvider bool
	}
}

// Reader defines the functions which are needed to return a country by an IP.
type Reader interface {
	Country(net.IP) (*Country, error)
	// Close may be called on shutdown of the overall app.
	Close() error
}

// IsAllowedFunc checks in middleware WithIsCountryAllowedByIP if the country is
// allowed to process the request. The StringSlice contains a list of ISO country
// names fetched from the config.ScopedGetter.
type IsAllowedFunc func(s *store.Store, c *Country, allowedCountries []string, r *http.Request) bool
