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

package ddl_test

import (
	"bytes"
	"context"
	"testing"

	"hash/fnv"

	"github.com/corestoreio/csfw/sql/ddl"
	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/csfw/util/byteconv"
	"github.com/corestoreio/csfw/util/cstesting"
)

func BenchmarkTableName(b *testing.B) {
	b.ReportAllocs()

	b.Run("Short with prefix suffix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := ddl.TableName("xyz_", "catalog_product_€ntity", "int")
			if want := `xyz_catalog_product_ntity_int`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
	b.Run("Short with prefix without suffix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := ddl.TableName("xyz_", "catalog_product_€ntity")
			if want := `xyz_catalog_product_ntity`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
	b.Run("Short without prefix and suffix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := ddl.TableName("", "catalog_product_€ntity")
			if want := `catalog_product_ntity`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
	b.Run("abbreviated", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := ddl.TableName("", "catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity")
			if want := `cat_prd_ntity_cat_prd_ntity_cstr_cat_prd_ntity`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
	b.Run("hashed", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := ddl.TableName("", "catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity_catalog_product_€ntity_catalog_product_€ntity_customer_catalog_product_€ntity")
			if want := `t_bb0f749c31c69ed73ad028cb61f43745`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
}

func BenchmarkIndexName(b *testing.B) {
	b.ReportAllocs()
	b.Run("unique short", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := ddl.IndexName("unique", "sales_invoiced_aggregated_order", "period", "store_id", "order_status")
			if want := `SALES_INVOICED_AGGREGATED_ORDER_PERIOD_STORE_ID_ORDER_STATUS`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})

	b.Run("unique abbreviated", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := ddl.IndexName("unique", "sales_invoiced_aggregated_order", "customer", "store_id", "order_status", "order_type")
			if want := `SALES_INVOICED_AGGRED_ORDER_CSTR_STORE_ID_ORDER_STS_ORDER_TYPE`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})

	b.Run("unique hashed", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			have := ddl.IndexName("unique", "sales_invoiced_aggregated_order", "period", "store_id", "order_status", "order_type", "order_date")
			if want := `UNQ_26EE326A968C157BC5004C8206E082E2`; have != want {
				b.Fatalf("\nHave %q\nWant %q", have, want)
			}
		}
	})
}

var benchmarkIsValidIdentifier error

func BenchmarkIsValidIdentifier(b *testing.B) {
	const id = `$catalog_product_3ntity`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkIsValidIdentifier = dml.IsValidIdentifier(id)
	}
	if benchmarkIsValidIdentifier != nil {
		b.Fatalf("%+v", benchmarkIsValidIdentifier)
	}
}

var benchmarkColumnsJoinFields string
var benchmarkColumnsJoinFieldsWant = "category_id|product_id|position"
var benchmarkColumnsJoinFieldsData = ddl.Columns{
	&ddl.Column{
		Field:      "category_id",
		ColumnType: "int(10) unsigned",
		Key:        "",
		Default:    dml.MakeNullString("0"),
		Extra:      "",
	},
	&ddl.Column{
		Field:      "product_id",
		ColumnType: "int(10) unsigned",
		Key:        "",
		Default:    dml.MakeNullString("0"),
		Extra:      "",
	},
	&ddl.Column{
		Field:      "position",
		ColumnType: "int(10) unsigned",
		Null:       "YES",
		Key:        "",
		Extra:      "",
	},
}

func BenchmarkColumnsJoinFields(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkColumnsJoinFields = benchmarkColumnsJoinFieldsData.JoinFields("|")
	}
	if benchmarkColumnsJoinFields != benchmarkColumnsJoinFieldsWant {
		b.Errorf("\nWant: %s\nHave: %s\n", benchmarkColumnsJoinFieldsWant, benchmarkColumnsJoinFields)
	}
}

var benchmarkLoadColumns map[string]ddl.Columns
var benchmarkLoadColumnsHashWant = []byte{0x95, 0x7d, 0x84, 0x3c, 0x1f, 0x8, 0x94, 0x2e}

func BenchmarkLoadColumns(b *testing.B) {
	const tn = "eav_attribute"
	ctx := context.TODO()
	db := cstesting.MustConnectDB(b)
	defer cstesting.Close(b, db)

	byteconv.UseStdLib = false

	var err error
	b.ReportAllocs()
	b.ResetTimer()
	b.Run("RowConvert", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchmarkLoadColumns, err = ddl.LoadColumns(ctx, db.DB, tn)
			if err != nil {
				b.Error(err)
			}
		}
		hashHave, err := benchmarkLoadColumns[tn].Hash(fnv.New64a())
		if err != nil {
			b.Error(err)
		}
		if 0 != bytes.Compare(hashHave, benchmarkLoadColumnsHashWant) {
			b.Errorf("\nHave %#v\nWant %#v\n", hashHave, benchmarkLoadColumnsHashWant)
		}
	})
}

// BenchmarkVariables-4   	    2000	   1046318 ns/op	   28401 B/op	    1121 allocs/op <= 186 rows
// BenchmarkVariables-4   	    2000	    651096 ns/op	     769 B/op	      21 allocs/op <= one row!
// BenchmarkVariables-4   	    2000	   1027245 ns/op	   22417 B/op	     935 allocs/op <= pre alloc slice
// BenchmarkVariables-4   	    2000	   1008059 ns/op	   19506 B/op	     750 allocs/op
func BenchmarkVariables(b *testing.B) {

	ctx := context.TODO()

	db := cstesting.MustConnectDB(b)
	defer cstesting.Close(b, db)

	vars := ddl.NewVariables("innodb%")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := dml.Load(ctx, db.DB, vars, vars)
		if err != nil {
			b.Fatalf("%+v", err)
		}
	}

	if "Barracuda" != vars.Data["innodb_file_format"] {
		b.Fatalf("storage_engine variable should be Barracuda, got: %q", vars.Data["innodb_file_format"])
	}
	if ld := len(vars.Data); 186 != ld {
		b.Fatalf("InnoDB Variables should count 186 but found %d", ld)
	}
}
