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
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io"
	"path/filepath"
	"text/template"

	"github.com/corestoreio/csfw/sql/ddl"
	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/slices"
	"github.com/corestoreio/csfw/util/strs"
	"github.com/corestoreio/errors"
)

type Tables struct {
	Package     string // Name of the package
	ImportPaths []string
	Tables      []*table
	template.FuncMap
	// ColumnAliases holds for a given key, the column name, its multiple aliases.
	// For example customer_entity.entity_id can also be sales_order.customer_id.
	// The alias would be just: entity_id:[]string{"customer_id"}.
	// tableName[columnName][]Aliases
	ColumnAliases map[string]map[string][]string
	// UniquifiedColumns defines a list of column names for which a
	// dedicated function gets generated to extract all values from the
	// collection into its own primitive slice. Only non blob/text columns can
	// be generated.
	UniquifiedColumns map[string][]string // tableName->column names
	// UniquifiedColumnMaxLength defines the maximal length a column can have to generate
	// its own extractor function. You must list the column in
	// UniquifiedColumns. Default value 256.
	UniquifiedColumnMaxLength int64 // columns longer than 255 characters won't have a dedicated method receiver
	DisableFileHeader         bool
	// tpl contains a parsed template to render a single table.
	tpl *template.Template
}

type Option func(*Tables) error

func WithColumnAliases(tableName, columnName string, aliases ...string) Option {
	return func(ts *Tables) error {
		if ts.ColumnAliases == nil {
			ts.ColumnAliases = make(map[string]map[string][]string)
		}
		if ts.ColumnAliases[tableName] == nil {
			ts.ColumnAliases[tableName] = make(map[string][]string)
		}
		ts.ColumnAliases[tableName][columnName] = aliases
		return nil
	}
}

func WithUniquifiedColumns(tableName string, columnNames ...string) Option {
	return func(ts *Tables) error {
		if ts.UniquifiedColumns == nil {
			ts.UniquifiedColumns = make(map[string][]string)
		}
		ts.UniquifiedColumns[tableName] = columnNames
		return nil
	}
}

func WithTable(tableName string, columns ddl.Columns) Option {
	return func(ts *Tables) error {
		ts.Tables = append(ts.Tables, &table{
			Name:    tableName,
			Columns: columns,
		})
		return nil
	}
}

func WithLoadColumns(ctx context.Context, db dml.Querier, tables ...string) Option {
	return func(ts *Tables) error {
		tables, err := ddl.LoadColumns(ctx, db, tables...)
		if err != nil {
			return errors.WithStack(err)
		}
		// fight with the randomized map elements to retain the order from the
		// SQL query for column loading.
		sortedKeys := make(slices.String, 0, len(tables))
		for k := range tables {
			sortedKeys = append(sortedKeys, k)
		}
		sortedKeys.Sort()
		for _, tblName := range sortedKeys {
			ts.Tables = append(ts.Tables, &table{
				Name:    tblName,
				Columns: tables[tblName],
			})
		}
		return nil
	}
}

func NewTables(packageName string, opts ...Option) (ts *Tables, err error) {
	ts = &Tables{
		Package: packageName,
		ImportPaths: []string{
			"database/sql",
			"github.com/corestoreio/csfw/sql/dml",
			"github.com/corestoreio/errors",
			"time",
		},
		FuncMap: make(template.FuncMap, 10),
	}
	ts.FuncMap["ToGoCamelCase"] = strs.ToGoCamelCase // net_http->NetHTTP entity_id->EntityID
	ts.FuncMap["GoTypeNull"] = toGoTypeNull
	ts.FuncMap["GoType"] = toGoType
	ts.FuncMap["GoFuncNull"] = toGoFuncNull
	ts.FuncMap["GoFunc"] = toGoFunc
	ts.FuncMap["GoPrimitive"] = toGoPrimitive
	ts.FuncMap["ColumnAliases"] = func(columnName string) []string {
		return []string{"PLACEHOLDER"}
	}
	ts.tpl, err = template.New("entity").Funcs(ts.FuncMap).Parse(TplDBAC)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, opt := range opts {
		if err := opt(ts); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	var charMaxLength int64 = 256
	if ts.UniquifiedColumnMaxLength > 0 {
		charMaxLength = ts.UniquifiedColumnMaxLength
	}
	for _, t := range ts.Tables {
		t.ColumnAliases = ts.ColumnAliases[t.Name]
		t.Package = ts.Package
		t.FilterUniquifiedColumns = func(c *ddl.Column) bool {
			// if slice UniquifiedColumns contains the column name X and is not
			// longer than 256 chars or is not a text/blob field then allowed to
			// generated the method receiver.
			return slices.String(ts.UniquifiedColumns[t.Name]).Contains(c.Field) && (!c.CharMaxLength.Valid ||
				(c.CharMaxLength.Valid && c.CharMaxLength.Int64 < charMaxLength))
		}
	}
	return ts, nil
}

// findUsedPackages poor mans Go file parsing by just checking if bytes.Contains
// report true. should be rewritten to real go code parser because we ignore
// here comments in checking, so false positive might get returned.
func (ts *Tables) findUsedPackages(file []byte) []string {
	ret := make([]string, 0, len(ts.ImportPaths))
	for _, path := range ts.ImportPaths {
		_, pkg := filepath.Split(path)
		if bytes.Contains(file, append([]byte(pkg), '.')) {
			ret = append(ret, path)
		}
	}
	return ret
}

// WriteTo executes the template parser and writes the result into w.
func (ts *Tables) WriteTo(w io.Writer) (int64, error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	for _, t := range ts.Tables {
		ts.FuncMap["ColumnAliases"] = func(columnName string) []string {
			return t.ColumnAliases[columnName]
		}
		if err := t.writeTo(buf, ts.tpl.Funcs(ts.FuncMap)); err != nil {
			return 0, errors.NewWriteFailed(err, "[dmlgen] For Table %q", t.Name)
		}
	}

	if !ts.DisableFileHeader {
		// now figure out all used package names in the buffer.
		fmt.Fprintf(w, "package %s\n\nimport (\n", ts.Package)
		for _, path := range ts.findUsedPackages(buf.Bytes()) {
			fmt.Fprintf(w, "\t%q\n", path)
		}
		fmt.Fprintf(w, "\n)\n")
	}

	fmted, err := format.Source(buf.Bytes())
	if err != nil {
		return 0, errors.WithStack(err)
	}
	buf.Reset()
	buf.Write(fmted)
	return buf.WriteTo(w)
}

// table writes one database table into Go source code.
type table struct {
	Package string      // Name of the package
	Name    string      // Name of the table
	Columns ddl.Columns // all columns of the table
	// ColumnAliases holds for a given key, the column name, its multiple aliases.
	// For example customer_entity.entity_id can also be sales_order.customer_id.
	// The alias would be just: entity_id:[]string{"customer_id"}.
	ColumnAliases           map[string][]string
	FilterUniquifiedColumns func(*ddl.Column) bool
}

// WriteTo implements io.WriterTo and writes the generated source code into w.
func (t *table) writeTo(w io.Writer, tpl *template.Template) error {

	data := struct {
		Package                  string
		Collection               string
		Entity                   string
		TableName                string
		Columns                  ddl.Columns
		Tick                     string
		ExtractColumns           ddl.Columns // contains a single PK and/or UNQ key
		ExtractUniquifiedColumns ddl.Columns // those columns have duplicate values and the dups get uniquified
	}{
		Package:    t.Package,
		Collection: strs.ToGoCamelCase(t.Name) + "Collection",
		Entity:     strs.ToGoCamelCase(t.Name),
		TableName:  t.Name,
		Columns:    t.Columns,
		Tick:       "`",
	}

	if pk := t.Columns.PrimaryKeys(); len(pk) == 1 {
		data.ExtractColumns = append(data.ExtractColumns, pk...)
	} else {
		data.ExtractUniquifiedColumns = append(data.ExtractUniquifiedColumns, pk...)
	}
	if uk := t.Columns.UniqueKeys(); len(uk) == 1 {
		data.ExtractColumns = append(data.ExtractColumns, uk...)
	} else {
		data.ExtractUniquifiedColumns = append(data.ExtractUniquifiedColumns, uk...)
	}

	// possibility of duplicate entries in this slice which gets filtered out in
	// the Filter call ;-)
	data.ExtractUniquifiedColumns = append(data.ExtractUniquifiedColumns, t.Columns.ColumnsNoPK()...)
	data.ExtractUniquifiedColumns = data.ExtractUniquifiedColumns.Filter(t.FilterUniquifiedColumns)

	return tpl.Execute(w, data)
}
