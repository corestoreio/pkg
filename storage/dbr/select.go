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
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Select contains the clauses for a SELECT statement. Wildcard `SELECT *`
// statements are not really supported.
// http://stackoverflow.com/questions/3639861/why-is-select-considered-harmful
type Select struct {
	// ID of the SELECT statement. Used in logging and during performance
	// monitoring. If empty the generated SQL string gets used which can might
	// contain sensitive information which should not get logged.
	// TODO implement
	ID  string
	Log log.Logger // Log optional logger
	// DB gets required once the Load*() functions will be used.
	DB Querier

	RawFullSQL   string
	RawArguments Arguments // Arguments used by RawFullSQL

	// Record if set retrieves the necessary arguments from the interface.
	Record ArgumentAssembler

	// Columns represents a slice of names and its optional aliases. Wildcard
	// `SELECT *` statements are not really supported:
	// http://stackoverflow.com/questions/3639861/why-is-select-considered-harmful
	Columns aliases

	//TODO: create a possibility of the Select type which has a half-pre-rendered
	// SQL statement where a developer can only modify or append WHERE clauses.
	// especially useful during code generation

	// Table table name and optional alias name to SELECT from.
	Table alias

	WhereFragments    WhereFragments
	JoinFragments     JoinFragments
	GroupBys          aliases
	HavingFragments   WhereFragments
	OrderBys          aliases
	LimitCount        uint64
	OffsetCount       uint64
	LimitValid        bool
	OffsetValid       bool
	IsDistinct        bool // See Distinct()
	IsStraightJoin    bool // See StraightJoin()
	IsSQLNoCache      bool // See SQLNoCache()
	IsForUpdate       bool // See ForUpdate()
	IsLockInShareMode bool // See LockInShareMode()
	IsInterpolate     bool // See Interpolate()
	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	// UseBuildCache if `true` the final build query including place holders
	// will be cached in a private field. Each time a call to function ToSQL
	// happens, the arguments will be re-evaluated and returned or interpolated.
	UseBuildCache bool
	cacheSQL      []byte
	cacheArgs     Arguments // like a buffer, gets reused
	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners SelectListeners
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
	// previousError any error occurred during construction the SQL statement
	previousError error
}

// NewSelect creates a new Select object.
func NewSelect(columns ...string) *Select {
	s := new(Select)
	s.Columns = s.Columns.appendColumns(columns, false)
	return s
}

// NewSelectFromSub creates a new derived table (Subquery in the FROM Clause)
// using the provided sub-select in the FROM part together with an alias name.
// Appends the arguments of the sub-select to the parent *Select pointer
// arguments list. SQL result may look like:
//		SELECT a,b FROM (SELECT x,y FROM `product` AS `p`) AS `t`
// https://dev.mysql.com/doc/refman/5.7/en/derived-tables.html
func NewSelectFromSub(subSelect *Select, aliasName string) *Select {
	s := &Select{
		Table: alias{
			Select: subSelect,
			Alias:  aliasName,
		},
	}
	return s
}

// Select creates a new Select which selects from the provided columns.
// Columns won't get quoted.
func (c *Connection) Select(columns ...string) *Select {
	s := &Select{
		Log: c.Log,
	}
	s.Columns = s.Columns.appendColumns(columns, false)
	s.DB = c.DB
	return s
}

// SelectBySQL creates a new Select for the given SQL string and arguments
func (c *Connection) SelectBySQL(sql string, args ...Argument) *Select {
	s := &Select{
		Log:          c.Log,
		RawFullSQL:   sql,
		RawArguments: args,
	}
	s.DB = c.DB
	return s
}

// Select creates a new Select that select that given columns bound to the transaction
func (tx *Tx) Select(columns ...string) *Select {
	s := &Select{
		Log: tx.Logger,
	}
	s.Columns = s.Columns.appendColumns(columns, false)
	s.DB = tx.Tx
	return s
}

// SelectBySQL creates a new Select for the given SQL string and arguments bound to the transaction
func (tx *Tx) SelectBySQL(sql string, args ...Argument) *Select {
	s := &Select{
		Log:          tx.Logger,
		RawFullSQL:   sql,
		RawArguments: args,
	}
	s.DB = tx.Tx
	return s
}

// WithDB sets the database query object.
func (b *Select) WithDB(db Querier) *Select {
	b.DB = db
	return b
}

// Distinct marks the statement at a DISTINCT SELECT. It specifies removal of
// duplicate rows from the result set.
func (b *Select) Distinct() *Select {
	b.IsDistinct = true
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

// Count resets the columns to COUNT(*) as `counted`
func (b *Select) Count() *Select {
	b.Columns = aliases{
		MakeAliasExpr("COUNT(*)", "counted"),
	}
	return b
}

// From sets the table to SELECT FROM. If second argument will be provided this
// at then considered at the alias. SELECT ... FROM table AS alias.
func (b *Select) From(from ...string) *Select {
	b.Table = MakeAlias(from...)
	return b
}

// AddColumns appends more columns to the Columns slice. If a single string gets
// passed with comma separated values, this string gets split by the comma and
// its values appended to the Columns slice. Columns won't get quoted.
// 		AddColumns("a","b") 		// `a`,`b`
// 		AddColumns("a,b","z","c,d")	// `a,b`,`z`,`c,d` <- invalid SQL!
//		AddColumns("t1.name","t1.sku","price") // `t1`.`name`, `t1`.`sku`,`price`
func (b *Select) AddColumns(cols ...string) *Select {
	b.Columns = b.Columns.appendColumns(cols, false)
	return b
}

// AddColumnsAlias expects a balanced slice of "Column1, Alias1, Column2,
// Alias2" and adds both to the Columns slice.
//		AddColumnsAlias("t1.name","t1Name","t1.sku","t1SKU") // `t1`.`name` AS `t1Name`, `t1`.`sku` AS `t1SKU`
// 		AddColumnsAlias("(e.price*x.tax*t.weee)", "final_price") // `(e.price*x.tax*t.weee)` AS `final_price`
func (b *Select) AddColumnsAlias(columnAliases ...string) *Select {
	if (len(columnAliases) % 2) == 1 {
		b.previousError = errors.NewMismatchf("[dbr] Expecting a balanced slice! Got: %v", columnAliases)
	} else {
		b.Columns = b.Columns.appendColumnsAliases(columnAliases, false)
	}
	return b
}

// AddColumnsExprAlias expects a balanced slice of "expression, AliasName" and
// adds both concatenated and quoted to the Columns slice.
// 		AddColumnsExprAlias("(e.price*x.tax*t.weee)", "final_price") // (e.price*x.tax*t.weee) AS `final_price`
func (b *Select) AddColumnsExprAlias(expressionAliases ...string) *Select {
	if (len(expressionAliases) % 2) == 1 {
		b.previousError = errors.NewMismatchf("[dbr] Expecting a balanced slice! Got: %v", expressionAliases)
	} else {
		b.Columns = b.Columns.appendColumnsAliases(expressionAliases, true)
	}
	return b
}

// AddRecord pulls in values to match Columns from the record generator.
func (b *Select) AddRecord(rec ArgumentAssembler) *Select {
	b.Record = rec
	return b
}

// AddArguments adds more arguments to the Argument field of the Select type.
// You must call this function directly after you have used e.g.
// AddColumnsExprAlias with place holders.
func (b *Select) AddArguments(args ...Argument) *Select {
	b.RawArguments = append(b.RawArguments, args...)
	return b
}

// Where appends a WHERE clause to the statement for the given string and args
// or map of column/value pairs.
func (b *Select) Where(c ...ConditionArg) *Select {
	b.WhereFragments = b.WhereFragments.append(c...)
	return b
}

// GroupBy appends columns to group the statement. The column gets always
// quoted. MySQL does not sort the results set. To avoid the overhead of sorting
// that GROUP BY produces this function should add an ORDER BY NULL with
// function `OrderByDeactivated`.
func (b *Select) GroupBy(columns ...string) *Select {
	b.GroupBys = b.GroupBys.appendColumns(columns, false)
	return b
}

// GroupByAsc sorts the groups in ascending order. No need to add an ORDER BY
// clause. When you use ORDER BY or GROUP BY to sort a column in a SELECT, the
// server sorts values using only the initial number of bytes indicated by the
// max_sort_length system variable.
func (b *Select) GroupByAsc(groups ...string) *Select {
	b.GroupBys = b.GroupBys.appendColumns(groups, false).applySort(len(groups), sortAscending)
	return b
}

// GroupByDesc sorts the groups in descending order. No need to add an ORDER BY
// clause. When you use ORDER BY or GROUP BY to sort a column in a SELECT, the
// server sorts values using only the initial number of bytes indicated by the
// max_sort_length system variable.
func (b *Select) GroupByDesc(groups ...string) *Select {
	b.GroupBys = b.GroupBys.appendColumns(groups, false).applySort(len(groups), sortDescending)
	return b
}

// GroupByExpr adds a custom SQL expression to the GROUP BY clause. Does not
// quote the strings nor add an ORDER BY NULL.
func (b *Select) GroupByExpr(groups ...string) *Select {
	b.GroupBys = b.GroupBys.appendColumns(groups, true)
	return b
}

// Having appends a HAVING clause to the statement
func (b *Select) Having(c ...ConditionArg) *Select {
	b.HavingFragments = b.HavingFragments.append(c...)
	return b
}

// OrderByDeactivated deactivates ordering of the result set by applying ORDER
// BY NULL to the SELECT statement. Very useful for GROUP BY queries.
func (b *Select) OrderByDeactivated() *Select {
	b.OrderBys = aliases{MakeAliasExpr("NULL")}
	return b
}

// OrderBy appends columns to the ORDER BY statement for ascending sorting.
// Columns are getting quoted. When you use ORDER BY or GROUP BY to sort a
// column in a SELECT, the server sorts values using only the initial number of
// bytes indicated by the max_sort_length system variable.
func (b *Select) OrderBy(columns ...string) *Select {
	b.OrderBys = b.OrderBys.appendColumns(columns, false)
	return b
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// Columns are getting quoted. When you use ORDER BY or GROUP BY to sort a
// column in a SELECT, the server sorts values using only the initial number of
// bytes indicated by the max_sort_length system variable.
func (b *Select) OrderByDesc(columns ...string) *Select {
	b.OrderBys = b.OrderBys.appendColumns(columns, false).applySort(len(columns), sortDescending)
	return b
}

// OrderByExpr adds a custom SQL expression to the ORDER BY clause. Does not
// quote the strings.
func (b *Select) OrderByExpr(columns ...string) *Select {
	b.OrderBys = b.OrderBys.appendColumns(columns, true)
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
// prepared statements. ToSQLs second argument `Arguments` will then be nil.
func (b *Select) Interpolate() *Select {
	b.IsInterpolate = true
	return b
}

func (b *Select) join(j string, t alias, on ...ConditionArg) *Select {
	jf := &joinFragment{
		JoinType: j,
		Table:    t,
	}
	jf.OnConditions = jf.OnConditions.append(on...)
	b.JoinFragments = append(b.JoinFragments, jf)
	return b
}

// Join creates an INNER join construct. By default, the onConditions are glued
// together with AND.
func (b *Select) Join(table alias, onConditions ...ConditionArg) *Select {
	return b.join("INNER", table, onConditions...)
}

// LeftJoin creates a LEFT join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) LeftJoin(table alias, onConditions ...ConditionArg) *Select {
	return b.join("LEFT", table, onConditions...)
}

// RightJoin creates a RIGHT join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) RightJoin(table alias, onConditions ...ConditionArg) *Select {
	return b.join("RIGHT", table, onConditions...)
}

// OuterJoin creates an OUTER join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) OuterJoin(table alias, onConditions ...ConditionArg) *Select {
	return b.join("OUTER", table, onConditions...)
}

// CrossJoin creates a CROSS join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) CrossJoin(table alias, onConditions ...ConditionArg) *Select {
	return b.join("CROSS", table, onConditions...)
}

// ToSQL converts the select statement into a string and returns its arguments.
func (b *Select) ToSQL() (string, Arguments, error) {
	return toSQL(b, b.IsInterpolate)
}

// argumentCapacity returns the total possible guessed size of a new Arguments
// slice. Use as the cap parameter in a call to `make`.
func (b *Select) argumentCapacity() int {
	return len(b.RawArguments) + len(b.JoinFragments) + len(b.WhereFragments)
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

func (b *Select) hasBuildCache() bool {
	return b.UseBuildCache
}

// ToSQL serialized the Select to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Select) toSQL(w queryWriter) error {
	if b.previousError != nil {
		return errors.Wrap(b.previousError, "[dbr] Select.toSQL")
	}
	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return errors.Wrap(err, "[dbr] Select.Listeners.dispatch")
	}

	if b.RawFullSQL != "" {
		_, err := w.WriteString(b.RawFullSQL)
		return err
	}

	if b.Table.Name == "" && b.Table.Select == nil {
		return errors.NewEmptyf("[dbr] Select: Table is missing")
	}
	if len(b.Columns) == 0 {
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
	if err := b.Columns.fQuoteAs(w); err != nil {
		return errors.Wrap(err, "[dbr] Select.toSQL.Columns.fQuoteAs")
	}

	w.WriteString(" FROM ")
	if err := b.Table.FquoteAs(w); err != nil {
		return errors.Wrap(err, "[dbr] Select.toSQL.Table.FquoteAs")
	}

	if len(b.JoinFragments) > 0 {
		for _, f := range b.JoinFragments {
			w.WriteByte(' ')
			w.WriteString(f.JoinType)
			w.WriteString(" JOIN ")
			f.Table.FquoteAs(w)
			if err := f.OnConditions.write(w, 'j'); err != nil {
				return errors.Wrap(err, "[dbr] Select.toSQL.write")
			}
		}
	}

	if err := b.WhereFragments.write(w, 'w'); err != nil {
		return errors.Wrap(err, "[dbr] Select.toSQL.write")
	}

	if len(b.GroupBys) > 0 {
		w.WriteString(" GROUP BY ")
		for i, c := range b.GroupBys {
			if i > 0 {
				w.WriteString(", ")
			}
			if err := c.FquoteAs(w); err != nil {
				return errors.Wrap(err, "[dbr] Select.toSQL.GroupBys")
			}
		}
	}

	if err := b.HavingFragments.write(w, 'h'); err != nil {
		return errors.Wrap(err, "[dbr] Select.toSQL.HavingFragments.write")
	}

	sqlWriteOrderBy(w, b.OrderBys, false)
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
	if b.previousError != nil {
		return nil, errors.Wrap(b.previousError, "[dbr] Select.toSQL")
	}

	if b.RawFullSQL != "" {
		return b.RawArguments, nil
	}

	// not sure if copying is necessary but leaves at least b.Arguments in pristine
	// condition
	if cap(args) == 0 {
		args = make(Arguments, 0, b.argumentCapacity())
	}
	args = append(args, b.RawArguments...)

	if args, err = b.Columns.appendArgs(args); err != nil {
		return nil, errors.Wrap(err, "[dbr] Select.toSQL.Columns.fQuoteAs")
	}

	if args, err = b.Table.appendArgs(args); err != nil {
		return nil, errors.Wrap(err, "[dbr] Select.toSQL.Table.FquoteAs")
	}

	var pap []int
	if len(b.JoinFragments) > 0 {
		for _, f := range b.JoinFragments {
			args, err = f.Table.appendArgs(args)
			if err != nil {
				return nil, errors.Wrap(err, "[dbr] Select.toSQL.write")
			}

			if args, pap, err = f.OnConditions.appendArgs(args, 'j'); err != nil {
				return nil, errors.Wrap(err, "[dbr] Select.toSQL.write")
			}
			if args, err = appendAssembledArgs(pap, b.Record, args, SQLStmtSelect|SQLPartJoin, f.OnConditions.Conditions()); err != nil {
				return nil, errors.Wrap(err, "[dbr] Select.toSQL.appendAssembledArgs")
			}
		}
	}

	if args, pap, err = b.WhereFragments.appendArgs(args, 'w'); err != nil {
		return nil, errors.Wrap(err, "[dbr] Select.toSQL.write")
	}
	if args, err = appendAssembledArgs(pap, b.Record, args, SQLStmtSelect|SQLPartWhere, b.WhereFragments.Conditions()); err != nil {
		return nil, errors.Wrap(err, "[dbr] Select.toSQL.appendAssembledArgs")
	}

	if args, pap, err = b.HavingFragments.appendArgs(args, 'h'); err != nil {
		return nil, errors.Wrap(err, "[dbr] Select.toSQL.HavingFragments.write")
	}
	if args, err = appendAssembledArgs(pap, b.Record, args, SQLStmtSelect|SQLPartHaving, b.HavingFragments.Conditions()); err != nil {
		return nil, errors.Wrap(err, "[dbr] Select.toSQL.appendAssembledArgs")
	}
	return args, nil
}
