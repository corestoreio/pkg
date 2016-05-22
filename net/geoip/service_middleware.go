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

	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/util/errors"
)

// newContextCountryByIP searches the country for an IP address and puts the country
// into a new context.
func (s *Service) newContextCountryByIP(r *http.Request) (context.Context, *Country, error) {

	ip := httputil.GetRealIP(r)
	if ip == nil {
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.newContextCountryByIP.GetRemoteAddr", "err", errors.NotFound(errCannotGetRemoteAddr), "req", r)
		}
		return nil, nil, errors.NewNotFoundf(errCannotGetRemoteAddr)
	}

	c, err := s.GeoIP.Country(ip)
	if err != nil {
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.newContextCountryByIP.GeoIP.Country", "err", err, "remoteAddr", ip, "req", r)
		}
		return nil, nil, errors.NewFatal(err, "[geoip] getting country")
	}
	return WithContextCountry(r.Context(), c), c, nil
}

// WithCountryByIP is a simple middleware which detects the country via an IP
// address. With the detected country a new tree context.Context gets created.
// Use FromContextCountry() to extract the country or an error.
func (s *Service) WithCountryByIP() mw.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, _, err := s.newContextCountryByIP(r)
			if err != nil {
				ctx = withContextError(r.Context(), errors.Wrap(err, "[geoip] newContextCountryByIP"))
			}
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WithIsCountryAllowedByIP queries the AllowedCountries slice
// to retrieve a list of countries for a scope and then uses the function
// IsAllowedFunc to check if a country is allowed for an IP address.
// Use FromContextCountry() to extract the country or an error.
func (s *Service) WithIsCountryAllowedByIP() mw.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			requestedStore, err := store.FromContextRequestedStore(r.Context())
			if err != nil {
				err = errors.Wrap(err, "[geoip] FromContextProvider")
				h.ServeHTTP(w, r.WithContext(withContextError(r.Context(), err)))
				return
			}

			ctx, c, err := s.newContextCountryByIP(r)
			if err != nil {
				err = errors.Wrap(err, "[geoip] newContextCountryByIP")
				h.ServeHTTP(w, r.WithContext(withContextError(r.Context(), err)))
				return
			}

			// the scpCfg depends on how you have initialized the storeService during app boot.
			// requestedStore.Website.Config is the reason that all options only support
			// website scope and not group or store scope.
			scpCfg, err := s.configByScopedGetter(requestedStore)
			if err != nil {
				if s.defaultScopeCache.log.IsDebug() {
					s.defaultScopeCache.log.Debug("Service.WithCORS.configByScopedGetter", "err", err, "requestedStore", requestedStore, "req", r)
				}
				err = errors.Wrap(err, "[mwcors] ConfigByScopedGetter")
				h.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			if scpCfg.checkAllow(requestedStore, c, r) {
				if s.Log.IsDebug() {
					s.Log.Debug("geoip.WithIsCountryAllowedByIP.checkAllow.true", "requestedStore", requestedStore, "country", c)
				}
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			// access denied
			if s.Log.IsDebug() {
				s.Log.Debug("geoip.WithIsCountryAllowedByIP.checkAllow.false", "requestedStore", requestedStore, "country", c)
			}
			scpCfg.alternativeHandler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WithInitStoreByCountryIP initializes a store scope via the IP address which
// is bound to a country. todo(CS) IDEA
func (s *Service) WithInitStoreByCountryIP() mw.Middleware {
	return nil
}
