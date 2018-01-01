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
)

func setupResurrect(t testing.TB) {

}

func TestResurrectStmt_Query_Execution(t *testing.T) {
	db, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, db, dbMock)

	dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT SLEEP(1)")).ExpectQuery().WithArgs().WillReturnRows(&sqlmock.Rows{})

	debugLogBuf := new(log.MutexBuffer)

	l := logw.NewLog(
		logw.WithDebug(debugLogBuf, "testDebug: ", golog.Lshortfile),
		logw.WithLevel(logw.LevelDebug),
	)

	uniID := new(int32)
	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 1))
	}

	if err := db.Options(
		dml.WithPreparedStatement("sleep1", dml.QuerySQL("SELECT SLEEP(1)"), 10*time.Millisecond),
		dml.WithLogger(l, uniqueIDFunc),
	); err != nil {
		t.Fatal(err)
	}

	stmt, err := db.Stmt("NotFound")
	assert.Nil(t, stmt, "stmt should be nil")
	require.True(t, errors.Is(err, errors.NotFound), "Should be a NotFound error, got: %+v", err)

	stmt, err = db.Stmt("sleep1")
	require.NoError(t, err)
	rows, err := stmt.Query(context.TODO())
	require.NoError(t, err)

	dmltest.Close(t, rows)

	require.NoError(t, rows.Err())

	t.Log("\n", debugLogBuf.String())
}
