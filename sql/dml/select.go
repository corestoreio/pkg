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
	"fmt"

	"github.com/corestoreio/errors"
)

// Select contains the clauses for a SELECT statement. Wildcard `SELECT *`
// statements are not really supported.
// http://stackoverflow.com/questions/3639861/why-is-select-considered-harmful
type Select struct {
	BuilderBase
	BuilderConditional

	// Columns represents a slice of names and its optional identifiers. Wildcard
	// `SELECT *` statements are not really supported:
	// http://stackoverflow.com/questions/3639861/why-is-select-considered-harmful
	Columns ids

	// TODO: create a possibility of the Select type which has a half-pre-rendered
	// SQL statement where a developer can only modify or append WHERE clauses.
	// especially useful during code generation

	GroupBys             ids
	Havings              Conditions
	IsStar               bool // IsStar generates a SELECT * FROM query
	IsCountStar          bool // IsCountStar retains the column names but executes a COUNT(*) query.
	IsDistinct           bool // See Distinct()
	IsStraightJoin       bool // See StraightJoin()
	IsSQLNoCache         bool // See SQLNoCache()
	IsForUpdate          bool // See ForUpdate()
	IsLockInShareMode    bool // See LockInShareMode()
	IsOrderByDeactivated bool // See OrderByDeactivated()
	IsOrderByRand        bool // enables the original slow ORDER BY RAND() clause
	OffsetCount          uint64
}

// NewSelect creates a new Select object.
func NewSelect(columns ...string) *Select {
	s := new(Select)
	if len(columns) == 1 && columns[0] == "*" {
		s.Star()
	} else {
		s.Columns = s.Columns.AppendColumns(false, columns...)
	}
	return s
}

// NewSelectWithDerivedTable creates a new derived table (Subquery in the FROM
// Clause) using the provided sub-select in the FROM part together with an alias
// name. Appends the arguments of the sub-select to the parent *Select pointer
// arguments list. SQL result may look like:
//		SELECT a,b FROM (SELECT x,y FROM `product` AS `p`) AS `t`
// https://dev.mysql.com/doc/refman/5.7/en/derived-tables.html
func NewSelectWithDerivedTable(subSelect *Select, aliasName string) *Select {
	return &Select{
		BuilderBase: BuilderBase{
			Table: id{
				DerivedTable: subSelect,
				Aliased:      aliasName,
			},
		},
	}
}

// Distinct marks the statement at a DISTINCT SELECT. It specifies removal of
// duplicate rows from the result set.
func (b *Select) Distinct() *Select {
	b.IsDistinct = true
	return b
}

// Unsafe see BuilderBase.IsUnsafe which weakens security when building the SQL
// string. This function must be called before calling any other function.
func (b *Select) Unsafe() *Select {
	b.IsUnsafe = true
	return b
}

// StraightJoin forces the optimizer to join the tables in the order in which
// they are listed in the FROM clause. You can use this to speed up a query if
// the optimizer joins the tables in nonoptimal order.
func (b *Select) StraightJoin() *Select {
	b.IsStraightJoin = true
	return b
}

// SQLNoCache tells the server that it does not use the query cache. It neither
// checks the query cache to see whether the result is already cached, nor does
// it cache the query result.
func (b *Select) SQLNoCache() *Select {
	b.IsSQLNoCache = true
	return b
}

// ForUpdate sets for index records the search encounters, locks the rows and
// any associated index entries, the same as if you issued an UPDATE statement
// for those rows. Other transactions are blocked from updating those rows, from
// doing SELECT ... LOCK IN SHARE MODE, or from reading the data in certain
// transaction isolation levels. Consistent reads ignore any locks set on the
// records that exist in the read view. (Old versions of a record cannot be
// locked; they are reconstructed by applying undo logs on an in-memory copy of
// the record.)
// Note: Locking of rows for update using SELECT FOR UPDATE only applies when
// autocommit is disabled (either by beginning transaction with START
// TRANSACTION or by setting autocommit to 0. If autocommit is enabled, the rows
// matching the specification are not locked.
// https://dev.mysql.com/doc/refman/5.5/en/innodb-locking-reads.html
func (b *Select) ForUpdate() *Select {
	b.IsForUpdate = true
	return b
}

// LockInShareMode sets a shared mode lock on any rows that are read. Other
// sessions can read the rows, but cannot modify them until your transaction
// commits. If any of these rows were changed by another transaction that has
// not yet committed, your query waits until that transaction ends and then uses
// the latest values.
// https://dev.mysql.com/doc/refman/5.5/en/innodb-locking-reads.html
func (b *Select) LockInShareMode() *Select {
	b.IsLockInShareMode = true
	return b
}

// Count executes a COUNT(*) as `counted` query without touching or changing the
// currently set columns.
func (b *Select) Count() *Select {
	b.IsCountStar = true
	return b
}

// Star creates a SELECT * FROM query. Such queries are discouraged from using.
func (b *Select) Star() *Select {
	b.IsStar = true
	return b
}

// From sets the table for the SELECT FROM part.
func (b *Select) From(from string) *Select {
	b.Table = MakeIdentifier(from)
	return b
}

// FromAlias sets the table and its alias name for a `SELECT ... FROM table AS
// alias` query.
func (b *Select) FromAlias(from, alias string) *Select {
	b.Table = MakeIdentifier(from).Alias(alias)
	return b
}

// AddColumns appends more columns to the Columns slice. If a column name is not
// valid identifier that column gets switched into an expression.
// 		AddColumns("a","b") 		// `a`,`b`
// 		AddColumns("a,b","z","c,d")	// a,b,`z`,c,d
//		AddColumns("t1.name","t1.sku","price") // `t1`.`name`, `t1`.`sku`,`price`
func (b *Select) AddColumns(cols ...string) *Select {
	b.Columns = b.Columns.AppendColumns(b.IsUnsafe, cols...)
	return b
}

// AddColumnsAliases expects a balanced slice of "Column1, Alias1, Column2,
// Alias2" and adds both to the Columns slice. An imbalanced slice will cause a
// panic. If a column name is not valid identifier that column gets switched
// into an expression.
//		AddColumnsAliases("t1.name","t1Name","t1.sku","t1SKU") // `t1`.`name` AS `t1Name`, `t1`.`sku` AS `t1SKU`
// 		AddColumnsAliases("(e.price*x.tax*t.weee)", "final_price") // error: `(e.price*x.tax*t.weee)` AS `final_price`
func (b *Select) AddColumnsAliases(columnAliases ...string) *Select {
	b.Columns = b.Columns.AppendColumnsAliases(b.IsUnsafe, columnAliases...)
	return b
}

// AddColumnsConditions adds a condition as a column to the statement. The
// operator field gets ignored. DBR in the condition gets applied to the
// RawArguments field to maintain the correct order of arguments.
// 		AddColumnsConditions(Expr("(e.price*x.tax*t.weee)").Alias("final_price")) // (e.price*x.tax*t.weee) AS `final_price`
func (b *Select) AddColumnsConditions(expressions ...*Condition) *Select {
	b.Columns, b.ärgErr = b.Columns.appendConditions(expressions)
	return b
}

// Where appends a WHERE clause to the statement for the given string and args
// or map of column/value pairs.
func (b *Select) Where(wf ...*Condition) *Select {
	b.Wheres = append(b.Wheres, wf...)
	return b
}

// When applies the function `fn` query changes if the given "test" is true.
// Providing the optional second function, uses it as the default value, if test
// is false. `defaultFn` can be nil.
func (b *Select) When(test bool, fn func(*Select), defaultFn func(*Select)) *Select {
	// TODO add this to other DML types
	switch {
	case test:
		fn(b) // test is true, applies callback
	case defaultFn != nil:
		defaultFn(b) // default value, if test is false
	}
	return b
}

// Unless applies the function `fn` query changes if the given "test" is false.
// Providing the optional second function, uses it as the default value, if test
// is false. `defaultFn` can be nil.
func (b *Select) Unless(test bool, fn func(*Select), defaultFn func(*Select)) *Select {
	// TODO add this to other DML types
	return b.When(!test, fn, defaultFn)
}

// GroupBy appends columns to group the statement. A column gets always quoted
// if it is a valid identifier otherwise it will be treated as an expression.
// MySQL does not sort the results set. To avoid the overhead of sorting that
// GROUP BY produces this function should add an ORDER BY NULL with function
// `OrderByDeactivated`.
func (b *Select) GroupBy(columns ...string) *Select {
	b.GroupBys = b.GroupBys.AppendColumns(b.IsUnsafe, columns...)
	return b
}

// GroupByAsc sorts the groups in ascending order. A column gets always quoted
// if it is a valid identifier otherwise it will be treated as an expression. No
// need to add an ORDER BY clause. When you use ORDER BY or GROUP BY to sort a
// column in a SELECT, the server sorts values using only the initial number of
// bytes indicated by the max_sort_length system variable.
func (b *Select) GroupByAsc(columns ...string) *Select {
	b.GroupBys = b.GroupBys.AppendColumns(b.IsUnsafe, columns...).applySort(len(columns), sortAscending)
	return b
}

// GroupByDesc sorts the groups in descending order. A column gets always quoted
// if it is a valid identifier otherwise it will be treated as an expression. No
// need to add an ORDER BY clause. When you use ORDER BY or GROUP BY to sort a
// column in a SELECT, the server sorts values using only the initial number of
// bytes indicated by the max_sort_length system variable.
func (b *Select) GroupByDesc(columns ...string) *Select {
	b.GroupBys = b.GroupBys.AppendColumns(b.IsUnsafe, columns...).applySort(len(columns), sortDescending)
	return b
}

// Having appends a HAVING clause to the statement
func (b *Select) Having(wf ...*Condition) *Select {
	b.Havings = append(b.Havings, wf...)
	return b
}

// OrderByDeactivated deactivates ordering of the result set by applying ORDER
// BY NULL to the SELECT statement. Very useful for GROUP BY queries.
func (b *Select) OrderByDeactivated() *Select {
	b.IsOrderByDeactivated = true
	return b
}

// OrderBy appends columns to the ORDER BY statement for ascending sorting. A
// column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. When you use ORDER BY or GROUP BY to sort a column
// in a SELECT, the server sorts values using only the initial number of bytes
// indicated by the max_sort_length system variable.
// A column name can also contain the suffix words " ASC" or " DESC" to indicate
// the sorting. This avoids using the method OrderByDesc when sorting certain
// columns descending.
func (b *Select) OrderBy(columns ...string) *Select {
	b.OrderBys = b.OrderBys.AppendColumns(b.IsUnsafe, columns...)
	return b
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// A column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. When you use ORDER BY or GROUP BY to sort a column
// in a SELECT, the server sorts values using only the initial number of bytes
// indicated by the max_sort_length system variable.
func (b *Select) OrderByDesc(columns ...string) *Select {
	b.OrderBys = b.OrderBys.AppendColumns(b.IsUnsafe, columns...).applySort(len(columns), sortDescending)
	return b
}

// OrderByRandom sorts the table randomly by not using ORDER BY RAND() rather
// using a JOIN with the single primary key column. This function overwrites
// previously set ORDER BY statements and the field LimitCount. The generated
// SQL by this function is about 3-4 times faster than ORDER BY RAND(). The
// generated SQL does not work for all queries. The underlying SQL statement
// might change without notice.
func (b *Select) OrderByRandom(idColumnName string, limit uint64) *Select {
	// Source https://stackoverflow.com/a/36013954 ;-)
	b.OrderByRandColumnName = idColumnName
	b.LimitCount = limit
	return b
}

// Limit sets a limit for the statement; overrides any existing LIMIT.
// Don't build a pagination with offset or you go straight to hell.
func (b *Select) Limit(offset uint64, limit uint64) *Select {
	b.OffsetCount = offset
	b.LimitCount = limit
	b.LimitValid = true
	return b
}

// Paginate sets LIMIT/OFFSET for the statement based on the given page/perPage
// Assumes page/perPage are valid. Page and perPage must be >= 1.
// Deprecated see a talk from Marcus Wienand - Modern SQL
func (b *Select) Paginate(page, perPage uint64) *Select {
	b.Limit((page-1)*perPage, perPage)
	return b
}

// Join creates an INNER join construct. By default, the onConditions are glued
// together with AND.
func (b *Select) Join(table id, onConditions ...*Condition) *Select {
	b.join("INNER", table, onConditions...)
	return b
}

// LeftJoin creates a LEFT join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) LeftJoin(table id, onConditions ...*Condition) *Select {
	b.join("LEFT", table, onConditions...)
	return b
}

// RightJoin creates a RIGHT join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) RightJoin(table id, onConditions ...*Condition) *Select {
	b.join("RIGHT", table, onConditions...)
	return b
}

// OuterJoin creates an OUTER join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) OuterJoin(table id, onConditions ...*Condition) *Select {
	b.join("OUTER", table, onConditions...)
	return b
}

// CrossJoin creates a CROSS join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) CrossJoin(table id, onConditions ...*Condition) *Select {
	b.join("CROSS", table, onConditions...)
	return b
}

// ToSQL generates the SQL string and might caches it internally, if not
// disabled.
func (b *Select) ToSQL() (string, []any, error) {
	rawSQL, err := b.buildToSQL(b)
	return rawSQL, nil, err
}

// ToSQL serialized the Select to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Select) toSQL(w *bytes.Buffer, placeHolders []string) (_placeHolders []string, err error) {
	if b.BuilderBase.ärgErr != nil {
		return nil, b.BuilderBase.ärgErr
	}
	if len(b.Columns) == 0 && !b.IsCountStar && !b.IsStar {
		return nil, errors.Empty.Newf("[dml] Select: no columns specified")
	}

	w.WriteString("SELECT ")
	if b.IsDistinct {
		w.WriteString("DISTINCT ")
	}
	if b.IsStraightJoin {
		w.WriteString("STRAIGHT_JOIN ")
	}
	if b.IsSQLNoCache {
		w.WriteString("SQL_NO_CACHE ")
	}

	switch {
	case b.IsStar:
		w.WriteByte('*')
	case b.IsCountStar:
		w.WriteString("COUNT(*) AS ")
		Quoter.quote(w, "counted")
	default:
		if placeHolders, err = b.Columns.writeQuoted(w, placeHolders); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if !b.Table.isEmpty() {
		w.WriteString(" FROM ")
		if placeHolders, err = b.Table.writeQuoted(w, placeHolders); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	joins := b.Joins
	if b.OrderByRandColumnName != "" {
		// This ORDER BY RAND() statement enables a 3-4 better processing in the
		// server.
		countSel := NewSelect().AddColumnsConditions(
			Expr(fmt.Sprintf("((%d / COUNT(*)) * 10)", b.LimitCount)),
		).From(b.Table.Name).Where(b.Wheres...)

		idSel := NewSelect(b.OrderByRandColumnName).From(b.Table.Name).
			Where(Expr("RAND()").Less().Sub(countSel)).
			Where(b.Wheres...).
			Limit(0, b.LimitCount)
		idSel.IsOrderByRand = true

		joins = append(joins, &join{
			Table: id{
				DerivedTable: idSel,
				Aliased:      "rand" + b.Table.Name,
			},
			On: Conditions{Columns(b.OrderByRandColumnName)},
		})
	}

	for _, f := range joins {
		w.WriteByte(' ')
		w.WriteString(f.JoinType)
		w.WriteString(" JOIN ")
		if placeHolders, err = f.Table.writeQuoted(w, placeHolders); err != nil {
			return nil, errors.WithStack(err)
		}
		if placeHolders, err = f.On.write(w, 'j', placeHolders, b.isWithDBR); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if placeHolders, err = b.Wheres.write(w, 'w', placeHolders, b.isWithDBR); err != nil {
		return nil, errors.WithStack(err)
	}

	if len(b.GroupBys) > 0 {
		w.WriteString(" GROUP BY ")
		for i, c := range b.GroupBys {
			if i > 0 {
				w.WriteString(", ")
			}
			if placeHolders, err = c.writeQuoted(w, placeHolders); err != nil {
				return nil, errors.WithStack(err)
			}
		}
	}

	if placeHolders, err = b.Havings.write(w, 'h', placeHolders, b.isWithDBR); err != nil {
		return nil, errors.WithStack(err)
	}

	switch {
	case b.IsOrderByDeactivated:
		w.WriteString(" ORDER BY NULL")
	case b.IsOrderByRand:
		w.WriteString(" ORDER BY RAND()")
	default:
		sqlWriteOrderBy(w, b.OrderBys, false)
	}

	sqlWriteLimitOffset(w, b.LimitValid, true, b.OffsetCount, b.LimitCount)

	switch {
	case b.IsLockInShareMode:
		w.WriteString(" LOCK IN SHARE MODE")
	case b.IsForUpdate:
		w.WriteString(" FOR UPDATE")
	}
	return placeHolders, err
}

// Clone creates a clone of the current object, leaving fields DB and Log
// untouched.
func (b *Select) Clone() *Select {
	if b == nil {
		return nil
	}
	c := *b
	c.BuilderBase = b.BuilderBase.Clone()
	c.BuilderConditional = b.BuilderConditional.Clone()
	c.Columns = b.Columns.Clone()
	c.GroupBys = b.GroupBys.Clone()
	c.Havings = b.Havings.Clone()
	return &c
}
