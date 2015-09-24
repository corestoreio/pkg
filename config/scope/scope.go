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
	AbsentID Scope = iota // order of the constants is used for comparison
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

	// Scope used in constants where default is the lowest and store the highest. Func String() attached.
	// Part of Perm.
	Scope uint8

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

const scopeGroupName = "ScopeAbsentScopeDefaultScopeWebsiteScopeGroupScopeStore"

var scopeGroupIndex = [...]uint8{0, 11, 23, 35, 45, 55}

// String human readable name of Group. For Marshaling see Perm
func (i Scope) String() string {
	if i+1 >= Scope(len(scopeGroupIndex)) {
		return fmt.Sprintf("ScopeGroup(%d)", i)
	}
	return scopeGroupName[scopeGroupIndex[i]:scopeGroupIndex[i+1]]
}

// GroupNames returns a slice containing all constant names
func GroupNames() (r utils.StringSlice) {
	return r.SplitStringer8(scopeGroupName, scopeGroupIndex[:]...)
}

func FromString(s string) Scope {
	switch s {
	case RangeWebsites:
		return WebsiteID
	case RangeStores:
		return StoreID
	}
	return DefaultID
}
