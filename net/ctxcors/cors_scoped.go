// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ctxcors

import (
	"sync"

	"github.com/corestoreio/csfw/config"
)

const (
	PathCorsExposedHeaders = "web/cors/exposed_headers"
	PathCorsAllowedOrigins = "web/cors/allowed_origins"
)

// corsScopeCache
type corsScopeCache struct {
	config config.Getter

	// rwmu protects the map
	rwmu sync.RWMutex
	// storage key is the website ID and value the current cors config
	storage map[int64]*Cors
}

// get uses a read lock to check if a Cors exists for a website ID. returns nil
// if there is no Cors pointer. aim is: multiple goroutines can read from the
// map while adding new Cors pointers can only be done by one goroutine.
func (cs *corsScopeCache) get(websiteID int64) *Cors {
	cs.rwmu.RLock()
	defer cs.rwmu.RUnlock()
	if c, ok := cs.storage[websiteID]; ok {
		return c
	}
	return nil
}

// create creates a new Cors type and returns it.
func (cs *corsScopeCache) insert(websiteID int64) *Cors {
	cs.rwmu.Lock()
	defer cs.rwmu.Unlock()

	// pulls the options from the scoped reader

	c := New()
	cs.storage[websiteID] = c
	return c
}
