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
	IDScopeAbsent ScopeID = iota // order of the constants is used for comparison
	IDScopeDefault
	IDScopeWebsite
	IDScopeGroup
	IDScopeStore
)

const (
	// StringScopeDefault defines the global scope. Stored in table core_config_data.scope.
	StringScopeDefault = "default"
	// StringScopeWebsites defines the website scope which has default as parent and stores as child.
	//  Stored in table core_config_data.scope.
	StringScopeWebsites = "websites"
	// StringScopeStores defines the store scope which has default and websites as parent.
	//  Stored in table core_config_data.scope.
	StringScopeStores = "stores"
)

type (

	// ScopeID used in constants where default is the lowest and store the highest. Func String() attached.
	// Part of ScopePerm.
	ScopeID uint8

	// Retriever implements how to get the website or store ID.
	// Duplicated to avoid import cycles. :-(
	Retriever interface {
		ID() int64
	}
)

const _ScopeID_name = "ScopeAbsentScopeDefaultScopeWebsiteScopeGroupScopeStore"

var _ScopeID_index = [...]uint8{0, 11, 23, 35, 45, 55}

// String human readable name of ScopeID. For Marshaling see ScopePerm
func (i ScopeID) String() string {
	if i+1 >= ScopeID(len(_ScopeID_index)) {
		return fmt.Sprintf("ScopeID(%d)", i)
	}
	return _ScopeID_name[_ScopeID_index[i]:_ScopeID_index[i+1]]
}

// ScopeIDNames returns a slice containing all constant names
func ScopeIDNames() (r utils.StringSlice) {
	return r.SplitStringer8(_ScopeID_name, _ScopeID_index[:]...)
}
