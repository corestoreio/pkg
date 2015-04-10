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

package eav

import (
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

// GetAttributeSelectSql generates the select query to retrieve full attribute configuration
// Implements the scope on a SQL query basis so that attribute functions does not need to deal with it.
// Tests see the tools package
// @see magento2/app/code/Magento/Eav/Model/Resource/Attribute/Collection.php::_initSelect()
func GetAttributeSelectSql(dbrSess dbr.SessionRunner, aat EntityTypeAdditionalAttributeTabler, entityTypeID, websiteId int64) (*dbr.SelectBuilder, error) {

	/*
			@todo
		   SELECT
		     `main_table`.`attribute_id`,
		     `main_table`.`entity_type_id`,
		     `main_table`.`attribute_code`,
		     `main_table`.`backend_model`,
		     `main_table`.`backend_type`,
		     `main_table`.`backend_table`,
		     `main_table`.`frontend_model`,
		     `main_table`.`frontend_input`,
		     `main_table`.`frontend_label`,
		     `main_table`.`frontend_class`,
		     `main_table`.`source_model`,
		     `main_table`.`is_user_defined`,
		     `main_table`.`is_unique`,
		     `main_table`.`note`,
		     `additional_table`.`input_filter`,
		     `additional_table`.`validate_rules`,
		     `additional_table`.`is_system`,
		     `additional_table`.`sort_order`,
		     `additional_table`.`data_model`,
		     `additional_table`.`is_used_for_customer_segment`,
		     IFNULL(`scope_table`.`is_required`, `main_table`.`is_required`)               AS `is_required`,
		     IFNULL(`scope_table`.`default_value`, `main_table`.`default_value`)           AS `default_value`,
		     IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`)           AS `is_visible`,
		     IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count`
		   FROM `eav_attribute` AS `main_table` INNER JOIN `customer_eav_attribute` AS `additional_table`
		       ON (`additional_table`.`attribute_id` = `main_table`.`attribute_id`) AND (`main_table`.`entity_type_id` = 1)
		     LEFT JOIN `customer_eav_attribute_website` AS `scope_table`
		       ON (`scope_table`.`attribute_id` = `main_table`.`attribute_id`) AND (`scope_table`.`website_id` = 4)
	*/

	ta, err := GetTableStructure(TableIndexAttribute)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	taa, err := aat.TableAdditionalAttribute()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	selectSql := dbrSess.
		Select(ta.AllColumnAliasQuote(csdb.MainTable)...).
		From(ta.Name, csdb.MainTable).
		Join(
		dbr.JoinTable(taa.Name, "additional_table"),
		taa.ColumnAliasQuote("additional_table"),
		dbr.JoinOn("`additional_table`.`attribute_id` = `main_table`.`attribute_id`"),
		dbr.JoinOn("`main_table`.`entity_type_id` = ?", entityTypeID),
	)

	tew, err := aat.TableEavWebsite()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	if tew != nil {
		const scopeTable = "scope_table"
		l := len(tew.Columns) * 2
		cols := make([]string, l)
		j := 0
		for i := 0; i < l; i = i + 2 {
			cols[i] = scopeTable + "." + tew.Columns[j] // real column name
			cols[i+1] = "scope_" + tew.Columns[j]       // alias column name
			j++
		}

		selectSql.
			LeftJoin(
			dbr.JoinTable(tew.Name, "scope_table"),
			dbr.ColumnAlias(cols...),
			dbr.JoinOn("`scope_table`.`attribute_id` = `main_table`.`attribute_id`"),
			dbr.JoinOn("`scope_table`.`website_id` = ?", websiteId),
		)
	}
	return selectSql, nil
}
