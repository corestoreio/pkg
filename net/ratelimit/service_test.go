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
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
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

	var runTest = func(logBuf *log.MutexBuffer, scp scope.Scope, id int64) func(t *testing.T) {
		return func(t *testing.T) {
			srv, err := ratelimit.New(
				ratelimit.WithLogger(logw.NewLog(logw.WithWriter(logBuf), logw.WithLevel(logw.LevelDebug))),
				//ratelimit.WithLogger(logw.NewLog(logw.WithWriter(ioutil.Discard), logw.WithLevel(logw.LevelDebug))),
				ratelimit.WithVaryBy(scp, id, pathGetter{}),
				ratelimit.WithRateLimiter(scp, id, stubLimiter{}),
				ratelimit.WithErrorHandler(scp, id, func(err error) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(500)
					})
				}),
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

			handler := srv.WithRateLimit()(finalHandler(t))

			runHTTPTestCases(t, handler, []httpTestCase{
				{"ok", 200, map[string]string{"X-Ratelimit-Limit": "1", "X-Ratelimit-Remaining": "2", "X-Ratelimit-Reset": "60"}},
				{"error", 500, map[string]string{}},
				{"limit", 429, map[string]string{"Retry-After": "60"}},
			})
		}
	}

	logBuf0 := new(log.MutexBuffer)
	t.Run("Scope Store Fallback to Default", runTest(logBuf0, scope.Default, 0))
	//t.Log("FallBack", logBuf0)
	cstesting.ContainsCount(t, logBuf0.String(), `Service.ConfigByScopedGetter.Fallback`, 1)

	logBuf1 := new(log.MutexBuffer)
	t.Run("Scope Store Fallback to Website", runTest(logBuf1, scope.Website, 1))
	//t.Log("FallBack", logBuf1)

	var logCheck1 = `Service.ConfigByScopedGetter.Fallback requested_scope: "Scope(Store) ID(1)" requested_fallback_scope: "Scope(Website) ID(1)" responded_scope: "Scope(Website) ID(1)`
	cstesting.ContainsCount(t, logBuf1.String(), logCheck1, 1)

	logBuf2 := new(log.MutexBuffer)
	t.Run("Scope Store No Fallback", runTest(logBuf2, scope.Store, 1))
	//t.Log("FallBackNope", logBuf2)

	var logCheck2 = `Service.ConfigByScopedGetter.IsValid requested_scope: "Scope(Store) ID(1)" requested_fallback_scope: "Scope(Absent) ID(0)" responded_scope: "Scope(Store) ID(1)"`
	cstesting.ContainsCount(t, logBuf2.String(), logCheck2, 150)

}

func TestService_WithDeniedHandler(t *testing.T) {

	var ac = new(int32)

	srv, err := ratelimit.New(
		ratelimit.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(500)
				// panic(fmt.Sprintf("Should not get called. Error Handler\n\n%+v", err))
			})
		}),
		ratelimit.WithVaryBy(scope.Default, 0, pathGetter{}),
		ratelimit.WithRateLimiter(scope.Default, 0, stubLimiter{}),
		ratelimit.WithDeniedHandler(scope.Default, 0, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			atomic.AddInt32(ac, 1)
			http.Error(w, "custom limit exceeded", 400)
		})),
	)
	if err != nil {
		t.Fatal(err)
	}
	srv.ErrorHandler = func(err error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(fmt.Sprintf("Should not get called. Root Error Handler\n\n%+v", err))
		})
	}

	handler := srv.WithRateLimit()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		panic("Should not get called this next handler")
	}))

	runHTTPTestCases(t, handler, []httpTestCase{
		{"limit", 400, map[string]string{}},
		{"error", 500, map[string]string{}},
	})
	if have, want := *ac, int32(runHTTPTestCasesUsers*runHTTPTestCasesLoops); have != want {
		t.Errorf("WithDeniedHandler call failed: Have: %d Want: %d", have, want)
	}
}

func TestService_RequestedStore_NotFound(t *testing.T) {
	srv, err := ratelimit.New(ratelimit.WithErrorHandler(scope.Default, 0,
		func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called. Scoped Error Handler\n\n%+v", err))
			})
		},
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
	srv.WithRateLimit()(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("Should not get called!")
	})).ServeHTTP(nil, req)
}

func TestService_ScopedConfig_NotFound(t *testing.T) {
	srv, err := ratelimit.New(ratelimit.WithErrorHandler(scope.Default, 0,
		func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called. Scoped Error Handler\n\n%+v", err))
			})
		},
	))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if err := srv.FlushCache(); err != nil {
		t.Fatalf("%+v", err)
	}
	srv.ErrorHandler = func(err error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Nil(t, w)
			assert.NotNil(t, r)
			if !errors.IsNotValid(err) {
				t.Errorf("Have: %+v Want: A not found error behaviour", err)
			}
		})
	}

	storeSrv := storemock.NewEurozzyService(scope.MustSetByCode(scope.Website, "euro"))
	req, _ := http.NewRequest("GET", "https://corestore.io", nil)
	st, err := storeSrv.Store(scope.MockID(1)) // German Store
	if err != nil {
		t.Fatalf("%+v", err)
	}
	st.Config = cfgmock.NewService().NewScoped(st.WebsiteID(), st.StoreID())
	req = req.WithContext(store.WithContextRequestedStore(req.Context(), st))

	srv.WithRateLimit()(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("Should not get called!")
	})).ServeHTTP(nil, req)
}

func TestService_WithDisabled(t *testing.T) {

	srv, err := ratelimit.New(
		ratelimit.WithVaryBy(scope.Default, 0, pathGetter{}),
		ratelimit.WithRateLimiter(scope.Default, 0, stubLimiter{}),
		ratelimit.WithDisable(scope.Default, 0, true),
	)
	if err != nil {
		t.Fatal(err)
	}

	handler := srv.WithRateLimit()(finalHandler(t))
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

		storeSrv := storemock.NewEurozzyService(scope.MustSetByCode(scope.Website, "euro"))
		req, _ := http.NewRequest("GET", c.path, nil)
		req.Header.Set("X-Forwarded-For", "2a02:d200::")
		st, err := storeSrv.Store(scope.MockID(1)) // German Store
		if err != nil {
			t.Fatalf("%+v", err)
		}
		st.Config = cfgmock.NewService().NewScoped(st.WebsiteID(), st.StoreID())
		req = req.WithContext(store.WithContextRequestedStore(req.Context(), st))

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
