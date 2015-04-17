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

var	storeManager = store.NewStoreManager(
		store.NewStoreBucket(
			store.TableStoreSlice{
				{{ range $k,$v := .Stores }}{{ $v | printf "%#v" }},
				{{end}}
			},
		),
		store.NewGroupBucket(
			store.TableGroupSlice{
				{{ range $k,$v := .Groups }}{{ $v | printf "%#v" }},
				{{end}}
			},
		),
		store.NewWebsiteBucket(
			store.TableWebsiteSlice{
				{{ range $k,$v := .Websites }}{{ $v | printf "%#v" }},
				{{end}}
			},
		),
	)
`
