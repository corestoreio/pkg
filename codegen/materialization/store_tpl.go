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

const tplMaterializationStore = `package {{ .PackageName }}

import (
	"database/sql"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
)

const (
    {{ range $k,$v := .Stores }} // Store{{prepareVarIndex $k $v.Code.String}} is the index to {{$v.Name}} ID: {{$v.StoreID}}
    Store{{prepareVarIndex $k $v.Code.String}} {{ if eq $k 0 }}store.StoreIndex = iota + 1{{end}}
{{ end }} // Store_Max end of index, not available.
	StoreIndexZZZ
)
const (
    {{ range $k,$v := .Groups }} // Group{{prepareVarIndex $k $v.Name}} is the index to {{$v.Name}} ID: {{$v.GroupID}}
    Group{{prepareVarIndex $k $v.Name}} {{ if eq $k 0 }}store.GroupIndex = iota + 1{{end}}
{{ end }} // Group_Max end of index, not available.
	GroupIndexZZZ
)
const (
    {{ range $k,$v := .Websites }} // Website{{prepareVarIndex $k $v.Code.String}} is the index to {{$v.Name.String}} ID: {{$v.WebsiteID}}
    Website{{prepareVarIndex $k $v.Code.String}} {{ if eq $k 0 }}store.WebsiteIndex = iota + 1{{end}}
{{ end }} // Website_Max end of index, not available.
	WebsiteIndexZZZ
)

var	storeManager = store.NewStoreManager(
		store.NewStoreBucket(
			store.TableStoreSlice{
				{{ range $k,$v := .Stores }}Store{{prepareVarIndex $k $v.Code.String}}: {{ $v | printf "%#v" }},
				{{end}}
			},
			store.StoreIndexIDMap{
				{{ range $k,$v := .Stores }} {{$v.StoreID}}: Store{{prepareVarIndex $k $v.Code.String}},
			{{end}}},
			store.StoreIndexCodeMap{
				{{ range $k,$v := .Stores }} "{{$v.Code.String}}": Store{{prepareVarIndex $k $v.Code.String}},
			{{end}}},
		),
		store.NewGroupBucket(
			store.TableGroupSlice{
				{{ range $k,$v := .Groups }}Group{{prepareVarIndex $k $v.Name}}: {{ $v | printf "%#v" }},
				{{end}}
			},
			store.GroupIndexIDMap{
				{{ range $k,$v := .Groups }} {{$v.GroupID}}: Group{{prepareVarIndex $k $v.Name}},
			{{end}}},
		),
		store.NewWebsiteBucket(
			store.TableWebsiteSlice{
				{{ range $k,$v := .Websites }}Website{{prepareVarIndex $k $v.Code.String}}: {{ $v | printf "%#v" }},
				{{end}}
			},
			store.WebsiteIndexIDMap{
				{{ range $k,$v := .Websites }} {{$v.WebsiteID}}: Website{{prepareVarIndex $k $v.Code.String}},
			{{end}}},
			store.WebsiteIndexCodeMap{
				{{ range $k,$v := .Websites }} "{{$v.Code.String}}": Website{{prepareVarIndex $k $v.Code.String}},
			{{end}}},
		),
	)
`
