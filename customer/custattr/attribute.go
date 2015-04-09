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
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

type (
	// AttributeSlice implements eav.AttributeSliceGetter @todo website must be present in the slice
	AttributeSlice []Attributer

	// Attributer defines the minimal requirements for a customer attribute. This interface consists
	// of two more tables: customer_eav_attribute and customer_eav_attribute_website. Developers
	// can also extend these tables to add more columns. These columns will be automatically transformed
	// into functions. as soon as there is a scope field and method receiver must cover that scope field.
	Attributer interface {
		eav.Attributer

		InputFilter() string
		Validate() bool // @todo convert php serialize string into a Go type and do only validation here
		IsSystem() bool
		SortOrder() int64
		DataModel() eav.AttributeDataModeller

		// IsVisible must cover also field scopeIsVisible
		IsVisible() bool
		// MultilineCount must also cover field scopeMultilineCount
		MultilineCount() int64
		// IsRequired must also cover field scopeIsRequired
		IsRequired() bool
		// DefaultValue must also cover field scopeDefaultValue
		DefaultValue() string
	}

	// Customer defines attribute properties for a customer and an address. You can use this struct to
	// embed into your own struct for maybe overriding some method receivers.
	Customer struct {
		*eav.Attribute
		isVisible      bool
		inputFilter    string
		multilineCount int64
		validateRules  string
		isSystem       bool
		sortOrder      int64
		dataModel      eav.AttributeDataModeller
		// every scope property must be a dbr.Null type because only nulls are not set otherwise
		// it would override the non-scope value
		scopeIsVisible      dbr.NullBool
		scopeIsRequired     dbr.NullBool
		scopeDefaultValue   dbr.NullString
		scopeMultilineCount dbr.NullInt64
	}

	// internal wrapper for attribute collection c, getter g and entity type id and to override New() method receiver
	catHandler struct {
		eav.Handler
	}
)

var (
	// verify if interfaces has been implemented
	_ eav.EntityTypeAttributeModeller = (*catHandler)(nil)
	_ eav.EntityAttributeCollectioner = (*catHandler)(nil)
	// Check if Attributer interface has been successfully implemented
	_ Attributer = (*Customer)(nil)

	// aa address attribute
	aa = &catHandler{
		Handler: eav.Handler{},
	}
	// ca customer attribute
	ca = &catHandler{
		Handler: eav.Handler{},
	}
)

func (a *Customer) IsVisible() bool {
	if a.scopeIsVisible.Valid {
		return a.scopeIsVisible.Bool
	}
	return a.isVisible
}
func (a *Customer) MultilineCount() int64 {
	if a.scopeMultilineCount.Valid {
		return a.scopeMultilineCount.Int64
	}
	return a.multilineCount
}
func (a *Customer) IsRequired() bool {
	if a.scopeIsRequired.Valid {
		return a.scopeIsRequired.Bool
	}
	return a.Attribute.IsRequired()
}
func (a *Customer) DefaultValue() string {
	if a.scopeDefaultValue.Valid {
		return a.scopeDefaultValue.String
	}
	return a.Attribute.DefaultValue()
}
func (a *Customer) InputFilter() string {
	return a.inputFilter
}
func (a *Customer) Validate() bool {
	return false // a.validateRules
}
func (a *Customer) IsSystem() bool {
	return a.isSystem
}
func (a *Customer) SortOrder() int64 {
	return a.sortOrder
}
func (a *Customer) DataModel() eav.AttributeDataModeller {
	return a.dataModel
}

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
