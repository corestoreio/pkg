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
	"fmt"
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/null"
	"github.com/stretchr/testify/assert"
)

// Check that type adheres to interfaces
var _ fmt.Stringer = (*csdb.Columns)(nil)
var _ fmt.GoStringer = (*csdb.Columns)(nil)

func TestGetColumnsMage19(t *testing.T) {

	if _, err := csdb.GetDSN(); errors.IsNotFound(err) {
		t.Skip("Skipping because no DSN found.")
	}

	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()
	sess := dbc.NewSession()

	tests := []struct {
		table          string
		want           string
		wantErr        error
		wantJoinFields string
	}{
		{"core_config_data",
			"csdb.Column{Field:null.StringFrom(`config_id`), Type:null.StringFrom(`int(10) unsigned`), Null:null.StringFrom(`NO`), Key:null.StringFrom(`PRI`), Default:null.String{}, Extra:null.StringFrom(`auto_increment`)},\ncsdb.Column{Field:null.StringFrom(`scope`), Type:null.StringFrom(`varchar(8)`), Null:null.StringFrom(`NO`), Key:null.StringFrom(`MUL`), Default:null.StringFrom(`default`), Extra:null.StringFrom(``)},\ncsdb.Column{Field:null.StringFrom(`scope_id`), Type:null.StringFrom(`int(11)`), Null:null.StringFrom(`NO`), Key:null.StringFrom(``), Default:null.StringFrom(`0`), Extra:null.StringFrom(``)},\ncsdb.Column{Field:null.StringFrom(`path`), Type:null.StringFrom(`varchar(255)`), Null:null.StringFrom(`NO`), Key:null.StringFrom(``), Default:null.StringFrom(`general`), Extra:null.StringFrom(``)},\ncsdb.Column{Field:null.StringFrom(`value`), Type:null.StringFrom(`text`), Null:null.StringFrom(`YES`), Key:null.StringFrom(``), Default:null.String{}, Extra:null.StringFrom(``)}\n",
			nil,
			"config_id_scope_scope_id_path_value",
		},
		{"catalog_category_product",
			"csdb.Column{Field:null.StringFrom(`category_id`), Type:null.StringFrom(`int(10) unsigned`), Null:null.StringFrom(`NO`), Key:null.StringFrom(`PRI`), Default:null.StringFrom(`0`), Extra:null.StringFrom(``)},\ncsdb.Column{Field:null.StringFrom(`product_id`), Type:null.StringFrom(`int(10) unsigned`), Null:null.StringFrom(`NO`), Key:null.StringFrom(`PRI`), Default:null.StringFrom(`0`), Extra:null.StringFrom(``)},\ncsdb.Column{Field:null.StringFrom(`position`), Type:null.StringFrom(`int(11)`), Null:null.StringFrom(`NO`), Key:null.StringFrom(``), Default:null.StringFrom(`0`), Extra:null.StringFrom(``)}\n",
			nil,
			"category_id_product_id_position",
		},
		{"non_existent",
			"",
			errors.New("non_existent"),
			"",
		},
	}

	for i, test := range tests {
		cols1, err1 := csdb.GetColumns(sess, test.table)
		if test.wantErr != nil {
			assert.Error(t, err1, "Index %d", i)
			assert.Contains(t, err1.Error(), test.wantErr.Error(), "Index %d", i)
		} else {
			assert.NoError(t, err1, "Index %d", i)
			assert.Equal(t, test.want, fmt.Sprintf("%#v\n", cols1), "Index %d", i)
			assert.Equal(t, test.wantJoinFields, cols1.JoinFields("_"), "Index %d", i)
		}
	}
}

func TestColumns(t *testing.T) {

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
			"csdb.Column{Field:null.StringFrom(`category_id`), Type:null.StringFrom(`int(10) unsigned`), Null:null.StringFrom(`NO`), Key:null.StringFrom(`MUL`), Default:null.StringFrom(`0`), Extra:null.StringFrom(``)},\ncsdb.Column{Field:null.StringFrom(`path`), Type:null.StringFrom(`varchar(255)`), Null:null.StringFrom(`YES`), Key:null.StringFrom(`MUL`), Default:null.String{}, Extra:null.StringFrom(``)}",
		},
		{
			mustStructure(table2).Columns.PrimaryKeys().Len(),
			1,
			mustStructure(table2).Columns.GoString(),
			"csdb.Column{Field:null.StringFrom(`category_id`), Type:null.StringFrom(`int(10) unsigned`), Null:null.StringFrom(`NO`), Key:null.StringFrom(`PRI`), Default:null.StringFrom(`0`), Extra:null.StringFrom(``)},\ncsdb.Column{Field:null.StringFrom(`path`), Type:null.StringFrom(`varchar(255)`), Null:null.StringFrom(`YES`), Key:null.String{}, Default:null.String{}, Extra:null.StringFrom(``)}",
		},
		{
			mustStructure(table4).Columns.UniqueKeys().Len(), 1,
			mustStructure(table4).Columns.GoString(),
			"csdb.Column{Field:null.StringFrom(`user_id`), Type:null.StringFrom(`int(10) unsigned`), Null:null.StringFrom(`NO`), Key:null.StringFrom(`PRI`), Default:null.String{}, Extra:null.StringFrom(`auto_increment`)},\ncsdb.Column{Field:null.StringFrom(`email`), Type:null.StringFrom(`varchar(128)`), Null:null.StringFrom(`YES`), Key:null.StringFrom(``), Default:null.String{}, Extra:null.StringFrom(``)},\ncsdb.Column{Field:null.StringFrom(`username`), Type:null.StringFrom(`varchar(40)`), Null:null.StringFrom(`YES`), Key:null.StringFrom(`UNI`), Default:null.String{}, Extra:null.StringFrom(``)}",
		},
		{mustStructure(table4).Columns.PrimaryKeys().Len(), 1, "", ""},
	}

	for i, test := range tests {
		assert.Equal(t, test.want, test.have, "Incorrect length at index %d", i)
		assert.Equal(t, test.wantS, test.haveS, "Index %d", i)
	}

	tsN := mustStructure(table4).Columns.ByName("user_id_not_found")
	assert.NotNil(t, tsN)
	assert.False(t, tsN.Field.Valid)

	ts4 := mustStructure(table4).Columns.ByName("user_id")
	assert.True(t, ts4.Field.Valid)
	assert.True(t, ts4.IsAutoIncrement())

	ts4b := mustStructure(table4).Columns.ByName("email")
	assert.True(t, ts4b.Field.Valid)
	assert.True(t, ts4b.IsNull())

	assert.True(t, mustStructure(table4).Columns.First().IsPK())
	emptyTS := &csdb.Table{}
	assert.False(t, emptyTS.Columns.First().IsPK())

	hash, err := mustStructure(table3).Columns.Hash()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xc7, 0xbc, 0x3b, 0xa5, 0x8f, 0x1e, 0x59, 0x3e}, hash)
}

func TestColumnsMap(t *testing.T) {

	cols := csdb.Columns{
		0: csdb.Column{
			Field:   null.StringFrom(`category_id131`),
			Type:    null.StringFrom(`int(10) unsigned`),
			Null:    null.StringFrom(`NO`),
			Key:     null.StringFrom(`PRI`),
			Default: null.StringFrom(`0`),
			Extra:   null.StringFrom(``),
		},
		1: csdb.Column{
			Field:   null.StringFrom(`is_root_category180`),
			Type:    null.StringFrom(`smallint(2) unsigned`),
			Null:    null.StringFrom(`YES`),
			Key:     null.StringFrom(``),
			Default: null.StringFrom(`0`),
			Extra:   null.StringFrom(``),
		},
	}
	colsHave := cols.Map(func(c csdb.Column) csdb.Column {
		c.Field.String = c.Field.String + "2"
		return c
	})

	colsWant := csdb.Columns{
		csdb.Column{Field: null.StringFrom(`category_id1312`), Type: null.StringFrom(`int(10) unsigned`), Null: null.StringFrom(`NO`), Key: null.StringFrom(`PRI`), Default: null.StringFrom(`0`), Extra: null.StringFrom(``)},
		csdb.Column{Field: null.StringFrom(`is_root_category1802`), Type: null.StringFrom(`smallint(2) unsigned`), Null: null.StringFrom(`YES`), Key: null.StringFrom(``), Default: null.StringFrom(`0`), Extra: null.StringFrom(``)},
	}

	assert.Exactly(t, colsWant, colsHave)
}

func TestColumnsFilter(t *testing.T) {

	cols := csdb.Columns{
		0: csdb.Column{
			Field:   null.StringFrom(`category_id131`),
			Type:    null.StringFrom(`int(10) unsigned`),
			Null:    null.StringFrom(`NO`),
			Key:     null.StringFrom(`PRI`),
			Default: null.StringFrom(`0`),
			Extra:   null.StringFrom(``),
		},
		1: csdb.Column{
			Field:   null.StringFrom(`is_root_category180`),
			Type:    null.StringFrom(`smallint(2) unsigned`),
			Null:    null.StringFrom(`YES`),
			Key:     null.StringFrom(``),
			Default: null.StringFrom(`0`),
			Extra:   null.StringFrom(``),
		},
	}
	colsHave := cols.Filter(func(c csdb.Column) bool {
		return c.Field.String == "is_root_category180"
	})
	colsWant := csdb.Columns{
		csdb.Column{Field: null.StringFrom(`is_root_category180`), Type: null.StringFrom(`smallint(2) unsigned`), Null: null.StringFrom(`YES`), Key: null.StringFrom(``), Default: null.StringFrom(`0`), Extra: null.StringFrom(``)},
	}

	assert.Exactly(t, colsWant, colsHave)
}

func TestGetGoPrimitive(t *testing.T) {

	tests := []struct {
		c           csdb.Column
		useNullType bool
		want        string
	}{
		{
			csdb.Column{
				Field:   null.StringFrom(`category_id131`),
				Type:    null.StringFrom(`int(10) unsigned`),
				Null:    null.StringFrom(`NO`),
				Key:     null.StringFrom(`PRI`),
				Default: null.StringFrom(`0`),
				Extra:   null.StringFrom(``),
			},
			false,
			"int64",
		},
		{
			csdb.Column{
				Field:   null.StringFrom(`category_id143`),
				Type:    null.StringFrom(`int(10) unsigned`),
				Null:    null.StringFrom(`YES`),
				Key:     null.StringFrom(`PRI`),
				Default: null.StringFrom(`0`),
				Extra:   null.StringFrom(``),
			},
			false,
			"int64",
		},
		{
			csdb.Column{
				Field:   null.StringFrom(`category_id155`),
				Type:    null.StringFrom(`int(10) unsigned`),
				Null:    null.StringFrom(`YES`),
				Key:     null.StringFrom(`PRI`),
				Default: null.StringFrom(`0`),
				Extra:   null.StringFrom(``),
			},
			true,
			"null.Int64",
		},

		{
			csdb.Column{
				Field:   null.StringFrom(`is_root_category155`),
				Type:    null.StringFrom(`smallint(2) unsigned`),
				Null:    null.StringFrom(`YES`),
				Key:     null.StringFrom(``),
				Default: null.StringFrom(`0`),
				Extra:   null.StringFrom(``),
			},
			false,
			"bool",
		},
		{
			csdb.Column{
				Field:   null.StringFrom(`is_root_category180`),
				Type:    null.StringFrom(`smallint(2) unsigned`),
				Null:    null.StringFrom(`YES`),
				Key:     null.StringFrom(``),
				Default: null.StringFrom(`0`),
				Extra:   null.StringFrom(``),
			},
			true,
			"null.Bool",
		},

		{
			csdb.Column{
				Field:   null.StringFrom(`product_name193`),
				Type:    null.StringFrom(`varchar(255)`),
				Null:    null.StringFrom(`YES`),
				Key:     null.StringFrom(``),
				Default: null.StringFrom(`0`),
				Extra:   null.StringFrom(``),
			},
			true,
			"null.String",
		},
		{
			csdb.Column{
				Field: null.StringFrom(`product_name193`),
				Type:  null.StringFrom(`varchar(255)`),
				Null:  null.StringFrom(`YES`),
			},
			false,
			"string",
		},

		{
			csdb.Column{
				Field: null.StringFrom(`price`),
				Type:  null.StringFrom(`decimal(12,4)`),
				Null:  null.StringFrom(`YES`),
			},
			false,
			"money.Money",
		},
		{
			csdb.Column{
				Field: null.StringFrom(`shipping_adjustment_230`),
				Type:  null.StringFrom(`decimal(12,4)`),
				Null:  null.StringFrom(`YES`),
			},
			true,
			"money.Money",
		},
		{
			csdb.Column{
				Field: null.StringFrom(`grand_absolut_233`),
				Type:  null.StringFrom(`decimal(12,4)`),
				Null:  null.StringFrom(`YES`),
			},
			true,
			"money.Money",
		},
		{
			csdb.Column{
				Field:   null.StringFrom(`some_currencies_242`),
				Type:    null.StringFrom(`decimal(12,4)`),
				Null:    null.StringFrom(`NO`),
				Default: null.StringFrom(`0.0000`),
			},
			true,
			"money.Money",
		},

		{
			csdb.Column{
				Field:   null.StringFrom(`weight_252`),
				Type:    null.StringFrom(`decimal(10,4)`),
				Null:    null.StringFrom(`YES`),
				Default: null.StringFrom(`0.0000`),
			},
			true,
			"null.Float64",
		},
		{
			csdb.Column{
				Field:   null.StringFrom(`weight_263`),
				Type:    null.StringFrom(`double(10,4)`),
				Null:    null.StringFrom(`YES`),
				Default: null.StringFrom(`0.0000`),
			},
			false,
			"float64",
		},

		{
			csdb.Column{
				Field:   null.StringFrom(`created_at_274`),
				Type:    null.StringFrom(`date`),
				Null:    null.StringFrom(`YES`),
				Default: null.StringFrom(`0000-00-00`),
			},
			false,
			"time.Time",
		},
		{
			csdb.Column{
				Field:   null.StringFrom(`created_at_274`),
				Type:    null.StringFrom(`date`),
				Null:    null.StringFrom(`YES`),
				Default: null.StringFrom(`0000-00-00`),
			},
			true,
			"null.Time",
		},
	}

	for i, test := range tests {
		have := test.c.GetGoPrimitive(test.useNullType)
		assert.Equal(t, test.want, have, "(%d) Test: %#v", i, test)
	}

}

var benchmarkGetColumns csdb.Columns
var benchmarkGetColumnsHashWant = []byte{0x3b, 0x2d, 0xdd, 0xf4, 0x4e, 0x2b, 0x3a, 0xd0}

// BenchmarkGetColumns-4       	5000	    395152 ns/op	   21426 B/op	     179 allocs/op
func BenchmarkGetColumns(b *testing.B) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	sess := dbc.NewSession()
	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkGetColumns, err = csdb.GetColumns(sess, "eav_attribute")
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
	csdb.Column{
		Field:   null.StringFrom("category_id"),
		Type:    null.StringFrom("int(10) unsigned"),
		Null:    null.StringFrom("NO"),
		Key:     null.StringFromPtr(nil),
		Default: null.StringFrom("0"),
		Extra:   null.StringFrom(""),
	},
	csdb.Column{
		Field:   null.StringFrom("product_id"),
		Type:    null.StringFrom("int(10) unsigned"),
		Null:    null.StringFrom("NO"),
		Key:     null.StringFrom(""),
		Default: null.StringFrom("0"),
		Extra:   null.StringFrom(""),
	},
	csdb.Column{
		Field:   null.StringFrom("position"),
		Type:    null.StringFrom("int(10) unsigned"),
		Null:    null.StringFrom("YES"),
		Key:     null.StringFrom(""),
		Default: null.String{},
		Extra:   null.StringFrom(""),
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
