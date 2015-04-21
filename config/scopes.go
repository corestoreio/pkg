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
	// DataScopeDefault defines the global scope. Stored in table core_config_data.scope.
	DataScopeDefault = "default"
	// DataScopeWebsites defines the website scope which has default as parent and stores as child.
	//  Stored in table core_config_data.scope.
	DataScopeWebsites = "websites"
	// DataScopeStores defines the store scope which has default and websites as parent.
	//  Stored in table core_config_data.scope.
	DataScopeStores = "stores"

	ScopeDefault ScopeID = iota
	ScopeWebsite
	ScopeGroup
	ScopeStore

	PATH_SINGLE_STORE_MODE_ENABLED = "general/single_store_mode/enabled"
)

type (
	// ScopeID used in constants where default is the lowest and store the highest
	ScopeID int

	// ScopePool reads from consul or etcd
	ScopePool interface {
	}

	ScopeReader interface {
		// GetString retrieves a config value by path and scope
		ReadString(path string, scope ScopeID, scopeCode string /*null*/) string

		// IsSetFlag retrieves a config flag by path and scope
		IsSetFlag(path string, scope ScopeID, scopeCode string) bool
	}

	ScopeWriter interface {
		// SetString sets config value in the corresponding config scope
		WriteString(path, value string, scope ScopeID, scopeCode string /*null*/)
	}
)

var (
	cfgDefault = viper.New()
	cfgwebsite = viper.New()
	cfgStore   = viper.New()
)

func init() {
	cfgDefault.SetDefault(PATH_SINGLE_STORE_MODE_ENABLED, false)
	cfgStore.SetDefault(PATH_SINGLE_STORE_MODE_ENABLED, false)
}
