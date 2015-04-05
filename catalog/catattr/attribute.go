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
)

var (
	categoryCollection AttributeSlice
	categoryGetter     eav.AttributeGetter
	productCollection  AttributeSlice
	productGetter      eav.AttributeGetter
)

func SetCategoryCollection(ac AttributeSlice) {
	if len(ac) == 0 {
		panic("AttributeSlice is empty")
	}
	categoryCollection = ac
}

func SetProductCollection(ac AttributeSlice) {
	if len(ac) == 0 {
		panic("AttributeSlice is empty")
	}
	productCollection = ac
}

func SetCategoryGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	categoryGetter = g
}

func SetProductGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	productGetter = g
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

// GetCategory uses an AttributeIndex to return an attribute or an error.
// One should not modify the attribute object.
func GetCategory(i eav.AttributeIndex) (Attributer, error) {
	if int(i) < len(categoryCollection) {
		return categoryCollection[i], nil
	}
	return nil, eav.ErrAttributeNotFound
}

// GetProduct uses an AttributeIndex to return an attribute or an error.
// One should not modify the attribute object.
func GetProduct(i eav.AttributeIndex) (Attributer, error) {
	if int(i) < len(productCollection) {
		return productCollection[i], nil
	}
	return nil, eav.ErrAttributeNotFound
}

// GetCategoryByID returns an category attribute by its id
func GetCategoryByID(id int64) (Attributer, error) {
	return categoryCollection.byID(categoryGetter, id)
}

// GetCategoryByCode returns an category attribute by its code
func GetCategoryByCode(code string) (Attributer, error) {
	return categoryCollection.byCode(categoryGetter, code)
}

// GetProductByID returns a product attribute by its id
func GetProductByID(id int64) (Attributer, error) {
	return productCollection.byID(productGetter, id)
}

// GetProductByCode returns a product attribute by its code
func GetProductByCode(code string) (Attributer, error) {
	return productCollection.byCode(productGetter, code)
}

// GetCategories returns a copy of the main category attribute slice.
// One should not modify the slice and its content.
func GetCategories() AttributeSlice {
	return categoryCollection
}

// GetProducts returns a copy of the main product attribute slice.
// One should not modify the slice and its content.
func GetProducts() AttributeSlice {
	return productCollection
}
