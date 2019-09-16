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

package lru

import (
	"testing"
)

type MyValue []byte

func (mv MyValue) Size() int {
	return cap(mv)
}

func BenchmarkGet(b *testing.B) {
	cache := New(64 * 1024 * 1024)
	value := make(MyValue, 1000)
	cache.Set("stuff1", value)
	value2 := make(MyValue, 1100)
	cache.Set("stuff2", value2)

	b.Run("one key", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				val, ok := cache.Get("stuff1")
				if !ok {
					panic("error")
				}
				_ = val
			}
		})
	})

	b.Run("two keys", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := "stuff1"
				if i%2 == 0 {
					key = "stuff2"
				}
				val, ok := cache.Get(key)
				if !ok {
					panic("error")
				}
				_ = val
				i++
			}
		})
	})
	b.Run("set get, two keys", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := "stuff1"
				if i%2 == 0 {
					key = "stuff2"
				}
				val, ok := cache.Get(key)
				if !ok {
					panic("error")
				}
				_ = val
				if i%4 == 0 {
					cache.Set("stuff3", value)
				}
				i++
			}
		})
	})
}
