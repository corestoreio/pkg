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

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/store"
	"github.com/juju/errors"
	"golang.org/x/net/context"
)

// ErrTokenBlacklisted returned by the middleware if the token can be found
// within the black list.
var ErrTokenBlacklisted = errors.New("Token has been black listed")

// ErrTokenInvalid returned by the middleware to make understandable that
// a token has been invalidated.
var ErrTokenInvalid = errors.New("Token has become invalid")

// SetHeaderAuthorization convenience function to set the Authorization Bearer
// Header on a request for a given token.
func SetHeaderAuthorization(req *http.Request, token []byte) {
	req.Header.Set("Authorization", "Bearer "+string(token))
}

// WithInitToken represent a middleware handler which parses and validates a
// token and adds it to the context. For a POST or a PUT request, it also parses the
// request body as a form. The extracted valid
// token will be added to the context or if an error has occurred, that error will
// be added to the context. The extracted token will be checked
// against the Blacklist.
//
// Tip: Instead of passing the token as an HTML Header you can also add the token
// to a form (multipart/form-data) with an input name of access_token. If the
// token cannot be found within the Header the fallback triggers the lookup within the form.
//func (s *Service) WithInitToken() ctxhttp.Middleware {
//
//	return func(h ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
//		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//
//		}
//	}
//}

// WithInitTokenAndStore  represent a middleware handler which parses and validates a
// token, adds the token to the context and initializes the requested store
// and scope.is a middleware which initializes a request based store
// via a JSON Web Token.
// Extracts the store.Provider and csjwt.Token from context.Context. If the requested
// store is different than the initialized requested store than the new requested
// store will be saved in the context.
func (s *Service) WithInitTokenAndStore() ctxhttp.Middleware {

	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			errHandler := s.defaultScopeCache.errorHandler

			storeService, requestedStore, err := store.FromContextProvider(ctx)
			if err != nil {
				return errHandler.ServeHTTPContext(withContextError(ctx, err), w, r)
			}

			// the scpCfg depends on how you have initialized the storeService during app boot.
			// requestedStore.Website.Config is the reason that all options only support
			// website scope and not group or store scope.
			scpCfg, err := s.getConfigByScopedGetter(requestedStore.Website.Config)
			if err != nil {
				err = errors.Mask(err)
				return errHandler.ServeHTTPContext(withContextError(ctx, err), w, r)
			}

			if scpCfg.errorHandler != nil {
				errHandler = scpCfg.errorHandler
			}

			token, err := scpCfg.parseFromRequest(r)
			if err != nil {
				err = errors.Mask(err)
				return errHandler.ServeHTTPContext(withContextError(ctx, err), w, r)
			}

			if false == token.Valid {
				return errHandler.ServeHTTPContext(withContextError(ctx, ErrTokenInvalid), w, r)
			}

			if s.Blacklist.Has(token.Raw) {
				return errHandler.ServeHTTPContext(withContextError(ctx, ErrTokenBlacklisted), w, r)
			}

			// add token to the context
			ctx = withContext(ctx, token)

			scopeOption, err := ScopeOptionFromClaim(token.Claims)

			if err == store.ErrStoreNotFound {
				// move on when the store code cannot be found in the token.
				return hf.ServeHTTPContext(ctx, w, r)
			}

			if err != nil {
				// invalid syntax of store code
				return errHandler.ServeHTTPContext(withContextError(ctx, err), w, r)
			}

			if PkgLog.IsDebug() {
				PkgLog.Debug("ctxjwt.Service.WithInitTokenAndStore.FromClaim", "token", token, "ScopeOption", scopeOption)
			}

			newRequestedStore, err := storeService.RequestedStore(scopeOption)
			if err != nil {
				err = errors.Mask(err)
				return errHandler.ServeHTTPContext(withContextError(ctx, err), w, r)
			}

			if newRequestedStore.StoreID() != requestedStore.StoreID() {
				// this should not lead to a bug because the previously set store.Provider and requestedStore
				// will still exists and have not been/cannot be removed.
				ctx = store.WithContextProvider(ctx, storeService, newRequestedStore)
			}
			// yay! we made it! the token and the requested store is valid!
			return hf.ServeHTTPContext(ctx, w, r)
		}
	}
}
