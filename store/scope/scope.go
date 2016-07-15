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
	"bytes"
	"fmt"
)

const maxUint8 = 1<<8 - 1

// Scope used in constants where default is the lowest and store the highest.
// Func String() attached. Part of type Perm.
type Scope uint8

// Those constants define the overall scopes. The hierarchical order is always:
// 		Absent -> Default -> Website -> Group -> Store
// These internal IDs may change without notice.
const (
	Absent Scope = iota // must start with 0
	Default
	Website
	Group
	Store
	maxScope
)

// Scoper specifies how to return the scope to which an ID belongs to. ID is one
// of a website, group or store ID as defined in their database tables. As long
// as we do not have returned a Default scope the ID (int64) must always be
// greater zero.
type Scoper interface {
	Scope() (Scope, int64)
}

const _ScopeName = "AbsentDefaultWebsiteGroupStore"

var _ScopeIndex = [...]uint8{0, 6, 13, 20, 25, 30}

// String human readable name of Group. For Marshaling see Perm
func (i Scope) String() string {
	if i+1 >= Scope(len(_ScopeIndex)) {
		return fmt.Sprintf("Scope(%d)", i)
	}
	return _ScopeName[_ScopeIndex[i]:_ScopeIndex[i+1]]
}

// StrScope converts the underlying scope ID to one of the three available scope
// strings in database table core_config_data.
func (i Scope) StrScope() string {
	return FromScope(i).String()
}

// Bytes returns the StrScope as byte slice from a Scope. The returned byte
// slice is owned by this package. You must copy it for further use.
func (i Scope) Bytes() []byte {
	switch i {
	case Website:
		return bWebsites
	case Store:
		return bStores
	}
	return bDefault
}

// StrScope represents a string scope from table core_config_data column scope
// with special functions attached, mainly for path generation
type StrScope string

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
	StrDefault  StrScope = strDefault
	StrWebsites StrScope = strWebsites
	StrStores   StrScope = strStores
)

// String returns the scope as string
func (s StrScope) String() string {
	return string(s)
}

// Scope returns the underlying scope
func (s StrScope) Scope() Scope {
	switch s {
	case StrWebsites:
		return Website
	case StrStores:
		return Store
	}
	return Default
}

// FromString returns the scope ID from a string: default, websites or stores.
// Opposite of FromScope
func FromString(s string) Scope {
	switch StrScope(s) {
	case StrWebsites:
		return Website
	case StrStores:
		return Store
	}
	return Default
}

// FromScope returns the string representation for a scope ID. Opposite of
// FromString.
func FromScope(scopeID Scope) StrScope {
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

// FromBytes returns the scope ID from a byte slice: default, websites or
// stores. Opposite of FromScope
func FromBytes(b []byte) Scope {
	switch {
	case bytes.Compare(bWebsites, b) == 0:
		return Website
	case bytes.Compare(bStores, b) == 0:
		return Store
	}
	return Default
}

// ValidBytes checks if b is a valid byte Scope of either StrDefault,
// StrWebsites or StrStores. Case-sensitive.
func ValidBytes(b []byte) bool {
	return bytes.Compare(bDefault, b) == 0 || bytes.Compare(bWebsites, b) == 0 || bytes.Compare(bStores, b) == 0
}

// ValidParent validates if the parent scope is within the hierarchical chain:
// default -> website -> store.
func ValidParent(current Scope, parent Scope) bool {
	return (parent == Default && current == Default) ||
		(parent == Default && current == Website) ||
		(parent == Website && current == Store)
}
