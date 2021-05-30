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

// +build ignore

package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmlgen"
)

func main() {
	const pkgPath = "github.com/corestoreio/pkg/store"

	dbcp := dml.MustConnectAndVerify(dml.WithDSNFromEnv(""))
	defer mustCheckCloseErr(dbcp)

	ctx := context.Background()

	g, err := dmlgen.NewGenerator(pkgPath,
		dmlgen.WithTablesFromDB(ctx, dbcp, "store_website", "store_group", "store"),
		dmlgen.WithTableConfig(
			"store_website", &dmlgen.TableConfig{
				Encoders:        []string{"easyjson", "protobuf"},
				FeaturesExclude: dmlgen.FeatureDB,
				StructTags:      []string{"max_len"},
			}),
		dmlgen.WithTableConfig(
			"store_group", &dmlgen.TableConfig{
				Encoders:        []string{"easyjson", "protobuf"},
				FeaturesExclude: dmlgen.FeatureDB,
				StructTags:      []string{"max_len"},
			}),
		dmlgen.WithTableConfig(
			"store", &dmlgen.TableConfig{
				Encoders:         []string{"easyjson", "protobuf"},
				FeaturesExclude:  dmlgen.FeatureDB,
				StructTags:       []string{"max_len"},
				CustomStructTags: []string{"StoreGroup", `faker:"-"`, "StoreWebsite", `faker:"-"`},
			}),
		dmlgen.WithForeignKeyRelationships(ctx, dbcp.DB,
			"store_group.group_id", "store.group_id",
		),
		dmlgen.WithProtobuf(&dmlgen.SerializerConfig{
			PackageImportPath: pkgPath,
		}),
	)

	mustCheckErr(err)
	// 	g.TestSQLDumpGlobPath = "test_*_tables.sql"
	writeFile("entities_gen.go", g.GenerateGo)
	writeFile("entities_gen.proto", g.GenerateSerializer)

	// write MySQL/MariaDB DB code
	g, err = dmlgen.NewGenerator(pkgPath,
		dmlgen.WithTablesFromDB(ctx, dbcp, "store_website", "store_group", "store"),

		dmlgen.WithBuildTags("csall db"),

		dmlgen.WithTableConfig(
			"store_website", &dmlgen.TableConfig{
				FeaturesInclude: dmlgen.FeatureDB,
			}),
		dmlgen.WithTableConfig(
			"store_group", &dmlgen.TableConfig{
				FeaturesInclude: dmlgen.FeatureDB,
			}),
		dmlgen.WithTableConfig(
			"store", &dmlgen.TableConfig{
				FeaturesInclude: dmlgen.FeatureDB,
			}),
		// Protobuf needed here to adjust the DB/Go types to protobuf.
		dmlgen.WithProtobuf(&dmlgen.SerializerConfig{}),
	)
	mustCheckErr(err)
	// 	g.TestSQLDumpGlobPath = "test_*_tables.sql"
	writeFile("entities_db_gen.go", g.GenerateGo)

	mustCheckErr(dmlgen.RunProtoc("./", &dmlgen.ProtocOptions{
		GRPC: true,
	}))
}

func writeFile(outFile string, wFn func(io.Writer, io.Writer) error) {
	testF := ioutil.Discard
	if strings.HasSuffix(outFile, ".go") {
		testFile := strings.Replace(outFile, ".go", "_test.go", 1)

		ft, err := os.Create(testFile)
		mustCheckErr(err)
		defer mustCheckCloseErr(ft)
		testF = ft
	}

	f, err := os.Create(outFile)
	mustCheckErr(err)
	defer mustCheckCloseErr(f)
	err = wFn(f, testF)
	mustCheckErr(err)
}

func mustCheckErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v\n", err))
	}
}

func mustCheckCloseErr(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(fmt.Sprintf("%+v\n", err))
	}
}
