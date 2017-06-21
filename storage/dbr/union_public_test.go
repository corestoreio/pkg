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
		dbr.NewSelect("value").From("eavInt").Where(dbr.Column("b", dbr.Equal.Float64(3.14159))),
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
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()
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

	u := dbr.NewUnion(
		dbr.NewSelect("a").AddColumnsAlias("d", "b").From("tableAD"),
		dbr.NewSelect("a", "b").From("tableAB").Where(dbr.Column("b", dbr.Equal.Float64(3.14159))),
	).
		OrderBy("a").OrderByDesc("b").OrderByExpr(`concat("c",b,"d")`).
		PreserveResultSet()
	u.UseBuildCache = true

	t.Run("Error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer func() {
			dbMock.ExpectClose()
			require.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()
		dbMock.ExpectPrepare(
			cstesting.SQLMockQuoteMeta("(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`) UNION (SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = ?)) ORDER BY `_preserve_result_set`, `a` ASC, `b` DESC, concat(\"c\",b,\"d\")"),
		).
			WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		u.WithDB(dbc.DB)

		stmt, err := u.Prepare(context.TODO())
		require.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

	t.Run("Prepared", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer func() {
			dbMock.ExpectClose()
			require.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()
		dbMock.ExpectPrepare(
			cstesting.SQLMockQuoteMeta("(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`) UNION (SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = ?)) ORDER BY `_preserve_result_set`, `a` ASC, `b` DESC, concat(\"c\",b,\"d\")"),
		)

		u.WithDB(dbc.DB)

		stmt, err := u.Prepare(context.TODO())
		require.NotNil(t, stmt)
		assert.NoError(t, err)
	})
}

func TestUnion_Load(t *testing.T) {
	t.Parallel()

	u := dbr.NewUnion(
		dbr.NewSelect("a").AddColumnsAlias("d", "b").From("tableAD"),
		dbr.NewSelect("a", "b").From("tableAB").Where(dbr.Column("b", dbr.Equal.Float64(3.14159))),
	).
		OrderBy("a").OrderByDesc("b").OrderByExpr(`concat("c",b,"d")`)

	t.Run("error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()
		dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta("(SELECT `a`, `d` AS `b` FROM `tableAD`) UNION (SELECT `a`, `b` FROM `tableAB` WHERE (`b` = ?)) ORDER BY `a` ASC, `b` DESC, concat(\"c\",b,\"d\")")).
			WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		rows, err := u.WithDB(dbc.DB).Load(context.TODO(), nil)
		assert.Exactly(t, int64(0), rows)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}
