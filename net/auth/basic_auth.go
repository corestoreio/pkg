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

package auth

import (
	"fmt"
	"net/http"

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/hashpool"
)

// BasicAuthFunc defines a function to validate basic auth credentials.
type BasicAuthFunc func(string, string) bool

func basicAuthValidator(hashName, username, password string) (BasicAuthFunc, error) {
	ht, err := hashpool.FromRegistry(hashName)
	if err != nil {
		return nil, errors.Wrapf(err, "[auth] Failed to find %q in hashpool. Please register it.", hashName)
	}
	return func(givenUser, givenPass string) bool {
		return ht.EqualPairs([]byte(username), []byte(givenUser), []byte(password), []byte(givenPass))
	}, nil
}

func basicAuth(baf BasicAuthFunc) ProviderFunc {
	var errInvalidData = errors.NewUnauthorizedf("[auth] Invalid username or password")
	return func(scopeID scope.TypeID, r *http.Request) (callNext bool, err error) {
		givenUser, givenPass, ok := r.BasicAuth()
		if !ok {
			return true, errors.Wrapf(errInvalidData, "[auth] Basic Auth not found in request. Scope(%s)", scopeID)
		}
		if !baf(givenUser, givenPass) {
			return true, errors.Wrapf(errInvalidData, "[auth] Username or password incorrect. Scope(%s)", scopeID)
		}
		return true, nil
	}
}
func basicAuthHandler(realm string) mw.ErrorHandler {
	return func(_ error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm=%q`, realm))
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}
