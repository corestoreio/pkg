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
	"fmt"
	"net/http"
	"testing"

	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

type stubLimiter struct{}

func (sl stubLimiter) RateLimit(key string, quantity int) (bool, throttled.RateLimitResult, error) {
	return false, throttled.RateLimitResult{}, nil
}

func TestCalculateRate(t *testing.T) {
	tests := []struct {
		duration   rune
		requests   int
		wantRate   throttled.Rate
		wantErrBhf errors.BehaviourFunc
	}{
		{'s', 11, throttled.PerSec(11), nil},
		{'i', 22, throttled.PerMin(22), nil},
		{'h', 33, throttled.PerHour(33), nil},
		{'d', 44, throttled.PerDay(44), nil},
		{'y', 55, throttled.Rate{}, errors.IsNotValid},
	}
	for _, test := range tests {
		haveR, err := calculateRate(test.duration, test.requests)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(err), "%+v", err)
		}
		assert.Exactly(t, test.wantRate, haveR)
	}
}

func TestWithDefaultConfig(t *testing.T) {

	s := MustNew(WithDefaultConfig(scope.Store.WithID(33)))
	s33 := scope.Store.WithID(33)
	want33 := newScopedConfig(s33, scope.DefaultTypeID)
	want0 := newScopedConfig(scope.DefaultTypeID, scope.DefaultTypeID)

	// poor mans comparison function. better solution? Before suggesting please test it :-)
	assert.Exactly(t, fmt.Sprintf("%#v", want33), fmt.Sprintf("%#v", s.scopeCache[s33]))
	assert.Exactly(t, fmt.Sprintf("%#v", want0), fmt.Sprintf("%#v", s.scopeCache[scope.DefaultTypeID]))
}

func TestWithVaryBy(t *testing.T) {
	vb := new(VaryBy)
	s33 := scope.MakeTypeID(scope.Store, 33)

	t.Run("Ok", func(t *testing.T) {
		s := MustNew(
			WithDefaultConfig(scope.Store.WithID(33)),
			WithVaryBy(vb, scope.Store.WithID(33)),
			WithVaryBy(vb, scope.Default.WithID(0)),
		)
		assert.Exactly(t, vb, s.scopeCache[s33].VaryByer)
		assert.Exactly(t, vb, s.scopeCache[scope.DefaultTypeID].VaryByer)
	})

	//TODO	move the following test into scopedservice package

	t.Run("OverwrittenByWithDefaultConfig", func(t *testing.T) {
		s := MustNew(
			WithVaryBy(vb, scope.Store.WithID(33)),
			WithDefaultConfig(scope.Store.WithID(33)),
		)
		// WithDefaultConfig overwrites the previously set VaryBy
		assert.Exactly(t, emptyVaryBy{}, s.scopeCache[s33].VaryByer)
	})
}

func TestWithRateLimiter(t *testing.T) {
	rsl := stubLimiter{}
	w2 := scope.MakeTypeID(scope.Website, 2)

	t.Run("Ok", func(t *testing.T) {
		s := MustNew(
			WithDefaultConfig(scope.Website.WithID(2)),
			WithRateLimiter(rsl, scope.Website.WithID(2)),
			WithRateLimiter(rsl, scope.Default.WithID(0)),
		)
		assert.Exactly(t, rsl, s.scopeCache[w2].RateLimiter)
		assert.Exactly(t, rsl, s.scopeCache[scope.DefaultTypeID].RateLimiter)
	})
	t.Run("OverwrittenByWithDefaultConfig", func(t *testing.T) {
		s := MustNew(
			WithRateLimiter(rsl, scope.Website.WithID(2)),
			WithDefaultConfig(scope.Website.WithID(2)),
		)
		// WithDefaultConfig overwrites the previously set RateLimiter
		assert.Nil(t, s.scopeCache[w2].RateLimiter)
		_, err := s.ConfigByScopeID(w2, 0)
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})
}

func TestWithDeniedHandler(t *testing.T) {
	dh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInsufficientStorage)
	})
	w2 := scope.MakeTypeID(scope.Website, 2)

	t.Run("Ok", func(t *testing.T) {
		s := MustNew(
			WithDefaultConfig(scope.Website.WithID(2)),
			WithDeniedHandler(dh, scope.Website.WithID(2)),
			WithDeniedHandler(dh, scope.Default.WithID(0)),
		)
		cstesting.EqualPointers(t, dh, s.scopeCache[w2].DeniedHandler)
		cstesting.EqualPointers(t, dh, s.scopeCache[scope.DefaultTypeID].DeniedHandler)
	})
	t.Run("OverwrittenByWithDefaultConfig", func(t *testing.T) {
		s := MustNew(
			WithDeniedHandler(dh, scope.Website.WithID(2)),
			WithDefaultConfig(scope.Website.WithID(2)),
		)
		// WithDefaultConfig overwrites the previously set RateLimiter
		cstesting.EqualPointers(t, DefaultDeniedHandler, s.scopeCache[w2].DeniedHandler)
		_, err := s.ConfigByScopeID(w2, 0)
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})
}

func TestWithGCRAStore(t *testing.T) {
	w2 := scope.Website.WithID(2)

	memStore, err := memstore.New(40)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("CalcError", func(t *testing.T) {
		s, err := New(WithGCRAStore(nil, 's', 33, -1, scope.Website.WithID(2)))
		assert.Nil(t, s)
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})

	t.Run("Ok", func(t *testing.T) {
		s := MustNew(
			WithDefaultConfig(scope.Website.WithID(2)),
			WithGCRAStore(memStore, 's', 100, 10, scope.Website.WithID(2)),
			WithGCRAStore(memStore, 'h', 100, 10, scope.Default.WithID(0)),
		)
		assert.NotNil(t, s.scopeCache[w2].RateLimiter)
		assert.NotNil(t, s.scopeCache[scope.DefaultTypeID].RateLimiter)
	})

	t.Run("OverwrittenByWithDefaultConfig", func(t *testing.T) {
		s := MustNew(
			WithGCRAStore(memStore, 's', 100, 10, scope.Website.WithID(2)),
			WithDefaultConfig(scope.Website.WithID(2)),
		)
		assert.Nil(t, s.scopeCache[w2].RateLimiter)
		_, err := s.ConfigByScopeID(w2, 0)
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})
}
