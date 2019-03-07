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

// +build csall bigcache

package storage

import (
	"github.com/allegro/bigcache"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
)

type bcStorage struct {
	bc *bigcache.BigCache
}

// NewBigCache creates a new cache with a minimum size set to 512KB. If the size
// is set relatively large, you should call `debug.SetGCPercent()`, set it to a
// much smaller value to limit the memory consumption and GC pause time.
//
// Bigcache delivers under high concurrent and parallel load better results than
// a simple key value mutex protected map.
//
// Maybe implements synchronization with MySQL core_config_data table. Converts
// all values to byte slices.
func NewBigCache(config bigcache.Config) (config.Storager, error) {
	bc, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, errors.Fatal.New(err, "[config/storage] NewBigCache")
	}
	return &bcStorage{
		bc: bc,
	}, nil
}

// Set writes a key with its value into the storage. The value
// gets converted to a byte slice.
func (s *bcStorage) Set(p *config.Path, value []byte) error {
	return s.bc.Set(p.String(), value)
}

// Get returns a value from the cache.
func (s *bcStorage) Get(p *config.Path) (v []byte, found bool, err error) {

	val, err := s.bc.Get(p.String())
	isNotFound := err == bigcache.ErrEntryNotFound
	if err != nil && !isNotFound {
		return nil, false, err
	}
	if isNotFound {
		return nil, false, nil
	}

	return val, true, nil
}
