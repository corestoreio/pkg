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
	"testing"

	"github.com/corestoreio/csfw/net/ctxcors"
	"golang.org/x/net/context"
)

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
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := testHandler(ctx, res, req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDefault(b *testing.B) {
	res := FakeResponse{http.Header{}}
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Origin", "somedomain.com")
	c, err := ctxcors.New()
	if err != nil {
		b.Fatal(err)
	}
	handler := c.WithCORS()(testHandler)

	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := handler(ctx, res, req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAllowedOrigin(b *testing.B) {
	res := FakeResponse{http.Header{}}
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Origin", "somedomain.com")
	c, err := ctxcors.New(ctxcors.WithAllowedOrigins("somedomain.com"))
	if err != nil {
		b.Fatal(err)
	}
	handler := c.WithCORS()(testHandler)

	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := handler(ctx, res, req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPreflight(b *testing.B) {
	res := FakeResponse{http.Header{}}
	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	req.Header.Add("Access-Control-Request-Method", "GET")
	c, err := ctxcors.New()
	if err != nil {
		b.Fatal(err)
	}
	handler := c.WithCORS()(testHandler)

	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := handler(ctx, res, req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPreflightHeader(b *testing.B) {
	res := FakeResponse{http.Header{}}
	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "Accept")
	c, err := ctxcors.New()
	if err != nil {
		b.Fatal(err)
	}
	handler := c.WithCORS()(testHandler)

	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := handler(ctx, res, req); err != nil {
			b.Fatal(err)
		}
	}
}
