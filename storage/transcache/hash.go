// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package transcache

// Hasher is responsible for generating unsigned, 64 bit hash of provided string. Hasher should minimize collisions
// (generating same hash for different strings) and while performance is also important fast functions are preferable (i.e.
// you can use FarmHash family).
type Hasher interface {
	Sum64([]byte) uint64
}

func newDefaultHasher() Hasher {
	return fnv64a{}
}

const (
	offset64 uint64 = 14695981039346656037
	prime64  uint64 = 1099511628211
)

type fnv64a struct{}

func (f fnv64a) Sum64(key []byte) uint64 {
	// Copyright 2011 The Go Authors. All rights reserved.
	// Use of this source code is governed by a BSD-style
	// license that can be found in the LICENSE file.

	// Package fnv implements FNV-1 and FNV-1a, non-cryptographic hash functions
	// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
	// See
	// https://en.wikipedia.org/wiki/Fowler-Noll-Vo_hash_function.
	var hash = offset64
	for _, c := range key {
		hash ^= uint64(c)
		hash *= prime64
	}
	return hash
}
