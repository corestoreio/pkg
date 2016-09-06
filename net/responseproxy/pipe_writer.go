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
	"sync"

	"github.com/corestoreio/csfw/util/errors"
)

// PipedWriter is a proxy around an http.ResponseWriter that allows you write
// the response body into a pipe instead of directly writing to the client. You
// must take care that the pipe output (the reader) gets written to the client.
// Calls to Flush() do not flush the pipe. Headers gets written directly to the
// client.
type PipedWriter interface {
	http.ResponseWriter
	// Close must be called to terminate the internal goroutine.
	Close() error
	// Unwrap returns the original proxied target.
	Unwrap() http.ResponseWriter
}

// WrapPiped wraps an http.ResponseWriter, returning a proxy which pipes the
// writes to your ReaderFrom. You must call "defer Close()" to quite the internal
// goroutine.
func WrapPiped(iorf io.ReaderFrom, w http.ResponseWriter) PipedWriter {
	_, cn := w.(http.CloseNotifier)
	_, fl := w.(http.Flusher)
	_, hj := w.(http.Hijacker)
	_, rf := w.(io.ReaderFrom)

	pr, pw := io.Pipe()
	mw := pipedWriter{ResponseWriter: w, pw: pw}
	mw.wg.Add(1)
	go func() {
		if _, err := iorf.ReadFrom(pr); err != nil {
			mw.err = err
			pr.CloseWithError(err)
		} else {
			pr.Close()
		}
		mw.wg.Done()
	}()

	if cn && fl && hj && rf {
		return &pipedFancyWriter{mw}
	}
	if fl {
		return &pipedFlushWriter{mw}
	}
	return &mw
}

// bufferedWriter wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface.
type pipedWriter struct {
	http.ResponseWriter
	pw  *io.PipeWriter
	wg  sync.WaitGroup
	err error
}

// Write does not write to the client instead it writes in the underlying buffer.
func (b *pipedWriter) Write(p []byte) (int, error) {
	b.WriteHeader(http.StatusOK)
	return b.pw.Write(p)
}

// Unwrap returns the original underlying ResponseWriter.
func (b *pipedWriter) Unwrap() http.ResponseWriter {
	return b.ResponseWriter
}

func (b *pipedWriter) Close() error {
	if err := b.pw.Close(); err != nil {
		return errors.Wrap(err, "[responseproxy] pipeWriter.pw.Close")
	}
	b.wg.Wait()
	return errors.Wrap(b.err, "[responseproxy] pipeWriter.Close")
}

// pipedFancyWriter is a writer that additionally satisfies http.CloseNotifier,
// http.Flusher, http.Hijacker, and io.ReaderFrom. It exists for the common case
// of wrapping the http.ResponseWriter that package http gives you, in order to
// make the proxied object support the full method set of the proxied object.
type pipedFancyWriter struct {
	pipedWriter
}

func (f *pipedFancyWriter) CloseNotify() <-chan bool {
	cn := f.pipedWriter.ResponseWriter.(http.CloseNotifier)
	return cn.CloseNotify()
}
func (f *pipedFancyWriter) Flush() {
	fl := f.pipedWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}
func (f *pipedFancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := f.pipedWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

// ReadFrom writes r into the underlying pipe
func (f *pipedFancyWriter) ReadFrom(r io.Reader) (int64, error) {
	return io.Copy(&f.pipedWriter, r)
}

var _ http.CloseNotifier = &pipedFancyWriter{}
var _ http.Flusher = &pipedFancyWriter{}
var _ http.Hijacker = &pipedFancyWriter{}
var _ io.ReaderFrom = &pipedFancyWriter{}
var _ http.Flusher = &flushWriter{}

// pipedFlushWriter implements only http.Flusher mostly used
type pipedFlushWriter struct {
	pipedWriter
}

func (f *pipedFlushWriter) Flush() {
	fl := f.pipedWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}
