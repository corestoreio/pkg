// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// +build !mage1,!mage2

// Only include this file IF no specific build tag for mage has been set

package ccd

// Auto generated via tableToStruct

import (
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
)

const (
	TableNameCoreConfigData = `core_config_data`
)

func NewTableCollection() *ddl.Tables {
	ddl.MustNewTables(
		ddl.WithTable(
			TableNameCoreConfigData,
			&ddl.Column{Field: `config_id`, ColumnType: `int(10) unsigned`, Null: `NO`, Key: `PRI`, Extra: `auto_increment`},
			&ddl.Column{Field: `scope`, ColumnType: `varchar(8)`, Null: `NO`, Key: `MUL`, Default: dml.MakeNullString(`default`), Extra: ""},
			&ddl.Column{Field: `scope_id`, ColumnType: `int(11)`, Null: `NO`, Key: "", Default: dml.MakeNullString(`0`), Extra: ""},
			&ddl.Column{Field: `path`, ColumnType: `varchar(255)`, Null: `NO`, Key: "", Default: dml.MakeNullString(`general`), Extra: ""},
			&ddl.Column{Field: `value`, ColumnType: `text`, Null: `YES`, Key: ``, Extra: ""},
		),
	)
}

// TableCoreConfigDataSlice represents a collection type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigDataSlice []*TableCoreConfigData

// TableCoreConfigData represents a type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigData struct {
	ConfigID int64        `json:",omitempty"`            // config_id int(10) unsigned NOT NULL PRI  auto_increment
	Path     config.Path  `db:"path" json:",omitempty"`  // path varchar(255) NOT NULL  DEFAULT 'general'
	Value    config.Value `db:"value" json:",omitempty"` // value text NULL
}

// MapColumns implements interface ColumnMapper only partially.
func (p *TableCoreConfigData) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		// bug check for null in Path
		var pth dml.NullString
		pth.String = p.Path.Data
		pth.Valid = p.Path.Valid
		return cm.Uint64(&p.ConfigID).String(&p.Scope).Int64(&p.ScopeID).NullString(&pth).NullString(&p.Value).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "config_id": // customer_id is an alias
			cm.Uint64(&p.ConfigID)
		case "scope":
			cm.String(&p.Scope)
		case "scope_id":
			cm.Int64(&p.ScopeID)
		case "path":
			var pth dml.NullString
			cm.NullString(&pth)
			p.Path.Data = pth.String
			p.Path.Valid = pth.Valid
		case "value":
			cm.NullString(&p.Value)

		default:
			return errors.NotFound.Newf("[dml_test] customerEntity Column %q not found", c)
		}
	}
	return cm.Err()
}
