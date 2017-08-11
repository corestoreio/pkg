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

	"github.com/corestoreio/csfw/util/bufferpool"
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
	DB execPreparer

	// TODO: add UPDATE JOINS SQLStmtUpdateJoin

	// RecordColumns only applicable in case when Record field has been set.
	// `RecordColumns` contains the lis of column names which gets passed to the
	// ArgumentsAppender function. If empty `RecordColumns` then the names gets
	// collected from `SetClauses`
	RecordColumns []string
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

// Update creates a new Update for the given table
func (c *Connection) Update(table string) *Update {
	return &Update{
		BuilderBase: BuilderBase{
			Table: MakeIdentifier(table),
			Log:   c.Log,
		},
		DB: c.DB,
	}
}

// Update creates a new Update for the given table bound to a transaction
func (tx *Tx) Update(table string) *Update {
	return &Update{
		BuilderBase: BuilderBase{
			Table: MakeIdentifier(table),
			Log:   tx.Logger,
		},
		DB: tx.Tx,
	}
}

// Alias sets an alias for the table name.
func (b *Update) Alias(alias string) *Update {
	b.Table.Aliased = alias
	return b
}

// WithDB sets the database query object.
func (b *Update) WithDB(db execPreparer) *Update {
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
// implementation. Mostly used with the type UpdateMulti.
func (b *Update) AddColumns(columnNames ...string) *Update {
	for _, col := range columnNames {
		b.SetClauses = append(b.SetClauses, Column(col))
	}
	return b
}

// SetRecord sets a new argument generator type. See the example for more
// details.
func (b *Update) SetRecord(rec ArgumentsAppender) *Update {
	b.Record = rec
	return b
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
func (b *Update) appendArgs(args Arguments) (Arguments, error) {

	if b.RawFullSQL != "" {
		return b.RawArguments, nil
	}

	if cap(args) == 0 {
		args = make(Arguments, 0, len(b.SetClauses)+len(b.Wheres))
	}
	var err error
	if b.Record != nil {
		if len(b.RecordColumns) == 0 {
			b.RecordColumns = b.SetClauses.leftHands(b.RecordColumns)
		}
		args, err = b.Record.AppendArguments(SQLStmtUpdate|SQLPartSet, args, b.RecordColumns)
		if err != nil {
			return nil, errors.WithStack(err)
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
	if args, err = appendAssembledArgs(pap, b.Record, args, SQLStmtUpdate|SQLPartWhere, b.Wheres.intersectConditions(placeHolderColumns)); err != nil {
		return nil, errors.WithStack(err)
	}

	return args, nil
}

// Exec interpolates and executes the statement represented by the Update
// object. It returns the raw database/sql Result and an error if there was one.
func (b *Update) Exec(ctx context.Context) (sql.Result, error) {
	result, err := Exec(ctx, b.DB, b)
	return result, errors.WithStack(err)
}

// Prepare creates a new prepared statement represented by the Update object. It
// returns the raw database/sql Stmt and an error if there was one.
func (b *Update) Prepare(ctx context.Context) (*sql.Stmt, error) {
	stmt, err := Prepare(ctx, b.DB, b)
	return stmt, errors.WithStack(err)
}

// UpdateMulti allows to run an UPDATE statement multiple times with different
// records in an optionally transaction. If you enable the interpolate feature on
// the Update object the interpolated SQL string will be send each time to the
// SQL server otherwise a prepared statement will be created. Create a single
// Update object without the SET columns and without arguments. Add a WHERE
// clause with common conditions and conditions with place holders where the
// value/s get derived from the ArgumentsAppender. The empty WHERE arguments
// trigger the placeholder and the correct operator. The values itself will be
// provided either through the Records slice or via RecordChan.
type UpdateMulti struct {
	Log log.Logger
	// DB can be either a *sql.DB (connection pool), a *sql.Conn (a single
	// dedicated database session) or a *sql.Tx (an in-progress database
	// transaction).
	DB execPreparer

	// IsTransaction set to true to enable running the UPDATE queries in a
	// transaction.
	IsTransaction bool
	// IsolationLevel defines the transaction isolation level.
	sql.IsolationLevel
	// Tx knows how to start a transaction. Must be set if transactions hasn't
	// been disabled.
	Tx *sql.DB // todo can be removed and also the sql.IsolationLevel
	// Update represents the template UPDATE statement.
	Update *Update

	// ColumnAliases provides a special feature, if set, that instead of the
	// column names, the identifiers will be passed to the ArgumentGenerater.Record
	// function. The alias slice must have the same length as the columns slice.
	// Despite setting `ColumnAliases` the Update.SetClauses.Columns must be
	// provided to create a valid SQL statement.
	ColumnAliases []string
}

// NewUpdateMulti creates new UPDATE statement which runs multiple times for a
// specific Update statement.
func NewUpdateMulti(tpl *Update) *UpdateMulti {
	return &UpdateMulti{
		Update: tpl,
	}
}

// UpdateMulti creates a new UpdateMulti for the UPDATE template object.
func (c *Connection) UpdateMulti(tpl *Update) *UpdateMulti {
	return &UpdateMulti{
		Log:    c.Log,
		DB:     c.DB,
		Update: tpl,
	}
}

// UpdateMulti creates a new UpdateMulti for the given UPDATE template object
// bound to a transaction.
func (tx *Tx) UpdateMulti(tpl *Update) *UpdateMulti {
	return &UpdateMulti{
		Log:    tx.Logger,
		DB:     tx.Tx,
		Update: tpl,
	}
}

// WithDB sets the database query object.
func (b *UpdateMulti) WithDB(db execPreparer) *UpdateMulti {
	b.DB = db
	return b
}

// Transaction enables transaction usage and sets an optional isolation level.
// If not set the default database isolation level gets used.
func (b *UpdateMulti) Transaction(level ...sql.IsolationLevel) *UpdateMulti {
	b.IsTransaction = true
	if len(level) == 1 {
		b.IsolationLevel = level[0]
	}
	return b
}

func (b *UpdateMulti) validate() error {
	if len(b.Update.SetClauses) == 0 {
		return errors.NewEmptyf("[dbr] UpdateMulti: Columns are empty")
	}
	if len(b.ColumnAliases) > 0 && len(b.ColumnAliases) != len(b.Update.SetClauses) {
		return errors.NewMismatchf("[dbr] UpdateMulti: ColumnAliases slice and Columns slice must have the same length")
	}
	return nil
}

func txUpdateMultiRollback(tx txer, previousErr error, msg string, args ...interface{}) ([]sql.Result, error) {
	if err := tx.Rollback(); err != nil {
		eArg := []interface{}{previousErr}
		return nil, errors.Wrapf(err, "[dbr] UpdateMulti.Tx.Rollback. Previous Error: %s. "+msg, append(eArg, args...)...)
	}
	return nil, errors.Wrapf(previousErr, msg, args...)
}

// Exec runs multiple UPDATE queries for different records in serial order. The
// returned result slice indexes are same index as for the Records slice.
func (b *UpdateMulti) Exec(ctx context.Context, records ...ArgumentsAppender) ([]sql.Result, error) {
	if err := b.validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.UpdateMulti.Exec.Timing",
			log.Stringer("sql", b.Update),
			log.Int("records", len(records)))
	}

	isInterpolate := b.Update.IsInterpolate
	b.Update.IsInterpolate = false

	sqlBuf := bufferpool.Get()
	defer bufferpool.Put(sqlBuf)

	err := b.Update.toSQL(sqlBuf)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	exec := b.DB
	var tx txer = txMock{}
	if b.IsTransaction {
		tx, err = b.Tx.BeginTx(ctx, &sql.TxOptions{
			Isolation: b.IsolationLevel,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "[dbr] UpdateMulti.Exec.Tx.BeginTx. with Query: %q", b.Update)
		}
		exec = tx
	}

	var stmt *sql.Stmt
	if !isInterpolate {
		stmt, err = exec.PrepareContext(ctx, sqlBuf.String())
		if err != nil {
			return txUpdateMultiRollback(tx, err, "[dbr] UpdateMulti.Exec.Prepare. with Query: %q", b.Update)
		}
		defer stmt.Close() // todo check for error
	}

	if len(b.ColumnAliases) > 0 {
		b.Update.RecordColumns = b.ColumnAliases
	}

	args := make(Arguments, 0, (len(records)+len(b.Update.Wheres))*3) // 3 just a guess
	results := make([]sql.Result, len(records))

	var ipBuf *bytes.Buffer // ip = interpolate buffer
	if isInterpolate {
		ipBuf = bufferpool.Get()
		defer bufferpool.Put(ipBuf)
	}

	for i, rec := range records {
		b.Update.Record = rec
		args, err = b.Update.appendArgs(args)
		if err != nil {
			return txUpdateMultiRollback(tx, err, "[dbr] UpdateMulti.Exec.Interpolate. Index %d with Query: %q", i, sqlBuf)
		}
		if isInterpolate {
			if err = writeInterpolate(ipBuf, sqlBuf.Bytes(), args); err != nil {
				return txUpdateMultiRollback(tx, err, "[dbr] UpdateMulti.Exec.Interpolate. Index %d with Query: %q", i, sqlBuf)
			}

			results[i], err = exec.ExecContext(ctx, ipBuf.String())
			if err != nil {
				return txUpdateMultiRollback(tx, err, "[dbr] UpdateMulti.Exec.Exec. Index %d with Query: %q", i, sqlBuf)
			}
			ipBuf.Reset()
		} else {
			results[i], err = stmt.ExecContext(ctx, args.Interfaces()...)
			if err != nil {
				return txUpdateMultiRollback(tx, err, "[dbr] UpdateMulti.Exec.Stmt.Exec. Index %d with Query: %q", i, sqlBuf)
			}
		}
		args = args[:0] // reset for re-usage
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.WithStack(err)
	}

	return results, nil
}

// ExecChan executes incoming Records and writes the output into the provided
// channels. It closes the channels once the queries have been sent.
// All queries will run parallel, except when using a transaction.
//func (b *UpdateMulti) ExecChan(ctx context.Context, records <-chan ArgumentsAppender, results chan<- sql.Result, errs chan<- error) {
//	defer close(errs)
//	defer close(errs)
//	if err := b.validate(); err != nil {
//		errs <- errors.WithStack(err)
//		return
//	}
//
//	// RecordChan waits for incoming records to send them to the prepared
//	// statement. If the channel gets closed the transaction gets terminated and
//	// the UPDATE statement removed.
//
//	//g, ctx := errgroup.WithContext(ctx)
//	//
//	//g.Go()
//	//
//	//go func() {
//	//	g.Wait()
//	//	close(b.RecordChan)
//	//}()
//	//
//	//if err := g.Wait(); err != nil {
//	//	errChan <- errors.WithStack(err)
//	//}
//
//	// This could run in parallel but it depends if each exec gets a
//	// different connection. In a transaction only serial processing is
//	// possible because a Go transaction gets bound to one connection.
//
//}
