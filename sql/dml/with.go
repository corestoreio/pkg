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

	"github.com/corestoreio/errors"
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

// Clone creates a cloned object of the current one.
func (cte WithCTE) Clone() WithCTE {
	cte.Columns = cloneStringSlice(cte.Columns)
	cte.Select = cte.Select.Clone()
	cte.Union = cte.Union.Clone()
	return cte
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
// http://dankleiman.com/2018/02/06/3-ways-to-level-up-your-sql-as-a-software-engineer/
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

// ToSQL converts the select statement into a string and returns its arguments.
func (b *With) ToSQL() (string, []interface{}, error) {
	rawSQL, err := b.buildToSQL(b)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	return rawSQL, nil, nil
}

func (b *With) toSQL(w *bytes.Buffer, placeHolders []string) (_ []string, err error) {
	w.WriteString("WITH ")
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
			placeHolders, err = sc.Select.toSQL(w, placeHolders)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		case sc.Union != nil:
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
		placeHolders, err = b.TopLevel.Select.toSQL(w, placeHolders)
		return placeHolders, errors.WithStack(err)

	case b.TopLevel.Union != nil:
		placeHolders, err = b.TopLevel.Union.toSQL(w, placeHolders)
		return placeHolders, errors.WithStack(err)

	case b.TopLevel.Update != nil:
		placeHolders, err = b.TopLevel.Update.toSQL(w, placeHolders)
		return placeHolders, errors.WithStack(err)

	case b.TopLevel.Delete != nil:
		placeHolders, err = b.TopLevel.Delete.toSQL(w, placeHolders)
		return placeHolders, errors.WithStack(err)
	}
	return nil, errors.Empty.Newf("[dml] Type With misses a top level statement")
}

// Clone creates a clone of the current object, leaving fields DB and Log
// untouched.
func (b *With) Clone() *With {
	if b == nil {
		return nil
	}

	c := *b
	c.BuilderBase = b.BuilderBase.Clone()
	if ls := len(b.Subclauses); ls > 0 {
		c.Subclauses = make([]WithCTE, ls)
		for i, s := range b.Subclauses {
			c.Subclauses[i] = s.Clone()
		}
	}
	c.TopLevel.Select = b.TopLevel.Select.Clone()
	c.TopLevel.Union = b.TopLevel.Union.Clone()
	c.TopLevel.Update = b.TopLevel.Update.Clone()
	c.TopLevel.Delete = b.TopLevel.Delete.Clone()
	return &c
}
