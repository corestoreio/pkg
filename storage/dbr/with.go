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
	"bytes"
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
	BuilderBase
	// DB can be either a *sql.DB (connection pool), a *sql.Conn (a single
	// dedicated database session) or a *sql.Tx (an in-progress database
	// transaction).
	DB QueryPreparer

	Subclauses []WithCTE
	// TopLevel a union type which allows only one of the fields to be set.
	TopLevel struct {
		Select *Select
		Union  *Union
		Update *Update
		Delete *Delete
	}
	IsRecursive bool // See Recursive()
}

// NewWith creates a new WITH statement with multiple common table expressions
// (CTE).
func NewWith(expressions ...WithCTE) *With {
	return &With{
		Subclauses: expressions,
	}
}

func withInitLog(l log.Logger, expressions []WithCTE, id string) log.Logger {
	if l != nil {
		tables := make([]string, len(expressions))
		for i, w := range expressions {
			tables[i] = w.Name
		}
		l = l.With(log.String("with_cte_id", id), log.Strings("tables", tables...))
	}
	return l
}

// With creates a new With statement.
func (c *ConnPool) With(expressions ...WithCTE) *With {
	id := c.makeUniqueID()
	return &With{
		BuilderBase: BuilderBase{
			id:  id,
			Log: withInitLog(c.Log, expressions, id),
		},
		Subclauses: expressions,
		DB:         c.DB,
	}
}

// With creates a new With statement bound to a single connection.
func (c *Conn) With(expressions ...WithCTE) *With {
	id := c.makeUniqueID()
	return &With{
		BuilderBase: BuilderBase{
			id:  id,
			Log: withInitLog(c.Log, expressions, id),
		},
		Subclauses: expressions,
		DB:         c.DB,
	}
}

// With creates a new With that select that given columns bound to the transaction
func (tx *Tx) With(expressions ...WithCTE) *With {
	id := tx.makeUniqueID()
	return &With{
		BuilderBase: BuilderBase{
			id:  id,
			Log: withInitLog(tx.Log, expressions, id),
		},
		Subclauses: expressions,
		DB:         tx.DB,
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
// prepared statements. ToSQLs second argument `args` will then be nil.
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

// BuildCache if `true` the final build query including place holders will be
// cached in a private field. Each time a call to function ToSQL happens, the
// arguments will be re-evaluated and returned or interpolated.
func (b *With) BuildCache() *With {
	b.IsBuildCache = true
	return b
}

func (b *With) hasBuildCache() bool {
	return b.IsBuildCache
}

func (b *With) toSQL(w *bytes.Buffer) error {

	w.WriteString("WITH ")
	writeStmtID(w, b.id)
	if b.IsRecursive {
		w.WriteString("RECURSIVE ")
	}

	for i, sc := range b.Subclauses {
		Quoter.quote(w, sc.Name)
		if len(sc.Columns) > 0 {
			w.WriteRune(' ')
			w.WriteRune('(')
			for j, c := range sc.Columns {
				if j > 0 {
					w.WriteRune(',')
				}
				Quoter.quote(w, c)
			}
			w.WriteRune(')')
		}
		w.WriteString(" AS (")
		switch {
		case sc.Select != nil:
			sc.Select.IsInterpolate = b.IsInterpolate
			sc.Select.IsBuildCache = b.IsBuildCache
			if err := sc.Select.toSQL(w); err != nil {
				return errors.WithStack(err)
			}
		case sc.Union != nil:
			sc.Union.IsInterpolate = b.IsInterpolate
			sc.Union.IsBuildCache = b.IsBuildCache
			if err := sc.Union.toSQL(w); err != nil {
				return errors.WithStack(err)
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
		b.TopLevel.Select.IsBuildCache = b.IsBuildCache
		return errors.WithStack(b.TopLevel.Select.toSQL(w))

	case b.TopLevel.Union != nil:
		b.TopLevel.Union.IsInterpolate = b.IsInterpolate
		b.TopLevel.Union.IsBuildCache = b.IsBuildCache
		return errors.WithStack(b.TopLevel.Union.toSQL(w))

	case b.TopLevel.Update != nil:
		b.TopLevel.Update.IsInterpolate = b.IsInterpolate
		b.TopLevel.Update.IsBuildCache = b.IsBuildCache
		return errors.WithStack(b.TopLevel.Update.toSQL(w))

	case b.TopLevel.Delete != nil:
		b.TopLevel.Delete.IsInterpolate = b.IsInterpolate
		b.TopLevel.Delete.IsBuildCache = b.IsBuildCache
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

func (b *With) bindRecord(records []QualifiedRecord) {
	// Current pattern: To whom it may concern.

	for _, sc := range b.Subclauses {
		switch {
		case sc.Select != nil:
			sc.Select.bindRecord(records)
		case sc.Union != nil:
			sc.Union.bindRecord(records)
		}
	}

	switch {
	case b.TopLevel.Select != nil:
		b.TopLevel.Select.bindRecord(records)
	case b.TopLevel.Union != nil:
		b.TopLevel.Union.bindRecord(records)
	case b.TopLevel.Update != nil:
		b.TopLevel.Update.bindRecord(records)
	case b.TopLevel.Delete != nil:
		b.TopLevel.Delete.bindRecord(records)
	}
}

// Query executes a query and returns many rows.
func (b *With) Query(ctx context.Context) (*sql.Rows, error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Query", log.Stringer("sql", b))
	}
	rows, err := Query(ctx, b.DB, b)
	return rows, errors.WithStack(err)
}

// Load loads data from a query into an object. You must set DB.QueryContext on
// the With object or it just panics. Load can load a single row or n-rows.
func (b *With) Load(ctx context.Context, s Scanner) (rowCount int64, err error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Load", log.Int64("row_count", rowCount), log.Stringer("sql", b))
	}
	rowCount, err = Load(ctx, b.DB, b, s)
	return rowCount, errors.WithStack(err)
}

// Prepare prepares a SQL statement. Sets IsInterpolate to false.
func (b *With) Prepare(ctx context.Context) (*StmtWith, error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Prepare", log.Stringer("sql", b))
	}
	stmt, err := Prepare(ctx, b.DB, b)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	const cap = 10 // just a guess; needs to be more precise but later.
	return &StmtWith{
		StmtBase: StmtBase{
			stmt:       stmt,
			argsCache:  make(Arguments, 0, cap),
			argsRaw:    make([]interface{}, 0, cap),
			bindRecord: b.bindRecord,
			log:        b.Log,
		},
		with: b,
	}, nil
}

// StmtWith wraps a *sql.Stmt with a specific SQL query. To create a
// StmtWith call the Prepare function of type Union. StmtWith is not safe
// for concurrent use, despite the underlying *sql.Stmt is. Don't forget to call
// Close!
type StmtWith struct {
	StmtBase
	with *With
}

// WithArgs sets the interfaced arguments for the execution with Query+. It
// internally resets previously applied arguments.
func (st *StmtWith) WithArgs(args ...interface{}) *StmtWith {
	st.withArgs(args)
	return st
}

// WithArguments sets the arguments for the execution with Query+. It internally
// resets previously applied arguments.
func (st *StmtWith) WithArguments(args Arguments) *StmtWith {
	st.withArguments(args)
	return st
}

// WithRecords sets the records for the execution with Query+. It internally
// resets previously applied arguments.
func (st *StmtWith) WithRecords(records ...QualifiedRecord) *StmtWith {
	st.withRecords(st.with.appendArgs, records...)
	return st
}
