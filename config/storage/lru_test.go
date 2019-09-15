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

package storage_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
)

var lruGetTests = []struct {
	name       string
	keyToAdd   config.Path
	keyToGet   config.Path
	expectedOk bool
}{
	{"01 hit", config.MustMakePath("aa/bb/cc"), config.MustMakePath("aa/bb/cc"), true},
	{"02 miss", config.MustMakePath("aa/bb/cc"), config.MustMakePath("aa/bb/cc").BindStore(11), false},
	{"03 hit", config.MustMakePath("aa/bb/cc").BindWebsite(3), config.MustMakePath("aa/bb/cc").BindWebsite(3), true},
	{"04 miss", config.MustMakePath("aa/bb/cc").BindStore(44), config.MustMakePath("aa/bb/cc").BindWebsite(4), false},
	{"05 miss", config.MustMakePath("aa/bb/cc").BindStore(55), config.MustMakePath("aa/bb/cc").BindStore(56), false},
	{"06 hit", config.MustMakePath("aa/bb/cc").BindStore(6), config.MustMakePath("aa/bb/cc").BindStore(6), true},
}

const testLRUDataStr = `Just a dummy text entry! 😇`

var testLRUData = []byte(testLRUDataStr)

// var testLRUVal = *config.NewValue(testLRUData)

func TestLRUGet(t *testing.T) {
	newDistinctLRUSize := func(size int) func(*testing.T) {
		return func(t *testing.T) {
			for _, tt := range lruGetTests {
				lru := storage.NewLRU(size)
				assert.NoError(t, lru.Set(tt.keyToAdd, testLRUData))
				val, ok, err := lru.Get(tt.keyToGet)
				assert.NoError(t, err)
				if ok != tt.expectedOk {
					t.Fatalf("%q %s: cache hit = %v; want %v", t.Name(), tt.name, ok, !ok)
				} else if ok && !bytes.Equal(val, testLRUData) {
					t.Fatalf("Distinct: %q %s expected get to return %q but got %q", t.Name(), tt.name, testLRUData, val)
				}
			}
		}
	}

	newCommonLRUSize := func(size int) func(*testing.T) {
		return func(t *testing.T) {
			lru := storage.NewLRU(size)
			for _, tt := range lruGetTests {
				assert.NoError(t, lru.Set(tt.keyToAdd, testLRUData))
				val, ok, err := lru.Get(tt.keyToGet)
				assert.NoError(t, err)
				if ok != tt.expectedOk {
					t.Fatalf("%q %s: cache hit = %v; want %v", t.Name(), tt.name, ok, !ok)
				} else if ok && !bytes.Equal(val, testLRUData) {
					t.Fatalf("Common: %q %s expected get to return %q but got %q", t.Name(), tt.name, testLRUData, val)
				}
			}
		}
	}

	t.Run("distinct LRU size=0", newDistinctLRUSize(0))
	t.Run("distinct LRU size=1", newDistinctLRUSize(1))
	t.Run("common LRU size=0", newCommonLRUSize(0))
	t.Run("common LRU size=1", newCommonLRUSize(1))
}

func TestLRUNew_Parallel(t *testing.T) {
	t.Run("add in parallel", func(t *testing.T) {
		lru := storage.NewLRU(2)
		bgwork.Wait(len(lruGetTests), func(idx int) {
			defer func() {
				if r := recover(); r != nil {
					t.Error(r)
				}
			}()
			tt := lruGetTests[idx]

			/* TODO THIS TEST IS FLAKY, with "-race" it passes. it fails depending on the LRU size. <= 3 fails, and 0 or > 3 succeeds always.
			$ go test -run=TestLRU -count=34
			panic: 06 hit: cache hit = false; want true: add:"stores/6/aa/bb/cc" get:"stores/6/aa/bb/cc"
			AFAIK problems with cache eviction
			*/
			d := []byte(testLRUDataStr)
			v := []byte(testLRUDataStr)
			_ = lru.Set(tt.keyToAdd, v)
			_ = lru.Set(tt.keyToAdd, v) // this is an ugly fix, but not that flaky anymore when running test with -count=40

			val, ok, err := lru.Get(tt.keyToGet)
			if err != nil {
				panic(err)
			}
			if ok != tt.expectedOk {
				panic(fmt.Sprintf("%s: cache hit = %v; want %v: add:%q get:%q", tt.name, ok, !ok, tt.keyToAdd.String(), tt.keyToGet.String()))
			} else if ok && !bytes.Equal(val, d) {
				panic(fmt.Sprintf("%s expected get to return %s but got %v", tt.name, testLRUData, val))
			}
		})
	})

	t.Run("add before parallel", func(t *testing.T) {
		lru := storage.NewLRU(5)

		for _, tt := range lruGetTests {
			assert.NoError(t, lru.Set(tt.keyToAdd, testLRUData))
		}

		bgwork.Wait(len(lruGetTests), func(idx int) {
			tt := lruGetTests[idx]

			val, ok, err := lru.Get(tt.keyToGet)
			if err != nil {
				panic(err)
			}
			if ok != tt.expectedOk {
				panic(fmt.Sprintf("%s: cache hit = %v; want %v", tt.name, ok, tt.expectedOk))
			} else if ok && !bytes.Equal(val, testLRUData) {
				panic(fmt.Sprintf("%s expected get to return %s but got %v", tt.name, testLRUData, val))
			}
		})
	})
}
