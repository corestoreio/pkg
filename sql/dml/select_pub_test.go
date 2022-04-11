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
	"sync/atomic"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
)

type fakePerson struct {
	ID         int
	FirstName  string
	LastName   string
	Sex        string
	BirthDate  time.Time
	Weight     int
	Height     int
	UpdateTime time.Time
}

// MapColumns implements interface ColumnMapper only partially.
func (p *fakePerson) MapColumns(cm *dml.ColumnMap) error {
	for cm.Next(8) {
		switch c := cm.Column(); c {
		case "id", "0":
			cm.Int(&p.ID)
		case "first_name", "1":
			cm.String(&p.FirstName)
		case "last_name", "2":
			cm.String(&p.LastName)
		case "sex", "3":
			cm.String(&p.Sex)
		case "birth_date", "4":
			cm.Time(&p.BirthDate)
		case "weight", "5":
			cm.Int(&p.Weight)
		case "height", "6":
			cm.Int(&p.Height)
		case "update_time", "7":
			cm.Time(&p.UpdateTime)
		default:
			return fmt.Errorf("[dml_test] fakePerson Column %q not found", c)
		}
	}
	return cm.Err()
}

type fakePersons struct {
	Data []fakePerson
}

func (cc *fakePersons) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {

	case dml.ColumnMapScan:
		if cm.Count == 0 {
			cc.Data = cc.Data[:0]
		}
		var p fakePerson
		if err := p.MapColumns(cm); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, p)

	default:
		return fmt.Errorf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

func TestSelect_QueryContext(t *testing.T) {
	t.Run("ToSQL Error because empty select", func(t *testing.T) {
		sel := &dml.Select{}
		rows, err := sel.WithDBR(dbMock{}).QueryContext(context.Background())
		assert.Nil(t, rows)
		assert.Error(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		sel := &dml.Select{
			BuilderBase: dml.BuilderBase{
				Table: dml.MakeIdentifier("tableX"),
			},
		}
		sel.AddColumns("a", "b")
		selDBR := sel.WithDBR(dbMock{
			error: errors.New("Who closed myself?"),
		})

		rows, err := selDBR.QueryContext(context.Background())
		assert.Nil(t, rows)
		assert.Error(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		smr := sqlmock.NewRows([]string{"a"}).AddRow("row1").AddRow("row2")
		dbMock.ExpectQuery("SELECT `a` FROM `tableX`").WillReturnRows(smr)

		sel := dbc.WithQueryBuilder(dml.NewSelect("a").From("tableX"))
		rows, err := sel.QueryContext(context.Background())
		assert.NoError(t, err)
		defer dmltest.Close(t, rows)

		var xx []string
		for rows.Next() {
			var x string
			assert.NoError(t, rows.Scan(&x))
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
	ConfigID int64       `json:",omitempty"` // config_id int(10) unsigned NOT NULL PRI  auto_increment
	Scope    string      `json:",omitempty"` // scope varchar(8) NOT NULL MUL DEFAULT 'default'
	ScopeID  int64       `json:",omitempty"` // scope_id int(11) NOT NULL  DEFAULT '0'
	Path     string      `json:",omitempty"` // path varchar(255) NOT NULL  DEFAULT 'general'
	Value    null.String `json:",omitempty"` // value text NULL
}

func (p *TableCoreConfigData) MapColumns(cm *dml.ColumnMap) error {
	for cm.Next(5) {
		switch c := cm.Column(); c {
		case "config_id", "0":
			cm.Int64(&p.ConfigID)
		case "scope", "1":
			cm.String(&p.Scope)
		case "scope_id", "2":
			cm.Int64(&p.ScopeID)
		case "path", "3":
			cm.String(&p.Path)
		case "value", "4":
			cm.NullString(&p.Value)
		default:
			return fmt.Errorf("[dml] Field %q not found", c)
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
		return fmt.Errorf("[dml] Unknown Mode: %q", string(m))
	}
	return nil
}

func (ps *TableCoreConfigDataSlice) Close() error {
	return ps.err
}

func TestSelect_Load(t *testing.T) {
	t.Run("success no condition", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT * FROM `core_config_data`")).
			WillReturnRows(dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data.csv")))

		ccd := &TableCoreConfigDataSlice{}

		_, err := dbc.WithQueryBuilder(dml.NewSelect("*").From("core_config_data")).Load(context.Background(), ccd)
		assert.NoError(t, err)

		buf := new(bytes.Buffer)
		je := json.NewEncoder(buf)

		for _, c := range ccd.Data {
			if err := je.Encode(c); err != nil {
				t.Fatalf("%+v", err)
			}
		}
		assert.Exactly(t, "{\"ConfigID\":2,\"Scope\":\"default\",\"Path\":\"web/unsecure/base_url\",\"Value\":\"http://mgeto2.local/\"}\n{\"ConfigID\":3,\"Scope\":\"website\",\"ScopeID\":11,\"Path\":\"general/locale/code\",\"Value\":\"en_US\"}\n{\"ConfigID\":4,\"Scope\":\"default\",\"Path\":\"general/locale/timezone\",\"Value\":\"Europe/Berlin\"}\n{\"ConfigID\":5,\"Scope\":\"default\",\"Path\":\"currency/options/base\",\"Value\":\"EUR\"}\n{\"ConfigID\":15,\"Scope\":\"store\",\"ScopeID\":33,\"Path\":\"design/head/includes\",\"Value\":\"\\u003clink  rel=\\\"stylesheet\\\" type=\\\"text/css\\\" href=\\\"{{MEDIA_URL}}styles.css\\\" /\\u003e\"}\n{\"ConfigID\":16,\"Scope\":\"default\",\"Path\":\"admin/security/use_case_sensitive_login\",\"Value\":null}\n{\"ConfigID\":17,\"Scope\":\"default\",\"Path\":\"admin/security/session_lifetime\",\"Value\":\"90000\"}\n",
			buf.String())
	})

	t.Run("success In with one ARG", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `config_id` FROM `core_config_data` WHERE (`config_id` IN (?)")).
			WillReturnRows(dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_ints.csv"))).
			WithArgs(199)

		var dst []int64
		dst, err := dbc.WithQueryBuilder(dml.NewSelect("config_id").From("core_config_data").Where(
			dml.Column("config_id").In().PlaceHolder(),
		)).ExpandPlaceHolders().LoadInt64s(context.Background(), dst, []int64{199})
		assert.NoError(t, err)

		// wrong result set for correct query. maybe some one can fix the returned data.
		assert.Exactly(t, []int64{2, 3, 4, 16, 17}, dst)
	})
	t.Run("success In with two ARGs", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `config_id` FROM `core_config_data` WHERE (`config_id` IN (?,?))")).
			WillReturnRows(dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_ints.csv"))).
			WithArgs(199, 217)

		var dst []int64
		dst, err := dbc.WithQueryBuilder(dml.NewSelect("config_id").From("core_config_data").Where(
			dml.Column("config_id").In().PlaceHolder(),
		)).ExpandPlaceHolders().LoadInt64s(context.Background(), dst, []int64{199, 217})
		assert.NoError(t, err)

		assert.Exactly(t, []int64{2, 3, 4, 16, 17}, dst)
	})

	t.Run("row error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		r := sqlmock.NewRows([]string{"config_id"}).FromCSVString("222\n333\n").
			RowError(1, fmt.Errorf("Con failed"))
		dbMock.ExpectQuery("SELECT").WillReturnRows(r)

		ccd := &TableCoreConfigDataSlice{}
		_, err := dbc.WithQueryBuilder(dml.NewSelect("config_id").From("core_config_data")).Load(context.Background(), ccd)
		assert.Error(t, err)
	})

	t.Run("ioClose error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		r := sqlmock.NewRows([]string{"config_id"}).FromCSVString("222\n333\n").AddRow("3456")
		dbMock.ExpectQuery("SELECT").WillReturnRows(r)

		ccd := &TableCoreConfigDataSlice{
			err: fmt.Errorf("Somewhere exists a duplicate entry"),
		}
		_, err := dbc.WithQueryBuilder(dml.NewSelect("config_id").From("core_config_data")).Load(context.Background(), ccd)
		assert.Error(t, err)
	})
}

func TestSelect_Prepare(t *testing.T) {
	t.Run("Prepare Error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectPrepare("SELECT `a`, `b` FROM `tableX`").WillReturnError(fmt.Errorf("Who closed myself?"))

		stmt := dbc.WithPrepare(context.Background(), dml.NewSelect("a", "b").From("tableX"))
		assert.NotNil(t, stmt)

		assert.Error(t, stmt.PreviousError())
	})

	t.Run("Prepare IN", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `id` FROM `tableX` WHERE (`id` IN (?,?))"))
		prep.ExpectQuery().WithArgs(3739, 3740).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(78))

		stmt := dbc.WithPrepare(context.Background(), dml.NewSelect("id").From("tableX").Where(dml.Column("id").In().PlaceHolders(2)))
		ints, err := stmt.LoadInt64s(context.Background(), nil, 3739, 3740)
		assert.NoError(t, err)
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

		stmt := dbc.WithPrepare(context.Background(),
			dml.NewSelect("name", "email").From("dml_person").Where(dml.Column("id").PlaceHolder()),
		)
		defer dmltest.Close(t, stmt)

		t.Run("Context", func(t *testing.T) {
			rows, err := stmt.QueryContext(context.Background(), 6789)
			assert.NoError(t, err)
			defer dmltest.Close(t, rows)

			cols, err := rows.Columns()
			assert.NoError(t, err)
			assert.Exactly(t, []string{"name", "email"}, cols)
		})

		t.Run("RowContext", func(t *testing.T) {
			row := stmt.QueryRowContext(context.Background(), 6790)
			n, e := "", ""
			assert.NoError(t, row.Scan(&n, &e))

			assert.Exactly(t, "Peter Gopher2", n)
			assert.Exactly(t, "peter@gopher.go2", e)
		})
	})

	t.Run("Records in final args", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `name`, `email` FROM `dml_person` WHERE (`id` = ?) AND (`name` = ?)"))
		prep.ExpectQuery().WithArgs(4211, "Peter Gopher").
			WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher", "peter@gopher.go"))

		stmt := dbc.WithPrepare(context.Background(), dml.NewSelect("name", "email").From("dml_person").
			Where(dml.Column("id").PlaceHolder(), dml.Column("name").PlaceHolder()))

		t.Run("Context", func(t *testing.T) {
			p := &dmlPerson{
				ID:   4211,
				Name: "Peter Gopher",
			}

			rows, err := stmt.QueryContext(context.Background(), dml.Qualify("", p))
			assert.NoError(t, err)
			defer dmltest.Close(t, rows)

			cols, err := rows.Columns()
			assert.NoError(t, err)
			assert.Exactly(t, []string{"name", "email"}, cols)
		})
	})

	t.Run("QueryContext", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `name`, `email` FROM `dml_person` WHERE (`id` = ?)"))

		stmt := dbc.WithPrepare(context.Background(), dml.NewSelect("name", "email").From("dml_person").
			Where(dml.Column("id").PlaceHolder()))
		defer dmltest.Close(t, stmt)

		const iterations = 3

		t.Run("WithArguments", func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				prep.ExpectQuery().WithArgs(6899).
					WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher", "peter@gopher.go"))
			}
			// use loop with Query+ and add args before

			for i := 0; i < iterations; i++ {
				rows, err := stmt.QueryContext(context.Background(), 6899)
				assert.NoError(t, err)

				cols, err := rows.Columns()
				assert.NoError(t, err)
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

			for i := 0; i < iterations; i++ {
				rows, err := stmt.QueryContext(context.Background(), dml.Qualify("", p))
				assert.NoError(t, err)

				cols, err := rows.Columns()
				assert.NoError(t, err)
				assert.Exactly(t, []string{"name", "email"}, cols)
				dmltest.Close(t, rows)
			}
		})

		t.Run("WithRecords_Error", func(t *testing.T) {
			p := &TableCoreConfigDataSlice{err: fmt.Errorf("Found a duplicate")}
			rows, err := stmt.QueryContext(context.Background(), dml.Qualify("", p))
			assert.Error(t, err)
			assert.Nil(t, rows)
		})
	})

	t.Run("Load", func(t *testing.T) {
		t.Run("multi rows", func(t *testing.T) {
			dbc, dbMock := dmltest.MockDB(t)
			defer dmltest.MockClose(t, dbc, dbMock)

			prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `config_id`, `scope_id`, `path` FROM `core_config_data` WHERE (`config_id` IN ?)"))

			stmt := dbc.WithPrepare(context.Background(), dml.NewSelect("config_id", "scope_id", "path").From("core_config_data").
				Where(dml.Column("config_id").In().PlaceHolder()))
			defer dmltest.Close(t, stmt)

			columns := []string{"config_id", "scope_id", "path"}

			prep.ExpectQuery().WithArgs(345).
				WillReturnRows(sqlmock.NewRows(columns).AddRow(3, 4, "a/b/c").AddRow(4, 4, "a/b/d"))

			ccd := &TableCoreConfigDataSlice{}

			rc, err := stmt.Load(context.Background(), ccd, 345)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(2), rc)

			assert.Exactly(t, "&{3  4 a/b/c null}", fmt.Sprintf("%v", ccd.Data[0]))
			assert.Exactly(t, "&{4  4 a/b/d null}", fmt.Sprintf("%v", ccd.Data[1]))
		})

		t.Run("Int64", func(t *testing.T) {
			dbc, dbMock := dmltest.MockDB(t)
			defer dmltest.MockClose(t, dbc, dbMock)

			prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `scope_id` FROM `core_config_data` WHERE (`config_id` = ?)"))

			stmt := dbc.WithPrepare(context.Background(), dml.NewSelect("scope_id").From("core_config_data").
				Where(dml.Column("config_id").PlaceHolder()))
			defer dmltest.Close(t, stmt)

			columns := []string{"scope_id"}

			prep.ExpectQuery().WithArgs(346).WillReturnRows(sqlmock.NewRows(columns).AddRow(35))

			val, found, err := stmt.LoadNullInt64(context.Background(), 346)
			assert.NoError(t, err)
			assert.True(t, found)
			assert.Exactly(t, null.MakeInt64(35), val)
		})

		t.Run("Int64s", func(t *testing.T) {
			dbc, dbMock := dmltest.MockDB(t)
			defer dmltest.MockClose(t, dbc, dbMock)

			prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `scope_id` FROM `core_config_data` WHERE (`config_id` IN ?)"))

			stmt := dbc.WithPrepare(context.Background(), dml.NewSelect("scope_id").From("core_config_data").
				Where(dml.Column("config_id").In().PlaceHolder()))
			defer dmltest.Close(t, stmt)

			columns := []string{"scope_id"}

			prep.ExpectQuery().WithArgs(346, 347).WillReturnRows(sqlmock.NewRows(columns).AddRow(36).AddRow(37))

			val, err := stmt.LoadInt64s(context.Background(), nil, []int64{346, 347})
			assert.NoError(t, err)
			assert.Exactly(t, []int64{36, 37}, val)
		})
	})
}

func TestSelect_Argument_Iterate(t *testing.T) {
	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)
	defer dmltest.SQLDumpLoad(t, "testdata/person_ffaker*", &dmltest.SQLDumpOptions{DSN: dbc.DSN()}).Deferred()

	rowCount, found, err := dbc.WithQueryBuilder(dml.NewSelect().Count().From("dml_fake_person")).LoadNullInt64(context.Background())
	assert.NoError(t, err)
	assert.True(t, found)
	if rowCount.Int64 < 10000 {
		t.Skipf("dml_fake_person table contains less than 10k items, seems not to be installed. Got %d items", rowCount.Int64)
	}

	t.Run("serial", func(t *testing.T) {
		t.Run("error in mapper", func(t *testing.T) {
			err := dbc.WithQueryBuilder(
				dml.NewSelect().AddColumns("id", "weight", "height", "update_time").From("dml_fake_person").
					Limit(0, 5).OrderBy("id"),
			).IterateSerial(context.Background(), func(cm *dml.ColumnMap) error {
				return fmt.Errorf("Mapping blocked")
			})
			assert.Error(t, err)
		})

		t.Run("serial serial", func(t *testing.T) {
			const rowCount = 500
			selExec := dbc.WithQueryBuilder(dml.NewSelect().From("dml_fake_person").AddColumns("id", "weight", "height", "update_time").
				OrderByRandom("id", rowCount))

			var counter int
			err := selExec.IterateSerial(context.Background(), func(cm *dml.ColumnMap) error {
				fp := &fakePerson{}
				if err := fp.MapColumns(cm); err != nil {
					return err
				}

				if fp.Weight < 1 || fp.Height < 1 || fp.ID < 0 || fp.UpdateTime.IsZero() {
					return fmt.Errorf("failed to load fakePerson: one of the four fields (id,weight,height,update_time) is empty: %#v", fp)
				}
				counter++
				return nil
			})
			assert.NoError(t, err)
			assert.Exactly(t, rowCount, counter, "Should have loaded %d rows", rowCount)
		})

		t.Run("serial parallel", func(t *testing.T) {
			// Do not run such a construct in production.

			const limit = 5
			const concurrencyLevel = 10

			processFakePerson := func(selProc *dml.DBR, i int) {
				// testing.T does not work in parallel context, so we use panic :-(

				fp := &fakePerson{}
				var counter int
				err := selProc.IterateSerial(context.Background(), func(cm *dml.ColumnMap) error {
					if err := fp.MapColumns(cm); err != nil {
						return err
					}
					// fmt.Printf("%d: %#v\n", i, fp)
					if fp.Weight < 1 || fp.Height < 1 || fp.ID < i || fp.UpdateTime.IsZero() {
						return fmt.Errorf("failed to load fakePerson: one of the four fields (id,weight,height,update_time) is empty: %#v", fp)
					}
					counter++
					return nil
				}, i, i+5)
				ifErrPanic(err)
				ifNotEqualPanic(counter, limit, "Should have loaded this amount of rows")
			}
			// dbc.Schema()+". that is a hack :-(
			fpSel := dbc.WithQueryBuilder(dml.NewSelect("id", "weight", "height", "update_time").From(dbc.Schema()+".dml_fake_person").
				Where(
					dml.Column("id").Between().PlaceHolder(),
				).
				Limit(0, limit).OrderBy("id")).Interpolate()

			t.Run("WG01 (query, conn pool)", func(t *testing.T) {
				bgwork.Wait(concurrencyLevel, func(index int) {
					// Every goroutine creates its own underlying connection to the
					// SQL server. This makes sense because each dataset is unique.
					processFakePerson(fpSel, index*concurrencyLevel)
				})
			})

			t.Run("WG02 (prepared,multi conn)", func(t *testing.T) {
				stmt, err := fpSel.Prepare(context.Background())
				assert.NoError(t, err)
				defer dmltest.Close(t, stmt)

				bgwork.Wait(concurrencyLevel, func(index int) {
					// Every goroutine creates its own underlying connection to the
					// SQL server and prepares its own statement in the SQL server
					// despite having one single pointer to *sql.Stmt.
					processFakePerson(stmt, index*concurrencyLevel)
				})
			})
		})
	})

	const concurrencyLevel = 4

	t.Run("parallel", func(t *testing.T) {
		t.Run("error wrong concurrency level", func(t *testing.T) {
			err := dbc.WithQueryBuilder(dml.NewSelect().From("dml_fake_person").AddColumns("id", "weight", "height", "update_time").
				Limit(0, 50).OrderBy("id")).IterateParallel(context.Background(), 0, func(cm *dml.ColumnMap) error {
				return nil
			})
			assert.Error(t, err)
		})

		t.Run("error in mapper of all workers", func(t *testing.T) {
			err := dbc.WithQueryBuilder(dml.NewSelect().From("dml_fake_person").AddColumns("id", "weight", "height", "update_time").
				Limit(0, 50).OrderBy("id")).IterateParallel(context.Background(), concurrencyLevel, func(cm *dml.ColumnMap) error {
				return fmt.Errorf("Mapping blocked error")
			})

			assert.Error(t, err)
		})

		t.Run("successful 40 rows fibonacci", func(t *testing.T) {
			sel01 := dml.NewSelect("id", "weight", "height", "update_time").From(dbc.Schema()+".dml_fake_person").
				Where(dml.Column("id").LessOrEqual().Int(60)).
				Limit(0, 40)
			sel01.IsOrderByRand = true

			rowsLoadedCounter := new(int32)
			err := dbc.WithQueryBuilder(sel01).IterateParallel(context.Background(), concurrencyLevel, func(cm *dml.ColumnMap) error {
				var fp fakePerson

				if err := fp.MapColumns(cm); err != nil {
					return err
				}

				if fp.Weight < 1 || fp.Height < 1 || fp.ID < 1 || fp.UpdateTime.IsZero() {
					return fmt.Errorf("failed to load fakePerson: one of the four fields (id,weight,height,update_time) is empty: %#v", fp)
				}

				if fp.ID < 41 {
					_ = fib(uint(fp.ID))
					// fi := fib(uint(fp.ID))
					// println("a 591", int(cm.Count), fp.ID, fi)
				} /* else {
					println("a 597", int(cm.Count), fp.ID)
				} */

				// fmt.Printf("%d: FIB:%d: fakePerson ID %d\n", cm.Count, fib(uint(fp.ID)), fp.ID)
				atomic.AddInt32(rowsLoadedCounter, 1)
				return nil
			})
			assert.NoError(t, err)
			assert.Exactly(t, int32(40), *rowsLoadedCounter, "Should load this amount of rows from the database server.")
		})
	})
}

func fib(n uint) uint {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func TestSelect_Clone(t *testing.T) {
	dbc, dbMock := dmltest.MockDB(t, dml.WithLogger(log.BlackHole{}, func() string { return "uniqueID" }))
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("nil", func(t *testing.T) {
		var s *dml.Select
		s2 := s.Clone()
		assert.Nil(t, s)
		assert.Nil(t, s2)
	})

	t.Run("non-nil", func(t *testing.T) {
		s := dml.NewSelect().FromAlias("dml_people", "p1").AddColumns("p1.*").
			AddColumnsAliases("p2.name", "p2Name", "p2.email", "p2Email").
			RightJoin(
				dml.MakeIdentifier("dml_people").Alias("p2"),
				dml.Columns("id", "email"),
			).
			Where(
				dml.Column("name").Like().PlaceHolder(),
			).
			OrderBy("email").
			GroupBy("last_name").
			Having(
				dml.Column("income").LessOrEqual().PlaceHolder(),
			)

		s2 := s.Clone()
		notEqualPointers(t, s, s2)
		notEqualPointers(t, s.BuilderConditional.Wheres, s2.BuilderConditional.Wheres)
		notEqualPointers(t, s.BuilderConditional.Joins, s2.BuilderConditional.Joins)
		notEqualPointers(t, s.BuilderConditional.OrderBys, s2.BuilderConditional.OrderBys)
		notEqualPointers(t, s.GroupBys, s2.GroupBys)
		notEqualPointers(t, s.Havings, s2.Havings)
		// assert.Exactly(t, s.db, s2.db) // test it via fmt.Sprintf?
	})
}

func TestSelect_When_Unless(t *testing.T) {
	t.Run("true and no default", func(t *testing.T) {
		s := dml.NewSelect("entity_id").From("catalog_product_entity")
		s.When(true, func(s2 *dml.Select) {
			s2.Where(dml.Column("sku").Like().Str("4711"))
		}, nil)
		assert.Exactly(t, "SELECT `entity_id` FROM `catalog_product_entity` WHERE (`sku` LIKE '4711')", s.String())
	})
	t.Run("false and no default", func(t *testing.T) {
		s := dml.NewSelect("entity_id").From("catalog_product_entity")
		s.When(false, func(s2 *dml.Select) {
			s2.Where(dml.Column("sku").Like().Str("4711"))
		}, nil)
		assert.Exactly(t, "SELECT `entity_id` FROM `catalog_product_entity`", s.String())
	})
	t.Run("Unless true and default", func(t *testing.T) {
		s := dml.NewSelect("entity_id").From("catalog_product_entity")
		s.Unless(true, func(s2 *dml.Select) {
			s2.Where(dml.Column("sku").Like().Str("4712"))
		}, func(s2 *dml.Select) {
			s2.Where(dml.Column("sku").Like().Str("4713"))
		})
		assert.Exactly(t, "SELECT `entity_id` FROM `catalog_product_entity` WHERE (`sku` LIKE '4713')", s.String())
	})
}

func TestPrepareWithDBR(t *testing.T) {
	t.Run("all dml types", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `a1` FROM `a`")).WillBeClosed()
		dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("(SELECT * FROM `b1`) UNION (SELECT * FROM `b2`)")).WillBeClosed()
		dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("UPDATE `c` SET `c1`=?")).WillBeClosed()
		dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("WITH `one` AS (SELECT 1) SELECT * FROM `one`")).WillBeClosed()
		dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("DELETE FROM `e`")).WillBeClosed()
		dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `d` (`d1`) VALUES (?)")).WillBeClosed()

		ctx := context.Background()

		for _, a := range []*dml.DBR{
			dbc.WithPrepare(ctx, dml.NewSelect().From("a").AddColumns("a1")),
			dbc.WithPrepare(ctx, dml.NewUnion(dml.NewSelect("*").From("b1"), dml.NewSelect("*").From("b2"))),
			dbc.WithPrepare(ctx, dml.NewUpdate("c").AddColumns("c1")),
			dbc.WithPrepare(ctx, dml.NewWith(
				dml.WithCTE{Name: "one", Select: dml.NewSelect().Unsafe().AddColumns("1")},
			).Select(dml.NewSelect().Star().From("one"))),
			dbc.WithPrepare(ctx, dml.NewDelete("e")),
			dbc.WithPrepare(ctx, dml.NewInsert("d").AddColumns("d1").BuildValues()),
		} {
			dmltest.Close(t, a)
		}
	})

	t.Run("select error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		ctx := context.Background()

		err := dbc.WithPrepare(ctx, dml.NewSelect().From("a")).Close()
		assert.Error(t, err)
	})

	t.Run("union error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		ctx := context.Background()

		err := dbc.WithPrepare(ctx, dml.NewUnion(dml.NewSelect("*").From("b"))).Close()
		assert.Error(t, err)
	})
}

func TestDBR_ExpandTuples(t *testing.T) {
	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)
	ctx := context.Background()

	dbr := dbc.WithQueryBuilder(dml.NewSelect().Star().From("core_config_data").Where(
		dml.Columns("entity_id", "attribute_id", "store_id", "source_id").In().Tuples(),
	)).ExpandPlaceHolders()

	t.Run("1,4 tuple, no interpolate", func(t *testing.T) {
		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT * FROM `core_config_data` WHERE ((`entity_id`, `attribute_id`, `store_id`, `source_id`) IN ((?,?,?,?)))")).
			WithArgs(1, 2, 3, 4).
			WillReturnRows(dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data.csv")))

		ccd := &TableCoreConfigDataSlice{}
		_, err := dbr.Load(ctx, ccd, 1, 2, 3, 4)
		assert.NoError(t, err)
	})

	t.Run("2,4 tuple, no interpolate", func(t *testing.T) {
		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT * FROM `core_config_data` WHERE ((`entity_id`, `attribute_id`, `store_id`, `source_id`) IN ((?,?,?,?),(?,?,?,?)))")).
			WithArgs("b1", 2, 3, 4, "a11", 22, 33, 44).
			WillReturnRows(dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data.csv")))

		ccd := &TableCoreConfigDataSlice{}
		_, err := dbr.Load(ctx, ccd, "b1", 2, 3, 4, "a11", 22, 33, 44)
		assert.NoError(t, err)
	})

	t.Run("1,4 tuple, with interpolate", func(t *testing.T) {
		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT * FROM `core_config_data` WHERE ((`entity_id`, `attribute_id`, `store_id`, `source_id`) IN ((1,2,3,4)))")).
			WithArgs().
			WillReturnRows(dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data.csv")))

		ccd := &TableCoreConfigDataSlice{}
		_, err := dbr.Interpolate().Load(ctx, ccd, 1, 2, 3, 4)
		assert.NoError(t, err)
	})
}
