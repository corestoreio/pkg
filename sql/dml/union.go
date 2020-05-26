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
	"strconv"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// Union represents a UNION SQL statement. UNION is used to combine the result
// from multiple SELECT statements into a single result set.
// With template usage enabled, it builds multiple select statements joined by
// UNION and all based on a common template.
type Union struct {
	// templateStmtCount only used in case a UNION statement acts as a template.
	// Create one SELECT statement and by setting the data for
	// Union.StringReplace function additional SELECT statements are getting
	// created. Now the arguments must be multiplied by the number of new
	// created SELECT statements. This value gets stored in templateStmtCount.
	// An example exists in TestUnionTemplate_ReuseArgs.
	templateStmtCount int
	BuilderBase
	Selects     []*Select
	OrderBys    ids
	IsAll       bool // IsAll enables UNION ALL
	IsIntersect bool // See Intersect()
	IsExcept    bool // See Except()

	// When using Union as a template, only one *Select is required.
	oldNew [][]string // use for string replacement with `repls` field
	repls  []*strings.Replacer
}

// NewUnion creates a new Union object. If using as a template, only one *Select
// object can be provided.
func NewUnion(selects ...*Select) *Union {
	return &Union{
		Selects: selects,
	}
}

// Unsafe see BuilderBase.IsUnsafe which weakens security when building the SQL
// string. This function must be called before calling any other function.
func (u *Union) Unsafe() *Union {
	u.IsUnsafe = true
	return u
}

// Append adds more *Select objects to the Union object. When using Union as a
// template only one *Select object can be provided.
func (u *Union) Append(selects ...*Select) *Union {
	u.Selects = append(u.Selects, selects...)
	return u
}

// All returns all rows. The default behavior for UNION is that duplicate rows
// are removed from the result. Enabling ALL returns all rows.
func (u *Union) All() *Union {
	u.IsAll = true
	return u
}

// PreserveResultSet enables the correct ordering of the result set from the
// Select statements. UNION by default produces an unordered set of rows. To
// cause rows in a UNION result to consist of the sets of rows retrieved by each
// SELECT one after the other, select an additional column in each SELECT to use
// as a sort column and add an ORDER BY following the last SELECT.
func (u *Union) PreserveResultSet() *Union {
	if len(u.Selects) > 1 {
		for i, s := range u.Selects {
			s.AddColumnsConditions(Expr(strconv.Itoa(i)).Alias("_preserve_result_set"))
		}
		u.OrderBys = append(ids{MakeIdentifier("_preserve_result_set")}, u.OrderBys...)
		return u
	}

	// Panics without any *Select in the slice. Programmer error.
	u.Selects[0].AddColumnsConditions(Expr("{preserveResultSet}").Alias("_preserve_result_set"))
	u.OrderBys = append(ids{MakeIdentifier("_preserve_result_set")}, u.OrderBys...)
	for i := 0; i < u.templateStmtCount; i++ {
		u.oldNew[i] = append(u.oldNew[i], "{preserveResultSet}", strconv.Itoa(i))
	}
	return u
}

// OrderBy appends a column to ORDER the statement ascending. A column gets
// always quoted if it is a valid identifier otherwise it will be treated as an
// expression. MySQL might order the result set in a temporary table, which is
// slow. Under different conditions sorting can skip the temporary table.
// https://dev.mysql.com/doc/relnotes/mysql/5.7/en/news-5-7-3.html
// A column name can also contain the suffix words " ASC" or " DESC" to indicate
// the sorting. This avoids using the method OrderByDesc when sorting certain
// columns descending.
func (u *Union) OrderBy(columns ...string) *Union {
	u.OrderBys = u.OrderBys.AppendColumns(u.IsUnsafe, columns...)
	return u
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// A column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. When you use ORDER BY or GROUP BY to sort a column
// in a DELETE, the server sorts values using only the initial number of bytes
// indicated by the max_sort_length system variable.
func (u *Union) OrderByDesc(columns ...string) *Union {
	u.OrderBys = u.OrderBys.AppendColumns(u.IsUnsafe, columns...).applySort(len(columns), sortDescending)
	return u
}

// Intersect switches the query type from UNION to INTERSECT. The result of an
// intersect is the intersection of right and left SELECT results, i.e. only
// records that are present in both result sets will be included in the result
// of the operation. INTERSECT has higher precedence than UNION and EXCEPT. If
// possible it will be executed linearly but if not it will be translated to a
// subquery in the FROM clause. Only supported in MariaDB >=10.3
func (u *Union) Intersect() *Union {
	u.IsIntersect = true
	return u
}

// Except switches the query from UNION to EXCEPT. The result of EXCEPT is all
// records of the left SELECT result except records which are in right SELECT
// result set, i.e. it is subtraction of two result sets. EXCEPT and UNION have
// the same operation precedence. Only supported in MariaDB >=10.3
func (u *Union) Except() *Union {
	u.IsExcept = true
	return u
}

// StringReplace is only applicable when using *Union as a template.
// StringReplace replaces the `key` with one of the `values`. Each value defines
// a generated SELECT query. Repeating calls of StringReplace must provide the
// same amount of `values` as the first  or an index of bound stack trace
// happens. This function is just a simple string replacement. Make sure that
// your key does not match other parts of the SQL query.
func (u *Union) StringReplace(key string, values ...string) *Union {
	if len(u.Selects) > 1 {
		return u
	}
	if u.templateStmtCount == 0 {
		u.templateStmtCount = len(values)
		u.oldNew = make([][]string, u.templateStmtCount)
		u.repls = make([]*strings.Replacer, u.templateStmtCount)
	}
	for i := 0; i < u.templateStmtCount; i++ {
		// The following block has been put on each line because the (index out of
		// bound) stack trace will show exactly what you have made wrong =>
		// Providing in the 2nd call of StringReplace too few `values`
		// arguments.
		u.oldNew[i] = append(u.oldNew[i], key,
			values[i])
	}
	return u
}

// ToSQL converts the statements into a string and returns its arguments.
func (u *Union) ToSQL() (string, []interface{}, error) {
	rawSQL, err := u.buildToSQL(u)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	return rawSQL, nil, nil
}

// ToSQL generates the SQL string and its arguments. Calls to this function are
// idempotent.
func (u *Union) toSQL(w *bytes.Buffer, placeHolders []string) (_ []string, err error) {
	if len(u.Selects) > 1 {
		for i, s := range u.Selects {
			if i > 0 {
				sqlWriteUnionAll(w, u.IsAll, u.IsIntersect, u.IsExcept)
			}
			w.WriteByte('(')

			placeHolders, err = s.toSQL(w, placeHolders)
			if err != nil {
				return nil, errors.Wrapf(err, "[dml] Union.ToSQL at Select index %d", i)
			}
			w.WriteByte(')')
		}
		sqlWriteOrderBy(w, u.OrderBys, true)
		return placeHolders, nil
	}

	bufSel0 := bufferpool.Get()
	placeHolders, err = u.Selects[0].toSQL(bufSel0, placeHolders)
	selStr := bufSel0.String()
	bufferpool.Put(bufSel0)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := 0; i < u.templateStmtCount; i++ {
		repl := u.repls[i]
		if repl == nil {
			repl = strings.NewReplacer(u.oldNew[i]...)
			u.repls[i] = repl
		}
		if i > 0 {
			sqlWriteUnionAll(w, u.IsAll, u.IsIntersect, u.IsExcept)
		}
		w.WriteByte('(')
		repl.WriteString(w, selStr)
		w.WriteByte(')')
	}

	if w.Len() == 0 {
		return nil, errors.Empty.Newf("[dml] No SQL string generated. Number of select stmts: %d", len(u.Selects))
	}

	sqlWriteOrderBy(w, u.OrderBys, true)
	return placeHolders, nil
}

// Clone creates a clone of the current object, leaving fields DB and Log
// untouched. Additionally the fields for replacing strings also won't get
// copied.
func (u *Union) Clone() *Union {
	if u == nil {
		return nil
	}

	c := *u
	c.BuilderBase = u.BuilderBase.Clone()
	if ls := len(u.Selects); ls > 0 {
		c.Selects = make([]*Select, ls)
		for i, s := range u.Selects {
			c.Selects[i] = s.Clone()
		}
	}
	c.OrderBys = u.OrderBys.Clone()
	return &c
}
