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
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWith_Query(t *testing.T) {
	t.Parallel()

	t.Run("error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("WITH `sel` AS (SELECT 1) SELECT * FROM `sel`")).
			WillReturnError(errors.AlreadyClosed.Newf("Who closed myself?"))

		sel := dml.NewWith(dml.WithCTE{Name: "sel", Select: dml.NewSelect().Unsafe().AddColumns("1")}).
			Select(dml.NewSelect().Star().From("sel")).
			WithDB(dbc.DB)
		rows, err := sel.WithArgs().QueryContext(context.TODO())
		assert.Nil(t, rows)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)

	})
}

func TestWith_Load(t *testing.T) {
	t.Parallel()

	t.Run("error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("WITH `sel` AS (SELECT 1) SELECT * FROM `sel`")).
			WillReturnError(errors.AlreadyClosed.Newf("Who closed myself?"))

		sel := dml.NewWith(dml.WithCTE{Name: "sel", Select: dml.NewSelect().Unsafe().AddColumns("1")}).
			Select(dml.NewSelect().Star().From("sel")).
			WithDB(dbc.DB)
		rows, err := sel.WithArgs().Load(context.TODO(), nil)
		assert.Exactly(t, uint64(0), rows)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})
}

func TestNewWith(t *testing.T) {
	t.Parallel()

	t.Run("one CTE", func(t *testing.T) {
		cte := dml.NewWith(
			dml.WithCTE{Name: "one", Select: dml.NewSelect().Unsafe().AddColumns("1")},
		).Select(dml.NewSelect().Star().From("one"))
		compareToSQL(t, cte, errors.NoKind,
			"WITH `one` AS (SELECT 1)\nSELECT * FROM `one`",
			"WITH `one` AS (SELECT 1)\nSELECT * FROM `one`",
		)
	})
	t.Run("one CTE recursive", func(t *testing.T) {
		cte := dml.NewWith(
			dml.WithCTE{
				Name:    "cte",
				Columns: []string{"n"},
				Union: dml.NewUnion(
					dml.NewSelect().Unsafe().AddColumns("1"),
					dml.NewSelect().Unsafe().AddColumns("n+1").From("cte").Where(dml.Column("n").Less().Int(5)),
				).All(),
			},
		).Recursive().Select(dml.NewSelect().Star().From("cte"))
		compareToSQL(t, cte, errors.NoKind,
			"WITH RECURSIVE `cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT n+1 FROM `cte` WHERE (`n` < 5)))\nSELECT * FROM `cte`",
			"WITH RECURSIVE `cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT n+1 FROM `cte` WHERE (`n` < 5)))\nSELECT * FROM `cte`",
		)
	})

	t.Run("two CTEs", func(t *testing.T) {
		cte := dml.NewWith(
			dml.WithCTE{Name: "intermed", Select: dml.NewSelect().Star().From("test").Where(dml.Column("x").GreaterOrEqual().Int(5))},
			dml.WithCTE{Name: "derived", Select: dml.NewSelect().Star().From("intermed").Where(dml.Column("x").Less().Int(10))},
		).Select(dml.NewSelect().Star().From("derived"))
		compareToSQL(t, cte, errors.NoKind,
			"WITH `intermed` AS (SELECT * FROM `test` WHERE (`x` >= 5)),\n`derived` AS (SELECT * FROM `intermed` WHERE (`x` < 10))\nSELECT * FROM `derived`",
			"WITH `intermed` AS (SELECT * FROM `test` WHERE (`x` >= 5)),\n`derived` AS (SELECT * FROM `intermed` WHERE (`x` < 10))\nSELECT * FROM `derived`",
		)
	})
	t.Run("multi column", func(t *testing.T) {
		cte := dml.NewWith(
			dml.WithCTE{Name: "multi", Columns: []string{"x", "y"}, Select: dml.NewSelect().Unsafe().AddColumns("1", "2")},
		).Select(dml.NewSelect("x", "y").From("multi"))
		compareToSQL(t, cte, errors.NoKind,
			"WITH `multi` (`x`,`y`) AS (SELECT 1, 2)\nSELECT `x`, `y` FROM `multi`",
			"",
		)
	})

	t.Run("DELETE", func(t *testing.T) {
		cte := dml.NewWith(
			dml.WithCTE{Name: "check_vals", Columns: []string{"val"}, Select: dml.NewSelect().Unsafe().AddColumns("123")},
		).Delete(dml.NewDelete("test").Where(dml.Column("val").In().Sub(dml.NewSelect("val").From("check_vals"))))

		compareToSQL(t, cte, errors.NoKind,
			"WITH `check_vals` (`val`) AS (SELECT 123)\nDELETE FROM `test` WHERE (`val` IN (SELECT `val` FROM `check_vals`))",
			"WITH `check_vals` (`val`) AS (SELECT 123)\nDELETE FROM `test` WHERE (`val` IN (SELECT `val` FROM `check_vals`))",
		)
	})
	t.Run("UPDATE", func(t *testing.T) {
		cte := dml.NewWith(
			dml.WithCTE{Name: "my_cte", Columns: []string{"n"}, Union: dml.NewUnion(
				dml.NewSelect().Unsafe().AddColumns("1"),
				dml.NewSelect().Unsafe().AddColumns("1+n").From("my_cte").Where(dml.Column("n").Less().Int(6)),
			).All()},
			// UPDATE statement is wrong because we're missing a JOIN which is not yet implemented.
		).Update(dml.NewUpdate("numbers").Set(dml.Column("n").Int(0)).Where(dml.Expr("n=my_cte.n*my_cte.n"))).
			Recursive()

		compareToSQL(t, cte, errors.NoKind,
			"WITH RECURSIVE `my_cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT 1+n FROM `my_cte` WHERE (`n` < 6)))\nUPDATE `numbers` SET `n`=0 WHERE (n=my_cte.n*my_cte.n)",
			"WITH RECURSIVE `my_cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT 1+n FROM `my_cte` WHERE (`n` < 6)))\nUPDATE `numbers` SET `n`=0 WHERE (n=my_cte.n*my_cte.n)",
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
		cte := dml.NewWith(
			dml.WithCTE{Name: "check_vals", Columns: []string{"val"}, Select: dml.NewSelect().AddColumns("123")},
		)
		compareToSQL(t, cte, errors.Empty,
			"",
			"",
		)
	})
}

func TestWith_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("WITH `sel` AS (SELECT 1) SELECT * FROM `sel`")).
			WillReturnError(errors.AlreadyClosed.Newf("Who closed myself?"))

		sel := dml.NewWith(dml.WithCTE{Name: "sel", Select: dml.NewSelect().Unsafe().AddColumns("1")}).
			Select(dml.NewSelect().Star().From("sel")).
			WithDB(dbc.DB)
		stmt, err := sel.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})

	t.Run("Query", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("WITH RECURSIVE `cte` (`n`) AS ((SELECT `a`, `d` AS `b` FROM `tableAD`) UNION ALL (SELECT `a`, `b` FROM `tableAB` WHERE (`b` = ?))) SELECT * FROM `cte`"))
		prep.ExpectQuery().WithArgs(6889).
			WillReturnRows(sqlmock.NewRows([]string{"a", "b"}).AddRow("Peter Gopher", "peter@gopher.go"))

		prep.ExpectQuery().WithArgs(6890).
			WillReturnRows(sqlmock.NewRows([]string{"a", "b"}).AddRow("Peter Gopher2", "peter@gopher.go2"))

		stmt, err := dml.NewWith(
			dml.WithCTE{
				Name:    "cte",
				Columns: []string{"n"},
				Union: dml.NewUnion(
					dml.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
					dml.NewSelect("a", "b").From("tableAB").Where(dml.Column("b").PlaceHolder()),
				).All(),
			},
		).
			Recursive().
			Select(dml.NewSelect().Star().From("cte")).
			WithDB(dbc.DB).
			Prepare(context.TODO())

		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		t.Run("Context", func(t *testing.T) {

			rows, err := stmt.WithArgs().QueryContext(context.TODO(), 6889)
			require.NoError(t, err)
			defer rows.Close()

			cols, err := rows.Columns()
			require.NoError(t, err)
			assert.Exactly(t, []string{"a", "b"}, cols)
		})

		t.Run("RowContext", func(t *testing.T) {

			row := stmt.WithArgs().QueryRowContext(context.TODO(), 6890)
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

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("WITH RECURSIVE `cte` (`n`) AS ((SELECT `name`, `d` AS `email` FROM `dml_person`) UNION ALL (SELECT `name`, `email` FROM `dml_person2` WHERE (`id` = ?))) SELECT * FROM `cte`"))

		stmt, err := dml.NewWith(
			dml.WithCTE{
				Name:    "cte",
				Columns: []string{"n"},
				Union: dml.NewUnion(
					dml.NewSelect("name").AddColumnsAliases("d", "email").From("dml_person"),
					dml.NewSelect("name", "email").From("dml_person2").Where(dml.Column("id").PlaceHolder()),
				).All(),
			},
		).
			Recursive().
			Select(dml.NewSelect().Star().From("cte")).
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
			// use loop with Query and add args before
			stmtA := stmt.WithArgs().Int(6899)

			for i := 0; i < iterations; i++ {
				rows, err := stmtA.QueryContext(context.TODO())
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
			stmtA := stmt.WithArgs().Record("", p)

			for i := 0; i < iterations; i++ {
				rows, err := stmtA.QueryContext(context.TODO())
				require.NoError(t, err)

				cols, err := rows.Columns()
				require.NoError(t, err)
				assert.Exactly(t, []string{"name", "email"}, cols)
				rows.Close()
			}
		})

		t.Run("WithRecords Error", func(t *testing.T) {
			p := &TableCoreConfigDataSlice{err: errors.Duplicated.Newf("Found a duplicate")}

			stmtA := stmt.WithArgs().Record("", p)
			rows, err := stmtA.QueryContext(context.TODO())
			assert.True(t, errors.Duplicated.Match(err), "%+v", err)
			assert.Nil(t, rows)
		})
	})
}

func TestWith_Clone(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t, dml.WithLogger(log.BlackHole{}, func() string { return "uniqueID" }))
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("nil", func(t *testing.T) {
		var d *dml.With
		d2 := d.Clone()
		assert.Nil(t, d)
		assert.Nil(t, d2)
	})

	t.Run("non-nil Union", func(t *testing.T) {
		cte := dml.NewWith(
			dml.WithCTE{
				Name:    "cte",
				Columns: []string{"n"},
				Union: dml.NewUnion(
					dml.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD").Where(dml.Column("b").PlaceHolder()),
					dml.NewSelect("a", "b").From("tableAB").Where(dml.Column("b").Like().NamedArg("nArg2")),
				).All(),
			},
		).
			Recursive().
			Select(dml.NewSelect().Star().From("cte").Where(dml.Column("a").GreaterOrEqual().PlaceHolder()))

		cte2 := cte.Clone()
		notEqualPointers(t, cte, cte2)
		notEqualPointers(t, cte.Subclauses, cte2.Subclauses)
		notEqualPointers(t, cte.Subclauses[0].Union, cte2.Subclauses[0].Union)
		notEqualPointers(t, cte.Subclauses[0].Columns, cte2.Subclauses[0].Columns)
		notEqualPointers(t, cte.TopLevel.Select, cte2.TopLevel.Select)
		assert.Nil(t, cte2.Subclauses[0].Select)
		assert.Nil(t, cte2.TopLevel.Update)
		assert.Nil(t, cte2.TopLevel.Delete)
		assert.Nil(t, cte2.TopLevel.Union)

		assert.Exactly(t, cte.DB, cte2.DB)
		assert.Exactly(t, cte.Log, cte2.Log)
	})
	// Add more tests for the different fields ... one day.
}
