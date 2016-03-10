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

// Only include this file IF no specific build tag for mage has been set

package cfgpath_test

import (
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
)

const (
	tableIndexCoreConfigData csdb.Index = iota // Table: core_config_data
	tableIndexZZZ                              // the maximum index, which is not available.
)

var tableCollection csdb.TableManager

func init() {
	tableCollection = csdb.MustNewTableService(
		csdb.WithTable(
			tableIndexCoreConfigData,
			"core_config_data",
			csdb.Column{Field: dbr.NewNullString(`config_id`), Type: dbr.NewNullString(`int(10) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`PRI`), Default: dbr.NewNullString(nil), Extra: dbr.NewNullString(`auto_increment`)},
			csdb.Column{Field: dbr.NewNullString(`scope`), Type: dbr.NewNullString(`varchar(8)`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`MUL`), Default: dbr.NewNullString(`default`), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`scope_id`), Type: dbr.NewNullString(`int(11)`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(``), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`cfgpath`), Type: dbr.NewNullString(`varchar(255)`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(``), Default: dbr.NewNullString(`general`), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`value`), Type: dbr.NewNullString(`text`), Null: dbr.NewNullString(`YES`), Key: dbr.NewNullString(``), Default: dbr.NewNullString(nil), Extra: dbr.NewNullString(``)},
		),
	)
}

// TableCoreConfigDataSlice represents a collection type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigDataSlice []*TableCoreConfigData

// TableCoreConfigData represents a type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigData struct {
	ConfigID int64          `db:"config_id" json:",omitempty"` // config_id int(10) unsigned NOT NULL PRI  auto_increment
	Scope    string         `db:"scope" json:",omitempty"`     // scope varchar(8) NOT NULL MUL DEFAULT 'default'
	ScopeID  int64          `db:"scope_id" json:",omitempty"`  // scope_id int(11) NOT NULL  DEFAULT '0'
	Path     cfgpath.Route  `db:"path" json:",omitempty"`      // path varchar(255) NOT NULL  DEFAULT 'general'
	Value    dbr.NullString `db:"value" json:",omitempty"`     // value text NULL
}
