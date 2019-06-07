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

package catalog

import (
	"github.com/corestoreio/pkg/eav"
	"github.com/corestoreio/pkg/storage/csdb"
)

type (
	CategoryModel struct {
		// TBD
	}
)

var (
	_ eav.EntityTypeModeller                  = (*CategoryModel)(nil)
	_ eav.EntityTypeTabler                    = (*CategoryModel)(nil)
	_ eav.EntityTypeAdditionalAttributeTabler = (*CategoryModel)(nil)
)

func (c *CategoryModel) TBD() {

}

func (c *CategoryModel) TableNameBase() string {
	return TableCollection.Name(TableIndexCategoryEntity)
}

func (c *CategoryModel) TableNameValue(i eav.ValueIndex) string {
	s, err := GetCategoryValueStructure(i)
	if err != nil {
		return ""
	}
	return s.Name
}

// TableAdditionalAttribute needed for eav.EntityTypeAdditionalAttributeTabler
func (c *CategoryModel) TableAdditionalAttribute() (*csdb.Table, error) {
	return TableCollection.Structure(TableIndexEAVAttribute)
}

// TableEavWebsite needed for eav.EntityTypeAdditionalAttributeTabler
func (c *CategoryModel) TableEavWebsite() (*csdb.Table, error) {
	return nil, nil
}

func Category() *CategoryModel {
	return &CategoryModel{}
}
