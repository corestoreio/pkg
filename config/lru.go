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

package config

import (
	"container/list"
	"sync"

	"github.com/corestoreio/errors"
)

// lruCache is an LRU cache. It is safe for concurrent access.
// This type does not get exported.
type lruCache struct {
	// maxEntries is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	maxEntries int

	mu    sync.Mutex
	ll    *list.List
	cache map[Path]*list.Element
}

// newLRU creates a new lruCache. If maxEntries is zero, the cache has no limit
// and it's assumed that eviction is done by the caller. This type does not get
// exported.
func newLRU(maxEntries int) *lruCache {
	if maxEntries == 0 {
		maxEntries = 1024
	}
	return &lruCache{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[Path]*list.Element, maxEntries+1),
	}
}

// Add adds a value to the cache. Panics on nil Path.
func (c *lruCache) Add(key Path, value Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !key.Equal(&value.Path) {
		panic(errors.Mismatch.Newf("[config] Key Path %q does not match Value path %q", key.String(), value.Path.String()))
	}

	if ee, ok := c.cache[key]; ok {
		ee.Value = value
		c.ll.MoveToFront(ee)
		return
	}
	ele := c.ll.PushFront(value)
	c.cache[key] = ele
	if c.maxEntries > 0 && c.ll.Len() > c.maxEntries {
		c.removeOldest()
	}
}

// Get looks up a key's value from the cache.
func (c *lruCache) Get(key Path) (value Value, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(Value), true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *lruCache) Remove(key Path) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *lruCache) RemoveOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.removeOldest()
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
	delete(c.cache, e.Value.(Value).Path)
}

// Len returns the number of items in the cache.
func (c *lruCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}

// Clear purges all stored items from the cache.
func (c *lruCache) Clear() {
	c.mu.Lock()
	c.ll = list.New()
	me := c.maxEntries
	if me == 0 {
		me = 1024
	}
	c.cache = make(map[Path]*list.Element, me+1)
	c.mu.Unlock()
}
