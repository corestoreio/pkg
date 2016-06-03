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

package scope

import "fmt"

// MaxStoreID maximum allowed ID from package store. Doesn't matter whether we
// have a website, group or store ID. int24 (8388607) size at the moment.
const MaxStoreID int64 = 1<<23 - 1

// DefaultHash default Hash value for Default Scope and ID 0. Avoids typing
// 		scope.NewHash(DefaultID,0)
const DefaultHash Hash = Hash(Default)<<24 | 0

// Hash defines a merged Scope with its ID. The ID can either be from a website,
// group or store.
type Hash uint32

// If we have need for more store IDs then we can change the underlying types here.

// String human readable output
func (h Hash) String() string {
	scp, id := h.Unpack()
	return fmt.Sprintf("Scope(%s) ID(%d)", scp, id)
}

// Unpack extracts a Scope and its ID from a hash. Returned ID can be -1 when
// the Hash contains invalid data. An ID of -1 is considered an error.
func (h Hash) Unpack() (s Scope, id int64) {

	prospectS := h >> 24
	if prospectS > maxUint8 || prospectS < 0 {
		return Absent, -1
	}
	s = Scope(prospectS)

	h64 := int64(h)
	prospectID := h64 ^ (h64>>24)<<24
	if prospectID > MaxStoreID || prospectID < 0 {
		return Absent, -1
	}

	id = int64(prospectID)
	return
}

// HashMaxSegments maximum supported segments or also known as shards. This
// constant can be used to create the segmented array in other packages.
const HashMaxSegments uint16 = 256

const hashBitAnd Hash = Hash(HashMaxSegments) - 1

// Segment generates an 0 < ID <= 255 from a hash. Only used within an array
// index to optimize map[] usage in high concurrent situations. Also known as
// shard. An array of N shards is created, each shard contains its own instance
// of the cache with a lock. When an item with unique key needs to be cached a
// shard for it is chosen at first by the function Segment(). After that the
// cache lock is acquired and a write to the cache takes place. Reads are
// analogue.
func (h Hash) Segment() uint8 {
	return uint8(h & hashBitAnd)
}

// NewHash creates a new merged value. An error is equal to returning 0. An
// error occurs when id is greater than MaxStoreID or smaller 0. An errors
// occurs when the Scope is Default and id anything else than 0.
func NewHash(s Scope, id int64) Hash {
	if id > MaxStoreID || id < 0 {
		return 0
	}
	if s < Website {
		id = 0
	}
	return Hash(s)<<24 | Hash(id)
}
