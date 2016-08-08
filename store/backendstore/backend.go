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
	"io"
	"text/template"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/util/errors"
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

	TemplateAddress *template.Template
}

const TemplateAddressText = `"{{name}}\n{{street_line1}}\n{{with street_line2}}{{.}}\n{{end}}
            {{city}}, {{region}} {{postcode}},\n{{country}}`

// New initializes the backend configuration models containing the cfgpath.Route
// variable to the appropriate entries in the storage. The argument SectionSlice
// and opts will be applied to all models.
func New(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Configuration {
	be := &Configuration{}

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	be.GeneralSingleStoreModeEnabled = cfgmodel.NewBool(`general/single_store_mode/enabled`, append(opts, cfgmodel.WithSource(source.EnableDisable))...)

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

type addressData struct {
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

func (c *Configuration) scopedAddressData(sg config.Scoped) (*addressData, error) {
	name, _, err := c.GeneralStoreInformationName.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationName")
	}
	phone, _, err := c.GeneralStoreInformationPhone.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationPhone")
	}
	hours, _, err := c.GeneralStoreInformationHours.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationHours")
	}
	country, _, err := c.GeneralStoreInformationCountryID.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationCountryID")
	}
	region, _, err := c.GeneralStoreInformationRegionID.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationRegionID")
	}
	postCode, _, err := c.GeneralStoreInformationPostcode.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationPostcode")
	}
	city, _, err := c.GeneralStoreInformationCity.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationCity")
	}
	sl1, _, err := c.GeneralStoreInformationStreetLine1.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationStreetLine1")
	}
	sl2, _, err := c.GeneralStoreInformationStreetLine2.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationStreetLine2")
	}
	vat, _, err := c.GeneralStoreInformationMerchantVatNumber.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendstore] GeneralStoreInformationMerchantVatNumber")
	}

	return &addressData{
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

func (c *Configuration) FormatAddressText(name string, w io.Writer, sg config.Scoped) error {
	if name == "" {
		name = "address"
	}
	if c.TemplateAddress == nil {
		// lazy init or user already created a custom template object
		var err error
		c.TemplateAddress, err = template.New(name).Parse(TemplateAddressText)
		if err != nil {
			return errors.Wrap(err, "[backendstore] Template.Parse")
		}
	}
	data, err := c.scopedAddressData(sg)
	if err != nil {
		return errors.Wrap(err, "[backendstore] Template.scopedAddressData")
	}
	return errors.Wrap(c.TemplateAddress.ExecuteTemplate(w, name, data), "[backendstore] Template.Execute")
}
