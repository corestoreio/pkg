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

package ctxmw_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/net/ctxmw"
	"golang.org/x/net/context"
)

var testHandler = func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
	_, err := w.Write([]byte("bar"))
	return err
}

func assertHeaders(t *testing.T, resHeaders http.Header, reqHeaders map[string]string) {
	for name, value := range reqHeaders {
		if actual := strings.Join(resHeaders[name], ", "); actual != value {
			t.Errorf("Invalid header `%s', wanted `%s', got `%s'", name, value, actual)
		}
	}
}

func TestCorsNoConfig(t *testing.T) {
	s := ctxmw.NewCors(nil)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	ctx := context.Background()

	s.WithCORS()(testHandler)(ctx, res, req) // yay that looks terrible!

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

func TestCorsMatchAllOrigin(t *testing.T) {
	s := ctxmw.NewCors(ctxmw.WithCorsAllowedOrigins("*"))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Origin", "http://foobar.com")
	ctx := context.Background()

	s.WithCORS()(testHandler)(ctx, res, req)

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

//func TestAllowedOrigin(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://foobar.com"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foobar.com")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin",
//		"Access-Control-Allow-Origin":      "http://foobar.com",
//		"Access-Control-Allow-Methods":     "",
//		"Access-Control-Allow-Headers":     "",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
//
//func TestWildcardOrigin(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://*.bar.com"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foo.bar.com")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin",
//		"Access-Control-Allow-Origin":      "http://foo.bar.com",
//		"Access-Control-Allow-Methods":     "",
//		"Access-Control-Allow-Headers":     "",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
//
//func TestDisallowedOrigin(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://foobar.com"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://barbaz.com")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin",
//		"Access-Control-Allow-Origin":      "",
//		"Access-Control-Allow-Methods":     "",
//		"Access-Control-Allow-Headers":     "",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
//
//func TestDisallowedWildcardOrigin(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://*.bar.com"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foo.baz.com")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin",
//		"Access-Control-Allow-Origin":      "",
//		"Access-Control-Allow-Methods":     "",
//		"Access-Control-Allow-Headers":     "",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
//
//func TestAllowedOriginFunc(t *testing.T) {
//	r, _ := regexp.Compile("^http://foo")
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowOriginFunc: func(o string) bool {
//			return r.MatchString(o)
//		},
//	})
//
//	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
//
//	res := httptest.NewRecorder()
//	req.Header.Set("Origin", "http://foobar.com")
//	s.Handler(testHandler).ServeHTTP(res, req)
//	assertHeaders(t, res.Header(), map[string]string{
//		"Access-Control-Allow-Origin": "http://foobar.com",
//	})
//
//	res = httptest.NewRecorder()
//	req.Header.Set("Origin", "http://barfoo.com")
//	s.Handler(testHandler).ServeHTTP(res, req)
//	assertHeaders(t, res.Header(), map[string]string{
//		"Access-Control-Allow-Origin": "",
//	})
//}
//
//func TestAllowedMethod(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://foobar.com"},
//		AllowedMethods: []string{"PUT", "DELETE"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foobar.com")
//	req.Header.Add("Access-Control-Request-Method", "PUT")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
//		"Access-Control-Allow-Origin":      "http://foobar.com",
//		"Access-Control-Allow-Methods":     "PUT",
//		"Access-Control-Allow-Headers":     "",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
//
//func TestDisallowedMethod(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://foobar.com"},
//		AllowedMethods: []string{"PUT", "DELETE"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foobar.com")
//	req.Header.Add("Access-Control-Request-Method", "PATCH")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
//		"Access-Control-Allow-Origin":      "",
//		"Access-Control-Allow-Methods":     "",
//		"Access-Control-Allow-Headers":     "",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
//
//func TestAllowedHeader(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://foobar.com"},
//		AllowedHeaders: []string{"X-Header-1", "x-header-2"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foobar.com")
//	req.Header.Add("Access-Control-Request-Method", "GET")
//	req.Header.Add("Access-Control-Request-Headers", "X-Header-2, X-HEADER-1")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
//		"Access-Control-Allow-Origin":      "http://foobar.com",
//		"Access-Control-Allow-Methods":     "GET",
//		"Access-Control-Allow-Headers":     "X-Header-2, X-Header-1",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
//
//func TestAllowedWildcardHeader(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://foobar.com"},
//		AllowedHeaders: []string{"*"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foobar.com")
//	req.Header.Add("Access-Control-Request-Method", "GET")
//	req.Header.Add("Access-Control-Request-Headers", "X-Header-2, X-HEADER-1")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
//		"Access-Control-Allow-Origin":      "http://foobar.com",
//		"Access-Control-Allow-Methods":     "GET",
//		"Access-Control-Allow-Headers":     "X-Header-2, X-Header-1",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
//
//func TestDisallowedHeader(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://foobar.com"},
//		AllowedHeaders: []string{"X-Header-1", "x-header-2"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foobar.com")
//	req.Header.Add("Access-Control-Request-Method", "GET")
//	req.Header.Add("Access-Control-Request-Headers", "X-Header-3, X-Header-1")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
//		"Access-Control-Allow-Origin":      "",
//		"Access-Control-Allow-Methods":     "",
//		"Access-Control-Allow-Headers":     "",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
//
//func TestOriginHeader(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://foobar.com"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foobar.com")
//	req.Header.Add("Access-Control-Request-Method", "GET")
//	req.Header.Add("Access-Control-Request-Headers", "origin")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
//		"Access-Control-Allow-Origin":      "http://foobar.com",
//		"Access-Control-Allow-Methods":     "GET",
//		"Access-Control-Allow-Headers":     "Origin",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
//
//func TestExposedHeader(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins: []string{"http://foobar.com"},
//		ExposedHeaders: []string{"X-Header-1", "x-header-2"},
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foobar.com")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin",
//		"Access-Control-Allow-Origin":      "http://foobar.com",
//		"Access-Control-Allow-Methods":     "",
//		"Access-Control-Allow-Headers":     "",
//		"Access-Control-Allow-Credentials": "",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "X-Header-1, X-Header-2",
//	})
//}
//
//func TestAllowedCredentials(t *testing.T) {
//	s := ctxmw.New(ctxmw.CorsOptions{
//		AllowedOrigins:   []string{"http://foobar.com"},
//		AllowCredentials: true,
//	})
//
//	res := httptest.NewRecorder()
//	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
//	req.Header.Add("Origin", "http://foobar.com")
//	req.Header.Add("Access-Control-Request-Method", "GET")
//
//	s.Handler(testHandler).ServeHTTP(res, req)
//
//	assertHeaders(t, res.Header(), map[string]string{
//		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
//		"Access-Control-Allow-Origin":      "http://foobar.com",
//		"Access-Control-Allow-Methods":     "GET",
//		"Access-Control-Allow-Headers":     "",
//		"Access-Control-Allow-Credentials": "true",
//		"Access-Control-Max-Age":           "",
//		"Access-Control-Expose-Headers":    "",
//	})
//}
