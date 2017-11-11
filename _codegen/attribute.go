// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package codegen

import (
	"database/sql"

	"github.com/corestoreio/cspkg/storage/csdb"
	"github.com/corestoreio/errors"
)

// Depends on generated code from tableToStruct.
type AddAttrTables struct {
	EntityTypeCode string
	db             *sql.DB
}

// Implements interface eav.EntityTypeAdditionalAttributeTabler
func NewAddAttrTables(db *sql.DB, code string) *AddAttrTables {
	return &AddAttrTables{
		EntityTypeCode: code,
		db:             db,
	}
}

// Implements interface eav.EntityTypeAdditionalAttributeTabler
func (aa *AddAttrTables) TableAdditionalAttribute() (*csdb.Table, error) {
	if t, ok := ConfigEntityType[aa.EntityTypeCode]; ok {
		if t.TempAdditionalAttributeTable != "" {
			return aa.newTableStructure(t.TempAdditionalAttributeTable)
		}
		return nil, nil
	}
	return nil, errors.NewNotFoundf("Table %q", aa.EntityTypeCode)
}

// Implements interface eav.EntityTypeAdditionalAttributeTabler
func (aa *AddAttrTables) TableEavWebsite() (*csdb.Table, error) {
	if t, ok := ConfigEntityType[aa.EntityTypeCode]; ok {
		if t.TempAdditionalAttributeTableWebsite != "" {
			return aa.newTableStructure(t.TempAdditionalAttributeTableWebsite)
		}
		return nil, nil
	}
	return nil, errors.NewNotFoundf("Table %q", aa.EntityTypeCode)
}

func (aa *AddAttrTables) newTableStructure(tableName string) (*csdb.Table, error) {
	tableName = ReplaceTablePrefix(tableName)
	cols, err := GetColumns(aa.db, tableName)
	if err != nil {
		return nil, errors.Wrapf(err, "Table %q", tableName)
	}
	return csdb.NewTable(tableName, cols.CopyToCSDB()...), nil
}
