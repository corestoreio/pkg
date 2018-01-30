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
			builderCommon: builderCommon{
				id:  id,
				Log: withInitLog(c.Log, expressions, id),
				DB:  c.DB,
			},
		},
		Subclauses: expressions,
	}
}

// With creates a new With statement bound to a single connection.
func (c *Conn) With(expressions ...WithCTE) *With {
	id := c.makeUniqueID()
	return &With{
		BuilderBase: BuilderBase{
			builderCommon: builderCommon{
				id:  id,
				Log: withInitLog(c.Log, expressions, id),
				DB:  c.DB,
			},
		},
		Subclauses: expressions,
	}
}

// With creates a new With that select that given columns bound to the transaction
func (tx *Tx) With(expressions ...WithCTE) *With {
	id := tx.makeUniqueID()
	return &With{
		BuilderBase: BuilderBase{
			builderCommon: builderCommon{
				id:  id,
				Log: withInitLog(tx.Log, expressions, id),
				DB:  tx.DB,
			},
		},
		Subclauses: expressions,
	}
}

// WithDB sets the database query object.
func (b *With) WithDB(db QueryExecPreparer) *With {
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

// WithArgs builds the SQL string and sets the optional interfaced arguments for
// the later execution. It copies the underlying connection and structs.
func (b *With) WithArgs(args ...interface{}) *Arguments {
	b.source = dmlSourceWith
	return b.withArgs(b, args...)
}

// ToSQL converts the select statement into a string and returns its arguments.
func (b *With) ToSQL() (string, []interface{}, error) {
	b.source = dmlSourceWith
	rawSQL, err := b.buildToSQL(b)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	return string(rawSQL), nil, nil
}

func (b *With) writeBuildCache(sql []byte) {
	b.Subclauses = nil
	b.TopLevel.Select = nil
	b.TopLevel.Union = nil
	b.TopLevel.Update = nil
	b.TopLevel.Delete = nil
	b.cachedSQL = sql
}

func (b *With) readBuildCache() (sql []byte) {
	return b.cachedSQL
}

// DisableBuildCache if enabled it does not cache the SQL string as a final
// rendered byte slice. Allows you to rebuild the query with different
// statements.
func (b *With) DisableBuildCache() *With {
	b.IsBuildCacheDisabled = true
	return b
}

func (b *With) toSQL(w *bytes.Buffer, placeHolders []string) (_ []string, err error) {

	w.WriteString("WITH ")
	writeStmtID(w, b.id)
	if b.IsRecursive {
		w.WriteString("RECURSIVE ")
	}

	//for i, ls := 0, len(b.Subclauses); i < ls && err == nil; i++ {
	//	sc := b.Subclauses[i]
	//	}
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
			sc.Select.IsBuildCacheDisabled = b.IsBuildCacheDisabled
			placeHolders, err = sc.Select.toSQL(w, placeHolders)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		case sc.Union != nil:
			sc.Union.IsBuildCacheDisabled = b.IsBuildCacheDisabled
			placeHolders, err = sc.Union.toSQL(w, placeHolders)
			if err != nil {
				return nil, errors.WithStack(err)
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
		b.TopLevel.Select.IsBuildCacheDisabled = b.IsBuildCacheDisabled
		placeHolders, err = b.TopLevel.Select.toSQL(w, placeHolders)
		return placeHolders, errors.WithStack(err)

	case b.TopLevel.Union != nil:
		b.TopLevel.Union.IsBuildCacheDisabled = b.IsBuildCacheDisabled
		placeHolders, err = b.TopLevel.Union.toSQL(w, placeHolders)
		return placeHolders, errors.WithStack(err)

	case b.TopLevel.Update != nil:
		b.TopLevel.Update.IsBuildCacheDisabled = b.IsBuildCacheDisabled
		placeHolders, err = b.TopLevel.Update.toSQL(w, placeHolders)
		return placeHolders, errors.WithStack(err)

	case b.TopLevel.Delete != nil:
		b.TopLevel.Delete.IsBuildCacheDisabled = b.IsBuildCacheDisabled
		placeHolders, err = b.TopLevel.Delete.toSQL(w, placeHolders)
		return placeHolders, errors.WithStack(err)
	}
	return nil, errors.Empty.Newf("[dml] Type With misses a top level statement")
}

// Prepare executes the statement represented by the `With` to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or recs in the `With` are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement. The returned Stmter is not safe for concurrent use, despite
// the underlying *sql.Stmt is.
func (b *With) Prepare(ctx context.Context) (*Stmt, error) {
	return b.prepare(ctx, b.DB, b, dmlSourceWith)
}
