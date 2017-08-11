// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// These interfaces are private on purpose. No need to export them.

// preparer prepares a query.
type preparer interface {
	// PrepareContext - the provided context is used for the preparation of the
	// statement, not for the execution of the statement.
	// PrepareContext creates a prepared statement for later queries or
	// executions. Multiple queries or executions may be run concurrently from
	// the returned statement. The caller must call the statement's Close method
	// when the statement is no longer needed.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// execer can execute a non-returning query.
type execer interface {
	// ExecContext executes a query that doesn't return rows.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// execPreparer a composite interface which can execute and prepare a query.
type execPreparer interface {
	preparer
	execer
}

// querier can execute a returning query.
type querier interface {
	// QueryContext executes a query that returns rows, typically a SELECT. The
	// args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// queryPreparer can execute a returning query and prepare a returning query.
type queryPreparer interface {
	preparer
	querier
}

// txer is an in-progress database transaction.
//
// A transaction must end with a call to Commit or Rollback.
//
// After a call to Commit or Rollback, all operations on the transaction fail
// with ErrTxDone.
//
// The statements prepared for a transaction by calling the transaction'ab
// Prepare or Stmt methods are closed by the call to Commit or Rollback.
type txer interface {
	Commit() error
	Rollback() error
	Stmt(stmt *sql.Stmt) *sql.Stmt
	execer
	preparer
	querier
}

var _ txer = (*txMock)(nil)

// txMock does nothing and returns always nil
type txMock struct{}

func (txMock) Commit() error                                                       { return nil }
func (txMock) Rollback() error                                                     { return nil }
func (txMock) Stmt(stmt *sql.Stmt) *sql.Stmt                                       { return nil }
func (txMock) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) { return nil, nil }
func (txMock) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}
func (txMock) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}
