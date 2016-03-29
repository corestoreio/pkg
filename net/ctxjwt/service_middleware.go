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
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
)

// SetHeaderAuthorization convenience function to set the Authorization Bearer
// Header on a request for a given token.
func SetHeaderAuthorization(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

// WithParseAndValidate represent a middleware handler. For POST or
// PUT requests, it also parses the request body as a form. The extracted valid
// token will be added to the context. The extracted token will be checked
// against the Blacklist. errHandler is an optional argument. Only the first
// item in the slice will be considered. Default errHandler is:
//
//		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
//
// ProTip: Instead of passing the token as an HTML Header you can also add the token
// to a form (multipart/form-data) with an input name of access_token. If the
// token cannot be found within the Header the fallback triggers the lookup within the form.
func (s *Service) WithParseAndValidate() ctxhttp.Middleware {

	return func(h ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			_, reqStore, err := store.FromContextProvider(ctx)
			if err != nil {
				return err
			}

			scpCfg, err := s.getConfigByScopedGetter(reqStore.Website.Config)
			if err != nil {
				return err
			}

			token, err := jwt.ParseFromRequest(r, scpCfg.keyFunc)

			var inBL bool
			if token != nil {
				inBL = s.Blacklist.Has([]byte(token.Raw))
			}

			if token != nil && err == nil && token.Valid && !inBL {
				return h.ServeHTTPContext(WithContext(ctx, token), w, r)
			}

			return scpCfg.errorHandler(WithContextError(ctx, err), w, r)
		}
	}
}
