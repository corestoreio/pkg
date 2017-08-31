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

// Update contains the clauses for an UPDATE statement
type Update struct {
	BuilderBase
	BuilderConditional
	// DB can be either a *sql.DB (connection pool), a *sql.Conn (a single
	// dedicated database session) or a *sql.Tx (an in-progress database
	// transaction).
	DB ExecPreparer

	// TODO: add UPDATE JOINS SQLStmtUpdateJoin

	// SetClausAliases only applicable in case when Record field has been set or
	// ExecMulti gets used. `SetClausAliases` contains the lis of column names
	// which gets passed to the ArgumentsAppender function. If empty
	// `SetClausAliases` collects the column names from the `SetClauses`. The
	// alias slice must have the same length as the columns slice. Despite
	// setting `SetClausAliases` the SetClauses.Columns must be provided to
	// create a valid SQL statement.
	SetClausAliases []string
	// SetClauses contains the column/argument association. For each column
	// there must be one argument.
	SetClauses Conditions
	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners UpdateListeners
}

// NewUpdate creates a new Update object.
func NewUpdate(table string) *Update {
	return &Update{
		BuilderBase: BuilderBase{
			Table: MakeIdentifier(table),
		},
	}
}

func newUpdate(db ExecPreparer, idFn uniqueIDFn, l log.Logger, table string) *Update {
	id := idFn()
	if l != nil {
		l = l.With(log.String("update_id", id), log.String("table", table))
	}
	return &Update{
		BuilderBase: BuilderBase{
			id:    id,
			Table: MakeIdentifier(table),
			Log:   l,
		},
		DB: db,
	}
}

// Update creates a new Update for the given table with a random connection from
// the pool.
func (c *ConnPool) Update(table string) *Update {
	return newUpdate(c.DB, c.makeUniqueID, c.Log, table)
}

// Update creates a new Update for the given table bound to a single connection.
func (c *Conn) Update(table string) *Update {
	return newUpdate(c.DB, c.makeUniqueID, c.Log, table)
}

// Update creates a new Update for the given table bound to a transaction.
func (tx *Tx) Update(table string) *Update {
	return newUpdate(tx.DB, tx.makeUniqueID, tx.Log, table)
}

// Alias sets an alias for the table name.
func (b *Update) Alias(alias string) *Update {
	b.Table.Aliased = alias
	return b
}

// WithDB sets the database query object.
func (b *Update) WithDB(db ExecPreparer) *Update {
	b.DB = db
	return b
}

// Unsafe see BuilderBase.IsUnsafe which weakens security when building the SQL
// string. This function must be called before calling any other function.
func (b *Update) Unsafe() *Update {
	b.IsUnsafe = true
	return b
}

// Set appends a column/value pair for the statement.
func (b *Update) Set(c ...*Condition) *Update {
	b.SetClauses = append(b.SetClauses, c...)
	return b
}

// AddColumns adds columns which values gets later derived from an
// ArgumentsAppender. Those columns will get passed to the ArgumentsAppender
// implementation. Mostly used with the type Update.
func (b *Update) AddColumns(columnNames ...string) *Update {
	for _, col := range columnNames {
		b.SetClauses = append(b.SetClauses, Column(col))
	}
	return b
}

// BindRecord binds the qualified record to the main table/view, or any other
// table/view/alias used in the query, for assembling and appending arguments.
// An ArgumentsAppender gets called if it matches the qualifier, in this case
// the current table name or its alias.
func (b *Update) BindRecord(records ...QualifiedRecord) *Update {
	b.bindRecord(records)
	return b
}

func (b *Update) bindRecord(records []QualifiedRecord) {
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

// Where appends a WHERE clause to the statement
func (b *Update) Where(wf ...*Condition) *Update {
	b.Wheres = append(b.Wheres, wf...)
	return b
}

// OrderBy appends columns to the ORDER BY statement for ascending sorting. A
// column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. When you use ORDER BY or GROUP BY to sort a column
// in a UPDATE, the server sorts values using only the initial number of bytes
// indicated by the max_sort_length system variable.
func (b *Update) OrderBy(columns ...string) *Update {
	b.OrderBys = b.OrderBys.AppendColumns(b.IsUnsafe, columns...)
	return b
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// A column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. When you use ORDER BY or GROUP BY to sort a column
// in a UPDATE, the server sorts values using only the initial number of bytes
// indicated by the max_sort_length system variable.
func (b *Update) OrderByDesc(columns ...string) *Update {
	b.OrderBys = b.OrderBys.AppendColumns(b.IsUnsafe, columns...).applySort(len(columns), sortDescending)
	return b
}

// Limit sets a limit for the statement; overrides any existing LIMIT
func (b *Update) Limit(limit uint64) *Update {
	b.LimitCount = limit
	b.LimitValid = true
	return b
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `args` will then be nil.
func (b *Update) Interpolate() *Update {
	b.IsInterpolate = true
	return b
}

// ToSQL converts the select statement into a string and returns its arguments.
func (b *Update) ToSQL() (string, []interface{}, error) {
	return toSQL(b, b.IsInterpolate, _isNotPrepared)
}

func (b *Update) writeBuildCache(sql []byte) {
	b.cacheSQL = sql
}

func (b *Update) readBuildCache() (sql []byte, _ Arguments, err error) {
	if b.cacheSQL == nil {
		return nil, nil, nil
	}
	b.cacheArgs, err = b.appendArgs(b.cacheArgs[:0])
	return b.cacheSQL, b.cacheArgs, err
}

// BuildCache if `true` the final build query including place holders will be
// cached in a private field. Each time a call to function ToSQL happens, the
// arguments will be re-evaluated and returned or interpolated.
func (b *Update) BuildCache() *Update {
	b.IsBuildCache = true
	return b
}

func (b *Update) hasBuildCache() bool {
	return b.IsBuildCache
}

// ToSQL serialized the Update to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Update) toSQL(buf *bytes.Buffer) error {

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return errors.WithStack(err)
	}

	if b.RawFullSQL != "" {
		buf.WriteString(b.RawFullSQL)
		return nil
	}

	if len(b.Table.Name) == 0 {
		return errors.NewEmptyf("[dbr] Update: Table at empty")
	}
	if len(b.SetClauses) == 0 {
		return errors.NewEmptyf("[dbr] Update: No columns specified")
	}

	buf.WriteString("UPDATE ")
	writeStmtID(buf, b.id)
	b.Table.WriteQuoted(buf)
	buf.WriteString(" SET ")

	if err := b.SetClauses.writeSetClauses(buf); err != nil {
		return errors.WithStack(err)
	}

	// Write WHERE clause if we have any fragments
	if err := b.Wheres.write(buf, 'w'); err != nil {
		return errors.WithStack(err)
	}

	sqlWriteOrderBy(buf, b.OrderBys, false)
	sqlWriteLimitOffset(buf, b.LimitValid, b.LimitCount, false, 0)
	return nil
}

// ToSQL serialized the Update to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Update) appendArgs(args Arguments) (_ Arguments, err error) {

	if b.RawFullSQL != "" {
		return b.RawArguments, nil
	}

	if cap(args) == 0 {
		args = make(Arguments, 0, len(b.SetClauses)+len(b.Wheres))
	}

	if b.ArgumentsAppender != nil {
		if len(b.SetClausAliases) == 0 {
			b.SetClausAliases = b.SetClauses.leftHands(b.SetClausAliases)
		}

		qualifier := b.Table.mustQualifier() // if this panics, you have different problems.

		if aa, ok := b.ArgumentsAppender[qualifier]; ok {
			var argCol [1]string
			for _, col := range b.SetClausAliases {
				argCol[0] = col
				args, err = aa.AppendArgs(args, argCol[:])
				if err != nil {
					return nil, errors.Wrapf(err, "[dbr] Update.appendArgs.AppendArgs at qualifier %q and column %q", qualifier, col)
				}
			}
		}
	}

	args, _, err = b.SetClauses.appendArgs(args, appendArgsSET)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Write WHERE clause if we have any fragments
	args, pap, err := b.Wheres.appendArgs(args, appendArgsWHERE)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	placeHolderColumns := make([]string, 0, len(b.Wheres)) // can be reused once we implement more features of the DELETE statement, like JOINs.
	if boundedCols := b.Wheres.intersectConditions(placeHolderColumns); len(boundedCols) > 0 {
		if args, err = appendArgs(pap, b.ArgumentsAppender, args, b.Table.mustQualifier(), boundedCols); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, nil
}

func (b *Update) validate() error {
	if len(b.SetClauses) == 0 {
		return errors.NewEmptyf("[dbr] Update: Columns are empty")
	}
	if len(b.SetClausAliases) > 0 && len(b.SetClausAliases) != len(b.SetClauses) {
		return errors.NewMismatchf("[dbr] Update: ColumnAliases slice and Columns slice must have the same length")
	}
	return nil
}

// Exec interpolates and executes the statement represented by the Update
// object. It returns the raw database/sql Result and an error if there was one.
func (b *Update) Exec(ctx context.Context) (sql.Result, error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Exec", log.Stringer("sql", b))
	}
	if err := b.validate(); err != nil {
		return nil, errors.WithStack(err)
	}
	result, err := Exec(ctx, b.DB, b)
	return result, errors.WithStack(err)
}

// Prepare executes the statement represented by the Update to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Update are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement.
func (b *Update) Prepare(ctx context.Context) (*StmtUpdate, error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Prepare", log.Stringer("sql", b))
	}
	if err := b.validate(); err != nil {
		return nil, errors.WithStack(err)
	}
	stmt, err := Prepare(ctx, b.DB, b)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cap := len(b.SetClauses) + len(b.Wheres)
	return &StmtUpdate{
		StmtBase: StmtBase{
			id:         b.id,
			stmt:       stmt,
			argsCache:  make(Arguments, 0, cap),
			argsRaw:    make([]interface{}, 0, cap),
			bindRecord: b.bindRecord,
			log:        b.Log,
		},
		upd: b,
	}, nil
}

// StmtUpdate wraps a *sql.Stmt with a specific SQL query. To create a
// StmtUpdate call the Prepare function of type Update. StmtUpdate is not safe
// for concurrent use, despite the underlying *sql.Stmt is. Don't forget to call
// Close!
type StmtUpdate struct {
	StmtBase
	upd *Update
}

// WithArguments sets the arguments for the execution with ExecContext. It
// internally resets previously applied arguments.
func (st *StmtUpdate) WithArguments(args Arguments) *StmtUpdate {
	st.withArguments(args)
	return st
}

// WithRecords sets the records for the execution with ExecContext. It
// internally resets previously applied arguments.
func (st *StmtUpdate) WithRecords(records ...QualifiedRecord) *StmtUpdate {
	st.withRecords(st.upd.appendArgs, records...)
	return st
}
