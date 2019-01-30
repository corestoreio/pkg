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
	"context"
	"encoding/json"
	"testing"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

func init() {
	null.JSONMarshalFn = json.Marshal
}

func TestLoadForeignKeys_Integration(t *testing.T) {
	t.Parallel()

	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)
	defer dmltest.SQLDumpLoad(t, "testdata/testLoadForeignKeys*.sql", nil).Deferred()

	t.Run("x859admin_user", func(t *testing.T) {

		tc, err := ddl.LoadKeyColumnUsage(context.TODO(), dbc.DB, "x859admin_user")
		assert.NoError(t, err)
		assert.Len(t, tc, 1, "Number of returned entries should be as stated")

		fkCols := tc["x859admin_user.user_id"]
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

		dataJSON, err := json.Marshal(tc["x910cms_block.block_id"].Data)
		assert.NoError(t, err)
		assert.Regexp(t,
			"[{\"ConstraintCatalog\":\"def\",\"ConstraintSchema\":\"[^\"]+\",\"ConstraintName\":\"CMS_BLOCK_STORE_BLOCK_ID_CMS_BLOCK_BLOCK_ID\",\"TableCatalog\":\"def\",\"TableSchema\":\"[^\"]+\",\"TableName\":\"x910cms_block_store\",\"ColumnName\":\"block_id\",\"OrdinalPosition\":1,\"PositionInUniqueConstraint\":1,\"ReferencedTableSchema\":\"[^\"]+\",\"ReferencedTableName\":\"x910cms_block\",\"ReferencedColumnName\":\"block_id\"}]",
			string(dataJSON),
		)

		dataJSON, err = json.Marshal(tc["x910cms_page.page_id"].Data)
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

		fkCols, ok := tc["x910catalog_eav_attribute.attribute_id"]
		assert.False(t, ok)
		assert.Nil(t, fkCols.Data)
	})
}
