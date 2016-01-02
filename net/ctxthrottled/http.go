// Copyright (c) 2014, Martin Angers and Contributors.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// * Neither the name of the author nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"golang.org/x/net/context"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

const (
	// PathRateLimitBurst defines the number of requests that
	// will be allowed to exceed the rate in a single burst and must be
	// greater than or equal to zero.
	// Scope Global, Type Int.
	PathRateLimitBurst = `corestore/ctxthrottled/burst`
	// PathRateLimitRequests number of requests allowed per time period
	// Scope Global, Type Int.
	PathRateLimitRequests = `corestore/ctxthrottled/requests`
	// PathRateLimitDuration per second (s), minute (i), hour (h), day (d)
	// Scope Global, Type String.
	PathRateLimitDuration = `corestore/ctxthrottled/duration`
)

// DefaultDeniedHandler is the default DeniedHandler for an
// HTTPRateLimit. It returns a 429 status code with a generic
// message.
var DefaultDeniedHandler = ctxhttp.Handler(ctxhttp.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
	http.Error(w, "limit exceeded", 429)
	return nil
}))

// DefaultBurst defines the number of requests that
// will be allowed to exceed the rate in a single burst and must be
// greater than or equal to zero.
var DefaultBurst int = 5

// DefaultRequests number of requests allowed per time period
var DefaultRequests int = 100

// DefaultDuration per second (s), minute (i), hour (h), day (d)
var DefaultDuration string = "h"

// HTTPRateLimit faciliates using a Limiter to limit HTTP requests.
type HTTPRateLimit struct {
	// Config is the config.Service with PubSub
	Config config.GetterPubSuber

	// configScope config.ScopedReader todo for later: rate limit on a per website level

	// DeniedHandler is called if the request is disallowed. If it is
	// nil, the DefaultDeniedHandler variable is used.
	DeniedHandler ctxhttp.Handler

	// Limiter is call for each request to determine whether the
	// request is permitted and update internal state. It must be set.
	RateLimiter throttled.RateLimiter

	// VaryBy is called for each request to generate a key for the
	// limiter. If it is nil, all requests use an empty string key.
	VaryBy interface {
		Key(*http.Request) string
	}
}

func (t *HTTPRateLimit) quota() throttled.RateQuota {
	var burst, request int
	var duration string

	if burst, _ = t.Config.Int(config.Path(PathRateLimitBurst)); burst < 0 {
		burst = DefaultBurst
	}
	if request, _ = t.Config.Int(config.Path(PathRateLimitRequests)); request == 0 {
		request = DefaultRequests
	}
	if duration, _ = t.Config.String(config.Path(PathRateLimitDuration)); duration == "" {
		duration = DefaultDuration
	}

	var r throttled.Rate
	switch duration {
	case "s": // second
		r = throttled.PerSec(request)
	case "i": // minute
		r = throttled.PerMin(request)
	case "h": // hour
		r = throttled.PerHour(request)
	case "d": // day
		r = throttled.PerDay(request)
	default:
		r = throttled.PerHour(request)
	}

	return throttled.RateQuota{r, burst}
}

// WithRateLimit wraps an ctxhttp.Handler to limit incoming requests.
// Requests that are not limited will be passed to the handler
// unchanged.  Limited requests will be passed to the DeniedHandler.
// X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset and
// Retry-After headers will be written to the response based on the
// values in the RateLimitResult.
func (t *HTTPRateLimit) WithRateLimit(rlStore throttled.GCRAStore, h ctxhttp.Handler) ctxhttp.Handler {
	if t.Config == nil {
		t.Config = config.DefaultService
	}
	if t.DeniedHandler == nil {
		t.DeniedHandler = DefaultDeniedHandler
	}

	if t.RateLimiter == nil {
		if rlStore == nil {
			var err error
			rlStore, err = memstore.New(65536)
			if err != nil {
				panic(err)
			}
		}

		var err error
		t.RateLimiter, err = throttled.NewGCRARateLimiter(rlStore, t.quota())
		if err != nil {
			panic(err)
		}
	}

	return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		var k string
		if t.VaryBy != nil {
			k = t.VaryBy.Key(r)
		}

		limited, context, err := t.RateLimiter.RateLimit(k, 1)

		if err != nil {
			return err
		}

		setRateLimitHeaders(w, context)

		if !limited {
			return h.ServeHTTPContext(ctx, w, r)
		}

		return t.DeniedHandler.ServeHTTPContext(ctx, w, r)
	})
}

func setRateLimitHeaders(w http.ResponseWriter, context throttled.RateLimitResult) {
	if v := context.Limit; v >= 0 {
		w.Header().Add("X-RateLimit-Limit", strconv.Itoa(v))
	}

	if v := context.Remaining; v >= 0 {
		w.Header().Add("X-RateLimit-Remaining", strconv.Itoa(v))
	}

	if v := context.ResetAfter; v >= 0 {
		vi := int(math.Ceil(v.Seconds()))
		w.Header().Add("X-RateLimit-Reset", strconv.Itoa(vi))
	}

	if v := context.RetryAfter; v >= 0 {
		vi := int(math.Ceil(v.Seconds()))
		w.Header().Add("Retry-After", strconv.Itoa(vi))
	}
}
