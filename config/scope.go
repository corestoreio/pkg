// Copyright 2015 CoreStore Authors
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

import (
	"fmt"

	"github.com/corestoreio/csfw/utils"
)

const (
	// Scope*ID defines the overall scopes in a configuration. If a Section/Group/Field
	// can be shown in the current selected scope.
	ScopeAbsentID ScopeGroup = iota // order of the constants is used for comparison
	ScopeDefaultID
	ScopeWebsiteID
	ScopeGroupID
	ScopeStoreID
)

const (
	// StringScopeDefault defines the global scope. Stored in table core_config_data.scope.
	ScopeRangeDefault = "default"
	// StringScopeWebsites defines the website scope which has default as parent and stores as child.
	//  Stored in table core_config_data.scope.
	ScopeRangeWebsites = "websites"
	// StringScopeStores defines the store scope which has default and websites as parent.
	//  Stored in table core_config_data.scope.
	ScopeRangeStores = "stores"
)

type (

	// ScopeGroup used in constants where default is the lowest and store the highest. Func String() attached.
	// Part of ScopePerm.
	ScopeGroup uint8

	// ScopeIDer implements how to get the ID. If ScopeIDer implements ScopeCoder
	// then ScopeCoder has precedence. ID can be any of the website, group or store IDs.
	ScopeIDer interface {
		ScopeID() int64
	}
	// ScopeCoder implements how to get an object by Code which can be website or store code.
	// Groups doesn't have codes.
	ScopeCoder interface {
		ScopeCode() string
	}
	// ID is convenience helper to satisfy the interface ScopeIDer.
	ScopeID int64
	// Code is convenience helper to satisfy the interface ScopeCoder and ScopeIDer.
	ScopeCode string
)

var _ ScopeIDer = ScopeID(0)
var _ ScopeCoder = ScopeCode("")

// ScopeID is convenience helper to satisfy the interface ScopeIDer
func (i ScopeID) ScopeID() int64 { return int64(i) }

// ScopeID is a noop method receiver to satisfy the interface ScopeIDer
func (c ScopeCode) ScopeID() int64 { return int64(0) }

// ScopeCode is convenience helper to satisfy the interface ScopeCoder
func (c ScopeCode) ScopeCode() string { return string(c) }

const scopeGroupName = "ScopeAbsentScopeDefaultScopeWebsiteScopeGroupScopeStore"

var scopeGroupIndex = [...]uint8{0, 11, 23, 35, 45, 55}

// String human readable name of ScopeGroup. For Marshaling see ScopePerm
func (i ScopeGroup) String() string {
	if i+1 >= ScopeGroup(len(scopeGroupIndex)) {
		return fmt.Sprintf("ScopeGroup(%d)", i)
	}
	return scopeGroupName[scopeGroupIndex[i]:scopeGroupIndex[i+1]]
}

// ScopeGroupNames returns a slice containing all constant names
func ScopeGroupNames() (r utils.StringSlice) {
	return r.SplitStringer8(scopeGroupName, scopeGroupIndex[:]...)
}

func GetScopeGroup(s string) ScopeGroup {
	switch s {
	case ScopeRangeWebsites:
		return ScopeWebsiteID
	case ScopeRangeStores:
		return ScopeStoreID
	}
	return ScopeDefaultID
}
