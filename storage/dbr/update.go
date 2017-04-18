package dbr

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Update contains the clauses for an UPDATE statement
type Update struct {
	Log log.Logger
	DB  struct {
		Preparer
		Execer
	}
	// TODO: add UPDATE JOINS

	RawFullSQL   string
	RawArguments Arguments

	Table alias
	// SetClauses contains the column/argument association. For each column
	// there must be one argument.
	SetClauses UpdatedColumns
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
	Listeners UpdateListeners
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
	// previousError any error occurred during construction the SQL statement
	previousError error
}

// NewUpdate creates a new object with a black hole logger.
func NewUpdate(table ...string) *Update {
	return &Update{
		Log:   log.BlackHole{},
		Table: MakeAlias(table...),
	}
}

// Update creates a new Update for the given table
func (c *Connection) Update(table ...string) *Update {
	u := &Update{
		Log:   c.Log,
		Table: MakeAlias(table...),
	}
	u.DB.Execer = c.DB
	return u
}

// UpdateBySQL creates a new Update for the given SQL string and arguments
func (c *Connection) UpdateBySQL(sql string, args ...Argument) *Update {
	u := &Update{
		Log:          c.Log,
		RawFullSQL:   sql,
		RawArguments: args,
	}
	u.DB.Execer = c.DB
	return u
}

// Update creates a new Update for the given table bound to a transaction
func (tx *Tx) Update(table ...string) *Update {
	u := &Update{
		Log:   tx.Logger,
		Table: MakeAlias(table...),
	}
	u.DB.Execer = tx.Tx
	return u
}

// UpdateBySQL creates a new Update for the given SQL string and arguments bound
// to a transaction
func (tx *Tx) UpdateBySQL(sql string, args ...Argument) *Update {
	u := &Update{
		Log:          tx.Logger,
		RawFullSQL:   sql,
		RawArguments: args,
	}
	u.DB.Execer = tx.Tx
	return u
}

// Set appends a column/value pair for the statement
func (b *Update) Set(column string, arg Argument) *Update {
	if b.previousError != nil {
		return b
	}
	b.SetClauses.Columns = append(b.SetClauses.Columns, column)
	b.SetClauses.Arguments = append(b.SetClauses.Arguments, arg)
	return b
}

// SetMap appends the elements of the map at column/value pairs for the
// statement.
func (b *Update) SetMap(clauses map[string]Argument) *Update {
	if b.previousError != nil {
		return b
	}
	for col, val := range clauses {
		b.Set(col, val)
	}
	return b
}

// Where appends a WHERE clause to the statement
func (b *Update) Where(args ...ConditionArg) *Update {
	if b.previousError != nil {
		return b
	}
	appendConditions(&b.WhereFragments, args...)
	return b
}

// OrderBy appends a column or an expression to ORDER the statement ascending.
func (b *Update) OrderBy(ord ...string) *Update {
	b.OrderBys = append(b.OrderBys, ord...)
	return b
}

// OrderByDesc appends a column or an expression to ORDER the statement
// descending.
func (b *Update) OrderByDesc(ord ...string) *Update {
	for _, o := range ord {
		b.OrderBys = append(b.OrderBys, o+" DESC")
	}
	return b
}

// Limit sets a limit for the statement; overrides any existing LIMIT
func (b *Update) Limit(limit uint64) *Update {
	b.LimitCount = limit
	b.LimitValid = true
	return b
}

// Offset sets an offset for the statement; overrides any existing OFFSET
func (b *Update) Offset(offset uint64) *Update {
	b.OffsetCount = offset
	b.OffsetValid = true
	return b
}

// ToSQL serialized the Update to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Update) ToSQL() (string, Arguments, error) {
	if b.previousError != nil {
		return "", nil, errors.Wrap(b.previousError, "[dbr] Update.ToSQL")
	}

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return "", nil, errors.Wrap(err, "[dbr] Update.Listeners.dispatch")
	}

	if b.RawFullSQL != "" {
		return b.RawFullSQL, b.RawArguments, nil
	}

	if len(b.Table.Expression) == 0 {
		return "", nil, errors.NewEmptyf("[dbr] Update: Table at empty")
	}
	if len(b.SetClauses.Columns) == 0 {
		return "", nil, errors.NewEmptyf("[dbr] Update: SetClauses are empty")
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	var args = make(Arguments, 0, len(b.SetClauses.Arguments)+len(b.WhereFragments))

	buf.WriteString("UPDATE ")
	b.Table.QuoteAsWriter(buf)
	buf.WriteString(" SET ")

	// Build SET clause SQL with placeholders and add values to args
	for i, c := range b.SetClauses.Columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		Quoter.quoteAs(buf, c)
		buf.WriteByte('=')
		if i < len(b.SetClauses.Arguments) {
			arg := b.SetClauses.Arguments[i]
			if e, ok := arg.(*expr); ok {
				e.writeTo(buf, 0)
				args = append(args, e.Arguments...)
			} else {
				buf.WriteByte('?')
				args = append(args, arg)
			}
		} else {
			buf.WriteByte('?')
		}
	}

	// Write WHERE clause if we have any fragments
	if len(b.WhereFragments) > 0 {
		if err := writeWhereFragmentsToSQL(b.WhereFragments, buf, &args, 'w'); err != nil {
			return "", nil, errors.Wrap(err, "[dbr] Update.ToSQL.writeWhereFragmentsToSQL")
		}
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

// Exec executes the statement represented by the Update object. It returns the
// raw database/sql Result and an error if there was one.
func (b *Update) Exec(ctx context.Context) (sql.Result, error) {
	rawSQL, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Update.Exec.ToSQL")
	}

	fullSQL, err := Preprocess(rawSQL, args...)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Update.Exec.Preprocess")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Update.Exec.Timing", log.String("sql", fullSQL))
	}

	result, err := b.DB.ExecContext(ctx, fullSQL)
	if err != nil {
		return result, errors.Wrap(err, "[dbr] Update.Exec.Exec")
	}

	return result, nil
}

// Prepare creates a new prepared statement represented by the Update object. It
// returns the raw database/sql Stmt and an error if there was one.
func (b *Update) Prepare(ctx context.Context) (*sql.Stmt, error) {
	rawSQL, _, err := b.ToSQL() // TODO create a ToSQL version without any arguments
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Update.Prepare.ToSQL")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Update.Prepare.Timing", log.String("sql", rawSQL))
	}

	stmt, err := b.DB.PrepareContext(ctx, rawSQL)
	return stmt, errors.Wrap(err, "[dbr] Update.Prepare.Prepare")
}

// UpdatedColumns contains the column/argument association for either the SET
// clause in an UPDATE statement or to be used in an INSERT ... ON DUPLICATE KEY
// statement. For each column there must be one argument which can either be nil
// or has an actual value.
//
// When using the ON DUPLICATE KEY feature in the Insert builder:
//
// The function dbr.Expr is supported and allows SQL
// constructs like:
// 		`columnA`=VALUES(`columnB`)+2
// by writing the Go code:
//		ib.AddOnDuplicateKey("columnA", Expr("VALUES(`columnB`)+?", ArgInt(2)))
// Omitting the argument and using the keyword nil will turn this Go code:
//		ib.AddOnDuplicateKey("columnA", nil)
// into that SQL:
// 		`columnA`=VALUES(`columnA`)
// Same applies as when the columns gets only assigned without any arguments:
//		ib.OnDuplicateKey.Columns = []string{"name","sku"}
// will turn into:
// 		`name`=VALUES(`name`), `sku`=VALUES(`sku`)
// Type `UpdatedColumns` gets used in type `Update` with field
// `SetClauses` and in type `Insert` with field OnDuplicateKey.
type UpdatedColumns struct {
	Columns   []string
	Arguments Arguments
}

func (uc UpdatedColumns) writeOnDuplicateKey(w queryWriter, args *Arguments) error {
	if len(uc.Columns) == 0 {
		return nil
	}

	useArgs := len(uc.Arguments) == len(uc.Columns)

	w.WriteString(" ON DUPLICATE KEY UPDATE ")
	for i, c := range uc.Columns {
		if i > 0 {
			w.WriteString(", ")
		}
		Quoter.quote(w, c)
		w.WriteRune('=')
		if useArgs {
			// todo remove continue
			if e, ok := uc.Arguments[i].(*expr); ok {
				_ = e.writeTo(w, 0)
				*args = append(*args, uc.Arguments[i])
				continue
			}
			if uc.Arguments[i] == nil {
				w.WriteString("VALUES(")
				Quoter.quote(w, c)
				w.WriteRune(')')
				continue
			}
			w.WriteRune('?')
			*args = append(*args, uc.Arguments[i])
		} else {
			w.WriteString("VALUES(")
			Quoter.quote(w, c)
			w.WriteRune(')')
		}
	}
	return nil
}

// UpdateMulti allows to run an UPDATE statement multiple times with different
// values either in a transaction or as a preprocessed SQL string. Create one
// update statement without the SET arguments but with empty WHERE arguments.
// The empty WHERE arguments trigger the placeholder and the correct operator.
// The values itself will be provided either through the Records slice or via
// RecordChan.
type UpdateMulti struct {
	// UsePreprocess if true, preprocesses the SQL statement with the provided
	// arguments. The placeholders gets replaced with the current values. SQL
	// injections cannot not be possible *cough* *cough*.
	UsePreprocess bool
	// UseTransaction set to true to enable running the UPDATE queries in a
	// transaction.
	UseTransaction bool
	// IsolationLevel defines the transaction isolation level.
	sql.IsolationLevel
	// Tx knows how to start a transaction. Must be set if transactions hasn't
	// been disabled.
	Tx TxBeginner
	// Update represents the main UPDATE statement
	Update *Update

	// Alias provides a special feature that instead of the column name, the
	// alias will be passed to the ArgumentGenerater.Record function. If the alias
	// slice is empty the column names get passed. Otherwise the alias slice
	// must have the same length as the columns slice.
	Alias   []string
	Records []ArgumentGenerater
	// RecordChan waits for incoming records to send them to the prepared
	// statement. If the channel gets closed the transaction gets terminated and
	// the UPDATE statement removed.
	RecordChan <-chan ArgumentGenerater
}

// NewUpdateMulti creates new UPDATE statement which runs multiple times for a
// specific table.
func NewUpdateMulti(table ...string) *UpdateMulti {
	return &UpdateMulti{
		Update: NewUpdate(table...),
	}
}

// AddRecords pulls in values to match Columns from the record.
func (b *UpdateMulti) AddRecords(recs ...ArgumentGenerater) *UpdateMulti {
	b.Records = append(b.Records, recs...)
	return b
}

func (b *UpdateMulti) validate() error {
	if len(b.Update.SetClauses.Columns) == 0 {
		return errors.NewEmptyf("[dbr] UpdateMulti: Columns are empty")
	}
	if len(b.Alias) > 0 && len(b.Alias) != len(b.Update.SetClauses.Columns) {
		return errors.NewMismatchf("[dbr] UpdateMulti: Alias slice and Columns slice must have the same length")
	}
	if len(b.Records) == 0 && b.RecordChan == nil {
		return errors.NewEmptyf("[dbr] UpdateMulti: Records empty or RecordChan is nil")
	}
	return nil
}

func txUpdateMultiRollback(tx Txer, previousErr error, msg string, args ...interface{}) ([]sql.Result, error) {
	if err := tx.Rollback(); err != nil {
		eArg := []interface{}{previousErr}
		return nil, errors.Wrapf(err, "[dbr] UpdateMulti.Tx.Rollback. Previous Error: %s. "+msg, append(eArg, args...)...)
	}
	return nil, errors.Wrapf(previousErr, msg, args...)
}

// Exec creates a transaction
func (b *UpdateMulti) Exec(ctx context.Context) ([]sql.Result, error) {
	if err := b.validate(); err != nil {
		return nil, errors.Wrap(err, "[dbr] UpdateMulti.Exec")
	}

	rawSQL, _, err := b.Update.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] UpdateMulti.Exec.ToSQL")
	}

	if b.Update.Log != nil && b.Update.Log.IsInfo() {
		defer log.WhenDone(b.Update.Log).Info("dbr.UpdateMulti.Exec.Timing",
			log.String("sql", rawSQL), log.Int("records", len(b.Records)))
	}

	exec := b.Update.DB.Execer
	prep := b.Update.DB.Preparer
	var tx Txer = txMock{}
	if b.UseTransaction {
		// TODO fix context and make it set-able via outside
		tx, err = b.Tx.BeginTx(context.TODO(), &sql.TxOptions{
			Isolation: b.IsolationLevel,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "[dbr] UpdateMulti.Exec.Tx.BeginTx. with Query: %q", rawSQL)
		}
		exec = tx
		prep = tx
	}

	var stmt *sql.Stmt
	if !b.UsePreprocess {
		var err error
		stmt, err = prep.PrepareContext(ctx, rawSQL)
		if err != nil {
			return txUpdateMultiRollback(tx, err, "[dbr] UpdateMulti.Exec.Prepare. with Query: %q", rawSQL)
		}
		defer stmt.Close()
	}

	where := make([]string, len(b.Update.WhereFragments))
	for i, w := range b.Update.WhereFragments {
		where[i] = w.Condition
	}

	var results = make([]sql.Result, len(b.Records))
	for i, rec := range b.Records {
		cols := b.Update.SetClauses.Columns
		if len(b.Alias) > 0 {
			cols = b.Alias
		}

		args, err := rec.GenerateArguments(StatementTypeUpdate, cols, where)
		if err != nil {
			return txUpdateMultiRollback(tx, err, "[dbr] UpdateMulti.Exec.Record. Index %d with Query: %q", i, rawSQL)
		}

		if b.UsePreprocess {
			fullSQL, err := Preprocess(rawSQL, args...)
			if err != nil {
				return txUpdateMultiRollback(tx, err, "[dbr] UpdateMulti.Exec.Preprocess. Index %d with Query: %q", i, rawSQL)
			}

			results[i], err = exec.ExecContext(ctx, fullSQL)
			if err != nil {
				return txUpdateMultiRollback(tx, err, "[dbr] UpdateMulti.Exec.Exec. Index %d with Query: %q", i, rawSQL)
			}
		} else {
			results[i], err = stmt.ExecContext(ctx, args.Interfaces()...)
			if err != nil {
				return txUpdateMultiRollback(tx, err, "[dbr] UpdateMulti.Exec.Stmt.Exec. Index %d with Query: %q", i, rawSQL)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "[dbr] UpdateMulti.Tx.Commit. Query: %q", rawSQL)
	}

	return results, nil
}

// ExecChan executes incoming Records and writes the output into the provided
// channels. It closes the channels once the queries have been sent.
func (b *UpdateMulti) ExecChan(resChan chan<- sql.Result, errChan chan<- error) {
	defer close(resChan)
	defer close(errChan)
	if err := b.validate(); err != nil {
		errChan <- errors.Wrap(err, "[dbr] UpdateMulti.Exec")
		return
	}

	panic("TODO(CyS) implement")

}
