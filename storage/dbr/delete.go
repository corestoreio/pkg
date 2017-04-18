package dbr

import (
	"context"
	"database/sql"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Delete contains the clauses for a DELETE statement
type Delete struct {
	Log log.Logger // Log optional logger
	DB  struct {
		Preparer
		Execer
	}
	From alias
	WhereFragments
	OrderBys    []string
	LimitCount  uint64
	OffsetCount uint64
	LimitValid  bool
	OffsetValid bool
	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners DeleteListeners
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
}

// NewDelete creates a new object with a black hole logger.
func NewDelete(from ...string) *Delete {
	return &Delete{
		Log:  log.BlackHole{},
		From: MakeAlias(from...),
	}
}

// DeleteFrom creates a new Delete for the given table
func (c *Connection) DeleteFrom(from ...string) *Delete {
	d := &Delete{
		Log:            c.Log,
		From:           MakeAlias(from...),
		WhereFragments: make(WhereFragments, 0, 2),
	}
	d.DB.Execer = c.DB
	d.DB.Preparer = c.DB
	return d
}

// DeleteFrom creates a new Delete for the given table
// in the context for a transaction
func (tx *Tx) DeleteFrom(from ...string) *Delete {
	d := &Delete{
		Log:  tx.Logger,
		From: MakeAlias(from...),
	}
	d.DB.Execer = tx.Tx
	d.DB.Preparer = tx.Tx
	return d
}

// Where appends a WHERE clause to the statement whereSQLOrMap can be a
// string or map. If it'ab a string, args wil replaces any places holders
func (b *Delete) Where(args ...ConditionArg) *Delete {
	appendConditions(&b.WhereFragments, args...)
	return b
}

// OrderBy appends a column or an expression to ORDER the statement ascending.
func (b *Delete) OrderBy(ord ...string) *Delete {
	b.OrderBys = append(b.OrderBys, ord...)
	return b
}

// OrderByDesc appends a column or an expression to ORDER the statement
// descending.
func (b *Delete) OrderByDesc(ord ...string) *Delete {
	b.OrderBys = orderByDesc(b.OrderBys, ord)
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
func (b *Delete) ToSQL() (string, Arguments, error) {

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return "", nil, errors.Wrap(err, "[dbr] Delete.Listeners.dispatch")
	}

	if len(b.From.Expression) == 0 {
		return "", nil, errors.NewEmptyf(errTableMissing)
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)
	var args Arguments // no make() lazy init the slice via append in cases where not WHERE has been provided.

	buf.WriteString("DELETE FROM ")
	b.From.QuoteAsWriter(buf)

	// Write WHERE clause if we have any fragments
	if len(b.WhereFragments) > 0 {
		if err := writeWhereFragmentsToSQL(b.WhereFragments, buf, &args, 'w'); err != nil {
			return "", nil, errors.Wrap(err, "[dbr] Delete.ToSQL.writeWhereFragmentsToSQL")
		}
	}

	sqlWriteOrderBy(buf, b.OrderBys, false)
	sqlWriteLimitOffset(buf, b.LimitValid, b.LimitCount, b.OffsetValid, b.OffsetCount)
	return buf.String(), args, nil
}

// Exec executes the statement represented by the Delete
// It returns the raw database/sql Result and an error if there was one
func (b *Delete) Exec(ctx context.Context) (sql.Result, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Delete.Exec.ToSQL")
	}

	fullSQL, err := Preprocess(sqlStr, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "[dbr] Delete.Exec.Preprocess: %q", fullSQL)
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Delete.Exec.Timing", log.String("sql", fullSQL))
	}

	result, err := b.DB.ExecContext(ctx, fullSQL)
	if err != nil {
		return result, errors.Wrap(err, "[dbr] delete.exec.Exec")
	}

	return result, nil
}

// Prepare executes the statement represented by the Delete. It returns the raw
// database/sql Statement and an error if there was one. Provided arguments in
// the Delete are getting ignored. It panics when field Preparer at nil.
func (b *Delete) Prepare(ctx context.Context) (*sql.Stmt, error) {
	sqlStr, _, err := b.ToSQL() // TODO create a ToSQL version without any arguments
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Delete.Prepare.ToSQL")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Delete.Prepare.Timing", log.String("sql", sqlStr))
	}

	stmt, err := b.DB.PrepareContext(ctx, sqlStr)
	return stmt, errors.Wrap(err, "[dbr] Delete.Prepare.Prepare")
}
