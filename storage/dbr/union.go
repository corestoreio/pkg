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
	Log log.Logger // Log optional logger
	// DB gets required once the Load*() functions will be used.
	DB QueryPreparer

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
	IsIntersect   bool // See Intersect()
	IsExcept      bool // See Except()

	// When using Union as a template, only one *Select is required.
	oldNew    [][]string //use for string replacement with `repls` field
	repls     []*strings.Replacer
	stmtCount int
}

// NewUnion creates a new Union object. If using as a template, only one *Select
// object can be provided.
func NewUnion(selects ...*Select) *Union {
	return &Union{
		Selects: selects,
	}
}

// Union creates a new Union which selects from the provided columns.
// Columns won't get quoted.
func (c *Connection) Union(selects ...*Select) *Union {
	return &Union{
		Log:     c.Log,
		Selects: selects,
		DB:      c.DB,
	}
}

// Union creates a new Union that select that given columns bound to the transaction
func (tx *Tx) Union(selects ...*Select) *Union {
	return &Union{
		Log:     tx.Logger,
		Selects: selects,
		DB:      tx.Tx,
	}
}

// WithDB sets the database query object.
func (b *Union) WithDB(db QueryPreparer) *Union {
	b.DB = db
	return b
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
		u.OrderBys = append(aliases{MakeNameAlias("_preserve_result_set", "")}, u.OrderBys...)
		return u
	}

	// Panics without any *Select in the slice. Programmer error.
	u.Selects[0].AddColumnsExprAlias("{preserveResultSet}", "_preserve_result_set")
	u.OrderBys = append(aliases{MakeNameAlias("_preserve_result_set", "")}, u.OrderBys...)
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

// Intersect switches the query type from UNION to INTERSECT. The result of an
// intersect is the intersection of right and left SELECT results, i.e. only
// records that are present in both result sets will be included in the result
// of the operation. INTERSECT has higher precedence than UNION and EXCEPT. If
// possible it will be executed linearly but if not it will be translated to a
// subquery in the FROM clause. Only supported in MariaDB >=10.3
func (u *Union) Intersect() *Union {
	u.IsIntersect = true
	return u
}

// Except switches the query from UNION to EXCEPT. The result of EXCEPT is all
// records of the left SELECT result except records which are in right SELECT
// result set, i.e. it is subtraction of two result sets. EXCEPT and UNION have
// the same operation precedence. Only supported in MariaDB >=10.3
func (u *Union) Except() *Union {
	u.IsExcept = true
	return u
}

// StringReplace is only applicable when using *Union as a template.
// StringReplace replaces the `key` with one of the `values`. Each value defines
// a generated SELECT query. Repeating calls of StringReplace must provide the
// same amount of `values` as the first  or an index of bound stack trace
// happens. This function is just a simple string replacement. Make sure that
// your key does not match other parts of the SQL query.
func (u *Union) StringReplace(key string, values ...string) *Union {
	if len(u.Selects) > 1 {
		return u
	}
	if u.stmtCount == 0 {
		u.stmtCount = len(values)
		u.oldNew = make([][]string, u.stmtCount)
		u.repls = make([]*strings.Replacer, u.stmtCount)
	}
	for i := 0; i < u.stmtCount; i++ {
		// The following block has been put on each line because the (index out of
		// bound) stack trace will show exactly what you have made wrong =>
		// Providing in the 2nd call of StringReplace too few `values`
		// arguments.
		u.oldNew[i] = append(u.oldNew[i], key,
			values[i])
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
func (u *Union) ToSQL() (string, []interface{}, error) {
	return toSQL(u, u.IsInterpolate, _isNotPrepared)
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

	if len(u.Selects) > 1 {
		for i, s := range u.Selects {
			if i > 0 {
				sqlWriteUnionAll(w, u.IsAll, u.IsIntersect, u.IsExcept)
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
		return errors.WithStack(err)
	}

	for i := 0; i < u.stmtCount; i++ {
		repl := u.repls[i]
		if repl == nil {
			repl = strings.NewReplacer(u.oldNew[i]...)
			u.repls[i] = repl
		}
		if i > 0 {
			sqlWriteUnionAll(w, u.IsAll, u.IsIntersect, u.IsExcept)
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
		return nil, errors.WithStack(err)
	}
	return u.MultiplyArguments(args...), nil
}

// Query executes a query and returns many rows.
func (b *Union) Query(ctx context.Context) (*sql.Rows, error) {
	rows, err := Query(ctx, b.DB, b)
	return rows, errors.WithStack(err)
}

// Prepare prepares a SQL statement. Sets IsInterpolate to false.
func (b *Union) Prepare(ctx context.Context) (*sql.Stmt, error) {
	stmt, err := Prepare(ctx, b.DB, b)
	return stmt, errors.WithStack(err)
}

// Load loads data from a query into an object. You must set DB.QueryContext on
// the Union object or it just panics. Load can load a single row or n-rows.
func (b *Union) Load(ctx context.Context, s Scanner) (rowCount int64, err error) {
	rowCount, err = Load(ctx, b.DB, b, s)
	return rowCount, errors.WithStack(err)
}
