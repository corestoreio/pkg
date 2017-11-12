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

// +build !mage1,!mage2

// Only include this file IF no specific build tag for mage has been set

package ccd

// Auto generated via tableToStruct

import (
	"github.com/corestoreio/pkg/storage/csdb"
	"github.com/corestoreio/pkg/util/null"
)

// TableIndex... is the index to a table. These constants are guaranteed
// to stay the same for all Magento versions. Please access a table via this
// constant instead of the raw table name. TableIndex iotas must start with 0.
const (
	TableIndexCoreConfigData = iota + 1 // Table: core_config_data
	TableIndexZZZ                       // the maximum index, which is not available.
)

func init() {
	TableCollection = csdb.MustNewTables(
		csdb.WithTable(
			TableIndexCoreConfigData,
			"core_config_data",
			&csdb.Column{Field: (`config_id`), ColumnType: (`int(10) unsigned`), Null: (`NO`), Key: (`PRI`), Extra: (`auto_increment`)},
			&csdb.Column{Field: (`scope`), ColumnType: (`varchar(8)`), Null: (`NO`), Key: (`MUL`), Default: null.StringFrom(`default`), Extra: (``)},
			&csdb.Column{Field: (`scope_id`), ColumnType: (`int(11)`), Null: (`NO`), Key: (``), Default: null.StringFrom(`0`), Extra: (``)},
			&csdb.Column{Field: (`path`), ColumnType: (`varchar(255)`), Null: (`NO`), Key: (``), Default: null.StringFrom(`general`), Extra: (``)},
			&csdb.Column{Field: (`value`), ColumnType: (`text`), Null: (`YES`), Key: (``), Extra: (``)},
		),
	)
}

// TableCoreConfigDataSlice represents a collection type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigDataSlice []*TableCoreConfigData

// TableCoreConfigData represents a type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigData struct {
	ConfigID int64       `db:"config_id" json:",omitempty"` // config_id int(10) unsigned NOT NULL PRI  auto_increment
	Scope    string      `db:"scope" json:",omitempty"`     // scope varchar(8) NOT NULL MUL DEFAULT 'default'
	ScopeID  int64       `db:"scope_id" json:",omitempty"`  // scope_id int(11) NOT NULL  DEFAULT '0'
	Path     string      `db:"path" json:",omitempty"`      // path varchar(255) NOT NULL  DEFAULT 'general'
	Value    null.String `db:"value" json:",omitempty"`     // value text NULL
}
