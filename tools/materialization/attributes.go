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
	"io/ioutil"
	"path"
	"strings"

	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/tools"
	"github.com/juju/errgo"
)

// materializeAttributes ...
// Depends on generated code from tableToStruct.
func materializeAttributes(ctx *context) {
	defer ctx.wg.Done()

	etc, err := getEntityTypeData(ctx.dbrConn.NewSession(nil))
	tools.LogFatal(err)
	for _, et := range etc {
		ctx.et = et
		tools.LogFatal(generateAttributeCode(ctx))
	}
}

func getImportPath(et *eav.EntityType) string {
	if etConfig, ok := tools.ConfigEntityType[et.EntityTypeCode]; ok {
		return etConfig.ImportPath
	}
	return ""
}

func getEAVPackage(et *eav.EntityType) string {
	if etConfig, ok := tools.ConfigMaterializationAttributes[et.EntityTypeCode]; ok {
		return etConfig.EAVPackage
	}
	return ""
}

func getOutputFile(et *eav.EntityType) string {
	if etConfig, ok := tools.ConfigMaterializationAttributes[et.EntityTypeCode]; ok {
		return etConfig.OutputFile
	}
	panic("You must specify an output file")
}

func getPackage(et *eav.EntityType) string {
	return path.Base(getImportPath(et))
}

// getName generates a nice struct name with a removed package name to avoid stutter but
// only removes the package name if the entity_type_code contains an underscore
// Depends on generated code from tableToStruct.
func getName(ctx *context, suffix ...string) string {
	structBaseName := ctx.et.EntityTypeCode
	if strings.Contains(ctx.et.EntityTypeCode, "_") {
		structBaseName = strings.Replace(ctx.et.EntityTypeCode, getPackage(ctx.et)+"_", "", -1)
	}
	return structBaseName + "_" + strings.Join(suffix, "_")
}

// Depends on generated code from tableToStruct.
func generateAttributeCode(ctx *context) error {

	dbrSelect, err := eav.GetAttributeSelectSql(
		ctx.dbrConn.NewSession(nil),
		tools.NewAddAttrTables(ctx.db, ctx.et.EntityTypeCode),
		ctx.et.EntityTypeID,
		0, // @todo get all website IDs
	)
	if err != nil {
		return err
	}
	dbrSelect.OrderDir("main_table.attribute_code", true)
	columns, err := tools.SQLQueryToColumns(ctx.db, dbrSelect)
	if err != nil {
		return err
	}

	tools.LogFatal(columns.MapSQLToGoType(tools.EavAttributeColumnNameToInterface))

	name := getName(ctx, "attribute")
	structCode, err := tools.ColumnsToStructCode(name, columns, tplTypeDefinition)
	if err != nil {
		println(string(structCode))
		return err
	}

	attributeCollection, err := tools.GetSQL(ctx.db, dbrSelect)
	if err != nil {
		return err
	}

	// @todo ValidateRules field must be converted from PHP serialized string to JSON
	pkg := getPackage(ctx.et)
	importPaths := tools.PrepareForTemplate(columns, attributeCollection, tools.ConfigAttributeModel, pkg)

	data := struct {
		TypeDefinition string
		Attributes     []tools.StringEntities
		Name           string
		ImportPaths    []string
		PackageName    string
		EAVPackage     string
	}{
		TypeDefinition: string(structCode),
		Attributes:     attributeCollection,
		Name:           name,
		ImportPaths:    importPaths,
		PackageName:    pkg,
		EAVPackage:     getEAVPackage(ctx.et),
	}

	code, err := tools.GenerateCode("", tplTypeDefinitionFile, data)
	if err != nil {
		println(string(code))
		return err
	}
	// @todo better path OR we expect the full path from the config
	//	path := fmt.Sprintf(
	//		"%s%s%sgenerated_%s.go",
	//		ctx.goSrcPath,
	//		getImportPath(ctx.et),
	//		string(os.PathSeparator),
	//		getName(ctx, "attribute"),
	//	)
	return errgo.Mask(ioutil.WriteFile(getOutputFile(ctx.et), code, 0600))
}
