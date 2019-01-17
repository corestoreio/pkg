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
	goparser "go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/pkg/util/slices"
	"github.com/corestoreio/pkg/util/strs"
	"github.com/mailru/easyjson/bootstrap"
	"github.com/mailru/easyjson/parser"
)

// Reason for using this package "shuLhan/go-bindata": at the moment well
// maintained and does not include net/http.
// Run: $ go get github.com/shuLhan/go-bindata
//go:generate go-bindata -o bindata.go -pkg dmlgen ./_tpl/...

// Initial idea and prototyping for code generation.
// TODO DML gen must take care of the types in myreplicator.RowsEvent.decodeValue (not possible)

// TODO generate a hidden type which contains the original data to detect
//  changes and store only changed fields. Investigate other implementations.
//  With the e.g. new field "originalData" pointing to the struct we create a
//  recursion, also we need to track which field needs to be updated after a
//  change. Alternative implementation: a byte slice using protobuf encoding
//  to store the current data and when saving to DB happens decode the data and
//  compare the fields, then return the change field names, except autoinc
//  fields.

// Tables can generated Go source for for database tables once correctly
// configured.
type Tables struct {
	Package            string // Name of the package
	PackageImportPath  string // Name of the package
	ImportPaths        []string
	ImportPathsTesting []string
	// Tables uses the table name as map key and the table description as value.
	Tables map[string]*table
	template.FuncMap
	DisableFileHeader   bool
	DisableTableSchemas bool
	// Serializer defines the de/serializing method to use. Either
	// `empty`=default (none) or proto=protocol buffers or fbs=flatbuffers. JSON
	// marshaling is not affected with this config. Only one serializer can be
	// defined because some Go types depends on the types of serializer. E.g.
	// protobuf cannot use int16/int8 whereas in flatbuffers they are available.
	// Using a common denominator for every serializer like an int32 can cause
	// data loss when communicating with the database table which has different
	// (smaller) column types than an int32.
	Serializer string
	// SerializerHeaderOptions defines custom headers to use in the .proto or .fbs file.
	// For proto, sane defaults are available.
	SerializerHeaderOptions []string
	// TestSQLDumpGlobPath contains the path and glob pattern to load a SQL dump
	// containing the table schemas to run integration tests. If empty no dumps
	// get loaded and the test program assumes that the tables already exists.
	TestSQLDumpGlobPath string
	tpls                *template.Template
	lastError           error

	// customCode injects custom code to manipulate testing and other generate
	// code blocks.
	customCode map[string]string
}

// Option represents a sortable option for the NewTables function. Each option
// function can be applied in a mixed order.
type Option struct {
	// sortOrder specifies the precedence of an option.
	sortOrder int
	fn        func(*Tables) error
}

// TableConfig used in conjunction with WithTableConfig to apply different
// configurations for a generated struct and its struct collection.
type TableConfig struct {
	// Encoders add method receivers for, each struct, compatible with the
	// interface declarations in the various encoding packages. Supported
	// encoder names are: json, binary, and protobuf. Text includes JSON. Binary
	// includes Gob.
	Encoders []string
	// StructTags enables struct tags proactively for the whole struct. Allowed
	// values are: bson, db, env, json, protobuf, toml, yaml and xml. For bson,
	// json, yaml and xml the omitempty attribute has been set. If you need a
	// different struct tag for a specifiv column you must set the option
	// CustomStructTags.
	StructTags []string
	// CustomStructTags allows to specify custom struct tags for a specific
	// column. The slice must be balanced, means index i sets to the column name
	// and index i+1 to the desired struct tag.
	//		[]string{"column_a",`json: ",omitempty"`,"column_b","`xml:,omitempty`"}
	CustomStructTags []string // balanced slice
	// Comment adds custom comments to each struct type. Useful when relying on
	// 3rd party JSON marshaler code generators like easyjson or ffjson. If
	// comment spans over multiple lines each line will be checked if it starts
	// with the comment identifier (//). If not, the identifier will be
	// prepended.
	Comment string
	// ColumnAliases specifies different names used for a column. For example
	// customer_entity.entity_id can also be sales_order.customer_id, hence a
	// Foreign Key. The alias would be just: entity_id:[]string{"customer_id"}.
	ColumnAliases map[string][]string // key=column name value a list of aliases
	// UniquifiedColumns specifies columns which are non primary/unique key one
	// but should have a dedicated function to extract their unique primitive
	// values as a slice. Not allowed are text, blob and binary.
	UniquifiedColumns []string
	// PrivateFields list struct field names which should be private to avoid
	// accidentally leaking through encoders. Appropriate getter/setter methods
	// get generated.
	PrivateFields []string
	lastErr       error
}

func (to *TableConfig) applyEncoders(ts *Tables, t *table) {
	for i := 0; i < len(to.Encoders) && to.lastErr == nil; i++ {
		switch enc := to.Encoders[i]; enc {
		case "json":
			t.HasJsonMarshaler = true // for now does nothing
		case "easyjson":
			t.HasEasyJsonMarshaler = true
		case "binary":
			t.HasBinaryMarshaler = true
		case "protobuf", "fbs":
			t.HasSerializer = true // for now leave it in. maybe later PB gets added to the struct tags.
		default:
			to.lastErr = errors.NotSupported.Newf("[dmlgen] WithTableConfig: Table %q Encoder %q not supported", t.TableName, enc)
		}
	}
}

func (to *TableConfig) applyStructTags(t *table) {
	for h := 0; h < len(t.Columns) && to.lastErr == nil; h++ {
		c := t.Columns[h]
		var buf strings.Builder
		for lst, i := len(to.StructTags), 0; i < lst && to.lastErr == nil; i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			// Maybe some types in the struct for a table don't need at
			// all an omitempty so build in some logic which creates the
			// tags more thoughtfully.
			switch tagName := to.StructTags[i]; tagName {
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
			case "max_len":
				l := c.CharMaxLength.Int64
				if c.Precision.Valid {
					l = c.Precision.Int64
				}
				if l > 0 {
					fmt.Fprintf(&buf, `max_len:"%d"`, l)
				}
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
				to.lastErr = errors.NotSupported.Newf("[dmlgen] WithTableConfig: Table %q Tag %q not supported", t.TableName, tagName)
			}
		}
		c.StructTag = buf.String()
	} // end Columns loop
}

func (to *TableConfig) applyCustomStructTags(t *table) {
	for i := 0; i < len(to.CustomStructTags) && to.lastErr == nil; i = i + 2 {
		found := false
		for _, c := range t.Columns {
			if c.Field == to.CustomStructTags[i] {
				c.StructTag = to.CustomStructTags[i+1]
				found = true
			}
		}
		if !found {
			to.lastErr = errors.NotFound.Newf("[dmlgen] WithTableConfig:CustomStructTags: For table %q the Column %q cannot be found.",
				t.TableName, to.CustomStructTags[i])
		}
	}
}

func (to *TableConfig) applyPrivateFields(t *table) {
	if len(to.PrivateFields) > 0 && t.privateFields == nil {
		t.privateFields = make(map[string]bool)
	}
	for _, pf := range to.PrivateFields {
		t.privateFields[pf] = true
	}
}

func (to *TableConfig) applyComments(t *table) {
	var buf strings.Builder
	lines := strings.Split(to.Comment, "\n")
	for i := 0; i < len(lines) && to.lastErr == nil && to.Comment != ""; i++ {
		line := lines[i]
		if !strings.HasPrefix(line, "//") {
			buf.WriteString("// ")
		}
		buf.WriteString(line)
		buf.WriteByte('\n')
	}
	t.Comment = strings.TrimSpace(buf.String())
}

func (to *TableConfig) applyColumnAliases(t *table) {
	if to.lastErr != nil {
		return
	}
	// With iteration looping it this way, we can easily proof if the
	// developer has correctly written the column name. We might have
	// more work here but the developer has a better experience when a
	// column can't be found.
	for colName, aliases := range to.ColumnAliases {
		found := false
		for _, col := range t.Columns {
			if col.Field == colName {
				found = true
				col.Aliases = aliases
			}
		}
		if !found {
			to.lastErr = errors.NotFound.Newf("[dmlgen] WithTableConfig:ColumnAliases: For table %q the Column %q cannot be found.",
				t.TableName, colName)
			return
		}
	}
}

// skips text and blob and varbinary and json and geo
func (to *TableConfig) applyUniquifiedColumns(t *table) {
	for i := 0; i < len(to.UniquifiedColumns) && to.lastErr == nil; i++ {
		cn := to.UniquifiedColumns[i]
		found := false
		for _, c := range t.Columns {
			if c.Field == cn && false == c.IsBlobDataType() {
				c.Uniquified = true
				found = true
			}
		}
		if !found {
			to.lastErr = errors.NotFound.Newf("[dmlgen] WithTableConfig:UniquifiedColumns: For table %q the Column %q cannot be found in the list of available columns or its data type is not allowed.",
				t.TableName, cn)
		}
	}
}

// WithTableConfig applies options to a table, identified by the table name used
// as map key. Options are custom struct or different encoders.
func WithTableConfig(tableName string, opt *TableConfig) (o Option) {
	// Panic as early as possible.
	if len(opt.CustomStructTags)%2 == 1 {
		panic(errors.Fatal.Newf("[dmlgen] WithTableConfig: Table %q option CustomStructTags must be a balanced slice.", tableName))
	}
	o.sortOrder = 150
	o.fn = func(ts *Tables) (err error) {
		t, ok := ts.Tables[tableName]
		if t == nil || !ok {
			return errors.NotFound.Newf("[dmlgen] WithTableConfig: Table %q not found.", tableName)
		}
		opt.applyEncoders(ts, t)
		opt.applyStructTags(t)
		opt.applyCustomStructTags(t)
		opt.applyPrivateFields(t)
		opt.applyComments(t)
		opt.applyColumnAliases(t)
		opt.applyUniquifiedColumns(t)
		return opt.lastErr
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

// WithReferenceEntitiesByForeignKeys analyses the foreign keys which points to
// a table and adds them as a struct field name. For example:
// customer_address_entity.parent_id is a foreign key to
// customer_entity.entity_id hence the generated struct CustomerEntity has a new
// field which gets named CustomerAddressEntityCollection, pointing to type
// CustomerAddressEntityCollection. structFieldNameMapperFn can be nil, if so
// the name of the Go collection type gets used as field name.
func WithReferenceEntitiesByForeignKeys(ctx context.Context, db dml.Querier, structFieldNameMapperFn func(string) string) (opt Option) {
	// initial implementation. might not work correctly with the EAV tables as
	// it might to point to to many tables/structs.

	if structFieldNameMapperFn == nil {
		structFieldNameMapperFn = func(s string) string { return s }
	}

	opt.sortOrder = 210 // must run at the end or where the end is near ;-)
	opt.fn = func(ts *Tables) error {

		tblFks, err := ddl.LoadKeyColumnUsage(ctx, db, ts.sortedTableNames()...)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, kcuc := range tblFks { // kcuc = keyColumnUsageCollection
			for _, kcuce := range kcuc.Data {
				if kcuce.ReferencedTableName.Valid {
					t := ts.Tables[kcuce.ReferencedTableName.String]
					t.ReferencedCollections = append(t.ReferencedCollections,
						structFieldNameMapperFn(kcuce.TableName)+" "+strs.ToGoCamelCase(kcuce.TableName)+"Collection")
				}
			}
		}
		return nil
	}
	return
}

// WithTable sets a table and its columns. Allows to overwrite a table fetched
// with function WithLoadColumns. Argument `actions` can only be set to "overwrite".
func WithTable(tableName string, columns ddl.Columns, actions ...string) (opt Option) {
	opt.sortOrder = 10
	opt.fn = func(ts *Tables) error {
		isOverwrite := len(actions) > 0 && actions[0] == "overwrite"
		t, ok := ts.Tables[tableName]
		if ok && isOverwrite {
			for ci, ct := range t.Columns {
				for _, cc := range columns {
					if ct.Field == cc.Field {
						t.Columns[ci] = cc
					}
				}
			}
		} else {
			t = &table{
				TableName: tableName,
				Columns:   columns,
			}
		}
		ts.Tables[tableName] = t
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

// WithProtobuf enables protocol buffers as a serialization method. Argument
// headerOptions is optional.
func WithProtobuf(headerOptions ...string) (opt Option) {
	opt.sortOrder = 110
	opt.fn = func(ts *Tables) error {
		ts.Serializer = "protobuf"
		if len(headerOptions) == 0 {
			ts.SerializerHeaderOptions = []string{
				"(gogoproto.typedecl_all) = false",
				"(gogoproto.goproto_getters_all) = false",
				"(gogoproto.unmarshaler_all) = true",
				"(gogoproto.marshaler_all) = true",
				"(gogoproto.sizer_all) = true",
				"(gogoproto.goproto_unrecognized_all) = false",
			}
		}
		return nil
	}
	return
}

// WithFlatbuffers enables flatbuffers (FBS) as a serialization method. Argument
// headerOptions is optional.
func WithFlatbuffers(headerOptions ...string) (opt Option) {
	opt.sortOrder = 111
	opt.fn = func(ts *Tables) error {
		ts.Serializer = "fbs"
		if len(headerOptions) == 0 {
			// TODO find sane defaults
			// ts.SerializerHeaderOptions = []string{}
		}
		return nil
	}
	return
}

// WithCustomCode inserts at the marker position your custom Go code. For
// available markers search the .go.tpl files for the function call
// `CustomCode`. An example got written in TestGenerate_Tables_Protobuf_Json. If
// the marker does not exists or has a typo, no error gets reported and no code
// gets written.
func WithCustomCode(marker, code string) (opt Option) {
	opt.sortOrder = 112
	opt.fn = func(ts *Tables) error {
		if ts.customCode == nil {
			ts.customCode = make(map[string]string)
		}
		ts.customCode[marker] = code
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

// NewTables creates a new instance of the SQL table code generator. The order
// of the applied options does not matter as they are getting sorted internally.
func NewTables(packageImportPath string, opts ...Option) (*Tables, error) {
	_, pkg := filepath.Split(packageImportPath)
	ts := &Tables{
		Tables:            make(map[string]*table),
		Package:           pkg,
		PackageImportPath: packageImportPath,
		ImportPaths: []string{
			"context",
			"database/sql",
			"encoding/json",
			"sort",
			"time",

			"github.com/corestoreio/errors",
			"github.com/corestoreio/pkg/sql/ddl",
			"github.com/corestoreio/pkg/sql/dml",
			"github.com/corestoreio/pkg/storage/null",
		},
		ImportPathsTesting: []string{
			"testing",
			"context",
			"sort",
			"github.com/corestoreio/pkg/sql/dmltest",
			"github.com/corestoreio/pkg/util/assert",
			"github.com/corestoreio/pkg/util/pseudo",
		},
		FuncMap: make(template.FuncMap, 20),
	}

	sort.Slice(opts, func(i, j int) bool {
		return opts[i].sortOrder < opts[j].sortOrder // ascending 0-9 sorting ;-)
	})

	for _, opt := range opts {
		if err := opt.fn(ts); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	ts.FuncMap["CustomCode"] = func(marker string) string { return ts.customCode[marker] }
	ts.FuncMap["GoCamel"] = strs.ToGoCamelCase // net_http->NetHTTP entity_id->EntityID
	ts.FuncMap["GoCamelMaybePrivate"] = func(s string) string { return s }
	ts.FuncMap["GoTypeNull"] = func(c *ddl.Column) string { return ts.mySQLToGoType(c, true) }
	ts.FuncMap["GoType"] = func(c *ddl.Column) string { return ts.mySQLToGoType(c, false) }
	ts.FuncMap["GoFuncNull"] = func(c *ddl.Column) string { return ts.mySQLToGoDmlColumnMap(c, true) }
	ts.FuncMap["GoFunc"] = func(c *ddl.Column) string { return ts.mySQLToGoDmlColumnMap(c, false) }
	ts.FuncMap["IsFieldPublic"] = func(string) bool { return false }
	ts.FuncMap["IsFieldPrivate"] = func(string) bool { return false }
	ts.FuncMap["GoPrimitiveNull"] = ts.toGoPrimitiveFromNull
	ts.FuncMap["SerializerType"] = func(c *ddl.Column) string {
		pt := ts.toSerializerType(c, true)
		if strings.IndexByte(pt, '/') > 0 { // slash identifies an import path
			return "bytes"
		}
		return pt
	}
	ts.FuncMap["SerializerCustomType"] = func(c *ddl.Column) string {
		pt := ts.toSerializerType(c, true)
		var buf strings.Builder
		if pt == "google.protobuf.Timestamp" {
			fmt.Fprint(&buf, ",(gogoproto.stdtime)=true")
		}
		if pt == "bytes" {
			return "" // bytes can be null
		}
		if c.IsNull() || strings.IndexByte(pt, '.') > 0 /*whenever it is a custom type like null. or google.proto.timestamp*/ {
			// Indeed nullable Go Types must be not-nullable in HasSerializer because we
			// have a non-pointer struct type which contains the field Valid.
			// HasSerializer treats nullable fields as pointer fields, but that is
			// ridiculous.
			fmt.Fprint(&buf, ",(gogoproto.nullable)=false")
		}
		return buf.String()
	}

	ts.tpls = template.New("InitialParseGlob").Funcs(ts.FuncMap)
	for _, f := range AssetNames() {
		data, err := Asset(f)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if _, err := ts.tpls.New(filepath.Base(f)).Parse(string(data)); err != nil {
			return nil, errors.Wrapf(err, "[dmlgen] Failed to parse %q", f)
		}
	}

	for _, t := range ts.Tables {
		t.Package = ts.Package
	}
	return ts, nil
}

// findUsedPackages checks for needed packages which we must import.
func (ts *Tables) findUsedPackages(file []byte) ([]string, error) {

	af, err := goparser.ParseFile(token.NewFileSet(), "cs_virtual_file.go", append([]byte("package temporarily_main\n\n"), file...), 0)
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

// GenerateSerializer writes the protocol buffer specifications into `w` and its test
// sources into wTest, if there are any tests.
func (ts *Tables) GenerateSerializer(w io.Writer, wTest io.Writer) error {
	switch ts.Serializer {
	case "protobuf", "fbs":
	// supported
	case "", "default", "none":
		return nil // do nothing
	default:
		return errors.NotAcceptable.Newf("[dmlgen] Serializer %q not supported.", ts.Serializer)
	}

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if !ts.DisableFileHeader {
		if err := ts.tpls.Funcs(ts.FuncMap).ExecuteTemplate(buf, ts.Serializer+"_10_header.go.tpl", ts); err != nil {
			return errors.WriteFailed.New(err, "[dmlgen] For file header")
		}
	}

	tableFileName := ts.Serializer + "_20_message.go.tpl"
	for _, tblName := range ts.sortedTableNames() {
		t := ts.Tables[tblName] // must panic if table name not found
		ts.FuncMap["IsFieldPublic"] = t.IsFieldPublic
		if t.HasSerializer {
			if err := t.writeTo(buf, ts.tpls.Lookup(tableFileName).Funcs(ts.FuncMap)); err != nil {
				return errors.WriteFailed.New(err, "[dmlgen] For Table %q", t.TableName)
			}
		}
	}
	_, err := buf.WriteTo(w)
	return err
}

func (ts *Tables) execTpl(w io.Writer, t *table, tplName string) {
	if ts.lastError != nil {
		return
	}
	if err := t.writeTo(w, ts.tpls.Lookup(tplName)); err != nil {
		ts.lastError = errors.WriteFailed.New(err, "[dmlgen] With template %q for Table %q", tplName, t.TableName)
	}
}

// GenerateGo writes the Go source code into `w` and the test code into wTest.
func (ts *Tables) GenerateGo(w io.Writer, wTest io.Writer) error {
	buf := new(bytes.Buffer)
	bufTest := new(bytes.Buffer)
	ts.tpls = ts.tpls.Funcs(ts.FuncMap)

	sortedTableNames := ts.sortedTableNames()
	if !ts.DisableTableSchemas { // Writes the table DDL function
		tables := make([]*table, len(ts.Tables))
		for i, tblname := range sortedTableNames {
			tables[i] = ts.Tables[tblname] // must panic if table name not found
		}
		data := struct {
			Package             string // Name of the package
			Tables              []*table
			TableNames          []string
			TestSQLDumpGlobPath string
		}{
			Package:             ts.Package,
			Tables:              tables,
			TableNames:          sortedTableNames,
			TestSQLDumpGlobPath: ts.TestSQLDumpGlobPath,
		}
		if err := ts.tpls.ExecuteTemplate(buf, "10_tables.go.tpl", data); err != nil {
			return errors.WriteFailed.New(err, "[dmlgen] For Tables %v", tables)
		}
		if err := ts.tpls.ExecuteTemplate(bufTest, "90_test.go.tpl", data); err != nil {
			return errors.WriteFailed.New(err, "[dmlgen] For Tables %v", tables)
		}
	}

	// deal with random map to guarantee the persistent code generation.
	for _, tblname := range sortedTableNames {
		t, ok := ts.Tables[tblname]
		if !ok || t == nil {
			return errors.NotFound.Newf("[dmlgen] Table %q not found", tblname)
		}
		ts.FuncMap["IsFieldPrivate"] = t.IsFieldPrivate
		ts.FuncMap["GoCamelMaybePrivate"] = t.GoCamelMaybePrivate
		ts.tpls = ts.tpls.Funcs(ts.FuncMap)

		ts.execTpl(buf, t, "20_entity.go.tpl")
		if !t.DisableCollectionMethods {
			ts.execTpl(buf, t, "30_collection_methods.go.tpl")
		}
		if t.HasBinaryMarshaler {
			ts.execTpl(buf, t, "40_binary.go.tpl")
		}
		if ts.lastError != nil {
			return ts.lastError
		}
	}

	if !ts.DisableFileHeader {
		// now figure out all used package names in the buffer.
		fmt.Fprintf(w, "// Auto generated via github.com/corestoreio/pkg/sql/dmlgen\n\npackage %s\n\nimport (\n", ts.Package)
		// println(buf.String())
		pkgs, err := ts.findUsedPackages(buf.Bytes())
		if err != nil {
			_, _ = w.Write(buf.Bytes()) // write malformed data for debugging reasons.
			return errors.WithStack(err)
		}
		for _, path := range pkgs {
			fmt.Fprintf(w, "\t%q\n", path)
		}
		fmt.Fprintf(w, "\n)\n")

		// now figure out all used package names in the buffer.
		fmt.Fprintf(wTest, "// Auto generated via github.com/corestoreio/pkg/sql/dmlgen\n\npackage %s\n\nimport (\n", ts.Package)

		for _, path := range ts.ImportPathsTesting {
			fmt.Fprintf(wTest, "\t%q\n", path)
		}
		fmt.Fprintf(wTest, "\n)\n")
	}

	fmted, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.WithStack(err)
	}
	buf.Reset()
	buf.Write(fmted)
	if _, err = buf.WriteTo(w); err != nil {
		return errors.WithStack(err)
	}

	fmted, err = format.Source(bufTest.Bytes())
	if err != nil {
		return errors.WithStack(err)
	}
	bufTest.Reset()
	bufTest.Write(fmted)
	_, err = bufTest.WriteTo(wTest)
	return err
}

// table writes one database table into Go source code.
type table struct {
	Package   string      // Name of the package
	TableName string      // Name of the table
	Comment   string      // Comment above the struct type declaration
	Columns   ddl.Columns // all columns of a table
	// ReferencedCollections, map key is the name of the struct field, map value
	// the target Go collection type.
	ReferencedCollections    []string
	HasJsonMarshaler         bool
	HasEasyJsonMarshaler     bool
	HasBinaryMarshaler       bool
	HasSerializer            bool // writes the .proto file if true
	DisableCollectionMethods bool
	// PrivateFields key=snake case name of the DB column, value=true, the field must be private
	privateFields map[string]bool
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

func (t *table) IsFieldPublic(dbColumnName string) bool {
	return t.privateFields == nil || !t.privateFields[dbColumnName]
}

func (t *table) IsFieldPrivate(dbColumnName string) bool {
	return t.privateFields != nil && t.privateFields[dbColumnName]
}

func (t *table) GoCamelMaybePrivate(s string) string {
	su := strs.ToGoCamelCase(s)
	if t.IsFieldPublic(s) {
		return su
	}
	sr := []rune(su)
	sr[0] = unicode.ToLower(sr[0])
	return string(sr)
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
		return errors.Wrapf(err, "[dmlgen] Can't access proto files in path %q", path)
	}

	// To generate PHP Code replace `gogo_out` with `php_out`.
	// Java bit similar. Java has ~15k LOC, Go ~3.7k
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
			return errors.WriteFailed.Newf("[dmlgen] protoc Error: %s", text)
		}
	}

	// what a hack: find all *.pb.go files and remove `import null
	// "github.com/corestoreio/pkg/storage/null"` because no other way to get
	// rid of the unused import or reference that import somehow in the
	// generated file :-( Once there's a better solution, remove this code.
	pbGoFiles, err := filepath.Glob(path + "*.pb.go")
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] Can't access pb.go files in path %q", path)
	}
	removeImport := []byte("import null \"github.com/corestoreio/pkg/storage/null\"\n")
	for _, file := range pbGoFiles {
		fContent, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.WithStack(err)
		}
		fContent = bytes.Replace(fContent, removeImport, nil, -1)
		if err := ioutil.WriteFile(file, fContent, 0644); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// GenerateJSON creates the easysjon code for a specific file or a whole
// directory. argument `g` can be nil.
func GenerateJSON(fname string, g *bootstrap.Generator) (err error) {
	fInfo, err := os.Stat(fname)
	if err != nil {
		return err
	}

	p := parser.Parser{}
	if err := p.Parse(fname, fInfo.IsDir()); err != nil {
		return errors.CorruptData.Newf("[dmlgen] Error parsing failed %v: %v", fname, err)
	}

	var outName string
	if fInfo.IsDir() {
		outName = filepath.Join(fname, p.PkgName+"_easyjson.go")
	} else {
		if s := strings.TrimSuffix(fname, ".go"); s == fname {
			return errors.New("Filename must end in '.go'")
		} else {
			outName = s + "_easyjson.go"
		}
	}

	var trimmedBuildTags string
	// if *buildTags != "" {
	// 	trimmedBuildTags = strings.TrimSpace(*buildTags)
	// }
	if g == nil {
		g = &bootstrap.Generator{
			BuildTags:             trimmedBuildTags,
			PkgPath:               p.PkgPath,
			PkgName:               p.PkgName,
			Types:                 p.StructNames,
			SnakeCase:             true,
			LowerCamelCase:        true,
			NoStdMarshalers:       false,
			DisallowUnknownFields: false,
			OmitEmpty:             true,
			LeaveTemps:            false,
			OutName:               outName,
			StubsOnly:             false,
			NoFormat:              false,
		}
	}
	if err := g.Run(); err != nil {
		return errors.Fatal.Newf("[dmlgen] easyJSON: Bootstrap failed: %v", err)
	}
	return nil
}
