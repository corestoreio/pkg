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

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
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
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	return withDefaultConfig(scp, id)
}

// WithVaryBy allows to set a custom key producer. VaryByer is called for each
// request to generate a key for the limiter. If it is nil, the middleware
// panics. The default VaryByer returns an empty string so that all requests
// uses the same key. VaryByer must be thread safe.
func WithVaryBy(scp scope.Scope, id int64, vb VaryByer) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.VaryByer = vb
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithRateLimiter creates a rate limiter for a specific scope with its ID.
// The rate limiter is already warmed up.
func WithRateLimiter(scp scope.Scope, id int64, rl throttled.RateLimiter) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.RateLimiter = rl
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithDeniedHandler sets a custom denied handler for a specific scope. The
// default denied handler returns a simple:
//		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
func WithDeniedHandler(scp scope.Scope, id int64, next http.Handler) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.DeniedHandler = next
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithDisable allows to disable a rate limit or enable it if set to false.
func WithDisable(scp scope.Scope, id int64, isDisabled bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.Disabled = isDisabled
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithLogger applies a logger to the default scope which gets inherited to
// subsequent scopes. Mainly used for debugging. Convenience helper function.
func WithLogger(l log.Logger) Option {
	return func(s *Service) error {
		s.Log = l
		return nil
	}
}

// WithGCRAStore creates a new GCRA rate limiter with a custom storage backend.
// Duration: (s second,i minute,h hour,d day)
// GCRA => https://en.wikipedia.org/wiki/Generic_cell_rate_algorithm
func WithGCRAStore(scp scope.Scope, id int64, store throttled.GCRAStore, duration rune, requests, burst int) Option {
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
		return WithRateLimiter(scp, id, rl)(s)
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
