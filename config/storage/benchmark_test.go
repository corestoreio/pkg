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

package storage_test

import (
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/store/scope"
)

var benchmarkToEnvVar string

// BenchmarkToEnvVar-4   	 3000000	       570 ns/op	     163 B/op	       6 allocs/op
func BenchmarkToEnvVar(b *testing.B) {
	p := config.MustNewPathWithScope(scope.Store.WithID(543), "aa/bb/cc")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkToEnvVar = storage.ToEnvVar(p)
	}
}

var benchmarkFromEnvVar *config.Path

// BenchmarkFromEnvVar-4   	 3000000	       501 ns/op	     184 B/op	       7 allocs/op
func BenchmarkFromEnvVar(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		if benchmarkFromEnvVar, err = storage.FromEnvVar(storage.Prefix, "CONFIG__WEBSITES__321__AA__BB__CC"); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLRUNew_Parallel/single-4         	 2000000	       854 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLRUNew_Parallel/parallel-4       	 1000000	      1274 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLRUNew_Parallel(b *testing.B) {
	const cacheSize = 5
	b.Run("single", func(b *testing.B) {
		lru := storage.NewLRU(cacheSize)
		for _, tt := range lruGetTests {
			if err := lru.Set(tt.keyToAdd, testLRUData); err != nil {
				b.Fatal(err)
			}
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, tt := range lruGetTests {
				_, ok, _ := lru.Get(tt.keyToGet)
				if ok != tt.expectedOk {
					b.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
				}
			}
		}
	})
	b.Run("parallel", func(b *testing.B) {
		lru := storage.NewLRU(cacheSize)
		for _, tt := range lruGetTests {
			if err := lru.Set(tt.keyToAdd, testLRUData); err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var val2 []byte
			for pb.Next() {
				for _, tt := range lruGetTests {
					val, ok, _ := lru.Get(tt.keyToGet)
					if ok != tt.expectedOk {
						b.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
					}
					val2 = val
				}
			}
			_ = val2 // ;-)
		})
	})
}
