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
	goparser "go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/pkg/util/codegen"
	"github.com/corestoreio/pkg/util/slices"
	"github.com/corestoreio/pkg/util/strs"
	"github.com/mailru/easyjson/bootstrap"
	"github.com/mailru/easyjson/parser"
)

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
	BuildTags          []string
	ImportPaths        []string
	ImportPathsTesting []string
	// Tables uses the table name as map key and the table description as value.
	Tables              map[string]*table
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
	// 		[]string{"column_a",`json: ",omitempty"`,"column_b","`xml:,omitempty`"}
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
	// DisableCollectionMethods if set, suppresses the generation of the
	// collection related functions.
	DisableCollectionMethods bool
	lastErr                  error
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
				// 	wire encoding
				// 	protocol tag number
				// 	opt,req,rep for optional, required, or repeated
				// 	packed whether the encoding is "packed" (optional; repeated primitives only)
				// 	name= the original declared name
				// 	enum= the name of the enum type if it is an enum-typed field.
				// 	proto3 if this field is in a proto3 message
				// 	def= string representation of the default value, if any.
				// The default value must be in a representation that can be used at run-time
				// to generate the default value. Thus bools become 0 and 1, for instance.

				// CYS: not quite sure if struct tags are really needed
				// pbType := "TODO"
				// customType := ",customtype=github.com/gogo/protobuf/test.TODO"
				// fmt.Fprintf(&buf, `protobuf:"%s,%d,opt,name=%s%s"`, pbType, c.Pos, c.Field, customType)
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

// WithTableConfig applies options to an existing table, identified by the table
// name used as map key. Options are custom struct or different encoders.
// Returns a not-found error if the table cannot be found in the `Tables` map.
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
		t.DisableCollectionMethods = opt.DisableCollectionMethods
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
// with function WithTablesFromDB. Argument `options` can be set to "overwrite"
// and/or "view". Each option is its own slice entry.
func WithTable(tableName string, columns ddl.Columns, options ...string) (opt Option) {
	checkAutoIncrement := func(previousSetting uint8) uint8 {
		if previousSetting > 0 {
			return previousSetting
		}
		for _, o := range options {
			if strings.ToLower(o) == "view" {
				return 1 // nope
			}
		}
		for _, c := range columns {
			if c.IsAutoIncrement() {
				return 2 // yes
			}
		}
		return 1 // nope
	}

	opt.sortOrder = 10
	opt.fn = func(ts *Tables) error {
		isOverwrite := len(options) > 0 && options[0] == "overwrite"
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
		t.HasAutoIncrement = checkAutoIncrement(t.HasAutoIncrement)
		ts.Tables[tableName] = t
		return nil
	}
	return
}

// WithTablesFromDB queries the information_schema table and loads the column
// definition of the provided `tables` slice. It adds the tables to the `Tables`
// map. Once added a call to WithTableConfig can add additional configurations.
func WithTablesFromDB(ctx context.Context, db *dml.ConnPool, tables ...string) (opt Option) {
	opt.sortOrder = 1
	opt.fn = func(ts *Tables) error {
		tables, err := ddl.LoadColumns(ctx, db.DB, tables...)
		if err != nil {
			return errors.WithStack(err)
		}
		views, err := db.WithRawSQL("SELECT `TABLE_NAME` FROM `information_schema`.`VIEWS` WHERE `TABLE_SCHEMA`=DATABASE()").LoadStrings(ctx, nil)
		if err != nil {
			return errors.WithStack(err)
		}

		checkAutoIncrement := func(tblName string) uint8 {
			for _, v := range views {
				if v == tblName {
					return 1 // nope
				}
			}
			for _, c := range tables[tblName] {
				if c.IsAutoIncrement() {
					return 2 // yes
				}
			}
			return 1 // nope
		}

		for tblName := range tables {
			ts.Tables[tblName] = &table{
				TableName:        tblName,
				Columns:          tables[tblName],
				HasAutoIncrement: checkAutoIncrement(tblName),
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

// WithBuildTags adds your build tags to the file header. Each argument
// represents a build tag line.
func WithBuildTags(lines ...string) (opt Option) {
	opt.sortOrder = 112
	opt.fn = func(ts *Tables) error {
		ts.BuildTags = append(ts.BuildTags, lines...)
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
			"time",
			"github.com/corestoreio/pkg/sql/ddl",
			"github.com/corestoreio/pkg/sql/dml",
			"github.com/corestoreio/pkg/sql/dmltest",
			"github.com/corestoreio/pkg/util/assert",
			"github.com/corestoreio/pkg/util/pseudo",
		},
	}

	sort.Slice(opts, func(i, j int) bool {
		return opts[i].sortOrder < opts[j].sortOrder // ascending 0-9 sorting ;-)
	})

	for _, opt := range opts {
		if err := opt.fn(ts); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// ts.FuncMap["CustomCode"] = func(marker string) string { return ts.customCode[marker] }
	// ts.FuncMap["GoCamel"] = strs.ToGoCamelCase // net_http->NetHTTP entity_id->EntityID
	// ts.FuncMap["GoCamelMaybePrivate"] = func(s string) string { return s }

	// ts.FuncMap["GoFuncNull"] = func(c *ddl.Column) string { return ts.mySQLToGoDmlColumnMap(c, true) }
	// ts.FuncMap["GoFunc"] = func(c *ddl.Column) string { return ts.mySQLToGoDmlColumnMap(c, false) }
	// ts.FuncMap["IsFieldPublic"] = func(string) bool { return false }
	// ts.FuncMap["IsFieldPrivate"] = func(string) bool { return false }
	// ts.FuncMap["GoPrimitiveNull"] = ts.toGoPrimitiveFromNull
	// ts.FuncMap["SerializerType"] = func(c *ddl.Column) string {
	// 	pt := ts.toSerializerType(c, true)
	// 	if strings.IndexByte(pt, '/') > 0 { // slash identifies an import path
	// 		return "bytes"
	// 	}
	// 	return pt
	// }
	// ts.FuncMap["SerializerCustomType"] = func(c *ddl.Column) string {
	// 	pt := ts.toSerializerType(c, true)
	// 	var buf strings.Builder
	// 	if pt == "google.protobuf.Timestamp" {
	// 		fmt.Fprint(&buf, ",(gogoproto.stdtime)=true")
	// 	}
	// 	if pt == "bytes" {
	// 		return "" // bytes can be null
	// 	}
	// 	if c.IsNull() || strings.IndexByte(pt, '.') > 0 /*whenever it is a custom type like null. or google.proto.timestamp*/ {
	// 		// Indeed nullable Go Types must be not-nullable in HasSerializer because we
	// 		// have a non-pointer struct type which contains the field Valid.
	// 		// HasSerializer treats nullable fields as pointer fields, but that is
	// 		// ridiculous.
	// 		fmt.Fprint(&buf, ",(gogoproto.nullable)=false")
	// 	}
	// 	return buf.String()
	// }

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

	// if !ts.DisableFileHeader {
	// 	if err := ts.tpls.Funcs(ts.FuncMap).ExecuteTemplate(buf, ts.Serializer+"_10_header.go.tpl", ts); err != nil {
	// 		return errors.WriteFailed.New(err, "[dmlgen] For file header")
	// 	}
	// }
	//
	// tableFileName := ts.Serializer + "_20_message.go.tpl"
	// for _, tblName := range ts.sortedTableNames() {
	// 	t := ts.Tables[tblName] // must panic if table name not found
	// 	ts.FuncMap["IsFieldPublic"] = t.IsFieldPublic
	// 	if t.HasSerializer {
	// 		if err := t.writeTo(buf, ts.tpls.Lookup(tableFileName).Funcs(ts.FuncMap)); err != nil {
	// 			return errors.WriteFailed.New(err, "[dmlgen] For Table %q", t.TableName)
	// 		}
	// 	}
	// }
	_, err := buf.WriteTo(w)
	return err
}

// GenerateGo writes the Go source code into `w` and the test code into wTest.
func (ts *Tables) GenerateGo(wMain io.Writer, wTest io.Writer) error {
	// TODO this God-function got created during migration from .go.tpl files to inline source code. Needs refactoring.

	mainGen := codegen.NewGo(ts.Package)
	testGen := codegen.NewGo(ts.Package)

	mainGen.BuildTags = ts.BuildTags
	testGen.BuildTags = ts.BuildTags

	sortedTableNames := ts.sortedTableNames()
	if !ts.DisableTableSchemas { // Writes the table DDL function
		tables := make([]*table, len(ts.Tables))
		var tableNames []string
		var tableCreateStmt []string
		for i, tblname := range sortedTableNames {
			tables[i] = ts.Tables[tblname] // must panic if table name not found
			constName := `TableName` + strs.ToGoCamelCase(tblname)
			mainGen.AddConstString(constName, tblname)
			tableNames = append(tableNames, tblname)
			tableCreateStmt = append(tableCreateStmt, constName, `""`)
		}

		mainGen.C(`NewTables returns a goified version of the MySQL/MariaDB table schema for the tables: `,
			strings.Join(tableNames, ", "), ` Auto generated by dmlgen.`)
		mainGen.Pln(`func NewTables(ctx context.Context, opts ...ddl.TableOption) (tm *ddl.Tables,err error) {`)
		{
			mainGen.In()
			mainGen.Pln(`if tm, err = ddl.NewTables(`)
			{
				mainGen.In()
				mainGen.Pln(`append(opts, ddl.WithCreateTable(ctx,`, strings.Join(tableCreateStmt, ","), `))...,`)
				mainGen.Out()
			}
			mainGen.Out()
		}
		// gofmt will later remove the semicolons and formats it correctly
		mainGen.Pln(`); err != nil { return nil, errors.WithStack(err); }; return tm, nil; }`)

		testGen.Pln(`func TestNewTables(t *testing.T) {`)
		testGen.Pln(`db := dmltest.MustConnectDB(t)`)
		testGen.Pln(`defer dmltest.Close(t, db)`)

		if ts.TestSQLDumpGlobPath != "" {
			testGen.Pln(`defer dmltest.SQLDumpLoad(t,`, strconv.Quote(ts.TestSQLDumpGlobPath), `, &dmltest.SQLDumpOptions{
				SkipDBCleanup: true,
			}).Deferred()`)
		}

		testGen.Pln(`ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)`)
		testGen.Pln(`defer cancel()`)
		testGen.Pln(`tbls, err := NewTables(ctx, ddl.WithConnPool(db))`)
		testGen.Pln(`assert.NoError(t, err)`)

		testGen.Pln(`tblNames := tbls.Tables()`)
		testGen.Pln(`sort.Strings(tblNames)`)
		testGen.Pln(`assert.Exactly(t, `, fmt.Sprintf("%#v", tableNames), `, tblNames)`)

		testGen.Pln(`err = tbls.Validate(ctx)`)
		testGen.Pln(`assert.NoError(t, err)`)
		testGen.Pln(`var ps *pseudo.Service`)
		testGen.Pln(`ps = pseudo.MustNewService(0, &pseudo.Options{Lang: "de",FloatMaxDecimals:6},`)
		testGen.In()
		testGen.Pln(`pseudo.WithTagFakeFunc("website_id", func(maxLen int) (interface{}, error) {`)
		testGen.Pln(`    return 1, nil`)
		testGen.Pln(`}),`)

		testGen.Pln(`pseudo.WithTagFakeFunc("store_id", func(maxLen int) (interface{}, error) {`)
		testGen.Pln(`    return 1, nil`)
		testGen.Pln(`}),`)
		testGen.Pln(ts.customCode["pseudo.MustNewService.Option"])
		testGen.Out()
		testGen.Pln(`)`)

		for _, table := range tables {

			testGen.Pln(`t.Run("` + strs.ToGoCamelCase(table.TableName) + `_Entity", func(t *testing.T) {`)
			testGen.Pln(`tbl := tbls.MustTable(TableName`+strs.ToGoCamelCase(table.TableName), `)`)

			testGen.Pln(`entSELECT := tbl.SelectByPK("*")`)
			testGen.C(`WithArgs generates the cached SQL string with empty key "".`)
			testGen.Pln(`entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders()`)

			testGen.Pln(`entSELECT.WithCacheKey("select_10").Wheres.Reset()`)
			testGen.Pln(`_, _, err := entSELECT.Where(`)

			for _, c := range table.Columns {
				if c.IsPK() && c.IsAutoIncrement() {
					testGen.Pln(`dml.Column(`, strconv.Quote(c.Field), `).LessOrEqual().Int(10),`)
				}
			}

			testGen.Pln(`).ToSQL() // ToSQL generates the new cached SQL string with key select_10`)
			testGen.Pln(`assert.NoError(t, err)`)
			testGen.Pln(`entCol := New`+table.CollectionName(), `()`)

			if table.HasAutoIncrement < 2 {
				testGen.C(`this table/view does not support auto_increment`)
				testGen.Pln(`rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)`)
				testGen.Pln(`assert.NoError(t, err)`)
				testGen.Pln(`t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())`)
				testGen.Pln(`t.Logf("Collection load rowCount: %d", rowCount)`)
			} else {
				testGen.Pln(`entINSERT := tbl.Insert().BuildValues()`)
				testGen.Pln(`entINSERTStmtA := entINSERT.PrepareWithArgs(ctx)`)

				testGen.Pln(`for i := 0; i < 9; i++ {`)
				{
					testGen.In()
					testGen.Pln(`entIn := new(`, strs.ToGoCamelCase(table.TableName), `)`)
					testGen.Pln(`if err := ps.FakeData(entIn); err != nil {`)
					{
						testGen.In()
						testGen.Pln(`t.Errorf("IDX[%d]: %+v", i, err)`)
						testGen.Pln(`return`)
						testGen.Out()
					}
					testGen.Pln(`}`)

					testGen.Pln(`lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.` + strs.ToGoCamelCase(table.TableName) + `_Entity")(entINSERTStmtA.Record("", entIn).ExecContext(ctx))`)
					testGen.Pln(`entINSERTStmtA.Reset()`)

					testGen.Pln(`entOut := new(`, strs.ToGoCamelCase(table.TableName), `)`)
					testGen.Pln(`rowCount, err := entSELECTStmtA.Int64s(lID).Load(ctx, entOut)`)
					testGen.Pln(`assert.NoError(t, err)`)
					testGen.Pln(`assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)`)

					for _, c := range table.Columns {
						fn := table.GoCamelMaybePrivate(c.Field)
						switch {
						case c.IsString():
							testGen.Pln(`assert.ExactlyLength(t,`, c.CharMaxLength.Int64, `, `, `&entIn.`, fn, `,`, `&entOut.`, fn, `,`, `"IDX%d:`, fn, `should match", lID)`)
						case !c.IsSystemVersioned():
							testGen.Pln(`assert.Exactly(t, entIn.`, fn, `,`, `entOut.`, fn, `,`, `"IDX%d:`, fn, `should match", lID)`)
						default:
							testGen.C(`ignoring:`, c.Field)
						}
					}
					testGen.Out()
				}
				testGen.Pln(`}`) // endfor
				testGen.Pln(`dmltest.Close(t, entINSERTStmtA)`)

				testGen.Pln(`rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)`)
				testGen.Pln(`assert.NoError(t, err)`)
				testGen.Pln(`t.Logf("Collection load rowCount: %d", rowCount)`)

				testGen.Pln(`entINSERTStmtA = entINSERT.WithCacheKey("row_count_%d", len(entCol.Data)).Replace().SetRowCount(len(entCol.Data)).PrepareWithArgs(ctx)`)
				testGen.Pln(`lID := dmltest.CheckLastInsertID(t, "Error: `, table.CollectionName(), `")(entINSERTStmtA.Record("", entCol).ExecContext(ctx))`)
				testGen.Pln(`dmltest.Close(t, entINSERTStmtA)`)
				testGen.Pln(`t.Logf("Last insert ID into: %d", lID)`)
				testGen.Pln(`t.Logf("INSERT queries: %#v", entINSERT.CachedQueries())`)
				testGen.Pln(`t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())`)
			}

			testGen.Pln(`})`)

		} // end for tables
		testGen.Pln(`}`) // end TestNewTables
	}

	// deal with random map to guarantee the persistent code generation.
	for _, tblname := range sortedTableNames {
		t, ok := ts.Tables[tblname]
		if !ok || t == nil {
			return errors.NotFound.Newf("[dmlgen] Table %q not found", tblname)
		}

		mainGen.C(t.EntityName(), ` represents a single row for DB table `, t.TableName, `.`)
		mainGen.C(`Auto generated. `)
		if t.Comment != "" {
			mainGen.C(t.Comment)
		}
		if t.HasEasyJsonMarshaler {
			mainGen.Pln(`//easyjson:json`)
		}

		// Generate table structs
		mainGen.Pln(`type `, t.EntityName(), ` struct {`)
		{
			mainGen.In()
			for _, c := range t.Columns {
				structTag := ""
				if c.StructTag != "" {
					structTag += "`" + c.StructTag + "`"
				}
				mainGen.Pln(t.GoCamelMaybePrivate(c.Field), ts.GoTypeNull(c), structTag, c.GoComment())
			}
			for _, c := range t.ReferencedCollections {
				mainGen.Pln(c)
			}
			mainGen.Out()
		}
		mainGen.Pln(`}`)

		// Generates the Getter/Setter for private fields
		for _, c := range t.Columns {
			if t.IsFieldPrivate(c.Field) {
				mainGen.C(`Set`, strs.ToGoCamelCase(c.Field), ` sets the data for a private and security sensitive field.`)
				mainGen.Pln(`func (e *`, t.EntityName(), `) Set`+strs.ToGoCamelCase(c.Field), `(d `, ts.GoTypeNull(c), `) *`, t.EntityName(), ` {`)
				{
					mainGen.In()
					mainGen.Pln(`e.`, t.GoCamelMaybePrivate(c.Field), ` = d`)
					mainGen.Pln(`return e`)
					mainGen.Out()
				}
				mainGen.Pln(`}`)

				mainGen.C(`Get`, strs.ToGoCamelCase(c.Field), ` returns the data from a private and security sensitive field.`)
				mainGen.Pln(`func (e *`, t.EntityName(), `) Get`+strs.ToGoCamelCase(c.Field), `() `, ts.GoTypeNull(c), `{`)
				{
					mainGen.In()
					mainGen.Pln(`return e.`, t.GoCamelMaybePrivate(c.Field))
					mainGen.Out()
				}
				mainGen.Pln(`}`)
			}
		}

		mainGen.Pln(`// Empty empties all the fields of the current object. Also known as Reset.`)
		mainGen.Pln(`func (e *`, t.EntityName(), `) Empty() *`, t.EntityName(), ` { *e = `, t.EntityName(), `{}; return e }`)

		mainGen.C(t.CollectionName(), `represents a collection type for DB table`, t.TableName)
		mainGen.C(`Not thread safe. Auto generated.`)
		if t.Comment != "" {
			mainGen.C(t.Comment)
		}
		if t.HasEasyJsonMarshaler {
			mainGen.Pln(`//easyjson:json`) // do not use C() because it adds a whitespace between "//" and "e"
		}
		mainGen.Pln(`type `, t.CollectionName(), ` struct {`)
		{
			mainGen.In()
			mainGen.Pln(`Data []*`, t.EntityName(), codegen.EncloseBT(`json:"data,omitempty"`))
			mainGen.Pln(`BeforeMapColumns	func(uint64, *`, t.EntityName(), `) error`, codegen.EncloseBT(`json:"-"`))
			mainGen.Pln(`AfterMapColumns 	func(uint64, *`, t.EntityName(), `) error `, codegen.EncloseBT(`json:"-"`))
			mainGen.Out()
		}
		mainGen.Pln(`}`)

		mainGen.C(`New`+t.CollectionName(), ` creates a new initialized collection. Auto generated.`)
		// TODO(idea): use a global pool which can register for each type the
		// before/after mapcolumn function so that the dev does not need to
		// assign each time. think if it's worth such a pattern.
		mainGen.Pln(`func New`+t.CollectionName(), `() *`, t.CollectionName(), ` {`)
		{
			mainGen.In()
			mainGen.Pln(`return &`, t.CollectionName(), `{`)
			{
				mainGen.In()
				mainGen.Pln(`Data: make([]*`, t.EntityName(), `, 0, 5),`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)

		// Generate functions to access SQL
		mainGen.C(`AssignLastInsertID updates the increment ID field with the last inserted ID from an INSERT operation.`,
			`Implements dml.InsertIDAssigner. Auto generated.`)
		mainGen.Pln(`func (e *`, t.EntityName(), `) AssignLastInsertID(id int64) {`)
		{
			mainGen.In()
			t.Columns.Each(func(c *ddl.Column) {
				if c.IsPK() && c.IsAutoIncrement() {
					mainGen.Pln(`e.`, t.GoCamelMaybePrivate(c.Field), ` = `, ts.GoType(c), `(id)`)
				}
			})
			mainGen.Out()
		}
		mainGen.Pln(`}`)

		mainGen.C(`MapColumns implements interface ColumnMapper only partially. Auto generated.`)
		mainGen.Pln(`func (e *`, t.EntityName(), `) MapColumns(cm *dml.ColumnMap) error {`)
		{
			mainGen.In()
			mainGen.Pln(`if cm.Mode() == dml.ColumnMapEntityReadAll {`)
			{
				mainGen.In()
				mainGen.P(`return cm`)
				t.Columns.Each(func(c *ddl.Column) {
					mainGen.P(`.`, ts.GoFuncNull(c), `(&e.`, t.GoCamelMaybePrivate(c.Field), `)`)
				})
				mainGen.Pln(`.Err()`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
			mainGen.Pln(`for cm.Next() {`)
			{
				mainGen.In()
				mainGen.Pln(`switch c := cm.Column(); c {`)
				{
					mainGen.In()
					t.Columns.Each(func(c *ddl.Column) {
						// mainGen.Pln(`case "` + c.Field + `"{{range .Aliases}},"{{.}}"{{end}}:`)
						mainGen.P(`case`, strconv.Quote(c.Field))
						for _, a := range c.Aliases {
							mainGen.P(`,`, strconv.Quote(a))
						}
						mainGen.Pln(`:`)
						mainGen.Pln(`cm.`, ts.GoFuncNull(c), `(&e.`, t.GoCamelMaybePrivate(c.Field), `)`)
					})
					mainGen.Pln(`default:`)
					mainGen.Pln(`return errors.NotFound.Newf("[`+ts.Package+`]`, t.EntityName(), `Column %q not found", c)`)
					mainGen.Out()
				}
				mainGen.Pln(`}`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
			mainGen.Pln(`return errors.WithStack(cm.Err())`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)

		mainGen.C(`AssignLastInsertID traverses through the slice and sets a decrementing new ID to each entity.`)
		mainGen.Pln(`func (cc *`, t.CollectionName(), `) AssignLastInsertID(id int64) {`)
		{
			mainGen.In()
			mainGen.Pln(`var j int64`)
			mainGen.Pln(`for i := len(cc.Data) - 1; i >= 0; i-- {`)
			{
				mainGen.In()
				mainGen.Pln(`cc.Data[i].AssignLastInsertID(id - j)`)
				mainGen.Pln(`j++`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)

		mainGen.Pln(`func (cc *`, t.CollectionName(), `) scanColumns(cm *dml.ColumnMap,e *`, t.EntityName(), `, idx uint64) error {
			if cc.BeforeMapColumns != nil {
				if err := cc.BeforeMapColumns(idx, e); err != nil {
					return errors.WithStack(err)
				}
			}
			if err := e.MapColumns(cm); err != nil {
				return errors.WithStack(err)
			}
			if cc.AfterMapColumns != nil {
				if err := cc.AfterMapColumns(idx, e); err != nil {
					return errors.WithStack(err)
				}
			}
			return nil
		}`)

		mainGen.C(`MapColumns implements dml.ColumnMapper interface. Auto generated.`)
		mainGen.Pln(`func (cc *`, t.CollectionName(), `) MapColumns(cm *dml.ColumnMap) error {`)
		{
			mainGen.Pln(`switch m := cm.Mode(); m {
						case dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
							for i, e := range cc.Data {
								if err := cc.scanColumns(cm, e, uint64(i)); err != nil {
									return errors.WithStack(err)
								}
							}`)

			mainGen.Pln(`case dml.ColumnMapScan:
							if cm.Count == 0 {
								cc.Data = cc.Data[:0]
							}
							e := new(`, t.EntityName(), `)
							if err := cc.scanColumns(cm, e, cm.Count); err != nil {
								return errors.WithStack(err)
							}
							cc.Data = append(cc.Data, e)`)

			mainGen.Pln(`case dml.ColumnMapCollectionReadSet:
							for cm.Next() {
								switch c := cm.Column(); c {`)

			t.Columns.UniqueColumns().Each(func(c *ddl.Column) {
				if !c.IsFloat() {
					mainGen.P(`case`, strconv.Quote(c.Field))
					for _, a := range c.Aliases {
						mainGen.P(`,`, strconv.Quote(a))
					}
					mainGen.Pln(`:`)
					mainGen.Pln(`cm = cm.`, ts.GoFuncNull(c)+`s(cc.`, strs.ToGoCamelCase(c.Field)+`s()...)`)
				}
			})
			mainGen.Pln(`default:
				return errors.NotFound.Newf("[`+t.Package+`]`, t.CollectionName(), `Column %q not found", c)
			}
		} // end for cm.Next

	default:
		return errors.NotSupported.Newf("[`+t.Package+`] Unknown Mode: %q", string(m))
	}
	return cm.Err()`)
		}
		mainGen.Pln(`}`) // end func MapColumns

		// Generates functions to return all data as a slice from unique/primary
		// columns.
		for _, c := range t.Columns.UniqueColumns() {
			gtn := ts.GoTypeNull(c)
			goCamel := strs.ToGoCamelCase(c.Field)
			mainGen.C(goCamel + `s returns a slice with the data or appends it to a slice.`)
			mainGen.C(`Auto generated.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) `, goCamel+`s(ret ...`+gtn, `) []`+gtn, ` {`)
			{
				mainGen.In()
				mainGen.Pln(`if ret == nil {`)
				{
					mainGen.In()
					mainGen.Pln(`ret = make([]`+gtn, `, 0, len(cc.Data))`)
					mainGen.Out()
				}
				mainGen.Pln(`}`)
				mainGen.Pln(`for _, e := range cc.Data {`)
				{
					mainGen.In()
					mainGen.Pln(`ret = append(ret, e.`+goCamel, `)`)
					mainGen.Out()
				}
				mainGen.Pln(`}`)
				mainGen.Pln(`return ret`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
		}

		// Generates functions to return data with removed duplicates from any
		// column which has set the flag Uniquified.
		for _, c := range t.Columns.UniquifiedColumns() {
			goType := ts.GoType(c)
			goCamel := strs.ToGoCamelCase(c.Field)

			mainGen.C(goCamel+`s belongs to the column`, strconv.Quote(c.Field), `and returns a slice or appends to a slice only`,
				`unique values of that column. The values will be filtered internally in a Go map. No DB query gets`,
				`executed. Auto generated.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) Unique`+goCamel+`s(ret ...`, goType, `) []`, goType, ` {`)
			{
				mainGen.In()
				mainGen.Pln(`if ret == nil {
					ret = make([]`, goType, `, 0, len(cc.Data))
				}`)

				// TODO: a reusable map and use different algorithms depending on
				// the size of the cc.Data slice. Sometimes a for/for loop runs
				// faster than a map.
				goPrimNull := ts.toGoPrimitiveFromNull(c)
				mainGen.Pln(`dupCheck := make(map[`, goType, `]bool, len(cc.Data))`)
				mainGen.Pln(`for _, e := range cc.Data {`)
				{
					mainGen.In()
					mainGen.Pln(`if !dupCheck[e.`+goPrimNull, `] {`)
					{
						mainGen.In()
						mainGen.Pln(`ret = append(ret, e.`, goPrimNull, `)`)
						mainGen.Pln(`dupCheck[e.`+goPrimNull, `] = true`)
						mainGen.Out()
					}
					mainGen.Pln(`}`)
					mainGen.Out()
				}
				mainGen.Pln(`}`)
				mainGen.Pln(`return ret`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
		}

		if !t.DisableCollectionMethods {
			mainGen.C(`Filter filters the current slice by predicate f without memory allocation. Auto generated via dmlgen.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) Filter(f func(*`, t.EntityName(), `) bool) *`, t.CollectionName(), ` {`)
			{
				mainGen.In()
				mainGen.Pln(`b,i := cc.Data[:0],0`)
				mainGen.Pln(`for _, e := range cc.Data {`)
				{
					mainGen.In()
					mainGen.Pln(`if f(e) {`)
					{
						mainGen.Pln(`b = append(b, e)`)
						mainGen.Pln(`cc.Data[i] = nil // this avoids the memory leak`)
					}
					mainGen.Pln(`}`) // endif
					mainGen.Pln(`i++`)
				}
				mainGen.Out()
				mainGen.Pln(`}`) // for loop
				mainGen.Pln(`cc.Data = b`)
				mainGen.Pln(`return cc`)
				mainGen.Out()
			}
			mainGen.Pln(`}`) // function

			mainGen.C(`Each will run function f on all items in []*`, t.EntityName(), `. Auto generated via dmlgen.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) Each(f func(*`, t.EntityName(), `)) *`, t.CollectionName(), ` {`)
			{
				mainGen.Pln(`for i := range cc.Data {`)
				{
					mainGen.Pln(`f(cc.Data[i])`)
				}
				mainGen.Pln(`}`)
				mainGen.Pln(`return cc`)
			}
			mainGen.Pln(`}`)

			mainGen.C(`Cut will remove items i through j-1. Auto generated via dmlgen.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) Cut(i, j int) *`, t.CollectionName(), ` {`)
			{
				mainGen.In()
				mainGen.Pln(`z := cc.Data // copy slice header`)
				mainGen.Pln(`copy(z[i:], z[j:])`)
				mainGen.Pln(`for k, n := len(z)-j+i, len(z); k < n; k++ {`)
				{
					mainGen.In()
					mainGen.Pln(`z[k] = nil // this avoids the memory leak`)
					mainGen.Out()
				}
				mainGen.Pln(`}`)
				mainGen.Pln(`z = z[:len(z)-j+i]`)
				mainGen.Pln(`cc.Data = z`)
				mainGen.Pln(`return cc`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)

			mainGen.C(`Swap will satisfy the sort.Interface. Auto generated via dmlgen.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) Swap(i, j int) { cc.Data[i], cc.Data[j] = cc.Data[j], cc.Data[i] }`)

			mainGen.C(`Delete will remove an item from the slice. Auto generated via dmlgen.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) Delete(i int) *`, t.CollectionName(), ` {`)
			{
				mainGen.Pln(`z := cc.Data // copy the slice header`)
				mainGen.Pln(`end := len(z) - 1`)
				mainGen.Pln(`cc.Swap(i, end)`)
				mainGen.Pln(`copy(z[i:], z[i+1:])`)
				mainGen.Pln(`z[end] = nil // this should avoid the memory leak`)
				mainGen.Pln(`z = z[:end]`)
				mainGen.Pln(`cc.Data = z`)
				mainGen.Pln(`return cc`)
			}
			mainGen.Pln(`}`)

			mainGen.C(`Insert will place a new item at position i. Auto generated via dmlgen.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) Insert(n *`, t.EntityName(), `, i int) *`, t.CollectionName(), ` {`)
			{
				mainGen.Pln(`z := cc.Data // copy the slice header`)
				mainGen.Pln(`z = append(z, &`+t.EntityName(), `{})`)
				mainGen.Pln(`copy(z[i+1:], z[i:])`)
				mainGen.Pln(`z[i] = n`)
				mainGen.Pln(`cc.Data = z`)
				mainGen.Pln(`return cc`)
			}
			mainGen.Pln(`}`)

			mainGen.C(`Append will add a new item at the end of *`, t.CollectionName(), `. Auto generated via dmlgen.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) Append(n ...*`, t.EntityName(), `) *`, t.CollectionName(), ` {`)
			{
				mainGen.Pln(`cc.Data = append(cc.Data, n...)`)
				mainGen.Pln(`return cc`)
			}
			mainGen.Pln(`}`)
		}

		if t.HasBinaryMarshaler {
			mainGen.C(`UnmarshalBinary implements encoding.BinaryUnmarshaler.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) UnmarshalBinary(data []byte) error {`)
			{
				mainGen.Pln(`return cc.Unmarshal(data) // Implemented via github.com/gogo/protobuf`)
			}
			mainGen.Pln(`}`)

			mainGen.C(`MarshalBinary implements encoding.BinaryMarshaler.`)
			mainGen.Pln(`func (cc *`, t.CollectionName(), `) MarshalBinary() (data []byte, err error) {`)
			{
				mainGen.Pln(`return cc.Marshal()  // Implemented via github.com/gogo/protobuf`)
			}
			mainGen.Pln(`}`)
		}
		if ts.lastError != nil {
			return ts.lastError
		}
	}

	// now figure out all used package names in the buffer.

	pkgs, err := ts.findUsedPackages(mainGen.Bytes())
	if err != nil {
		wMain.Write(mainGen.Bytes()) // write for debug reasons
		return errors.WithStack(err)
	}
	mainGen.AddImports(pkgs...)

	testGen.AddImports(ts.ImportPathsTesting...)

	if err := mainGen.GenerateFile(wMain); err != nil {
		return err
	}
	if err := testGen.GenerateFile(wTest); err != nil {
		return err
	}
	return nil
}

func (ts *Tables) GoTypeNull(c *ddl.Column) string { return ts.mySQLToGoType(c, true) }
func (ts *Tables) GoType(c *ddl.Column) string     { return ts.mySQLToGoType(c, false) }
func (ts *Tables) GoFuncNull(c *ddl.Column) string { return ts.mySQLToGoDmlColumnMap(c, true) }

// table writes one database table into Go source code.
type table struct {
	Package          string      // Name of the package
	TableName        string      // Name of the table
	Comment          string      // Comment above the struct type declaration
	Columns          ddl.Columns // all columns of a table
	HasAutoIncrement uint8       // 0=nil,1=false (has NO auto increment),2=true has auto increment
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

func (t *table) IsFieldPublic(dbColumnName string) bool {
	return t.privateFields == nil || !t.privateFields[dbColumnName]
}

func (t *table) IsFieldPrivate(dbColumnName string) bool {
	return t.privateFields != nil && t.privateFields[dbColumnName]
}

func (t *table) GoCamelMaybePrivate(fieldName string) string {
	su := strs.ToGoCamelCase(fieldName)
	if t.IsFieldPublic(fieldName) {
		return su
	}
	sr := []rune(su)
	sr[0] = unicode.ToLower(sr[0])
	return string(sr)
}

func (t *table) CollectionName() string {
	return strs.ToGoCamelCase(t.TableName) + "Collection"
}

func (t *table) EntityName() string {
	return strs.ToGoCamelCase(t.TableName)
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
