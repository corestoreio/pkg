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

import "github.com/corestoreio/csfw/eav"

type (
	AddressModel struct {
		// TBD
	}

	AddressAttributeModel struct {
	}
)

var (
	_ eav.EntityTypeModeller                  = (*AddressModel)(nil)
	_ eav.EntityTypeTabler                    = (*AddressModel)(nil)
	_ eav.EntityTypeAttributeModeller         = (*AddressModel)(nil)
	_ eav.EntityTypeAdditionalAttributeTabler = (*AddressModel)(nil)
	_ eav.EntityAttributeCollectioner         = (*AddressModel)(nil)
	_ eav.EntityTypeIncrementModeller         = (*AddressModel)(nil)
)

func AddressAttribute() *AddressAttributeModel {
	return &AddressAttributeModel{}
}

func (c AddressAttributeModel) TBD() {

}

func (c *AddressModel) TBD() {

}

func (c *AddressModel) AttributeCollection() {

}

func (c *AddressModel) TableNameBase() string {
	return GetTableName(TableCustomerEntity)
}

func (c *AddressModel) TableNameValue(i eav.ValueIndex) string {
	s, err := GetCustomerAddressValueStructure(i)
	if err != nil {
		return ""
	}
	return s.Name
}

func (c *AddressModel) TableNameAdditionalAttribute() string {
	return GetTableName(TableCustomerEavAttribute)
}

func Address() *AddressModel {
	return &AddressModel{}
}
