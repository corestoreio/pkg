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

package store

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
)

const (
	// SingleStoreModeEnabled if true then single store mode enabled
	// This flag only shows that admin does not want to show certain
	// UI components in the backend (like store switchers etc)
	// If there is only one store view but it does not check the store view collection. WTF?
	PathSingleStoreModeEnabled = "general/single_store_mode/enabled"
	PathStoreStoreName         = "general/store_information/name"
	PathStoreStorePhone        = "general/store_information/phone"
	PathStoreInURL             = "web/url/use_store"
	PathUseRewrites            = "web/seo/use_rewrites"
	PathSecureInFrontend       = "web/secure/use_in_frontend"

	PathUnsecureBaseURL = "web/unsecure/base_url"
	PathSecureBaseURL   = "web/secure/base_url"

	//	PathSecureBaseLinkUrl   = "web/secure/base_link_url"
	//	PathUnsecureBaseLinkUrl = "web/unsecure/base_link_url"

	PathSecureBaseStaticURL   = "web/secure/base_static_url"
	PathUnsecureBaseStaticURL = "web/unsecure/base_static_url"

	PathSecureBaseMediaURL   = "web/secure/base_media_url"
	PathUnsecureBaseMediaURL = "web/unsecure/base_media_url"

	// This defines the base currency scope ("Currency Setup" > "Currency Options" > "Base Currency").
	// can be 0 = Global or 1 = Website
	PathPriceScope = "catalog/price/scope"

	PlaceholderBaseURL         = config.LeftDelim + "base_url" + config.RightDelim
	PlaceholderBaseURLSecure   = config.LeftDelim + "secure_base_url" + config.RightDelim
	PlaceholderBaseURLUnSecure = config.LeftDelim + "unsecure_base_url" + config.RightDelim
)

var (
	// configReader stores the reader. Should not be used. Access it via mustConfig()
	configReader config.Reader
	// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
	TableCollection csdb.TableStructureSlice
)

// SetConfig sets the internal variable to the current scope config reader.
// ScopeReader will be used across all functions in this package.
func SetConfigReader(c config.Reader) {
	if c == nil || configReader != nil {
		panic("config.ScopeReader cannot be nil or already set")
	}
	configReader = c
}

// mustReadConfig internally used
func mustReadConfig() config.Reader {
	if configReader == nil {
		panic("config.ScopeReader cannot be nil")
	}
	return configReader
}

// GetDefaultConfiguration in conjunction with config.Scope.ApplyDefaults function to
// set the default configuration value for a package.
func GetDefaultConfiguration() config.DefaultMap {
	return config.DefaultMap{
		PathSingleStoreModeEnabled: false,
		PathStoreInURL:             false,
		PathUnsecureBaseURL:        PlaceholderBaseURL,
		PathSecureBaseURL:          PlaceholderBaseURLUnSecure,
		PathSecureInFrontend:       false,

		PathPriceScope: PriceScopeGlobal,
	}
}
