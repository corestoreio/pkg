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
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/cspkg/net/cors"
	"github.com/corestoreio/cspkg/util/cstesting"
)

const testHandlerBodyData = `foobar`

func testHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(testHandlerBodyData))
	}
}

func assertHeaders(t *testing.T, resHeaders http.Header, reqHeaders map[string]string) {
	for name, value := range reqHeaders {
		if actual := strings.Join(resHeaders[name], ", "); actual != value {
			//var buf [1024]byte
			//written := runtime.Stack(buf[:], false)
			t.Errorf("Invalid header %q, wanted %q, got %q", name, value, actual) // string(buf[:written])
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
}

func TestMatchAllOrigin(t *testing.T, s *cors.Service, req *http.Request) {

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
			"Access-Control-Expose-Headers":    "",
		})
	}
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
}

func TestAllowedOrigin(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")

	hpu := cstesting.NewHTTPParallelUsers(10, 4, 300, time.Millisecond)
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
}

func TestDisallowedOrigin(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://barbaz.com")

	hpu := cstesting.NewHTTPParallelUsers(10, 5, 300, time.Millisecond)
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
}

func TestAllowedOriginFunc(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Set("Origin", "http://foobar.com")
	res := httptest.NewRecorder()
	s.WithCORS(testHandler()).ServeHTTP(res, req)
	assertHeaders(t, res.Header(), map[string]string{
		"Access-Control-Allow-Origin": "http://foobar.com",
	})

	res = httptest.NewRecorder()
	req.Header.Set("Origin", "http://barfoo.com")
	s.WithCORS(testHandler()).ServeHTTP(res, req)
	assertHeaders(t, res.Header(), map[string]string{
		"Access-Control-Allow-Origin": "",
	})
}

func TestAllowedMethodNoPassthrough(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "PUT")

	// todo: regression test because once time.MicroSecond will be used this test fails
	// but after a couple of days refactoring i can't find it.
	hpu := cstesting.NewHTTPParallelUsers(10, 10, 200, time.Millisecond)
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

		if have, want := rec.Body.String(), ""; have != want {
			t.Errorf("Have: %v Want: %v", have, want)
		}
	}
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
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
		if have, want := rec.Body.String(), testHandlerBodyData; have != want {
			t.Errorf("Have: %v Want: %v", have, want)
		}
	}
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
}

func TestExposedHeader(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")

	hpu := cstesting.NewHTTPParallelUsers(10, 4, 300, time.Millisecond)
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
}

func TestAllowedCredentials(t *testing.T, s *cors.Service, req *http.Request) {

	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")

	hpu := cstesting.NewHTTPParallelUsers(5, 5, 200, time.Millisecond)
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
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
	hpu.ServeHTTP(req, s.WithCORS(testHandler()))
}
