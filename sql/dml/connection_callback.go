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

	"github.com/corestoreio/errors"
)

// wrapDriver is used to create a new instrumented driver, it takes a vendor specific
// driver, and a call back instance to produce a new driver instance. It's usually
// used inside a sql.Register() statement
func wrapDriver(drv driver.Driver, cb DriverCallBack) driver.Driver {
	return cbDriver{drv: drv, cb: cb}
}

// DriverCallBack defines the call back signature used in every driver function.
// The returned function gets called in a defer. `fnName` states the name of the
// parent function like PrepareContext or Query, etc. The call to the first
// function can be used to e.g. start a timer. The call to second function can
// log the query and its args and also measure the time spend. The error as
// first argument in the returned function comes from the parent called function
// and should be returned or wrapped into a new one. `namedArgs` contains the,
// sometimes, named arguments. It can also be nil. context.Context can be added
// later.
type DriverCallBack func(fnName string) func(err error, query string, args []driver.NamedValue) error

// cbDriver implements a database/sql/driver.Driver
type cbDriver struct {
	drv driver.Driver
	cb  DriverCallBack
}

func (drv cbDriver) Open(name string) (driver.Conn, error) {
	conn, err := drv.drv.Open(name)
	if err != nil {
		return nil, err
	}
	fc, ok := conn.(fullConner)
	if !ok {
		return nil, errors.NotSupported.Newf("[dml] Driver does not support all required interfaces (fullConner)")
	}
	return cbConn{fc, drv.cb}, nil
}

type fullConner interface {
	driver.Conn
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
	driver.Pinger
	driver.ConnBeginTx
	// driver.ResetSessioner // later
}

type cbConn struct {
	Conn fullConner
	cb   DriverCallBack
}

func (c cbConn) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	fn := c.cb("Conn.PrepareContext")
	defer func() {
		if errFn := fn(err, query, nil); err == nil && errFn != nil {
			err = errFn
		}
	}()
	if stmt, err = c.Conn.PrepareContext(ctx, query); err != nil {
		return nil, err
	}
	if fStmt, ok := stmt.(fullStmter); ok {
		stmt = &cbStmt{Stmt: fStmt, cb: c.cb, query: query}
	} else {
		err = errors.NotSupported.Newf("[dml] Driver does not support all required interfaces (fullStmter)")
	}
	return stmt, err
}

func (c cbConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (res driver.Result, err error) {
	fn := c.cb("Conn.ExecContext")
	defer func() {
		if errFn := fn(err, query, args); err == nil && errFn != nil {
			err = errFn
		}
	}()
	res, err = c.Conn.ExecContext(ctx, query, args)
	return // do not write `return c.Conn.ExecContext` because of the defer
}

func (c cbConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rws driver.Rows, err error) {
	fn := c.cb("Conn.QueryContext")
	defer func() {
		if errFn := fn(err, query, args); err == nil && errFn != nil {
			err = errFn
		}
	}()
	rws, err = c.Conn.QueryContext(ctx, query, args)
	return // do not write `return c.Conn.QueryContext` because of the defer
}

func (c cbConn) Prepare(query string) (stmt driver.Stmt, err error) {
	fn := c.cb("Conn.Prepare")
	defer func() {
		if errFn := fn(err, query, nil); err == nil && errFn != nil {
			err = errFn
		}
	}()
	if stmt, err = c.Conn.Prepare(query); err != nil {
		return nil, err
	}

	if fStmt, ok := stmt.(fullStmter); ok {
		stmt = &cbStmt{Stmt: fStmt, cb: c.cb, query: query}
	} else {
		err = errors.NotSupported.Newf("[dml] Driver does not support all required interfaces (fullStmter)")
	}
	return stmt, err
}

func (c cbConn) Close() (err error) {
	fn := c.cb("Conn.Close")
	defer func() {
		if errFn := fn(err, "", nil); err == nil && errFn != nil {
			err = errFn
		}
	}()
	err = c.Conn.Close()
	return
}

func (c cbConn) Begin() (tx driver.Tx, err error) {
	fn := c.cb("Conn.Begin")
	defer func() {
		if errFn := fn(err, "", nil); err == nil && errFn != nil {
			err = errFn
		}
	}()
	tx, err = c.Conn.Begin()
	return
}

func (c cbConn) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	fn := c.cb("Conn.BeginTx")
	defer func() {
		if errFn := fn(err, "", nil); err == nil && errFn != nil {
			err = errFn
		}
	}()
	tx, err = c.Conn.BeginTx(ctx, opts)
	return
}

func (c cbConn) Ping(ctx context.Context) (err error) {
	fn := c.cb("Conn.Ping")
	defer func() {
		if errFn := fn(err, "", nil); err == nil && errFn != nil {
			err = errFn
		}
	}()
	err = c.Conn.Ping(ctx)
	return
}

type fullStmter interface {
	driver.Stmt
	driver.StmtExecContext
	driver.StmtQueryContext
}

// Stmt implements a database/sql/driver.Stmt
type cbStmt struct {
	Stmt  fullStmter
	cb    DriverCallBack
	query string
}

func (stmt *cbStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (res driver.Result, err error) {
	fn := stmt.cb("Stmt.ExecContext")
	defer func() {
		if errFn := fn(err, stmt.query, args); err == nil && errFn != nil {
			err = errFn
		}
	}()
	res, err = stmt.Stmt.ExecContext(ctx, args)
	return
}

func (stmt *cbStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rws driver.Rows, err error) {
	fn := stmt.cb("Stmt.QueryContext")
	defer func() {
		if errFn := fn(err, stmt.query, args); err == nil && errFn != nil {
			err = errFn
		}
	}()
	rws, err = stmt.Stmt.QueryContext(ctx, args)
	return
}

func (stmt *cbStmt) Close() (err error) {
	fn := stmt.cb("Stmt.Close")
	defer func() {
		if errFn := fn(err, stmt.query, nil); err == nil && errFn != nil {
			err = errFn
		}
	}()
	err = stmt.Stmt.Close()
	return
}

func (stmt *cbStmt) NumInput() int { return stmt.Stmt.NumInput() }

func driverValueToNamed(args []driver.Value) []driver.NamedValue {
	ret := make([]driver.NamedValue, len(args))
	for i, a := range args {
		ret[i].Ordinal = i + 1
		ret[i].Value = a
	}
	return ret
}

func (stmt *cbStmt) Exec(args []driver.Value) (res driver.Result, err error) {
	fn := stmt.cb("Stmt.Exec")
	defer func() {
		if errFn := fn(err, stmt.query, driverValueToNamed(args)); err == nil && errFn != nil {
			err = errFn
		}
	}()
	res, err = stmt.Stmt.Exec(args)
	return
}

func (stmt *cbStmt) Query(args []driver.Value) (rws driver.Rows, err error) {
	fn := stmt.cb("Stmt.Query")
	defer func() {
		if errFn := fn(err, stmt.query, driverValueToNamed(args)); err == nil && errFn != nil {
			err = errFn
		}
	}()
	rws, err = stmt.Stmt.Query(args)
	return
}
