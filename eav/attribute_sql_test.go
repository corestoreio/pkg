package eav_test

import (
	"testing"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
)

func TestGetAttributeSelectSql(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()

	dbrSess := dbc.NewSession()
	dbrSelect, err := eav.GetAttributeSelectSql(dbrSess, codegen.NewAddAttrTables(dbc.DB, "customer"), 1, 4)
	if err != nil {
		t.Error(err)
	} else {
		sql, _ := dbrSelect.ToSql()
		assert.Equal(t, "SELECT `main_table`.`attribute_id`, `main_table`.`entity_type_id`, `main_table`.`attribute_code`, `main_table`.`backend_model`, `main_table`.`backend_type`, `main_table`.`backend_table`, `main_table`.`frontend_model`, `main_table`.`frontend_input`, `main_table`.`frontend_label`, `main_table`.`frontend_class`, `main_table`.`source_model`, `main_table`.`is_user_defined`, `main_table`.`is_unique`, `main_table`.`note`, `additional_table`.`input_filter`, `additional_table`.`validate_rules`, `additional_table`.`is_system`, `additional_table`.`sort_order`, `additional_table`.`data_model`, IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`) AS `is_visible`, IFNULL(`scope_table`.`is_required`, `main_table`.`is_required`) AS `is_required`, IFNULL(`scope_table`.`default_value`, `main_table`.`default_value`) AS `default_value`, IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count` FROM `eav_attribute` AS `main_table` INNER JOIN `customer_eav_attribute` AS `additional_table` ON (`additional_table`.`attribute_id` = `main_table`.`attribute_id`) AND (`main_table`.`entity_type_id` = ?) LEFT JOIN `customer_eav_attribute_website` AS `scope_table` ON (`scope_table`.`attribute_id` = `main_table`.`attribute_id`) AND (`scope_table`.`website_id` = ?)", sql)
	}

	// @todo error is that we have column attribute_model in the select list but it should not occur
	// because in codegen it is defined that this column has no usage so we can skip it.
}
