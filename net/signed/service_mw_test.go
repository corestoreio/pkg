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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/net/signed"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var testData = []byte(`“The most important property of a program is whether it accomplishes the intention of its user.” ― C.A.R. Hoare`)

func TestService_WithResponseSignature_MissingContext(t *testing.T) {

	var serviceErrorHandlerCalled = new(int32)

	srv := signed.MustNew(
		signed.WithDebugLog(ioutil.Discard),
		signed.WithRootConfig(cfgmock.NewService()),
		signed.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
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
		signed.WithDebugLog(ioutil.Discard),
		signed.WithDisable(scope.Store, 2, true),
		signed.WithRootConfig(cfgmock.NewService()),
		signed.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
			})
		}),
		signed.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
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
		signed.WithTrailer(scope.Website, 1, false),
		signed.WithDebugLog(ioutil.Discard),
		signed.WithContentHMACSHA256(scope.Website, 1, key),
		signed.WithRootConfig(cfgmock.NewService()),
		signed.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
			})
		}),
		signed.WithErrorHandler(scope.Website, 1, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
			})
		}),
		signed.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
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
		assert.Empty(t, w.Header().Get("Trailer"), `Should contain a trailer => w.Header().Get("Trailer")`)
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
		signed.WithDebugLog(ioutil.Discard),
		signed.WithTrailer(scope.Store, 2, true),
		signed.WithContentHMACBlake2b256(scope.Store, 2, key),
		signed.WithRootConfig(cfgmock.NewService()),
		signed.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
			})
		}),
		signed.WithErrorHandler(scope.Store, 2, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
			})
		}),
		signed.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
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

func TestService_WithRequestSignatureValidation_Full_Roundtrip(t *testing.T) {

	key := []byte(`My guinea p1g run5 acro55 my keyb0ard`)

	srv := signed.MustNew(
		signed.WithDebugLog(ioutil.Discard),
		signed.WithContentHMACSHA256(scope.Website, 1, key),
		signed.WithRootConfig(cfgmock.NewService()),
		signed.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
			})
		}),
		signed.WithErrorHandler(scope.Website, 1, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
			})
		}),
		signed.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
			})
		}),
	)

	var finalHandlerCalled = new(int32)
	initReq := httptest.NewRequest("GET", "https://corestore.io/product/id/3456", nil)
	initReq = initReq.WithContext(scope.WithContext(initReq.Context(), 1, 2))
	initResp := httptest.NewRecorder()
	// mw.Chain to be used to validate the correct API signature of WithResponseSignature function
	mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusAccepted)
			w.Write(testData)
			atomic.AddInt32(finalHandlerCalled, 1)
		}),
		srv.WithResponseSignature,
	).ServeHTTP(initResp, initReq)

	if initResp.Code != http.StatusAccepted || initResp.Header().Get(signed.ContentHMAC) == "" {
		t.Fatalf("Status %d\nHeader %v\nBody: %s", initResp.Code, initResp.HeaderMap, initResp.Body)
	}
	if have, want := *finalHandlerCalled, int32(1); have != want {
		t.Errorf("WithResponseSignature NextHandler call failed: Have: %d Want: %d", have, want)
	}
	atomic.StoreInt32(finalHandlerCalled, 0) // reset internal counter

	// keep this at 1,1 because hpu.ServeHTTP reuses the request for all calls to ServeHTTP,
	// but instead hpu.ServeHTTP must create for each request a new one.
	hpu := cstesting.NewHTTPParallelUsers(5, 5, 100, time.Millisecond)
	hpu.AssertResponse = func(w *httptest.ResponseRecorder) {
		assert.Exactly(t, http.StatusPartialContent, w.Code)
		assert.Exactly(t, `OK!`, w.Body.String())
	}
	hpu.ServeHTTPNewRequest(func() *http.Request {
		// create a new request. the idea is that this request comes from an
		// untrusted 3rd party service. we send the previous received body back to
		// the microservice. variable w refers to the previous made request to fetch
		// the data.
		r2 := httptest.NewRequest("POST", "https://corestore.io/checkout/cart/add", bytes.NewReader(initResp.Body.Bytes())) // reader to avoid races
		r2 = r2.WithContext(scope.WithContext(r2.Context(), 1, 2))
		r2.Header.Set(signed.ContentHMAC, initResp.Header().Get(signed.ContentHMAC))
		return r2
	}, mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			assert.Exactly(t, `sha256 7dace9827fd7aa3c83eee3776a81d03653ba1e272c98809f0752d9ded4561419`, rq.Header.Get(signed.ContentHMAC))
			w.WriteHeader(http.StatusPartialContent)
			w.Write([]byte(`OK!`))

			// read body twice (1. time in the middleware and 2nd time here) to check for
			// the copied io.ReadCloser in the r.Body.
			body, err := ioutil.ReadAll(rq.Body)
			if err != nil {
				t.Fatalf("failed to read the request body: %s", err)
			}
			assert.Exactly(t, string(testData), string(body))
			atomic.AddInt32(finalHandlerCalled, 1)
		}),
		srv.WithRequestSignatureValidation,
	))

	if have, want := *finalHandlerCalled, int32(25); have != want {
		t.Errorf("NextHandler call failed: Have: %d Want: %d", have, want)
	}
}

func TestService_WithRequestSignatureValidation(t *testing.T) {

	const hmacHeaderValue = `sha256 7dace9827fd7aa3c83eee3776a81d03653ba1e272c98809f0752d9ded4561419`
	key := []byte(`My guinea p1g run5 acro55 my keyb0ard`)
	var finalHandlerCalled = new(int32)

	tester := func(req *http.Request, opts ...signed.Option) func(*testing.T) {
		defer atomic.StoreInt32(finalHandlerCalled, 0)

		srv := signed.MustNew(
			signed.WithDebugLog(ioutil.Discard),
			signed.WithContentHMACSHA256(scope.Website, 1, key),
			signed.WithRootConfig(cfgmock.NewService()),
			signed.WithErrorHandler(scope.Default, 0, func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic(fmt.Sprintf("Should not get called\n%+v", err))
				})
			}),
			signed.WithErrorHandler(scope.Website, 1, func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic(fmt.Sprintf("Should not get called\n%+v", err))
				})
			}),
			signed.WithServiceErrorHandler(func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic(fmt.Sprintf("Should not get called\n%+v", err))
				})
			}),
		)

		if err := srv.Options(opts...); err != nil {
			t.Fatalf("%+v", err)
		}

		return func(t *testing.T) {

			rec := httptest.NewRecorder()
			mw.Chain(
				http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
					assert.Exactly(t, hmacHeaderValue, rq.Header.Get(signed.ContentHMAC))
					w.WriteHeader(http.StatusPartialContent)
					w.Write([]byte(`OK!`))
					atomic.AddInt32(finalHandlerCalled, 1)
				}),
				srv.WithRequestSignatureValidation,
			).ServeHTTP(rec, req)

			if rec.Code == http.StatusLoopDetected {
				assert.Exactly(t, http.StatusLoopDetected, rec.Code)
			} else {
				assert.Exactly(t, http.StatusPartialContent, rec.Code)
				assert.Exactly(t, `OK!`, rec.Body.String())
			}
			if have, want := *finalHandlerCalled, int32(1); have != want {
				t.Errorf("NextHandler call failed: Have: %d Want: %d", have, want)
			}
		}
	}

	t.Run("Valid", tester(
		func() *http.Request {
			r := httptest.NewRequest("POST", "https://corestore.io/checkout/cart/add", bytes.NewReader(testData))
			r = r.WithContext(scope.WithContext(r.Context(), 1, 2))
			r.Header.Set(signed.ContentHMAC, hmacHeaderValue)
			return r
		}(),
	))

	t.Run("Signature Not Found", tester(
		func() *http.Request {
			r := httptest.NewRequest("POST", "https://corestore.io/checkout/cart/add", nil)
			r = r.WithContext(scope.WithContext(r.Context(), 1, 2))
			return r
		}(),
		signed.WithErrorHandler(scope.Website, 1, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusLoopDetected)
				assert.Error(t, err)
				assert.True(t, errors.IsNotFound(err), "%+v", err)
				atomic.AddInt32(finalHandlerCalled, 1)
			})
		}),
	))

	t.Run("Unkown Method POST", tester(
		func() *http.Request {
			r := httptest.NewRequest("POST", "https://corestore.io/checkout/cart/add", nil)
			r = r.WithContext(scope.WithContext(r.Context(), 1, 2))
			r.Header.Set(signed.ContentHMAC, hmacHeaderValue)
			return r
		}(),
		signed.WithAllowedMethods(scope.Website, 1, "PUT", "PATCH"),
		signed.WithErrorHandler(scope.Website, 1, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusLoopDetected)
				assert.Error(t, err)
				assert.True(t, errors.IsNotValid(err), "%+v", err)
				assert.Contains(t, err.Error(), `"POST" not allowed in list: ["PUT" "PATCH"]`)
				atomic.AddInt32(finalHandlerCalled, 1)
			})
		}),
	))

	t.Run("Error Reading Body", tester(
		func() *http.Request {
			r := httptest.NewRequest("PUT", "https://corestore.io/checkout/cart/add", readerError{})
			r = r.WithContext(scope.WithContext(r.Context(), 1, 2))
			r.Header.Set(signed.ContentHMAC, hmacHeaderValue)
			return r
		}(),
		signed.WithErrorHandler(scope.Website, 1, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusLoopDetected)
				assert.Error(t, err)
				assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
				assert.Contains(t, err.Error(), `Reader already closed`)
				atomic.AddInt32(finalHandlerCalled, 1)
			})
		}),
	))

	t.Run("Disabled", tester(
		func() *http.Request {
			r := httptest.NewRequest("PUT", "https://corestore.io/checkout/cart/add", bytes.NewReader(testData))
			r = r.WithContext(scope.WithContext(r.Context(), 1, 2))
			r.Header.Set(signed.ContentHMAC, hmacHeaderValue)
			return r
		}(),
		signed.WithDisable(scope.Website, 1, true),
		signed.WithErrorHandler(scope.Website, 1, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
			})
		}),
	))

	t.Run("Config not valid", tester(
		func() *http.Request {
			r := httptest.NewRequest("PUT", "https://corestore.io/checkout/cart/add", bytes.NewReader(testData))
			r = r.WithContext(scope.WithContext(r.Context(), 1, 2))
			r.Header.Set(signed.ContentHMAC, hmacHeaderValue)
			return r
		}(),
		signed.WithDisable(scope.Website, 1, false),
		signed.WithAllowedMethods(scope.Website, 1), // empty list of allowed methods triggers an error
		signed.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusLoopDetected)
				assert.Error(t, err)
				assert.True(t, errors.IsNotValid(err), "%+v", err)
				assert.Contains(t, err.Error(), `ScopedConfig Scope(Website) ID(1) is invalid. IsNil(HeaderParseWriter=false) AllowedMethods: []`)
				atomic.AddInt32(finalHandlerCalled, 1)
			})
		}),
		signed.WithErrorHandler(scope.Website, 1, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called\n%+v", err))
			})
		}),
	))
}

type readerError struct{}

func (readerError) Read(p []byte) (int, error) {
	return 0, errors.NewAlreadyClosedf("Reader already closed")
}
