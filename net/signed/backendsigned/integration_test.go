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

package backendsigned_test

import (
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/signed"
	"github.com/corestoreio/csfw/net/signed/backendsigned"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestBackend_Path_Errors(t *testing.T) {
	tests := []struct {
		toPath func(s scope.Scope, scopeID int64) string
		val    interface{}
		errBhf errors.BehaviourFunc
	}{
		0: {backend.Disabled.MustFQ, struct{}{}, errors.IsNotValid},
		1: {backend.InTrailer.MustFQ, struct{}{}, errors.IsNotValid},
		2: {backend.AllowedMethods.MustFQ, struct{}{}, errors.IsNotValid},
		3: {backend.Key.MustFQ, struct{}{}, errors.IsNotValid},
		4: {backend.Algorithm.MustFQ, struct{}{}, errors.IsNotValid},
		5: {backend.HTTPHeaderType.MustFQ, struct{}{}, errors.IsNotValid},
		6: {backend.KeyID.MustFQ, struct{}{}, errors.IsNotValid},
	}
	for i, test := range tests {

		scpFnc := backendsigned.PrepareOptions(backend)
		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			test.toPath(scope.Website, 2): test.val,
		})
		cfgScp := cfgSrv.NewScoped(2, 0)

		_, err := signed.New(scpFnc(cfgScp)...)
		assert.True(t, test.errBhf(err), "Index %d Error: %+v", i, err)
	}
}

//func TestBackend_GCRA_Not_Registered(t *testing.T) {
//	testBackendConfiguration(t, "panic",
//		cfgmock.PathValue{},
//		func(rec *httptest.ResponseRecorder) {
//			assert.Exactly(t, http.StatusTeapot, rec.Code)
//		},
//		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			panic("Should not get called")
//		}),
//		false, // do not test logger
//		signed.WithVaryBy(scope.Website, 1, pathGetter{}),
//		signed.WithRateLimiter(scope.Website, 1, stubLimiter{}),
//	)
//}
//
//func TestBackend_WithDisable(t *testing.T) {
//
//	testBackendConfiguration(t, "panic",
//		cfgmock.PathValue{
//			backend.Disabled.MustFQ(scope.Website, 1): 1,
//		},
//		func(rec *httptest.ResponseRecorder) {
//			assert.Exactly(t, http.StatusTeapot, rec.Code)
//		},
//		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.WriteHeader(http.StatusTeapot)
//		}),
//		true, // do test logger
//		signed.WithVaryBy(scope.Website, 1, pathGetter{}),
//		signed.WithRateLimiter(scope.Website, 1, stubLimiter{}),
//	)
//}
//
//func TestBackend_WithGCRAMemStore(t *testing.T) {
//	var countDenied = new(int32)
//	var countAllowed = new(int32)
//
//	backend.Register(memstore.NewOptionFactory(backend))
//	defer backend.Deregister(memstore.OptionName)
//
//	testBackendConfiguration(t,
//		"http://corestore.io",
//		cfgmock.PathValue{
//			backend.RateLimitGCRAName.MustFQ(scope.Website, 1):                 "memstore",
//			backend.Disabled.MustFQ(scope.Website, 1):                          0,
//			backend.RateLimitStorageGcraMaxMemoryKeys.MustFQ(scope.Website, 1): 50,
//			backend.RateLimitBurst.MustFQ(scope.Website, 1):                    3,
//			backend.RateLimitRequests.MustFQ(scope.Website, 1):                 1,
//			backend.RateLimitDuration.MustFQ(scope.Website, 1):                 "i",
//		},
//		func(rec *httptest.ResponseRecorder) {
//			//t.Logf("Code %d Remain:%s Limit:%s Reset:%s",
//			//	rec.Code,
//			//	rec.Header().Get("X-RateLimit-Remaining"),
//			//	rec.Header().Get("X-RateLimit-Limit"),
//			//	rec.Header().Get("X-RateLimit-Reset"),
//			//)
//			if rec.Code != http.StatusTeapot && rec.Code != http.StatusConflict {
//				t.Fatalf("Unexpected http code: %d", rec.Code)
//			}
//		},
//		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.WriteHeader(http.StatusTeapot)
//			atomic.AddInt32(countAllowed, 1)
//		}),
//		true, // do test logger
//		signed.WithVaryBy(scope.Website, 1, pathGetter{}),
//		signed.WithDeniedHandler(scope.Website, 1, http.HandlerFunc(func(w http.ResponseWriter, rec *http.Request) {
//			atomic.AddInt32(countDenied, 1)
//			http.Error(w, "custom limit exceeded", http.StatusConflict)
//		})),
//	)
//
//	if have, want := atomic.LoadInt32(countDenied), int32(5); have != want {
//		t.Errorf("Denied Have: %v Want: %v", have, want)
//	}
//	if have, want := atomic.LoadInt32(countAllowed), int32(4); have != want {
//		t.Errorf("Allowed Have: %v Want: %v", have, want)
//	}
//}
//
//func testBackendConfiguration(
//	t *testing.T, httpRequestURL string,
//	pv cfgmock.PathValue,
//	assertResponse func(*httptest.ResponseRecorder),
//	nextH http.Handler,
//	testLogger bool,
//	opts ...signed.Option,
//) {
//
//	var logBuf log.MutexBuffer
//	const httpUsers = 3
//	const httpLoops = 3
//
//	var baseOpts = []signed.Option{
//		signed.WithRootConfig(cfgmock.NewService(pv)),
//		signed.WithLogger(logw.NewLog(logw.WithWriter(&logBuf), logw.WithLevel(logw.LevelDebug))),
//		signed.WithOptionFactory(backendsigned.PrepareOptions(backend)),
//	}
//
//	srv := signed.MustNew(append(baseOpts, opts...)...)
//
//	srv.ErrorHandler = func(err error) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.WriteHeader(http.StatusTeapot)
//			assert.True(t, errors.IsNotFound(err), "%+v", err)
//		})
//	}
//
//	req := func() *http.Request {
//		req, _ := http.NewRequest("GET", httpRequestURL, nil)
//		req.RemoteAddr = "2a02:d180::"
//		return req.WithContext(scope.WithContext(req.Context(), 1, 2)) // website=euro store=at
//	}()
//
//	hpu := cstesting.NewHTTPParallelUsers(httpUsers, httpLoops, 600, time.Millisecond)
//	hpu.AssertResponse = assertResponse
//
//	// Food for the race detector
//	hpu.ServeHTTP(req, srv.WithRateLimit(nextH))
//
//	if testLogger {
//		// Min 20 calls IsValid
//		// Exactly one call to optionInflight.Do
//		if have, want := strings.Count(logBuf.String(), `Service.ConfigByScopedGetter.IsValid`), (httpUsers*httpLoops)-1; have < want {
//			t.Errorf("Service.ConfigByScopedGetter.IsValid: Have: %d < Want: %d", have, want)
//		}
//		if have, want := strings.Count(logBuf.String(), `Service.ConfigByScopedGetter.Inflight.Do`), 1; have != want {
//			t.Errorf("Service.ConfigByScopedGetter.Inflight.Do: Have: %d Want: %d", have, want)
//		}
//		// println("\n", logBuf.String(), "\n")
//	}
//}
