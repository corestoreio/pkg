// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"testing"

	"errors"
	"fmt"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
)

func TestGetColumns(t *testing.T) {
	// @todo fix test data retrieving from database ...
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	sess := dbc.NewSession()

	tests := []struct {
		table          string
		want           string
		wantErr        error
		wantJoinFields string
	}{
		{"core_config_data",
			"csdb.Column{Field:csdb.NewNullString(`config_id`), Type:csdb.NewNullString(`int(10) unsigned`), Null:csdb.NewNullString(`NO`), Key:csdb.NewNullString(`PRI`), Default:csdb.NewNullString(nil), Extra:csdb.NewNullString(`auto_increment`)},\ncsdb.Column{Field:csdb.NewNullString(`scope`), Type:csdb.NewNullString(`varchar(8)`), Null:csdb.NewNullString(`NO`), Key:csdb.NewNullString(`MUL`), Default:csdb.NewNullString(`default`), Extra:csdb.NewNullString(``)},\ncsdb.Column{Field:csdb.NewNullString(`scope_id`), Type:csdb.NewNullString(`int(11)`), Null:csdb.NewNullString(`NO`), Key:csdb.NewNullString(``), Default:csdb.NewNullString(`0`), Extra:csdb.NewNullString(``)},\ncsdb.Column{Field:csdb.NewNullString(`path`), Type:csdb.NewNullString(`varchar(255)`), Null:csdb.NewNullString(`NO`), Key:csdb.NewNullString(``), Default:csdb.NewNullString(`general`), Extra:csdb.NewNullString(``)},\ncsdb.Column{Field:csdb.NewNullString(`value`), Type:csdb.NewNullString(`text`), Null:csdb.NewNullString(`YES`), Key:csdb.NewNullString(``), Default:csdb.NewNullString(nil), Extra:csdb.NewNullString(``)}\n",
			nil,
			"config_id_scope_scope_id_path_value",
		},
		{"catalog_category_product",
			"csdb.Column{Field:csdb.NewNullString(`category_id`), Type:csdb.NewNullString(`int(10) unsigned`), Null:csdb.NewNullString(`NO`), Key:csdb.NewNullString(`PRI`), Default:csdb.NewNullString(`0`), Extra:csdb.NewNullString(``)},\ncsdb.Column{Field:csdb.NewNullString(`product_id`), Type:csdb.NewNullString(`int(10) unsigned`), Null:csdb.NewNullString(`NO`), Key:csdb.NewNullString(`PRI`), Default:csdb.NewNullString(`0`), Extra:csdb.NewNullString(``)},\ncsdb.Column{Field:csdb.NewNullString(`position`), Type:csdb.NewNullString(`int(11)`), Null:csdb.NewNullString(`NO`), Key:csdb.NewNullString(``), Default:csdb.NewNullString(`0`), Extra:csdb.NewNullString(``)}\n",
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
			//t.Logf("%s\n%#v\n", err1.Error(), err1.(errgo.Locationer).Location())
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
			"csdb.Column{Field:csdb.NewNullString(`category_id`), Type:csdb.NewNullString(`int(10) unsigned`), Null:csdb.NewNullString(`NO`), Key:csdb.NewNullString(`MUL`), Default:csdb.NewNullString(`0`), Extra:csdb.NewNullString(``)},\ncsdb.Column{Field:csdb.NewNullString(`path`), Type:csdb.NewNullString(`varchar(255)`), Null:csdb.NewNullString(`YES`), Key:csdb.NewNullString(`MUL`), Default:csdb.NewNullString(nil), Extra:csdb.NewNullString(``)}",
		},
		{
			mustStructure(table2).Columns.PrimaryKeys().Len(),
			1,
			mustStructure(table2).Columns.GoString(),
			"csdb.Column{Field:csdb.NewNullString(`category_id`), Type:csdb.NewNullString(`int(10) unsigned`), Null:csdb.NewNullString(`NO`), Key:csdb.NewNullString(`PRI`), Default:csdb.NewNullString(`0`), Extra:csdb.NewNullString(``)},\ncsdb.Column{Field:csdb.NewNullString(`path`), Type:csdb.NewNullString(`varchar(255)`), Null:csdb.NewNullString(`YES`), Key:csdb.NewNullString(nil), Default:csdb.NewNullString(nil), Extra:csdb.NewNullString(``)}",
		},
		{
			mustStructure(table4).Columns.UniqueKeys().Len(), 1,
			mustStructure(table4).Columns.GoString(),
			"csdb.Column{Field:csdb.NewNullString(`user_id`), Type:csdb.NewNullString(`int(10) unsigned`), Null:csdb.NewNullString(`NO`), Key:csdb.NewNullString(`PRI`), Default:csdb.NewNullString(nil), Extra:csdb.NewNullString(`auto_increment`)},\ncsdb.Column{Field:csdb.NewNullString(`email`), Type:csdb.NewNullString(`varchar(128)`), Null:csdb.NewNullString(`YES`), Key:csdb.NewNullString(``), Default:csdb.NewNullString(nil), Extra:csdb.NewNullString(``)},\ncsdb.Column{Field:csdb.NewNullString(`username`), Type:csdb.NewNullString(`varchar(40)`), Null:csdb.NewNullString(`YES`), Key:csdb.NewNullString(`UNI`), Default:csdb.NewNullString(nil), Extra:csdb.NewNullString(``)}",
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

func TestGetGoPrimitive(t *testing.T) {

	tests := []struct {
		c           csdb.Column
		useNullType bool
		want        string
	}{
		{
			csdb.Column{
				Field:   csdb.NewNullString(`category_id131`),
				Type:    csdb.NewNullString(`int(10) unsigned`),
				Null:    csdb.NewNullString(`NO`),
				Key:     csdb.NewNullString(`PRI`),
				Default: csdb.NewNullString(`0`),
				Extra:   csdb.NewNullString(``),
			},
			false,
			"int64",
		},
		{
			csdb.Column{
				Field:   csdb.NewNullString(`category_id143`),
				Type:    csdb.NewNullString(`int(10) unsigned`),
				Null:    csdb.NewNullString(`YES`),
				Key:     csdb.NewNullString(`PRI`),
				Default: csdb.NewNullString(`0`),
				Extra:   csdb.NewNullString(``),
			},
			false,
			"int64",
		},
		{
			csdb.Column{
				Field:   csdb.NewNullString(`category_id155`),
				Type:    csdb.NewNullString(`int(10) unsigned`),
				Null:    csdb.NewNullString(`YES`),
				Key:     csdb.NewNullString(`PRI`),
				Default: csdb.NewNullString(`0`),
				Extra:   csdb.NewNullString(``),
			},
			true,
			"csdb.NullInt64",
		},

		{
			csdb.Column{
				Field:   csdb.NewNullString(`is_root_category155`),
				Type:    csdb.NewNullString(`smallint(2) unsigned`),
				Null:    csdb.NewNullString(`YES`),
				Key:     csdb.NewNullString(``),
				Default: csdb.NewNullString(`0`),
				Extra:   csdb.NewNullString(``),
			},
			false,
			"bool",
		},
		{
			csdb.Column{
				Field:   csdb.NewNullString(`is_root_category180`),
				Type:    csdb.NewNullString(`smallint(2) unsigned`),
				Null:    csdb.NewNullString(`YES`),
				Key:     csdb.NewNullString(``),
				Default: csdb.NewNullString(`0`),
				Extra:   csdb.NewNullString(``),
			},
			true,
			"csdb.NullBool",
		},

		{
			csdb.Column{
				Field:   csdb.NewNullString(`product_name193`),
				Type:    csdb.NewNullString(`varchar(255)`),
				Null:    csdb.NewNullString(`YES`),
				Key:     csdb.NewNullString(``),
				Default: csdb.NewNullString(`0`),
				Extra:   csdb.NewNullString(``),
			},
			true,
			"csdb.NullString",
		},
		{
			csdb.Column{
				Field: csdb.NewNullString(`product_name193`),
				Type:  csdb.NewNullString(`varchar(255)`),
				Null:  csdb.NewNullString(`YES`),
			},
			false,
			"string",
		},

		{
			csdb.Column{
				Field: csdb.NewNullString(`price`),
				Type:  csdb.NewNullString(`decimal(12,4)`),
				Null:  csdb.NewNullString(`YES`),
			},
			false,
			"money.Money",
		},
		{
			csdb.Column{
				Field: csdb.NewNullString(`shipping_adjustment_230`),
				Type:  csdb.NewNullString(`decimal(12,4)`),
				Null:  csdb.NewNullString(`YES`),
			},
			true,
			"money.Money",
		},
		{
			csdb.Column{
				Field: csdb.NewNullString(`grand_absolut_233`),
				Type:  csdb.NewNullString(`decimal(12,4)`),
				Null:  csdb.NewNullString(`YES`),
			},
			true,
			"money.Money",
		},
		{
			csdb.Column{
				Field:   csdb.NewNullString(`some_currencies_242`),
				Type:    csdb.NewNullString(`decimal(12,4)`),
				Null:    csdb.NewNullString(`NO`),
				Default: csdb.NewNullString(`0.0000`),
			},
			true,
			"money.Money",
		},

		{
			csdb.Column{
				Field:   csdb.NewNullString(`weight_252`),
				Type:    csdb.NewNullString(`decimal(10,4)`),
				Null:    csdb.NewNullString(`YES`),
				Default: csdb.NewNullString(`0.0000`),
			},
			true,
			"csdb.NullFloat64",
		},
		{
			csdb.Column{
				Field:   csdb.NewNullString(`weight_263`),
				Type:    csdb.NewNullString(`double(10,4)`),
				Null:    csdb.NewNullString(`YES`),
				Default: csdb.NewNullString(`0.0000`),
			},
			false,
			"float64",
		},

		{
			csdb.Column{
				Field:   csdb.NewNullString(`created_at_274`),
				Type:    csdb.NewNullString(`date`),
				Null:    csdb.NewNullString(`YES`),
				Default: csdb.NewNullString(`0000-00-00`),
			},
			false,
			"time.Time",
		},
		{
			csdb.Column{
				Field:   csdb.NewNullString(`created_at_274`),
				Type:    csdb.NewNullString(`date`),
				Null:    csdb.NewNullString(`YES`),
				Default: csdb.NewNullString(`0000-00-00`),
			},
			true,
			"csdb.NullTime",
		},
	}

	for _, test := range tests {
		have := test.c.GetGoPrimitive(test.useNullType)
		assert.Equal(t, test.want, have, "Test: %#v", test)
	}

}

var benchmarkGetColumns csdb.Columns
var benchmarkGetColumnsHashWant = []byte{0x3b, 0x2d, 0xdd, 0xf4, 0x4e, 0x2b, 0x3a, 0xd0}

// BenchmarkGetColumns-4	    1000	   3376128 ns/op	   24198 B/op	     196 allocs/op
// BenchmarkGetColumns-4	    1000	   1185381 ns/op	   21861 B/op	     179 allocs/op
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
		Field:   csdb.NewNullString("category_id"),
		Type:    csdb.NewNullString("int(10) unsigned"),
		Null:    csdb.NewNullString("NO"),
		Key:     csdb.NewNullString(nil),
		Default: csdb.NewNullString("0"),
		Extra:   csdb.NewNullString(""),
	},
	csdb.Column{
		Field:   csdb.NewNullString("product_id"),
		Type:    csdb.NewNullString("int(10) unsigned"),
		Null:    csdb.NewNullString("NO"),
		Key:     csdb.NewNullString(""),
		Default: csdb.NewNullString("0"),
		Extra:   csdb.NewNullString(""),
	},
	csdb.Column{
		Field:   csdb.NewNullString("position"),
		Type:    csdb.NewNullString("int(10) unsigned"),
		Null:    csdb.NewNullString("YES"),
		Key:     csdb.NewNullString(""),
		Default: csdb.NullString{},
		Extra:   csdb.NewNullString(""),
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
