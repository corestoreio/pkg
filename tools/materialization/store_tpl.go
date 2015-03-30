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

const tplMaterializationStore = tools.Copyright + `
package {{ .PackageName }}

import (
	"database/sql"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
)

const (
    {{ range $k,$v := .Stores }} // Store{{prepareVarIndex $k $v.Code.String}} is the index to {{$v.Name}} ID: {{$v.StoreID}}
    Store{{prepareVarIndex $k $v.Code.String}} {{ if eq $k 0 }}store.StoreIndex = iota{{end}}
{{ end }} // Store_Max end of index, not available.
	Store999Max

    {{ range $k,$v := .Groups }} // Group{{prepareVarIndex $k $v.Name}} is the index to {{$v.Name}} ID: {{$v.GroupID}}
    Group{{prepareVarIndex $k $v.Name}} {{ if eq $k 0 }}store.GroupIndex = iota{{end}}
{{ end }} // Group_Max end of index, not available.
	Group999Max

    {{ range $k,$v := .Websites }} // Website{{prepareVarIndex $k $v.Code.String}} is the index to {{$v.Name.String}} ID: {{$v.WebsiteID}}
    Website{{prepareVarIndex $k $v.Code.String}} {{ if eq $k 0 }}store.WebsiteIndex = iota{{end}}
{{ end }} // Website_Max end of index, not available.
	Website999Max
)

type si struct {}

func (si) ByID(id int64) (store.StoreIndex, error){
	switch id {
	{{ range $k,$v := .Stores }} case {{$v.StoreID}}:
		return Store{{prepareVarIndex $k $v.Code.String}}, nil
	{{end}}
	default:
		return -1, store.ErrStoreNotFound
	}
}

func (si) ByCode(code string) (store.StoreIndex, error){
	switch code {
	{{ range $k,$v := .Stores }} case "{{$v.Code.String}}":
		return Store{{prepareVarIndex $k $v.Code.String}}, nil
	{{end}}
	default:
		return -1, store.ErrStoreNotFound
	}
}

func init(){
	store.SetStoreGetter(si{})
	store.SetStoreCollection(store.StoreSlice{
		{{ range $k,$v := .Stores }}Store{{prepareVarIndex $k $v.Code.String}}: {{ $v | printf "%#v" }},
		{{end}}
	})
}

type gi struct {}

func (gi) ByID(id int64) (store.GroupIndex, error){
	switch id {
	{{ range $k,$v := .Groups }} case {{$v.GroupID}}:
		return Group{{prepareVarIndex $k $v.Name}}, nil
	{{end}}
	default:
		return -1, store.ErrGroupNotFound
	}
}

func init(){
	store.SetGroupGetter(gi{})
	store.SetGroupCollection(store.GroupSlice{
		{{ range $k,$v := .Groups }}Group{{prepareVarIndex $k $v.Name}}: {{ $v | printf "%#v" }},
		{{end}}
	})
}

type wi struct {}

func (wi) ByID(id int64) (store.WebsiteIndex, error){
	switch id {
	{{ range $k,$v := .Websites }} case {{$v.WebsiteID}}:
		return Website{{prepareVarIndex $k $v.Code.String}}, nil
	{{end}}
	default:
		return -1, store.ErrWebsiteNotFound
	}
}

func (wi) ByCode(code string) (store.WebsiteIndex, error){
	switch code {
	{{ range $k,$v := .Websites }} case "{{$v.Code.String}}":
		return Website{{prepareVarIndex $k $v.Code.String}}, nil
	{{end}}
	default:
		return -1, store.ErrWebsiteNotFound
	}
}

func init(){
	store.SetWebsiteGetter(wi{})
	store.SetWebsiteCollection(store.WebsiteSlice{
		{{ range $k,$v := .Websites }}Website{{prepareVarIndex $k $v.Code.String}}: {{ $v | printf "%#v" }},
		{{end}}
	})
}
`
