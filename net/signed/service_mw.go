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
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/zenazn/goji/web/mutil"
)

// todo: refactor to use Service type and backendsigned package

// WithCompressor is a middleware applies the GZIP or deflate algorithm on
// the bytes writer. GZIP or deflate usage depends on the HTTP Accept
// Encoding header. Flush(), Hijack() and CloseNotify() interfaces will be
// preserved. No header set, no compression takes place. GZIP has priority
// before deflate.
func (s *Service) WithResponseSignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scpCfg := s.configByContext(r.Context())
		if err := scpCfg.IsValid(); err != nil {
			s.Log.Info("signed.Service.WithRateLimit.configByContext.Error", log.Err(err))
			if s.Log.IsDebug() {
				s.Log.Debug("signed.Service.WithRateLimit.configByContext", log.Err(err), log.HTTPRequest("request", r))
			}
			s.ErrorHandler(errors.Wrap(err, "signed.Service.WithRateLimit.configFromContext")).ServeHTTP(w, r)
			return
		}
		if scpCfg.Disabled {
			if s.Log.IsDebug() {
				s.Log.Debug("signed.Service.WithRateLimit.Disabled", log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			next.ServeHTTP(w, r)
			return
		}

		alg := scpCfg.hashPool.Get()
		defer scpCfg.hashPool.Put(alg)

		if scpCfg.InTrailer {
			buf := bufferpool.Get()

			lw := mutil.WrapWriter(w)
			lw.Tee(alg) // write also to alg

			lw.Header().Set("Trailer", scpCfg.HeaderName)

			next.ServeHTTP(lw, r)

			tmp := alg.Sum(buf.Bytes()) // append to buffer
			buf.Reset()
			_, _ = buf.Write(tmp)

			sig := Signature{
				EncodeFn:  scpCfg.EncodeFn,
				KeyID:     "test",
				Algorithm: scpCfg.AlgorithmName,
			}
			sig.Write(w, buf.Bytes())

			bufferpool.Put(buf)
		}

	})
}

func (s *Service) WithRequestSignatureValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scpCfg := s.configByContext(r.Context())
		if err := scpCfg.IsValid(); err != nil {
			s.Log.Info("signed.Service.WithRateLimit.configByContext.Error", log.Err(err))
			if s.Log.IsDebug() {
				s.Log.Debug("signed.Service.WithRateLimit.configByContext", log.Err(err), log.HTTPRequest("request", r))
			}
			s.ErrorHandler(errors.Wrap(err, "signed.Service.WithRateLimit.configFromContext")).ServeHTTP(w, r)
			return
		}
		if scpCfg.Disabled {
			if s.Log.IsDebug() {
				s.Log.Debug("signed.Service.WithRateLimit.Disabled", log.Stringer("scope", scpCfg.ScopeHash), log.Object("scpCfg", scpCfg), log.HTTPRequest("request", r))
			}
			next.ServeHTTP(w, r)
			return
		}
		// check if encoded form ... or body is www form ... or raw ... then verify
		panic("todo implement it")
	})
}
