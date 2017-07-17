// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbr

import (
	"database/sql/driver"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

// https://www.adampalmer.me/iodigitalsec/2013/08/18/mysql_real_escape_string-wont-magically-solve-your-sql-injection-problems/

const (
	sqlStrNull = "NULL"
	sqlStar    = "*"
)

// SQL statement types and parts used as bit flag e.g. hint in
// ArgumentsAppender.AppendArguments.
const (
	SQLStmtInsert int = 1 << iota
	SQLStmtSelect
	SQLStmtUpdate
	SQLStmtDelete

	SQLPartJoin
	SQLPartWhere
	SQLPartHaving
	SQLPartSet
	SQLPartValues
)

// ArgumentsAppender assembles arguments for CRUD statements. The `stmtType`
// variable contains a flag from the constants SQLStmt* and SQLPart* to allow
// the knowledge in which case the function AppendArguments gets called. Any new
// arguments must be append to variable `args` and then returned. The readonly
// variable `columns` contains the name of the requested columns. E.g. if the
// first requested column names `id` then the first appended value must be an
// integer. Variable `columns` can additionally contain the names and/or
// expressions used in the WHERE, JOIN or HAVING clauses, if applicable for the
// SQL statement type. In case where stmtType has been set to
// SQLStmtInsert|SQLPartValues, the `columns` slice can be empty which means
// that all arguments are requested.
type ArgumentsAppender interface {
	AppendArguments(stmtType int, args Arguments, columns []string) (Arguments, error)
}

// Argument an interface which knows how to handle primitive arguments in case
// of interface conversion and transforming into the arguments textual
// representation. Later gets used during interpolation of SQL strings to avoid
// round trips to the database, use it wisely.
type Argument interface {
	// toIFace appends the argument or arguments to interface slice and returns it.
	toIFace([]interface{}) []interface{}
	// writeTo writes the value correctly escaped to the queryWriter. It must
	// avoid SQL injections.
	writeTo(w queryWriter, position int) error
	// len returns the length of the available arguments.
	len() int
}

// Arguments representing multiple Argument objects.
type Arguments []Argument

// len calculates the total length of all arguments
func (as Arguments) len() (tl int) {
	for _, a := range as {
		tl += a.len()
	}
	return
}

// Interfaces converts the underlying concrete types into an interface slice.
// Each entry in the interface is guaranteed to be one of the following arguments:
// []byte, bool, float64, int64, string or time.Time. Use driver.IsValue() for a
// check.
func (as Arguments) Interfaces() []interface{} {
	if len(as) == 0 {
		return nil
	}
	ret := make([]interface{}, 0, len(as))
	for _, a := range as {
		ret = a.toIFace(ret)
	}
	return ret
}

func iFaceToArgs(values ...interface{}) Arguments {
	args := make(Arguments, 0, len(values))
	for _, val := range values {
		switch v := val.(type) {
		case float32:
			args = append(args, Float64(float64(v)))
		case float64:
			args = append(args, Float64(v))
		case int64:
			args = append(args, Int64(v))
		case int:
			args = append(args, Int64(int64(v)))
		case int32:
			args = append(args, Int64(int64(v)))
		case int16:
			args = append(args, Int64(int64(v)))
		case int8:
			args = append(args, Int64(int64(v)))
		case uint32:
			args = append(args, Int64(int64(v)))
		case uint16:
			args = append(args, Int64(int64(v)))
		case uint8:
			args = append(args, Int64(int64(v)))
		case bool:
			args = append(args, Bool(v))
		case string:
			args = append(args, String(v))
		case []byte:
			args = append(args, Bytes(v))
		case time.Time:
			args = append(args, MakeTime(v))
		case *time.Time:
			if v != nil {
				args = append(args, MakeTime(*v))
			}
		case nil:
			args = append(args, NullValue())
		default:
			panic(errors.NewNotSupportedf("[dbr] iFaceToArgs type %#v not yet supported", v))
		}
	}
	return args
}

// DriverValues allows to use any type which implements driver.Valuer interface.
// Implements interface Argument.
type DriverValues []driver.Valuer

func (a DriverValues) toIFace(args []interface{}) []interface{} {
	for _, val := range a {
		v, err := val.Value()
		if err != nil {
			panic(err) // TODO(CyS) fix evil implementation of panic and remove panic
		}
		args = append(args, v)
	}
	return args
}

func writeDriverValuer(w queryWriter, v driver.Valuer) error {
	if v == nil {
		_, err := w.WriteString("NULL")
		return err
	}
	val, err := v.Value()
	if err != nil {
		return errors.Wrapf(err, "[dbr] DriverValues.WriteTo: %#v", v)
	}
	switch t := val.(type) {
	case nil:
		_, err = w.WriteString("NULL")
	case int:
		err = writeInt64(w, int64(t))
	case int64:
		err = writeInt64(w, t)
	case float64:
		err = writeFloat64(w, t)
	case string:
		dialect.EscapeString(w, t)
	case bool:
		dialect.EscapeBool(w, t)
	case time.Time:
		dialect.EscapeTime(w, t)
	case []byte:
		dialect.EscapeBinary(w, t)
	default:
		return errors.NewNotSupportedf("[dbr] DriverValues.WriteTo Type not yet supported: %#v", v)
	}
	return err
}

func (a DriverValues) writeTo(w queryWriter, pos int) error { return writeDriverValuer(w, a[pos]) }
func (a DriverValues) len() int                             { return len(a) }

// Times implements interface Argument to handle multiple time.Time arguments. The
// time.Time value gets correctly encoded in the MySQL/MariaDB format.
type Times []time.Time

func (a Times) toIFace(args []interface{}) []interface{} {
	for _, v := range a {
		args = append(args, v)
	}
	return args
}

func (a Times) writeTo(w queryWriter, pos int) error { dialect.EscapeTime(w, a[pos]); return nil }
func (a Times) len() int                             { return len(a) }

// Time implements interface Argument.
type Time struct{ time.Time }

func (a Time) toIFace(args []interface{}) []interface{} { return append(args, a.Time) }
func (a Time) writeTo(w queryWriter, _ int) error       { dialect.EscapeTime(w, a.Time); return nil }
func (a Time) len() int                                 { return 1 }
func (a Time) Value() (driver.Value, error) {
	return a.Time, nil
}

// MakeTime implements interface Argument and creates a new time value.
func MakeTime(t time.Time) Time { return Time{Time: t} }

// BytesSlice implements interface Argument. The slice can handle multiple
// []byte slices. Providing a nil returns a NULL type. Detects between valid
// UTF-8 strings and binary data. Later gets hex encoded.
type BytesSlice [][]byte

func (a BytesSlice) toIFace(args []interface{}) []interface{} {
	for _, v := range a {
		args = append(args, []byte(v))
	}
	return args
}

func (a BytesSlice) writeTo(w queryWriter, pos int) (err error) {
	if !utf8.Valid(a[pos]) {
		dialect.EscapeBinary(w, a[pos])
	} else {
		dialect.EscapeString(w, string(a[pos]))
	}
	return nil
}

func (a BytesSlice) len() int { return len(a) }

// ArgBytes implements interface Argument. Providing a nil returns a NULL type.
// Detects between valid UTF-8 strings and binary data. Later gets hex encoded.
type Bytes []byte

func (a Bytes) toIFace(args []interface{}) []interface{} { return append(args, []byte(a)) }

func (a Bytes) writeTo(w queryWriter, _ int) (err error) {
	if !utf8.Valid(a) {
		dialect.EscapeBinary(w, a)
	} else {
		dialect.EscapeString(w, string(a))
	}
	return nil
}

func (a Bytes) len() int { return 1 }

// Argument implements the driver Valuer interface.
func (a Bytes) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return []byte(a), nil
}

type nullValue rune

func (i nullValue) toIFace(args []interface{}) []interface{} { return append(args, nil) }
func (i nullValue) writeTo(w queryWriter, _ int) (err error) {
	_, err = w.WriteString("NULL")
	return err
}
func (i nullValue) len() int { return 1 }

// NullValue treats the argument as a SQL `IS NULL` or `NULL`. IN clause not
// supported. Implements interface Argument.
func NullValue() Argument { return nullValue(0) }

// String implements interface Argument. String also checks for valid UTF-8
// strings.
type String string

func (a String) toIFace(args []interface{}) []interface{} { return append(args, string(a)) }

func (a String) writeTo(w queryWriter, _ int) error {
	if !utf8.ValidString(string(a)) {
		return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", a)
	}
	dialect.EscapeString(w, string(a))
	return nil
}

func (a String) len() int { return 1 }
func (a String) Value() (driver.Value, error) {
	return string(a), nil
}

// Strings implements interface Argument and handles multiple string arguments.
// Strings also checks for valid UTF-8 strings.
type Strings []string

func (a Strings) toIFace(args []interface{}) []interface{} {
	for _, v := range a {
		args = append(args, v)
	}
	return args
}

func (a Strings) writeTo(w queryWriter, pos int) error {
	if !utf8.ValidString(a[pos]) {
		return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", a[pos])
	}
	dialect.EscapeString(w, a[pos])
	return nil
}

func (a Strings) len() int {
	return len(a)
}

// Bool implements interface Argument.
type Bool bool

func (a Bool) toIFace(args []interface{}) []interface{} { return append(args, a == true) }
func (a Bool) writeTo(w queryWriter, _ int) error       { dialect.EscapeBool(w, a == true); return nil }
func (a Bool) len() int                                 { return 1 }
func (a Bool) Value() (driver.Value, error) {
	return a == true, nil
}

// Bools implements interface Argument and handles multiple bool arguments.
type Bools []bool

func (a Bools) toIFace(args []interface{}) []interface{} {
	for _, v := range a {
		args = append(args, v)
	}
	return args
}

func (a Bools) writeTo(w queryWriter, pos int) error { dialect.EscapeBool(w, a[pos]); return nil }
func (a Bools) len() int                             { return len(a) }

// Int implements interface Argument.
type Int int

func (a Int) toIFace(args []interface{}) []interface{} { return append(args, int64(a)) }
func (a Int) writeTo(w queryWriter, _ int) error       { return writeInt64(w, int64(a)) }
func (a Int) len() int                                 { return 1 }
func (a Int) Value() (driver.Value, error) {
	return int64(a), nil
}

// Ints implements interface Argument and handles multiple int arguments.
type Ints []int

func (a Ints) toIFace(args []interface{}) []interface{} {
	for _, v := range a {
		args = append(args, int64(v))
	}
	return args
}

func (a Ints) writeTo(w queryWriter, pos int) error { return writeInt64(w, int64(a[pos])) }
func (a Ints) len() int                             { return len(a) }

// Int64 implements interface Argument.
type Int64 int64

func (a Int64) toIFace(args []interface{}) []interface{} { return append(args, int64(a)) }
func (a Int64) writeTo(w queryWriter, _ int) error       { return writeInt64(w, int64(a)) }
func (a Int64) len() int                                 { return 1 }
func (a Int64) Value() (driver.Value, error) {
	return int64(a), nil
}

// Uint64 implements interface Argument and driver.Argument, later get encoded
// in a byte slice. Full max uint64 support. The downside shows that the byte
// encoded uint64 gets transferred as a string to MySQL/MariaDB and the DB must
// type cast the string into a type bigint.
type Uint64 uint64

func (a Uint64) toIFace(args []interface{}) []interface{} {
	return append(args, strconv.AppendUint([]byte{}, uint64(a), 10))
}
func (a Uint64) writeTo(w queryWriter, _ int) error { return writeUint64(w, uint64(a)) }
func (a Uint64) len() int                           { return 1 }
func (a Uint64) Value() (driver.Value, error) {
	return strconv.AppendUint([]byte{}, uint64(a), 10), nil
}

// Int64s implements interface Argument and handles multiple int64 arguments.
type Int64s []int64

func (a Int64s) toIFace(args []interface{}) []interface{} {
	for _, v := range a {
		args = append(args, v)
	}
	return args
}

func (a Int64s) writeTo(w queryWriter, pos int) error { return writeInt64(w, a[pos]) }
func (a Int64s) len() int                             { return len(a) }

// Float64 implements interface Argument.
type Float64 float64

func (a Float64) toIFace(args []interface{}) []interface{} { return append(args, float64(a)) }

func (a Float64) writeTo(w queryWriter, _ int) error { return writeFloat64(w, float64(a)) }
func (a Float64) len() int                           { return 1 }
func (a Float64) Value() (driver.Value, error) {
	return float64(a), nil
}

// Float64s implements interface Argument and handles multiple float64 arguments.
type Float64s []float64

func (a Float64s) toIFace(args []interface{}) []interface{} {
	for _, v := range a {
		args = append(args, v)
	}
	return args
}

func (a Float64s) writeTo(w queryWriter, pos int) error { return writeFloat64(w, a[pos]) }
func (a Float64s) len() int                             { return len(a) }

type expr struct {
	SQL string
	Arguments
}

// ExpressionValue implements a SQL fragment with placeholders, and a slice of
// arguments to replace them with. Mostly used in UPDATE statements. Implements
// interface Argument.
func ExpressionValue(sql string, args ...Argument) Argument {
	return &expr{SQL: sql, Arguments: args}
}

func (e *expr) toIFace(args []interface{}) []interface{} {
	for _, a := range e.Arguments {
		args = a.toIFace(args)
	}
	return args
}

func (e *expr) writeTo(w queryWriter, _ int) error { w.WriteString(e.SQL); return nil }
func (e *expr) len() int                           { return 1 }

// placeHolderOp identifies place holder arguments. Those arguments will get
// assembled from an external type.
type placeHolderOp rune

// toIFace does not append anything because the placeHolderOp acts as an identifier.
func (i placeHolderOp) toIFace(args []interface{}) []interface{} { return args }
func (i placeHolderOp) len() int                                 { return 1 }
func (i placeHolderOp) writeTo(w queryWriter, _ int) (err error) {
	_, err = w.WriteString("?")
	return err
}
