package dbr

import (
	"database/sql"
	"strconv"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
)

// Delete contains the clauses for a DELETE statement
type Delete struct {
	log.Logger // optional
	Execer
	Preparer

	From alias
	WhereFragments
	OrderBys    []string
	LimitCount  uint64
	LimitValid  bool
	OffsetCount uint64
	OffsetValid bool

	// DeleteListeners allows to dispatch certain functions in different
	// situations.
	DeleteListeners
}

// NewDelete creates a new object with a black hole logger.
func NewDelete(from ...string) *Delete {
	return &Delete{
		Logger: log.BlackHole{},
		From:   MakeAlias(from...),
	}
}

// DeleteFrom creates a new Delete for the given table
func (sess *Session) DeleteFrom(from ...string) *Delete {
	return &Delete{
		Logger:         sess.Logger,
		Execer:         sess.cxn.DB,
		Preparer:       sess.cxn.DB,
		From:           MakeAlias(from...),
		WhereFragments: make(WhereFragments, 0, 2),
	}
}

// DeleteFrom creates a new Delete for the given table
// in the context for a transaction
func (tx *Tx) DeleteFrom(from ...string) *Delete {
	return &Delete{
		Logger:         tx.Logger,
		Execer:         tx.Tx,
		Preparer:       tx.Tx,
		From:           MakeAlias(from...),
		WhereFragments: make(WhereFragments, 0, 2),
	}
}

// Where appends a WHERE clause to the statement whereSQLOrMap can be a
// string or map. If it's a string, args wil replaces any places holders
func (b *Delete) Where(args ...ConditionArg) *Delete {
	b.WhereFragments = append(b.WhereFragments, newWhereFragments(args...)...)
	return b
}

// OrderBy appends an ORDER BY clause to the statement
func (b *Delete) OrderBy(ord string) *Delete {
	b.OrderBys = append(b.OrderBys, ord)
	return b
}

// OrderDir appends an ORDER BY clause with a direction to the statement
func (b *Delete) OrderDir(ord string, isAsc bool) *Delete {
	if isAsc {
		b.OrderBys = append(b.OrderBys, ord+" ASC")
	} else {
		b.OrderBys = append(b.OrderBys, ord+" DESC")
	}
	return b
}

// Limit sets a LIMIT clause for the statement; overrides any existing LIMIT
func (b *Delete) Limit(limit uint64) *Delete {
	b.LimitCount = limit
	b.LimitValid = true
	return b
}

// Offset sets an OFFSET clause for the statement; overrides any existing OFFSET
func (b *Delete) Offset(offset uint64) *Delete {
	b.OffsetCount = offset
	b.OffsetValid = true
	return b
}

// ToSQL serialized the Delete to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Delete) ToSQL() (string, []interface{}, error) {

	if err := b.DeleteListeners.dispatch(OnBeforeToSQL, b); err != nil {
		return "", nil, errors.Wrap(err, "[dbr] Delete.Listeners.dispatch")
	}

	if len(b.From.Expression) == 0 {
		return "", nil, errors.NewEmptyf(errTableMissing)
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)
	var args []interface{}

	buf.WriteString("DELETE FROM ")
	buf.WriteString(b.From.QuoteAs())

	// Write WHERE clause if we have any fragments
	if len(b.WhereFragments) > 0 {
		buf.WriteString(" WHERE ")
		writeWhereFragmentsToSQL(b.WhereFragments, buf, &args)
	}

	// Ordering and limiting
	if len(b.OrderBys) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, s := range b.OrderBys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}

	if b.LimitValid {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.FormatUint(b.LimitCount, 10))
	}

	if b.OffsetValid {
		buf.WriteString(" OFFSET ")
		buf.WriteString(strconv.FormatUint(b.OffsetCount, 10))
	}
	return buf.String(), args, nil
}

// Exec executes the statement represented by the Delete
// It returns the raw database/sql Result and an error if there was one
func (b *Delete) Exec() (sql.Result, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Delete.Exec.ToSQL")
	}

	fullSQL, err := Preprocess(sqlStr, args)
	if err != nil {
		return nil, errors.Wrapf(err, "[dbr] Delete.Exec.Preprocess: %q", fullSQL)
	}

	if b.Logger != nil && b.Logger.IsInfo() {
		defer log.WhenDone(b.Logger).Info("dbr.Delete.Exec.Timing", log.String("sql", fullSQL))
	}

	result, err := b.Execer.Exec(fullSQL)
	if err != nil {
		return result, errors.Wrap(err, "[dbr] delete.exec.Exec")
	}

	return result, nil
}

// Prepare executes the statement represented by the Delete. It returns the raw
// database/sql Statement and an error if there was one. Provided arguments in
// the Delete are getting ignored. It panics when field Preparer is nil.
func (b *Delete) Prepare() (*sql.Stmt, error) {
	sqlStr, _, err := b.ToSQL() // TODO create a ToSQL version without any arguments
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Delete.Prepare.ToSQL")
	}

	if b.Logger != nil && b.Logger.IsInfo() {
		defer log.WhenDone(b.Logger).Info("dbr.Delete.Prepare.Timing", log.String("sql", sqlStr))
	}

	stmt, err := b.Preparer.Prepare(sqlStr)
	return stmt, errors.Wrap(err, "[dbr] Delete.Prepare.Prepare")
}
