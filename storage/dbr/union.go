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
	"strconv"
	"strings"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Union represents a UNION SQL statement. UNION is used to combine the result
// from multiple SELECT statements into a single result set.
// With template usage enabled, it builds multiple select statements joined by
// UNION and all based on a common template.
type Union struct {
	// UseBuildCache if `true` the final build query including place holders
	// will be cached in a private field. Each time a call to function ToSQL
	// happens, the arguments will be re-evaluated and returned or interpolated.
	UseBuildCache bool
	cacheSQL      []byte
	cacheArgs     Arguments // like a buffer, gets reused

	Selects       []*Select
	OrderBys      aliases
	IsAll         bool // IsAll enables UNION ALL
	IsInterpolate bool // See Interpolate()

	// When using Union as a template, only one *Select is required.
	oldNew        [][]string //use for string replacement with `repls` field
	repls         []*strings.Replacer
	stmtCount     int
	previousError error
}

// NewUnion creates a new Union object. If using as a template, only one *Select
// object can be provided.
func NewUnion(selects ...*Select) *Union {
	return &Union{
		Selects: selects,
	}
}

// Append adds more *Select objects to the Union object. When using Union as a
// template only one *Select object can be provided.
func (u *Union) Append(selects ...*Select) *Union {
	u.Selects = append(u.Selects, selects...)
	return u
}

// All returns all rows. The default behavior for UNION is that duplicate rows
// are removed from the result. Enabling ALL returns all rows.
func (u *Union) All() *Union {
	u.IsAll = true
	return u
}

// PreserveResultSet enables the correct ordering of the result set from the
// Select statements. UNION by default produces an unordered set of rows. To
// cause rows in a UNION result to consist of the sets of rows retrieved by each
// SELECT one after the other, select an additional column in each SELECT to use
// as a sort column and add an ORDER BY following the last SELECT.
func (u *Union) PreserveResultSet() *Union {
	if len(u.Selects) > 1 {
		for i, s := range u.Selects {
			s.AddColumnsExprAlias(strconv.Itoa(i), "_preserve_result_set")
		}
		u.OrderBys = append(aliases{MakeAlias("_preserve_result_set")}, u.OrderBys...)
		return u
	}

	// Panics without any *Select in the slice. Programmer error.
	u.Selects[0].AddColumnsExprAlias("{preserveResultSet}", "_preserve_result_set")
	u.OrderBys = append(aliases{MakeAlias("_preserve_result_set")}, u.OrderBys...)
	for i := 0; i < u.stmtCount; i++ {
		u.oldNew[i] = append(u.oldNew[i], "{preserveResultSet}", strconv.Itoa(i))
	}
	return u
}

// OrderBy appends a column to ORDER the statement ascending. Columns are
// getting quoted. MySQL might order the result set in a temporary table, which
// is slow. Under different conditions sorting can skip the temporary table.
// https://dev.mysql.com/doc/relnotes/mysql/5.7/en/news-5-7-3.html
func (u *Union) OrderBy(columns ...string) *Union {
	u.OrderBys = u.OrderBys.appendColumns(columns, false).applySort(len(columns), sortAscending)
	return u
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// Columns are getting quoted. When you use ORDER BY or GROUP BY to sort a
// column in a DELETE, the server sorts values using only the initial number of
// bytes indicated by the max_sort_length system variable.
func (u *Union) OrderByDesc(columns ...string) *Union {
	u.OrderBys = u.OrderBys.appendColumns(columns, false).applySort(len(columns), sortDescending)
	return u
}

// OrderByExpr adds a custom SQL expression to the ORDER BY clause. Does not
// quote the strings.
func (u *Union) OrderByExpr(columns ...string) *Union {
	u.OrderBys = u.OrderBys.appendColumns(columns, true)
	return u
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `Arguments` will then be nil.
func (u *Union) Interpolate() *Union {
	u.IsInterpolate = true
	return u
}

// StringReplace is only applicable when using *Union as a template.
// StringReplace replaces the `key` with one of the `values`. Each value defines
// a generated SELECT query. Repeating calls of StringReplace must provide the
// same amount of `values` as the first call. This function is just a simple
// string replacement. Make sure that your key does not match other parts of the
// SQL query.
func (u *Union) StringReplace(key string, values ...string) *Union {
	if len(u.Selects) > 1 {
		return u
	}
	if u.stmtCount == 0 {
		u.stmtCount = len(values)
		u.oldNew = make([][]string, u.stmtCount)
		u.repls = make([]*strings.Replacer, u.stmtCount)
	}
	if len(values) != u.stmtCount {
		u.previousError = errors.NewNotValidf("[dbr] Union.StringReplace: Argument count for values too short. Have %d Want %d", len(values), u.stmtCount)
		return u
	}
	for i := 0; i < u.stmtCount; i++ {
		u.oldNew[i] = append(u.oldNew[i], key, values[i])
	}
	return u
}

// MultiplyArguments is only applicable when using *Union as a template.
// MultiplyArguments repeats the `args` variable n-times to match the number of
// generated SELECT queries in the final UNION statement. It should be called
// after all calls to `StringReplace` have been made.
func (u *Union) MultiplyArguments(args ...Argument) Arguments {
	if len(u.Selects) > 1 {
		return args
	}
	ret := make(Arguments, len(args)*u.stmtCount)
	lArgs := len(args)
	for i := 0; i < u.stmtCount; i++ {
		copy(ret[i*lArgs:], args)
	}
	return ret
}

// ToSQL converts the statements into a string and returns its arguments.
func (u *Union) ToSQL() (string, Arguments, error) {
	return toSQL(u, u.IsInterpolate)
}

func (u *Union) writeBuildCache(sql []byte) {
	u.cacheSQL = sql
}

func (u *Union) readBuildCache() (sql []byte, _ Arguments, err error) {
	if u.cacheSQL == nil {
		return nil, nil, nil
	}
	u.cacheArgs, err = u.appendArgs(u.cacheArgs[:0])
	return u.cacheSQL, u.cacheArgs, err
}

func (u *Union) hasBuildCache() bool {
	return u.UseBuildCache
}

// ToSQL generates the SQL string and its arguments. Calls to this function are
// idempotent.
func (u *Union) toSQL(w queryWriter) error {
	if u.previousError != nil {
		return u.previousError
	}
	if len(u.Selects) > 1 {
		for i, s := range u.Selects {
			if i > 0 {
				sqlWriteUnionAll(w, u.IsAll)
			}
			w.WriteByte('(')

			if err := s.toSQL(w); err != nil {
				return errors.Wrapf(err, "[dbr] Union.ToSQL at Select index %d", i)
			}
			w.WriteByte(')')
		}
		sqlWriteOrderBy(w, u.OrderBys, true)
		return nil
	}

	bufS1 := bufferpool.Get()
	err := u.Selects[0].toSQL(bufS1)
	selStr := bufS1.String()
	bufferpool.Put(bufS1)
	if err != nil {
		return errors.Wrap(err, "[dbr] Union.ToSQL: toSQL template")
	}

	for i := 0; i < u.stmtCount; i++ {
		repl := u.repls[i]
		if repl == nil {
			repl = strings.NewReplacer(u.oldNew[i]...)
			u.repls[i] = repl
		}
		if i > 0 {
			sqlWriteUnionAll(w, u.IsAll)
		}
		w.WriteByte('(')
		repl.WriteString(w, selStr)
		w.WriteByte(')')
	}
	sqlWriteOrderBy(w, u.OrderBys, true)
	return nil
}

func (u *Union) makeArguments() Arguments {
	var argCap int
	for _, s := range u.Selects {
		argCap += s.argumentCapacity()
	}
	return make(Arguments, 0, len(u.Selects)*argCap)
}

func (u *Union) appendArgs(args Arguments) (_ Arguments, err error) {
	if u.previousError != nil {
		return nil, u.previousError
	}
	if cap(args) == 0 {
		args = u.makeArguments()
	}
	if len(u.Selects) > 1 {
		for i, s := range u.Selects {
			args, err = s.appendArgs(args)
			if err != nil {
				return nil, errors.Wrapf(err, "[dbr] Union.ToSQL at Select index %d", i)
			}
		}
		return args, nil
	}
	args, err = u.Selects[0].appendArgs(args)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Union.ToSQL: toSQL template")
	}
	return u.MultiplyArguments(args...), nil
}

// Exec executes the statement represented by the Union. It returns the raw
// database/sql Result or an error if there was one. It expects the *sql.DB
// object on the first []*Select index.
func (u *Union) Query(ctx context.Context) (*sql.Rows, error) {
	sqlStr, args, err := u.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Union.Exec.ToSQL")
	}

	s1 := u.Selects[0]
	if s1.Log != nil && s1.Log.IsInfo() {
		defer log.WhenDone(s1.Log).Info("dbr.Union.Exec.Timing", log.String("sql", sqlStr))
	}

	rows, err := s1.DB.QueryContext(ctx, sqlStr, args.Interfaces()...)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] delete.exec.Exec")
	}

	return rows, nil
}

// Prepare executes the statement represented by the Union. It returns the raw
// database/sql Statement and an error if there was one. Provided arguments in
// the Union are getting ignored. It panics when field Preparer at nil. It
// expects the *sql.DB object on the first []*Select index.
func (u *Union) Prepare(ctx context.Context) (*sql.Stmt, error) {
	sqlStr, err := toSQLPrepared(u)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Union.Prepare.toSQLPrepared")
	}

	s1 := u.Selects[0]
	if s1.Log != nil && s1.Log.IsInfo() {
		defer log.WhenDone(s1.Log).Info("dbr.Union.Prepare.Timing", log.String("sql", sqlStr))
	}

	stmt, err := s1.DB.PrepareContext(ctx, sqlStr)
	return stmt, errors.Wrap(err, "[dbr] Union.Prepare.Prepare")
}
