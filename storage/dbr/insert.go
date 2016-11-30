package dbr

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
)

// InsertBuilder contains the clauses for an INSERT statement
type InsertBuilder struct {
	log.Logger
	Execer

	Into string
	Cols []string
	Vals [][]interface{}
	Recs []interface{}
	Maps map[string]interface{}
}

var _ queryBuilder = (*InsertBuilder)(nil)

// InsertInto instantiates a InsertBuilder for the given table
func (sess *Session) InsertInto(into string) *InsertBuilder {
	return &InsertBuilder{
		Logger: sess.Logger,
		Execer: sess.cxn.DB,
		Into:   into,
	}
}

// InsertInto instantiates a InsertBuilder for the given table bound to a transaction
func (tx *Tx) InsertInto(into string) *InsertBuilder {
	return &InsertBuilder{
		Logger: tx.Logger,
		Execer: tx.Tx,
		Into:   into,
	}
}

// Columns appends columns to insert in the statement
func (b *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	b.Cols = columns
	return b
}

// Values appends a set of values to the statement.
// Pro Tip: Use Values() and not Record() to avoid reflection.
// Only this function will consider the driver.Valuer interface when you pass
// a pointer to the value.
func (b *InsertBuilder) Values(vals ...interface{}) *InsertBuilder {
	if err := argsValuer(&vals); err != nil {
		// b.EventErrKv("dbr.insertbuilder.values", err, kvs{"args": fmt.Sprint(vals)})
		panic(err) // todo remove panic
	}
	b.Vals = append(b.Vals, vals)
	return b
}

// Record pulls in values to match Columns from the record. Uses reflection.
func (b *InsertBuilder) Record(record interface{}) *InsertBuilder {
	b.Recs = append(b.Recs, record)
	return b
}

// Record pulls in values to match Columns from the record
func (b *InsertBuilder) Map(m map[string]interface{}) *InsertBuilder {
	b.Maps = m
	return b
}

// Pair adds a key/value pair to the statement. Uses not reflection.
func (b *InsertBuilder) Pair(column string, value interface{}) *InsertBuilder {
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

// ToSql serialized the InsertBuilder to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *InsertBuilder) ToSql() (string, []interface{}, error) {
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
		return b.MapToSql(buf)
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

// MapToSql serialized the InsertBuilder to a SQL string
// It goes through the Maps param and combined its keys/values into the SQL query string
// It returns the string with placeholders and a slice of query arguments
func (b *InsertBuilder) MapToSql(w QueryWriter) (string, []interface{}, error) {

	keys := make([]string, len(b.Maps))
	vals := make([]interface{}, len(b.Maps))
	i := 0
	for k, v := range b.Maps {
		keys[i] = k
		if dbVal, ok := v.(driver.Valuer); ok {
			if val, err := dbVal.Value(); err == nil {
				vals[i] = val
			} else {
				return "", nil, errors.Wrap(err, "[dbr] MapToSql -> driver.Valuer")
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

// Exec executes the statement represented by the InsertBuilder
// It returns the raw database/sql Result and an error if there was one.
// Regarding LastInsertID(): If you insert multiple rows using a single
// INSERT statement, LAST_INSERT_ID() returns the value generated for
// the first inserted row only. The reason for this is to make it possible to
// reproduce easily the same INSERT statement against some other server.
func (b *InsertBuilder) Exec(ctx context.Context) (sql.Result, error) {
	sql, args, err := b.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] insert.exec.tosql")
	}

	fullSql, err := Preprocess(sql, args)
	if err != nil {
		return nil, errors.Wrap(err, "dbr.insert.exec.interpolate")
	}

	if b.Logger.IsInfo() {
		defer log.WhenDone(b.Logger).Info("dbr.InsertBuilder.ExecContext.timing", log.String("sql", fullSql))
	}

	result, err := b.ExecContext(ctx, fullSql)
	if err != nil {
		return result, errors.Wrap(err, "[dbr] insert.exec.exec")
	}

	return result, nil
}
