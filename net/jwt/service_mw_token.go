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
	loghttp "github.com/corestoreio/csfw/log/http"
	"github.com/corestoreio/csfw/util/errors"
)

// WithToken parses and validates a token depending on the scope. A check to the
// blacklist will be performed. The token gets added to the context for further
// processing for the next middlewares. This function depends on the runMode and
// its scope which must exists in the requests context. WithToken() does not
// change the scope of the previously initialized runMode and its scope.
func (s *Service) WithToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scpCfg, err := s.configByContext(r.Context())
		if err != nil {
			s.Log.Info("jwt.Service.WithToken.configByContext.Error", log.Err(err))
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithToken.configByContext", log.Err(err), loghttp.Request("request", r))
			}
			s.ErrorHandler(errors.Wrap(err, "jwt.Service.WithToken.configFromContext")).ServeHTTP(w, r)
			return
		}
		if scpCfg.Disabled {
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithToken.Disabled", log.Stringer("scope", scpCfg.ScopeID), log.Object("scpCfg", scpCfg), loghttp.Request("request", r))
			}
			next.ServeHTTP(w, r)
			return
		}

		token, err := scpCfg.ParseFromRequest(s.Blacklist, r)
		if err != nil {
			s.Log.Info("jwt.Service.WithToken.ParseFromRequest.Error", log.Err(err))
			if s.Log.IsDebug() {
				s.Log.Debug("jwt.Service.WithToken.ParseFromRequest", log.Err(err), log.Marshal("token", token), log.Stringer("scope", scpCfg.ScopeID), log.Object("scpCfg", scpCfg), loghttp.Request("request", r))
			}
			// todo what should be done when the token has expired?
			scpCfg.UnauthorizedHandler(errors.Wrap(err, "[jwt] WithToken.ParseFromRequest")).ServeHTTP(w, r)
			return
		}

		// add token to the context
		ctx := withContext(r.Context(), token)

		// continue without changing the scope
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
