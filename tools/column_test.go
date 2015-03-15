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

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
)

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
	}{
		{
			table:    "catalog_product_entity_decimal",
			expErr:   false,
			expCount: 5,
		},
		{
			table:    "customer_entity",
			expErr:   false,
			expCount: 11,
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
		if !test.expErr && err != nil {
			t.Error(err)
		}
		assert.Len(t, cols, test.expCount)
	}
}
