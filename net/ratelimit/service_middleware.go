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

package ratelimit

import (
	"math"
	"net/http"
	"strconv"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/util/errors"
	"gopkg.in/throttled/throttled.v2"
)

// WithRateLimit wraps an http.Handler to limit incoming requests. Requests that
// are not limited will be passed to the handler unchanged.  Limited requests
// will be passed to the DeniedHandler. X-RateLimit-Limit,
// X-RateLimit-Remaining, X-RateLimit-Reset and Retry-After headers will be
// written to the response based on the values in the RateLimitResult. The next
// handler may check an error with FromContextRateLimit().
func (s *Service) WithRateLimit() mw.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			requestedStore, err := store.FromContextRequestedStore(r.Context())
			if err != nil {
				err = errors.Wrap(err, "[ratelimit] FromContextRequestedStore")
				h.ServeHTTP(w, wrapContextError(r, err))
				return
			}

			// requestedStore.Config contains the scope for store and then
			// website or finally can fall back to default scope.
			scpCfg := s.configByScopedGetter(requestedStore.Config)
			if err := scpCfg.isValid(); err != nil {
				if s.Log.IsDebug() {
					s.Log.Debug("Service.WithRateLimit.configByScopedGetter.Error",
						log.Err(err),
						log.Stringer("scope", scpCfg.scopeHash),
						log.Marshal("requestedStore", requestedStore),
						log.HTTPRequest("request", r),
					)
				}
				err = errors.Wrap(err, "[ratelimit] ConfigByScopedGetter")
				h.ServeHTTP(w, wrapContextError(r, err))
				return
			}

			if scpCfg.disabled {
				h.ServeHTTP(w, r)
				return
			}

			isLimited, rlResult, err := scpCfg.requestRateLimit(r)
			if s.Log.IsDebug() {
				s.Log.Debug("Service.WithRateLimit.configByScopedGetter.RateLimit",
					log.Err(err),
					log.Bool("is_limited", isLimited),
					log.Object("rate_limit_result", rlResult),
					log.Stringer("scope", scpCfg.scopeHash),
					log.Marshal("requestedStore", requestedStore),
					log.HTTPRequest("request", r),
				)
			}
			if err != nil {
				err = errors.Wrap(err, "[ratelimit] scpCfg.RateLimit")
				h.ServeHTTP(w, wrapContextError(r, err))
				return
			}

			setRateLimitHeaders(w, rlResult)
			next := scpCfg.deniedHandler
			if !isLimited {
				next = h
			}
			next.ServeHTTP(w, r)
		})
	}
}

func setRateLimitHeaders(w http.ResponseWriter, rlr throttled.RateLimitResult) {
	if v := rlr.Limit; v >= 0 {
		w.Header().Add("X-RateLimit-Limit", strconv.Itoa(v))
	}

	if v := rlr.Remaining; v >= 0 {
		w.Header().Add("X-RateLimit-Remaining", strconv.Itoa(v))
	}

	if v := rlr.ResetAfter; v >= 0 {
		vi := int(math.Ceil(v.Seconds()))
		w.Header().Add("X-RateLimit-Reset", strconv.Itoa(vi))
	}

	if v := rlr.RetryAfter; v >= 0 {
		vi := int(math.Ceil(v.Seconds()))
		w.Header().Add("Retry-After", strconv.Itoa(vi))
	}
}
