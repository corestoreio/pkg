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

import "gopkg.in/throttled/throttled.v2"

// Option can be used as an argument in NewService to configure a token service.
type Option func(a *HTTPRateLimit)

// WithVaryBy ...
func WithVaryBy(vb VaryByer) Option {
	return func(s *HTTPRateLimit) {
		s.VaryByer = vb
	}
}

// WithRateLimiterForWebsite ...
func WithRateLimiterForWebsite(websiteID int64, rl throttled.RateLimiter) Option {
	return func(s *HTTPRateLimit) {
		s.mu.Lock()
		s.scopedRLs[websiteID] = rl
		s.mu.Unlock()
	}
}

// WithRateLimiterFactory ...
func WithRateLimiterFactory(rlf RateLimiterFactory) Option {
	return func(s *HTTPRateLimit) {
		s.RateLimiterFactory = rlf
	}
}
