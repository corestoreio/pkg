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

// MaxStoreID maximum allowed store ID
const MaxStoreID = 1<<23 - 1

// Hash defines a merged Scope with its ID. The ID can either be from
// a website, group or store.
type Hash uint32

// If we have need for more store IDs then we can change the underlying types here.

// String human readable output
func (h Hash) String() string {
	scp, id := h.Unpack()
	return fmt.Sprintf("Scope(%s) ID(%d)", scp, id)
}

// Unpack extracts a Scope and its ID from a hash.
// Returned ID can be -1 when the Hash contains invalid data.
// An ID of -1 is considered an error.
func (h Hash) Unpack() (s Scope, id int) {

	prospectS := h >> 24
	if prospectS > maxUint8 || prospectS < 0 {
		return AbsentID, -1
	}
	s = Scope(prospectS)

	prospectID := h ^ (h>>24)<<24
	if prospectID > MaxStoreID || prospectID < 0 {
		return AbsentID, -1
	}

	id = int(prospectID)
	return
}

// NewHash creates a new merged value. An error is equal to returning 0.
// An error occurs when id is greater than MaxStoreID or smaller 0.
func NewHash(s Scope, id int) Hash {
	if id > MaxStoreID || id < 0 {
		return 0
	}
	return Hash(s)<<24 | Hash(id)
}
