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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/util/slices"
	"github.com/corestoreio/pkg/util/strs"
)

// Initial idea and prototyping for code generation.

const pkgPath = `src/github.com/corestoreio/pkg/sql/dmlgen`

// Tables can generated Go source for for database tables once correctly
// configured.
type Tables struct {
	Package     string // Name of the package
	ImportPaths []string
	// Tables uses the table name as map key and the table description as value.
	Tables map[string]*table
	template.FuncMap
	DisableFileHeader   bool
	DisableTableSchemas bool
	GogoProtoOptions    []string
	// goTpl contains a parsed template to render a single table.
	tpls       *template.Template
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
// in the various encoding packages. Supported encoder names are: text, binary,
// and protobuf. Text includes JSON. Binary includes Gob.
func WithEncoder(tableName string, encoderNames ...string) (opt Option) {
	opt.sortOrder = 90 // must run before custom struct tags
	opt.fn = func(ts *Tables) (err error) {
		if ts.Tables[tableName] == nil {
			return errors.NewNotFoundf("[dmlgen] WithEncoder: Table %q not found.", tableName)
		}
		for _, enc := range encoderNames {
			switch enc {
			case "text":
				ts.Tables[tableName].TextMarshaler = true
			case "binary":
				ts.Tables[tableName].BinaryMarshaler = true
			case "protobuf":
				// github.com/gogo/protobuf/protoc-gen-gogo/generator/generator.go#L1629 Generator.goTag
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
// values are: bson, db, env, json, protobuf, toml, yaml and xml. For bson,
// json, yaml and xml the omitempty attribute has been set. If you need a
// different struct tag for a specifiv column you must set the option
// WithCustomStructTags. It doesn't matter in which order you apply the options
// ;-)
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
					case "protobuf":
						// github.com/gogo/protobuf/protoc-gen-gogo/generator/generator.go#L1629 Generator.goTag
						// The tag is a string like "varint,2,opt,name=fieldname,def=7" that
						// identifies details of the field for the protocol buffer marshaling and unmarshaling
						// code.  The fields are:
						//	wire encoding
						//	protocol tag number
						//	opt,req,rep for optional, required, or repeated
						//	packed whether the encoding is "packed" (optional; repeated primitives only)
						//	name= the original declared name
						//	enum= the name of the enum type if it is an enum-typed field.
						//	proto3 if this field is in a proto3 message
						//	def= string representation of the default value, if any.
						// The default value must be in a representation that can be used at run-time
						// to generate the default value. Thus bools become 0 and 1, for instance.

						// CYS: not quite sure if struct tags are really needed
						//pbType := "TODO"
						//customType := ",customtype=github.com/gogo/protobuf/test.TODO"
						//fmt.Fprintf(&buf, `protobuf:"%s,%d,opt,name=%s%s"`, pbType, c.Pos, c.Field, customType)
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
	if len(columnNameStructTag)%2 == 1 {
		panic(errors.NewFatalf("[dmlgen] WithCustomStructTags: Argument columnNameStructTag must be a balanced slice."))
	}
	opt.sortOrder = 100
	opt.fn = func(ts *Tables) error {
		t, ok := ts.Tables[tableName]
		if !ok {
			return errors.NewNotFoundf("[dmlgen] WithCustomStructTags: Table %q cannot be found.", tableName)
		}
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
			return errors.NewNotFoundf("[dmlgen] WithCustomStructTags: For table %q one column in %v cannot be found.", tableName, columnNameStructTag)
		}
		return nil
	}
	return
}

// WithStructComment adds custom comments to each struct type. Useful when
// relying on 3rd party JSON marshaler code generators like easyjson or ffjson.
// If comment spans over multiple lines each line will be checked if it starts
// with the comment identifier (//). If not the identifier will be prepended.
func WithStructComment(tableNameComment ...string) (opt Option) {
	if len(tableNameComment)%2 == 1 {
		panic(errors.NewFatalf("[dmlgen] WithStructComment: Argument tableNameComment must be a balanced slice."))
	}
	opt.sortOrder = 105
	opt.fn = func(ts *Tables) error {
		for i := 0; i < len(tableNameComment); i = i + 2 {
			tableName := tableNameComment[i]
			t, ok := ts.Tables[tableName]
			if !ok {
				return errors.NewNotFoundf("[dmlgen] WithStructComment: Table %q cannot be found.", tableName)
			}
			var buf strings.Builder
			for _, line := range strings.Split(tableNameComment[i+1], "\n") {
				if !strings.HasPrefix(line, "//") {
					buf.WriteString("// ")
				}
				buf.WriteString(line)
				buf.WriteByte('\n')
			}
			t.Comment = buf.String()
		}
		return nil
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
			"github.com/corestoreio/pkg/sql/dml",
			"github.com/corestoreio/pkg/sql/ddl",
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

	if len(ts.GogoProtoOptions) == 0 {
		ts.GogoProtoOptions = []string{
			"(gogoproto.typedecl_all) = false",
			"(gogoproto.goproto_getters_all) = false",
			"(gogoproto.unmarshaler_all) = true",
			"(gogoproto.marshaler_all) = true",
			"(gogoproto.sizer_all) = true",
		}
	}

	sort.Slice(opts, func(i, j int) bool {
		return opts[i].sortOrder < opts[j].sortOrder // ascending 0-9 sorting ;-)
	})

	for _, opt := range opts {
		if err := opt.fn(ts); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	var err error
	glob := filepath.Join(build.Default.GOPATH, pkgPath, "code_*.go.tpl")
	ts.tpls, err = template.New("InitialParseGlob").Funcs(ts.FuncMap).ParseGlob(glob)
	if err != nil {
		return nil, errors.WithStack(err)
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
		ts.tpls.Funcs(ts.FuncMap).ExecuteTemplate(buf, "code_proto_header.go.tpl", ts)
	}

	for _, tblname := range ts.sortedTableNames() {
		t := ts.Tables[tblname] // must panic if table name not found
		if err := t.writeTo(buf, ts.tpls.Lookup("code_proto.go.tpl").Funcs(ts.FuncMap)); err != nil {
			return errors.NewWriteFailed(err, "[dmlgen] For Table %q", t.TableName)
		}
	}
	_, err := buf.WriteTo(w)
	return err
}

// WriteGo writes the Go source code into `w`.
func (ts *Tables) WriteGo(w io.Writer) error {
	buf := new(bytes.Buffer)
	tpl := ts.tpls.Funcs(ts.FuncMap)

	sortedTableNames := ts.sortedTableNames()
	if !ts.DisableTableSchemas { // Writes the table DDL function
		tables := make([]*table, len(ts.Tables))
		for i, tblname := range sortedTableNames {
			tables[i] = ts.Tables[tblname] // must panic if table name not found
		}
		data := struct {
			Package    string // Name of the package
			Tables     []*table
			TableNames []string
		}{
			Package:    ts.Package,
			Tables:     tables,
			TableNames: sortedTableNames,
		}
		if err := tpl.ExecuteTemplate(buf, "code_tables.go.tpl", data); err != nil {
			return errors.NewWriteFailed(err, "[dmlgen] For Tables %v", tables)
		}
	}

	// deal with random map to guarantee the persistent code generation.
	for _, tblname := range sortedTableNames {
		t := ts.Tables[tblname] // must panic if table name not found

		if err := t.writeTo(buf, tpl.Lookup("code_entity.go.tpl")); err != nil {
			return errors.NewWriteFailed(err, "[dmlgen] For Table %q", t.TableName)
		}
		if err := t.writeTo(buf, tpl.Lookup("code_collection.go.tpl")); err != nil {
			return errors.NewWriteFailed(err, "[dmlgen] For Table %q", t.TableName)
		}
		if t.TextMarshaler {
			if err := t.writeTo(buf, tpl.Lookup("code_text.go.tpl")); err != nil {
				return errors.NewWriteFailed(err, "[dmlgen] For Table %q", t.TableName)
			}
		}
		if t.BinaryMarshaler {
			if err := t.writeTo(buf, tpl.Lookup("code_binary.go.tpl")); err != nil {
				return errors.NewWriteFailed(err, "[dmlgen] For Table %q", t.TableName)
			}
		}
	}

	if !ts.DisableFileHeader {
		// now figure out all used package names in the buffer.
		fmt.Fprintf(w, "// Auto generated via github.com/corestoreio/pkg/sql/dmlgen\n\npackage %s\n\nimport (\n", ts.Package)
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
	Comment         string      // Comment above the struct type declaration
	Columns         ddl.Columns // all columns of a table
	TextMarshaler   bool
	BinaryMarshaler bool
	Protobuf        bool // writes the .proto file if true
}

// WriteTo implements io.WriterTo and writes the generated source code into w.
func (t *table) writeTo(w io.Writer, tpl *template.Template) error {

	data := struct {
		table
		Collection string
		Entity     string
	}{
		table:      *t,
		Collection: strs.ToGoCamelCase(t.TableName) + "Collection",
		Entity:     strs.ToGoCamelCase(t.TableName),
	}

	return tpl.Execute(w, data)
}

// GenerateProto searches all *.proto files in the given path and calls protoc
// to generate the Go source code.
func GenerateProto(path string) error {

	path = filepath.Clean(path)
	if ps := string(os.PathSeparator); !strings.HasSuffix(path, ps) {
		path += ps
	}

	protoFiles, err := filepath.Glob(path + "*.proto")
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] Can't find proto files in path %q", path)
	}

	args := []string{
		"--gogo_out", "Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:.",
		"--proto_path", fmt.Sprintf("%s/src/:%s/src/github.com/gogo/protobuf/protobuf/:.", build.Default.GOPATH, build.Default.GOPATH),
	}
	args = append(args, protoFiles...)

	cmd := exec.Command("protoc", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] %s", out)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		text := scanner.Text()
		if !strings.Contains(text, "WARNING") {
			return errors.NewWriteFailedf("[dmlgen] protoc Error: %s", text)
		}
	}
	return nil
}
