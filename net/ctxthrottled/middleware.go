// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/utils/log"
	"golang.org/x/net/context"
	"gopkg.in/throttled/throttled.v2"
	"net/http"
)

// WithLimiter a middleware for rate limiting
func WithLimiter(cr config.ReaderPubSuber, rlStore throttled.GCRAStore) ctxhttp.Middleware {

	// Having the GetBool command here, means you must restart the app to take
	// changes in effect. @todo refactor and use pub/sub to automatically change

	// "gopkg.in/throttled/throttled.v2/store/memstore"
	//	store, err := memstore.New(65536)
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	// Maximum burst of 5 which refills at 20 tokens per minute.
	quota := throttled.RateQuota{throttled.PerMin(20), 5}

	rateLimiter, err := throttled.NewGCRARateLimiter(rlStore, quota)
	if err != nil {
		log.Fatal(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true},
	}

	httpRateLimiter.RateLimit(myHandler)

	//	redirectCode := http.StatusMovedPermanently
	//	if rc, err := cr.GetInt(config.Path(PathRedirectToBase)); rc != redirectCode && false == config.NotKeyNotFoundError(err) {
	//		redirectCode = http.StatusFound
	//	}

	return func(h ctxhttp.Handler) ctxhttp.Handler {
		return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}
