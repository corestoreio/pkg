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

	"github.com/corestoreio/errors"
)

// Scanner allows a type to load data from database query. It's used in the
// rows.Next() for-loop.
type Scanner interface {
	// RowScan implementation must use function `Scan` to scan the values of the
	// query into its own type. See database/sql package for examples.
	RowScan(*sql.Rows) error
}

// RowCloser allows to execute special functions after the scanning has
// happened. Should only be implemented in a custom type, if the interface
// Scanner has been implemented too. Not every type might need a RowCloser.
type RowCloser interface {
	// RowClose gets called at the very end and even after rows.Close. Allows
	// to implement special functions, like unlocking a mutex or updating
	// internal structures or resetting internal type containers.
	RowClose() error
}

// Exec executes the statement represented by the QueryBuilder. It returns the
// raw database/sql Result or an error if there was one. Regarding
// LastInsertID(): If you insert multiple rows using a single INSERT statement,
// LAST_INSERT_ID() returns the value generated for the first inserted row only.
// The reason for this is to make it possible to reproduce easily the same
// INSERT statement against some other server.
func Exec(ctx context.Context, db Execer, b QueryBuilder) (sql.Result, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result, err := db.ExecContext(ctx, sqlStr, args...)
	return result, errors.Wrapf(err, "[dbr] Exec.ExecContext with query %q", sqlStr)
}

// Prepare prepares a SQL statement. Sets IsInterpolate to false.
func Prepare(ctx context.Context, db Preparer, b QueryBuilder) (*sql.Stmt, error) {
	var sqlStr string
	var err error
	if qb, ok := b.(queryBuilder); ok { // Interface upgrade
		sqlStr, _, err = toSQL(qb, _isNotInterpolate, _isPrepared)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		sqlStr, _, err = b.ToSQL()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	stmt, err := db.PrepareContext(ctx, sqlStr)
	return stmt, errors.Wrapf(err, "[dbr] Prepare.PrepareContext with query %q", sqlStr)
}

// Query executes a query and returns many rows.
func Query(ctx context.Context, db Querier, b QueryBuilder) (*sql.Rows, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err := db.QueryContext(ctx, sqlStr, args...)
	return rows, errors.Wrapf(err, "[dbr] Query.QueryContext with query %q", sqlStr)
}

// Load loads data from a query into `s`. Load supports loading of up to n-rows.
// Load checks if a type implements RowCloser interface.
func Load(ctx context.Context, db Querier, b QueryBuilder, s Scanner) (rowCount int64, err error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	r, err := db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "[dbr] Load.QueryContext with query %q", sqlStr)
	}
	defer func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dbr] Load.QueryContext.Rows.Close")
		}
		if rc, ok := s.(RowCloser); ok {
			if err2 := rc.RowClose(); err2 != nil && err == nil {
				err = errors.Wrap(err2, "[dbr] Load.QueryContext.Scanner.RowClose")
			}
		}
	}()

	for r.Next() {
		err = s.RowScan(r)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		rowCount++
	}
	if err = r.Err(); err != nil {
		return rowCount, errors.WithStack(err)
	}
	return rowCount, err
}
