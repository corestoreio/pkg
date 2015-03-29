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

const (
    {{ range $k,$v := .Stores }} // Store{{prepareVarIndex $k $v.Code.String}} is the index to {{$v.Name}} ID: {{$v.StoreID}}
    Store{{prepareVarIndex $k $v.Code.String}} {{ if eq $k 0 }}StoreIndex = iota{{end}}
{{ end }} // Store_Max end of index, not available.
	Store999Max

    {{ range $k,$v := .Groups }} // Group{{prepareVarIndex $k $v.Name}} is the index to {{$v.Name}} ID: {{$v.GroupID}}
    Group{{prepareVarIndex $k $v.Name}} {{ if eq $k 0 }}GroupIndex = iota{{end}}
{{ end }} // Group_Max end of index, not available.
	Group999Max

    {{ range $k,$v := .Websites }} // Website{{prepareVarIndex $k $v.Code.String}} is the index to {{$v.Name.String}} ID: {{$v.WebsiteID}}
    Website{{prepareVarIndex $k $v.Code.String}} {{ if eq $k 0 }}WebsiteIndex = iota{{end}}
{{ end }} // Website_Max end of index, not available.
	Website999Max
)

func init(){
	storeCollection = StoreSlice{
		{{ range $k,$v := .Stores }}Store{{prepareVarIndex $k $v.Code.String}}: {{ $v | printf "%#v" }},
		{{end}}
	}
	groupCollection = GroupSlice{
		{{ range $k,$v := .Groups }}Group{{prepareVarIndex $k $v.Name}}: {{ $v | printf "%#v" }},
		{{end}}
	}
	websiteCollection = WebsiteSlice{
		{{ range $k,$v := .Websites }}Website{{prepareVarIndex $k $v.Code.String}}: {{ $v | printf "%#v" }},
		{{end}}
	}
}

func GetStoreByID(id int64) (*Store, error) {
	switch id {
	{{ range $k,$v := .Stores }} case {{$v.StoreID}}:
		return storeCollection[Store{{prepareVarIndex $k $v.Code.String}}], nil
	{{end}}
	default:
		return nil, ErrStoreNotFound
	}
}

func GetStoreByCode(code string) (*Store, error) {
	switch code {
	{{ range $k,$v := .Stores }} case "{{$v.Code.String}}":
		return storeCollection[Store{{prepareVarIndex $k $v.Code.String}}], nil
	{{end}}
	default:
		return nil, ErrStoreNotFound
	}
}

// @todo add GetGroup/s and GetWebsite/s

`
