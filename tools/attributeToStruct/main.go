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
		if err != nil {
			tools.LogFatal(err)
		}
		// create table from select, get all columns from that table
		// create the struct with interfaces

		// create the same as with CSEntityTypeSlice and CSEntityType but only for attribtues
		// auto create the structs containing the Go interfaces and then put in the data
		// write ann into the concrete package.

		//		structCode, err := tools.QueryToStruct(db, et.EntityTypeCode+"EavAttributeSelect", dbrSelect)
		//		if err != nil {
		//			tools.LogFatal(err)
		//		}
		//		fmt.Printf("\n%s\n", structCode)
		// now aggregate structCode and write then all into the generated files in a package
		// use the data from JSON mapping
	}

	//ioutil.WriteFile(*outputFile, formatted, 0600)
}

/*
to retrieve the attributes. The eav library must implement:

EAV -> Create queries for AttributeSets and AttributeGroups
    ->
*/
