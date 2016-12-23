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
	"net/http"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/errors"
	"gopkg.in/throttled/throttled.v2"
)

// WithDefaultConfig applies the default ratelimit configuration settings based
// for a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- Denied Handler: http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
//		- VaryByer: returns an empty key
// Example:
//		s := MustNewService(WithDefaultConfig(scope.Store,1), WithVaryBy(scope.Store, 1, myVB))
func WithDefaultConfig(id scope.TypeID) Option {
	return withDefaultConfig(id)
}

// WithVaryBy allows to set a custom key producer. VaryByer is called for each
// request to generate a key for the limiter. If it is nil, the middleware
// panics. The default VaryByer returns an empty string so that all requests
// uses the same key. VaryByer must be thread safe.
func WithVaryBy(vb VaryByer, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.VaryByer = vb
		return s.updateScopedConfig(sc)
	}
}

// WithRateLimiter creates a rate limiter for a specific scope with its ID.
// The rate limiter is already warmed up.
func WithRateLimiter(rl throttled.RateLimiter, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.RateLimiter = rl
		return s.updateScopedConfig(sc)
	}
}

// WithDeniedHandler sets a custom denied handler for a specific scope. The
// default denied handler returns a simple:
//		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
func WithDeniedHandler(next http.Handler, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.DeniedHandler = next
		return s.updateScopedConfig(sc)
	}
}

// WithGCRAStore creates a new GCRA rate limiter with a custom storage backend.
// Duration: (s second,i minute,h hour,d day)
// GCRA => https://en.wikipedia.org/wiki/Generic_cell_rate_algorithm
func WithGCRAStore(store throttled.GCRAStore, duration rune, requests, burst int, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {

		cr, err := calculateRate(duration, requests)
		if err != nil {
			return errors.Wrap(err, "[ratelimit] WithGCRAStore.calculateRate")
		}

		rq := throttled.RateQuota{
			MaxRate:  cr,
			MaxBurst: burst,
		}

		rl, err := throttled.NewGCRARateLimiter(store, rq)
		if err != nil {
			return errors.NewNotValidf("[ratelimit] throttled.NewGCRARateLimiter: %s", err)
		}
		return WithRateLimiter(rl, scopeIDs...)(s)
	}
}

// calculateRate calculates the rate depending on the duration (s second,i minute,h hour,d day) and the
// maximum requests. Invalid duration returns a NotValid error.
func calculateRate(duration rune, requests int) (r throttled.Rate, err error) {
	switch duration {
	case 's': // second
		r = throttled.PerSec(requests)
	case 'i': // minute
		r = throttled.PerMin(requests)
	case 'h': // hour
		r = throttled.PerHour(requests)
	case 'd': // day
		r = throttled.PerDay(requests)
	default:
		err = errors.NewNotValidf(errUnknownDurationRune, string(duration), requests)
	}
	return
}
