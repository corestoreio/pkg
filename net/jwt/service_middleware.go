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

package jwt

import (
	"net/http"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/util/errors"
)

// SetHeaderAuthorization convenience function to set the Authorization Bearer
// Header on a request for a given token.
func SetHeaderAuthorization(req *http.Request, token []byte) {
	req.Header.Set("Authorization", "Bearer "+string(token))
}

// WithInitTokenAndStore represent a middleware handler which parses and
// validates a token, adds the token to the context and initializes the
// requested store and scope.is a middleware which initializes a request based
// store via a JSON Web Token. Extracts the store.Provider and csjwt.Token from
// context.Context. If the requested store is different than the initialized
// requested store than the new requested store will be saved in the context.
func (s *Service) WithInitTokenAndStore(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		scpCfg := s.configFromContext(w, r)
		if scpCfg.IsValid() != nil {
			// every error gets previously logged in the configFromContext() function.
			return
		}
		if scpCfg.Disabled {
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithInitTokenAndStore.Disabled", log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			hf.ServeHTTP(w, r)
			return
		}

		token, err := scpCfg.ParseFromRequest(r)
		if err != nil {
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithInitTokenAndStore.ParseFromRequest", log.Err(err), log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			scpCfg.ErrorHandler(errors.Wrap(err, "[jwt] ParseFromRequest")).ServeHTTP(w, r)
			return
		}
		if s.Blacklist.Has(token.Raw) {
			err = errors.NewNotValidf(errTokenBlacklisted)
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithInitTokenAndStore.Blacklist.Has", log.Err(err), log.Marshal("token", token), log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			// consider your ErrorHandler before leaking sensitive information.
			scpCfg.ErrorHandler(err).ServeHTTP(w, r)
			return
		}

		// add token to the context
		ctx := withContext(r.Context(), token)

		scopeOption, err := ScopeOptionFromClaim(token.Claims)
		switch {
		case err != nil && errors.IsNotFound(err):
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithInitTokenAndStore.ScopeOptionFromClaim.notFound", log.Err(err), log.Marshal("token", token), log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			// move on when the store code cannot be found in the token.
			// todo(CS) this should be an error or make it configurable that either error or just go on
			hf.ServeHTTP(w, r.WithContext(ctx))
			return

		case err != nil:
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithInitTokenAndStore.ScopeOptionFromClaim.error", log.Err(err), log.Marshal("token", token), log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			// invalid syntax of store code
			scpCfg.ErrorHandler(err).ServeHTTP(w, r)
			return

		case scopeOption.StoreCode() == requestedStore.StoreCode():
			// move on when there is no change between scopeOption and requestedStore, skip the lookup in func RequestedStore()
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithInitTokenAndStore.ScopeOptionFromClaim.StoreCodeEqual", log.Err(err), log.Marshal("token", token), log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			hf.ServeHTTP(w, r.WithContext(ctx))
			return

		case s.StoreService == nil:
			// when StoreService has not been set, do not change the store despite there is another requested one.
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithInitTokenAndStore.ScopeOptionFromClaim.StoreServiceIsNil", log.Err(err), log.Marshal("token", token), log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			hf.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		newRequestedStore, err := s.StoreService.RequestedStore(scopeOption)
		if err != nil {
			err = errors.Wrap(err, "[jwt] storeService.RequestedStore")
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithInitTokenAndStore.StoreService.RequestedStore", log.Err(err), log.Marshal("token", token), log.Marshal("newRequestedStore", newRequestedStore), log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			scpCfg.ErrorHandler(err).ServeHTTP(w, r)
			return
		}

		if newRequestedStore.ID() != requestedStore.StoreID() {
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithInitTokenAndStore.SetRequestedStore", log.Err(err), log.Marshal("token", token), log.Marshal("newRequestedStore", newRequestedStore), log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			// this should not lead to a bug because the previously set store.Provider and requestedStore
			// will still exists and have not been/cannot be removed.
			ctx = store.WithContextRequestedStore(ctx, newRequestedStore)
		}
		// yay! we made it! the token and the requested store are valid!
		hf.ServeHTTP(w, r.WithContext(ctx))
	})
}
