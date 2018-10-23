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
	"fmt"
	"sync"

	"github.com/corestoreio/errors"
)

// Option provides convenience helper functions to apply various options while
// creating a new Service type.
type Option struct {
	sortOrder int
	fn        func(*Service) error
}

type options []Option

func (o options) Len() int           { return len(o) }
func (o options) Less(i, j int) bool { return o[i].sortOrder < o[j].sortOrder }
func (o options) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

// WithEncoder sets a custom encoder and decoder.
func WithEncoder(codec Codecer) Option {
	return Option{
		fn: func(p *Service) error {
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
		fn: func(p *Service) error {
			p.codec = newPooledCodec(codec, primeObjects...)
			return nil
		},
	}
}

// WithAddStorage sets a custom cache type. Examples in the subpackages.
func WithAddStorage(c Storager) Option {
	return Option{
		fn: func(p *Service) error {
			p.cache[len(p.cache)+1] = c
			return nil
		},
	}
}

// WithSimpleSlowCacheMap creates an in-memory map map[string]string as cache
// backend.
func WithSimpleSlowCacheMap() Option {
	return Option{
		fn: func(p *Service) error {
			p.cache[len(p.cache)+1] = &mapCache{
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

func (mc *mapCache) Set(_ context.Context, items *Items) (err error) {
	mc.Lock()
	defer mc.Unlock()
	keys, values, err := items.Encode(nil, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	for i, key := range keys {
		mc.items[key] = string(values[i])
	}
	return nil
}

func (mc *mapCache) Get(_ context.Context, keys []string) (values [][]byte, err error) {
	mc.RLock()
	defer mc.RUnlock()
	for _, key := range keys {
		if v, ok := mc.items[key]; ok {
			values = append(values, []byte(v))
		} else {
			return nil, ErrKeyNotFound(key)
		}
	}
	return values, nil
}

func (mc *mapCache) Delete(_ context.Context, keys []string) (err error) {
	mc.Lock()
	defer mc.Unlock()
	for _, key := range keys {
		delete(mc.items, key)
	}
	return nil
}
func (mc *mapCache) Close() error { return nil }

// ErrKeyNotFound returned by a backend cache to indicate that a key can't be
// found.
type ErrKeyNotFound string

func (e ErrKeyNotFound) ErrorKind() errors.Kind {
	return errors.NotFound
}
func (e ErrKeyNotFound) Error() string {
	return fmt.Sprintf("[objcache] Key %q not found", e)
}
