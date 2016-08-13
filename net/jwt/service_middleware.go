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
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
)

// SetHeaderAuthorization convenience function to set the Authorization Bearer
// Header on a request for a given token.
func SetHeaderAuthorization(req *http.Request, token []byte) {
	req.Header.Set("Authorization", "Bearer "+string(token))
}

// WithToken parses and validates a token depending on the scope. A check to the
// blacklist will be performed. The token gets added to the context for further
// processing for the next middlewares. This function depends on the runMode and
// its scope which must exists in the requests context. WithToken() does not
// change the scope of the previously initialized runMode and its scope.
func (s *Service) WithToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scpCfg := s.configFromContext(r.Context())
		if err := scpCfg.IsValid(); err != nil {
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithToken.configFromContext", log.Err(err), log.HTTPRequest("request", r))
			}
			s.ErrorHandler(errors.Wrap(err, "jwt.Service.WithToken.configFromContext")).ServeHTTP(w, r)
			return
		}
		if scpCfg.Disabled {
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithToken.Disabled", log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			next.ServeHTTP(w, r)
			return
		}

		token, err := scpCfg.ParseFromRequest(r)
		if err != nil {
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithToken.ParseFromRequest", log.Err(err), log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			// todo what should be done when the token has expired?
			scpCfg.UnauthorizedHandler(errors.Wrap(err, "[jwt] ParseFromRequest")).ServeHTTP(w, r)
			return
		}
		if s.Blacklist.Has(token.Raw) {
			err = errors.NewNotValidf(errTokenBlacklisted)
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithToken.Blacklist.Has", log.Err(err), log.Marshal("token", token), log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			// consider your ErrorHandler before leaking sensitive information.
			scpCfg.UnauthorizedHandler(err).ServeHTTP(w, r)
			return
		}

		// add token to the context
		ctx := withContext(r.Context(), token)

		// continue without changing the scope
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithRunMode sets the initial runMode, loads the token configuration, parses
// and validates a token and if the token contains a new store code it changes
// the scope for the context.
//
// RunMode optional custom runMode otherwise falls back to scope.DefaultRunMode
// which selects the default website with its default store. To use the admin
// area enable scope.Store and ID 0.
//
// Finder selects the new store ID and website ID based on the store code. It
// changes the scope in the context.
func (s *Service) WithRunMode(rm scope.RunMode, sf store.Finder) mw.Middleware {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// set run mode and add it to the context
			runMode := rm.CalculateMode(r)
			r = r.WithContext(scope.WithContextRunMode(r.Context(), runMode))

			// find the default store ID for the runMode
			storeID, websiteID, err := sf.DefaultStoreID(runMode)
			if err != nil {
				if s.Log.IsDebug() {
					s.Log.Debug("jwt.Service.WithRunMode.DefaultStoreID.Error", log.Err(err),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				s.ErrorHandler(errors.Wrap(err, "[store] WithRunMode.DefaultStoreID")).ServeHTTP(w, r)
				return
			}

			// load default scoped configuration and call next handler if disabled
			defaultScpCfg := s.configFromScope(websiteID, storeID) // scope of the DefaultStore selected by the run mode.
			if err := defaultScpCfg.IsValid(); err != nil {
				if s.Log.IsDebug() {
					s.Log.Debug("jwt.Service.WithRunMode.ConfigFromScope.Error", log.Err(err),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				s.ErrorHandler(errors.Wrap(err, "[jwt] ConfigByScopedGetter")).ServeHTTP(w, r)
				return
			}
			if defaultScpCfg.Disabled {
				if s.Log.IsDebug() {
					s.Log.Debug("jwt.Service.WithRunMode.Disabled", log.Stringer("scope", defaultScpCfg.ScopeHash), log.Object("scpCfg", defaultScpCfg),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				next.ServeHTTP(w, r)
				return
			}

			// extract token from request and validate, if invalid call scoped error handler
			token, err := defaultScpCfg.ParseFromRequest(r)
			if err != nil {
				if s.Log.IsDebug() {
					s.Log.Debug("jwt.Service.WithRunMode.ParseFromRequest", log.Err(err), log.Stringer("scope", defaultScpCfg.ScopeHash), log.Object("scpCfg", defaultScpCfg),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				defaultScpCfg.UnauthorizedHandler(errors.Wrap(err, "[jwt] ParseFromRequest")).ServeHTTP(w, r)
				return
			}
			if s.Blacklist.Has(token.Raw) {
				err = errors.NewNotValidf(errTokenBlacklisted)
				if s.Log.IsDebug() {
					s.Log.Debug("jwt.Service.WithRunMode.Blacklist.Has", log.Err(err), log.Marshal("token", token), log.Stringer("scope", defaultScpCfg.ScopeHash), log.Object("scpCfg", defaultScpCfg),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				// consider your ErrorHandler before leaking sensitive information.
				defaultScpCfg.UnauthorizedHandler(err).ServeHTTP(w, r)
				return
			}

			// add token to the context
			ctx := withContext(r.Context(), token)

			// extracts the store code from the token.
			reqCode := codeFromToken(token, defaultScpCfg.StoreCodeFieldName)
			if reqCode == "" {
				// no code found in token so call next handler and add the scope to the context
				r = r.WithContext(scope.WithContext(r.Context(), websiteID, storeID))
				if s.Log.IsDebug() {
					s.Log.Debug("jwt.Service.WithRunMode.NextHandler.WithoutCode", log.Marshal("token", token),
						log.Stringer("scope", defaultScpCfg.ScopeHash), log.Object("scpCfg", defaultScpCfg),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// convert the code string into its internal ID depending on the scope.
			newStoreID, newWebsiteID, err := sf.StoreIDbyCode(runMode, reqCode)
			if err != nil && !errors.IsNotFound(err) {
				if s.Log.IsDebug() {
					s.Log.Debug("jwt.Service.WithRunMode.IDbyCode.Error", log.Err(err), log.String("http_store_code", reqCode),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				defaultScpCfg.ErrorHandler(errors.Wrap(err, "[store] WithRunMode.IDbyCode")).ServeHTTP(w, r)
				return
			}
			if err == nil {
				storeID = newStoreID
				websiteID = newWebsiteID
			}

			r = r.WithContext(scope.WithContext(r.Context(), websiteID, storeID))

			// check if the new store ID is allowed
			isStoreAllowed, _, err := sf.IsAllowedStoreID(runMode, storeID)
			if err != nil {
				if s.Log.IsDebug() {
					s.Log.Debug("jwt.Service.WithRunMode.IsAllowedStoreID.Error", log.Err(err),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID),
						log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				defaultScpCfg.ErrorHandler(errors.Wrap(err, "[store] WithRunMode.AllowedStoreIDs")).ServeHTTP(w, r)
				return
			}

			// not found, not active, whatever, we cannot proceed.
			if !isStoreAllowed {
				if s.Log.IsDebug() {
					s.Log.Debug("jwt.Service.WithRunMode.StoreNotAllowed",
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID),
						log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}

				defaultScpCfg.UnauthorizedHandler(errors.NewUnauthorizedf("[store] RunMode %s with requested Store ID %d cannot be authorized", runMode, storeID)).ServeHTTP(w, r)
				return
			}

			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithRunMode.NextHandler.WithCode",
					log.Int64("store_id", storeID), log.Int64("website_id", websiteID),
					log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func codeFromToken(token csjwt.Token, storeCodeFieldName string) string {
	// extracts the store code from the token.
	key := store.CodeFieldName
	if storeCodeFieldName != "" {
		key = storeCodeFieldName
	}
	tokenStoreCode, _ := token.Claims.Get(key)
	return conv.ToString(tokenStoreCode)
}
