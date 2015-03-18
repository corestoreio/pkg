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

package tools

import (
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
)

func TestQueryToStruct(t *testing.T) {
	db := csdb.MustConnectTest()
	defer db.Close()

	structCode, err := QueryToStruct(db, "CatalogProductEavAttributeJoin", testQryEavAttributeJoin)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, testQryEavAttributeJoinExpected, string(structCode))
}

const testQryEavAttributeJoin = `SELECT main_table.attribute_id, main_table.entity_type_id,
    main_table.attribute_code, main_table.attribute_model, main_table.backend_model,
    main_table.backend_type, main_table.backend_table, main_table.frontend_model,
    main_table.frontend_input, main_table.frontend_label, main_table.frontend_class,
    main_table.source_model, main_table.is_required, main_table.is_user_defined,
    main_table.default_value, main_table.is_unique, main_table.note,
    additional_table.frontend_input_renderer, additional_table.is_global,
    additional_table.is_visible, additional_table.is_searchable,
    additional_table.is_filterable, additional_table.is_comparable,
    additional_table.is_visible_on_front, additional_table.is_html_allowed_on_front,
    additional_table.is_used_for_price_rules, additional_table.is_filterable_in_search,
    additional_table.used_in_product_listing, additional_table.used_for_sort_by,
    additional_table.is_configurable, additional_table.apply_to,
    additional_table.is_visible_in_advanced_search, additional_table.position,
    additional_table.is_wysiwyg_enabled, additional_table.is_used_for_promo_rules,
    additional_table.search_weight
    FROM eav_attribute AS main_table
        INNER JOIN catalog_eav_attribute AS additional_table
            ON (additional_table.attribute_id = main_table.attribute_id) AND (main_table.entity_type_id = 4)`

const testQryEavAttributeJoinExpected = `
type (
	// CatalogProductEavAttributeJoinSlice contains pointers to CatalogProductEavAttributeJoin types
	CatalogProductEavAttributeJoinSlice []*CatalogProductEavAttributeJoin
	// CatalogProductEavAttributeJoin a type for a MySQL Query
	CatalogProductEavAttributeJoin struct {
		AttributeID               int64          ` + "`" + `db:"attribute_id"` + "`" + `                  // attribute_id smallint(5) unsigned NOT NULL  DEFAULT '0'
		EntityTypeID              int64          ` + "`" + `db:"entity_type_id"` + "`" + `                // entity_type_id smallint(5) unsigned NOT NULL  DEFAULT '0'
		AttributeCode             dbr.NullString ` + "`" + `db:"attribute_code"` + "`" + `                // attribute_code varchar(255) NULL
		AttributeModel            dbr.NullString ` + "`" + `db:"attribute_model"` + "`" + `               // attribute_model varchar(255) NULL
		BackendModel              dbr.NullString ` + "`" + `db:"backend_model"` + "`" + `                 // backend_model varchar(255) NULL
		BackendType               string         ` + "`" + `db:"backend_type"` + "`" + `                  // backend_type varchar(8) NOT NULL  DEFAULT 'static'
		BackendTable              dbr.NullString ` + "`" + `db:"backend_table"` + "`" + `                 // backend_table varchar(255) NULL
		FrontendModel             dbr.NullString ` + "`" + `db:"frontend_model"` + "`" + `                // frontend_model varchar(255) NULL
		FrontendInput             dbr.NullString ` + "`" + `db:"frontend_input"` + "`" + `                // frontend_input varchar(50) NULL
		FrontendLabel             dbr.NullString ` + "`" + `db:"frontend_label"` + "`" + `                // frontend_label varchar(255) NULL
		FrontendClass             dbr.NullString ` + "`" + `db:"frontend_class"` + "`" + `                // frontend_class varchar(255) NULL
		SourceModel               dbr.NullString ` + "`" + `db:"source_model"` + "`" + `                  // source_model varchar(255) NULL
		IsRequired                bool           ` + "`" + `db:"is_required"` + "`" + `                   // is_required smallint(5) unsigned NOT NULL  DEFAULT '0'
		IsUserDefined             bool           ` + "`" + `db:"is_user_defined"` + "`" + `               // is_user_defined smallint(5) unsigned NOT NULL  DEFAULT '0'
		DefaultValue              dbr.NullString ` + "`" + `db:"default_value"` + "`" + `                 // default_value text NULL
		IsUnique                  bool           ` + "`" + `db:"is_unique"` + "`" + `                     // is_unique smallint(5) unsigned NOT NULL  DEFAULT '0'
		Note                      dbr.NullString ` + "`" + `db:"note"` + "`" + `                          // note varchar(255) NULL
		FrontendInputRenderer     dbr.NullString ` + "`" + `db:"frontend_input_renderer"` + "`" + `       // frontend_input_renderer varchar(255) NULL
		IsGlobal                  bool           ` + "`" + `db:"is_global"` + "`" + `                     // is_global smallint(5) unsigned NOT NULL  DEFAULT '1'
		IsVisible                 bool           ` + "`" + `db:"is_visible"` + "`" + `                    // is_visible smallint(5) unsigned NOT NULL  DEFAULT '1'
		IsSearchable              bool           ` + "`" + `db:"is_searchable"` + "`" + `                 // is_searchable smallint(5) unsigned NOT NULL  DEFAULT '0'
		IsFilterable              bool           ` + "`" + `db:"is_filterable"` + "`" + `                 // is_filterable smallint(5) unsigned NOT NULL  DEFAULT '0'
		IsComparable              bool           ` + "`" + `db:"is_comparable"` + "`" + `                 // is_comparable smallint(5) unsigned NOT NULL  DEFAULT '0'
		IsVisibleOnFront          bool           ` + "`" + `db:"is_visible_on_front"` + "`" + `           // is_visible_on_front smallint(5) unsigned NOT NULL  DEFAULT '0'
		IsHtmlAllowedOnFront      bool           ` + "`" + `db:"is_html_allowed_on_front"` + "`" + `      // is_html_allowed_on_front smallint(5) unsigned NOT NULL  DEFAULT '0'
		IsUsedForPriceRules       bool           ` + "`" + `db:"is_used_for_price_rules"` + "`" + `       // is_used_for_price_rules smallint(5) unsigned NOT NULL  DEFAULT '0'
		IsFilterableInSearch      bool           ` + "`" + `db:"is_filterable_in_search"` + "`" + `       // is_filterable_in_search smallint(5) unsigned NOT NULL  DEFAULT '0'
		UsedInProductListing      int64          ` + "`" + `db:"used_in_product_listing"` + "`" + `       // used_in_product_listing smallint(5) unsigned NOT NULL  DEFAULT '0'
		UsedForSortBy             int64          ` + "`" + `db:"used_for_sort_by"` + "`" + `              // used_for_sort_by smallint(5) unsigned NOT NULL  DEFAULT '0'
		IsConfigurable            bool           ` + "`" + `db:"is_configurable"` + "`" + `               // is_configurable smallint(5) unsigned NOT NULL  DEFAULT '1'
		ApplyTo                   dbr.NullString ` + "`" + `db:"apply_to"` + "`" + `                      // apply_to varchar(255) NULL
		IsVisibleInAdvancedSearch bool           ` + "`" + `db:"is_visible_in_advanced_search"` + "`" + ` // is_visible_in_advanced_search smallint(5) unsigned NOT NULL  DEFAULT '0'
		Position                  int64          ` + "`" + `db:"position"` + "`" + `                      // position int(11) NOT NULL  DEFAULT '0'
		IsWysiwygEnabled          bool           ` + "`" + `db:"is_wysiwyg_enabled"` + "`" + `            // is_wysiwyg_enabled smallint(5) unsigned NOT NULL  DEFAULT '0'
		IsUsedForPromoRules       bool           ` + "`" + `db:"is_used_for_promo_rules"` + "`" + `       // is_used_for_promo_rules smallint(5) unsigned NOT NULL  DEFAULT '0'
		SearchWeight              int64          ` + "`" + `db:"search_weight"` + "`" + `                 // search_weight smallint(5) unsigned NOT NULL  DEFAULT '1'
	}
)
`
