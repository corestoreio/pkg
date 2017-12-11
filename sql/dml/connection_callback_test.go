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

package dml

import (
	"context"
	"database/sql/driver"
	"strings"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ fullConner = (*cbConn)(nil)
var _ driver.Conn = (*cbConn)(nil)
var _ driver.ConnPrepareContext = (*cbConn)(nil)
var _ driver.ExecerContext = (*cbConn)(nil)
var _ driver.QueryerContext = (*cbConn)(nil)
var _ driver.Pinger = (*cbConn)(nil)
var _ driver.ConnBeginTx = (*cbConn)(nil)
var _ fullStmter = (*cbStmt)(nil)

//var _ driver.ResetSessioner = (*cbConn)(nil)

func TestWrapDriver_Connection_Error(t *testing.T) {
	t.Parallel()

	getCon := func(t *testing.T, errCon SqlErrDriverCon) driver.Conn {
		wrappedDrv := wrapDriver(
			SqlErrDriver{Con: errCon}, // AirCon ;-)
			func(fnName string) func(error, string, []driver.NamedValue) error {
				return func(error, string, []driver.NamedValue) error {
					return errors.NewAlreadyClosedf("Connection closed")
				}
			})
		con, err := wrappedDrv.Open("nvr mind")
		require.NoError(t, err)
		return con
	}

	t.Run("PrepareContext", func(t *testing.T) {
		con := getCon(t, SqlErrDriverCon{})
		_, err := con.(driver.ConnPrepareContext).PrepareContext(nil, "")
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("PrepareContext Original Error", func(t *testing.T) {
		con := getCon(t, SqlErrDriverCon{
			PrepareError: errors.NewWriteFailedf("Should not get overwritten"),
		})
		_, err := con.(driver.ConnPrepareContext).PrepareContext(nil, "")
		assert.True(t, errors.IsWriteFailed(err), "%s", err)
	})
	t.Run("Prepare", func(t *testing.T) {
		con := getCon(t, SqlErrDriverCon{})
		_, err := con.Prepare("")
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("Close", func(t *testing.T) {
		con := getCon(t, SqlErrDriverCon{})
		err := con.Close()
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("Begin", func(t *testing.T) {
		con := getCon(t, SqlErrDriverCon{})
		_, err := con.Begin()
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("BeginTx", func(t *testing.T) {
		con := getCon(t, SqlErrDriverCon{})
		_, err := con.(driver.ConnBeginTx).BeginTx(nil, driver.TxOptions{})
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("Ping", func(t *testing.T) {
		con := getCon(t, SqlErrDriverCon{})
		err := con.(driver.Pinger).Ping(nil)
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("ExecContext", func(t *testing.T) {
		con := getCon(t, SqlErrDriverCon{})
		_, err := con.(driver.ExecerContext).ExecContext(nil, "", nil)
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("ExecContext Original Error", func(t *testing.T) {
		con := getCon(t, SqlErrDriverCon{
			ExecError: errors.NewWriteFailedf("Should not get overwritten"),
		})
		_, err := con.(driver.ExecerContext).ExecContext(nil, "", nil)
		assert.True(t, errors.IsWriteFailed(err), "%s", err)
	})
	t.Run("QueryContext", func(t *testing.T) {
		con := getCon(t, SqlErrDriverCon{})
		_, err := con.(driver.QueryerContext).QueryContext(nil, "", nil)
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
}

func TestWrapDriver_Stmt_Error(t *testing.T) {
	t.Parallel()

	getStmt := func(t *testing.T) driver.Stmt {
		wrappedDrv := wrapDriver(SqlErrDriver{}, func(fnName string) func(error, string, []driver.NamedValue) error {
			return func(err error, _ string, _ []driver.NamedValue) error {
				if strings.HasPrefix(fnName, "Stmt.") {
					err = errors.NewAbortedf("Connection closed")
				}
				return err
			}
		})
		con, err := wrappedDrv.Open("nvr mind")
		require.NoError(t, err)
		stmt, err := con.Prepare("")
		require.NoError(t, err)
		return stmt
	}

	t.Run("Exec", func(t *testing.T) {
		con := getStmt(t)
		_, err := con.Exec(nil)
		assert.True(t, errors.IsAborted(err), "%s", err)
	})
	t.Run("Query", func(t *testing.T) {
		con := getStmt(t)
		_, err := con.Query(nil)
		assert.True(t, errors.IsAborted(err), "%s", err)
	})
	t.Run("Close", func(t *testing.T) {
		con := getStmt(t)
		err := con.Close()
		assert.True(t, errors.IsAborted(err), "%s", err)
	})

	t.Run("ExecContext", func(t *testing.T) {
		con := getStmt(t)
		_, err := con.(driver.StmtExecContext).ExecContext(nil, nil)
		assert.True(t, errors.IsAborted(err), "%s", err)
	})
	t.Run("QueryContext", func(t *testing.T) {
		con := getStmt(t)
		_, err := con.(driver.StmtQueryContext).QueryContext(nil, nil)
		assert.True(t, errors.IsAborted(err), "%s", err)
	})
}

// The next structs can be migrated to the cstesting package once needed.

type SqlErrDriver struct {
	OpenError error
	Con       SqlErrDriverCon
}

func (md SqlErrDriver) Open(name string) (driver.Conn, error) {
	return md.Con, md.OpenError
}

type SqlErrDriverCon struct {
	PrepareError  error
	ExecError     error
	QueryError    error
	PingError     error
	CloseError    error
	BeginError    error
	TxCommitErr   error
	TxRollbackErr error
	Stmt          SqlErrDriverStmt
	Tx            SqlErrDriverTx
}

func (mc SqlErrDriverCon) Prepare(query string) (driver.Stmt, error) {
	return mc.Stmt, mc.PrepareError
}

func (mc SqlErrDriverCon) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return mc.Stmt, mc.PrepareError
}
func (mc SqlErrDriverCon) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return nil, mc.ExecError
}
func (mc SqlErrDriverCon) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return nil, mc.QueryError
}
func (mc SqlErrDriverCon) Ping(ctx context.Context) (err error) {
	return mc.PingError
}

func (mc SqlErrDriverCon) Close() error              { return mc.CloseError }
func (mc SqlErrDriverCon) Begin() (driver.Tx, error) { return mc.Tx, mc.BeginError }
func (mc SqlErrDriverCon) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return mc.Tx, mc.BeginError
}

type SqlErrDriverTx struct {
	CommitErr   error
	RollbackErr error
}

func (mt SqlErrDriverTx) Commit() error   { return mt.CommitErr }
func (mt SqlErrDriverTx) Rollback() error { return mt.RollbackErr }

type SqlErrDriverStmt struct {
	CloseError error
	ExecError  error
	QueryError error
}

func (mt SqlErrDriverStmt) Close() error                                    { return mt.CloseError }
func (mt SqlErrDriverStmt) NumInput() int                                   { return 0 }
func (mt SqlErrDriverStmt) Exec(args []driver.Value) (driver.Result, error) { return nil, mt.ExecError }
func (mt SqlErrDriverStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (res driver.Result, err error) {
	return nil, mt.ExecError
}
func (mt SqlErrDriverStmt) Query(args []driver.Value) (driver.Rows, error) { return nil, mt.QueryError }
func (mt SqlErrDriverStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rws driver.Rows, err error) {
	return nil, mt.QueryError
}
