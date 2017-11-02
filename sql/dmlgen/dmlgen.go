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
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
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

// Tables can generated Go source for for database tables once correctly
// configured.
type Tables struct {
	Package     string // Name of the package
	ImportPaths []string
	// Tables uses the table name as map key and the table description as value.
	Tables map[string]*table
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

// WithStructTags enables struct tags proactively for the whole struct. Allowed
// values are: bson, db, env, json, toml, yaml and xml. For bson, json, yaml and
// xml the omitempty attribute has been set. If you need a different struct tag
// for a specifiv column you must set the option WithCustomStructTags. It
// doesn't matter in which order you apply the options ;-)
func WithStructTags(tableName string, tagNames ...string) (opt Option) {
	opt.sortOrder = 90 // must run before custom struct tags
	opt.fn = func(ts *Tables) (err error) {
		if t, ok := ts.Tables[tableName]; ok {
			for _, c := range t.Columns {
				var buf bytes.Buffer // TODO use new string.Builder in Go 1.10 or 1.11
				for i, tagName := range tagNames {
					if i > 0 {
						buf.WriteByte(' ')
					}
					// Maybe some types in the struct for a table don't need at
					// all an omitempty so build in some logic which creates the
					// tags more thoughtfully.
					switch tagName {
					case "bson":
						fmt.Fprintf(&buf, `bson:"%s,omitempty"`, c.Field)
					case "db":
						fmt.Fprintf(&buf, `db:"%s"`, c.Field)
					case "env":
						fmt.Fprintf(&buf, `env:"%s"`, c.Field)
					case "json":
						fmt.Fprintf(&buf, `json:"%s,omitempty"`, c.Field)
					case "toml":
						fmt.Fprintf(&buf, `toml:"%s"`, c.Field)
					case "yaml":
						fmt.Fprintf(&buf, `yaml:"%s,omitempty"`, c.Field)
					case "xml":
						fmt.Fprintf(&buf, `xml:"%s,omitempty"`, c.Field)
					default:
						return errors.NewNotSupportedf("[dmlgen] WithStructTags: tag %q not supported", tagName)
					}
				}
				c.StructTag = buf.String()
			}
		} else {
			err = errors.NewNotFoundf("[dmlgen] WithStructTags: Table %q cannot be found.", tableName)
		}
		return err
	}
	return
}

// WithCustomStructTags allows to specify custom struct tags for a specific column.
// The argument `columnNameStructTag` must be a balanced slice where index i
// sets to the column name and index i+1 to the desired struct tag.
//		dmlgen.WithCustomStructTags("table_name","column_a",`json: ",omitempty"`,"column_b","`xml:,omitempty`")
// It doesn't matter in which order you apply the options ;-)
func WithCustomStructTags(tableName string, columnNameStructTag ...string) (opt Option) {
	// Maybe create a new function option called WithStructTag(tableName string, json,xml,yaml,protobuf bool)
	if len(columnNameStructTag)%2 == 1 {
		panic(errors.NewFatalf("[dmlgen] WithCustomStructTags: Argument columnNameStructTag must be a balanced slice."))
	}
	opt.sortOrder = 100
	opt.fn = func(ts *Tables) (err error) {
		if t, ok := ts.Tables[tableName]; ok {
			var found int
			for _, c := range t.Columns {
				for i := 0; i < len(columnNameStructTag); i = i + 2 {
					if c.Field == columnNameStructTag[i] {
						c.StructTag = columnNameStructTag[i+1]
						found++
					}
				}
			}
			if found != len(columnNameStructTag)/2 {
				err = errors.NewNotFoundf("[dmlgen] WithCustomStructTags: For table %q one column in %v cannot be found.", tableName, columnNameStructTag)
			}
		} else {
			err = errors.NewNotFoundf("[dmlgen] WithCustomStructTags: Table %q cannot be found.", tableName)
		}
		return err
	}
	return
}

// WithColumnAliases specifies different names used for a column. For example
// customer_entity.entity_id can also be sales_order.customer_id, hence a
// Foreign Key. The alias would be just: entity_id:[]string{"customer_id"}.
func WithColumnAliases(tableName, columnName string, aliases ...string) (opt Option) {
	opt.sortOrder = 110
	opt.fn = func(ts *Tables) (err error) {
		if t, ok := ts.Tables[tableName]; ok {
			found := false
			for _, c := range t.Columns {
				if c.Field == columnName {
					c.Aliases = aliases
					found = true
				}
			}
			if !found {
				err = errors.NewNotFoundf("[dmlgen] WithColumnAliases: For table %q the column %q has not been found.", tableName, columnName)
			}
		} else {
			err = errors.NewNotFoundf("[dmlgen] WithColumnAliases: Table %q cannot be found.", tableName)
		}
		return err
	}
	return
}

// WithUniquifiedColumns specifies columns which are non primary/unique key one
// but should have a dedicated function to extract their unique primitive values
// as a slice.
func WithUniquifiedColumns(tableName string, columnNames ...string) (opt Option) {
	opt.sortOrder = 120
	opt.fn = func(ts *Tables) (err error) {
		if t, ok := ts.Tables[tableName]; ok {
			var found int
			for _, c := range t.Columns {
				for _, cn := range columnNames {
					if c.Field == cn {
						c.Uniquified = true
						found++
					}
				}
			}
			if len(columnNames) != found {
				err = errors.NewNotFoundf("[dmlgen] WithUniquifiedColumns: For table %q one column of %v cannot been found.", tableName, columnNames)
			}
		} else {
			err = errors.NewNotFoundf("[dmlgen] WithColumnAliases: Table %q cannot be found.", tableName)
		}
		return err
	}
	return
}

// WithTable sets a table and its columns. Allows to overwrite a table fetched
// with function WithLoadColumns.
func WithTable(tableName string, columns ddl.Columns) (opt Option) {
	opt.sortOrder = 10
	opt.fn = func(ts *Tables) error {
		ts.Tables[tableName] = &table{
			Name:    tableName,
			Columns: columns,
		}
		return nil
	}
	return
}

// WithLoadColumns queries the information_schema table and loads the column
// definition of the provided `tables` slice.
func WithLoadColumns(ctx context.Context, db dml.Querier, tables ...string) (opt Option) {
	opt.sortOrder = 1
	opt.fn = func(ts *Tables) error {
		tables, err := ddl.LoadColumns(ctx, db, tables...)
		if err != nil {
			return errors.WithStack(err)
		}
		for tblName := range tables {
			ts.Tables[tblName] = &table{
				Name:    tblName,
				Columns: tables[tblName],
			}
		}
		return nil
	}
	return
}

func (ts *Tables) sortedTableNames() []string {
	sortedKeys := make(slices.String, 0, len(ts.Tables))
	for k := range ts.Tables {
		sortedKeys = append(sortedKeys, k)
	}
	sortedKeys.Sort()
	return sortedKeys
}

// NewTables creates a new instance of the SQL table code generator.
func NewTables(packageName string, opts ...Option) (ts *Tables, err error) {
	ts = &Tables{
		Tables:  make(map[string]*table),
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
	ts.tpl, err = template.New("entity").Funcs(ts.FuncMap).Parse(tplDBAC)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sOpts := optionSorter(opts)
	sort.Slice(sOpts, func(i, j int) bool {
		return sOpts[i].sortOrder < sOpts[j].sortOrder // ascending 0-9 sorting ;-)
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

// findUsedPackages checks for needed packages which we must import.
func (ts *Tables) findUsedPackages(file []byte) ([]string, error) {

	af, err := parser.ParseFile(token.NewFileSet(), "virtual_file.go", append([]byte("package temporarily_main\n\n"), file...), 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	idents := map[string]struct{}{}
	ast.Inspect(af, func(n ast.Node) bool {
		switch nt := n.(type) {
		case *ast.Ident:
			idents[nt.Name] = struct{}{} // will contain too much info
			// we only need to know: pkg.TYPE
		}
		return true
	})

	ret := make([]string, 0, len(ts.ImportPaths))
	for _, path := range ts.ImportPaths {
		_, pkg := filepath.Split(path)
		if _, ok := idents[pkg]; ok {
			ret = append(ret, path)
		}
	}
	return ret, nil
}

// WriteTo executes the template parser and writes the result into w.
func (ts *Tables) WriteTo(w io.Writer) (int64, error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	// deal with random map to guarantee the persistent code generation.
	for _, tblname := range ts.sortedTableNames() {
		t := ts.Tables[tblname] // must panic if table name not found
		if err := t.writeTo(buf, ts.tpl.Funcs(ts.FuncMap)); err != nil {
			return 0, errors.NewWriteFailed(err, "[dmlgen] For Table %q", t.Name)
		}
	}

	if !ts.DisableFileHeader {
		// now figure out all used package names in the buffer.
		fmt.Fprintf(w, "// Auto generated by dmlgen\n\npackage %s\n\nimport (\n", ts.Package)
		pkgs, err := ts.findUsedPackages(buf.Bytes())
		if err != nil {
			return 0, errors.WithStack(err)
		}
		for _, path := range pkgs {
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
