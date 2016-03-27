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

package ctxthrottled

import (
	"math"
	"net/http"
	"strconv"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/juju/errors"
	"golang.org/x/net/context"
	"gopkg.in/throttled/throttled.v2"
)

// WithRateLimit wraps an ctxhttp.Handler to limit incoming requests.
// Requests that are not limited will be passed to the handler
// unchanged.  Limited requests will be passed to the DeniedHandler.
// X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset and
// Retry-After headers will be written to the response based on the
// values in the RateLimitResult.
func (hrl *HTTPRateLimit) WithRateLimit() ctxhttp.Middleware {

	return func(h ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			_, reqStore, err := store.FromContextProvider(ctx)
			if err != nil {
				return err
			}

			var rl throttled.RateLimiter
			var ok bool

			scopeHash := scope.NewHash(scope.WebsiteID, reqStore.WebsiteID())
			hrl.mu.RLock()
			rl, ok = hrl.scopedRLs[scopeHash]
			hrl.mu.RUnlock()

			if !ok {
				var err error
				rl, err = hrl.RateLimiterFactory(hrl.Backend, reqStore.Website.Config)
				if err != nil {
					return errors.Mask(err)
				}
				hrl.mu.Lock()
				hrl.scopedRLs[scopeHash] = rl
				hrl.mu.Unlock()
			}

			var k string
			if hrl.VaryByer != nil {
				k = hrl.Key(r)
			}

			limited, context, err := rl.RateLimit(k, 1)
			if err != nil {
				return err
			}

			setRateLimitHeaders(w, context)

			if !limited {
				return h.ServeHTTPContext(ctx, w, r)
			}

			return hrl.DeniedHandler(ctx, w, r)
		}
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
