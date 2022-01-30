/*
Copyright 2017 Google Inc.

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

// Package cache implements a LRU cache.
//
// The implementation borrows heavily from SmallLRUCache
// (originally by Nathan Schrenk). The object maintains a doubly-linked list of
// elements. When an element is accessed, it is promoted to the head of the
// list. When space is needed, the element at the tail of the list
// (the least recently used element) is evicted.
package lru

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// Cache is a typical LRU cache implementation.  If the cache
// reaches the capacity, the least recently used item is deleted from
// the cache. Note the capacity is not the number of items, but the
// total sum of the Size() of each item.
type Cache[K comparable] struct {
	mu sync.Mutex

	// list & table contain *entry objects.
	list  *list.List
	table map[K]*list.Element

	size      int64
	capacity  int64
	evictions int64
}

// Value is the interface values that go into Cache need to satisfy
type Value interface {
	// Size returns how big this value is. If you want to just track
	// the cache by number of objects, you may return the size as 1.
	Size() int
}

// Item is what is stored in the cache
type Item[K comparable] struct {
	Key   K
	Value Value
}

type entry[K comparable] struct {
	key          K
	value        Value
	size         int64
	timeAccessed time.Time
}

// New creates a new empty cache with the given capacity.
func New[K comparable](capacity int64) *Cache[K] {
	return &Cache[K]{
		list:     list.New(),
		table:    make(map[K]*list.Element, capacity),
		capacity: capacity,
	}
}

// Get returns a value from the cache, and marks the entry as most
// recently used.
func (lru *Cache[K]) Get(key K) (v Value, ok bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element, ok := lru.table[key]
	if element == nil {
		return nil, false
	}
	lru.moveToFront(element)
	return element.Value.(*entry[K]).value, true
}

// Peek returns a value from the cache without changing the LRU order.
func (lru *Cache[K]) Peek(key K) (v Value, ok bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element := lru.table[key]
	if element == nil {
		return nil, false
	}
	return element.Value.(*entry[K]).value, true
}

// Set sets a value in the cache.
func (lru *Cache[K]) Set(key K, value Value) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element := lru.table[key]; element != nil {
		lru.updateInplace(element, value)
	} else {
		lru.addNew(key, value)
	}
}

// SetIfAbsent will set the value in the cache if not present. If the
// value exists in the cache, we don't set it.
func (lru *Cache[K]) SetIfAbsent(key K, value Value) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element := lru.table[key]; element != nil {
		lru.moveToFront(element)
	} else {
		lru.addNew(key, value)
	}
}

// Delete removes an entry from the cache, and returns if the entry existed.
func (lru *Cache[K]) Delete(key K) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element := lru.table[key]
	if element == nil {
		return false
	}

	lru.list.Remove(element)
	delete(lru.table, key)
	lru.size -= element.Value.(*entry[K]).size
	return true
}

// Clear will clear the entire cache.
func (lru *Cache[K]) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.list.Init()
	lru.table = make(map[K]*list.Element, lru.capacity)
	lru.size = 0
}

// SetCapacity will set the capacity of the cache. If the capacity is
// smaller, and the current cache size exceed that capacity, the cache
// will be shrank.
func (lru *Cache[K]) SetCapacity(capacity int64) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.capacity = capacity
	lru.checkCapacity()
}

// Stats returns a few stats on the cache.
func (lru *Cache[K]) Stats() (length, size, capacity, evictions int64, oldest time.Time) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if lastElem := lru.list.Back(); lastElem != nil {
		oldest = lastElem.Value.(*entry[K]).timeAccessed
	}
	return int64(lru.list.Len()), lru.size, lru.capacity, lru.evictions, oldest
}

// StatsJSON returns stats as a JSON object in a string.
func (lru *Cache[K]) StatsJSON() string {
	if lru == nil {
		return "{}"
	}
	l, s, c, e, o := lru.Stats()
	return fmt.Sprintf("{\"Length\": %v, \"Size\": %v, \"Capacity\": %v, \"Evictions\": %v, \"OldestAccess\": \"%v\"}", l, s, c, e, o)
}

// Length returns how many elements are in the cache
func (lru *Cache[K]) Length() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return int64(lru.list.Len())
}

// Size returns the sum of the objects' Size() method.
func (lru *Cache[K]) Size() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.size
}

// Capacity returns the cache maximum capacity.
func (lru *Cache[K]) Capacity() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.capacity
}

// Evictions returns the eviction count.
func (lru *Cache[K]) Evictions() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.evictions
}

// Oldest returns the insertion time of the oldest element in the cache,
// or a IsZero() time if cache is empty.
func (lru *Cache[K]) Oldest() (oldest time.Time) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if lastElem := lru.list.Back(); lastElem != nil {
		oldest = lastElem.Value.(*entry[K]).timeAccessed
	}
	return
}

// Keys returns all the keys for the cache, ordered from most recently
// used to last recently used.
func (lru *Cache[K]) Keys() []K {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	keys := make([]K, 0, lru.list.Len())
	for e := lru.list.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(*entry[K]).key)
	}
	return keys
}

// Items returns all the values for the cache, ordered from most recently
// used to last recently used.
func (lru *Cache[K]) Items() []Item[K] {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	items := make([]Item[K], 0, lru.list.Len())
	for e := lru.list.Front(); e != nil; e = e.Next() {
		v := e.Value.(*entry[K])
		items = append(items, Item[K]{Key: v.key, Value: v.value})
	}
	return items
}

func (lru *Cache[K]) updateInplace(element *list.Element, value Value) {
	valueSize := int64(value.Size())
	sizeDiff := valueSize - element.Value.(*entry[K]).size
	element.Value.(*entry[K]).value = value
	element.Value.(*entry[K]).size = valueSize
	lru.size += sizeDiff
	lru.moveToFront(element)
	lru.checkCapacity()
}

func (lru *Cache[K]) moveToFront(element *list.Element) {
	lru.list.MoveToFront(element)
	element.Value.(*entry[K]).timeAccessed = time.Now()
}

func (lru *Cache[K]) addNew(key K, value Value) {
	newEntry := &entry[K]{key: key, value: value, size: int64(value.Size()), timeAccessed: time.Now()}
	element := lru.list.PushFront(newEntry)
	lru.table[key] = element
	lru.size += newEntry.size
	lru.checkCapacity()
}

func (lru *Cache[K]) checkCapacity() {
	// Partially duplicated from Delete
	for lru.size > lru.capacity {
		delElem := lru.list.Back()
		delValue := delElem.Value.(*entry[K])
		lru.list.Remove(delElem)

		delete(lru.table, delValue.key)
		lru.size -= delValue.size
		lru.evictions++
	}
}
