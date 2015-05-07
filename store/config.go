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

import "github.com/corestoreio/csfw/config"

const (
	// SingleStoreModeEnabled if true then single store mode enabled
	// This flag only shows that admin does not want to show certain
	// UI components in the backend (like store switchers etc)
	// If there is only one store view but it does not check the store view collection. WTF?
	PathSingleStoreModeEnabled = "general/single_store_mode/enabled"
	PathStoreStoreName         = "general/store_information/name"
	PathStoreStorePhone        = "general/store_information/phone"
	PathStoreInUrl             = "web/url/use_store"
	PathUseRewrites            = "web/seo/use_rewrites"

	PathUnsecureBaseUrl = "web/unsecure/base_url"
	PathSecureBaseUrl   = "web/secure/base_url"

	PathSecureInFrontend = "web/secure/use_in_frontend"
	//	PathSecureInAdminhtml      = "web/secure/use_in_adminhtml"

	PathSecureBaseLinkUrl   = "web/secure/base_link_url"
	PathUnsecureBaseLinkUrl = "web/unsecure/base_link_url"

	PathSecureBaseStaticUrl   = "web/secure/base_static_url"
	PathUnsecureBaseStaticUrl = "web/unsecure/base_static_url"

	PathSecureBaseMediaUrl   = "web/secure/base_media_url"
	PathUnsecureBaseMediaUrl = "web/unsecure/base_media_url"

	// This defines the base currency scope ("Currency Setup" > "Currency Options" > "Base Currency").
	// can be 0 = Global or 1 = Website
	PathPriceScope = "catalog/price/scope"

	BaseUrlPlaceholder = "{{base_url}}"
)

// configReader stores the reader. Should not be used. Access it via mustConfig()
var configReader config.ScopeReader

// SetConfig sets the internal variable to the current scope config reader.
// ScopeReader will be used across all functions in this package.
func SetConfigReader(c config.ScopeReader) {
	if c == nil {
		panic("config.ScopeReader cannot be nil")
	}
	configReader = c
}

// mustReadConfig internally used
func mustReadConfig() config.ScopeReader {
	if configReader == nil {
		panic("config.ScopeReader cannot be nil")
	}
	return configReader
}

func GetDefaultConfiguration() config.DefaultMap {
	return config.DefaultMap{
		PathSingleStoreModeEnabled: false,
		PathStoreInUrl:             false,
		PathUnsecureBaseUrl:        BaseUrlPlaceholder,
		PathSecureBaseUrl:          "{{unsecure_base_url}}",
		PathSecureInFrontend:       false,

		PathUnsecureBaseLinkUrl: "{{unsecure_base_url}}", // switch to static
		PathSecureBaseLinkUrl:   "{{secure_base_url}}",

		PathPriceScope: PriceScopeGlobal,
	}
}
