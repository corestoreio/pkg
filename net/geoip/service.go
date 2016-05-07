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

	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
)

// Service represents a service manager
type Service struct {
	// AllowedCountries a model containing a path to the configuration which
	// countries are allowed within a scope. Current implementation triggers
	// for each HTTP request a configuration lookup which can be a bottle neck.
	AllowedCountries cfgmodel.StringCSV
	// GeoIP searches the country for an IP address
	GeoIP Reader
	// IsAllowed checks in middleware WithIsCountryAllowedByIP if the country is
	// allowed to process the request.
	IsAllowed IsAllowedFunc

	// optionError use by functional option arguments to indicate that one
	// option has triggered an error and hence the other can options can
	// skip their process.
	optionError error

	// AltH are alternative handlers if the current request is not allowed for
	// a country.
	// IDs and AltH slices must have both the same length because with the ID
	// found in IDs slice we take the index key and access the appropriate handler in AltH.
	websiteIDs  util.Int64Slice
	websiteAltH []http.Handler
	storeIDs    util.Int64Slice
	storeAltH   []http.Handler
	// Log used for debugging. Defaults to black hole. Panics if nil.
	Log log.Logger
}

// NewService creates a new GeoIP service to be used as a middleware.
func NewService(opts ...Option) (*Service, error) {
	s := &Service{
		Log: log.BlackHole{}, // disabled info and debug logging
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.optionError != nil {
		return nil, s.optionError
	}
	if s.GeoIP == nil {
		return nil, errors.NewFatalf("[geoip] Missing GeoIP Reader")
	}
	if s.IsAllowed == nil {
		s.IsAllowed = defaultIsCountryAllowed
	}
	return s, nil
}

// newContextCountryByIP searches the country for an IP address and puts the country
// into a new context.
func (s *Service) newContextCountryByIP(r *http.Request) (context.Context, *Country, error) {

	ip := httputil.GetRemoteAddr(r)
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
				ctx = withContextError(r.Context(), err)
			}
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WithIsCountryAllowedByIP queries the AllowedCountries configuration model
// to retrieve a list of countries for a scope and then uses the function
// IsAllowedFunc to check if a country is allowed for an IP address.
// Use FromContextCountry() to extract the country or an error.
func (s *Service) WithIsCountryAllowedByIP() mw.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			_, requestedStore, err := store.FromContextProvider(r.Context())
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

			allowedCountries, err := s.AllowedCountries.Get(requestedStore.Config)
			if err != nil {
				if s.Log.IsDebug() {
					s.Log.Debug("geoip.WithIsCountryAllowedByIP.AllowedCountries", "err", err, "requestedStore", requestedStore, "country", c)
				}
				err = errors.NewFatal(err, "[geoip] AllowedCountries.Get")
				h.ServeHTTP(w, r.WithContext(withContextError(r.Context(), err)))
				return
			}

			if false == s.IsAllowed(requestedStore, c, allowedCountries, r) {
				h = s.altHandlerByID(requestedStore)
			}

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WithInitStoreByCountryIP initializes a store scope via the IP address which
// is bound to a country. todo(CS) IDEA
func (s *Service) WithInitStoreByCountryIP() mw.Middleware {
	return nil
}

// altHandlerByID searches in the hierarchical order of store -> website -> default.
// the next alternative handler IF a country is not allowed as defined in function
// type IsAllowedFunc.
func (s *Service) altHandlerByID(st *store.Store) http.Handler {

	if s.storeIDs != nil && s.storeAltH != nil {
		return findHandlerByID(scope.Store, st.StoreID(), s.storeIDs, s.storeAltH)
	}
	if s.websiteIDs != nil && s.websiteAltH != nil {
		return findHandlerByID(scope.Website, st.Website.WebsiteID(), s.websiteIDs, s.websiteAltH)
	}
	return DefaultAlternativeHandler
}

// findHandlerByID returns the Handler for the searchID. If not found
// or slices have an indifferent length or something is nil it will
// return the DefaultErrorHandler.
func findHandlerByID(so scope.Scope, id int64, idsIdx util.Int64Slice, handlers []http.Handler) http.Handler {

	if len(idsIdx) != len(handlers) {
		return DefaultAlternativeHandler
	}
	index := idsIdx.Index(id)
	if index < 0 {
		return DefaultAlternativeHandler
	}
	prospect := handlers[index]
	if nil == prospect {
		return DefaultAlternativeHandler
	}

	return prospect
}

func defaultIsCountryAllowed(_ *store.Store, c *Country, allowedCountries []string, r *http.Request) bool {
	var ac util.StringSlice = allowedCountries
	if false == ac.Contains(c.Country.IsoCode) {
		return false
	}
	return true
}
