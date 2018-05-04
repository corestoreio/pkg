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

package cfgbigcache

import (
	"strings"

	"github.com/allegro/bigcache"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
)

const kvMapScopeSep = '~'

// Storage wrapper around the freecache.Cache type
type Storage struct {
	Cache *bigcache.BigCache
}

// New creates a new cache with a minimum size set to 512KB.
// If the size is set relatively large, you should call
// `debug.SetGCPercent()`, set it to a much smaller value
// to limit the memory consumption and GC pause time.
func New(config bigcache.Config) (*Storage, error) {
	bc, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, errors.Fatal.New(err, "[bigcache] NewBigCache")
	}

	return &Storage{
		Cache: bc,
	}, nil
}

// Set writes a key with its value into the storage. The value
// gets converted to a byte slice.
func (s *Storage) Set(scp scope.TypeID, path string, value []byte) error {
	var key strings.Builder
	key.WriteString(scp.ToIntString())
	key.WriteByte(kvMapScopeSep)
	key.WriteString(path)
	return s.Cache.Set(key.String(), value)
}

// Get may return a ErrKeyNotFound error
func (s *Storage) Value(scp scope.TypeID, path string) (v []byte, found bool, err error) {
	var key strings.Builder
	key.WriteString(scp.ToIntString())
	key.WriteByte(kvMapScopeSep)
	key.WriteString(path)

	val, err := s.Cache.Get(key.String())
	_, isNotFound := (err).(*bigcache.EntryNotFoundError)
	if err != nil && !isNotFound {
		return nil, false, err
	}
	if isNotFound {
		return nil, false, nil
	}

	return val, true, nil
}
