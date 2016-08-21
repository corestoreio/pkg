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
	"github.com/corestoreio/csfw/util/errors"
)

// WithCORS to be used as a middleware. This middleware expects to find a
// scope.FromContext().
func (s *Service) WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		scpCfg := s.configByContext(r.Context())
		if err := scpCfg.IsValid(); err != nil {
			s.Log.Info("cors.Service.WithCORS.configByContext.Error", log.Err(err))
			if s.Log.IsDebug() {
				s.Log.Debug("cors.Service.WithCORS.configByContext", log.Err(err), log.HTTPRequest("request", r))
			}
			s.ErrorHandler(errors.Wrap(err, "cors.Service.WithCORS.configFromContext")).ServeHTTP(w, r)
			return
		}

		if s.Log.IsInfo() {
			s.Log.Info("cors.Service.WithCORS.handleActualRequest", log.String("method", r.Method), log.Object("scopedConfig", scpCfg))
		}

		if r.Method == methodOptions {
			if s.Log.IsDebug() {
				s.Log.Debug("cors.Service.WithCORS.handlePreflight", log.String("method", r.Method), log.Bool("OptionsPassthrough", scpCfg.OptionsPassthrough))
			}
			scpCfg.handlePreflight(w, r)
			// Preflight requests are standalone and should stop the chain as
			// some other middleware may not handle OPTIONS requests correctly.
			// One typical example is authentication middleware ; OPTIONS
			// requests won't carry authentication headers (see #1)
			if scpCfg.OptionsPassthrough {
				next.ServeHTTP(w, r)
			}
			return
		}
		scpCfg.handleActualRequest(w, r)
		next.ServeHTTP(w, r)
	})
}
