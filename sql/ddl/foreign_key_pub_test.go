// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"encoding/json"
	"testing"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

func init() {
	null.MustSetJSONMarshaler(json.Marshal, json.Unmarshal)
}

func TestLoadForeignKeys_Integration(t *testing.T) {
	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)
	defer dmltest.SQLDumpLoad(t, "testdata/testLoadForeignKeys*.sql", nil).Deferred()

	t.Run("x859admin_user", func(t *testing.T) {
		tc, err := ddl.LoadKeyColumnUsage(context.TODO(), dbc.DB, "x859admin_user")
		assert.NoError(t, err)
		assert.Len(t, tc, 1, "Number of returned entries should be as stated")

		fkCols := tc["x859admin_passwords"]
		assert.NotNil(t, fkCols.Data)

		dataJSON, err := json.Marshal(fkCols.Data)
		assert.NoError(t, err)
		assert.Regexp(t,
			"[{\"ConstraintCatalog\":\"def\",\"ConstraintSchema\":\"[^\"]+\",\"ConstraintName\":\"ADMIN_PASSWORDS_USER_ID_ADMIN_USER_USER_ID\",\"TableCatalog\":\"def\",\"TableSchema\":\"[^\"]+\",\"TableName\":\"x859admin_passwords\",\"ColumnName\":\"user_id\",\"OrdinalPosition\":1,\"PositionInUniqueConstraint\":1,\"ReferencedTableSchema\":\"[^\"]+\",\"ReferencedTableName\":\"x859admin_user\",\"ReferencedColumnName\":\"user_id\"}]",
			string(dataJSON),
		)
	})

	t.Run("x910cms_block and x910cms_page", func(t *testing.T) {
		tc, err := ddl.LoadKeyColumnUsage(context.TODO(), dbc.DB, "x910cms_block", "x910cms_page")
		assert.NoError(t, err)
		assert.Len(t, tc, 2, "Number of returned entries should be as stated")

		dataJSON, err := json.Marshal(tc["x910cms_block_store"].Data)
		assert.NoError(t, err)
		assert.Regexp(t,
			"[{\"ConstraintCatalog\":\"def\",\"ConstraintSchema\":\"[^\"]+\",\"ConstraintName\":\"CMS_BLOCK_STORE_BLOCK_ID_CMS_BLOCK_BLOCK_ID\",\"TableCatalog\":\"def\",\"TableSchema\":\"[^\"]+\",\"TableName\":\"x910cms_block_store\",\"ColumnName\":\"block_id\",\"OrdinalPosition\":1,\"PositionInUniqueConstraint\":1,\"ReferencedTableSchema\":\"[^\"]+\",\"ReferencedTableName\":\"x910cms_block\",\"ReferencedColumnName\":\"block_id\"}]",
			string(dataJSON),
		)

		dataJSON, err = json.Marshal(tc["x910cms_page_store"].Data)
		assert.NoError(t, err)
		assert.Regexp(t,
			"[{\"ConstraintCatalog\":\"def\",\"ConstraintSchema\":\"[^\"]+\",\"ConstraintName\":\"CMS_PAGE_STORE_PAGE_ID_CMS_PAGE_PAGE_ID\",\"TableCatalog\":\"def\",\"TableSchema\":\"[^\"]+\",\"TableName\":\"x910cms_page_store\",\"ColumnName\":\"page_id\",\"OrdinalPosition\":1,\"PositionInUniqueConstraint\":1,\"ReferencedTableSchema\":\"[^\"]+\",\"ReferencedTableName\":\"x910cms_page\",\"ReferencedColumnName\":\"page_id\"}]",
			string(dataJSON),
		)
	})

	t.Run("catalog_eav_attribute", func(t *testing.T) {
		tc, err := ddl.LoadKeyColumnUsage(context.TODO(), dbc.DB, "x910catalog_eav_attribute")
		assert.NoError(t, err)
		assert.Len(t, tc, 0, "Number of returned entries should be as stated")

		fkCols, ok := tc["x910catalog_eav_attribute"]
		assert.False(t, ok)
		assert.Nil(t, fkCols.Data)
	})
}

func TestLoadKeyRelationships(t *testing.T) {
	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)
	defer dmltest.SQLDumpLoad(t, "testdata/testLoadForeignKeys*.sql", nil).Deferred()
	// dmltest.SQLDumpLoad(t, "testdata/testLoadForeignKeys*.sql", nil)

	ctx := context.Background()

	tc, err := ddl.LoadKeyColumnUsage(context.TODO(), dbc.DB)
	assert.NoError(t, err)

	krs, err := ddl.GenerateKeyRelationships(ctx, dbc.DB, tc)
	assert.NoError(t, err)

	var buf bytes.Buffer
	krs.Debug(&buf)
	t.Log("\n", buf.String())
	assert.LenBetween(t, buf.String(), 100, 2300)

	t.Run("ManyToMany", func(t *testing.T) {
		targetTable, targetColumn := krs.ManyToManyTarget("athlete_team_member", "team_id", "athlete_team", "team_id")
		assert.Exactly(t, "athlete", targetTable, "targetTable.targetColumn: %q.%q", targetTable, targetColumn)
		assert.Exactly(t, "athlete_id", targetColumn)
		targetTable, targetColumn = krs.ManyToManyTarget("athlete_team_member", "athlete_id", "athlete", "athlete_id")
		assert.Exactly(t, "athlete_team", targetTable)
		assert.Exactly(t, "team_id", targetColumn)
		targetTable, targetColumn = krs.ManyToManyTarget("athlete_team_member", "athlete_idx", "athlete", "athlete_idx")
		assert.Exactly(t, "", targetTable)
		assert.Exactly(t, "", targetColumn)
	})

	t.Run("OneToX", func(t *testing.T) {
		tests := []struct {
			checkFn                                          func(referencedTable, referencedColumn, table, column string) bool
			referencedTable, referencedColumn, table, column string
			want                                             bool
		}{
			{krs.IsOneToOne, "x859admin_passwords", "user_id", "x859admin_user", "user_id", true},
			{krs.IsOneToOne, "x859admin_user", "user_id", "x859admin_passwords", "user_id", false},
			{krs.IsOneToMany, "x859admin_user", "user_id", "x859admin_passwords", "user_id", true},
			{krs.IsOneToMany, "x859admin_passwords", "user_id", "x859admin_user", "user_id", false},

			{krs.IsOneToOne, "x910cms_page_store", "page_id", "x910cms_page", "page_id", true},
			{krs.IsOneToMany, "x910cms_page", "page_id", "x910cms_page_store", "page_id", true},

			{krs.IsOneToMany, "store_group", "website_id", "store", "website_id", false}, // no FK constraint
			{krs.IsOneToOne, "store_group", "website_id", "store", "website_id", false},  // no FK constraint
			{krs.IsOneToMany, "store", "website_id", "store_group", "website_id", false}, // no FK constraint
			{krs.IsOneToOne, "store", "website_id", "store_group", "website_id", false},

			{krs.IsOneToOne, "store_group", "website_id", "store_website", "website_id", true},
			{krs.IsOneToMany, "store_website", "website_id", "store_group", "website_id", true}, // reversed above

			{krs.IsOneToOne, "catalog_category_entity", "entity_id", "sequence_catalog_category", "sequence_value", true},
			// reversed must also be true for oneToOne because sequence_catalog_category contains only one column
			{krs.IsOneToOne, "sequence_catalog_category", "sequence_value", "catalog_category_entity", "entity_id", true},
			{krs.IsOneToMany, "sequence_catalog_category", "sequence_value", "catalog_category_entity", "entity_id", false},
			{krs.IsOneToMany, "catalog_category_entity", "entity_id", "sequence_catalog_category", "sequence_value", false},
		}

		for i, test := range tests {
			assert.Exactly(t,
				test.want,
				test.checkFn(test.referencedTable, test.referencedColumn, test.table, test.column),
				"IDX %d %q.%q => %q.%q",
				i, test.referencedTable, test.referencedColumn, test.table, test.column,
			)
		}
	})
}
