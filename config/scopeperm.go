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

package config

import "github.com/corestoreio/csfw/utils"

// ScopePerm is a bit set and used for permissions, ScopeGroup is not a part of this bit set.
// Type ScopeGroup is a subpart of ScopePerm
type ScopePerm uint64

// ScopePermAll convenient helper variable contains all scope permission levels
var ScopePermAll = ScopePerm(1<<ScopeDefaultID | 1<<ScopeWebsiteID | 1<<ScopeStoreID)

// NewScopePerm returns a new permission container
func NewScopePerm(scopes ...ScopeGroup) ScopePerm {
	p := ScopePerm(0)
	p.Set(scopes...)
	return p
}

// All applies all scopes
func (bits *ScopePerm) All() ScopePerm {
	bits.Set(ScopeDefaultID, ScopeWebsiteID, ScopeStoreID)
	return *bits
}

// Set takes a variadic amount of ScopeGroup to set them to ScopeBits
func (bits *ScopePerm) Set(scopes ...ScopeGroup) ScopePerm {
	for _, i := range scopes {
		*bits = *bits | (1 << i) // (1 << power = 2^power)
	}
	return *bits
}

// Has checks if ScopeGroup is in ScopeBits
func (bits ScopePerm) Has(s ScopeGroup) bool {
	var one ScopeGroup = 1 // ^^
	return (bits & ScopePerm(one<<s)) != 0
}

// Human readable representation of the permissions
func (bits ScopePerm) Human() utils.StringSlice {
	var ret utils.StringSlice
	var i uint
	for i = 0; i < 64; i++ {
		bit := ((bits & (1 << i)) != 0)
		if bit {
			ret.Append(ScopeGroup(i).String())
		}
	}
	return ret
}

// MarshalJSON implements marshaling into an array or null if no bits are set. @todo UnMarshal
func (bits ScopePerm) MarshalJSON() ([]byte, error) {
	if bits == 0 {
		return []byte("null"), nil
	}
	return []byte(`["` + bits.Human().Join(`","`) + `"]`), nil
}
