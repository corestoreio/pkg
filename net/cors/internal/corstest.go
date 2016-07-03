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

package internal

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/cors"
	"github.com/corestoreio/csfw/util/cstesting"
)

func testHandler(fa interface {
	Fatalf(format string, args ...interface{})
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := cors.FromContext(r.Context()); err != nil {
			fa.Fatalf("%+v", err)
		}
		_, _ = w.Write([]byte("bar"))
	}
}

func assertHeaders(t *testing.T, resHeaders http.Header, reqHeaders map[string]string) {
	for name, value := range reqHeaders {
		if actual := strings.Join(resHeaders[name], ", "); actual != value {
			var buf [1024]byte
			written := runtime.Stack(buf[:], false)
			t.Errorf("Invalid header %q, wanted %q, got %q\n%s\n\n", name, value, actual, string(buf[:written]))
		}
	}
}

func TestNoConfig(t *testing.T, s *cors.Service, req *http.Request) {

	hpu := cstesting.NewHTTPParallelUsers(10, 3, 300, time.Microsecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin",
			"Access-Control-Allow-Origin":      "",
			"Access-Control-Allow-Methods":     "",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestMatchAllOrigin(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")

	hpu := cstesting.NewHTTPParallelUsers(10, 4, 300, time.Microsecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin",
			"Access-Control-Allow-Origin":      "http://foobar.com",
			"Access-Control-Allow-Methods":     "",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestAllowedOrigin(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")

	hpu := cstesting.NewHTTPParallelUsers(10, 2, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin",
			"Access-Control-Allow-Origin":      "http://foobar.com",
			"Access-Control-Allow-Methods":     "",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestWildcardOrigin(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foo.bar.com")

	hpu := cstesting.NewHTTPParallelUsers(2, 10, 300, time.Microsecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin",
			"Access-Control-Allow-Origin":      "http://foo.bar.com",
			"Access-Control-Allow-Methods":     "",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestDisallowedOrigin(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://barbaz.com")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin",
			"Access-Control-Allow-Origin":      "",
			"Access-Control-Allow-Methods":     "",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestDisallowedWildcardOrigin(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foo.baz.com")

	hpu := cstesting.NewHTTPParallelUsers(10, 3, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin",
			"Access-Control-Allow-Origin":      "",
			"Access-Control-Allow-Methods":     "",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestAllowedOriginFunc(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Set("Origin", "http://foobar.com")
	res := httptest.NewRecorder()
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

func TestAllowedMethod(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "PUT")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 10, time.Microsecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			"Access-Control-Allow-Origin":      "http://foobar.com",
			"Access-Control-Allow-Methods":     "PUT",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestAllowedMethodPassthrough(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "PUT")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 100, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			"Access-Control-Allow-Origin":      "http://foobar.com",
			"Access-Control-Allow-Methods":     "PUT",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestDisallowedMethod(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "PATCH")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			"Access-Control-Allow-Origin":      "",
			"Access-Control-Allow-Methods":     "",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestAllowedHeader(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "X-Header-2, X-HEADER-1")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			"Access-Control-Allow-Origin":      "http://foobar.com",
			"Access-Control-Allow-Methods":     "GET",
			"Access-Control-Allow-Headers":     "X-Header-2, X-Header-1",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestAllowedWildcardHeader(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "X-Header-2, X-HEADER-1")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			"Access-Control-Allow-Origin":      "http://foobar.com",
			"Access-Control-Allow-Methods":     "GET",
			"Access-Control-Allow-Headers":     "X-Header-2, X-Header-1",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestDisallowedHeader(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "X-Header-3, X-Header-1")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			"Access-Control-Allow-Origin":      "",
			"Access-Control-Allow-Methods":     "",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestOriginHeader(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			"Access-Control-Allow-Origin":      "http://foobar.com",
			"Access-Control-Allow-Methods":     "GET",
			"Access-Control-Allow-Headers":     "Origin",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestExposedHeader(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin",
			"Access-Control-Allow-Origin":      "http://foobar.com",
			"Access-Control-Allow-Methods":     "",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "X-Header-1, X-Header-2",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestAllowedCredentials(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			"Access-Control-Allow-Origin":      "http://foobar.com",
			"Access-Control-Allow-Methods":     "GET",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Max-Age":           "",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}

func TestMaxAge(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 300, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assertHeaders(t, rec.Header(), map[string]string{
			"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			"Access-Control-Allow-Origin":      "http://foobar.com",
			"Access-Control-Allow-Methods":     "GET",
			"Access-Control-Allow-Headers":     "",
			"Access-Control-Allow-Credentials": "",
			"Access-Control-Max-Age":           "30",
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS()(testHandler(t)))
}
