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
	IDScopeAbsent ScopeGroup = iota // order of the constants is used for comparison
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

	// ScopeGroup used in constants where default is the lowest and store the highest. Func String() attached.
	// Part of ScopePerm.
	ScopeGroup uint8

	// Retriever implements how to get the ID. If Retriever implements CodeRetriever
	// then CodeRetriever has precedence. ID can be any of the website, group or store IDs.
	Retriever interface {
		ScopeID() int64
	}
	// CodeRetriever implements how to get an object by Code which can be website or store code.
	// Groups doesn't have codes.
	CodeRetriever interface {
		ScopeCode() string
	}
	// ID is convenience helper to satisfy the interface Retriever.
	ScopeID int64
	// Code is convenience helper to satisfy the interface CodeRetriever and Retriever.
	ScopeCode string
)

var _ Retriever = ScopeID(0)
var _ CodeRetriever = ScopeCode("")

// ScopeID is convenience helper to satisfy the interface Retriever
func (i ScopeID) ScopeID() int64 { return int64(i) }

// ScopeID is a noop method receiver to satisfy the interface Retriever
func (c ScopeCode) ScopeID() int64 { return int64(0) }

// ScopeCode is convenience helper to satisfy the interface CodeRetriever
func (c ScopeCode) ScopeCode() string { return string(c) }

const _ScopeGroup_name = "ScopeAbsentScopeDefaultScopeWebsiteScopeGroupScopeStore"

var _ScopeGroup_index = [...]uint8{0, 11, 23, 35, 45, 55}

// String human readable name of ScopeGroup. For Marshaling see ScopePerm
func (i ScopeGroup) String() string {
	if i+1 >= ScopeGroup(len(_ScopeGroup_index)) {
		return fmt.Sprintf("ScopeGroup(%d)", i)
	}
	return _ScopeGroup_name[_ScopeGroup_index[i]:_ScopeGroup_index[i+1]]
}

// ScopeGroupNames returns a slice containing all constant names
func ScopeGroupNames() (r utils.StringSlice) {
	return r.SplitStringer8(_ScopeGroup_name, _ScopeGroup_index[:]...)
}

func GetScopeGroup(s string) ScopeGroup {
	switch s {
	case StringScopeWebsites:
		return IDScopeWebsite
	case StringScopeStores:
		return IDScopeStore
	}
	return IDScopeDefault
}
