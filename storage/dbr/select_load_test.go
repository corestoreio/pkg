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

package dbr_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelect_Rows(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		sel := &dbr.Select{}
		sel.AddColumns("a", "b")
		rows, err := sel.Rows(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Query Error", func(t *testing.T) {
		sel := &dbr.Select{
			Table: dbr.MakeAlias("tableX"),
		}
		sel.AddColumns("a", "b")
		sel.DB.Querier = dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		}

		rows, err := sel.Rows(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

	t.Run("success", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()
		smr := sqlmock.NewRows([]string{"a"}).AddRow("row1").AddRow("row2")
		dbMock.ExpectQuery("SELECT `a` FROM `tableX`").WillReturnRows(smr)

		sel := dbr.NewSelect("a").From("tableX")
		sel.DB.Querier = dbc.DB
		rows, err := sel.Rows(context.TODO())
		require.NoError(t, err, "%+v", err)
		defer func() {
			if err := rows.Close(); err != nil {
				t.Fatal(err)
			}
		}()

		var xx []string
		for rows.Next() {
			var x string
			require.NoError(t, rows.Scan(&x))
			xx = append(xx, x)
		}
		assert.Exactly(t, []string{"row1", "row2"}, xx)
	})
}

func TestSelect_Prepare(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		sel := dbr.NewSelect("a", "b")
		stmt, err := sel.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()
		dbMock.ExpectPrepare("SELECT `a`, `b` FROM `tableX`").WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		sel := dbr.NewSelect("a", "b").From("tableX")
		sel.DB.Preparer = dbc.DB
		stmt, err := sel.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

}

// TableCoreConfigDataSlice used in benchmarks
type TableCoreConfigDataSlice []*TableCoreConfigData

// TableCoreConfigDatas represents a collection type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigDatas struct {
	DataCap  int
	Data     []*TableCoreConfigData
	scanArgs []interface{}
	dto      TableCoreConfigData
}

// TableCoreConfigData represents a type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigData struct {
	ConfigID int64          `json:",omitempty"` // config_id int(10) unsigned NOT NULL PRI  auto_increment
	Scope    string         `json:",omitempty"` // scope varchar(8) NOT NULL MUL DEFAULT 'default'
	ScopeID  int64          `json:",omitempty"` // scope_id int(11) NOT NULL  DEFAULT '0'
	Path     string         `json:",omitempty"` // path varchar(255) NOT NULL  DEFAULT 'general'
	Value    dbr.NullString `json:",omitempty"` // value text NULL
}

func (ps *TableCoreConfigDatas) ScanRow(idx int, columns []string, scan func(dest ...interface{}) error) error {
	const fieldCount = 5 //  5 == number of struct fields
	if idx == 0 && nil == ps.Data {
		cap := ps.DataCap
		if cap == 0 {
			cap = 10
		}
		ps.Data = make([]*TableCoreConfigData, 0, cap)
		ps.scanArgs = make([]interface{}, 0, fieldCount)
		for _, c := range columns {
			switch c {
			case "config_id":
				ps.scanArgs = append(ps.scanArgs, &ps.dto.ConfigID)
			case "scope":
				ps.scanArgs = append(ps.scanArgs, &ps.dto.Scope)
			case "scope_id":
				ps.scanArgs = append(ps.scanArgs, &ps.dto.ScopeID)
			case "path":
				ps.scanArgs = append(ps.scanArgs, &ps.dto.Path)
			case "value":
				ps.scanArgs = append(ps.scanArgs, &ps.dto.Value)
			default:
				return errors.NewNotFoundf("[dbr_test] Column %q not found", c)
			}
		}
	}

	if err := scan(ps.scanArgs...); err != nil {
		return errors.Wrap(err, "[dbr_test] dbrPersons.ScanRow")
	}

	// We can copy here easy by assigning the value to a new variable and then
	// creating a pointer. If the the TableCoreConfigData struct would contain
	// pointers to other structs, then the copying would be more complex to
	// avoid a shallow copy.
	ccd := ps.dto
	ps.Data = append(ps.Data, &ccd)
	return nil
}

func TestSelect_Load(t *testing.T) {

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	dbMock.ExpectQuery("SELECT").WillReturnRows(cstesting.MustMockRows(cstesting.WithFile("testdata/core_config_data.csv")))
	s := dbr.NewSelect("*").From("core_config_data")
	s.DB.Querier = dbc.DB

	ccd := &TableCoreConfigDatas{}

	_, err := s.Load(context.TODO(), ccd)
	assert.NoError(t, err, "%+v", err)

	buf := new(bytes.Buffer)
	je := json.NewEncoder(buf)

	for _, c := range ccd.Data {
		if err := je.Encode(c); err != nil {
			t.Fatalf("%+v", err)
		}
	}
	assert.Equal(t, "{\"ConfigID\":2,\"Scope\":\"default\",\"Path\":\"web/unsecure/base_url\",\"Value\":\"http://mgeto2.local/\"}\n{\"ConfigID\":3,\"Scope\":\"website\",\"ScopeID\":11,\"Path\":\"general/locale/code\",\"Value\":\"en_US\"}\n{\"ConfigID\":4,\"Scope\":\"default\",\"Path\":\"general/locale/timezone\",\"Value\":\"Europe/Berlin\"}\n{\"ConfigID\":5,\"Scope\":\"default\",\"Path\":\"currency/options/base\",\"Value\":\"EUR\"}\n{\"ConfigID\":15,\"Scope\":\"store\",\"ScopeID\":33,\"Path\":\"design/head/includes\",\"Value\":\"\\u003clink  rel=\\\"stylesheet\\\" type=\\\"text/css\\\" href=\\\"{{MEDIA_URL}}styles.css\\\" /\\u003e\"}\n{\"ConfigID\":16,\"Scope\":\"default\",\"Path\":\"admin/security/use_case_sensitive_login\",\"Value\":null}\n{\"ConfigID\":17,\"Scope\":\"default\",\"Path\":\"admin/security/session_lifetime\",\"Value\":\"90000\"}\n",
		buf.String())
}
