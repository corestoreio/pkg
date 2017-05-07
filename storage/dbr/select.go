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
	DB struct {
		Querier
		Preparer
	}

	RawFullSQL string
	Arguments

	// Columns represents a slice of names and its optional aliases. Wildcard
	// `SELECT *` statements are not really supported:
	// http://stackoverflow.com/questions/3639861/why-is-select-considered-harmful
	Columns aliases

	//TODO: create a possibility of the Select type which has a half-pre-rendered
	// SQL statement where a developer can only modify or append WHERE clauses.
	// especially useful during code generation

	// Table table name and optional alias name to SELECT from.
	Table alias

	WhereFragments
	JoinFragments
	GroupBys          []string // TODO use aliases
	HavingFragments   WhereFragments
	OrderBys          []string // TODO use aliases
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
	return &Select{
		Columns: appendColumns(nil, columns),
	}
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
		Log:     c.Log,
		Columns: appendColumns(nil, columns),
	}
	s.DB.Querier = c.DB
	s.DB.Preparer = c.DB
	return s
}

// SelectBySQL creates a new Select for the given SQL string and arguments
func (c *Connection) SelectBySQL(sql string, args ...Argument) *Select {
	s := &Select{
		Log:        c.Log,
		RawFullSQL: sql,
		Arguments:  args,
	}
	s.DB.Querier = c.DB
	s.DB.Preparer = c.DB
	return s
}

// Select creates a new Select that select that given columns bound to the transaction
func (tx *Tx) Select(columns ...string) *Select {
	s := &Select{
		Log:     tx.Logger,
		Columns: appendColumns(nil, columns),
	}
	s.DB.Querier = tx.Tx
	s.DB.Preparer = tx.Tx
	return s
}

// SelectBySQL creates a new Select for the given SQL string and arguments bound to the transaction
func (tx *Tx) SelectBySQL(sql string, args ...Argument) *Select {
	s := &Select{
		Log:        tx.Logger,
		RawFullSQL: sql,
		Arguments:  args,
	}
	s.DB.Querier = tx.Tx
	s.DB.Preparer = tx.Tx
	return s
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
	b.Columns = appendColumns(b.Columns, cols)
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
		b.Columns = appendColumnsAliases(b.Columns, columnAliases, false)
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
		b.Columns = appendColumnsAliases(b.Columns, expressionAliases, true)
	}
	return b
}

// AddArguments adds more arguments to the Argument field of the Select type.
// You must call this function directly after you have used e.g.
// AddColumnsExprAlias with place holders.
func (b *Select) AddArguments(args ...Argument) *Select {
	b.Arguments = append(b.Arguments, args...)
	return b
}

// Where appends a WHERE clause to the statement for the given string and args
// or map of column/value pairs.
func (b *Select) Where(c ...ConditionArg) *Select {
	b.WhereFragments = appendConditions(b.WhereFragments, c...)
	return b
}

// GroupBy appends a column or an expression to group the statement.
func (b *Select) GroupBy(groups ...string) *Select {
	b.GroupBys = append(b.GroupBys, groups...)
	return b
}

// Having appends a HAVING clause to the statement
func (b *Select) Having(c ...ConditionArg) *Select {
	b.HavingFragments = appendConditions(b.HavingFragments, c...)
	return b
}

// OrderBy appends a column or an expression to ORDER the statement ascending.
func (b *Select) OrderBy(ord ...string) *Select {
	b.OrderBys = append(b.OrderBys, ord...)
	return b
}

// OrderByDesc appends a column or an expression to ORDER the statement
// descending.
func (b *Select) OrderByDesc(ord ...string) *Select {
	b.OrderBys = orderByDesc(b.OrderBys, ord)
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

// ToSQL converts the select statement into a string and returns its arguments.
func (b *Select) ToSQL() (string, Arguments, error) {
	return toSQL(b, b.IsInterpolate)
}

// ToSQL serialized the Select to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Select) toSQL(w queryWriter) (Arguments, error) {
	if b.previousError != nil {
		return nil, errors.Wrap(b.previousError, "[dbr] Select.toSQL")
	}
	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return nil, errors.Wrap(err, "[dbr] Select.Listeners.dispatch")
	}
	// TODO(CyS) implement SQL string cache. If cache set to true, then the
	// finalized query will be written in the empty RawFullSQL field. if cache
	// has been set to false, then query gets regenerated.

	if b.RawFullSQL != "" {
		w.WriteString(b.RawFullSQL)
		return b.Arguments, nil
	}

	if b.Table.Name == "" && b.Table.Select == nil {
		return nil, errors.NewEmptyf(errTableMissing)
	}
	if len(b.Columns) == 0 {
		return nil, errors.NewEmptyf(errColumnsMissing)
	}

	// not sure if copying is necessary but leaves at least b.Arguments in pristine
	// condition
	var args = make(Arguments, len(b.Arguments), len(b.Arguments)+len(b.JoinFragments)+len(b.WhereFragments))
	copy(args, b.Arguments)

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
	{
		var err error
		if args, err = b.Columns.fQuoteAs(w, args); err != nil {
			return nil, errors.Wrap(err, "[dbr] Selec.toSQL.Columns.fQuoteAs")
		}
	}

	w.WriteString(" FROM ")
	tArgs, err := b.Table.FquoteAs(w)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Selec.toSQL.Table.FquoteAs")
	}
	args = append(args, tArgs...)

	if len(b.JoinFragments) > 0 {
		for _, f := range b.JoinFragments {
			w.WriteByte(' ')
			w.WriteString(f.JoinType)
			w.WriteString(" JOIN ")
			f.Table.FquoteAs(w)
			if args, err = writeWhereFragmentsToSQL(f.OnConditions, w, args, 'j'); err != nil {
				return nil, errors.Wrap(err, "[dbr] Select.toSQL.writeWhereFragmentsToSQL")
			}
		}
	}

	if args, err = writeWhereFragmentsToSQL(b.WhereFragments, w, args, 'w'); err != nil {
		return nil, errors.Wrap(err, "[dbr] Select.toSQL.writeWhereFragmentsToSQL")
	}

	if len(b.GroupBys) > 0 {
		w.WriteString(" GROUP BY ")
		for i, s := range b.GroupBys {
			if i > 0 {
				w.WriteString(", ")
			}
			w.WriteString(s)
		}
	}

	if args, err = writeWhereFragmentsToSQL(b.HavingFragments, w, args, 'h'); err != nil {
		return nil, errors.Wrap(err, "[dbr] Select.toSQL.writeWhereFragmentsToSQL")
	}

	sqlWriteOrderBy(w, b.OrderBys, false)
	sqlWriteLimitOffset(w, b.LimitValid, b.LimitCount, b.OffsetValid, b.OffsetCount)
	switch {
	case b.IsLockInShareMode:
		w.WriteString(" LOCK IN SHARE MODE")
	case b.IsForUpdate:
		w.WriteString(" FOR UPDATE")
	}
	return args, nil
}
