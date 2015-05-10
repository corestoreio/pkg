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

// Package custattr_test with address attribute data contains test data which has
// been previously auto generated. So do not touch 8-)
package custattr_test

import (
	"github.com/corestoreio/csfw/customer/custattr"
	"github.com/corestoreio/csfw/eav"
)

type (
	// @todo website must be present in the slice
	// CustomerAddressAttribute a data container for attributes. You can use this struct to
	// embed into your own struct for maybe overriding some method receivers.
	customerAddressAttribute struct {
		*eav.Attribute
		isVisible                bool
		inputFilter              string
		multilineCount           int64
		validateRules            string
		isSystem                 bool
		sortOrder                int64
		dataModel                eav.AttributeDataModeller
		isUsedForCustomerSegment bool
		scopeIsVisible           bool
		scopeIsRequired          bool
		scopeDefaultValue        string
		scopeMultilineCount      int64
	}
)

func (a *customerAddressAttribute) IsVisible() bool {
	return a.isVisible
}
func (a *customerAddressAttribute) InputFilter() string {
	return a.inputFilter
}
func (a *customerAddressAttribute) MultilineCount() int64 {
	return a.multilineCount
}
func (a *customerAddressAttribute) ValidateRules() string {
	return a.validateRules
}
func (a *customerAddressAttribute) IsSystem() bool {
	return a.isSystem
}
func (a *customerAddressAttribute) SortOrder() int64 {
	return a.sortOrder
}
func (a *customerAddressAttribute) DataModel() eav.AttributeDataModeller {
	return a.dataModel
}
func (a *customerAddressAttribute) IsUsedForCustomerSegment() bool {
	return a.isUsedForCustomerSegment
}
func (a *customerAddressAttribute) ScopeIsVisible() bool {
	return a.scopeIsVisible
}
func (a *customerAddressAttribute) ScopeIsRequired() bool {
	return a.scopeIsRequired
}
func (a *customerAddressAttribute) ScopeDefaultValue() string {
	return a.scopeDefaultValue
}
func (a *customerAddressAttribute) ScopeMultilineCount() int64 {
	return a.scopeMultilineCount
}
func (a *customerAddressAttribute) Validate() bool {
	return false
}

// Check if Attributer interface has been successfully implemented
var _ custattr.Attributer = (*customerAddressAttribute)(nil)

const (
	CustomerAddressAttributeCity eav.AttributeIndex = iota
	CustomerAddressAttributeCompany
	CustomerAddressAttributeCountryID
	CustomerAddressAttributeFax
	CustomerAddressAttributeFirstname
	CustomerAddressAttributeLastname
	CustomerAddressAttributeMiddlename
	CustomerAddressAttributePostcode
	CustomerAddressAttributePrefix
	CustomerAddressAttributeRegion
	CustomerAddressAttributeRegionID
	CustomerAddressAttributeStreet
	CustomerAddressAttributeSuffix
	CustomerAddressAttributeTelephone
	CustomerAddressAttributeVatID
	CustomerAddressAttributeVatIsValid
	CustomerAddressAttributeVatRequestDate
	CustomerAddressAttributeVatRequestID
	CustomerAddressAttributeVatRequestSuccess

	CustomerAddressAttributeZZZ
)

type siCustomerAddressAttribute struct{}

func (siCustomerAddressAttribute) ByID(id int64) (eav.AttributeIndex, error) {
	switch id {
	case 26:
		return CustomerAddressAttributeCity, nil
	case 24:
		return CustomerAddressAttributeCompany, nil
	case 27:
		return CustomerAddressAttributeCountryID, nil
	case 32:
		return CustomerAddressAttributeFax, nil
	case 20:
		return CustomerAddressAttributeFirstname, nil
	case 22:
		return CustomerAddressAttributeLastname, nil
	case 21:
		return CustomerAddressAttributeMiddlename, nil
	case 30:
		return CustomerAddressAttributePostcode, nil
	case 19:
		return CustomerAddressAttributePrefix, nil
	case 28:
		return CustomerAddressAttributeRegion, nil
	case 29:
		return CustomerAddressAttributeRegionID, nil
	case 25:
		return CustomerAddressAttributeStreet, nil
	case 23:
		return CustomerAddressAttributeSuffix, nil
	case 31:
		return CustomerAddressAttributeTelephone, nil
	case 36:
		return CustomerAddressAttributeVatID, nil
	case 37:
		return CustomerAddressAttributeVatIsValid, nil
	case 39:
		return CustomerAddressAttributeVatRequestDate, nil
	case 38:
		return CustomerAddressAttributeVatRequestID, nil
	case 40:
		return CustomerAddressAttributeVatRequestSuccess, nil

	default:
		return eav.AttributeIndex(0), eav.ErrAttributeNotFound
	}
}

func (siCustomerAddressAttribute) ByCode(code string) (eav.AttributeIndex, error) {
	switch code {
	case "city":
		return CustomerAddressAttributeCity, nil
	case "company":
		return CustomerAddressAttributeCompany, nil
	case "country_id":
		return CustomerAddressAttributeCountryID, nil
	case "fax":
		return CustomerAddressAttributeFax, nil
	case "firstname":
		return CustomerAddressAttributeFirstname, nil
	case "lastname":
		return CustomerAddressAttributeLastname, nil
	case "middlename":
		return CustomerAddressAttributeMiddlename, nil
	case "postcode":
		return CustomerAddressAttributePostcode, nil
	case "prefix":
		return CustomerAddressAttributePrefix, nil
	case "region":
		return CustomerAddressAttributeRegion, nil
	case "region_id":
		return CustomerAddressAttributeRegionID, nil
	case "street":
		return CustomerAddressAttributeStreet, nil
	case "suffix":
		return CustomerAddressAttributeSuffix, nil
	case "telephone":
		return CustomerAddressAttributeTelephone, nil
	case "vat_id":
		return CustomerAddressAttributeVatID, nil
	case "vat_is_valid":
		return CustomerAddressAttributeVatIsValid, nil
	case "vat_request_date":
		return CustomerAddressAttributeVatRequestDate, nil
	case "vat_request_id":
		return CustomerAddressAttributeVatRequestID, nil
	case "vat_request_success":
		return CustomerAddressAttributeVatRequestSuccess, nil

	default:
		return eav.AttributeIndex(0), eav.ErrAttributeNotFound
	}
}

var _ eav.AttributeGetter = (*siCustomerAddressAttribute)(nil)

func init() {
	custattr.SetAddressGetter(siCustomerAddressAttribute{})
	custattr.SetAddressCollection(custattr.AttributeSlice{

		CustomerAddressAttributeCity: &customerAddressAttribute{
			Attribute: eav.NewAttribute("city", // attribute_code
				26,        // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				2,         // entity_type_id
				"",        // frontend_class
				"text",    // frontend_input
				"City",    // frontend_label
				nil,       // frontend_model
				true,      // is_required
				false,     // is_unique
				false,     // is_user_defined
				"",        // note
				nil,       // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: true,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                80,
			validateRules:            "a:2:{s:15:\"max_text_length\";i:255;s:15:\"min_text_length\";i:1;}",
		},

		CustomerAddressAttributeCompany: &customerAddressAttribute{
			Attribute: eav.NewAttribute("company", // attribute_code
				24,        // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				2,         // entity_type_id
				"",        // frontend_class
				"text",    // frontend_input
				"Company", // frontend_label
				nil,       // frontend_model
				false,     // is_required
				false,     // is_unique
				false,     // is_user_defined
				"",        // note
				nil,       // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: true,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                60,
			validateRules:            "a:2:{s:15:\"max_text_length\";i:255;s:15:\"min_text_length\";i:1;}",
		},

		CustomerAddressAttributeCountryID: &customerAddressAttribute{
			Attribute: eav.NewAttribute("country_id", // attribute_code
				27,        // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				2,         // entity_type_id
				"",        // frontend_class
				"select",  // frontend_input
				"Country", // frontend_label
				nil,       // frontend_model
				true,      // is_required
				false,     // is_unique
				false,     // is_user_defined
				"",        // note
				custattr.AddressSourceCountry().Config(eav.AttributeSourceIdx(CustomerAddressAttributeCountryID)), // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: true,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                90,
			validateRules:            "",
		},

		CustomerAddressAttributeFax: &customerAddressAttribute{
			Attribute: eav.NewAttribute("fax", // attribute_code
				32,        // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				2,         // entity_type_id
				"",        // frontend_class
				"text",    // frontend_input
				"Fax",     // frontend_label
				nil,       // frontend_model
				false,     // is_required
				false,     // is_unique
				false,     // is_user_defined
				"",        // note
				nil,       // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: false,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                130,
			validateRules:            "a:2:{s:15:\"max_text_length\";i:255;s:15:\"min_text_length\";i:1;}",
		},

		CustomerAddressAttributeFirstname: &customerAddressAttribute{
			Attribute: eav.NewAttribute("firstname", // attribute_code
				20,           // attribute_id
				nil,          // backend_model
				"",           // backend_table
				"varchar",    // backend_type
				"",           // default_value
				2,            // entity_type_id
				"",           // frontend_class
				"text",       // frontend_input
				"First Name", // frontend_label
				nil,          // frontend_model
				true,         // is_required
				false,        // is_unique
				false,        // is_user_defined
				"",           // note
				nil,          // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: true,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                20,
			validateRules:            "a:2:{s:15:\"max_text_length\";i:255;s:15:\"min_text_length\";i:1;}",
		},

		CustomerAddressAttributeLastname: &customerAddressAttribute{
			Attribute: eav.NewAttribute("lastname", // attribute_code
				22,          // attribute_id
				nil,         // backend_model
				"",          // backend_table
				"varchar",   // backend_type
				"",          // default_value
				2,           // entity_type_id
				"",          // frontend_class
				"text",      // frontend_input
				"Last Name", // frontend_label
				nil,         // frontend_model
				true,        // is_required
				false,       // is_unique
				false,       // is_user_defined
				"",          // note
				nil,         // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: true,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                40,
			validateRules:            "a:2:{s:15:\"max_text_length\";i:255;s:15:\"min_text_length\";i:1;}",
		},

		CustomerAddressAttributeMiddlename: &customerAddressAttribute{
			Attribute: eav.NewAttribute("middlename", // attribute_code
				21,                    // attribute_id
				nil,                   // backend_model
				"",                    // backend_table
				"varchar",             // backend_type
				"",                    // default_value
				2,                     // entity_type_id
				"",                    // frontend_class
				"text",                // frontend_input
				"Middle Name/Initial", // frontend_label
				nil,   // frontend_model
				false, // is_required
				false, // is_unique
				false, // is_user_defined
				"",    // note
				nil,   // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 false,
			isUsedForCustomerSegment: false,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                30,
			validateRules:            "",
		},

		CustomerAddressAttributePostcode: &customerAddressAttribute{
			Attribute: eav.NewAttribute("postcode", // attribute_code
				30,                // attribute_id
				nil,               // backend_model
				"",                // backend_table
				"varchar",         // backend_type
				"",                // default_value
				2,                 // entity_type_id
				"",                // frontend_class
				"text",            // frontend_input
				"Zip/Postal Code", // frontend_label
				nil,               // frontend_model
				true,              // is_required
				false,             // is_unique
				false,             // is_user_defined
				"",                // note
				nil,               // source_model

			),
			dataModel:                custattr.AddressDataPostcode().Config(eav.AttributeDataIdx(CustomerAddressAttributePostcode)),
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: true,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                110,
			validateRules:            "a:0:{}",
		},

		CustomerAddressAttributePrefix: &customerAddressAttribute{
			Attribute: eav.NewAttribute("prefix", // attribute_code
				19,        // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				2,         // entity_type_id
				"",        // frontend_class
				"text",    // frontend_input
				"Prefix",  // frontend_label
				nil,       // frontend_model
				false,     // is_required
				false,     // is_unique
				false,     // is_user_defined
				"",        // note
				nil,       // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 false,
			isUsedForCustomerSegment: false,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                10,
			validateRules:            "",
		},

		CustomerAddressAttributeRegion: &customerAddressAttribute{
			Attribute: eav.NewAttribute("region", // attribute_code
				28, // attribute_id
				custattr.AddressBackendRegion().Config(eav.AttributeBackendIdx(CustomerAddressAttributeRegion)), // backend_model
				"",               // backend_table
				"varchar",        // backend_type
				"",               // default_value
				2,                // entity_type_id
				"",               // frontend_class
				"text",           // frontend_input
				"State/Province", // frontend_label
				nil,              // frontend_model
				false,            // is_required
				false,            // is_unique
				false,            // is_user_defined
				"",               // note
				nil,              // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: false,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                100,
			validateRules:            "",
		},

		CustomerAddressAttributeRegionID: &customerAddressAttribute{
			Attribute: eav.NewAttribute("region_id", // attribute_code
				29,               // attribute_id
				nil,              // backend_model
				"",               // backend_table
				"int",            // backend_type
				"",               // default_value
				2,                // entity_type_id
				"",               // frontend_class
				"hidden",         // frontend_input
				"State/Province", // frontend_label
				nil,              // frontend_model
				false,            // is_required
				false,            // is_unique
				false,            // is_user_defined
				"",               // note
				custattr.AddressSourceRegion().Config(eav.AttributeSourceIdx(CustomerAddressAttributeRegionID)), // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: true,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                100,
			validateRules:            "",
		},

		CustomerAddressAttributeStreet: &customerAddressAttribute{
			Attribute: eav.NewAttribute("street", // attribute_code
				25, // attribute_id
				custattr.AddressBackendStreet().Config(eav.AttributeBackendIdx(CustomerAddressAttributeStreet)), // backend_model
				"",               // backend_table
				"text",           // backend_type
				"",               // default_value
				2,                // entity_type_id
				"",               // frontend_class
				"multiline",      // frontend_input
				"Street Address", // frontend_label
				nil,              // frontend_model
				true,             // is_required
				false,            // is_unique
				false,            // is_user_defined
				"",               // note
				nil,              // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: true,
			isVisible:                true,
			multilineCount:           2,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                70,
			validateRules:            "a:2:{s:15:\"max_text_length\";i:255;s:15:\"min_text_length\";i:1;}",
		},

		CustomerAddressAttributeSuffix: &customerAddressAttribute{
			Attribute: eav.NewAttribute("suffix", // attribute_code
				23,        // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				2,         // entity_type_id
				"",        // frontend_class
				"text",    // frontend_input
				"Suffix",  // frontend_label
				nil,       // frontend_model
				false,     // is_required
				false,     // is_unique
				false,     // is_user_defined
				"",        // note
				nil,       // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 false,
			isUsedForCustomerSegment: false,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                50,
			validateRules:            "",
		},

		CustomerAddressAttributeTelephone: &customerAddressAttribute{
			Attribute: eav.NewAttribute("telephone", // attribute_code
				31,          // attribute_id
				nil,         // backend_model
				"",          // backend_table
				"varchar",   // backend_type
				"",          // default_value
				2,           // entity_type_id
				"",          // frontend_class
				"text",      // frontend_input
				"Telephone", // frontend_label
				nil,         // frontend_model
				true,        // is_required
				false,       // is_unique
				false,       // is_user_defined
				"",          // note
				nil,         // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: true,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                120,
			validateRules:            "a:2:{s:15:\"max_text_length\";i:255;s:15:\"min_text_length\";i:1;}",
		},

		CustomerAddressAttributeVatID: &customerAddressAttribute{
			Attribute: eav.NewAttribute("vat_id", // attribute_code
				36,           // attribute_id
				nil,          // backend_model
				"",           // backend_table
				"varchar",    // backend_type
				"",           // default_value
				2,            // entity_type_id
				"",           // frontend_class
				"text",       // frontend_input
				"VAT number", // frontend_label
				nil,          // frontend_model
				false,        // is_required
				false,        // is_unique
				false,        // is_user_defined
				"",           // note
				nil,          // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: false,
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                140,
			validateRules:            "",
		},

		CustomerAddressAttributeVatIsValid: &customerAddressAttribute{
			Attribute: eav.NewAttribute("vat_is_valid", // attribute_code
				37,                    // attribute_id
				nil,                   // backend_model
				"",                    // backend_table
				"int",                 // backend_type
				"",                    // default_value
				2,                     // entity_type_id
				"",                    // frontend_class
				"text",                // frontend_input
				"VAT number validity", // frontend_label
				nil,   // frontend_model
				false, // is_required
				false, // is_unique
				false, // is_user_defined
				"",    // note
				nil,   // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: false,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},

		CustomerAddressAttributeVatRequestDate: &customerAddressAttribute{
			Attribute: eav.NewAttribute("vat_request_date", // attribute_code
				39,        // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				2,         // entity_type_id
				"",        // frontend_class
				"text",    // frontend_input
				"VAT number validation request date", // frontend_label
				nil,   // frontend_model
				false, // is_required
				false, // is_unique
				false, // is_user_defined
				"",    // note
				nil,   // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: false,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},

		CustomerAddressAttributeVatRequestID: &customerAddressAttribute{
			Attribute: eav.NewAttribute("vat_request_id", // attribute_code
				38,        // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				2,         // entity_type_id
				"",        // frontend_class
				"text",    // frontend_input
				"VAT number validation request ID", // frontend_label
				nil,   // frontend_model
				false, // is_required
				false, // is_unique
				false, // is_user_defined
				"",    // note
				nil,   // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: false,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},

		CustomerAddressAttributeVatRequestSuccess: &customerAddressAttribute{
			Attribute: eav.NewAttribute("vat_request_success", // attribute_code
				40,     // attribute_id
				nil,    // backend_model
				"",     // backend_table
				"int",  // backend_type
				"",     // default_value
				2,      // entity_type_id
				"",     // frontend_class
				"text", // frontend_input
				"VAT number validation request success", // frontend_label
				nil,   // frontend_model
				false, // is_required
				false, // is_unique
				false, // is_user_defined
				"",    // note
				nil,   // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 true,
			isUsedForCustomerSegment: false,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},
	})
}
