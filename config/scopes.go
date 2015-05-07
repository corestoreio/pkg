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

import "github.com/spf13/viper"

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

	// PS defines the path separator for the configuration path
	PS = "/"

	LeftDelim  = "{{"
	RightDelim = "}}"
)

const (
	// UrlTypeWeb main base url
	UrlTypeWeb UrlType = iota + 1
	// UrlTypeStatic url to the static assets
	UrlTypeStatic
	UrlTypeLink
)

type (
	// UrlType defines the type of the URL. Used in const declaration.
	// @see https://github.com/magento/magento2/blob/0.74.0-beta7/lib/internal/Magento/Framework/UrlInterface.php#L13
	UrlType int

	// ScopeID used in constants where default is the lowest and store the highest
	ScopeID int

	// DefaultMap contains the default aka global configuration of a package
	DefaultMap map[string]interface{}

	// Retriever implements how to get the ID. If Retriever implements CodeRetriever
	// then CodeRetriever has precedence. ID can be any of the website, group or store IDs.
	// Duplicated to avoid import cycles. @todo refactor
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
		WriteString(path, value string, scope ScopeID, r ...Retriever)
	}

	Scope struct {
		Config *viper.Viper
	}
)

func NewScope() *Scope {
	return &Scope{
		Config: viper.New(),
	}
}

// ApplyDefaults reads the map and applies the keys and values to the default configuration
func (sp *Scope) ApplyDefaults(m DefaultMap) *Scope {
	// mutex necessary?
	for k, v := range m {
		sp.Config.SetDefault(DataScopeDefault+PS+k, v)
	}
	return sp
}
