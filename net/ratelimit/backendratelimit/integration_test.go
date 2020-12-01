// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package backendratelimit_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/net/mw"
	"github.com/corestoreio/pkg/net/ratelimit"
	"github.com/corestoreio/pkg/net/ratelimit/memstore"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/cstesting"
	"github.com/throttled/throttled/v2"
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

func TestBackend_GCRA_Not_Registered(t *testing.T) {
	testBackendConfiguration(t, "panic",
		cfgmock.PathValue{},
		func(rec *httptest.ResponseRecorder) {
			assert.Exactly(t, http.StatusTeapot, rec.Code)
		},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("Should not get called")
		}),
		false, // do not test logger
		ratelimit.WithVaryBy(pathGetter{}, scope.Website.WithID(1)),
		ratelimit.WithRateLimiter(stubLimiter{}, scope.Website.WithID(1)),
	)
}

func TestBackend_WithDisable(t *testing.T) {
	testBackendConfiguration(t, "panic",
		cfgmock.PathValue{
			backend.Disabled.MustFQWebsite(1): 1,
		},
		func(rec *httptest.ResponseRecorder) {
			assert.Exactly(t, http.StatusTeapot, rec.Code)
		},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}),
		true, // do test logger
		ratelimit.WithVaryBy(pathGetter{}, scope.Website.WithID(1)),
		ratelimit.WithRateLimiter(stubLimiter{}, scope.Website.WithID(1)),
	)
}

func TestBackend_WithGCRAMemStore(t *testing.T) {
	countDenied := new(int32)
	countAllowed := new(int32)

	backend.Register(memstore.NewOptionFactory(backend.Burst, backend.Requests, backend.Duration, backend.StorageGCRAMaxMemoryKeys))
	defer backend.Deregister(memstore.OptionName)

	deniedH := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(countDenied, 1)
		http.Error(w, "custom limit exceeded", http.StatusConflict)
	})

	// fmt.Printf("deniedH: Default: %#v CustomTest: %#v\n", ratelimit.DefaultDeniedHandler, deniedH)

	testBackendConfiguration(t,
		"http://corestore.io",
		cfgmock.PathValue{
			backend.GCRAName.MustFQWebsite(1):                 "memstore",
			backend.Disabled.MustFQWebsite(1):                 0,
			backend.StorageGCRAMaxMemoryKeys.MustFQWebsite(1): 50,
			backend.Burst.MustFQWebsite(1):                    3,
			backend.Requests.MustFQWebsite(1):                 1,
			backend.Duration.MustFQWebsite(1):                 "i",
		},
		func(rec *httptest.ResponseRecorder) {
			//fmt.Printf("Code %d Remain:%s Limit:%s Reset:%s\n",
			//	rec.Code,
			//	rec.Header().Get("X-RateLimit-Remaining"),
			//	rec.Header().Get("X-RateLimit-Limit"),
			//	rec.Header().Get("X-RateLimit-Reset"),
			//)

			//t.Logf("Code %d Remain:%s Limit:%s Reset:%s",
			//	rec.Code,
			//	rec.Header().Get("X-RateLimit-Remaining"),
			//	rec.Header().Get("X-RateLimit-Limit"),
			//	rec.Header().Get("X-RateLimit-Reset"),
			//)
			if rec.Code != http.StatusTeapot && rec.Code != http.StatusConflict {
				t.Fatalf("Unexpected http code: %d", rec.Code)
			}
		},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
			atomic.AddInt32(countAllowed, 1)
		}),
		true, // do test logger
		ratelimit.WithVaryBy(pathGetter{}, scope.Website.WithID(1)),
		ratelimit.WithErrorHandler(mw.ErrorWithPanic, scope.Website.WithID(1)),
		ratelimit.WithDeniedHandler(deniedH, scope.Website.WithID(1)),
	)

	if have, want := atomic.LoadInt32(countDenied), int32(5); have != want {
		t.Errorf("Denied Have: %v Want: %v", have, want)
	}
	if have, want := atomic.LoadInt32(countAllowed), int32(4); have != want {
		t.Errorf("Allowed Have: %v Want: %v", have, want)
	}
}

func testBackendConfiguration(
	t *testing.T, httpRequestURL string,
	pv cfgmock.PathValue,
	assertResponse func(*httptest.ResponseRecorder),
	nextH http.Handler,
	testLogger bool,
	opts ...ratelimit.Option,
) {
	var logBuf log.MutexBuffer
	const httpUsers = 3
	const httpLoops = 3

	baseOpts := []ratelimit.Option{
		ratelimit.WithRootConfig(cfgmock.NewService(pv)),
		ratelimit.WithDebugLog(&logBuf),
		ratelimit.WithOptionFactory(backend.PrepareOptionFactory()),
	}

	srv := ratelimit.MustNew(append(baseOpts, opts...)...)

	srv.ErrorHandler = func(err error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
			assert.True(t, errors.IsNotFound(err), "%+v", err)
		})
	}

	// buf := new(bytes.Buffer)
	// srv.DebugCache(buf)
	// println("DebugCache: ", buf.String())

	req := func() *http.Request {
		req, _ := http.NewRequest("GET", httpRequestURL, nil)
		req.RemoteAddr = "2a02:d180::"
		return req.WithContext(scope.WithContext(req.Context(), 1, 2)) // website=euro store=at
	}()

	hpu := cstesting.NewHTTPParallelUsers(httpUsers, httpLoops, 600, time.Millisecond)
	hpu.AssertResponse = assertResponse

	// Food for the race detector
	hpu.ServeHTTP(req, srv.WithRateLimit(nextH))

	if testLogger {
		// Min 20 calls IsValid
		// Exactly one call to optionInflight.Do
		if have, want := strings.Count(logBuf.String(), `Service.ConfigByScopedGetter.IsValid`), (httpUsers*httpLoops)-1; have < want {
			t.Errorf("Service.ConfigByScopedGetter.IsValid: Have: %d < Want: %d", have, want)
		}
		if have, want := strings.Count(logBuf.String(), `Service.ConfigByScopedGetter.Inflight.Do`), 1; have != want {
			t.Errorf("Service.ConfigByScopedGetter.Inflight.Do: Have: %d Want: %d", have, want)
		}
		// println("\n", logBuf.String(), "\n")
	}
}

func TestBackend_Path_Errors(t *testing.T) {
	tests := []struct {
		cfgPath string
		val     interface{}
		errBhf  errors.BehaviourFunc
	}{
		{backend.Disabled.MustFQWebsite(2), struct{}{}, errors.IsNotValid},
		{backend.GCRAName.MustFQWebsite(2), struct{}{}, errors.IsNotValid},
	}
	for i, test := range tests {

		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			test.cfgPath: test.val,
		})
		cfgScp := cfgSrv.NewScoped(2, 0)

		optFnc := backend.PrepareOptionFactory()
		_, err := ratelimit.New(optFnc(cfgScp)...)
		assert.True(t, test.errBhf(err), "Index %d Error: %+v", i, err)
	}
}
