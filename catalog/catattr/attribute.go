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

	// Catalog a data container for attributes. You can use this struct to
	// embed into your own struct for maybe overriding some method receivers.
	Catalog struct {
		*eav.Attribute
		frontendInputRenderer     eav.FrontendInputRendererIFace
		isGlobal                  bool
		isVisible                 bool
		isSearchable              bool
		isFilterable              bool
		isComparable              bool
		isVisibleOnFront          bool
		isHtmlAllowedOnFront      bool
		isUsedForPriceRules       bool
		isFilterableInSearch      bool
		usedInProductListing      bool
		usedForSortBy             bool
		isConfigurable            bool
		applyTo                   string
		isVisibleInAdvancedSearch bool
		position                  int64
		isWysiwygEnabled          bool
		isUsedForPromoRules       bool
		searchWeight              int64
	}

	// internal wrapper to override New() method receiver
	catHandler struct {
		eav.Handler
	}
)

var (
	// verify if interfaces has been implemented
	_ eav.EntityTypeAttributeModeller     = (*catHandler)(nil)
	_ eav.EntityTypeAttributeCollectioner = (*catHandler)(nil)

	// ca category attribute
	ca = &catHandler{
		Handler: eav.Handler{},
	}
	// pa product attribute
	pa = &catHandler{
		Handler: eav.Handler{},
	}
	// Check if Attributer interface has been successfully implemented
	_ Attributer = (*Catalog)(nil)
)

func (a *Catalog) FrontendInputRenderer() eav.FrontendInputRendererIFace {
	return a.frontendInputRenderer
}
func (a *Catalog) IsGlobal() bool {
	return a.isGlobal
}
func (a *Catalog) IsVisible() bool {
	return a.isVisible
}
func (a *Catalog) IsSearchable() bool {
	return a.isSearchable
}
func (a *Catalog) IsFilterable() bool {
	return a.isFilterable
}
func (a *Catalog) IsComparable() bool {
	return a.isComparable
}
func (a *Catalog) IsVisibleOnFront() bool {
	return a.isVisibleOnFront
}
func (a *Catalog) IsHtmlAllowedOnFront() bool {
	return a.isHtmlAllowedOnFront
}
func (a *Catalog) IsUsedForPriceRules() bool {
	return a.isUsedForPriceRules
}
func (a *Catalog) IsFilterableInSearch() bool {
	return a.isFilterableInSearch
}
func (a *Catalog) UsedInProductListing() bool {
	return a.usedInProductListing
}
func (a *Catalog) UsedForSortBy() bool {
	return a.usedForSortBy
}
func (a *Catalog) IsConfigurable() bool {
	return a.isConfigurable
}
func (a *Catalog) ApplyTo() string {
	return a.applyTo
}
func (a *Catalog) IsVisibleInAdvancedSearch() bool {
	return a.isVisibleInAdvancedSearch
}
func (a *Catalog) Position() int64 {
	return a.position
}
func (a *Catalog) IsWysiwygEnabled() bool {
	return a.isWysiwygEnabled
}
func (a *Catalog) IsUsedForPromoRules() bool {
	return a.isUsedForPromoRules
}
func (a *Catalog) SearchWeight() int64 {
	return a.searchWeight
}

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
