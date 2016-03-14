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
	"net/http"
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
	"gopkg.in/throttled/throttled.v2"

	"golang.org/x/net/context"
)

type RateLimiterFactory func(*PkgBackend, config.ScopedGetter) (throttled.RateLimiter, error)

// VaryByer is called for each request to generate a key for the
// limiter. If it is nil, all requests use an empty string key.
type VaryByer interface {
	Key(*http.Request) string
}

// HTTPRateLimit faciliates using a Limiter to limit HTTP requests.
type HTTPRateLimit struct {
	me *cserr.MultiErr
	// Backend configuration, if nil everything panics.
	be *PkgBackend

	// DeniedHandler can be customized instead of showing a HTTP status 429
	// error page once the HTTPRateLimit has been reached.
	// It will be called if the request gets over the limit.
	DeniedHandler ctxhttp.HandlerFunc

	// RateLimitFactory creates a new rate limiter for each scope
	RateLimiterFactory

	// VaryByer is called for each request to generate a key for the
	// limiter. If it is nil, all requests use an empty string key.
	VaryByer

	mu sync.RWMutex
	// scopedRLs internal cache of already created rate limiter with their
	// storage. ID relates to the website ID.
	// Due to the overall nature I assume that the rate limit is the bottleneck
	// for an application instead of this mutex protected map.
	scopedRLs map[int64]throttled.RateLimiter
}

func NewHTTPRateLimit(be *PkgBackend, opts ...Option) (*HTTPRateLimit, error) {
	if be == nil {
		return nil, errors.New("PkgBackend cannot be nil")
	}

	rl := &HTTPRateLimit{
		be:        be,
		scopedRLs: make(map[int64]throttled.RateLimiter),
		DeniedHandler: func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return nil
		},
	}

	if err := rl.Options(opts...); err != nil {
		return nil, err
	}

	if rl.RateLimiterFactory == nil {
		rl.RateLimiterFactory = NewGCRAMemStore(MemStoreMaxKeys)
	}
	return rl, nil
}

// Options applies option at creation time or refreshes them.
func (s *HTTPRateLimit) Options(opts ...Option) error {
	for _, opt := range opts {
		opt(s)
	}
	if s.me.HasErrors() {
		return s
	}
	return nil
}
