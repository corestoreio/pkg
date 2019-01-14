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
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dmlgen"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

var _ = null.JSONMarshalFn

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

func writeFile(t *testing.T, outFile string, wFn func(io.Writer, io.Writer) error) {
	testF := ioutil.Discard
	if strings.HasSuffix(outFile, ".go") {
		testFile := strings.Replace(outFile, ".go", "_test.go", 1)

		ft, err := os.Create(testFile)
		assert.NoError(t, err)
		defer dmltest.Close(t, ft)
		testF = ft
	}

	f, err := os.Create(outFile)
	assert.NoError(t, err)
	defer dmltest.Close(t, f)
	err = wFn(f, testF)
	assert.NoError(t, err, "%+v", err)
}

// TestGenerate_Tables_Protobuf_Json writes a Go and Proto file to the testdata
// directory for manual review for different tables. This test also analyzes the
// foreign keys pointing to customer_entity.
func TestGenerate_Tables_Protobuf_Json(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	// defer dmltest.SQLDumpLoad(t, "testdata/test_*.sql", nil)()
	dmltest.SQLDumpLoad(t, "testdata/test_*.sql", nil)

	ctx := context.Background()
	ts, err := dmlgen.NewTables("github.com/corestoreio/pkg/sql/dmlgen/testdata",

		dmlgen.WithProtobuf(),

		dmlgen.WithLoadColumns(ctx, db.DB, "dmlgen_types", "core_config_data", "customer_entity", "customer_address_entity"),
		dmlgen.WithTableOption(
			"customer_entity", &dmlgen.TableOption{
				Encoders:   []string{"json", "protobuf"},
				StructTags: []string{"max_len"},
			}),
		dmlgen.WithTableOption(
			"customer_address_entity", &dmlgen.TableOption{
				Encoders:   []string{"json", "protobuf"},
				StructTags: []string{"max_len"},
			}),

		dmlgen.WithTableOption(
			"core_config_data", &dmlgen.TableOption{
				Encoders: []string{"easyjson", "protobuf"},
				CustomStructTags: []string{
					"path", `json:"x_path" xml:"y_path" max_len:"255"`,
					"scope_id", `json:"scope_id" xml:"scope_id"`,
				},
				StructTags: []string{"json", "max_len"},
				ColumnAliases: map[string][]string{
					"path": {"storage_location", "config_directory"},
				},
				UniquifiedColumns: []string{"path"},
			}),

		dmlgen.WithTable("core_config_data", ddl.Columns{
			&ddl.Column{Field: "path", Pos: 5, Default: null.MakeString("'general'"), Null: "NO", DataType: "varchar", CharMaxLength: null.MakeInt64(255), ColumnType: "varchar(255)", Comment: "Config Path overwritten"},
		}, "overwrite"),

		dmlgen.WithTableOption(
			"dmlgen_types", &dmlgen.TableOption{
				Encoders:          []string{"easyjson", "binary", "protobuf"},
				StructTags:        []string{"json", "protobuf", "max_len"},
				UniquifiedColumns: []string{"col_varchar_100", "price_12_4a", "col_int_1", "col_int_2", "has_smallint_5", "col_date_2"},
				Comment:           "Just another comment.",
			}),

		dmlgen.WithColumnAliasesFromForeignKeys(ctx, db.DB),
		dmlgen.WithReferenceEntitiesByForeignKeys(ctx, db.DB, func(tableName string) string {
			switch tableName {
			case "customer_address_entity":
				tableName = "Addresses"
			}
			return tableName
		}),

		dmlgen.WithCustomCode("pseudo.MustNewService.Option", `
		pseudo.WithTagFakeFunc("CustomerAddressEntity.ParentID", func(maxLen int) (interface{}, error) {
			return nil, nil 
		}),
		pseudo.WithTagFakeFunc("col_date1", func(maxLen int) (interface{}, error) {
			if ps.Intn(1000)%3 == 0 {
				return nil, nil
			}
			return ps.Dob18(), nil
		}),
		pseudo.WithTagFakeFunc("col_date2", func(maxLen int) (interface{}, error) {
			return ps.Dob18().MarshalText()
		}),
		pseudo.WithTagFakeFunc("col_decimal101", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.1f", ps.Price()), nil
		}),
		pseudo.WithTagFakeFunc("price124b", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.4f", ps.Price()), nil
		}),
		pseudo.WithTagFakeFunc("col_decimal123", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.3f", ps.Float64()), nil
		}),
		pseudo.WithTagFakeFunc("col_decimal206", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.6f", ps.Float64()), nil
		}),
		pseudo.WithTagFakeFunc("col_decimal2412", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.12f", ps.Float64()), nil
		}),
		pseudo.WithTagFakeFuncAlias(
			"col_decimal124", "price124b",
			"price124a", "price124b",
			"col_float", "col_decimal206",
		),
`),
	)
	assert.NoError(t, err)

	ts.ImportPathsTesting = append(ts.ImportPathsTesting, "fmt") // only needed for pseudo functional options.
	ts.TestSQLDumpGlobPath = "test_*_tables.sql"

	writeFile(t, "testdata/output_gen.go", ts.GenerateGo)
	writeFile(t, "testdata/output_gen.proto", ts.GenerateSerializer)
	// Generates for all proto files the Go source code.
	err = dmlgen.GenerateProto("./testdata")
	assert.NoError(t, err, "%+v", err)
	err = dmlgen.GenerateJSON("./testdata", nil)
	assert.NoError(t, err, "%+v", err)
}

func TestInfoSchemaForeignKeys(t *testing.T) {

	t.Skip("One time test. Use when needed to regenerate the code")

	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	ts, err := dmlgen.NewTables("testdata",
		dmlgen.WithTableOption("KEY_COLUMN_USAGE", &dmlgen.TableOption{
			Encoders:          []string{"json", "binary"},
			UniquifiedColumns: []string{"TABLE_NAME", "COLUMN_NAME"},
		}),
		dmlgen.WithLoadColumns(context.Background(), db.DB, "KEY_COLUMN_USAGE"),
	)
	assert.NoError(t, err)

	writeFile(t, "testdata/KEY_COLUMN_USAGE_gen.go", ts.GenerateGo)
}

func TestWithCustomStructTags(t *testing.T) {
	t.Parallel()

	t.Run("unbalanced should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.Fatal.Match(err), "%s", err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()

		tbl, err := dmlgen.NewTables("testdata",
			dmlgen.WithTable("table", ddl.Columns{&ddl.Column{Field: "config_id"}}),
			dmlgen.WithTableOption("table", &dmlgen.TableOption{
				CustomStructTags: []string{"unbalanced"},
			}),
		)
		assert.Nil(t, tbl)
		assert.NoError(t, err)
	})

	t.Run("table not found", func(t *testing.T) {
		tbls, err := dmlgen.NewTables("test",
			dmlgen.WithTableOption("tableNOTFOUND", &dmlgen.TableOption{
				CustomStructTags: []string{"column", "db:..."},
			}),
		)
		assert.Nil(t, tbls)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

	t.Run("column not found", func(t *testing.T) {
		tbls, err := dmlgen.NewTables("test",
			dmlgen.WithTableOption("core_config_data", &dmlgen.TableOption{
				CustomStructTags: []string{"scope_id", "toml:..."},
			}),
			dmlgen.WithTable("core_config_data", ddl.Columns{
				&ddl.Column{Field: "config_id"},
			}),
		)
		assert.Nil(t, tbls)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})
}

func TestWithStructTags(t *testing.T) {
	t.Parallel()

	t.Run("table not found", func(t *testing.T) {
		tbls, err := dmlgen.NewTables("test",
			dmlgen.WithTableOption("tableNOTFOUND", &dmlgen.TableOption{
				StructTags: []string{"unbalanced"},
			}),
		)
		assert.Nil(t, tbls)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

	t.Run("struct tag not supported", func(t *testing.T) {
		tbls, err := dmlgen.NewTables("test",
			dmlgen.WithTableOption("core_config_data", &dmlgen.TableOption{
				StructTags: []string{"hjson"},
			}),
			dmlgen.WithTable("core_config_data", ddl.Columns{
				&ddl.Column{Field: "config_id"},
			}),
		)
		assert.Nil(t, tbls)
		assert.True(t, errors.NotSupported.Match(err), "%+v", err)
	})

	t.Run("al available struct tags", func(t *testing.T) {
		tbls, err := dmlgen.NewTables("test",
			dmlgen.WithTableOption("core_config_data", &dmlgen.TableOption{
				StructTags: []string{"bson", "db", "env", "json", "toml", "yaml", "xml"},
			}),
			dmlgen.WithTable("core_config_data", ddl.Columns{
				&ddl.Column{Field: "config_id"},
			}),
		)
		assert.NoError(t, err)
		have := tbls.Tables["core_config_data"].Columns.ByField("config_id").GoString()
		assert.Exactly(t, "&ddl.Column{Field: \"config_id\", StructTag: \"bson:\\\"config_id,omitempty\\\" db:\\\"config_id\\\" env:\\\"config_id\\\" json:\\\"config_id,omitempty\\\" toml:\\\"config_id\\\" yaml:\\\"config_id,omitempty\\\" xml:\\\"config_id,omitempty\\\"\", }", have)
	})
}

func TestWithColumnAliases(t *testing.T) {
	t.Parallel()

	t.Run("table not found", func(t *testing.T) {
		tbls, err := dmlgen.NewTables("test",
			dmlgen.WithTableOption("tableNOTFOUND", &dmlgen.TableOption{
				ColumnAliases: map[string][]string{"column": {"alias"}},
			}),
		)
		assert.Nil(t, tbls)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

	t.Run("column not found", func(t *testing.T) {
		tbls, err := dmlgen.NewTables("test",
			dmlgen.WithTableOption("tableNOTFOUND", &dmlgen.TableOption{
				ColumnAliases: map[string][]string{"scope_id": {"scopeID"}},
			}),
			dmlgen.WithTable("core_config_data", ddl.Columns{
				&ddl.Column{Field: "config_id"},
			}),
		)
		assert.Nil(t, tbls)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})
}

func TestWithUniquifiedColumns(t *testing.T) {
	t.Parallel()

	t.Run("column not found", func(t *testing.T) {
		tbls, err := dmlgen.NewTables("test",
			dmlgen.WithTableOption("core_config_data", &dmlgen.TableOption{
				UniquifiedColumns: []string{"scope_id", "scopeID"},
			}),

			dmlgen.WithTable("core_config_data", ddl.Columns{
				&ddl.Column{Field: "config_id"},
			}),
		)
		assert.Nil(t, tbls)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})
}
