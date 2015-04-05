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

// package custattr handles all customer and address related attributes. The name custattr has been chosen
// to be unique so that one can use goimports without conflicts.
package custattr

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
	addressCollection  AttributeSlice
	addressGetter      eav.AttributeGetter
	customerCollection AttributeSlice
	customerGetter     eav.AttributeGetter
)

func SetAddressCollection(ac AttributeSlice) {
	if len(ac) == 0 {
		panic("AttributeSlice is empty")
	}
	addressCollection = ac
}

func SetCustomerCollection(ac AttributeSlice) {
	if len(ac) == 0 {
		panic("AttributeSlice is empty")
	}
	customerCollection = ac
}

func SetAddressGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	addressGetter = g
}

func SetCustomerGetter(g eav.AttributeGetter) {
	if g == nil {
		panic("AttributeGetter cannot be nil")
	}
	customerGetter = g
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

// GetAddress uses an AttributeIndex to return an attribute or an error.
// One should not modify the attribute object.
func GetAddress(i eav.AttributeIndex) (Attributer, error) {
	if int(i) < len(addressCollection) {
		return addressCollection[i], nil
	}
	return nil, eav.ErrAttributeNotFound
}

// GetCustomer uses an AttributeIndex to return an attribute or an error.
// One should not modify the attribute object.
func GetCustomer(i eav.AttributeIndex) (Attributer, error) {
	if int(i) < len(customerCollection) {
		return customerCollection[i], nil
	}
	return nil, eav.ErrAttributeNotFound
}

// GetAdressByID returns an address attribute by its id
func GetAdressByID(id int64) (Attributer, error) {
	return addressCollection.byID(addressGetter, id)
}

// GetAddressByCode returns an address attribute by its code
func GetAddressByCode(code string) (Attributer, error) {
	return addressCollection.byCode(addressGetter, code)
}

// GetCustomerByID returns a customer attribute by its id
func GetCustomerByID(id int64) (Attributer, error) {
	return customerCollection.byID(customerGetter, id)
}

// GetCustomerByCode returns a customer attribute by its code
func GetCustomerByCode(code string) (Attributer, error) {
	return customerCollection.byCode(customerGetter, code)
}

// GetAddresses returns a copy of the main address attribute slice.
// One should not modify the slice and its content.
func GetAddresses() AttributeSlice {
	return addressCollection
}

// GetCustomers returns a copy of the main customer attribute slice.
// One should not modify the slice and its content.
func GetCustomers() AttributeSlice {
	return customerCollection
}
