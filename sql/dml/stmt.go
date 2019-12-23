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
	"database/sql"

	"github.com/corestoreio/errors"
)

type stmtWrapper struct {
	stmt interface {
		ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
		ioCloser
	}
}

func (sw stmtWrapper) PrepareContext(_ context.Context, sql string) (*sql.Stmt, error) {
	if sql != "" {
		panic(errors.NotAllowed.Newf("[dml] Argument `sql` with %q not allowed because this is a prepared statement", sql))
	}
	return nil, errors.NotImplemented.Newf("[dml] A *sql.Stmt cannot prepare anything")
}

func (sw stmtWrapper) ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	if sql != "" {
		panic(errors.NotAllowed.Newf("[dml] Argument `sql` with %q not allowed because this is a prepared statement", sql))
	}
	return sw.stmt.ExecContext(ctx, args...)
}

func (sw stmtWrapper) QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	if sql != "" {
		panic(errors.NotAllowed.Newf("[dml] Argument `sql` with %q not allowed because this is a prepared statement", sql))
	}
	return sw.stmt.QueryContext(ctx, args...)
}

func (sw stmtWrapper) QueryRowContext(ctx context.Context, sql string, args ...interface{}) *sql.Row {
	if sql != "" {
		panic(errors.NotAllowed.Newf("[dml] Argument `sql` with %q not allowed because this is a prepared statement", sql))
	}
	return sw.stmt.QueryRowContext(ctx, args...)
}

func (sw stmtWrapper) Close() error {
	return sw.stmt.Close()
}

// Stmt wraps a *sql.Stmt (a prepared statement) with a specific SQL query. To
// create a Stmt call the Prepare function of a specific DML type. Stmt is not
// yet safe for concurrent use, despite the underlying *sql.Stmt is. Don't
// forget to call Close!
type Stmt struct {
	base builderCommon
	Stmt *sql.Stmt
}

// WithArgs creates a new argument handler. Not safe for concurrent use.
func (st *Stmt) WithArgs() *DBR {
	var args [defaultArgumentsCapacity]argument
	a := &DBR{
		base:       st.base,
		arguments:  args[:0],
		isPrepared: true,
	}
	a.base.DB = stmtWrapper{stmt: st.Stmt}
	return a
}

// Close closes the statement in the database and frees its resources.
func (st *Stmt) Close() error { return st.Stmt.Close() }
