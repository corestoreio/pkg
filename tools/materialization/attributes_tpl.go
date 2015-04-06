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

import "github.com/corestoreio/csfw/tools"

/* @todo
   Data will be "carved in stone" because it only changes during development.
   - attribute_set related tables: eav_attribute_set, eav_entity_attribute, eav_attribute_group, etc
   - label and option tables will not be hard coded
*/

const tplTypeDefinition = `
type (
    // @todo website must be present in the slice
    // {{ .Name | prepareVar }} a data container for attributes. You can use this struct to
    // embed into your own struct for maybe overriding some method receivers.
    {{ .Name | prepareVar | toLowerFirst }} struct {
        *eav.Attribute
        {{ range .Columns }}{{ .GoName | toLowerFirst }} {{ .GoType }}
        {{ end }} }
)

{{ range .Columns }} func (a *{{ $.Name | prepareVar | toLowerFirst }}) {{ .GoName }}() {{ .GoType }}{
    return a.{{ .GoName | toLowerFirst }}
}
{{ end }}

// Check if Attributer interface has been successfully implemented
var _ {{ .AttrPkg }}.Attributer = (*{{ .Name | prepareVar | toLowerFirst }})(nil)

`

const tplTypeDefinitionFile = tools.Copyright + `
package {{ .PackageName }}
    import (
        "github.com/corestoreio/csfw/eav"
        "{{ .AttrPkgImp }}"
        {{ range .ImportPaths }}"{{ . }}"
        {{ end }} )

{{ .TypeDefinition }}

const (
    {{ range $k, $row := .Attributes }}{{ $.Name | prepareVar }}{{ index $row "attribute_code" | prepareVar }} {{ if eq $k 0 }} eav.AttributeIndex = iota + 1{{ end }}
    {{ end }}
    {{ $.Name | prepareVar }}ZZZ
)

type si{{ $.Name | prepareVar }} struct {}

func (si{{ $.Name | prepareVar }}) ByID(id int64) (eav.AttributeIndex, error){
	switch id {
	{{ range $k, $row := .Attributes }} case {{ index $row "attribute_id" }}:
		return {{ $.Name | prepareVar }}{{ index $row "attribute_code" | prepareVar }}, nil
	{{ end }}
	default:
		return eav.AttributeIndex(0), eav.ErrAttributeNotFound
	}
}

func (si{{ $.Name | prepareVar }}) ByCode(code string) (eav.AttributeIndex, error){
	switch code {
	{{ range $k, $row := .Attributes }} case {{ index $row "attribute_code" }}:
		return {{ $.Name | prepareVar }}{{ index $row "attribute_code" | prepareVar }}, nil
	{{ end }}
	default:
		return eav.AttributeIndex(0), eav.ErrAttributeNotFound
	}
}

var _ eav.AttributeGetter = (*si{{ $.Name | prepareVar }})(nil)

func init(){
    {{ .AttrPkg }}.{{ .FuncGetter }}(&si{{ $.Name | prepareVar }}{})
    {{ .AttrPkg }}.{{ .FuncCollection }}({{ .AttrPkg }}.AttributeSlice{
        {{ range $row := .Attributes }}
        {{ $const := sprintf "%s%s" (prepareVar $.Name) (prepareVar (index $row "attribute_code")) }}
        {{ $const }}: {{ if ne $.MyStruct "" }} &{{ $.MyStruct }} {
        {{ end }} &{{ $.Name | prepareVar | toLowerFirst }} {
            Attribute: eav.NewAttribute({{ range $k,$v := $row }} {{ if (isEavAttr $k) }} {{ setAttrIdx $v $const }}, // {{ $k }}
            {{ end }}{{ end }}
            ),
            {{ range $k,$v := $row }} {{ if not (isEavAttr $k) }} {{ $k | prepareVar | toLowerFirst }}: {{ setAttrIdx $v $const }},
            {{ end }}{{ end }}
        },
        {{ if ne $.MyStruct "" }} }, {{ end }}
        {{ end }}
    })
}
`

/*
@todo eav.NewAttribute() we may run into trouble regarding the args to eav.NewAttribute() because $row is a map. fix that
*/
