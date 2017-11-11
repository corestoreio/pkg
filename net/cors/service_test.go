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
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/corestoreio/cspkg/config/cfgmock"
	"github.com/corestoreio/cspkg/net/cors"
	corstest "github.com/corestoreio/cspkg/net/cors/internal"
	"github.com/corestoreio/cspkg/net/mw"
	"github.com/corestoreio/cspkg/store/scope"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func reqWithStore(method string) *http.Request {
	req, err := http.NewRequest(method, "http://corestore.io/foo", nil)
	if err != nil {
		panic(err)
	}
	// @see storemock.MustNewStoreAU
	// 2 = website OZ
	// 5 = australia
	return req.WithContext(scope.WithContext(req.Context(), 2, 5))
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
	_ = cors.MustNew(cors.WithSettings(cors.Settings{}, scope.Website.Pack(2)))
}

func TestNoConfig(t *testing.T) {
	s := cors.MustNew(cors.WithRootConfig(cfgmock.NewService()))
	req := reqWithStore("GET")
	corstest.TestNoConfig(t, s, req)
}

func TestService_Options_Scope_Website(t *testing.T) {

	var newSrv = func(opts ...cors.Option) *cors.Service {
		s := cors.MustNew(
			cors.WithRootConfig(cfgmock.NewService()),
			cors.WithServiceErrorHandler(mw.ErrorWithPanic),
		)
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
			newSrv(cors.WithSettings(cors.Settings{AllowedOrigins: []string{"*"}}, scope.Website.Pack(2))),
			reqWithStore("GET"),
			corstest.TestMatchAllOrigin,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{AllowedOrigins: []string{"http://foobar.com"}}, scope.Website.Pack(2))),
			reqWithStore("GET"),
			corstest.TestAllowedOrigin,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{AllowedOrigins: []string{"http://*.bar.com"}}, scope.Website.Pack(2))),
			reqWithStore("GET"),
			corstest.TestWildcardOrigin,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{AllowedOrigins: []string{"http://foobar.com"}}, scope.Website.Pack(2))),
			reqWithStore("GET"),
			corstest.TestDisallowedOrigin,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{AllowedOrigins: []string{"http://*.bar.com"}}, scope.Website.Pack(2))),
			reqWithStore("GET"),
			corstest.TestDisallowedWildcardOrigin,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{
				AllowOriginFunc: func(o string) bool {
					r, _ := regexp.Compile("^http://foo") // don't do this on production systems! pre-compile before use!
					return r.MatchString(o)
				},
			}, scope.Website.Pack(2))),
			reqWithStore("GET"),
			corstest.TestAllowedOriginFunc,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{
				AllowedOrigins: []string{"http://foobar.com"},
				AllowedMethods: []string{"PUT", "DELETE"},
			}, scope.Website.Pack(2))),
			reqWithStore("OPTIONS"),
			corstest.TestAllowedMethodNoPassthrough,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{
				AllowedOrigins:     []string{"http://foobar.com"},
				AllowedMethods:     []string{"PUT", "DELETE"},
				OptionsPassthrough: true,
			}, scope.Website.Pack(2))),
			reqWithStore("OPTIONS"),
			corstest.TestAllowedMethodPassthrough,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{
				AllowedOrigins: []string{"http://foobar.com"},
				AllowedHeaders: []string{"X-Header-1", "x-header-2"},
			}, scope.Website.Pack(2))),
			reqWithStore("OPTIONS"),
			corstest.TestAllowedHeader,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{
				AllowedOrigins: []string{"http://foobar.com"},
				ExposedHeaders: []string{"X-Header-1", "x-header-2"},
			}, scope.Website.Pack(2))),
			reqWithStore("GET"),
			corstest.TestExposedHeader,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{
				AllowedOrigins:   []string{"http://foobar.com"},
				AllowCredentials: true,
			}, scope.Website.Pack(2))),
			reqWithStore("OPTIONS"),
			corstest.TestAllowedCredentials,
		},
		{
			newSrv(cors.WithSettings(cors.Settings{
				AllowedOrigins: []string{"http://foobar.com"},
				MaxAge:         "30",
			}, scope.Website.Pack(2))),
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

func getBaseCorsService(opts ...cors.Option) *cors.Service {
	s := cors.MustNew(
		cors.WithRootConfig(cfgmock.NewService()),
	)
	if err := s.Options(opts...); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	return s
}

func TestMatchAllOrigin(t *testing.T) {
	s := getBaseCorsService(cors.WithSettings(cors.Settings{AllowedOrigins: []string{"*"}}))
	req := reqWithStore("GET")
	corstest.TestMatchAllOrigin(t, s, req)
}

func TestAllowedOrigin(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{AllowedOrigins: []string{"http://foobar.com"}}),
	)
	req := reqWithStore("GET")
	corstest.TestAllowedOrigin(t, s, req)
}

func TestWildcardOrigin(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{AllowedOrigins: []string{"http://*.bar.com"}}),
	)
	req := reqWithStore("GET")
	corstest.TestWildcardOrigin(t, s, req)
}

func TestDisallowedOrigin(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{AllowedOrigins: []string{"http://foobar.com"}}),
	)
	req := reqWithStore("GET")
	corstest.TestDisallowedOrigin(t, s, req)
}

func TestDisallowedWildcardOrigin(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{AllowedOrigins: []string{"http://*.bar.com"}}),
	)
	req := reqWithStore("GET")
	corstest.TestDisallowedWildcardOrigin(t, s, req)
}

func TestAllowedOriginFunc(t *testing.T) {
	r, _ := regexp.Compile("^http://foo")
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{AllowOriginFunc: func(o string) bool {
			return r.MatchString(o)
		}}),
	)
	req := reqWithStore("GET")
	corstest.TestAllowedOriginFunc(t, s, req)
}

func TestAllowedMethod(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins: []string{"http://foobar.com"},
			AllowedMethods: []string{"PUT", "DELETE"},
		}),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedMethodNoPassthrough(t, s, req)
}

func TestAllowedMethodPassthrough(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins:     []string{"http://foobar.com"},
			AllowedMethods:     []string{"PUT", "DELETE"},
			OptionsPassthrough: true,
		}),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedMethodPassthrough(t, s, req)
}

func TestDisallowedMethod(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins: []string{"http://foobar.com"},
			AllowedMethods: []string{"PUT", "DELETE"},
		}),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestDisallowedMethod(t, s, req)
}

func TestAllowedHeader(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins: []string{"http://foobar.com"},
			AllowedHeaders: []string{"X-Header-1", "x-header-2"},
		}),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedHeader(t, s, req)
}

func TestAllowedWildcardHeader(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins: []string{"http://foobar.com"},
			AllowedHeaders: []string{"*"},
		}),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedWildcardHeader(t, s, req)
}

func TestDisallowedHeader(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins: []string{"http://foobar.com"},
			AllowedHeaders: []string{"X-Header-1", "x-header-2"},
		}),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestDisallowedHeader(t, s, req)
}

func TestOriginHeader(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins: []string{"http://foobar.com"},
		}),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestOriginHeader(t, s, req)
}

func TestExposedHeader(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins: []string{"http://foobar.com"},
			ExposedHeaders: []string{"X-Header-1", "x-header-2"},
		}),
	)
	req := reqWithStore("GET")
	corstest.TestExposedHeader(t, s, req)
}

func TestExposedHeader_MultiScope(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins: []string{"http://foobar.com"},
			ExposedHeaders: []string{"X-Header-1", "x-header-2"},
		}),
		cors.WithSettings(cors.Settings{
			AllowCredentials: false,
		}, scope.Website.Pack(1)),
	)

	reqDefault := httptest.NewRequest("GET", "http://corestore.io/reqDefault", nil)
	reqDefault = reqDefault.WithContext(
		scope.WithContext(reqDefault.Context(), 144, 244), // falls back to default scope
	)
	corstest.TestExposedHeader(t, s, reqDefault)

	// 2nd test
	reqWebsite, _ := http.NewRequest("OPTIONS", "http://corestore.io/reqWebsite", nil)
	reqWebsite = reqWebsite.WithContext(
		scope.WithContext(reqDefault.Context(), 1, 2), // website EURO and store AT
	)
	if err := s.Options(cors.WithSettings(cors.Settings{
		AllowCredentials: true,
	}, scope.Website.Pack(1))); err != nil {
		t.Errorf("%+v", err)
	}
	corstest.TestAllowedCredentials(t, s, reqWebsite)
}

func TestAllowedCredentials(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins:   []string{"http://foobar.com"},
			AllowCredentials: true,
		}),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedCredentials(t, s, req)
}

func TestMaxAge(t *testing.T) {
	s := getBaseCorsService(
		cors.WithSettings(cors.Settings{
			AllowedOrigins: []string{"http://foobar.com"},
			MaxAge:         "30", // seconds
		}),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestMaxAge(t, s, req)
}

func TestWithCORS_Error_Missing_ScopeContext(t *testing.T) {
	var serviceErrorHandlerCalled bool
	s := getBaseCorsService(
		cors.WithServiceErrorHandler(func(err error) http.Handler {
			assert.True(t, errors.IsNotFound(err), "%+v", err)
			serviceErrorHandlerCalled = true
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			})
		}),
	)

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should be called")
	})

	countryHandler := s.WithCORS(finalHandler)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://corestore.io", nil)
	countryHandler.ServeHTTP(rec, req)
	assert.Exactly(t, http.StatusTeapot, rec.Code)
	assert.Empty(t, rec.Body.String())
	assert.True(t, serviceErrorHandlerCalled, "serviceErrorHandlerCalled")
}

func TestWithCORS_Error_InvalidConfig(t *testing.T) {
	s := getBaseCorsService(cors.WithSettings(cors.Settings{
		AllowedMethods: []string{"a", "b", "b", "b", "b", "b", "b", "b", "b", "b", "b", "b", "b"},
	}))

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	countryHandler := s.WithCORS(finalHandler)
	rec := httptest.NewRecorder()
	req := reqWithStore("GET")
	countryHandler.ServeHTTP(rec, req)
	assert.Exactly(t, http.StatusServiceUnavailable, rec.Code)
}
