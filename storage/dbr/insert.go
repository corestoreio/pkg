package dbr

import (
	"database/sql"
	"database/sql/driver"
	"reflect"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Insert contains the clauses for an INSERT statement
type Insert struct {
	Log log.Logger // Log optional logger
	Execer
	Preparer

	Into string
	Cols []string
	Vals [][]interface{}
	Recs []interface{}
	Maps map[string]interface{}

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
}

// NewInsert creates a new object with a black hole logger.
func NewInsert(into string) *Insert {
	return &Insert{
		Log:  log.BlackHole{},
		Into: into,
	}
}

// InsertInto instantiates a Insert for the given table
func (sess *Session) InsertInto(into string) *Insert {
	return &Insert{
		Log:      sess.Logger,
		Execer:   sess.cxn.DB,
		Preparer: sess.cxn.DB,
		Into:     into,
	}
}

// InsertInto instantiates a Insert for the given table bound to a transaction
func (tx *Tx) InsertInto(into string) *Insert {
	return &Insert{
		Log:      tx.Logger,
		Execer:   tx.Tx,
		Preparer: tx.Tx,
		Into:     into,
	}
}

// Columns appends columns to insert in the statement.
func (b *Insert) Columns(columns ...string) *Insert {
	b.Cols = append(b.Cols, columns...)
	return b
}

// Values appends a set of values to the statement. Pro Tip: Use Values() and
// not Record() to avoid reflection. Only this function will consider the
// driver.Valuer interface when you pass a pointer to the value.
func (b *Insert) Values(vals ...interface{}) *Insert {
	if err := argsValuer(&vals); err != nil {
		panic(err) // todo remove panic
	}
	b.Vals = append(b.Vals, vals)
	return b
}

// Record pulls in values to match Columns from the record. Uses reflection.
func (b *Insert) Record(record interface{}) *Insert {
	b.Recs = append(b.Recs, record)
	return b
}

// Map pulls in values to match Columns from the record. Calling multiple
// times will add new map entries to the Insert map.
func (b *Insert) Map(m map[string]interface{}) *Insert {
	if b.Maps == nil {
		b.Maps = make(map[string]interface{})
	}
	for col, val := range m {
		b.Maps[col] = val
	}
	return b
}

// Pair adds a key/value pair to the statement. Uses not reflection.
func (b *Insert) Pair(column string, value interface{}) *Insert {
	if dbVal, ok := value.(driver.Valuer); ok {
		if val, err := dbVal.Value(); err == nil {
			value = val // overrides the current value ...
		} else {
			panic(err) // todo remove panic
		}
	}

	b.Cols = append(b.Cols, column)
	lenVals := len(b.Vals)
	if lenVals == 0 {
		args := []interface{}{value}
		b.Vals = [][]interface{}{args}
	} else if lenVals == 1 {
		b.Vals[0] = append(b.Vals[0], value)
	} else {
		panic("pair only allows you to specify 1 record to insret") // todo remove panic
	}
	return b
}

// ToSQL serialized the Insert to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Insert) ToSQL() (string, []interface{}, error) {

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
	buf.WriteString(b.Into)
	buf.WriteString(" (")

	if len(b.Maps) != 0 {
		return b.MapToSQL(buf)
	}

	var args []interface{}
	var placeholder = bufferpool.Get() // Build the placeholder like "(?,?,?)"
	defer bufferpool.Put(placeholder)

	// Simultaneously write the cols to the sql buffer, and build a placeholder
	placeholder.WriteRune('(')
	for i, c := range b.Cols {
		if i > 0 {
			buf.WriteRune(',')
			placeholder.WriteRune(',')
		}
		Quoter.writeQuotedColumn(c, buf)
		placeholder.WriteRune('?')
	}
	buf.WriteString(") VALUES ")
	placeholder.WriteRune(')')
	placeholderStr := placeholder.String()

	// Go thru each value we want to insert. Write the placeholders, and collect args
	for i, row := range b.Vals {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(placeholderStr)
		args = append(args, row...)
	}
	anyVals := len(b.Vals) > 0

	// Go thru the records. Write the placeholders, and do reflection on the records to extract args
	for i, rec := range b.Recs {
		if i > 0 || anyVals {
			buf.WriteRune(',')
		}
		buf.WriteString(placeholderStr)

		ind := reflect.Indirect(reflect.ValueOf(rec))
		vals, err := valuesFor(ind.Type(), ind, b.Cols)
		if err != nil {
			return "", nil, errors.Wrap(err, "[dbr] valuesFor")
		}
		args = append(args, vals...)
	}

	return buf.String(), args, nil
}

// MapToSQL serialized the Insert to a SQL string
// It goes through the Maps param and combined its keys/values into the SQL query string
// It returns the string with placeholders and a slice of query arguments
func (b *Insert) MapToSQL(w QueryWriter) (string, []interface{}, error) {

	keys := make([]string, len(b.Maps))
	vals := make([]interface{}, len(b.Maps))
	i := 0
	for k, v := range b.Maps {
		keys[i] = k
		if dbVal, ok := v.(driver.Valuer); ok {
			if val, err := dbVal.Value(); err == nil {
				vals[i] = val
			} else {
				return "", nil, errors.Wrap(err, "[dbr] MapToSQL -> driver.Valuer")
			}
		} else {
			vals[i] = v
		}
		i++
	}
	var args []interface{}
	var placeholder = bufferpool.Get() // Build the placeholder like "(?,?,?)"
	defer bufferpool.Put(placeholder)

	placeholder.WriteRune('(')
	for i, c := range keys {
		if i > 0 {
			w.WriteRune(',')
			placeholder.WriteRune(',')
		}
		Quoter.writeQuotedColumn(c, w)
		placeholder.WriteRune('?')
	}
	w.WriteString(") VALUES ")
	placeholder.WriteRune(')')
	w.WriteString(placeholder.String())

	args = append(args, vals...)

	return w.String(), args, nil
}

// Exec executes the statement represented by the Insert
// It returns the raw database/sql Result and an error if there was one.
// Regarding LastInsertID(): If you insert multiple rows using a single
// INSERT statement, LAST_INSERT_ID() returns the value generated for
// the first inserted row only. The reason for this is to make it possible to
// reproduce easily the same INSERT statement against some other server.
func (b *Insert) Exec() (sql.Result, error) {
	sql, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.Exec.ToSQL")
	}

	fullSQL, err := Preprocess(sql, args)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.Exec.Preprocess")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Insert.Exec.Timing", log.String("sql", fullSQL))
	}

	result, err := b.Execer.Exec(fullSQL)
	if err != nil {
		return result, errors.Wrap(err, "[dbr] Insert.Exec.Exec")
	}

	return result, nil
}

// Prepare creates a prepared statement
func (b *Insert) Prepare() (*sql.Stmt, error) {
	rawSQL, _, err := b.ToSQL() // TODO create a ToSQL version without any arguments
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.Exec.ToSQL")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Insert.Prepare.Timing", log.String("sql", rawSQL))
	}

	stmt, err := b.Preparer.Prepare(rawSQL)
	return stmt, errors.Wrap(err, "[dbr] Insert.Prepare.Prepare")
}
