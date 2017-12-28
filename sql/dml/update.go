// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	// SetClausAliases only applicable in case when field QualifiedRecords has
	// been set or ExecMulti gets used. `SetClausAliases` contains the lis of
	// column names which gets passed to the ColumnMapper. If empty,
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
			builderCommon: builderCommon{
				id:  id,
				Log: l,
			},
			Table: MakeIdentifier(table),
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

// AddColumns adds columns which values gets later derived from a ColumnMapper.
// Those columns will get passed to the ColumnMapper implementation.
func (b *Update) AddColumns(columnNames ...string) *Update {
	for _, col := range columnNames {
		b.SetClauses = append(b.SetClauses, Column(col))
	}
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

// WithArgs sets the interfaced arguments for the execution with Query+. It
// internally resets previously applied arguments. This function does not
// support interpolation.
func (b *Update) WithArgs(args ...interface{}) *Update {
	b.withArgs(args)
	return b
}

// WithArguments sets the arguments for the execution with Query+. It internally
// resets previously applied arguments. This function supports interpolation.
func (b *Update) WithArguments(args Arguments) *Update {
	b.withArguments(args)
	return b
}

// WithRecords binds the qualified record to the main table/view, or any other
// table/view/alias used in the query, for assembling and appending arguments.
// The ColumnMapper gets called if it matches the qualifier, in this case the
// current table name or its alias.
func (b *Update) WithRecords(records ...QualifiedRecord) *Update {
	b.withRecords(records)
	return b
}

// ToSQL converts the select statement into a string and returns its arguments.
func (b *Update) ToSQL() (string, []interface{}, error) {
	return b.buildArgsAndSQL(b)
}

func (b *Update) writeBuildCache(sql []byte) {
	b.BuilderConditional = BuilderConditional{}
	b.SetClausAliases = nil
	b.SetClauses = nil
	b.cachedSQL = sql
}

func (b *Update) readBuildCache() (sql []byte) {
	return b.cachedSQL
}

// DisableBuildCache if enabled it does not cache the SQL string as a final
// rendered byte slice. Allows you to rebuild the query with different
// statements.
func (b *Update) DisableBuildCache() *Update {
	b.IsBuildCacheDisabled = true
	return b
}

// ToSQL serialized the Update to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Update) toSQL(buf *bytes.Buffer, placeHolders []string) ([]string, error) {
	b.defaultQualifier = b.Table.qualifier()

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return nil, errors.WithStack(err)
	}

	if b.RawFullSQL != "" {
		buf.WriteString(b.RawFullSQL)
		return placeHolders, nil
	}

	if len(b.Table.Name) == 0 {
		return nil, errors.Empty.Newf("[dml] Update: Table at empty")
	}
	if len(b.SetClauses) == 0 {
		return nil, errors.Empty.Newf("[dml] Update: No columns specified")
	}

	buf.WriteString("UPDATE ")
	writeStmtID(buf, b.id)
	_, _ = b.Table.writeQuoted(buf, nil)
	buf.WriteString(" SET ")

	placeHolders, err := b.SetClauses.writeSetClauses(buf, placeHolders)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Write WHERE clause if we have any fragments
	placeHolders, err = b.Wheres.write(buf, 'w', placeHolders)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlWriteOrderBy(buf, b.OrderBys, false)
	sqlWriteLimitOffset(buf, b.LimitValid, b.LimitCount, false, 0)
	return placeHolders, nil
}

func (b *Update) validate() error {
	if len(b.cachedSQL) > 1 { // already validated
		return nil
	}
	if len(b.SetClauses) == 0 {
		return errors.Empty.Newf("[dml] Update: Columns are empty")
	}
	if len(b.SetClausAliases) > 0 && len(b.SetClausAliases) != len(b.SetClauses) {
		return errors.Mismatch.Newf("[dml] Update: ColumnAliases slice and Columns slice must have the same length")
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
// of the statement. The returned Stmter is not safe for concurrent use, despite
// the underlying *sql.Stmt is.
func (b *Update) Prepare(ctx context.Context) (Stmter, error) {
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
	return &stmtBase{
		builderCommon: builderCommon{
			id:               b.id,
			argsArgs:         make(Arguments, 0, cap),
			argsRaw:          make([]interface{}, 0, cap),
			defaultQualifier: b.Table.qualifier(),
			qualifiedColumns: b.qualifiedColumns,
			Log:              b.Log,
		},
		source: dmlTypeUpdate,
		stmt:   stmt,
	}, nil
}
