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

// Generates code for all attribute types
package main

import "github.com/corestoreio/csfw/tools"

const tplTypeDefinitions = `
type (
    // {{.Name | prepareVar}}Slice contains pointers to {{.Name | prepareVar}} types
    {{.Name | prepareVar}}Slice []*{{.Name | prepareVar}}
    // {{.Name | prepareVar}} a data container for the data from a MySQL query
    {{.Name | prepareVar}} struct {
        {{ range .Columns }}{{.GoName}} {{.GoType}}
        {{ end }} }
)
`

const tplFileBody = tools.Copyright + `
package {{ .PackageName }}
{{ if gt (len .ImportPaths) 0 }}
    import (
    {{ range .ImportPaths }} "{{.}}"
    {{ end }} )
{{ end }}

{{.TypeDefinition}}

var private{{.Name | prepareVar}}Collection = {{.Name | prepareVar}}Slice{
        {{ range $row := .Attributes }} &{{$.Name | prepareVar}} {
            {{ range $k,$v := $row }} {{ $k | prepareVar }}: {{ $v }},
            {{ end }}
        },
        {{ end }}
    }
`
