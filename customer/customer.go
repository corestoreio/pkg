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
	// see table customer_eav_attribute.data_model
	DataModeller interface {
		TBD()
	}

	CustomerModel struct {
		// TBD
	}

	AttributeModel struct {
	}
)

var (
	_ eav.EntityTypeModeller                  = (*CustomerModel)(nil)
	_ eav.EntityTypeTabler                    = (*CustomerModel)(nil)
	_ eav.EntityTypeAttributeModeller         = (*CustomerModel)(nil)
	_ eav.EntityTypeAdditionalAttributeTabler = (*CustomerModel)(nil)
	_ eav.EntityAttributeCollectioner         = (*CustomerModel)(nil)
	_ eav.EntityTypeIncrementModeller         = (*CustomerModel)(nil)
)

func Attribute() *AttributeModel {
	return &AttributeModel{}
}

func (c AttributeModel) TBD() {

}

func (c *CustomerModel) TBD() {

}
func (c *CustomerModel) AttributeCollection() {

}

func (c *CustomerModel) TableNameBase() string {
	return GetTableName(TableEntity)
}

func (c *CustomerModel) TableNameValue(i eav.ValueIndex) string {
	switch i {
	case eav.EntityTypeDateTime:
		return GetTableName(TableEntityDatetime)
	case eav.EntityTypeDecimal:
		return GetTableName(TableEntityDecimal)
	case eav.EntityTypeInt:
		return GetTableName(TableEntityInt)
	case eav.EntityTypeText:
		return GetTableName(TableEntityText)
	case eav.EntityTypeVarchar:
		return GetTableName(TableEntityVarchar)
	}
	return ""
}

func (c *CustomerModel) TableNameAdditionalAttribute() string {
	return GetTableName(TableEavAttribute)
}

func Customer() *CustomerModel {
	return &CustomerModel{}
}
