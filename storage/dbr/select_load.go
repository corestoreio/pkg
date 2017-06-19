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
	"github.com/corestoreio/log"
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

// Load loads data from a query into `s`. Load supports up to n-rows.
func Load(ctx context.Context, b QueryBuilder, db Querier, s Scanner) (rowCount int, err error) {
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

// Query executes a query and returns many rows.
func (b *Select) Query(ctx context.Context) (*sql.Rows, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if b.Log != nil && b.Log.IsInfo() {
		// we might log sensitive data
		defer log.WhenDone(b.Log).Info("dbr.Select.Query.Timing", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args.Interfaces()...)
	return rows, errors.WithStack(err)
}

// Prepare prepares a SQL statement. Sets IsInterpolate to false.
func (b *Select) Prepare(ctx context.Context) (*sql.Stmt, error) {
	sqlStr, err := toSQLPrepared(b)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	stmt, err := b.DB.PrepareContext(ctx, sqlStr)
	return stmt, errors.WithStack(err)
}

// Load loads data from a query into an object. You must set DB.QueryContext on
// the Select object or it just panics. Load can load a single row or n-rows.
func (b *Select) Load(ctx context.Context, s Scanner) (rowCount int, err error) {
	return Load(ctx, b, b.DB, s)
}

// The partially duplicated code in the Load[a-z0-9]+ functions can be optimized
// later. The Scanner interface should not be used for loading primitive types
// as the Scanner interface shall only be used with larger structs, means
// structs with at least two fields.

// IDEA:
//func (b *Select) LoadPairInt64(ctx context.Context) (col1 []int64,col2 []int64,err error) {
//
//}

// LoadInt64 executes the Select and returns the value at an int64. It returns a
// NotFound error if the query returns nothing.
func (b *Select) LoadInt64(ctx context.Context) (int64, error) {
	sqlStr, tArg, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadInt64", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, tArg.Interfaces()...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer rows.Close() // ToDo(CyS) See Load() if returned error should be checked. hard to tests :-(

	var value int64
	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return 0, errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if !found {
		err = errors.NewNotFoundf("[dbr] LoadInt64 value not found")
	}
	return value, err
}

// LoadInt64s executes the Select and returns the value as a slice of int64s.
func (b *Select) LoadInt64s(ctx context.Context) ([]int64, error) {
	sqlStr, tArg, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadInt64s", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, tArg.Interfaces()...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	values := make([]int64, 0, 10)
	for rows.Next() {
		var value int64
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, nil
}

// LoadUint64 executes the Select and returns the value at an uint64. It returns
// a NotFound error if the query returns nothing. This function comes in handy
// when performing a COUNT(*) query. See function `Select.Count`.
func (b *Select) LoadUint64(ctx context.Context) (uint64, error) {
	sqlStr, tArg, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadUint64", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, tArg.Interfaces()...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer rows.Close()

	var value uint64
	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return 0, errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if !found {
		err = errors.NewNotFoundf("[dbr] LoadUint64 value not found")
	}
	return value, err
}

// LoadUint64s executes the Select and returns the value at a slice of uint64s.
func (b *Select) LoadUint64s(ctx context.Context) ([]uint64, error) {
	sqlStr, tArg, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadUint64s", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, tArg.Interfaces()...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	values := make([]uint64, 0, 10)
	for rows.Next() {
		var value uint64
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, nil
}

// LoadFloat64 executes the Select and returns the value at an float64. It
// returns a NotFound error if the query returns nothing.
func (b *Select) LoadFloat64(ctx context.Context) (float64, error) {
	sqlStr, tArg, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadFloat64", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, tArg.Interfaces()...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer rows.Close()

	var value float64
	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return 0, errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if !found {
		err = errors.NewNotFoundf("[dbr] LoadFloat64 value not found")
	}
	return value, err
}

// LoadFloat64s executes the Select and returns the value at a slice of float64s.
func (b *Select) LoadFloat64s(ctx context.Context) ([]float64, error) {
	sqlStr, tArg, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadFloat64s", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, tArg.Interfaces()...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	values := make([]float64, 0, 10)
	for rows.Next() {
		var value float64
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, nil
}

// LoadString executes the Select and returns the value as a string. It
// returns a NotFound error if the row amount is not equal one.
func (b *Select) LoadString(ctx context.Context) (string, error) {
	sqlStr, tArg, err := b.ToSQL()
	if err != nil {
		return "", errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadInt64", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, tArg.Interfaces()...)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer rows.Close()

	var value string
	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return "", errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return "", errors.WithStack(err)
	}
	if !found {
		err = errors.NewNotFoundf("[dbr] LoadInt64 value not found")
	}
	return value, err
}

// LoadStrings executes the Select and returns a slice of strings.
func (b *Select) LoadStrings(ctx context.Context) ([]string, error) {
	sqlStr, tArg, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadStrings", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, tArg.Interfaces()...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	values := make([]string, 0, 10)
	for rows.Next() {
		var value string
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, nil
}
