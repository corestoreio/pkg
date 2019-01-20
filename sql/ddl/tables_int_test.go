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

package ddl

import (
	"context"
	"sort"
	"testing"

	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
)

func TestIsCreateStmt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		tableName string
		stmt      string
		ok        bool
	}{
		{"tableA", "CREATE TABLE `tableA` ( `config_id` int(10) );", true},
		{"tableA2", "ALTER TABLE `tableA2` ( `config_id` int(10) );", false},
		{"tableB", "CREATE TABLE tableB ( `config_id` int(10) );", true},
		{"tableV1", "CREATE VIEW tableV1 ( `config_id` int(10) );", true},
		{"tableV2", "CREATE VIEW `tableV2` ( `config_id` int(10) );", true},
		{"tableV3", "create view `tableV2` ( `config_id` int(10) );", false},
		{"tableC", "CREATE TABLE\nIF NOT\tEXISTS\ntableC ( `config_id` int(10) );", true},
		{"tableC", "CREATE\tVIEW\nIF NOT\tEXISTS\ntableC ( `config_id` int(10) );", true},
	}

	for i, test := range tests {
		haveOK := isCreateStmt(test.tableName, test.stmt) && isCreateStmtBytes([]byte(test.tableName), []byte(test.stmt))
		assert.Exactly(t, test.ok, haveOK, "Index %d %q", i, test.tableName)
	}
}

func TestWithTableDB(t *testing.T) {
	t.Parallel()
	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE").
		WillReturnRows(
			dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_columns.csv")))

	ts := MustNewTables(
		WithDB(dbc.DB),
		WithTable("tableA"),
		WithCreateTable(context.TODO(), "tableB", ""),
	) // +=2

	assert.Exactly(t, dbc.DB, ts.MustTable("tableA").dcp.DB)
	assert.Exactly(t, dbc.DB, ts.MustTable("tableB").dcp.DB)
	haveTS := ts.Tables()
	sort.Strings(haveTS)
	assert.Exactly(t, []string{"tableA", "tableB"}, haveTS)
}
