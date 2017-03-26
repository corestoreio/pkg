// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package csdb_test

import (
	"strings"
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestIsValidIdentifier(t *testing.T) {
	t.Parallel()

	t.Run("Names", func(t *testing.T) {
		const errDummy = errors.Error("Dummy")
		tests := []struct {
			have string
			want error
		}{
			{"$catalog_product_3ntity", nil},
			{"`catalog_product_3ntity", errDummy},
			{"", errDummy},
			{strings.Repeat("a", 65), errDummy},
			{strings.Repeat("a", 64), nil},
		}
		for i, test := range tests {
			haveErr := csdb.IsValidIdentifier(test.have)
			if test.want != nil {
				assert.True(t, errors.IsNotValid(haveErr), "Index %d", i)
			} else {
				assert.NoError(t, haveErr, "Index %d", i)
			}
		}
	})
	t.Run("No args", func(t *testing.T) {
		haveErr := csdb.IsValidIdentifier()
		assert.True(t, errors.IsNotValid(haveErr), "%+v", haveErr)
	})
	t.Run("Multiple args but last with error", func(t *testing.T) {
		haveErr := csdb.IsValidIdentifier("customer", "product", "namecatalog_category_anc_categs_index_tmpcatalog_category_anc_categs")
		assert.True(t, errors.IsNotValid(haveErr), "%+v", haveErr)
	})
}

var benchmarkIsValidIdentifier error

func BenchmarkIsValidIdentifier(b *testing.B) {
	const id = `$catalog_product_3ntity`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkIsValidIdentifier = csdb.IsValidIdentifier(id)
	}
	if benchmarkIsValidIdentifier != nil {
		b.Fatalf("%+v", benchmarkIsValidIdentifier)
	}
}

func TestTableName(t *testing.T) {
	t.Parallel()

	t.Run("Short with prefix in table name", func(t *testing.T) {
		have := csdb.TableName("catalog_", "catalog_product_€ntity", "int")
		assert.Exactly(t, `catalog_product_ntity_int`, have)
	})
	t.Run("Short with prefix suffix", func(t *testing.T) {
		have := csdb.TableName("xyz_", "catalog_product_€ntity", "int")
		assert.Exactly(t, `xyz_catalog_product_ntity_int`, have)
	})
	t.Run("Short with prefix without suffix", func(t *testing.T) {
		have := csdb.TableName("xyz_", "catalog_product_€ntity")
		assert.Exactly(t, `xyz_catalog_product_ntity`, have)
	})
	t.Run("Short without prefix and suffix", func(t *testing.T) {
		have := csdb.TableName("", "catalog_product_€ntity")
		assert.Exactly(t, `catalog_product_ntity`, have)
	})

	t.Run("abbreviated", func(t *testing.T) {
		have := csdb.TableName("", "catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity")
		assert.Exactly(t, `cat_prd_ntity_cat_prd_ntity_cstr_cat_prd_ntity`, have)
	})
	t.Run("hashed", func(t *testing.T) {
		have := csdb.TableName("", "catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity_catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity")
		assert.Exactly(t, `t_bb0f749c31c69ed73ad028cb61f43745`, have)
	})
}

func TestIndexName(t *testing.T) {
	t.Parallel()

	t.Run("unique short", func(t *testing.T) {
		have := csdb.IndexName("unique", "sales_invoiced_aggregated_order", "period", "store_id", "order_status")
		assert.Exactly(t, `SALES_INVOICED_AGGREGATED_ORDER_PERIOD_STORE_ID_ORDER_STATUS`, have)
	})

	t.Run("unique abbreviated", func(t *testing.T) {
		have := csdb.IndexName("unique", "sales_invoiced_aggregated_order", "customer", "store_id", "order_status", "order_type")
		assert.Exactly(t, `SALES_INVOICED_AGGRED_ORDER_CSTR_STORE_ID_ORDER_STS_ORDER_TYPE`, have)
	})

	t.Run("unique hashed", func(t *testing.T) {
		have := csdb.IndexName("unique", "sales_invoiced_aggregated_order", "period", "store_id", "order_status", "order_type", "order_date")
		assert.Exactly(t, `UNQ_26EE326A968C157BC5004C8206E082E2`, have)
	})

	t.Run("fulltext short", func(t *testing.T) {
		have := csdb.IndexName("fulltext", "catalog_product_entity_int", "entity_id", "attribute_id")
		assert.Exactly(t, `CATALOG_PRODUCT_ENTITY_INT_ENTITY_ID_ATTRIBUTE_ID`, have)
	})
	t.Run("fulltext abbreviated", func(t *testing.T) {
		have := csdb.IndexName("fulltext", "catalog_product_entity_int", "entity_id", "attribute_id", "status_id", "value_id", "options_id")
		assert.Exactly(t, `CAT_PRD_ENTT_INT_ENTT_ID_ATTR_ID_STS_ID_VAL_ID_OPTS_ID`, have)
	})
	t.Run("fulltext hashed", func(t *testing.T) {
		have := csdb.IndexName("fulltext", "catalog_product_entity_int", "entity_id", "attribute_id", "status_id", "value_id", "options_id", "category_id", "customer_id", "group_id")
		assert.Exactly(t, `FTI_22F7610CC009B1EA375A74E4EE8BA2A1`, have)
	})

	t.Run("index short", func(t *testing.T) {
		have := csdb.IndexName("index", "catalog_product_entity_int", "entity_id", "attribute_id")
		assert.Exactly(t, `CATALOG_PRODUCT_ENTITY_INT_ENTITY_ID_ATTRIBUTE_ID`, have)
	})
	t.Run("index abbreviated", func(t *testing.T) {
		have := csdb.IndexName("index", "catalog_product_entity_int", "entity_id", "attribute_id", "entity_id", "attribute_id", "entity_id", "attribute_id")
		assert.Exactly(t, `CAT_PRD_ENTT_INT_ENTT_ID_ATTR_ID_ENTT_ID_ATTR_ID_ENTT_ID_ATTR_ID`, have)
	})
	t.Run("index hashed", func(t *testing.T) {
		have := csdb.IndexName("index", "catalog_product_entity_int", "entity_id", "attribute_id", "entity_id", "attribute_id", "entity_id", "attribute_id", "entity_id", "attribute_id", "entity_id", "attribute_id", "entity_id", "attribute_id")
		assert.Exactly(t, `IDX_B0CA6EB77B502EF1D87FDC50544DC34D`, have)
	})

	t.Run("unknown", func(t *testing.T) {
		have := csdb.IndexName("rablablablub", "catalog_product_entity_int", "entity_id", "attribute_id")
		assert.Exactly(t, `CATALOG_PRODUCT_ENTITY_INT_ENTITY_ID_ATTRIBUTE_ID`, have)
	})
	t.Run("empty", func(t *testing.T) {
		have := csdb.IndexName("rablablablub", "catalog_product_entity_int", "entity_id", "attribute_id")
		assert.Exactly(t, `CATALOG_PRODUCT_ENTITY_INT_ENTITY_ID_ATTRIBUTE_ID`, have)
	})
}

func TestTriggerName(t *testing.T) {
	t.Parallel()

	t.Run("unique short", func(t *testing.T) {
		have := csdb.TriggerName("sales_invoiced_aggregated_order", "before", "update")
		assert.Exactly(t, `sales_invoiced_aggregated_order_before_update`, have)
	})

	t.Run("unique abbreviated", func(t *testing.T) {
		have := csdb.TriggerName("sales_invoiced_aggregated_order_aggregated_order_customer", "before", "update")
		assert.Exactly(t, `sales_invoiced_aggred_order_aggred_order_cstr_before_update`, have)
	})

	t.Run("unique hashed", func(t *testing.T) {
		have := csdb.TriggerName("sales_invoiced_aggregated_ordersales_invoiced_aggregated_ordersales_invoiced_aggregated_ordersales_invoiced_aggregated_ordersales_invoiced_aggregated_ordersales_invoiced_aggregated_ordersales_invoiced_aggregated_order", "before", "update")
		assert.Exactly(t, `trg_fcd07948031b28888ab2b4959c097ef4`, have)
	})
}

func TestForeignKeyName(t *testing.T) {
	t.Parallel()

	t.Run("unique short", func(t *testing.T) {
		have := csdb.ForeignKeyName("catalog_product", "parent_id", "catalog_product_entity", "entity_id")
		assert.Exactly(t, `CATALOG_PRODUCT_PARENT_ID_CATALOG_PRODUCT_ENTITY_ENTITY_ID`, have)
	})

	t.Run("unique abbreviated", func(t *testing.T) {
		have := csdb.ForeignKeyName("catalog_product_bundle_option", "parent_id", "catalog_product_entity", "entity_id")
		assert.Exactly(t, `CAT_PRD_BNDL_OPT_PARENT_ID_CAT_PRD_ENTT_ENTT_ID`, have)
	})

	t.Run("unique hashed", func(t *testing.T) {
		have := csdb.ForeignKeyName("catalog_product_bundle_optioncatalog_product_bundle_optioncatalog_product_bundle_option", "parent_id", "catalog_product_entitycatalog_product_entitycatalog_product_entity", "entity_id")
		assert.Exactly(t, `FK_8FC0390CFB720B81470F95BE9E5A8584`, have)
	})
}

func BenchmarkIndexName(b *testing.B) {
	b.ReportAllocs()
	b.Run("unique short", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := csdb.IndexName("unique", "sales_invoiced_aggregated_order", "period", "store_id", "order_status")
			if want := `SALES_INVOICED_AGGREGATED_ORDER_PERIOD_STORE_ID_ORDER_STATUS`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})

	b.Run("unique abbreviated", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := csdb.IndexName("unique", "sales_invoiced_aggregated_order", "customer", "store_id", "order_status", "order_type")
			if want := `SALES_INVOICED_AGGRED_ORDER_CSTR_STORE_ID_ORDER_STS_ORDER_TYPE`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})

	b.Run("unique hashed", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := csdb.IndexName("unique", "sales_invoiced_aggregated_order", "period", "store_id", "order_status", "order_type", "order_date")
			if want := `UNQ_26EE326A968C157BC5004C8206E082E2`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
}

func BenchmarkTableName(b *testing.B) {
	b.ReportAllocs()

	b.Run("Short with prefix suffix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := csdb.TableName("xyz_", "catalog_product_€ntity", "int")
			if want := `xyz_catalog_product_ntity_int`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
	b.Run("Short with prefix without suffix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := csdb.TableName("xyz_", "catalog_product_€ntity")
			if want := `xyz_catalog_product_ntity`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
	b.Run("Short without prefix and suffix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := csdb.TableName("", "catalog_product_€ntity")
			if want := `catalog_product_ntity`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
	b.Run("abbreviated", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := csdb.TableName("", "catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity")
			if want := `cat_prd_ntity_cat_prd_ntity_cstr_cat_prd_ntity`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
	b.Run("hashed", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := csdb.TableName("", "catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity_catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity")
			if want := `t_bb0f749c31c69ed73ad028cb61f43745`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})

}
