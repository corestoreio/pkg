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

package containable

import (
	"sync"
	"time"

	"sync/atomic"
)

// Container allows to check if a value, identified by a key, has been previously
// seen. Must be thread safe.
type Container interface {
	// Set adds the ID to the container and may perform a purge operation. Set
	// must make sure to copy away the bytes or hash them.
	Set(id []byte, expires time.Duration) error
	// Has checks if an ID has been stored in the container and may delete the
	// ID if the expiration time is up.
	Has(id []byte) bool
}

// Mock implements interface Container and allows mocking it in tests.
type Mock struct {
	SetFn func(hash []byte, ttl time.Duration) error
	HasFn func(hash []byte) bool
}

func (cm Mock) Set(hash []byte, ttl time.Duration) error { return cm.SetFn(hash, ttl) }
func (cm Mock) Has(hash []byte) bool                     { return cm.HasFn(hash) }

// InMemory creates an in-memory map which holds as a key the ID and as value an
// expiration duration. Once a Set() operation will be called the ID list get
// purged. The map type has been optimized for less GC and can hold millions of
// entries.
type InMemory struct {
	mu sync.RWMutex
	// keys contains a map consisting only of integers which skips scanning a
	// map by the GC.
	keys map[string]int64 // int64 unix timestamp
	// Map access for map[string([]byte)] has been optmized in ~Go 1.6
	shouldPurge uint32 // internal counter
}

const purgeEveryNTimes uint32 = 5

// NewInMemory creates a new in memory map.
func NewInMemory() *InMemory {
	return &InMemory{
		keys: make(map[string]int64),
	}
}

// Has checks if an ID has been stored in the map and may delete the ID if
// expiration time is up.
func (bl *InMemory) Has(id []byte) bool {

	bl.mu.RLock()
	ts, ok := bl.keys[string(id)]
	bl.mu.RUnlock()

	if !ok {
		return false
	}
	isValid := time.Now().Unix() < ts

	if false == isValid {
		bl.mu.Lock()
		delete(bl.keys, string(id))
		bl.mu.Unlock()
	}
	return isValid
}

// Set adds an ID to the map and may perform a purge operation every fifth
// access time.
func (bl *InMemory) Set(id []byte, expires time.Duration) error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	if atomic.AddUint32(&bl.shouldPurge, 1)%purgeEveryNTimes == 0 {
		now := time.Now().Unix()
		for k, v := range bl.keys {
			if now > v {
				delete(bl.keys, k)
			}
		}
	}

	bl.keys[string(id)] = time.Now().Add(expires).Unix()
	return nil
}

// Len returns the number of entries in the blacklist
func (bl *InMemory) Len() int {
	bl.mu.RLock()
	l := len(bl.keys)
	bl.mu.RUnlock()
	return l
}
