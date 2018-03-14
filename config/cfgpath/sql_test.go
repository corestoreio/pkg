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

package cfgpath_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/strs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tableCollection *ddl.Tables

const (
	TableNameCoreConfigData = `core_config_data`
)

func init() {
	tableCollection = ddl.MustNewTables(
		ddl.WithTable(
			TableNameCoreConfigData,
			&ddl.Column{Field: `config_id`, ColumnType: `int(10) unsigned`, Null: `NO`, Key: `PRI`, Extra: `auto_increment`},
			&ddl.Column{Field: `scope`, ColumnType: `varchar(8)`, Null: `NO`, Key: `MUL`, Default: dml.MakeNullString(`default`), Extra: ``},
			&ddl.Column{Field: `scope_id`, ColumnType: `int(11)`, Null: `NO`, Key: ``, Default: dml.MakeNullString(`0`), Extra: ``},
			&ddl.Column{Field: `path`, ColumnType: `varchar(255)`, Null: `NO`, Key: ``, Default: dml.MakeNullString(`general`), Extra: ``},
			&ddl.Column{Field: `value`, ColumnType: `text`, Null: `YES`, Key: ``, Extra: ``},
		),
	)
}

// TableCoreConfigData represents a type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigData struct {
	ConfigID uint64         `json:",omitempty"` // config_id int(10) unsigned NOT NULL PRI  auto_increment
	Scope    string         `json:",omitempty"` // scope varchar(8) NOT NULL MUL DEFAULT 'default'
	ScopeID  int64          `json:",omitempty"` // scope_id int(11) NOT NULL  DEFAULT '0'
	Path     cfgpath.Route  `json:",omitempty"` // path varchar(255) NOT NULL  DEFAULT 'general'
	Value    dml.NullString `json:",omitempty"` // value text NULL
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

// TestIntegrationSQLType is not a real test for the type Route
func TestIntegrationSQLType(t *testing.T) {

	dbCon, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbCon, dbMock)

	var testPath = `system/full_page_cache/varnish/` + strs.RandAlnum(5)
	//var insPath = cfgpath.MakeRoute(testPath)
	var insVal = time.Now().Unix()

	dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `config_id`, `scope`, `scope_id`, `path`, `value` FROM `core_config_data` AS `main_table`")).
		WithArgs(testPath).
		WillReturnRows(
			sqlmock.NewRows([]string{"config_id", "scope", "scope_id", "path", "value"}).
				AddRow(1, "default", 0, testPath, fmt.Sprintf("%d", insVal)),
		)

	var ccd TableCoreConfigData
	tbl := tableCollection.MustTable(TableNameCoreConfigData)
	rc, err := tbl.SelectAll().WithDB(dbCon.DB).WithArgs().String(testPath).Load(context.TODO(), &ccd)
	require.NoError(t, err)
	assert.Exactly(t, uint64(1), rc)

	assert.Exactly(t, testPath, ccd.Path.String())
	haveI64, err := strconv.ParseInt(ccd.Value.String, 10, 64)
	assert.NoError(t, err)
	assert.Exactly(t, insVal, haveI64)
}
