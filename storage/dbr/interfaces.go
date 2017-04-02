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

package dbr

import (
	"context"
	"database/sql"
)

// DBer is a composition of multiple interfaces to describe the common needed
// behaviour for querying a database. This interface is context independent.
type DBer interface {
	Preparer
	Execer
	Querier
	QueryRower
}

// Preparer creates a new prepared statement.
type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

// Querier can execute a SELECT query which can return many rows.
type Querier interface {
	// Query executes a query that returns rows, typically a SELECT. The
	// args are for any placeholder parameters in the query.
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// Execer can execute all other queries except SELECT.
type Execer interface {
	// Exec executes a query that doesn't return rows.
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// QueryRower executes a SELECT query which returns one row.
type QueryRower interface {
	// QueryRow executes a query that is expected to return at most one
	// row. QueryRow always returns a non-nil value. Errors are deferred
	// until Row'ab Scan method is called.
	QueryRow(query string, args ...interface{}) *sql.Row
}

type wrapDBContext struct {
	context.Context
	db interface {
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}
}

func (wc wrapDBContext) Prepare(query string) (*sql.Stmt, error) {
	return wc.db.PrepareContext(wc.Context, query)
}
func (wc wrapDBContext) Exec(query string, args ...interface{}) (sql.Result, error) {
	return wc.db.ExecContext(wc.Context, query, args...)
}
func (wc wrapDBContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return wc.db.QueryContext(wc.Context, query, args...)
}
func (wc wrapDBContext) QueryRow(query string, args ...interface{}) *sql.Row {
	return wc.db.QueryRowContext(wc.Context, query, args...)
}

// WrapDBContext wraps a context to be used in non context-aware functions. In
// case of a prepared statement: The provided context is used for the
// preparation of the statement, not for the execution of the statement.
func WrapDBContext(ctx context.Context, db interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}) DBer {
	return wrapDBContext{
		Context: ctx,
		db:      db,
	}
}

// WrapStmtContext wraps a context to be used in non context-aware prepared
// statement functions.
func WrapStmtContext(ctx context.Context, stmt interface {
	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
}) Stmter {
	return wrapStmtContext{
		Context: ctx,
		stmt:    stmt,
	}
}

type wrapStmtContext struct {
	context.Context
	stmt interface {
		ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
	}
}

func (wc wrapStmtContext) Exec(args ...interface{}) (sql.Result, error) {
	return wc.stmt.ExecContext(wc.Context, args...)
}
func (wc wrapStmtContext) Query(args ...interface{}) (*sql.Rows, error) {
	return wc.stmt.QueryContext(wc.Context, args...)
}
func (wc wrapStmtContext) QueryRow(args ...interface{}) *sql.Row {
	return wc.stmt.QueryRowContext(wc.Context, args...)
}

// Stmter is a composition of multiple interfaces to describe the common needed
// behaviour for querying a database within a prepared statement. This interface
// is context independent.
type Stmter interface {
	StmtExecer
	StmtQueryer
	StmtQueryRower
}

// StmtExecer executes a prepared non-SELECT statement
type StmtExecer interface {
	Exec(args ...interface{}) (sql.Result, error)
}

// StmtQueryer executes a prepared e.g. SELECT statement which can return many
// rows.
type StmtQueryer interface {
	Query(args ...interface{}) (*sql.Rows, error)
}

// StmtQueryRower executes a prepared e.g. SELECT statement which can return one
// row.
type StmtQueryRower interface {
	QueryRow(args ...interface{}) *sql.Row
}

// Txer is an in-progress database transaction.
//
// A transaction must end with a call to Commit or Rollback.
//
// After a call to Commit or Rollback, all operations on the transaction fail
// with ErrTxDone.
//
// The statements prepared for a transaction by calling the transaction'ab
// Prepare or Stmt methods are closed by the call to Commit or Rollback.
type Txer interface {
	Commit() error
	Rollback() error
	Stmt(stmt *sql.Stmt) *sql.Stmt
	DBer
}
