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

// package custattr_test with customer attribute data contains test data which has
// been previously auto generated. So do not touch 8-)
package custattr_test

import (
	"github.com/corestoreio/csfw/customer/custattr"
	"github.com/corestoreio/csfw/eav"
)

type (
	// @todo website must be present in the slice
	// CustomerAttribute a data container for attributes. You can use this struct to
	// embed into your own struct for maybe overriding some method receivers.
	customerAttribute struct {
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

func (a *customerAttribute) IsVisible() bool {
	return a.isVisible
}
func (a *customerAttribute) InputFilter() string {
	return a.inputFilter
}
func (a *customerAttribute) MultilineCount() int64 {
	return a.multilineCount
}
func (a *customerAttribute) ValidateRules() string {
	return a.validateRules
}
func (a *customerAttribute) IsSystem() bool {
	return a.isSystem
}
func (a *customerAttribute) SortOrder() int64 {
	return a.sortOrder
}
func (a *customerAttribute) DataModel() eav.AttributeDataModeller {
	return a.dataModel
}
func (a *customerAttribute) IsUsedForCustomerSegment() bool {
	return a.isUsedForCustomerSegment
}
func (a *customerAttribute) ScopeIsVisible() bool {
	return a.scopeIsVisible
}
func (a *customerAttribute) ScopeIsRequired() bool {
	return a.scopeIsRequired
}
func (a *customerAttribute) ScopeDefaultValue() string {
	return a.scopeDefaultValue
}
func (a *customerAttribute) ScopeMultilineCount() int64 {
	return a.scopeMultilineCount
}

// Check if Attributer interface has been successfully implemented
var _ custattr.Attributer = (*customerAttribute)(nil)

const (
	CustomerAttributeConfirmation eav.AttributeIndex = iota
	CustomerAttributeCreatedAt
	CustomerAttributeCreatedIn
	CustomerAttributeDefaultBilling
	CustomerAttributeDefaultShipping
	CustomerAttributeDisableAutoGroupChange
	CustomerAttributeDob
	CustomerAttributeEmail
	CustomerAttributeFirstname
	CustomerAttributeGender
	CustomerAttributeGroupID
	CustomerAttributeLastname
	CustomerAttributeMiddlename
	CustomerAttributePasswordHash
	CustomerAttributePrefix
	CustomerAttributeRewardUpdateNotification
	CustomerAttributeRewardWarningNotification
	CustomerAttributeRpToken
	CustomerAttributeRpTokenCreatedAt
	CustomerAttributeStoreID
	CustomerAttributeSuffix
	CustomerAttributeTaxvat
	CustomerAttributeWebsiteID

	CustomerAttributeZZZ
)

type siCustomerAttribute struct{}

func (siCustomerAttribute) ByID(id int64) (eav.AttributeIndex, error) {
	switch id {
	case 16:
		return CustomerAttributeConfirmation, nil
	case 17:
		return CustomerAttributeCreatedAt, nil
	case 3:
		return CustomerAttributeCreatedIn, nil
	case 13:
		return CustomerAttributeDefaultBilling, nil
	case 14:
		return CustomerAttributeDefaultShipping, nil
	case 35:
		return CustomerAttributeDisableAutoGroupChange, nil
	case 11:
		return CustomerAttributeDob, nil
	case 9:
		return CustomerAttributeEmail, nil
	case 5:
		return CustomerAttributeFirstname, nil
	case 18:
		return CustomerAttributeGender, nil
	case 10:
		return CustomerAttributeGroupID, nil
	case 7:
		return CustomerAttributeLastname, nil
	case 6:
		return CustomerAttributeMiddlename, nil
	case 12:
		return CustomerAttributePasswordHash, nil
	case 4:
		return CustomerAttributePrefix, nil
	case 149:
		return CustomerAttributeRewardUpdateNotification, nil
	case 150:
		return CustomerAttributeRewardWarningNotification, nil
	case 33:
		return CustomerAttributeRpToken, nil
	case 34:
		return CustomerAttributeRpTokenCreatedAt, nil
	case 2:
		return CustomerAttributeStoreID, nil
	case 8:
		return CustomerAttributeSuffix, nil
	case 15:
		return CustomerAttributeTaxvat, nil
	case 1:
		return CustomerAttributeWebsiteID, nil

	default:
		return eav.AttributeIndex(0), eav.ErrAttributeNotFound
	}
}

func (siCustomerAttribute) ByCode(code string) (eav.AttributeIndex, error) {
	switch code {
	case "confirmation":
		return CustomerAttributeConfirmation, nil
	case "created_at":
		return CustomerAttributeCreatedAt, nil
	case "created_in":
		return CustomerAttributeCreatedIn, nil
	case "default_billing":
		return CustomerAttributeDefaultBilling, nil
	case "default_shipping":
		return CustomerAttributeDefaultShipping, nil
	case "disable_auto_group_change":
		return CustomerAttributeDisableAutoGroupChange, nil
	case "dob":
		return CustomerAttributeDob, nil
	case "email":
		return CustomerAttributeEmail, nil
	case "firstname":
		return CustomerAttributeFirstname, nil
	case "gender":
		return CustomerAttributeGender, nil
	case "group_id":
		return CustomerAttributeGroupID, nil
	case "lastname":
		return CustomerAttributeLastname, nil
	case "middlename":
		return CustomerAttributeMiddlename, nil
	case "password_hash":
		return CustomerAttributePasswordHash, nil
	case "prefix":
		return CustomerAttributePrefix, nil
	case "reward_update_notification":
		return CustomerAttributeRewardUpdateNotification, nil
	case "reward_warning_notification":
		return CustomerAttributeRewardWarningNotification, nil
	case "rp_token":
		return CustomerAttributeRpToken, nil
	case "rp_token_created_at":
		return CustomerAttributeRpTokenCreatedAt, nil
	case "store_id":
		return CustomerAttributeStoreID, nil
	case "suffix":
		return CustomerAttributeSuffix, nil
	case "taxvat":
		return CustomerAttributeTaxvat, nil
	case "website_id":
		return CustomerAttributeWebsiteID, nil

	default:
		return eav.AttributeIndex(0), eav.ErrAttributeNotFound
	}
}

var _ eav.AttributeGetter = (*siCustomerAttribute)(nil)

func init() {
	custattr.SetCustomerGetter(siCustomerAttribute{})
	custattr.SetCustomerCollection(custattr.AttributeSlice{

		CustomerAttributeConfirmation: &customerAttribute{
			Attribute: eav.NewAttribute("confirmation", // attribute_code
				16,             // attribute_id
				nil,            // backend_model
				"",             // backend_table
				"varchar",      // backend_type
				"",             // default_value
				1,              // entity_type_id
				"",             // frontend_class
				"text",         // frontend_input
				"Is Confirmed", // frontend_label
				nil,            // frontend_model
				false,          // is_required
				false,          // is_unique
				false,          // is_user_defined
				"",             // note
				nil,            // source_model

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

		CustomerAttributeCreatedAt: &customerAttribute{
			Attribute: eav.NewAttribute("created_at", // attribute_code
				17,           // attribute_id
				nil,          // backend_model
				"",           // backend_table
				"static",     // backend_type
				"",           // default_value
				1,            // entity_type_id
				"",           // frontend_class
				"datetime",   // frontend_input
				"Created At", // frontend_label
				nil,          // frontend_model
				false,        // is_required
				false,        // is_unique
				false,        // is_user_defined
				"",           // note
				nil,          // source_model

			),
			dataModel:                nil,
			inputFilter:              "datetime",
			isSystem:                 false,
			isUsedForCustomerSegment: true,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},

		CustomerAttributeCreatedIn: &customerAttribute{
			Attribute: eav.NewAttribute("created_in", // attribute_code
				3,              // attribute_id
				nil,            // backend_model
				"",             // backend_table
				"varchar",      // backend_type
				"",             // default_value
				1,              // entity_type_id
				"",             // frontend_class
				"text",         // frontend_input
				"Created From", // frontend_label
				nil,            // frontend_model
				false,          // is_required
				false,          // is_unique
				false,          // is_user_defined
				"",             // note
				nil,            // source_model

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
			sortOrder:                20,
			validateRules:            "",
		},

		CustomerAttributeDefaultBilling: &customerAttribute{
			Attribute: eav.NewAttribute("default_billing", // attribute_code
				13, // attribute_id
				custattr.CustomerBackendBilling().Config(eav.AttributeBackendIdx(CustomerAttributeDefaultBilling)), // backend_model
				"",     // backend_table
				"int",  // backend_type
				"",     // default_value
				1,      // entity_type_id
				"",     // frontend_class
				"text", // frontend_input
				"Default Billing Address", // frontend_label
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
			isUsedForCustomerSegment: true,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},

		CustomerAttributeDefaultShipping: &customerAttribute{
			Attribute: eav.NewAttribute("default_shipping", // attribute_code
				14, // attribute_id
				custattr.CustomerBackendShipping().Config(eav.AttributeBackendIdx(CustomerAttributeDefaultShipping)), // backend_model
				"",     // backend_table
				"int",  // backend_type
				"",     // default_value
				1,      // entity_type_id
				"",     // frontend_class
				"text", // frontend_input
				"Default Shipping Address", // frontend_label
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
			isUsedForCustomerSegment: true,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},

		CustomerAttributeDisableAutoGroupChange: &customerAttribute{
			Attribute: eav.NewAttribute("disable_auto_group_change", // attribute_code
				35, // attribute_id
				custattr.CustomerBackendDataBoolean().Config(eav.AttributeBackendIdx(CustomerAttributeDisableAutoGroupChange)), // backend_model
				"",        // backend_table
				"static",  // backend_type
				"",        // default_value
				1,         // entity_type_id
				"",        // frontend_class
				"boolean", // frontend_input
				"Disable Automatic Group Change Based on VAT ID", // frontend_label
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
			isVisible:                true,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                28,
			validateRules:            "",
		},

		CustomerAttributeDob: &customerAttribute{
			Attribute: eav.NewAttribute("dob", // attribute_code
				11, // attribute_id
				eav.AttributeBackendDatetime().Config(eav.AttributeBackendIdx(CustomerAttributeDob)), // backend_model
				"",              // backend_table
				"datetime",      // backend_type
				"",              // default_value
				1,               // entity_type_id
				"",              // frontend_class
				"date",          // frontend_input
				"Date Of Birth", // frontend_label
				eav.AttributeFrontendDatetime().Config(eav.AttributeFrontendIdx(CustomerAttributeDob)), // frontend_model
				false, // is_required
				false, // is_unique
				false, // is_user_defined
				"",    // note
				nil,   // source_model

			),
			dataModel:                nil,
			inputFilter:              "date",
			isSystem:                 false,
			isUsedForCustomerSegment: true,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                90,
			validateRules:            "a:1:{s:16:\"input_validation\";s:4:\"date\";}",
		},

		CustomerAttributeEmail: &customerAttribute{
			Attribute: eav.NewAttribute("email", // attribute_code
				9,        // attribute_id
				nil,      // backend_model
				"",       // backend_table
				"static", // backend_type
				"",       // default_value
				1,        // entity_type_id
				"",       // frontend_class
				"text",   // frontend_input
				"Email",  // frontend_label
				nil,      // frontend_model
				true,     // is_required
				false,    // is_unique
				false,    // is_user_defined
				"",       // note
				nil,      // source_model

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
			validateRules:            "a:1:{s:16:\"input_validation\";s:5:\"email\";}",
		},

		CustomerAttributeFirstname: &customerAttribute{
			Attribute: eav.NewAttribute("firstname", // attribute_code
				5,            // attribute_id
				nil,          // backend_model
				"",           // backend_table
				"varchar",    // backend_type
				"",           // default_value
				1,            // entity_type_id
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
			sortOrder:                40,
			validateRules:            "a:2:{s:15:\"max_text_length\";i:255;s:15:\"min_text_length\";i:1;}",
		},

		CustomerAttributeGender: &customerAttribute{
			Attribute: eav.NewAttribute("gender", // attribute_code
				18,       // attribute_id
				nil,      // backend_model
				"",       // backend_table
				"int",    // backend_type
				"",       // default_value
				1,        // entity_type_id
				"",       // frontend_class
				"select", // frontend_input
				"Gender", // frontend_label
				nil,      // frontend_model
				false,    // is_required
				false,    // is_unique
				false,    // is_user_defined
				"",       // note
				eav.AttributeSourceTable().Config(eav.AttributeSourceIdx(CustomerAttributeGender)), // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 false,
			isUsedForCustomerSegment: true,
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                110,
			validateRules:            "a:0:{}",
		},

		CustomerAttributeGroupID: &customerAttribute{
			Attribute: eav.NewAttribute("group_id", // attribute_code
				10,       // attribute_id
				nil,      // backend_model
				"",       // backend_table
				"static", // backend_type
				"",       // default_value
				1,        // entity_type_id
				"",       // frontend_class
				"select", // frontend_input
				"Group",  // frontend_label
				nil,      // frontend_model
				true,     // is_required
				false,    // is_unique
				false,    // is_user_defined
				"",       // note
				custattr.CustomerSourceGroup().Config(eav.AttributeSourceIdx(CustomerAttributeGroupID)), // source_model

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
			sortOrder:                25,
			validateRules:            "",
		},

		CustomerAttributeLastname: &customerAttribute{
			Attribute: eav.NewAttribute("lastname", // attribute_code
				7,           // attribute_id
				nil,         // backend_model
				"",          // backend_table
				"varchar",   // backend_type
				"",          // default_value
				1,           // entity_type_id
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
			sortOrder:                60,
			validateRules:            "a:2:{s:15:\"max_text_length\";i:255;s:15:\"min_text_length\";i:1;}",
		},

		CustomerAttributeMiddlename: &customerAttribute{
			Attribute: eav.NewAttribute("middlename", // attribute_code
				6,                     // attribute_id
				nil,                   // backend_model
				"",                    // backend_table
				"varchar",             // backend_type
				"",                    // default_value
				1,                     // entity_type_id
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
			sortOrder:                50,
			validateRules:            "",
		},

		CustomerAttributePasswordHash: &customerAttribute{
			Attribute: eav.NewAttribute("password_hash", // attribute_code
				12, // attribute_id
				custattr.CustomerBackendPassword().Config(eav.AttributeBackendIdx(CustomerAttributePasswordHash)), // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				1,         // entity_type_id
				"",        // frontend_class
				"hidden",  // frontend_input
				"",        // frontend_label
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
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},

		CustomerAttributePrefix: &customerAttribute{
			Attribute: eav.NewAttribute("prefix", // attribute_code
				4,         // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				1,         // entity_type_id
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
			sortOrder:                30,
			validateRules:            "",
		},

		CustomerAttributeRewardUpdateNotification: &customerAttribute{
			Attribute: eav.NewAttribute("reward_update_notification", // attribute_code
				149,    // attribute_id
				nil,    // backend_model
				"",     // backend_table
				"int",  // backend_type
				"",     // default_value
				1,      // entity_type_id
				"",     // frontend_class
				"text", // frontend_input
				"",     // frontend_label
				nil,    // frontend_model
				true,   // is_required
				false,  // is_unique
				false,  // is_user_defined
				"",     // note
				nil,    // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 false,
			isUsedForCustomerSegment: false,
			isVisible:                false,
			multilineCount:           1,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},

		CustomerAttributeRewardWarningNotification: &customerAttribute{
			Attribute: eav.NewAttribute("reward_warning_notification", // attribute_code
				150,    // attribute_id
				nil,    // backend_model
				"",     // backend_table
				"int",  // backend_type
				"",     // default_value
				1,      // entity_type_id
				"",     // frontend_class
				"text", // frontend_input
				"",     // frontend_label
				nil,    // frontend_model
				true,   // is_required
				false,  // is_unique
				false,  // is_user_defined
				"",     // note
				nil,    // source_model

			),
			dataModel:                nil,
			inputFilter:              "",
			isSystem:                 false,
			isUsedForCustomerSegment: false,
			isVisible:                false,
			multilineCount:           1,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},

		CustomerAttributeRpToken: &customerAttribute{
			Attribute: eav.NewAttribute("rp_token", // attribute_code
				33,        // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				1,         // entity_type_id
				"",        // frontend_class
				"hidden",  // frontend_input
				"",        // frontend_label
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
			isVisible:                false,
			multilineCount:           0,
			scopeDefaultValue:        "",
			scopeIsRequired:          false,
			scopeIsVisible:           false,
			scopeMultilineCount:      0,
			sortOrder:                0,
			validateRules:            "",
		},

		CustomerAttributeRpTokenCreatedAt: &customerAttribute{
			Attribute: eav.NewAttribute("rp_token_created_at", // attribute_code
				34,         // attribute_id
				nil,        // backend_model
				"",         // backend_table
				"datetime", // backend_type
				"",         // default_value
				1,          // entity_type_id
				"",         // frontend_class
				"date",     // frontend_input
				"",         // frontend_label
				nil,        // frontend_model
				false,      // is_required
				false,      // is_unique
				false,      // is_user_defined
				"",         // note
				nil,        // source_model

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
			validateRules:            "a:1:{s:16:\"input_validation\";s:4:\"date\";}",
		},

		CustomerAttributeStoreID: &customerAttribute{
			Attribute: eav.NewAttribute("store_id", // attribute_code
				2, // attribute_id
				custattr.CustomerBackendStore().Config(eav.AttributeBackendIdx(CustomerAttributeStoreID)), // backend_model
				"",          // backend_table
				"static",    // backend_type
				"",          // default_value
				1,           // entity_type_id
				"",          // frontend_class
				"select",    // frontend_input
				"Create In", // frontend_label
				nil,         // frontend_model
				true,        // is_required
				false,       // is_unique
				false,       // is_user_defined
				"",          // note
				custattr.CustomerSourceStore().Config(eav.AttributeSourceIdx(CustomerAttributeStoreID)), // source_model

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

		CustomerAttributeSuffix: &customerAttribute{
			Attribute: eav.NewAttribute("suffix", // attribute_code
				8,         // attribute_id
				nil,       // backend_model
				"",        // backend_table
				"varchar", // backend_type
				"",        // default_value
				1,         // entity_type_id
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
			sortOrder:                70,
			validateRules:            "",
		},

		CustomerAttributeTaxvat: &customerAttribute{
			Attribute: eav.NewAttribute("taxvat", // attribute_code
				15,               // attribute_id
				nil,              // backend_model
				"",               // backend_table
				"varchar",        // backend_type
				"",               // default_value
				1,                // entity_type_id
				"",               // frontend_class
				"text",           // frontend_input
				"Tax/VAT Number", // frontend_label
				nil,              // frontend_model
				false,            // is_required
				false,            // is_unique
				false,            // is_user_defined
				"",               // note
				nil,              // source_model

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
			sortOrder:                100,
			validateRules:            "a:1:{s:15:\"max_text_length\";i:255;}",
		},

		CustomerAttributeWebsiteID: &customerAttribute{
			Attribute: eav.NewAttribute("website_id", // attribute_code
				1, // attribute_id
				custattr.CustomerBackendWebsite().Config(eav.AttributeBackendIdx(CustomerAttributeWebsiteID)), // backend_model
				"",                     // backend_table
				"static",               // backend_type
				"",                     // default_value
				1,                      // entity_type_id
				"",                     // frontend_class
				"select",               // frontend_input
				"Associate to Website", // frontend_label
				nil,   // frontend_model
				true,  // is_required
				false, // is_unique
				false, // is_user_defined
				"",    // note
				custattr.CustomerSourceWebsite().Config(eav.AttributeSourceIdx(CustomerAttributeWebsiteID)), // source_model

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
			sortOrder:                10,
			validateRules:            "",
		},
	})
}
