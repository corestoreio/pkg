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
	"fmt"
	"os"
	"strings"
	"sync"

	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"text/template"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
)

type generator struct {
	tts            codegen.TableToStruct
	dbrConn        *dbr.Connection
	outfile        *os.File
	tables         []string
	eavValueTables codegen.TypeCodeValueTable
	wg             *sync.WaitGroup
	// existingMethodSets contains all existing method sets for the Table* types
	existingMethodSets *duplicateChecker
}

func newGenerator(tts codegen.TableToStruct, dbrConn *dbr.Connection, wg *sync.WaitGroup) *generator {
	wg.Add(1)
	g := &generator{
		tts:                tts,
		dbrConn:            dbrConn,
		wg:                 wg,
		existingMethodSets: newDuplicateChecker(),
	}
	g.analyzePackage()

	var err error
	g.outfile, err = os.OpenFile(g.tts.OutputFile.String(), os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	codegen.LogFatal(err)
	g.appendToFile(tplCopyPkg, struct{ Package string }{Package: g.tts.Package}, nil)

	g.tables, g.eavValueTables = g.initTables()
	return g
}

func (g *generator) run() {
	defer g.wg.Done()
	g.runHeader()
	g.runTable()
	g.runEAValueTables()
	codegen.LogFatal(g.outfile.Close())
}

// analyzePackage extracts from all types the method receivers and type names. If we found existing
// functions we will add a MethodRecvPrefix to the generated functions to avoid conflicts.
func (g *generator) analyzePackage() {
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
		if fName == g.tts.OutputFile.String() {
			// skip the generated file or we have recursion 8-)
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

func (g *generator) initTables() ([]string, codegen.TypeCodeValueTable) {
	tables, err := codegen.GetTables(g.dbrConn.DB, codegen.ReplaceTablePrefix(g.tts.SQLQuery))
	codegen.LogFatal(err)

	var eavValueTables codegen.TypeCodeValueTable
	if len(g.tts.EntityTypeCodes) > 0 && g.tts.EntityTypeCodes[0] != "" {
		var err error
		eavValueTables, err = codegen.GetEavValueTables(g.dbrConn, g.tts.EntityTypeCodes)
		codegen.LogFatal(err)

		for _, vTables := range eavValueTables {
			for t := range vTables {
				if false == isDuplicate(tables, t) {
					tables = append(tables, t)
				}
			}
		}
	}
	return tables, eavValueTables
}

func (g *generator) runHeader() {
	type Table struct {
		NameRaw string
		Name    string
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
		var name = getTableName(table)
		data.Tables = append(data.Tables, Table{name, codegen.PrepareVar(g.tts.Package, name)})
	}
	g.appendToFile(tplHeader, data, nil)
}

func (g *generator) runTable() {

	type OneTable struct {
		Package          string
		Tick             string
		NameRaw          string
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

		var name = getTableName(table)
		data := OneTable{
			Package:   g.tts.Package,
			Tick:      "`",
			NameRaw:   name,
			Struct:    TypePrefix + codegen.PrepareVar(g.tts.Package, name),
			Slice:     TypePrefix + codegen.PrepareVar(g.tts.Package, name) + "Slice",
			Table:     table,
			GoColumns: columns,
			Columns:   columns.CopyToCSDB(),
		}

		if data.Columns.PrimaryKeys().Len() > 0 {
			data.FindByPk = "FindBy" + codegen.Camelize(data.Columns.PrimaryKeys().JoinFields("_"))
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
			"findBy": func(s string) string {
				return "FindBy" + codegen.Camelize(s)
			},
			"dbrType": func(c csdb.Column) string {
				switch {
				case false == c.IsNull():
					return ""
				case c.IsString():
					return ".String" // dbr.NullString
				case c.IsMoney():
					return "" // money.Currency
				case c.IsFloat():
					return ".Float64" // dbr.NullFloat64
				case c.IsInt():
					return ".Int64" // dbr.NullInt64
				}
				return ""
			},
		}

		g.appendToFile(tplTable, data, tplFuncs)
	}
}

func (g *generator) runEAValueTables() {
	if len(g.eavValueTables) == 0 {
		return
	}

	data := struct {
		TypeCodeValueTables codegen.TypeCodeValueTable
	}{
		TypeCodeValueTables: g.eavValueTables,
	}

	g.appendToFile(tplEAValues, data, nil)
}
