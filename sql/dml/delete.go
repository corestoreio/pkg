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
	BuilderBase
	BuilderConditional
	// MultiTables specifies the additional tables to delete from. Use function
	// `FromTables` to conveniently set it.
	MultiTables ids
	// Returning allows from MariaDB 10.0.5, it is possible to return a
	// resultset of the deleted rows for a single table to the client by using
	// the syntax DELETE ... RETURNING select_expr [, select_expr2 ...]] Any of
	// SQL expression that can be calculated from a single row fields is
	// allowed. Subqueries are allowed. The AS keyword is allowed, so it is
	// possible to use aliases. The use of aggregate functions is not allowed.
	// RETURNING cannot be used in multi-table DELETEs.
	Returning *Select
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

func newDeleteFrom(db QueryExecPreparer, cCom *connCommon, from string) *Delete {
	id := cCom.makeUniqueID()
	l := cCom.Log
	from = cCom.mapTableName(from)
	if l != nil {
		l = l.With(log.String("delete_id", id), log.String("table", from))
	}
	return &Delete{
		BuilderBase: BuilderBase{
			builderCommon: builderCommon{
				id:  id,
				Log: l,
				DB:  db,
			},
			Table: MakeIdentifier(from),
		},
		BuilderConditional: BuilderConditional{
			Wheres: make(Conditions, 0, 2),
		},
	}
}

// DeleteFrom creates a new Delete for the given table. Mapping the table name
// is supported.
func (c *ConnPool) DeleteFrom(from string) *Delete {
	return newDeleteFrom(c.DB, &c.connCommon, from)
}

// DeleteFrom creates a new Delete for the given table in the context for a
// single database connection. Mapping the table name is supported.
func (c *Conn) DeleteFrom(from string) *Delete {
	return newDeleteFrom(c.DB, &c.connCommon, from)
}

// DeleteFrom creates a new Delete for the given table in the context for a
// transaction. Mapping the table name is supported.
func (tx *Tx) DeleteFrom(from string) *Delete {
	return newDeleteFrom(tx.DB, &tx.connCommon, from)
}

// FromTables specifies additional tables to delete from besides the default table.
func (b *Delete) FromTables(tables ...string) *Delete {
	// DELETE [LOW_PRIORITY] [QUICK] [IGNORE]
	// tbl_name[.*] [, tbl_name[.*]] ...	<-- MultiTables/FromTables
	// FROM table_references
	//[WHERE where_condition]
	for _, t := range tables {
		b.MultiTables = append(b.MultiTables, MakeIdentifier(t))
	}
	return b
}

// Join creates an INNER join construct. By default, the onConditions are glued
// together with AND. Same Source and Target Table: Until MariaDB 10.3.1,
// deleting from a table with the same source and target was not possible. From
// MariaDB 10.3.1, this is now possible. For example:
//		DELETE FROM t1 WHERE c1 IN (SELECT b.c1 FROM t1 b WHERE b.c2=0);
func (b *Delete) Join(table id, onConditions ...*Condition) *Delete {
	b.join("INNER", table, onConditions...)
	return b
}

// LeftJoin creates a LEFT join construct. By default, the onConditions are
// glued together with AND.
func (b *Delete) LeftJoin(table id, onConditions ...*Condition) *Delete {
	b.join("LEFT", table, onConditions...)
	return b
}

// RightJoin creates a RIGHT join construct. By default, the onConditions are
// glued together with AND.
func (b *Delete) RightJoin(table id, onConditions ...*Condition) *Delete {
	b.join("RIGHT", table, onConditions...)
	return b
}

// OuterJoin creates an OUTER join construct. By default, the onConditions are
// glued together with AND.
func (b *Delete) OuterJoin(table id, onConditions ...*Condition) *Delete {
	b.join("OUTER", table, onConditions...)
	return b
}

// CrossJoin creates a CROSS join construct. By default, the onConditions are
// glued together with AND.
func (b *Delete) CrossJoin(table id, onConditions ...*Condition) *Delete {
	b.join("CROSS", table, onConditions...)
	return b
}

// Alias sets an alias for the table name.
func (b *Delete) Alias(alias string) *Delete {
	b.Table.Aliased = alias
	return b
}

// WithDB sets the database query object. DB can be either a *sql.DB (connection
// pool), a *sql.Conn (a single dedicated database session) or a *sql.Tx (an
// in-progress database transaction).
func (b *Delete) WithDB(db QueryExecPreparer) *Delete {
	b.DB = db
	return b
}

// Unsafe see BuilderBase.IsUnsafe which weakens security when building the SQL
// string. This function must be called before calling any other function.
func (b *Delete) Unsafe() *Delete {
	b.IsUnsafe = true
	return b
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
// A column name can also contain the suffix words " ASC" or " DESC" to indicate
// the sorting. This avoids using the method OrderByDesc when sorting certain
// columns descending.
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

// WithArgs returns a new Artisan type to support multiple executions of the
// underlying SQL statement and reuse of memory allocations for the arguments.
// WithArgs builds the SQL string in a thread safe way. It copies the underlying
// connection and settings from the current DML type (Delete, Insert, Select,
// Update, Union, With, etc.). The field DB can still be overwritten.
// Interpolation does not support the raw interfaces. It's an architecture bug
// to use WithArgs inside a loop. WithArgs does support thread safety and can be
// used in parallel. Each goroutine must have its own dedicated *Artisan
// pointer.
func (b *Delete) WithArgs() *Artisan {
	return b.withArtisan(b)
}

// ToSQL generates the SQL string and might caches it internally, if not
// disabled. The returned interface slice is always nil.
func (b *Delete) ToSQL() (string, []interface{}, error) {
	b.source = dmlSourceDelete
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
func (b *Delete) WithCacheKey(key string, args ...interface{}) *Delete {
	b.withCacheKey(key, args...)
	return b
}

// ToSQL serialized the Delete to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Delete) toSQL(w *bytes.Buffer, placeHolders []string) (_ []string, err error) {
	b.source = dmlSourceDelete
	b.defaultQualifier = b.Table.qualifier()

	if b.Table.Name == "" {
		return nil, errors.Empty.Newf("[dml] Delete: Table is missing")
	}

	w.WriteString("DELETE ")
	writeStmtID(w, b.id)

	for i, mt := range b.MultiTables {
		if i == 0 {
			if b.Table.Aliased != "" {
				Quoter.WriteIdentifier(w, b.Table.Aliased)
			} else {
				Quoter.WriteIdentifier(w, b.Table.Name)
			}
			w.WriteByte(',')
		}
		if i > 0 {
			w.WriteByte(',')
		}
		placeHolders, err = mt.writeQuoted(w, placeHolders)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if len(b.MultiTables) > 0 {
		w.WriteByte(' ')
		if b.Returning != nil {
			return nil, errors.NotAllowed.Newf("[dml] MariaDB does not support RETURNING in multi-table DELETEs")
		}
	}

	w.WriteString("FROM ")
	placeHolders, err = b.Table.writeQuoted(w, placeHolders)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, f := range b.Joins {
		w.WriteByte(' ')
		w.WriteString(f.JoinType)
		w.WriteString(" JOIN ")
		if placeHolders, err = f.Table.writeQuoted(w, placeHolders); err != nil {
			return nil, errors.WithStack(err)
		}
		if placeHolders, err = f.On.write(w, 'j', placeHolders); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	placeHolders, err = b.Wheres.write(w, 'w', placeHolders)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlWriteOrderBy(w, b.OrderBys, false)
	sqlWriteLimitOffset(w, b.LimitValid, false, 0, b.LimitCount)

	if b.Returning != nil {
		w.WriteString(" RETURNING ")
		placeHolders, err = b.Returning.toSQL(w, placeHolders)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return placeHolders, nil
}

// Prepare executes the statement represented by the Delete to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Delete are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement. If debug mode for logging has been enabled it logs the
// duration taken and the SQL string. The returned Stmter is not safe for
// concurrent use, despite the underlying *sql.Stmt is.
func (b *Delete) Prepare(ctx context.Context) (*Stmt, error) {
	return b.prepare(ctx, b.DB, b, dmlSourceDelete)
}

// PrepareWithArgs same as Prepare but forwards the possible error of creating a
// prepared statement into the Artisan type. Reduces boilerplate code. You must
// call Artisan.Close to deallocate the prepared statement in the SQL server.
func (b *Delete) PrepareWithArgs(ctx context.Context) *Artisan {
	stmt, err := b.prepare(ctx, b.DB, b, dmlSourceDelete)
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
func (b *Delete) Clone() *Delete {
	if b == nil {
		return nil
	}
	c := *b
	c.BuilderBase = b.BuilderBase.Clone()
	c.BuilderConditional = b.BuilderConditional.Clone()
	c.MultiTables = b.MultiTables.Clone()
	c.Returning = b.Returning.Clone()
	return &c
}
