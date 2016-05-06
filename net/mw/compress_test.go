// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package mw_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/assert"
)

var testJson string

func init() {
	testJsonB, err := ioutil.ReadFile(filepath.Join("testdata", "config.json"))
	if err != nil {
		panic(err)
	}
	testJson = string(testJsonB)
}

func testCompressReqRes() (w *httptest.ResponseRecorder, r *http.Request) {
	var err error
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		panic(err)
	}
	return
}

func TestWithCompressorNone(t *testing.T) {
	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Empty(t, w.Header().Get(httputil.ContentEncoding))
		assert.Empty(t, w.Header().Get(httputil.Vary))

	}, mw.WithCompressor())

	w, r := testCompressReqRes()
	finalCH.ServeHTTP(w, r)
}

func TestWithCompressorGZIPHeader(t *testing.T) {
	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Exactly(t, httputil.CompressGZIP, w.Header().Get(httputil.ContentEncoding))
		assert.Exactly(t, httputil.AcceptEncoding, w.Header().Get(httputil.Vary))

	}, mw.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(httputil.AcceptEncoding, "deflate, gzip")
	finalCH.ServeHTTP(w, r)
}

func TestWithCompressorDeflateHeader(t *testing.T) {
	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Exactly(t, httputil.CompressDeflate, w.Header().Get(httputil.ContentEncoding))
		assert.Exactly(t, httputil.AcceptEncoding, w.Header().Get(httputil.Vary))

	}, mw.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(httputil.AcceptEncoding, "deflate")
	finalCH.ServeHTTP(w, r)
}

func TestWithCompressorDeflateConcrete(t *testing.T) {
	testWithCompressorConcrete(t, httputil.CompressDeflate, func(r io.Reader) string {
		fr := flate.NewReader(r)
		defer func() {
			if err := fr.Close(); err != nil {
				t.Fatal(err)
			}
		}()
		var un = make([]byte, len(testJson))
		rl, err := fr.Read(un)
		if err != nil {
			t.Error(err)
		}
		if rl != len(testJson) {
			t.Errorf("Read only %d from expected %d bytes. Buffer size: %d", rl, len(testJson), len(un))
		}
		return string(un)
	})
}

func TestWithCompressorGZIPConcrete(t *testing.T) {
	testWithCompressorConcrete(t, httputil.CompressGZIP, func(r io.Reader) string {
		zr, err := gzip.NewReader(r)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := zr.Close(); err != nil {
				t.Fatal(err)
			}
		}()
		var un bytes.Buffer
		if _, err := zr.WriteTo(&un); err != nil {
			t.Fatal(err)
		}
		return un.String()
	})
}

func testWithCompressorConcrete(t *testing.T, header string, uncompressor func(io.Reader) string) {

	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := httputil.NewPrinter(w, r).WriteString(http.StatusOK, testJson); err != nil {
			t.Fatal(err)
		}
	}, mw.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(httputil.AcceptEncoding, header)
	finalCH.ServeHTTP(w, r)
	assert.False(t, len(testJson) < len(w.Body.Bytes()))

	uncompressedBody := uncompressor(w.Body)

	if testJson != uncompressedBody {
		t.Errorf("Want: %d\n\nHave: %d\n", len(testJson), len(uncompressedBody))
		t.Logf("Want: %s\n\nHave: %s\n", testJson, uncompressedBody)
	}
	assert.Exactly(t, header, w.Header().Get(httputil.ContentEncoding))
	assert.Exactly(t, httputil.AcceptEncoding, w.Header().Get(httputil.Vary))
	assert.Exactly(t, httputil.TextPlain, w.Header().Get(httputil.ContentType))

}

// BenchmarkWithCompressorGZIP_1024B-4	   20000	     81916 ns/op	    1330 B/op	       5 allocs/op
func BenchmarkWithCompressorGZIP_1024B(b *testing.B) {

	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := httputil.NewPrinter(w, r).WriteString(http.StatusOK, testJson); err != nil {
			b.Fatal(err)
		}
	}, mw.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(httputil.AcceptEncoding, httputil.CompressGZIP)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		finalCH.ServeHTTP(w, r)
		w.Body.Reset()
	}
}
