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

import (
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/bufferpool"
)

// Perm is a bit set and used for permissions. Uint16 should be big enough.
type Perm uint16

// PermStore convenient helper contains all scope permission levels.
// The official core_config_data table and its classes to not support the
// GroupID scope, so that is the reason why PermStore does not have a GroupID.
const PermStore Perm = 1<<DefaultID | 1<<WebsiteID | 1<<StoreID

// PermWebsite convenient helper contains default and website scope permission levels.
const PermWebsite Perm = 1<<DefaultID | 1<<WebsiteID

// PermDefault convenient helper contains default scope permission level.
const PermDefault Perm = 1 << DefaultID

// PermStoreReverse convenient helper to enforce hierarchy levels.
// Only used in config.scopedService implementation.
const PermStoreReverse Perm = 1 << StoreID

// PermWebsiteReverse convenient helper to enforce hierarchy levels
// Only used in config.scopedService implementation.
const PermWebsiteReverse Perm = 1<<StoreID | 1<<WebsiteID

// All applies DefaultID, WebsiteID and StoreID scopes
func (bits Perm) All() Perm {
	return bits.Set(DefaultID, WebsiteID, StoreID)
}

// Set takes a variadic amount of Group to set them to Bits
func (bits Perm) Set(scopes ...Scope) Perm {
	for _, i := range scopes {
		bits = bits | (1 << i) // (1 << power = 2^power)
	}
	return bits
}

// Top returns the highest stored scope within a Perm.
// A Perm can consists of 3 scopes: 1. Default -> 2. Website -> 3. Store
// Highest scope for a Perm with all scopes is: Store.
func (bits Perm) Top() Scope {
	switch {
	case bits.Has(StoreID):
		return StoreID
	case bits.Has(WebsiteID):
		return WebsiteID
	}
	return DefaultID
}

// Has checks if a give scope exists within a Perm. Only the
// first argument is supported. Providing no argument assumes
// the scope.DefaultID.
func (bits Perm) Has(s ...Scope) bool {
	scp := DefaultID
	if len(s) > 0 {
		scp = s[0]
	}
	return (bits & Perm(1<<scp)) != 0
}

// Human readable representation of the permissions
func (bits Perm) Human() util.StringSlice {
	var ret util.StringSlice
	for i := uint(0); i < 64; i++ {
		bit := ((bits & (1 << i)) != 0)
		if bit {
			ret.Append(Scope(i).String())
		}
	}
	return ret
}

// String readable representation of the permissions
func (bits Perm) String() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	for i := uint(0); i < 64; i++ {
		if (bits & (1 << i)) != 0 {
			_, _ = buf.WriteString(Scope(i).String())
			_ = buf.WriteByte(',')
		}
	}
	buf.Truncate(buf.Len() - 1) // remove last colon
	return buf.String()

}

var nullByte = []byte("null")

// MarshalJSON implements marshaling into an array or null if no bits are set.
// Returns null when Perm is empty aka zero. null and 0 are considered the
// same for a later unmarshalling.
// @todo UnMarshal
func (bits Perm) MarshalJSON() ([]byte, error) {
	if bits == 0 {
		return nullByte, nil
	}
	return []byte(`["` + bits.Human().Join(`","`) + `"]`), nil
}
