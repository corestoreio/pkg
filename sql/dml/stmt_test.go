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
	"strings"
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

	dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT SLEEP(1)")).ExpectQuery().WithArgs().WillReturnRows(&sqlmock.Rows{})

	dmltest.FatalIfError(t, dbCon.Options(
		dml.WithPreparedStatement("sleep1", dml.QuerySQL("SELECT SLEEP(1)"), 10*time.Millisecond),
	))

	stmt, err := dbCon.Stmt("NotFound")
	assert.Nil(t, stmt, "stmt should be nil")
	require.True(t, errors.Is(err, errors.NotFound), "Should be a NotFound error, got: %+v", err)

	stmt, err = dbCon.Stmt("sleep1")
	require.NoError(t, err)
	rows, err := stmt.Query(context.TODO())
	require.NoError(t, err)
	dmltest.Close(t, rows)

	dmltest.Close(t, stmt)
	require.NoError(t, rows.Err())

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

	dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT SLEEP(2)"))
	dbMock.
		ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT SLEEP(2)")).
		ExpectQuery().WithArgs().WillReturnRows(&sqlmock.Rows{})

	dmltest.FatalIfError(t, dbCon.Options(
		dml.WithPreparedStatement("sleep1", dml.QuerySQL("SELECT SLEEP(2)"), 10*time.Millisecond),
	))

	time.Sleep(100 * time.Millisecond)

	stmt, err := dbCon.Stmt("sleep1")
	require.NoError(t, err)
	rows, err := stmt.Query(context.TODO())

	require.NoError(t, err)
	dmltest.Close(t, rows)

	dmltest.Close(t, stmt)

	require.NoError(t, rows.Err())

	dmltest.MockClose(t, dbCon, dbMock)

	assertContainsCount(t, logBuf, "reduxStmt.rePrepare.stmt.preparing", 2)
	assertContainsCount(t, logBuf, "reduxStmt.closeStmtCon.stmt.close", 2)
	assertContainsCount(t, logBuf, "reduxStmt.closeStmtCon.con.close", 2)
	assertContainsCount(t, logBuf, "reduxStmt.idleDaemon.stmt.closing", 1)
	assertContainsCount(t, logBuf, "reduxStmt.idleDaemon.ticker.stopped", 1)

	//assert.Contains(t, logBuf.String(), "reduxStmt.rePrepare.stmt.preparing")
	//assert.Contains(t, logBuf.String(), "reduxStmt.rePrepare.stmt.prepared")
	//assert.Contains(t, logBuf.String(), "Query conn_pool_id")
	//assert.Contains(t, logBuf.String(), `reduxStmt.closeStmtCon.stmt.close conn_pool_id: "UNIQ01" name: "sleep1"`)
	//assert.Contains(t, logBuf.String(), `reduxStmt.closeStmtCon.con.close conn_pool_id: "UNIQ01" name: "sleep1"`)

	//t.Log("\n", logBuf.String())
}

func assertContainsCount(t testing.TB, s fmt.Stringer, contains string, occurrences int) {
	if assert.Contains(t, s.String(), contains) {
		assert.Exactly(t, occurrences, strings.Count(s.String(), contains), "Should contain %d times", occurrences)
	}
}
