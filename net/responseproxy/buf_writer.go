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

package responseproxy

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

// BufferedWriter is a proxy around an http.ResponseWriter that allows you write
// the response body into a buffer instead of directly writing to the client.
// You must take care that the buffer gets written to the client. Calls to
// Flush() do not flush the buffer. Headers gets written directly to the client.
type BufferedWriter interface {
	http.ResponseWriter
	// Unwrap returns the original proxied target.
	Unwrap() http.ResponseWriter
}

// WrapBuffered wraps an http.ResponseWriter, returning a proxy which only writes
// into the provided io.Writer.
func WrapBuffered(buf io.Writer, w http.ResponseWriter) BufferedWriter {
	_, cn := w.(http.CloseNotifier)
	_, fl := w.(http.Flusher)
	_, hj := w.(http.Hijacker)
	_, rf := w.(io.ReaderFrom)

	bw := bufferedWriter{
		ResponseWriter: w,
		buf:            buf,
	}
	if cn && fl && hj && rf {
		return &bufferedFancyWriter{bw}
	}
	if fl {
		return &bufferedFlushWriter{bw}
	}
	return &bw
}

// bufferedWriter wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface.
type bufferedWriter struct {
	http.ResponseWriter
	buf io.Writer
}

// Write does not write to the client instead it writes in the underlying buffer.
func (b *bufferedWriter) Write(buf []byte) (int, error) {
	b.WriteHeader(http.StatusOK)
	return b.buf.Write(buf)
}

// Unwrap returns the original underlying ResponseWriter.
func (b *bufferedWriter) Unwrap() http.ResponseWriter {
	return b.ResponseWriter
}

// bufferedFancyWriter is a writer that additionally satisfies
// http.CloseNotifier, http.Flusher, http.Hijacker, and io.ReaderFrom. It exists
// for the common case of wrapping the http.ResponseWriter that package http
// gives you, in order to make the proxied object support the full method set of
// the proxied object.
type bufferedFancyWriter struct {
	bufferedWriter
}

func (f *bufferedFancyWriter) CloseNotify() <-chan bool {
	cn := f.bufferedWriter.ResponseWriter.(http.CloseNotifier)
	return cn.CloseNotify()
}
func (f *bufferedFancyWriter) Flush() {
	fl := f.bufferedWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}
func (f *bufferedFancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := f.bufferedWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

// ReadFrom writes r into the underlying buffer
func (f *bufferedFancyWriter) ReadFrom(r io.Reader) (int64, error) {
	return io.Copy(&f.bufferedWriter, r)
}

var _ http.CloseNotifier = &bufferedFancyWriter{}
var _ http.Flusher = &bufferedFancyWriter{}
var _ http.Hijacker = &bufferedFancyWriter{}
var _ io.ReaderFrom = &bufferedFancyWriter{}
var _ http.Flusher = &flushWriter{}

// bufferedFlushWriter implements only http.Flusher mostly used
type bufferedFlushWriter struct {
	bufferedWriter
}

func (f *bufferedFlushWriter) Flush() {
	fl := f.bufferedWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}
