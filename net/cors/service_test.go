// Copyright (c) 2014 Olivier Poitrey <rs@dailymotion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cors_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/cors"
	corstest "github.com/corestoreio/csfw/net/cors/internal"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func reqWithStore(method string) *http.Request {
	req, err := http.NewRequest(method, "http://corestore.io/foo", nil)
	if err != nil {
		panic(err)
	}

	return req.WithContext(
		store.WithContextRequestedStore(req.Context(), storemock.MustNewStoreAU(cfgmock.NewService())),
	)
}

func withError() cors.Option {
	return func(s *cors.Service) error {
		return errors.NewNotValidf("Paaaaaaaniic!")
	}
}

func TestMustNew_Default(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotValid(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	_ = cors.MustNew(withError())
}

func TestMustNew_Store(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotSupported(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	_ = cors.MustNew(func(s *cors.Service) error {
		return errors.NewNotSupportedf("Not supported!")
	})
}

func TestMustNew_NoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			t.Fatalf("Expecting NOT a Panic with error: %s", err)
		}
	}()
	_ = cors.MustNew(cors.WithMaxAge(scope.Website, 2, -2*time.Second))
}

func TestNoConfig(t *testing.T) {
	s := cors.MustNew()
	req := reqWithStore("GET")
	corstest.TestNoConfig(t, s, req)
}

func TestService_Options_Scope_Website(t *testing.T) {

	var newSrv = func(opts ...cors.Option) *cors.Service {
		s := cors.MustNew(cors.WithLogger(log.BlackHole{})) // why can't i append opts... after WithLogger() ?
		if err := s.Options(opts...); err != nil {
			t.Fatal(err)
		}
		return s
	}

	tests := []struct {
		srv    *cors.Service
		req    *http.Request
		tester func(t *testing.T, s *cors.Service, req *http.Request)
	}{
		{
			newSrv(cors.WithAllowedOrigins(scope.Website, 2, "*")),
			reqWithStore("GET"),
			corstest.TestMatchAllOrigin,
		},
		{
			newSrv(cors.WithAllowedOrigins(scope.Website, 2, "http://foobar.com")),
			reqWithStore("GET"),
			corstest.TestAllowedOrigin,
		},
		{
			newSrv(cors.WithAllowedOrigins(scope.Website, 2, "http://*.bar.com")),
			reqWithStore("GET"),
			corstest.TestWildcardOrigin,
		},
		{
			newSrv(cors.WithAllowedOrigins(scope.Website, 2, "http://foobar.com")),
			reqWithStore("GET"),
			corstest.TestDisallowedOrigin,
		},
		{
			newSrv(cors.WithAllowedOrigins(scope.Website, 2, "http://*.bar.com")),
			reqWithStore("GET"),
			corstest.TestDisallowedWildcardOrigin,
		},
		{
			newSrv(cors.WithAllowOriginFunc(scope.Website, 2, func(o string) bool {
				r, _ := regexp.Compile("^http://foo") // don't do this on production systems!
				return r.MatchString(o)
			})),
			reqWithStore("GET"),
			corstest.TestAllowedOriginFunc,
		},
		{
			newSrv(cors.WithAllowedOrigins(scope.Website, 2, "http://foobar.com"), cors.WithAllowedMethods(scope.Website, 2, "PUT", "DELETE")),
			reqWithStore("OPTIONS"),
			corstest.TestAllowedMethod,
		},
		{
			newSrv(cors.WithAllowedMethods(scope.Website, 2, "PUT", "DELETE"), cors.WithAllowedOrigins(scope.Website, 2, "http://foobar.com"), cors.WithOptionsPassthrough(scope.Website, 2, true)),
			reqWithStore("OPTIONS"),
			corstest.TestAllowedMethodPassthrough,
		},
		{
			newSrv(cors.WithAllowedOrigins(scope.Website, 2, "http://foobar.com"), cors.WithAllowedHeaders(scope.Website, 2, "X-Header-1", "x-header-2")),
			reqWithStore("OPTIONS"),
			corstest.TestAllowedHeader,
		},
		{
			newSrv(cors.WithAllowedOrigins(scope.Website, 2, "http://foobar.com"), cors.WithExposedHeaders(scope.Website, 2, "X-Header-1", "x-header-2")),
			reqWithStore("GET"),
			corstest.TestExposedHeader,
		},
		{
			newSrv(cors.WithAllowedOrigins(scope.Website, 2, "http://foobar.com"), cors.WithAllowCredentials(scope.Website, 2, true)),
			reqWithStore("OPTIONS"),
			corstest.TestAllowedCredentials,
		},
		{
			newSrv(cors.WithMaxAge(scope.Website, 2, time.Second*30), cors.WithAllowedOrigins(scope.Website, 2, "http://foobar.com")),
			reqWithStore("OPTIONS"),
			corstest.TestMaxAge,
		},
	}
	for _, test := range tests {
		// for debugging comment this out to see the index which fails
		// t.Logf("Running Index %d Tester %q", i, runtime.FuncForPC(reflect.ValueOf(test.tester).Pointer()).Name())
		test.tester(t, test.srv, test.req)
	}
}

func TestMatchAllOrigin(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "*"),
	)
	req := reqWithStore("GET")
	corstest.TestMatchAllOrigin(t, s, req)
}

func TestAllowedOrigin(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
	)
	req := reqWithStore("GET")
	corstest.TestAllowedOrigin(t, s, req)
}

func TestWildcardOrigin(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://*.bar.com"),
	)
	req := reqWithStore("GET")
	corstest.TestWildcardOrigin(t, s, req)
}

func TestDisallowedOrigin(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
	)
	req := reqWithStore("GET")
	corstest.TestDisallowedOrigin(t, s, req)
}

func TestDisallowedWildcardOrigin(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://*.bar.com"),
	)
	req := reqWithStore("GET")
	corstest.TestDisallowedWildcardOrigin(t, s, req)
}

func TestAllowedOriginFunc(t *testing.T) {
	r, _ := regexp.Compile("^http://foo")
	s := cors.MustNew(
		cors.WithAllowOriginFunc(scope.Default, 0, func(o string) bool {
			return r.MatchString(o)
		}),
	)
	req := reqWithStore("GET")
	corstest.TestAllowedOriginFunc(t, s, req)
}

func TestAllowedMethod(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		cors.WithAllowedMethods(scope.Default, 0, "PUT", "DELETE"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedMethod(t, s, req)
}

func TestAllowedMethodPassthrough(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		cors.WithAllowedMethods(scope.Default, 0, "PUT", "DELETE"),
		cors.WithOptionsPassthrough(scope.Default, 0, true),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedMethodPassthrough(t, s, req)
}

func TestDisallowedMethod(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		cors.WithAllowedMethods(scope.Default, 0, "PUT", "DELETE"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestDisallowedMethod(t, s, req)
}

func TestAllowedHeader(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		cors.WithAllowedHeaders(scope.Default, 0, "X-Header-1", "x-header-2"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedHeader(t, s, req)
}

func TestAllowedWildcardHeader(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		cors.WithAllowedHeaders(scope.Default, 0, "*"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedWildcardHeader(t, s, req)
}

func TestDisallowedHeader(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		cors.WithAllowedHeaders(scope.Default, 0, "X-Header-1", "x-header-2"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestDisallowedHeader(t, s, req)
}

func TestOriginHeader(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestOriginHeader(t, s, req)
}

func TestExposedHeader(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		cors.WithExposedHeaders(scope.Default, 0, "X-Header-1", "x-header-2"),
	)

	req := reqWithStore("GET")
	corstest.TestExposedHeader(t, s, req)
}

func TestExposedHeader_MultiScope(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		cors.WithExposedHeaders(scope.Default, 0, "X-Header-1", "x-header-2"),
		cors.WithAllowCredentials(scope.Website, 1, false),
	)

	reqDefault, _ := http.NewRequest("GET", "http://corestore.io/reqDefault", nil)
	reqDefault = reqDefault.WithContext(
		store.WithContextRequestedStore(reqDefault.Context(), storemock.MustNewStoreAU(cfgmock.NewService())),
	)
	corstest.TestExposedHeader(t, s, reqDefault)

	eur := storemock.NewEurozzyService(scope.Option{Website: scope.MockID(1)}, store.WithStorageConfig(cfgmock.NewService()))
	atStore, atErr := eur.Store(scope.MockID(2)) // ID = 2 store Austria
	reqWebsite, _ := http.NewRequest("OPTIONS", "http://corestore.io/reqWebsite", nil)
	reqWebsite = reqWebsite.WithContext(
		store.WithContextRequestedStore(reqWebsite.Context(), atStore, atErr),
	)
	if err := s.Options(cors.WithAllowCredentials(scope.Website, 1, true)); err != nil {
		t.Errorf("%+v", err)
	}
	corstest.TestAllowedCredentials(t, s, reqWebsite)
}

func TestAllowedCredentials(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		cors.WithAllowCredentials(scope.Default, 0, true),
	)

	req := reqWithStore("OPTIONS")
	corstest.TestAllowedCredentials(t, s, req)
}

func TestMaxAge(t *testing.T) {
	s := cors.MustNew(
		cors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		cors.WithMaxAge(scope.Default, 0, time.Second*30),
	)

	req := reqWithStore("OPTIONS")
	corstest.TestMaxAge(t, s, req)
}

func TestWithCORS_Error_StoreManager(t *testing.T) {
	s := cors.MustNew()

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := cors.FromContext(r.Context())
		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	})

	countryHandler := s.WithCORS()(finalHandler)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	countryHandler.ServeHTTP(rec, req)
}

func TestWithCORS_Error_InvalidConfig(t *testing.T) {
	s := cors.MustNew(cors.WithAllowedMethods(scope.Default, 0))

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := cors.FromContext(r.Context())
		assert.True(t, errors.IsNotValid(err), "Error: %s", err)
	})

	countryHandler := s.WithCORS()(finalHandler)
	rec := httptest.NewRecorder()
	req := reqWithStore("GET")
	countryHandler.ServeHTTP(rec, req)
}
