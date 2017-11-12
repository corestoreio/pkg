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

package codegen

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/corestoreio/pkg/eav"
	"github.com/corestoreio/pkg/storage/csdb"
	"github.com/corestoreio/pkg/storage/dbr"
	"github.com/stretchr/testify/assert"
)

func init() {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	eav.TableCollection.Init(dbc.NewSession())
}

func TestValueSuffixes(t *testing.T) {
	vs := ValueSuffixes{"One", "two", "Baum"}
	assert.True(t, vs.contains("One"))
	assert.False(t, vs.contains("Strauch"))
	assert.Equal(t, vs.String(), "One, two, Baum")
}

func TestTypeCodeValueTable(t *testing.T) {
	tcMap := make(TypeCodeValueTable, 1)
	assert.True(t, tcMap.Empty())
	tcMap["catalog"] = make(map[string]string)
	assert.False(t, tcMap.Empty())
}

func TestGetTables(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()

	tests := []struct {
		query    string
		expErr   bool
		expCount int
	}{
		{
			query:    "catalog_product_entity%",
			expErr:   false,
			expCount: 11,
		},
		{
			query: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
						    TABLE_SCHEMA = DATABASE() AND
						    TABLE_NAME LIKE '%directory%' GROUP BY TABLE_NAME;`,
			expErr:   false,
			expCount: 5,
		},
		{
			query:    "' catalog_product_entity",
			expErr:   true,
			expCount: 0,
		},
	}

	for _, test := range tests {
		tables, err := GetTables(dbc.NewSession(), test.query)
		if test.expErr {
			assert.Error(t, err)
		}
		if !test.expErr && err != nil {
			t.Error(err)
		}
		assert.True(t, len(tables) >= test.expCount, "have %d min want %d", len(tables), test.expCount)
	}
}

type dataGetEavValueTables struct {
	haveETC   []string // have entity type codes
	wantErr   bool
	wantCVMap TypeCodeValueTable
}

func TestGetEavValueTables(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()

	// @todo for mage2 we need more tests for custom named entity value tables (which is pretty rare)

	// getDataGetEavValueTables depends on the build tag mage1 or mage2
	var tests []dataGetEavValueTables = getDataGetEavValueTables() // type hint needed for Intellij

	for i, test := range tests {
		tcMap, err := GetEavValueTables(dbc, test.haveETC)
		if test.wantErr {
			assert.Error(t, err, "Index %d", i)
		}
		if !test.wantErr && err != nil {
			t.Error(err)
		}

		assert.EqualValues(t, test.wantCVMap, tcMap, "Index %d", i)
		assert.Len(t, tcMap, len(test.wantCVMap), "Index %d", i)
	}
}

func TestColumnComment(t *testing.T) {
	c := column{
		Column: csdb.Column{
			Field:      dbr.NewNullString("entity_id"),
			ColumnType: dbr.NewNullString("varchar"),
			Null:       dbr.NewNullString("YES"),
			Key:        dbr.NewNullString("PRI"),
			Default:    dbr.NewNullString("0"),
			Extra:      dbr.NewNullString("unsigned"),
		},
	}
	assert.Equal(t, "// entity_id varchar NULL PRI DEFAULT '0' unsigned", c.Comment())
}

func TestGetColumns(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()

	tests := []struct {
		table    string
		expErr   bool
		expCount int
		colName  string
	}{
		{
			table:    "eav_attribute", // table is in mage 1 and 2 equal
			expErr:   false,
			expCount: 16,
			colName:  "attribute_id",
		},
		{
			table:    "catalog_product_entity_decimal",
			expErr:   false,
			expCount: 5,
			colName:  "value_id",
		},
		{
			table:    "customer_entity", // Magento2 has more columns because to get rid of EAV tables regarding performance.
			expErr:   false,
			expCount: 11,
			colName:  "entity_id",
		},
		{
			table:    "catalog_category_entity_datetime",
			expErr:   false,
			expCount: 5,
			colName:  "entity_id",
		},
		{
			table:    "customer_address_entity_decimal",
			expErr:   false,
			expCount: 4,
			colName:  "entity_id",
		},
		{
			table:    "', customer_entity",
			expErr:   true,
			expCount: 0,
		},
	}

	for _, test := range tests {
		cols, err := GetColumns(dbc.DB, test.table)
		if test.expErr {
			assert.Error(t, err)
		}
		if err2 := cols.MapSQLToGoDBRType(); err2 != nil {
			t.Error(err2)
		}
		if !test.expErr && err != nil {
			t.Error(err)
		}
		assert.True(t, len(cols) >= test.expCount, "For table %s", test.table)

		col := cols.GetByName(test.colName)
		if test.colName != "" {
			assert.Equal(t, col.Field.String, test.colName)
		} else {
			assert.NotNil(t, col)
		}

		if test.table == "eav_attribute" {
			if err := cols.MapSQLToGoType(EavAttributeColumnNameToInterface); err != nil {
				t.Error(err)
			}
			hits := 0
			for _, col := range cols {
				if interf, ok := EavAttributeColumnNameToInterface[col.Field.String]; ok {
					assert.Equal(t, interf, col.GoType)
					hits++
				}
				if strings.Contains(col.GoType, "dbr") {
					t.Errorf("%s contains dbr but it should not. %#v", col.GoType, col)
				}
			}
			assert.Equal(t, 3, hits)
		}

	}
}

func TestGetFieldNames(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()

	tests := []struct {
		table  string
		pkOnly bool
		count  int
	}{
		{
			table:  "eav_attribute",
			pkOnly: false,
			count:  15,
		},
		{
			table:  "catalog_product_entity_decimal",
			pkOnly: true,
			count:  1,
		},
	}

	for _, test := range tests {
		cols, err := GetColumns(dbc.DB, test.table)
		if err != nil {
			t.Error(err)
		}
		fields := cols.GetFieldNames(test.pkOnly)
		assert.Len(t, fields, test.count)
	}
}

// depends on generated code
func TestSQLQueryToColumnsToStruct(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()

	dbrSess := dbc.NewSession()
	dbrSelect, err := eav.GetAttributeSelectSql(dbrSess, NewAddAttrTables(dbc.DB, "catalog_product"), 4, 0)
	if err != nil {
		t.Error(err)
	}

	colSliceDbr, err := SQLQueryToColumns(dbc.DB, dbrSelect)
	if err != nil {
		t.Error(err)
	}

	assert.True(t, len(colSliceDbr) >= 18, "len(colSliceDbr) == %d, should have min 18", len(colSliceDbr))

	for _, col := range colSliceDbr {
		assert.True(t, col.Field.Valid, fmt.Sprintf("%#v", col))
		assert.True(t, col.ColumnType.Valid, fmt.Sprintf("%#v", col))
	}

	columns2, err2 := SQLQueryToColumns(dbc.DB, nil, "SELECT * FROM `catalog_product_option`", " ORDER BY option_id DESC")
	if err2 != nil {
		t.Error(err2)
	}
	assert.Len(t, columns2, 10)
	for _, col := range columns2 {
		assert.True(t, col.Field.Valid, fmt.Sprintf("%#v", col))
		assert.True(t, col.ColumnType.Valid, fmt.Sprintf("%#v", col))
	}

	colSliceDbr.MapSQLToGoDBRType()
	code, err := ColumnsToStructCode(nil, "testStruct", colSliceDbr)
	if err != nil {
		t.Error(err, "\n", string(code))
	}

	checkContains := [][]byte{
		[]byte(`TeststructSlice`),
		[]byte(`dbr.NullString`),
		[]byte("`db:\"is_visible_in_advanced_search\"`"),
	}
	for _, s := range checkContains {
		if false == bytes.Contains(code, s) {
			t.Errorf("%s\ndoes not contain %s", code, s)
		}
	}
}

func TestGetSQLPrepareForTemplate(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()

	sel := dbc.NewSession().Select("*").From("cataloginventory_stock").OrderBy("stock_id")
	resultSlice2, err := LoadStringEntities(dbc.DB, sel)
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, resultSlice2, 1) // 1 row
	for _, row := range resultSlice2 {
		assert.True(t, len(row["stock_id"]) > 0, "Incorrect length of stock_id", fmt.Sprintf("%#v", row))
	}

	// advanced test

	dbrSess := dbc.NewSession()
	dbrSelect, err := eav.GetAttributeSelectSql(dbrSess, NewAddAttrTables(dbc.DB, "catalog_product"), 4, 0)
	if err != nil {
		t.Error(err)
	}

	attributeResultSlice, err := LoadStringEntities(dbc.DB, dbrSelect)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, len(attributeResultSlice) >= 59, "attributeResultSlice should have at least 59 entries, but has %d", len(attributeResultSlice))

	for _, row := range attributeResultSlice {
		assert.True(t, len(row["attribute_id"]) > 0, "Incorrect length of attribute_id: %#v", row)
	}

	colSliceDbr, err := SQLQueryToColumns(dbc.DB, dbrSelect)
	if err != nil {
		t.Fatal(err)
	}

	for _, col := range colSliceDbr {
		assert.Empty(t, col.GoType)
		assert.Empty(t, col.GoName)
	}

	var unchanged = make(map[string]string)
	for _, s := range attributeResultSlice {
		assert.True(t, len(s["is_wysiwyg_enabled"]) == 1, "Should contain 0 or 1 as string: %s", s["is_wysiwyg_enabled"])
		assert.True(t, len(s["used_in_product_listing"]) == 1, "Should contain 0 or 1 as string: %s", s["used_in_product_listing"])
		assert.False(t, strings.ContainsRune(s["attribute_code"], '"'), "Should not contain double quotes for escaping: %s", s["attribute_code"])
		unchanged[s["attribute_id"]] = s["entity_type_id"]
	}

	importPaths1 := PrepareForTemplate(colSliceDbr, attributeResultSlice, ConfigAttributeModel, "catalog")
	assert.True(t, len(importPaths1) > 1, "Should output multiple import paths: %#v", importPaths1)

	for _, s := range attributeResultSlice {
		assert.True(t, len(s["is_wysiwyg_enabled"]) >= 4, "Should contain false or true as string: %s", s["is_wysiwyg_enabled"])
		assert.True(t, len(s["used_in_product_listing"]) >= 4, "Should contain false or true as string: %s", s["used_in_product_listing"])
		assert.True(t, strings.ContainsRune(s["attribute_code"], '"'), "Should contain double quotes for escaping: %s", s["attribute_code"])
		assert.Equal(t, unchanged[s["attribute_id"]], s["entity_type_id"], "Columns: %#v", s)
		assert.True(t, len(s["frontend_model"]) >= 3, "Should contain nil or a Go func: %s", s["frontend_model"])
		assert.True(t, len(s["backend_model"]) >= 3, "Should contain nil or a Go func: %s", s["backend_model"])
		assert.True(t, len(s["source_model"]) >= 3, "Should contain nil or a Go func: %s", s["source_model"])
	}
}

var benchmarkGetTables []string

// BenchmarkGetTables-4	    2000	    865974 ns/op	   31042 B/op	     683 allocs/op
func BenchmarkGetTables(b *testing.B) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	b.ReportAllocs()
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkGetTables, err = GetTables(dbc.NewSession())
		if err != nil {
			b.Error(err)
		}
		if len(benchmarkGetTables) < 200 {
			b.Errorf("There should be at least 200 tables in the database. Got: %d", len(benchmarkGetTables))
		}
	}
}
