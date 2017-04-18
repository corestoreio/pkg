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

	Into    string
	Columns []string
	Values  Arguments

	Records []ArgumentGenerater
	Maps    map[string]Argument

	// OnDuplicateKey updates the referenced columns. See documentation for type
	// `UpdatedColumns`. For more details
	// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
	OnDuplicateKey UpdatedColumns

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

// AddColumns appends columns to insert in the statement.
func (b *Insert) AddColumns(columns ...string) *Insert {
	b.Columns = append(b.Columns, columns...)
	return b
}

// AddValues appends a set of values to the statement. Pro Tip: Use Values() and
// not Record() to avoid reflection. Only this function will consider the
// driver.Valuer interface when you pass a pointer to the value. Values must be
// balanced to the number of columns. You can even provide more values, like
// records. see BenchmarkInsertValuesSQL
func (b *Insert) AddValues(vals ...Argument) *Insert {
	b.Values = append(b.Values, vals...)
	return b
}

// AddRecords pulls in values to match Columns from the record generator.
func (b *Insert) AddRecords(recs ...ArgumentGenerater) *Insert {
	b.Records = append(b.Records, recs...)
	return b
}

// AddOnDuplicateKey has some hidden features for best flexibility. You can only
// set the Columns itself to allow the following SQL construct:
//		`columnA`=VALUES(`columnA`)
// Means columnA gets automatically mapped to the VALUES column name.
func (b *Insert) AddOnDuplicateKey(column string, arg Argument) *Insert {
	b.OnDuplicateKey.Columns = append(b.OnDuplicateKey.Columns, column)
	b.OnDuplicateKey.Arguments = append(b.OnDuplicateKey.Arguments, arg)
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
	for _, c := range b.Columns {
		if c == column {
			b.previousError = errors.NewAlreadyExistsf("[dbr] Column %q has already been added", c)
			return b
		}
	}
	b.Columns = append(b.Columns, column)
	b.Values = append(b.Values, arg)
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

	sqlWriteInsertInto(buf, b.Into)
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
	if len(b.Columns) == 0 && len(b.Maps) == 0 {
		return "", nil, errors.NewEmptyf(errColumnsMissing)
	} else if len(b.Maps) == 0 {
		if len(b.Values) == 0 && len(b.Records) == 0 {
			return "", nil, errors.NewEmptyf(errRecordsMissing)
		}
		if len(b.Columns) == 0 && (len(b.Values) > 0 || len(b.Records) > 0) {
			return "", nil, errors.NewEmptyf(errColumnsMissing)
		}
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	sqlWriteInsertInto(buf, b.Into)
	buf.WriteString(" (")

	if len(b.Maps) != 0 {
		args, err := b.mapToSQL(buf)
		return buf.String(), args, err
	}

	var ph = bufferpool.Get() // Build the ph like "(?,?,?)"

	// Simultaneously write the cols to the sql buffer, and build a ph
	ph.WriteRune('(')
	for i, c := range b.Columns {
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
	for i := 0; i < len(b.Values); i = i + len(b.Columns) {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(placeholderStr)
	}

	args := make(Arguments, len(b.Values), len(b.Values)+len(b.Records)) // sneaky ;-)
	copy(args, b.Values)

	if b.Records == nil {
		if err := b.OnDuplicateKey.writeOnDuplicateKey(buf, &args); err != nil {
			return "", nil, errors.Wrap(err, "[dbr] Insert.OnDuplicateKey.writeOnDuplicateKey")
		}
		return buf.String(), args, nil
	}

	for i, rec := range b.Records {
		a2, err := rec.GenerateArguments(StatementTypeInsert, b.Columns, nil)
		if err != nil {
			return "", nil, errors.Wrap(err, "[dbr] Insert.ToSQL.Record")
		}

		args = append(args, a2...)
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(placeholderStr)
	}

	if err := b.OnDuplicateKey.writeOnDuplicateKey(buf, &args); err != nil {
		return "", nil, errors.Wrap(err, "[dbr] Insert.OnDuplicateKey.writeOnDuplicateKey")
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
