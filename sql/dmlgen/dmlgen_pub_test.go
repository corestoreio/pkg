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

package dmlgen_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/corestoreio/csfw/sql/ddl"
	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/csfw/sql/dmlgen"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
SELECT
  concat('col_',
         replace(
             replace(
                 replace(
                     replace(COLUMN_TYPE, '(', '_')
                     , ')', '')
                 , ' ', '_')
             , ',', '_')
  )
    AS ColName,
  COLUMN_TYPE,
  IF(IS_NULLABLE = 'NO', 'NOT NULL', ''),
  ' DEFAULT',
  COLUMN_DEFAULT,
  ','
FROM information_schema.COLUMNS
WHERE
  table_schema = 'magento22' AND
  column_type IN (SELECT column_type
                  FROM information_schema.`COLUMNS`
                  GROUP BY column_type)
GROUP BY COLUMN_TYPE
ORDER BY COLUMN_TYPE
*/

var _ io.WriterTo = (*dmlgen.Tables)(nil)

func TestNewTables(t *testing.T) {
	t.Parallel()

	const outFile = "testdata/core_config_data_gen.go"
	os.Remove(outFile)

	ts, err := dmlgen.NewTables("testdata",
		dmlgen.WithStructTags("core_config_data",
			"path", `json:"x_path" xml:"y_path"`,
			"scope_id", `json:"scope_id" xml:"scope_id"`,
		),
		dmlgen.WithColumnAliases("core_config_data", "path", "storage_location", "config_directory"),
		dmlgen.WithUniquifiedColumns("core_config_data", "path"),
		dmlgen.WithTable("core_config_data", ddl.Columns{
			&ddl.Column{Field: "config_id", Pos: 1, Null: "NO", DataType: "int", Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: "int(10) unsigned", Key: "PRI", Extra: "auto_increment", Comment: "Config Id"},
			&ddl.Column{Field: "scope", Pos: 2, Default: dml.MakeNullString("'default'"), Null: "NO", DataType: "varchar", CharMaxLength: dml.MakeNullInt64(8), ColumnType: "varchar(8)", Key: "MUL", Comment: "Config Scope"},
			&ddl.Column{Field: "scope_id", Pos: 3, Default: dml.MakeNullString("0"), Null: "NO", DataType: "int", Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: "int(11)", Comment: "Config Scope Id"},
			&ddl.Column{Field: "path", Pos: 4, Default: dml.MakeNullString("'general'"), Null: "NO", DataType: "varchar", CharMaxLength: dml.MakeNullInt64(255), ColumnType: "varchar(255)", Comment: "Config Path"},
			&ddl.Column{Field: "value", Pos: 5, Default: dml.MakeNullString("NULL"), Null: "YES", DataType: "text", CharMaxLength: dml.MakeNullInt64(65535), ColumnType: "text", Comment: "Config Value"},
		}),
	)
	require.NoError(t, err)

	f, err := os.Create(outFile)
	if err != nil {
		t.Fatal(err)
	}
	defer cstesting.Close(t, f)

	_, err = ts.WriteTo(f)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func TestTables_WithAllTypes(t *testing.T) {
	t.Parallel()

	db, mock := cstesting.MockDB(t)
	defer cstesting.MockClose(t, db, mock)

	mock.ExpectQuery("SELECT.+").WillReturnRows(cstesting.MustMockRows(
		cstesting.WithFile("testdata/dmlgen_types.csv"),
	))

	const outFile = "testdata/dmlgen_types_gen.go"
	os.Remove(outFile)
	f, err := os.Create(outFile)
	if err != nil {
		t.Fatal(err)
	}
	defer cstesting.Close(t, f)

	ts, err := dmlgen.NewTables("testdata",
		dmlgen.WithUniquifiedColumns("dmlgen_types", "col_longtext_2", "col_int_1", "col_int_2", "has_smallint_5", "col_date_2", "col_blob"),
		dmlgen.WithLoadColumns(context.Background(), db.DB, "dmlgen_types"),
	)
	require.NoError(t, err)

	_, err = ts.WriteTo(f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithStructTags(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.IsFatal(err), "%s", err)
			} else {
				t.Errorf("Panic should contain an error but got:\n%+v", r)
			}
		} else {
			t.Error("Expecting a panic but got nothing")
		}
	}()

	dmlgen.WithStructTags("table", "unbalanced")
}
