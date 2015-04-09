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

package testgen

import (
	"testing"

	"github.com/corestoreio/csfw/catalog/catattr"
	"github.com/corestoreio/csfw/customer/custattr"
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/stretchr/testify/assert"
)

func TestGetAttributeSelect(t *testing.T) {
	db := csdb.MustConnectTest()
	defer db.Close()
	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)
	et, err := eav.GetEntityTypeCollection().GetByCode("catalog_product")
	if err != nil {
		t.Error(err)
	}
	websiteID := int64(1)
	dbrSelect, err := eav.GetAttributeSelectSql(dbrSess, et.AdditionalAttributeTable, et.EntityTypeID, websiteID)
	if err != nil {
		t.Error(err)
	}
	sql, _ := dbrSelect.ToSql()

	assert.Equal(
		t,
		"SELECT `main_table`.`attribute_id`, `main_table`.`entity_type_id`, `main_table`.`attribute_code`, `main_table`.`backend_model`, `main_table`.`backend_type`, `main_table`.`backend_table`, `main_table`.`frontend_model`, `main_table`.`frontend_input`, `main_table`.`frontend_label`, `main_table`.`frontend_class`, `main_table`.`source_model`, `main_table`.`is_required`, `main_table`.`is_user_defined`, `main_table`.`default_value`, `main_table`.`is_unique`, `main_table`.`note`, `additional_table`.`frontend_input_renderer`, `additional_table`.`is_global`, `additional_table`.`is_visible`, `additional_table`.`is_searchable`, `additional_table`.`is_filterable`, `additional_table`.`is_comparable`, `additional_table`.`is_visible_on_front`, `additional_table`.`is_html_allowed_on_front`, `additional_table`.`is_used_for_price_rules`, `additional_table`.`is_filterable_in_search`, `additional_table`.`used_in_product_listing`, `additional_table`.`used_for_sort_by`, `additional_table`.`is_configurable`, `additional_table`.`apply_to`, `additional_table`.`is_visible_in_advanced_search`, `additional_table`.`position`, `additional_table`.`is_wysiwyg_enabled`, `additional_table`.`is_used_for_promo_rules`, `additional_table`.`search_weight` FROM `eav_attribute` AS `main_table` INNER JOIN `catalog_eav_attribute` AS `additional_table` ON (`additional_table`.`attribute_id` = `main_table`.`attribute_id`) AND (`main_table`.`entity_type_id` = ?)",
		sql,
	)

	et, err = eav.GetEntityTypeCollection().GetByCode("customer")
	if err != nil {
		t.Error(err)
	}
	dbrSelect, err = eav.GetAttributeSelectSql(dbrSess, et.AdditionalAttributeTable, et.EntityTypeID, websiteID)
	if err != nil {
		t.Error(err)
	}
	sql, _ = dbrSelect.ToSql()
	assert.Equal(
		t,
		"SELECT `main_table`.`attribute_id`, `main_table`.`entity_type_id`, `main_table`.`attribute_code`, `main_table`.`backend_model`, `main_table`.`backend_type`, `main_table`.`backend_table`, `main_table`.`frontend_model`, `main_table`.`frontend_input`, `main_table`.`frontend_label`, `main_table`.`frontend_class`, `main_table`.`source_model`, `main_table`.`is_required`, `main_table`.`is_user_defined`, `main_table`.`default_value`, `main_table`.`is_unique`, `main_table`.`note`, `additional_table`.`is_visible`, `additional_table`.`input_filter`, `additional_table`.`multiline_count`, `additional_table`.`validate_rules`, `additional_table`.`is_system`, `additional_table`.`sort_order`, `additional_table`.`data_model`, `additional_table`.`is_used_for_customer_segment`, `scope_table`.`is_visible` AS `scope_is_visible`, `scope_table`.`is_required` AS `scope_is_required`, `scope_table`.`default_value` AS `scope_default_value`, `scope_table`.`multiline_count` AS `scope_multiline_count` FROM `eav_attribute` AS `main_table` INNER JOIN `customer_eav_attribute` AS `additional_table` ON (`additional_table`.`attribute_id` = `main_table`.`attribute_id`) AND (`main_table`.`entity_type_id` = ?) LEFT JOIN `customer_eav_attribute_website` AS `scope_table` ON (`scope_table`.`attribute_id` = `main_table`.`attribute_id`) AND (`scope_table`.`website_id` = ?)",
		sql,
	)
}

// @todo implement this test also for EntityAttributeCollection
func TestNewAttributeCustomer(t *testing.T) {
	//t.Skip("@todo implement New() attribute functionality")
	for _, et := range []string{"customer", "customer_address", "catalog_product", "catalog_category", "@todo demo"} {
		cu, err := eav.GetEntityTypeByCode(et)
		assert.NoError(t, err)
		ca := cu.AttributeModel.New()

		switch ca.(type) {
		case custattr.Attributer:
			// do here customer attribute stuff
			break
		case catattr.Attributer:
			// do here catalog stuff
			break
		case eav.Attributer:
			// do here default attribute stuff @todo add a sample EAV model
			break
		default:
			t.Errorf("%s Attribute model does not implement interface eav.Attributer and its child interfaces", et)
		}

	}
}
