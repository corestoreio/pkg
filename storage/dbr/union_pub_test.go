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
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnion_Query(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL Error", func(t *testing.T) {
		u := dbr.NewUnion(
			dbr.NewSelect(),
			dbr.NewSelect(),
		)
		rows, err := u.Query(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.IsEmpty(err))
	})

	u := dbr.NewUnion(
		dbr.NewSelect("value").From("eavChar"),
		dbr.NewSelect("value").From("eavInt").Where(dbr.Column("b").Float64(3.14159)),
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

	u := dbr.NewUnion(
		dbr.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
		dbr.NewSelect("a", "b").From("tableAB").Where(dbr.Column("b").Float64(3.14159)),
	).Unsafe().
		OrderBy("a").OrderByDesc("b").OrderBy(`concat("c",b,"d")`)

	t.Run("error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta("(SELECT `a`, `d` AS `b` FROM `tableAD`) UNION (SELECT `a`, `b` FROM `tableAB` WHERE (`b` = ?)) ORDER BY `a`, `b` DESC, concat(\"c\",b,\"d\")")).
			WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		rows, err := u.WithDB(dbc.DB).Load(context.TODO(), nil)
		assert.Exactly(t, int64(0), rows)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestUnion_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL Error", func(t *testing.T) {
		u := dbr.NewUnion(
			dbr.NewSelect(),
			dbr.NewSelect(),
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

		u := dbr.NewUnion(
			dbr.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
			dbr.NewSelect("a", "b").From("tableAB").Where(dbr.Column("b").Float64(3.14159)),
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

		stmt, err := dbr.NewUnion(
			dbr.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
			dbr.NewSelect("a", "b").From("tableAB").Where(dbr.Column("b").PlaceHolder()),
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

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("(SELECT `name`, `d` AS `email` FROM `dbr_person`) UNION (SELECT `name`, `email` FROM `dbr_person2` WHERE (`id` = ?))"))

		stmt, err := dbr.NewUnion(
			dbr.NewSelect("name").AddColumnsAliases("d", "email").From("dbr_person"),
			dbr.NewSelect("name", "email").From("dbr_person2").Where(dbr.Column("id").PlaceHolder()),
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
			// use loop with Do and add args before
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

}
