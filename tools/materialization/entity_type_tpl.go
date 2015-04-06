// Copyright 2015 CoreStore Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	 http://www.apache.org/licenses/LICENSE-2.0
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
   - DONE: entity_type with translation of some columns to the Go type
   - eav_attribute full config and from that one the flat table structure
*/

const tplEav = tools.Copyright + `
package {{ .Package }}

// Package {{ .Package }} file is auto generated

import (
	"github.com/corestoreio/csfw/eav"{{ range .ImportPaths }}
	"{{ . }}"{{end}}
)

func init(){
	eav.SetEntityTypeCollection(eav.CSEntityTypeSlice{
		{{ range .ETypeData }} &eav.CSEntityType {
			EntityTypeID: {{ .EntityTypeID }},
			EntityTypeCode: "{{ .EntityTypeCode }}",
			EntityModel: {{ extractFuncType .EntityModel }},
			AttributeModel: {{ extractFuncType .AttributeModel.String }},
			EntityTable: {{ extractFuncType .EntityTable.String }},
			ValueTablePrefix: "{{ .ValueTablePrefix.String }}",
			IsDataSharing: {{ .IsDataSharing }},
			DataSharingKey: "{{ .DataSharingKey.String }}",
			DefaultAttributeSetID: {{ .DefaultAttributeSetID }},
			{{ if ne "" .IncrementModel.String }}IncrementModel: {{ extractFuncType .IncrementModel.String }},{{ end }}
			IncrementPerStore: {{ .IncrementPerStore }},
			IncrementPadLength: {{ .IncrementPadLength }},
			IncrementPadChar: "{{ .IncrementPadChar }}",
			AdditionalAttributeTable: {{ extractFuncType .AdditionalAttributeTable.String }},
			EntityAttributeCollection: {{ extractFuncType .EntityAttributeCollection.String }},
		},
		{{ end }}
	})
}
`
