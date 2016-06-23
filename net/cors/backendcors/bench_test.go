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

package backendcors_test

import (
	"net/http"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/cors"
	"github.com/corestoreio/csfw/store/scope"
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

func BenchmarkExposedHeader_WebsiteScope_AllowOriginRegex(b *testing.B) {

	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowOriginRegex.MustFQ(scope.Website, 2): "^http://foo",
	}))
	req.Header.Set("Origin", "http://foobar.com")
	//req.Header.Set("Origin", "http://barfoo.com") // not allowed
	// res := httptest.NewRecorder() // only for debugging

	handler := s.WithCORS()(testHandler(b))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		res := FakeResponse{http.Header{}}
		for pb.Next() {
			handler.ServeHTTP(res, req)
		}
	})
}
