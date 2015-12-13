// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package diff_test

import (
	"testing"

	"github.com/corestoreio/csfw/util/diff"
	"github.com/corestoreio/csfw/util/sqlbeautifier"
	"github.com/stretchr/testify/assert"
)

func TestUnified(t *testing.T) {
	tests := []struct {
		haveA   string
		haveB   string
		want    string
		wantErr error
	}{
		{
			sqlbeautifier.MustFromString("SELECT `main_table`.`attribute_id`, `main_table`.`entity_type_id`, `main_table`.`attribute_code`, `main_table`.`attribute_model`, `main_table`.`backend_model`, `main_table`.`backend_type`, `main_table`.`backend_table`, `main_table`.`frontend_model`, `main_table`.`frontend_input`, `main_table`.`frontend_label`, `main_table`.`frontend_class`, `main_table`.`source_model`, `main_table`.`is_user_defined`, `main_table`.`is_unique`, `main_table`.`note`, `additional_table`.`input_filter`, `additional_table`.`validate_rules`, `additional_table`.`is_system`, `additional_table`.`sort_order`, `additional_table`.`data_model`, IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`) AS `is_visible`, IFNULL(`scope_table`.`is_required`, `main_table`.`is_required`) AS `is_required`, IFNULL(`scope_table`.`default_value`, `main_table`.`default_value`) AS `default_value`, IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count` FROM `eav_attribute` AS `main_table` INNER JOIN `customer_eav_attribute` AS `additional_table` ON (additional_table.attribute_id = main_table.attribute_id) AND (main_table.entity_type_id = ?) LEFT JOIN `customer_eav_attribute_website` AS `scope_table` ON (scope_table.attribute_id = main_table.attribute_id) AND (scope_table.website_id = ?) WHERE `multiline_count` > 0 ORDER BY `main_table`.`attribute_id`"),
			sqlbeautifier.MustFromString("SELECT `main_table`.`attribute_id`, `main_table`.`entity_type_id`, `main_table`.`attribute_model`, `main_table`.`backend_model`, `main_table`.`backend_type`, `main_table`.`backend_table`, `main_table`.`frontend_model`, `main_table`.`frontend_input`, `main_table`.`frontend_label`, `main_table`.`frontend_class`, `main_table`.`source_model`, `main_table`.`is_user_defined`, `main_table`.`is_unique`, `main_table`.`note`, `additional_table`.`input_filter`, `additional_table`.`validate_rules`, `additional_table`.`is_system`, `additional_table`.`sort_order`, `additional_table`.`data_model`, IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`) AS `is_visible`, IFNULL(`scope_table`.`is_required`, `main_table`.`is_required`) AS `is_required`, IFNULL(`scope_table`.`default_value`, `main_table`.`default_value`) AS `default_value`, IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count`   FROM\n`eav_attribute` AS `main_table` INNER JOIN `customer_eav_attribute` AS `additional_table` ON (additional_table.attribute_id = main_table.attribute_id) AND (main_table.entity_type_id = ?) LEFT JOIN `customer_eav_attribute_website` AS `scope_table` ON (scope_table.attribute_id = main_table.attribute_id) WHERE `multiline_count` > 0 ORDER BY `main_table`.`attribute_id`"),
			"--- Original\n+++ Current\n@@ -3,3 +3,2 @@\n \tmain_table.entity_type_id, \n-\tmain_table.attribute_code, \n \tmain_table.attribute_model, \n@@ -44,4 +43,3 @@\n \tcustomer_eav_attribute_website as scope_table on \n-\t\t(scope_table.attribute_id = main_table.attribute_id) and \n-\t\t(scope_table.website_id = :v2)\n+\t\t(scope_table.attribute_id = main_table.attribute_id)\n \n",
			nil,
		},
	}
	for i, test := range tests {
		haveDiff, haveErr := diff.Unified(test.haveA, test.haveB)
		if test.wantErr != nil {
			assert.Error(t, haveErr, "Index %d", i)
			continue
		}
		if false == assert.Exactly(t, test.want, haveDiff) {
			t.Errorf("Error @ Index %d\n%s\n", i, haveDiff)
		}
	}
}
