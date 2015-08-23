// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// @todo hide password and other sensitive fields in JSON struct tags

const tplCopyPkg = `// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package {{ .Package }}
`

const tplHeader = `
// Auto generated via tableToStruct

import (
	"time"
    {{ if .HasTypeCodeValueTables }}
	"github.com/corestoreio/csfw/eav"{{end}}
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
)

var _ = (*time.Time)(nil)

// TableIndex... is the index to a table. These constants are guaranteed
// to stay the same for all Magento versions. Please access a table via this
// constant instead of the raw table name. TableIndex iotas must start with 0.
const (
    {{ range $k,$v := .Tables }}TableIndex{{$v.Name}} {{ if eq $k 0 }}csdb.Index = iota{{ end }} // Table: {{$v.NameRaw}}
{{ end }}	TableIndexZZZ  // the maximum index, which is not available.
)

func init(){
    TableCollection = csdb.NewTableManager(
    {{ range $k,$v := .Tables }} csdb.AddTableByName(TableIndex{{.Name}}, "{{.NameRaw}}"),
    {{ end }} )
    // Don't forget to call TableCollection.ReInit(...) in your code to load the column definitions.
}`

const tplTable = `
// Table{{.Name}} and Table{{.Name}}Slice, a type for DB table {{ .NameRaw }}
type (
    Table{{.Name}}Slice []*Table{{.Name}}
    Table{{.Name}} struct {
        {{ range .Columns }}{{.GoName}} {{.GoType}} {{ $.Tick }}db:"{{.Field.String}}" json:",omitempty"{{ $.Tick }} {{.Comment}}
        {{ end }} }
)

// {{.MethodRecvPrefix}}Load fills this slice with data from the database
func (s *Table{{.Name}}Slice) {{.MethodRecvPrefix}}Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return csdb.LoadSlice(dbrSess, TableCollection, TableIndex{{.Name}}, &(*s), cbs...)
}

// Len returns the length
func (s Table{{.Name}}Slice) Len() int { return len(s) }
`

const tplEAValues = `
{{range $typeCode,$valueTables := .TypeCodeValueTables}}
// Get{{ $typeCode | prepareVar }}ValueStructure returns for an EAV index the table structure.
// Important also if you have custom value tables
func Get{{ $typeCode | prepareVar }}ValueStructure(i eav.ValueIndex) (*csdb.Table, error) {
	switch i {
	{{range $vt,$v := $valueTables }}case eav.EntityType{{ $v | prepareVar }}:
		return TableCollection.Structure(TableIndex{{ $vt | prepareVar }})
    {{end}}	}
	return nil, eav.ErrEntityTypeValueNotFound
}
{{end}}
`

//// find by primary key
//{{ if .columns.PrimaryKeys.Length eq 1 }}
//{{ $pkColumn := .columns.PrimaryKeys.First }}
//// FindByID returns a Table{{.Name}} if found by id or an error
//func (s Table{{.Name}}Slice) FindBy{{ $pkColumn | prepareVar }}(id int64) (*Table{{.Name}}, error) {
//for _, u := range s {
//{{/* @todo check if column PK is really an int64 column */}}
//if u != nil && u.{{ $pkColumn | prepareVar }} == id {
//return u, nil
//}
//}
//return nil, errgo.Newf("ID %d in slice Table{{.Name}}Slice not found", id)
//}
//{{ end }}

//// FindByUsername returns a Table{{.Name}} if found by code or an error
//func (s Table{{.Name}}Slice) FindByUsername(username string) (*Table{{.Name}}, error) {
//for _, u := range s {
//if u != nil && u.Username.Valid && u.Username.String == username {
//return u, nil
//}
//}
//return nil, ErrUserNotFound
//}
//
//// Filter returns a new slice filtered by predicate f
//func (s Table{{.Name}}Slice) Filter(f func(*Table{{.Name}}) bool) Table{{.Name}}Slice {
//var tws Table{{.Name}}Slice
//for _, w := range s {
//if w != nil && f(w) {
//tws = append(tws, w)
//}
//}
//return tws
//}
//
//// IDs returns an Int64Slice with all primary key ids
//func (s Table{{.Name}}Slice) IDs() utils.Int64Slice {
//id := make(utils.Int64Slice, len(s))
//for i, r := range s {
//id[i] = r.{{ .primary key }}
//}
//return id
//}
