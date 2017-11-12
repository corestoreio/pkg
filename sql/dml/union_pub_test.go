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

package dml_test

import (
	"bytes"
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log/logw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnion_Query(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL Error", func(t *testing.T) {
		u := dml.NewUnion(
			dml.NewSelect(),
			dml.NewSelect(),
		)
		rows, err := u.Query(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.IsEmpty(err))
	})

	u := dml.NewUnion(
		dml.NewSelect("value").From("eavChar"),
		dml.NewSelect("value").From("eavInt").Where(dml.Column("b").Float64(3.14159)),
	)

	t.Run("Error", func(t *testing.T) {
		u.WithDB(dbMock{
			error: errors.NewConnectionFailedf("Who closed myself?"),
		})
		rows, err := u.Query(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.IsConnectionFailed(err), "%+v", err)
	})

	t.Run("Success", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		smr := sqlmock.NewRows([]string{"value"}).AddRow("row1").AddRow("row2")
		dbMock.ExpectQuery(
			cstesting.SQLMockQuoteMeta("(SELECT `value` FROM `eavChar`) UNION (SELECT `value` FROM `eavInt` WHERE (`b` = ?))"),
		).WillReturnRows(smr)

		u.WithDB(dbc.DB)

		rows, err := u.Query(context.TODO())
		require.NoError(t, err, "%+v", err)

		var xx []string
		for rows.Next() {
			var x string
			require.NoError(t, rows.Scan(&x))
			xx = append(xx, x)
		}
		assert.Exactly(t, []string{"row1", "row2"}, xx)
		require.NoError(t, rows.Close())
	})
}

func TestUnion_Load(t *testing.T) {
	t.Parallel()

	u := dml.NewUnion(
		dml.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
		dml.NewSelect("a", "b").From("tableAB").Where(dml.Column("b").Float64(3.14159)),
	).Unsafe().
		OrderBy("a").OrderByDesc("b").OrderBy(`concat("c",b,"d")`)

	t.Run("error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta("(SELECT `a`, `d` AS `b` FROM `tableAD`) UNION (SELECT `a`, `b` FROM `tableAB` WHERE (`b` = ?)) ORDER BY `a`, `b` DESC, concat(\"c\",b,\"d\")")).
			WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		rows, err := u.WithDB(dbc.DB).Load(context.TODO(), nil)
		assert.Exactly(t, uint64(0), rows)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestUnion_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL Error", func(t *testing.T) {
		u := dml.NewUnion(
			dml.NewSelect(),
			dml.NewSelect(),
		)
		stmt, err := u.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectPrepare(
			cstesting.SQLMockQuoteMeta("(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`) UNION (SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = ?)) ORDER BY `_preserve_result_set`, `a`, `b` DESC, concat(\"c\",b,\"d\")"),
		).
			WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		u := dml.NewUnion(
			dml.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
			dml.NewSelect("a", "b").From("tableAB").Where(dml.Column("b").Float64(3.14159)),
		).
			Unsafe().
			OrderBy("a").OrderByDesc("b").OrderBy(`concat("c",b,"d")`).
			PreserveResultSet().BuildCache().WithDB(dbc.DB)

		stmt, err := u.Prepare(context.TODO())
		require.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

	t.Run("Query", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("(SELECT `a`, `d` AS `b` FROM `tableAD`) UNION (SELECT `a`, `b` FROM `tableAB` WHERE (`b` = ?))"))
		prep.ExpectQuery().WithArgs(6889).
			WillReturnRows(sqlmock.NewRows([]string{"a", "b"}).AddRow("Peter Gopher", "peter@gopher.go"))

		prep.ExpectQuery().WithArgs(6890).
			WillReturnRows(sqlmock.NewRows([]string{"a", "b"}).AddRow("Peter Gopher2", "peter@gopher.go2"))

		stmt, err := dml.NewUnion(
			dml.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
			dml.NewSelect("a", "b").From("tableAB").Where(dml.Column("b").PlaceHolder()),
		).
			BuildCache().WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		t.Run("Context", func(t *testing.T) {

			rows, err := stmt.Query(context.TODO(), 6889)
			require.NoError(t, err)
			defer rows.Close()

			cols, err := rows.Columns()
			require.NoError(t, err)
			assert.Exactly(t, []string{"a", "b"}, cols)
		})

		t.Run("RowContext", func(t *testing.T) {

			row := stmt.QueryRow(context.TODO(), 6890)
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

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("(SELECT `name`, `d` AS `email` FROM `dml_people`) UNION (SELECT `name`, `email` FROM `dml_people2` WHERE (`id` = ?))"))

		stmt, err := dml.NewUnion(
			dml.NewSelect("name").AddColumnsAliases("d", "email").From("dml_people"),
			dml.NewSelect("name", "email").From("dml_people2").Where(dml.Column("id").PlaceHolder()),
		).
			BuildCache().WithDB(dbc.DB).
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
			// use loop with Query and add args before
			stmt.WithArguments(dml.MakeArgs(1).Int(6899))

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

			p := &dmlPerson{ID: 6900}
			stmt.WithRecords(dml.Qualify("", p))

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
			p := &TableCoreConfigDataSlice{err: errors.NewDuplicatedf("Found a duplicate")}
			stmt.WithRecords(dml.Qualify("", p))
			rows, err := stmt.Query(context.TODO())
			assert.True(t, errors.IsDuplicated(err), "%+v", err)
			assert.Nil(t, rows)
		})
	})
}

func TestUnion_WithLogger(t *testing.T) {
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
	require.NoError(t, rConn.Options(dml.WithLogger(lg, uniqueIDFunc)))

	t.Run("ConnPool", func(t *testing.T) {
		u := rConn.Union(
			dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
			dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(6, 8)),
		)

		t.Run("Query", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			rows, err := u.Interpolate().Query(context.TODO())
			require.NoError(t, err)
			require.NoError(t, rows.Close())

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" union_id: \"UNIQ02\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID:UNIQ02*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8)))\"\n",
				buf.String())
		})

		t.Run("Load", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			p := &dmlPerson{}
			_, err := u.Interpolate().Load(context.TODO(), p)
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Load conn_pool_id: \"UNIQ01\" union_id: \"UNIQ02\" tables: \"dml_people, dml_people\" duration: 0 row_count: 0x0 sql: \"(SELECT /*ID:UNIQ02*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8)))\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := u.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ01\" union_id: \"UNIQ02\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID:UNIQ02*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (?,?)))\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			tx, err := rConn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.NoError(t, tx.Wrap(func() error {
				rows, err := tx.Union(
					dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
					dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(7, 9)),
				).Interpolate().Query(context.TODO())

				require.NoError(t, rows.Close())
				return err
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\"\nDEBUG Query conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\" union_id: \"UNIQ04\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID:UNIQ04*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (7,9)))\"\nDEBUG Commit conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\" duration: 0\n",
				buf.String())
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		require.NoError(t, err)

		u := conn.Union(
			dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
			dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(61, 81)),
		)
		t.Run("Query", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()

			rows, err := u.Interpolate().Query(context.TODO())
			require.NoError(t, err)
			require.NoError(t, rows.Close())

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" union_id: \"UNIQ06\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID:UNIQ06*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (61,81)))\"\n",
				buf.String())
		})

		t.Run("Load", func(t *testing.T) {
			defer func() {
				buf.Reset()
				u.IsInterpolate = false
			}()
			p := &dmlPerson{}
			_, err := u.Interpolate().Load(context.TODO(), p)
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Load conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" union_id: \"UNIQ06\" tables: \"dml_people, dml_people\" duration: 0 row_count: 0x0 sql: \"(SELECT /*ID:UNIQ06*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (61,81)))\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := u.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" union_id: \"UNIQ06\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID:UNIQ06*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (?,?)))\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			tx, err := conn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.NoError(t, tx.Wrap(func() error {
				rows, err := tx.Union(
					dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
					dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(71, 91)),
				).Interpolate().Query(context.TODO())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\"\nDEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\" union_id: \"UNIQ08\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID:UNIQ08*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (71,91)))\"\nDEBUG Commit conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()
			tx, err := conn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.Error(t, tx.Wrap(func() error {
				rows, err := tx.Union(
					dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
					dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().PlaceHolder()),
				).Interpolate().Query(context.TODO())
				if err != nil {
					return err
				}
				return rows.Close()
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\"\nDEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\" union_id: \"UNIQ10\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID:UNIQ10*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (?)))\"\nDEBUG Rollback conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\" duration: 0\n",
				buf.String())
		})
	})
}
