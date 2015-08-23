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
	"io"
	"os"
	"strings"
	"sync"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/codegen/codecgen"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"regexp"
	"runtime"
)

var mapm1m2Mu sync.Mutex // protects map TableMapMagento1To2

func main() {
	dbc, err := csdb.Connect()
	codegen.LogFatal(err)
	defer dbc.Close()
	var wg sync.WaitGroup
	fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
	for _, tStruct := range codegen.ConfigTableToStruct {
		go newGenerator(tStruct, dbc, &wg).run()
	}
	fmt.Printf("CPUs: %d\tGoroutines: %d\tGo Version %s\n", runtime.NumCPU(), runtime.NumGoroutine(), runtime.Version())
	wg.Wait()
}

type generator struct {
	tStruct        *codegen.TableToStruct
	dbrConn        *dbr.Connection
	w              io.WriteCloser
	tables         []string
	eavValueTables codegen.TypeCodeValueTable
	wg             *sync.WaitGroup
}

func newGenerator(tStruct *codegen.TableToStruct, dbrConn *dbr.Connection, wg *sync.WaitGroup) *generator {
	wg.Add(1)
	g := &generator{
		tStruct: tStruct,
		dbrConn: dbrConn,
		wg:      wg,
	}
	var err error
	g.w, err = os.OpenFile(g.tStruct.OutputFile.String(), os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	codegen.LogFatal(err)

	g.tables, g.eavValueTables = g.initTables()
	return g
}

func (g *generator) run() {
	defer g.wg.Done()
	g.runHeader()
	g.runTable()
	g.runEAValueTables()
	g.runCodec()
	codegen.LogFatal(g.w.Close())
}

func (g *generator) appendToFile(tpl string, data interface{}) {

	formatted, err := codegen.GenerateCode(g.tStruct.Package, tpl, data, nil)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		codegen.LogFatal(err)
	}

	if _, err := g.w.Write(formatted); err != nil {
		codegen.LogFatal(err)
	}
}

func (g *generator) initTables() ([]string, codegen.TypeCodeValueTable) {
	tables, err := codegen.GetTables(g.dbrConn.DB, codegen.ReplaceTablePrefix(g.tStruct.SQLQuery))
	codegen.LogFatal(err)

	var eavValueTables codegen.TypeCodeValueTable
	if len(g.tStruct.EntityTypeCodes) > 0 && g.tStruct.EntityTypeCodes[0] != "" {
		var err error
		eavValueTables, err = codegen.GetEavValueTables(g.dbrConn, g.tStruct.EntityTypeCodes)
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
		Package: g.tStruct.Package,
		Tick:    "`",
		HasTypeCodeValueTables: len(g.eavValueTables) > 0,
	}

	for _, table := range g.tables {
		var name = getTableName(table)
		data.Tables = append(data.Tables, Table{name, codegen.PrepareVar(g.tStruct.Package, name)})
	}
	g.appendToFile(tplHeader, data)
}

func (g *generator) runTable() {

	type OneTable struct {
		Package string
		Tick    string
		NameRaw string
		Name    string
		Table   string
		Columns codegen.Columns
	}

	for _, table := range g.tables {

		columns, err := codegen.GetColumns(g.dbrConn.DB, table)
		codegen.LogFatal(err)
		codegen.LogFatal(columns.MapSQLToGoDBRType())

		var name = getTableName(table)
		data := OneTable{
			Package: g.tStruct.Package,
			Tick:    "`",
			NameRaw: name,
			Name:    codegen.PrepareVar(g.tStruct.Package, name),
			Table:   table,
			Columns: columns,
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

func (g *generator) runCodec() {

	if err := codecgen.Generate(
		g.tStruct.OutputFile.AppendName("_codec").String(), // outfile
		"", // buildTag
		codecgen.GenCodecPath,
		false, // use unsafe
		"",
		regexp.MustCompile("Table.*"), // Prefix of generated structs and slices
		true, // delete temp files
		g.tStruct.OutputFile.String(), // read from file
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
