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

	IsDistinct bool
	Columns    []string

	// FromTable table name and optional alias name to SELECT from.
	FromTable alias

	WhereFragments
	JoinFragments
	GroupBys        []string
	HavingFragments WhereFragments
	OrderBys        []string
	LimitCount      uint64
	LimitValid      bool
	OffsetCount     uint64
	OffsetValid     bool

	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners SelectListeners
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
		Log:       log.BlackHole{},
		FromTable: MakeAlias(from...),
	}
}

// NewSelectFromSub creates a new SELECT pointer using the provided sub-select
// in the FROM part together with an alias name. Appends the arguments of the
// sub-select to the parent *Select pointer arguments list. SQL result may look
// like:
//		SELECT a,b FROM (SELECT x,y FROM `product` AS `p`) AS `t`
func NewSelectFromSub(subSelect *Select, aliasName string) *Select {
	s := &Select{
		Log: log.BlackHole{},
		FromTable: alias{
			Select: subSelect,
			Alias:  aliasName,
		},
	}
	return s
}

// Select creates a new Select that select that given columns
func (c *Connection) Select(cols ...string) *Select {
	s := &Select{
		Log:     c.Log,
		Columns: cols,
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
func (tx *Tx) Select(cols ...string) *Select {
	s := &Select{
		Log:     tx.Logger,
		Columns: cols,
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

// Distinct marks the statement at a DISTINCT SELECT
func (b *Select) Distinct() *Select {
	b.IsDistinct = true
	return b
}

// From sets the table to SELECT FROM. If second argument will be provided this
// at then considered at the alias. SELECT ... FROM table AS alias.
func (b *Select) From(from ...string) *Select {
	b.FromTable = MakeAlias(from...)
	return b
}

// AddColumns appends more columns to the Columns slice. If a single string gets
// passed with comma separated values, this string gets split by the command and
// its values appended to the Columns slice.
func (b *Select) AddColumns(cols ...string) *Select {
	if len(cols) > 0 && strings.IndexByte(cols[0], ',') > 0 {
		cols = strings.Split(cols[0], ",")
		for i, c := range cols {
			cols[i] = strings.TrimSpace(c)
		}
	}
	b.Columns = append(b.Columns, cols...)
	return b
}

// todo
//func (b *Select) AddColumnsQuoted(cols ...string) *Select {
//	if len(cols) > 0 && strings.IndexByte(cols[0], ',') > 0 {
//		cols = strings.Split(cols[0], ",")
//		for i, c := range cols {
//			cols[i] = strings.TrimSpace(c)
//		}
//	}
//	b.Columns = append(b.Columns, cols...)
//	return b
//}

// AddColumnsAliases expects a balanced slice of ColumnName, AliasName and adds
// both concatenated and quoted to the Columns slice.
func (b *Select) AddColumnsAliases(colsAlias ...string) *Select {
	for i := 0; i < len(colsAlias); i = i + 2 {
		b.Columns = append(b.Columns, Quoter.Alias(colsAlias[i], colsAlias[i+1]))
	}
	return b
}

// Where appends a WHERE clause to the statement for the given string and args
// or map of column/value pairs
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

// OrderBy appends a column or an expression to ORDER the statement by
func (b *Select) OrderBy(ord ...string) *Select {
	b.OrderBys = append(b.OrderBys, ord...)
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

	if b.FromTable.Expression == "" && b.FromTable.Select == nil {
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

	for i, s := range b.Columns {
		if i > 0 {
			w.WriteString(", ")
		}
		w.WriteString(s)
	}

	if len(b.JoinFragments) > 0 {
		for _, f := range b.JoinFragments {
			for _, c := range f.Columns {
				w.WriteString(", ")
				w.WriteString(c)
			}
		}
	}

	w.WriteString(" FROM ")
	tArgs, err := b.FromTable.QuoteAsWriter(w)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Selec.toSQL.FromTable.QuoteAsWriter")
	}
	args = append(args, tArgs...)

	if len(b.JoinFragments) > 0 {
		for _, f := range b.JoinFragments {
			w.WriteRune(' ')
			w.WriteString(f.JoinType)
			w.WriteString(" JOIN ")
			f.Table.QuoteAsWriter(w)
			w.WriteString(" ON ")
			writeWhereFragmentsToSQL(f.OnConditions, w, &args)
		}
	}

	if len(b.WhereFragments) > 0 {
		w.WriteString(" WHERE ")
		writeWhereFragmentsToSQL(b.WhereFragments, w, &args)
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
		w.WriteString(" HAVING ")
		writeWhereFragmentsToSQL(b.HavingFragments, w, &args)
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
