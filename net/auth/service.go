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

//go:generate go run ../internal/scopedservice/main_copy.go "$GOPACKAGE"

package auth

import (
	"net/http"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// Service describes
type Service struct {
	service
}

// New creates a new Cors handler with the provided options.
func New(opts ...Option) (*Service, error) {
	s := &Service{
		scopeCache: make(map[scope.Hash]ScopedConfig),
	}
	if err := s.Options(WithDefaultConfig(scope.Default, 0)); err != nil {
		return nil, errors.Wrap(err, "[mwauth] Options WithDefaultConfig")
	}
	if err := s.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[mwauth] Options Any Config")
	}
	return s, nil
}

// WithAuthentication to be used as a middleware for net.Handler. The applied
// configuration is used for the all store scopes or if the PkgBackend has been
// provided then on a website specific level. Middleware expects to find in a
// context a store.FromContextProvider().
func (s *Service) WithAuthentication() mw.Middleware {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			requestedStore, err := store.FromContextRequestedStore(ctx)
			if err != nil {
				if s.defaultScopeCache.log.IsDebug() {
					s.defaultScopeCache.log.Debug("Service.WithCORS.FromContextProvider", log.Err(err), log.HTTPRequest("request", r))
				}
				err = errors.Wrap(err, "[mwauth] FromContextRequestedStore")
				h.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			// the scpCfg depends on how you have initialized the storeService during app boot.
			// requestedStore.Website.Config is the reason that all options only support
			// website scope and not group or store scope.
			/* scpCfg */ _, err = s.configByScopedGetter(requestedStore.Website.Config) // TODO support ALL scopes, @see package geoip
			if err != nil {
				if s.defaultScopeCache.log.IsDebug() {
					s.defaultScopeCache.log.Debug("Service.WithCORS.configByScopedGetter", log.Err(err), log.Marshal("requestedStore", requestedStore), log.HTTPRequest("request", r))
				}
				err = errors.Wrap(err, "[mwauth] ConfigByScopedGetter")
				h.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

// WithCookieValidation validates a secure cookie which has been set by
// the middleware WithAuthentication()
func (s *Service) WithCookieValidation() mw.Middleware {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("Paaaaannniiiiiccc")
		})
	}
}
