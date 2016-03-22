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
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
	"gopkg.in/throttled/throttled.v2"
)

type stubLimiter struct{}

func (sl stubLimiter) RateLimit(key string, quantity int) (bool, throttled.RateLimitResult, error) {
	return false, throttled.RateLimitResult{}, nil
}

func TestWithScopedRateLimiter(t *testing.T) {
	t.Parallel()

	hashedScoped := scope.NewHash(scope.StoreID, 33)
	wantSL := stubLimiter{}
	rl, err := NewHTTPRateLimit(WithScopedRateLimiter(scope.StoreID, 33, wantSL))
	if err != nil {
		t.Fatal(err)
	}
	if haveSL, ok := rl.scopedRLs[hashedScoped]; ok {
		assert.Exactly(t, wantSL, haveSL)
	} else {
		t.Fatalf("Cannot find scoped rate limiter in map with hash key %d", hashedScoped)
	}
}
