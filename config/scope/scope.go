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

	"github.com/corestoreio/csfw/utils"
)

const (
	// *ID defines the overall scopes in a configuration. If a Section/Group/Field
	// can be shown in the current selected scope.
	AbsentID Group = iota // order of the constants is used for comparison
	DefaultID
	WebsiteID
	GroupID
	StoreID
)

const (
	// RangeDefault defines the global scope. Stored in table core_config_data.scope.
	RangeDefault = "default"
	// RangeWebsites defines the website scope which has default as parent and stores as child.
	//  Stored in table core_config_data.scope.
	RangeWebsites = "websites"
	// RangeStores defines the store scope which has default and websites as parent.
	//  Stored in table core_config_data.scope.
	RangeStores = "stores"
)

type (

	// Group used in constants where default is the lowest and store the highest. Func String() attached.
	// Part of Perm.
	Group uint8

	// IDer implements how to get the ID. If IDer implements Coder
	// then Coder has precedence. ID can be any of the website, group or store IDs.
	IDer interface {
		ScopeID() int64
	}
	// Coder implements how to get an object by Code which can be website or store code.
	// Groups doesn't have codes.
	Coder interface {
		ScopeCode() string
	}
	// ID is convenience helper to satisfy the interface IDer.
	ID int64
	// Code is convenience helper to satisfy the interface Coder and IDer.
	Code string
)

var _ IDer = ID(0)
var _ Coder = Code("")

// ScopeID is convenience helper to satisfy the interface IDer
func (i ID) ScopeID() int64 { return int64(i) }

// ScopeID is a noop method receiver to satisfy the interface IDer
func (c Code) ScopeID() int64 { return int64(0) }

// ScopeCode is convenience helper to satisfy the interface Coder
func (c Code) ScopeCode() string { return string(c) }

const scopeGroupName = "ScopeAbsentScopeDefaultScopeWebsiteScopeGroupScopeStore"

var scopeGroupIndex = [...]uint8{0, 11, 23, 35, 45, 55}

// String human readable name of Group. For Marshaling see Perm
func (i Group) String() string {
	if i+1 >= Group(len(scopeGroupIndex)) {
		return fmt.Sprintf("ScopeGroup(%d)", i)
	}
	return scopeGroupName[scopeGroupIndex[i]:scopeGroupIndex[i+1]]
}

// GroupNames returns a slice containing all constant names
func GroupNames() (r utils.StringSlice) {
	return r.SplitStringer8(scopeGroupName, scopeGroupIndex[:]...)
}

func GetGroup(s string) Group {
	switch s {
	case RangeWebsites:
		return WebsiteID
	case RangeStores:
		return StoreID
	}
	return DefaultID
}
