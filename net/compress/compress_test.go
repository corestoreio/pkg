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

package compress_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	csnet "github.com/corestoreio/pkg/net"
	"github.com/corestoreio/pkg/net/compress"
	"github.com/corestoreio/pkg/net/mw"
	"github.com/corestoreio/pkg/net/response"
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

		assert.Empty(t, w.Header().Get(csnet.ContentEncoding))
		assert.Empty(t, w.Header().Get(csnet.Vary))

	}, compress.WithCompressor())

	w, r := testCompressReqRes()
	finalCH.ServeHTTP(w, r)
}

func TestWithCompressorGZIPHeader(t *testing.T) {
	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Exactly(t, csnet.CompressGZIP, w.Header().Get(csnet.ContentEncoding))
		assert.Exactly(t, csnet.AcceptEncoding, w.Header().Get(csnet.Vary))

	}, compress.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(csnet.AcceptEncoding, "deflate, gzip")
	finalCH.ServeHTTP(w, r)
}

func TestWithCompressorDeflateHeader(t *testing.T) {
	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Exactly(t, csnet.CompressDeflate, w.Header().Get(csnet.ContentEncoding))
		assert.Exactly(t, csnet.AcceptEncoding, w.Header().Get(csnet.Vary))

	}, compress.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(csnet.AcceptEncoding, "deflate")
	finalCH.ServeHTTP(w, r)
}

func TestWithCompressorDeflateConcrete(t *testing.T) {
	testWithCompressorConcrete(t, csnet.CompressDeflate, func(r io.Reader) string {
		fr := flate.NewReader(r)
		defer func() {
			if err := fr.Close(); err != nil {
				t.Fatal(err)
			}
		}()
		var un = make([]byte, len(testJson))
		rl, err := fr.Read(un)
		if err != nil && err != io.EOF {
			t.Error(err)
		}
		if rl != len(testJson) {
			t.Errorf("Read only %d from expected %d bytes. Buffer size: %d", rl, len(testJson), len(un))
		}
		return string(un)
	})
}

func TestWithCompressorGZIPConcrete(t *testing.T) {
	testWithCompressorConcrete(t, csnet.CompressGZIP, func(r io.Reader) string {
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
		if err := response.NewPrinter(w, r).WriteString(http.StatusOK, testJson); err != nil {
			t.Fatal(err)
		}
	}, compress.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(csnet.AcceptEncoding, header)
	finalCH.ServeHTTP(w, r)
	assert.False(t, len(testJson) < len(w.Body.Bytes()))

	uncompressedBody := uncompressor(w.Body)

	if testJson != uncompressedBody {
		t.Errorf("Want: %d\n\nHave: %d\n", len(testJson), len(uncompressedBody))
		t.Logf("Want: %s\n\nHave: %s\n", testJson, uncompressedBody)
	}
	assert.Exactly(t, header, w.Header().Get(csnet.ContentEncoding))
	assert.Exactly(t, csnet.AcceptEncoding, w.Header().Get(csnet.Vary))
	assert.Exactly(t, csnet.TextPlain, w.Header().Get(csnet.ContentType))

}

// BenchmarkWithCompressorGZIP_1024B-4	   20000	     81916 ns/op	    1330 B/op	       5 allocs/op
func BenchmarkWithCompressorGZIP_1024B(b *testing.B) {

	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := response.NewPrinter(w, r).WriteString(http.StatusOK, testJson); err != nil {
			b.Fatal(err)
		}
	}, compress.WithCompressor())

	w, r := testCompressReqRes()
	r.Header.Set(csnet.AcceptEncoding, csnet.CompressGZIP)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		finalCH.ServeHTTP(w, r)
		w.Body.Reset()
	}
}
