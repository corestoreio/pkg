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
	ScopeAbsent ScopeID = iota // must start from 0 because 0 means not set
	ScopeDefault
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
	URLTypeAbsent URLType = iota
	// UrlTypeWeb defines the ULR type to generate the main base URL.
	URLTypeWeb
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
	ScopeID uint8

	// Retriever implements how to get the ID. If Retriever implements CodeRetriever
	// then CodeRetriever has precedence. ID can be any of the website, group or store IDs.
	// Duplicated to avoid import cycles.
	Retriever interface {
		ID() int64
	}

	Reader interface {
		// ReadString retrieves a config value by path, ScopeID and/or ID
		ReadString(path string, scope ScopeID, r ...Retriever) string

		// IsSetFlag retrieves a config flag by path, ScopeID and/or ID
		IsSetFlag(path string, scope ScopeID, r ...Retriever) bool
	}

	Writer interface {
		// SetString sets config value in the corresponding config scope
		Write(path, value interface{}, scope ScopeID, r ...Retriever)
	}

	// Scope main configuration struct which includes Viper, unhappy with the name Scope
	Scope struct {
		*viper.Viper
	}
)

// NewScope creates the main new configuration for all scopes: default, website and store
func NewScope() *Scope {
	s := &Scope{
		Viper: viper.New(),
	}
	s.SetDefault(PathCSBaseURL, CSBaseURL)
	return s
}

// ApplyDefaults reads the map and applies the keys and values to the default configuration
func (sp *Scope) ApplyDefaults(ss Sectioner) *Scope {
	// mutex necessary?
	ctxLog := logger.WithField("Scope", "ApplyDefaults")
	for k, v := range ss.Defaults() {
		ctxLog.Debug(k, v)
		sp.SetDefault(DataScopeDefault+"/"+k, v)
	}
	return sp
}

const _ScopeID_name = "ScopeAbsentScopeDefaultScopeWebsiteScopeGroupScopeStore"

var _ScopeID_index = [...]uint8{0, 11, 23, 35, 45, 55}

// String human readable name of ScopeID
func (i ScopeID) String() string {
	if i+1 >= ScopeID(len(_ScopeID_index)) {
		return fmt.Sprintf("ScopeID(%d)", i)
	}
	return _ScopeID_name[_ScopeID_index[i]:_ScopeID_index[i+1]]
}

// ScopeIDNames returns a slice containing all constant names
func ScopeIDNames() (r utils.StringSlice) {
	return r.SplitStringer8(_Visible_name, _Visible_index[:]...)
}
