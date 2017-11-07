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
	"strings"
	"text/template"

	"github.com/corestoreio/csfw/sql/ddl"
	"github.com/corestoreio/csfw/sql/dml"
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
	// goTpl contains a parsed template to render a single table.
	goTpl      *template.Template
	protoTpl   *template.Template
	writeProto bool
}

// Option represents a sortable option for the NewTables function. Each option
// function can be applied in a mixed order.
type Option struct {
	// sortOrder specifies the precedence of an option.
	sortOrder int
	fn        func(*Tables) error
}

// WithEncoder adds method receivers compatible with the interface declarations
// in the various encoding packages. Supported encoder names are: json, binary,
// gob and proto. More to follow.
func WithEncoder(tableName string, encoderNames ...string) (opt Option) {
	opt.sortOrder = 90 // must run before custom struct tags
	opt.fn = func(ts *Tables) (err error) {
		if ts.Tables[tableName] == nil {
			return errors.NewNotFoundf("[dmlgen] WithEncoder: Table %q not found.", tableName)
		}
		for _, enc := range encoderNames {
			switch enc {
			case "json":
				ts.Tables[tableName].JsonMarshaler = true
			case "binary":
				ts.Tables[tableName].BinaryMarshaler = true
			case "gob":
				ts.Tables[tableName].GobEncoding = true
			case "proto":
				ts.writeProto = true
				ts.Tables[tableName].Protobuf = true // for now leave it in. maybe later PB gets added to the struct tags.
			default:
				return errors.NewNotSupportedf("[dmlgen] WithMarshaler: encoder %q not supported", enc)
			}
		}
		return nil
	}
	return
}

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

// WithColumnAliasesFromForeignKeys extracts similar column names from foreign
// key definitions. For the list of tables and their primary/unique keys, this
// function searches the foreign keys to other tables and uses the column name
// as the alias. For example the table `customer_entity` and its PK column
// `entity_id` has a foreign key in table `sales_order` whose name is
// `customer_id`. When generation code for customer_entity, the column entity_id
// can be used additionally with the name customer_id, hence customer_id is the
// alias.
func WithColumnAliasesFromForeignKeys(ctx context.Context, db dml.Querier) (opt Option) {
	opt.sortOrder = 200 // must run at the end or where the end is near ;-)
	opt.fn = func(ts *Tables) error {

		tblFks, err := ddl.LoadKeyColumnUsage(ctx, db, ts.sortedTableNames()...)
		if err != nil {
			return errors.WithStack(err)
		}
		for tblPkCol, kcuc := range tblFks {
			// tblPkCol == REFERENCED_TABLE_NAME.REFERENCED_COLUMN_NAME
			// REFERENCED_TABLE_NAME is contained in sortedTableNames()
			dotPos := strings.IndexByte(tblPkCol, '.')
			refTable := tblPkCol[:dotPos]
			refColumn := tblPkCol[dotPos+1:]

			t := ts.Tables[refTable]
			for _, c := range t.Columns {
				// TODO: optimize this and rethink method receivers like Each, on the collection.
				if c.Field == refColumn {
					unique := map[string]bool{refColumn: true} // refColumn already seen because field name
					for _, kcu := range kcuc.Data {
						if kcu.ReferencedColumnName.String == refColumn && !unique[kcu.ColumnName] {
							c.Aliases = append(c.Aliases, kcu.ColumnName)
							unique[kcu.ColumnName] = true
						}
					}
				}
			}
		}
		return nil
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
			TableName: tableName,
			Columns:   columns,
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
				TableName: tblName,
				Columns:   tables[tblName],
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
func NewTables(packageName string, opts ...Option) (*Tables, error) {
	ts := &Tables{
		Tables:  make(map[string]*table),
		Package: packageName,
		ImportPaths: []string{
			"database/sql",
			"encoding/json",
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
	ts.FuncMap["ProtoType"] = toProtoType
	ts.FuncMap["ProtoCustomType"] = toProtoCustomType

	sort.Slice(opts, func(i, j int) bool {
		return opts[i].sortOrder < opts[j].sortOrder // ascending 0-9 sorting ;-)
	})

	for _, opt := range opts {
		if err := opt.fn(ts); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	var err error
	ts.goTpl, err = template.New("go_entity").Funcs(ts.FuncMap).Parse(tplDBAC)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if ts.writeProto {
		ts.protoTpl, err = template.New("proto_entity").Funcs(ts.FuncMap).Parse(tplProto)
		if err != nil {
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

// WriteProto writes the protocol buffer specifications into `w`.
func (ts *Tables) WriteProto(w io.Writer) error {
	buf := new(bytes.Buffer)
	if !ts.writeProto {
		return errors.NewNotAcceptablef("[dmlgen] Protocol buffer generation not enabled.")
	}

	if !ts.DisableFileHeader {
		fmt.Fprintf(buf, `// Auto generated via github.com/corestoreio/csfw/sql/dmlgen
syntax = "proto3";
package %s;
import "github.com/gogo/protobuf/gogoproto/gogo.proto";

option (gogoproto.typedecl_all) = false;
option (gogoproto.unmarshaler_all) = true;
option (gogoproto.marshaler_all) = true;
option (gogoproto.sizer_all) = true;

`, ts.Package)
	}

	for _, tblname := range ts.sortedTableNames() {
		t := ts.Tables[tblname] // must panic if table name not found
		if err := t.writeTo(buf, ts.protoTpl.Funcs(ts.FuncMap)); err != nil {
			return errors.NewWriteFailed(err, "[dmlgen] For Table %q", t.TableName)
		}
	}
	_, err := buf.WriteTo(w)
	return err
}

// WriteGo writes the Go source code into `w`.
func (ts *Tables) WriteGo(w io.Writer) error {
	buf := new(bytes.Buffer)

	// deal with random map to guarantee the persistent code generation.
	for _, tblname := range ts.sortedTableNames() {
		t := ts.Tables[tblname] // must panic if table name not found
		if err := t.writeTo(buf, ts.goTpl.Funcs(ts.FuncMap)); err != nil {
			return errors.NewWriteFailed(err, "[dmlgen] For Table %q", t.TableName)
		}
	}

	if !ts.DisableFileHeader {
		// now figure out all used package names in the buffer.
		fmt.Fprintf(w, "// Auto generated via github.com/corestoreio/csfw/sql/dmlgen\n\npackage %s\n\nimport (\n", ts.Package)
		pkgs, err := ts.findUsedPackages(buf.Bytes())
		if err != nil {
			return errors.WithStack(err)
		}
		for _, path := range pkgs {
			fmt.Fprintf(w, "\t%q\n", path)
		}
		fmt.Fprint(w, "\n)\n")
	}

	fmted, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.WithStack(err)
	}
	buf.Reset()
	buf.Write(fmted)
	_, err = buf.WriteTo(w)
	return err
}

// table writes one database table into Go source code.
type table struct {
	Package         string      // Name of the package
	TableName       string      // Name of the table
	Columns         ddl.Columns // all columns of the table
	JsonMarshaler   bool
	BinaryMarshaler bool
	GobEncoding     bool
	Protobuf        bool
}

// WriteTo implements io.WriterTo and writes the generated source code into w.
func (t *table) writeTo(w io.Writer, tpl *template.Template) error {

	data := struct {
		table
		Collection string
		Entity     string
		Tick       string
	}{
		table:      *t,
		Collection: strs.ToGoCamelCase(t.TableName) + "Collection",
		Entity:     strs.ToGoCamelCase(t.TableName),
		Tick:       "`",
	}

	return tpl.Execute(w, data)
}
