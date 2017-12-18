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

package ddl_test

import (
	"testing"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableName(t *testing.T) {
	t.Parallel()

	t.Run("Short with prefix in table name", func(t *testing.T) {
		have := ddl.TableName("catalog_", "catalog_product_€ntity", "int")
		assert.Exactly(t, `catalog_product_ntity_int`, have)
	})
	t.Run("Short with prefix suffix", func(t *testing.T) {
		have := ddl.TableName("xyz_", "catalog_product_€ntity", "int")
		assert.Exactly(t, `xyz_catalog_product_ntity_int`, have)
	})
	t.Run("Short with prefix without suffix", func(t *testing.T) {
		have := ddl.TableName("xyz_", "catalog_product_€ntity")
		assert.Exactly(t, `xyz_catalog_product_ntity`, have)
	})
	t.Run("Short without prefix and suffix", func(t *testing.T) {
		have := ddl.TableName("", "catalog_product_€ntity")
		assert.Exactly(t, `catalog_product_ntity`, have)
	})

	t.Run("abbreviated", func(t *testing.T) {
		have := ddl.TableName("", "catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity")
		assert.Exactly(t, `cat_prd_ntity_cat_prd_ntity_cstr_cat_prd_ntity`, have)
	})
	t.Run("hashed", func(t *testing.T) {
		have := ddl.TableName("", "catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity_catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity")
		assert.Exactly(t, `t_bb0f749c31c69ed73ad028cb61f43745`, have)
	})
}

func TestIndexName(t *testing.T) {
	t.Parallel()

	t.Run("unique short", func(t *testing.T) {
		have := ddl.IndexName("unique", "sales_invoiced_aggregated_order", "period", "store_id", "order_status")
		assert.Exactly(t, `SALES_INVOICED_AGGREGATED_ORDER_PERIOD_STORE_ID_ORDER_STATUS`, have)
	})

	t.Run("unique abbreviated", func(t *testing.T) {
		have := ddl.IndexName("unique", "sales_invoiced_aggregated_order", "customer", "store_id", "order_status", "order_type")
		assert.Exactly(t, `SALES_INVOICED_AGGRED_ORDER_CSTR_STORE_ID_ORDER_STS_ORDER_TYPE`, have)
	})

	t.Run("unique hashed", func(t *testing.T) {
		have := ddl.IndexName("unique", "sales_invoiced_aggregated_order", "period", "store_id", "order_status", "order_type", "order_date")
		assert.Exactly(t, `UNQ_26EE326A968C157BC5004C8206E082E2`, have)
	})

	t.Run("fulltext short", func(t *testing.T) {
		have := ddl.IndexName("fulltext", "catalog_product_entity_int", "entity_id", "attribute_id")
		assert.Exactly(t, `CATALOG_PRODUCT_ENTITY_INT_ENTITY_ID_ATTRIBUTE_ID`, have)
	})
	t.Run("fulltext abbreviated", func(t *testing.T) {
		have := ddl.IndexName("fulltext", "catalog_product_entity_int", "entity_id", "attribute_id", "status_id", "value_id", "options_id")
		assert.Exactly(t, `CAT_PRD_ENTT_INT_ENTT_ID_ATTR_ID_STS_ID_VAL_ID_OPTS_ID`, have)
	})
	t.Run("fulltext hashed", func(t *testing.T) {
		have := ddl.IndexName("fulltext", "catalog_product_entity_int", "entity_id", "attribute_id", "status_id", "value_id", "options_id", "category_id", "customer_id", "group_id")
		assert.Exactly(t, `FTI_22F7610CC009B1EA375A74E4EE8BA2A1`, have)
	})

	t.Run("index short", func(t *testing.T) {
		have := ddl.IndexName("index", "catalog_product_entity_int", "entity_id", "attribute_id")
		assert.Exactly(t, `CATALOG_PRODUCT_ENTITY_INT_ENTITY_ID_ATTRIBUTE_ID`, have)
	})
	t.Run("index abbreviated", func(t *testing.T) {
		have := ddl.IndexName("index", "catalog_product_entity_int", "entity_id", "attribute_id", "entity_id", "attribute_id", "entity_id", "attribute_id")
		assert.Exactly(t, `CAT_PRD_ENTT_INT_ENTT_ID_ATTR_ID_ENTT_ID_ATTR_ID_ENTT_ID_ATTR_ID`, have)
	})
	t.Run("index hashed", func(t *testing.T) {
		have := ddl.IndexName("index", "catalog_product_entity_int", "entity_id", "attribute_id", "entity_id", "attribute_id", "entity_id", "attribute_id", "entity_id", "attribute_id", "entity_id", "attribute_id", "entity_id", "attribute_id")
		assert.Exactly(t, `IDX_B0CA6EB77B502EF1D87FDC50544DC34D`, have)
	})

	t.Run("unknown", func(t *testing.T) {
		have := ddl.IndexName("rablablablub", "catalog_product_entity_int", "entity_id", "attribute_id")
		assert.Exactly(t, `CATALOG_PRODUCT_ENTITY_INT_ENTITY_ID_ATTRIBUTE_ID`, have)
	})
	t.Run("empty", func(t *testing.T) {
		have := ddl.IndexName("rablablablub", "catalog_product_entity_int", "entity_id", "attribute_id")
		assert.Exactly(t, `CATALOG_PRODUCT_ENTITY_INT_ENTITY_ID_ATTRIBUTE_ID`, have)
	})
}

func TestTriggerName(t *testing.T) {
	t.Parallel()

	t.Run("unique short", func(t *testing.T) {
		have := ddl.TriggerName("sales_invoiced_aggregated_order", "before", "update")
		assert.Exactly(t, `sales_invoiced_aggregated_order_before_update`, have)
	})

	t.Run("unique abbreviated", func(t *testing.T) {
		have := ddl.TriggerName("sales_invoiced_aggregated_order_aggregated_order_customer", "before", "update")
		assert.Exactly(t, `sales_invoiced_aggred_order_aggred_order_cstr_before_update`, have)
	})

	t.Run("unique hashed", func(t *testing.T) {
		have := ddl.TriggerName("sales_invoiced_aggregated_ordersales_invoiced_aggregated_ordersales_invoiced_aggregated_ordersales_invoiced_aggregated_ordersales_invoiced_aggregated_ordersales_invoiced_aggregated_ordersales_invoiced_aggregated_order", "before", "update")
		assert.Exactly(t, `trg_fcd07948031b28888ab2b4959c097ef4`, have)
	})
}

func TestForeignKeyName(t *testing.T) {
	t.Parallel()

	t.Run("unique short", func(t *testing.T) {
		have := ddl.ForeignKeyName("catalog_product", "parent_id", "catalog_product_entity", "entity_id")
		assert.Exactly(t, `CATALOG_PRODUCT_PARENT_ID_CATALOG_PRODUCT_ENTITY_ENTITY_ID`, have)
	})

	t.Run("unique abbreviated", func(t *testing.T) {
		have := ddl.ForeignKeyName("catalog_product_bundle_option", "parent_id", "catalog_product_entity", "entity_id")
		assert.Exactly(t, `CAT_PRD_BNDL_OPT_PARENT_ID_CAT_PRD_ENTT_ENTT_ID`, have)
	})

	t.Run("unique hashed", func(t *testing.T) {
		have := ddl.ForeignKeyName("catalog_product_bundle_optioncatalog_product_bundle_optioncatalog_product_bundle_option", "parent_id", "catalog_product_entitycatalog_product_entitycatalog_product_entity", "entity_id")
		assert.Exactly(t, `FK_8FC0390CFB720B81470F95BE9E5A8584`, have)
	})
}

func TestAlias(t *testing.T) {
	t.Parallel()
	var testData = map[string]string{
		"passwords":                                            "pa",
		"p":                                                    "p",
		"admin_passwords":                                      "adpa",
		"admin_system_messages":                                "adsyme",
		"admin_user":                                           "adus",
		"admin_user_session":                                   "adusse",
		"adminnotification_inbox":                              "adin",
		"authorization_role":                                   "auro",
		"authorization_rule":                                   "auru",
		"cache":                                                "ca",
		"cache_tag":                                            "cata",
		"captcha_log":                                          "calo",
		"catalog_category_entity":                              "cacaen",
		"catalog_category_entity_datetime":                     "cacaenda",
		"catalog_category_entity_decimal":                      "cacaende",
		"catalog_category_entity_int":                          "cacaenin",
		"catalog_category_entity_text":                         "cacaente", // ;-) a poo duck
		"catalog_category_entity_varchar":                      "cacaenva",
		"catalog_category_product":                             "cacapr",
		"catalog_category_product_index":                       "cacaprin",
		"catalog_category_product_index_replica":               "cacaprinre",
		"catalog_category_product_index_tmp":                   "cacaprintm",
		"catalog_compare_item":                                 "cacoit",
		"catalog_eav_attribute":                                "caeaat",
		"catalog_product_entity_media_gallery_value_to_entity": "caprenmegavatoen",
		"cms_block":               "cmbl",
		"importexport_importdata": "imim",
		"quote":                   "qu",
	}

	for have, want := range testData {
		require.Exactly(t, want, ddl.Shorten(have), "Input: %q", have)
	}
}
