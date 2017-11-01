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
	"sort"
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
	DisableFileHeader bool
	// tpl contains a parsed template to render a single table.
	tpl *template.Template
}

// Option represents a sortable option for the NewTables function. Each option
// function can be applied in a mixed order.
type Option struct {
	sortOrder int
	fn        func(*Tables) error
}

// optionSorter to satisfy the sort.Slice function
type optionSorter []Option

func WithStructTags(tableName string, columnNameStructTag ...string) (opt Option) {
	if len(columnNameStructTag)%2 == 1 {
		panic(errors.NewFatalf("[dmlgen] WithStructTags: Argument columnNameStructTag must be a balanced slice."))
	}
	opt.sortOrder = 100
	opt.fn = func(ts *Tables) (err error) {
		var found int
		for _, t := range ts.Tables {
			if t.Name == tableName {
				for _, c := range t.Columns {
					for i := 0; i < len(columnNameStructTag); i = i + 2 {
						if c.Field == columnNameStructTag[i] {
							c.StructTag = columnNameStructTag[i+1]
							found++
						}
					}
				}
			}
		}
		if found != len(columnNameStructTag)/2 {
			err = errors.NewNotFoundf("[dmlgen] WithStructTags For table %q one column in %v cannot be found.", tableName, columnNameStructTag)
		}
		return err
	}
	return
}

func WithColumnAliases(tableName, columnName string, aliases ...string) (opt Option) {
	opt.sortOrder = 110
	opt.fn = func(ts *Tables) (err error) {
		found := false
		for _, t := range ts.Tables {
			if t.Name == tableName {
				for _, c := range t.Columns {
					if c.Field == columnName {
						c.Aliases = aliases
						found = true
					}
				}
			}
		}
		if !found {
			err = errors.NewNotFoundf("[dmlgen] WithColumnAliases: For table %q the column %q has not been found.", tableName, columnName)
		}
		return err
	}
	return
}

func WithUniquifiedColumns(tableName string, columnNames ...string) (opt Option) {
	opt.sortOrder = 120
	opt.fn = func(ts *Tables) (err error) {
		var found int
		// yay three for loops! But doesn't matter in this case as we're not
		// in a performance critical code.
		for _, t := range ts.Tables {
			if t.Name == tableName {
				for _, c := range t.Columns {
					for _, cn := range columnNames {
						if c.Field == cn {
							c.Uniquified = true
							found++
						}
					}
				}
			}
		}
		if len(columnNames) != found {
			err = errors.NewNotFoundf("[dmlgen] WithUniquifiedColumns: For table %q one column of %v cannot been found.", tableName, columnNames)
		}
		return err
	}
	return
}

func WithTable(tableName string, columns ddl.Columns) (opt Option) {
	opt.sortOrder = 10
	opt.fn = func(ts *Tables) error {
		ts.Tables = append(ts.Tables, &table{
			Name:    tableName,
			Columns: columns,
		})
		return nil
	}
	return
}

func WithLoadColumns(ctx context.Context, db dml.Querier, tables ...string) (opt Option) {
	opt.sortOrder = 1
	opt.fn = func(ts *Tables) error {
		tables, err := ddl.LoadColumns(ctx, db, tables...)
		if err != nil {
			return errors.WithStack(err)
		}
		// fight with the randomized map elements to retain the sortOrder from the
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
	return
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
	ts.tpl, err = template.New("entity").Funcs(ts.FuncMap).Parse(TplDBAC)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sOpts := optionSorter(opts)
	sort.Slice(sOpts, func(i, j int) bool {
		return sOpts[i].sortOrder < sOpts[j].sortOrder // ascending a-z sorting ;-)
	})

	for _, opt := range sOpts {
		if err := opt.fn(ts); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	for _, t := range ts.Tables {
		t.Package = ts.Package
	}
	return ts, nil
}

// findUsedPackages poor mans Go file parsing by just checking if bytes.Contains
// report true. should be rewritten to a real go code parser because we ignore
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
		if err := t.writeTo(buf, ts.tpl.Funcs(ts.FuncMap)); err != nil {
			return 0, errors.NewWriteFailed(err, "[dmlgen] For Table %q", t.Name)
		}
	}

	if !ts.DisableFileHeader {
		// now figure out all used package names in the buffer.
		fmt.Fprintf(w, "// Auto generated by dmlgen\n\npackage %s\n\nimport (\n", ts.Package)
		for _, path := range ts.findUsedPackages(buf.Bytes()) {
			fmt.Fprintf(w, "\t%q\n", path)
		}
		fmt.Fprint(w, "\n)\n")
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
}

// WriteTo implements io.WriterTo and writes the generated source code into w.
func (t *table) writeTo(w io.Writer, tpl *template.Template) error {

	data := struct {
		Package    string
		Collection string
		Entity     string
		TableName  string
		Columns    ddl.Columns
		Tick       string
	}{
		Package:    t.Package,
		Collection: strs.ToGoCamelCase(t.Name) + "Collection",
		Entity:     strs.ToGoCamelCase(t.Name),
		TableName:  t.Name,
		Columns:    t.Columns,
		Tick:       "`",
	}

	return tpl.Execute(w, data)
}
