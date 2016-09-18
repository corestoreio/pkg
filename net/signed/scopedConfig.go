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
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/corestoreio/csfw/net/responseproxy"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/hashpool"
)

// DefaultHashName identifies the default hash when creating a new scoped
// configuration. You must register this name before using this package via:
//		hashpool.Register(`sha256`, sha256.New)
// If you would like to use different hashes you must registered them also in
// the hashpool package.
const DefaultHashName = `sha256`

// ScopedConfig scoped based configuration and should not be embedded into your
// own types. Call ScopedConfig.ScopeHash to know to which scope this
// configuration has been bound to.
type ScopedConfig struct {
	scopedConfigGeneric
	// start of package specific config values
	hashPool hashpool.Tank

	// Disabled set to true to disable content signing.
	Disabled bool
	// InTrailer set to true and the signature will be added to the HTTP Trailer for
	// responses.
	InTrailer bool
	// HeaderParseWriter see description of interface HeaderParseWriter
	HeaderParseWriter
	// AllowedMethods list of allowed HTTP methods. Must be upper case.
	AllowedMethods []string
	// TransparentCacher stores the calculated hashes in memory with a TTL. The hash
	// won't get written into the HTTP response. If enable you must set the
	// Cacher field in the Service struct.
	TransparentCacher Cacher
	// TransparentTTL defines the time to live for a hash within the Cacher
	// interface.
	TransparentTTL time.Duration
}

// newScopedConfig creates a new object with the minimum needed configuration.
// Acts also as WithDefaultConfig()
// Default settings: InTrailer activated, Content-HMAC header with sha256,
// allowed HTTP methods set to POST, PUT, PATCH and password for the HMAC SHA
// 256 from a cryptographically random source with a length of 64 bytes.
func newScopedConfig() *ScopedConfig {
	sc := &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(),
		InTrailer:           true,
		HeaderParseWriter:   NewContentHMAC(DefaultHashName),
		AllowedMethods:      []string{"POST", "PUT", "PATCH"},
	}
	key := make([]byte, 64) // 64 character password
	if _, err := rand.Read(key); err != nil {
		sc.lastErr = errors.Wrap(err, "[signed] newScopedConfig: Failed to cread from crypto/rand.Read")
		// don't init hashpool and let app panic
	} else {
		sc.hashPoolInit(DefaultHashName, key)
	}
	return sc
}

func (sc *ScopedConfig) hashPoolInit(name string, key []byte) {
	sc.hashPool, sc.lastErr = hashpool.FromRegistryHMAC(name, key)
	sc.lastErr = errors.Wrapf(sc.lastErr, "[signed] The hash %q has not yet been registered via hashpool.Register() function.", name)
}

// IsValid a configuration for a scope is only then valid when several fields
// are not empty: HeaderParseWriter and AllowedMethods OR disabled for current
// scope.
func (sc *ScopedConfig) IsValid() error {
	if sc.lastErr != nil {
		return errors.Wrap(sc.lastErr, "[signed] scopedConfig.isValid has an lastErr")
	}
	if sc.Disabled {
		return nil
	}
	if sc.ScopeHash == 0 || sc.HeaderParseWriter == nil || len(sc.AllowedMethods) == 0 {
		return errors.NewNotValidf(errScopedConfigNotValid, sc.ScopeHash, sc.HeaderParseWriter == nil, sc.AllowedMethods)
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

	if k := sc.HeaderParseWriter.HeaderKey(); k != "" {
		wt.Header().Add("Trailer", k)
	}
	next.ServeHTTP(wt, r)

	buf := bufferpool.Get()
	sc.HeaderParseWriter.Write(w, h.Sum(buf.Bytes()))
	bufferpool.Put(buf)
}

// the write to w gets buffered and we calculate the checksum of the buffer and
// then flush the buffer to the client.
func (sc *ScopedConfig) writeBuffered(next http.Handler, w http.ResponseWriter, r *http.Request) {
	h := sc.hashPool.Get()
	defer sc.hashPool.Put(h)

	rwBuf := bufferpool.Get()
	hBuf := bufferpool.Get()
	defer bufferpool.Put(hBuf)
	defer bufferpool.Put(rwBuf)

	next.ServeHTTP(responseproxy.WrapBuffered(rwBuf, w), r)

	// calculate the hash based on the buffered response body

	if _, err := h.Write(rwBuf.Bytes()); err != nil {
		sc.ErrorHandler(errors.Wrap(err, "[signed] ScopedConfig.writeBuffered failed to io.Copy")).ServeHTTP(w, r)
		return
	}

	sc.HeaderParseWriter.Write(w, h.Sum(hBuf.Bytes()))
	if _, err := io.Copy(w, rwBuf); err != nil {
		sc.ErrorHandler(errors.Wrap(err, "[signed] ScopedConfig.writeBuffered failed to io.Copy")).ServeHTTP(w, r)
	}
}

// CalculateHash calculates the hash sum from the request body. The full body
// gets read into a buffer. This buffer gets assigned to the r.Body to make a
// read possible for the next consumer.
func (sc *ScopedConfig) CalculateHash(r *http.Request) ([]byte, error) {

	h := sc.hashPool.Get()
	defer sc.hashPool.Put(h)
	defer r.Body.Close()

	// copy the body so that the next consumer can read it.
	body := new(bytes.Buffer)
	buf := make([]byte, 4096) // maybe make it configurable ...
	for {
		n, err := r.Body.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "[signed] ValidateBody HTTP.Body.Read")
		}
		if _, err := h.Write(buf[:n]); err != nil {
			return nil, errors.Wrap(err, "[signed] ValidateBody Hash.Write")
		}
		_, _ = body.Write(buf[:n])
	}

	r.Body = ioutil.NopCloser(body)

	buf = buf[:0]
	return h.Sum(buf), nil
}

func (sc *ScopedConfig) isMethodAllowed(reqMethod string) bool {
	for _, m := range sc.AllowedMethods {
		if reqMethod == m {
			return true
		}
	}
	return false
}

// ValidateBody uses the HTTPParser to extract the hash signature. It then
// hashes the body and compares the hash of the body with the hash value found
// in the HTTP header. Hash comparison via constant time.
func (sc *ScopedConfig) ValidateBody(r *http.Request) error {

	if !sc.isMethodAllowed(r.Method) {
		return errors.NewNotValidf(errScopedConfigMethodNotAllowed, r.Method, sc.AllowedMethods)
	}

	hashSum, err := sc.CalculateHash(r)
	if err != nil {
		return errors.Wrap(err, "[signed] ScopedConfig.ValidateBody.calculateHash")
	}

	// check if we're using transparent hashing and store the hash values in a
	// cache. constant time not implement and responsibility of the Cacher
	// implementation.
	if sc.TransparentCacher != nil {
		if sc.TransparentCacher.Has(hashSum) {
			return nil
		}
		return errors.NewNotValidf(errScopedConfigCacheNotFound, hashSum)
	}

	// parse the header to find the signature
	reqHashSum, err := sc.HeaderParseWriter.Parse(r)
	if err != nil {
		return errors.Wrap(err, "[signed] ValidateBody HTTPParser.Parse")
	}
	if !hmac.Equal(reqHashSum, hashSum) {
		return errors.NewNotValidf(errScopedConfigSignatureNoMatch, reqHashSum, hashSum)
	}

	return nil
}
