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

package dml

import (
	"context"
	"database/sql"
)

// Preparer prepares a query in the server. The underlying type can be either a
// *sql.DB (connection pool), a *sql.Conn (a single dedicated database session)
// or a *sql.Tx (an in-progress database transaction).
type Preparer interface {
	// PrepareContext - the provided context is used for the preparation of the
	// statement, not for the execution of the statement.
	// PrepareContext creates a prepared statement for later queries or
	// executions. Multiple queries or executions may be run concurrently from
	// the returned statement. The caller must call the statement's Close method
	// when the statement is no longer needed.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// Execer can execute a non-returning query. The underlying type can be either a
// *sql.DB (connection pool), a *sql.Conn (a single dedicated database session)
// or a *sql.Tx (an in-progress database transaction).
type Execer interface {
	// ExecContext executes a query that doesn't return rows.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// ExecPreparer a composite interface which can execute and prepare a query. The
// underlying type can be either a *sql.DB (connection pool), a *sql.Conn (a
// single dedicated database session) or a *sql.Tx (an in-progress database
// transaction).
type ExecPreparer interface {
	Preparer
	Execer
}

// Querier can execute a returning query. The underlying type can be either a
// *sql.DB (connection pool), a *sql.Conn (a single dedicated database session)
// or a *sql.Tx (an in-progress database transaction).
type Querier interface {
	// QueryContext executes a query that returns rows, typically a SELECT. The
	// args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// QueryPreparer can execute a returning query and prepare a returning query.
// The underlying type can be either a *sql.DB (connection pool), a *sql.Conn (a
// single dedicated database session) or a *sql.Tx (an in-progress database
// transaction).
type QueryPreparer interface {
	Preparer
	Querier
}
