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
	"github.com/corestoreio/pkg/util/codegen"
	"github.com/corestoreio/pkg/util/strs"
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
	if err != nil {
		if fe, ok := err.(*codegen.FormatError); ok {
			t.Errorf("Formatting failed: %s\n%s", fe.Error(), fe.Code)
			return
		}
	}
	assert.NoError(t, err, "%+v", err)
}

// TestNewGenerator_Protobuf_Json writes a Go and Proto file to the
// dmltestgenerated directory for manual review for different tables. This test
// also analyzes the foreign keys pointing to customer_entity.
func TestNewGenerator_Protobuf_Json(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	defer dmltest.SQLDumpLoad(t, "testdata/test_*.sql", nil).Deferred()
	// dmltest.SQLDumpLoad(t, "testdata/test_*.sql", nil)

	ctx := context.Background()
	g, err := dmlgen.NewGenerator("github.com/corestoreio/pkg/sql/dmlgen/dmltestgenerated",

		dmlgen.WithBuildTags("!ignore", "!ignored"),
		dmlgen.WithProtobuf(&dmlgen.SerializerConfig{}),

		dmlgen.WithTablesFromDB(ctx, db,
			"dmlgen_types", "core_configuration", "customer_entity", "customer_address_entity",
			"catalog_product_index_eav_decimal_idx", "sales_order_status_state",
			"view_customer_no_auto_increment", "view_customer_auto_increment",
		),

		dmlgen.WithTableConfig(
			"customer_entity", &dmlgen.TableConfig{
				Encoders:      []string{"json", "protobuf"},
				StructTags:    []string{"max_len"},
				PrivateFields: []string{"password_hash"},
			}),
		dmlgen.WithTableConfig(
			"customer_address_entity", &dmlgen.TableConfig{
				Encoders:   []string{"json", "protobuf"},
				StructTags: []string{"max_len"},
			}),

		dmlgen.WithTableConfig("catalog_product_index_eav_decimal_idx", &dmlgen.TableConfig{}),
		dmlgen.WithTableConfig("sales_order_status_state", &dmlgen.TableConfig{
			Encoders:   []string{"json", "protobuf"},
			StructTags: []string{"max_len"},
		}),
		dmlgen.WithTableConfig("view_customer_no_auto_increment", &dmlgen.TableConfig{
			Encoders:   []string{"json", "protobuf"},
			StructTags: []string{"max_len"},
		}),
		dmlgen.WithTableConfig("view_customer_auto_increment", &dmlgen.TableConfig{
			Encoders:   []string{"json", "protobuf"},
			StructTags: []string{"max_len"},
		}),

		dmlgen.WithTableConfig(
			"core_configuration", &dmlgen.TableConfig{
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

		dmlgen.WithTable("core_configuration", ddl.Columns{
			&ddl.Column{Field: "path", Pos: 5, Default: null.MakeString("'general'"), Null: "NO", DataType: "varchar", CharMaxLength: null.MakeInt64(255), ColumnType: "varchar(255)", Comment: "Config Path overwritten"},
		}, "overwrite"),

		dmlgen.WithTableConfig(
			"dmlgen_types", &dmlgen.TableConfig{
				Encoders:          []string{"easyjson", "protobuf"},
				StructTags:        []string{"json", "protobuf", "max_len"},
				UniquifiedColumns: []string{"col_varchar_100", "price_a_12_4", "col_int_1", "col_int_2", "has_smallint_5", "col_date_2"},
				Comment:           "Just another comment.",
			}),

		// dmlgen.WithColumnAliasesFromForeignKeys(ctx, db.DB),
		dmlgen.WithForeignKeyRelationships(ctx, db.DB, dmlgen.ForeignKeyOptions{
			ExcludeRelationships: []string{"customer_address_entity.parent_id", "customer_entity.entity_id"},
		},
		),

		dmlgen.WithCustomCode("pseudo.MustNewService.Option", `
		pseudo.WithTagFakeFunc("dmltestgenerated.CustomerAddressEntity.ParentID", func(maxLen int) (interface{}, error) {
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
		pseudo.WithTagFakeFunc("price_b124", func(maxLen int) (interface{}, error) {
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
			"col_decimal124", "price_b124",
			"price_a124", "price_b124",
			"col_float", "col_decimal206",
		),
`),
	)
	assert.NoError(t, err)

	g.ImportPathsTesting = append(g.ImportPathsTesting, "fmt") // only needed for pseudo functional options.
	g.TestSQLDumpGlobPath = "../testdata/test_*_tables.sql"

	writeFile(t, "dmltestgenerated/output_gen.go", g.GenerateGo)
	writeFile(t, "dmltestgenerated/output_gen.proto", g.GenerateSerializer)
	// // Generates for all proto files the Go source code.

	assert.NoError(t, dmlgen.GenerateProto("./dmltestgenerated", &dmlgen.ProtocOptions{
		ProtoGen: "gogo",
	}))
	assert.NoError(t, dmlgen.GenerateJSON("./dmltestgenerated", "", nil))
}

func TestInfoSchemaForeignKeys(t *testing.T) {
	t.Skip("One time test. Use when needed to regenerate the code")

	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	ts, err := dmlgen.NewGenerator("dmltestgenerated",
		dmlgen.WithTableConfig("KEY_COLUMN_USAGE", &dmlgen.TableConfig{
			Encoders:          []string{"json", "binary"},
			UniquifiedColumns: []string{"TABLE_NAME", "COLUMN_NAME"},
		}),
		dmlgen.WithTablesFromDB(context.Background(), db, "KEY_COLUMN_USAGE"),
	)
	assert.NoError(t, err)

	writeFile(t, "dmltestgenerated/KEY_COLUMN_USAGE_gen.go", ts.GenerateGo)
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

		tbl, err := dmlgen.NewGenerator("dmltestgenerated",
			dmlgen.WithTable("table", ddl.Columns{&ddl.Column{Field: "config_id"}}),
			dmlgen.WithTableConfig("table", &dmlgen.TableConfig{
				CustomStructTags: []string{"unbalanced"},
			}),
		)
		assert.Nil(t, tbl)
		assert.NoError(t, err)
	})

	t.Run("table not found", func(t *testing.T) {
		tbls, err := dmlgen.NewGenerator("test",
			dmlgen.WithTableConfig("tableNOTFOUND", &dmlgen.TableConfig{
				CustomStructTags: []string{"column", "db:..."},
			}),
		)
		assert.Nil(t, tbls)
		assert.ErrorIsKind(t, errors.NotFound, err)
	})
}

func TestWithStructTags(t *testing.T) {
	t.Parallel()

	t.Run("table not found", func(t *testing.T) {
		tbls, err := dmlgen.NewGenerator("test",
			dmlgen.WithTableConfig("tableNOTFOUND", &dmlgen.TableConfig{
				StructTags: []string{"unbalanced"},
			}),
		)
		assert.Nil(t, tbls)
		assert.ErrorIsKind(t, errors.NotFound, err)
	})

	t.Run("struct tag not supported", func(t *testing.T) {
		tbls, err := dmlgen.NewGenerator("test",
			dmlgen.WithTableConfig("core_configuration", &dmlgen.TableConfig{
				StructTags: []string{"hjson"},
			}),
			dmlgen.WithTable("core_configuration", ddl.Columns{
				&ddl.Column{Field: "config_id"},
			}),
		)
		assert.Nil(t, tbls)
		assert.True(t, errors.NotSupported.Match(err), "%+v", err)
	})

	t.Run("al available struct tags", func(t *testing.T) {
		tbls, err := dmlgen.NewGenerator("test",
			dmlgen.WithTableConfig("core_configuration", &dmlgen.TableConfig{
				StructTags: []string{"bson", "db", "env", "json", "toml", "yaml", "xml"},
			}),
			dmlgen.WithTable("core_configuration", ddl.Columns{
				&ddl.Column{Field: "config_id"},
			}),
		)
		assert.NoError(t, err)
		have := tbls.Tables["core_configuration"].Columns.ByField("config_id").GoString()
		assert.Exactly(t, "&ddl.Column{Field: \"config_id\", StructTag: \"bson:\\\"config_id,omitempty\\\" db:\\\"config_id\\\" env:\\\"config_id\\\" json:\\\"config_id,omitempty\\\" toml:\\\"config_id\\\" yaml:\\\"config_id,omitempty\\\" xml:\\\"config_id,omitempty\\\"\", }", have)
	})
}

func TestWithColumnAliases(t *testing.T) {
	t.Parallel()

	t.Run("table not found", func(t *testing.T) {
		tbls, err := dmlgen.NewGenerator("test",
			dmlgen.WithTableConfig("tableNOTFOUND", &dmlgen.TableConfig{
				ColumnAliases: map[string][]string{"column": {"alias"}},
			}),
		)
		assert.Nil(t, tbls)
		assert.ErrorIsKind(t, errors.NotFound, err)
	})

	t.Run("column not found", func(t *testing.T) {
		tbls, err := dmlgen.NewGenerator("test",
			dmlgen.WithTableConfig("tableNOTFOUND", &dmlgen.TableConfig{
				ColumnAliases: map[string][]string{"scope_id": {"scopeID"}},
			}),
			dmlgen.WithTable("core_configuration", ddl.Columns{
				&ddl.Column{Field: "config_id"},
			}),
		)
		assert.Nil(t, tbls)
		assert.ErrorIsKind(t, errors.NotFound, err)
	})
}

func TestWithUniquifiedColumns(t *testing.T) {
	t.Parallel()

	t.Run("column not found", func(t *testing.T) {
		tbls, err := dmlgen.NewGenerator("test",
			dmlgen.WithTableConfig("core_configuration", &dmlgen.TableConfig{
				UniquifiedColumns: []string{"scope_id", "scopeID"},
			}),

			dmlgen.WithTable("core_configuration", ddl.Columns{
				&ddl.Column{Field: "config_id"},
			}),
		)
		assert.Nil(t, tbls)
		assert.ErrorIsKind(t, errors.NotFound, err)
	})
}

func TestNewGenerator_NoDB(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	defer dmltest.SQLDumpLoad(t, "testdata/test_*.sql", nil).Deferred()
	// dmltest.SQLDumpLoad(t, "testdata/test_*.sql", nil)

	ctx := context.Background()
	ts, err := dmlgen.NewGenerator("github.com/corestoreio/pkg/sql/dmlgen/dmltestgenerated2",

		dmlgen.WithTablesFromDB(ctx, db,
			"core_configuration", "sales_order_status_state", "view_customer_auto_increment",
		),

		dmlgen.WithTableConfigDefault(dmlgen.TableConfig{
			Encoders:        []string{"json", "protobuf"},
			StructTags:      []string{"max_len"},
			FeaturesExclude: dmlgen.FeatureDB | dmlgen.FeatureCollectionUniquifiedGetters,
		}),

		dmlgen.WithTableConfig("view_customer_auto_increment", &dmlgen.TableConfig{
			StructTags:      []string{"yaml"},
			FeaturesInclude: dmlgen.FeatureCollectionFilter | dmlgen.FeatureCollectionEach,
		}),
	)
	assert.NoError(t, err)

	ts.ImportPathsTesting = append(ts.ImportPathsTesting, "fmt") // only needed for pseudo functional options.
	ts.TestSQLDumpGlobPath = "../testdata/test_*_tables.sql"

	writeFile(t, "dmltestgenerated2/no_db_gen.go", ts.GenerateGo)
}

func TestNewGenerator_ReversedForeignKeys(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	defer dmltest.SQLDumpLoad(t, "testdata/test_*.sql", nil).Deferred()
	// dmltest.SQLDumpLoad(t, "testdata/test_*.sql", nil)

	ctx := context.Background()
	ts, err := dmlgen.NewGenerator("github.com/corestoreio/pkg/sql/dmlgen/dmltestgenerated3",

		dmlgen.WithTablesFromDB(ctx, db,
			"store", "store_group", "store_website",
			"customer_entity", "customer_address_entity",
			"catalog_category_entity", "sequence_catalog_category",
		),

		dmlgen.WithTableConfigDefault(dmlgen.TableConfig{
			Encoders:        []string{"json", "protobuf"},
			StructTags:      []string{"max_len"},
			FeaturesInclude: dmlgen.FeatureEntityStruct | dmlgen.FeatureCollectionStruct | dmlgen.FeatureEntityRelationships,
		}),
		// Just an empty TableConfig to trigger the default config update for
		// this table. Hacky for now.
		dmlgen.WithTableConfig("customer_address_entity", &dmlgen.TableConfig{}),

		dmlgen.WithTableConfig("customer_entity", &dmlgen.TableConfig{
			FieldMapFn: func(dbIdentifier string) (fieldName string) {
				switch dbIdentifier {
				case "customer_address_entity":
					return "Address"
				}
				return strs.ToGoCamelCase(dbIdentifier)
			},
		}),

		dmlgen.WithTableConfig("store", &dmlgen.TableConfig{
			CustomStructTags: []string{
				"store_website", `json:"-"`,
				"store_group", `json:"-"`,
			},
		}),

		dmlgen.WithForeignKeyRelationships(ctx, db.DB, dmlgen.ForeignKeyOptions{
			// IncludeRelationShips: []string{"what are the names?"},
			ExcludeRelationships: []string{
				"store_website.website_id", "customer_entity.website_id",
				"store.store_id", "customer_entity.store_id",
				"customer_entity.store_id", "store.store_id",
				"customer_entity.website_id", "store_website.website_id",
				"customer_address_entity.parent_id", "customer_entity.entity_id",
			},
		},
		),
	)
	assert.NoError(t, err)

	ts.ImportPathsTesting = append(ts.ImportPathsTesting, "fmt") // only needed for pseudo functional options.
	ts.TestSQLDumpGlobPath = "../testdata/test_*_tables.sql"

	writeFile(t, "dmltestgenerated3/rev_fk_gen.go", ts.GenerateGo)
}

func TestNewGenerator_MToMForeignKeys(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	defer dmltest.SQLDumpLoad(t, "testdata/test_*.sql", nil).Deferred()
	// dmltest.SQLDumpLoad(t, "testdata/test_*.sql", nil)

	ctx := context.Background()
	ts, err := dmlgen.NewGenerator("github.com/corestoreio/pkg/sql/dmlgen/dmltestgeneratedMToM",

		dmlgen.WithTablesFromDB(ctx, db,
			//"athlete_team_member",
			"athlete_team", "athlete",
			"customer_entity", "customer_address_entity",
		),

		dmlgen.WithTableConfigDefault(dmlgen.TableConfig{
			StructTags:      []string{"max_len"},
			FeaturesInclude: dmlgen.FeatureEntityStruct | dmlgen.FeatureCollectionStruct | dmlgen.FeatureEntityRelationships,
		}),
		// Just an empty TableConfig to trigger the default config update for
		// this table. Hacky for now.
		dmlgen.WithTableConfig("athlete_team", &dmlgen.TableConfig{}),
		dmlgen.WithTableConfig("athlete", &dmlgen.TableConfig{}),

		dmlgen.WithTableConfig("customer_entity", &dmlgen.TableConfig{
			FieldMapFn: func(dbIdentifier string) (fieldName string) {
				switch dbIdentifier {
				case "customer_address_entity":
					return "Address"
				}
				return strs.ToGoCamelCase(dbIdentifier)
			},
		}),

		dmlgen.WithForeignKeyRelationships(ctx, db.DB, dmlgen.ForeignKeyOptions{
			// IncludeRelationShips: []string{"what are the names?"},
			ExcludeRelationships: []string{
				//"athlete.athlete_id", "athlete_team_member.athlete_id",

				"athlete_team.team_id", "athlete_team_member.team_id",
				//"athlete_team.team_id", "athlete.athlete_id",
				"athlete_team_member.*", "*.*", // do not print relations for the relation table itself.

				"customer_address_entity.parent_id", "customer_entity.entity_id",
			},
		},
		),
	)
	assert.NoError(t, err)

	ts.ImportPathsTesting = append(ts.ImportPathsTesting, "fmt") // only needed for pseudo functional options.
	ts.TestSQLDumpGlobPath = "../testdata/test_*_tables.sql"

	writeFile(t, "dmltestgeneratedMToM/fkm2n_gen.go", ts.GenerateGo)
}
