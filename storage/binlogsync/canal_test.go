// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package binlogsync_test

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/binlogsync"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCanal_Option_With_DB_Error(t *testing.T) {
	t.Run("MySQL Ping", func(t *testing.T) {
		dsn := &mysql.Config{
			User:   "root",
			Passwd: "",
			Net:    "x'",
			Addr:   "tcp127",
			DBName: "",
		}
		c, err := binlogsync.NewCanal(dsn, binlogsync.WithMySQL())
		assert.Nil(t, c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), `unknown network x'`)
	})

	t.Run("DB Ping", func(t *testing.T) {
		dsn := &mysql.Config{
			User:   "root",
			Passwd: "",
			Net:    "x'",
			Addr:   "tcp127",
			DBName: "",
		}
		db, err := sql.Open("mysql", "root:root@localhost/db")
		c, err := binlogsync.NewCanal(dsn, binlogsync.WithDB(db))
		assert.Nil(t, c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), `default addr for network 'localhost' unknown`)
	})
}

func TestNewCanal_FailedMasterStatus(t *testing.T) {
	dsn := &mysql.Config{
		User:   "root",
		Passwd: "",
		Net:    "x'err",
		Addr:   "localhost:3306",
		DBName: "TestDB",
	}
	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	wantErr := errors.NewAlreadyClosedf("I'm already closed")
	dbMock.ExpectQuery(`SHOW MASTER STATUS`).
		WillReturnError(wantErr)

	c, err := binlogsync.NewCanal(dsn, binlogsync.WithDB(dbc.DB))
	assert.Nil(t, c)
	assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
}

func TestNewCanal_CheckBinlogRowFormat_Wrong(t *testing.T) {
	dsn := &mysql.Config{
		User:   "root",
		Passwd: "",
		Net:    "x'err",
		Addr:   "localhost:3306",
		DBName: "TestDB",
	}
	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	dbMock.ExpectQuery(`SHOW MASTER STATUS`).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"File", "Position", "Binlog_Do_DB", "Binlog_Ignore_DB", "Executed_Gtid_Set"}).
				FromCSVString(`mysqlbin.log:0001,4711,,,`),
		)
	dbMock.ExpectQuery(`SHOW SESSION VARIABLES LIKE`).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"Variable_Name", "Value"}).
				FromCSVString(`binlog_format,a cat`),
		)

	c, err := binlogsync.NewCanal(dsn, binlogsync.WithDB(dbc.DB))
	assert.Nil(t, c)
	assert.True(t, errors.IsNotSupported(err), "%+v", err)
	assert.Contains(t, err.Error(), `a cat`)
}

func TestNewCanal_CheckBinlogRowFormat_Error(t *testing.T) {
	dsn := &mysql.Config{
		User:   "root",
		Passwd: "",
		Net:    "x'err",
		Addr:   "localhost:3306",
		DBName: "TestDB",
	}
	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	dbMock.ExpectQuery(`SHOW MASTER STATUS`).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"File", "Position", "Binlog_Do_DB", "Binlog_Ignore_DB", "Executed_Gtid_Set"}).
				FromCSVString(`mysqlbin.log:0001,4711,,,`),
		)
	wantErr := errors.NewNotImplementedf("MySQL Syntax not implemted")
	dbMock.ExpectQuery(`SHOW SESSION VARIABLES LIKE`).
		WillReturnError(wantErr)

	c, err := binlogsync.NewCanal(dsn, binlogsync.WithDB(dbc.DB))
	assert.Nil(t, c)
	assert.True(t, errors.IsNotImplemented(err), "%+v", err)
	assert.Contains(t, err.Error(), `MySQL Syntax not implemted`)
}

func newTestCanal(t *testing.T) (*binlogsync.Canal, sqlmock.Sqlmock, func()) {
	dsn := &mysql.Config{
		User:   "root",
		Passwd: "",
		Net:    "x'err",
		Addr:   "localhost:3306",
		DBName: "TestDB",
	}
	dbc, dbMock := cstesting.MockDB(t)
	deferred := func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}

	dbMock.ExpectQuery(`SHOW MASTER STATUS`).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"File", "Position", "Binlog_Do_DB", "Binlog_Ignore_DB", "Executed_Gtid_Set"}).
				FromCSVString(`mysqlbin.log:0001,4711,,,`),
		)
	dbMock.ExpectQuery(`SHOW SESSION VARIABLES LIKE`).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"Variable_Name", "Value"}).
				FromCSVString(`binlog_format,row`),
		)

	c, err := binlogsync.NewCanal(dsn, binlogsync.WithDB(dbc.DB))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	return c, dbMock, deferred
}

func TestNewCanal_SuccessfulStart(t *testing.T) {
	c, _, deferred := newTestCanal(t)
	defer deferred()

	cp := c.SyncedPosition()
	assert.Exactly(t, `mysqlbin.log:0001`, cp.File)
	assert.Exactly(t, uint(4711), cp.Position)
}

func TestCanal_FindTable(t *testing.T) {
	c, dbMock, deferred := newTestCanal(t)
	defer deferred()

	dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta(csdb.DMLLoadColumns)).
		WithArgs("core_config_data").
		WillReturnRows(
			cstesting.MustMockRows(cstesting.WithFile("testdata/core_config_data_columns.csv")))

	// food for the race detector and to test that only one query will be executed.
	const iterations = 10
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()

			if i < 4 {
				time.Sleep(time.Microsecond * time.Duration(i*10))
			}

			tbl, err := c.FindTable(context.Background(), 31, "core_config_data")
			if err != nil {
				t.Fatalf("%+v", err)
			}
			assert.Exactly(t,
				c.DSN.DBName,
				tbl.Schema,
			)
			assert.Exactly(t, `core_config_data`, tbl.Name)
			assert.Exactly(t, []string{"config_id", "scope", "scope_id", "path", "value"}, tbl.Columns.FieldNames())

		}(&wg, i)
	}
	wg.Wait()

}

func TestCanal_FindTable_Error(t *testing.T) {
	c, dbMock, deferred := newTestCanal(t)
	defer deferred()

	wantErr := errors.NewUnauthorizedf("Du kommst da nicht rein")
	dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta(csdb.DMLLoadColumns)).
		WithArgs("core_config_data").
		WillReturnError(wantErr)

	tbl, err := c.FindTable(context.Background(), 31, "core_config_data")
	assert.Exactly(t, csdb.Table{}, tbl)
	assert.True(t, errors.IsUnauthorized(err), "%+v", err)
}
