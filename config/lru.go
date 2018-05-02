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
)

// lruCache is an LRU cache. It is safe for concurrent access.
// This type does not get exported.
type lruCache struct {
	// maxEntries is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	maxEntries int

	// OnEvicted optionally specificies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key Path, value Value)

	mu    sync.Mutex
	ll    *list.List
	cache map[Path]*list.Element
}

type entry struct {
	key   Path
	value Value
}

// newLRU creates a new lruCache. If maxEntries is zero, the cache has no limit
// and it's assumed that eviction is done by the caller. This type does not get
// exported.
func newLRU(maxEntries int) *lruCache {
	return &lruCache{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[Path]*list.Element, maxEntries),
	}
}

// Add adds a value to the cache.
func (c *lruCache) Add(key Path, value Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ev := ee.Value.(*entry)
		ev.value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
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
		return ele.Value.(*entry).value, true
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
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
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
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.ll = list.New()
	c.cache = make(map[Path]*list.Element)
	c.mu.Unlock()
}
