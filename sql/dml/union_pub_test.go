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
	"github.com/corestoreio/pkg/util/assert"
)

func TestUnion_Query(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL Error", func(t *testing.T) {
		u := dml.NewUnion(
			dml.NewSelect(),
			dml.NewSelect(),
		)
		rows, err := u.WithDBR().QueryContext(context.TODO())
		assert.Nil(t, rows)
		assert.ErrorIsKind(t, errors.Empty, err)
	})

	u := dml.NewUnion(
		dml.NewSelect("value").From("eavChar"),
		dml.NewSelect("value").From("eavInt").Where(dml.Column("b").Float64(3.14159)),
	)

	t.Run("Error", func(t *testing.T) {
		u.WithDB(dbMock{
			error: errors.ConnectionFailed.Newf("Who closed myself?"),
		})
		rows, err := u.WithDBR().QueryContext(context.TODO())
		assert.Nil(t, rows)
		assert.ErrorIsKind(t, errors.ConnectionFailed, err)
	})

	t.Run("Success", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		smr := sqlmock.NewRows([]string{"value"}).AddRow("row1").AddRow("row2")
		dbMock.ExpectQuery(
			dmltest.SQLMockQuoteMeta("(SELECT `value` FROM `eavChar`) UNION (SELECT `value` FROM `eavInt` WHERE (`b` = 3.14159))"),
		).WillReturnRows(smr)

		u.WithDB(dbc.DB)

		rows, err := u.WithDBR().QueryContext(context.TODO())
		assert.NoError(t, err)

		var xx []string
		for rows.Next() {
			var x string
			assert.NoError(t, rows.Scan(&x))
			xx = append(xx, x)
		}
		assert.Exactly(t, []string{"row1", "row2"}, xx)
		assert.NoError(t, rows.Close())
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
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta(
			"(SELECT `a`, `d` AS `b` FROM `tableAD`) UNION (SELECT `a`, `b` FROM `tableAB` WHERE (`b` = 3.14159)) ORDER BY `a`, `b` DESC, concat(\"c\",b,\"d\")")).
			WillReturnError(errors.AlreadyClosed.Newf("Who closed myself?"))

		rows, err := u.WithDB(dbc.DB).WithDBR().Load(context.TODO(), nil)
		assert.Exactly(t, uint64(0), rows)
		assert.ErrorIsKind(t, errors.AlreadyClosed, err)
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
		assert.ErrorIsKind(t, errors.Empty, err)
	})

	t.Run("Error", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectPrepare(
			dmltest.SQLMockQuoteMeta(
				"(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`) UNION (SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = 3.14159)) ORDER BY `_preserve_result_set`, `a`, `b` DESC, concat(\"c\",b,\"d\")"),
		).
			WillReturnError(errors.AlreadyClosed.Newf("Who closed myself?"))

		u := dml.NewUnion(
			dml.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
			dml.NewSelect("a", "b").From("tableAB").Where(dml.Column("b").Float64(3.14159)),
		).
			Unsafe().
			OrderBy("a").OrderByDesc("b").OrderBy(`concat("c",b,"d")`).
			PreserveResultSet().WithDB(dbc.DB)

		stmt, err := u.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.ErrorIsKind(t, errors.AlreadyClosed, err)
	})

	t.Run("Query", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("(SELECT `a`, `d` AS `b` FROM `tableAD`) UNION (SELECT `a`, `b` FROM `tableAB` WHERE (`b` = ?))"))
		prep.ExpectQuery().WithArgs(6889).
			WillReturnRows(sqlmock.NewRows([]string{"a", "b"}).AddRow("Peter Gopher", "peter@gopher.go"))

		prep.ExpectQuery().WithArgs(6890).
			WillReturnRows(sqlmock.NewRows([]string{"a", "b"}).AddRow("Peter Gopher2", "peter@gopher.go2"))

		stmt, err := dml.NewUnion(
			dml.NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
			dml.NewSelect("a", "b").From("tableAB").Where(dml.Column("b").PlaceHolder()),
		).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		assert.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			assert.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		t.Run("Context", func(t *testing.T) {
			rows, err := stmt.WithDBR().QueryContext(context.TODO(), 6889)
			assert.NoError(t, err)
			defer rows.Close()

			cols, err := rows.Columns()
			assert.NoError(t, err)
			assert.Exactly(t, []string{"a", "b"}, cols)
		})

		t.Run("RowContext", func(t *testing.T) {
			row := stmt.WithDBR().QueryRowContext(context.TODO(), 6890)
			assert.NoError(t, err)
			n, e := "", ""
			assert.NoError(t, row.Scan(&n, &e))

			assert.Exactly(t, "Peter Gopher2", n)
			assert.Exactly(t, "peter@gopher.go2", e)
		})
	})

	t.Run("Exec", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("(SELECT `name`, `d` AS `email` FROM `dml_people`) UNION (SELECT `name`, `email` FROM `dml_people2` WHERE (`id` = ?))"))

		stmt, err := dml.NewUnion(
			dml.NewSelect("name").AddColumnsAliases("d", "email").From("dml_people"),
			dml.NewSelect("name", "email").From("dml_people2").Where(dml.Column("id").PlaceHolder()),
		).
			WithDB(dbc.DB).
			Prepare(context.TODO())

		assert.NoError(t, err, "failed creating a prepared statement")
		defer dmltest.Close(t, stmt)

		const iterations = 3

		t.Run("WithArguments", func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				prep.ExpectQuery().WithArgs(6899).
					WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher", "peter@gopher.go"))
			}
			// use loop with Query and add args before
			stmtA := stmt.WithDBR()

			for i := 0; i < iterations; i++ {
				rows, err := stmtA.QueryContext(context.TODO(), 6899)
				assert.NoError(t, err)

				cols, err := rows.Columns()
				assert.NoError(t, err)
				assert.Exactly(t, []string{"name", "email"}, cols)
				rows.Close()
			}
		})

		t.Run("WithRecords", func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				prep.ExpectQuery().
					WillReturnRows(sqlmock.NewRows([]string{"name", "email"}).AddRow("Peter Gopher2", "peter@gopher.go2"))
			}

			p := &dmlPerson{ID: 6900}
			stmtA := stmt.WithDBR()

			for i := 0; i < iterations; i++ {
				rows, err := stmtA.QueryContext(context.TODO(), dml.Qualify("", p))
				assert.NoError(t, err)

				cols, err := rows.Columns()
				assert.NoError(t, err)
				assert.Exactly(t, []string{"name", "email"}, cols)
				rows.Close()
			}
		})

		t.Run("WithRecords Error", func(t *testing.T) {
			p := &TableCoreConfigDataSlice{err: errors.Duplicated.Newf("Found a duplicate")}
			rows, err := stmt.WithDBR().QueryContext(context.TODO(), dml.Qualify("", p))
			assert.ErrorIsKind(t, errors.Duplicated, err)
			assert.Nil(t, rows)
		})
	})
}

func TestUnion_Clone(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t, dml.WithLogger(log.BlackHole{}, func() string { return "uniqueID" }))
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("nil", func(t *testing.T) {
		var d *dml.Union
		d2 := d.Clone()
		assert.Nil(t, d)
		assert.Nil(t, d2)
	})

	t.Run("non-nil", func(t *testing.T) {
		u := dml.NewUnion(
			dml.NewSelect("a", "b").From("tableAD").Where(dml.Column("a").Like().PlaceHolder()),
			dml.NewSelect("a", "b").From("tableAB").Where(dml.Column("c").Between().PlaceHolder()),
		).All().OrderBy("a").OrderByDesc("b").PreserveResultSet()

		u2 := u.Clone()
		notEqualPointers(t, u, u2)
		notEqualPointers(t, u, u2)
		notEqualPointers(t, u.Selects, u2.Selects)
		notEqualPointers(t, u.OrderBys, u2.OrderBys)
		notEqualPointers(t, u.Selects[0], u2.Selects[0])
		notEqualPointers(t, u.Selects[1], u2.Selects[1])
		assert.Exactly(t, u.DB, u2.DB)
		assert.Exactly(t, u.Log, u2.Log)
	})
}
