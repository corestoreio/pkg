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
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log/logw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync/atomic"
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

var _ dbr.Scanner = (*TableCoreConfigDataSlice)(nil)

var _ dbr.RowCloser = (*TableCoreConfigDataSlice)(nil)

// TableCoreConfigDataSlice represents a collection type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigDataSlice struct {
	Convert dbr.RowConvert
	Data    []*TableCoreConfigData
	err     error // just for testing not needed otherwise
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
	return ps.err
}

var _ dbr.ArgumentsAppender = (*TableCoreConfigDataSlice)(nil)

func (ps TableCoreConfigDataSlice) AppendArgs(args dbr.Arguments, _ []string) (dbr.Arguments, error) {
	return args, ps.err
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
			err: errors.NewDuplicatedf("Somewhere exists a duplicate entry"),
		}
		_, err := s.Load(context.TODO(), ccd)
		assert.True(t, errors.IsDuplicated(err), "%+v", err)
	})
}

func TestSelect_Prepare(t *testing.T) {
	t.Parallel()

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

	t.Run("Query", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("SELECT `name`, `email` FROM `dbr_person` WHERE (`id` = ?)"))
		prep.ExpectQuery().WithArgs(6789).
			WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher", "peter@gopher.go"))

		prep.ExpectQuery().WithArgs(6790).
			WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher2", "peter@gopher.go2"))

		stmt, err := dbr.NewSelect("name", "email").From("dbr_person").
			Where(dbr.Column("id").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		t.Run("Context", func(t *testing.T) {

			rows, err := stmt.Query(context.TODO(), 6789)
			require.NoError(t, err)
			defer rows.Close()

			cols, err := rows.Columns()
			require.NoError(t, err)
			assert.Exactly(t, []string{"name", "email"}, cols)
		})

		t.Run("RowContext", func(t *testing.T) {

			row := stmt.QueryRow(context.TODO(), 6790)
			require.NoError(t, err)
			n, e := "", ""
			require.NoError(t, row.Scan(&n, &e))

			assert.Exactly(t, "Peter Gopher2", n)
			assert.Exactly(t, "peter@gopher.go2", e)
		})
	})

	t.Run("Exec", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("SELECT `name`, `email` FROM `dbr_person` WHERE (`id` = ?)"))

		stmt, err := dbr.NewSelect("name", "email").From("dbr_person").
			Where(dbr.Column("id").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		const iterations = 3

		t.Run("WithArguments", func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				prep.ExpectQuery().WithArgs(6899).
					WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher", "peter@gopher.go"))
			}
			// use loop with Query+ and add args before
			stmt.WithArguments(dbr.MakeArgs(1).Int(6899))

			for i := 0; i < iterations; i++ {
				rows, err := stmt.Query(context.TODO())
				require.NoError(t, err)

				cols, err := rows.Columns()
				require.NoError(t, err)
				assert.Exactly(t, []string{"name", "email"}, cols)
				rows.Close()
			}
		})

		t.Run("WithRecords", func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				prep.ExpectQuery().WithArgs(6900).
					WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher2", "peter@gopher.go2"))
			}

			p := &dbrPerson{ID: 6900}
			stmt.WithRecords(dbr.Qualify("", p))

			for i := 0; i < iterations; i++ {
				rows, err := stmt.Query(context.TODO())
				require.NoError(t, err)

				cols, err := rows.Columns()
				require.NoError(t, err)
				assert.Exactly(t, []string{"name", "email"}, cols)
				rows.Close()
			}
		})

		t.Run("WithRecords Error", func(t *testing.T) {
			p := TableCoreConfigDataSlice{err: errors.NewDuplicatedf("Found a duplicate")}
			stmt.WithRecords(dbr.Qualify("", p))
			rows, err := stmt.Query(context.TODO())
			assert.True(t, errors.IsDuplicated(err), "%+v", err)
			assert.Nil(t, rows)
		})
	})

	t.Run("Load", func(t *testing.T) {

		t.Run("multi rows", func(t *testing.T) {
			dbc, dbMock := cstesting.MockDB(t)
			defer cstesting.MockClose(t, dbc, dbMock)

			prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("SELECT `config_id`, `scope_id`, `path` FROM `core_config_data` WHERE (`config_id` IN (?))"))

			stmt, err := dbr.NewSelect("config_id", "scope_id", "path").From("core_config_data").
				Where(dbr.Column("config_id").In().PlaceHolder()).
				WithDB(dbc.DB).
				Prepare(context.TODO())
			require.NoError(t, err, "failed creating a prepared statement")
			defer func() {
				require.NoError(t, stmt.Close(), "Close on a prepared statement")
			}()

			columns := []string{"config_id", "scope_id", "path"}

			prep.ExpectQuery().WithArgs(345).
				WillReturnRows(sqlmock.NewRows(columns).AddRow(3, 4, "a/b/c").AddRow(4, 4, "a/b/d"))

			ccd := &TableCoreConfigDataSlice{}

			rc, err := stmt.WithArgs(345).Load(context.TODO(), ccd)
			require.NoError(t, err)
			assert.Exactly(t, int64(2), rc)

			assert.Exactly(t, "&{3  4 a/b/c {{ false}}}", fmt.Sprintf("%v", ccd.Data[0]))
			assert.Exactly(t, "&{4  4 a/b/d {{ false}}}", fmt.Sprintf("%v", ccd.Data[1]))
		})

		t.Run("Int64", func(t *testing.T) {
			dbc, dbMock := cstesting.MockDB(t)
			defer cstesting.MockClose(t, dbc, dbMock)

			prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("SELECT `scope_id` FROM `core_config_data` WHERE (`config_id` = ?)"))

			stmt, err := dbr.NewSelect("scope_id").From("core_config_data").
				Where(dbr.Column("config_id").PlaceHolder()).
				WithDB(dbc.DB).
				Prepare(context.TODO())
			require.NoError(t, err, "failed creating a prepared statement")
			defer func() {
				require.NoError(t, stmt.Close(), "Close on a prepared statement")
			}()

			columns := []string{"scope_id"}

			prep.ExpectQuery().WithArgs(346).WillReturnRows(sqlmock.NewRows(columns).AddRow(35))

			val, err := stmt.WithArguments(dbr.MakeArgs(1).Int64(346)).LoadInt64(context.TODO())
			require.NoError(t, err)
			assert.Exactly(t, int64(35), val)
		})

		t.Run("Int64s", func(t *testing.T) {
			dbc, dbMock := cstesting.MockDB(t)
			defer cstesting.MockClose(t, dbc, dbMock)

			prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("SELECT `scope_id` FROM `core_config_data` WHERE (`config_id` IN (?))"))

			stmt, err := dbr.NewSelect("scope_id").From("core_config_data").
				Where(dbr.Column("config_id").In().PlaceHolder()).
				WithDB(dbc.DB).
				Prepare(context.TODO())
			require.NoError(t, err, "failed creating a prepared statement")
			defer func() {
				require.NoError(t, stmt.Close(), "Close on a prepared statement")
			}()

			columns := []string{"scope_id"}

			prep.ExpectQuery().WithArgs(346, 347).WillReturnRows(sqlmock.NewRows(columns).AddRow(36).AddRow(37))

			val, err := stmt.WithArguments(dbr.MakeArgs(1).Int64s(346, 347)).LoadInt64s(context.TODO())
			require.NoError(t, err)
			assert.Exactly(t, []int64{36, 37}, val)
		})

	})
}

func TestSelect_WithLogger(t *testing.T) {
	uniID := new(int32)
	rConn := createRealSession(t)
	defer cstesting.Close(t, rConn)

	var uniqueIDFunc = func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 1))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	require.NoError(t, rConn.Options(dbr.WithLogger(lg, uniqueIDFunc)))

	t.Run("ConnPool", func(t *testing.T) {
		u := rConn.SelectFrom("dbr_people").AddColumns("name", "email").Where(dbr.Column("id").In().Int64s(6, 8))

		t.Run("Query", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			rows, err := u.Interpolate().Query(context.TODO())
			require.NoError(t, err)
			require.NoError(t, rows.Close())

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (6,8))\"\n",
				buf.String())
		})

		t.Run("Load", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			p := &dbrPerson{}
			_, err := u.Interpolate().Load(context.TODO(), p)
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Load conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (6,8))\"\n",
				buf.String())
		})

		t.Run("LoadInt64", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			_, err := u.Interpolate().LoadInt64(context.TODO())
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}

			assert.Exactly(t, "DEBUG LoadInt64 conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (6,8))\"\n",
				buf.String())
		})

		t.Run("LoadInt64s", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			_, err := u.Interpolate().LoadInt64s(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG LoadInt64s conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 row_count: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (6,8))\"\n",
				buf.String())
		})

		t.Run("LoadUint64", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			_, err := u.Interpolate().LoadUint64(context.TODO())
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
			assert.Exactly(t, "DEBUG LoadUint64 conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (6,8))\"\n",
				buf.String())
		})

		t.Run("LoadUint64s", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			_, err := u.Interpolate().LoadUint64s(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG LoadUint64s conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 row_count: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (6,8))\"\n",
				buf.String())
		})

		t.Run("LoadFloat64", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			_, err := u.Interpolate().LoadFloat64(context.TODO())
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
			assert.Exactly(t, "DEBUG LoadFloat64 conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (6,8))\"\n",
				buf.String())
		})

		t.Run("LoadFloat64s", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			_, err := u.Interpolate().LoadFloat64s(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG LoadFloat64s conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (6,8))\"\n",
				buf.String())
		})

		t.Run("LoadString", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			_, err := u.Interpolate().LoadString(context.TODO())
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}

			assert.Exactly(t, "DEBUG LoadString conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (6,8))\"\n",
				buf.String())
		})

		t.Run("LoadStrings", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			_, err := u.Interpolate().LoadStrings(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG LoadStrings conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 row_count: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (6,8))\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := u.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ02*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (?,?))\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			tx, err := rConn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.NoError(t, tx.Wrap(func() error {
				rows, err := tx.SelectFrom("dbr_people").
					AddColumns("name", "email").Where(dbr.Column("id").In().Int64s(7, 9)).
					Interpolate().Query(context.TODO())
				require.NoError(t, err)
				return rows.Close()
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\"\nDEBUG Query conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\" select_id: \"UNIQ04\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ04*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (7,9))\"\nDEBUG Commit conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\" duration: 0\n",
				buf.String())
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		require.NoError(t, err)

		u := conn.SelectFrom("dbr_people", "dp2").AddColumns("name", "email").Where(dbr.Column("id").In().Int64s(61, 81))

		t.Run("Query", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()

			rows, err := u.Interpolate().Query(context.TODO())
			require.NoError(t, err)
			require.NoError(t, rows.Close())

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ06*/ `name`, `email` FROM `dbr_people` AS `dp2` WHERE (`id` IN (61,81))\"\n",
				buf.String())
		})

		t.Run("Load", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			p := &dbrPerson{}
			_, err := u.Interpolate().Load(context.TODO(), p)
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Load conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ06*/ `name`, `email` FROM `dbr_people` AS `dp2` WHERE (`id` IN (61,81))\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := u.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ06*/ `name`, `email` FROM `dbr_people` AS `dp2` WHERE (`id` IN (?,?))\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			tx, err := conn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.NoError(t, tx.Wrap(func() error {
				rows, err := tx.SelectFrom("dbr_people").AddColumns("name", "email").Where(dbr.Column("id").In().Int64s(71, 91)).
					Interpolate().Query(context.TODO())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\"\nDEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\" select_id: \"UNIQ08\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ08*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (71,91))\"\nDEBUG Commit conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()
			tx, err := conn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.Error(t, tx.Wrap(func() error {
				rows, err := tx.SelectFrom("dbr_people").AddColumns("name", "email").Where(dbr.Column("id").In().PlaceHolder()).
					Interpolate().Query(context.TODO())
				if err != nil {
					return err
				}
				return rows.Close()
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\"\nDEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\" select_id: \"UNIQ10\" table: \"dbr_people\" duration: 0 sql: \"SELECT /*ID:UNIQ10*/ `name`, `email` FROM `dbr_people` WHERE (`id` IN (?))\"\nDEBUG Rollback conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\" duration: 0\n",
				buf.String())
		})
	})
}
