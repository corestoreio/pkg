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

	ctx := context.TODO()

	getCon := func(t *testing.T, errCon SQLErrDriverCon) driver.Conn {
		wrappedDrv := wrapDriver(
			SQLErrDriver{Con: errCon}, // AirCon ;-)
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
		con := getCon(t, SQLErrDriverCon{})
		_, err := con.(driver.ConnPrepareContext).PrepareContext(ctx, "")
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("PrepareContext Original Error", func(t *testing.T) {
		con := getCon(t, SQLErrDriverCon{
			PrepareError: errors.NewWriteFailedf("Should not get overwritten"),
		})
		_, err := con.(driver.ConnPrepareContext).PrepareContext(ctx, "")
		assert.True(t, errors.IsWriteFailed(err), "%s", err)
	})
	t.Run("Prepare", func(t *testing.T) {
		con := getCon(t, SQLErrDriverCon{})
		_, err := con.Prepare("")
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("Close", func(t *testing.T) {
		con := getCon(t, SQLErrDriverCon{})
		err := con.Close()
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("Begin", func(t *testing.T) {
		con := getCon(t, SQLErrDriverCon{})
		_, err := con.Begin()
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("BeginTx", func(t *testing.T) {
		con := getCon(t, SQLErrDriverCon{})
		_, err := con.(driver.ConnBeginTx).BeginTx(ctx, driver.TxOptions{})
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("Ping", func(t *testing.T) {
		con := getCon(t, SQLErrDriverCon{})
		err := con.(driver.Pinger).Ping(ctx)
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("ExecContext", func(t *testing.T) {
		con := getCon(t, SQLErrDriverCon{})
		_, err := con.(driver.ExecerContext).ExecContext(ctx, "", nil)
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
	t.Run("ExecContext Original Error", func(t *testing.T) {
		con := getCon(t, SQLErrDriverCon{
			ExecError: errors.NewWriteFailedf("Should not get overwritten"),
		})
		_, err := con.(driver.ExecerContext).ExecContext(ctx, "", nil)
		assert.True(t, errors.IsWriteFailed(err), "%s", err)
	})
	t.Run("QueryContext", func(t *testing.T) {
		con := getCon(t, SQLErrDriverCon{})
		_, err := con.(driver.QueryerContext).QueryContext(ctx, "", nil)
		assert.True(t, errors.IsAlreadyClosed(err), "%s", err)
	})
}

func TestWrapDriver_Stmt_Error(t *testing.T) {
	t.Parallel()

	getStmt := func(t *testing.T) driver.Stmt {
		wrappedDrv := wrapDriver(SQLErrDriver{}, func(fnName string) func(error, string, []driver.NamedValue) error {
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

	ctx := context.TODO()

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
		_, err := con.(driver.StmtExecContext).ExecContext(ctx, nil)
		assert.True(t, errors.IsAborted(err), "%s", err)
	})
	t.Run("QueryContext", func(t *testing.T) {
		con := getStmt(t)
		_, err := con.(driver.StmtQueryContext).QueryContext(ctx, nil)
		assert.True(t, errors.IsAborted(err), "%s", err)
	})
}

// The next structs can be migrated to the cstesting package once needed.

type SQLErrDriver struct {
	OpenError error
	Con       SQLErrDriverCon
}

func (md SQLErrDriver) Open(name string) (driver.Conn, error) {
	return md.Con, md.OpenError
}

type SQLErrDriverCon struct {
	PrepareError  error
	ExecError     error
	QueryError    error
	PingError     error
	CloseError    error
	BeginError    error
	TxCommitErr   error
	TxRollbackErr error
	Stmt          SQLErrDriverStmt
	Tx            SQLErrDriverTx
}

func (mc SQLErrDriverCon) Prepare(query string) (driver.Stmt, error) {
	return mc.Stmt, mc.PrepareError
}

func (mc SQLErrDriverCon) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return mc.Stmt, mc.PrepareError
}
func (mc SQLErrDriverCon) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return nil, mc.ExecError
}
func (mc SQLErrDriverCon) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return nil, mc.QueryError
}
func (mc SQLErrDriverCon) Ping(ctx context.Context) (err error) {
	return mc.PingError
}

func (mc SQLErrDriverCon) Close() error              { return mc.CloseError }
func (mc SQLErrDriverCon) Begin() (driver.Tx, error) { return mc.Tx, mc.BeginError }
func (mc SQLErrDriverCon) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return mc.Tx, mc.BeginError
}

type SQLErrDriverTx struct {
	CommitErr   error
	RollbackErr error
}

func (mt SQLErrDriverTx) Commit() error   { return mt.CommitErr }
func (mt SQLErrDriverTx) Rollback() error { return mt.RollbackErr }

type SQLErrDriverStmt struct {
	CloseError error
	ExecError  error
	QueryError error
}

func (mt SQLErrDriverStmt) Close() error                                    { return mt.CloseError }
func (mt SQLErrDriverStmt) NumInput() int                                   { return 0 }
func (mt SQLErrDriverStmt) Exec(args []driver.Value) (driver.Result, error) { return nil, mt.ExecError }
func (mt SQLErrDriverStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (res driver.Result, err error) {
	return nil, mt.ExecError
}
func (mt SQLErrDriverStmt) Query(args []driver.Value) (driver.Rows, error) { return nil, mt.QueryError }
func (mt SQLErrDriverStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rws driver.Rows, err error) {
	return nil, mt.QueryError
}
