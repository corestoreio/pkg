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
	golog "log"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/log/logw"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var reduxUniID = new(int32)

func reduxUniqueIDFunc() string {
	return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(reduxUniID, 1))
}

func setupResurrect(t testing.TB) (*dml.ConnPool, sqlmock.Sqlmock, fmt.Stringer) {

	debugLogBuf := new(log.MutexBuffer)
	l := logw.NewLog(
		logw.WithDebug(debugLogBuf, "", golog.Lshortfile),
		logw.WithLevel(logw.LevelDebug),
	)

	dbCon, dbMock := dmltest.MockDB(t, dml.WithLogger(l, reduxUniqueIDFunc))
	return dbCon, dbMock, debugLogBuf
}

// Yeah these tests are pretty bad to check the log if the program works
// correctly. Can be refactored later.

func TestReduxStmt_Query_Execution(t *testing.T) {
	// do not run parallel due to the atomic counter.
	dbCon, dbMock, logBuf := setupResurrect(t)
	defer dmltest.MockClose(t, dbCon, dbMock)

	prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT SLEEP(1)"))
	prep.ExpectQuery().WillReturnRows(&sqlmock.Rows{})
	prep.ExpectQuery().WillReturnRows(&sqlmock.Rows{})

	dmltest.FatalIfError(t, dbCon.Options(
		dml.WithPreparedStatement("sleep1", dml.QuerySQL("SELECT SLEEP(1)"), 10*time.Millisecond),
	))

	stmt, err := dbCon.Stmt("NotFound")
	assert.Nil(t, stmt, "stmt should be nil")
	require.True(t, errors.Is(err, errors.NotFound), "Should be a NotFound error, got: %+v", err)

	stmt, err = dbCon.Stmt("sleep1")
	require.NoError(t, err)
	for i := 1; i <= 2; i++ {
		rows, err := stmt.Query(context.TODO())
		require.NoError(t, err, "iteration %d", i)
		dmltest.Close(t, rows)
		require.NoError(t, rows.Err())
	}
	dmltest.Close(t, stmt)

	assert.Contains(t, logBuf.String(), "reduxStmt.rePrepare.stmt.preparing")
	assert.Contains(t, logBuf.String(), "reduxStmt.rePrepare.stmt.prepared")
	assert.Contains(t, logBuf.String(), "Query conn_pool_id")
	assert.Contains(t, logBuf.String(), `reduxStmt.closeStmtCon.stmt.close conn_pool_id: "UNIQ`)
	assert.Contains(t, logBuf.String(), `reduxStmt.closeStmtCon.con.close conn_pool_id: "UNIQ`)
	assert.Contains(t, logBuf.String(), `name: "sleep1"`)

	//t.Log("\n", logBuf.String())
}

func TestReduxStmt_Query_Resurrect(t *testing.T) {

	// do not run parallel due to the atomic counter.
	dbCon, dbMock, logBuf := setupResurrect(t)

	//dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT SLEEP(2)"))
	dbMock.
		ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT SLEEP(2)")).
		ExpectQuery().WithArgs().WillReturnRows(&sqlmock.Rows{})
	dbMock.
		ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT SLEEP(2)")).
		ExpectQuery().WithArgs().WillReturnRows(&sqlmock.Rows{})

	dmltest.FatalIfError(t, dbCon.Options(
		dml.WithPreparedStatement("sleep2", dml.QuerySQL("SELECT SLEEP(2)"), 3*time.Millisecond),
	))

	stmt, err := dbCon.Stmt("sleep2")
	require.NoError(t, err)

	rows, err := stmt.Query(context.TODO())
	require.NoError(t, err)
	require.NoError(t, rows.Err())
	dmltest.Close(t, rows)

	time.Sleep(100 * time.Millisecond) // until idleDaemon closes the statement due to inactivity

	rows, err = stmt.Query(context.TODO())
	require.NoError(t, err)
	require.NoError(t, rows.Err())
	dmltest.Close(t, rows)

	dmltest.Close(t, stmt)

	dmltest.MockClose(t, dbCon, dbMock)

	lbs := logBuf.String()
	assertContainsCount(t, lbs, "reduxStmt.rePrepare.stmt.preparing", 2)
	assertContainsCount(t, lbs, "reduxStmt.closeStmtCon.stmt.close", 2)
	assertContainsCount(t, lbs, "reduxStmt.closeStmtCon.con.close", 2)
	assertContainsCount(t, lbs, "reduxStmt.idleDaemon.stmt.closing", 1)

	time.Sleep(2 * time.Millisecond) // stupid wait until the terminated idleDaemon has written to the log file
	lbs = logBuf.String()
	assertContainsCount(t, lbs, "reduxStmt.idleDaemon.ticker.stopped", 1)

	//t.Log("\n", logBuf.String())
	// no more failures of this test: $ go test -count=217 -run=TestReduxStmt_Query_Resurrect -failfast [-race]
}

func assertContainsCount(t testing.TB, s string, contains string, occurrences int) {
	if assert.Contains(t, s, contains) {
		assert.Exactly(t, occurrences, strings.Count(s, contains), "Should contain %d times: %q", occurrences, contains)
	}
}

func TestReduxStmt_Stmt(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.Is(err, errors.NotSupported), "should have kind errors.NotSupported\n%+v", err)
			} else {
				t.Errorf("Panic should contain an error but got:\n%+v", r)
			}
		} else {
			t.Error("Expecting a panic but got nothing")
		}
	}()

	// do not run parallel due to the atomic counter.
	dbCon, dbMock, _ := setupResurrect(t)

	dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT SLEEP(3)"))

	dmltest.FatalIfError(t, dbCon.Options(
		dml.WithPreparedStatement("sleep3", dml.QuerySQL("SELECT SLEEP(3)"), 3*time.Millisecond),
	))

	stmt, err := dbCon.Stmt("sleep3")
	require.NoError(t, err)
	_ = stmt.Stmt()
}

func TestRedux_ConnPool_StmtPrepare_LoadInt64(t *testing.T) {
	t.Parallel()

	// do not run parallel due to the atomic counter.
	dbCon, dbMock, _ := setupResurrect(t)

	r := sqlmock.NewRows([]string{"SELECT_4"}).AddRow(4)
	dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT 4")).
		ExpectQuery().WithArgs().WillReturnRows(r)

	stmt, err := dbCon.StmtPrepare("select4", dml.QuerySQL("SELECT 4"), 50*time.Millisecond)
	require.NoError(t, err)
	val, err := stmt.LoadInt64(context.TODO())
	require.NoError(t, err)
	assert.Exactly(t, int64(4), val)
}

func TestRedux_ConnPool_StmtPrepare_Exec(t *testing.T) {
	t.Parallel()

	// do not run parallel due to the atomic counter.
	dbCon, dbMock, _ := setupResurrect(t)

	prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT 4"))
	prep.ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(22, 33))

	haveErr := errors.NotExists.Newf("Flat earth")
	prep.ExpectExec().WillReturnError(haveErr)

	stmt, err := dbCon.StmtPrepare("select4", dml.QuerySQL("SELECT 4"), 50*time.Millisecond)
	require.NoError(t, err)
	val, err := stmt.Exec(context.TODO())
	require.NoError(t, err)

	id, err := val.LastInsertId()
	require.NoError(t, err)
	assert.Exactly(t, int64(22), id)

	id, err = val.RowsAffected()
	require.NoError(t, err)
	assert.Exactly(t, int64(33), id)

	val, err = stmt.Exec(context.TODO())
	assert.Nil(t, val)
	assert.True(t, errors.Is(err, errors.NotExists), "Want errors.NotExists\n%+v", err)
}

func TestRedux_ConnPool_StmtPrepare_LoadInt64s(t *testing.T) {
	t.Parallel()

	// do not run parallel due to the atomic counter.
	dbCon, dbMock, _ := setupResurrect(t)

	r := sqlmock.NewRows([]string{"SELECT_4"}).AddRow(4).AddRow(5)
	prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT 4"))
	prep.ExpectQuery().WithArgs().WillReturnRows(r)

	haveErr := errors.NotExists.Newf("Flat earth")
	prep.ExpectQuery().WithArgs().WillReturnError(haveErr)

	stmt, err := dbCon.StmtPrepare("select4", dml.QuerySQL("SELECT 4"), 50*time.Millisecond)
	require.NoError(t, err)
	val, err := stmt.LoadInt64s(context.TODO())
	require.NoError(t, err)
	assert.Exactly(t, []int64{4, 5}, val)

	val, err = stmt.LoadInt64s(context.TODO())
	assert.Nil(t, val)
	assert.True(t, errors.Is(err, errors.NotExists), "Want errors.NotExists\n%+v", err)
}

func TestRedux_ConnPool_StmtPrepare_LoadInt64s_Error(t *testing.T) {
	t.Parallel()

	// do not run parallel due to the atomic counter.
	dbCon, dbMock, _ := setupResurrect(t)

	haveErr := errors.NotExists.Newf("Flat earth")
	dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT 4")).
		ExpectQuery().WithArgs().WillReturnError(haveErr)

	stmt, err := dbCon.StmtPrepare("select4", dml.QuerySQL("SELECT 4"), 50*time.Millisecond)
	require.NoError(t, err)
	val, err := stmt.LoadInt64s(context.TODO())
	assert.Nil(t, val)
	assert.True(t, errors.Is(err, errors.NotExists), "Want errors.NotExists\n%+v", err)

}
