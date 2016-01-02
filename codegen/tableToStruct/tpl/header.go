// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package tpl

const Header = `
// Auto generated via tableToStruct

import (
	"sort"
	"time"
    {{ if .HasTypeCodeValueTables }}
	"github.com/corestoreio/csfw/eav"{{end}}
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/storage/money"
)

var (
	_ = (*sort.IntSlice)(nil)
	_ = (*time.Time)(nil)
	_ = (*money.Money)(nil)
)

// TableIndex... is the index to a table. These constants are guaranteed
// to stay the same for all Magento versions. Please access a table via this
// constant instead of the raw table name. TableIndex iotas must start with 0.
const (
    {{ range $k,$v := .Tables }}TableIndex{{$v.Name}} {{ if eq $k 0 }}csdb.Index = iota{{ end }} // Table: {{$v.TableName}}
{{ end }}	TableIndexZZZ  // the maximum index, which is not available.
)

func init(){
    TableCollection = csdb.NewTableManager(
    {{ range $k,$v := .Tables }} csdb.AddTableByName(TableIndex{{.Name}}, "{{.TableName}}"),
    {{ end }} )
    // Don't forget to call TableCollection.ReInit(...) in your code to load the column definitions.
}`
