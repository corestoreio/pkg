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

package catalog

import (
	"github.com/corestoreio/cspkg/eav"
	"github.com/corestoreio/cspkg/storage/csdb"
)

type (
	ProductModel struct {
		// TBD
	}
)

var (
	_ eav.EntityTypeModeller                  = (*ProductModel)(nil)
	_ eav.EntityTypeTabler                    = (*ProductModel)(nil)
	_ eav.EntityTypeAdditionalAttributeTabler = (*ProductModel)(nil)
)

func (c *ProductModel) TBD() {

}

func (c *ProductModel) TableNameBase() string {
	return TableCollection.Name(TableIndexProductEntity)
}

func (c *ProductModel) TableNameValue(i eav.ValueIndex) string {
	s, err := GetProductValueStructure(i)
	if err != nil {
		return ""
	}
	return s.Name
}

// TableAdditionalAttribute needed for EntityTypeAdditionalAttributeTabler
func (c *ProductModel) TableAdditionalAttribute() (*csdb.Table, error) {
	return TableCollection.Structure(TableIndexEAVAttribute)
}

// TableEavWebsite needed for EntityTypeAdditionalAttributeTabler
func (c *ProductModel) TableEavWebsite() (*csdb.Table, error) {
	return nil, nil
}

func Product() *ProductModel {
	return &ProductModel{}
}
