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
	"database/sql"
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
		rows, err := sel.Query(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Error", func(t *testing.T) {
		sel := &dbr.Select{
			BuilderBase: dbr.BuilderBase{
				Table: dbr.MakeIdentifier("tableX"),
			},
		}
		sel.AddColumns("a", "b")
		sel.DB = dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		}

		rows, err := sel.Query(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

	t.Run("Success", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		smr := sqlmock.NewRows([]string{"a"}).AddRow("row1").AddRow("row2")
		dbMock.ExpectQuery("SELECT `a` FROM `tableX`").WillReturnRows(smr)

		sel := dbr.NewSelect("a").From("tableX")
		sel.DB = dbc.DB
		rows, err := sel.Query(context.TODO())
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
		sel := dbr.NewSelect()
		stmt, err := sel.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectPrepare("SELECT `a`, `b` FROM `tableX`").WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		sel := dbr.NewSelect("a", "b").From("tableX").WithDB(dbc.DB)
		stmt, err := sel.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

}

var _ dbr.Scanner = (*TableCoreConfigDataSlice)(nil)
var _ dbr.RowCloser = (*TableCoreConfigDataSlice)(nil)

// TableCoreConfigDataSlice represents a collection type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigDataSlice struct {
	Convert dbr.RowConvert
	Data    []*TableCoreConfigData
	scanErr error // just for testing not needed otherwise
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

func (ps *TableCoreConfigDataSlice) RowScan(r *sql.Rows) error {
	if err := ps.Convert.Scan(r); err != nil {
		return errors.WithStack(err)
	}
	var o TableCoreConfigData
	for i, col := range ps.Convert.Columns {
		b := ps.Convert.Index(i)
		var err error
		switch col {
		case "config_id":
			o.ConfigID, err = b.Int64()
		case "scope":
			o.Scope, err = b.Str()
		case "scope_id":
			o.ScopeID, err = b.Int64()
		case "path":
			o.Path, err = b.Str()
		case "value":
			o.Value.NullString, err = b.NullString()
		}
		if err != nil {
			return errors.Wrapf(err, "[dbr] Failed to convert value at row % with column index %d", ps.Convert.Count, i)
		}
	}
	ps.Data = append(ps.Data, &o)
	return nil
}

func (ps *TableCoreConfigDataSlice) RowClose() error {
	return ps.scanErr
}

func TestSelect_Load(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery("SELECT").WillReturnRows(cstesting.MustMockRows(cstesting.WithFile("testdata/core_config_data.csv")))
		s := dbr.NewSelect("*").From("core_config_data")
		s.DB = dbc.DB

		ccd := &TableCoreConfigDataSlice{}

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
	})

	t.Run("row error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		r := sqlmock.NewRows([]string{"config_id"}).FromCSVString("222\n333\n").
			RowError(1, errors.NewConnectionFailedf("Con failed"))
		dbMock.ExpectQuery("SELECT").WillReturnRows(r)
		s := dbr.NewSelect("config_id").From("core_config_data")
		s.DB = dbc.DB

		ccd := &TableCoreConfigDataSlice{}
		_, err := s.Load(context.TODO(), ccd)
		assert.True(t, errors.IsConnectionFailed(err), "%+v", err)
	})

	t.Run("RowClose error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		r := sqlmock.NewRows([]string{"config_id"}).FromCSVString("222\n333\n").AddRow("3456")
		dbMock.ExpectQuery("SELECT").WillReturnRows(r)
		s := dbr.NewSelect("config_id").From("core_config_data")
		s.DB = dbc.DB

		ccd := &TableCoreConfigDataSlice{
			scanErr: errors.NewDuplicatedf("Somewhere exists a duplicate entry"),
		}
		_, err := s.Load(context.TODO(), ccd)
		assert.True(t, errors.IsDuplicated(err), "%+v", err)
	})
}
