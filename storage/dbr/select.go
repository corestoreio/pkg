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

	"github.com/corestoreio/errors"
)

// Select contains the clauses for a SELECT statement. Wildcard `SELECT *`
// statements are not really supported.
// http://stackoverflow.com/questions/3639861/why-is-select-considered-harmful
type Select struct {
	BuilderBase
	BuilderConditional
	// DB can be either a *sql.DB (connection pool), a *sql.Conn (a single
	// dedicated database session) or a *sql.Tx (an in-progress database
	// transaction).
	DB QueryPreparer

	// Columns represents a slice of names and its optional identifiers. Wildcard
	// `SELECT *` statements are not really supported:
	// http://stackoverflow.com/questions/3639861/why-is-select-considered-harmful
	Columns identifiers

	//TODO: create a possibility of the Select type which has a half-pre-rendered
	// SQL statement where a developer can only modify or append WHERE clauses.
	// especially useful during code generation

	GroupBys             identifiers
	Havings              Conditions
	IsStar               bool // IsStar generates a SELECT * FROM query
	IsCountStar          bool // IsCountStar retains the column names but executes a COUNT(*) query.
	IsDistinct           bool // See Distinct()
	IsStraightJoin       bool // See StraightJoin()
	IsSQLNoCache         bool // See SQLNoCache()
	IsForUpdate          bool // See ForUpdate()
	IsLockInShareMode    bool // See LockInShareMode()
	IsOrderByDeactivated bool // See OrderByDeactivated()
	OffsetValid          bool
	OffsetCount          uint64
	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners SelectListeners
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
			Table: identifier{
				DerivedTable: subSelect,
				Aliased:      aliasName,
			},
		},
	}
}

// Select creates a new Select which selects from the provided columns.
// Columns won't get quoted.
func (c *Connection) Select(columns ...string) *Select {
	s := &Select{}
	s.BuilderBase.Log = c.Log
	if len(columns) == 1 && columns[0] == "*" {
		s.Star()
	} else {
		s.Columns = s.Columns.AppendColumns(false, columns...)
	}
	s.DB = c.DB
	return s
}

// SelectBySQL creates a new Select for the given SQL string and arguments.
func (c *Connection) SelectBySQL(sql string, args Arguments) *Select {
	s := &Select{
		BuilderBase: BuilderBase{
			Log:          c.Log,
			RawFullSQL:   sql,
			RawArguments: args,
		},
	}
	s.DB = c.DB
	return s
}

// Select creates a new Select that select that given columns bound to the
// transaction.
func (tx *Tx) Select(columns ...string) *Select {
	s := &Select{}
	s.BuilderBase.Log = tx.Logger
	if len(columns) == 1 && columns[0] == "*" {
		s.Star()
	} else {
		s.Columns = s.Columns.AppendColumns(false, columns...)
	}
	s.DB = tx.Tx
	return s
}

// SelectBySQL creates a new Select for the given SQL string and arguments bound
// to the transaction.
func (tx *Tx) SelectBySQL(sql string, args Arguments) *Select {
	s := &Select{
		BuilderBase: BuilderBase{
			Log:          tx.Logger,
			RawFullSQL:   sql,
			RawArguments: args,
		},
	}
	s.DB = tx.Tx
	return s
}

// WithDB sets the database query object.
func (b *Select) WithDB(db QueryPreparer) *Select {
	b.DB = db
	return b
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
// operator field gets ignored. Arguments in the condition gets applied to the
// RawArguments field to maintain the correct order of arguments.
// 		AddColumnsConditions(Expr("(e.price*x.tax*t.weee)").Alias("final_price")) // (e.price*x.tax*t.weee) AS `final_price`
func (b *Select) AddColumnsConditions(expressions ...*Condition) *Select {
	b.Columns, b.RawArguments = b.Columns.AppendConditions(expressions, b.RawArguments)
	return b
}

// BindRecord binds the qualified record to the main table/view, or any other
// table/view/alias used in the query, for assembling and appending arguments.
// An ArgumentsAppender gets called if it matches the qualifier, in this case
// the current table name or its alias.
func (b *Select) BindRecord(records ...QualifiedRecord) *Select {
	b.bindRecord(records...)
	return b
}

func (b *Select) bindRecord(records ...QualifiedRecord) {
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

// Where appends a WHERE clause to the statement for the given string and args
// or map of column/value pairs.
func (b *Select) Where(wf ...*Condition) *Select {
	b.Wheres = append(b.Wheres, wf...)
	return b
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

// Limit sets a limit for the statement; overrides any existing LIMIT
func (b *Select) Limit(limit uint64) *Select {
	b.LimitCount = limit
	b.LimitValid = true
	return b
}

// Offset sets an offset for the statement; overrides any existing OFFSET
func (b *Select) Offset(offset uint64) *Select {
	b.OffsetCount = offset
	b.OffsetValid = true
	return b
}

// Paginate sets LIMIT/OFFSET for the statement based on the given page/perPage
// Assumes page/perPage are valid. Page and perPage must be >= 1
func (b *Select) Paginate(page, perPage uint64) *Select {
	b.Limit(perPage)
	b.Offset((page - 1) * perPage)
	return b
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `args` will then be nil.
func (b *Select) Interpolate() *Select {
	b.IsInterpolate = true
	return b
}

// Join creates an INNER join construct. By default, the onConditions are glued
// together with AND.
func (b *Select) Join(table identifier, onConditions ...*Condition) *Select {
	b.join("INNER", table, onConditions...)
	return b
}

// LeftJoin creates a LEFT join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) LeftJoin(table identifier, onConditions ...*Condition) *Select {
	b.join("LEFT", table, onConditions...)
	return b
}

// RightJoin creates a RIGHT join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) RightJoin(table identifier, onConditions ...*Condition) *Select {
	b.join("RIGHT", table, onConditions...)
	return b
}

// OuterJoin creates an OUTER join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) OuterJoin(table identifier, onConditions ...*Condition) *Select {
	b.join("OUTER", table, onConditions...)
	return b
}

// CrossJoin creates a CROSS join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) CrossJoin(table identifier, onConditions ...*Condition) *Select {
	b.join("CROSS", table, onConditions...)
	return b
}

// ToSQL converts the select statement into a string and returns its arguments.
func (b *Select) ToSQL() (string, []interface{}, error) {
	return toSQL(b, b.IsInterpolate, _isNotPrepared)
}

// argumentCapacity returns the total possible guessed size of a new args
// slice. Use as the cap parameter in a call to `make`.
func (b *Select) argumentCapacity() int {
	return len(b.RawArguments) + (len(b.Joins)+len(b.Wheres))*2
}

func (b *Select) writeBuildCache(sql []byte) {
	b.cacheSQL = sql
}

func (b *Select) readBuildCache() (sql []byte, _ Arguments, err error) {
	if b.cacheSQL == nil {
		return nil, nil, nil
	}
	b.cacheArgs, err = b.appendArgs(b.cacheArgs[:0])
	return b.cacheSQL, b.cacheArgs, err
}

// BuildCache enables that the final SQL string including place holders will be
// cached in a private field. Each time a call to function ToSQL happens, the
// arguments will be re-evaluated and returned or interpolated together with the
// SQL string.
func (b *Select) BuildCache() *Select {
	b.IsBuildCache = true
	return b
}

func (b *Select) hasBuildCache() bool {
	return b.IsBuildCache
}

// ToSQL serialized the Select to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Select) toSQL(w *bytes.Buffer) error {
	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return errors.WithStack(err)
	}

	if b.RawFullSQL != "" {
		_, err := w.WriteString(b.RawFullSQL)
		return err
	}

	if len(b.Columns) == 0 && !b.IsCountStar && !b.IsStar {
		return errors.NewEmptyf("[dbr] Select: no columns specified")
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
		if err := b.Columns.WriteQuoted(w); err != nil {
			return errors.WithStack(err)
		}
	}

	if !b.Table.isEmpty() {
		w.WriteString(" FROM ")
		if err := b.Table.WriteQuoted(w); err != nil {
			return errors.WithStack(err)
		}
	}

	for _, f := range b.Joins {
		w.WriteByte(' ')
		w.WriteString(f.JoinType)
		w.WriteString(" JOIN ")
		f.Table.WriteQuoted(w)
		if err := f.On.write(w, 'j'); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := b.Wheres.write(w, 'w'); err != nil {
		return errors.WithStack(err)
	}

	if len(b.GroupBys) > 0 {
		w.WriteString(" GROUP BY ")
		for i, c := range b.GroupBys {
			if i > 0 {
				w.WriteString(", ")
			}
			if err := c.WriteQuoted(w); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	if err := b.Havings.write(w, 'h'); err != nil {
		return errors.WithStack(err)
	}

	if b.IsOrderByDeactivated {
		w.WriteString(" ORDER BY NULL")
	} else {
		sqlWriteOrderBy(w, b.OrderBys, false)
	}

	sqlWriteLimitOffset(w, b.LimitValid, b.LimitCount, b.OffsetValid, b.OffsetCount)
	switch {
	case b.IsLockInShareMode:
		w.WriteString(" LOCK IN SHARE MODE")
	case b.IsForUpdate:
		w.WriteString(" FOR UPDATE")
	}
	return nil
}

// ToSQL serialized the Select to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Select) appendArgs(args Arguments) (_ Arguments, err error) {
	if b.RawFullSQL != "" {
		return b.RawArguments, nil
	}

	// not sure if copying is necessary but leaves at least b.args in pristine
	// condition
	if cap(args) == 0 {
		args = make(Arguments, 0, b.argumentCapacity())
	}
	args = append(args, b.RawArguments...)

	if args, err = b.Columns.appendArgs(args); err != nil {
		return nil, errors.WithStack(err)
	}

	if args, err = b.Table.appendArgs(args); err != nil {
		return nil, errors.WithStack(err)
	}

	placeHolderColumns := make([]string, 0, len(b.Joins)+len(b.Wheres)+len(b.Havings))
	var pap []int
	if len(b.Joins) > 0 {
		for _, f := range b.Joins {
			args, err = f.Table.appendArgs(args)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			if args, pap, err = f.On.appendArgs(args, appendArgsJOIN); err != nil {
				return nil, errors.WithStack(err)
			}
			// TODO: think about caching all calls to intersectConditions
			if boundCols := f.On.intersectConditions(placeHolderColumns); len(boundCols) > 0 {
				defaultQualifier := b.Table.mustQualifier()
				if args, err = appendArgs(pap, b.ArgumentsAppender, args, defaultQualifier, boundCols); err != nil {
					return nil, errors.WithStack(err)
				}
			}
		}
		placeHolderColumns = placeHolderColumns[:0]
	}

	if args, pap, err = b.Wheres.appendArgs(args, appendArgsWHERE); err != nil {
		return nil, errors.WithStack(err)
	}
	if boundCols := b.Wheres.intersectConditions(placeHolderColumns); len(boundCols) > 0 {
		defaultQualifier := b.Table.mustQualifier()
		if args, err = appendArgs(pap, b.ArgumentsAppender, args, defaultQualifier, boundCols); err != nil {
			return nil, errors.WithStack(err)
		}
		placeHolderColumns = placeHolderColumns[:0]
	}

	if args, pap, err = b.Havings.appendArgs(args, appendArgsHAVING); err != nil {
		return nil, errors.WithStack(err)
	}
	if boundCols := b.Havings.intersectConditions(placeHolderColumns); len(boundCols) > 0 {
		defaultQualifier := b.Table.mustQualifier()
		if args, err = appendArgs(pap, b.ArgumentsAppender, args, defaultQualifier, boundCols); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, nil
}
