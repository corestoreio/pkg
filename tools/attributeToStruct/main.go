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
	"flag"
	"fmt"
	"os"

	"github.com/corestoreio/csfw/concrete"
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/tools"
)

var (
	//pkg        = flag.String("p", "", "Package name in template")
	run        = flag.Bool("run", false, "If true program runs")
	outputFile = flag.String("o", "", "Output file name")
)

func main() {
	flag.Parse()

	if false == *run || *outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	db, dbrConn, err := csdb.Connect()
	tools.LogFatal(err)
	defer db.Close()

	// mapping see: tools.JSONMappingEntityTypes and tools.JSONMappingEAVAttributeModels

	for _, et := range concrete.CSEntityTypeCollection {
		dbrSelect, err := eav.GetAttributeSelectSql(dbrConn.NewSession(nil), et, 0)
		tools.LogFatal(err)

		c, err := tools.SQLQueryToColumns(db, dbrSelect)
		tools.LogFatal(err)

		tools.LogFatal(tools.MapSQLToGoType(c, tools.EavAttributeColumnNameToInterface))
		structName := "CS_" + et.EntityTypeCode
		structCode, err := tools.ColumnsToStructCode(structName, c, tplQueryStruct)
		tools.LogFatal(err)

		fmt.Printf("\n%s\n", structCode)

		attributeCollection, err := tools.GetSQL(db, dbrSelect)
		tools.LogFatal(err)

		// iterate over attributeCollection and escape or not the values to be used as string, int, bool or Go func

		data := struct {
			QueryStruct []byte
			Attributes  []tools.StringEntities
			Name        string
		}{
			QueryStruct: structCode,
			Attributes:  attributeCollection,
			Name:        structName,
		}

		code, err := tools.GenerateCode("packageNameTODO", tplQueryData, data)
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
