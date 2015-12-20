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

package config

import "github.com/corestoreio/csfw/config/scope"

type (

	// deprecated

	// ModelConstructor implements different fields/functions which can be differently used
	// by the FieldSourceModeller or FieldBackendModeller types.
	// Nearly all functions will return not nil. The Construct() function takes what it needs.
	ModelConstructor struct {
		// WebsiteID contains a website/store ID or nil (=default scope) both can be nil or just one
		ScopeWebsite scope.WebsiteIDer
		ScopeStore   scope.StoreIDer
		// Config returns the configuration Getter and never nil
		Config Getter
		// @todo more fields to be added, depends on the overall requirements of all Magento models.
	}

	// BackendModeller defines how to save and load the data @todo rethink AddData
	// In Magento slang: beforeSave() and afterLoad().
	// The Construct() must be used because NOT all fields of ModelConstructor are available during
	// init process and can of course change during the running app. Also to prevent circular dependencies.
	BackendModeller interface {
		// not sure Construct(ModelConstructor) error
		AddData(interface{})
		Save() error
	}
)
