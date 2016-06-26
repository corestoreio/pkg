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
	"testing"

	"fmt"

	"net/http"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/throttled/throttled.v2"
)

type stubLimiter struct{}

func (sl stubLimiter) RateLimit(key string, quantity int) (bool, throttled.RateLimitResult, error) {
	return false, throttled.RateLimitResult{}, nil
}

func TestCalculateRate(t *testing.T) {
	tests := []struct {
		duration rune
		requests int
		wantRate throttled.Rate
	}{
		{'s', 11, throttled.PerSec(11)},
		{'i', 22, throttled.PerMin(22)},
		{'h', 33, throttled.PerHour(33)},
		{'d', 44, throttled.PerDay(44)},
		{'y', 55, throttled.PerHour(55)},
	}
	for _, test := range tests {
		haveR := CalculateRate(test.duration, test.requests)
		assert.Exactly(t, test.wantRate, haveR)
	}
}

func TestWithDefaultConfig(t *testing.T) {

	s := MustNew(WithDefaultConfig(scope.Store, 33))
	s33 := scope.NewHash(scope.Store, 33)
	want33 := defaultScopedConfig(s33)
	want0 := defaultScopedConfig(scope.DefaultHash)

	// poor mans comparison function. better solution? Before suggesting please test it :-)
	assert.Exactly(t, fmt.Sprintf("%#v", want33), fmt.Sprintf("%#v", s.scopeCache[s33]))
	assert.Exactly(t, fmt.Sprintf("%#v", want0), fmt.Sprintf("%#v", s.defaultScopeCache))
}

func TestWithVaryBy(t *testing.T) {
	vb := new(VaryBy)
	s33 := scope.NewHash(scope.Store, 33)

	t.Run("Ok", func(t *testing.T) {
		s := MustNew(
			WithDefaultConfig(scope.Store, 33),
			WithVaryBy(scope.Store, 33, vb),
			WithVaryBy(scope.Default, 0, vb),
		)
		assert.Exactly(t, vb, s.scopeCache[s33].VaryByer)
		assert.Exactly(t, vb, s.defaultScopeCache.VaryByer)
	})
	t.Run("OverwrittenByWithDefaultConfig", func(t *testing.T) {
		s := MustNew(
			WithVaryBy(scope.Store, 33, vb),
			WithDefaultConfig(scope.Store, 33),
		)
		// WithDefaultConfig overwrites the previously set VaryBy
		assert.Exactly(t, emptyVaryBy{}, s.scopeCache[s33].VaryByer)
	})
}

func TestWithRateLimiter(t *testing.T) {
	rsl := stubLimiter{}
	w2 := scope.NewHash(scope.Website, 2)

	t.Run("Ok", func(t *testing.T) {
		s := MustNew(
			WithDefaultConfig(scope.Website, 2),
			WithRateLimiter(scope.Website, 2, rsl),
			WithRateLimiter(scope.Default, 0, rsl),
		)
		assert.Exactly(t, rsl, s.scopeCache[w2].RateLimiter)
		assert.Exactly(t, rsl, s.defaultScopeCache.RateLimiter)
	})
	t.Run("OverwrittenByWithDefaultConfig", func(t *testing.T) {
		s := MustNew(
			WithRateLimiter(scope.Website, 2, rsl),
			WithDefaultConfig(scope.Website, 2),
		)
		// WithDefaultConfig overwrites the previously set RateLimiter
		assert.Nil(t, s.scopeCache[w2].RateLimiter)
		err := s.getConfigByScopeID(w2, false).isValid()
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})
}

func TestWithDeniedHandler(t *testing.T) {
	dh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInsufficientStorage)
	})
	w2 := scope.NewHash(scope.Website, 2)

	t.Run("Ok", func(t *testing.T) {
		s := MustNew(
			WithDefaultConfig(scope.Website, 2),
			WithDeniedHandler(scope.Website, 2, dh),
			WithDeniedHandler(scope.Default, 0, dh),
		)
		cstesting.EqualPointers(t, dh, s.scopeCache[w2].deniedHandler)
		cstesting.EqualPointers(t, dh, s.defaultScopeCache.deniedHandler)
	})
	t.Run("OverwrittenByWithDefaultConfig", func(t *testing.T) {
		s := MustNew(
			WithDeniedHandler(scope.Website, 2, dh),
			WithDefaultConfig(scope.Website, 2),
		)
		// WithDefaultConfig overwrites the previously set RateLimiter
		cstesting.EqualPointers(t, defaultDeniedHandler, s.scopeCache[w2].deniedHandler)
		err := s.getConfigByScopeID(w2, false).isValid()
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})
}
