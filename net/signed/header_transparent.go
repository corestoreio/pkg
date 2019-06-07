// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"net/http"
	"time"
)

// Cacher defines a custom cache type in conjunction with Transparent hashing.
// It will be used as an underlying storage for the hash values. In middleware
// WithResponseSignature a hash gets added to the Cacher with a time to live
// (TTL). The middleware WithRequestSignatureValidation checks if the hashed
// body matches the cache. The Cacher implementation itself must handle the ttl
// and the constant time comparison.
type Cacher interface {
	Set(hash []byte, ttl time.Duration) error
	Has(hash []byte) bool
}

// Transparent stores the calculated hashes in memory with a TTL using the
// Cacher interface. The hash won't get written into the HTTP response.
type Transparent struct {
	TTL time.Duration
	// Cacher stores hashes for a limited time. Can be nil. Must be set when
	// applying the functional option WithTransparentHashing().
	cache Cacher
}

// MakeTransparent creates a new hash writer. Parse is a noop.
func MakeTransparent(c Cacher, ttl time.Duration) Transparent {
	return Transparent{
		TTL:   ttl,
		cache: c,
	}
}

// HeaderKey returns an empty string.
func (t Transparent) HeaderKey() string { return "" }

// Write sets the signature into the cache with the TTL.
func (t Transparent) Write(_ http.ResponseWriter, signature []byte) {
	_ = t.cache.Set(signature, t.TTL)
}

// Parse returns always nil,nil.
func (t Transparent) Parse(_ *http.Request) ([]byte, error) {
	return nil, nil
}
