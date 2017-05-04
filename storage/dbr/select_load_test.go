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

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestSelect_Rows(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		sel := &dbr.Select{}
		sel.Columns = []string{"a", "b"}
		rows, err := sel.Rows(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Query Error", func(t *testing.T) {
		sel := &dbr.Select{
			Table:   dbr.MakeAlias("tableX"),
			Columns: []string{"a", "b"},
		}
		sel.DB.Querier = dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		}

		rows, err := sel.Rows(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestSelect_Row(t *testing.T) {

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()
	dbMock.ExpectQuery("SELECT a, b FROM `tableX`").WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

	sel := &dbr.Select{
		Table:   dbr.MakeAlias("tableX"),
		Columns: []string{"a", "b"},
	}
	sel.DB.QueryRower = dbc.DB
	row := sel.Row(context.TODO())
	var x string
	err := row.Scan(&x)
	assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
}

func TestSelect_Prepare(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		sel := &dbr.Select{}
		sel.Columns = []string{"a", "b"}
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
		dbMock.ExpectPrepare("SELECT a, b FROM `tableX`").WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		sel := &dbr.Select{
			Table:   dbr.MakeAlias("tableX"),
			Columns: []string{"a", "b"},
		}
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
	Data []*TableCoreConfigData
	sel  *dbr.Select
}

func newTableCoreConfigDatas() *TableCoreConfigDatas {
	return &TableCoreConfigDatas{
		sel: dbr.NewSelect("*").From("core_config_data"),
	}
}

// TableCoreConfigData represents a type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigData struct {
	ConfigID int64          `db:"config_id" json:",omitempty"` // config_id int(10) unsigned NOT NULL PRI  auto_increment
	Scope    string         `db:"scope" json:",omitempty"`     // scope varchar(8) NOT NULL MUL DEFAULT 'default'
	ScopeID  int64          `db:"scope_id" json:",omitempty"`  // scope_id int(11) NOT NULL  DEFAULT '0'
	Path     string         `db:"path" json:",omitempty"`      // path varchar(255) NOT NULL  DEFAULT 'general'
	Value    dbr.NullString `db:"value" json:",omitempty"`     // value text NULL
}

// ScanArgs implement Loader interface
func (s *TableCoreConfigDatas) ScanArgs(columns []string) []interface{} {
	s.Data = make([]*TableCoreConfigData, 0, 10)
	var c TableCoreConfigData
	return []interface{}{&c.ConfigID, &c.Scope, &c.ScopeID, &c.Path, &c.Value}
}

// Row implement Loader interface
func (s *TableCoreConfigDatas) Row(idx int64, values []interface{}) error {
	c := &TableCoreConfigData{}
	if v, ok := values[0].(*int64); ok {
		c.ConfigID = *v
	} else {
		return errors.NewNotValidf("[pkg] Field ConfigID. Failed to assert to *int64. Got: %#v", values[0])
	}
	if v, ok := values[1].(*string); ok {
		c.Scope = *v
	} else {
		return errors.NewNotValidf("[pkg] Field Scope. Failed to assert to *string. Got: %#v", values[1])
	}
	if v, ok := values[2].(*int64); ok {
		c.ScopeID = *v
	} else {
		return errors.NewNotValidf("[pkg] Field ConfigID. Failed to assert to *int64. Got: %#v", values[2])
	}
	if v, ok := values[3].(*string); ok {
		c.Path = *v
	} else {
		return errors.NewNotValidf("[pkg] Field Path. Failed to assert to *string. Got: %#v", values[3])
	}
	if v, ok := values[4].(*dbr.NullString); ok {
		c.Value = *v
	} else {
		return errors.NewNotValidf("[pkg] Field Path. Failed to assert to *NullString. Got: %#v", values[4])
	}

	s.Data = append(s.Data, c)

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
