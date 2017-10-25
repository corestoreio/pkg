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
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/corestoreio/csfw/sql/ddl"
	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/csfw/sql/dmlgen"
	"github.com/corestoreio/csfw/util/cstesting"
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

var _ io.WriterTo = (*dmlgen.Table)(nil)

func writeGoFileHeader(w io.Writer, imports []string) {
	w.Write([]byte("package testdata\n\nimport (\n"))
	for _, i := range imports {
		fmt.Fprintf(w, "\t%q\n", i)
	}
	w.Write([]byte("\n)\n"))
}

func TestTable_WriteTo(t *testing.T) {
	t.Parallel()

	const outFile = "testdata/core_config_data.gen.go"

	os.Remove(outFile)

	tbl := &dmlgen.Table{
		Package: "testdata",
		Name:    "core_config_data",
		Columns: ddl.Columns{
			&ddl.Column{Field: "config_id", Pos: 1, Null: "NO", DataType: "int", Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: "int(10) unsigned", Key: "PRI", Extra: "auto_increment", Comment: "Config Id"},
			&ddl.Column{Field: "scope", Pos: 2, Default: dml.MakeNullString("'default'"), Null: "NO", DataType: "varchar", CharMaxLength: dml.MakeNullInt64(8), ColumnType: "varchar(8)", Key: "MUL", Comment: "Config Scope"},
			&ddl.Column{Field: "scope_id", Pos: 3, Default: dml.MakeNullString("0"), Null: "NO", DataType: "int", Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: "int(11)", Comment: "Config Scope Id"},
			&ddl.Column{Field: "path", Pos: 4, Default: dml.MakeNullString("'general'"), Null: "NO", DataType: "varchar", CharMaxLength: dml.MakeNullInt64(255), ColumnType: "varchar(255)", Comment: "Config Path"},
			&ddl.Column{Field: "value", Pos: 5, Default: dml.MakeNullString("NULL"), Null: "YES", DataType: "text", CharMaxLength: dml.MakeNullInt64(65535), ColumnType: "text", Comment: "Config Value"},
		},
		ColumnAliases: map[string][]string{
			"path": {"storage_location", "config_directory"}, // just some values
		},
		AllowedDuplicateValueColumns: []string{"path"},
	}
	f, err := os.Create(outFile)
	if err != nil {
		t.Fatal(err)
	}
	defer cstesting.Close(t, f)

	writeGoFileHeader(f, dmlgen.Imports["table"])

	_, err = tbl.WriteTo(f)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func TestTable_WithAllTypes(t *testing.T) {
	t.Parallel()

	db, mock := cstesting.MockDB(t)
	defer cstesting.MockClose(t, db, mock)

	mock.ExpectQuery("SELECT.+").WillReturnRows(cstesting.MustMockRows(
		cstesting.WithFile("testdata/dmlgen_types.csv"),
	))

	colMap, err := ddl.LoadColumns(context.Background(), db.DB, "dmlgen_types")
	require.NoError(t, err)

	const outFile = "testdata/dmlgen_types.gen.go"
	os.Remove(outFile)
	f, err := os.Create(outFile)
	if err != nil {
		t.Fatal(err)
	}
	defer cstesting.Close(t, f)

	writeGoFileHeader(f, dmlgen.Imports["table"])
	tbl := &dmlgen.Table{
		Package: "testdata",
		Name:    "dmlgen_types",
		Columns: colMap["dmlgen_types"],
		AllowedDuplicateValueColumns: []string{"col_longtext_2", "col_int_1", "col_int_2", "has_smallint_5", "col_date_2", "col_blob"},
	}
	_, err = tbl.WriteTo(f)
	if err != nil {
		t.Fatal(err)
	}
}
