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

type (
	Option struct {
		Value, Label string
	}

	// Element @see magento2/app/code/Magento/Config/etc/system_file.xsd @todo add later more fields IF necessary
	Element struct {
		Label string
		// Show: eg: ScopeDefault & ScopeWebsite & ScopeDefault: showInDefault="1" showInWebsite="1" showInStore="1"
		Show       uint
		SortOrder  int
		ConfigPath string
		// SourceModel defines how to retrieve all option values
		SourceModel interface {
			Options() []Option
		}
		// BackendModel defines @todo think about AddData
		BackendModel interface {
			AddData(interface{})
			Save() error
		}
	}
	ElementSlice []*Element
)
