// Copyright (c) 2014 Olivier Poitrey <rs@dailymotion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cors

import (
	"net/http"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/util/errors"
)

// WithCORS to be used as a middleware for net.Handler. The applied
// configuration is used for the all store scopes or if the PkgBackend has been
// provided then on a website specific level. Middleware expects to find in a
// context a store.FromContextProvider().
func (s *Service) WithCORS() mw.Middleware {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			requestedStore, err := store.FromContextRequestedStore(ctx)
			if err != nil {
				if s.defaultScopeCache.log.IsDebug() {
					s.defaultScopeCache.log.Debug("Service.WithCORS.FromContextProvider", log.Err(err), log.HTTPRequest("request", r))
				}
				err = errors.Wrap(err, "[cors] FromContextRequestedStore")
				h.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			// the scpCfg depends on how you have initialized the storeService during app boot.
			// requestedStore.Website.Config is the reason that all options only support
			// website scope and not group or store scope.
			scpCfg := s.configByScopedGetter(requestedStore.Website.Config)

			if err := scpCfg.isValid(); err != nil {
				if s.defaultScopeCache.log.IsDebug() {
					s.defaultScopeCache.log.Debug("Service.WithCORS.configByScopedGetter", log.Err(err), log.Marshal("requestedStore", requestedStore), log.HTTPRequest("request", r))
				}
				err = errors.Wrap(err, "[cors] ConfigByScopedGetter")
				h.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			if s.defaultScopeCache.log.IsInfo() {
				s.defaultScopeCache.log.Info("Service.WithCORS.handleActualRequest", log.String("method", r.Method), log.Object("scopedConfig", scpCfg))
			}

			if r.Method == methodOptions {
				if s.defaultScopeCache.log.IsDebug() {
					s.defaultScopeCache.log.Debug("Service.WithCORS.handlePreflight", log.String("method", r.Method), log.Bool("OptionsPassthrough", scpCfg.optionsPassthrough))
				}
				scpCfg.handlePreflight(w, r)
				// Preflight requests are standalone and should stop the chain as some other
				// middleware may not handle OPTIONS requests correctly. One typical example
				// is authentication middleware ; OPTIONS requests won't carry authentication
				// headers (see #1)
				if scpCfg.optionsPassthrough {
					h.ServeHTTP(w, r)
				}
				return
			}
			scpCfg.handleActualRequest(w, r)
			h.ServeHTTP(w, r)
		})
	}
}
