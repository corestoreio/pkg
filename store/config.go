// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package store

// import (
// 	"github.com/corestoreio/pkg/config"
// 	"github.com/corestoreio/pkg/store/scope"
// )
//
// const (
// 	// GeneralStoreInformationName => Store Name.
// 	ConfigPathGeneralStoreInformationName = `general/store_information/name`
// 	// GeneralStoreInformationPhone => Store Phone Number.
// 	ConfigPathGeneralStoreInformationPhone = `general/store_information/phone`
// 	// GeneralStoreInformationHours => Store Hours of Operation.
// 	ConfigPathGeneralStoreInformationHours = `general/store_information/hours`
// 	// GeneralStoreInformationCountryID => Country
// 	ConfigPathGeneralStoreInformationCountryID = `general/store_information/country_id`
// 	// GeneralStoreInformationRegionID => Region/State
// 	ConfigPathGeneralStoreInformationRegionID = `general/store_information/region_id`
// 	// GeneralStoreInformationPostcode => ZIP/Postal Code.
// 	ConfigPathGeneralStoreInformationPostcode = `general/store_information/postcode`
// 	// GeneralStoreInformationCity => City.
// 	ConfigPathGeneralStoreInformationCity = `general/store_information/city`
// 	// GeneralStoreInformationStreetLine1 => Street Address.
// 	ConfigPathGeneralStoreInformationStreetLine1 = `general/store_information/street_line1`
// 	// GeneralStoreInformationStreetLine2 => Street Address Line 2.
// 	ConfigPathGeneralStoreInformationStreetLine2 = `general/store_information/street_line2`
// 	// GeneralStoreInformationMerchantVatNumber => VAT Number.
// 	ConfigPathGeneralStoreInformationMerchantVatNumber = `general/store_information/merchant_vat_number`
// )
//
// // NewConfigStructure global configuration structure for this package. Used in
// // frontend (to display the user all the settings) and in backend (scope checks
// // and default values). See the source code of this function for the overall
// // available sections, groups and fields.
// func NewConfigStructure() (config.Sections, error) {
// 	return config.MakeSectionsValidated(
// 		&config.Section{
// 			ID:        "general",
// 			Label:     `General`,
// 			SortOrder: 10,
// 			Scopes:    scope.PermStore,
// 			Groups: config.MakeGroups(
// 				&config.Group{
// 					ID:        "store_information",
// 					Label:     `Store Information`,
// 					SortOrder: 100,
// 					Scopes:    scope.PermStore,
// 					Fields: config.MakeFields(
// 						&config.Field{
// 							// Path: general/store_information/name
// 							ID:        "name",
// 							Label:     `Store Name`,
// 							Type:      config.TypeText,
// 							SortOrder: 10,
// 							Visible:   true,
// 							Scopes:    scope.PermStore,
// 						},
//
// 						&config.Field{
// 							// Path: general/store_information/phone
// 							ID:        "phone",
// 							Label:     `Store Phone Number`,
// 							Type:      config.TypeText,
// 							SortOrder: 20,
// 							Visible:   true,
// 							Scopes:    scope.PermStore,
// 						},
//
// 						&config.Field{
// 							// Path: general/store_information/hours
// 							ID:        "hours",
// 							Label:     `Store Hours of Operation`,
// 							Type:      config.TypeText,
// 							SortOrder: 22,
// 							Visible:   true,
// 							Scopes:    scope.PermStore,
// 						},
//
// 						&config.Field{
// 							// Path: general/store_information/country_id
// 							ID:         "country_id",
// 							Label:      `Country`,
// 							Type:       config.TypeSelect,
// 							SortOrder:  25,
// 							Visible:    true,
// 							Scopes:     scope.PermWebsite,
// 							CanBeEmpty: true,
// 						},
//
// 						&config.Field{
// 							// Path: general/store_information/region_id
// 							ID:        "region_id",
// 							Label:     `Region/State`,
// 							Type:      config.TypeText,
// 							SortOrder: 27,
// 							Visible:   true,
// 							Scopes:    scope.PermWebsite,
// 						},
//
// 						&config.Field{
// 							// Path: general/store_information/postcode
// 							ID:        "postcode",
// 							Label:     `ZIP/Postal Code`,
// 							Type:      config.TypeText,
// 							SortOrder: 30,
// 							Visible:   true,
// 							Scopes:    scope.PermWebsite,
// 						},
//
// 						&config.Field{
// 							// Path: general/store_information/city
// 							ID:        "city",
// 							Label:     `City`,
// 							Type:      config.TypeText,
// 							SortOrder: 45,
// 							Visible:   true,
// 							Scopes:    scope.PermWebsite,
// 						},
//
// 						&config.Field{
// 							// Path: general/store_information/street_line1
// 							ID:        "street_line1",
// 							Label:     `Street Address`,
// 							Type:      config.TypeText,
// 							SortOrder: 55,
// 							Visible:   true,
// 							Scopes:    scope.PermWebsite,
// 						},
//
// 						&config.Field{
// 							// Path: general/store_information/street_line2
// 							ID:        "street_line2",
// 							Label:     `Street Address Line 2`,
// 							Type:      config.TypeText,
// 							SortOrder: 60,
// 							Visible:   true,
// 							Scopes:    scope.PermWebsite,
// 						},
//
// 						&config.Field{
// 							// Path: general/store_information/merchant_vat_number
// 							ID:         "merchant_vat_number",
// 							Label:      `VAT Number`,
// 							Type:       config.TypeText,
// 							SortOrder:  61,
// 							Visible:    true,
// 							Scopes:     scope.PermWebsite,
// 							CanBeEmpty: true,
// 						},
// 					),
// 				},
// 			),
// 		},
// 	)
// }
