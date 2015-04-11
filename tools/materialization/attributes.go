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
	"text/template"

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

func getAttrPkg(et *eav.TableEntityType) string {
	if etConfig, ok := tools.ConfigMaterializationAttributes[et.EntityTypeCode]; ok {
		return path.Base(etConfig.AttrPkgImp)
	}
	return ""
}

func getOutputFile(et *eav.TableEntityType) string {
	if etConfig, ok := tools.ConfigMaterializationAttributes[et.EntityTypeCode]; ok {
		return etConfig.OutputFile
	}
	panic("You must specify an output file")
}

func getPackage(et *eav.TableEntityType) string {
	if etConfig, ok := tools.ConfigMaterializationAttributes[et.EntityTypeCode]; ok {
		return etConfig.Package
	}
	panic("You must specify a package name")
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

// stripCoreAttributeColumns returns a copy of columns and removes all core/default eav_attribute columns
func stripCoreAttributeColumns(cols tools.Columns) tools.Columns {
	ret := make(tools.Columns, 0, len(cols))
	for _, col := range cols {
		if tools.EAVAttributeCoreColumns.Include(col.Field.String) {
			continue
		}
		f := false
		for _, et := range tools.ConfigEntityType {
			if et.AttributeCoreColumns.Include(col.Field.String) {
				f = true
				break
			}
		}
		if f == false {
			ret = append(ret, col)
		}
	}
	return ret
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
	typeTplData := map[string]interface{}{
		"AttrPkg":    getAttrPkg(ctx.et),
		"AttrStruct": tools.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].AttrStruct,
	}
	structCode, err := tools.ColumnsToStructCode(typeTplData, name, stripCoreAttributeColumns(columns), tplTypeDefinition)
	if err != nil {
		println(string(structCode))
		return err
	}

	attributeCollection, err := tools.LoadStringEntities(ctx.db, dbrSelect)
	if err != nil {
		return err
	}

	// @todo ValidateRules field must be converted from PHP serialized string to JSON
	pkg := getPackage(ctx.et)
	importPaths := tools.PrepareForTemplate(columns, attributeCollection, tools.ConfigAttributeModel, pkg)
	data := map[string]interface{}{
		"TypeDefinition": string(structCode),
		"Attributes":     attributeCollection,
		"Name":           name,
		"MyStruct":       tools.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].MyStruct,
		"ImportPaths":    importPaths,
		"PackageName":    pkg,
		"AttrPkg":        getAttrPkg(ctx.et),
		"AttrPkgImp":     tools.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].AttrPkgImp,
		"AttrStruct":     tools.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].AttrStruct,
		"FuncCollection": tools.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].FuncCollection,
		"FuncGetter":     tools.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].FuncGetter,
	}
	funcMap := template.FuncMap{
		// isEavAttr checks if the attribute/column name belongs to table eav_attribute
		"isEavAttr": func(a string) bool { return tools.EAVAttributeCoreColumns.Include(a) },
		// isEavEntityAttr checks if the attribute/column belongs to (customer|catalog|etc)_eav_attribute
		"isEavEntityAttr": func(a string) bool {
			if et, ok := tools.ConfigEntityType[ctx.et.EntityTypeCode]; ok {
				return false == tools.EAVAttributeCoreColumns.Include(a) && et.AttributeCoreColumns.Include(a)
			}
			return false
		},
		"isUnknownAttr": func(a string) bool {
			if et, ok := tools.ConfigEntityType[ctx.et.EntityTypeCode]; ok {
				return false == tools.EAVAttributeCoreColumns.Include(a) && false == et.AttributeCoreColumns.Include(a)
			}
			return false
		},
		"setAttrIdx": func(value, constName string) string {
			return strings.Replace(value, "{{.AttributeIndex}}", constName, -1)
		},
	}

	code, err := tools.GenerateCode("", tplTypeDefinitionFile, data, funcMap)
	if err != nil {
		println(string(code))
		return err
	}

	return errgo.Mask(ioutil.WriteFile(getOutputFile(ctx.et), code, 0600))
}
