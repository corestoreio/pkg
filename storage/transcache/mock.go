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

package transcache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sync"
	"sync/atomic"

	"github.com/corestoreio/errors"
)

// Mock used mainly for testing. Fully concurrent safe. Does not implement
// Encoder and Decoder and stores the values in a map[string][]byte with values
// encoded via gob.
type Mock struct {
	mu       sync.RWMutex
	cache    map[string][]byte
	SetErr   error
	GetErr   error
	setCount *int32
	getCount *int32
}

// NewMock creates a new Transcacher compatible mock and initializes the
// underlying map and synchronization types.
func NewMock() *Mock {
	return &Mock{
		cache:    make(map[string][]byte),
		setCount: new(int32),
		getCount: new(int32),
	}
}

// Set writes a src into the cache or returns the error defined in the field
// SetErr.
func (mc *Mock) Set(key []byte, src interface{}) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	if mc.SetErr != nil {
		return mc.SetErr
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return errors.NewFatal(err, "[transcache] Set.Gob.Encode")
	}
	mc.cache[string(key)] = buf.Bytes()

	atomic.AddInt32(mc.setCount, 1)
	return nil
}

// Get looks up the key in the cache and parses the value into dst (destination)
// or returns an error as defined in GetErr. If the key cannot be found an error
// behaviour of NotFound will get returned. Dst must be a pointer.
func (mc *Mock) Get(key []byte, dst interface{}) error {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	if mc.GetErr != nil {
		return mc.GetErr
	}

	if raw, ok := mc.cache[string(key)]; ok {
		if err := gob.NewDecoder(bytes.NewReader(raw)).Decode(dst); err != nil {
			return errors.NewFatal(json.Unmarshal(raw, dst), "[transcache] Get.Gob.Decode")
		}
		atomic.AddInt32(mc.getCount, 1)
		return nil
	}
	return errors.NewNotFoundf("[transcache] Key %q not found", string(key))
}

// SetCount returns the cache set count
func (mc *Mock) SetCount() int {
	return int(atomic.LoadInt32(mc.setCount))
}

// GetCount returns the cache hits.
func (mc *Mock) GetCount() int {
	return int(atomic.LoadInt32(mc.getCount))
}
