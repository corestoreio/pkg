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
	"github.com/juju/errors"
)

// scopeCache creates a new Cors type for a website configuration. Why are
// we not using the store view based configuration? because that is too low level
// and store views are mainly used for languages. but we can change that or
// make it configurable.
type scopeCache struct {
	parent *Cors

	// rwmu protects the map
	rwmu sync.RWMutex
	// storage key is the ID and value the current cors config
	storage map[scope.Hash]*Cors
	// under very very high load this map will become a bottle neck so we
	// should switch to a lock free data structure.
}

func newScopeCache(parent *Cors) *scopeCache {
	return &scopeCache{
		parent:  parent,
		storage: make(map[scope.Hash]*Cors),
	}
}

// get uses a read lock to check if a Cors exists for an ID. returns nil
// if there is no Cors pointer. aim is: multiple goroutines can read from the
// map while adding new Cors pointers can only be done by one goroutine.
func (cs *scopeCache) get(s scope.Scope, id int64) *Cors {
	cs.rwmu.RLock()
	defer cs.rwmu.RUnlock()
	if c, ok := cs.storage[scope.NewHash(s, id)]; ok {
		return c
	}
	return nil
}

// create creates a new Cors type for a scope and returns it.
func (cs *scopeCache) insert(sg config.ScopedGetter) (*Cors, error) {
	cs.rwmu.Lock()
	defer cs.rwmu.Unlock()

	var c *Cors
	var err error
	if cs.parent.Backend != nil {
		c, err = New(WithBackendApplied(cs.parent.Backend, sg))
	} else {
		c, err = New(WithLogger(cs.parent.Log)) // inherit more?
	}
	if err != nil {
		return nil, errors.Mask(err)
	}

	c.scopedTo = scope.NewHash(sg.Scope())

	cs.storage[c.scopedTo] = c
	return c, nil
}
