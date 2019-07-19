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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Update contains the logic for an UPDATE statement.
// TODO: add UPDATE JOINS
type Update struct {
	BuilderBase
	BuilderConditional
	// SetClauses contains the column/argument association. For each column
	// there must be one argument.
	SetClauses Conditions
}

// NewUpdate creates a new Update object.
func NewUpdate(table string) *Update {
	return &Update{
		BuilderBase: BuilderBase{
			Table: MakeIdentifier(table),
		},
	}
}

func newUpdate(db QueryExecPreparer, cComm *connCommon, table string) *Update {
	id := cComm.makeUniqueID()
	l := cComm.Log
	table = cComm.mapTableName(table)
	if l != nil {
		l = l.With(log.String("update_id", id), log.String("table", table))
	}
	return &Update{
		BuilderBase: BuilderBase{
			builderCommon: builderCommon{
				id:  id,
				Log: l,
				DB:  db,
			},
			Table: MakeIdentifier(table),
		},
	}
}

// Update creates a new Update for the given table with a random connection from
// the pool.
func (c *ConnPool) Update(table string) *Update {
	return newUpdate(c.DB, &c.connCommon, table)
}

// Update creates a new Update for the given table bound to a single connection.
func (c *Conn) Update(table string) *Update {
	return newUpdate(c.DB, &c.connCommon, table)
}

// Update creates a new Update for the given table bound to a transaction.
func (tx *Tx) Update(table string) *Update {
	return newUpdate(tx.DB, &tx.connCommon, table)
}

// Alias sets an alias for the table name.
func (b *Update) Alias(alias string) *Update {
	b.Table.Aliased = alias
	return b
}

// WithDB sets the database query object.
func (b *Update) WithDB(db QueryExecPreparer) *Update {
	b.DB = db
	return b
}

// Unsafe see BuilderBase.IsUnsafe which weakens security when building the SQL
// string. This function must be called before calling any other function.
func (b *Update) Unsafe() *Update {
	b.IsUnsafe = true
	return b
}

// AddClauses appends a column/value pair for the statement.
func (b *Update) AddClauses(c ...*Condition) *Update {
	b.SetClauses = append(b.SetClauses, c...)
	return b
}

// AddColumns adds columns whose values gets later derived from a ColumnMapper.
// Those columns will get passed to the ColumnMapper implementation.
func (b *Update) AddColumns(columnNames ...string) *Update {
	for _, col := range columnNames {
		b.SetClauses = append(b.SetClauses, Column(col))
	}
	return b
}

// SetColumns resets the SetClauses slice and adds the columns. Same behaviour
// as AddColumns.
func (b *Update) SetColumns(columnNames ...string) *Update {
	b.SetClauses = b.SetClauses[:0]
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

// WithArgs returns a new type to support multiple executions of the underlying
// SQL statement and reuse of memory allocations for the arguments. WithArgs
// builds the SQL string in a thread safe way. It copies the underlying
// connection and settings from the current DML type (Delete, Insert, Select,
// Update, Union, With, etc.). The field DB can still be overwritten.
// Interpolation does not support the raw interfaces. It's an architecture bug
// to use WithArgs inside a loop. WithArgs does support thread safety and can be
// used in parallel. Each goroutine must have its own dedicated *Artisan
// pointer.
func (b *Update) WithArgs() *Artisan {
	return b.withArtisan(b)
}

// ToSQL converts the select statement into a string and returns its arguments.
func (b *Update) ToSQL() (string, []interface{}, error) {
	b.source = dmlSourceUpdate
	rawSQL, err := b.buildToSQL(b)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	return string(rawSQL), nil, nil
}

// WithCacheKey sets the currently used cache key when generating a SQL string.
// By setting a different cache key, a previous generated SQL query is
// accessible again. New cache keys allow to change the generated query of the
// current object. E.g. different where clauses or different row counts in
// INSERT ... VALUES statements. The empty string defines the default cache key.
// If the `args` argument contains values, then fmt.Sprintf gets used.
func (b *Update) WithCacheKey(key string, args ...interface{}) *Update {
	b.withCacheKey(key, args...)
	return b
}

// ToSQL serialized the Update to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Update) toSQL(buf *bytes.Buffer, placeHolders []string) ([]string, error) {
	b.defaultQualifier = b.Table.qualifier()
	b.source = dmlSourceUpdate

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
	sqlWriteLimitOffset(buf, b.LimitValid, false, 0, b.LimitCount)
	return placeHolders, nil
}

// Prepare executes the statement represented by the Update to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Update are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement. The returned Stmter is not safe for concurrent use, despite
// the underlying *sql.Stmt is.
func (b *Update) Prepare(ctx context.Context) (*Stmt, error) {
	return b.prepare(ctx, b.DB, b, dmlSourceUpdate)
}

// PrepareWithArgs same as Prepare but forwards the possible error of creating a
// prepared statement into the Artisan type. Reduces boilerplate code. You must
// call Artisan.Close to deallocate the prepared statement in the SQL server.
func (b *Update) PrepareWithArgs(ctx context.Context) *Artisan {
	stmt, err := b.prepare(ctx, b.DB, b, dmlSourceUpdate)
	if err != nil {
		a := &Artisan{
			base: builderCommon{
				ärgErr: errors.WithStack(err),
			},
		}
		return a
	}
	return stmt.WithArgs()
}

// Clone creates a clone of the current object, leaving fields DB and Log
// untouched.
func (b *Update) Clone() *Update {
	if b == nil {
		return nil
	}
	c := *b
	c.BuilderBase = b.BuilderBase.Clone()
	c.BuilderConditional = b.BuilderConditional.Clone()
	c.SetClauses = b.SetClauses.Clone()
	return &c
}
