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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/materialized"
	"github.com/corestoreio/csfw/tools"
	"github.com/juju/errgo"
)

func materializeAttributes(ctx *context) {
	defer ctx.wg.Done()
	for _, et := range materialized.GetEntityTypeCollection() {
		ctx.et = et
		tools.LogFatal(generateAttributeCode(ctx))

		// EAV -> Create queries for AttributeSets and AttributeGroups
	}
}

// getName
func getName(ctx *context, suffix ...string) string {
	pkg := path.Base(ctx.et.ImportPath)
	structBaseName := ctx.et.EntityTypeCode
	if strings.Contains(ctx.et.EntityTypeCode, "_") {
		structBaseName = strings.Replace(ctx.et.EntityTypeCode, pkg+"_", "", -1)
	}
	return structBaseName + "_" + strings.Join(suffix, "_")
}

func generateAttributeCode(ctx *context) error {
	dbrSelect, err := eav.GetAttributeSelectSql(ctx.dbrConn.NewSession(nil), ctx.et, 0)
	if err != nil {
		return err
	}

	columns, err := tools.SQLQueryToColumns(ctx.db, dbrSelect)
	if err != nil {
		return err
	}

	tools.LogFatal(columns.MapSQLToGoType(tools.EavAttributeColumnNameToInterface))

	name := "CS_" + getName(ctx, "attribute")
	structCode, err := tools.ColumnsToStructCode(name, columns, tplTypeDefinitions)
	if err != nil {
		return err
	}

	attributeCollection, err := tools.GetSQL(ctx.db, dbrSelect)
	if err != nil {
		return err
	}

	// @todo ValidateRules field must be converted from PHP serialized string to JSON
	importPaths := tools.PrepareForTemplate(columns, attributeCollection, ctx.modelMap)

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
		PackageName:    path.Base(ctx.et.ImportPath),
	}

	code, err := tools.GenerateCode("", tplFileBody, data)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%s%s%sgenerated_%s.go",
		ctx.goSrcPath,
		ctx.et.ImportPath,
		string(os.PathSeparator),
		getName(ctx, "attribute"),
	)
	return errgo.Mask(ioutil.WriteFile(path, code, 0600))
}
