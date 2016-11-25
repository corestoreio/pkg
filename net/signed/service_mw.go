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

package signed

import (
	"net/http"

	"github.com/corestoreio/csfw/log"
	loghttp "github.com/corestoreio/csfw/log/http"
	"github.com/corestoreio/csfw/util/errors"
)

// WithResponseSignature hashes the data written to http.ResponseWriter and adds
// the hash to the HTTP header. For large data sets use the option InTrailer to
// provide stream based hashing but the hash gets written into the HTTP trailer.
// Not all clients can read the HTTP trailer.
func (s *Service) WithResponseSignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scpCfg, err := s.configByContext(r.Context())
		if err != nil {
			if s.Log.IsDebug() {
				s.Log.Debug("signed.Service.WithResponseSignature.configByContext", log.Err(err), loghttp.Request("request", r))
			}
			s.ErrorHandler(errors.Wrap(err, "signed.Service.WithResponseSignature.configFromContext")).ServeHTTP(w, r)
			return
		}
		if scpCfg.Disabled {
			if s.Log.IsDebug() {
				s.Log.Debug("signed.Service.WithResponseSignature.Disabled", log.Stringer("scope", scpCfg.ScopeID), log.Object("scpCfg", scpCfg), loghttp.Request("request", r))
			}
			next.ServeHTTP(w, r)
			return
		}

		if scpCfg.InTrailer {
			// direct output to the client and the signature will be inserted
			// after the body has been written. ideal for streaming but not all
			// clients can process a trailer.
			scpCfg.writeTrailer(next, w, r)
			return
		}
		// the write to w gets buffered and we calculate the checksum of the
		// buffer and then flush the buffer to the client.
		scpCfg.writeBuffered(next, w, r)
	})
}

// WithRequestSignatureValidation extracts from the header or trailer the hash
// value and hashes the body of the incoming request and compares those two
// hashes. On success the next handler will be called otherwise the scope based
// error handler.
func (s *Service) WithRequestSignatureValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scpCfg, err := s.configByContext(r.Context())
		if err != nil {
			if s.Log.IsDebug() {
				s.Log.Debug("signed.Service.WithRequestSignatureValidation.configByContext", log.Err(err), log.Stringer("scope", scpCfg.ScopeID), log.Object("scpCfg", scpCfg), loghttp.Request("request", r))
			}
			s.ErrorHandler(errors.Wrap(err, "signed.Service.WithRequestSignatureValidation.configFromContext")).ServeHTTP(w, r)
			return
		}
		if scpCfg.Disabled {
			if s.Log.IsDebug() {
				s.Log.Debug("signed.Service.WithRequestSignatureValidation.Disabled", log.Stringer("scope", scpCfg.ScopeID), log.Object("scpCfg", scpCfg), loghttp.Request("request", r))
			}
			next.ServeHTTP(w, r)
			return
		}

		if err := scpCfg.ValidateBody(r); err != nil {
			if s.Log.IsInfo() {
				s.Log.Info("signed.Service.WithRequestSignatureValidation.ValidateBody", log.Err(err))
			}
			if s.Log.IsDebug() {
				s.Log.Debug("signed.Service.WithRequestSignatureValidation.ValidateBody", log.Err(err), log.Stringer("scope", scpCfg.ScopeID), log.Object("scpCfg", scpCfg), loghttp.Request("request", r))
			}
			scpCfg.ErrorHandler(err).ServeHTTP(w, r)
			return
		}
		// signature validated and trusted
		next.ServeHTTP(w, r)
	})
}
