package dbr

import (
	"strconv"
	"strings"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Select contains the clauses for a SELECT statement
type Select struct {
	Log log.Logger // Log optional logger
	// DB gets required once the Load*() functions will be used.
	DB struct {
		Querier
		QueryRower
		Preparer
	}

	RawFullSQL string
	Arguments

	Columns []string

	// Table table name and optional alias name to SELECT from.
	Table alias

	WhereFragments
	JoinFragments
	GroupBys        []string
	HavingFragments WhereFragments
	OrderBys        []string
	LimitCount      uint64
	OffsetCount     uint64
	LimitValid      bool
	OffsetValid     bool
	IsDistinct      bool // See Distinct()
	IsStraightJoin  bool // See StraightJoin()
	IsSQLNoCache    bool // See IsSQLNoCache()
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
}

// NewSelect creates a new Select object with a black hole logger and selecting
// from the specified columns. The provided columns won't get quoted.
func NewSelect(columns ...string) *Select {
	return &Select{
		Log:     log.BlackHole{},
		Columns: columns,
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
		Log: log.BlackHole{},
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
		Columns: columns,
	}
	s.DB.Querier = c.DB
	s.DB.QueryRower = c.DB
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
	s.DB.QueryRower = c.DB
	s.DB.Preparer = c.DB
	return s
}

// Select creates a new Select that select that given columns bound to the transaction
func (tx *Tx) Select(columns ...string) *Select {
	s := &Select{
		Log:     tx.Logger,
		Columns: columns,
	}
	s.DB.Querier = tx.Tx
	s.DB.QueryRower = tx.Tx
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
	s.DB.QueryRower = tx.Tx
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

// From sets the table to SELECT FROM. If second argument will be provided this
// at then considered at the alias. SELECT ... FROM table AS alias.
func (b *Select) From(from ...string) *Select {
	b.Table = MakeAlias(from...)
	return b
}

func splitColumns(cols []string) []string {
	for i := 0; i < len(cols); i++ {
		if c := cols[i]; strings.IndexByte(c, ',') > 0 {
			cs := strings.Split(c, ",")
			for j, c2 := range cs {
				cs[j] = strings.TrimSpace(c2)
			}
			cols = append(cols[:i], append(cs, cols[i+1:]...)...)
		}
	}
	return cols
}

// AddColumns appends more columns to the Columns slice. If a single string gets
// passed with comma separated values, this string gets split by the comma and
// its values appended to the Columns slice. Columns won't get quoted.
// 		AddColumns("a","b") 		// []string{"a","b"}
// 		AddColumns("a,b","z","c,d")	// []string{"a","b","z","c","d"}
func (b *Select) AddColumns(cols ...string) *Select {
	b.Columns = append(b.Columns, splitColumns(cols)...)
	return b
}

// AddColumnsQuoted appends more columns to the Columns slice and quotes them.
// Give "t1.name" gets translated to "`t1`.`name`". Comma separated input is
// supported for each slice item:
//		AddColumnsQuoted("t1.name","t1.sku","price") // []string{"`t1`.`name`", "`t1`.`sku`","`price`"}
//		AddColumnsQuoted("t1.name,t1.sku")	// []string{"`t1`.`name`", "`t1`.`sku`"}
func (b *Select) AddColumnsQuoted(cols ...string) *Select {
	cols = splitColumns(cols)
	for i, c := range cols {
		cols[i] = Quoter.QuoteAs(c)
	}
	b.Columns = append(b.Columns, cols...)
	return b
}

// AddColumnsQuotedAlias expects a balanced slice of "ColumnName, AliasName" and
// adds both concatenated and quoted to the Columns slice. It panics when the
// provided `columnAliases` seems not be balanced.
//		AddColumnsQuotedAlias("t1.name","t1Name","t1.sku","t1SKU") // []string{"`t1`.`name` AS `t1Name`", "`t1`.`sku` AS `t1SKU`"}
func (b *Select) AddColumnsQuotedAlias(columnAliases ...string) *Select {
	columnAliases = splitColumns(columnAliases)
	for i := 0; i < len(columnAliases); i = i + 2 {
		b.Columns = append(b.Columns, Quoter.QuoteAs(columnAliases[i], columnAliases[i+1]))
	}
	return b
}

// AddColumnsExprAlias expects a balanced slice of "expression, AliasName" and
// adds both concatenated and quoted to the Columns slice. It panics when the
// provided `expressionAlias` seems not be balanced.
// 		AddColumnsExprAlias("(e.price*x.tax*t.weee)", "final_price") // (e.price*x.tax*t.weee) AS `final_price`
func (b *Select) AddColumnsExprAlias(expressionAliases ...string) *Select {
	for i := 0; i < len(expressionAliases); i = i + 2 {
		b.Columns = append(b.Columns, Quoter.ExprAlias(expressionAliases[i], expressionAliases[i+1]))
	}
	return b
}

// Where appends a WHERE clause to the statement for the given string and args
// or map of column/value pairs.
func (b *Select) Where(c ...ConditionArg) *Select {
	appendConditions(&b.WhereFragments, c...)
	return b
}

// GroupBy appends a column or an expression to group the statement.
func (b *Select) GroupBy(groups ...string) *Select {
	b.GroupBys = append(b.GroupBys, groups...)
	return b
}

// Having appends a HAVING clause to the statement
func (b *Select) Having(c ...ConditionArg) *Select {
	appendConditions(&b.HavingFragments, c...)
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
	for _, o := range ord {
		b.OrderBys = append(b.OrderBys, o+" DESC")
	}
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

// ToSQL converts the select statement into a string and returns its arguments.
func (b *Select) ToSQL() (string, Arguments, error) {
	var w = bufferpool.Get()
	defer bufferpool.Put(w)
	args, err := b.toSQL(w)
	return w.String(), args, err
}

// ToSQL serialized the Select to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Select) toSQL(w queryWriter) (Arguments, error) {

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

	if b.Table.Expression == "" && b.Table.Select == nil {
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

	for i, s := range b.Columns {
		if i > 0 {
			w.WriteString(", ")
		}
		w.WriteString(s)
	}

	w.WriteString(" FROM ")
	tArgs, err := b.Table.QuoteAsWriter(w)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Selec.toSQL.Table.QuoteAsWriter")
	}
	args = append(args, tArgs...)

	if len(b.JoinFragments) > 0 {
		for _, f := range b.JoinFragments {
			w.WriteRune(' ')
			w.WriteString(f.JoinType)
			w.WriteString(" JOIN ")
			f.Table.QuoteAsWriter(w)
			if err := writeWhereFragmentsToSQL(f.OnConditions, w, &args, 'j'); err != nil {
				return nil, errors.Wrap(err, "[dbr] Select.toSQL.writeWhereFragmentsToSQL")
			}
		}
	}

	if len(b.WhereFragments) > 0 {
		if err := writeWhereFragmentsToSQL(b.WhereFragments, w, &args, 'w'); err != nil {
			return nil, errors.Wrap(err, "[dbr] Select.toSQL.writeWhereFragmentsToSQL")
		}
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

	if len(b.HavingFragments) > 0 {
		if err := writeWhereFragmentsToSQL(b.HavingFragments, w, &args, 'h'); err != nil {
			return nil, errors.Wrap(err, "[dbr] Select.toSQL.writeWhereFragmentsToSQL")
		}
	}

	if len(b.OrderBys) > 0 {
		w.WriteString(" ORDER BY ")
		for i, s := range b.OrderBys {
			if i > 0 {
				w.WriteString(", ")
			}
			w.WriteString(s)
		}
	}

	if b.LimitValid {
		w.WriteString(" LIMIT ")
		w.WriteString(strconv.FormatUint(b.LimitCount, 10))
	}

	if b.OffsetValid {
		w.WriteString(" OFFSET ")
		w.WriteString(strconv.FormatUint(b.OffsetCount, 10))
	}
	return args, nil
}
