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

package hashpool

import (
	"crypto/hmac"
	"hash"
	"sync"

	"github.com/corestoreio/csfw/util/errors"
)

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fnv implements FNV-1 and FNV-1a, non-cryptographic hash functions
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See
// https://en.wikipedia.org/wiki/Fowler-Noll-Vo_hash_function.
const offset64 hash64 = 14695981039346656037
const prime64 hash64 = 1099511628211

// http://programmers.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed
// maybe switch to siphash if there are too many collisions ...

type hash64 uint64

// fnv64a
func (h hash64) writeBytes(data []byte) hash64 {
	if h == 0 {
		h = offset64
	}
	for _, c := range data {
		h ^= hash64(c)
		h *= prime64
	}
	return h
}

// fnv64a
func (h hash64) writeStr(data string) hash64 {
	if h == 0 {
		h = offset64
	}
	for i := 0; i < len(data); i++ {
		h ^= hash64(data[i])
		h *= prime64
	}
	return h
}

// END Copyright 2011 The Go Authors. All rights reserved.

type htVal struct {
	hashedName hash64
	Tank
	hf func() hash.Hash
}

func makeHtVal(name hash64, hh func() hash.Hash) htVal {
	return htVal{
		hashedName: name,
		Tank:       New(hh),
		hf:         hh,
	}
}

var db = &struct {
	sync.RWMutex
	// ht == hash tank
	ht map[hash64]htVal
}{
	ht: make(map[hash64]htVal),
}

// Register registers a new hash function in the global map. Returns an error
// behaviour of AlreadyExists if you try to register the same name twice. Safe
// for concurrent use.
func Register(name string, hh func() hash.Hash) error {
	index := hash64(0).writeStr(name)
	db.Lock()
	defer db.Unlock()
	if _, ok := db.ht[index]; ok {
		return errors.NewAlreadyExistsf("[hashpool] %q has already been registered", name)
	}
	db.ht[index] = makeHtVal(index, hh)
	return nil
}

// Deregister removes a previously registered hash and all its HMAC hashes. Safe
// for concurrent use.
func Deregister(name string) {
	index := hash64(0).writeStr(name)
	db.Lock()
	delete(db.ht, index)
	for k, v := range db.ht {
		// find all hmac hashes with the same name
		if v.hashedName == index {
			delete(db.ht, k)
		}
	}
	db.Unlock()
}

// FromRegistry returns a hash pool for the given name, otherwise a NotFound
// error behaviour. Safe for concurrent use.
func FromRegistry(name string) (Tank, error) {
	db.RLock()
	defer db.RUnlock()
	ht, ok := db.ht[hash64(0).writeStr(name)]
	if !ok {
		return Tank{}, errors.NewNotFoundf("[hashpool] Unknown Hash %q. Not yet registered?", name)
	}
	return ht.Tank, nil
}

// MustFromRegistry same as FromRegistry but panics when the name does not
// exists in the internal map. Safe for concurrent use.
func MustFromRegistry(name string) Tank {
	t, err := FromRegistry(name)
	if err != nil {
		panic(err)
	}
	return t
}

// FromRegisterHMAC returns a pool with a HMAC hash from a previously registered
// hash. Key aka. password. The newly generated hash pool will be internally
// cached with the identifier of name plus key.
func FromRegistryHMAC(name string, key []byte) (Tank, error) {
	nameIndex := hash64(0).writeStr(name)

	hashTnk, ok := db.ht[nameIndex]
	if !ok {
		return Tank{}, errors.NewNotFoundf("[hashpool] Unknown Hash %q. Not yet registered?", name)
	}

	index := nameIndex.writeBytes(key)

	db.RLock()
	ht, ok := db.ht[index]
	db.RUnlock()
	if !ok {
		ht = makeHtVal(nameIndex, func() hash.Hash {
			return hmac.New(hashTnk.hf, key)
		})
		db.Lock()
		db.ht[index] = ht
		db.Unlock()
	}
	return ht.Tank, nil
}
