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

package backendratelimit_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/net/ratelimit/backendratelimit"
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

func TestBackend_WithDisable(t *testing.T) {
	testBackendConfiguration(t, "panic",
		cfgmock.PathValue{
			backend.RateLimitDisabled.MustFQ(scope.Website, 1): 1,
		},
		func(rec *httptest.ResponseRecorder) {
			assert.Exactly(t, http.StatusTeapot, rec.Code)
		},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := ratelimit.FromContextRateLimit(r.Context()); err != nil {
				t.Errorf("%+v", err)
			}
			w.WriteHeader(http.StatusTeapot)
		}),
		ratelimit.WithVaryBy(scope.Website, 1, pathGetter{}),
		ratelimit.WithRateLimiter(scope.Website, 1, stubLimiter{}),
	)
}

func TestBackend_WithGCRAMemStore(t *testing.T) {
	var countDenied = new(int32)
	var countAllowed = new(int32)
	testBackendConfiguration(t,
		"http://corestore.io",
		cfgmock.PathValue{
			backend.RateLimitDisabled.MustFQ(scope.Website, 1):                 0,
			backend.RateLimitStorageGcraMaxMemoryKeys.MustFQ(scope.Website, 1): 50,
			backend.RateLimitBurst.MustFQ(scope.Website, 1):                    3,
			backend.RateLimitRequests.MustFQ(scope.Website, 1):                 1,
			backend.RateLimitDuration.MustFQ(scope.Website, 1):                 "i",
		},
		func(rec *httptest.ResponseRecorder) {
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
			if err := ratelimit.FromContextRateLimit(r.Context()); err != nil {
				t.Errorf("%+v", err)
			}
			w.WriteHeader(http.StatusTeapot)
			atomic.AddInt32(countAllowed, 1)
		}),
		ratelimit.WithVaryBy(scope.Website, 1, pathGetter{}),
		ratelimit.WithDeniedHandler(scope.Website, 1, http.HandlerFunc(func(w http.ResponseWriter, rec *http.Request) {
			atomic.AddInt32(countDenied, 1)
			http.Error(w, "custom limit exceeded", http.StatusConflict)
		})),
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
	opts ...ratelimit.Option,
) {

	var logBuf log.MutexBuffer
	const httpUsers = 3
	const httpLoops = 3

	var baseOpts = []ratelimit.Option{
		ratelimit.WithLogger(logw.NewLog(logw.WithWriter(&logBuf), logw.WithLevel(logw.LevelDebug))),
		ratelimit.WithOptionFactory(backendratelimit.PrepareOptions(backend)),
	}

	srv := ratelimit.MustNew(append(baseOpts, opts...)...)

	req := func() *http.Request {
		o, err := scope.SetByCode(scope.Website, "euro")
		if err != nil {
			t.Fatal(err)
		}
		storeSrv := storemock.NewEurozzyService(o, store.WithStorageConfig(
			cfgmock.NewService(cfgmock.WithPV(pv)),
		))
		req, _ := http.NewRequest("GET", httpRequestURL, nil)
		req.RemoteAddr = "2a02:d180::"
		atSt, err := storeSrv.Store(scope.MockID(2)) // Austria Store
		if err != nil {
			t.Fatalf("%+v", err)
		}
		return req.WithContext(store.WithContextRequestedStore(req.Context(), atSt))
	}()

	hpu := cstesting.NewHTTPParallelUsers(httpUsers, httpLoops, 600, time.Millisecond)
	hpu.AssertResponse = assertResponse

	// Food for the race detector
	hpu.ServeHTTP(req, srv.WithRateLimit()(nextH))

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

func TestBackend_Path_Errors(t *testing.T) {

	tests := []struct {
		toPath func(s scope.Scope, scopeID int64) string
		val    interface{}
		errBhf errors.BehaviourFunc
	}{
		{backend.RateLimitDisabled.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.RateLimitBurst.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.RateLimitRequests.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.RateLimitDuration.MustFQ, "[a-z+", errors.IsFatal},
		{backend.RateLimitStorageGcraMaxMemoryKeys.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.RateLimitStorageGCRARedis.MustFQ, struct{}{}, errors.IsNotValid},
	}
	for i, test := range tests {

		scpFnc := backendratelimit.PrepareOptions(backend)
		cfgSrv := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			test.toPath(scope.Website, 2): test.val,
		}))
		cfgScp := cfgSrv.NewScoped(2, 0)

		_, err := ratelimit.New(scpFnc(cfgScp)...)
		assert.True(t, test.errBhf(err), "Index %d Error: %+v", i, err)
	}
}
