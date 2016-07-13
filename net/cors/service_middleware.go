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

package cors

import (
	"net/http"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/mw"
)

// WithCORS to be used as a middleware for net.Handler. The applied
// configuration is used for the all store scopes or if the PkgBackend has been
// provided then on a website specific level. Middleware expects to find in a
// context a store.FromContextProvider().
func (s *Service) WithCORS() mw.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			scpCfg := s.configFromContext(w, r)
			if scpCfg.IsValid() != nil {
				// every error gets previously logged in the configFromContext() function.
				return
			}

			if s.Log.IsInfo() {
				s.Log.Info("Service.WithCORS.handleActualRequest", log.String("method", r.Method), log.Object("scopedConfig", scpCfg))
			}

			if r.Method == methodOptions {
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithCORS.handlePreflight", log.String("method", r.Method), log.Bool("OptionsPassthrough", scpCfg.optionsPassthrough))
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
