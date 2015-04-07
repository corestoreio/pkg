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

// package catattr handles all product and category related attributes. The name custattr has been chosen
// to be unique so that one can use goimports without conflicts.
package catattr

import (
	"github.com/corestoreio/csfw/eav"
	"github.com/juju/errgo"
)

type (
	// @todo website must be present in the slice
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
	// internal wrapper for attribute collection c, getter g and entity type id.
	container struct {
		entityTyeID int64
		c           AttributeSlice
		g           eav.AttributeGetter
	}
)

var (
	// ca category attribute
	ca = &container{}
	// pa product attribute
	pa = &container{}
	// verify if interfaces has been implemented
	_ eav.EntityTypeAttributeModeller = (*container)(nil)
	_ eav.EntityAttributeCollectioner = (*container)(nil)
)

func SetCategoryCollection(s AttributeSlice) {
	if len(s) == 0 {
		panic("AttributeSlice is empty")
	}
	ca.c = s
}

func SetCategoryGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	ca.g = g
}

func SetProductCollection(s AttributeSlice) {
	if len(s) == 0 {
		panic("AttributeSlice is empty")
	}
	pa.c = s
}

func SetProductGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	pa.g = g
}

func (s AttributeSlice) byID(g eav.AttributeGetter, id int64) (Attributer, error) {
	if g == nil {
		panic("AttributeGetter is nil")
	}
	i, err := g.ByID(id)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

func (s AttributeSlice) byCode(g eav.AttributeGetter, code string) (Attributer, error) {
	if g == nil {
		panic("AttributeGetter is nil")
	}
	i, err := g.ByCode(code)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

func Product(i int64) *container {
	pa.entityTyeID = i
	return pa
}

func Category(i int64) *container {
	ca.entityTyeID = i
	return ca
}

// New creates a new attribute and returns interface custattr.Attributer
func (h *container) New() interface{} {
	return nil
}

// Get uses an AttributeIndex to return an attribute or an error.
// Use type assertion to convert to Attributer.
func (h *container) Get(i eav.AttributeIndex) (interface{}, error) {
	if int(i) < len(h.c) {
		return h.c[i], nil
	}
	return nil, eav.ErrAttributeNotFound
}

// GetByID returns an address attribute by its id
// Use type assertion to convert to Attributer.
func (h *container) GetByID(id int64) (interface{}, error) {
	return h.c.byID(h.g, id)
}

// GetByCode returns an address attribute by its code
// Use type assertion to convert to Attributer.
func (h *container) GetByCode(code string) (interface{}, error) {
	return h.c.byCode(h.g, code)
}

// Collection returns the full attribute collection AttributeSlice.
// You must use type assertion to convert to custattr.AttributeSlice.
func (h *container) Collection() interface{} {
	return h.c
}
