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

package dml_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelect_QueryContext(t *testing.T) {

	t.Run("ToSQL Error because empty select", func(t *testing.T) {
		sel := (&dml.Select{}).WithDB(dbMock{})
		rows, err := sel.WithArgs().QueryContext(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.Empty.Match(err))
	})

	t.Run("Error", func(t *testing.T) {
		sel := &dml.Select{
			BuilderBase: dml.BuilderBase{
				Table: dml.MakeIdentifier("tableX"),
			},
		}
		sel.AddColumns("a", "b")
		sel.DB = dbMock{
			error: errors.AlreadyClosed.Newf("Who closed myself?"),
		}

		rows, err := sel.WithArgs().QueryContext(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})

	t.Run("Success", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		smr := sqlmock.NewRows([]string{"a"}).AddRow("row1").AddRow("row2")
		dbMock.ExpectQuery("SELECT `a` FROM `tableX`").WillReturnRows(smr)

		sel := dml.NewSelect("a").From("tableX")
		sel.DB = dbc.DB
		rows, err := sel.WithArgs().QueryContext(context.TODO())
		require.NoError(t, err, "%+v", err)
		defer dmltest.Close(t, rows)

		var xx []string
		for rows.Next() {
			var x string
			require.NoError(t, rows.Scan(&x))
			xx = append(xx, x)
		}
		assert.Exactly(t, []string{"row1", "row2"}, xx)
	})
}

var (
	_ dml.ColumnMapper = (*TableCoreConfigDataSlice)(nil)
	_ dml.ColumnMapper = (*TableCoreConfigData)(nil)
	_ io.Closer        = (*TableCoreConfigDataSlice)(nil)
)

// TableCoreConfigDataSlice represents a collection type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigDataSlice struct {
	Data []*TableCoreConfigData
	err  error // just for testing not needed otherwise
}

// TableCoreConfigData represents a type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigData struct {
	ConfigID int64          `json:",omitempty"` // config_id int(10) unsigned NOT NULL PRI  auto_increment
	Scope    string         `json:",omitempty"` // scope varchar(8) NOT NULL MUL DEFAULT 'default'
	ScopeID  int64          `json:",omitempty"` // scope_id int(11) NOT NULL  DEFAULT '0'
	Path     string         `json:",omitempty"` // path varchar(255) NOT NULL  DEFAULT 'general'
	Value    dml.NullString `json:",omitempty"` // value text NULL
}

func (p *TableCoreConfigData) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Int64(&p.ConfigID).String(&p.Scope).Int64(&p.ScopeID).String(&p.Path).NullString(&p.Value).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "config_id":
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
			return errors.NotFound.Newf("[dml] Field %q not found", c)
		}
	}
	return cm.Err()
}

func (ps *TableCoreConfigDataSlice) MapColumns(cm *dml.ColumnMap) error {
	if ps.err != nil {
		return ps.err
	}
	switch m := cm.Mode(); m {
	case dml.ColumnMapScan:
		// case for scanning when loading certain rows, hence we write data from
		// the DB into the struct in each for-loop.
		if cm.Count == 0 {
			ps.Data = ps.Data[:0]
		}
		p := new(TableCoreConfigData)
		if err := p.MapColumns(cm); err != nil {
			return errors.WithStack(err)
		}
		ps.Data = append(ps.Data, p)
	case dml.ColumnMapCollectionReadSet, dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
		// noop not needed
	default:
		return errors.NotSupported.Newf("[dml] Unknown Mode: %q", string(m))
	}
	return nil
}

func (ps *TableCoreConfigDataSlice) Close() error {
	return ps.err
}

func TestSelect_Load(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery("SELECT").WillReturnRows(dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data.csv")))
		s := dml.NewSelect("*").From("core_config_data")
		s.DB = dbc.DB

		ccd := &TableCoreConfigDataSlice{}

		_, err := s.WithArgs().Load(context.TODO(), ccd)
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
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		r := sqlmock.NewRows([]string{"config_id"}).FromCSVString("222\n333\n").
			RowError(1, errors.ConnectionFailed.Newf("Con failed"))
		dbMock.ExpectQuery("SELECT").WillReturnRows(r)
		s := dml.NewSelect("config_id").From("core_config_data")
		s.DB = dbc.DB

		ccd := &TableCoreConfigDataSlice{}
		_, err := s.WithArgs().Load(context.TODO(), ccd)
		assert.True(t, errors.ConnectionFailed.Match(err), "%+v", err)
	})

	t.Run("ioClose error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		r := sqlmock.NewRows([]string{"config_id"}).FromCSVString("222\n333\n").AddRow("3456")
		dbMock.ExpectQuery("SELECT").WillReturnRows(r)
		s := dml.NewSelect("config_id").From("core_config_data")
		s.DB = dbc.DB

		ccd := &TableCoreConfigDataSlice{
			err: errors.Duplicated.Newf("Somewhere exists a duplicate entry"),
		}
		_, err := s.WithArgs().Load(context.TODO(), ccd)
		assert.True(t, errors.Duplicated.Match(err), "%+v", err)
	})
}

func TestSelect_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL Error", func(t *testing.T) {
		sel := dml.NewSelect()
		stmt, err := sel.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.Empty.Match(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectPrepare("SELECT `a`, `b` FROM `tableX`").WillReturnError(errors.AlreadyClosed.Newf("Who closed myself?"))

		sel := dml.NewSelect("a", "b").From("tableX").WithDB(dbc.DB)
		stmt, err := sel.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})

	t.Run("Prepare IN", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `id` FROM `tableX` WHERE (`id` IN (?,?))"))
		prep.ExpectQuery().WithArgs(3739, 3740).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(78))

		sel := dml.NewSelect("id").From("tableX").Where(dml.Column("id").In().PlaceHolders(2)).WithDB(dbc.DB)
		stmt, err := sel.Prepare(context.TODO())
		require.NoError(t, err)
		ints, err := stmt.WithArgs(3739, 3740).LoadInt64s(context.TODO())
		require.NoError(t, err)
		assert.Exactly(t, []int64{78}, ints)
	})

	t.Run("Query", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `name`, `email` FROM `dml_person` WHERE (`id` = ?)"))
		prep.ExpectQuery().WithArgs(6789).
			WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher", "peter@gopher.go"))

		prep.ExpectQuery().WithArgs(6790).
			WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher2", "peter@gopher.go2"))

		stmt, err := dml.NewSelect("name", "email").From("dml_person").
			Where(dml.Column("id").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer dmltest.Close(t, stmt)

		t.Run("Context", func(t *testing.T) {

			rows, err := stmt.WithArgs().QueryContext(context.TODO(), 6789)
			require.NoError(t, err)
			defer dmltest.Close(t, rows)

			cols, err := rows.Columns()
			require.NoError(t, err)
			assert.Exactly(t, []string{"name", "email"}, cols)
		})

		t.Run("RowContext", func(t *testing.T) {
			row := stmt.WithArgs().QueryRowContext(context.TODO(), 6790)
			require.NoError(t, err)
			n, e := "", ""
			require.NoError(t, row.Scan(&n, &e))

			assert.Exactly(t, "Peter Gopher2", n)
			assert.Exactly(t, "peter@gopher.go2", e)
		})
	})

	t.Run("Exec", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `name`, `email` FROM `dml_person` WHERE (`id` = ?)"))

		stmt, err := dml.NewSelect("name", "email").From("dml_person").
			Where(dml.Column("id").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer dmltest.Close(t, stmt)

		const iterations = 3

		t.Run("WithArguments", func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				prep.ExpectQuery().WithArgs(6899).
					WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher", "peter@gopher.go"))
			}
			// use loop with Query+ and add args before
			stmtA := stmt.WithArgs().Int(6899)

			for i := 0; i < iterations; i++ {
				rows, err := stmtA.QueryContext(context.TODO())
				require.NoError(t, err)

				cols, err := rows.Columns()
				require.NoError(t, err)
				assert.Exactly(t, []string{"name", "email"}, cols)
				dmltest.Close(t, rows)
			}
		})

		t.Run("WithRecords_OK", func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				prep.ExpectQuery().WithArgs(6900).
					WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher2", "peter@gopher.go2"))
			}

			p := &dmlPerson{ID: 6900}
			stmtA := stmt.WithArgs().Record("", p)

			for i := 0; i < iterations; i++ {
				rows, err := stmtA.QueryContext(context.TODO())
				require.NoError(t, err)

				cols, err := rows.Columns()
				require.NoError(t, err)
				assert.Exactly(t, []string{"name", "email"}, cols)
				dmltest.Close(t, rows)
			}
		})

		t.Run("WithRecords_Error", func(t *testing.T) {
			p := &TableCoreConfigDataSlice{err: errors.Duplicated.Newf("Found a duplicate")}
			stmtA := stmt.WithArgs().Record("", p)
			rows, err := stmtA.QueryContext(context.TODO())
			assert.True(t, errors.Duplicated.Match(err), "%+v", err)
			assert.Nil(t, rows)
		})
	})

	t.Run("Load", func(t *testing.T) {

		t.Run("multi rows", func(t *testing.T) {
			dbc, dbMock := dmltest.MockDB(t)
			defer dmltest.MockClose(t, dbc, dbMock)

			prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `config_id`, `scope_id`, `path` FROM `core_config_data` WHERE (`config_id` IN ?)"))

			stmt, err := dml.NewSelect("config_id", "scope_id", "path").From("core_config_data").
				Where(dml.Column("config_id").In().PlaceHolder()).
				WithDB(dbc.DB).
				Prepare(context.TODO())
			require.NoError(t, err, "failed creating a prepared statement")
			defer dmltest.Close(t, stmt)

			columns := []string{"config_id", "scope_id", "path"}

			prep.ExpectQuery().WithArgs(345).
				WillReturnRows(sqlmock.NewRows(columns).AddRow(3, 4, "a/b/c").AddRow(4, 4, "a/b/d"))

			ccd := &TableCoreConfigDataSlice{}

			rc, err := stmt.WithArgs(345).Load(context.TODO(), ccd)
			require.NoError(t, err)
			assert.Exactly(t, uint64(2), rc)

			assert.Exactly(t, "&{3  4 a/b/c {{ false}}}", fmt.Sprintf("%v", ccd.Data[0]))
			assert.Exactly(t, "&{4  4 a/b/d {{ false}}}", fmt.Sprintf("%v", ccd.Data[1]))
		})

		t.Run("Int64", func(t *testing.T) {
			dbc, dbMock := dmltest.MockDB(t)
			defer dmltest.MockClose(t, dbc, dbMock)

			prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `scope_id` FROM `core_config_data` WHERE (`config_id` = ?)"))

			stmt, err := dml.NewSelect("scope_id").From("core_config_data").
				Where(dml.Column("config_id").PlaceHolder()).
				WithDB(dbc.DB).
				Prepare(context.TODO())
			require.NoError(t, err, "failed creating a prepared statement")
			defer dmltest.Close(t, stmt)

			columns := []string{"scope_id"}

			prep.ExpectQuery().WithArgs(346).WillReturnRows(sqlmock.NewRows(columns).AddRow(35))

			val, err := stmt.WithArgs().Int64(346).LoadInt64(context.TODO())
			require.NoError(t, err)
			assert.Exactly(t, int64(35), val)
		})

		t.Run("Int64s", func(t *testing.T) {
			dbc, dbMock := dmltest.MockDB(t)
			defer dmltest.MockClose(t, dbc, dbMock)

			prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `scope_id` FROM `core_config_data` WHERE (`config_id` IN ?)"))

			stmt, err := dml.NewSelect("scope_id").From("core_config_data").
				Where(dml.Column("config_id").In().PlaceHolder()).
				WithDB(dbc.DB).
				Prepare(context.TODO())
			require.NoError(t, err, "failed creating a prepared statement")
			defer dmltest.Close(t, stmt)

			columns := []string{"scope_id"}

			prep.ExpectQuery().WithArgs(346, 347).WillReturnRows(sqlmock.NewRows(columns).AddRow(36).AddRow(37))

			val, err := stmt.WithArgs().Int64s(346, 347).LoadInt64s(context.TODO())
			require.NoError(t, err)
			assert.Exactly(t, []int64{36, 37}, val)
		})

	})
}

func TestSelect_Argument_Iterate(t *testing.T) {
	dbc := createRealSession(t)
	defer dmltest.Close(t, dbc)

	rowCount, err := dbc.SelectFrom("dml_fake_person").Count().WithArgs().LoadInt64(context.Background())
	require.NoError(t, err)
	if rowCount < 10000 {
		t.Fatalf("dml_fake_person table contains less than 10k items, seems not to be installed. Got %d items", rowCount)
	}

	t.Run("error in mapper", func(t *testing.T) {
		err := dbc.SelectFrom("dml_fake_person").AddColumns("id", "weight", "height", "update_time").
			Limit(5).
			OrderBy("id").WithArgs().Iterate(context.Background(), func(cm *dml.ColumnMap) error {
			return errors.Blocked.Newf("Mapping blocked")
		})
		assert.True(t, errors.Is(err, errors.Blocked), "Error should have kind errors.Blocked")
	})

	t.Run("no rows but callback gets called", func(t *testing.T) {
		err := dbc.SelectFrom("dml_fake_person").AddColumns("id", "weight", "height", "update_time").
			Where(dml.Column("id").Int(-1000)).
			OrderBy("id").WithArgs().Iterate(context.Background(), func(cm *dml.ColumnMap) error {
			return errors.NotAcceptable.Newf("Mapping blocked")
		})
		assert.True(t, errors.Is(err, errors.NotAcceptable), "Error should have kind errors.NotAcceptable")
	})

	t.Run("iterate serial serial", func(t *testing.T) {
		selExec := dbc.SelectFrom("dml_fake_person").AddColumns("id", "weight", "height", "update_time").
			Limit(500).OrderBy("id").WithArgs()

		err := selExec.Iterate(context.Background(), func(cm *dml.ColumnMap) error {

			fp := &fakePerson{}
			if err := fp.MapColumns(cm); err != nil {
				return err
			}

			if fp.Weight < 1 || fp.Height < 1 || fp.ID < 0 || fp.UpdateTime.IsZero() {
				return errors.NotValid.Newf("failed to load fakePerson: one of the four fields (id,weight,height,update_time) is empty: %#v", fp)
			}

			return nil
		})
		require.NoError(t, err)
	})

	t.Run("iterate serial parallel WG (query)", func(t *testing.T) {
		const limit = 5
		sel := dbc.SelectFrom("dml_fake_person").AddColumns("id", "weight", "height", "update_time").
			Where(
				dml.Column("id").Between().PlaceHolder(),
			).
			Limit(limit).OrderBy("id")

		const iterations = 10
		var wg sync.WaitGroup
		wg.Add(iterations)
		for i := 0; i < iterations; i++ {

			go func(wg *sync.WaitGroup, i int) {
				defer wg.Done()

				err := sel.WithArgs().Int(i).Int(i+5).Iterate(context.Background(), func(cm *dml.ColumnMap) error {

					fp := &fakePerson{}
					if err := fp.MapColumns(cm); err != nil {
						return err
					}
					fmt.Printf("%d: %#v\n", i, fp)
					if fp.Weight < 1 || fp.Height < 1 || fp.ID < i || fp.UpdateTime.IsZero() {
						return errors.NotValid.Newf("failed to load fakePerson: one of the four fields (id,weight,height,update_time) is empty: %#v", fp)
					}

					return nil
				})
				require.NoError(t, err)

			}(&wg, i*10)
		}
		wg.Wait()
	})

	t.Run("iterate serial parallel WG (prepared)", func(t *testing.T) {

	})
}
