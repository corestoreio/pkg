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

func TestWith_Query(t *testing.T) {
	t.Parallel()

	t.Run("error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta("WITH `sel` AS (SELECT 1) SELECT * FROM `sel`")).
			WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		sel := dbr.NewWith(dbr.WithCTE{Name: "sel", Select: dbr.NewSelect().Unsafe().AddColumns("1")}).
			Select(dbr.NewSelect().Star().From("sel")).
			WithDB(dbc.DB)
		rows, err := sel.Query(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)

	})
}

func TestWith_Load(t *testing.T) {
	t.Parallel()

	t.Run("error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta("WITH `sel` AS (SELECT 1) SELECT * FROM `sel`")).
			WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		sel := dbr.NewWith(dbr.WithCTE{Name: "sel", Select: dbr.NewSelect().Unsafe().AddColumns("1")}).
			Select(dbr.NewSelect().Star().From("sel")).
			WithDB(dbc.DB)
		rows, err := sel.Load(context.TODO(), nil)
		assert.Exactly(t, int64(0), rows)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestNewWith(t *testing.T) {
	t.Parallel()

	t.Run("one CTE", func(t *testing.T) {
		cte := dbr.NewWith(
			dbr.WithCTE{Name: "one", Select: dbr.NewSelect().Unsafe().AddColumns("1")},
		).Select(dbr.NewSelect().Star().From("one"))
		compareToSQL(t, cte, nil,
			"WITH\n`one` AS (SELECT 1)\nSELECT * FROM `one`",
			"WITH\n`one` AS (SELECT 1)\nSELECT * FROM `one`",
		)
	})
	t.Run("one CTE recursive", func(t *testing.T) {
		cte := dbr.NewWith(
			dbr.WithCTE{
				Name:    "cte",
				Columns: []string{"n"},
				Union: dbr.NewUnion(
					dbr.NewSelect().Unsafe().AddColumns("1"),
					dbr.NewSelect().Unsafe().AddColumns("n+1").From("cte").Where(dbr.Column("n").Less().Int(5)),
				).All(),
			},
		).Recursive().Select(dbr.NewSelect().Star().From("cte"))
		compareToSQL(t, cte, nil,
			"WITH RECURSIVE\n`cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT n+1 FROM `cte` WHERE (`n` < ?)))\nSELECT * FROM `cte`",
			"WITH RECURSIVE\n`cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT n+1 FROM `cte` WHERE (`n` < 5)))\nSELECT * FROM `cte`",
			int64(5),
		)
	})

	t.Run("two CTEs", func(t *testing.T) {
		cte := dbr.NewWith(
			dbr.WithCTE{Name: "intermed", Select: dbr.NewSelect().Star().From("test").Where(dbr.Column("x").GreaterOrEqual().Int(5))},
			dbr.WithCTE{Name: "derived", Select: dbr.NewSelect().Star().From("intermed").Where(dbr.Column("x").Less().Int(10))},
		).Select(dbr.NewSelect().Star().From("derived"))
		compareToSQL(t, cte, nil,
			"WITH\n`intermed` AS (SELECT * FROM `test` WHERE (`x` >= ?)),\n`derived` AS (SELECT * FROM `intermed` WHERE (`x` < ?))\nSELECT * FROM `derived`",
			"WITH\n`intermed` AS (SELECT * FROM `test` WHERE (`x` >= 5)),\n`derived` AS (SELECT * FROM `intermed` WHERE (`x` < 10))\nSELECT * FROM `derived`",
			int64(5), int64(10),
		)
	})
	t.Run("multi column", func(t *testing.T) {
		cte := dbr.NewWith(
			dbr.WithCTE{Name: "multi", Columns: []string{"x", "y"}, Select: dbr.NewSelect().Unsafe().AddColumns("1", "2")},
		).Select(dbr.NewSelect("x", "y").From("multi"))
		compareToSQL(t, cte, nil,
			"WITH\n`multi` (`x`,`y`) AS (SELECT 1, 2)\nSELECT `x`, `y` FROM `multi`",
			"",
		)
	})

	t.Run("DELETE", func(t *testing.T) {
		cte := dbr.NewWith(
			dbr.WithCTE{Name: "check_vals", Columns: []string{"val"}, Select: dbr.NewSelect().Unsafe().AddColumns("123")},
		).Delete(dbr.NewDelete("test").Where(dbr.Column("val").In().Sub(dbr.NewSelect("val").From("check_vals"))))

		compareToSQL(t, cte, nil,
			"WITH\n`check_vals` (`val`) AS (SELECT 123)\nDELETE FROM `test` WHERE (`val` IN (SELECT `val` FROM `check_vals`))",
			"WITH\n`check_vals` (`val`) AS (SELECT 123)\nDELETE FROM `test` WHERE (`val` IN (SELECT `val` FROM `check_vals`))",
		)
	})
	t.Run("UPDATE", func(t *testing.T) {
		cte := dbr.NewWith(
			dbr.WithCTE{Name: "my_cte", Columns: []string{"n"}, Union: dbr.NewUnion(
				dbr.NewSelect().Unsafe().AddColumns("1"),
				dbr.NewSelect().Unsafe().AddColumns("1+n").From("my_cte").Where(dbr.Column("n").Less().Int(6)),
			).All()},
			// UPDATE statement is wrong because we're missing a JOIN which is not yet implemented.
		).Update(dbr.NewUpdate("numbers").Set(dbr.Column("n").Int(0)).Where(dbr.Expr("n=my_cte.n*my_cte.n"))).
			Recursive()

		compareToSQL(t, cte, nil,
			"WITH RECURSIVE\n`my_cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT 1+n FROM `my_cte` WHERE (`n` < ?)))\nUPDATE `numbers` SET `n`=? WHERE (n=my_cte.n*my_cte.n)",
			"WITH RECURSIVE\n`my_cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT 1+n FROM `my_cte` WHERE (`n` < 6)))\nUPDATE `numbers` SET `n`=0 WHERE (n=my_cte.n*my_cte.n)",
			int64(6), int64(0),
		)
		//WITH RECURSIVE my_cte(n) AS
		//(
		//	SELECT 1
		//UNION ALL
		//SELECT 1+n FROM my_cte WHERE n<6
		//)
		//UPDATE numbers, my_cte
		//# Change to 0...
		//	SET numbers.n=0
		//# ... the numbers which are squares, i.e. 1 and 4
		//WHERE numbers.n=my_cte.n*my_cte.n;
	})

	t.Run("error EMPTY top clause", func(t *testing.T) {
		cte := dbr.NewWith(
			dbr.WithCTE{Name: "check_vals", Columns: []string{"val"}, Select: dbr.NewSelect().AddColumns("123")},
		)
		compareToSQL(t, cte, errors.IsEmpty,
			"",
			"",
		)
	})
}

func TestWith_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("WITH `sel` AS (SELECT 1) SELECT * FROM `sel`")).
			WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		sel := dbr.NewWith(dbr.WithCTE{Name: "sel", Select: dbr.NewSelect().Unsafe().AddColumns("1")}).
			Select(dbr.NewSelect().Star().From("sel")).
			WithDB(dbc.DB)
		stmt, err := sel.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

	t.Run("Query", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("WITH RECURSIVE `cte` (`n`) AS ((SELECT `a`, `d` AS `b` FROM `tableAD`) UNION ALL (SELECT `a`, `b` FROM `tableAB` WHERE (`b` = ?))) SELECT * FROM `cte`"))
		prep.ExpectQuery().WithArgs(6889).
			WillReturnRows(sqlmock.NewRows([]string{"a", "b"}).AddRow("Peter Gopher", "peter@gopher.go"))

		prep.ExpectQuery().WithArgs(6890).
			WillReturnRows(sqlmock.NewRows([]string{"a", "b"}).AddRow("Peter Gopher2", "peter@gopher.go2"))

		stmt, err := dbr.NewWith(
			dbr.WithCTE{
				Name:    "cte",
				Columns: []string{"n"},
				Union: dbr.NewUnion(
					dbr.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
					dbr.NewSelect("a", "b").From("tableAB").Where(dbr.Column("b").PlaceHolder()),
				).All(),
			},
		).
			Recursive().
			Select(dbr.NewSelect().Star().From("cte")).
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

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("WITH RECURSIVE `cte` (`n`) AS ((SELECT `name`, `d` AS `email` FROM `dbr_person`) UNION ALL (SELECT `name`, `email` FROM `dbr_person2` WHERE (`id` = ?))) SELECT * FROM `cte`"))

		stmt, err := dbr.NewWith(
			dbr.WithCTE{
				Name:    "cte",
				Columns: []string{"n"},
				Union: dbr.NewUnion(
					dbr.NewSelect("name").AddColumnsAliases("d", "email").From("dbr_person"),
					dbr.NewSelect("name", "email").From("dbr_person2").Where(dbr.Column("id").PlaceHolder()),
				).All(),
			},
		).
			Recursive().
			Select(dbr.NewSelect().Star().From("cte")).
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
