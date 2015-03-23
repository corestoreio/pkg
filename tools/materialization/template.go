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

// Generates code for all EAV types
package main

import "github.com/corestoreio/csfw/tools"

/*
   Data will be "carved in stone" because it only changes during development.
   - DONE: entity_type with translation of some columns to the Go type
   - attribute_set related tables: eav_attribute_set, eav_entity_attribute, eav_attribute_group, etc
   - label and option tables will not be hard coded
   - eav_attribute full config and from that one the flat table structure
*/

const tplEav = tools.Copyright + `
// Package {{ .Package }} file is auto generated via eavToStruct
package {{ .Package }}
import (
    "github.com/corestoreio/csfw/eav"
    {{ range .EntityTypeMap }}{{if ne .ImportPath "" }}"{{.ImportPath}}"
{{end}}{{end}}
)

func init(){
    csEntityTypeCollection = eav.CSEntityTypeSlice{
        {{ range .ETypeData }} &eav.CSEntityType {
            EntityTypeID: {{ .EntityTypeID }},
            EntityTypeCode: "{{ .EntityTypeCode }}",
            EntityModel: {{ .EntityModel }},
            AttributeModel: {{ .AttributeModel.String }},
            EntityTable: {{ .EntityTable.String }},
            ValueTablePrefix: "{{ .ValueTablePrefix.String }}",
            IsDataSharing: {{ .IsDataSharing }},
            DataSharingKey: "{{ .DataSharingKey.String }}",
            DefaultAttributeSetID: {{ .DefaultAttributeSetID }},
            {{ if ne "" .IncrementModel.String }}IncrementModel: {{ .IncrementModel.String }},{{ end }}
            IncrementPerStore: {{ .IncrementPerStore }},
            IncrementPadLength: {{ .IncrementPadLength }},
            IncrementPadChar: "{{ .IncrementPadChar }}",
            AdditionalAttributeTable: {{ .AdditionalAttributeTable.String }},
            EntityAttributeCollection: {{ .EntityAttributeCollection.String }},
        },
        {{ end }}
    }
}
`

const tplTypeDefinition = `
type (
    // {{.Name | prepareVar}}Slice contains pointers to {{.Name | prepareVar}} types
    // @todo website must be present in the slice
    {{.Name | prepareVar}}Slice []*{{.Name | prepareVar}}
    // {{.Name | prepareVar}} a data container for the data from a MySQL query
    {{.Name | prepareVar}} struct {
        {{ range .Columns }}{{.GoName}} {{.GoType}}
        {{ end }} }
)
`

const tplTypeDefinitionFile = tools.Copyright + `
package {{ .PackageName }}
{{ if gt (len .ImportPaths) 0 }}
    import (
    {{ range .ImportPaths }} "{{.}}"
    {{ end }} )
{{ end }}

{{.TypeDefinition}}
type index{{ .Name | prepareVar }} int

const (
    {{ range $k, $row := .Attributes }}{{$.Name | prepareVar}}{{index $row "attribute_code" | prepareVar}} {{ if eq $k 0 }} index{{ $.Name | prepareVar }} = iota {{ end }}
    {{end}}
)

var private{{.Name | prepareVar}}Collection = {{.Name | prepareVar}}Slice{
        {{ range $row := .Attributes }} {{$.Name | prepareVar}}{{index $row "attribute_code" | prepareVar}}: &{{$.Name | prepareVar}} {
            {{ range $k,$v := $row }} {{ $k | prepareVar }}: {{ $v }},
            {{ end }}
        },
        {{ end }}
    }

// Get{{.Name | prepareVar}}Collection returns a slice with all attributes and its configuration
func Get{{.Name | prepareVar}}Collection() {{.Name | prepareVar}}Slice {
    return private{{.Name | prepareVar}}Collection
}

// Get{{.Name | prepareVar}} returns a single attribute
func Get{{.Name | prepareVar}}(i index{{ .Name | prepareVar }}) (*{{.Name | prepareVar}}, error) {
    return private{{.Name | prepareVar}}Collection[i], nil
}

`
