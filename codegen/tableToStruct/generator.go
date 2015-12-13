// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/codegen/tableToStruct/tpl"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/log"
)

type generator struct {
	tts             codegen.TableToStruct
	dbrConn         *dbr.Connection
	outfile         *os.File
	tables          []string         // all available tables for which we should at least generate a type definition
	whiteListTables util.StringSlice // table name in this slice is allowed for generic functions
	eavValueTables  codegen.TypeCodeValueTable
	wg              *sync.WaitGroup
	// existingMethodSets contains all existing method sets from a package for the Table* types
	existingMethodSets *duplicateChecker
	isMagento1         bool
	isMagento2         bool
}

func newGenerator(tts codegen.TableToStruct, dbrConn *dbr.Connection, wg *sync.WaitGroup) *generator {
	wg.Add(1)
	return &generator{
		tts:                tts,
		dbrConn:            dbrConn,
		wg:                 wg,
		existingMethodSets: newDuplicateChecker(),
	}
}

func (g *generator) run() {
	defer log.WhenDone(PkgLog).Info("Stats", "Package", g.tts.Package)
	defer g.wg.Done()
	g.analyzePackage()

	var err error
	g.outfile, err = os.OpenFile(g.tts.OutputFile.String(), os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	codegen.LogFatal(err)
	g.appendToFile(tpl.Copy, struct{ Package string }{Package: g.tts.Package}, nil)

	g.initTables()
	g.runHeader()
	g.runTable()
	g.runEAValueTables()
	codegen.LogFatal(g.outfile.Close())
}

func (g *generator) setMagentoVersion(magento1, magento2 bool) *generator {
	g.isMagento1 = magento1
	g.isMagento2 = magento2
	return g
}

// analyzePackage extracts from all types the method receivers and type names. If we found existing
// functions we will add a MethodRecvPrefix to the generated functions to avoid conflicts.
func (g *generator) analyzePackage() {
	defer log.WhenDone(PkgLog).Info("Stats", "Package", g.tts.Package, "Step", "AnalyzePackage")
	fset := token.NewFileSet()

	path := filepath.Dir(g.tts.OutputFile.String())
	pkgs, err := parser.ParseDir(fset, path, nil, parser.AllErrors)
	codegen.LogFatal(err)

	var astPkg *ast.Package
	var ok bool
	if astPkg, ok = pkgs[g.tts.Package]; !ok {
		fmt.Printf("Package %s not found in path %s. Skipping.", g.tts.Package, path)
		return
	}

	for fName, astFile := range astPkg.Files {
		if fName == g.tts.OutputFile.String() || strings.Contains(fName, "_fallback.go") {
			// skip the generated file (tables_generated.go) or we have recursion 8-)
			// skip also the _fallback.go files which contain the generated code
			// because those files get included if no build tag has been specified.
			continue
		}
		ast.Inspect(astFile, func(n ast.Node) bool {
			switch stmt := n.(type) {
			case *ast.FuncDecl:
				if stmt.Recv != nil { // we have a method receiver and not a normal function
					switch t := stmt.Recv.List[0].Type.(type) {
					case *ast.Ident: // non-pointer-type
						if strings.Index(t.Name, TypePrefix) == 0 {
							g.existingMethodSets.add(t.Name + stmt.Name.Name) // e.g.: TableWebsiteSliceLoad where Load is the function name
						}
					case *ast.StarExpr: // pointer-type
						switch t2 := t.X.(type) {
						case *ast.Ident:
							if strings.Index(t2.Name, TypePrefix) == 0 {
								g.existingMethodSets.add(t2.Name + stmt.Name.Name) // e.g.: *TableWebsiteSliceLoad where Load is the function name
							}
						}
					}
				}
			}
			return true
		})
	}
}

func (g *generator) appendToFile(tpl string, data interface{}, addFM template.FuncMap) {
	formatted, err := codegen.GenerateCode(g.tts.Package, tpl, data, addFM)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		codegen.LogFatal(err)
	}

	if _, err := g.outfile.Write(formatted); err != nil {
		codegen.LogFatal(err)
	}
	codegen.LogFatal(g.outfile.Sync()) // flush immediately to disk to prevent a race condition
}

func (g *generator) initTables() {
	defer log.WhenDone(PkgLog).Info("Stats", "Package", g.tts.Package, "Step", "InitTables")
	var err error
	g.tables, err = codegen.GetTables(g.dbrConn.NewSession(), codegen.ReplaceTablePrefix(g.tts.SQLQuery))
	codegen.LogFatal(err)

	if len(g.tts.EntityTypeCodes) > 0 && g.tts.EntityTypeCodes[0] != "" {
		g.eavValueTables, err = codegen.GetEavValueTables(g.dbrConn, g.tts.EntityTypeCodes)
		codegen.LogFatal(err)

		for _, vTables := range g.eavValueTables {
			for t := range vTables {
				if false == isDuplicate(g.tables, t) {
					g.tables = append(g.tables, t)
				}
			}
		}
	}

	if g.tts.GenericsWhiteList == "" {
		return // do nothing because nothing defined, neither custom SQL nor to copy from SQLQuery field
	}
	if false == dbr.Stmt.IsSelect(g.tts.GenericsWhiteList) {
		// copy result from tables because select key word not found
		g.whiteListTables = g.tables
		return
	}

	g.whiteListTables, err = codegen.GetTables(g.dbrConn.NewSession(), codegen.ReplaceTablePrefix(g.tts.GenericsWhiteList))
	codegen.LogFatal(err)
}

func (g *generator) runHeader() {
	defer log.WhenDone(PkgLog).Info("Stats", "Package", g.tts.Package, "Step", "RunHeader")
	type Table struct {
		Name      string
		TableName string
	}

	data := struct {
		Package, Tick          string
		HasTypeCodeValueTables bool
		Tables                 []Table
	}{
		Package: g.tts.Package,
		Tick:    "`",
		HasTypeCodeValueTables: len(g.eavValueTables) > 0,
	}

	for _, table := range g.tables {
		data.Tables = append(data.Tables, Table{
			Name:      g.getConsistentName(table),
			TableName: g.getMagento2TableName(table),
		})
	}
	g.appendToFile(tpl.Header, data, nil)
}

func (g *generator) runTable() {
	defer log.WhenDone(PkgLog).Info("Stats", "Package", g.tts.Package, "Step", "RunTable")
	type OneTable struct {
		Package          string
		Tick             string
		Name             string
		TableName        string
		Struct           string
		Slice            string
		Table            string
		GoColumns        codegen.Columns
		Columns          csdb.Columns
		MethodRecvPrefix string
		FindByPk         string
	}

	for _, table := range g.tables {

		columns, err := codegen.GetColumns(g.dbrConn.DB, table)
		codegen.LogFatal(err)
		codegen.LogFatal(columns.MapSQLToGoDBRType())

		name := g.getConsistentName(table)
		data := OneTable{
			Package:   g.tts.Package,
			Tick:      "`",
			Name:      name,
			TableName: g.getMagento2TableName(table), // original table name!
			Struct:    TypePrefix + name,             // getTableConstantName
			Slice:     TypePrefix + name + "Slice",   // getTableConstantName
			Table:     table,
			GoColumns: columns,
			Columns:   columns.CopyToCSDB(),
		}

		if data.Columns.PrimaryKeys().Len() > 0 {
			data.FindByPk = "FindBy" + util.UnderscoreCamelize(data.Columns.PrimaryKeys().JoinFields("_"))
		}

		tplFuncs := template.FuncMap{
			"typePrefix": func(name string) string {
				// if the method already exists in package then add the prefix parent
				// to avoid duplicate function names.
				search := data.Slice + name
				if g.existingMethodSets.has(search) {
					return MethodRecvPrefix + name
				}
				return name
			},
			"findBy":  findBy,
			"dbrType": dbrType,
		}

		g.appendToFile(g.getGenericTemplate(table), data, tplFuncs)
	}
}

func (g *generator) getGenericTemplate(tableName string) string {
	var finalTpl bytes.Buffer

	// at least we need a type definition
	if _, err := finalTpl.WriteString(tpl.Type); err != nil {
		codegen.LogFatal(err)
	}

	if false == g.whiteListTables.Include(tableName) {
		return finalTpl.String()
	}
	isAll := (g.tts.GenericsFunctions & tpl.OptAll) == tpl.OptAll

	if isAll || (g.tts.GenericsFunctions&tpl.OptSQL) == tpl.OptSQL {
		_, err := finalTpl.WriteString(tpl.SQL)
		codegen.LogFatal(err)
	}
	if isAll || (g.tts.GenericsFunctions&tpl.OptFindBy) == tpl.OptFindBy {
		_, err := finalTpl.WriteString(tpl.FindBy)
		codegen.LogFatal(err)
	}
	if isAll || (g.tts.GenericsFunctions&tpl.OptSort) == tpl.OptSort {
		_, err := finalTpl.WriteString(tpl.Sort)
		codegen.LogFatal(err)
	}
	if isAll || (g.tts.GenericsFunctions&tpl.OptSliceFunctions) == tpl.OptSliceFunctions {
		_, err := finalTpl.WriteString(tpl.SliceFunctions)
		codegen.LogFatal(err)
	}
	if isAll || (g.tts.GenericsFunctions&tpl.OptExtractFromSlice) == tpl.OptExtractFromSlice {
		_, err := finalTpl.WriteString(tpl.ExtractFromSlice)
		codegen.LogFatal(err)
	}
	return finalTpl.String()
}

func (g *generator) runEAValueTables() {
	if len(g.eavValueTables) == 0 {
		return
	}
	defer log.WhenDone(PkgLog).Info("Stats", "Package", g.tts.Package, "Step", "RunEAValueTables")

	data := struct {
		TypeCodeValueTables codegen.TypeCodeValueTable
	}{
		TypeCodeValueTables: g.eavValueTables,
	}

	g.appendToFile(tpl.EAValueStructure, data, nil)
}

// protects map TableMapMagento1To2. Code smell global variable to protect a global var.
var mapm1m2Mu sync.Mutex

// getMagento2TableName takes care of the correct table name as some tables
// have different names in Magento2 than in Magento1.
func (g *generator) getMagento2TableName(table string) string {
	mapm1m2Mu.Lock()
	defer mapm1m2Mu.Unlock()

	if g.isMagento1 && !g.isMagento2 {
		return table
	}
	if mappedName, ok := codegen.TableMapMagento1To2[strings.Replace(table, codegen.TablePrefix, "", 1)]; ok {
		return mappedName
	}

	return table
}

// getConsistentName generates a consistent name for all TableIndex*, Table*Slice
// and Table* types. This names is guaranteed to be the same whether we run
// Magento 1 or 2.
func (g *generator) getConsistentName(table string) string {
	mapm1m2Mu.Lock()
	defer mapm1m2Mu.Unlock()
	// 1. retrieve the mapped name used in Magento2
	if mappedName, ok := codegen.TableMapMagento1To2[strings.Replace(table, codegen.TablePrefix, "", 1)]; ok {
		table = mappedName
	}
	// 2. Remove the package name from the table name
	return codegen.PrepareVar(g.tts.Package, table)
}
