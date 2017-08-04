// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package eav_test

import (
	"testing"

	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/stretchr/testify/assert"
)

func TestIfNull(t *testing.T) {
	tests := []struct {
		alias      string
		columnName string
		defaultVal string
		want       string
	}{
		{
			"manufacturer", "value", "",
			"IFNULL(`manufacturerStore`.`value`,IFNULL(`manufacturerGroup`.`value`,IFNULL(`manufacturerWebsite`.`value`,IFNULL(`manufacturerDefault`.`value`,'')))) AS `manufacturer`",
		},
		{
			"manufacturer", "value", "0",
			"IFNULL(`manufacturerStore`.`value`,IFNULL(`manufacturerGroup`.`value`,IFNULL(`manufacturerWebsite`.`value`,IFNULL(`manufacturerDefault`.`value`,0)))) AS `manufacturer`",
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, eav.IfNull(test.alias, test.columnName, test.defaultVal), "Index %d", i)
	}
}

func TestSelect_Join_EAVIfNull(t *testing.T) {
	t.Parallel()
	const want = "SELECT IFNULL(`manufacturerStore`.`value`,IFNULL(`manufacturerGroup`.`value`,IFNULL(`manufacturerWebsite`.`value`,IFNULL(`manufacturerDefault`.`value`,'')))) AS `manufacturer`, cpe.* FROM `catalog_product_entity` AS `cpe` LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerDefault` ON (manufacturerDefault.scope = 0) AND (manufacturerDefault.scope_id = 0) AND (manufacturerDefault.attribute_id = 83) AND (manufacturerDefault.value IS NOT NULL) LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerWebsite` ON (manufacturerWebsite.scope = 1) AND (manufacturerWebsite.scope_id = 10) AND (manufacturerWebsite.attribute_id = 83) AND (manufacturerWebsite.value IS NOT NULL) LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerGroup` ON (manufacturerGroup.scope = 2) AND (manufacturerGroup.scope_id = 20) AND (manufacturerGroup.attribute_id = 83) AND (manufacturerGroup.value IS NOT NULL) LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerStore` ON (manufacturerStore.scope = 2) AND (manufacturerStore.scope_id = 20) AND (manufacturerStore.attribute_id = 83) AND (manufacturerStore.value IS NOT NULL)"

	s := dbr.NewSelect(eav.IfNull("manufacturer", "value", "''"), "cpe.*").
		From("catalog_product_entity", "cpe").
		LeftJoin(
			dbr.MakeIdentifier("catalog_product_entity_varchar", "manufacturerDefault"),
			dbr.Column("manufacturerDefault.scope = 0"),
			dbr.Column("manufacturerDefault.scope_id = 0"),
			dbr.Column("manufacturerDefault.attribute_id = 83"),
			dbr.Column("manufacturerDefault.value IS NOT NULL"),
		).
		LeftJoin(
			dbr.MakeIdentifier("catalog_product_entity_varchar", "manufacturerWebsite"),
			dbr.Column("manufacturerWebsite.scope = 1"),
			dbr.Column("manufacturerWebsite.scope_id = 10"),
			dbr.Column("manufacturerWebsite.attribute_id = 83"),
			dbr.Column("manufacturerWebsite.value IS NOT NULL"),
		).
		LeftJoin(
			dbr.MakeIdentifier("catalog_product_entity_varchar", "manufacturerGroup"),
			dbr.Column("manufacturerGroup.scope = 2"),
			dbr.Column("manufacturerGroup.scope_id = 20"),
			dbr.Column("manufacturerGroup.attribute_id = 83"),
			dbr.Column("manufacturerGroup.value IS NOT NULL"),
		).
		LeftJoin(
			dbr.MakeIdentifier("catalog_product_entity_varchar", "manufacturerStore"),
			dbr.Column("manufacturerStore.scope = 2"),
			dbr.Column("manufacturerStore.scope_id = 20"),
			dbr.Column("manufacturerStore.attribute_id = 83"),
			dbr.Column("manufacturerStore.value IS NOT NULL"),
		)

	sql, _, err := s.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t,
		want,
		sql,
	)
}
