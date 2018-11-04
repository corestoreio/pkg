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
	PrimeObjects   []interface{}
	DefaultExpires time.Duration
}

// NewCacheSimpleInmemory creates an in-memory map map[string]string as cache
// backend which supports expiration.
func NewCacheSimpleInmemory() (Storager, error) {
	mc := &mapCache{}
	return mc, nil
}

type mapCacheItem struct {
	value      string
	expiration time.Time
}

type mapCache struct {
	items sync.Map
}

func (mc *mapCache) Put(_ context.Context, keys []string, values [][]byte, expirations []time.Duration) (err error) {
	hasExp := len(expirations) > 0
	n := now()
	for i, key := range keys {
		var e time.Time
		if hasExp {
			if ed := expirations[i]; ed > 0 {
				e = n.Add(ed)
			}
		}
		mc.items.Store(key, &mapCacheItem{value: string(values[i]), expiration: e})
	}
	return nil
}

func (mc *mapCache) Get(_ context.Context, keys []string) (values [][]byte, err error) {
	n := now()
	for _, key := range keys {
		val, ok := mc.items.Load(key)
		if v, ok2 := val.(*mapCacheItem); ok2 && ok && (v.expiration.IsZero() || v.expiration.After(n)) {
			values = append(values, []byte(v.value))
		} else {
			values = append(values, nil)
		}
	}
	return values, nil
}

func (mc *mapCache) Delete(_ context.Context, keys []string) (err error) {
	for _, key := range keys {
		mc.items.Delete(key)
	}
	return nil
}
func (mc *mapCache) Truncate(ctx context.Context) (err error) {
	mc.items.Range(func(key, value interface{}) bool {
		value = nil
		mc.items.Delete(key)
		return true
	})
	mc.items = sync.Map{}
	return nil
}
func (mc *mapCache) Close() error { return nil }

// NewBlackHoleClient creates a black hole client for testing with the ability
// to return errors.
func NewBlackHoleClient(optionalTestErr error) NewStorageFn {
	return func() (Storager, error) { return blackHole{err: optionalTestErr}, nil }
}

type blackHole struct {
	err error
}

func (mc blackHole) Put(_ context.Context, keys []string, values [][]byte, expirations []time.Duration) (err error) {
	return mc.err
}

func (mc blackHole) Get(_ context.Context, keys []string) (values [][]byte, err error) {
	return nil, mc.err
}

func (mc blackHole) Delete(_ context.Context, keys []string) (err error) { return mc.err }
func (mc blackHole) Truncate(ctx context.Context) (err error)            { return mc.err }
func (mc blackHole) Close() error                                        { return mc.err }

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

// MakeBinary creates a binary type for using in Put/Get functions when the code
// needs a `set` algorithm. For example checking if a JWT exists in the
// blacklist. Function IsValid returns true if the key exists.
func MakeBinary() binary { return binary('1') }
