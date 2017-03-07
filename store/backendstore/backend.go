// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package backendstore

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgsource"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/errors"
)

// Configuration just exported for the sake of documentation. See fields for
// more information. Please call the New() function for creating a new
// Configuration object. Only the New() function will set the paths to the
// fields.
type Configuration struct {

	// GeneralStoreInformationName => Store Name.
	// Path: general/store_information/name
	GeneralStoreInformationName cfgmodel.Str

	// GeneralStoreInformationPhone => Store Phone Number.
	// Path: general/store_information/phone
	GeneralStoreInformationPhone cfgmodel.Str

	// GeneralStoreInformationHours => Store Hours of Operation.
	// Path: general/store_information/hours
	GeneralStoreInformationHours cfgmodel.Str

	// GeneralStoreInformationCountryID => Country. You must set the
	// cfgmodel.MapIntResolver after calling New() of this package.
	// Path: general/store_information/country_id
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	GeneralStoreInformationCountryID cfgmodel.MapIntStr

	// GeneralStoreInformationRegionID => Region/State. You must set the
	// cfgmodel.MapIntResolver after calling New() of this package.
	// Path: general/store_information/region_id
	GeneralStoreInformationRegionID cfgmodel.MapIntStr

	// GeneralStoreInformationPostcode => ZIP/Postal Code.
	// Path: general/store_information/postcode
	GeneralStoreInformationPostcode cfgmodel.Str

	// GeneralStoreInformationCity => City.
	// Path: general/store_information/city
	GeneralStoreInformationCity cfgmodel.Str

	// GeneralStoreInformationStreetLine1 => Street Address.
	// Path: general/store_information/street_line1
	GeneralStoreInformationStreetLine1 cfgmodel.Str

	// GeneralStoreInformationStreetLine2 => Street Address Line 2.
	// Path: general/store_information/street_line2
	GeneralStoreInformationStreetLine2 cfgmodel.Str

	// GeneralStoreInformationMerchantVatNumber => VAT Number.
	// Path: general/store_information/merchant_vat_number
	GeneralStoreInformationMerchantVatNumber cfgmodel.Str

	// GeneralSingleStoreModeEnabled => Enable Single-Store Mode. This setting
	// will not be taken into account if system has more than one store view.
	// Path: general/single_store_mode/enabled
	GeneralSingleStoreModeEnabled cfgmodel.Bool
}

// New initializes the backend configuration models containing the cfgpath.Route
// variable to the appropriate entries in the storage. The argument SectionSlice
// and opts will be applied to all models.
func New(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Configuration {
	be := &Configuration{}

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	be.GeneralSingleStoreModeEnabled = cfgmodel.NewBool(`general/single_store_mode/enabled`, append(opts, cfgmodel.WithSource(cfgsource.EnableDisable))...)

	be.GeneralStoreInformationName = cfgmodel.NewStr(`general/store_information/name`, opts...)
	be.GeneralStoreInformationPhone = cfgmodel.NewStr(`general/store_information/phone`, opts...)
	be.GeneralStoreInformationHours = cfgmodel.NewStr(`general/store_information/hours`, opts...)
	be.GeneralStoreInformationCountryID = cfgmodel.NewMapIntStr(`general/store_information/country_id`, opts...)
	be.GeneralStoreInformationRegionID = cfgmodel.NewMapIntStr(`general/store_information/region_id`, opts...)
	be.GeneralStoreInformationPostcode = cfgmodel.NewStr(`general/store_information/postcode`, opts...)
	be.GeneralStoreInformationCity = cfgmodel.NewStr(`general/store_information/city`, opts...)
	be.GeneralStoreInformationStreetLine1 = cfgmodel.NewStr(`general/store_information/street_line1`, opts...)
	be.GeneralStoreInformationStreetLine2 = cfgmodel.NewStr(`general/store_information/street_line2`, opts...)
	be.GeneralStoreInformationMerchantVatNumber = cfgmodel.NewStr(`general/store_information/merchant_vat_number`, opts...)
	return be
}

// StoreInformation defines the address data for a merchant. Might be usable in
// e.g. text/template or html/template.
type StoreInformation struct {
	ScopeID     scope.TypeID
	Name        string
	Phone       string
	Hours       string
	Country     string
	Region      string
	PostCode    string
	City        string
	StreetLine1 string
	StreetLine2 string
	Vat         string
}

// StoreInformation reads the store information from the configuration depending
// on the scope. Might be usable in e.g. text/template or html/template. Does
// not yet cache internally per scope the data.
func (c *Configuration) StoreInformation(sg config.Scoped) (*StoreInformation, error) {
	name, err := c.GeneralStoreInformationName.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationName")
	}
	phone, err := c.GeneralStoreInformationPhone.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationPhone")
	}
	hours, err := c.GeneralStoreInformationHours.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationHours")
	}
	country, err := c.GeneralStoreInformationCountryID.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationCountryID")
	}
	region, err := c.GeneralStoreInformationRegionID.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationRegionID")
	}
	postCode, err := c.GeneralStoreInformationPostcode.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationPostcode")
	}
	city, err := c.GeneralStoreInformationCity.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationCity")
	}
	sl1, err := c.GeneralStoreInformationStreetLine1.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationStreetLine1")
	}
	sl2, err := c.GeneralStoreInformationStreetLine2.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationStreetLine2")
	}
	vat, err := c.GeneralStoreInformationMerchantVatNumber.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationMerchantVatNumber")
	}

	return &StoreInformation{
		ScopeID:     sg.ScopeID(),
		Name:        name,
		Phone:       phone,
		Hours:       hours,
		Country:     country,
		Region:      region,
		PostCode:    postCode,
		City:        city,
		StreetLine1: sl1,
		StreetLine2: sl2,
		Vat:         vat,
	}, nil
}
