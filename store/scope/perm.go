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

// Perm is a bit set and used for permissions, Group is not a part of this bit set.
// Type Group is a subpart of Perm
type Perm uint64

// PermStore convenient helper variable contains all scope permission levels.
// The official core_config_data table and its classes to not support the
// GroupID scope, so that is the reasion why PermStore does not have a GroupID.
var PermStore = NewPerm(DefaultID, WebsiteID, StoreID)

// PermGroup convenient helper variable contains default, website and group scope permission levels.
// Not officially supported by M1 and M2.
var PermGroup = NewPerm(DefaultID, WebsiteID, GroupID)

// PermWebsite convenient helper variable contains default and website scope permission levels.
var PermWebsite = NewPerm(DefaultID, WebsiteID)

// PermDefault convenient helper variable contains default scope permission level.
var PermDefault = NewPerm(DefaultID)

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
			buf.WriteString(Scope(i).String())
			buf.WriteByte(',')
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
