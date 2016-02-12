package dbr

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/juju/errors"
)

// InsertBuilder contains the clauses for an INSERT statement
type InsertBuilder struct {
	*Session
	runner

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
		Session: sess,
		runner:  sess.cxn.DB,
		Into:    into,
	}
}

// InsertInto instantiates a InsertBuilder for the given table bound to a transaction
func (tx *Tx) InsertInto(into string) *InsertBuilder {
	return &InsertBuilder{
		Session: tx.Session,
		runner:  tx.Tx,
		Into:    into,
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
		b.EventErrKv("dbr.insertbuilder.values", err, kvs{"args": fmt.Sprint(vals)})
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
			panic(err)
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
		panic("pair only allows you to specify 1 record to insret")
	}
	return b
}

// ToSql serialized the InsertBuilder to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *InsertBuilder) ToSql() (string, []interface{}, error) {
	if len(b.Into) == 0 {
		return "", nil, ErrMissingTable
	}
	if len(b.Cols) == 0 && len(b.Maps) == 0 {
		return "", nil, errors.New("no columns or map specified")
	} else if len(b.Maps) == 0 {
		if len(b.Vals) == 0 && len(b.Recs) == 0 {
			return "", nil, errors.New("no values or records specified")
		}
		if len(b.Cols) == 0 && (len(b.Vals) > 0 || len(b.Recs) > 0) {
			return "", nil, errors.New("no columns specified")
		}
	}

	var sql = bufferpool.Get()

	sql.WriteString("INSERT INTO ")
	sql.WriteString(b.Into)
	sql.WriteString(" (")

	if len(b.Maps) != 0 {
		return b.MapToSql(sql)
	}
	defer bufferpool.Put(sql)

	var args []interface{}
	var placeholder = bufferpool.Get() // Build the placeholder like "(?,?,?)"
	defer bufferpool.Put(placeholder)

	// Simulataneously write the cols to the sql buffer, and build a placeholder
	placeholder.WriteRune('(')
	for i, c := range b.Cols {
		if i > 0 {
			sql.WriteRune(',')
			placeholder.WriteRune(',')
		}
		Quoter.writeQuotedColumn(c, sql)
		placeholder.WriteRune('?')
	}
	sql.WriteString(") VALUES ")
	placeholder.WriteRune(')')
	placeholderStr := placeholder.String()

	// Go thru each value we want to insert. Write the placeholders, and collect args
	for i, row := range b.Vals {
		if i > 0 {
			sql.WriteRune(',')
		}
		sql.WriteString(placeholderStr)

		for _, v := range row {
			args = append(args, v)
		}
	}
	anyVals := len(b.Vals) > 0

	// Go thru the records. Write the placeholders, and do reflection on the records to extract args
	for i, rec := range b.Recs {
		if i > 0 || anyVals {
			sql.WriteRune(',')
		}
		sql.WriteString(placeholderStr)

		ind := reflect.Indirect(reflect.ValueOf(rec))
		vals, err := b.valuesFor(ind.Type(), ind, b.Cols)
		if err != nil {
			return "", nil, errors.Mask(err)
		}
		for _, v := range vals {
			args = append(args, v)
		}
	}

	return sql.String(), args, nil
}

// MapToSql serialized the InsertBuilder to a SQL string
// It goes through the Maps param and combined its keys/values into the SQL query string
// It returns the string with placeholders and a slice of query arguments
func (b *InsertBuilder) MapToSql(sql *bytes.Buffer) (string, []interface{}, error) {
	defer bufferpool.Put(sql)
	keys := make([]string, len(b.Maps))
	vals := make([]interface{}, len(b.Maps))
	i := 0
	for k, v := range b.Maps {
		keys[i] = k
		if dbVal, ok := v.(driver.Valuer); ok {
			if val, err := dbVal.Value(); err == nil {
				vals[i] = val
			} else {
				return "", nil, errors.Mask(err)
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
			sql.WriteRune(',')
			placeholder.WriteRune(',')
		}
		Quoter.writeQuotedColumn(c, sql)
		placeholder.WriteRune('?')
	}
	sql.WriteString(") VALUES ")
	placeholder.WriteRune(')')
	sql.WriteString(placeholder.String())

	for _, row := range vals {
		args = append(args, row)
	}

	return sql.String(), args, nil
}

// Exec executes the statement represented by the InsertBuilder
// It returns the raw database/sql Result and an error if there was one.
// Regarding LastInsertID(): If you insert multiple rows using a single
// INSERT statement, LAST_INSERT_ID() returns the value generated for
// the first inserted row only. The reason for this is to make it possible to
// reproduce easily the same INSERT statement against some other server.
func (b *InsertBuilder) Exec() (sql.Result, error) {
	sql, args, err := b.ToSql()
	if err != nil {
		return nil, b.EventErrKv("dbr.insert.exec.tosql", err, nil)
	}

	fullSql, err := Preprocess(sql, args)
	if err != nil {
		return nil, b.EventErrKv("dbr.insert.exec.interpolate", err, kvs{"sql": sql, "args": fmt.Sprint(args)})
	}

	// Start the timer:
	startTime := time.Now()
	defer func() { b.TimingKv("dbr.insert", time.Since(startTime).Nanoseconds(), kvs{"sql": fullSql}) }()

	result, err := b.runner.Exec(fullSql)
	if err != nil {
		return result, b.EventErrKv("dbr.insert.exec.exec", err, kvs{"sql": fullSql})
	}

	// If the structure has an "Id" field which is an int64, set it from the LastInsertId(). Otherwise, don't bother.
	if len(b.Recs) == 1 {
		rec := b.Recs[0]
		val := reflect.Indirect(reflect.ValueOf(rec))
		if val.Kind() == reflect.Struct && val.CanSet() {
			// @todo important: make Id configurable to match all magento internal ID columns
			if idField := val.FieldByName("Id"); idField.IsValid() && idField.Kind() == reflect.Int64 {
				if lastID, err := result.LastInsertId(); err == nil {
					idField.Set(reflect.ValueOf(lastID))
				} else {
					b.EventErrKv("dbr.insert.exec.last_inserted_id", err, kvs{"sql": fullSql})
				}
			}
		}
	}

	return result, nil
}
