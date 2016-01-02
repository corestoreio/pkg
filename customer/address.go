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

package customer

import (
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
)

type (
	AddressModel struct {
		// TBD
	}
)

var (
	_ eav.EntityTypeModeller                  = (*AddressModel)(nil)
	_ eav.EntityTypeTabler                    = (*AddressModel)(nil)
	_ eav.EntityTypeAdditionalAttributeTabler = (*AddressModel)(nil)
	_ eav.EntityTypeIncrementModeller         = (*AddressModel)(nil)
)

func (c *AddressModel) TBD() {

}

func (c *AddressModel) TableNameBase() string {
	return TableCollection.Name(TableIndexAddressEntity)
}

func (c *AddressModel) TableNameValue(i eav.ValueIndex) string {
	s, err := GetAddressValueStructure(i)
	if err != nil {
		return ""
	}
	return s.Name
}

// EntityTypeAdditionalAttributeTabler
func (c *AddressModel) TableAdditionalAttribute() (*csdb.Table, error) {
	return TableCollection.Structure(TableIndexEAVAttribute)
}

// EntityTypeAdditionalAttributeTabler
func (c *AddressModel) TableEavWebsite() (*csdb.Table, error) {
	return TableCollection.Structure(TableIndexEAVAttributeWebsite)
}

func Address() *AddressModel {
	return &AddressModel{}
}
