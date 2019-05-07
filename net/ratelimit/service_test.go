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
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/net/mw"
	"github.com/corestoreio/pkg/net/ratelimit"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/util/assert"
	"gopkg.in/throttled/throttled.v2"
)

const errMessage = `stubLimiter TEST error`

type stubLimiter struct{}

func (sl stubLimiter) RateLimit(key string, quantity int) (bool, throttled.RateLimitResult, error) {
	switch key {
	case "limit":
		return true, throttled.RateLimitResult{-1, -1, -1, time.Minute}, nil
	case "error":
		return false, throttled.RateLimitResult{}, errors.NewFatalf(errMessage)
	case "panic":
		panic("RateLimit should not be called")
	default:
		return false, throttled.RateLimitResult{1, 2, time.Minute, -1}, nil
	}
}

type pathGetter struct{}

func (pathGetter) Key(r *http.Request) string {
	return strings.TrimLeft(r.URL.Path, "/")
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
	_ = ratelimit.MustNew(ratelimit.WithGCRAStore(nil, 'h', 2, -1, scope.DefaultTypeID))
}

type httpTestCase struct {
	path    string
	code    int
	headers map[string]string
}

var finalHandler = func(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`Everything OK dude!`))
	})
}

// TestService_WithRateLimit_ScopeStore1 runs without the backend configuration
// because we're using WithRateLimiter() and WithVaryBy() to set a rate limiter
// for a specific WebsiteID(1). Despite the request will come in with StoreID(1)
// we must fall back to the websiteID(1) to fetch there the configuration.
func TestService_WithRateLimit_StoreFallbackToWebsite(t *testing.T) {

	var runTest = func(logBuf io.Writer, scopeID scope.TypeID) func(t *testing.T) {
		return func(t *testing.T) {
			errH := mw.ErrorHandler(func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(500)
				})
			})

			srv, err := ratelimit.New(
				ratelimit.WithRootConfig(cfgmock.NewService()),
				ratelimit.WithDebugLog(logBuf),
				//ratelimit.WithLogger(logw.NewLog(logw.WithWriter(ioutil.Discard), logw.WithLevel(logw.LevelDebug))),
				ratelimit.WithVaryBy(pathGetter{}, scopeID),
				ratelimit.WithRateLimiter(stubLimiter{}, scopeID),
				ratelimit.WithErrorHandler(errH, scopeID),
			)
			if err != nil {
				t.Fatal(err)
			}
			srv.ErrorHandler = func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(500)
					panic(fmt.Sprintf("Root Scope\n%+v", err))
				})
			}

			handler := srv.WithRateLimit(finalHandler(t))

			runHTTPTestCases(t, handler, []httpTestCase{
				{"ok", 200, map[string]string{"X-Ratelimit-Limit": "1", "X-Ratelimit-Remaining": "2", "X-Ratelimit-Reset": "60"}},
				{"error", 500, map[string]string{}},
				{"limit", 429, map[string]string{"Retry-After": "60"}},
			})

			scpCfg, err := srv.ConfigByScopeID(scopeID, 0)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, scopeID, scpCfg.ScopeID, "ScopeID")
			cstesting.EqualPointers(t, errH, scpCfg.ErrorHandler)
		}
	}

	logBuf0 := new(log.MutexBuffer)
	t.Run("Scope Store Fallback to Default", runTest(logBuf0, scope.DefaultTypeID))
	//t.Log("FallBack", logBuf0)
	//cstesting.ContainsCount(t, logBuf0.String(), `Service.ConfigByScopedGetter.Fallback`, 1)
	logBuf0.Reset()

	t.Run("Scope Store Fallback to Website", runTest(logBuf0, scope.Website.WithID(1)))
	////t.Log("FallBack", logBuf1)
	//
	//var logCheck1 = `Service.ConfigByScopedGetter.Fallback requested_scope: "Scope(Store) ID(1)" requested_fallback_scope: "Scope(Website) ID(1)" responded_scope: "Scope(Website) ID(1)`
	//cstesting.ContainsCount(t, logBuf1.String(), logCheck1, 1)
	//
	//logBuf2 := new(log.MutexBuffer)
	t.Run("Scope Store No Fallback", runTest(logBuf0, scope.Store.WithID(1)))
	////t.Log("FallBackNope", logBuf2)
	//
	//var logCheck2 = `Service.ConfigByScopedGetter.IsValid requested_scope: "Scope(Store) ID(1)" requested_fallback_scope: "Scope(Absent) ID(0)" responded_scope: "Scope(Store) ID(1)"`
	//cstesting.ContainsCount(t, logBuf2.String(), logCheck2, 150)

}

func TestService_WithDeniedHandler(t *testing.T) {

	var deniedHandlerCalled = new(int32)

	srv, err := ratelimit.New(
		ratelimit.WithRootConfig(cfgmock.NewService()),
		ratelimit.WithErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(500)
				// panic(fmt.Sprintf("Should not get called. Error Handler\n\n%+v", err))
			})
		}, scope.DefaultTypeID),
		ratelimit.WithVaryBy(pathGetter{}, scope.DefaultTypeID),
		ratelimit.WithRateLimiter(stubLimiter{}, scope.DefaultTypeID),
		ratelimit.WithDeniedHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			atomic.AddInt32(deniedHandlerCalled, 1)
			http.Error(w, "custom limit exceeded", 400)
		}), scope.DefaultTypeID),
	)
	if err != nil {
		t.Fatal(err)
	}
	srv.ErrorHandler = func(err error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(fmt.Sprintf("Should not get called. Root Error Handler\n\n%+v", err))
		})
	}

	handler := srv.WithRateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called this next handler")
	}))

	runHTTPTestCases(t, handler, []httpTestCase{
		{"limit", 400, map[string]string{}},
		{"error", 500, map[string]string{}},
	})
	if have, want := *deniedHandlerCalled, int32(runHTTPTestCasesUsers*runHTTPTestCasesLoops); have != want {
		t.Errorf("WithDeniedHandler call failed: Have: %d Want: %d", have, want)
	}
}

func TestService_RequestedStore_NotFound(t *testing.T) {
	srv, err := ratelimit.New(ratelimit.WithErrorHandler(
		func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called. Scoped Error Handler\n\n%+v", err))
			})
		},
		scope.DefaultTypeID,
	))
	if err != nil {
		t.Fatal(err)
	}
	srv.ErrorHandler = func(err error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Nil(t, w)
			assert.NotNil(t, r)
			if !errors.IsNotFound(err) {
				t.Errorf("Have: %+v Want: A not found error behaviour", err)
			}
		})
	}

	req := httptest.NewRequest("GET", "/", nil)
	srv.WithRateLimit(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("Should not get called!")
	})).ServeHTTP(nil, req)
}

func TestService_ScopedConfig_NotFound(t *testing.T) {
	srv, err := ratelimit.New(
		ratelimit.WithRootConfig(cfgmock.NewService()),
		ratelimit.WithErrorHandler(
			func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic(fmt.Sprintf("Should not get called. Scoped Error Handler\n\n%+v", err))
				})
			},
			scope.DefaultTypeID),
	)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if err := srv.ClearCache(); err != nil {
		t.Fatalf("%+v", err)
	}
	srv.ErrorHandler = func(err error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Nil(t, w)
			assert.NotNil(t, r)
			if !errors.IsNotFound(err) {
				t.Errorf("Have: %+v Want: A not found error behaviour", err)
			}
		})
	}

	req := httptest.NewRequest("GET", "https://corestore.io", nil)
	req = req.WithContext(scope.WithContext(req.Context(), 1, 1)) // website=euro store=german

	srv.WithRateLimit(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("Should not get called!")
	})).ServeHTTP(nil, req)
}

func TestService_WithDisabled(t *testing.T) {

	srv, err := ratelimit.New(
		ratelimit.WithRootConfig(cfgmock.NewService()),
		ratelimit.WithVaryBy(pathGetter{}, scope.DefaultTypeID),
		ratelimit.WithRateLimiter(stubLimiter{}, scope.DefaultTypeID),
		ratelimit.WithDisable(true, scope.DefaultTypeID),
	)
	if err != nil {
		t.Fatal(err)
	}

	handler := srv.WithRateLimit(finalHandler(t))
	runHTTPTestCases(t, handler, []httpTestCase{
		{"panic", 200, map[string]string{}},
	})
}

const (
	runHTTPTestCasesUsers = 10
	runHTTPTestCasesLoops = 5
)

func runHTTPTestCases(t *testing.T, h http.Handler, cs []httpTestCase) {
	for i, c := range cs {

		req := httptest.NewRequest("GET", "/"+c.path, nil)
		req.Header.Set("X-Forwarded-For", "2a02:d200::")
		req = req.WithContext(scope.WithContext(req.Context(), 1, 1))

		hpu := cstesting.NewHTTPParallelUsers(runHTTPTestCasesUsers, runHTTPTestCasesLoops, 200, time.Millisecond)
		hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
			if have, want := rec.Code, c.code; have != want {
				t.Errorf("Expected request %d at %s to return %d but got %d",
					i, c.path, want, have)
			}

			for name, want := range c.headers {
				if have := rec.HeaderMap.Get(name); have != want {
					t.Errorf("Expected request %d at %s to have header '%s: %s' but got '%s'",
						i, c.path, name, want, have)
				}
			}
		}
		hpu.ServeHTTP(req, h)
	}
}
