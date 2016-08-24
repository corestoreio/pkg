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
	"context"
	"net/http"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/net/request"
	"github.com/corestoreio/csfw/util/errors"
)

// CountryByIP searches a country by an IP address and returns the found
// country. It only needs the functional options WithGeoIP*().
func (s *Service) CountryByIP(r *http.Request) (*Country, error) {

	ip := request.RealIP(r, request.IPForwardedTrust) // todo make IPForwardedTrust configurable
	if ip == nil {
		nf := errors.NewNotFoundf(errCannotGetRemoteAddr)
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.newContextCountryByIP.GetRemoteAddr", log.Err(nf), log.HTTPRequest("request", r))
		}
		return nil, nf
	}

	c, err := s.geoIP.Country(ip)
	if err != nil {
		if s.Log.IsDebug() {
			s.Log.Debug(
				"geoip.Service.newContextCountryByIP.GeoIP.Country",
				log.Err(err), log.Stringer("remote_addr", ip), log.HTTPRequest("request", r))
		}
		return nil, errors.Wrap(err, "[geoip] getting country")
	}
	return c, nil
}

// newContextCountryByIP searches a country by an IP address and puts the country
// into a new context.
func (s *Service) newContextCountryByIP(r *http.Request) (context.Context, *Country, error) {
	c, err := s.CountryByIP(r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[geoip] CountryByIP")
	}
	return withContextCountry(r.Context(), c), c, nil
}

// WithCountryByIP is a simple middleware which detects the country via an IP
// address. With the detected country a new tree context.Context gets created.
// Use FromContextCountry() to extract the country or an error. If you don't
// like the middleware consider using the function CountryByIP().
func (s *Service) WithCountryByIP() mw.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, _, err := s.newContextCountryByIP(r)
			if err != nil {
				s.ErrorHandler(errors.Wrap(err, "[geoip] newContextCountryByIP")).ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

// WithIsCountryAllowedByIP queries the AllowedCountries slice to retrieve a
// list of countries for a scope and then uses the function IsAllowedFunc to
// check if a country is allowed for an IP address. If a country should not
// access the next handler within the middleware chain it will call an
// alternative handler to e.g. show a different page or perform a redirect. Use
// FromContextCountry() to extract the country or an error. Tis middleware
// allows geo blocking.
func (s *Service) WithIsCountryAllowedByIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		scpCfg := s.configByContext(r.Context())
		if err := scpCfg.IsValid(); err != nil {
			s.Log.Info("geoip.Service.WithIsCountryAllowedByIP.configByContext.Error", log.Err(err))
			if s.Log.IsDebug() {
				s.Log.Debug("geoip.Service.WithIsCountryAllowedByIP.configByContext", log.Err(err), log.HTTPRequest("request", r))
			}
			s.ErrorHandler(errors.Wrap(err, "geoip.Service.WithIsCountryAllowedByIP.configFromContext")).ServeHTTP(w, r)
			return
		}
		if scpCfg.Disabled {
			if s.Log.IsDebug() {
				s.Log.Debug("geoip.Service.WithIsCountryAllowedByIP.Disabled", log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			next.ServeHTTP(w, r)
			return
		}

		ctx, c, err := s.newContextCountryByIP(r)
		if err != nil {
			err = errors.Wrap(err, "[geoip] newContextCountryByIP")
			scpCfg.ErrorHandler(err).ServeHTTP(w, r)
			return
		}

		if err := scpCfg.checkAllow(scpCfg.ScopeHash, c); err != nil {
			// access denied
			if s.Log.IsDebug() {
				s.Log.Debug("geoip.WithIsCountryAllowedByIP.checkAllow.false", log.Err(err), log.Stringer("scope", scpCfg.ScopeHash), log.String("countryISO", c.Country.IsoCode), log.Strings("allowedCountries", scpCfg.AllowedCountries...))
			}
			err = errors.Wrap(err, "[geoip] WithIsCountryAllowedByIP.CheckAllow")
			scpCfg.AlternativeHandler(err).ServeHTTP(w, r)
			return
		}

		// access granted
		if s.Log.IsDebug() {
			s.Log.Debug("Service.WithIsCountryAllowedByIP.checkAllow.true", log.Stringer("scope", scpCfg.ScopeHash), log.String("countryISO", c.Country.IsoCode), log.Strings("allowedCountries", scpCfg.AllowedCountries...))
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithInitStoreByCountryIP initializes a store scope via the IP address which
// is bound to a country. todo(CS) idea
func (s *Service) WithInitStoreByCountryIP() mw.Middleware {
	// - define a mapping for a store assigned to countries ISO codes
	// - load that store default but allow a user to switch
	// - force set a store to a country and the user cannot switch.
	return nil
}
