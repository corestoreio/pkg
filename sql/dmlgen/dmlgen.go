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
	"context"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"io"
	"os"
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
		t.featuresInclude = opt.FeaturesInclude | g.defaultTableConfig.FeaturesInclude
		t.featuresExclude = opt.FeaturesExclude | g.defaultTableConfig.FeaturesExclude
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
		if db != nil {
			g.kcu, err = ddl.LoadKeyColumnUsage(ctx, db, g.sortedTableNames()...)
			if err != nil {
				return errors.WithStack(err)
			}
			g.kcuRev = ddl.ReverseKeyColumnUsage(g.kcu)

			g.krs, err = ddl.GenerateKeyRelationships(ctx, db, g.kcu)
			if err != nil {
				return errors.WithStack(err)
			}
		} else if isDebug() {
			println("DEBUG[WithForeignKeyRelationships] db dml.Querier is nil. Did not load LoadKeyColumnUsage and GenerateKeyRelationships.")
		}

		// TODO(CSC) maybe excludeRelationships can contain a wild card to disable
		//  all embedded structs from/to a type. e.g. "customer_entity.website_id",
		//  "store_website.website_id", for CustomerEntity would become
		//  "*.website_id", "store_website.website_id", to disable all tables which
		//  have a foreign key to store_website.
		// TODO optimize code.

		g.krsExclude = make(map[string]bool, len(o.ExcludeRelationships)/2)
		for i := 0; i < len(o.ExcludeRelationships); i += 2 {
			mainTable := strings.Split(o.ExcludeRelationships[i], ".") // mainTable.mainColumn
			mainTab := mainTable[0]
			mainCol := mainTable[1]
			referencedTable := strings.Split(o.ExcludeRelationships[i+1], ".") // referencedTable.referencedColumn
			referencedTab := referencedTable[0]
			referencedCol := referencedTable[1]

			var buf strings.Builder
			buf.WriteString(mainTab)
			if mainCol == "*" {
				g.krsExclude[buf.String()] = true
			}
			buf.WriteByte('.')
			buf.WriteString(mainCol)
			if referencedTab == "*" && referencedCol == "*" {
				g.krsExclude[buf.String()] = true
			}

			buf.WriteByte(':')
			buf.WriteString(referencedTab)
			if referencedCol == "*" {
				g.krsExclude[buf.String()] = true
			}
			buf.WriteByte('.')
			buf.WriteString(referencedCol)
			g.krsExclude[buf.String()] = true
		}

		if len(o.IncludeRelationShips) > 0 {
			g.krsInclude = make(map[string]bool, len(o.IncludeRelationShips)/2)
			for i := 0; i < len(o.IncludeRelationShips); i += 2 {
				mainTable := strings.Split(o.IncludeRelationShips[i], ".") // mainTable.mainColumn
				mainTab := mainTable[0]
				mainCol := mainTable[1]
				referencedTable := strings.Split(o.IncludeRelationShips[i+1], ".") // referencedTable.referencedColumn
				referencedTab := referencedTable[0]
				referencedCol := referencedTable[1]

				var buf strings.Builder
				buf.WriteString(mainTab)
				if mainCol == "*" {
					g.krsInclude[buf.String()] = true
				}
				buf.WriteByte('.')
				buf.WriteString(mainCol)
				if referencedTab == "*" && referencedCol == "*" {
					g.krsInclude[buf.String()] = true
				}
				buf.WriteByte(':')
				buf.WriteString(referencedTab)
				if referencedCol == "*" {
					g.krsInclude[buf.String()] = true
				}
				buf.WriteByte('.')
				buf.WriteString(referencedCol)
				g.krsInclude[buf.String()] = true
			}
		}

		if isDebug() {
			g.krs.Debug(os.Stdout)
			debugMapSB("krsInclude", g.krsInclude)
			debugMapSB("krsExclude", g.krsExclude)
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
			for ci, ct := range t.Table.Columns {
				for _, cc := range columns {
					if ct.Field == cc.Field {
						t.Table.Columns[ci] = cc
					}
				}
			}
		} else {
			t = &Table{
				Table: ddl.NewTable(tableName, columns...),
				debug: isDebug(),
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
	checkAutoIncrement := func(oneTbl *ddl.Table) uint8 {
		if oneTbl.IsView() {
			return 1 // nope
		}
		for _, c := range oneTbl.Columns {
			if c.IsAutoIncrement() {
				return 2 // yes
			}
		}
		return 1 // nope
	}

	opt.sortOrder = 1
	opt.fn = func(g *Generator) error {
		nt, err := ddl.NewTables(ddl.WithLoadTables(ctx, db.DB, tables...))
		if err != nil {
			return errors.WithStack(err)
		}

		if len(tables) == 0 {
			tables = nt.Tables() // use all tables from the DB
		}

		for _, tblName := range tables {
			oneTbl := nt.MustTable(tblName)
			g.Tables[tblName] = &Table{
				Table:            oneTbl,
				HasAutoIncrement: checkAutoIncrement(oneTbl),
				debug:            isDebug(),
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
			"github.com/corestoreio/pkg/util/cstrace",
			"go.opentelemetry.io/otel/api/trace",
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

func (g *Generator) sortedTableNames() []string {
	sortedKeys := make(slices.String, 0, len(g.Tables))
	for k := range g.Tables {
		sortedKeys = append(sortedKeys, k)
	}
	sortedKeys.Sort()
	return sortedKeys
}

func (g *Generator) isAllowedRelationship(table1, column1, table2, column2 string) bool {
	// krsExclude
	{
		var buf strings.Builder
		buf.WriteString(table1)
		if g.krsExclude[buf.String()] { // tableName
			return false
		}
		buf.WriteByte('.')
		buf.WriteString(column1)
		if g.krsExclude[buf.String()] { // tableName.columnName
			return false
		}
		buf.WriteByte(':')
		buf.WriteString(table2)
		if g.krsExclude[buf.String()] { // tableName.columnName.referencedTableName
			return false
		}
		buf.WriteByte('.')
		buf.WriteString(column2)
		if g.krsExclude[buf.String()] { // tableName.columnName.referencedTableName.referencedColumnName
			return false
		}
	}
	if g.krsInclude == nil { // must be a nil check, only when not nil, then more restrictions apply.
		return true
	}

	// krsInclude block
	{
		var buf strings.Builder
		buf.WriteString(table1)
		if g.krsInclude[buf.String()] { // tableName
			return true
		}
		buf.WriteByte('.')
		buf.WriteString(column1)
		if g.krsInclude[buf.String()] { // tableName.columnName
			return true
		}

		buf.WriteByte(':')
		buf.WriteString(table2)
		if g.krsInclude[buf.String()] { // tableName.columnName.referencedTableName
			return true
		}
		buf.WriteByte('.')
		buf.WriteString(column2)
		return g.krsInclude[buf.String()] // tableName.columnName.referencedTableName.referencedColumnName
	}
}

func (g *Generator) hasFeature(includes, excludes, features FeatureToggle, mode ...rune) bool {
	if includes == 0 {
		includes = g.defaultTableConfig.FeaturesInclude
	}
	if excludes == 0 {
		excludes = g.defaultTableConfig.FeaturesExclude
	}
	// hasFeature runs with default mode: OR
	return hasFeature(includes, excludes, features, mode...) > 0
}

// findUsedPackages checks for needed packages which we must import.
func (g *Generator) findUsedPackages(file []byte,predefinedImportPaths []string) ([]string, error) {
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

	ret := make([]string, 0, len(predefinedImportPaths))
	for _, path := range predefinedImportPaths {
		_, pkg := filepath.Split(path)
		if _, ok := idents[pkg]; ok {
			ret = append(ret, path)
		}
	}
	return ret, nil
}

// serializerType converts the column type to the supported type of the current
// serializer. For now supports only protobuf.
func (g *Generator) serializerType(c *ddl.Column) string {
	pt := g.toSerializerType(c, true)
	if strings.IndexByte(pt, '/') > 0 { // slash identifies an import path
		return "bytes"
	}
	return pt
}

// serializerCustomType switches the default type from function serializerType
// to the new type. For now supports only protobuf.
func (g *Generator) serializerCustomType(c *ddl.Column) []string {
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

		proto.C(t.EntityName(), `represents a single row for`, t.Table.Name, `DB table. Auto generated.`)
		if t.Table.TableComment != "" {
			proto.C("Table comment:", t.Table.TableComment)
		}
		proto.Pln(`message`, t.EntityName(), `{`)
		{
			proto.In()
			var lastColumnPos uint64
			t.Table.Columns.Each(func(c *ddl.Column) {
				if t.IsFieldPublic(c.Field) {
					serType := g.serializerType(c)
					if !hasTimestampField && strings.HasPrefix(serType, "google.protobuf.Timestamp") {
						hasTimestampField = true
					}
					var optionConcret string
					if options := g.serializerCustomType(c); len(options) > 0 {
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

				if kcuc, ok := g.kcu[t.Table.Name]; ok { // kcu = keyColumnUsage && kcuc = keyColumnUsageCollection
					for _, kcuce := range kcuc.Data {
						if !kcuce.ReferencedTableName.Valid {
							continue
						}

						// case ONE-TO-MANY
						isOneToMany := g.krs.IsOneToMany(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						isRelationAllowed := g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						hasTable := g.Tables[kcuce.ReferencedTableName.Data] != nil
						if isOneToMany && hasTable && isRelationAllowed {
							proto.Pln(collectionName(kcuce.ReferencedTableName.Data), fieldMapFn(collectionName(kcuce.ReferencedTableName.Data)),
								"=", lastColumnPos, ";",
								"// 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
							lastColumnPos++
						}

						// case ONE-TO-ONE
						isOneToOne := g.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						if isOneToOne && hasTable && isRelationAllowed {
							proto.Pln(strs.ToGoCamelCase(kcuce.ReferencedTableName.Data), fieldMapFn(strs.ToGoCamelCase(kcuce.ReferencedTableName.Data)),
								"=", lastColumnPos, ";",
								"// 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
							lastColumnPos++
						}
					}
				}

				// TODO reversed M:N might be buggy as this code is not equal to the table.go code.
				if kcuc, ok := g.kcuRev[t.Table.Name]; ok { // kcu = keyColumnUsage && kcuc = keyColumnUsageCollection
					for _, kcuce := range kcuc.Data {
						if !kcuce.ReferencedTableName.Valid {
							continue
						}

						// case ONE-TO-MANY
						isOneToMany := g.krs.IsOneToMany(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						isRelationAllowed := g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						hasTable := g.Tables[kcuce.ReferencedTableName.Data] != nil
						if isOneToMany && hasTable && isRelationAllowed {
							proto.Pln(collectionName(kcuce.ReferencedTableName.Data), fieldMapFn(collectionName(kcuce.ReferencedTableName.Data)),
								"=", lastColumnPos, ";",
								"// Reversed 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
							lastColumnPos++
						}

						// case ONE-TO-ONE
						isOneToOne := g.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						if isOneToOne && hasTable && isRelationAllowed {
							proto.Pln(strs.ToGoCamelCase(kcuce.ReferencedTableName.Data), fieldMapFn(strs.ToGoCamelCase(kcuce.ReferencedTableName.Data)),
								"=", lastColumnPos, ";",
								"// Reversed 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
							lastColumnPos++
						}

						// case MANY-TO-MANY
						targetTbl, targetColumn := g.krs.ManyToManyTarget(kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						if targetTbl != "" && targetColumn != "" {
							isRelationAllowed = g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, targetTbl, targetColumn)
						}

						// case MANY-TO-MANY
						if isRelationAllowed && targetTbl != "" && targetColumn != "" {
							proto.Pln(collectionName(targetTbl), fieldMapFn(collectionName(targetTbl)),
								"=", lastColumnPos, ";",
								"// Reversed M:N", kcuce.TableName+"."+kcuce.ColumnName, "via", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data,
								"=>", targetTbl+"."+targetColumn)
							lastColumnPos++
						}

					}
				}
			}
			proto.Out()
		}
		proto.Pln(`}`)

		proto.C(t.CollectionName(), `represents multiple rows for the`, t.Table.Name, `DB table. Auto generated.`)
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

	g.fnCreateDBM(mainGen, tables)
	g.fnTestMainOther(testGen, tables)
	g.fnTestMainDB(testGen, tables)

	// deal with random map to guarantee the persistent code generation.
	for _, t := range tables {
		t.entityStruct(mainGen, g)

		t.fnEntityCopy(mainGen, g)
		t.fnEntityDBAssignLastInsertID(mainGen, g)
		t.fnEntityDBMapColumns(mainGen, g)
		t.fnEntityDBMHandler(mainGen, g)
		t.fnEntityEmpty(mainGen, g)
		t.fnEntityIsSet(mainGen, g)
		t.fnEntityGetSetPrivateFields(mainGen, g)
		t.fnEntityValidate(mainGen, g)
		t.fnEntityWriteTo(mainGen, g)

		t.collectionStruct(mainGen, g)

		t.fnCollectionAppend(mainGen, g)
		t.fnCollectionBinaryMarshaler(mainGen, g)
		t.fnCollectionCut(mainGen, g)
		t.fnCollectionDBAssignLastInsertID(mainGen, g)
		t.fnCollectionDBMapColumns(mainGen, g)
		t.fnCollectionDBMHandler(mainGen, g)
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
	pkgs, err := g.findUsedPackages(mainGen.Bytes(),g.ImportPaths)
	if err != nil {
		_, _ = wMain.Write(mainGen.Bytes()) // write for debug reasons
		return errors.WithStack(err)
	}
	mainGen.AddImports(pkgs...)

	if err := mainGen.GenerateFile(wMain); err != nil {
		return errors.WithStack(err)
	}

	pkgs, err = g.findUsedPackages(testGen.Bytes(),g.ImportPathsTesting)
	if err != nil {
		_, _ = wMain.Write(testGen.Bytes()) // write for debug reasons
		return errors.WithStack(err)
	}
	testGen.AddImports(pkgs...)
	//testGen.AddImports(g.ImportPathsTesting...)

	if err := testGen.GenerateFile(wTest); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (g *Generator) fnCreateDBM(mainGen *codegen.Go, tbls tables) {
	if !tbls.hasFeature(g, FeatureDB|FeatureDBTracing|FeatureDBSelect|FeatureDBDelete|
		FeatureDBInsert|FeatureDBUpdate|FeatureDBUpsert) {
		return
	}

	// TODO add some feature switches to include/excluded parts.
	var tableNames []string
	var tableCreateStmt []string
	var tableConstants []string
	for _, tblname := range g.sortedTableNames() {
		constName := `TableName` + strs.ToGoCamelCase(tblname)
		tableConstants = append(tableConstants, fmt.Sprintf("%s = %q", constName, tblname))
		tableNames = append(tableNames, tblname)
		tableCreateStmt = append(tableCreateStmt, constName, `""`)
	}
	mainGen.WriteConstants(tableConstants...)

	mainGen.Pln(`var dbmEmptyOpts = []dml.DBRFunc{func(dbr *dml.DBR) {
			// do nothing because Clone gets called automatically
		}}
		func dbmNoopResultCheckFn(_ sql.Result, err error) error { return err }
`)

	// <event functions>
	mainGen.C(`Event functions are getting dispatched during before or after handling a collection or an entity.
Context is always non-nil but either collection or entity pointer will be set.`)
	mainGen.Pln(`type (`)
	for _, tbl := range tbls {
		mainGen.Pln(`Event` + tbl.EntityName() + `Fn func(context.Context, *` + tbl.CollectionName() + `, *` + tbl.EntityName() + `) error`)
	}
	mainGen.Pln(`)`)
	// </event functions>

	// <DBM option struct>
	mainGen.C(`DBMOption provides various options to the DBM object.`)
	mainGen.Pln(`type DBMOption struct {`)
	{
		mainGen.Pln(tbls.hasFeature(g, FeatureDBTracing), `Trace                trace.Tracer`)
		mainGen.Pln(`TableOptions         []ddl.TableOption // gets applied at the beginning`)
		mainGen.Pln(`TableOptionsAfter    []ddl.TableOption // gets applied at the end`)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBSelect), `InitSelectFn         func(*dml.Select) *dml.Select`)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBUpdate), `InitUpdateFn         func(*dml.Update) *dml.Update`)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBDelete), `InitDeleteFn         func(*dml.Delete) *dml.Delete`)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBInsert|FeatureDBUpsert), `InitInsertFn         func(*dml.Insert) *dml.Insert`)
		for _, tbl := range tbls {
			mainGen.Pln(`event` + tbl.EntityName() + `Func [dml.EventFlagMax][]Event` + tbl.EntityName() + `Fn`)
		}
	}
	mainGen.Pln(`}`)
	// </DBM option struct>

	// <event adder>
	for _, tbl := range tbls {
		mainGen.C(codegen.SkipWS(`AddEvent`, tbl.EntityName()), `adds a specific defined event call back to the DBM.
It panics if the event argument is larger than dml.EventFlagMax.`)
		mainGen.Pln(`func (o *DBMOption) `, codegen.SkipWS(`AddEvent`, tbl.EntityName(), `(`), `event dml.EventFlag,`,
			`fn Event`+tbl.EntityName()+`Fn) *DBMOption {`)
		{
			mainGen.In()
			mainGen.Pln(`o.event` + tbl.EntityName() + `Func[event] = append(o.event` + tbl.EntityName() + `Func[event], fn)`)
			mainGen.Pln(`return o`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
	}
	// </event adder>

	mainGen.C(`DBM defines the DataBaseManagement object for the tables `, tableNames)
	mainGen.Pln(`type DBM struct { *ddl.Tables; option DBMOption }`)

	// <event dispatcher>
	for _, tbl := range tbls {
		mainGen.Pln(codegen.SkipWS(`func (dbm DBM) event`, tbl.EntityName(), `Func(ctx context.Context, ef dml.EventFlag, ec `, codegen.SkipWS(`*`, tbl.CollectionName()), `, e `, codegen.SkipWS(`*`, tbl.EntityName()), `) error`), ` {`)
		{
			mainGen.In()
			mainGen.Pln(`if len(dbm.option.`, codegen.SkipWS(`event`, tbl.EntityName(), `Func`), `[ef]) == 0 || dml.EventsAreSkipped(ctx) {`)
			mainGen.In()
			{
				mainGen.Pln(`return nil`)
			}
			mainGen.Out()
			mainGen.Pln(`}`)

			mainGen.Pln(`for _, fn := range dbm.option.`, codegen.SkipWS(`event`, tbl.EntityName(), `Func`), `[ef] {`)
			{
				mainGen.In()
				mainGen.Pln(`if err := fn(ctx, ec, e); err != nil {
				return errors.WithStack(err)
			}`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
			mainGen.Pln(`return nil`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
	} // </event dispatcher>

	mainGen.C(`NewDBManager returns a goified version of the MySQL/MariaDB table schema for the tables: `, tableNames, `Auto generated by dmlgen.`)
	mainGen.Pln(`func NewDBManager(ctx context.Context, dbmo *DBMOption) (*DBM, error) {`)
	{
		mainGen.In()
		mainGen.Pln(`tbls, err := ddl.NewTables(append([]ddl.TableOption{ddl.WithCreateTable(ctx, `, tableCreateStmt, `)},dbmo.TableOptions...)...)`)
		mainGen.Pln(`if err != nil { return nil, errors.WithStack(err); }`)

		mainGen.Pln(tbls.hasFeature(g, FeatureDBSelect),
			`	if dbmo.InitSelectFn == nil { dbmo.InitSelectFn = func(s *dml.Select) *dml.Select { return s; }; } `)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBUpdate),
			`	if dbmo.InitUpdateFn == nil { dbmo.InitUpdateFn = func(s *dml.Update) *dml.Update { return s; }; } `)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBDelete),
			`	if dbmo.InitDeleteFn == nil { dbmo.InitDeleteFn = func(s *dml.Delete) *dml.Delete { return s; }; } `)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBInsert|FeatureDBUpsert),
			`	if dbmo.InitInsertFn == nil { dbmo.InitInsertFn = func(s *dml.Insert) *dml.Insert { return s; }; } `)

		{
			mainGen.Pln(`err = tbls.Options(`)
			for _, tbl := range tbls {
				tbl.fnDBMOptionsSQLBuildQueries(mainGen, g)
			}
			mainGen.Pln(`)`) // end options
			mainGen.Pln(`if err != nil { return nil, err }`)
			mainGen.Pln(`if err := tbls.Options(dbmo.TableOptionsAfter...); err != nil { return nil, err }`)
		}
		mainGen.Out()
	}

	mainGen.Pln(tbls.hasFeature(g, FeatureDBTracing), `	if dbmo.Trace == nil { dbmo.Trace = trace.NoopTracer{}; }`)
	mainGen.Pln(`return &DBM{	Tables: tbls, option: *dbmo, }, nil }`)
}

func (g *Generator) fnTestMainOther(testGen *codegen.Go, tbls tables) {
	// Test Header
	lenBefore := testGen.Len()
	var codeWritten int

	testGen.Pln(`func TestNewDBManagerNonDB_` + tbls.nameID() + `(t *testing.T) {`)
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
	testGen.Pln(`}`) // end TestNewDBManager

	if codeWritten == 0 {
		testGen.Truncate(lenBefore)
	}
}

func (g *Generator) fnTestMainDB(testGen *codegen.Go, tbls tables) {
	if !tbls.hasFeature(g, FeatureDB) {
		return
	}

	// Test Header
	testGen.Pln(`func TestNewDBManagerDB_` + tbls.nameID() + `(t *testing.T) {`)
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
		testGen.Pln(`tbls, err := NewDBManager(ctx, &DBMOption{TableOptions: []ddl.TableOption{ddl.WithConnPool(db)}} )`)
		testGen.Pln(`assert.NoError(t, err)`)

		testGen.Pln(`tblNames := tbls.Tables.Tables()`)
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
	testGen.Pln(`}`) // end TestNewDBManager
}

func isDebug() bool {
	return os.Getenv("DEBUG") != ""
}

func debugMapSB(name string, m map[string]bool) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %q\n", name, k)
	}
}
