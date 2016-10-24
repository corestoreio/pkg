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
	"bytes"
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/null"
	"github.com/stretchr/testify/assert"
)

// Check that type adheres to interfaces
var _ fmt.Stringer = (*csdb.Columns)(nil)
var _ fmt.GoStringer = (*csdb.Columns)(nil)
var _ sort.Interface = (*csdb.Columns)(nil)

func TestGetColumnsMage19(t *testing.T) {
	t.Parallel()
	if _, err := csdb.GetDSN(); errors.IsNotFound(err) {
		t.Skipf("Skipping because environment variable %q not found.", csdb.EnvDSN)
	}

	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()

	tests := []struct {
		table          string
		want           string
		wantErr        errors.BehaviourFunc
		wantJoinFields string
	}{
		{"core_config_data",
			"csdb.Columns{\n&csdb.Column{Field:\"config_id\", Pos:1, Default:null.String{}, Null:\"NO\", DataType:\"int\", CharMaxLength:null.Int64{}, Precision:null.Int64From(10), Scale:null.Int64From(0), TypeRaw:\"int(10) unsigned\", Key:\"PRI\", Extra:\"auto_increment\", Comment:\"Config Id\"},\n&csdb.Column{Field:\"scope\", Pos:2, Default:null.StringFrom(`default`), Null:\"NO\", DataType:\"varchar\", CharMaxLength:null.Int64From(8), Precision:null.Int64{}, Scale:null.Int64{}, TypeRaw:\"varchar(8)\", Key:\"MUL\", Extra:\"\", Comment:\"Config Scope\"},\n&csdb.Column{Field:\"scope_id\", Pos:3, Default:null.StringFrom(`0`), Null:\"NO\", DataType:\"int\", CharMaxLength:null.Int64{}, Precision:null.Int64From(10), Scale:null.Int64From(0), TypeRaw:\"int(11)\", Key:\"\", Extra:\"\", Comment:\"Config Scope Id\"},\n&csdb.Column{Field:\"path\", Pos:4, Default:null.StringFrom(`general`), Null:\"NO\", DataType:\"varchar\", CharMaxLength:null.Int64From(255), Precision:null.Int64{}, Scale:null.Int64{}, TypeRaw:\"varchar(255)\", Key:\"\", Extra:\"\", Comment:\"Config Path\"},\n&csdb.Column{Field:\"value\", Pos:5, Default:null.String{}, Null:\"YES\", DataType:\"text\", CharMaxLength:null.Int64From(65535), Precision:null.Int64{}, Scale:null.Int64{}, TypeRaw:\"text\", Key:\"\", Extra:\"\", Comment:\"Config Value\"},\n}\n",
			nil,
			"config_id_scope_scope_id_path_value",
		},
		{"catalog_category_product",
			"csdb.Columns{\n&csdb.Column{Field:\"category_id\", Pos:1, Default:null.StringFrom(`0`), Null:\"NO\", DataType:\"int\", CharMaxLength:null.Int64{}, Precision:null.Int64From(10), Scale:null.Int64From(0), TypeRaw:\"int(10) unsigned\", Key:\"PRI\", Extra:\"\", Comment:\"Category ID\"},\n&csdb.Column{Field:\"product_id\", Pos:2, Default:null.StringFrom(`0`), Null:\"NO\", DataType:\"int\", CharMaxLength:null.Int64{}, Precision:null.Int64From(10), Scale:null.Int64From(0), TypeRaw:\"int(10) unsigned\", Key:\"PRI\", Extra:\"\", Comment:\"Product ID\"},\n&csdb.Column{Field:\"position\", Pos:3, Default:null.StringFrom(`0`), Null:\"NO\", DataType:\"int\", CharMaxLength:null.Int64{}, Precision:null.Int64From(10), Scale:null.Int64From(0), TypeRaw:\"int(11)\", Key:\"\", Extra:\"\", Comment:\"Position\"},\n}\n",
			nil,
			"category_id_product_id_position",
		},
		{"non_existent",
			"",
			errors.IsNotFound,
			"",
		},
	}

	for i, test := range tests {
		cols1, err := csdb.LoadColumns(context.TODO(), dbc.DB, test.table)
		if test.wantErr != nil {
			assert.Error(t, err, "Index %d => %+v", i, err)
			assert.True(t, test.wantErr(err), "Index %d", i)
		} else {
			assert.NoError(t, err, "Index %d => %+v", i, err)
			assert.Equal(t, test.want, fmt.Sprintf("%#v\n", cols1), "Index %d", i)
			assert.Equal(t, test.wantJoinFields, cols1.JoinFields("_"), "Index %d", i)
		}
	}
}

func TestColumns(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have  int
		want  int
		haveS string
		wantS string
	}{
		{
			mustStructure(table1).Columns.PrimaryKeys().Len(),
			0,
			mustStructure(table1).Columns.GoString(),
			"csdb.Columns{\n&csdb.Column{Field:\"category_id\", Pos:0, Default:null.StringFrom(`0`), Null:\"\", DataType:\"\", CharMaxLength:null.Int64{}, Precision:null.Int64{}, Scale:null.Int64{}, TypeRaw:\"int(10) unsigned\", Key:\"MUL\", Extra:\"\", Comment:\"\"},\n&csdb.Column{Field:\"path\", Pos:0, Default:null.String{}, Null:\"YES\", DataType:\"\", CharMaxLength:null.Int64{}, Precision:null.Int64{}, Scale:null.Int64{}, TypeRaw:\"varchar(255)\", Key:\"MUL\", Extra:\"\", Comment:\"\"},\n}",
		},
		{
			mustStructure(table2).Columns.PrimaryKeys().Len(),
			1,
			mustStructure(table2).Columns.GoString(),
			"csdb.Columns{\n&csdb.Column{Field:\"category_id\", Pos:0, Default:null.StringFrom(`0`), Null:\"\", DataType:\"\", CharMaxLength:null.Int64{}, Precision:null.Int64{}, Scale:null.Int64{}, TypeRaw:\"int(10) unsigned\", Key:\"PRI\", Extra:\"\", Comment:\"\"},\n&csdb.Column{Field:\"path\", Pos:0, Default:null.String{}, Null:\"YES\", DataType:\"\", CharMaxLength:null.Int64{}, Precision:null.Int64{}, Scale:null.Int64{}, TypeRaw:\"varchar(255)\", Key:\"\", Extra:\"\", Comment:\"\"},\n}",
		},
		{
			mustStructure(table4).Columns.UniqueKeys().Len(), 1,
			mustStructure(table4).Columns.GoString(),
			"csdb.Columns{\n&csdb.Column{Field:\"user_id\", Pos:0, Default:null.String{}, Null:\"\", DataType:\"\", CharMaxLength:null.Int64{}, Precision:null.Int64{}, Scale:null.Int64{}, TypeRaw:\"int(10) unsigned\", Key:\"PRI\", Extra:\"auto_increment\", Comment:\"\"},\n&csdb.Column{Field:\"email\", Pos:0, Default:null.String{}, Null:\"YES\", DataType:\"\", CharMaxLength:null.Int64{}, Precision:null.Int64{}, Scale:null.Int64{}, TypeRaw:\"varchar(128)\", Key:\"\", Extra:\"\", Comment:\"\"},\n&csdb.Column{Field:\"username\", Pos:0, Default:null.String{}, Null:\"YES\", DataType:\"\", CharMaxLength:null.Int64{}, Precision:null.Int64{}, Scale:null.Int64{}, TypeRaw:\"varchar(40)\", Key:\"UNI\", Extra:\"\", Comment:\"\"},\n}",
		},
		{mustStructure(table4).Columns.PrimaryKeys().Len(), 1, "", ""},
	}

	for i, test := range tests {
		assert.Equal(t, test.want, test.have, "Incorrect length at index %d", i)
		assert.Equal(t, test.wantS, test.haveS, "Index %d", i)
	}

	tsN := mustStructure(table4).Columns.ByName("user_id_not_found")
	assert.NotNil(t, tsN)
	assert.Empty(t, tsN.Field)

	ts4 := mustStructure(table4).Columns.ByName("user_id")
	assert.NotEmpty(t, ts4.Field)
	assert.True(t, ts4.IsAutoIncrement())

	ts4b := mustStructure(table4).Columns.ByName("email")
	assert.NotEmpty(t, ts4b.Field)
	assert.True(t, ts4b.IsNull())

	assert.True(t, mustStructure(table4).Columns.First().IsPK())
	emptyTS := &csdb.Table{}
	assert.False(t, emptyTS.Columns.First().IsPK())

	hash, err := mustStructure(table3).Columns.Hash()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x3b, 0x72, 0x14, 0x1d, 0x3f, 0x61, 0xf, 0x5b}, hash)
}

func TestColumnsMap(t *testing.T) {
	t.Parallel()
	cols := csdb.Columns{
		&csdb.Column{
			Field:   (`category_id131`),
			TypeRaw: (`int(10) unsigned`),
			Key:     (`PRI`),
			Default: null.StringFrom(`0`),
			Extra:   (``),
		},
		&csdb.Column{
			Field:   (`is_root_category180`),
			TypeRaw: (`smallint(2) unsigned`),
			Null:    csdb.ColumnNull,
			Key:     (``),
			Default: null.StringFrom(`0`),
			Extra:   (``),
		},
	}
	colsHave := cols.Map(func(c *csdb.Column) *csdb.Column {
		c.Field = c.Field + "2"
		return c
	})

	colsWant := csdb.Columns{
		&csdb.Column{Field: (`category_id1312`),
			TypeRaw: (`int(10) unsigned`), Key: (`PRI`),
			Default: null.StringFrom(`0`), Extra: (``)},
		&csdb.Column{Field: (`is_root_category1802`),
			TypeRaw: (`smallint(2) unsigned`), Null: csdb.ColumnNull,
			Key: (``), Default: null.StringFrom(`0`), Extra: (``)},
	}

	assert.Exactly(t, colsWant, colsHave)
}

func TestColumnsFilter(t *testing.T) {
	t.Parallel()
	cols := csdb.Columns{
		&csdb.Column{
			Field:   (`category_id131`),
			TypeRaw: (`int(10) unsigned`),
			Key:     (`PRI`),
			Default: null.StringFrom(`0`),
			Extra:   (``),
		},
		&csdb.Column{
			Field:   (`is_root_category180`),
			TypeRaw: (`smallint(2) unsigned`),
			Null:    csdb.ColumnNull,
			Key:     (``),
			Default: null.StringFrom(`0`),
			Extra:   (``),
		},
	}
	colsHave := cols.Filter(func(c *csdb.Column) bool {
		return c.Field == "is_root_category180"
	})
	colsWant := csdb.Columns{
		&csdb.Column{Field: (`is_root_category180`), TypeRaw: (`smallint(2) unsigned`), Null: csdb.ColumnNull, Key: (``), Default: null.StringFrom(`0`), Extra: (``)},
	}

	assert.Exactly(t, colsWant, colsHave)
}

func TestGetGoPrimitive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		c           csdb.Column
		useNullType bool
		want        string
	}{
		{
			csdb.Column{
				Field:    (`category_id131`),
				DataType: `int`,
				Key:      (`PRI`),
				Default:  null.StringFrom(`0`),
				Extra:    (``),
			},
			false,
			"int64",
		},
		{
			csdb.Column{
				Field:    (`category_id143`),
				DataType: (`int`),
				Null:     csdb.ColumnNull,
				Key:      (`PRI`),
				Default:  null.StringFrom(`0`),
				Extra:    (``),
			},
			false,
			"int64",
		},
		{
			csdb.Column{
				Field:    (`category_id155`),
				DataType: (`int`),
				Null:     csdb.ColumnNull,
				Key:      (`PRI`),
				Default:  null.StringFrom(`0`),
				Extra:    (``),
			},
			true,
			"null.Int64",
		},

		{
			csdb.Column{
				Field:    (`is_root_category155`),
				DataType: (`smallint`),
				Null:     csdb.ColumnNull,
				Key:      (``),
				Default:  null.StringFrom(`0`),
				Extra:    (``),
			},
			false,
			"bool",
		},
		{
			csdb.Column{
				Field:    (`is_root_category180`),
				DataType: (`smallint`),
				Null:     csdb.ColumnNull,
				Key:      (``),
				Default:  null.StringFrom(`0`),
				Extra:    (``),
			},
			true,
			"null.Bool",
		},

		{
			csdb.Column{
				Field:    (`product_name193`),
				DataType: (`varchar`),
				Null:     csdb.ColumnNull,
				Key:      (``),
				Default:  null.StringFrom(`0`),
				Extra:    (``),
			},
			true,
			"null.String",
		},
		{
			csdb.Column{
				Field:    (`product_name193`),
				DataType: (`varchar`),
				Null:     csdb.ColumnNull,
			},
			false,
			"string",
		},

		{
			csdb.Column{
				Field:    (`price`),
				DataType: (`decimal`),
				Null:     csdb.ColumnNull,
			},
			false,
			"money.Money",
		},
		{
			csdb.Column{
				Field:    (`shipping_adjustment_230`),
				DataType: (`decimal`),
				Null:     csdb.ColumnNull,
			},
			true,
			"money.Money",
		},
		{
			csdb.Column{
				Field:    (`grand_absolut_233`),
				DataType: (`decimal`),
				Null:     csdb.ColumnNull,
			},
			true,
			"money.Money",
		},
		{
			csdb.Column{
				Field:    (`some_currencies_242`),
				DataType: (`decimal`),
				Default:  null.StringFrom(`0.0000`),
			},
			true,
			"money.Money",
		},

		{
			csdb.Column{
				Field:    (`weight_252`),
				DataType: (`decimal`),
				Null:     csdb.ColumnNull,
				Default:  null.StringFrom(`0.0000`),
			},
			true,
			"null.Float64",
		},
		{
			csdb.Column{
				Field:    (`weight_263`),
				DataType: (`double`),
				Null:     csdb.ColumnNull,
				Default:  null.StringFrom(`0.0000`),
			},
			false,
			"float64",
		},

		{
			csdb.Column{
				Field:    (`created_at_274`),
				DataType: (`date`),
				Null:     csdb.ColumnNull,
				Default:  null.StringFrom(`0000-00-00`),
			},
			false,
			"time.Time",
		},
		{
			csdb.Column{
				Field:    (`created_at_274`),
				DataType: (`date`),
				Null:     csdb.ColumnNull,
				Default:  null.StringFrom(`0000-00-00`),
			},
			true,
			"null.Time",
		},
	}

	for i, test := range tests {
		var have string
		if test.useNullType {
			have = test.c.GoPrimitiveNull()
		} else {
			have = test.c.GoPrimitive()
		}
		assert.Equal(t, test.want, have, "Index %d", i)
	}
}

var adminUserColumns = csdb.Columns{
	&csdb.Column{Field: "user_id", Pos: 1, Default: null.String{}, Null: "NO", DataType: "int", CharMaxLength: null.Int64{}, Precision: null.Int64From(10), Scale: null.Int64From(0), TypeRaw: "int(10) unsigned", Key: "PRI", Extra: "auto_increment", Comment: "User ID"},
	&csdb.Column{Field: "firstname", Pos: 2, Default: null.String{}, Null: "YES", DataType: "varchar", CharMaxLength: null.Int64From(32), Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "varchar(32)", Key: "", Extra: "", Comment: "User First Name"},
	&csdb.Column{Field: "lastname", Pos: 3, Default: null.String{}, Null: "YES", DataType: "varchar", CharMaxLength: null.Int64From(32), Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "varchar(32)", Key: "", Extra: "", Comment: "User Last Name"},
	&csdb.Column{Field: "email", Pos: 4, Default: null.String{}, Null: "YES", DataType: "varchar", CharMaxLength: null.Int64From(128), Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "varchar(128)", Key: "", Extra: "", Comment: "User Email"},
	&csdb.Column{Field: "username", Pos: 5, Default: null.String{}, Null: "YES", DataType: "varchar", CharMaxLength: null.Int64From(40), Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "varchar(40)", Key: "UNI", Extra: "", Comment: "User Login"},
	&csdb.Column{Field: "password", Pos: 6, Default: null.String{}, Null: "NO", DataType: "varchar", CharMaxLength: null.Int64From(255), Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "varchar(255)", Key: "", Extra: "", Comment: "User Password"},
	&csdb.Column{Field: "created", Pos: 7, Default: null.StringFrom(`0000-00-00 00:00:00`), Null: "NO", DataType: "timestamp", CharMaxLength: null.Int64{}, Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "timestamp", Key: "", Extra: "", Comment: "User Created Time"},
	&csdb.Column{Field: "modified", Pos: 8, Default: null.StringFrom(`CURRENT_TIMESTAMP`), Null: "NO", DataType: "timestamp", CharMaxLength: null.Int64{}, Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "timestamp", Key: "", Extra: "on update CURRENT_TIMESTAMP", Comment: "User Modified Time"},
	&csdb.Column{Field: "logdate", Pos: 9, Default: null.String{}, Null: "YES", DataType: "timestamp", CharMaxLength: null.Int64{}, Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "timestamp", Key: "", Extra: "", Comment: "User Last Login Time"},
	&csdb.Column{Field: "lognum", Pos: 10, Default: null.StringFrom(`0`), Null: "NO", DataType: "smallint", CharMaxLength: null.Int64{}, Precision: null.Int64From(5), Scale: null.Int64From(0), TypeRaw: "smallint(5) unsigned", Key: "", Extra: "", Comment: "User Login Number"},
	&csdb.Column{Field: "reload_acl_flag", Pos: 11, Default: null.StringFrom(`0`), Null: "NO", DataType: "smallint", CharMaxLength: null.Int64{}, Precision: null.Int64From(5), Scale: null.Int64From(0), TypeRaw: "smallint(6)", Key: "", Extra: "", Comment: "Reload ACL"},
	&csdb.Column{Field: "is_active", Pos: 12, Default: null.StringFrom(`1`), Null: "NO", DataType: "smallint", CharMaxLength: null.Int64{}, Precision: null.Int64From(5), Scale: null.Int64From(0), TypeRaw: "smallint(6)", Key: "", Extra: "", Comment: "User Is Active"},
	&csdb.Column{Field: "extra", Pos: 13, Default: null.String{}, Null: "YES", DataType: "text", CharMaxLength: null.Int64From(65535), Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "text", Key: "", Extra: "", Comment: "User Extra Data"},
	&csdb.Column{Field: "rp_token", Pos: 14, Default: null.String{}, Null: "YES", DataType: "text", CharMaxLength: null.Int64From(65535), Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "text", Key: "", Extra: "", Comment: "Reset Password Link Token"},
	&csdb.Column{Field: "rp_token_created_at", Pos: 15, Default: null.String{}, Null: "YES", DataType: "timestamp", CharMaxLength: null.Int64{}, Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "timestamp", Key: "", Extra: "", Comment: "Reset Password Link Token Creation Date"},
	&csdb.Column{Field: "interface_locale", Pos: 16, Default: null.StringFrom(`en_US`), Null: "NO", DataType: "varchar", CharMaxLength: null.Int64From(16), Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "varchar(16)", Key: "", Extra: "", Comment: "Backend interface locale"},
	&csdb.Column{Field: "failures_num", Pos: 17, Default: null.StringFrom(`0`), Null: "YES", DataType: "smallint", CharMaxLength: null.Int64{}, Precision: null.Int64From(5), Scale: null.Int64From(0), TypeRaw: "smallint(6)", Key: "", Extra: "", Comment: "Failure Number"},
	&csdb.Column{Field: "first_failure", Pos: 18, Default: null.String{}, Null: "YES", DataType: "timestamp", CharMaxLength: null.Int64{}, Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "timestamp", Key: "", Extra: "", Comment: "First Failure"},
	&csdb.Column{Field: "lock_expires", Pos: 19, Default: null.String{}, Null: "YES", DataType: "timestamp", CharMaxLength: null.Int64{}, Precision: null.Int64{}, Scale: null.Int64{}, TypeRaw: "timestamp", Key: "", Extra: "", Comment: "Expiration Lock Dates"},
}

func TestColumnsSort(t *testing.T) {
	t.Parallel()
	//sort.Reverse(adminUserColumns) doesn't work and not yet needed
	sort.Sort(adminUserColumns)
	assert.Exactly(t, `user_id`, adminUserColumns.First().Field)
}

func TestColumn_IsUnsigned(t *testing.T) {
	t.Parallel()
	assert.True(t, adminUserColumns.ByName("lognum").IsUnsigned())
	assert.False(t, adminUserColumns.ByName("reload_acl_flag").IsUnsigned())
}

func TestColumn_IsDate(t *testing.T) {
	t.Parallel()
	assert.True(t, adminUserColumns.ByName("logdate").IsDate())
	assert.False(t, adminUserColumns.ByName("reload_acl_flag").IsDate())
}

func TestColumn_IsCurrentTimestamp(t *testing.T) {
	t.Parallel()
	assert.True(t, adminUserColumns.ByName("modified").IsCurrentTimestamp())
	assert.False(t, adminUserColumns.ByName("reload_acl_flag").IsCurrentTimestamp())
}

var benchmarkGetColumns csdb.Columns
var benchmarkGetColumnsHashWant = []byte{0x3b, 0x2d, 0xdd, 0xf4, 0x4e, 0x2b, 0x3a, 0xd0}

// BenchmarkGetColumns-4       	5000	    395152 ns/op	   21426 B/op	     179 allocs/op
func BenchmarkGetColumns(b *testing.B) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	var err error
	ctx := context.TODO()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkGetColumns, err = csdb.LoadColumns(ctx, dbc.DB, "eav_attribute")
		if err != nil {
			b.Error(err)
		}
	}
	hashHave, err := benchmarkGetColumns.Hash()
	if err != nil {
		b.Error(err)
	}
	if 0 != bytes.Compare(hashHave, benchmarkGetColumnsHashWant) {
		b.Errorf("\nHave %#v\nWant %#v\n", hashHave, benchmarkGetColumnsHashWant)
	}
	//	b.Log(benchmarkGetColumns.GoString())
}

var benchmarkColumnsJoinFields string
var benchmarkColumnsJoinFieldsWant = "category_id|product_id|position"
var benchmarkColumnsJoinFieldsData = csdb.Columns{
	&csdb.Column{
		Field:   "category_id",
		TypeRaw: ("int(10) unsigned"),
		Key:     "",
		Default: null.StringFrom("0"),
		Extra:   (""),
	},
	&csdb.Column{
		Field:   "product_id",
		TypeRaw: ("int(10) unsigned"),
		Key:     (""),
		Default: null.StringFrom("0"),
		Extra:   (""),
	},
	&csdb.Column{
		Field:   "position",
		TypeRaw: ("int(10) unsigned"),
		Null:    csdb.ColumnNull,
		Key:     (""),
		Extra:   (""),
	},
}

// BenchmarkColumnsJoinFields-4	 2000000	       625 ns/op	     176 B/op	       5 allocs/op <- Go 1.5
func BenchmarkColumnsJoinFields(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkColumnsJoinFields = benchmarkColumnsJoinFieldsData.JoinFields("|")
	}
	if benchmarkColumnsJoinFields != benchmarkColumnsJoinFieldsWant {
		b.Errorf("\nWant: %s\nHave: %s\n", benchmarkColumnsJoinFieldsWant, benchmarkColumnsJoinFields)
	}
}
