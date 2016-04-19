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

package blacklist

import (
	"hash"
	"hash/fnv"
	"sync"
	"time"
)

// Map creates an in-memory map which holds as a key the
// tokens and as value the token expiration duration. Once a Set() operation
// will be called the tokens list get purged. Don't use this feature in
// production as the underlying mutex will become a bottleneck with higher
// throughput, but still faster as a connection to Redis ;-)
type Map struct {
	mu sync.RWMutex
	hash.Hash64
	tokens map[uint64]time.Time
}

// NewMap creates a new blacklist map.
func NewMap() *Map {
	return &Map{
		Hash64: fnv.New64a(),
		tokens: make(map[uint64]time.Time),
	}
}

// hash generates a hash value of a byte slice. not concurrent save
func (bl *Map) hash(token []byte) uint64 {
	bl.Hash64.Reset()
	_, _ = bl.Hash64.Write(token)
	return bl.Hash64.Sum64()

}

// Has checks if a token has been stored in the blacklist and may
// delete the token if expiration time is up.
func (bl *Map) Has(token []byte) bool {

	bl.mu.RLock()
	h := bl.hash(token)
	d, ok := bl.tokens[h]
	bl.mu.RUnlock()

	if !ok {
		return false
	}
	isValid := time.Since(d) < 0

	if false == isValid {
		bl.mu.Lock()
		delete(bl.tokens, h)
		bl.mu.Unlock()
	}
	return isValid
}

// Set adds a token to the blacklist and may perform a
// purge operation. Set should be called when you log out a user.
// Set must make sure to copy away the token bytes or hash them.
func (bl *Map) Set(token []byte, expires time.Duration) error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	h := bl.hash(token)

	for k, v := range bl.tokens {
		if time.Since(v) > 0 {
			delete(bl.tokens, k)
		}
	}
	bl.tokens[h] = time.Now().Add(expires)
	return nil
}

// Len returns the number of entries in the blacklist
func (bl *Map) Len() int {
	bl.mu.RLock()
	l := len(bl.tokens)
	bl.mu.RUnlock()
	return l
}
