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

import (
	"github.com/corestoreio/csfw/eav"
	"github.com/juju/errgo"
)

type (
	AttributeSlice []Attributer

	// Attributer defines the minimal requirements for a customer attribute. This interface consists
	// of two more tables: customer_eav_attribute and customer_eav_attribute_website. Developers
	// can also extend these tables to add more columns. These columns will be automatically transformed
	// into functions.
	Attributer interface {
		eav.Attributer

		IsVisible() bool
		InputFilter() string
		MultilineCount() int64
		ValidateRules() string
		IsSystem() bool
		SortOrder() int64
		DataModel() eav.AttributeDataModeller
		IsUsedForCustomerSegment() bool

		ScopeIsVisible() bool
		ScopeIsRequired() bool
		ScopeDefaultValue() string
		ScopeMultilineCount() int64
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
