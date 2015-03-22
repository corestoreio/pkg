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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/tools"
	"github.com/juju/errgo"
)

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

func getPackage(et *eav.EntityType) string {
	return path.Base(getImportPath(et))
}

// getName generates a nice struct name with a removed package name to avoid stutter but
// only removes the package name if the entity_type_code contains an underscore
func getName(ctx *context, suffix ...string) string {
	structBaseName := ctx.et.EntityTypeCode
	if strings.Contains(ctx.et.EntityTypeCode, "_") {
		structBaseName = strings.Replace(ctx.et.EntityTypeCode, getPackage(ctx.et)+"_", "", -1)
	}
	return structBaseName + "_" + strings.Join(suffix, "_")
}

func generateAttributeCode(ctx *context) error {
	// @todo get all website IDs

	dbrSelect, err := eav.GetAttributeSelectSql(ctx.dbrConn.NewSession(nil), newAddAttrTables(ctx), ctx.et.EntityTypeID, 0)
	if err != nil {
		return err
	}

	columns, err := tools.SQLQueryToColumns(ctx.db, dbrSelect)
	if err != nil {
		return err
	}

	tools.LogFatal(columns.MapSQLToGoType(tools.EavAttributeColumnNameToInterface))

	name := "CS_" + getName(ctx, "attribute")
	structCode, err := tools.ColumnsToStructCode(name, columns, tplTypeDefinition)
	if err != nil {
		return err
	}

	attributeCollection, err := tools.GetSQL(ctx.db, dbrSelect)
	if err != nil {
		return err
	}

	// @todo ValidateRules field must be converted from PHP serialized string to JSON
	importPaths := tools.PrepareForTemplate(columns, attributeCollection, tools.ConfigAttributeModel)

	data := struct {
		TypeDefinition string
		Attributes     []tools.StringEntities
		Name           string
		ImportPaths    []string
		PackageName    string
	}{
		TypeDefinition: string(structCode),
		Attributes:     attributeCollection,
		Name:           name,
		ImportPaths:    importPaths,
		PackageName:    getPackage(ctx.et),
	}

	code, err := tools.GenerateCode("", tplTypeDefinitionFile, data)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%s%s%sgenerated_%s.go",
		ctx.goSrcPath,
		getImportPath(ctx.et),
		string(os.PathSeparator),
		getName(ctx, "attribute"),
	)
	return errgo.Mask(ioutil.WriteFile(path, code, 0600))
}

type addAttrTables struct {
	*eav.EntityType
	db *sql.DB
}

// Implements interface eav.EntityTypeAdditionalAttributeTabler
func newAddAttrTables(ctx *context) *addAttrTables {
	return &addAttrTables{
		EntityType: ctx.et,
		db:         ctx.db,
	}
}

// Implements interface eav.EntityTypeAdditionalAttributeTabler
func (aa *addAttrTables) TableAdditionalAttribute() (*csdb.TableStructure, error) {
	if t, ok := tools.ConfigEntityType[aa.EntityTypeCode]; ok {
		if t.TempAdditionalAttributeTable != "" {
			return aa.newTableStructure(t.TempAdditionalAttributeTable)
		}
		return nil, nil
	}
	return nil, errgo.Newf("Table for %s not found", aa.EntityTypeCode)
}

// Implements interface eav.EntityTypeAdditionalAttributeTabler
func (aa *addAttrTables) TableEavWebsite() (*csdb.TableStructure, error) {
	if t, ok := tools.ConfigEntityType[aa.EntityTypeCode]; ok {
		if t.TempAdditionalAttributeTableWebsite != "" {
			return aa.newTableStructure(t.TempAdditionalAttributeTableWebsite)
		}
		return nil, nil
	}
	return nil, errgo.Newf("Table for %s not found", aa.EntityTypeCode)
}

func (aa *addAttrTables) newTableStructure(tableName string) (*csdb.TableStructure, error) {
	cols, err := tools.GetColumns(aa.db, tableName)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return csdb.NewTableStructure(tableName, cols.GetFieldNames(true), cols.GetFieldNames(false)), nil
}
