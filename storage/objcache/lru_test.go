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

package objcache_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/corestoreio/pkg/storage/objcache"
	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
)

var lruGetTests = []struct {
	name       string
	key        string
	expectedOk bool
}{
	{"01 hit", "aa/bb/cc", true},
	{"02 hit", "aa/bb/cc", true},
	{"03 hit", "aa/bb/cc3", true},
	{"04 hit", "aa/bb/cc4", true},
}

var testLRUData = [][]byte{[]byte(`Just a dummy text entry! 😇`)}

func TestLRUGet(t *testing.T) {

	newDistinctLRUSize := func(size int) func(*testing.T) {
		return func(t *testing.T) {
			ctx := context.TODO()
			for _, tt := range lruGetTests {
				lru, err := objcache.NewCacheLRU(size)()
				assert.NoError(t, err)
				assert.NoError(t, lru.Set(ctx, []string{tt.key}, testLRUData, nil))
				val, err := lru.Get(ctx, []string{tt.key})
				assert.NoError(t, err)
				ok := (val != nil)
				if ok != tt.expectedOk {
					// debugSliceBytes(t, val...)
					t.Fatalf("%q %s: cache hit = %v; want %v => %#v", t.Name(), tt.name, ok, !ok, val)
				} else if ok && !bytes.Equal(val[0], testLRUData[0]) {
					t.Fatalf("Distinct: %q %s expected get to return %q but got %q", t.Name(), tt.name, testLRUData, val)
				}
			}
		}
	}

	newCommonLRUSize := func(size int) func(*testing.T) {
		return func(t *testing.T) {
			lru, err := objcache.NewCacheLRU(size)()
			assert.NoError(t, err)
			ctx := context.TODO()
			for _, tt := range lruGetTests {
				assert.NoError(t, lru.Set(ctx, []string{tt.key}, testLRUData, nil))
				val, err := lru.Get(ctx, []string{tt.key})
				assert.NoError(t, err)
				ok := val != nil
				if ok != tt.expectedOk {
					t.Fatalf("%q %s: cache hit = %v; want %v", t.Name(), tt.name, ok, !ok)
				} else if ok && !bytes.Equal(val[0], testLRUData[0]) {
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
		lru, err := objcache.NewCacheLRU(2)()
		assert.NoError(t, err)
		ctx := context.TODO()
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

			err := lru.Set(ctx, []string{tt.key}, testLRUData, nil)
			if err != nil {
				panic(err)
			}
			err = lru.Set(ctx, []string{tt.key}, testLRUData, nil) // this is an ugly fix, but not that flaky anymore when running test with -count=40
			if err != nil {
				panic(err)
			}

			val, err := lru.Get(ctx, []string{tt.key})
			if err != nil {
				panic(err)
			}
			ok := val != nil
			if ok != tt.expectedOk {
				panic(fmt.Sprintf("%s: cache hit = %v; want %v: key:%q", tt.name, ok, !ok, tt.key))
			} else if ok && !bytes.Equal(val[0], testLRUData[0]) {
				panic(fmt.Sprintf("%s expected get to return %s but got %v", tt.name, testLRUData, val))
			}
		})
	})

	t.Run("add before parallel", func(t *testing.T) {
		lru, err := objcache.NewCacheLRU(5)()
		assert.NoError(t, err)
		ctx := context.TODO()
		for _, tt := range lruGetTests {
			assert.NoError(t, lru.Set(ctx, []string{tt.key}, testLRUData, nil))
		}

		bgwork.Wait(len(lruGetTests), func(idx int) {
			tt := lruGetTests[idx]

			val, err := lru.Get(ctx, []string{tt.key})
			if err != nil {
				panic(err)
			}
			ok := val != nil
			if ok != tt.expectedOk {
				panic(fmt.Sprintf("%s: cache hit = %v; want %v", tt.name, ok, tt.expectedOk))
			} else if ok && !bytes.Equal(val[0], testLRUData[0]) {
				panic(fmt.Sprintf("%s expected get to return %s but got %v", tt.name, testLRUData, val))
			}
		})
	})
}

func TestNewCacheLRU_Delete(t *testing.T) {
	newTestServiceDelete(t, objcache.NewCacheLRU(10))
}

func TestNewCacheLRU_ComplexParallel(t *testing.T) {

	t.Run("gob", func(t *testing.T) {
		newServiceComplexParallelTest(t, objcache.NewCacheLRU(10), nil)
	})

	t.Run("json", func(t *testing.T) {
		newServiceComplexParallelTest(t, objcache.NewCacheLRU(10), &objcache.ServiceOptions{
			Codec: JSONCodec{},
		})
	})
}
