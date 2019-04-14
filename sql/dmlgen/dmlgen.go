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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/util/codegen"
	"github.com/corestoreio/pkg/util/slices"
	"github.com/corestoreio/pkg/util/strs"
	"github.com/mailru/easyjson/bootstrap"
	"github.com/mailru/easyjson/parser"
)

// Initial idea and prototyping for code generation.
// TO DO DML gen must take care of the types in myreplicator.RowsEvent.decodeValue (not possible)

// TODO generate a hidden type which contains the original data to detect
//  changes and store only changed fields. Investigate other implementations.
//  With the e.g. new field "originalData" pointing to the struct we create a
//  recursion, also we need to track which field needs to be updated after a
//  change. Alternative implementation: a byte slice using protobuf encoding
//  to store the current data and when saving to DB happens decode the data and
//  compare the fields, then return the change field names, except autoinc
//  fields.

// Generator can generated Go source for for database tables once correctly
// configured.
type Generator struct {
	Package            string // Name of the package
	PackageImportPath  string // Name of the package
	BuildTags          []string
	ImportPaths        []string
	ImportPathsTesting []string
	// Tables uses the table name as map key and the table description as value.
	Tables map[string]*Table
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
	defaultTableConfig  TableConfig
	// customCode injects custom code to manipulate testing and other generate
	// code blocks. A few implementations now but more can be added later.
	customCode map[string]func(*Generator, *Table, io.Writer)

	kcu    map[string]ddl.KeyColumnUsageCollection
	kcuRev map[string]ddl.KeyColumnUsageCollection // rev = reversed relationship to find OneToMany
	krs    *ddl.KeyRelationShips
	// "mainTable.mainColumn" : "referencedTable.referencedColumn"
	// or skips the reversed relationship
	// "referencedTable.referencedColumn": "mainTable.mainColumn"
	krsSkip map[string]string
}

// Option represents a sortable option for the NewGenerator function. Each option
// function can be applied in a mixed order.
type Option struct {
	// sortOrder specifies the precedence of an option.
	sortOrder int
	fn        func(*Generator) error
}

// TableConfig used in conjunction with WithTableConfig and
// WithTableConfigDefault to apply different configurations for a generated
// struct and its struct collection.
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
	// FeaturesInclude if set includes only those features, otherwise
	// everything. Some features can only be included on Default level and not
	// on a per table level.
	FeaturesInclude FeatureToggle
	// FeaturesExclude excludes features while field FeaturesInclude is empty.
	// Some features can only be excluded on Default level and not on a per
	// table level.
	FeaturesExclude FeatureToggle
	// FieldMapFn can map a dbIdentifier (database identifier) of the current
	// table to a new name. dbIdentifier is in most cases the column name and in
	// cases of foreign keys, it is the table name.
	FieldMapFn func(dbIdentifier string) (newName string)
	lastErr    error
}

func (to *TableConfig) applyEncoders(t *Table) {
	for i := 0; i < len(to.Encoders) && to.lastErr == nil; i++ {
		switch enc := to.Encoders[i]; enc {
		case "json":
			t.HasJSONMarshaler = true // for now does nothing
		case "easyjson":
			t.HasEasyJSONMarshaler = true
		case "binary":
			t.HasBinaryMarshaler = true
		case "protobuf", "fbs":
			t.HasSerializer = true // for now leave it in. maybe later PB gets added to the struct tags.
		default:
			to.lastErr = errors.NotSupported.Newf("[dmlgen] WithTableConfig: Table %q Encoder %q not supported", t.TableName, enc)
		}
	}
}

func (to *TableConfig) applyStructTags(t *Table) {
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

func (to *TableConfig) applyCustomStructTags(t *Table) {
	for i := 0; i < len(to.CustomStructTags) && to.lastErr == nil; i += 2 {
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

func (to *TableConfig) applyPrivateFields(t *Table) {
	if len(to.PrivateFields) > 0 && t.privateFields == nil {
		t.privateFields = make(map[string]bool)
	}
	for _, pf := range to.PrivateFields {
		t.privateFields[pf] = true
	}
}

func (to *TableConfig) applyComments(t *Table) {
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

func (to *TableConfig) applyColumnAliases(t *Table) {
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
func (to *TableConfig) applyUniquifiedColumns(t *Table) {
	for i := 0; i < len(to.UniquifiedColumns) && to.lastErr == nil; i++ {
		cn := to.UniquifiedColumns[i]
		found := false
		for _, c := range t.Columns {
			if c.Field == cn && !c.IsBlobDataType() {
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

// WithTableConfigDefault sets for all tables the same configuration but can be
// overwritten on a per table basis with function WithTableConfig.
func WithTableConfigDefault(opt TableConfig) (o Option) {
	o.sortOrder = 149 // must be applied before WithTableConfig
	o.fn = func(ts *Generator) (err error) {
		ts.defaultTableConfig = opt
		return opt.lastErr
	}
	return o
}

func defaultFieldMapFn(s string) string {
	return strs.ToGoCamelCase(s)
}

// WithTableConfig applies options to an existing table, identified by the table
// name used as map key. Options are custom struct or different encoders.
// Returns a not-found error if the table cannot be found in the `Generator` map.
func WithTableConfig(tableName string, opt *TableConfig) (o Option) {
	// Panic as early as possible.
	if len(opt.CustomStructTags)%2 == 1 {
		panic(errors.Fatal.Newf("[dmlgen] WithTableConfig: Table %q option CustomStructTags must be a balanced slice.", tableName))
	}
	o.sortOrder = 150
	o.fn = func(ts *Generator) (err error) {
		t, ok := ts.Tables[tableName]
		if t == nil || !ok {
			return errors.NotFound.Newf("[dmlgen] WithTableConfig: Table %q not found.", tableName)
		}
		opt.applyEncoders(t)
		opt.applyStructTags(t)
		opt.applyCustomStructTags(t)
		opt.applyPrivateFields(t)
		opt.applyComments(t)
		opt.applyColumnAliases(t)
		opt.applyUniquifiedColumns(t)
		t.featuresInclude = opt.FeaturesInclude
		t.featuresExclude = opt.FeaturesExclude
		t.fieldMapFn = opt.FieldMapFn
		if t.fieldMapFn == nil {
			t.fieldMapFn = defaultFieldMapFn
		}
		return opt.lastErr
	}
	return o
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
	opt.fn = func(ts *Generator) error {

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
	return opt
}

// WithForeignKeyRelationships analyzes the foreign keys which points to a table
// and adds them as a struct field name. For example:
// customer_address_entity.parent_id is a foreign key to
// customer_entity.entity_id hence the generated struct CustomerEntity has a new
// field which gets named CustomerAddressEntityCollection, pointing to type
// CustomerAddressEntityCollection. skipRelationships must be a balanced slice
// in the notation of "table1.column1","table2.column2". For example:
// 		"customer_entity.store_id", "store.store_id"
// which means that the struct CustomerEntity won't have a field to the Store
// struct (1:1 relationship). The reverse case can also be added
// 		"store.store_id", "customer_entity.store_id"
// which means that the Store struct won't have a field pointing to the
// CustomerEntityCollection (1:M relationship).
func WithForeignKeyRelationships(ctx context.Context, db dml.Querier, skipRelationships ...string) (opt Option) {
	opt.sortOrder = 210 // must run at the end or where the end is near ;-)
	opt.fn = func(ts *Generator) (err error) {

		if len(skipRelationships)%2 == 1 {
			return errors.Fatal.Newf("[dmlgen] skipRelationships must be balanced slice. Read the doc.")
		}

		ts.kcu, err = ddl.LoadKeyColumnUsage(ctx, db, ts.sortedTableNames()...)
		if err != nil {
			return errors.WithStack(err)
		}
		ts.kcuRev = ddl.ReverseKeyColumnUsage(ts.kcu)

		ts.krs, err = ddl.GenerateKeyRelationships(ctx, db, ts.kcu)
		if err != nil {
			return errors.WithStack(err)
		}

		ts.krsSkip = make(map[string]string, len(skipRelationships)/2)
		for i := 0; i < len(skipRelationships); i += 2 {
			ts.krsSkip[skipRelationships[i]] = skipRelationships[i+1]
		}

		var buf bytes.Buffer
		ts.krs.Debug(&buf)
		println(buf.String())

		return nil
	}
	return opt
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
	opt.fn = func(ts *Generator) error {
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
			t = &Table{
				TableName: tableName,
				Columns:   columns,
			}
		}
		t.HasAutoIncrement = checkAutoIncrement(t.HasAutoIncrement)
		ts.Tables[tableName] = t
		return nil
	}
	return opt
}

// WithTablesFromDB queries the information_schema table and loads the column
// definition of the provided `tables` slice. It adds the tables to the `Generator`
// map. Once added a call to WithTableConfig can add additional configurations.
func WithTablesFromDB(ctx context.Context, db *dml.ConnPool, tables ...string) (opt Option) {
	opt.sortOrder = 1
	opt.fn = func(ts *Generator) error {
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
			ts.Tables[tblName] = &Table{
				TableName:        tblName,
				Columns:          tables[tblName],
				HasAutoIncrement: checkAutoIncrement(tblName),
			}
		}
		return nil
	}
	return opt
}

// WithProtobuf enables protocol buffers as a serialization method. Argument
// headerOptions is optional.
func WithProtobuf(headerOptions ...string) (opt Option) {
	opt.sortOrder = 110
	opt.fn = func(ts *Generator) error {
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
	return opt
}

// WithFlatbuffers enables flatbuffers (FBS) as a serialization method. Argument
// headerOptions is optional.
func WithFlatbuffers(headerOptions ...string) (opt Option) {
	opt.sortOrder = 111
	opt.fn = func(ts *Generator) error {
		ts.Serializer = "fbs"
		// if len(headerOptions) == 0 {
		// TODO find sane defaults
		// ts.SerializerHeaderOptions = []string{}
		// }
		return nil
	}
	return opt
}

// WithBuildTags adds your build tags to the file header. Each argument
// represents a build tag line.
func WithBuildTags(lines ...string) (opt Option) {
	opt.sortOrder = 112
	opt.fn = func(ts *Generator) error {
		ts.BuildTags = append(ts.BuildTags, lines...)
		return nil
	}
	return opt
}

// WithCustomCode inserts at the marker position your custom Go code. For
// available markers search the .go.tpl files for the function call
// `CustomCode`. An example got written in TestGenerate_Tables_Protobuf_Json. If
// the marker does not exists or has a typo, no error gets reported and no code
// gets written.
func WithCustomCode(marker, code string) (opt Option) {
	opt.sortOrder = 112
	opt.fn = func(ts *Generator) error {
		if ts.customCode == nil {
			ts.customCode = make(map[string]func(*Generator, *Table, io.Writer))
		}
		ts.customCode[marker] = func(_ *Generator, _ *Table, w io.Writer) {
			w.Write([]byte(code))
		}
		return nil
	}
	return opt
}

// WithCustomCodeFunc same as WithCustomCode but allows access to meta data. The
// func fn takes as first argument the main Generator where access to package
// global configuration is possible. If the scope of the marker is within a
// table, then argument t gets set, otherwise it is nil. The output must be
// written to w.
func WithCustomCodeFunc(marker string, fn func(g *Generator, t *Table, w io.Writer)) (opt Option) {
	opt.sortOrder = 113
	opt.fn = func(ts *Generator) error {
		if ts.customCode == nil {
			ts.customCode = make(map[string]func(*Generator, *Table, io.Writer))
		}
		ts.customCode[marker] = fn
		return nil
	}
	return opt
}

func (ts *Generator) sortedTableNames() []string {
	sortedKeys := make(slices.String, 0, len(ts.Tables))
	for k := range ts.Tables {
		sortedKeys = append(sortedKeys, k)
	}
	sortedKeys.Sort()
	return sortedKeys
}

// NewGenerator creates a new instance of the SQL table code generator. The order
// of the applied options does not matter as they are getting sorted internally.
func NewGenerator(packageImportPath string, opts ...Option) (*Generator, error) {
	_, pkg := filepath.Split(packageImportPath)
	ts := &Generator{
		Tables:            make(map[string]*Table),
		Package:           pkg,
		PackageImportPath: packageImportPath,
		ImportPaths: []string{
			"context",
			"database/sql",
			"encoding/json",
			"fmt",
			"io",
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

	for _, t := range ts.Tables {
		t.Package = ts.Package
	}
	return ts, nil
}

func (ts *Generator) skipRelationship(table1, column1, table2, column2 string) bool {
	// println("skip", "key", table1+"."+column1, ": ", ts.krsSkip[table1+"."+column1], "==", table2+"."+column2)
	return ts.krsSkip[table1+"."+column1] == table2+"."+column2
}

func (ts *Generator) hasFeature(tableInclude, tableExclude, feature FeatureToggle) bool {
	if tableInclude == 0 {
		tableInclude = ts.defaultTableConfig.FeaturesInclude
	}
	if tableExclude == 0 {
		tableExclude = ts.defaultTableConfig.FeaturesExclude
	}

	switch {
	case tableInclude == 0 && tableExclude == 0:
		return true
	case tableInclude > 0 && (tableInclude&feature) != 0 && tableExclude == 0:
		return true
	case tableInclude == 0 && tableExclude > 0 && (tableExclude&feature) != 0:
		return false
	case tableInclude == 0 && tableExclude > 0 && (tableExclude&feature) == 0:
		return true
	default:
		return false
	}
}

// findUsedPackages checks for needed packages which we must import.
func (ts *Generator) findUsedPackages(file []byte) ([]string, error) {

	af, err := goparser.ParseFile(token.NewFileSet(), "cs_virtual_file.go", append([]byte("package temporarily_main\n\n"), file...), 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	idents := map[string]struct{}{}
	ast.Inspect(af, func(n ast.Node) bool {
		if nt, ok := n.(*ast.Ident); ok {
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

// SerializerType converts the column type to the supported type of the current
// serializer. For now supports only protobuf.
func (ts *Generator) SerializerType(c *ddl.Column) string {
	pt := ts.toSerializerType(c, true)
	if strings.IndexByte(pt, '/') > 0 { // slash identifies an import path
		return "bytes"
	}
	return pt
}

// SerializerCustomType switches the default type from function SerializerType
// to the new type. For now supports only protobuf.
func (ts *Generator) SerializerCustomType(c *ddl.Column) string {
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

// GenerateSerializer writes the protocol buffer specifications into `w` and its test
// sources into wTest, if there are any tests.
func (ts *Generator) GenerateSerializer(wMain, wTest io.Writer) error {
	switch ts.Serializer {
	case "protobuf":
		if err := ts.generateProto(wMain); err != nil {
			return errors.WithStack(err)
		}
	case "fbs":

	case "", "default", "none":
		return nil // do nothing
	default:
		return errors.NotAcceptable.Newf("[dmlgen] Serializer %q not supported.", ts.Serializer)
	}

	return nil
}

func (ts *Generator) generateProto(w io.Writer) error {
	cg := codegen.NewProto(ts.Package)
	cg.Pln(`import "github.com/gogo/protobuf/gogoproto/gogo.proto";`)
	cg.Pln(`import "google/protobuf/timestamp.proto";`)
	cg.Pln(`import "github.com/corestoreio/pkg/storage/null/null.proto";`)
	cg.Pln(`option go_package = "` + ts.Package + `";`)
	for _, o := range ts.SerializerHeaderOptions {
		cg.Pln(`option ` + o + `;`)
	}

	for _, tblname := range ts.sortedTableNames() {
		t := ts.Tables[tblname] // must panic if table name not found

		cg.C(t.EntityName(), `represents a single row for DB table`, t.TableName, `. Auto generated.`)
		cg.Pln(`message`, t.EntityName(), `{`)
		{
			cg.In()
			t.Columns.Each(func(c *ddl.Column) {
				if t.IsFieldPublic(c.Field) {
					cg.Pln(ts.SerializerType(c), c.Field+`=`, c.Pos, `[(gogoproto.customname)=`+strconv.Quote(strs.ToGoCamelCase(c.Field)), ts.SerializerCustomType(c)+`];`)
				}
			})
			cg.Out()
		}
		cg.Pln(`}`)

		cg.C(t.CollectionName(), `represents multiple rows for DB table`, t.TableName, `. Auto generated.`)
		cg.Pln(`message`, t.CollectionName(), `{`)
		{
			cg.In()
			cg.Pln(`repeated`, t.EntityName(), `Data = 1;`)
			cg.Out()
		}
		cg.Pln(`}`)
	}
	return cg.GenerateFile(w)
}

// GenerateGo writes the Go source code into `w` and the test code into wTest.
func (ts *Generator) GenerateGo(wMain, wTest io.Writer) error {

	mainGen := codegen.NewGo(ts.Package)
	testGen := codegen.NewGo(ts.Package)

	mainGen.BuildTags = ts.BuildTags
	testGen.BuildTags = ts.BuildTags

	tables := make([]*Table, len(ts.Tables))
	for i, tblname := range ts.sortedTableNames() {
		tables[i] = ts.Tables[tblname] // must panic if table name not found
	}

	ts.fnDBNewTables(mainGen)
	ts.fnTestMainOther(testGen, tables)
	ts.fnTestMainDB(testGen, tables)

	// deal with random map to guarantee the persistent code generation.
	for _, t := range tables {

		t.entityStruct(mainGen, ts)
		t.fnEntityGetSetPrivateFields(mainGen, ts)
		t.fnEntityEmpty(mainGen, ts)
		t.fnEntityCopy(mainGen, ts)
		t.fnEntityWriteTo(mainGen, ts)
		t.fnEntityDBAssignLastInsertID(mainGen, ts)
		t.fnEntityDBMapColumns(mainGen, ts)

		t.collectionStruct(mainGen, ts)
		t.fnCollectionDBAssignLastInsertID(mainGen, ts)
		t.fnCollectionDBMapColumns(mainGen, ts)
		t.fnCollectionUniqueGetters(mainGen, ts)
		t.fnCollectionUniquifiedGetters(mainGen, ts)
		t.fnCollectionWriteTo(mainGen, ts)
		t.fnCollectionFilter(mainGen, ts)
		t.fnCollectionEach(mainGen, ts)
		t.fnCollectionCut(mainGen, ts)
		t.fnCollectionSwap(mainGen, ts)
		t.fnCollectionDelete(mainGen, ts)
		t.fnCollectionInsert(mainGen, ts)
		t.fnCollectionAppend(mainGen, ts)

		if t.HasBinaryMarshaler {
			t.fnCollectionBinaryMarshaler(mainGen, ts)
		}
	}

	// now figure out all used package names in the buffer.
	pkgs, err := ts.findUsedPackages(mainGen.Bytes())
	if err != nil {
		_, _ = wMain.Write(mainGen.Bytes()) // write for debug reasons
		return errors.WithStack(err)
	}
	mainGen.AddImports(pkgs...)

	testGen.AddImports(ts.ImportPathsTesting...)

	if err := mainGen.GenerateFile(wMain); err != nil {
		return errors.WithStack(err)
	}
	if err := testGen.GenerateFile(wTest); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (ts *Generator) goTypeNull(c *ddl.Column) string { return ts.mySQLToGoType(c, true) }
func (ts *Generator) goType(c *ddl.Column) string     { return ts.mySQLToGoType(c, false) }
func (ts *Generator) goFuncNull(c *ddl.Column) string { return ts.mySQLToGoDmlColumnMap(c, true) }

func (ts *Generator) fnDBNewTables(mainGen *codegen.Go) {
	if !ts.hasFeature(0, 0, FeatureDB) {
		return
	}

	var tableNames []string
	var tableCreateStmt []string
	for _, tblname := range ts.sortedTableNames() {
		constName := `TableName` + strs.ToGoCamelCase(tblname)
		mainGen.AddConstString(constName, tblname)
		tableNames = append(tableNames, tblname)
		tableCreateStmt = append(tableCreateStmt, constName, `""`)
	}

	mainGen.C(`NewTables returns a goified version of the MySQL/MariaDB table schema for the tables: `,
		strings.Join(tableNames, ", "), `Auto generated by dmlgen.`)
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
}

func (ts *Generator) fnTestMainOther(testGen *codegen.Go, tables []*Table) {
	// Test Header
	testGen.Pln(`func TestNewTablesNonDB(t *testing.T) {`)
	{
		testGen.Pln(`ps := pseudo.MustNewService(0, &pseudo.Options{Lang: "de",FloatMaxDecimals:6})`)
		// If some features haven't been enabled, then there are no tests so
		// assign ps to underscore to avoid the unused variable error.
		// Alternatively figure out how not to print the whole test function at
		// all.
		testGen.Pln(`_ = ps`)

		for _, t := range tables {
			t.generateTestOther(testGen, ts)
		} // end for tables
	}
	testGen.Pln(`}`) // end TestNewTables
}

func (ts *Generator) fnTestMainDB(testGen *codegen.Go, tables []*Table) {
	if !ts.hasFeature(0, 0, FeatureDB) {
		return
	}

	// Test Header
	testGen.Pln(`func TestNewTablesDB(t *testing.T) {`)
	{
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
		testGen.Pln(`assert.Exactly(t, `, fmt.Sprintf("%#v", ts.sortedTableNames()), `, tblNames)`)

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
		if fn, ok := ts.customCode["pseudo.MustNewService.Option"]; ok {
			fn(ts, nil, testGen)
		}
		testGen.Out()
		testGen.Pln(`)`)

		for _, t := range tables {
			t.generateTestDB(testGen)
		} // end for tables
	}
	testGen.Pln(`}`) // end TestNewTables
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

	p := new(parser.Parser)
	if err := p.Parse(fname, fInfo.IsDir()); err != nil {
		return errors.CorruptData.Newf("[dmlgen] Error parsing failed %v: %v", fname, err)
	}

	var outName string
	if fInfo.IsDir() {
		outName = filepath.Join(fname, p.PkgName+"_easyjson.go")
	} else {
		if s := strings.TrimSuffix(fname, ".go"); s == fname {
			return errors.NotAcceptable.Newf("[dmlgen] GenerateJSON: Filename must end in '.go'")
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
