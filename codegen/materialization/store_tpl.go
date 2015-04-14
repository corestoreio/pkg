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
    {{ range $k,$v := .Stores }}s{{prepareVarIndex $k $v.Code.String}} {{ if eq $k 0 }}store.IDX = iota{{end}}
{{ end }} // Store_Max end of index, not available.
	sIndexZZZ
)
const (
    {{ range $k,$v := .Groups }}g{{prepareVarIndex $k $v.Name}} {{ if eq $k 0 }}store.IDX = iota{{end}}
{{ end }} // Group_Max end of index, not available.
	gIndexZZZ
)
const (
    {{ range $k,$v := .Websites }}w{{prepareVarIndex $k $v.Code.String}} {{ if eq $k 0 }}store.IDX = iota{{end}}
{{ end }} // Website_Max end of index, not available.
	wIndexZZZ
)

var	storeManager = store.NewStoreManager(
		store.NewStoreBucket(
			store.TableStoreSlice{
				{{ range $k,$v := .Stores }}s{{prepareVarIndex $k $v.Code.String}}: {{ $v | printf "%#v" }},
				{{end}}
			},
			store.StoreIndexIDMap{
				{{ range $k,$v := .Stores }} {{$v.StoreID}}: s{{prepareVarIndex $k $v.Code.String}},
			{{end}}},
			store.StoreIndexCodeMap{
				{{ range $k,$v := .Stores }} "{{$v.Code.String}}": s{{prepareVarIndex $k $v.Code.String}},
			{{end}}},
		),
		store.NewGroupBucket(
			store.TableGroupSlice{
				{{ range $k,$v := .Groups }}g{{prepareVarIndex $k $v.Name}}: {{ $v | printf "%#v" }},
				{{end}}
			},
			store.GroupIndexIDMap{
				{{ range $k,$v := .Groups }} {{$v.GroupID}}: g{{prepareVarIndex $k $v.Name}},
			{{end}}},
		),
		store.NewWebsiteBucket(
			store.TableWebsiteSlice{
				{{ range $k,$v := .Websites }}w{{prepareVarIndex $k $v.Code.String}}: {{ $v | printf "%#v" }},
				{{end}}
			},
			store.WebsiteIndexIDMap{
				{{ range $k,$v := .Websites }} {{$v.WebsiteID}}: w{{prepareVarIndex $k $v.Code.String}},
			{{end}}},
			store.WebsiteIndexCodeMap{
				{{ range $k,$v := .Websites }} "{{$v.Code.String}}": w{{prepareVarIndex $k $v.Code.String}},
			{{end}}},
		),
	)
`
