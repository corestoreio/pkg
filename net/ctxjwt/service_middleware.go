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

package ctxjwt

import (
	"net/http"

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/util/errors"
)

// SetHeaderAuthorization convenience function to set the Authorization Bearer
// Header on a request for a given token.
func SetHeaderAuthorization(req *http.Request, token []byte) {
	req.Header.Set("Authorization", "Bearer "+string(token))
}

// WithInitTokenAndStore  represent a middleware handler which parses and validates a
// token, adds the token to the context and initializes the requested store
// and scope.is a middleware which initializes a request based store
// via a JSON Web Token.
// Extracts the store.Provider and csjwt.Token from context.Context. If the requested
// store is different than the initialized requested store than the new requested
// store will be saved in the context.
func (s *Service) WithInitTokenAndStore() mw.Middleware {

	return func(hf http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			errHandler := hf
			if s.defaultScopeCache.ErrorHandler != nil {
				errHandler = s.defaultScopeCache.ErrorHandler
			}
			ctx := r.Context()

			storeService, requestedStore, err := store.FromContextProvider(ctx)
			if err != nil {
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithInitTokenAndStore.FromContextProvider", "err", err, "ctx", ctx, "req", r)
				}
				err = errors.Wrap(err, "[ctxjwt] FromContextProvider")
				errHandler.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			// the scpCfg depends on how you have initialized the storeService during app boot.
			// requestedStore.Website.Config is the reason that all options only support
			// website scope and not group or store scope.
			scpCfg, err := s.ConfigByScopedGetter(requestedStore.Website.Config)
			if err != nil {
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithInitTokenAndStore.ConfigByScopedGetter", "err", err, "requestedStore", requestedStore, "ctx", ctx, "req", r)
				}
				err = errors.Wrap(err, "[ctxjwt] ConfigByScopedGetter")
				errHandler.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			if scpCfg.ErrorHandler != nil {
				errHandler = scpCfg.ErrorHandler
			}

			token, err := scpCfg.ParseFromRequest(r)
			if err != nil {
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithInitTokenAndStore.ParseFromRequest", "err", err, "requestedStore", requestedStore, "scpCfg", scpCfg, "ctx", ctx, "req", r)
				}
				err = errors.Wrap(err, "[ctxjwt] ParseFromRequest")
				errHandler.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			if false == token.Valid {
				err := errors.NewNotValidf(errTokenInvalid)
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithInitTokenAndStore.token.valid", "err", err, "token", token, "requestedStore", requestedStore, "scpCfg", scpCfg, "ctx", ctx, "req", r)
				}
				errHandler.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			if s.Blacklist.Has(token.Raw) {
				err := errors.NewNotValidf(errTokenBlacklisted)
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithInitTokenAndStore.token.blacklist", "err", err, "token", token, "requestedStore", requestedStore, "scpCfg", scpCfg, "ctx", ctx, "req", r)
				}
				errHandler.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			// add token to the context
			ctx = withContext(ctx, token)

			scopeOption, err := ScopeOptionFromClaim(token.Claims)

			if err != nil && errors.IsNotFound(err) {
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithInitTokenAndStore.ScopeOptionFromClaim.notfound", "err", err, "token", token, "requestedStore", requestedStore, "scpCfg", scpCfg, "ctx", ctx, "req", r)
				}
				// move on when the store code cannot be found in the token.
				hf.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if err != nil {
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithInitTokenAndStore.ScopeOptionFromClaim.error", "err", err, "token", token, "requestedStore", requestedStore, "scpCfg", scpCfg, "ctx", ctx, "req", r)
				}
				// invalid syntax of store code
				errHandler.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			newRequestedStore, err := storeService.RequestedStore(scopeOption)
			if err != nil {
				err = errors.Wrap(err, "[ctxjwt] storeService.RequestedStore")
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithInitTokenAndStore.GetRequestedStore", "err", err, "token", token, "newRequestedStore", newRequestedStore, "scpCfg", scpCfg, "ctx", ctx, "req", r)
				}
				errHandler.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			if newRequestedStore.StoreID() != requestedStore.StoreID() {
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithInitTokenAndStore.SetRequestedStore", "token", token, "newRequestedStore", newRequestedStore, "requestedStore", requestedStore, "scpCfg", scpCfg, "ctx", ctx, "req", r)
				}
				// this should not lead to a bug because the previously set store.Provider and requestedStore
				// will still exists and have not been/cannot be removed.
				ctx = store.WithContextProvider(ctx, storeService, newRequestedStore)
			}
			// yay! we made it! the token and the requested store is valid!
			hf.ServeHTTP(w, r.WithContext(ctx))
			return
		})
	}
}
