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

package storage

import (
	"sync"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/util/errors"
)

// Storager is the underlying data storage for holding the keys and its values.
// Implementations can be spf13/viper or MySQL backed. Default Storager is a
// simple mutex protected map[string]interface{}. ProTip: If you use MySQL as
// Storager don't execute function ApplyCoreConfigData() The config.Writer calls
// the config.Storager functions and config.Storager must make sure of the
// correct type conversions to the supported type of the underlying storage
// engine.
type Storager interface {
	// Set sets a key with a value and returns on success nil or ErrKeyOverwritten,
	// on failure any other error
	Set(key cfgpath.Path, value interface{}) error
	// Get returns the raw value on success or may return a NotFound error
	// behaviour if an entry cannot be found or does not exists. Any other error
	// can also occur.
	Get(key cfgpath.Path) (interface{}, error)
	// AllKeys returns the fully qualified keys
	AllKeys() (cfgpath.PathSlice, error)
}

// NotFound error type which defines that a specific key cannot be found.
type NotFound struct{}

func (NotFound) Error() string  { return "[storage] Key not found" }
func (NotFound) NotFound() bool { return true } // @see errors.IsNotFound

type keyVal struct {
	k cfgpath.Path
	v interface{}
}

type kvmap struct {
	sync.Mutex
	kv map[uint32]keyVal // todo: create benchmark to check if worth switching to pointers
}

// NewKV creates a new simple key value storage using a map[string]interface{}
// without any persistence or sync to MySQL core_confing_data table
func NewKV() *kvmap {
	return &kvmap{
		kv: make(map[uint32]keyVal),
	}
}

// Set implements Storager interface
func (sp *kvmap) Set(key cfgpath.Path, value interface{}) error {
	sp.Lock()
	defer sp.Unlock()

	h32, err := key.Hash(-1)
	if err != nil {
		return errors.Wrap(err, "[storage] key.Hash")
	}
	sp.kv[h32] = keyVal{key, value}
	return nil
}

// Get implements Storager interface.
// Error behaviour: NotFound.
func (sp *kvmap) Get(key cfgpath.Path) (interface{}, error) {
	sp.Lock()
	defer sp.Unlock()

	h32, err := key.Hash(-1)
	if err != nil {
		return nil, errors.Wrap(err, "[storage] key.Hash")
	}
	if data, ok := sp.kv[h32]; ok {
		return data.v, nil
	}
	return nil, NotFound{}
}

// AllKeys implements Storager interface
func (sp *kvmap) AllKeys() (cfgpath.PathSlice, error) {
	sp.Lock()
	defer sp.Unlock()

	var ret = make(cfgpath.PathSlice, len(sp.kv), len(sp.kv))
	i := 0
	for _, kv := range sp.kv {
		ret[i] = kv.k
		i++
	}
	return ret, nil
}
