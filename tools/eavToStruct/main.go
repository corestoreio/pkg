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
	"github.com/corestoreio/csfw/tools"
	"github.com/corestoreio/csfw/tools/toolsdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/juju/errgo"
)

var (
	pkg        = flag.String("p", "", "Package name in template")
	run        = flag.Bool("run", false, "If true program runs")
	outputFile = flag.String("o", "", "Output file name")
	// @todo maybe use also an ENV var for the tableMap
	tableMap = flag.String("map", "", "JSON file for mapping entity types to real table names and Go interfaces. If empty fall back to default mapping.")
)

type (
	JsonEntityTypeMap map[string]*EntityTypeMap
	EntityTypeMap     struct {
		EntityTypeID              int64
		EntityTypeCode            string
		EntityModel               string `json:"entity_model"`
		AttributeModel            string `json:"attribute_model"`
		EntityTable               string `json:"entity_table"`
		ValueTablePrefix          string
		EntityIDField             string
		IsDataSharing             bool
		DataSharingKey            string
		DefaultAttributeSetID     int64
		IncrementModel            string `json:"increment_model"`
		IncrementPerStore         bool
		IncrementPadLength        int64
		IncrementPadChar          string
		AdditionalAttributeTable  string `json:"additional_attribute_table"`
		EntityAttributeCollection string `json:"entity_attribute_collection"`
	}
)

func main() {
	flag.Parse()

	if false == *run || *outputFile == "" || *pkg == "" {
		flag.Usage()
		os.Exit(1)
	}

	db, dbrSess, err := csdb.Connect()
	toolsdb.LogFatal(err)
	defer db.Close()

	type dataContainer struct {
		ETypeData     JsonEntityTypeMap
		Package, Tick string
	}

	etData, err := getEntityTypeData(dbrSess)
	toolsdb.LogFatal(err)

	tplData := &dataContainer{
		ETypeData: etData,
		Package:   *pkg,
		Tick:      "`",
	}

	formatted, err := tools.GenerateCode(tplEav, tplData)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		toolsdb.LogFatal(err)
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
		Select(s.Columns...).
		From(s.Name).
		LoadStructs(&entityTypeCollection)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	mapCollection, err := getMapping(*tableMap, defaultMapping)
	toolsdb.LogFatal(err)

	for typeCode, mapData := range mapCollection {
		// now map the values from entityTypeCollection into mapData
		et, err := entityTypeCollection.GetByCode(typeCode)
		toolsdb.LogFatal(err)
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
