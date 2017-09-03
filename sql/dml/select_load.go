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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Query executes a query and returns many rows.
func (b *Select) Query(ctx context.Context) (*sql.Rows, error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Query", log.Stringer("sql", b))
	}
	rows, err := Query(ctx, b.DB, b)
	return rows, errors.WithStack(err)
}

// Load loads data from a query into an object. You must set DB.QueryContext on
// the Select object or it just panics. Load can load a single row or n-rows.
func (b *Select) Load(ctx context.Context, s Scanner) (rowCount int64, err error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Load", log.Stringer("sql", b))
	}
	rowCount, err = Load(ctx, b.DB, b, s)
	return rowCount, errors.WithStack(err)
}

// Prepare executes the statement represented by the Select to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Select are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement.
func (b *Select) Prepare(ctx context.Context) (*StmtSelect, error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Prepare", log.Stringer("sql", b))
	}
	stmt, err := Prepare(ctx, b.DB, b)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	l := b.Log
	if l != nil {
		l = b.Log.With(log.Bool("is_prepared", true))
	}
	cap := b.argumentCapacity()
	return &StmtSelect{
		StmtBase: StmtBase{
			id:         b.id,
			stmt:       stmt,
			argsCache:  make(Arguments, 0, cap),
			argsRaw:    make([]interface{}, 0, cap),
			bindRecord: b.bindRecord,
			log:        l,
		},
		sel: b,
	}, nil
}

// StmtSelect wraps a *sql.Stmt with a specific SQL query. To create a
// StmtSelect call the Prepare function of type Select. StmtSelect is not safe
// for concurrent use, despite the underlying *sql.Stmt is. Don't forget to call
// Close!
type StmtSelect struct {
	StmtBase
	sel *Select
}

// WithArgs sets the interfaced arguments for the execution with Query+. It
// internally resets previously applied arguments.
func (st *StmtSelect) WithArgs(args ...interface{}) *StmtSelect {
	st.withArgs(args)
	return st
}

// WithArguments sets the arguments for the execution with Query+. It internally
// resets previously applied arguments.
func (st *StmtSelect) WithArguments(args Arguments) *StmtSelect {
	st.withArguments(args)
	return st
}

// WithRecords sets the records for the execution with Query+. It internally
// resets previously applied arguments.
func (st *StmtSelect) WithRecords(records ...QualifiedRecord) *StmtSelect {
	st.withRecords(st.sel.appendArgs, records...)
	return st
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
		defer log.WhenDone(b.Log).Debug("LoadInt64", log.String("sql", sqlStr))
	}
	return loadInt64(b.DB.QueryContext(ctx, sqlStr, args...))
}

func loadInt64(rows *sql.Rows, errIn error) (value int64, err error) {
	if errIn != nil {
		return 0, errors.WithStack(errIn)
	}

	defer func() {
		if cErr := rows.Close(); err == nil && cErr != nil {
			err = errors.WithStack(cErr)
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
		err = errors.NewNotFoundf("[dml] LoadInt64 value not found")
	}
	return value, err
}

// LoadInt64s executes the Select and returns the value as a slice of int64s.
func (b *Select) LoadInt64s(ctx context.Context) (ret []int64, err error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("LoadInt64s", log.Int("row_count", len(ret)), log.String("sql", sqlStr))
	}
	ret, err = loadInt64s(b.DB.QueryContext(ctx, sqlStr, args...))
	return ret, err // used in defer!
}

func loadInt64s(rows *sql.Rows, errIn error) (_ []int64, err error) {
	if errIn != nil {
		return nil, errors.WithStack(errIn)
	}
	defer func() {
		if cErr := rows.Close(); err == nil && cErr != nil {
			err = errors.WithStack(cErr)
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
func (b *Select) LoadUint64(ctx context.Context) (_ uint64, err error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("LoadUint64", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

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
		err = errors.NewNotFoundf("[dml] LoadUint64 value not found")
	}
	return value, err
}

// LoadUint64s executes the Select and returns the value at a slice of uint64s.
func (b *Select) LoadUint64s(ctx context.Context) (values []uint64, err error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("LoadUint64s", log.Int("row_count", len(values)), log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	values = make([]uint64, 0, 10)
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
func (b *Select) LoadFloat64(ctx context.Context) (_ float64, err error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("LoadFloat64", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

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
		err = errors.NewNotFoundf("[dml] LoadFloat64 value not found")
	}
	return value, err
}

// LoadFloat64s executes the Select and returns the value at a slice of float64s.
func (b *Select) LoadFloat64s(ctx context.Context) (_ []float64, err error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("LoadFloat64s", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

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
	return values, err
}

// LoadString executes the Select and returns the value as a string. It
// returns a NotFound error if the row amount is not equal one.
func (b *Select) LoadString(ctx context.Context) (_ string, err error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return "", errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("LoadString", log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

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
		err = errors.NewNotFoundf("[dml] LoadInt64 value not found")
	}
	return value, err
}

// LoadStrings executes the Select and returns a slice of strings.
func (b *Select) LoadStrings(ctx context.Context) (values []string, err error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Log != nil && b.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(b.Log).Debug("LoadStrings", log.Int("row_count", len(values)), log.String("sql", sqlStr))
	}

	rows, err := b.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	values = make([]string, 0, 10)
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
	return values, err
}
