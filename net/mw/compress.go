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

package mw

import (
	"bufio"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/corestoreio/csfw/net/httputil"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
)

var gzWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(ioutil.Discard)
	},
}

var defWriterPool = sync.Pool{
	New: func() interface{} {
		w, err := flate.NewWriter(ioutil.Discard, 2)
		if err != nil {
			panic(err)
		}
		return w
	},
}

type compressWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w compressWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w compressWriter) Write(b []byte) (int, error) {
	if w.Header().Get(httputil.ContentType) == "" {
		w.Header().Set(httputil.ContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w compressWriter) Flush() error {
	if f, ok := w.Writer.(*gzip.Writer); ok {
		return f.Flush()
	}
	if f, ok := w.Writer.(*flate.Writer); ok {
		return f.Flush()
	}
	return nil
}

func (w compressWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *compressWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// WithCompressor is a middleware applies the GZIP or deflate algorithm on
// the bytes writer. GZIP or deflate usage depends on the HTTP Accept
// Encoding header. Flush(), Hijack() and CloseNotify() interfaces will be
// preserved. No header set, no compression takes place. GZIP has priority
// before deflate.
func WithCompressor() Middleware {

	// todo(cs): maybe the sync.Pools can be put in here because then
	// the developer can set the deflate compression level.
	// todo(cs) handle compression depending on the website ...

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enc := r.Header.Get(httputil.AcceptEncoding)

			if strings.Contains(enc, httputil.CompressGZIP) {
				w.Header().Set(httputil.ContentEncoding, httputil.CompressGZIP)
				w.Header().Add(httputil.Vary, httputil.AcceptEncoding)

				zw := gzWriterPool.Get().(*gzip.Writer)
				zw.Reset(w)
				defer func() {
					zw.Close()
					gzWriterPool.Put(zw)
				}()
				cw := compressWriter{Writer: zw, ResponseWriter: w}
				h.ServeHTTP(cw, r)
				return
			}

			if strings.Contains(enc, httputil.CompressDeflate) {
				w.Header().Set(httputil.ContentEncoding, httputil.CompressDeflate)
				w.Header().Add(httputil.Vary, httputil.AcceptEncoding)

				zw := defWriterPool.Get().(*flate.Writer)
				zw.Reset(w)
				defer func() {
					zw.Close()
					defWriterPool.Put(zw)
				}()
				cw := compressWriter{Writer: zw, ResponseWriter: w}
				h.ServeHTTP(cw, r)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
