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
	"github.com/gocraft/dbr"
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
		prefix   string
		expErr   bool
		expCount int
	}{
		{
			prefix:   "catalog_product_entity",
			expErr:   false,
			expCount: 11,
		},
		{
			prefix:   "' catalog_product_entity",
			expErr:   true,
			expCount: 0,
		},
	}

	for _, test := range tests {
		tables, err := GetTables(db, test.prefix)
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
		tcMap, err := GetEavValueTables(dbrConn, test.prefix, test.entityTypeCodes)
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
