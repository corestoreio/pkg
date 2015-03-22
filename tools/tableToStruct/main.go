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
	"os"
	"strings"

	"fmt"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/tools"
)

type (
	dataContainer struct {
		Tables              []map[string]interface{}
		Package, Tick       string
		TypeCodeValueTables tools.TypeCodeValueTable
	}
)

var (
	pkg        = flag.String("p", "", "Package name in template")
	run        = flag.Bool("run", false, "If true program runs")
	outputFile = flag.String("o", "", "Output file name")

	prefixSearch   = flag.String("prefixSearch", "", "Search Table Prefix. Used in where condition to list tables")
	prefixName     = flag.String("prefixName", "", "Table name prefix") // @todo via env var !?
	entityTypeCode = flag.String("entityTypeCodes", "", "If provided then eav_entity_type.value_table_prefix will be evaluated for further tables. Use comma to separate codes.")
)

func main() {
	flag.Parse()

	if false == *run || *outputFile == "" || *pkg == "" {
		flag.Usage()
		os.Exit(1)
	}

	db, dbrConn, err := csdb.Connect()
	tools.LogFatal(err)
	defer db.Close()

	tplData := &dataContainer{
		Tables:  make([]map[string]interface{}, 0, 200),
		Package: *pkg,
		Tick:    "`",
	}

	tables, err := tools.GetTables(db, *prefixName+*prefixSearch)
	tools.LogFatal(err)

	entityTypeCodes := strings.Split(*entityTypeCode, ",")
	if len(entityTypeCodes) > 0 {
		tplData.TypeCodeValueTables, err = tools.GetEavValueTables(dbrConn, *prefixName, entityTypeCodes)
		tools.LogFatal(err)
		for _, vTables := range tplData.TypeCodeValueTables {
			for t, _ := range vTables {
				if false == isDuplicate(tables, t) {
					tables = append(tables, t)
				}
			}
		}
	}

	for _, table := range tables {

		if skipCatalogFlatTable(table) {
			continue
		}

		columns, err := tools.GetColumns(db, *prefixName+table)
		tools.LogFatal(err)
		tools.LogFatal(columns.MapSQLToGoDBRType())
		tplData.Tables = append(tplData.Tables, map[string]interface{}{
			"table":   table,
			"columns": columns,
		})
	}

	formatted, err := tools.GenerateCode(*pkg, tplCode, tplData)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		tools.LogFatal(err)
	}

	ioutil.WriteFile(*outputFile, formatted, 0600)
}

// isDuplicate slow duplicate checker ...
func isDuplicate(sl []string, st string) bool {
	for _, s := range sl {
		if s == st {
			return true
		}
	}
	return false
}

// shouldSkipTable checks if a table is a catalog*flat* table. These tables will get automatically created
// due to the variable attributes which are used as columns. And also dependent on the store count.
func skipCatalogFlatTable(table string) bool {
	return strings.Index(table, "catalog_") == 0 && strings.Index(table, "_flat_") > 6
}
