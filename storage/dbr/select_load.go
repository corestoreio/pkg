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

// Query executes a query and returns many rows.
func (b *Select) Query(ctx context.Context) (*sql.Rows, error) {
	rows, err := Query(ctx, b.DB, b)
	return rows, errors.WithStack(err)
}

// Load loads data from a query into an object. You must set DB.QueryContext on
// the Select object or it just panics. Load can load a single row or n-rows.
func (b *Select) Load(ctx context.Context, s Scanner) (rowCount int64, err error) {
	rowCount, err = Load(ctx, b.DB, b, s)
	return rowCount, errors.WithStack(err)
}

// Prepare executes the statement represented by the Select to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Select are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement.
func (b *Select) Prepare(ctx context.Context) (*StmtSelect, error) {
	stmt, err := Prepare(ctx, b.DB, b)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cap := b.argumentCapacity()
	return &StmtSelect{
		sel:       b,
		argsCache: make(Arguments, 0, cap),
		iFaces:    make([]interface{}, 0, cap),
		stmt:      stmt,
	}, nil
}

// StmtSelect wraps a *sql.Stmt with a specific SQL query. To create a
// StmtSelect call the Prepare function of type Select. StmtSelect is not safe
// for concurrent use, despite the underlying *sql.Stmt is. Don't forget to call
// Close!
type StmtSelect struct {
	sel       *Select
	stmt      *sql.Stmt
	argsCache Arguments
	iFaces    []interface{}
	ärgErr    error // Sorry Germans for that terrible pun #notSorry
}

// Close closes the underlying prepared statement.
func (st *StmtSelect) Close() error { return st.stmt.Close() }

// WithArgs sets the arguments for the execution with Exec. It internally resets
// previously applied arguments.
func (st *StmtSelect) WithArgs(args Arguments) *StmtSelect {
	st.argsCache = st.argsCache[:0]
	st.argsCache = append(st.argsCache, args...)
	return st
}

// WithRecords sets the records for the execution with Exec. It internally
// resets previously applied arguments.
func (st *StmtSelect) WithRecords(records ...QualifiedRecord) *StmtSelect {
	st.argsCache = st.argsCache[:0]
	st.sel.BindRecord(records...)
	st.argsCache, st.ärgErr = st.sel.appendArgs(st.argsCache)
	return st
}

// Do executes a query with the previous set arguments or records or without
// arguments. It does not reset the internal arguments, so multiple executions
// with the same arguments/records are possible. Number of previously applied
// arguments or records must be the same as in the defined SQL but
// With*().Exec() can be called in a loop, both are not thread safe.
func (st *StmtSelect) Do(ctx context.Context) (*sql.Rows, error) {
	if st.ärgErr != nil {
		return nil, st.ärgErr
	}
	st.iFaces = st.iFaces[:0]
	return st.stmt.QueryContext(ctx, st.argsCache.Interfaces(st.iFaces...)...)
}

// QueryContext traditional way, allocation heavy.
func (st *StmtSelect) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	return st.stmt.QueryContext(ctx, args...)
}

// QueryRowContext traditional way, allocation heavy.
func (st *StmtSelect) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	return st.stmt.QueryRowContext(ctx, args...)
}

// Load loads data from a query into an object. You must set DB.QueryContext on
// the Select object or it just panics. Load can load a single row or n-rows.
func (st *StmtSelect) Load(ctx context.Context, s Scanner) (rowCount int64, err error) {
	r, err := st.Do(ctx)
	rowCount, err = load(r, err, s)
	return rowCount, errors.WithStack(err)
}

// LoadInt64 executes the prepared statement and returns the value at an int64.
// It returns a NotFound error if the query returns nothing.
func (st *StmtSelect) LoadInt64(ctx context.Context) (int64, error) {
	if st.sel.Log != nil && st.sel.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(st.sel.Log).Debug("dbr.Select.StmtSelect.LoadInt64", log.Stringer("sql", st.sel))
	}
	return loadInt64(st.Do(ctx))
}

// LoadInt64s executes the Select and returns the value as a slice of int64s.
func (st *StmtSelect) LoadInt64s(ctx context.Context) ([]int64, error) {
	if st.sel.Log != nil && st.sel.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(st.sel.Log).Debug("dbr.Select.StmtSelect.LoadInt64s", log.Stringer("sql", st.sel))
	}
	return loadInt64s(st.Do(ctx))
}

// The partially duplicated code in the Load[a-z0-9]+ functions can be optimized
// later. The Scanner interface should not be used for loading primitive types
// as the Scanner interface shall only be used with larger structs, means
// structs with at least two fields.

// LoadInt64 executes the Select and returns the value at an int64. It returns a
// NotFound error if the query returns nothing.
func (b *Select) LoadInt64(ctx context.Context) (int64, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadInt64", log.String("sql", sqlStr))
	}
	return loadInt64(b.DB.QueryContext(ctx, sqlStr, args...))
}

func loadInt64(rows *sql.Rows, errIn error) (value int64, err error) {
	if errIn != nil {
		return 0, errors.WithStack(errIn)
	}

	defer func() {
		if err2 := rows.Close(); err == nil && err2 != nil {
			err = errors.WithStack(err)
		}
	}()

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
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadInt64s", log.String("sql", sqlStr))
	}
	return loadInt64s(b.DB.QueryContext(ctx, sqlStr, args...))
}

func loadInt64s(rows *sql.Rows, errIn error) (_ []int64, err error) {
	if errIn != nil {
		return nil, errors.WithStack(errIn)
	}
	defer func() {
		if err2 := rows.Close(); err == nil && err2 != nil {
			err = errors.WithStack(err)
		}
	}()

	values := make([]int64, 0, 16)
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
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadUint64", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
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
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadUint64s", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
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
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadFloat64", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
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
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadFloat64s", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
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
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return "", errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadInt64", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
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
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("dbr.Select.LoadStrings", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
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
