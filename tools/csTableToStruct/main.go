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

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"fmt"

	"github.com/corestoreio/csfw/tools"
	"github.com/corestoreio/csfw/tools/toolsdb"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cliDsn := flag.String("dsn", "test:test@tcp(localhost:3306)/test", "MySQL DSN data source name. Can also be provided via ENV with key CS_DSN")
	pkg := flag.String("p", "", "Package name in template")
	run := flag.Bool("run", false, "If true program runs")
	outputFile := flag.String("o", "", "Output file name")

	prefixSearch := flag.String("prefixSearch", "eav", "Search Table Prefix. Used in where condition to list tables")
	prefixName := flag.String("prefixName", "", "Table name prefix")
	flag.Parse()

	if false == *run || *outputFile == "" || *pkg == "" {
		flag.Usage()
		os.Exit(1)
	}
	db, err := toolsdb.Connect(*cliDsn)
	toolsdb.LogFatal(err)
	defer db.Close()

	type dataContainer struct {
		Tables        []map[string]interface{}
		Package, Tick string
	}

	tplData := &dataContainer{
		Tables:  make([]map[string]interface{}, 0, 200),
		Package: *pkg,
		Tick:    "`",
	}

	tables, err := toolsdb.GetTables(db, *prefixName+*prefixSearch)
	toolsdb.LogFatal(err)

	for _, table := range tables {
		columns, err := toolsdb.GetColumns(db, *prefixName+table)
		toolsdb.LogFatal(err)

		structNames := make([]string, len(columns))
		rawColumnNames := make([]string, len(columns))

		for i, c := range columns {
			structNames[i] = c.GoName
			rawColumnNames[i] = c.Field.String
		}

		tplData.Tables = append(tplData.Tables, map[string]interface{}{
			"table":    toolsdb.Camelize(strings.Replace(table, "eav_", "", 1)), //@todo strip table prefix
			"tableOrg": table,
			"columns":  columns,
			//"columnsSelect": "`" + strings.Join(rawColumnNames, "`, `") + "`",
			//"columnsScan":   "e." + strings.Join(structNames, ", e."),
		})
	}

	formatted, err := tools.GenerateCode(tplCode, tplData)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		toolsdb.LogFatal(err)
	}

	ioutil.WriteFile(*outputFile, formatted, 0600)
	log.Println("ok csTableToStruct")
}
