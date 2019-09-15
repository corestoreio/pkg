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

package storage

import (
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
)

type kvmap struct {
	sync.RWMutex
	kv map[cacheKey]string
}

// NewMap creates a new simple key value storage using a map. Mainly used for
// testing. fqPathValue must be a balanced slice where i=fully qualified path
// and i+1 the value.
func NewMap(fqPathValue ...string) config.Storager {
	m := &kvmap{
		kv: make(map[cacheKey]string),
	}
	if len(fqPathValue) == 0 {
		return m
	}

	p := new(config.Path)
	for i := 0; i < len(fqPathValue); i += 2 {
		fq := fqPathValue[i]
		if err := p.Parse(fq); err != nil {
			panic(errors.Fatal.New(err, "[config/storage] NewMap with path %q", fq))
		}
		m.kv[makeCacheKey(p.ScopeRoute())] = fqPathValue[i+1]
		p.Reset()
	}

	return m
}

// Set implements Storager interface
func (sp *kvmap) Set(p config.Path, value []byte) error {
	sp.Lock()
	sp.kv[makeCacheKey(p.ScopeRoute())] = string(value)
	sp.Unlock()
	return nil
}

// Get implements Storager interface and returns a byte slice available for
// modifications.
func (sp *kvmap) Get(p config.Path) (v []byte, found bool, err error) {
	sp.RLock()
	defer sp.RUnlock()
	data, found := sp.kv[makeCacheKey(p.ScopeRoute())]
	if found {
		return []byte(data), found, nil
	}
	return nil, false, nil
}

// Flush purges all stored items from the cache.
func (sp *kvmap) Flush() error {
	sp.Lock()
	prevLen := len(sp.kv)
	sp.kv = make(map[cacheKey]string, prevLen)
	sp.Unlock()
	return nil
}

// Keys returns a randomized slice of all keys. Useful while testing.
func (sp *kvmap) Keys(ret ...string) []string {
	sp.RLock()
	defer sp.RUnlock()
	if lkv := len(sp.kv); cap(ret) == 0 {
		ret = make([]string, 0, lkv)
	}
	for k := range sp.kv {
		ret = append(ret, k.scp.String()+string(config.PathSeparator)+k.route)
	}
	return ret
}
