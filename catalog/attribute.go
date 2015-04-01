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

package catalog

import (
	"github.com/corestoreio/csfw/eav"
)

type (
	indexCSCategoryAttribute int

	// Attributer defines the minimum fields needed.
	// Custom struct types will use this interface for embedding.
	// materializer creates all required methods plus of course the additional definied ones
	// then generate an empty struct which has the values coded in it!
	/*
		type myCategoryAttributeAllChildren struct{}
		func (myCategoryAttributeAllChildren) AttributeID() int64 { return 2 }
		func (myCategoryAttributeAllChildren) EntityTypeID() int64 { return 2 }
	*/

	// CSCategoryAttributeSlice contains pointers to CSCategoryAttribute types
	// @todo website must be present in the slice
	AttributeSlicer []Attributer

	// use eav.Attributer as embedded ... because eav_attribute :-)
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
		FrontendInputRenderer() string
		IsGlobal() bool
		IsVisible() bool
		IsSearchable() bool
		IsFilterable() bool
		IsComparable() bool
		IsVisibleOnFront() bool
		IsHtmlAllowedOnFront() bool
		IsUsedForPriceRules() bool
		IsFilterableInSearch() bool
		UsedInProductListing() bool
		UsedForSortBy() bool
		IsConfigurable() bool
		ApplyTo() string
		IsVisibleInAdvancedSearch() bool
		Position() int64
		IsWysiwygEnabled() bool
		IsUsedForPromoRules() bool
		SearchWeight() int64
	}
)
