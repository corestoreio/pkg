package dbr

import (
	"context"
	"database/sql"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Insert contains the clauses for an INSERT statement
type Insert struct {
	Log log.Logger // Log optional logger
	DB  struct {
		Preparer
		Execer
	}

	Into string
	Cols []string
	Vals Arguments

	Recs []ArgumentGenerater
	Maps map[string]Argument

	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners InsertListeners
	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
	// previousError any error occurred during construction the SQL statement
	previousError error
}

// NewInsert creates a new object with a black hole logger.
func NewInsert(into string) *Insert {
	return &Insert{
		Log:  log.BlackHole{},
		Into: into,
	}
}

// InsertInto instantiates a Insert for the given table
func (c *Connection) InsertInto(into string) *Insert {
	i := &Insert{
		Log:  c.Log,
		Into: into,
	}
	i.DB.Execer = c.DB
	i.DB.Preparer = c.DB
	return i
}

// InsertInto instantiates a Insert for the given table bound to a transaction
func (tx *Tx) InsertInto(into string) *Insert {
	i := &Insert{
		Log:  tx.Logger,
		Into: into,
	}
	i.DB.Execer = tx.Tx
	i.DB.Preparer = tx.Tx
	return i
}

// Columns appends columns to insert in the statement.
func (b *Insert) Columns(columns ...string) *Insert {
	b.Cols = append(b.Cols, columns...)
	return b
}

// Values appends a set of values to the statement. Pro Tip: Use Values() and
// not Record() to avoid reflection. Only this function will consider the
// driver.Valuer interface when you pass a pointer to the value. Values must be
// balanced to the number of columns. You can even provide more values, like
// records. see BenchmarkInsertValuesSQL
func (b *Insert) Values(vals ...Argument) *Insert {
	b.Vals = append(b.Vals, vals...)
	return b
}

// Record pulls in values to match Columns from the record. Think about a vector
// on how to use this.
func (b *Insert) Record(recs ...ArgumentGenerater) *Insert {
	b.Recs = append(b.Recs, recs...)
	return b
}

// Map pulls in values to match Columns from the record. Calling multiple
// times will add new map entries to the Insert map.
func (b *Insert) Map(m map[string]Argument) *Insert {
	if b.Maps == nil {
		b.Maps = make(map[string]Argument)
	}
	for col, val := range m {
		b.Maps[col] = val
	}
	return b
}

// Pair adds a key/value pair to the statement.
func (b *Insert) Pair(column string, arg Argument) *Insert {
	if b.previousError != nil {
		return b
	}
	for _, c := range b.Cols {
		if c == column {
			b.previousError = errors.NewAlreadyExistsf("[dbr] Column %q has already been added", c)
			return b
		}
	}
	b.Cols = append(b.Cols, column)
	b.Vals = append(b.Vals, arg)
	return b
}

// FromSelect creates an "INSERT INTO `table` SELECT ..." statement from a
// previously created SELECT statement.
func (b *Insert) FromSelect(s *Select) (string, Arguments, error) {
	if b.previousError != nil {
		return "", nil, errors.Wrap(b.previousError, "[dbr] Insert.ToSQL")
	}

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return "", nil, errors.Wrap(err, "[dbr] Insert.Listeners.dispatch")
	}

	if len(b.Into) == 0 {
		return "", nil, errors.NewEmptyf(errTableMissing)
	}

	sSQL, sArgs, err := s.ToSQL()
	if err != nil {
		return "", nil, errors.Wrap(err, "[dbr] Insert.FromSelect")
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	buf.WriteString("INSERT INTO ")
	Quoter.quote(buf, b.Into)
	buf.WriteByte(' ')
	buf.WriteString(sSQL)

	return buf.String(), sArgs, nil
}

// ToSQL serialized the Insert to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Insert) ToSQL() (string, Arguments, error) {
	if b.previousError != nil {
		return "", nil, errors.Wrap(b.previousError, "[dbr] Insert.ToSQL")
	}

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return "", nil, errors.Wrap(err, "[dbr] Insert.Listeners.dispatch")
	}

	if len(b.Into) == 0 {
		return "", nil, errors.NewEmptyf(errTableMissing)
	}
	if len(b.Cols) == 0 && len(b.Maps) == 0 {
		return "", nil, errors.NewEmptyf(errColumnsMissing)
	} else if len(b.Maps) == 0 {
		if len(b.Vals) == 0 && len(b.Recs) == 0 {
			return "", nil, errors.NewEmptyf(errRecordsMissing)
		}
		if len(b.Cols) == 0 && (len(b.Vals) > 0 || len(b.Recs) > 0) {
			return "", nil, errors.NewEmptyf(errColumnsMissing)
		}
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	buf.WriteString("INSERT INTO ")
	Quoter.quote(buf, b.Into)
	buf.WriteString(" (")

	if len(b.Maps) != 0 {
		args, err := b.mapToSQL(buf)
		return buf.String(), args, err
	}

	var ph = bufferpool.Get() // Build the ph like "(?,?,?)"

	// Simultaneously write the cols to the sql buffer, and build a ph
	ph.WriteRune('(')
	for i, c := range b.Cols {
		if i > 0 {
			buf.WriteRune(',')
			ph.WriteRune(',')
		}
		Quoter.quoteAs(buf, c)
		ph.WriteRune('?')
	}
	buf.WriteString(") VALUES ")
	ph.WriteRune(')')
	placeholderStr := ph.String()
	bufferpool.Put(ph)

	// Go thru each value we want to insert. Write the placeholders, and collect args
	for i := 0; i < len(b.Vals); i = i + len(b.Cols) {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(placeholderStr)
	}

	if b.Recs == nil {
		return buf.String(), b.Vals, nil
	}

	args := make(Arguments, len(b.Vals), len(b.Vals)+len(b.Recs)) // sneaky ;-)
	copy(args, b.Vals)

	for i, rec := range b.Recs {
		a2, err := rec.GenerateArguments(StatementTypeInsert, b.Cols, nil)
		if err != nil {
			return "", nil, errors.Wrap(err, "[dbr] Insert.ToSQL.Record")
		}

		args = append(args, a2...)
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(placeholderStr)
	}

	return buf.String(), args, nil
}

// mapToSQL serialized the Insert to a SQL string
// It goes through the Maps param and combined its keys/values into the SQL query string
// It returns the string with placeholders and a slice of query arguments
func (b *Insert) mapToSQL(w queryWriter) (Arguments, error) {
	if b.previousError != nil {
		return nil, errors.Wrap(b.previousError, "[dbr] Insert.ToSQL")
	}

	keys := make([]string, len(b.Maps))
	vals := make(Arguments, len(b.Maps))
	i := 0
	for k, v := range b.Maps {
		keys[i] = k
		vals[i] = v
		i++
	}
	var args Arguments
	var placeholder = bufferpool.Get() // Build the placeholder like "(?,?,?)"
	defer bufferpool.Put(placeholder)

	placeholder.WriteRune('(')
	for i, c := range keys {
		if i > 0 {
			w.WriteRune(',')
			placeholder.WriteRune(',')
		}
		Quoter.quoteAs(w, c)
		placeholder.WriteRune('?')
	}
	w.WriteString(") VALUES ")
	placeholder.WriteRune(')')
	w.WriteString(placeholder.String())

	args = append(args, vals...)

	return args, nil
}

// Exec executes the statement represented by the Insert
// It returns the raw database/sql Result and an error if there was one.
// Regarding LastInsertID(): If you insert multiple rows using a single
// INSERT statement, LAST_INSERT_ID() returns the value generated for
// the first inserted row only. The reason for this at to make it possible to
// reproduce easily the same INSERT statement against some other server.
func (b *Insert) Exec(ctx context.Context) (sql.Result, error) {
	sql, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.Exec.ToSQL")
	}

	fullSQL, err := Preprocess(sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.Exec.Preprocess")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Insert.Exec.Timing", log.String("sql", fullSQL))
	}

	result, err := b.DB.ExecContext(ctx, fullSQL)
	if err != nil {
		return result, errors.Wrap(err, "[dbr] Insert.Exec.Exec")
	}

	return result, nil
}

// Prepare creates a prepared statement
func (b *Insert) Prepare(ctx context.Context) (*sql.Stmt, error) {
	rawSQL, _, err := b.ToSQL() // TODO create a ToSQL version without any arguments
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.Exec.ToSQL")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Insert.Prepare.Timing", log.String("sql", rawSQL))
	}

	stmt, err := b.DB.PrepareContext(ctx, rawSQL)
	return stmt, errors.Wrap(err, "[dbr] Insert.Prepare.Prepare")
}
