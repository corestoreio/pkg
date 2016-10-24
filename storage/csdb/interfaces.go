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

package csdb

import (
	"context"
	"database/sql"
)

// Preparer defines the only needed function to create a new statement in the
// database. The provided context is used for the preparation of the statement,
// not for the execution of the statement.
type Preparer interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// QueryRower can query one row from a database.
type QueryRower interface {
	// QueryRowContext executes a query that is expected to return at most one
	// row. QueryRowContext always returns a non-nil value. Errors are deferred
	// until Row's Scan method is called.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Querier can query multiple rows in a database.
type Querier interface {
	// QueryContext executes a query that returns rows, typically a SELECT. The
	// args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}
