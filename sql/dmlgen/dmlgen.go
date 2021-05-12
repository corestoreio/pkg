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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/util/codegen"
	"github.com/corestoreio/pkg/util/slices"
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
		t.relationshipSeen = map[string]bool{}
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

func (g *Generator) findColumn(tableName, columnName string) *ddl.Column {
	if g.Tables[tableName] == nil {
		panic(fmt.Sprintf("table not found in Tables map: tableName:%q columnName:%q\n", tableName, columnName)) // TODO fix this, but how?
	}
	return g.Tables[tableName].Table.Columns.ByField(columnName)
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
	// deal with random map to guarantee the persistent code generation.
	for _, t := range tables {
		t.fnEntityRelationStruct(mainGen, g) // this must go first because to check if a table has relations
		t.fnEntityStruct(mainGen, g)
	}
	g.fnCreateDBM(mainGen, tables)
	g.fnTestMainOther(testGen, tables)
	g.fnTestMainDB(testGen, tables)

	// deal with random map to guarantee the persistent code generation.
	for _, t := range tables {
		//	t.fnEntityRelationStruct(mainGen, g)  // this must go first because to check if a table has relations
		t.fnEntityRelationMethods(mainGen, g) // this must go first because to check if a table has relations
		// t.fnEntityStruct(mainGen, g)

		t.fnEntityCopy(mainGen, g)
		t.fnEntityDBAssignLastInsertID(mainGen, g)
		t.fnEntityDBMapColumns(mainGen, g)
		t.fnEntityDBMHandler(mainGen, g)
		t.fnEntityEmpty(mainGen, g)
		t.fnEntityIsSet(mainGen, g)
		t.fnEntityGetSetPrivateFields(mainGen, g)
		t.fnEntityValidate(mainGen, g)
		t.fnEntityWriteTo(mainGen, g)

		t.fnCollectionStruct(mainGen, g)
		t.fnCollectionAppend(mainGen, g)
		t.fnCollectionClear(mainGen, g)
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
	pkgs, err := findUsedPackages(mainGen.Bytes(), g.ImportPaths)
	if err != nil {
		_, _ = wMain.Write(mainGen.Bytes()) // write for debug reasons
		return errors.WithStack(err)
	}
	mainGen.AddImports(pkgs...)

	if err := mainGen.GenerateFile(wMain); err != nil {
		return errors.WithStack(err)
	}

	pkgs, err = findUsedPackages(testGen.Bytes(), g.ImportPathsTesting)
	if err != nil {
		_, _ = wMain.Write(testGen.Bytes()) // write for debug reasons
		return errors.WithStack(err)
	}
	testGen.AddImports(pkgs...)
	// testGen.AddImports(g.ImportPathsTesting...)

	if err := testGen.GenerateFile(wTest); err != nil {
		return errors.WithStack(err)
	}
	return nil
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
