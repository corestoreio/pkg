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

package dmlgen

import (
	"testing"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

func TestToGoTypeNull(t *testing.T) {
	tests := []struct {
		c    *ddl.Column
		want string
	}{
		{&ddl.Column{Field: `category_id214`, DataType: `bigint`, ColumnType: `bigint unsigned`}, "uint64"},
		{&ddl.Column{Field: `category_id224`, DataType: `bigint`, ColumnType: `bigint`}, "int64"},
		{&ddl.Column{Field: `category_id225`, DataType: `bigint`, Null: "YES", ColumnType: `bigint unsigned`}, "null.Uint64"},
		{&ddl.Column{Field: `category_id225`, DataType: `bigint`, Null: "YES", ColumnType: `bigint`}, "null.Int64"},

		{&ddl.Column{Field: `category_id227a`, DataType: `int`, ColumnType: `int unsigned`}, "uint32"},
		{&ddl.Column{Field: `category_id227b`, DataType: `int`, ColumnType: `int unsigned`, Null: "YES"}, "null.Uint32"},
		{&ddl.Column{Field: `category_id227c`, DataType: `int`, ColumnType: `int`}, "int32"},
		{&ddl.Column{Field: `category_id227d`, DataType: `int`, Null: "YES"}, "null.Int32"},

		{&ddl.Column{Field: `is_root_cat269`, DataType: `smallint`, Null: "YES", Default: null.MakeString(`0`)}, "null.Bool"},
		{&ddl.Column{Field: `is_root_cat180`, DataType: `smallint`, Null: "YES", Default: null.MakeString(`0`)}, "null.Bool"},

		{&ddl.Column{Field: `root_int16_a39`, DataType: `smallint`, Null: "NO"}, "int16"},
		{&ddl.Column{Field: `root_int16_a40`, DataType: `smallint`, Null: "YES"}, "null.Int16"},
		{&ddl.Column{Field: `root_int16_a41`, DataType: `smallint`, ColumnType: `smallint unsigned`, Null: "NO"}, "uint16"},
		{&ddl.Column{Field: `root_int16_a42`, DataType: `smallint`, ColumnType: `smallint unsigned`, Null: "YES"}, "null.Uint16"},

		{&ddl.Column{Field: `root_int8_a50`, DataType: `tinyint`, ColumnType: `tinyint unsigned`, Null: "YES"}, "null.Uint8"},
		{&ddl.Column{Field: `root_int8_a50`, DataType: `tinyint`, ColumnType: `tinyint unsigned`, Null: "NO"}, "uint8"},
		{&ddl.Column{Field: `root_int8_a50`, DataType: `tinyint`, ColumnType: `tinyint`, Null: "NO"}, "int8"},
		{&ddl.Column{Field: `root_int8_a50`, DataType: `tinyint`, ColumnType: `tinyint`, Null: "YES"}, "null.Int8"},

		{&ddl.Column{Field: `product_name193`, DataType: `varchar`, Null: "YES", Default: null.MakeString(`0`)}, "null.String"},
		{&ddl.Column{Field: `product_name193`, DataType: `varchar`, Null: "YES"}, "null.String"},

		{&ddl.Column{Field: `_price_______`, DataType: `decimal`, Null: "YES"}, "null.Decimal"},
		{&ddl.Column{Field: `price`, DataType: `double`, Null: "NO"}, "null.Decimal"},
		{&ddl.Column{Field: `msrp`, DataType: `double`, Null: "NO"}, "null.Decimal"},
		{&ddl.Column{Field: `shipping_adjustment_230`, DataType: `decimal`, Null: "YES"}, "null.Decimal"},
		{&ddl.Column{Field: `shipping_adjustment_241`, DataType: `decimal`, Null: "NO"}, "null.Decimal"},
		{&ddl.Column{Field: `shipping_adjstment_252`, DataType: `decimal`, Null: "YES"}, "null.Decimal"},
		{&ddl.Column{Field: `rate__232`, DataType: `decimal`, Null: "NO"}, "null.Decimal"},
		{&ddl.Column{Field: `rate__233`, DataType: `decimal`, ColumnType: `float unsigned`, Null: "NO"}, "null.Decimal"},
		{&ddl.Column{Field: `grand_absot_233`, DataType: `decimal`, Null: "YES"}, "null.Decimal"},
		{&ddl.Column{Field: `some_currencies_242`, DataType: `decimal`, Default: null.MakeString(`0.0000`)}, "null.Decimal"},
		{&ddl.Column{Field: `weight_252`, DataType: `decimal`, Null: "YES", Default: null.MakeString(`0.0000`)}, "null.Decimal"},

		{&ddl.Column{Field: `weight_263`, DataType: `double`, Default: null.MakeString(`0.0000`)}, "float64"},
		{&ddl.Column{Field: `created_at_674`, DataType: `date`, Default: null.MakeString(`0000-00-00`)}, "time.Time"},
		{&ddl.Column{Field: `created_at_774`, DataType: `date`, Null: "YES", Default: null.MakeString(`0000-00-00`)}, "null.Time"},
		{&ddl.Column{Field: `created_at_874`, DataType: `datetime`, Null: "NO", Default: null.MakeString(`0000-00-00`)}, "time.Time"},
		{&ddl.Column{Field: `image001`, DataType: `varbinary`, Null: "NO"}, "[]byte"},
		{&ddl.Column{Field: `image002`, DataType: `varbinary`, Null: "YES"}, "[]byte"},
		{&ddl.Column{Field: `image003`, DataType: `blob`, Null: "YES"}, "[]byte"},
		{&ddl.Column{Field: `image004`, DataType: `longblob`, Null: "YES"}, "[]byte"},
		{&ddl.Column{Field: `image005`, DataType: `mediumblob`, Null: "YES"}, "[]byte"},
		{&ddl.Column{Field: `ok_dude1`, DataType: `bit`, Null: "NO"}, "bool"},
		{&ddl.Column{Field: `ok_dude2`, DataType: `bit`, Null: "YES"}, "null.Bool"},
		{&ddl.Column{Field: `description_001`, DataType: `varchar`, Null: "YES"}, "null.String"},
		{&ddl.Column{Field: `description_002`, DataType: `varchar`, Null: "NO"}, "string"},
		{&ddl.Column{Field: `description_003`, DataType: `char`, Null: "YES"}, "null.String"},
		{&ddl.Column{Field: `description_004`, DataType: `char`, Null: "NO"}, "string"},
	}
	ts := new(Generator)
	for _, test := range tests {
		have := ts.mySQLToGoType(test.c, true) // including null
		assert.Exactly(t, test.want, have, "%q", test.c.Field)
	}
}

func TestMySQLToGoDmlColumnMap(t *testing.T) {
	tests := []struct {
		c    *ddl.Column
		want string // The function names as mentioned in dml.ColumnMap.[TFunc]
	}{
		{&ddl.Column{Field: `category_id214`, DataType: `bigint`, ColumnType: `bigint unsigned`}, "Uint64"},
		{&ddl.Column{Field: `category_id224`, DataType: `bigint`, ColumnType: `bigint`}, "Int64"},
		{&ddl.Column{Field: `category_id225`, DataType: `bigint`, Null: "YES", ColumnType: `bigint unsigned`}, "NullUint64"},
		{&ddl.Column{Field: `category_id225`, DataType: `bigint`, Null: "YES", ColumnType: `bigint`}, "NullInt64"},

		{&ddl.Column{Field: `category_id227a`, DataType: `int`, ColumnType: `int unsigned`}, "Uint32"},
		{&ddl.Column{Field: `category_id227b`, DataType: `int`, ColumnType: `int unsigned`, Null: "YES"}, "NullUint32"},
		{&ddl.Column{Field: `category_id227c`, DataType: `int`, ColumnType: `int`}, "Int32"},
		{&ddl.Column{Field: `category_id227d`, DataType: `int`, Null: "YES"}, "NullInt32"},

		{&ddl.Column{Field: `is_root_cat269`, DataType: `smallint`, Null: "YES", Default: null.MakeString(`0`)}, "NullBool"},
		{&ddl.Column{Field: `is_root_cat180`, DataType: `smallint`, Null: "YES", Default: null.MakeString(`0`)}, "NullBool"},

		{&ddl.Column{Field: `root_int16_a39`, DataType: `smallint`, Null: "NO"}, "Int16"},
		{&ddl.Column{Field: `root_int16_a40`, DataType: `smallint`, Null: "YES"}, "NullInt16"},
		{&ddl.Column{Field: `root_int16_a41`, DataType: `smallint`, ColumnType: `smallint unsigned`, Null: "NO"}, "Uint16"},
		{&ddl.Column{Field: `root_int16_a42`, DataType: `smallint`, ColumnType: `smallint unsigned`, Null: "YES"}, "NullUint16"},

		{&ddl.Column{Field: `root_int8_a50`, DataType: `tinyint`, ColumnType: `tinyint unsigned`, Null: "YES"}, "NullUint8"},
		{&ddl.Column{Field: `root_int8_a50`, DataType: `tinyint`, ColumnType: `tinyint unsigned`, Null: "NO"}, "Uint8"},
		{&ddl.Column{Field: `root_int8_a50`, DataType: `tinyint`, ColumnType: `tinyint`, Null: "NO"}, "Int8"},
		{&ddl.Column{Field: `root_int8_a50`, DataType: `tinyint`, ColumnType: `tinyint`, Null: "YES"}, "NullInt8"},

		{&ddl.Column{Field: `product_name193`, DataType: `varchar`, Null: "YES"}, "NullString"},

		{&ddl.Column{Field: `category_id236`, DataType: `int`, Default: null.MakeString(`0`)}, "Int32"},
		{&ddl.Column{Field: `category_id247`, DataType: `int`, Null: "YES", Default: null.MakeString(`0`)}, "NullInt32"},
		{&ddl.Column{Field: `category_id258`, DataType: `int`, Null: "YES", Default: null.MakeString(`0`)}, "NullInt32"},
		{&ddl.Column{Field: `is_root_cat269`, DataType: `smallint`, Null: "YES", Default: null.MakeString(`0`)}, "NullBool"},
		{&ddl.Column{Field: `is_root_cat180`, DataType: `smallint`, Null: "YES", Default: null.MakeString(`0`)}, "NullBool"},
		{&ddl.Column{Field: `product_name193`, DataType: `varchar`, Null: "YES", Default: null.MakeString(`0`)}, "NullString"},
		{&ddl.Column{Field: `product_name193`, DataType: `varchar`, Null: "YES"}, "NullString"},
		{&ddl.Column{Field: `_price_______`, DataType: `decimal`, Null: "YES"}, "Decimal"},
		{&ddl.Column{Field: `price`, DataType: `double`, Null: "NO"}, "Decimal"},
		{&ddl.Column{Field: `msrp`, DataType: `double`, Null: "NO"}, "Decimal"},
		{&ddl.Column{Field: `shipping_adjustment_230`, DataType: `decimal`, Null: "YES"}, "Decimal"},
		{&ddl.Column{Field: `shipping_adjustment_241`, DataType: `decimal`, Null: "NO"}, "Decimal"},
		{&ddl.Column{Field: `shipping_adjstment_252`, DataType: `decimal`, Null: "YES"}, "Decimal"},
		{&ddl.Column{Field: `rate__232`, DataType: `decimal`, Null: "NO"}, "Decimal"},
		{&ddl.Column{Field: `rate__233`, DataType: `decimal`, ColumnType: `float unsigned`, Null: "NO"}, "Decimal"},
		{&ddl.Column{Field: `grand_absot_233`, DataType: `decimal`, Null: "YES"}, "Decimal"},
		{&ddl.Column{Field: `some_currencies_242`, DataType: `decimal`, Default: null.MakeString(`0.0000`)}, "Decimal"},
		{&ddl.Column{Field: `weight_252`, DataType: `decimal`, Null: "YES", Default: null.MakeString(`0.0000`)}, "Decimal"},
		{&ddl.Column{Field: `weight_263`, DataType: `double`, Default: null.MakeString(`0.0000`)}, "Float64"},
		{&ddl.Column{Field: `created_at_674`, DataType: `date`, Default: null.MakeString(`0000-00-00`)}, "Time"},
		{&ddl.Column{Field: `created_at_774`, DataType: `date`, Null: "YES", Default: null.MakeString(`0000-00-00`)}, "NullTime"},
		{&ddl.Column{Field: `created_at_874`, DataType: `datetime`, Null: "NO", Default: null.MakeString(`0000-00-00`)}, "Time"},
		{&ddl.Column{Field: `image001`, DataType: `varbinary`, Null: "NO"}, "Byte"},
		{&ddl.Column{Field: `image002`, DataType: `varbinary`, Null: "YES"}, "Byte"},
		{&ddl.Column{Field: `ok_dude1`, DataType: `bit`, Null: "NO"}, "Bool"},
		{&ddl.Column{Field: `ok_dude2`, DataType: `bit`, Null: "YES"}, "NullBool"},
		{&ddl.Column{Field: `description_001`, DataType: `varchar`, Null: "YES"}, "NullString"},
		{&ddl.Column{Field: `description_002`, DataType: `varchar`, Null: "NO"}, "String"},
		{&ddl.Column{Field: `description_003`, DataType: `char`, Null: "YES"}, "NullString"},
		{&ddl.Column{Field: `description_004`, DataType: `char`, Null: "NO"}, "String"},
	}
	ts := new(Generator)
	for i, test := range tests {
		have := ts.mySQLToGoDmlColumnMap(test.c, true) // including null
		assert.Exactly(t, test.want, have, "IDX:%d %#v", i, test.c)
	}
}

func TestToGoPrimitive(t *testing.T) {
	tests := []struct {
		want string // The function names as mentioned in dml.ColumnMap.[TFunc]
		c    *ddl.Column
	}{
		{"ColInt1.Int32", &ddl.Column{Field: "col_int_1", Pos: 19, Default: null.MakeString("NULL"), Null: "YES", DataType: "int", Precision: null.MakeInt64(10), Scale: null.MakeInt64(0), ColumnType: "int(10)", Uniquified: true, Generated: "NEVER"}},
		{"CategoryId214", &ddl.Column{Field: `category_id214`, DataType: `bigint`, ColumnType: `bigint unsigned`}},
		{"ID", &ddl.Column{Field: "id", Pos: 1, Null: "NO", DataType: "int", Precision: null.MakeInt64(10), Scale: null.MakeInt64(0), ColumnType: "int(11)", Key: "PRI", Extra: "auto_increment", Generated: "NEVER"}},
		{"ColBigint1.Int64", &ddl.Column{Field: "col_bigint_1", Pos: 2, Default: null.MakeString("NULL"), Null: "YES", DataType: "bigint", Precision: null.MakeInt64(19), Scale: null.MakeInt64(0), ColumnType: "bigint(20)", Generated: "NEVER"}},
		{"ColBigint2", &ddl.Column{Field: "col_bigint_2", Pos: 3, Default: null.MakeString("0"), Null: "NO", DataType: "bigint", Precision: null.MakeInt64(19), Scale: null.MakeInt64(0), ColumnType: "bigint(20)", Generated: "NEVER"}},
		{"ColBigint3.Uint64", &ddl.Column{Field: "col_bigint_3", Pos: 4, Default: null.MakeString("NULL"), Null: "YES", DataType: "bigint", Precision: null.MakeInt64(20), Scale: null.MakeInt64(0), ColumnType: "bigint(20) unsigned", Generated: "NEVER"}},
		{"ColBigint4", &ddl.Column{Field: "col_bigint_4", Pos: 5, Default: null.MakeString("0"), Null: "NO", DataType: "bigint", Precision: null.MakeInt64(20), Scale: null.MakeInt64(0), ColumnType: "bigint(20) unsigned", Generated: "NEVER"}},
		{"ColBlob", &ddl.Column{Field: "col_blob", Pos: 6, Default: null.MakeString("NULL"), Null: "YES", DataType: "blob", CharMaxLength: null.MakeInt64(65535), ColumnType: "blob", Generated: "NEVER"}},
		{"ColMediumText.Data", &ddl.Column{Field: "col_medium_text", Pos: 6, Default: null.MakeString("NULL"), Null: "YES", DataType: "mediumtext", CharMaxLength: null.MakeInt64(65535), ColumnType: "blob", Generated: "NEVER"}},
		{"ColDate1.Time", &ddl.Column{Field: "col_date_1", Pos: 7, Default: null.MakeString("NULL"), Null: "YES", DataType: "date", ColumnType: "date", Generated: "NEVER"}},
		{"ColDate2", &ddl.Column{Field: "col_date_2", Pos: 8, Default: null.MakeString("'0000-00-00'"), Null: "NO", DataType: "date", ColumnType: "date", Generated: "NEVER"}},
		{"ColDatetime1.Time", &ddl.Column{Field: "col_datetime_1", Pos: 9, Default: null.MakeString("NULL"), Null: "YES", DataType: "datetime", ColumnType: "datetime", Generated: "NEVER"}},
		{"ColDecimal100", &ddl.Column{Field: "col_decimal_10_0", Pos: 11, Default: null.MakeString("NULL"), Null: "YES", DataType: "decimal", Precision: null.MakeInt64(10), Scale: null.MakeInt64(0), ColumnType: "decimal(10,0) unsigned", Generated: "NEVER"}},
		{"ColDecimal124", &ddl.Column{Field: "col_decimal_12_4", Pos: 12, Default: null.MakeString("NULL"), Null: "YES", DataType: "decimal", Precision: null.MakeInt64(12), Scale: null.MakeInt64(4), ColumnType: "decimal(12,4)", Generated: "NEVER"}},
		{"Price124a", &ddl.Column{Field: "price_12_4a", Pos: 13, Default: null.MakeString("NULL"), Null: "YES", DataType: "decimal", Precision: null.MakeInt64(12), Scale: null.MakeInt64(4), ColumnType: "decimal(12,4)", Generated: "NEVER"}},
	}
	ts := new(Generator)
	for i, test := range tests {
		have := ts.toGoPrimitiveFromNull(test.c)
		assert.Exactly(t, test.want, have, "IDX:%d %#v", i, test.c)
	}
}
