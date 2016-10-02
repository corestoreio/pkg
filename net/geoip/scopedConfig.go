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
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// ScopedConfig scoped based configuration and should not be embedded into your
// own types. Call ScopedConfig.ScopeID to know to which scope this
// configuration has been bound to.
type ScopedConfig struct {
	scopedConfigGeneric

	// Disabled set to true to disable GeoIP Service for this scope.
	Disabled bool

	// AllowedCountries a slice which contains all allowed countries. An
	// incoming request for a scope checks if the country for an IP is contained
	// within this slice. Empty slice means that all countries are allowed. The
	// slice is owned by the callee.
	AllowedCountries []string
	// IsAllowedFunc checks in middleware WithIsCountryAllowedByIP if the country is
	// allowed to process the request.
	isAllowedFn IsAllowedFunc // func(s scope.Hash, c *Country, allowedCountries []string) error

	// AlternativeHandler if ip/country is denied we call this handler.
	AlternativeHandler mw.ErrorHandler
}

func newScopedConfig() *ScopedConfig {
	sc := &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(),
		isAllowedFn: func(_ scope.TypeID, c *Country, allowedCountries []string) error {
			for _, ac := range allowedCountries {
				if ac == c.Country.IsoCode { // case sensitive matching
					return nil
				}
			}
			return errors.NewUnauthorizedf(errUnAuthorizedCountry, c.Country.IsoCode, allowedCountries)
		},
		AlternativeHandler: DefaultAlternativeHandler,
	}
	return sc
}

// isValid a configuration for a scope is only then valid when the Key has been
// supplied, a non-nil signing method and a non-nil Verifier.
func (sc *ScopedConfig) isValid() error {
	if err := sc.isValidPreCheck(); err != nil {
		return errors.Wrap(err, "[cors] scopedConfig.isValid as an lastErr")
	}
	if sc.Disabled {
		return nil
	}
	if sc.isAllowedFn == nil || sc.AlternativeHandler == nil {
		return errors.NewNotValidf(errScopedConfigNotValid, sc.ScopeID, sc.isAllowedFn == nil, sc.AlternativeHandler == nil)
	}
	return nil
}

// IsAllowed checks if the country is allowed. An empty AllowedCountries fields
// allows all countries.
func (sc *ScopedConfig) IsAllowed(c *Country) error {
	// think about: either if no country has been set and allow to proceed or be
	// more strict and proceeding is not allowed except sea and air territories
	// ;-).
	if len(sc.AllowedCountries) == 0 {
		return nil
	}
	return sc.isAllowedFn(sc.ScopeID, c, sc.AllowedCountries)
}
