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

package ctxthrottled

import (
	"errors"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"golang.org/x/net/context"
	"gopkg.in/throttled/throttled.v2"
	"math"
	"net/http"
	"strconv"
)

var (
	// DefaultDeniedHandler is the default DeniedHandler for an
	// HTTPRateLimiter. It returns a 429 status code with a generic
	// message.
	DefaultDeniedHandler = ctxhttp.Handler(ctxhttp.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
		http.Error(w, "limit exceeded", 429)
		return nil
	}))

	// DefaultError is the default Error function for an HTTPRateLimiter.
	// It returns a 500 status code with a generic message.
	DefaultError = func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return nil
	}
)

// HTTPRateLimiter faciliates using a Limiter to limit HTTP requests.
type HTTPRateLimiter struct {
	// DeniedHandler is called if the request is disallowed. If it is
	// nil, the DefaultDeniedHandler variable is used.
	DeniedHandler ctxhttp.Handler

	// Error is called if the RateLimiter returns an error. If it is
	// nil, the DefaultErrorFunc is used.
	Error func(w http.ResponseWriter, r *http.Request, err error)

	// Limiter is call for each request to determine whether the
	// request is permitted and update internal state. It must be set.
	RateLimiter throttled.RateLimiter

	// VaryBy is called for each request to generate a key for the
	// limiter. If it is nil, all requests use an empty string key.
	VaryBy interface {
		Key(*http.Request) string
	}
}

// RateLimit wraps an ctxhttp.Handler to limit incoming requests.
// Requests that are not limited will be passed to the handler
// unchanged.  Limited requests will be passed to the DeniedHandler.
// X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset and
// Retry-After headers will be written to the response based on the
// values in the RateLimitResult.
func (t *HTTPRateLimiter) RateLimit(h ctxhttp.Handler) ctxhttp.Handler {
	return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if t.RateLimiter == nil {
			t.error(w, r, errors.New("You must set a RateLimiter on HTTPRateLimiter"))
		}

		var k string
		if t.VaryBy != nil {
			k = t.VaryBy.Key(r)
		}

		limited, context, err := t.RateLimiter.RateLimit(k, 1)

		if err != nil {
			t.error(w, r, err)
			return
		}

		setRateLimitHeaders(w, context)

		if !limited {
			h.ServeHTTPContext(ctx, w, r)
		} else {
			dh := t.DeniedHandler
			if dh == nil {
				dh = DefaultDeniedHandler
			}
			dh.ServeHTTPContext(ctx, w, r)
		}
	})
}

func (t *HTTPRateLimiter) error(w http.ResponseWriter, r *http.Request, err error) {
	e := t.Error
	if e == nil {
		e = DefaultError
	}
	e(w, r, err)
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
