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
	"time"

	"os"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
)

// Option can be used as an argument in NewService to configure it with
// different settings.
type Option func(*Service) error

// ScopedOptionFunc a closure around a scoped configuration to figure out which
// options should be returned depending on the scope brought to you during
// a request.
type ScopedOptionFunc func(config.ScopedGetter) []Option

// WithDefaultConfig applies the default GeoIP configuration settings based for
// a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- Alternative Handler: variable DefaultAlternativeHandler
//		- Logger black hole
//		- Check allow:
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		var err error
		if h == scope.DefaultHash {
			s.defaultScopeCache, err = defaultScopedConfig()
			return errors.Wrap(err, "[geoip] Default Scope with Default Config")
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		s.scopeCache[h], err = defaultScopedConfig()
		return errors.Wrapf(err, "[geoip] Scope %s with Default Config", h)
	}
}

// WithAlternativeHandler sets for a scope the alternative handler
// if an IP address has been access denied.
// Only to be used with function WithIsCountryAllowedByIP()
func WithAlternativeHandler(scp scope.Scope, id int64, altHndlr http.Handler) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.alternativeHandler = altHndlr
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.alternativeHandler = altHndlr

		if sc, ok := s.scopeCache[h]; ok {
			sc.alternativeHandler = scNew.alternativeHandler
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithAlternativeRedirect sets for a scope the error handler
// on a Service if an IP address has been access denied.
// Only to be used with function WithIsCountryAllowedByIP()
func WithAlternativeRedirect(scp scope.Scope, id int64, urlStr string, code int) Option {
	return WithAlternativeHandler(scp, id, http.RedirectHandler(urlStr, code))
}

// WithCheckAllow sets your custom function which checks if the country of an IP
// address should access to granted, or the next middleware handler in the chain
// gets called.
// Only to be used with function WithIsCountryAllowedByIP()
func WithCheckAllow(scp scope.Scope, id int64, f IsAllowedFunc) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.IsAllowedFunc = f
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.IsAllowedFunc = f

		if sc, ok := s.scopeCache[h]; ok {
			sc.IsAllowedFunc = scNew.IsAllowedFunc
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithAllowedCountryCodes sets a list of ISO countries to be validated against.
// Only to be used with function WithIsCountryAllowedByIP()
func WithAllowedCountryCodes(scp scope.Scope, id int64, isoCountryCodes ...string) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.allowedCountries = isoCountryCodes
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.allowedCountries = isoCountryCodes

		if sc, ok := s.scopeCache[h]; ok {
			sc.allowedCountries = scNew.allowedCountries
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithLogger applies a logger to the default scope which gets inherited to
// subsequent scopes.
// Mainly used for debugging.
func WithLogger(l log.Logger) Option {
	return func(s *Service) error {
		s.defaultScopeCache.log = l
		s.Log = l
		return nil
	}
}

// WithGeoIP2Reader creates a new GeoIP2.Reader. As long as there are no other
// readers this is a mandatory argument.
// Error behaviour: NotFound, NotValid
func WithGeoIP2File(filename string) Option {
	return func(s *Service) error {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return errors.NewNotFoundf("[geoip] File %s not found", filename)
		}

		var err error
		s.GeoIP, err = newMMDBByFile(filename)
		return errors.NewNotValid(err, "[geoip] Maxmind Open")
	}
}

// WithGeoIP2WebService uses for each incoming a request a lookup request to
// the Maxmind Webservice http://dev.maxmind.com/geoip/geoip2/web-services/
// and caches the result in Transcacher.
// Hint: use package storage/transcache.
func WithGeoIP2Webservice(t TransCacher, userID, licenseKey string, httpTimeout time.Duration) Option {
	return func(s *Service) error {
		s.GeoIP = newMMWS(t, userID, licenseKey, httpTimeout)
		return nil
	}
}

// WithOptionFactory applies a function which lazily loads the option depending
// on the incoming scope within a request. For example applies the backend
// configuration to the service.
// Once this option function has been set all other option functions are not really
// needed.
//	cfgStruct, err := backendgeoip.NewConfigStructure()
//	if err != nil {
//		panic(err)
//	}
//	pb := backendgeoip.New(cfgStruct)
//
//	cors := geoip.MustNewService(
//		geoip.WithOptionFactory(backendgeoip.PrepareOptions(pb)),
//	)
func WithOptionFactory(f ScopedOptionFunc) Option {
	return func(s *Service) error {
		s.scopedOptionFunc = f
		return nil
	}
}
