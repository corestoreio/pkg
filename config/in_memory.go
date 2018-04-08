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

package config

import (
	"strings"
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
)

const kvMapScopeSep = '~'

type kvmap struct {
	sync.RWMutex
	kv map[string][]byte
}

// NewInMemoryStore creates a new simple key value storage using a map[string]interface{}.
// Mainly used for testing.
func NewInMemoryStore() Storager {
	return &kvmap{
		kv: make(map[string][]byte),
	}
}

// Set implements Storager interface
func (sp *kvmap) Set(scp scope.TypeID, path string, value []byte) error {

	var key strings.Builder
	key.WriteString(scp.ToIntString())
	key.WriteByte(kvMapScopeSep)
	key.WriteString(path)

	sp.Lock()
	sp.kv[key.String()] = value
	sp.Unlock()
	return nil
}

// Get implements Storager interface.
func (sp *kvmap) Value(scp scope.TypeID, path string) (v []byte, ok bool, err error) {

	var key strings.Builder
	key.WriteString(scp.ToIntString())
	key.WriteByte(kvMapScopeSep)
	key.WriteString(path)

	sp.RLock()
	data, ok := sp.kv[key.String()]
	sp.RUnlock()
	if ok {
		return data, ok, nil
	}
	return nil, false, nil
}

// AllKeys implements Storager interface and return unsorted slices. Despite
// having two slices and they are unsorted the index of the TypeIDs slice does
// still refer to the same index of the paths slice.
func (sp *kvmap) AllKeys() (scps scope.TypeIDs, paths []string, err error) {
	sp.RLock()

	scps = make(scope.TypeIDs, len(sp.kv))
	paths = make([]string, len(sp.kv))
	i := 0
	for key := range sp.kv {
		idx := strings.IndexByte(key, kvMapScopeSep)
		scps[i], err = scope.MakeTypeIDString(key[:idx])
		if err != nil {
			return nil, nil, errors.Wrapf(err, "[config] InMemory Storage with key %q", key)
		}
		paths[i] = key[idx+1:]
		i++
	}
	sp.RUnlock()
	return
}
