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
	"net/http"
	"sync/atomic"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
)

// IsAllowedFunc checks in middleware WithIsCountryAllowedByIP if the country is
// allowed to process the request. The StringSlice contains a list of ISO
// country names fetched from the config.ScopedGetter. Return nil to indicate
// that the request can continue.
type IsAllowedFunc func(_ scope.TypeID, _ *Country, allowedCountries []string) error

// WithDefaultConfig applies the default GeoIP configuration settings based for
// a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- Alternative Handler: variable DefaultAlternativeHandler
//		- Logger black hole
//		- Check allow: If allowed countries are empty, all countries are allowed
func WithDefaultConfig(scopeIDs ...scope.TypeID) Option {
	return withDefaultConfig(scopeIDs...)
}

// WithAlternativeHandler sets for a scope the alternative handler
// if an IP address has been access denied.
// Only to be used with function WithIsCountryAllowedByIP()
func WithAlternativeHandler(altHndlr mw.ErrorHandler, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.AlternativeHandler = altHndlr
		return s.updateScopedConfig(sc)
	}
}

// WithAlternativeRedirect sets for a scope the error handler
// on a Service if an IP address has been access denied.
// Only to be used with function WithIsCountryAllowedByIP()
func WithAlternativeRedirect(urlStr string, code int, scopeIDs ...scope.TypeID) Option {
	return WithAlternativeHandler(func(_ error) http.Handler {
		return http.RedirectHandler(urlStr, code)
	}, scopeIDs...)
}

// WithCheckAllow sets your custom function which checks if the country of an IP
// address should access to granted, or the next middleware handler in the chain
// gets called.
// Only to be used with function WithIsCountryAllowedByIP()
func WithCheckAllow(f IsAllowedFunc, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.IsAllowedFunc = f
		return s.updateScopedConfig(sc)
	}
}

// WithAllowedCountryCodes sets a list of ISO countries to be validated against.
// Only to be used with function WithIsCountryAllowedByIP()
func WithAllowedCountryCodes(isoCountryCodes []string, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.AllowedCountries = isoCountryCodes
		return s.updateScopedConfig(sc)
	}
}

// WithCountryFinder applies a custom CountryRetriever. Sets the retriever atomically
// and only once.
func WithCountryFinder(cr Finder) Option {
	return func(s *Service) error {
		if s.isGeoIPLoaded() {
			if s.Log.IsDebug() {
				s.Log.Debug("geoip.WithGeoIP.geoIPDone", log.Int("done", 1))
			}
			return nil
		}
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		if s.geoIPLoaded == 0 {
			s.Finder = cr
			atomic.StoreUint32(&s.geoIPLoaded, 1)
			if s.Log.IsDebug() {
				s.Log.Debug("geoip.WithGeoIP.geoIPLoaded", log.Int("done", 0))
			}
		}
		return nil
	}
}
