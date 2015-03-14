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

package catalog

import "github.com/corestoreio/csfw/eav"

type (
	CategoryModel struct {
		// TBD
	}
)

var (
	_ eav.EntityTypeModeller                  = (*CategoryModel)(nil)
	_ eav.EntityTypeTabler                    = (*CategoryModel)(nil)
	_ eav.EntityTypeAttributeModeller         = (*CategoryModel)(nil)
	_ eav.EntityTypeAdditionalAttributeTabler = (*CategoryModel)(nil)
	_ eav.EntityAttributeCollectioner         = (*CategoryModel)(nil)
)

func (c *CategoryModel) TBD() {

}
func (c *CategoryModel) AttributeCollection() {

}

func (c *CategoryModel) TableNameBase() string {
	return GetTableName(TableCategoryEntity)
}

func (c *CategoryModel) TableNameValue(i eav.ValueIndex) string {
	switch i {
	case eav.EntityTypeDateTime:
		return GetTableName(TableCategoryEntityDatetime)
	case eav.EntityTypeDecimal:
		return GetTableName(TableCategoryEntityDecimal)
	case eav.EntityTypeInt:
		return GetTableName(TableCategoryEntityInt)
	case eav.EntityTypeText:
		return GetTableName(TableCategoryEntityText)
	case eav.EntityTypeVarchar:
		return GetTableName(TableCategoryEntityVarchar)
	}
	return ""
}

func (c *CategoryModel) TableNameAdditionalAttribute() string {
	return GetTableName(TableEavAttribute)
}

func Category() *CategoryModel {
	return &CategoryModel{}
}
