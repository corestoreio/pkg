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

package mycanal_test

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/sql/mycanal"
	"github.com/corestoreio/pkg/util/assert"
)

func TestCanal_Option_With_DB_Error(t *testing.T) {
	t.Run("MySQL Ping", func(t *testing.T) {
		c, err := mycanal.NewCanal(
			`root:@x'(tcp127)/?allowNativePasswords=false&parseTime=true&maxAllowedPacket=0`,
			mycanal.WithMySQL(),
			&mycanal.Options{})
		assert.Nil(t, c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), `unknown network x'`)
	})

	t.Run("DB Ping", func(t *testing.T) {
		db, err := sql.Open("mysql", "root:root@localhost/db")
		c, err := mycanal.NewCanal(
			`root:@x'(tcp127)/?allowNativePasswords=false&maxAllowedPacket=0`,
			mycanal.WithDB(db), &mycanal.Options{})
		assert.Nil(t, c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), `default addr for network 'localhost' unknown`)
	})
}

func TestNewCanal_FailedMasterStatus(t *testing.T) {

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	wantErr := errors.AlreadyClosed.Newf("I'm already closed")
	dbMock.ExpectQuery(`SHOW MASTER STATUS`).
		WillReturnError(wantErr)

	c, err := mycanal.NewCanal(
		`root:@x'err(localhost:3306)/TestDB?allowNativePasswords=false&maxAllowedPacket=0`,
		mycanal.WithDB(dbc.DB), &mycanal.Options{})
	assert.Nil(t, c)
	assert.True(t, errors.Is(err, errors.AlreadyClosed), "%+v", err)
}

func TestNewCanal_CheckBinlogRowFormat_Wrong(t *testing.T) {

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	dbMock.ExpectQuery(`SHOW MASTER STATUS`).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"File", "Position", "Binlog_Do_DB", "Binlog_Ignore_DB", "Executed_Gtid_Set"}).
				FromCSVString(`mysqlbin.log:0001,4711,,,`),
		)
	dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SHOW VARIABLES WHERE (`Variable_name` LIKE 'binlog_format')")).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"Variable_name", "Value"}).
				FromCSVString(`binlog_format,a cat`),
		)

	c, err := mycanal.NewCanal(
		`root:@x'err(localhost:3306)/TestDB?allowNativePasswords=false&maxAllowedPacket=0`,
		mycanal.WithDB(dbc.DB), &mycanal.Options{})
	assert.Nil(t, c)
	assert.True(t, errors.Is(err, errors.NotSupported), "%+v", err)
	assert.Contains(t, err.Error(), `a cat`)
}

func TestNewCanal_CheckBinlogRowFormat_Error(t *testing.T) {

	dbc, dbMock := dmltest.MockDB(t)
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
	wantErr := errors.NotImplemented.Newf("MySQL Syntax not implemted")
	dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SHOW VARIABLES WHERE (`Variable_name` LIKE 'binlog_format')")).
		WillReturnError(wantErr)

	c, err := mycanal.NewCanal(
		`root:@x'err(localhost:3306)/TestDB?allowNativePasswords=false&maxAllowedPacket=0`,
		mycanal.WithDB(dbc.DB), &mycanal.Options{})
	assert.Nil(t, c)
	assert.True(t, errors.Is(err, errors.NotImplemented), "%+v", err)
	assert.Contains(t, err.Error(), `MySQL Syntax not implemted`)
}

func newTestCanal(t *testing.T) (*mycanal.Canal, sqlmock.Sqlmock, func()) {

	dbc, dbMock := dmltest.MockDB(t)

	deferred := func() {
		dmltest.MockClose(t, dbc, dbMock)
	}

	dbMock.ExpectQuery(`SHOW MASTER STATUS`).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"File", "Position", "Binlog_Do_DB", "Binlog_Ignore_DB", "Executed_Gtid_Set"}).
				FromCSVString(`mysqlbin.log:0001,4711,,,`),
		)
	dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SHOW VARIABLES WHERE (`Variable_name` LIKE 'binlog_format')")).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"Variable_name", "Value"}).
				FromCSVString(`binlog_format,row`),
		)

	c, err := mycanal.NewCanal(`root:@x'err(localhost:3306)/TestDB?allowNativePasswords=false&maxAllowedPacket=0`,
		mycanal.WithDB(dbc.DB), &mycanal.Options{})
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

func TestCanal_FindTable_RaceFree(t *testing.T) {
	c, dbMock, deferred := newTestCanal(t)
	defer deferred()

	dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE.+TABLE_NAME IN \\('core_config_data'\\)").
		WithArgs().
		WillReturnRows(
			dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_columns.csv")))

	// food for the race detector and to test that only one query will be executed.
	const iterations = 10
	var wg sync.WaitGroup
	wg.Add(iterations)
	wantColumns := []string{"config_id", "scope", "scope_id", "path", "value"}
	for i := 0; i < iterations; i++ {
		go func(wg *sync.WaitGroup, i int) {
			// in those goroutines we can't use *testing.T because it causes
			// race conditions.
			defer wg.Done()

			if i < 4 {
				time.Sleep(time.Microsecond * time.Duration(i*10))
			}
			tbl, err := c.FindTable(context.Background(), "core_config_data")
			if err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			if want, have := `core_config_data`, tbl.Name; have != want {
				panic(fmt.Sprintf("have %q want %q", have, want))
			}
			if !reflect.DeepEqual(wantColumns, tbl.Columns.FieldNames()) {
				panic(fmt.Sprintf("have %q want %q", wantColumns, tbl.Columns.FieldNames()))
			}
		}(&wg, i)
	}
	wg.Wait()
}

func TestCanal_FindTable_Error(t *testing.T) {
	c, dbMock, deferred := newTestCanal(t)
	defer deferred()

	wantErr := errors.Unauthorized.Newf("Du kommst da nicht rein")

	dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE.. AND TABLE_NAME IN \\('core_config_data'\\)").
		WithArgs().
		WillReturnError(wantErr)

	tbl, err := c.FindTable(context.Background(), "core_config_data")
	assert.True(t, errors.Unauthorized.Match(err), "%+v", err)
	assert.Nil(t, tbl)
}
