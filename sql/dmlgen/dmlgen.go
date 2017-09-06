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

package dmlgen

import (
	"go/format"
	"io"
	"text/template"

	"github.com/corestoreio/csfw/sql/ddl"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/strs"
	"github.com/corestoreio/errors"
)

var Imports = map[string][]string{
	"table": {
		"database/sql",
		"github.com/corestoreio/csfw/sql/dml",
		"github.com/corestoreio/errors",
	},
}

// Table writes one database table into Go source code.
type Table struct {
	Package string
	Name    string
	Columns ddl.Columns
	template.FuncMap
}

// WriteTo implements io.WriterTo and writes the generated source code into w.
func (t *Table) WriteTo(w io.Writer) (n int64, err error) {

	if t.FuncMap == nil {
		t.FuncMap = make(template.FuncMap)
	}
	t.FuncMap["ToCamelCase"] = strs.ToCamelCase
	t.FuncMap["ToGoCamelCase"] = strs.ToGoCamelCase
	t.FuncMap["FieldName"] = func(c *ddl.Column) string {
		fn := strs.ToGoCamelCase(c.Field)
		if !c.IsNull() {
			return fn
		}
		return fn + ".Null" + strs.ToGoCamelCase(c.DataTypeSimple())
	}

	tplEntity, err := template.New("entity").Funcs(t.FuncMap).Parse(TplEntity)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	data := struct {
		Package   string
		Entity    string
		TableName string
		Columns   ddl.Columns
		Tick      string
	}{
		Package:   t.Package,
		Entity:    strs.ToGoCamelCase(t.Name),
		TableName: t.Name,
		Columns:   t.Columns,
		Tick:      "`",
	}
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := tplEntity.Execute(buf, data); err != nil {
		return 0, errors.WithStack(err)
	}

	//tplCollection, err := template.New("entity").Parse(TplCollection)
	//if err != nil {
	//	return 0, errors.WithStack(err)
	//}

	fmted, err := format.Source(buf.Bytes())
	if err != nil {
		return 0, errors.WithStack(err)
	}
	buf.Reset()
	buf.Write(fmted)
	return buf.WriteTo(w)
}
