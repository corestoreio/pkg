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

import (
	"fmt"
	"strconv"
	"strings"
)

// Scope used in constants where default is the lowest and store the highest.
// Func String() attached. Part of type Perm.
type Scope uint8

// *ID defines the overall scopes. The hierarchical order is always
// Default -> Website -> Group -> Store.
const (
	AbsentID Scope = iota // must start with 0
	DefaultID
	WebsiteID
	GroupID
	StoreID
)

// Interfaces for different scopes. Note that WebsiteIDer may have an underlying
// WebsiteCoder interface
type (
	// WebsiteIDer defines the scope of a website.
	WebsiteIDer interface {
		WebsiteID() int64
	}

	// GroupIDer defines the scope of a group.
	GroupIDer interface {
		GroupID() int64
	}

	// StoreIDer defines the scope of a store.
	StoreIDer interface {
		StoreID() int64
	}

	// GroupCoder not available because not existent.

	// WebsiteCoder defines the scope of a website by returning the store code.
	WebsiteCoder interface {
		WebsiteCode() string
	}

	// StoreCoder defines the scope of a store by returning the store code.
	StoreCoder interface {
		StoreCode() string
	}
)

const _ScopeName = "AbsentDefaultWebsiteGroupStore"

var _ScopeIndex = [...]uint8{0, 6, 13, 20, 25, 30}

// String human readable name of Group. For Marshaling see Perm
func (i Scope) String() string {
	if i+1 >= Scope(len(_ScopeIndex)) {
		return fmt.Sprintf("Scope(%d)", i)
	}
	return _ScopeName[_ScopeIndex[i]:_ScopeIndex[i+1]]
}

// PS path separator used in the database table core_config_data and in config.Service
const PS = "/"

// StrScope represents a string scope from table core_config_data column scope with
// special functions attached, mainly for path generation
type StrScope string

const (
	strDefault  = "default"
	strWebsites = "websites"
	strStores   = "stores"
)

// Str* constants are used in the database table core_config_data.
// StrDefault defines the global scope.
// StrWebsites defines the website scope which has default as parent and stores as child.
// StrStores defines the store scope which has default and websites as parent.
const (
	StrDefault  StrScope = strDefault
	StrWebsites StrScope = strWebsites
	StrStores   StrScope = strStores
)

// FQPath returns the fully qualified path. ID is an int string. Paths is either
// one path (system/smtp/host) including path separators or three
// parts ("system", "smtp", "host").
func (s StrScope) FQPath(scopeID string, paths ...string) string {
	return string(s) + PS + scopeID + PS + PathJoin(paths...)
}

// this "cache" should cover ~80% of all store setups
var int64Cache = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
}
var int64CacheLen = int64(len(int64Cache))

// FQPathInt64 same as FQPath() but for int64 scope IDs.
func (s StrScope) FQPathInt64(scopeID int64, paths ...string) string {
	scopeStr := "0"
	if scopeID > 0 {
		if scopeID <= int64CacheLen {
			scopeStr = int64Cache[scopeID]
		} else {
			scopeStr = strconv.FormatInt(scopeID, 10)
		}
	}
	return s.FQPath(scopeStr, paths...)
}

// String returns the scope as string
func (s StrScope) String() string {
	return string(s)
}

// FromString returns the scope ID from a string: default, websites or stores.
// Opposite of FromScope
func FromString(s string) Scope {
	switch StrScope(s) {
	case StrWebsites:
		return WebsiteID
	case StrStores:
		return StoreID
	}
	return DefaultID
}

// FromScope returns the string representation for a scope ID. Opposite of FromString.
func FromScope(scopeID Scope) StrScope {
	switch scopeID {
	case WebsiteID:
		return StrWebsites
	case StoreID:
		return StrStores
	}
	return StrDefault
}

// PathSplit splits a configuration path by the path separator PS.
func PathSplit(path string) []string {
	return strings.Split(path, PS)
}

// PathJoin joins configuration path parts by the path separator PS.
func PathJoin(path ...string) string {
	return strings.Join(path, PS)
}

func ReverseFQPath(fqPath string) (scope string, scopeID int64, path string, err error) {
	// todo optimize :-)
	paths := PathSplit(fqPath)
	if len(paths) < 5 {
		err = fmt.Errorf("Incorrect fully qualified path: %q", fqPath)
		return
	}
	scope = paths[0]
	scopeID, err = strconv.ParseInt(paths[1], 10, 64)
	path = PathJoin(paths[2:]...)
	return
}
