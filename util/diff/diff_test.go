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

package diff_test

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/diff"
)

func TestUnified(t *testing.T) {
	tests := []struct {
		haveA   string
		haveB   string
		want    string
		wantErr error
	}{
		{
			"SELECT\n`main_table`.`attribute_id`,\n`main_table`.`entity_type_id`,\n`main_table`.`attribute_code`,\n`main_table`.`attribute_model`,\n`main_table`.`backend_model`,\n`main_table`.`backend_type`,\n`main_table`.`backend_table`,\n`main_table`.`frontend_model`,\n`main_table`.`frontend_input`,\n`main_table`.`frontend_label`,\n`main_table`.`frontend_class`,\n`main_table`.`source_model`,\n`main_table`.`is_user_defined`,\n`main_table`.`is_unique`,\n`main_table`.`note`,\n`additional_table`.`input_filter`,\n`additional_table`.`validate_rules`,\n`additional_table`.`is_system`,\n`additional_table`.`sort_order`,\n`additional_table`.`data_model`,\nIFNULL(`scope_table`.`is_visible`,\n`additional_table`.`is_visible`)\nAS\n`is_visible`,\nIFNULL(`scope_table`.`is_required`,\n`main_table`.`is_required`)\nAS\n`is_required`,\nIFNULL(`scope_table`.`default_value`,\n`main_table`.`default_value`)\nAS\n`default_value`,\nIFNULL(`scope_table`.`multiline_count`,\n`additional_table`.`multiline_count`)\nAS\n`multiline_count`\nFROM\n`eav_attribute`\nAS\n`main_table`\nINNER\nJOIN\n`customer_eav_attribute`\nAS\n`additional_table`\nON\n(additional_table.attribute_id\n=\nmain_table.attribute_id)\nAND\n(main_table.entity_type_id\n=\n?)\nLEFT\nJOIN\n`customer_eav_attribute_website`\nAS\n`scope_table`\nON\n(scope_table.attribute_id\n=\nmain_table.attribute_id)\nAND\n(scope_table.website_id\n=\n?)\nWHERE\n`multiline_count`\n>\n0\nORDER\nBY\n`main_table`.`attribute_id`",
			"SELECT\n`main_table`.`attribute_id`,\n`main_table`.`entity_type_id`,\n`main_table`.`attribute_model`,\n`main_table`.`backend_model`,\n`main_table`.`backend_type`,\n`main_table`.`backend_table`,\n`main_table`.`frontend_model`,\n`main_table`.`frontend_input`,\n`main_table`.`frontend_label`,\n`main_table`.`frontend_class`,\n`main_table`.`source_model`,\n`main_table`.`is_user_defined`,\n`main_table`.`is_unique`,\n`main_table`.`note`,\n`additional_table`.`input_filter`,\n`additional_table`.`validate_rules`,\n`additional_table`.`is_system`,\n`additional_table`.`sort_order`,\n`additional_table`.`data_model`,\nIFNULL(`scope_table`.`is_visible`,\n`additional_table`.`is_visible`)\nAS\n`is_visible`,\nIFNULL(`scope_table`.`is_required`,\n`main_table`.`is_required`)\nAS\n`is_required`,\nIFNULL(`scope_table`.`default_value`,\n`main_table`.`default_value`)\nAS\n`default_value`,\nIFNULL(`scope_table`.`multiline_count`,\n`additional_table`.`multiline_count`)\nAS\n`multiline_count`\n\n\nFROM\n`eav_attribute`\nAS\n`main_table`\nINNER\nJOIN\n`customer_eav_attribute`\nAS\n`additional_table`\nON\n(additional_table.attribute_id\n=\nmain_table.attribute_id)\nAND\n(main_table.entity_type_id\n=\n?)\nLEFT\nJOIN\n`customer_eav_attribute_website`\nAS\n`scope_table`\nON\n(scope_table.attribute_id\n=\nmain_table.attribute_id)\nWHERE\n`multiline_count`\n>\n0\nORDER\nBY\n`main_table`.`attribute_id`",
			"--- Original\n+++ Current\n@@ -3,3 +3,2 @@\n `main_table`.`entity_type_id`,\n-`main_table`.`attribute_code`,\n `main_table`.`attribute_model`,\n@@ -37,2 +36,4 @@\n `multiline_count`\n+\n+\n FROM\n@@ -63,6 +64,2 @@\n main_table.attribute_id)\n-AND\n-(scope_table.website_id\n-=\n-?)\n WHERE\n",
			nil,
		},
	}
	for i, test := range tests {
		haveDiff, haveErr := diff.Unified(test.haveA, test.haveB)
		if test.wantErr != nil {
			assert.Error(t, haveErr, "Index %d", i)
			continue
		}
		assert.Exactly(t, test.want, haveDiff)
	}
}
