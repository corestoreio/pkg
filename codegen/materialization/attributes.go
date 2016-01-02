// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"text/template"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/juju/errgo"
)

// materializeAttributes ...
// Depends on generated code from tableToStruct.
func materializeAttributes(ctx *context) {
	defer ctx.wg.Done()

	// generators, order of execution is important
	var gs = []func(*context, map[string]interface{}) ([]byte, error){
		attrCopyright,
		attrImport,
		attrTypes,
		attrGetter,
		attrCollection,
	}

	etc, err := getEntityTypeData(ctx.dbc.NewSession(nil))
	codegen.LogFatal(err)
	for _, et := range etc {
		ctx.et = et
		ctx.aat = codegen.NewAddAttrTables(ctx.dbc.DB, ctx.et.EntityTypeCode)
		data := attrGenerateData(ctx)
		var cb bytes.Buffer // code buffer
		for _, g := range gs {
			code, err := g(ctx, data)
			if err != nil {
				println(string(code))
				codegen.LogFatal(err)
			}
			cb.Write(code)
		}
		codegen.LogFatal(ioutil.WriteFile(getOutputFile(ctx.et), cb.Bytes(), 0600))
	}
}

func attrGenerateData(ctx *context) map[string]interface{} {
	websiteID := int64(0) // always 0 because we're building the base struct
	columns := getAttrColumns(ctx, websiteID)
	attributeCollection, err := codegen.LoadStringEntities(ctx.dbc.DB, getAttrSelect(ctx, websiteID))
	codegen.LogFatal(err)

	pkg := getPackage(ctx.et)
	importPaths := codegen.PrepareForTemplate(columns, attributeCollection, codegen.ConfigAttributeModel, pkg)

	return map[string]interface{}{
		"AttrCol":        attributeCollection,
		"AttrPkg":        getAttrPkg(ctx.et),
		"AttrPkgImp":     codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].AttrPkgImp,
		"AttrStruct":     codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].AttrStruct,
		"FuncCollection": codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].FuncCollection,
		"FuncGetter":     codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].FuncGetter,
		"ImportPaths":    importPaths,
		"MyStruct":       codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].MyStruct,
		"Name":           getStructName(ctx, "attribute"),
		"PackageName":    pkg,
	}
}

func attrCopyright(ctx *context, _ map[string]interface{}) ([]byte, error) {
	return codegen.Copyright, nil
}

func attrImport(ctx *context, data map[string]interface{}) ([]byte, error) {
	return codegen.GenerateCode("", tplAttrImport, data, nil)
}

func attrTypes(ctx *context, data map[string]interface{}) ([]byte, error) {
	columns := getAttrColumns(ctx, 0) // always zero websiteID
	return codegen.ColumnsToStructCode(data, data["Name"].(string), stripCoreAttributeColumns(columns), tplAttrTypes)
}

func attrGetter(ctx *context, data map[string]interface{}) ([]byte, error) {
	return codegen.GenerateCode("", tplAttrGetter, data, nil)
}

// getAttributeValuesForWebsites creates a map where the key is the attribute ID and
// each part of the StringEntities slice are the full attribute values for a website ID.
func getAttributeValuesForWebsites(ctx *context) map[string][]codegen.StringEntities {

	var tws store.TableWebsiteSlice
	tws.Load(ctx.dbc.NewSession(nil), func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.Where("website_id > 0")
	})

	// key contains the attributeID as a string
	var aws = make(map[string][]codegen.StringEntities)
	tew, err := ctx.aat.TableEavWebsite()
	codegen.LogFatal(err)
	if tew != nil { // only for those who have a wbesite specific table
		for _, w := range tws {
			aCollection, err := codegen.LoadStringEntities(ctx.dbc.DB, getAttrSelect(ctx, w.WebsiteID))
			codegen.LogFatal(err)
			for _, row := range aCollection {
				if aid, ok := row["attribute_id"]; ok {
					if nil == aws[aid] {
						aws[aid] = make([]codegen.StringEntities, 0, 200) // up to 200 websites at once
					}
					aws[aid] = append(aws[aid], row)
				} else {
					codegen.LogFatal(errgo.Newf("Column attribute_id not found in collection %#v\n", aCollection))
				}
			}
		}
	}
	return aws
}

func attrCollection(ctx *context, data map[string]interface{}) ([]byte, error) {

	aws := getAttributeValuesForWebsites(ctx)

	fmt.Printf("\n%s : %#v\n\n", ctx.et.EntityTypeCode, aws)

	/*
		1. codegen: tplAttrWebsiteEavAttribute
			Need: values from eav_attribute and check from website table of an entity
		2. codegen: tplAttrWebsiteEntityAttribute and use the code from 1 to embed
		3. codegen: tplAttrCollection
	*/

	funcMap := template.FuncMap{
		// isEavAttr checks if the attribute/column name belongs to table eav_attribute
		"isEavAttr": func(a string) bool { return codegen.EAVAttributeCoreColumns.Include(a) },
		// isEavEntityAttr checks if the attribute/column belongs to (customer|catalog|etc)_eav_attribute
		"isEavEntityAttr": func(a string) bool {
			if et, ok := codegen.ConfigEntityType[ctx.et.EntityTypeCode]; ok {
				return false == codegen.EAVAttributeCoreColumns.Include(a) && et.AttributeCoreColumns.Include(a)
			}
			return false
		},
		"isUnknownAttr": func(a string) bool {
			if et, ok := codegen.ConfigEntityType[ctx.et.EntityTypeCode]; ok {
				return false == codegen.EAVAttributeCoreColumns.Include(a) && false == et.AttributeCoreColumns.Include(a)
			}
			return false
		},
		"setAttrIdx": func(value, constName string) string {
			return strings.Replace(value, "{{.AttributeIndex}}", constName, -1)
		},
		"printWebsiteEavAttribute": func(attrID string) string {
			if _, ok := aws[attrID]; ok {
				//				fmt.Printf("\n%#v\n\n", cols)
				return "/* found " + attrID + " */ nil"
			}
			return "nil"
		},
	}

	return codegen.GenerateCode("", tplAttrCollection, data, funcMap)
}

func getAttrSelect(ctx *context, websiteID int64) *dbr.SelectBuilder {

	dbrSelect, err := eav.GetAttributeSelectSql(
		ctx.dbc.NewSession(nil),
		ctx.aat,
		ctx.et.EntityTypeID,
		websiteID,
	)
	codegen.LogFatal(err)
	dbrSelect.OrderDir(csdb.MainTable+".attribute_code", true)

	tew, err := ctx.aat.TableEavWebsite()
	codegen.LogFatal(err)
	if websiteID > 0 && tew != nil {
		// only here in codegen used to detect any changes if an attribute value will be overridden by a website ID
		dbrSelect.Where(csdb.ScopeTable + ".website_id IS NOT NULL")
		dbrSelect.Columns = append(dbrSelect.Columns, csdb.ScopeTable+".website_id")
	}

	return dbrSelect
}

func getAttrColumns(ctx *context, websiteID int64) codegen.Columns {
	columns, err := codegen.SQLQueryToColumns(ctx.dbc.DB, getAttrSelect(ctx, websiteID))
	codegen.LogFatal(err)
	codegen.LogFatal(columns.MapSQLToGoType(codegen.EavAttributeColumnNameToInterface))
	return columns
}

func getAttrPkg(et *eav.TableEntityType) string {
	if etConfig, ok := codegen.ConfigMaterializationAttributes[et.EntityTypeCode]; ok {
		return path.Base(etConfig.AttrPkgImp)
	}
	return ""
}

func getOutputFile(et *eav.TableEntityType) string {
	if etConfig, ok := codegen.ConfigMaterializationAttributes[et.EntityTypeCode]; ok {
		return etConfig.OutputFile
	}
	panic("You must specify an output file")
}

func getPackage(et *eav.TableEntityType) string {
	if etConfig, ok := codegen.ConfigMaterializationAttributes[et.EntityTypeCode]; ok {
		return etConfig.Package
	}
	panic("You must specify a package name")
}

// getStructName generates a nice struct name with a removed package name to avoid stutter but
// only removes the package name if the entity_type_code contains an underscore
// Depends on generated code from tableToStruct.
func getStructName(ctx *context, suffix ...string) string {
	structBaseName := ctx.et.EntityTypeCode
	if strings.Contains(ctx.et.EntityTypeCode, "_") {
		structBaseName = strings.Replace(ctx.et.EntityTypeCode, getPackage(ctx.et)+"_", "", -1)
	}
	return structBaseName + "_" + strings.Join(suffix, "_")
}

// stripCoreAttributeColumns returns a copy of columns and removes all core/default eav_attribute columns
func stripCoreAttributeColumns(cols codegen.Columns) codegen.Columns {
	ret := make(codegen.Columns, 0, len(cols))
	for _, col := range cols {
		if codegen.EAVAttributeCoreColumns.Include(col.Field.String) {
			continue
		}
		f := false
		for _, et := range codegen.ConfigEntityType {
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
//func generateAttributeCode(ctx *context) error {

//	//	name := getStructName(ctx, "attribute")
//	typeTplData := map[string]interface{}{
//		"AttrPkg":    getAttrPkg(ctx.et),
//		"AttrStruct": codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].AttrStruct,
//	}
//	if err != nil {
//		println(string(structCode))
//		return err
//	}

//	attributeCollection, err := codegen.LoadStringEntities(ctx.db, dbrSelect)
//	if err != nil {
//		return err
//	}

// @todo ValidateRules field must be converted from PHP serialized string to JSON
//	pkg := getPackage(ctx.et)

//	data := map[string]interface{}{

//		"Attributes":     attributeCollection,
//		"Name":     name,
//		"ImportPaths":    importPaths,
//		"PackageName":    pkg,
//		"AttrPkg": getAttrPkg(ctx.et),
//		"AttrPkgImp":     codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].AttrPkgImp,
//		"MyStruct": codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].MyStruct,
//		"AttrStruct":     codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].AttrStruct,
//		"FuncCollection": codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].FuncCollection,
//		"FuncGetter":     codegen.ConfigMaterializationAttributes[ctx.et.EntityTypeCode].FuncGetter,
//	}
//	funcMap := template.FuncMap{
//		// isEavAttr checks if the attribute/column name belongs to table eav_attribute
//		"isEavAttr": func(a string) bool { return codegen.EAVAttributeCoreColumns.Include(a) },
//		// isEavEntityAttr checks if the attribute/column belongs to (customer|catalog|etc)_eav_attribute
//		"isEavEntityAttr": func(a string) bool {
//			if et, ok := codegen.ConfigEntityType[ctx.et.EntityTypeCode]; ok {
//				return false == codegen.EAVAttributeCoreColumns.Include(a) && et.AttributeCoreColumns.Include(a)
//			}
//			return false
//		},
//		"isUnknownAttr": func(a string) bool {
//			if et, ok := codegen.ConfigEntityType[ctx.et.EntityTypeCode]; ok {
//				return false == codegen.EAVAttributeCoreColumns.Include(a) && false == et.AttributeCoreColumns.Include(a)
//			}
//			return false
//		},
//		"setAttrIdx": func(value, constName string) string {
//			return strings.Replace(value, "{{.AttributeIndex}}", constName, -1)
//		},
//	}

//	code, err := codegen.GenerateCode("", "tplTypeDefinitionFile", data, funcMap)
//	if err != nil {
//		println(string(code))
//		return err
//	}
//
//	return errgo.Mask(ioutil.WriteFile(getOutputFile(ctx.et), code, 0600))
//}
