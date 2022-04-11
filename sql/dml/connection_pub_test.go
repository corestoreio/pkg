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
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
)

func TestTableNameMapper(t *testing.T) {
	dbc, dbMock := dmltest.MockDB(t, dml.ConnPoolOption{
		TableNameMapper: func(old string) string { return fmt.Sprintf("prefix_%s", old) },
	})
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("ConnPool", func(t *testing.T) {
		t.Run("DELETE", func(t *testing.T) {
			dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("DELETE FROM `prefix_tableZ`")).WillReturnResult(sqlmock.NewResult(0, 0))
			_, err := dbc.WithQueryBuilder(dml.NewDelete("tableZ")).ExecContext(context.TODO())
			assert.NoError(t, err)
		})
		t.Run("INSERT", func(t *testing.T) {
			dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("INSERT INTO `prefix_tableZ` (`a`) VALUES (?)")).WillReturnResult(sqlmock.NewResult(0, 0))
			_, err := dbc.WithQueryBuilder(dml.NewInsert("tableZ").AddColumns("a")).ExecContext(context.TODO())
			assert.NoError(t, err)
		})
		t.Run("UPDATE", func(t *testing.T) {
			dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("UPDATE `prefix_tableZ` SET `a`=?")).WillReturnResult(sqlmock.NewResult(0, 0))
			_, err := dbc.WithQueryBuilder(dml.NewUpdate("tableZ").AddColumns("a")).ExecContext(context.TODO())
			assert.NoError(t, err)
		})
		t.Run("SELECT", func(t *testing.T) {
			dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `a` FROM `prefix_tableZ`")).WillReturnRows(sqlmock.NewRows([]string{"a"}).AddRow(1))
			_, _, err := dbc.WithQueryBuilder(dml.NewSelect("a").From("tableZ")).LoadNullInt64(context.TODO())
			assert.NoError(t, err)
		})
	})

	t.Run("Conn", func(t *testing.T) {
		con, err := dbc.Conn(context.TODO())
		assert.NoError(t, err)
		defer dmltest.Close(t, con)

		t.Run("DELETE", func(t *testing.T) {
			dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("DELETE FROM `prefix_tableZ`")).WillReturnResult(sqlmock.NewResult(0, 0))
			_, err := con.WithQueryBuilder(dml.NewDelete("tableZ")).ExecContext(context.TODO())
			assert.NoError(t, err)
		})
		t.Run("INSERT", func(t *testing.T) {
			dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("INSERT INTO `prefix_tableZ` (`a`) VALUES (?)")).WillReturnResult(sqlmock.NewResult(0, 0))
			_, err := con.WithQueryBuilder(dml.NewInsert("tableZ").AddColumns("a")).ExecContext(context.TODO())
			assert.NoError(t, err)
		})
		t.Run("UPDATE", func(t *testing.T) {
			dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("UPDATE `prefix_tableZ` SET `a`=?")).WillReturnResult(sqlmock.NewResult(0, 0))
			_, err := con.WithQueryBuilder(dml.NewUpdate("tableZ").AddColumns("a")).ExecContext(context.TODO())
			assert.NoError(t, err)
		})
		t.Run("SELECT", func(t *testing.T) {
			dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `a` FROM `prefix_tableZ`")).WillReturnRows(sqlmock.NewRows([]string{"a"}).AddRow(1))
			_, _, err := con.WithQueryBuilder(dml.NewSelect("a").From("tableZ")).LoadNullInt64(context.TODO())
			assert.NoError(t, err)
		})
	})

	t.Run("Tx", func(t *testing.T) {
		dbMock.ExpectBegin()
		tx, err := dbc.BeginTx(context.TODO(), nil)
		assert.NoError(t, err)
		defer func() { dbMock.ExpectCommit(); assert.NoError(t, tx.Commit()) }()

		t.Run("DELETE", func(t *testing.T) {
			dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("DELETE FROM `prefix_tableZ`")).WillReturnResult(sqlmock.NewResult(0, 0))
			_, err := tx.WithQueryBuilder(dml.NewDelete("tableZ")).ExecContext(context.TODO())
			assert.NoError(t, err)
		})
		t.Run("INSERT", func(t *testing.T) {
			dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("INSERT INTO `prefix_tableZ` (`a`) VALUES (?)")).WillReturnResult(sqlmock.NewResult(0, 0))
			_, err := tx.WithQueryBuilder(dml.NewInsert("tableZ").AddColumns("a")).ExecContext(context.TODO())
			assert.NoError(t, err)
		})
		t.Run("UPDATE", func(t *testing.T) {
			dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("UPDATE `prefix_tableZ` SET `a`=?")).WillReturnResult(sqlmock.NewResult(0, 0))
			_, err := tx.WithQueryBuilder(dml.NewUpdate("tableZ").AddColumns("a")).ExecContext(context.TODO())
			assert.NoError(t, err)
		})
		t.Run("SELECT", func(t *testing.T) {
			dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `a` FROM `prefix_tableZ`")).WillReturnRows(sqlmock.NewRows([]string{"a"}).AddRow(1))
			_, _, err := tx.WithQueryBuilder(dml.NewSelect("a").From("tableZ")).LoadNullInt64(context.TODO())
			assert.NoError(t, err)
		})
	})
}

func TestConnPool_Schema(t *testing.T) {
	cp, err := dml.NewConnPool(dml.WithDSN("xuser:xpassw0rd@tcp(localhost:3307)/t3st?parseTime=true&loc=UTC"))
	assert.NoError(t, err)
	assert.Exactly(t, "t3st", cp.Schema())
	dmltest.Close(t, cp)

	cp, err = dml.NewConnPool()
	assert.NoError(t, err)
	assert.Exactly(t, "", cp.Schema())
	dmltest.Close(t, cp)
}

func TestTx_Wrap(t *testing.T) {
	t.Run("commit", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectBegin()
		dbMock.ExpectExec("UPDATE `tableX` SET `value`").WithArgs().WillReturnResult(sqlmock.NewResult(0, 9))
		dbMock.ExpectCommit()

		assert.NoError(t, dbc.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
			// this creates an interpolated statement
			res, err := tx.WithQueryBuilder(
				dml.NewUpdate("tableX").AddClauses(dml.Column("value").Int(5)).Where(dml.Column("scope").Str("default")),
			).ExecContext(context.TODO())
			if err != nil {
				return err
			}
			af, err := res.RowsAffected()
			if err != nil {
				return err
			}
			assert.Exactly(t, int64(9), af)
			return nil
		}))
	})

	t.Run("rollback", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectBegin()
		dbMock.ExpectExec("UPDATE `tableX` SET `value`").WithArgs().WillReturnError(errors.Aborted.Newf("Sorry dude"))
		dbMock.ExpectRollback()

		err := dbc.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
			// Interpolated statement
			res, err := tx.WithQueryBuilder(
				dml.NewUpdate("tableX").AddClauses(dml.Column("value").Int(5)).Where(dml.Column("scope").Str("default")),
			).ExecContext(context.TODO())
			assert.Nil(t, res)
			return err
		})
		assert.Error(t, err)
	})
}

func TestWithRawSQL(t *testing.T) {
	dbc, mock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, mock)

	t.Run("ConnPool", func(t *testing.T) {
		// compareToSQL(t,
		//	dbc.WithRawSQL("SELECT * FROM users WHERE x = ? AND y IN (?,?,?)").TestWithArgs(9, 5, 6, 7),
		//	errors.NoKind,
		//	"SELECT * FROM users WHERE x = ? AND y IN (?,?,?)",
		//	"",
		//	int64(9), int64(5), int64(6), int64(7),
		//)
		//
		// compareToSQL(t,
		//	dbc.WithRawSQL("SELECT * FROM users WHERE x = 1"),
		//	errors.NoKind,
		//	"SELECT * FROM users WHERE x = 1",
		//	"",
		//)
		// compareToSQL(t,
		//	dbc.WithRawSQL("SELECT * FROM users WHERE x = ? AND y IN ?").ExpandPlaceHolders().TestWithArgs(9, []int{5, 6, 7}),
		//	errors.NoKind,
		//	"SELECT * FROM users WHERE x = ? AND y IN (?,?,?)",
		//	"",
		//	int64(9), int64(5), int64(6), int64(7),
		//)
		compareToSQL(t,
			dbc.WithQueryBuilder(dml.QuerySQL("SELECT * FROM users WHERE x = ? AND y IN ?")).
				TestWithArgs(9, []int{5, 6, 7}), // .Interpolate() gets called automatically
			errors.NoKind,
			"SELECT * FROM users WHERE x = ? AND y IN ?",
			"SELECT * FROM users WHERE x = 9 AND y IN (5,6,7)",
			int64(9), int64(5), int64(6), int64(7),
		)
		// compareToSQL(t,
		//	dbc.WithRawSQL("wat").TestWithArgs(9, 5, 6, 7),
		//	errors.NoKind,
		//	"wat",
		//	"",
		//	int64(9), int64(5), int64(6), int64(7),
		//)
	})

	t.Run("ConnSingle", func(t *testing.T) {
		c, err := dbc.Conn(context.TODO())
		defer dmltest.Close(t, c)
		if err != nil {
			t.Fatal(err)
		}
		compareToSQL(t,
			c.WithQueryBuilder(dml.QuerySQL("SELECT * FROM users WHERE x = ? AND y IN ?")).
				Interpolate().TestWithArgs(9, []int{5, 6, 7}),
			errors.NoKind,
			"SELECT * FROM users WHERE x = ? AND y IN ?",
			"SELECT * FROM users WHERE x = 9 AND y IN (5,6,7)",
			int64(9), int64(5), int64(6), int64(7),
		)
		compareToSQL(t,
			c.WithQueryBuilder(dml.QuerySQL("wat")).TestWithArgs(9, 5, 6, 7),
			errors.NoKind,
			"wat",
			"",
			int64(9), int64(5), int64(6), int64(7),
		)
	})

	t.Run("Tx", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		tx, err := dbc.BeginTx(context.TODO(), nil)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { assert.NoError(t, tx.Commit()) }()
		compareToSQL(t,
			tx.WithQueryBuilder(dml.QuerySQL("SELECT * FROM users WHERE x = ? AND y IN ?")).
				Interpolate().TestWithArgs(9, []int{5, 6, 7}),
			errors.NoKind,
			"SELECT * FROM users WHERE x = ? AND y IN ?",
			"SELECT * FROM users WHERE x = 9 AND y IN (5,6,7)",
			int64(9), int64(5), int64(6), int64(7),
		)
		compareToSQL(t,
			tx.WithQueryBuilder(dml.QuerySQL("wat")).TestWithArgs(9, 5, 6, 7),
			errors.NoKind,
			"wat",
			"",
			int64(9), int64(5), int64(6), int64(7),
		)
	})
}

func TestWithExecSQLOnConn(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.TODO()
		dbc, mock := dmltest.MockDBCallBack(t,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("create table xx3").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec("create table xx4").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()

				mock.ExpectExec("drop table xx3").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
			},
			dml.WithExecSQLOnConnClose(ctx, "drop table xx3"),
			dml.WithExecSQLOnConnOpen(ctx, "create table xx3", "create table xx4"),
		)

		dmltest.MockClose(t, dbc, mock)
	})

	t.Run("transaction rollback", func(t *testing.T) {
		ctx := context.TODO()
		errDrop := errors.NotAcceptable.Newf("Ups")

		dbc, mock := dmltest.MockDBCallBack(t,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("create table xx3").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectExec("drop table xx3").WithArgs().WillReturnError(errDrop)
			},
			dml.WithExecSQLOnConnClose(ctx, "drop table xx3"),
			dml.WithExecSQLOnConnOpen(ctx, "create table xx3"),
		)

		err := dbc.Close()
		assert.ErrorIsKind(t, errors.NotAcceptable, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestConnPool_WithPrepare(t *testing.T) {
	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	dbMock.ExpectPrepare("DROP TABLE \\?").ExpectExec().WithArgs("tabA").WillReturnResult(sqlmock.NewResult(1, 1))

	a := dbc.WithPrepare(context.TODO(), dml.QuerySQL("DROP TABLE ?"))
	_, err := a.ExecContext(context.TODO(), "tabA")
	assert.NoError(t, err)
}

func TestTx_WithPrepare(t *testing.T) {
	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	dbMock.ExpectBegin()
	dbMock.ExpectPrepare("DROP TABLE \\?").ExpectExec().WithArgs("tabA").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	err := dbc.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
		a := tx.WithPrepare(context.TODO(), dml.QuerySQL("DROP TABLE ?"))
		_, err := a.ExecContext(context.TODO(), "tabA")
		return err
	})
	assert.NoError(t, err)
}

func TestWithCreateDatabase_GivenName(t *testing.T) {
	dbc, mock := dmltest.MockDBCallBack(t,
		func(mock sqlmock.Sqlmock) {
			mock.ExpectExec("SET NAMES 'utf8mb4'").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectExec("CREATE DATABASE IF NOT EXISTS `myTestDb`").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectExec("ALTER DATABASE `myTestDb` DEFAULT CHARACTER SET='utf8mb4' COLLATE='utf8mb4_unicode_ci'").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectClose()
		},
		dml.WithCreateDatabase(context.TODO(), "myTestDb"),
	)
	assert.NoError(t, dbc.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithCreateDatabase_RandomTest(t *testing.T) {
	dbc, mock := dmltest.MockDBCallBack(t,
		func(mock sqlmock.Sqlmock) {
			mock.ExpectExec("SET NAMES 'utf8mb4'").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectExec("CREATE DATABASE IF NOT EXISTS `test_[0-9]+`").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectExec("ALTER DATABASE `test_[0-9]+` DEFAULT CHARACTER SET='utf8mb4' COLLATE='utf8mb4_unicode_ci'").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectExec("DROP DATABASE IF EXISTS `test_[0-9]+`").WithArgs().WillReturnResult(sqlmock.NewResult(0, 0))
		},
		dml.WithDSN("sqlmock:sqlmock@tcp(127.0.0.2:3306)/random?parseTime=true&loc=UTC"),
		dml.WithCreateDatabase(context.TODO(), ""),
	)
	mock.ExpectClose()
	assert.NoError(t, dbc.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestConnPool_WithDisabledForeignKeyChecks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=0`).WillReturnResult(sqlmock.NewResult(0, 0))
		dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=1`).WillReturnResult(sqlmock.NewResult(0, 0))

		err := dbc.WithDisabledForeignKeyChecks(context.TODO(), func(conn *dml.Conn) error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("error in FOREIGN_KEY_CHECKS=0", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		wantErr := errors.AlreadyClosed.Newf("Closed")
		dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=0`).WillReturnError(wantErr)

		err := dbc.WithDisabledForeignKeyChecks(context.TODO(), func(conn *dml.Conn) error {
			return errors.Blocked.Newf("gets suppressed")
		})
		assert.ErrorIsKind(t, errors.AlreadyClosed, err)
	})

	t.Run("error in FOREIGN_KEY_CHECKS=0", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		wantErr := errors.AlreadyClosed.Newf("Closed")
		dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=0`).WillReturnResult(sqlmock.NewResult(0, 0))
		dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=1`).WillReturnError(wantErr)

		err := dbc.WithDisabledForeignKeyChecks(context.TODO(), func(conn *dml.Conn) error {
			return nil
		})
		assert.ErrorIsKind(t, errors.AlreadyClosed, err)
	})

	t.Run("error in FOREIGN_KEY_CHECKS=0", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		wantErr := errors.AlreadyClosed.Newf("Closed")
		dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=0`).WillReturnResult(sqlmock.NewResult(0, 0))
		dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=1`).WillReturnError(wantErr)

		err := dbc.WithDisabledForeignKeyChecks(context.TODO(), func(conn *dml.Conn) error {
			return errors.Blocked.Newf("gets NOT suppressed")
		})
		assert.ErrorIsKind(t, errors.Blocked, err)
	})
}
