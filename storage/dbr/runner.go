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
	// ScanRow implementation must use function `scan` to scan the values of the
	// query into its own type. See database/sql package for examples. `idx`
	// defines the current iteration number. `columns` specifies the list of
	// provided column names used in the query. This function signature shows
	// its strength in creating slices of values or iterating over a result set,
	// modifying values and saving it back somewhere.
	ScanRow(idx int, columns []string, scan func(dest ...interface{}) error) error
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
	result, err := db.ExecContext(ctx, sqlStr, args.Interfaces()...)
	return result, errors.WithStack(err)
}

// Prepare prepares a SQL statement. Sets IsInterpolate to false.
func Prepare(ctx context.Context, db Preparer, b QueryBuilder) (*sql.Stmt, error) {
	var sqlStr string
	var err error
	if qb, ok := b.(queryBuilder); ok { // Interface upgrade
		sqlStr, _, err = toSQL(qb, isNotInterpolate, isPrepared)
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
	return stmt, errors.WithStack(err)
}

// Query executes a query and returns many rows.
func Query(ctx context.Context, db Querier, b QueryBuilder) (*sql.Rows, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err := db.QueryContext(ctx, sqlStr, args.Interfaces()...)
	return rows, errors.WithStack(err)
}

// Load loads data from a query into `s`. Load supports up to n-rows.
func Load(ctx context.Context, db Querier, b QueryBuilder, s Scanner) (rowCount int, err error) {
	sqlStr, tArg, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	rows, err := db.QueryContext(ctx, sqlStr, tArg.Interfaces()...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer func() {
		// Not testable with the sqlmock package :-(
		if err2 := rows.Close(); err2 != nil && err == nil {
			err = errors.WithStack(err2)
		}
	}()

	columns, err := rows.Columns()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	for rows.Next() {
		err = s.ScanRow(rowCount, columns, rows.Scan)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		rowCount++
	}
	if err = rows.Err(); err != nil {
		return rowCount, errors.WithStack(err)
	}
	return rowCount, err
}
