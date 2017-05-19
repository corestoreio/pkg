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
type Delete struct {
	Log log.Logger // Log optional logger
	DB  struct {
		Preparer
		Execer
	}

	// TODO(CyS) add DELETE ... JOIN ... statement

	RawFullSQL   string
	RawArguments Arguments // Arguments used by RawFullSQL or BuildCache

	// Record if set retrieves the necessary arguments from the interface.
	Record ArgumentAssembler

	From alias
	WhereFragments
	OrderBys    aliases
	LimitCount  uint64
	OffsetCount uint64
	LimitValid  bool
	OffsetValid bool
	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	IsInterpolate      bool // See Interpolate()
	// UseBuildCache if set to true the final build query will be stored in
	// field private field `buildCache` and the arguments in field `Arguments`
	UseBuildCache bool
	buildCache    []byte
	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners DeleteListeners
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
}

// NewDelete creates a new object with a black hole logger.
func NewDelete(from ...string) *Delete {
	return &Delete{
		From: MakeAlias(from...),
	}
}

// DeleteFrom creates a new Delete for the given table
func (c *Connection) DeleteFrom(from ...string) *Delete {
	d := &Delete{
		Log:            c.Log,
		From:           MakeAlias(from...),
		WhereFragments: make(WhereFragments, 0, 2),
	}
	d.DB.Execer = c.DB
	d.DB.Preparer = c.DB
	return d
}

// DeleteFrom creates a new Delete for the given table
// in the context for a transaction
func (tx *Tx) DeleteFrom(from ...string) *Delete {
	d := &Delete{
		Log:  tx.Logger,
		From: MakeAlias(from...),
	}
	d.DB.Execer = tx.Tx
	d.DB.Preparer = tx.Tx
	return d
}

// AddRecord pulls in values to match Columns from the record generator.
func (b *Delete) AddRecord(rec ArgumentAssembler) *Delete {
	b.Record = rec
	return b
}

// Where appends a WHERE clause to the statement whereSQLOrMap can be a
// string or map. If it'ab a string, args wil replaces any places holders
func (b *Delete) Where(args ...ConditionArg) *Delete {
	b.WhereFragments = b.WhereFragments.append(args...)
	return b
}

// OrderBy appends columns to the ORDER BY statement for ascending sorting.
// Columns are getting quoted. When you use ORDER BY or GROUP BY to sort a
// column in a DELETE, the server sorts values using only the initial number of
// bytes indicated by the max_sort_length system variable.
func (b *Delete) OrderBy(columns ...string) *Delete {
	b.OrderBys = b.OrderBys.appendColumns(columns, false)
	return b
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// Columns are getting quoted. When you use ORDER BY or GROUP BY to sort a
// column in a DELETE, the server sorts values using only the initial number of
// bytes indicated by the max_sort_length system variable.
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

// Offset sets an OFFSET clause for the statement; overrides any existing OFFSET
func (b *Delete) Offset(offset uint64) *Delete {
	b.OffsetCount = offset
	b.OffsetValid = true
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
func (b *Delete) ToSQL() (string, Arguments, error) {
	return toSQL(b, b.IsInterpolate)
}

func (b *Delete) writeBuildCache(sql []byte, arguments Arguments) {
	b.buildCache = sql
	b.RawArguments = arguments
}

func (b *Delete) readBuildCache() (sql []byte, arguments Arguments) {
	return b.buildCache, b.RawArguments
}

func (b *Delete) hasBuildCache() bool {
	return b.UseBuildCache
}

// ToSQL serialized the Delete to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Delete) toSQL(buf queryWriter) (Arguments, error) {

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return nil, errors.Wrap(err, "[dbr] Delete.Listeners.dispatch")
	}

	if b.RawFullSQL != "" {
		buf.WriteString(b.RawFullSQL)
		return b.RawArguments, nil
	}

	if len(b.From.Name) == 0 {
		return nil, errors.NewEmptyf(errTableMissing)
	}

	buf.WriteString("DELETE FROM ")
	b.From.FquoteAs(buf)

	// Write WHERE clause if we have any fragments.
	// pap defines the pending argument positions. The pending arguments gets
	// assembled in the Record.AssembleArguments.
	args, pap, err := b.WhereFragments.write(buf, make(Arguments, 0, len(b.WhereFragments)), 'w')
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Delete.ToSQL.write")
	}
	if args, err = appendAssembledArgs(pap, b.Record, args, stmtTypeDelete, nil, b.WhereFragments.Conditions()); err != nil {
		return nil, errors.Wrap(err, "[dbr] Delete.toSQL.appendAssembledArgs")
	}

	sqlWriteOrderBy(buf, b.OrderBys, false)
	sqlWriteLimitOffset(buf, b.LimitValid, b.LimitCount, b.OffsetValid, b.OffsetCount)

	return args, nil
}

// Exec executes the statement represented by the Delete
// It returns the raw database/sql Result and an error if there was one
func (b *Delete) Exec(ctx context.Context) (sql.Result, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Delete.Exec.ToSQL")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Delete.Exec.Timing", log.String("sql", sqlStr))
	}

	result, err := b.DB.ExecContext(ctx, sqlStr, args.Interfaces()...)
	if err != nil {
		return result, errors.Wrap(err, "[dbr] delete.exec.Exec")
	}

	return result, nil
}

// Prepare executes the statement represented by the Delete. It returns the raw
// database/sql Statement and an error if there was one. Provided arguments in
// the Delete are getting ignored. It panics when field Preparer at nil.
func (b *Delete) Prepare(ctx context.Context) (*sql.Stmt, error) {
	sqlStr, _, err := b.ToSQL() // TODO create a ToSQL version without any arguments
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Delete.Prepare.ToSQL")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Delete.Prepare.Timing", log.String("sql", sqlStr))
	}

	stmt, err := b.DB.PrepareContext(ctx, sqlStr)
	return stmt, errors.Wrap(err, "[dbr] Delete.Prepare.Prepare")
}
