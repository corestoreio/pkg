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
	"database/sql"
	"io/ioutil"

	"fmt"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/tools"
)

type (
	dataContainer struct {
		Tables              []map[string]interface{}
		Package, Tick       string
		TypeCodeValueTables tools.TypeCodeValueTable
	}
)

func main() {
	db, dbrConn, err := csdb.Connect()
	tools.LogFatal(err)
	defer db.Close()
	for _, tStruct := range tools.ConfigTableToStruct {
		generateStructures(tStruct, db, dbrConn)
	}
}

func generateStructures(tStruct *tools.TableToStruct, db *sql.DB, dbrConn *dbr.Connection) {
	tplData := &dataContainer{
		Tables:  make([]map[string]interface{}, 0, 200),
		Package: tStruct.Package,
		Tick:    "`",
	}

	tables, err := tools.GetTables(db, tools.ReplaceTablePrefix(tStruct.QueryString))
	tools.LogFatal(err)

	if len(tStruct.EntityTypeCodes) > 0 && tStruct.EntityTypeCodes[0] != "" {
		tplData.TypeCodeValueTables, err = tools.GetEavValueTables(dbrConn, tStruct.EntityTypeCodes)
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

		columns, err := tools.GetColumns(db, table)
		tools.LogFatal(err)
		tools.LogFatal(columns.MapSQLToGoDBRType())
		tplData.Tables = append(tplData.Tables, map[string]interface{}{
			"table":   table,
			"columns": columns,
		})
	}

	formatted, err := tools.GenerateCode(tStruct.Package, tplCode, tplData)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		tools.LogFatal(err)
	}

	tools.LogFatal(ioutil.WriteFile(tStruct.OutputFile, formatted, 0600))
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
