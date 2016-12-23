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

package geoip_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/net/geoip/maxmindfile"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var _ io.Closer = (*geoip.Service)(nil)

func mustGetTestService(opts ...geoip.Option) (*geoip.Service, func()) {
	maxMindDB := filepath.Join("testdata", "GeoIP2-Country-Test.mmdb")
	s := geoip.MustNew(append(opts, geoip.WithRootConfig(cfgmock.NewService()), maxmindfile.WithCountryFinder(maxMindDB))...)
	return s, func() {
		if err := s.Close(); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
	}
}

var _ geoip.Finder = (*finderError)(nil)

type finderError struct{}

func (finderError) FindCountry(ipAddress net.IP) (*geoip.Country, error) {
	return nil, errors.NewNotImplementedf("Failed to read country from MMDB")
}
func (finderError) Close() error { return nil }

func TestService_WithCountryByIP_Error_GeoReader(t *testing.T) {
	s := geoip.MustNew()
	defer func() {
		if err := s.Close(); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
	}()

	var calledWithServiceErrorHandler bool
	if err := s.Options(
		geoip.WithCountryFinder(finderError{}),
		geoip.WithErrorHandler(func(_ error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		} /*Default Scope ;-) */),
		geoip.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusVariantAlsoNegotiates)
				assert.True(t, errors.IsNotImplemented(err), "%+v", err)
				calledWithServiceErrorHandler = true
			})
		}),
	); err != nil {
		t.Fatal(err)
	}

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	countryHandler := s.WithCountryByIP(finalHandler)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://corestore.io", nil)
	req.Header.Set("X-Forwarded-For", "2a02:d200::")
	countryHandler.ServeHTTP(rec, req)
	assert.Exactly(t, http.StatusVariantAlsoNegotiates, rec.Code)
	assert.True(t, calledWithServiceErrorHandler, "calledWithServiceErrorHandler")
}

func TestService_WithCountryByIP_OK(t *testing.T) {
	s, closeFn := mustGetTestService()
	defer closeFn()

	var calledWithCountryByIP bool
	countryHandler := s.WithCountryByIP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipc, ok := geoip.FromContextCountry(r.Context())
		if !ok {
			panic("Expecting a country in the context")
		}
		assert.NotNil(t, ipc)
		assert.Exactly(t, "2a02:d200::", ipc.IP.String())
		assert.Exactly(t, "FI", ipc.Country.IsoCode)
		calledWithCountryByIP = true
	}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://corestore.io", nil)
	req.Header.Set("X-Forwarded-For", "2a02:d200::")
	countryHandler.ServeHTTP(rec, req)
	assert.Exactly(t, http.StatusOK, rec.Code)
	assert.True(t, calledWithCountryByIP, "calledWithCountryByIP")
}

func TestService_WithIsCountryAllowedByIP_Missing_Scope(t *testing.T) {
	s, closeFn := mustGetTestService()
	defer closeFn()

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	countryHandler := s.WithIsCountryAllowedByIP(finalHandler)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://corestore.io", nil)
	countryHandler.ServeHTTP(rec, req)
	assert.Exactly(t, http.StatusServiceUnavailable, rec.Code)
}

func TestService_WithCountryByIP_IncorrectIP(t *testing.T) {
	s, closeFn := mustGetTestService(geoip.WithServiceErrorHandler(func(err error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusVariantAlsoNegotiates)
			assert.True(t, errors.IsNotFound(err), "%+v", err)
		})
	}),
	)
	defer closeFn()

	countryHandler := s.WithCountryByIP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://corestore.io", nil)
	req.Header.Set("X-Forwarded-For", "324.2334.432.534")
	req.RemoteAddr = "Cat Content"
	countryHandler.ServeHTTP(rec, req)
	assert.Exactly(t, http.StatusVariantAlsoNegotiates, rec.Code)
}

func TestService_WithIsCountryAllowedByIP_ErrorWithContextCountryByIP(t *testing.T) {
	s, closeFn := mustGetTestService()
	defer closeFn()

	countryHandler := s.WithIsCountryAllowedByIP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://corestore.io", nil)
	req.Header.Set("X-Forwarded-For", "2R02:d2'0.:")
	countryHandler.ServeHTTP(rec, req)
	assert.Exactly(t, http.StatusServiceUnavailable, rec.Code)
}

func TestService_WithIsCountryAllowedByIP_MultiScopes(t *testing.T) {
	var logBuf = new(log.MutexBuffer)
	s, closeFn := mustGetTestService(
		geoip.WithDebugLog(logBuf),
	)
	defer closeFn()

	var calledAltHndlr int32
	var altHndlr = func(err error) http.Handler {
		if !errors.IsUnauthorized(err) {
			panic(fmt.Sprintf("Expecting an IsUnauthorized error:\n%+v", err))
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			atomic.AddInt32(&calledAltHndlr, 1)
		})
	}

	var calledFinalTestHandler int32
	var finalTestHandler = func(i int, wantCountryISO string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted) // just write a random status
			ipc, ok := geoip.FromContextCountry(r.Context())
			if !ok {
				panic(fmt.Sprintf("Cannot country context in request: %d wantCountryISO: %s", i, wantCountryISO))
			}

			if assert.NotNil(t, ipc, "Index %d", i) { // avoid runtime panic due to nil
				assert.Exactly(t, wantCountryISO, ipc.Country.IsoCode, "Index %d", i)
			}
			atomic.AddInt32(&calledFinalTestHandler, 1)
		})
	}

	var calledErrorHandler int32
	if err := s.Options(
		geoip.WithAlternativeHandler(altHndlr, scope.Store.Pack(2)), // for test case 2
		geoip.WithAllowedCountryCodes([]string{"DE", "CH", "FI"}, scope.Store.Pack(2)),
		geoip.WithServiceErrorHandler(func(err error) http.Handler {
			assert.Error(t, err, "%+v", err)
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusGatewayTimeout)
				atomic.AddInt32(&calledErrorHandler, 1)
			})
		}),
		geoip.WithErrorHandler(func(err error) http.Handler {
			assert.Error(t, err, "%+v", err)
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusHTTPVersionNotSupported)
				atomic.AddInt32(&calledErrorHandler, 1)
			})
		}, scope.Store.Pack(2)),
	); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		req            func() *http.Request
		wantCountryISO string
		wantCode       int
	}{
		// scope configuration websiteID& storeID not found
		0: {func() *http.Request {
			return httptest.NewRequest("GET", "http://corestore.io", nil)
		},
			"", http.StatusGatewayTimeout},

		// IP detected as origin from Finland
		1: {func() *http.Request {
			req := httptest.NewRequest("GET", "http://corestore.io", nil)
			req.Header.Set("X-Forwarded-For", "2a02:d200::")
			return req.WithContext(scope.WithContext(req.Context(), 1, 1)) // euro website / german store
		},
			"FI", http.StatusAccepted},

		// IP detected as origin from AT and alternative handler for scope Store == 2 gets called but AT not allowed
		2: {func() *http.Request {
			req := httptest.NewRequest("GET", "http://corestore.io", nil)
			req.RemoteAddr = "2a02:da80::"
			return req.WithContext(scope.WithContext(req.Context(), 1, 2)) // euro website / austrian store
		},
			"AT", http.StatusBadGateway},

		// IP detection errors and an error gets attached to the context
		3: {func() *http.Request {
			req := httptest.NewRequest("GET", "http://corestore.io", nil)
			req.RemoteAddr = "Er00r"
			return req.WithContext(scope.WithContext(req.Context(), 1, 2)) // euro website / austrian store
		},
			"XX", http.StatusHTTPVersionNotSupported},

		// IP from Germany, scope config not available and hence fall back to default
		4: {func() *http.Request {
			req := httptest.NewRequest("GET", "http://corestore.io", nil)
			req.RemoteAddr = "2a02:e240::"
			return req.WithContext(scope.WithContext(req.Context(), 1, 1)) // euro website / german store
		},
			"DE", http.StatusAccepted},
	}
	for i, test := range tests {
		req := test.req() // within the loop we'll get a race condition
		//hpu := cstesting.NewHTTPParallelUsers(1, 1, 200, time.Millisecond)
		hpu := cstesting.NewHTTPParallelUsers(8, 15, 200, time.Millisecond)
		hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
			assert.Exactly(t, test.wantCode, rec.Code, "Index %d", i)
			//t.Log(i, rec.Code, http.StatusText(rec.Code))
		}
		hpu.ServeHTTP(req,
			s.WithIsCountryAllowedByIP(finalTestHandler(i, test.wantCountryISO)),
		)
	}
	assert.Exactly(t, int32(120), atomic.LoadInt32(&calledAltHndlr), "calledAltHndlr")
	assert.Exactly(t, int32(240), atomic.LoadInt32(&calledFinalTestHandler), "calledFinalTestHandler")
	assert.Exactly(t, int32(240), atomic.LoadInt32(&calledErrorHandler), "calledErrorHandler")
	// println("\n\n", logBuf.String(), "\n\n")
}
