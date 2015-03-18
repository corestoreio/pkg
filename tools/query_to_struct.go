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

import "database/sql"

// QueryToStruct generates a struct in Go code to match the query
func QueryToStruct(db *sql.DB, query string) ([]byte, error) {

	return
}

const tplQueryStruct = `
type (
    // {{.table | prepareVar}}Slice contains pointers to {{.table | prepareVar}} types
    {{.table | prepareVar}}Slice []*{{.table | prepareVar}}
    // {{.table | prepareVar}} a type for the MySQL table {{ .table }}
    {{.table | prepareVar}} struct {
        {{ range .columns }}{{.GoName}} {{.GoType}} {{ $.Tick }}db:"{{.Field.String}}"{{ $.Tick }} {{.Comment}}
        {{ end }} }
)
`
