package dbr

import (
	"strconv"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
)

// Select contains the clauses for a SELECT statement
type Select struct {
	log.Logger // optional
	// The next three fields depend on which method receiver you would like to
	// execute. Leaving them empty results in a panic.
	Querier
	QueryRower
	Preparer

	RawFullSQL   string
	RawArguments []interface{}

	IsDistinct bool
	Columns    []string
	FromTable  alias
	WhereFragments
	JoinFragments
	GroupBys        []string
	HavingFragments WhereFragments
	OrderBys        []string
	LimitCount      uint64
	LimitValid      bool
	OffsetCount     uint64
	OffsetValid     bool

	// SelectListeners allows to dispatch certain functions in different
	// situations.
	SelectListeners
	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
}

// NewSelect creates a new object with a black hole logger.
func NewSelect(from ...string) *Select {
	return &Select{
		Logger:    log.BlackHole{},
		FromTable: MakeAlias(from...),
	}
}

// Select creates a new Select that select that given columns
func (sess *Session) Select(cols ...string) *Select {
	return &Select{
		Logger:     sess.Logger,
		Querier:    sess.cxn.DB,
		QueryRower: sess.cxn.DB,
		Preparer:   sess.cxn.DB,
		Columns:    cols,
	}
}

// SelectBySQL creates a new Select for the given SQL string and arguments
func (sess *Session) SelectBySQL(sql string, args ...interface{}) *Select {
	return &Select{
		Logger:       sess.Logger,
		Querier:      sess.cxn.DB,
		QueryRower:   sess.cxn.DB,
		Preparer:     sess.cxn.DB,
		RawFullSQL:   sql,
		RawArguments: args,
	}
}

// Select creates a new Select that select that given columns bound to the transaction
func (tx *Tx) Select(cols ...string) *Select {
	return &Select{
		Logger:     tx.Logger,
		QueryRower: tx.Tx,
		Querier:    tx.Tx,
		Preparer:   tx.Tx,
		Columns:    cols,
	}
}

// SelectBySQL creates a new Select for the given SQL string and arguments bound to the transaction
func (tx *Tx) SelectBySQL(sql string, args ...interface{}) *Select {
	return &Select{
		Logger:       tx.Logger,
		QueryRower:   tx.Tx,
		Querier:      tx.Tx,
		Preparer:     tx.Tx,
		RawFullSQL:   sql,
		RawArguments: args,
	}
}

// Distinct marks the statement as a DISTINCT SELECT
func (b *Select) Distinct() *Select {
	b.IsDistinct = true
	return b
}

// From sets the table to SELECT FROM. If second argument will be provided this is
// then considered as the alias. SELECT ... FROM table AS alias.
func (b *Select) From(from ...string) *Select {
	b.FromTable = MakeAlias(from...)
	return b
}

// Where appends a WHERE clause to the statement for the given string and args
// or map of column/value pairs
func (b *Select) Where(args ...ConditionArg) *Select {
	b.WhereFragments = append(b.WhereFragments, newWhereFragments(args...)...)
	return b
}

// GroupBy appends a column to group the statement
func (b *Select) GroupBy(group string) *Select {
	b.GroupBys = append(b.GroupBys, group)
	return b
}

// Having appends a HAVING clause to the statement
func (b *Select) Having(args ...ConditionArg) *Select {
	b.HavingFragments = append(b.HavingFragments, newWhereFragments(args...)...)
	return b
}

// OrderBy appends a column to ORDER the statement by
func (b *Select) OrderBy(ord string) *Select {
	b.OrderBys = append(b.OrderBys, ord)
	return b
}

// OrderDir appends a column to ORDER the statement by with a given direction
func (b *Select) OrderDir(ord string, isAsc bool) *Select {
	if isAsc {
		b.OrderBys = append(b.OrderBys, ord+" ASC")
	} else {
		b.OrderBys = append(b.OrderBys, ord+" DESC")
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

// ToSQL serialized the Select to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Select) ToSQL() (string, []interface{}, error) {

	if err := b.SelectListeners.dispatch(b.Logger, OnBeforeToSQL, b); err != nil {
		return "", nil, errors.Wrap(err, "[dbr] Select.Listeners.dispatch")
	}

	if b.RawFullSQL != "" {
		return b.RawFullSQL, b.RawArguments, nil
	}

	if len(b.FromTable.Expression) == 0 {
		return "", nil, errors.NewEmptyf(errTableMissing)
	}
	if len(b.Columns) == 0 {
		return "", nil, errors.NewEmptyf(errColumnsMissing)
	}

	var sql = bufferpool.Get()
	defer bufferpool.Put(sql)

	var args []interface{}

	sql.WriteString("SELECT ")

	if b.IsDistinct {
		sql.WriteString("DISTINCT ")
	}

	for i, s := range b.Columns {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(s)
	}

	if len(b.JoinFragments) > 0 {
		for _, f := range b.JoinFragments {
			for _, c := range f.Columns {
				sql.WriteString(", ")
				sql.WriteString(c)
			}
		}
	}

	sql.WriteString(" FROM ")
	sql.WriteString(b.FromTable.QuoteAs())

	if len(b.JoinFragments) > 0 {
		for _, f := range b.JoinFragments {
			sql.WriteRune(' ')
			sql.WriteString(f.JoinType)
			sql.WriteString(" JOIN ")
			sql.WriteString(f.Table.QuoteAs())
			sql.WriteString(" ON ")
			writeWhereFragmentsToSQL(f.OnConditions, sql, &args)
		}
	}

	if len(b.WhereFragments) > 0 {
		sql.WriteString(" WHERE ")
		writeWhereFragmentsToSQL(b.WhereFragments, sql, &args)
	}

	if len(b.GroupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		for i, s := range b.GroupBys {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(s)
		}
	}

	if len(b.HavingFragments) > 0 {
		sql.WriteString(" HAVING ")
		writeWhereFragmentsToSQL(b.HavingFragments, sql, &args)
	}

	if len(b.OrderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		for i, s := range b.OrderBys {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(s)
		}
	}

	if b.LimitValid {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.FormatUint(b.LimitCount, 10))
	}

	if b.OffsetValid {
		sql.WriteString(" OFFSET ")
		sql.WriteString(strconv.FormatUint(b.OffsetCount, 10))
	}
	return sql.String(), args, nil
}
