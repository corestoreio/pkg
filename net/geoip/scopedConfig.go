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

// scopedConfig private internal scoped based configuration
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
	IsAllowedFunc // func(s scope.Hash, c *Country, allowedCountries []string) error

	// AlternativeHandler if ip/country is denied we call this handler.
	AlternativeHandler mw.ErrorHandler
}

func newScopedConfig() *ScopedConfig {
	sc := &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(),
		IsAllowedFunc: func(_ scope.TypeID, c *Country, allowedCountries []string) error {
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

// IsValid a configuration for a scope is only then valid when the Key has been
// supplied, a non-nil signing method and a non-nil Verifier.
func (sc *ScopedConfig) IsValid() error {
	if sc.lastErr != nil {
		return errors.Wrap(sc.lastErr, "[geoip] scopedConfig.isValid as an lastErr")
	}

	if sc.ScopeHash == 0 || sc.IsAllowedFunc == nil ||
		sc.AlternativeHandler == nil {
		return errors.NewNotValidf(errScopedConfigNotValid, sc.ScopeHash, sc.IsAllowedFunc == nil, sc.AlternativeHandler == nil)
	}
	return nil
}

func (sc *ScopedConfig) checkAllow(s scope.TypeID, c *Country) error {
	if len(sc.AllowedCountries) == 0 {
		return nil
	}
	return sc.IsAllowedFunc(s, c, sc.AllowedCountries)
}
