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

package ctxjwt

import (
	"sync"
	"time"

	"github.com/pborman/uuid"
)

// nullBL is the black hole black list
type nullBL struct{}

func (b nullBL) Set(_ string, _ time.Duration) error { return nil }
func (b nullBL) Has(_ string) bool                   { return false }

// SimpleMapBlackList creates an in-memory map which holds as a key the
// tokens and as value the token expiration duration. Once a Set() operation
// will be called the tokens list get purged. Don't use this feature in
// production as the underlying mutex will become a bottleneck with higher
// throughput.
type SimpleMapBlackList struct {
	// TODO: implement github.com/coocood/freecache
	mu     sync.RWMutex
	tokens map[string]time.Time
}

// NewSimpleMapBlackList creates a new blacklist map.
func NewSimpleMapBlackList() *SimpleMapBlackList {
	return &SimpleMapBlackList{
		tokens: make(map[string]time.Time),
	}
}

// Has checks if token is within the blacklist.
func (bl *SimpleMapBlackList) Has(token string) bool {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	d, ok := bl.tokens[token]
	if !ok {
		return false
	}
	isValid := time.Since(d) < 0

	if false == isValid {
		delete(bl.tokens, token)
	}
	return isValid
}

// Set adds a token to the map and performs a purge operation.
func (bl *SimpleMapBlackList) Set(token string, expires time.Duration) error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	for k, v := range bl.tokens {
		if time.Since(v) > 0 {
			delete(bl.tokens, k)
		}
	}
	bl.tokens[token] = time.Now().Add(expires)
	return nil
}

// Len returns the number of entries in the blacklist
func (bl *SimpleMapBlackList) Len() int {
	bl.mu.RLock()
	defer bl.mu.RUnlock()
	return len(bl.tokens)
}

// jti type to generate a JTI for a token, a unique ID
type jti struct{}

func (j jti) Get() string {
	return uuid.New()
}
