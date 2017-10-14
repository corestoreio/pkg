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
		"github.com/corestoreio/csfw/storage/money",
		"github.com/corestoreio/errors",
		"time",
	},
}

// Table writes one database table into Go source code.
type Table struct {
	Package string
	Name    string
	Columns ddl.Columns
	// ColumnAliases holds for a given key, the column name, its multiple aliases.
	// For example customer_entity.entity_id can also be sales_order.customer_id.
	// The alias would be just: entity_id:[]string{"customer_id"}.
	ColumnAliases map[string][]string
	template.FuncMap
	CharMaxLength int64 // default 256; all columns longer than 256 won't have a single function ...
}

// WriteTo implements io.WriterTo and writes the generated source code into w.
func (t *Table) WriteTo(w io.Writer) (n int64, err error) {

	if t.FuncMap == nil {
		t.FuncMap = make(template.FuncMap)
	}
	t.FuncMap["ToCamelCase"] = strs.ToCamelCase
	t.FuncMap["ToGoCamelCase"] = strs.ToGoCamelCase
	if _, ok := t.FuncMap["MySQLToGoType"]; !ok {
		t.FuncMap["MySQLToGoType"] = MySQLToGoType
	}
	if _, ok := t.FuncMap["GoTypeFuncName"]; !ok {
		t.FuncMap["GoTypeFuncName"] = GoTypeFuncName
	}
	if _, ok := t.FuncMap["ColumnAliases"]; !ok {
		t.FuncMap["ColumnAliases"] = func(columnName string) []string {
			return t.ColumnAliases[columnName]
		}
	}

	tplEntity, err := template.New("entity").Funcs(t.FuncMap).Parse(TplDBAC)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	// TODO: to generate the  entity slice extractor function ESEF figure out
	// which columns are PK and UNQ. But only use those which have a single
	// column as PK or UNQ. All other non unique columns must have a possibility
	// to remove duplicate values.

	data := struct {
		Package               string
		Collection            string
		Entity                string
		TableName             string
		Columns               ddl.Columns
		Tick                  string
		SingleKeyColumns      ddl.Columns // contains a single PK and/or UNQ key
		DuplicateValueColumns ddl.Columns // those columns have duplicate values
	}{
		Package:    t.Package,
		Collection: strs.ToGoCamelCase(t.Name) + "Collection",
		Entity:     strs.ToGoCamelCase(t.Name),
		TableName:  t.Name,
		Columns:    t.Columns,
		Tick:       "`",
	}

	if sk := t.Columns.PrimaryKeys(); len(sk) == 1 {
		data.SingleKeyColumns = append(data.SingleKeyColumns, sk...)
	} else {
		data.DuplicateValueColumns = append(data.DuplicateValueColumns, sk...)
	}
	if sk := t.Columns.UniqueKeys(); len(sk) == 1 {
		data.SingleKeyColumns = append(data.SingleKeyColumns, sk...)
	} else {
		data.DuplicateValueColumns = append(data.DuplicateValueColumns, sk...)
	}

	// possibility of duplicate entries in this slice.
	data.DuplicateValueColumns = append(data.DuplicateValueColumns, t.Columns.ColumnsNoPK()...)

	var charMaxLength int64 = 256
	if t.CharMaxLength > 0 {
		charMaxLength = t.CharMaxLength
	}
	data.DuplicateValueColumns = data.DuplicateValueColumns.Filter(func(c *ddl.Column) bool {
		return !c.CharMaxLength.Valid || (c.CharMaxLength.Valid && c.CharMaxLength.Int64 < charMaxLength)
	})

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if err := tplEntity.Execute(buf, data); err != nil {
		return 0, errors.WithStack(err)
	}

	fmted, err := format.Source(buf.Bytes())
	if err != nil {
		return 0, errors.WithStack(err)
	}
	buf.Reset()
	buf.Write(fmted)
	return buf.WriteTo(w)
}
