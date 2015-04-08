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

package custattr

import (
	"github.com/corestoreio/csfw/eav"
	"github.com/juju/errgo"
)

type (
	// AttributeSlice implements eav.AttributeSliceGetter @todo website must be present in the slice
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
	// internal wrapper for attribute collection c, getter g and entity type id.
	// internal wrapper to override New() method receiver
	catHandler struct {
		eav.Handler
	}
)

var (
	// verify if interfaces has been implemented
	_ eav.EntityTypeAttributeModeller = (*catHandler)(nil)
	_ eav.EntityAttributeCollectioner = (*catHandler)(nil)

	// aa address attribute
	aa = &catHandler{
		Handler: eav.Handler{},
	}
	// ca customer attribute
	ca = &catHandler{
		Handler: eav.Handler{},
	}
)

// SetAddressCollection requires a slice to set the address attribute collection
func SetAddressCollection(s eav.AttributeSliceGetter) {
	if s.Len() == 0 {
		panic("AttributeSlice is empty")
	}
	aa.C = s
}

// SetAddressGetter knows how to get an attribute using an index out of an attribute collection
func SetAddressGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	aa.G = g
}

// SetCustomerCollection requires a slice to set the customer attribute collection
func SetCustomerCollection(s eav.AttributeSliceGetter) {
	if s.Len() == 0 {
		panic("AttributeSlice is empty")
	}
	ca.C = s
}

// SetCustomerGetter knows how to get an attribute using an index out of an attribute collection
func SetCustomerGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	ca.G = g
}

// ByID returns an custattr.Attributer by int64 id. Use type assertion.
func (s AttributeSlice) ByID(g eav.AttributeGetter, id int64) (interface{}, error) {
	if g == nil {
		panic("AttributeGetter is nil")
	}
	i, err := g.ByID(id)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

// ByCode returns an custattr.Attributer by code. Use type assertion.
func (s AttributeSlice) ByCode(g eav.AttributeGetter, code string) (interface{}, error) {
	if g == nil {
		panic("AttributeGetter is nil")
	}
	i, err := g.ByCode(code)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

// Index returns the current cust.Attributer from index i. Use type assertion.
func (s AttributeSlice) Index(i eav.AttributeIndex) interface{} {
	return s[i]
}

// Len returns the length of a slice
func (s AttributeSlice) Len() int {
	return len(s)
}

func Customer(i int64) *catHandler {
	ca.EntityTyeID = i
	return ca
}

func Address(i int64) *catHandler {
	aa.EntityTyeID = i
	return aa
}

// New creates a new attribute and returns interface custattr.Attributer
// overrides eav.Handler's New() method receiver
func (h *catHandler) New() interface{} {
	return nil
}
