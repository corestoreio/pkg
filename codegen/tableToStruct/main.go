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
	"regexp"
	"runtime"
	"strings"
	"sync"

	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/codegen/codecgen"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils"
)

var mapm1m2Mu sync.Mutex // protects map TableMapMagento1To2
const MethodRecvPrefix = "parent"

// TypePrefix of the generated types e.g. TableStoreSlice, TableStore ...
// If you change this you must change all "Table" in the template.
const TypePrefix = "Table"

// generatedFunctions within the template if a package has already such a function
// then prefix MethodRecvPrefix to the generating the function so that in our code we
// can refer to the "parent" function. No composition possible.
var generatedFunctions = []string{"Load", "Len"}

func main() {
	dbc, err := csdb.Connect()
	codegen.LogFatal(err)
	defer dbc.Close()
	var wg sync.WaitGroup
	fmt.Printf("CPUs: %d\tGoroutines: %d\n", runtime.NumCPU(), runtime.NumGoroutine())
	for _, tStruct := range codegen.ConfigTableToStruct {
		go newGenerator(tStruct, dbc, &wg).run()
	}
	fmt.Printf("Goroutines: %d\tGo Version %s\n", runtime.NumGoroutine(), runtime.Version())
	wg.Wait()

	//	for _, ts := range codegen.ConfigTableToStruct {
	//		// due to a race condition the codec generator must run after the newGenerator() calls
	//		runCodec(ts.OutputFile.AppendName("_codec").String(), ts.OutputFile.String())
	//	}

}

type generator struct {
	tts            *codegen.TableToStruct
	dbrConn        *dbr.Connection
	outfile        *os.File
	tables         []string
	eavValueTables codegen.TypeCodeValueTable
	wg             *sync.WaitGroup
	existingFnc    utils.StringSlice // "receiver name"."function name" if in this slice then the generated
	// function will be private and you can refer from your function to the generated one
	// because sometimes e.g. Load() needs more SQL
}

func newGenerator(tts *codegen.TableToStruct, dbrConn *dbr.Connection, wg *sync.WaitGroup) *generator {
	wg.Add(1)
	g := &generator{
		tts:     tts,
		dbrConn: dbrConn,
		wg:      wg,
	}
	g.analyzePackage()

	var err error
	g.outfile, err = os.OpenFile(g.tts.OutputFile.String(), os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	codegen.LogFatal(err)
	g.appendToFile(tplCopyPkg, struct{ Package string }{Package: g.tts.Package})

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

// analyzePackage extracts from all types the method receivers and type names
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
			// skip the generated file
			continue
		}
		ast.Inspect(astFile, func(n ast.Node) bool {
			switch stmt := n.(type) {
			case *ast.FuncDecl:
				if stmt.Recv != nil { // we have a method receiver and not a normal function
					switch t := stmt.Recv.List[0].Type.(type) {
					case *ast.Ident: // non-pointer-type
						if strings.Index(t.Name, TypePrefix) == 0 {
							g.existingFnc.Append(t.Name + "." + stmt.Name.Name) // e.g.: TableWebsiteSliceLoad where Load is the function name
						}
					case *ast.StarExpr: // pointer-type
						switch t2 := t.X.(type) {
						case *ast.Ident:
							if strings.Index(t2.Name, TypePrefix) == 0 {
								g.existingFnc.Append(t2.Name + "." + stmt.Name.Name) // e.g.: *TableWebsiteSliceLoad where Load is the function name
							}
						}
					}
				}
			}
			return true
		})
	}
}

func (g *generator) appendToFile(tpl string, data interface{}) {

	formatted, err := codegen.GenerateCode(g.tts.Package, tpl, data, nil)
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
	g.appendToFile(tplHeader, data)
}

func (g *generator) runTable() {

	type OneTable struct {
		Package          string
		Tick             string
		NameRaw          string
		Name             string
		Table            string
		GoColumns        codegen.Columns
		Columns          csdb.Columns
		MethodRecvPrefix string
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
			Name:      codegen.PrepareVar(g.tts.Package, name),
			Table:     table,
			GoColumns: columns,
			Columns:   columns.CopyToCSDB(),
		}

		if g.existingFnc.Len() > 0 {
			has := false
			for _, fnn := range generatedFunctions {
				tsl := TypePrefix + data.Name + "Slice." + fnn
				tn := TypePrefix + data.Name + "." + fnn
				if g.existingFnc.Include(tsl) || g.existingFnc.Include(tn) {
					has = true
					break
				}
			}
			if has {
				data.MethodRecvPrefix = MethodRecvPrefix
			}
		}
		g.appendToFile(tplTable, data)
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

	g.appendToFile(tplEAValues, data)
}

// runCodec generates the codecs to be used later in JSON or msgpack or etc
func runCodec(outfile, readfile string) {

	if err := codecgen.Generate(
		outfile, // outfile
		"",      // buildTag
		codecgen.GenCodecPath,
		false, // use unsafe
		"",
		regexp.MustCompile(TypePrefix+".*"), // Prefix of generated structs and slices
		true,     // delete temp files
		readfile, // read from file
	); err != nil {
		fmt.Println("codecgen.Generate Error:")
		codegen.LogFatal(err)
	}
}

// isDuplicate slow duplicate checker ...
func isDuplicate(sl []string, st string) bool {
	for _, s := range sl {
		if s == st {
			return true
		}
	}
	return false
}

func getTableName(table string) (name string) {
	mapm1m2Mu.Lock()
	name = table
	if mappedName, ok := codegen.TableMapMagento1To2[strings.Replace(table, codegen.TablePrefix, "", 1)]; ok {
		name = mappedName
	}
	mapm1m2Mu.Unlock()
	return
}
