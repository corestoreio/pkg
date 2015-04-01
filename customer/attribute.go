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

package customer

import "github.com/corestoreio/csfw/eav"

type (
	indexCSAddressAttribute int

	AttributeSlicer []Attributer

	Attributer interface {
		AttributeID() int64
		EntityTypeID() int64
		AttributeCode() string
		AttributeModel() string
		BackendModel() eav.AttributeBackendModeller
		BackendType() string
		BackendTable() string
		FrontendModel() eav.AttributeFrontendModeller
		FrontendInput() string
		FrontendLabel() string
		FrontendClass() string
		SourceModel() eav.AttributeSourceModeller
		IsRequired() bool
		IsUserDefined() bool
		DefaultValue() string
		IsUnique() bool
		Note() string
		IsVisible() bool
		InputFilter() string
		MultilineCount() int64
		ValidateRules() string
		IsSystem() bool
		SortOrder() int64
		DataModel() string
		IsUsedForCustomerSegment() bool
		ScopeIsVisible() bool
		ScopeIsRequired() bool
		ScopeDefaultValue() string
		ScopeMultilineCount() int64
	}
)
