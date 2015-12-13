// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ctxmw_test

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxmw"
	"github.com/corestoreio/csfw/net/httputil"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func testCompressReqRes() (w *httptest.ResponseRecorder, r *http.Request) {
	var err error
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		panic(err)
	}
	return
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestWithCompressorNone(t *testing.T) {
	finalCH := ctxhttp.Chain(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		assert.Empty(t, w.Header().Get(httputil.ContentEncoding))
		assert.Empty(t, w.Header().Get(httputil.Vary))

		return nil
	}, ctxmw.WithCompressor())

	w, r := testCompressReqRes()
	if err := finalCH.ServeHTTPContext(context.TODO(), w, r); err != nil {
		t.Fatal(err)
	}
}

func TestWithCompressorGZIPHeader(t *testing.T) {
	finalCH := ctxhttp.Chain(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		assert.Exactly(t, httputil.CompressGZIP, w.Header().Get(httputil.ContentEncoding))
		assert.Exactly(t, httputil.AcceptEncoding, w.Header().Get(httputil.Vary))

		return nil
	}, ctxmw.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(httputil.AcceptEncoding, "deflate, gzip")
	if err := finalCH.ServeHTTPContext(context.TODO(), w, r); err != nil {
		t.Fatal(err)
	}
}

func TestWithCompressorDeflateHeader(t *testing.T) {
	finalCH := ctxhttp.Chain(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		assert.Exactly(t, httputil.CompressDeflate, w.Header().Get(httputil.ContentEncoding))
		assert.Exactly(t, httputil.AcceptEncoding, w.Header().Get(httputil.Vary))
		return nil
	}, ctxmw.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(httputil.AcceptEncoding, "deflate")
	if err := finalCH.ServeHTTPContext(context.TODO(), w, r); err != nil {
		t.Fatal(err)
	}
}

func TestWithCompressorGZIPConcrete(t *testing.T) {
	testWithCompressorConcrete(t, httputil.CompressGZIP, func(r io.Reader) string {
		zr, err := gzip.NewReader(r)
		assert.NoError(t, err)
		defer zr.Close()
		var un bytes.Buffer
		zr.WriteTo(&un)
		return un.String()
	})
}

func TestWithCompressorDeflateConcrete(t *testing.T) {
	testWithCompressorConcrete(t, httputil.CompressDeflate, func(r io.Reader) string {
		fr := flate.NewReader(r)
		defer fr.Close()
		var un = make([]byte, 1024)
		fr.Read(un)
		return string(un)
	})
}

func testWithCompressorConcrete(t *testing.T, header string, uncompressor func(io.Reader) string) {

	rawData := randSeq(1024)

	finalCH := ctxhttp.Chain(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return httputil.NewPrinter(w, r).WriteString(http.StatusOK, rawData)
	}, ctxmw.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(httputil.AcceptEncoding, header)
	if err := finalCH.ServeHTTPContext(context.TODO(), w, r); err != nil {
		t.Fatal(err)
	}
	assert.False(t, len(rawData) < len(w.Body.Bytes()))

	uncompressedBody := uncompressor(w.Body)

	assert.Exactly(t, rawData, uncompressedBody)
	assert.Exactly(t, header, w.Header().Get(httputil.ContentEncoding))
	assert.Exactly(t, httputil.AcceptEncoding, w.Header().Get(httputil.Vary))
	assert.Exactly(t, httputil.TextPlain, w.Header().Get(httputil.ContentType))

}

// BenchmarkWithCompressorGZIP_1024B-4	   20000	     81916 ns/op	    1330 B/op	       5 allocs/op
func BenchmarkWithCompressorGZIP_1024B(b *testing.B) {

	rawData := randSeq(1024)

	finalCH := ctxhttp.Chain(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return httputil.NewPrinter(w, r).WriteString(http.StatusOK, rawData)
	}, ctxmw.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(httputil.AcceptEncoding, httputil.CompressGZIP)

	ctx := context.TODO()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := finalCH.ServeHTTPContext(ctx, w, r); err != nil {
			b.Fatal(err)
		}
		w.Body.Reset()
	}
}
