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
	"strings"

	"database/sql"
	"fmt"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/util"
)

type OneTable struct {
	Package          string
	Tick             string
	Name             string
	TableName        string // original table name
	Struct           string
	Slice            string
	Table            string
	GoColumns        codegen.Columns
	Columns          csdb.Columns
	MethodRecvPrefix string
	FindByPk         string
}

func NewOneTable(db *sql.DB, mageVersion int, pkgName, table string) OneTable {
	ot := OneTable{}
	ot.initTableNames(mageVersion, pkgName, table)
	ot.initColumns(db, table)
	return ot
}

// initTable takes care of the correct table name as some tables
// have different names in Magento2 than in Magento1.
//
// Generates a consistent name in field Name for all TableIndex*, Table*Slice
// and Table* types. This names is guaranteed to be the same whether we run
// Magento 1 or 2.
func (ot *OneTable) initTableNames(mageVersion int, pkgName, table string) {
	ot.Package = pkgName
	ot.Tick = "`"
	ot.Table = table
	ot.TableName = table

	if mappedName, ok := codegen.TableMapMagento1To2[strings.Replace(table, codegen.TablePrefix, "", 1)]; ok && mageVersion == util.MagentoV2 {
		ot.TableName = mappedName
	}

	// generate consistent name
	ot.Name = table
	// 1. retrieve the mapped name used in Magento2
	if mappedName, ok := codegen.TableMapMagento1To2[strings.Replace(table, codegen.TablePrefix, "", 1)]; ok {
		ot.Name = mappedName
	}
	// 2. Remove the package name from the table name
	ot.Name = codegen.PrepareVar(pkgName, ot.Name)

	ot.Struct = fmt.Sprintf("%s%s", TypePrefix, ot.Name)
	ot.Slice = fmt.Sprintf("%s%sSlice", TypePrefix, ot.Name)
}

func (ot *OneTable) initColumns(db *sql.DB, table string) {
	columns, err := codegen.GetColumns(db, table)
	codegen.LogFatal(err)
	codegen.LogFatal(columns.MapSQLToGoDBRType())

	ot.GoColumns = columns
	ot.Columns = columns.CopyToCSDB()

	if ot.Columns.PrimaryKeys().Len() > 0 {
		ot.FindByPk = "FindBy" + util.UnderscoreCamelize(ot.Columns.PrimaryKeys().JoinFields("_"))
	}

}
