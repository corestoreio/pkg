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

// DeleteFrom creates a new Delete for the given table
func (c *Connection) DeleteFrom(from string) *Delete {
	return &Delete{
		BuilderBase: BuilderBase{
			Table: MakeIdentifier(from),
			Log:   c.Log,
		},
		BuilderConditional: BuilderConditional{
			Wheres: make(Conditions, 0, 2),
		},
		DB: c.DB,
	}
}

// DeleteFrom creates a new Delete for the given table
// in the context for a transaction
func (tx *Tx) DeleteFrom(from string) *Delete {
	return &Delete{
		BuilderBase: BuilderBase{
			Table: MakeIdentifier(from),
			Log:   tx.Logger,
		},
		BuilderConditional: BuilderConditional{
			Wheres: make(Conditions, 0, 2),
		},
		DB: tx.Tx,
	}
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

// SetRecord pulls in arguments to match Columns from the argument appender.
func (b *Delete) SetRecord(aa ArgumentsAppender) *Delete {
	b.Record = aa
	return b
}

// Where appends a WHERE clause to the statement whereSQLOrMap can be a string
// or map. If it'ab a string, args wil replaces any places holders.
func (b *Delete) Where(wf ...*Condition) *Delete {
	b.Wheres = append(b.Wheres, wf...)
	return b
}

// OrderBy appends columns to the ORDER BY statement for ascending sorting.
// Columns are getting quoted. When you use ORDER BY or GROUP BY to sort a
// column in a DELETE, the server sorts arguments using only the initial number
// of bytes indicated by the max_sort_length system variable.
func (b *Delete) OrderBy(columns ...string) *Delete {
	b.OrderBys = b.OrderBys.appendColumns(columns, false)
	return b
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// Columns are getting quoted. When you use ORDER BY or GROUP BY to sort a
// column in a DELETE, the server sorts arguments using only the initial number
// of bytes indicated by the max_sort_length system variable.
func (b *Delete) OrderByDesc(columns ...string) *Delete {
	b.OrderBys = b.OrderBys.appendColumns(columns, false).applySort(len(columns), sortDescending)
	return b
}

// OrderByExpr adds a custom SQL expression to the ORDER BY clause. Does not
// quote the strings.
func (b *Delete) OrderByExpr(columns ...string) *Delete {
	b.OrderBys = b.OrderBys.appendColumns(columns, true)
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
// prepared statements. ToSQLs second argument `Arguments` will then be nil.
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

	buf.WriteString("DELETE FROM ")
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

	args, pap, err := b.Wheres.appendArgs(args, appendArgsWHERE)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	placeHolderColumns := make([]string, 0, len(b.Wheres)) // can be reused once we implement more features of the DELETE statement, like JOINs.
	if args, err = appendAssembledArgs(pap, b.Record, args, SQLStmtDelete|SQLPartWhere, b.Wheres.intersectConditions(placeHolderColumns)); err != nil {
		return nil, errors.WithStack(err)
	}
	return args, nil
}

// Exec executes the statement represented by the Delete
// It returns the raw database/sql Result and an error if there was one
func (b *Delete) Exec(ctx context.Context) (sql.Result, error) {
	r, err := Exec(ctx, b.DB, b)
	return r, errors.WithStack(err)
}

// Prepare executes the statement represented by the Delete. It returns the raw
// database/sql Statement and an error if there was one. Provided arguments in
// the Delete are getting ignored. It panics when field Preparer at nil.
func (b *Delete) Prepare(ctx context.Context) (*sql.Stmt, error) {
	stmt, err := Prepare(ctx, b.DB, b)
	return stmt, errors.WithStack(err)
}
