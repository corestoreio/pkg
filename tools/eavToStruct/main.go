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

// Generates code for all EAV types
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"encoding/json"

	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/tools"
	_ "github.com/go-sql-driver/mysql"
	"github.com/juju/errgo"
)

const (
	envTableMap string = "CS_EAV_MAP"
)

var (
	pkg        = flag.String("p", "", "Package name in template")
	run        = flag.Bool("run", false, "If true program runs")
	outputFile = flag.String("o", "", "Output file name")
)

type (
	JsonEntityTypeMap map[string]*tools.EntityTypeMap
)

func main() {
	flag.Parse()

	if false == *run || *outputFile == "" || *pkg == "" {
		flag.Usage()
		os.Exit(1)
	}

	db, dbrConn, err := csdb.Connect()
	tools.LogFatal(err)
	defer db.Close()

	type dataContainer struct {
		ETypeData     JsonEntityTypeMap
		Package, Tick string
	}

	etData, err := getEntityTypeData(dbrConn.NewSession(nil))
	tools.LogFatal(err)

	tplData := &dataContainer{
		ETypeData: etData,
		Package:   *pkg,
		Tick:      "`",
	}

	formatted, err := tools.GenerateCode(*pkg, tplEav, tplData)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		tools.LogFatal(err)
	}

	ioutil.WriteFile(*outputFile, formatted, 0600)
}

func getEntityTypeData(dbrSess *dbr.Session) (JsonEntityTypeMap, error) {

	s, err := eav.GetTableStructure(eav.TableEntityType)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	var entityTypeCollection eav.EntityTypeSlice
	_, err = dbrSess.
		Select(s.AllColumnAliasQuote(s.Name)...).
		From(s.Name).
		LoadStructs(&entityTypeCollection)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	mapCollection, err := getMapping(os.Getenv(envTableMap), tools.JSONMappingEntityTypes)
	tools.LogFatal(err)

	for typeCode, mapData := range mapCollection {
		// now map the values from entityTypeCollection into mapData
		et, err := entityTypeCollection.GetByCode(typeCode)
		tools.LogFatal(err)
		mapData.EntityTypeID = et.EntityTypeID
		mapData.EntityTypeCode = et.EntityTypeCode
		mapData.ValueTablePrefix = et.ValueTablePrefix.String
		mapData.EntityIDField = et.EntityIDField.String
		mapData.IsDataSharing = et.IsDataSharing
		mapData.DataSharingKey = et.DataSharingKey.String
		mapData.DefaultAttributeSetID = et.DefaultAttributeSetID
		mapData.IncrementPerStore = et.IncrementPerStore
		mapData.IncrementPadLength = et.IncrementPadLength
		mapData.IncrementPadChar = et.IncrementPadChar
	}

	return mapCollection, nil
}

func getMapping(fileName string, rawJson []byte) (JsonEntityTypeMap, error) {
	var err error
	if fileName != "" && fileName[len(fileName)-5:] == ".json" { // check if file ext is .json
		rawJson, err = ioutil.ReadFile(fileName)
		if err != nil {
			return nil, errgo.Mask(err)
		}
	}
	mapping := make(JsonEntityTypeMap)
	err = json.Unmarshal(rawJson, &mapping)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return mapping, nil
}
