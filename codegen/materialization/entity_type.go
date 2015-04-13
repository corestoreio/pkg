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

package main

import (
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils"
	"github.com/juju/errgo"
)

// materializeEntityType writes the data from eav_entity_type into a Go file and transforms
// Magento classes and config strings into Go functions.
// Depends on generated code from tableToStruct.
func materializeEntityType(ctx *context) {
	defer ctx.wg.Done()
	type dataContainer struct {
		ETypeData     eav.TableEntityTypeSlice
		ImportPaths   []string
		Package, Tick string
	}

	etData, err := getEntityTypeData(ctx.dbrConn.NewSession(nil))
	codegen.LogFatal(err)

	tplData := &dataContainer{
		ETypeData:   etData,
		ImportPaths: getImportPaths(),
		Package:     codegen.ConfigMaterializationEntityType.Package,
		Tick:        "`",
	}

	addFM := template.FuncMap{
		"extractFuncType": codegen.ExtractFuncType,
	}

	formatted, err := codegen.GenerateCode(codegen.ConfigMaterializationEntityType.Package, tplEav, tplData, addFM)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		codegen.LogFatal(err)
	}

	codegen.LogFatal(ioutil.WriteFile(codegen.ConfigMaterializationEntityType.OutputFile, formatted, 0600))
}

// getEntityTypeData retrieves all EAV models from table eav_entity_type but only those listed in variable
// codegen.ConfigEntityType. It then applies the mapping data from codegen.ConfigEntityType to the entity_type struct.
// Depends on generated code from tableToStruct.
func getEntityTypeData(dbrSess *dbr.Session) (etc eav.TableEntityTypeSlice, err error) {

	s, err := eav.GetTableStructure(eav.TableIndexEntityType)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	_, err = dbrSess.
		Select(s.AllColumnAliasQuote(s.Name)...).
		From(s.Name).
		Where("entity_type_code IN ?", codegen.ConfigEntityType.Keys()).
		LoadStructs(&etc)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	for typeCode, mapData := range codegen.ConfigEntityType {
		// map the fields from the config struct to the data retrieved from the database.
		et, err := etc.GetByCode(typeCode)
		codegen.LogFatal(err)
		et.EntityModel = codegen.ParseString(mapData.EntityModel, et)
		et.AttributeModel.String = codegen.ParseString(mapData.AttributeModel, et)
		et.EntityTable.String = codegen.ParseString(mapData.EntityTable, et)
		et.IncrementModel.String = codegen.ParseString(mapData.IncrementModel, et)
		et.AdditionalAttributeTable.String = codegen.ParseString(mapData.AdditionalAttributeTable, et)
		et.EntityAttributeCollection.String = codegen.ParseString(mapData.EntityAttributeCollection, et)
	}

	return etc, nil
}

func getImportPaths() []string {
	var paths utils.stringSlice

	var getPath = func(s string) string {
		ps, err := codegen.ExtractImportPath(s)
		codegen.LogFatal(err)
		return ps
	}

	for _, et := range codegen.ConfigEntityType {
		paths.Append(
			getPath(et.EntityModel),
			getPath(et.AttributeModel),
			getPath(et.EntityTable),
			getPath(et.IncrementModel),
			getPath(et.AdditionalAttributeTable),
			getPath(et.EntityAttributeCollection),
		)
	}
	return paths.Unique().ToString()
}
