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

package signed

import (
	"io"
	"net/http"

	"crypto/hmac"

	"github.com/corestoreio/csfw/net/responseproxy"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/hashpool"
)

// ScopedConfig scoped based configuration and should not be embedded into your
// own types. Call ScopedConfig.ScopeHash to know to which scope this
// configuration has been bound to.
type ScopedConfig struct {
	scopedConfigGeneric

	// start of package specific config values

	// Disabled set to true to disable content signing.
	Disabled bool
	// InTrailer set to true and the signature will be added to the HTTP Trailer for
	// responses.
	InTrailer bool

	// HTTPWriter writes a signature into a header or trailer to the HTTP response.
	HTTPWriter
	// HTTPParser reads the header or trailer from a request and extracts the raw
	// signature.
	HTTPParser

	hashPool hashpool.Tank

	// AllowedMethods list of allowed HTTP methods. Must be upper case.
	AllowedMethods []string
}

// newScopedConfig creates a new object with the minimum needed configuration.
func newScopedConfig() *ScopedConfig {

	return &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(),
		AllowedMethods:      []string{"POST", "PUT", "PATCH"},
	}
}

// IsValid a configuration for a scope is only then valid when several fields
// are not empty: RateLimiter, DeniedHandler and VaryByer.
func (sc *ScopedConfig) IsValid() error {
	if sc.lastErr != nil {
		return errors.Wrap(sc.lastErr, "[signed] scopedConfig.isValid has an lastErr")
	}
	if sc.ScopeHash == 0 {
		return errors.NewNotValidf(errScopedConfigNotValid, sc.ScopeHash)
	}

	return nil
}

// direct output to the client and the signature will be inserted after the body
// has been written. ideal for streaming but not all clients can process a
// trailer.
func (sc *ScopedConfig) writeTrailer(next http.Handler, w http.ResponseWriter, r *http.Request) {
	h := sc.hashPool.Get()
	defer sc.hashPool.Put(h)

	wt := responseproxy.WrapTee(w)
	wt.Tee(h) // write also to hash

	wt.Header().Add("Trailer", sc.HTTPWriter.HeaderKey())

	next.ServeHTTP(wt, r)

	buf := bufferpool.Get()
	tmp := h.Sum(buf.Bytes()) // append to buffer
	buf.Reset()
	_, _ = buf.Write(tmp)
	sc.HTTPWriter.Write(w, buf.Bytes())
	bufferpool.Put(buf)
}

// the write to w gets buffered and we calculate the checksum of the buffer and
// then flush the buffer to the client.
func (sc *ScopedConfig) writeBuffered(next http.Handler, w http.ResponseWriter, r *http.Request) {
	h := sc.hashPool.Get()
	defer sc.hashPool.Put(h)

	wBuf := bufferpool.Get()
	hashBuf := bufferpool.Get()
	defer bufferpool.Put(hashBuf)
	defer bufferpool.Put(wBuf)

	next.ServeHTTP(responseproxy.WrapBuffered(wBuf, w), r)

	// calculate the hash based on the buffered response body

	if _, err := h.Write(wBuf.Bytes()); err != nil {
		sc.ErrorHandler(errors.Wrap(err, "[signed] ScopedConfig.writeBuffered failed to io.Copy")).ServeHTTP(w, r)
		return
	}

	sc.HTTPWriter.Write(w, h.Sum(hashBuf.Bytes()))
	if _, err := io.Copy(w, wBuf); err != nil {
		sc.ErrorHandler(errors.Wrap(err, "[signed] ScopedConfig.writeBuffered failed to io.Copy")).ServeHTTP(w, r)
	}
}

func (sc *ScopedConfig) ValidateBody(r *http.Request) error {

	signature, err := sc.HTTPParser.Parse(r)
	if err != nil {
		return errors.Wrap(err, "[signed] ValidateBody HTTPParser.Parse")
	}

	ok := false
	for _, m := range sc.AllowedMethods {
		if r.Method == m {
			ok = true
		}
	}
	if !ok {
		return errors.NewNotValidf("[signed] ValidateBody HTTP Method %q not allowed in list: %q", r.Method, sc.AllowedMethods)
	}

	buf := make([]byte, 1024)
	h := sc.hashPool.Get()
	defer sc.hashPool.Put(h)
	defer r.Body.Close() // what happens to any unread bytes of the body?

	for {
		n, err := r.Body.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "[signed] ValidateBody HTTP.Body.Read")
		}
		if _, err := h.Write(buf[:n]); err != nil {
			return errors.Wrap(err, "[signed] ValidateBody Hash.Write")
		}
	}
	hashBuf := bufferpool.Get()
	defer bufferpool.Put(hashBuf)

	if hs := h.Sum(hashBuf.Bytes()); !hmac.Equal(signature, hs) {
		return errors.NewNotValidf("[signed] ValidateBody. Signatures do not match. Have: %q Want: %q", signature, hs)
	}

	return nil
}
