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

// Package {{ .Package }} is auto generated via tableToStruct
package {{ .Package }}
import (
	"time"
    {{ if not .TypeCodeValueTables.Empty }}
	"github.com/corestoreio/csfw/eav"{{end}}
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
)

const (
    // TableNoop has index 0
    TableNoop csdb.Index = iota
    {{ range .Tables }} // Table{{.table | prepareVar}} is the index to {{.table}}
    Table{{.table | prepareVar}}
{{ end }} // TableMax represents the maximum index, which is not available.
TableMax
)

var (
    // Always reference these packages, just in case the auto-generated code
    // below doesn't.
    _ = time.Time{}

    tableMap = csdb.TableMap{
{{ range .Tables }}Table{{.table | prepareVar}} : csdb.NewTableStructure(
        "{{.table}}",
        []string{
        {{ range .columns }}{{ if eq .Key.String "PRI" }} "{{.Field.String}}",{{end}}
        {{ end }} },
        []string {
        {{ range .columns }}{{ if ne .Key.String "PRI" }} "{{.Field.String}}",{{end}}
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

{{ if not .TypeCodeValueTables.Empty }}
{{range $typeCode,$valueTables := .TypeCodeValueTables}}
// Get{{ $typeCode | prepareVar }}ValueStructure returns for an eav value index the table structure.
// Important also if you have custom value tables
func Get{{ $typeCode | prepareVar }}ValueStructure(i eav.ValueIndex) (*csdb.TableStructure, error) {
	switch i {
	{{range $vt,$v := $valueTables }}case eav.EntityType{{ $v | prepareVar }}:
		return GetTableStructure(Table{{ $vt | prepareVar }})
    {{end}}	}
	return nil, eav.ErrEntityTypeValueNotFound
}
{{end}}{{end}}

type (

{{ range .Tables }}
    // {{.table | prepareVar}}Slice contains pointers to {{.table | prepareVar}} types
    {{.table | prepareVar}}Slice []*{{.table | prepareVar}}
    // {{.table | prepareVar}} a type for the MySQL table {{ .table }}
    {{.table | prepareVar}} struct {
        {{ range .columns }}{{.GoName}} {{.GoType}} {{ $.Tick }}db:"{{.Field.String}}"{{ $.Tick }} {{.Comment}}
        {{ end }} }
{{ end }}
)
`
