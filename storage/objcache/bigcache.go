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

// +build bigcache csall

package objcache

import (
	"context"
	"time"

	"github.com/allegro/bigcache"
	"github.com/corestoreio/errors"
)

var errKeyNotFound = errors.NotFound.Newf(`[objcache] Key not found`)

// WithBigCache sets the bigcache as underlying storage engine to the
// This function allows to set custom configuration options to the bigcache
// instance.
// Default option: shards 256, LifeWindow 12 hours, Verbose false
//
// For more details: https://godoc.org/github.com/allegro/bigcache
func WithBigCache(c bigcache.Config) Option {
	def := bigcache.Config{
		// optimize this ...
		Shards:             256,
		LifeWindow:         time.Hour * 12,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            false,
		HardMaxCacheSize:   0,
	}
	if c.Shards > 0 {
		def = c
	}
	return Option{
		fn: func(p *Manager) error {
			c, err := bigcache.NewBigCache(def)
			if err != nil {
				return errors.WithStack(err)
			}
			p.cache = bigCacheWrapper{c}
			return nil
		},
	}
}

type bigCacheWrapper struct {
	*bigcache.BigCache
}

func (w bigCacheWrapper) Set(_ context.Context, key string, value []byte) error {
	if err := w.BigCache.Set(key, value); err != nil {
		// This error construct save some unneeded allocations.
		return errors.Wrapf(err, "[objcache] With key %q", key)
	}
	return nil
}

func (w bigCacheWrapper) Get(_ context.Context, key string) ([]byte, error) {
	v, err := w.BigCache.Get(key)
	if _, ok := err.(*bigcache.EntryNotFoundError); ok {
		return nil, errKeyNotFound
	}
	if err != nil {
		return nil, errors.Fatal.New(err, "[objcache] With key %q", key)
	}
	return v, nil
}

func (w bigCacheWrapper) Delete(_ context.Context, key string) (err error) {
	return w.BigCache.Delete(key)
}

func (w bigCacheWrapper) Close() error {
	return nil
}
