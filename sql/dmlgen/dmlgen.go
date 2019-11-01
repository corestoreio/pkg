// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	Package                     string // Name of the package
	PackageImportPath           string // Name of the package
	PackageSerializer           string // Name of the package, if empty uses field Package
	PackageSerializerImportPath string // Name of the package, if empty uses field PackageImportPath
	BuildTags                   []string
	ImportPaths                 []string
	ImportPathsTesting          []string
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
	krs    ddl.KeyRelationShips
	// "mainTable.mainColumn":"referencedTable.referencedColumn"
	// or skips the reversed relationship
	// "referencedTable.referencedColumn":"mainTable.mainColumn"
	krsExclude map[string]bool
	krsInclude map[string]bool
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
	// encoder names are: json resp. easyjson.
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
	// In case when foreign keys should be referenced:
	// 		[]string{"FieldNameX",`faker: "-"`,"FieldNameY","`xml:field_name_y,omitempty`"}
	// TODO CustomStructTags should be appended to StructTags
	CustomStructTags []string // balanced slice
	// AppendCustomStructTags []string // balanced slice TODO maybe this additionally
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

func (to *TableConfig) applyEncoders(t *Table, g *Generator) {
	encoders := make([]string, 0, len(to.Encoders)+len(g.defaultTableConfig.Encoders))
	encoders = append(encoders, to.Encoders...)
	encoders = append(encoders, g.defaultTableConfig.Encoders...)

	for i := 0; i < len(encoders) && to.lastErr == nil; i++ {
		switch enc := encoders[i]; enc {
		case "json", "easyjson":
			t.HasEasyJSONMarshaler = true
		case "protobuf", "fbs":
			t.HasSerializer = true // for now leave it in. maybe later PB gets added to the struct tags.
		default:
			to.lastErr = errors.NotSupported.Newf("[dmlgen] WithTableConfig: Table %q Encoder %q not supported", t.TableName, enc)
		}
	}
}

func (to *TableConfig) applyStructTags(t *Table, g *Generator) {
	for h := 0; h < len(t.Columns) && to.lastErr == nil; h++ {
		c := t.Columns[h]
		var buf strings.Builder

		structTags := make([]string, 0, len(to.StructTags)+len(g.defaultTableConfig.StructTags))
		structTags = append(structTags, to.StructTags...)
		structTags = append(structTags, g.defaultTableConfig.StructTags...)

		for lst, i := len(structTags), 0; i < lst && to.lastErr == nil; i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			// Maybe some types in the struct for a table don't need at
			// all an omitempty so build in some logic which creates the
			// tags more thoughtfully.
			switch tagName := structTags[i]; tagName {
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
		for _, c := range t.Columns {
			if c.Field == to.CustomStructTags[i] {
				c.StructTag = to.CustomStructTags[i+1]
			}
		}
		if t.customStructTagFields == nil {
			t.customStructTagFields = make(map[string]string)
		}

		// copy data to handle foreign keys if they should have struct tags.
		// as key use the kcuce.ReferencedTableName.String and value the struct tag itself.
		t.customStructTagFields[to.CustomStructTags[i]] = "`" + to.CustomStructTags[i+1] + "`"
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
// overwritten on a per table basis with function WithTableConfig. Current
// behaviour: Does not apply the default config to all tables but only to those
// which have set an empty
// `WithTableConfig("table_name",&dmlgen.TableConfig{})`.
func WithTableConfigDefault(opt TableConfig) (o Option) {
	o.sortOrder = 149 // must be applied before WithTableConfig
	o.fn = func(g *Generator) (err error) {
		g.defaultTableConfig = opt
		return opt.lastErr
	}
	return o
}

func defaultFieldMapFn(s string) string {
	return s
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
	o.fn = func(g *Generator) (err error) {
		t, ok := g.Tables[tableName]
		if t == nil || !ok {
			return errors.NotFound.Newf("[dmlgen] WithTableConfig: Table %q not found.", tableName)
		}
		opt.applyEncoders(t, g)
		opt.applyStructTags(t, g)
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

// ForeignKeyOptions applies to WithForeignKeyRelationships
type ForeignKeyOptions struct {
	IncludeRelationShips []string
	ExcludeRelationships []string
	// MToMRelations        bool
}

// WithForeignKeyRelationships analyzes the foreign keys which points to a table
// and adds them as a struct field name. For example:
// customer_address_entity.parent_id is a foreign key to
// customer_entity.entity_id hence the generated struct CustomerEntity has a new
// field which gets named CustomerAddressEntityCollection, pointing to type
// CustomerAddressEntityCollection. includeRelationShips and
// excludeRelationships must be a balanced slice in the notation of
// "table1.column1","table2.column2". For example:
// 		"customer_entity.store_id", "store.store_id" which means that the
// struct CustomerEntity won't have a field to the Store struct (1:1
// relationship) (in case of excluding). The reverse case can also be added
// 		"store.store_id", "customer_entity.store_id" which means that the
// Store struct won't or will have a field pointing to the
// CustomerEntityCollection (1:M relationship).
// Setting includeRelationShips to nil will include all relationships.
func WithForeignKeyRelationships(ctx context.Context, db dml.Querier, o ForeignKeyOptions) (opt Option) {
	opt.sortOrder = 210 // must run at the end or where the end is near ;-)
	opt.fn = func(g *Generator) (err error) {
		if len(o.ExcludeRelationships)%2 == 1 {
			return errors.Fatal.Newf("[dmlgen] excludeRelationships must be balanced slice. Read the doc.")
		}
		if len(o.IncludeRelationShips)%2 == 1 {
			return errors.Fatal.Newf("[dmlgen] includeRelationShips must be balanced slice. Read the doc.")
		}

		g.kcu, err = ddl.LoadKeyColumnUsage(ctx, db, g.sortedTableNames()...)
		if err != nil {
			return errors.WithStack(err)
		}
		g.kcuRev = ddl.ReverseKeyColumnUsage(g.kcu)

		g.krs, err = ddl.GenerateKeyRelationships(ctx, db, g.kcu)
		if err != nil {
			return errors.WithStack(err)
		}

		// TODO(CSC) maybe excludeRelationships can contain a wild card to disable
		//  all embedded structs from/to a type. e.g. "customer_entity.website_id",
		//  "store_website.website_id", for CustomerEntity would become
		//  "*.website_id", "store_website.website_id", to disable all tables which
		//  have a foreign key to store_website.

		g.krsExclude = make(map[string]bool, len(o.ExcludeRelationships)/2)
		for i := 0; i < len(o.ExcludeRelationships); i += 2 {
			var buf strings.Builder
			buf.WriteString(o.ExcludeRelationships[i])
			buf.WriteByte(':')
			buf.WriteString(o.ExcludeRelationships[i+1])
			g.krsExclude[buf.String()] = true
		}
		if len(o.IncludeRelationShips) > 0 {
			g.krsInclude = make(map[string]bool, len(o.IncludeRelationShips)/2)
			for i := 0; i < len(o.IncludeRelationShips); i += 2 {
				var buf strings.Builder
				buf.WriteString(o.IncludeRelationShips[i])
				buf.WriteByte(':')
				buf.WriteString(o.IncludeRelationShips[i+1])
				g.krsInclude[buf.String()] = true
			}
		}
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
	opt.fn = func(g *Generator) error {
		isOverwrite := len(options) > 0 && options[0] == "overwrite"
		t, ok := g.Tables[tableName]
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
				debug:     os.Getenv("DEBUG") != "",
			}
		}
		t.HasAutoIncrement = checkAutoIncrement(t.HasAutoIncrement)
		g.Tables[tableName] = t
		return nil
	}
	return opt
}

// WithTablesFromDB queries the information_schema table and loads the column
// definition of the provided `tables` slice. It adds the tables to the `Generator`
// map. Once added a call to WithTableConfig can add additional configurations.
func WithTablesFromDB(ctx context.Context, db *dml.ConnPool, tables ...string) (opt Option) {
	opt.sortOrder = 1
	opt.fn = func(g *Generator) error {
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
			g.Tables[tblName] = &Table{
				TableName:        tblName,
				Columns:          tables[tblName],
				HasAutoIncrement: checkAutoIncrement(tblName),
				debug:            os.Getenv("DEBUG") != "",
			}
		}
		return nil
	}
	return opt
}

// SerializerConfig applies various optional settings to WithProtobuf and/or
// WithFlatbuffers and/or WithTypeScript.
type SerializerConfig struct {
	PackageImportPath string
	Headers           []string
}

// WithProtobuf enables protocol buffers as a serialization method. Argument
// headerOptions is optional. Heads up: This function also sets the internal
// Serializer field and all types will get adjusted to the minimum protobuf
// types. E.g. uint32 minimum instead of uint8/uint16. So if the Generator gets
// created multiple times to separate the creation of code, the WithProtobuf
// function must get set for Generator objects. See package store.
func WithProtobuf(sc *SerializerConfig) (opt Option) {
	_, pkg := filepath.Split(sc.PackageImportPath)

	opt.sortOrder = 110
	opt.fn = func(g *Generator) error {
		g.Serializer = "protobuf"
		g.PackageSerializer = pkg
		g.PackageSerializerImportPath = sc.PackageImportPath
		if len(sc.Headers) == 0 {
			g.SerializerHeaderOptions = []string{
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
func WithFlatbuffers(sc *SerializerConfig) (opt Option) {
	panic("TODO implement WithFlatbuffers")
	// opt.sortOrder = 111
	// opt.fn = func(g *Generator) error {
	// 	g.Serializer = "fbs"
	// 	// if len(headerOptions) == 0 {
	// 	// g.SerializerHeaderOptions = []string{}
	// 	// }
	// 	return nil
	// }
	return opt
}

// WithBuildTags adds your build tags to the file header. Each argument
// represents a build tag line.
func WithBuildTags(lines ...string) (opt Option) {
	opt.sortOrder = 112
	opt.fn = func(g *Generator) error {
		g.BuildTags = append(g.BuildTags, lines...)
		return nil
	}
	return opt
}

// WithCustomCode inserts at the marker position your custom Go code. For
// available markers search these .go files for the map access of field
// *Generator.customCode. An example got written in
// TestGenerate_Tables_Protobuf_Json. If the marker does not exists or has a
// typo, no error gets reported and no code gets written.
func WithCustomCode(marker, code string) (opt Option) {
	opt.sortOrder = 112
	opt.fn = func(g *Generator) error {
		if g.customCode == nil {
			g.customCode = make(map[string]func(*Generator, *Table, io.Writer))
		}
		g.customCode[marker] = func(_ *Generator, _ *Table, w io.Writer) {
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
	opt.fn = func(g *Generator) error {
		if g.customCode == nil {
			g.customCode = make(map[string]func(*Generator, *Table, io.Writer))
		}
		g.customCode[marker] = fn
		return nil
	}
	return opt
}

func (g *Generator) sortedTableNames() []string {
	sortedKeys := make(slices.String, 0, len(g.Tables))
	for k := range g.Tables {
		sortedKeys = append(sortedKeys, k)
	}
	sortedKeys.Sort()
	return sortedKeys
}

// NewGenerator creates a new instance of the SQL table code generator. The order
// of the applied options does not matter as they are getting sorted internally.
func NewGenerator(packageImportPath string, opts ...Option) (*Generator, error) {
	_, pkg := filepath.Split(packageImportPath)
	g := &Generator{
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

			"github.com/corestoreio/errors",
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
		if err := opt.fn(g); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	for _, t := range g.Tables {
		t.Package = g.Package
	}
	return g, nil
}

func (g *Generator) isAllowedRelationship(table1, column1, table2, column2 string) bool {
	var buf strings.Builder
	buf.WriteString(table1)
	buf.WriteByte('.')
	buf.WriteString(column1)
	buf.WriteByte(':')
	buf.WriteString(table2)
	buf.WriteByte('.')
	buf.WriteString(column2)

	var allowed bool
	switch {
	case g.krsExclude[buf.String()]:
		allowed = false
	case g.krsInclude == nil:
		allowed = true
	default:
		allowed = g.krsInclude[buf.String()]
	}
	return allowed
}

func (g *Generator) hasFeature(tableInclude, tableExclude, feature FeatureToggle) bool {
	if tableInclude == 0 {
		tableInclude = g.defaultTableConfig.FeaturesInclude
	}
	if tableExclude == 0 {
		tableExclude = g.defaultTableConfig.FeaturesExclude
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
func (g *Generator) findUsedPackages(file []byte) ([]string, error) {
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

	ret := make([]string, 0, len(g.ImportPaths))
	for _, path := range g.ImportPaths {
		_, pkg := filepath.Split(path)
		if _, ok := idents[pkg]; ok {
			ret = append(ret, path)
		}
	}
	return ret, nil
}

// SerializerType converts the column type to the supported type of the current
// serializer. For now supports only protobuf.
func (g *Generator) SerializerType(c *ddl.Column) string {
	pt := g.toSerializerType(c, true)
	if strings.IndexByte(pt, '/') > 0 { // slash identifies an import path
		return "bytes"
	}
	return pt
}

// SerializerCustomType switches the default type from function SerializerType
// to the new type. For now supports only protobuf.
func (g *Generator) SerializerCustomType(c *ddl.Column) []string {
	pt := g.toSerializerType(c, true)
	var buf []string
	if pt == "google.protobuf.Timestamp" {
		buf = append(buf, "(gogoproto.stdtime)=true")
	}
	if pt == "bytes" {
		return nil // bytes can be null
	}
	if c.IsNull() || strings.IndexByte(pt, '.') > 0 /*whenever it is a custom type like null. or google.proto.timestamp*/ {
		// Indeed nullable Go Types must be not-nullable in HasSerializer because we
		// have a non-pointer struct type which contains the field Valid.
		// HasSerializer treats nullable fields as pointer fields, but that is
		// ridiculous.
		buf = append(buf, "(gogoproto.nullable)=false")
	}
	return buf
}

// GenerateSerializer writes the protocol buffer specifications into `w` and its test
// sources into wTest, if there are any tests.
func (g *Generator) GenerateSerializer(wMain, wTest io.Writer) error {
	switch g.Serializer {
	case "protobuf":
		if err := g.generateProto(wMain); err != nil {
			return errors.WithStack(err)
		}
	case "fbs":

	case "", "default", "none":
		return nil // do nothing
	default:
		return errors.NotAcceptable.Newf("[dmlgen] Serializer %q not supported.", g.Serializer)
	}

	return nil
}

func (g *Generator) generateProto(w io.Writer) error {
	pPkg := g.PackageSerializer
	if pPkg == "" {
		pPkg = g.Package
	}

	proto := codegen.NewProto(pPkg)
	proto.Pln(`import "github.com/gogo/protobuf/gogoproto/gogo.proto";`)

	const importTimeStamp = `import "google/protobuf/timestamp.proto";`
	proto.Pln(importTimeStamp)
	proto.Pln(`import "github.com/corestoreio/pkg/storage/null/null.proto";`)
	proto.Pln(`option go_package = "` + pPkg + `";`)

	for _, o := range g.SerializerHeaderOptions {
		proto.Pln(`option ` + o + `;`)
	}

	var hasTimestampField bool
	for _, tblname := range g.sortedTableNames() {
		t := g.Tables[tblname] // must panic if table name not found

		fieldMapFn := g.defaultTableConfig.FieldMapFn
		if fieldMapFn == nil {
			fieldMapFn = t.fieldMapFn
		}
		if fieldMapFn == nil {
			fieldMapFn = defaultFieldMapFn
		}

		proto.C(t.EntityName(), `represents a single row for DB table`, t.TableName, `. Auto generated.`)
		proto.Pln(`message`, t.EntityName(), `{`)
		{
			proto.In()
			var lastColumnPos uint64
			t.Columns.Each(func(c *ddl.Column) {
				if t.IsFieldPublic(c.Field) {
					serType := g.SerializerType(c)
					if !hasTimestampField && strings.HasPrefix(serType, "google.protobuf.Timestamp") {
						hasTimestampField = true
					}
					var optionConcret string
					if options := g.SerializerCustomType(c); len(options) > 0 {
						optionConcret = `[` + strings.Join(options, ",") + `]`
					}
					// extend here with a custom code option, if someone needs
					proto.Pln(serType, strs.ToGoCamelCase(c.Field), `=`, c.Pos, optionConcret+`;`)
					lastColumnPos = c.Pos
				}
			})
			lastColumnPos++

			if g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityRelationships) {

				// for debugging see Table.entityStruct function. This code is only different in the Pln function.

				if kcuc, ok := g.kcu[t.TableName]; ok { // kcu = keyColumnUsage && kcuc = keyColumnUsageCollection
					for _, kcuce := range kcuc.Data {
						if !kcuce.ReferencedTableName.Valid {
							continue
						}

						// case ONE-TO-MANY
						isOneToMany := g.krs.IsOneToMany(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						isRelationAllowed := g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						hasTable := g.Tables[kcuce.ReferencedTableName.String] != nil
						if isOneToMany && hasTable && isRelationAllowed {
							proto.Pln(collectionName(kcuce.ReferencedTableName.String), fieldMapFn(collectionName(kcuce.ReferencedTableName.String)),
								"=", lastColumnPos, ";",
								"// 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
							lastColumnPos++
						}

						// case ONE-TO-ONE
						isOneToOne := g.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						if isOneToOne && hasTable && isRelationAllowed {
							proto.Pln(strs.ToGoCamelCase(kcuce.ReferencedTableName.String), fieldMapFn(strs.ToGoCamelCase(kcuce.ReferencedTableName.String)),
								"=", lastColumnPos, ";",
								"// 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
							lastColumnPos++
						}
					}
				}

				if kcuc, ok := g.kcuRev[t.TableName]; ok { // kcu = keyColumnUsage && kcuc = keyColumnUsageCollection
					for _, kcuce := range kcuc.Data {
						if !kcuce.ReferencedTableName.Valid {
							continue
						}

						// case ONE-TO-MANY
						isOneToMany := g.krs.IsOneToMany(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						isRelationAllowed := g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						hasTable := g.Tables[kcuce.ReferencedTableName.String] != nil
						if isOneToMany && hasTable && isRelationAllowed {
							proto.Pln(collectionName(kcuce.ReferencedTableName.String), fieldMapFn(collectionName(kcuce.ReferencedTableName.String)),
								"=", lastColumnPos, ";",
								"// Reversed 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
							lastColumnPos++
						}

						// case ONE-TO-ONE
						isOneToOne := g.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						if isOneToOne && hasTable && isRelationAllowed {
							proto.Pln(strs.ToGoCamelCase(kcuce.ReferencedTableName.String), fieldMapFn(strs.ToGoCamelCase(kcuce.ReferencedTableName.String)),
								"=", lastColumnPos, ";",
								"// Reversed 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
							lastColumnPos++
						}
					}
				}
			}
			proto.Out()
		}
		proto.Pln(`}`)

		proto.C(t.CollectionName(), `represents multiple rows for DB table`, t.TableName, `. Auto generated.`)
		proto.Pln(`message`, t.CollectionName(), `{`)
		{
			proto.In()
			proto.Pln(`repeated`, t.EntityName(), `Data = 1;`)
			proto.Out()
		}
		proto.Pln(`}`)
	}

	if !hasTimestampField {
		// bit hacky to remove the import of timestamp proto but for now OK.
		removedImport := strings.ReplaceAll(proto.String(), importTimeStamp, "")
		proto.Reset()
		proto.WriteString(removedImport)
	}

	return proto.GenerateFile(w)
}

// GenerateGo writes the Go source code into `w` and the test code into wTest.
func (g *Generator) GenerateGo(wMain, wTest io.Writer) error {
	mainGen := codegen.NewGo(g.Package)
	testGen := codegen.NewGo(g.Package)

	mainGen.SecondLineComments = []string{"Generated by sql/dmlgen. DO NOT EDIT."}
	testGen.SecondLineComments = []string{"Generated by sql/dmlgen. DO NOT EDIT."}

	mainGen.BuildTags = g.BuildTags
	testGen.BuildTags = g.BuildTags

	tables := make([]*Table, len(g.Tables))
	for i, tblname := range g.sortedTableNames() {
		tables[i] = g.Tables[tblname] // must panic if table name not found
	}

	g.fnDBNewTables(mainGen, tables)
	g.fnTestMainOther(testGen, tables)

	g.fnTestMainDB(testGen, tables)

	// deal with random map to guarantee the persistent code generation.
	for _, t := range tables {
		t.entityStruct(mainGen, g)

		t.fnEntityCopy(mainGen, g)
		t.fnEntityDBAssignLastInsertID(mainGen, g)
		t.fnEntityDBMapColumns(mainGen, g)
		t.fnEntityEmpty(mainGen, g)
		t.fnEntityGetSetPrivateFields(mainGen, g)
		t.fnEntityValidate(mainGen, g)
		t.fnEntityWriteTo(mainGen, g)

		t.collectionStruct(mainGen, g)

		t.fnCollectionAppend(mainGen, g)
		t.fnCollectionBinaryMarshaler(mainGen, g)
		t.fnCollectionCut(mainGen, g)
		t.fnCollectionDBAssignLastInsertID(mainGen, g)
		t.fnCollectionDBMapColumns(mainGen, g)
		t.fnCollectionDelete(mainGen, g)
		t.fnCollectionEach(mainGen, g)
		t.fnCollectionFilter(mainGen, g)
		t.fnCollectionInsert(mainGen, g)
		t.fnCollectionSwap(mainGen, g)
		t.fnCollectionUniqueGetters(mainGen, g)
		t.fnCollectionUniquifiedGetters(mainGen, g)
		t.fnCollectionValidate(mainGen, g)
		t.fnCollectionWriteTo(mainGen, g)
	}

	// now figure out all used package names in the buffer.
	pkgs, err := g.findUsedPackages(mainGen.Bytes())
	if err != nil {
		_, _ = wMain.Write(mainGen.Bytes()) // write for debug reasons
		return errors.WithStack(err)
	}
	mainGen.AddImports(pkgs...)

	testGen.AddImports(g.ImportPathsTesting...)

	if err := mainGen.GenerateFile(wMain); err != nil {
		return errors.WithStack(err)
	}

	if err := testGen.GenerateFile(wTest); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (g *Generator) goTypeNull(c *ddl.Column) string { return g.mySQLToGoType(c, true) }
func (g *Generator) goType(c *ddl.Column) string     { return g.mySQLToGoType(c, false) }
func (g *Generator) goFuncNull(c *ddl.Column) string { return g.mySQLToGoDmlColumnMap(c, true) }

func (g *Generator) fnDBNewTables(mainGen *codegen.Go, tbls []*Table) {
	if !tables(tbls).hasFeature(g, FeatureDB) {
		return
	}

	var tableNames []string
	var tableCreateStmt []string
	for _, tblname := range g.sortedTableNames() {
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

func (g *Generator) fnTestMainOther(testGen *codegen.Go, tbls tables) {
	// Test Header
	lenBefore := testGen.Len()
	var codeWritten int

	testGen.Pln(`func TestNewTablesNonDB_` + tbls.nameID() + `(t *testing.T) {`)
	{
		testGen.Pln(`ps := pseudo.MustNewService(0, &pseudo.Options{Lang: "de",FloatMaxDecimals:6})`)
		// If some features haven't been enabled, then there are no tests so
		// assign ps to underscore to avoid the unused variable error.
		// Alternatively figure out how not to print the whole test function at
		// all.
		testGen.Pln(`_ = ps`)
		for _, t := range tbls {
			codeWritten += t.generateTestOther(testGen, g)
		}
	}
	testGen.Pln(`}`) // end TestNewTables

	if codeWritten == 0 {
		testGen.Truncate(lenBefore)
	}
}

func (g *Generator) fnTestMainDB(testGen *codegen.Go, tbls tables) {
	if !tables(tbls).hasFeature(g, FeatureDB) {
		return
	}

	// Test Header
	testGen.Pln(`func TestNewTablesDB_` + tbls.nameID() + `(t *testing.T) {`)
	{
		testGen.Pln(`db := dmltest.MustConnectDB(t)`)
		testGen.Pln(`defer dmltest.Close(t, db)`)

		if g.TestSQLDumpGlobPath != "" {
			testGen.Pln(`defer dmltest.SQLDumpLoad(t,`, strconv.Quote(g.TestSQLDumpGlobPath), `, &dmltest.SQLDumpOptions{
					SkipDBCleanup: true,
				}).Deferred()`)
		}

		testGen.Pln(`ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)`)
		testGen.Pln(`defer cancel()`)
		testGen.Pln(`tbls, err := NewTables(ctx, ddl.WithConnPool(db))`)
		testGen.Pln(`assert.NoError(t, err)`)

		testGen.Pln(`tblNames := tbls.Tables()`)
		testGen.Pln(`sort.Strings(tblNames)`)
		testGen.Pln(`assert.Exactly(t, `, fmt.Sprintf("%#v", g.sortedTableNames()), `, tblNames)`)

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
		if fn, ok := g.customCode["pseudo.MustNewService.Option"]; ok {
			fn(g, nil, testGen)
		}
		testGen.Out()
		testGen.Pln(`)`)

		for _, t := range tbls {
			t.generateTestDB(testGen)
		} // end for tables
	}
	testGen.Pln(`}`) // end TestNewTables
}

// ProtocOptions allows to modify the protoc CLI command.
type ProtocOptions struct {
	BuildTags          []string
	WorkingDirectory   string
	ProtoGen           string // default gofast, options: gogo, gogofast, gogofaster
	Debug              bool   // prints the final protoc command
	GRPC               bool
	GRPCGatewayOutMap  []string // GRPC must be enabled in the above field
	GRPCGatewayOutPath string   // GRPC must be enabled in the above field
	ProtoPath          []string
	GoGoOutPath        string
	GoGoOutMap         []string
	SwaggerOutPath     string
	CustomArgs         []string
	// TODO add validation plugin, either
	//  https://github.com/mwitkow/go-proto-validators as used in github.com/gogo/grpc-example/proto/example.proto
	//  This github.com/mwitkow/go-proto-validators seems dead.
	//  or https://github.com/envoyproxy/protoc-gen-validate
	//  Requirement: error messages must be translatable and maybe use an errors.Kind type
}

var defaultProtoPaths = make([]string, 0, 8)

func init() {
	preDefinedPaths := [...]string{
		build.Default.GOPATH + "/src/",
		build.Default.GOPATH + "/src/github.com/gogo/protobuf/protobuf/",
		build.Default.GOPATH + "/src/github.com/gogo/googleapis/",
		"vendor/github.com/grpc-ecosystem/grpc-gateway/",
		"vendor/github.com/gogo/googleapis/",
		"vendor/",
		".",
	}
	for _, pdp := range preDefinedPaths {
		if _, err := os.Stat(pdp); !os.IsNotExist(err) {
			defaultProtoPaths = append(defaultProtoPaths, pdp)
		}
	}
}

func (po *ProtocOptions) toArgs() []string {
	gogoOut := make([]string, 0, 4)
	if po.GRPC {
		gogoOut = append(gogoOut, "plugins=grpc")
		if po.GRPCGatewayOutMap == nil {
			po.GRPCGatewayOutMap = []string{
				"allow_patch_feature=false",
				"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types",
				"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types",
				"Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types",
				"Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api",
				"Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types",
			}
		}
		if po.GRPCGatewayOutPath == "" {
			po.GRPCGatewayOutPath = "."
		}
	}
	if po.GoGoOutPath == "" {
		po.GoGoOutPath = "."
	}
	if po.GoGoOutMap == nil {
		po.GoGoOutMap = []string{
			"Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api",
			"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types",
			"Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types",
			"Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types",
			"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types",
		}
	}
	gogoOut = append(gogoOut, po.GoGoOutMap...)

	if po.ProtoPath == nil {
		po.ProtoPath = append(po.ProtoPath, defaultProtoPaths...)
	}

	if po.ProtoGen == "" {
		po.ProtoGen = "gofast"
	} else {
		switch po.ProtoGen {
		case "gofast", "gogo", "gogofast", "gogofaster":
			// ok
		default:
			panic(fmt.Sprintf("[dmlgen] ProtoGen CLI command %q not supported, allowed: gofast, gogo, gogofast, gogofaster", po.ProtoGen))
		}
	}

	// To generate PHP Code replace `gogo_out` with `php_out`.
	// Java bit similar. Java has ~15k LOC, Go ~3.7k
	args := []string{
		"--" + po.ProtoGen + "_out", fmt.Sprintf("%s:%s", strings.Join(gogoOut, ","), po.GoGoOutPath),
		"--proto_path", strings.Join(po.ProtoPath, ":"),
	}
	if po.GRPC && len(po.GRPCGatewayOutMap) > 0 {
		args = append(args, "--grpc-gateway_out="+strings.Join(po.GRPCGatewayOutMap, ",")+":"+po.GRPCGatewayOutPath)
	}
	if po.SwaggerOutPath != "" {
		args = append(args, "--swagger_out="+po.SwaggerOutPath)
	}
	return append(args, po.CustomArgs...)
}

func (po *ProtocOptions) chdir() (deferred func(), _ error) {
	deferred = func() {}
	if po.WorkingDirectory != "" {
		oldWD, err := os.Getwd()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err := os.Chdir(po.WorkingDirectory); err != nil {
			return nil, errors.Wrapf(err, "[dmlgen] Failed to chdir to %q", po.WorkingDirectory)
		}
		deferred = func() {
			_ = os.Chdir(oldWD)
		}
	}
	return deferred, nil
}

// GenerateProto searches all *.proto files in the given path and calls protoc
// to generate the Go source code.
func GenerateProto(protoFilesPath string, po *ProtocOptions) error {
	restoreFn, err := po.chdir()
	if err != nil {
		return errors.WithStack(err)
	}
	defer restoreFn()

	protoFilesPath = filepath.Clean(protoFilesPath)
	if ps := string(os.PathSeparator); !strings.HasSuffix(protoFilesPath, ps) {
		protoFilesPath += ps
	}

	protoFiles, err := filepath.Glob(protoFilesPath + "*.proto")
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] Can't access proto files in path %q", protoFilesPath)
	}

	cmd := exec.Command("protoc", append(po.toArgs(), protoFiles...)...)
	if po.Debug {
		if po.WorkingDirectory == "" {
			po.WorkingDirectory = "."
		}
		// TODO fix ./dmlgen.go:1291:59: cmd.String undefined (type *exec.Cmd has no field or method String)
		// fmt.Printf("\ncd %s && %s\n\n", po.WorkingDirectory, cmd.String())
	}
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
	pbGoFiles, err := filepath.Glob(protoFilesPath + "*.pb.*go")
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] Can't access pb.go files in path %q", protoFilesPath)
	}

	removeImports := [][]byte{
		[]byte("import null \"github.com/corestoreio/pkg/storage/null\"\n"),
		[]byte("null \"github.com/corestoreio/pkg/storage/null\"\n"),
	}
	for _, file := range pbGoFiles {
		fContent, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, ri := range removeImports {
			fContent = bytes.Replace(fContent, ri, nil, -1)
		}

		var buf bytes.Buffer
		for _, bt := range po.BuildTags {
			fmt.Fprintf(&buf, "// +build %s\n", bt)
		}
		if buf.Len() > 0 {
			buf.WriteByte('\n')
			buf.Write(fContent)
			fContent = buf.Bytes()
		}

		if err := ioutil.WriteFile(file, fContent, 0644); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// GenerateJSON creates the easysjon code for a specific file or a whole
// directory. argument `g` can be nil.
func GenerateJSON(fileNameOrDirectory, buildTags string, g *bootstrap.Generator) error {
	fInfo, err := os.Stat(fileNameOrDirectory)
	if err != nil {
		return errors.WithStack(err)
	}

	p := new(parser.Parser)
	if err := p.Parse(fileNameOrDirectory, fInfo.IsDir()); err != nil {
		return errors.CorruptData.Newf("[dmlgen] Error parsing failed %q: %v", fileNameOrDirectory, err)
	}

	var outName string
	if fInfo.IsDir() {
		outName = filepath.Join(fileNameOrDirectory, p.PkgName+"_easyjson.go")
	} else {
		if s := strings.TrimSuffix(fileNameOrDirectory, ".go"); s == fileNameOrDirectory {
			return errors.NotAcceptable.Newf("[dmlgen] GenerateJSON: Filename must end in '.go'")
		} else {
			outName = s + "_easyjson.go"
		}
	}

	if len(p.StructNames) == 0 {
		return errors.NotFound.Newf("[dmlgen] Can't find any StructNames in the Go files of %q", fileNameOrDirectory)
	}

	if g == nil {
		g = &bootstrap.Generator{
			BuildTags:             buildTags,
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
	} else {
		g.Types = p.StructNames
	}
	if err := g.Run(); err != nil {
		return errors.Fatal.Newf("[dmlgen] easyJSON: Bootstrap failed: %v", err)
	}
	return nil
}
