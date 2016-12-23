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

package dbr_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ dbr.DBer = (*sql.DB)(nil)
var _ dbr.Preparer = (*sql.DB)(nil)
var _ dbr.Querier = (*sql.DB)(nil)
var _ dbr.Execer = (*sql.DB)(nil)
var _ dbr.QueryRower = (*sql.DB)(nil)

var _ dbr.Stmter = (*sql.Stmt)(nil)
var _ dbr.StmtQueryRower = (*sql.Stmt)(nil)
var _ dbr.StmtQueryer = (*sql.Stmt)(nil)
var _ dbr.StmtExecer = (*sql.Stmt)(nil)

var _ dbr.Txer = (*sql.Tx)(nil)

var _ dbr.Preparer = (*dbMock)(nil)
var _ dbr.Querier = (*dbMock)(nil)
var _ dbr.Execer = (*dbMock)(nil)

type dbMock struct {
	error
	prepareFn func(query string) (*sql.Stmt, error)
}

func (pm dbMock) Prepare(query string) (*sql.Stmt, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return pm.prepareFn(query)
}

func (pm dbMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return nil, nil
}

func (pm dbMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return nil, nil
}

func TestWrapDBContext(t *testing.T) {

	dbConn, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbConn.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	dbCTX := dbr.WrapDBContext(context.TODO(), dbConn.DB)

	dbMock.ExpectPrepare("INSERT INTO TABLE abc").WillReturnError(errors.New("Upssss"))
	stmt, err := dbCTX.Prepare("INSERT INTO TABLE abc")
	assert.Nil(t, stmt)
	assert.EqualError(t, err, "Upssss", "%+v", err)

	dbMock.ExpectQuery("SELECT a FROM tableX").WithArgs(1).WillReturnError(errors.New("Upssss"))
	rows, err := dbCTX.Query("SELECT a FROM tableX where b = ?", 1)
	assert.Nil(t, rows)
	assert.EqualError(t, err, "Upssss", "%+v", err)

	dbMock.ExpectQuery("SELECT a FROM tableY").WithArgs(1).WillReturnError(errors.New("Upssss"))
	row := dbCTX.QueryRow("SELECT a FROM tableY where b = ?", 1)
	err = row.Scan()
	assert.EqualError(t, err, "Upssss", "%+v", err)

	dbMock.ExpectExec("ALTER TABLE add").WithArgs(1).WillReturnError(errors.New("Upssss"))
	res, err := dbCTX.Exec("ALTER TABLE add a = ?", 1)
	assert.Nil(t, res)
	assert.EqualError(t, err, "Upssss", "%+v", err)
}

func TestWrapStmtContext(t *testing.T) {

	dbConn, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbConn.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	dbMock.ExpectPrepare("INSERT INTO TABLE abc").ExpectExec().WillReturnError(errors.New("Upssss Exec"))
	dbMock.ExpectPrepare("SELECT a FROM tableA WHERE a = ").ExpectQuery().WithArgs(123).WillReturnError(errors.New("Upssss Query"))
	dbMock.ExpectPrepare("SELECT b FROM tableB WHERE b = ").ExpectQuery().WithArgs(456).WillReturnError(errors.New("Upssss QueryRow"))

	stmt, err := dbConn.DB.Prepare("INSERT INTO TABLE abc")
	require.NoError(t, err, "%+v", err)
	stmtCTX := dbr.WrapStmtContext(context.TODO(), stmt)
	res, err := stmtCTX.Exec()
	assert.Nil(t, res)
	assert.EqualError(t, err, "Upssss Exec", "%+v", err)

	stmt, err = dbConn.DB.Prepare("SELECT a FROM tableA WHERE a = ?")
	require.NoError(t, err, "%+v", err)
	stmtCTX = dbr.WrapStmtContext(context.TODO(), stmt)
	rows, err := stmtCTX.Query(123)
	assert.Nil(t, rows)
	assert.EqualError(t, err, "Upssss Query", "%+v", err)

	stmt, err = dbConn.DB.Prepare("SELECT b FROM tableB WHERE b = ?")
	require.NoError(t, err, "%+v", err)
	stmtCTX = dbr.WrapStmtContext(context.TODO(), stmt)
	row := stmtCTX.QueryRow(456)
	err = row.Scan()
	assert.EqualError(t, err, "Upssss QueryRow", "%+v", err)
}
