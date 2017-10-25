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
	"github.com/corestoreio/csfw/util/slices"
	"github.com/corestoreio/csfw/util/strs"
	"github.com/corestoreio/errors"
)

var Imports = map[string][]string{
	"table": {
		"database/sql",
		"github.com/corestoreio/csfw/sql/dml",
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
	// CharMaxLength defines the maximal length a column can have to generate
	// its own extractor function. You must list the column in
	// AllowedDuplicateValueColumns. Default value 256.
	CharMaxLength int64
	// AllowedDuplicateValueColumns defines a list of column names for which a
	// dedicated function gets generated to extract all values from the
	// collection into its own primitive slice. Only non blob/text columns can
	// be generated.
	AllowedDuplicateValueColumns []string
}

func (t *Table) initFuncMap() {
	if t.FuncMap == nil {
		t.FuncMap = make(template.FuncMap, 10)
	}

	t.FuncMap["ToGoCamelCase"] = strs.ToGoCamelCase // net_http->NetHTTP entity_id->EntityID
	if _, ok := t.FuncMap["GoTypeNull"]; !ok {
		t.FuncMap["GoTypeNull"] = toGoTypeNull
	}
	if _, ok := t.FuncMap["GoType"]; !ok {
		t.FuncMap["GoType"] = toGoType
	}
	if _, ok := t.FuncMap["GoFuncNull"]; !ok {
		t.FuncMap["GoFuncNull"] = toGoFuncNull
	}
	if _, ok := t.FuncMap["GoFunc"]; !ok {
		t.FuncMap["GoFunc"] = toGoFunc
	}
	if _, ok := t.FuncMap["GoPrimitive"]; !ok {
		t.FuncMap["GoPrimitive"] = toGoPrimitive
	}
	if _, ok := t.FuncMap["ColumnAliases"]; !ok {
		t.FuncMap["ColumnAliases"] = func(columnName string) []string {
			return t.ColumnAliases[columnName]
		}
	}
}

// WriteTo implements io.WriterTo and writes the generated source code into w.
func (t *Table) WriteTo(w io.Writer) (int64, error) {
	t.initFuncMap()

	tplEntity, err := template.New("entity").Funcs(t.FuncMap).Parse(TplDBAC)
	if err != nil {
		return 0, errors.WithStack(err)
	}

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

	if pk := t.Columns.PrimaryKeys(); len(pk) == 1 {
		data.SingleKeyColumns = append(data.SingleKeyColumns, pk...)
	} else {
		data.DuplicateValueColumns = append(data.DuplicateValueColumns, pk...)
	}
	if uk := t.Columns.UniqueKeys(); len(uk) == 1 {
		data.SingleKeyColumns = append(data.SingleKeyColumns, uk...)
	} else {
		data.DuplicateValueColumns = append(data.DuplicateValueColumns, uk...)
	}

	// possibility of duplicate entries in this slice.
	data.DuplicateValueColumns = append(data.DuplicateValueColumns, t.Columns.ColumnsNoPK()...)

	var charMaxLength int64 = 256
	if t.CharMaxLength > 0 {
		charMaxLength = t.CharMaxLength
	}
	data.DuplicateValueColumns = data.DuplicateValueColumns.Filter(func(c *ddl.Column) bool {
		return slices.String(t.AllowedDuplicateValueColumns).Contains(c.Field) && (!c.CharMaxLength.Valid ||
			(c.CharMaxLength.Valid && c.CharMaxLength.Int64 < charMaxLength))
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
