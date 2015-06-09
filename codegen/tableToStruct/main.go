// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"strings"

	"fmt"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
)

type (
	dataContainer struct {
		Tables              []map[string]interface{}
		Package, Tick       string
		TypeCodeValueTables codegen.TypeCodeValueTable
	}
)

func main() {
	db, dbrConn, err := csdb.Connect()
	codegen.LogFatal(err)
	defer db.Close()
	for _, tStruct := range codegen.ConfigTableToStruct {
		generateStructures(tStruct, db, dbrConn)
	}
}

func generateStructures(tStruct *codegen.TableToStruct, db *sql.DB, dbrConn *dbr.Connection) {
	tplData := &dataContainer{
		Tables:  make([]map[string]interface{}, 0, 200),
		Package: tStruct.Package,
		Tick:    "`",
	}

	tables, err := codegen.GetTables(db, codegen.ReplaceTablePrefix(tStruct.SQLQuery))
	codegen.LogFatal(err)

	if len(tStruct.EntityTypeCodes) > 0 && tStruct.EntityTypeCodes[0] != "" {
		tplData.TypeCodeValueTables, err = codegen.GetEavValueTables(dbrConn, tStruct.EntityTypeCodes)
		codegen.LogFatal(err)

		for _, vTables := range tplData.TypeCodeValueTables {
			for t := range vTables {
				if false == isDuplicate(tables, t) {
					tables = append(tables, t)
				}
			}
		}
	}

	for _, table := range tables {

		columns, err := codegen.GetColumns(db, table)
		codegen.LogFatal(err)
		codegen.LogFatal(columns.MapSQLToGoDBRType())
		var name = table
		if mappedName, ok := codegen.TableMapMagento1To2[strings.Replace(table, codegen.TablePrefix, "", 1)]; ok {
			name = mappedName
		}
		tplData.Tables = append(tplData.Tables, map[string]interface{}{
			"name":    name,
			"table":   table,
			"columns": columns,
		})
	}

	formatted, err := codegen.GenerateCode(tStruct.Package, tplCode, tplData, nil)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		codegen.LogFatal(err)
	}

	codegen.LogFatal(ioutil.WriteFile(tStruct.OutputFile, formatted, 0600))
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
