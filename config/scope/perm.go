// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

import "github.com/corestoreio/csfw/util"

// Perm is a bit set and used for permissions, Group is not a part of this bit set.
// Type Group is a subpart of Perm
type Perm uint64

// PermAll convenient helper variable contains all scope permission levels
var PermAll = Perm(1<<DefaultID | 1<<WebsiteID | 1<<StoreID)

// NewPerm returns a new permission container
func NewPerm(scopes ...Scope) Perm {
	p := Perm(0)
	p.Set(scopes...)
	return p
}

// All applies all scopes
func (bits *Perm) All() Perm {
	bits.Set(DefaultID, WebsiteID, StoreID)
	return *bits
}

// Set takes a variadic amount of Group to set them to Bits
func (bits *Perm) Set(scopes ...Scope) Perm {
	for _, i := range scopes {
		*bits = *bits | (1 << i) // (1 << power = 2^power)
	}
	return *bits
}

// Has checks if Group is in Bits
func (bits Perm) Has(s Scope) bool {
	var one Scope = 1 // ^^
	return (bits & Perm(one<<s)) != 0
}

// Human readable representation of the permissions
func (bits Perm) Human() util.StringSlice {
	var ret util.StringSlice
	var i uint
	for i = 0; i < 64; i++ {
		bit := ((bits & (1 << i)) != 0)
		if bit {
			ret.Append(Scope(i).String())
		}
	}
	return ret
}

// MarshalJSON implements marshaling into an array or null if no bits are set. @todo UnMarshal
func (bits Perm) MarshalJSON() ([]byte, error) {
	if bits == 0 {
		return []byte("null"), nil
	}
	return []byte(`["` + bits.Human().Join(`","`) + `"]`), nil
}
