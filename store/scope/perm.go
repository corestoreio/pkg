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

import "github.com/corestoreio/csfw/util/bufferpool"

// Perm is a bit set and used for permissions. Uint16 should be big enough.
type Perm uint16

// PermStore convenient helper contains all scope permission levels.
// The official core_config_data table and its classes to not support the
// GroupID scope, so that is the reason why PermStore does not have a GroupID.
const PermStore Perm = 1<<Default | 1<<Website | 1<<Store

// PermWebsite convenient helper contains default and website scope permission levels.
const PermWebsite Perm = 1<<Default | 1<<Website

// PermDefault convenient helper contains default scope permission level.
const PermDefault Perm = 1 << Default

// PermStoreReverse convenient helper to enforce hierarchy levels.
// Only used in config.scopedService implementation.
const PermStoreReverse Perm = 1 << Store

// PermWebsiteReverse convenient helper to enforce hierarchy levels
// Only used in config.scopedService implementation.
const PermWebsiteReverse Perm = 1<<Store | 1<<Website

// All applies DefaultID, WebsiteID and StoreID scopes
func (bits Perm) All() Perm {
	return bits.Set(Default, Website, Store)
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
	case bits.Has(Store):
		return Store
	case bits.Has(Website):
		return Website
	}
	return Default
}

// Has checks if a give scope exists within a Perm. Only the
// first argument is supported. Providing no argument assumes
// the scope.DefaultID.
func (bits Perm) Has(s ...Scope) bool {
	scp := Default
	if len(s) > 0 {
		scp = s[0]
	}
	return (bits & Perm(1<<scp)) != 0
}

// Human readable representation of the permissions
func (bits Perm) Human() []string {
	var ret = make([]string, 0, maxScope)
	for i := uint(0); i < uint(maxScope); i++ {
		bit := ((bits & (1 << i)) != 0)
		if bit {
			ret = append(ret, (Scope(i).String()))
		}
	}
	return ret
}

// String readable representation of the permissions
func (bits Perm) String() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	for i := uint(0); i < uint(maxScope); i++ {
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
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	buf.WriteString(`["`)
	hm := bits.Human()
	lhm := len(hm) - 1
	for i, h := range hm {
		buf.WriteString(h)
		if i < lhm {
			buf.WriteString(`","`)
		}
	}
	buf.WriteString(`"]`)

	// seems redundant but we must copy the bytes aways because bufferpool.Put() resets the buffer
	return []byte(buf.String()), nil
}
