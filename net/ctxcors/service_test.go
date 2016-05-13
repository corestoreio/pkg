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

package ctxcors_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/ctxcors"
	"github.com/corestoreio/csfw/net/ctxcors/internal/corstest"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func testHandler(fa interface {
	Fatal(args ...interface{})
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ctxcors.FromContext(r.Context()); err != nil {
			fa.Fatal(errors.PrintLoc(err))
		}
		_, _ = w.Write([]byte("bar"))
	}
}

func assertHeaders(t *testing.T, resHeaders http.Header, reqHeaders map[string]string) {
	for name, value := range reqHeaders {
		if actual := strings.Join(resHeaders[name], ", "); actual != value {
			t.Errorf("Invalid header %q, wanted %q, got %q", name, value, actual)
		}
	}
}

func reqWithStore(method string) *http.Request {
	req, err := http.NewRequest(method, "http://corestore.io/foo", nil)
	if err != nil {
		panic(err)
	}

	return req.WithContext(
		store.WithContextRequestedStore(req.Context(), storemock.MustNewStoreAU(cfgmock.NewService())),
	)
}

func TestMustNew(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotValid(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	_ = ctxcors.MustNew(ctxcors.WithMaxAge(scope.Default, 0, -2*time.Second))
}

func TestNoConfig(t *testing.T) {
	s := ctxcors.MustNew()

	res := httptest.NewRecorder()

	s.WithCORS()(testHandler(t)).ServeHTTP(res, reqWithStore("GET"))

	assertHeaders(t, res.Header(), map[string]string{
		"Vary": "Origin",
		"Access-Control-Allow-Origin":      "",
		"Access-Control-Allow-Methods":     "",
		"Access-Control-Allow-Headers":     "",
		"Access-Control-Allow-Credentials": "",
		"Access-Control-Max-Age":           "",
		"Access-Control-Expose-Headers":    "",
	})
}

func TestMatchAllOrigin(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "*"),
	)

	res := httptest.NewRecorder()
	req := reqWithStore("GET")
	req.Header.Add("Origin", "http://foobar.com")

	s.WithCORS()(testHandler(t)).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"Vary": "Origin",
		"Access-Control-Allow-Origin":      "http://foobar.com",
		"Access-Control-Allow-Methods":     "",
		"Access-Control-Allow-Headers":     "",
		"Access-Control-Allow-Credentials": "",
		"Access-Control-Max-Age":           "",
		"Access-Control-Expose-Headers":    "",
	})
}

func TestAllowedOrigin(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
	)

	res := httptest.NewRecorder()
	req := reqWithStore("GET")
	req.Header.Add("Origin", "http://foobar.com")

	s.WithCORS()(testHandler(t)).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"Vary": "Origin",
		"Access-Control-Allow-Origin":      "http://foobar.com",
		"Access-Control-Allow-Methods":     "",
		"Access-Control-Allow-Headers":     "",
		"Access-Control-Allow-Credentials": "",
		"Access-Control-Max-Age":           "",
		"Access-Control-Expose-Headers":    "",
	})
}

func TestWildcardOrigin(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://*.bar.com"),
	)

	res := httptest.NewRecorder()
	req := reqWithStore("GET")
	req.Header.Add("Origin", "http://foo.bar.com")

	s.WithCORS()(testHandler(t)).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"Vary": "Origin",
		"Access-Control-Allow-Origin":      "http://foo.bar.com",
		"Access-Control-Allow-Methods":     "",
		"Access-Control-Allow-Headers":     "",
		"Access-Control-Allow-Credentials": "",
		"Access-Control-Max-Age":           "",
		"Access-Control-Expose-Headers":    "",
	})
}

func TestDisallowedOrigin(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
	)

	res := httptest.NewRecorder()
	req := reqWithStore("GET")
	req.Header.Add("Origin", "http://barbaz.com")

	s.WithCORS()(testHandler(t)).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"Vary": "Origin",
		"Access-Control-Allow-Origin":      "",
		"Access-Control-Allow-Methods":     "",
		"Access-Control-Allow-Headers":     "",
		"Access-Control-Allow-Credentials": "",
		"Access-Control-Max-Age":           "",
		"Access-Control-Expose-Headers":    "",
	})
}

func TestDisallowedWildcardOrigin(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://*.bar.com"),
	)

	res := httptest.NewRecorder()
	req := reqWithStore("GET")
	req.Header.Add("Origin", "http://foo.baz.com")

	s.WithCORS()(testHandler(t)).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"Vary": "Origin",
		"Access-Control-Allow-Origin":      "",
		"Access-Control-Allow-Methods":     "",
		"Access-Control-Allow-Headers":     "",
		"Access-Control-Allow-Credentials": "",
		"Access-Control-Max-Age":           "",
		"Access-Control-Expose-Headers":    "",
	})
}

func TestAllowedOriginFunc(t *testing.T) {
	r, _ := regexp.Compile("^http://foo")
	s := ctxcors.MustNew(
		ctxcors.WithAllowOriginFunc(scope.Default, 0, func(o string) bool {
			return r.MatchString(o)
		}),
	)

	req := reqWithStore("GET")

	res := httptest.NewRecorder()
	req.Header.Set("Origin", "http://foobar.com")
	s.WithCORS()(testHandler(t)).ServeHTTP(res, req)
	assertHeaders(t, res.Header(), map[string]string{
		"Access-Control-Allow-Origin": "http://foobar.com",
	})

	res = httptest.NewRecorder()
	req.Header.Set("Origin", "http://barfoo.com")
	s.WithCORS()(testHandler(t)).ServeHTTP(res, req)
	assertHeaders(t, res.Header(), map[string]string{
		"Access-Control-Allow-Origin": "",
	})
}

func TestAllowedMethod(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedMethods(scope.Default, 0, "PUT", "DELETE"),
	)

	res := httptest.NewRecorder()
	req := reqWithStore("OPTIONS")
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "PUT")

	s.WithCORS()(testHandler(t)).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		"Access-Control-Allow-Origin":      "http://foobar.com",
		"Access-Control-Allow-Methods":     "PUT",
		"Access-Control-Allow-Headers":     "",
		"Access-Control-Allow-Credentials": "",
		"Access-Control-Max-Age":           "",
		"Access-Control-Expose-Headers":    "",
	})
}

func TestAllowedMethodPassthrough(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedMethods(scope.Default, 0, "PUT", "DELETE"),
		ctxcors.WithOptionsPassthrough(scope.Default, 0, true),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedMethodPassthrough(t, s, req)
}

func TestDisallowedMethod(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedMethods(scope.Default, 0, "PUT", "DELETE"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestDisallowedMethod(t, s, req)
}

func TestAllowedHeader(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedHeaders(scope.Default, 0, "X-Header-1", "x-header-2"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedHeader(t, s, req)
}

func TestAllowedWildcardHeader(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedHeaders(scope.Default, 0, "*"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedWildcardHeader(t, s, req)
}

func TestDisallowedHeader(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedHeaders(scope.Default, 0, "X-Header-1", "x-header-2"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestDisallowedHeader(t, s, req)
}

func TestOriginHeader(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestOriginHeader(t, s, req)
}

func TestExposedHeader(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithExposedHeaders(scope.Default, 0, "X-Header-1", "x-header-2"),
	)

	req := reqWithStore("GET")
	corstest.TestExposedHeader(t, s, req)
}

func TestAllowedCredentials(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowCredentials(scope.Default, 0, true),
	)

	req := reqWithStore("OPTIONS")
	corstest.TestAllowedCredentials(t, s, req)
}

func TestMaxAge(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithMaxAge(scope.Default, 0, time.Second*30),
	)

	req := reqWithStore("OPTIONS")
	corstest.TestMaxAge(t, s, req)
}
