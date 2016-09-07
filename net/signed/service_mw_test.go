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

package signed_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/signed"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var testData = []byte(`“The most important property of a program is whether it accomplishes the intention of its user.” ― C.A.R. Hoare`)

const dataSHA256 = `keyId="test",algorithm="rot13",signature="cc7b14f207d3896a74ba4e4e965d49e6098af2191058edb9e9247caf0db8cd7b"`

func TestService_WithResponseSignature_MissingContext(t *testing.T) {

	var serviceErrorHandlerCalled = new(int32)

	srv := signed.MustNew(
		signed.WithRootConfig(cfgmock.NewService()),
		signed.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
		signed.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusExpectationFailed)
				assert.Error(t, err, "%+v", err)
				assert.True(t, errors.IsNotFound(err), "%+v", err)
				atomic.AddInt32(serviceErrorHandlerCalled, 1)
			})
		}),
	)

	handler := srv.WithResponseSignature(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called this next handler")
	}))

	r := httptest.NewRequest("/", "https://corestore.io", nil)

	hpu := cstesting.NewHTTPParallelUsers(5, 5, 100, time.Millisecond)
	hpu.AssertResponse = func(w *httptest.ResponseRecorder) {
		assert.Exactly(t, http.StatusExpectationFailed, w.Code)
		assert.Empty(t, w.Body.String())
	}
	hpu.ServeHTTP(r, handler)

	if have, want := *serviceErrorHandlerCalled, int32(25); have != want {
		t.Errorf("WithServiceErrorHandler call failed: Have: %d Want: %d", have, want)
	}
}

func TestService_WithResponseSignature_Disabled(t *testing.T) {

	var nextHandlerCalled = new(int32)

	srv := signed.MustNew(
		signed.WithDisable(scope.Store, 2, true),
		signed.WithRootConfig(cfgmock.NewService()),
		signed.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
		signed.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
	)

	handler := srv.WithResponseSignature(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		atomic.AddInt32(nextHandlerCalled, 1)
	}))

	r := httptest.NewRequest("/", "https://corestore.io", nil)
	r = r.WithContext(scope.WithContext(r.Context(), 1, 2))

	hpu := cstesting.NewHTTPParallelUsers(5, 5, 100, time.Millisecond)
	hpu.AssertResponse = func(w *httptest.ResponseRecorder) {
		assert.Exactly(t, http.StatusTeapot, w.Code)
		assert.Empty(t, w.Body.String())
	}
	hpu.ServeHTTP(r, handler)

	if have, want := *nextHandlerCalled, int32(25); have != want {
		t.Errorf("NextHandler call failed: Have: %d Want: %d", have, want)
	}
}

func TestService_WithResponseSignature_Buffered(t *testing.T) {

	var nextHandlerCalled = new(int32)
	key := []byte(`My guinea p1g runs acro55 my keyb0ard`)

	srv := signed.MustNew(
		signed.WithContentHMAC_SHA256(scope.Website, 1, key),
		signed.WithRootConfig(cfgmock.NewService()),
		signed.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
		signed.WithErrorHandler(scope.Website, 1, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
		signed.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
	)

	handler := srv.WithResponseSignature(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write(testData)
		atomic.AddInt32(nextHandlerCalled, 1)
	}))

	r := httptest.NewRequest("/", "https://corestore.io", nil)
	r = r.WithContext(scope.WithContext(r.Context(), 1, 2))

	hpu := cstesting.NewHTTPParallelUsers(5, 5, 100, time.Millisecond)
	hpu.AssertResponse = func(w *httptest.ResponseRecorder) {
		assert.Exactly(t, `sha256 41d1c5095693f329b0be01535af4069e6ecae899ede244eaf39c6f4f616307a6`, w.Header().Get(signed.ContentHMAC))
		assert.Exactly(t, http.StatusTeapot, w.Code)
		assert.Exactly(t, string(testData), w.Body.String())
	}
	hpu.ServeHTTP(r, handler)

	if have, want := *nextHandlerCalled, int32(25); have != want {
		t.Errorf("NextHandler call failed: Have: %d Want: %d", have, want)
	}
}

func TestService_WithResponseSignature_Trailer(t *testing.T) {

	var nextHandlerCalled = new(int32)
	key := []byte(`My gu1n34 p1g run5 acro55 my k3yb0ard`)

	srv := signed.MustNew(
		signed.WithTrailer(scope.Store, 2, true),
		signed.WithContentHMAC_Blake2b256(scope.Store, 2, key),
		signed.WithRootConfig(cfgmock.NewService()),
		signed.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
		signed.WithErrorHandler(scope.Store, 2, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
		signed.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
	)

	handler := srv.WithResponseSignature(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write(testData)
		atomic.AddInt32(nextHandlerCalled, 1)
	}))

	r := httptest.NewRequest("/", "https://corestore.io", nil)
	r = r.WithContext(scope.WithContext(r.Context(), 1, 2))

	hpu := cstesting.NewHTTPParallelUsers(5, 5, 100, time.Millisecond)
	hpu.AssertResponse = func(w *httptest.ResponseRecorder) {
		// ResponseRecorder cannot write the HTTP Trailer ...
		assert.Exactly(t, `blk2b256 5fa2a2c12bb66c830b84bb8b13e7ff0af0c6aa39236e3cf256c1e0eab16b4b05`, w.Header().Get(signed.ContentHMAC))
		assert.Exactly(t, http.StatusTeapot, w.Code)
		assert.Exactly(t, string(testData), w.Body.String())
		assert.Exactly(t, signed.ContentHMAC, w.Header().Get("Trailer"))
		//t.Logf("%#v", w.HeaderMap)
	}
	hpu.ServeHTTP(r, handler)

	if have, want := *nextHandlerCalled, int32(25); have != want {
		t.Errorf("NextHandler call failed: Have: %d Want: %d", have, want)
	}
}
