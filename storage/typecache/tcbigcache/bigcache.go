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

package tcbigcache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/corestoreio/csfw/storage/typecache"
	"github.com/corestoreio/csfw/util/errors"
)

var errKeyNotFound = errors.NewNotFoundf(`[tcbigcache] Key not found`)

// With sets the bigcache as underlying storage engine to the typecache.
// This function allows to set custom configuration options to the bigcache
// instance.
// Default option: shards 256, LifeWindow 12 hours, Verbose false
//
// For more details: https://godoc.org/github.com/allegro/bigcache
func With(c ...bigcache.Config) typecache.Option {
	def := bigcache.Config{
		// optimize this ...
		Shards:             256,
		LifeWindow:         time.Hour * 12,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            false,
		HardMaxCacheSize:   0,
	}
	if len(c) == 1 {
		def = c[0]
	}
	return func(p *typecache.Processor) error {
		c, err := bigcache.NewBigCache(def)
		p.Cache = wrapper{c}
		return errors.NewFatal(err, "[tcbigcache] bigcache.NewBigCache")
	}
}

type wrapper struct {
	*bigcache.BigCache
}

func (w wrapper) Set(key []byte, value []byte) error {
	return errors.Wrap(
		w.BigCache.Set(string(key), value),
		"[tcbigcache] wrapper.Set.Set")
}

func (w wrapper) Get(key []byte) ([]byte, error) {
	v, err := w.BigCache.Get(string(key))
	if _, ok := err.(*bigcache.EntryNotFoundError); ok {
		return nil, errKeyNotFound
	}
	if err != nil {
		return nil, errors.NewFatal(err, "[tcbigcache] wrapper.Get.Get")
	}
	return v, nil
	// just to sure to copy the data away
	//buf := make([]byte, len(v), len(v))
	//copy(buf, v)
	//return buf, nil
}

func (bw wrapper) Close() error {
	return nil
}
