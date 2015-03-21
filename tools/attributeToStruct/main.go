// Copyright 2015 CoreStore Authors
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

// Generates code for all EAV attribute types
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/corestoreio/csfw/concrete"
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/tools"
	"github.com/juju/errgo"
)

var (
	//pkg        = flag.String("p", "", "Package name in template")
	run        = flag.Bool("run", false, "If true program runs")
	outputFile = flag.String("o", "", "Output file name")
)

const (
	envModelMap string = "CS_ATTRIBUTE_MODEL_MAP"
)

type context struct {
	db       *sql.DB
	dbrConn  *dbr.Connection
	et       *eav.CSEntityType // will be updated each iteration
	modelMap tools.AttributeModelMap
}

func newContext() *context {
	db, dbrConn, err := csdb.Connect()
	tools.LogFatal(err)
	modelMap, err := getMapping(os.Getenv(envModelMap), tools.JSONMapAttributeModels)
	tools.LogFatal(err)

	return &context{
		db:       db,
		dbrConn:  dbrConn,
		modelMap: modelMap,
	}
}

func main() {
	flag.Parse()

	if false == *run || *outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx := newContext()
	defer ctx.db.Close()

	for _, et := range concrete.CSEntityTypeCollection {
		ctx.et = et
		code, err := prepareAttributeCode(ctx)
		//tools.LogFatal(err)
		if err != nil {
			fmt.Printf("\n%s\n", err)
		}
		fmt.Printf("\n%s\n+++++++++++++++++++++++++++++++++++++++++++++++++\n", code)
		// auto create the structs containing the Go interfaces and then put in the data
		// write ann into the concrete package.

		// now aggregate structCode and write then all into the generated files in a package
		// use the data from JSON mapping
		// EAV -> Create queries for AttributeSets and AttributeGroups
	}

	//ioutil.WriteFile(*outputFile, formatted, 0600)
}

func prepareAttributeCode(ctx *context) ([]byte, error) {
	dbrSelect, err := eav.GetAttributeSelectSql(ctx.dbrConn.NewSession(nil), ctx.et, 0)
	if err != nil {
		return nil, err
	}

	columns, err := tools.SQLQueryToColumns(ctx.db, dbrSelect)
	if err != nil {
		return nil, err
	}

	tools.LogFatal(columns.MapSQLToGoType(tools.EavAttributeColumnNameToInterface))
	structName := "cs_" + ctx.et.EntityTypeCode
	structCode, err := tools.ColumnsToStructCode(structName, columns, tplQueryStruct)
	if err != nil {
		return nil, err
	}

	attributeCollection, err := tools.GetSQL(ctx.db, dbrSelect)
	if err != nil {
		return nil, err
	}

	tools.PrepareForTemplate(columns, attributeCollection, ctx.modelMap)

	// iterate over attributeCollection and escape or not the values to be used as string, int, bool or Go func

	data := struct {
		QueryStruct string
		Attributes  []tools.StringEntities
		Name        string
	}{
		QueryStruct: string(structCode),
		Attributes:  attributeCollection,
		Name:        structName,
	}

	return tools.GenerateCode("packageNameTODO", tplQueryData, data)
}

func getMapping(fileName string, rawJson []byte) (tools.AttributeModelMap, error) {
	var err error
	if fileName != "" && fileName[len(fileName)-5:] == ".json" { // check if file ext is .json
		rawJson = nil
		rawJson, err = ioutil.ReadFile(fileName)
		if err != nil {
			return nil, errgo.Mask(err)
		}
	}
	mapping := make(tools.AttributeModelMap)
	err = json.Unmarshal(rawJson, &mapping)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return mapping, nil
}
