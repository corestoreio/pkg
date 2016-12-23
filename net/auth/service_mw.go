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
	"net/http"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	loghttp "github.com/corestoreio/log/http"
)

// WithAuthentication to be used as a middleware for net.Handler. The applied
// configuration is used for the all store scopes or if the PkgBackend has been
// provided then on a website specific level. Middleware expects to find in a
// context a store.FromContextProvider().
func (s *Service) WithAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scpCfg, err := s.configByContext(r.Context())
		if err != nil {
			if s.Log.IsDebug() {
				s.Log.Debug("auth.Service.WithAuthentication.configByContext", log.Err(err), loghttp.Request("request", r))
			}
			s.ErrorHandler(errors.Wrap(err, "jwt.Service.WithToken.configFromContext")).ServeHTTP(w, r)
			return
		}
		if scpCfg.Disabled {
			if s.Log.IsDebug() {
				s.Log.Debug("auth.Service.WithAuthentication.Disabled", log.Stringer("scope", scpCfg.ScopeID), log.Object("scpCfg", scpCfg), loghttp.Request("request", r))
			}
			next.ServeHTTP(w, r)
			return
		}
		if err := scpCfg.Authenticate(r); err != nil {
			if s.Log.IsDebug() {
				s.Log.Debug("auth.Service.Authenticate.Failed", log.Err(err), log.Stringer("scope", scpCfg.ScopeID), log.Object("scpCfg", scpCfg), loghttp.Request("request", r))
			}
			scpCfg.UnauthorizedHandler(errors.Wrap(err, "[auth] Authentication failed"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
