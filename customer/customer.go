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

import (
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
)

type (
	// see table customer_eav_attribute.data_model
	DataModeller interface {
		TBD()
	}

	CustomerModel struct {
		// TBD
	}
)

var (
	_ eav.EntityTypeModeller                  = (*CustomerModel)(nil)
	_ eav.EntityTypeTabler                    = (*CustomerModel)(nil)
	_ eav.EntityTypeAdditionalAttributeTabler = (*CustomerModel)(nil)
	_ eav.EntityTypeIncrementModeller         = (*CustomerModel)(nil)
)

func (c *CustomerModel) TBD() {

}

func (c *CustomerModel) TableNameBase() string {
	return GetTableName(TableIndexEntity)
}

func (c *CustomerModel) TableNameValue(i eav.ValueIndex) string {
	s, err := GetCustomerValueStructure(i)
	if err != nil {
		return ""
	}
	return s.Name
}

// EntityTypeAdditionalAttributeTabler
func (c *CustomerModel) TableAdditionalAttribute() (*csdb.TableStructure, error) {
	return GetTableStructure(TableIndexEAVAttribute)
}

// EntityTypeAdditionalAttributeTabler
func (c *CustomerModel) TableEavWebsite() (*csdb.TableStructure, error) {
	return GetTableStructure(TableIndexEAVAttributeWebsite)
}

func Customer() *CustomerModel {
	return &CustomerModel{}
}
