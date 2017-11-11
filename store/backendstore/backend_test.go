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

package backendstore_test

import (
	"testing"

	"github.com/corestoreio/cspkg/config"
	"github.com/corestoreio/cspkg/config/cfgmock"
	"github.com/corestoreio/cspkg/store/backendstore"
	"github.com/corestoreio/cspkg/store/scope"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

// backend overall backend models for all tests
var backend *backendstore.Configuration

// this would belong into the test suit setup
func init() {
	cfgStruct, err := backendstore.NewConfigStructure()
	if err != nil {
		panic(err)
	}
	backend = backendstore.New(cfgStruct)
}

type mockIntToStr struct {
	error
	string
}

func (mis mockIntToStr) IntToStr(_ config.Scoped, id int) (string, error) {
	switch id {
	case 144:
		return "Germany", nil
	case 5:
		return "Berlin", nil
	}
	return mis.string, mis.error
}

func TestConfiguration_AddressData(t *testing.T) {
	sg := cfgmock.NewService(cfgmock.PathValue{
		backend.GeneralStoreInformationName.MustFQWebsite(3):              `CoreStore SA`,
		backend.GeneralStoreInformationPhone.MustFQWebsite(3):             `32608`,
		backend.GeneralStoreInformationHours.MustFQWebsite(3):             `11am-7pm`,
		backend.GeneralStoreInformationCountryID.MustFQWebsite(3):         144,
		backend.GeneralStoreInformationRegionID.MustFQWebsite(3):          5,
		backend.GeneralStoreInformationPostcode.MustFQWebsite(3):          `10100`,
		backend.GeneralStoreInformationCity.MustFQWebsite(3):              `Shopville`,
		backend.GeneralStoreInformationStreetLine1.MustFQWebsite(3):       `Market Str 134`,
		backend.GeneralStoreInformationStreetLine2.MustFQWebsite(3):       `Booth 987`,
		backend.GeneralStoreInformationMerchantVatNumber.MustFQWebsite(3): `DE12345678`,
	}).NewScoped(3, 4)

	backend.GeneralStoreInformationCountryID.MapIntResolver = mockIntToStr{}
	backend.GeneralStoreInformationRegionID.MapIntResolver = mockIntToStr{}

	ad, err := backend.StoreInformation(sg)
	assert.NoError(t, err)
	want := &backendstore.StoreInformation{ScopeID: scope.Store.Pack(4), Name: "CoreStore SA", Phone: "32608", Hours: "11am-7pm", Country: "Germany", Region: "Berlin", PostCode: "10100", City: "Shopville", StreetLine1: "Market Str 134", StreetLine2: "Booth 987", Vat: "DE12345678"}
	assert.Exactly(t, want, ad)
}

func TestConfiguration_AddressData_Country_Error(t *testing.T) {
	sg := cfgmock.NewService(cfgmock.PathValue{}).NewScoped(3, 4)

	backend.GeneralStoreInformationCountryID.MapIntResolver = mockIntToStr{
		error: errors.NewNotSupportedf("Some countries are not supported"),
	}
	backend.GeneralStoreInformationRegionID.MapIntResolver = mockIntToStr{}

	ad, err := backend.StoreInformation(sg)
	assert.True(t, errors.IsNotSupported(err), "%+v", err)
	assert.Nil(t, ad)
}

func TestConfiguration_AddressData_Region_Error(t *testing.T) {
	sg := cfgmock.NewService(cfgmock.PathValue{}).NewScoped(3, 4)

	backend.GeneralStoreInformationRegionID.MapIntResolver = mockIntToStr{
		error: errors.NewUnauthorizedf("Some countries are not supported"),
	}
	backend.GeneralStoreInformationCountryID.MapIntResolver = mockIntToStr{}

	ad, err := backend.StoreInformation(sg)
	assert.True(t, errors.IsUnauthorized(err), "%+v", err)
	assert.Nil(t, ad)
}
