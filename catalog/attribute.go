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
	attributeCollection AttributeSlice
	attributeGetter     eav.AttributeGetter
)

func SetAttributeCollection(ac AttributeSlice) {
	if len(ac) == 0 {
		panic("AttributeSlice is empty")
	}
	attributeCollection = ac
}

func SetAttributeGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	attributeGetter = g
}

func (s AttributeSlice) ByID(id int64) (Attributer, error) {
	i, err := attributeGetter.ByID(id)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

func (s AttributeSlice) ByCode(code string) (Attributer, error) {
	i, err := attributeGetter.ByCode(code)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

// GetAttribute uses an AttributeIndex to return a attribute or an error.
// One should not modify the attribute object.
func GetAttribute(i eav.AttributeIndex) (Attributer, error) {
	if int(i) < len(attributeCollection) {
		return attributeCollection[i], nil
	}
	return nil, eav.ErrAttributeNotFound
}

func GetAttributeByID(id int64) (Attributer, error) {
	return attributeCollection.ByID(id)
}

func GetAttributeByCode(code string) (Attributer, error) {
	return attributeCollection.ByCode(code)
}

// GetAttributes returns a copy of the main slice of attributes.
// One should not modify the slice and its content.
func GetAttributes() AttributeSlice {
	return attributeCollection
}
