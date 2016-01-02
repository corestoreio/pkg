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

package ctxcors

import (
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// corsScopeCache creates a new Cors type for a website configuration. Why are
// we not using the store view based configuration? because that is too low level
// and store views are mainly used for languages. but we can change that or
// make it configurable.
type corsScopeCache struct {
	config config.Getter
	scope  scope.Scope
	parent *Cors

	// rwmu protects the map
	rwmu sync.RWMutex
	// storage key is the ID and value the current cors config
	storage map[int64]*Cors
}

func newCorsScopeCache(cg config.Getter, s scope.Scope, parent *Cors) *corsScopeCache {
	return &corsScopeCache{
		config:  cg,
		scope:   s,
		parent:  parent,
		storage: make(map[int64]*Cors),
	}
}

// get uses a read lock to check if a Cors exists for an ID. returns nil
// if there is no Cors pointer. aim is: multiple goroutines can read from the
// map while adding new Cors pointers can only be done by one goroutine.
func (cs *corsScopeCache) get(id int64) *Cors {
	cs.rwmu.RLock()
	defer cs.rwmu.RUnlock()
	if c, ok := cs.storage[id]; ok {
		return c
	}
	return nil
}

// create creates a new Cors type and returns it.
func (cs *corsScopeCache) insert(id int64) *Cors {
	cs.rwmu.Lock()
	defer cs.rwmu.Unlock()

	// pulls the options from the scoped reader
	//	headers, _ := cs.config.String(config.Path(PathCorsExposedHeaders), config.Scope(cs.scope, id))
	//	fields, err := PackageConfiguration.FindFieldByPath(PathCorsExposedHeaders)

	c := New(WithLogger(cs.parent.Log)) // inherit more?
	cs.storage[id] = c
	return c
}
