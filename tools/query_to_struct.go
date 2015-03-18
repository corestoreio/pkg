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

package tools

import (
	"database/sql"

	"strings"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

// QueryToStruct generates from a SQL query a Go type struct. Arguments/binds to a query are not considered
// dbSelect argument can be nil but then you must provide a query string
func QueryToStruct(db *sql.DB, name string, dbSelect *dbr.SelectBuilder, query ...string) ([]byte, error) {
	const tplQueryStruct = `
type (
    // {{.Name | prepareVar}}Slice contains pointers to {{.Name | prepareVar}} types
    {{.Name | prepareVar}}Slice []*{{.Name | prepareVar}}
    // {{.Name | prepareVar}} a type for a MySQL Query
    {{.Name | prepareVar}} struct {
        {{ range .Columns }}{{.GoName}} {{.GoType}} {{ $.Tick }}db:"{{.Field.String}}"{{ $.Tick }} {{.Comment}}
        {{ end }} }
)
`

	tableName := "csQryToStruct_" + name
	dropTable := func() {
		_, err := db.Exec("DROP TABLE IF EXISTS `" + tableName + "`")
		if err != nil {
			panic(err)
		}
	}
	dropTable()
	defer dropTable()
	qry := strings.Join(query, " ")
	if qry == "" && dbSelect != nil {
		qry, _ = dbSelect.ToSql() // discard the arguments
	}
	_, err := db.Exec("CREATE TABLE `" + tableName + "` AS " + qry)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	cols, err := GetColumns(db, tableName)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	tplData := struct {
		Name    string
		Columns []column
		Tick    string
	}{
		Name:    name,
		Columns: cols,
		Tick:    "`",
	}
	return GenerateCode("", tplQueryStruct, tplData)
}
