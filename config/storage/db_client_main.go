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

	dbcp := dml.MustConnectAndVerify(dml.WithDSNfromEnv(""))
	defer mustCheckCloseErr(dbcp)

	// we assume the table core_configuration already exists

	ctx := context.Background()
	ts, err := dmlgen.NewGenerator("github.com/corestoreio/pkg/config/storage",

		dmlgen.WithTablesFromDB(ctx, dbcp, "core_configuration"),
		dmlgen.WithBuildTags("csall db"),
		dmlgen.WithTableConfig(
			"core_configuration", &dmlgen.TableConfig{
				UniquifiedColumns: []string{"path"},
				// `max_len` defines for Faker package the maximum size/length for
				// a field. Only used during testing.
				StructTags: []string{"max_len"},
				// DisableCollectionMethods: true,
				FeaturesInclude: dmlgen.FeatureCollectionStruct | dmlgen.FeatureCollectionUniqueGetters | dmlgen.FeatureEntityStruct |
					dmlgen.FeatureDB | dmlgen.FeatureEntityWriteTo,
			}),
	)
	mustCheckErr(err)

	//	ts.TestSQLDumpGlobPath = "test_*_tables.sql"

	writeFile("db_client_schema_gen.go", ts.GenerateGo)
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
