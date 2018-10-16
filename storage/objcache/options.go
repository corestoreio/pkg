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

package objcache

import (
	"context"
	"sync"
)

// Option provides convenience helper functions to apply various options while
// creating a new Manager type.
type Option struct {
	sortOrder int
	fn        func(*Manager) error
}

type options []Option

func (o options) Len() int           { return len(o) }
func (o options) Less(i, j int) bool { return o[i].sortOrder < o[j].sortOrder }
func (o options) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

// WithEncoder sets a custom encoder and decoder.
func WithEncoder(codec Codecer) Option {
	return Option{
		fn: func(p *Manager) error {
			p.codec = codec
			return nil
		},
	}
}

// WithPooledEncoder creates new encoder/decoder with a sync.Pool to reuse the
// objects. Providing argument primeObjects causes the encoder/decode to prime
// the data which means that no type information will be stored in the cache.
// If you use gob you must use gob.Register() for your types.
func WithPooledEncoder(codec Codecer, primeObjects ...interface{}) Option {
	return Option{
		fn: func(p *Manager) error {
			p.codec = newPooledCodec(codec, primeObjects...)
			return nil
		},
	}
}

// WithCache sets a custom cache type. Examples in the subpackages.
func WithCache(c Storager) Option {
	return Option{
		fn: func(p *Manager) error {
			p.cache = c
			return nil
		},
	}
}

// WithSimpleSlowCacheMap creates an in-memory map map[string]string as cache
// backend.
func WithSimpleSlowCacheMap() Option {
	return Option{
		fn: func(p *Manager) error {
			p.cache = &mapCache{
				items: make(map[string]string),
			}
			return nil
		},
	}
}

type mapCache struct {
	sync.RWMutex
	items map[string]string
}

func (mc *mapCache) Set(_ context.Context, key string, value []byte) (err error) {
	mc.Lock()
	defer mc.Unlock()
	mc.items[string(key)] = string(value)
	return nil
}

func (mc *mapCache) Get(_ context.Context, key string) (value []byte, err error) {
	mc.RLock()
	defer mc.RUnlock()
	if v, ok := mc.items[string(key)]; ok {
		return []byte(v), nil
	}
	return nil, nil
}

func (mc *mapCache) Delete(_ context.Context, key string) (err error) {
	mc.Lock()
	defer mc.Unlock()
	delete(mc.items, key)
	return nil
}
func (mc *mapCache) Close() error { return nil }

// OpOption configures Operations like Get, Set, Delete.
type OpOption struct{}
