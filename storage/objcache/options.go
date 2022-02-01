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
	"time"

	"golang.org/x/exp/maps"
)

var now = time.Now

// ServiceOptions used when creating a NewService.
type ServiceOptions struct {
	// Codec optionally encodes and decodes an object. The default "codec"
	// checks if a type implements Marshal() or Unmarshal() interface or any of
	// the two interfaces from package "encoding".
	Codec Codecer
	// PrimeObjects creates new encoder/decoder with a sync.Pool to reuse the
	// objects. Setting PrimeObjects causes the encoder/decode to prime the data
	// which means that no type information will be stored in the cache. If you
	// use gob you must use still gob.Register() your types. TL;DR: Skips type
	// information in the cache.
	PrimeObjects   []any
	DefaultExpires time.Duration
}

// NewCacheSimpleInmemory creates an in-memory map map[string]string as cache
// backend which supports expiration.
func NewCacheSimpleInmemory[K comparable]() (Storager[K], error) {
	mc := &mapCache[K]{
		items: make(map[K]mapCacheItem),
	}
	return mc, nil
}

type mapCacheItem struct {
	value      []byte
	expiration time.Time
}

type mapCache[K comparable] struct {
	rwmu  sync.RWMutex
	items map[K]mapCacheItem
}

func (mc *mapCache[K]) Set(_ context.Context, keys []K, values [][]byte, expirations []time.Duration) (err error) {
	hasExp := len(expirations) > 0
	n := now()
	for i, key := range keys {
		var e time.Time
		if hasExp {
			if ed := expirations[i]; ed > 0 {
				e = n.Add(ed)
			}
		}
		mc.items[key] = mapCacheItem{value: values[i], expiration: e}
	}
	return nil
}

func (mc *mapCache[K]) Get(_ context.Context, keys []K) (values [][]byte, err error) {
	n := now()
	for _, key := range keys {
		v, ok := mc.items[key]
		if ok && (v.expiration.IsZero() || v.expiration.After(n)) {
			values = append(values, v.value)
		} else {
			values = append(values, nil)
		}
	}
	return values, nil
}

func (mc *mapCache[K]) Delete(_ context.Context, keys []K) error {
	for _, key := range keys {
		delete(mc.items, key)
	}
	return nil
}

func (mc *mapCache[K]) Truncate(ctx context.Context) (err error) {
	mc.rwmu.Lock()
	defer mc.rwmu.Unlock()
	maps.DeleteFunc(mc.items, func(k K, v mapCacheItem) bool {
		v.value = nil
		return true
	})
	return nil
}
func (mc *mapCache[K]) Close() error { return nil }

// NewBlackHoleClient creates a black hole client for testing with the ability
// to return errors.
func NewBlackHoleClient[K comparable](optionalTestErr error) NewStorageFn[K] {
	return func() (Storager[K], error) { return blackHole[K]{err: optionalTestErr}, nil }
}

type blackHole[K comparable] struct {
	err error
}

func (mc blackHole[K]) Set(_ context.Context, keys []K, values [][]byte, expirations []time.Duration) (err error) {
	return mc.err
}

func (mc blackHole[K]) Get(_ context.Context, keys []K) (values [][]byte, err error) {
	return nil, mc.err
}

func (mc blackHole[K]) Delete(_ context.Context, keys []K) (err error) { return mc.err }
func (mc blackHole[K]) Truncate(ctx context.Context) (err error)       { return mc.err }
func (mc blackHole[K]) Close() error                                   { return mc.err }

// binary a simple type to use the Service as a set-algorithm to e.g. check if a
// key exists.
type binary byte

func (n *binary) Unmarshal(data []byte) error {
	var val binary = '0'
	if len(data) == 1 && data[0] == '1' {
		val = '1'
	}
	*n = val
	return nil
}

func (n binary) Marshal() ([]byte, error) { return []byte{byte(n)}, nil }
func (n binary) IsValid() bool            { return n == '1' }

// MakeBinary creates a binary type for using in Set/Get functions when the code
// needs a `set` algorithm. For example checking if a JWT exists in the
// blockList. Function IsValid returns true if the key exists.
func MakeBinary() binary { return binary('1') }
