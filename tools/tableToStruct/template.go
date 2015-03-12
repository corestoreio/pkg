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

const tplCode = `// Copyright 2015 CoreStore Authors
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

// Package {{ .Package }} is auto generated via csTableToStruct
package {{ .Package }}
import (
	"time"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/gocraft/dbr"
)

const (
    // TableNoop has index 0
    TableNoop csdb.Index = iota
    {{ range .Tables }} // Table{{.table}} is the index to {{.tableOrg}}
    Table{{.table}}
{{ end }} // TableMax represents the maximum index, which is not available.
TableMax
)

var (
    // read only map
    tableMap = csdb.TableMap{
{{ range .Tables }}Table{{.table}} : csdb.NewTableStructure(
        "{{.tableOrg}}",
        []string{
        {{ range .columns }}{{ if eq .Key.String "PRI" }} "{{.Field.String}}",{{end}}
        {{ end }} },
        []string {
        {{ range .columns }} "{{.Field.String}}",
        {{ end }} },
    ),
    {{ end }}
    }
)

// GetTableStructure returns for a given index i the table structure or an error it not found.
func GetTableStructure(i csdb.Index) (*csdb.TableStructure, error) {
	return tableMap.Structure(i)
}

// GetTableName returns for a given index the table name. If not found an empty string.
func GetTableName(i csdb.Index) string {
	return tableMap.Name(i)
}

type (

{{ range .Tables }}
    // {{.table}}Slice contains pointers to {{.table}} types
    {{.table}}Slice []*{{.table}}
    // {{.table}} a type for the MySQL table {{ .tableOrg }}
    {{.table}} struct {
        {{ range .columns }}{{.GoName}} {{.GoType}} {{ $.Tick }}db:"{{.Field.String}}"{{ $.Tick }} {{.Comment}}
        {{ end }} }
{{ end }}
)
`

// maybe for later use include in tpl
//const tplBody = `
//func Select{{.table}}(db *sql.DB, sqlWhere ...string) ({{.table}}Slice, error) {
//	rows, err := db.Query("SELECT {{.columnsSelect}} FROM {{ quote .tableOrg }} "+strings.Join(sqlWhere," "))
//	if err != nil {
//		return nil,err
//	}
//	defer rows.Close()
//    var c = make({{.table}}Slice, 0, 200)
//	for rows.Next() {
//		e := &{{.table}}{}
//		err := rows.Scan({{.columnsScan}})
//		if err != nil {
//			return nil,err
//		}
//		c = append(c, e)
//	}
//	err = rows.Err()
//	if err != nil {
//		return nil,err
//	}
//	return c, nil
//}
//`
