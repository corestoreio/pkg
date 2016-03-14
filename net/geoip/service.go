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

	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storenet"
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
	"golang.org/x/net/context"
)

// ErrCannotGetRemoteAddr will be returned if there is an invalid or not found
// RemoteAddr in the request.
var ErrCannotGetRemoteAddr = errors.New("Cannot get request.RemoteAddr")

// Service represents a service manager
type Service struct {
	*cserr.MultiErr
	// AllowedCountries a model containing a path to the configuration which
	// countries are allowed within a scope. Current implementation triggers
	// for each HTTP request a configuration lookup which can be a bottle neck.
	AllowedCountries cfgmodel.StringCSV
	// GeoIP searches the country for an IP address
	GeoIP Reader
	// IsAllowed checks in middleware WithIsCountryAllowedByIP if the country is
	// allowed to process the request.
	IsAllowed IsAllowedFunc
	// AltH are alternative handlers if the current request is not allowed for
	// a country.
	// IDs and AltH slices must have both the same length because with the ID
	// found in IDs slice we take the index key and access the appropriate handler in AltH.
	websiteIDs  util.Int64Slice
	websiteAltH []ctxhttp.HandlerFunc
	storeIDs    util.Int64Slice
	storeAltH   []ctxhttp.HandlerFunc
}

// NewService creates a new GeoIP service to be used as a middleware.
func NewService(opts ...Option) (*Service, error) {
	s := new(Service)
	for _, opt := range opts {
		opt(s)
	}
	if s.HasErrors() {
		return nil, s
	}
	if s.GeoIP == nil {
		return nil, errors.New("Please provide a GeoIP Reader.")
	}
	if s.IsAllowed == nil {
		s.IsAllowed = defaultIsCountryAllowed
	}
	return s, nil
}

// newContextCountryByIP searches the country for an IP address and puts the country
// into a new context.
func (s *Service) newContextCountryByIP(ctx context.Context, r *http.Request) (context.Context, *Country, error) {

	ip := httputil.GetRemoteAddr(r)
	if ip == nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("geoip.WithCountryByIP.GetRemoteAddr", "err", ErrCannotGetRemoteAddr, "req", r)
		}
		return ctx, nil, ErrCannotGetRemoteAddr
	}

	c, err := s.GeoIP.Country(ip)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("geoip.WithCountryByIP.GeoIP.Country", "err", err, "remoteAddr", ip, "req", r)
		}
		return ctx, nil, errors.Mask(err)
	}
	return WithContextCountry(ctx, c), c, nil
}

// WithCountryByIP is a simple middleware which detects the country via an IP
// address. With the detected country a new tree context.Context gets created.
// Use FromContextCountry() to extract the country or an error.
func (s *Service) WithCountryByIP() ctxhttp.Middleware {
	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var err error
			ctx, _, err = s.newContextCountryByIP(ctx, r)
			if err != nil {
				ctx = WithContextError(ctx, err)
			}
			return hf(ctx, w, r)
		}
	}
}

// WithIsCountryAllowedByIP queries the AllowedCountries configuration model
// to retrieve a list of countries for a scope and then uses the function
// IsAllowedFunc to check if a country is allowed for an IP address.
// Use FromContextCountry() to extract the country or an error.
func (s *Service) WithIsCountryAllowedByIP() ctxhttp.Middleware {
	return func(h ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			_, requestedStore, err := storenet.FromContextProvider(ctx)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("geoip.WithCountryByIP.FromContextManagerReader", "err", err)
				}
				return errors.Mask(err)
			}

			var c *Country
			ctx, c, err = s.newContextCountryByIP(ctx, r)
			if err != nil {
				ctx = WithContextError(ctx, err)
				return h(ctx, w, r)
			}

			allowedCountries, err := s.AllowedCountries.Get(requestedStore.Config)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("geoip.WithCountryByIP.directory.AllowedCountries", "err", err, "st.Config", requestedStore.Config)
				}
				return errors.Mask(err)
			}

			if false == s.IsAllowed(requestedStore, c, allowedCountries, r) {
				h = s.altHandlerByID(requestedStore)
			}

			return h(ctx, w, r)
		}
	}
}

// DefaultAlternativeHandler gets called when detected Country cannot be found
// within the list of allowed countries. This handler can be overridden to provide
// a fallback for all scopes. To set a alternative handler for a website or store
// use the With*() options. This function gets called in WithIsCountryAllowedByIP.
//
// Status is StatusServiceUnavailable
var DefaultAlternativeHandler ctxhttp.HandlerFunc = defaultAlternativeHandler

var defaultAlternativeHandler = func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	return nil
}

// altHandlerByID searches in the hierarchical order of store -> website -> default.
// the next alternative handler IF a country is not allowed as defined in function
// type IsAllowedFunc.
func (s *Service) altHandlerByID(st *store.Store) ctxhttp.HandlerFunc {

	if s.storeIDs != nil && s.storeAltH != nil {
		return findHandlerByID(scope.StoreID, st.StoreID(), s.storeIDs, s.storeAltH)
	}
	if s.websiteIDs != nil && s.websiteAltH != nil {
		return findHandlerByID(scope.WebsiteID, st.Website.WebsiteID(), s.websiteIDs, s.websiteAltH)
	}
	return DefaultAlternativeHandler
}

// findHandlerByID returns the Handler for the searchID. If not found
// or slices have an indifferent length or something is nil it will
// return the DefaultErrorHandler.
func findHandlerByID(so scope.Scope, id int64, idsIdx util.Int64Slice, handlers []ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {

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

	if PkgLog.IsInfo() {
		PkgLog.Info("geoip.findHandlerByID.found", "scope", so.String(), "id", id, "idsIdx", idsIdx)
	}
	return prospect
}

func defaultIsCountryAllowed(_ *store.Store, c *Country, allowedCountries []string, r *http.Request) bool {
	var ac util.StringSlice = allowedCountries
	if false == ac.Contains(c.Country.IsoCode) {
		if PkgLog.IsInfo() {
			PkgLog.Info("geoip.checkAllow", "Country", c, "allowedCountries", allowedCountries, "request", r)
		}
		return false
	}
	return true
}
