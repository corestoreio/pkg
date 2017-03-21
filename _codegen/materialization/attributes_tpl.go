// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

/* @todo
   Data will be "carved in stone" because it only changes during development.
   - attribute_set related tables: eav_attribute_set, eav_entity_attribute, eav_attribute_group, etc
   - label and option tables will not be hard coded
*/
const tplAttrImport = `
package {{ .PackageName }}
    import (
        "github.com/corestoreio/csfw/eav"
        "{{ .AttrPkgImp }}"
        {{ range .ImportPaths }}"{{ . }}"
        {{ end }} )
`

const tplAttrTypes = `
const (
    {{ range $k, $row := .AttrCol }}{{ $.Name | prepareVar }}{{ index $row "attribute_code" | prepareVar }} {{ if eq $k 0 }} eav.AttributeIndex = iota + 1{{ end }}
    {{ end }}
    {{ $.Name | prepareVar }}ZZZ
)

type (
    // {{ .Name | prepareVar }} a data container for attributes. You can use this struct to
    // embed into your own struct for maybe overriding some method receivers.
    {{ .Name | prepareVar | toLowerFirst }} struct {
        // this type contains only custom columns ...
        *{{ .AttrPkg }}.{{ .AttrStruct }}
        // add only those columns which are not in variables AttributeCoreColumns for each attr package
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

const tplAttrGetter = `
func init(){
    {{ .AttrPkg }}.{{ .FuncGetter }}(eav.NewAttributeMapGet(
        map[int64]eav.AttributeIndex{
            {{ range $k, $row := .AttrCol }} {{ index $row "attribute_id" }}: {{ $.Name | prepareVar }}{{ index $row "attribute_code" | prepareVar }},
            {{ end }} },
        map[string]eav.AttributeIndex{
        {{ range $k, $row := .AttrCol }} {{ index $row "attribute_code" }}: {{ $.Name | prepareVar }}{{ index $row "attribute_code" | prepareVar }},
        {{ end }} },
    ))
}
`

const tplAttrCollection = `
func init(){
    {{ .AttrPkg }}.{{ .FuncCollection }}({{ .AttrPkg }}.AttributeSlice{
    {{ range $row := .AttrCol }}
        {{ $const := sprintf "%s%s" (prepareVar $.Name) (prepareVar (index $row "attribute_code")) }}
        {{ $const }}: {{ if ne $.MyStruct "" }} &{{ $.MyStruct }} {
        {{ end }} &{{ $.Name | prepareVar | toLowerFirst }} {
            {{ $.AttrStruct }}: {{ .WebsiteEntityAttribute }},
            {{ range $k,$v := $row }} {{ if (isUnknownAttr $k) }} {{ setAttrIdx $v $const }}, // {{ $k }}
            {{ end }}{{ end }}
        },
        {{ if ne $.MyStruct "" }} }, {{ end }}
    {{ end }}
    })
}
`

const tplAttrWebsiteEavAttribute = `
    eav.NewAttribute(
        {{ .websiteAttribute }},
        {{ .website_id }},
        {{ range $k,$v := $row }} {{ if (isEavAttr $k) }} {{ setAttrIdx $v $const }}, // {{ $k }}
        {{ end }}{{ end }}
    )
`
const tplAttrWebsiteEntityAttribute = `
    {{ $.AttrPkg }}.New{{ $.AttrStruct }}(
        {{ .WebsiteEavAttribute }}
        {{ .websiteAttribute }},
        {{ range $k,$v := $row }} {{ if (isEavEntityAttr $k) }} {{ setAttrIdx $v $const }}, // {{ $k }}
        {{ end }}{{ end }}
    )
`

/*
@todo eav.NewAttribute() we may run into trouble regarding the args to eav.NewAttribute() because $row is a map. fix that
*/
