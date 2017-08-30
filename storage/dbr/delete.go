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

// Delete contains the clauses for a DELETE statement.
//
// InnoDB Tables: If you are deleting many rows from a large table, you may
// exceed the lock table size for an InnoDB table. To avoid this problem, or
// simply to minimize the time that the table remains locked, the following
// strategy (which does not use DELETE at all) might be helpful:
//
// Select the rows not to be deleted into an empty table that has the same
// structure as the original table:
//	INSERT INTO t_copy SELECT * FROM t WHERE ... ;
// Use RENAME TABLE to atomically move the original table out of the way and
// rename the copy to the original name:
//	RENAME TABLE t TO t_old, t_copy TO t;
// Drop the original table:
//	DROP TABLE t_old;
// No other sessions can access the tables involved while RENAME TABLE executes,
// so the rename operation is not subject to concurrency problems.
// TODO(CyS) add DELETE ... JOIN ... statement SQLStmtDeleteJoin
type Delete struct {
	BuilderBase
	BuilderConditional
	// DB can be either a *sql.DB (connection pool), a *sql.Conn (a single
	// dedicated database session) or a *sql.Tx (an in-progress database
	// transaction).
	DB ExecPreparer
	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners DeleteListeners
}

// NewDelete creates a new Delete object.
func NewDelete(from string) *Delete {
	return &Delete{
		BuilderBase: BuilderBase{
			Table: MakeIdentifier(from),
		},
		BuilderConditional: BuilderConditional{
			Wheres: make(Conditions, 0, 2),
		},
	}
}

func newDeleteFrom(db ExecPreparer, l log.Logger, from, id string) *Delete {
	return &Delete{
		BuilderBase: BuilderBase{
			id:    id,
			Table: MakeIdentifier(from),
			Log:   l,
		},
		BuilderConditional: BuilderConditional{
			Wheres: make(Conditions, 0, 2),
		},
		DB: db,
	}
}

// DeleteFrom creates a new Delete for the given table
func (c *ConnPool) DeleteFrom(from string) *Delete {
	l := c.Log
	id := c.makeUniqueID()
	if l != nil {
		l = c.Log.With(log.String("deleteID", id), log.String("table", from))
	}
	return newDeleteFrom(c.DB, l, from, id)
}

// DeleteFrom creates a new Delete for the given table
// in the context for a single database connection.
func (c *Conn) DeleteFrom(from string) *Delete {
	l := c.Log
	id := c.makeUniqueID()
	if l != nil {
		l = c.Log.With(log.String("deleteID", id), log.String("table", from))
	}
	return newDeleteFrom(c.DB, l, from, id)
}

// DeleteFrom creates a new Delete for the given table
// in the context for a transaction
func (tx *Tx) DeleteFrom(from string) *Delete {
	l := tx.Log
	id := tx.makeUniqueID()
	if l != nil {
		l = tx.Log.With(log.String("deleteID", id), log.String("table", from))
	}
	return newDeleteFrom(tx.DB, l, from, id)
}

// Alias sets an alias for the table name.
func (b *Delete) Alias(alias string) *Delete {
	b.Table.Aliased = alias
	return b
}

// WithDB sets the database query object.
func (b *Delete) WithDB(db ExecPreparer) *Delete {
	b.DB = db
	return b
}

// Unsafe see BuilderBase.IsUnsafe which weakens security when building the SQL
// string. This function must be called before calling any other function.
func (b *Delete) Unsafe() *Delete {
	b.IsUnsafe = true
	return b
}

// BindRecord binds the qualified record to the main table/view, or any other
// table/view/alias used in the query, for assembling and appending arguments.
// An ArgumentsAppender gets called if it matches the qualifier, in this case
// the current table name or its alias.
func (b *Delete) BindRecord(records ...QualifiedRecord) *Delete {
	b.bindRecord(records)
	return b
}

func (b *Delete) bindRecord(records []QualifiedRecord) {
	if b.ArgumentsAppender == nil {
		b.ArgumentsAppender = make(map[string]ArgumentsAppender)
	}
	for _, rec := range records {
		q := rec.Qualifier
		if q == "" {
			q = b.Table.mustQualifier()
		}
		b.ArgumentsAppender[q] = rec.Record
	}
}

// Where appends a WHERE clause to the statement whereSQLOrMap can be a string
// or map. If it'ab a string, args wil replaces any places holders.
func (b *Delete) Where(wf ...*Condition) *Delete {
	b.Wheres = append(b.Wheres, wf...)
	return b
}

// OrderBy appends columns to the ORDER BY statement for ascending sorting. A
// column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. When you use ORDER BY or GROUP BY to sort a column
// in a DELETE, the server sorts arguments using only the initial number of
// bytes indicated by the max_sort_length system variable.
func (b *Delete) OrderBy(columns ...string) *Delete {
	b.OrderBys = b.OrderBys.AppendColumns(b.IsUnsafe, columns...)
	return b
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// A column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. When you use ORDER BY or GROUP BY to sort a column
// in a DELETE, the server sorts arguments using only the initial number of
// bytes indicated by the max_sort_length system variable.
func (b *Delete) OrderByDesc(columns ...string) *Delete {
	b.OrderBys = b.OrderBys.AppendColumns(b.IsUnsafe, columns...).applySort(len(columns), sortDescending)
	return b
}

// Limit sets a LIMIT clause for the statement; overrides any existing LIMIT
func (b *Delete) Limit(limit uint64) *Delete {
	b.LimitCount = limit
	b.LimitValid = true
	return b
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `args` will then be nil.
func (b *Delete) Interpolate() *Delete {
	b.IsInterpolate = true
	return b
}

// ToSQL serialized the Delete to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Delete) ToSQL() (string, []interface{}, error) {
	return toSQL(b, b.IsInterpolate, _isNotPrepared)
}

func (b *Delete) writeBuildCache(sql []byte) {
	b.cacheSQL = sql
}

func (b *Delete) readBuildCache() (sql []byte, _ Arguments, err error) {
	if b.cacheSQL == nil {
		return nil, nil, nil
	}
	b.cacheArgs, err = b.appendArgs(b.cacheArgs[:0])
	return b.cacheSQL, b.cacheArgs, err
}

// BuildCache if `true` the final build query including place holders will be
// cached in a private field. Each time a call to function ToSQL happens, the
// arguments will be re-evaluated and returned or interpolated.
func (b *Delete) BuildCache() *Delete {
	b.IsBuildCache = true
	return b
}

func (b *Delete) hasBuildCache() bool {
	return b.IsBuildCache
}

// ToSQL serialized the Delete to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Delete) toSQL(buf *bytes.Buffer) error {

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return errors.WithStack(err)
	}

	if b.RawFullSQL != "" {
		buf.WriteString(b.RawFullSQL)
		return nil
	}

	if b.Table.Name == "" {
		return errors.NewEmptyf("[dbr] Delete: Table is missing")
	}

	buf.WriteString("DELETE ")
	writeStmtID(buf, b.id)
	buf.WriteString("FROM ")
	b.Table.WriteQuoted(buf)

	// TODO(CyS) add SQLStmtDeleteJoin

	if err := b.Wheres.write(buf, 'w'); err != nil {
		return errors.WithStack(err)
	}

	sqlWriteOrderBy(buf, b.OrderBys, false)
	sqlWriteLimitOffset(buf, b.LimitValid, b.LimitCount, false, 0)

	return nil
}

// ToSQL serialized the Delete to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Delete) appendArgs(args Arguments) (_ Arguments, err error) {

	if b.RawFullSQL != "" {
		return b.RawArguments, nil
	}
	if cap(args) == 0 {
		args = make(Arguments, 0, len(b.Wheres))
	}
	args, err = b.Table.appendArgs(args)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO(CyS) add SQLStmtDeleteJoin
	placeHolderColumns := make([]string, 0, len(b.Wheres)) // can be reused once we implement more features of the DELETE statement, like JOINs.

	args, pap, err := b.Wheres.appendArgs(args, appendArgsWHERE)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if boundCols := b.Wheres.intersectConditions(placeHolderColumns); len(boundCols) > 0 {
		if args, err = appendArgs(pap, b.ArgumentsAppender, args, b.Table.mustQualifier(), boundCols); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, nil
}

// Exec executes the statement represented by the Delete. It returns the raw
// database/sql Result and an error if there was one. If info mode for logging
// has been enabled it logs the duration taken. In debug mode the SQL string.
func (b *Delete) Exec(ctx context.Context) (sql.Result, error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Exec", log.Stringer("sql", b))
	}
	r, err := Exec(ctx, b.DB, b)
	return r, errors.WithStack(err)
}

// Prepare executes the statement represented by the Delete to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Delete are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement. If debug mode for logging has been enabled it logs the
// duration taken and the SQL string.
func (b *Delete) Prepare(ctx context.Context) (*StmtDelete, error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Prepare", log.Stringer("sql", b))
	}
	sqlStmt, err := Prepare(ctx, b.DB, b)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cap := len(b.Wheres)
	return &StmtDelete{
		StmtBase: StmtBase{
			stmt:       sqlStmt,
			argsCache:  make(Arguments, 0, cap),
			argsRaw:    make([]interface{}, 0, cap),
			bindRecord: b.bindRecord,
		},
		del: b,
	}, nil
}

// StmtDelete wraps a *sql.Stmt with a specific SQL query. To create a
// StmtDelete call the Prepare function of type Delete. StmtDelete is not safe
// for concurrent use, despite the underlying *sql.Stmt is. Don't forget to call
// Close!
type StmtDelete struct {
	StmtBase
	del *Delete
}

// WithArguments sets the arguments for the execution with Exec. It internally resets
// previously applied arguments.
func (st *StmtDelete) WithArguments(args Arguments) *StmtDelete {
	st.withArguments(args)
	return st
}

// WithRecords sets the records for the execution with Exec. It internally
// resets previously applied arguments.
func (st *StmtDelete) WithRecords(records ...QualifiedRecord) *StmtDelete {
	st.withRecords(st.del.appendArgs, records...)
	return st
}
