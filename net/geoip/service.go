// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"golang.org/x/net/context"

	"errors"
	"net"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/directory"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/httputils"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/utils"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/oschwald/geoip2-golang"
)

// ErrCannotGetRemoteAddr will be returned if there is an invalid or not found
// RemoteAddr in the request.
var ErrCannotGetRemoteAddr = errors.New("Cannot get request.RemoteAddr")

// ErrCannotGetStoreManagerReader will be returned when the store.ManagerReader
// cannot be found in a context.
var ErrCannotGetStoreManagerReader = errors.New("Cannot get store.ManagerReader from context.Context")

// IPCountry contains the found country and the IP address
type IPCountry struct {
	*geoip2.Country
	IP net.IP
}

// Reader defines the functions which are needed to return a country by an IP.
type Reader interface {
	Country(net.IP) (*geoip2.Country, error)
	Close() error
}

// IsAllowedFunc checks in middleware WithIsCountryAllowedByIP if the country is
// allowed to process the request. The StringSlice contains a list of ISO country
// names fetched from the config.Reader for a specific store view including fallback.
type IsAllowedFunc func(*store.Store, *IPCountry, utils.StringSlice, *http.Request) bool

type Service struct {
	// GeoIP searches the country for an IP address
	GeoIP Reader
	// IsAllowed checks in middleware WithIsCountryAllowedByIP if the country is
	// allowed to process the request.
	IsAllowed  IsAllowedFunc
	lastErrors []error
	// IDs and AltH slices must have both the same length because with the ID found in IDs slice
	// we take the index key and access the appropriate handler in AltH.
	websiteIDs  utils.Int64Slice
	websiteAltH []ctxhttp.Handler
	groupIDs    utils.Int64Slice
	groupAltH   []ctxhttp.Handler
	storeIDs    utils.Int64Slice
	storeAltH   []ctxhttp.Handler
}

// NewService creates a new GeoIP service to be used as a middleware.
func NewService(opts ...Option) (*Service, error) {
	s := new(Service)
	for _, opt := range opts {
		opt(s)
	}
	if s.lastErrors != nil {
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

var _ error = (*Service)(nil)

// Error returns an error string
func (s *Service) Error() string {
	return utils.Errors(s.lastErrors...)
}

func (s *Service) GetCountryByIP(ip net.IP) (*IPCountry, error) {
	// todo maybe add caching layer
	c, err := s.GeoIP.Country(ip)
	return &IPCountry{
		Country: c,
	}, err
}

// newContextCountryByIP searches the country for an IP address and puts the country
// into a new context.
func (s *Service) newContextCountryByIP(ctx context.Context, r *http.Request) (context.Context, *IPCountry, error) {

	remoteAddr := httputils.GetRemoteAddr(r)
	if remoteAddr == nil {
		return ctx, nil, log.Error("geoip.WithCountryByIP.GetRemoteAddr", "err", ErrCannotGetRemoteAddr, "req", r)
	}

	c, err := s.GetCountryByIP(remoteAddr)
	if err != nil {
		return ctx, nil, log.Error("geoip.WithCountryByIP.GetCountryByIP", "err", err, "remoteAddr", remoteAddr, "req", r)
	}
	c.IP = remoteAddr
	return NewContextCountry(ctx, c), c, nil
}

// WithCountryByIP is a simple middleware which detects the country via an IP
// address. With the detected country a new tree context.Context gets created.
func (s *Service) WithCountryByIP() ctxhttp.Middleware {
	return func(h ctxhttp.Handler) ctxhttp.Handler {
		return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var err error
			ctx, _, err = s.newContextCountryByIP(ctx, r)
			if err != nil {
				ctx = NewContextWithError(ctx, err)
			}
			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}

// WithIsCountryAllowedByIP a more advanced function. It expects from the context
// the store.ManagerReader ...
func (s *Service) WithIsCountryAllowedByIP() ctxhttp.Middleware {
	return func(h ctxhttp.Handler) ctxhttp.Handler {
		return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			sm, ok := store.FromContextReader(ctx)
			if !ok {
				return log.Error("geoip.WithCountryByIP.FromContextManagerReader", "err", ErrCannotGetStoreManagerReader)
			}
			var err error
			var ipCountry *IPCountry
			ctx, ipCountry, err = s.newContextCountryByIP(ctx, r)
			if err != nil {
				ctx = NewContextWithError(ctx, err)
				return h.ServeHTTPContext(ctx, w, r)
			}

			st, err := sm.Store()
			if err != nil {
				return log.Error("geoip.WithCountryByIP.sm.Store", "err", err)
			}
			allowedCountries, err := directory.AllowedCountries(st.Config)
			if err != nil {
				return log.Error("geoip.WithCountryByIP.directory.AllowedCountries", "err", err, "st.Config", st.Config)
			}

			if false == s.IsAllowed(st, ipCountry, allowedCountries, r) {
				h = s.altHandlerByID(st)
			}

			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}

// DefaultAlternativeHandler gets called when detected IPCountry cannot be found
// within the list of allowed countries. This handler can be overridden depending
// on the overall scope (Website, Group or Store). This function gets called in
// WithIsCountryAllowedByIP.
//
// Status is StatusServiceUnavailable
var DefaultAlternativeHandler ctxhttp.Handler = defaultAlternativeHandler

var defaultAlternativeHandler ctxhttp.Handler = ctxhttp.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	return nil
})

// altHandlerByID searches in the hierarchical order of store -> group -> website
// the next alternative handler IF a country is not allowed as defined in function
// type IsAllowedFunc.
func (s *Service) altHandlerByID(st *store.Store) ctxhttp.Handler {

	if s.storeIDs != nil && s.storeAltH != nil {
		return findHandlerByID(scope.StoreID, st.StoreID(), s.storeIDs, s.storeAltH)
	}
	if s.groupIDs != nil && s.groupAltH != nil {
		return findHandlerByID(scope.GroupID, st.Group.GroupID(), s.groupIDs, s.groupAltH)
	}
	if s.websiteIDs != nil && s.websiteAltH != nil {
		return findHandlerByID(scope.WebsiteID, st.Website.WebsiteID(), s.websiteIDs, s.websiteAltH)
	}
	return DefaultAlternativeHandler
}

// altHandlerByID returns the Handler for the searchID. If not found
// or slices have an indifferent length or something is nil it will
// return the DefaultErrorHandler.
func findHandlerByID(so scope.Scope, id int64, idsIdx utils.Int64Slice, handlers []ctxhttp.Handler) ctxhttp.Handler {

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

	if log.IsInfo() {
		log.Info("geoip.findHandlerByID.found", "scope", so.String(), "id", id, "idsIdx", idsIdx)
	}
	return prospect
}

func defaultIsCountryAllowed(_ *store.Store, c *IPCountry, allowedCountries utils.StringSlice, r *http.Request) bool {
	if false == allowedCountries.Include(c.Country.Country.IsoCode) {
		if log.IsInfo() {
			log.Info("geoip.checkAllow", "IPCountry", c, "allowedCountries", allowedCountries, "request", r)
		}
		return false
	}
	return true
}
