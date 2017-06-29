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

// WithCTE defines a common table expression used in the type `With`.
type WithCTE struct {
	Name string
	// Columns, optionally, the number of names in the list must be the same as
	// the number of columns in the result set.
	Columns []string
	// Select clause as a common table expression. Has precedence over the Union field.
	Select *Select
	// Union clause as a common table expression. Select field pointer must be
	// nil to trigger SQL generation of this field.
	Union *Union
}

// With represents a common table expression. Common Table Expressions (CTEs)
// are a standard SQL feature, and are essentially temporary named result sets.
// Non-recursive CTES are basically 'query-local VIEWs'. One CTE can refer to
// another. The syntax is more readable than nested FROM (SELECT ...). One can
// refer to a CTE from multiple places. They are better than copy-pasting
// FROM(SELECT ...)
//
// Common Table Expression versus Derived Table: Better readability; Can be
// referenced multiple times; Can refer to other CTEs; Improved performance.
//
// https://dev.mysql.com/doc/refman/8.0/en/with.html
//
// https://mariadb.com/kb/en/mariadb/non-recursive-common-table-expressions-overview/
//
// http://mysqlserverteam.com/mysql-8-0-labs-recursive-common-table-expressions-in-mysql-ctes/
//
// Supported in: MySQL >=8.0.1 and MariaDb >=10.2
type With struct {
	Log log.Logger // Log optional logger
	// DB gets required once the Load*() functions will be used.
	DB QueryPreparer

	Subclauses []WithCTE
	// TopLevel a union type which allows only one of the fields to be set.
	TopLevel struct {
		Select *Select
		Union  *Union
		Update *Update
		Delete *Delete
	}
	IsRecursive   bool // See Recursive()
	IsInterpolate bool // See Interpolate()
	// UseBuildCache if `true` the final build query including place holders
	// will be cached in a private field. Each time a call to function ToSQL
	// happens, the arguments will be re-evaluated and returned or interpolated.
	UseBuildCache bool
	cacheArgs     Arguments // like a buffer, gets reused
	cacheSQL      []byte
}

// NewWith creates a new WITH statement with multiple common table expressions
// (CTE).
func NewWith(expressions ...WithCTE) *With {
	return &With{
		Subclauses: expressions,
	}
}

// With creates a new With which selects from the provided columns.
// Columns won't get quoted.
func (c *Connection) With(expressions ...WithCTE) *With {
	return &With{
		Log:        c.Log,
		Subclauses: expressions,
		DB:         c.DB,
	}
}

// With creates a new With that select that given columns bound to the transaction
func (tx *Tx) With(expressions ...WithCTE) *With {
	return &With{
		Log:        tx.Logger,
		Subclauses: expressions,
		DB:         tx.Tx,
	}
}

// WithDB sets the database query object.
func (b *With) WithDB(db QueryPreparer) *With {
	b.DB = db
	return b
}

// Select gets used in the top level statement.
func (b *With) Select(topLevel *Select) *With {
	b.TopLevel.Select = topLevel
	return b
}

// Update gets used in the top level statement.
func (b *With) Update(topLevel *Update) *With {
	b.TopLevel.Update = topLevel
	return b
}

// Delete gets used in the top level statement.
func (b *With) Delete(topLevel *Delete) *With {
	b.TopLevel.Delete = topLevel
	return b
}

// Union gets used in the top level statement.
func (b *With) Union(topLevel *Union) *With {
	b.TopLevel.Union = topLevel
	return b
}

// Recursive common table expressions are one having a subquery that refers to
// its own name. The WITH clause must begin with WITH RECURSIVE if any CTE in
// the WITH clause refers to itself. (If no CTE refers to itself, RECURSIVE is
// permitted but not required.) Common applications of recursive CTEs include
// series generation and traversal of hierarchical or tree-structured data. It
// is simpler, when experimenting with WITH RECURSIVE, to put this at the start
// of your session: `SET max_execution_time = 10000;` so that the runaway query
// aborts automatically after 10 seconds, if the WHERE clause wasn’t correct.
func (b *With) Recursive() *With {
	b.IsRecursive = true
	return b
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `Arguments` will then be nil.
func (b *With) Interpolate() *With {
	b.IsInterpolate = true
	return b
}

// ToSQL converts the select statement into a string and returns its arguments.
func (b *With) ToSQL() (string, []interface{}, error) {
	return toSQL(b, b.IsInterpolate, _isNotPrepared)
}

func (b *With) writeBuildCache(sql []byte) {
	b.cacheSQL = sql
}

func (b *With) readBuildCache() (sql []byte, _ Arguments, err error) {
	if b.cacheSQL == nil {
		return nil, nil, nil
	}
	b.cacheArgs, err = b.appendArgs(b.cacheArgs[:0])
	return b.cacheSQL, b.cacheArgs, err
}

func (b *With) hasBuildCache() bool {
	return b.UseBuildCache
}

func (b *With) toSQL(w queryWriter) error {

	w.WriteString("WITH")
	if b.IsRecursive {
		w.WriteString(" RECURSIVE")
	}
	w.WriteRune('\n')

	for i, sc := range b.Subclauses {
		Quoter.writeName(w, sc.Name)
		if len(sc.Columns) > 0 {
			w.WriteRune(' ')
			w.WriteRune('(')
			for j, c := range sc.Columns {
				if j > 0 {
					w.WriteRune(',')
				}
				Quoter.writeName(w, c)
			}
			w.WriteRune(')')
		}
		w.WriteString(" AS (")
		switch {
		case sc.Select != nil:
			sc.Select.IsInterpolate = b.IsInterpolate
			sc.Select.UseBuildCache = b.UseBuildCache
			if err := sc.Select.toSQL(w); err != nil {
				return errors.Wrap(err, "[dbr] sc.Select.toSQL")
			}
		case sc.Union != nil:
			sc.Union.IsInterpolate = b.IsInterpolate
			sc.Union.UseBuildCache = b.UseBuildCache
			if err := sc.Union.toSQL(w); err != nil {
				return errors.Wrap(err, "[dbr] sc.Union.toSQL")
			}
		}
		w.WriteRune(')')
		if i < len(b.Subclauses)-1 {
			w.WriteRune(',')
		}
		w.WriteRune('\n')
	}

	switch {
	case b.TopLevel.Select != nil:
		b.TopLevel.Select.IsInterpolate = b.IsInterpolate
		b.TopLevel.Select.UseBuildCache = b.UseBuildCache
		return errors.WithStack(b.TopLevel.Select.toSQL(w))

	case b.TopLevel.Union != nil:
		b.TopLevel.Union.IsInterpolate = b.IsInterpolate
		b.TopLevel.Union.UseBuildCache = b.UseBuildCache
		return errors.WithStack(b.TopLevel.Union.toSQL(w))

	case b.TopLevel.Update != nil:
		b.TopLevel.Update.IsInterpolate = b.IsInterpolate
		b.TopLevel.Update.UseBuildCache = b.UseBuildCache
		return errors.WithStack(b.TopLevel.Update.toSQL(w))

	case b.TopLevel.Delete != nil:
		b.TopLevel.Delete.IsInterpolate = b.IsInterpolate
		b.TopLevel.Delete.UseBuildCache = b.UseBuildCache
		return errors.WithStack(b.TopLevel.Delete.toSQL(w))
	}
	return errors.NewEmptyf("[dbr] Type With misses a top level statement")
}

func (b *With) appendArgs(args Arguments) (_ Arguments, err error) {
	for _, sc := range b.Subclauses {
		switch {
		case sc.Select != nil:
			args, err = sc.Select.appendArgs(args)
		case sc.Union != nil:
			args, err = sc.Union.appendArgs(args)
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	switch {
	case b.TopLevel.Select != nil:
		return b.TopLevel.Select.appendArgs(args)
	case b.TopLevel.Union != nil:
		return b.TopLevel.Union.appendArgs(args)
	case b.TopLevel.Update != nil:
		return b.TopLevel.Update.appendArgs(args)
	case b.TopLevel.Delete != nil:
		return b.TopLevel.Delete.appendArgs(args)
	}
	return nil, errors.NewEmptyf("[dbr] Type With misses a top level statement")
}

// Query executes a query and returns many rows.
func (b *With) Query(ctx context.Context) (*sql.Rows, error) {
	rows, err := Query(ctx, b.DB, b)
	return rows, errors.WithStack(err)
}

// Prepare prepares a SQL statement. Sets IsInterpolate to false.
func (b *With) Prepare(ctx context.Context) (*sql.Stmt, error) {
	stmt, err := Prepare(ctx, b.DB, b)
	return stmt, errors.WithStack(err)
}

// Load loads data from a query into an object. You must set DB.QueryContext on
// the With object or it just panics. Load can load a single row or n-rows.
func (b *With) Load(ctx context.Context, s Scanner) (rowCount int64, err error) {
	rowCount, err = Load(ctx, b.DB, b, s)
	return rowCount, errors.WithStack(err)
}
