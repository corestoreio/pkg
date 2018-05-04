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
	"sync"

	"github.com/corestoreio/pkg/store/scope"
)

type kvMapKey struct {
	scope.TypeID
	string
}

type kvmap struct {
	sync.RWMutex
	kv map[kvMapKey][]byte
}

// NewInMemoryStore creates a new simple key value storage using a map[string]interface{}.
// Mainly used for testing.
func NewInMemoryStore() Storager {
	return &kvmap{
		kv: make(map[kvMapKey][]byte),
	}
}

// Set implements Storager interface
func (sp *kvmap) Set(scp scope.TypeID, path string, value []byte) error {
	sp.Lock()
	sp.kv[kvMapKey{scp, path}] = value
	sp.Unlock()
	return nil
}

// Get implements Storager interface.
func (sp *kvmap) Value(scp scope.TypeID, path string) (v []byte, found bool, err error) {
	sp.RLock()
	defer sp.RUnlock()
	data, found := sp.kv[kvMapKey{scp, path}]
	return data, found, nil
}

// AllKeys implements Storager interface and return unsorted slices. Despite
// having two slices and they are unsorted the index of the TypeIDs slice does
// still refer to the same index of the paths slice.
func (sp *kvmap) AllKeys() (scps scope.TypeIDs, paths []string, err error) {
	sp.RLock()
	defer sp.RUnlock()

	scps = make(scope.TypeIDs, len(sp.kv))
	paths = make([]string, len(sp.kv))
	i := 0
	for key := range sp.kv {
		scps[i] = key.TypeID
		paths[i] = key.string
		i++
	}
	return
}
