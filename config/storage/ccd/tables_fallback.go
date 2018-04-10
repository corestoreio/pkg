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

package ccd

// Auto generated via tableToStruct

import (
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
)

const (
	TableNameCoreConfigData = `core_config_data`
)

// NewTableCollection creates a new Tables object for TableNameCoreConfigData.
func NewTableCollection(db dml.QueryExecPreparer) *ddl.Tables {
	return ddl.MustNewTables(
		ddl.WithTable(
			TableNameCoreConfigData,
			&ddl.Column{Field: `config_id`, ColumnType: `int(10) unsigned`, Null: `NO`, Key: `PRI`, Extra: `auto_increment`},
			&ddl.Column{Field: `scope`, ColumnType: `varchar(8)`, Null: `NO`, Key: `MUL`, Default: dml.MakeNullString(`default`), Extra: ""},
			&ddl.Column{Field: `scope_id`, ColumnType: `int(11)`, Null: `NO`, Key: "", Default: dml.MakeNullString(`0`), Extra: ""},
			&ddl.Column{Field: `path`, ColumnType: `varchar(255)`, Null: `NO`, Key: "", Default: dml.MakeNullString(`general`), Extra: ""},
			&ddl.Column{Field: `value`, ColumnType: `text`, Null: `YES`, Key: ``, Extra: ""},
		),
		ddl.WithDB(db),
	)
}

// TableCoreConfigData represents a type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigData struct {
	ConfigID int64          // config_id int(10) unsigned NOT NULL PRI  auto_increment
	Scope    string         // scope varchar(8) NOT NULL MUL DEFAULT 'default'
	ScopeID  int64          // scope_id int(11) NOT NULL  DEFAULT '0'
	Path     string         // path varchar(255) NOT NULL  DEFAULT 'general'
	Value    dml.NullString // value text NULL
}

// MapColumns implements interface ColumnMapper only partially.
func (p *TableCoreConfigData) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Int64(&p.ConfigID).String(&p.Scope).Int64(&p.ScopeID).String(&p.Path).NullString(&p.Value).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "config_id": // customer_id is an alias
			cm.Int64(&p.ConfigID)
		case "scope":
			cm.String(&p.Scope)
		case "scope_id":
			cm.Int64(&p.ScopeID)
		case "path":
			cm.String(&p.Path)
		case "value":
			cm.NullString(&p.Value)
		default:
			return errors.NotFound.Newf("[ccd] TableCoreConfigData Column %q not found", c)
		}
	}
	return cm.Err()
}
