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
	"github.com/spf13/viper"
)

const (
	ScopeDefault ScopeID = iota + 1
	ScopeWebsite
	ScopeGroup
	ScopeStore
)

const (
	// DataScopeDefault defines the global scope. Stored in table core_config_data.scope.
	DataScopeDefault = "default"
	// DataScopeWebsites defines the website scope which has default as parent and stores as child.
	//  Stored in table core_config_data.scope.
	DataScopeWebsites = "websites"
	// DataScopeStores defines the store scope which has default and websites as parent.
	//  Stored in table core_config_data.scope.
	DataScopeStores = "stores"

	LeftDelim  = "{{"
	RightDelim = "}}"

	CSBaseURL     = "http://localhost:9500/"
	PathCSBaseURL = "web/corestore/base_url"
)

const (
	// UrlTypeWeb defines the ULR type to generate the main base URL.
	URLTypeWeb URLType = iota + 1
	// UrlTypeStatic defines the url to the static assets like css, js or theme images
	URLTypeStatic
	// UrlTypeLink hmmm
	// UrlTypeLink
	// UrlTypeMedia defines the ULR type for generating URLs to product photos
	URLTypeMedia
)

type (
	// UrlType defines the type of the URL. Used in const declaration.
	// @see https://github.com/magento/magento2/blob/0.74.0-beta7/lib/internal/Magento/Framework/UrlInterface.php#L13
	URLType int

	// ScopeID used in constants where default is the lowest and store the highest. Func String() attached
	ScopeID uint
	// ScopeBits mostly used for permissions, ScopeGroup is not a part of this bit set.
	ScopeBits uint64

	// DefaultMap contains the default aka global configuration of a package
	DefaultMap map[string]interface{}

	// Retriever implements how to get the ID. If Retriever implements CodeRetriever
	// then CodeRetriever has precedence. ID can be any of the website, group or store IDs.
	// Duplicated to avoid import cycles.
	Retriever interface {
		ID() int64
	}

	ScopeReader interface {
		// ReadString retrieves a config value by path, ScopeID and/or ID
		ReadString(path string, scope ScopeID, r ...Retriever) string

		// IsSetFlag retrieves a config flag by path, ScopeID and/or ID
		IsSetFlag(path string, scope ScopeID, r ...Retriever) bool
	}

	ScopeWriter interface {
		// SetString sets config value in the corresponding config scope
		Write(path, value interface{}, scope ScopeID, r ...Retriever)
	}

	// Scope main configuration struct which includes Viper
	Scope struct {
		*viper.Viper
	}
)

// AllScopes convenient helper variable contains all scope permission levels
var AllScopes = ScopeBits(0).All()

// NewScope creates the main new configuration for all scopes: default, website and store
func NewScope() *Scope {
	s := &Scope{
		Viper: viper.New(),
	}
	s.SetDefault(PathCSBaseURL, CSBaseURL)
	return s
}

// ApplyDefaults reads the map and applies the keys and values to the default configuration
func (sp *Scope) ApplyDefaults(m DefaultMap) *Scope {
	// mutex necessary?
	for k, v := range m {
		sp.SetDefault(DataScopeDefault+"/"+k, v)
	}
	return sp
}

// All applies all scopes
func (bits *ScopeBits) All() ScopeBits {
	bits.Set(ScopeDefault, ScopeWebsite, ScopeStore)
	return *bits
}

// Set takes a variadic amount of ScopeID to set them to ScopeBits
func (bits *ScopeBits) Set(scopes ...ScopeID) ScopeBits {
	for _, i := range scopes {
		*bits = *bits | (1 << i) // (1 << power = 2^power)
	}
	return *bits
}

// Has checks if ScopeID is in ScopeBits
func (bits ScopeBits) Has(s ScopeID) bool {
	var one ScopeID = 1
	return (bits & ScopeBits(one<<s)) != 0
}

// Human
func (bits ScopeBits) Human() utils.StringSlice {
	var ret utils.StringSlice
	var i uint
	for i = 0; i < 64; i++ {
		bit := ((bits & (1 << i)) != 0)
		if bit {
			ret.Append(ScopeID(i).String())
		}
	}
	return ret
}

const _ScopeID_name = "ScopeDefaultScopeWebsiteScopeGroupScopeStore"

var _ScopeID_index = [...]uint8{12, 24, 34, 44}

func (i ScopeID) String() string {
	i -= 1
	if i >= ScopeID(len(_ScopeID_index)) {
		return fmt.Sprintf("ScopeID(%d)", i+1)
	}
	hi := _ScopeID_index[i]
	lo := uint8(0)
	if i > 0 {
		lo = _ScopeID_index[i-1]
	}
	return _ScopeID_name[lo:hi]
}
