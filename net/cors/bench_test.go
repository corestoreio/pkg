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
	"testing"

	"github.com/corestoreio/csfw/net/cors"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

func testHandler(fa interface {
	Fatal(args ...interface{})
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := cors.FromContext(r.Context()); err != nil {
			fa.Fatal(errors.PrintLoc(err))
		}
		_, _ = w.Write([]byte("bar"))
	}
}

type FakeResponse struct {
	header http.Header
}

func (r FakeResponse) Header() http.Header {
	return r.header
}

func (r FakeResponse) WriteHeader(n int) {
}

func (r FakeResponse) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func BenchmarkWithout(b *testing.B) {
	res := FakeResponse{http.Header{}}
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	h := testHandler(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h(res, req)
	}
}

func BenchmarkDefault(b *testing.B) {
	res := FakeResponse{http.Header{}}
	req := reqWithStore("GET")
	req.Header.Add("Origin", "somedomain.com")
	c, err := cors.New()
	if err != nil {
		b.Fatal(err)
	}
	handler := c.WithCORS()(testHandler(b))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(res, req)
	}
}

func BenchmarkAllowedOrigin(b *testing.B) {
	res := FakeResponse{http.Header{}}
	req := reqWithStore("GET")
	req.Header.Add("Origin", "somedomain.com")
	c, err := cors.New(cors.WithAllowedOrigins(scope.Default, 0, "somedomain.com"))
	if err != nil {
		b.Fatal(err)
	}
	handler := c.WithCORS()(testHandler(b))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(res, req)
	}
}

func BenchmarkPreflight(b *testing.B) {
	res := FakeResponse{http.Header{}}
	req := reqWithStore("OPTIONS")
	req.Header.Add("Access-Control-Request-Method", "GET")
	c, err := cors.New()
	if err != nil {
		b.Fatal(err)
	}
	handler := c.WithCORS()(testHandler(b))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(res, req)
	}
}

func BenchmarkPreflightHeader(b *testing.B) {
	res := FakeResponse{http.Header{}}
	req := reqWithStore("OPTIONS")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "Accept")
	c, err := cors.New()
	if err != nil {
		b.Fatal(err)
	}
	handler := c.WithCORS()(testHandler(b))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(res, req)
	}
}
