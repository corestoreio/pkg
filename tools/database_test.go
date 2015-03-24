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
	"database/sql"
	"testing"

	"strings"

	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/stretchr/testify/assert"
)

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
	db := csdb.MustConnectTest()
	defer db.Close()

	tests := []struct {
		query    string
		expErr   bool
		expCount int
	}{
		{
			query:    "SHOW TABLES LIKE 'catalog_product_entity%'",
			expErr:   false,
			expCount: 11,
		},
		{
			query:    "SHOW TABLES LIK ' catalog_product_entity",
			expErr:   true,
			expCount: 0,
		},
	}

	for _, test := range tests {
		tables, err := GetTables(db, test.query)
		if test.expErr {
			assert.Error(t, err)
		}
		if !test.expErr && err != nil {
			t.Error(err)
		}

		assert.Len(t, tables, test.expCount)
	}
}

func TestGetEavValueTables(t *testing.T) {
	db := csdb.MustConnectTest()
	defer db.Close()
	dbrConn := dbr.NewConnection(db, nil)

	tests := []struct {
		prefix          string // this is the global table name prefix
		entityTypeCodes []string
		expErr          bool
		expMap          TypeCodeValueTable
	}{
		{
			entityTypeCodes: []string{"catalog_category", "catalog_product"},
			expErr:          false,
			expMap:          TypeCodeValueTable{"catalog_category": map[string]string{"catalog_category_entity_datetime": "datetime", "catalog_category_entity_decimal": "decimal", "catalog_category_entity_int": "int", "catalog_category_entity_text": "text", "catalog_category_entity_varchar": "varchar"}, "catalog_product": map[string]string{"catalog_product_entity_datetime": "datetime", "catalog_product_entity_decimal": "decimal", "catalog_product_entity_int": "int", "catalog_product_entity_text": "text", "catalog_product_entity_varchar": "varchar"}},
		},
		{
			entityTypeCodes: []string{"customer_address", "customer"},
			expErr:          false,
			expMap:          TypeCodeValueTable{"customer_address": map[string]string{"customer_address_entity_text": "text", "customer_address_entity_varchar": "varchar", "customer_address_entity_datetime": "datetime", "customer_address_entity_decimal": "decimal", "customer_address_entity_int": "int"}, "customer": map[string]string{"csCustomer_value_decimal": "decimal", "csCustomer_value_int": "int", "csCustomer_value_text": "text", "csCustomer_value_varchar": "varchar", "csCustomer_value_datetime": "datetime"}},
		},
		{
			entityTypeCodes: []string{"catalog_address"},
			expErr:          false,
			expMap:          TypeCodeValueTable{"catalog_address": map[string]string{}},
		},
	}

	for _, test := range tests {
		tcMap, err := GetEavValueTables(dbrConn, test.entityTypeCodes)
		if test.expErr {
			assert.Error(t, err)
		}
		if !test.expErr && err != nil {
			t.Error(err)
		}

		assert.EqualValues(t, test.expMap, tcMap)
		assert.Len(t, tcMap, len(test.expMap))

	}

}

func TestColumnComment(t *testing.T) {
	c := column{
		Field: sql.NullString{
			String: "entity_id",
			Valid:  true,
		},
		Type: sql.NullString{
			String: "varchar",
			Valid:  true,
		},
		Null: sql.NullString{
			String: "YES",
			Valid:  true,
		},
		Key: sql.NullString{
			String: "PRI",
			Valid:  true,
		},
		Default: sql.NullString{
			String: "0",
			Valid:  true,
		},
		Extra: sql.NullString{
			String: "unsigned",
			Valid:  true,
		},
	}
	assert.Equal(t, "// entity_id varchar NULL PRI DEFAULT '0' unsigned", c.Comment())
}

func TestGetColumns(t *testing.T) {
	db := csdb.MustConnectTest()
	defer db.Close()

	tests := []struct {
		table    string
		expErr   bool
		expCount int
		colName  string
	}{
		{
			table:    "eav_attribute", // table is in mage 1 and 2 equal
			expErr:   false,
			expCount: 17,
			colName:  "attribute_id",
		},
		{
			table:    "catalog_product_entity_decimal",
			expErr:   false,
			expCount: 5,
			colName:  "value_id",
		},
		{
			table:    "customer_entity",
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
		cols, err := GetColumns(db, test.table)
		if test.expErr {
			assert.Error(t, err)
		}
		if err2 := cols.MapSQLToGoDBRType(); err2 != nil {
			t.Error(err2)
		}
		if !test.expErr && err != nil {
			t.Error(err)
		}
		assert.Len(t, cols, test.expCount)

		col := cols.getByName(test.colName)
		if test.colName != "" {
			assert.Equal(t, col.Field.String, test.colName)
		} else {
			assert.Nil(t, col)
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
			assert.Equal(t, len(EavAttributeColumnNameToInterface), hits)
		}

	}
}

func TestGetFieldNames(t *testing.T) {
	db := csdb.MustConnectTest()
	defer db.Close()

	tests := []struct {
		table  string
		pkOnly bool
		count  int
	}{
		{
			table:  "eav_attribute",
			pkOnly: false,
			count:  16,
		},
		{
			table:  "catalog_product_entity_decimal",
			pkOnly: true,
			count:  1,
		},
	}

	for _, test := range tests {
		cols, err := GetColumns(db, test.table)
		if err != nil {
			t.Error(err)
		}
		fields := cols.GetFieldNames(test.pkOnly)
		assert.Len(t, fields, test.count)
	}
}

// depends on generated code
func TestSQLQueryToColumns(t *testing.T) {
	db := csdb.MustConnectTest()
	defer db.Close()

	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)

	// we assume that catalog_product has entity_type_id 1 in the database ... which it has :-)
	dbrSelect, err := eav.GetAttributeSelectSql(dbrSess, NewAddAttrTables(db, "catalog_product"), 1, 0)
	if err != nil {
		t.Error(err)
	}
	dbrSelect.OrderDir("main_table.attribute_code", true)
	columns, err := SQLQueryToColumns(db, dbrSelect)
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, columns, 36)
	// @todo deeper checks
}
