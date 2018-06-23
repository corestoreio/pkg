// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"bytes"
	"fmt"

	"github.com/corestoreio/errors"
)

const maxUint8 = 1<<8 - 1

// Type or also known as Scope defines the hierarchy of the overall CoreStore
// library. The hierarchy chain travels from Default->Website->Group->Store. The
// type relates to the database tables `website`, `store_group` and `store`.
// Type is a part of type Perm.
type Type uint8

// Those constants define the overall scopes. The hierarchical order is always:
// 		Absent -> Default -> Website -> Group -> Store
// These internal IDs may change without notice.
const (
	Absent Type = iota // must start with 0
	Default
	Website
	Group
	Store
	maxType
)

var (
	jsonDefault = []byte(`"Default"`)
	jsonWebsite = []byte(`"Website"`)
	jsonGroup   = []byte(`"Group"`)
	jsonStore   = []byte(`"Store"`)

	sbWebsite = []byte(`Website`)
	sbGroup   = []byte(`Group`)
	sbStore   = []byte(`Store`)
)

const _TypeName = "AbsentDefaultWebsiteGroupStore"

var _TypeIndex = [...]uint8{0, 6, 13, 20, 25, 30}

// String human readable name of a Type. For Marshaling see Perm.
func (s Type) String() string {
	if s+1 >= Type(len(_TypeIndex)) {
		return fmt.Sprintf("Type(%d)", s)
	}
	return _TypeName[_TypeIndex[s]:_TypeIndex[s+1]]
}

// StrType converts the underlying Type to one of the three available type
// strings from the database table `core_config_data`.
func (s Type) StrType() string {
	return FromType(s).String()
}

// MarshalJSON implements the Marshaler interface. The returned byte slice is
// owned by the callee. You must copy it for further use.
func (s Type) MarshalJSON() ([]byte, error) {
	var ret []byte
	switch s {
	case Website:
		ret = jsonWebsite
	case Group:
		ret = jsonGroup
	case Store:
		ret = jsonStore
	default:
		ret = jsonDefault
	}
	return ret, nil
}

// UnmarshalJSON implements the Unmarshaler interface
func (s *Type) UnmarshalJSON(b []byte) error {
	*s = FromBytes(b)
	return nil
}

// StrBytes returns the TypeStr as byte slice from a Type. The returned byte
// slice is owned by the callee. You must copy it for further use.
func (s Type) StrBytes() []byte {
	switch s {
	case Website:
		return bWebsites
	case Store:
		return bStores
	}
	return bDefault
}

// IsWebSiteOrStore returns true if the type is either Website or Store.
func (s Type) IsWebSiteOrStore() bool {
	return s == Website || s == Store
}

// WithID calls MakeTypeID for your convenience. It packs the id into a new value
// containing the Type and its ID ;-).
func (s Type) WithID(id int64) TypeID {
	return MakeTypeID(s, id)
}

// IsValid checks if the type is within the scope Default, Website, Group or
// Store.
func (s Type) IsValid() error {
	if s >= maxType {
		return errors.NotValid.Newf("[scope] Invalid Type: %s", s)
	}
	return nil
}

// TypeStr represents a string Type from table `core_config_data` column scope
// with special functions attached, mainly for path generation
type TypeStr string

const (
	strDefault  = "default"
	strWebsites = "websites"
	strStores   = "stores"
)

var (
	bDefault  = []byte(strDefault)
	bWebsites = []byte(strWebsites)
	bStores   = []byte(strStores)
)

// Str* constants are used in the database table core_config_data. StrDefault
// defines the global scope. StrWebsites defines the website scope which has
// default as parent and stores as child. StrStores defines the store scope
// which has default and websites as parent.
const (
	StrDefault  TypeStr = strDefault
	StrWebsites TypeStr = strWebsites
	StrStores   TypeStr = strStores
)

// String returns the scope as string
func (s TypeStr) String() string {
	return string(s)
}

// Type returns the underlying type.
func (s TypeStr) Type() Type {
	switch s {
	case StrWebsites:
		return Website
	case StrStores:
		return Store
	}
	return Default
}

// FromString returns the Type from a string: default, websites or stores.
// Opposite of FromType.
func FromString(s string) Type {
	switch TypeStr(s) {
	case StrWebsites:
		return Website
	case StrStores:
		return Store
	}
	return Default
}

// FromType returns the string representation for a Type. Opposite of
// FromString.
func FromType(scopeID Type) TypeStr {
	switch scopeID {
	case Website:
		return StrWebsites
	case Store:
		return StrStores
	}
	return StrDefault
}

// Valid checks if s is a valid StrScope of either StrDefault, StrWebsites or
// StrStores. Case-sensitive. Input should all be lowercase.
func Valid(s string) bool {
	switch s {
	case strWebsites, strStores, strDefault:
		return true
	}
	return false
}

// FromBytes returns the Type from a byte slice. Supported values are
// default, websites, stores, Default, Website, Group and store. Case sensitive.
func FromBytes(b []byte) Type {
	switch {
	case bytes.Equal(bWebsites, b):
		return Website
	case bytes.Equal(bStores, b):
		return Store

	case bytes.Equal(jsonWebsite, b):
		return Website
	case bytes.Equal(jsonGroup, b):
		return Group
	case bytes.Equal(jsonStore, b):
		return Store

	case bytes.Equal(sbWebsite, b):
		return Website
	case bytes.Equal(sbGroup, b):
		return Group
	case bytes.Equal(sbStore, b):
		return Store
	}
	return Default
}

// ValidBytes checks if b is a valid byte Type of either StrDefault,
// StrWebsites or StrStores. Case-sensitive.
func ValidBytes(b []byte) bool {
	return bytes.Equal(bDefault, b) || bytes.Equal(bWebsites, b) || bytes.Equal(bStores, b)
}

// ValidParent validates if the parent scope is within the hierarchical chain:
// default -> website -> store.
func ValidParent(current Type, parent Type) bool {
	return (parent == Default && current == Default) ||
		(parent == Default && current == Website) ||
		(parent == Website && current == Store)
}
