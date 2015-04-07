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

// Package catattr handles all product and category related attributes. The name catattr has been chosen
// to be unique so that one can use goimports without conflicts.
package catattr

import (
	"github.com/corestoreio/csfw/eav"
	"github.com/juju/errgo"
)

type (
	// AttributeSlice implements eav.AttributeSliceGetter @todo website must be present in the slice
	AttributeSlice []Attributer

	// Attributer defines the minimal requirements for a catalog attribute. This interface consists
	// of one more tables: catalog_eav_attribute. Developers can also extend this table to add more columns.
	// These columns will be automatically transformed into more functions.
	Attributer interface {
		eav.Attributer

		FrontendInputRenderer() eav.FrontendInputRendererIFace
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
		// IsConfigurable() bool not used anymore in Magento2
		ApplyTo() string
		IsVisibleInAdvancedSearch() bool
		Position() int64
		IsWysiwygEnabled() bool
		IsUsedForPromoRules() bool
		SearchWeight() int64
	}
	// internal wrapper to override New() method receiver
	catHandler struct {
		eav.Handler
	}
)

var (
	// verify if interfaces has been implemented
	_ eav.EntityTypeAttributeModeller = (*catHandler)(nil)
	_ eav.EntityAttributeCollectioner = (*catHandler)(nil)

	// ca category attribute
	ca = &catHandler{
		Handler: eav.Handler{},
	}
	// pa product attribute
	pa = &catHandler{
		Handler: eav.Handler{},
	}
)

// SetCategoryCollection requires a slice to set the category attribute collection
func SetCategoryCollection(s eav.AttributeSliceGetter) {
	if s.Len() == 0 {
		panic("AttributeSlice is empty")
	}
	ca.C = s
}

// SetCategoryGetter knows how to get an attribute using an index out of an attribute collection
func SetCategoryGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	ca.G = g
}

// SetProductCollection requires a slice to set the product attribute collection
func SetProductCollection(s eav.AttributeSliceGetter) {
	if s.Len() == 0 {
		panic("AttributeSlice is empty")
	}
	pa.C = s
}

// SetProductGetter knows how to get an attribute using an index out of an attribute collection
func SetProductGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	pa.G = g
}

// ByID returns an catattr.Attributer by int64 id. Use type assertion.
func (s AttributeSlice) ByID(g eav.AttributeGetter, id int64) (interface{}, error) {
	if g == nil {
		panic("AttributeGetter is nil")
	}
	i, err := g.ByID(id)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s.Index(i), nil
}

// ByCode returns an catattr.Attributer by code. Use type assertion.
func (s AttributeSlice) ByCode(g eav.AttributeGetter, code string) (interface{}, error) {
	if g == nil {
		panic("AttributeGetter is nil")
	}
	i, err := g.ByCode(code)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s.Index(i), nil
}

// Index returns the current catattr.Attributer from index i. Use type assertion.
func (s AttributeSlice) Index(i eav.AttributeIndex) interface{} {
	return s[i]
}

// Len returns the length of a slice
func (s AttributeSlice) Len() int {
	return len(s)
}

func Product(i int64) *catHandler {
	pa.EntityTyeID = i
	return pa
}

func Category(i int64) *catHandler {
	ca.EntityTyeID = i
	return ca
}

// New creates a new attribute and returns interface catattr.Attributer
// overrides eav.Handler's New() method receiver
func (h *catHandler) New() interface{} {
	return nil
}
