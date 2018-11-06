/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package objcache

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type liElem struct {
	key  string
	bVal []byte
}

// lruCache is an LRU cache. It is safe for concurrent access.
// This type does not get exported.
type lruCache struct {
	// maxEntries is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	maxEntries int

	mu    sync.Mutex
	ll    *list.List
	cache map[string]*list.Element
}

// NewCacheLRU creates a new lruCache. If maxEntries is zero, the cache has no limit
// and it's assumed that eviction is done by the caller. This type does not get
// exported.
// WithLRU provides the `lru` cache which implements a fixed-size thread safe
// LRU cache. If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
// lru cache to limit the amount of request to the backend service. Use
// function WithLRU to enable it and set the correct max size of the LRU
// cache. For now this algorithm should be good enough. Can be refactored
// any time later.
// Can be replaced by https://github.com/goburrow/cache TinyLFU
// LRU cache backend does not support expiration.
func NewCacheLRU(maxEntries int) NewStorageFn {
	return func() (Storager, error) {
		if maxEntries == 0 {
			maxEntries = 1024
		}
		return &lruCache{
			maxEntries: maxEntries,
			ll:         list.New(),
			cache:      make(map[string]*list.Element, maxEntries+1),
		}, nil
	}
}

// Add adds a value to the cache. Panics on nil Path.
func (c *lruCache) Set(_ context.Context, keys []string, values [][]byte, _ []time.Duration) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, key := range keys {
		if ee, ok := c.cache[key]; ok {
			ee.Value = liElem{key: key, bVal: values[i]}
			c.ll.MoveToFront(ee)
			return nil
		}
		ele := c.ll.PushFront(liElem{key: key, bVal: values[i]})
		c.cache[key] = ele
		if c.maxEntries > 0 && c.ll.Len() > c.maxEntries {
			c.removeOldest()
		}
	}
	return nil
}

// Get looks up a key's value from the cache.
func (c *lruCache) Get(_ context.Context, keys []string) (values [][]byte, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		if ele, ok := c.cache[key]; ok {
			c.ll.MoveToFront(ele)
			values = append(values, ele.Value.(liElem).bVal)
		} else {
			values = append(values, nil)
		}
	}
	return
}

func (c *lruCache) removeOldest() {
	ele := c.ll.Back()
	if ele == nil {
		return
	}
	c.removeElement(ele)
}

func (c *lruCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	le := e.Value.(liElem)
	delete(c.cache, le.key)
}

// Flush purges all stored items from the cache.
func (c *lruCache) Truncate(_ context.Context) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ll = list.New()
	me := c.maxEntries
	if me == 0 {
		me = 1024
	}
	c.cache = make(map[string]*list.Element, me+1)
	return nil
}

func (c *lruCache) Delete(_ context.Context, keys []string) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		if ee, ok := c.cache[key]; ok {
			c.ll.Remove(ee)
			c.cache[key] = nil
			delete(c.cache, key)
		}
	}
	return nil
}

func (c *lruCache) Close() error {
	return c.Truncate(context.TODO())
}
