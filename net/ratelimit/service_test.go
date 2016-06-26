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

package ratelimit_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/throttled/throttled.v2"
)

type stubLimiter struct {
}

func (sl stubLimiter) RateLimit(key string, quantity int) (bool, throttled.RateLimitResult, error) {
	switch key {
	case "limit":
		return true, throttled.RateLimitResult{-1, -1, -1, time.Minute}, nil
	case "error":
		return false, throttled.RateLimitResult{}, errors.New("stubLimiter error")
	default:
		return false, throttled.RateLimitResult{1, 2, time.Minute, -1}, nil
	}
}

type pathGetter struct{}

func (pathGetter) Key(r *http.Request) string {
	return r.URL.Path
}

var _ ratelimit.VaryByer = (*pathGetter)(nil)

func TestMustNew_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
			} else {
				t.Fatal("Expecting an error")
			}
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	_ = ratelimit.MustNew(ratelimit.WithGCRAStore(scope.Default, 0, nil, 'h', 2, -1))
}

//
//type httpTestCase struct {
//	path    string
//	code    int
//	headers map[string]string
//}
//
//var finalHandler200 = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) error {
//	w.WriteHeader(200)
//	return nil
//})
//
//func TestHTTPRateLimit_WithoutConfig(t *testing.T) {
//	t.Parallel()
//	// this test case runs without the backend configuration because
//	// we're using WithScopedRateLimiter() to set a rate limiter for a specific
//	// website (ID 1). In real life you must create a rate limiter for each website
//	// or we can implement a configurable pass-through option which bypasses the RL. that
//	// means if a rate limiter cannot be found the next HTTP handler gets called.
//
//	limiter, err := ratelimit.NewService(
//		ratelimit.WithVaryBy(pathGetter{}),
//		ratelimit.WithScopedRateLimiter(scope.Website, 1, stubLimiter{}), // 1 = NewEurozzyService() website euro
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	ctx := store.WithContextProvider(
//		context.Background(),
//		storemock.NewEurozzyService(scope.MustSetByCode(scope.Website, "euro")),
//	)
//
//	handler := limiter.WithRateLimit()(finalHandler200)
//
//	runHTTPTestCases(t, ctx, handler, []httpTestCase{
//		{"ok", 200, map[string]string{"X-Ratelimit-Limit": "1", "X-Ratelimit-Remaining": "2", "X-Ratelimit-Reset": "60"}},
//		{"error", 500, map[string]string{}},
//		{"limit", 429, map[string]string{"Retry-After": "60"}},
//	})
//}
//
//func TestHTTPRateLimit_CustomHandlers(t *testing.T) {
//	t.Parallel()
//	limiter, err := ratelimit.NewService(
//		ratelimit.WithVaryBy(pathGetter{}),
//		ratelimit.WithRateLimiterFactory(newStubLimiter(nil)),
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	limiter.DeniedHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) error {
//		http.Error(w, "custom limit exceeded", 400)
//		return nil
//	})
//
//	ctx := store.WithContextProvider(
//		context.Background(),
//		storemock.NewEurozzyService(
//			scope.MustSetByCode(scope.Website, "euro"),
//		),
//	)
//
//	handler := limiter.WithRateLimit()(finalHandler200)
//
//	runHTTPTestCases(t, ctx, handler, []httpTestCase{
//		{"limit", 400, map[string]string{}},
//		{"error", 500, map[string]string{}},
//	})
//}
//
//func TestHTTPRateLimit_RateLimiterFactoryError(t *testing.T) {
//	t.Parallel()
//
//	testedErr := errors.New("RateLimiterFactory Error")
//
//	limiter, err := ratelimit.New(
//		ratelimit.WithVaryBy(pathGetter{}),
//		ratelimit.WithRateLimiterFactory(newStubLimiter(testedErr)),
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	ctx := store.WithContextProvider(
//		context.Background(),
//		storemock.NewEurozzyService(
//			scope.MustSetByCode(scope.Website, "euro"),
//		),
//	)
//
//	req, err := http.NewRequest("GET", "/", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = limiter.WithRateLimit()(finalHandler200).ServeHTTP( nil, req)
//	assert.EqualError(t, err, testedErr.Error())
//}
//
//func TestHTTPRateLimit_MissingContext(t *testing.T) {
//	t.Parallel()
//
//	limiter, err := ratelimit.New(
//		ratelimit.WithVaryBy(pathGetter{}),
//		ratelimit.WithRateLimiterFactory(newStubLimiter(nil)),
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	ctx := store.WithContextProvider(
//		context.Background(),
//		nil,
//	)
//
//	req, err := http.NewRequest("GET", "/", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = limiter.WithRateLimit()(finalHandler200).ServeHTTPContext(ctx, nil, req)
//	assert.EqualError(t, err, store.ErrContextProviderNotFound.Error())
//}
//
//func runHTTPTestCases(t *testing.T, ctx context.Context, h ctxhttp.Handler, cs []httpTestCase) {
//	for i, c := range cs {
//		req, err := http.NewRequest("GET", c.path, nil)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		rr := httptest.NewRecorder()
//		if err := h.ServeHTTPContext(ctx, rr, req); err != nil {
//			http.Error(rr, err.Error(), http.StatusInternalServerError)
//		}
//
//		if have, want := rr.Code, c.code; have != want {
//			t.Errorf("Expected request %d at %s to return %d but got %d",
//				i, c.path, want, have)
//		}
//
//		for name, want := range c.headers {
//			if have := rr.HeaderMap.Get(name); have != want {
//				t.Errorf("Expected request %d at %s to have header '%s: %s' but got '%s'",
//					i, c.path, name, want, have)
//			}
//		}
//	}
//}
