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

package dml

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

// https://www.adampalmer.me/iodigitalsec/2013/08/18/mysql_real_escape_string-wont-magically-solve-your-sql-injection-problems/

const (
	sqlStrNullUC = "NULL"
	sqlStrNullLC = "null"
	sqlStar      = "*"
)

var (
	sqlBytesNullUC = []byte(sqlStrNullUC)
	sqlBytesNullLC = []byte(sqlStrNullLC)
)

// QualifiedRecord is a ColumnMapper with a qualifier. A QualifiedRecord gets
// used as arguments to ExecRecord or BindRecord in the SQL statement. If you
// use an alias for the main table/view you must set the alias as the qualifier.
type QualifiedRecord struct {
	_Named_Fields_Required struct{}

	// Qualifier is the name of the table or view or procedure or can be their
	// alias. It must be a valid MySQL/MariaDB identifier.
	//
	// If empty, the main table or its alias of a query will be used. We call it
	// the default qualifier. Each query can only contain one default qualifier.
	// If you provide multiple default qualifier, the last one wins and
	// overwrites the previous.
	Qualifier string
	Record    ColumnMapper
}

// Qualify provides a more concise way to create QualifiedRecord values.
func Qualify(q string, record ColumnMapper) QualifiedRecord {
	return QualifiedRecord{Qualifier: q, Record: record}
}

// argument is union type for different Go primitives and their slice
// representation. argument must be used as a pointer because it slows
// everything down. Check the benchmarks.
type argument struct {
	// isSet indicates if an argument is really set, because `argument` gets
	// used as an embedded non-pointer type in type Condition.
	isSet bool
	// name for named place holders sql.NamedArg. Write a converter to and from
	// sql.NamedArg
	name  string
	value interface{}
}

type placeHolder int8

func (arg *argument) set(v interface{}) {
	arg.isSet = true
	arg.value = v
}

func (arg *argument) len() (l int) {
	switch v := arg.value.(type) {
	case nil, int, int64, uint64, float64, bool, string, []byte, time.Time, placeHolder:
		l = 1
	case []int:
		l = len(v)
	case []int64:
		l = len(v)
	case []uint64:
		l = len(v)
	case []uint:
		l = len(v)
	case []float64:
		l = len(v)
	case []bool:
		l = len(v)
	case []string:
		l = len(v)
	case [][]byte:
		l = len(v)
	case []time.Time:
		l = len(v)
	case []NullString:
		l = len(v)
	case []NullInt64:
		l = len(v)
	case []NullFloat64:
		l = len(v)
	case []NullBool:
		l = len(v)
	case []NullTime:
		l = len(v)
	default:
		panic("Unsupported type")
	}
	// default is 0
	return
}

func (arg argument) writeTo(w *bytes.Buffer, pos int) (err error) {
	switch v := arg.value.(type) {
	case int:
		err = writeInt64(w, int64(v))
	case []int:
		err = writeInt64(w, int64(v[pos]))
	case int64:
		err = writeInt64(w, v)
	case []int64:
		err = writeInt64(w, v[pos])
	case []NullInt64:
		if s := v[pos]; s.Valid {
			return writeInt64(w, s.Int64)
		}
		_, err = w.WriteString(sqlStrNullUC)

	case uint64:
		err = writeUint64(w, v)
	case []uint64:
		err = writeUint64(w, v[pos])
	case []uint:
		err = writeUint64(w, uint64(v[pos]))

	case float64:
		err = writeFloat64(w, v)
	case []float64:
		err = writeFloat64(w, v[pos])
	case []NullFloat64:
		if s := v[pos]; s.Valid {
			return writeFloat64(w, s.Float64)
		}
		_, err = w.WriteString(sqlStrNullUC)

	case bool:
		dialect.EscapeBool(w, v)
	case []bool:
		dialect.EscapeBool(w, v[pos])
	case []NullBool:
		if s := v[pos]; s.Valid {
			dialect.EscapeBool(w, s.Bool)
			return nil
		}
		_, err = w.WriteString(sqlStrNullUC)

		// TODO(CyS) Cut the printed string in errors if it's longer than XX chars
	case string:
		if !utf8.ValidString(v) {
			return errors.NewNotValidf("[dml] Argument.WriteTo: String is not UTF-8: %q", v)
		}
		dialect.EscapeString(w, v)
	case []string:
		if !utf8.ValidString(v[pos]) {
			return errors.NewNotValidf("[dml] Argument.WriteTo: String is not UTF-8: %q", v[pos])
		}
		dialect.EscapeString(w, v[pos])
	case []NullString:
		if s := v[pos]; s.Valid {
			if !utf8.ValidString(s.String) {
				return errors.NewNotValidf("[dml] Argument.WriteTo: String is not UTF-8: %q", s.String)
			}
			dialect.EscapeString(w, s.String)
		} else {
			_, err = w.WriteString(sqlStrNullUC)
		}

	case []byte:
		err = writeBytes(w, v)

	case [][]byte:
		err = writeBytes(w, v[pos])

	case time.Time:
		dialect.EscapeTime(w, v)
	case []time.Time:
		dialect.EscapeTime(w, v[pos])
	case []NullTime:
		if nt := v[pos]; nt.Valid {
			dialect.EscapeTime(w, nt.Time)
		} else {
			_, err = w.WriteString(sqlStrNullUC)
		}

	case nil:
		_, err = w.WriteString(sqlStrNullUC)
	case placeHolder:
		err = w.WriteByte(placeHolderRune)

	default:
		panic(errors.NewNotSupportedf("[dml] Unsupported field type: %d", arg.value))
	}
	return err
}

func (arg argument) GoString() string {
	buf := new(bytes.Buffer)
	if arg.name != "" {
		fmt.Fprintf(buf, ".Name(%q)", arg.name)
	}
	switch v := arg.value.(type) {
	case int:
		fmt.Fprintf(buf, ".Int(%d)", v)
	case []int:
		fmt.Fprintf(buf, ".Ints(%#v...)", v)

	case int64:
		fmt.Fprintf(buf, ".Int64(%d)", v)
	case []int64:
		fmt.Fprintf(buf, ".Int64s(%#v...)", v)
	case []NullInt64:
		buf.WriteString(".NullInt64(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(nv.GoString())
		}
		buf.WriteString(")")

	case uint64:
		fmt.Fprintf(buf, ".Uint64(%d)", v)
	case []uint64:
		fmt.Fprintf(buf, ".Uint64s(%#v...)", v)
	case []uint:
		fmt.Fprintf(buf, ".Uints(%#v...)", v)

	case float64:
		fmt.Fprintf(buf, ".Float64(%f)", v)
	case []float64:
		fmt.Fprintf(buf, ".Float64s(%#v...)", v) // the lazy way; prints `[]float64{2.76, 3.141}...` but should `2.76, 3.141`
	case []NullFloat64:
		buf.WriteString(".NullFloat64(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(nv.GoString())
		}
		buf.WriteString(")")

	case bool:
		fmt.Fprintf(buf, ".Bool(%v)", v)
	case []bool:
		fmt.Fprintf(buf, ".Bools(%#v...)", v)
	case []NullBool:
		buf.WriteString(".NullBool(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(nv.GoString())
		}
		buf.WriteString(")")

	case string:
		fmt.Fprintf(buf, ".String(%q)", v)
	case []string:
		buf.WriteString(".Strings(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "%q", nv)
		}
		buf.WriteString(")")
	case []NullString:
		buf.WriteString(".NullString(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(nv.GoString())
		}
		buf.WriteString(")")

	case []byte:
		fmt.Fprintf(buf, ".Bytes(%#v)", v)
	case [][]byte:
		buf.WriteString(".BytesSlice(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "%#v", nv)
		}
		buf.WriteString(")")

	case time.Time:
		fmt.Fprintf(buf, ".Time(time.Unix(%d,%d))", v.Unix(), v.Nanosecond())
	case []time.Time:
		buf.WriteString(".Times(")
		for i, t := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "time.Unix(%d,%d)", t.Unix(), t.Nanosecond())
		}
		buf.WriteString(")")
	case []NullTime:
		buf.WriteString(".NullTime(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(nv.GoString())
		}
		buf.WriteString(")")

	case nil:
		fmt.Fprint(buf, ".Null()")
	case placeHolder:
		// noop
	default:
		panic(errors.NewNotSupportedf("[dml] Unsupported field type: %T", arg.value))
	}
	return buf.String()
}

// Arguments a collection of primitive types or slices of primitive types. The
// method receiver functions have the same names as in type RowConvert.
type Arguments []argument

// MakeArgs creates a new argument slice with the desired capacity.
func MakeArgs(cap int) Arguments {
	return make(Arguments, 0, cap)
}

// MapColumns allows to merge one argument slice with another depending on the
// matched columns. Each argument in the slice must be a named argument.
// Implements interface ColumnMapper.
func (a Arguments) MapColumns(cm *ColumnMap) error {
	if cm.Mode() == ColumnMapEntityReadAll {
		cm.Args = append(cm.Args, a...)
		return cm.Err()
	}
	for cm.Next() {
		// now a bit slow ... but will be refactored later with constant time
		// access, but first benchmark it. This for loop can be the 3rd one in the
		// overall chain.
		c := cm.Column()
		for _, arg := range a {
			// Case sensitive comparison
			if c != "" && arg.name == c {
				cm.Args = append(cm.Args, arg)
				break
			}
		}
	}
	return cm.Err()
}

func (a Arguments) GoString() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "dml.MakeArgs(%d)", len(a))
	for _, arg := range a {
		buf.WriteString(arg.GoString())
	}
	return buf.String()
}

// Len returns the total length of all arguments.
func (a Arguments) Len() int {
	var l int
	for _, arg := range a {
		l += arg.len()
	}
	return l
}

// Write writes all arguments into buf and separates by a comma.
func (a Arguments) Write(buf *bytes.Buffer) error {
	buf.WriteByte('(')
	for j, arg := range a {
		l := arg.len()
		for i := 0; i < l; i++ {
			if i > 0 || j > 0 {
				buf.WriteByte(',')
			}
			if err := arg.writeTo(buf, i); err != nil {
				return errors.Wrapf(err, "[dml] args write failed at pos %d with argument %#v", j, arg)
			}
		}
	}
	return buf.WriteByte(')')
}

// Interfaces creates an interface slice with flatend values. Each type is one
// of the allowed types in driver.Value.
func (a Arguments) Interfaces(args ...interface{}) []interface{} {
	const maxInt64 = 1<<63 - 1
	if len(a) == 0 {
		return nil
	}
	if args == nil {
		args = make([]interface{}, 0, 2*len(a))
	}

	for _, arg := range a {
		switch vv := arg.value.(type) {

		case bool, string, []byte, time.Time, float64, int64, nil:
			args = append(args, arg.value)

		case int:
			args = append(args, int64(vv))
		case []int:
			for _, v := range vv {
				args = append(args, int64(v))
			}

		case []int64:
			for _, v := range vv {
				args = append(args, v)
			}
		case []NullInt64:
			for _, v := range vv {
				if v.Valid {
					args = append(args, v.Int64)
				} else {
					args = append(args, nil)
				}
			}

			// Get send as text in a byte slice. The MySQL/MariaDB Server type
			// casts it into a bigint. If you change this, a test will fail.
		case uint64:
			if vv > maxInt64 {
				args = append(args, strconv.AppendUint([]byte{}, vv, 10))
			} else {
				args = append(args, int64(vv))
			}

		case []uint64:
			for _, v := range vv {
				if v > maxInt64 {
					args = append(args, strconv.AppendUint([]byte{}, v, 10))
				} else {
					args = append(args, int64(v))
				}
			}
		case []uint:
			for _, v := range vv {
				if v > maxInt64 {
					args = append(args, strconv.AppendUint([]byte{}, uint64(v), 10))
				} else {
					args = append(args, int64(v))
				}
			}

		case []float64:
			for _, v := range vv {
				args = append(args, v)
			}
		case []NullFloat64:
			for _, v := range vv {
				if v.Valid {
					args = append(args, v.Float64)
				} else {
					args = append(args, nil)
				}
			}

		case []bool:
			for _, v := range vv {
				args = append(args, v)
			}
		case []NullBool:
			for _, v := range vv {
				if v.Valid {
					args = append(args, v.Bool)
				} else {
					args = append(args, nil)
				}
			}

		case []string:
			for _, v := range vv {
				args = append(args, v)
			}
		case []NullString:
			for _, v := range vv {
				if v.Valid {
					args = append(args, v.String)
				} else {
					args = append(args, nil)
				}
			}

		case [][]byte:
			for _, v := range vv {
				args = append(args, v)
			}

		case []time.Time:
			for _, v := range vv {
				args = append(args, v)
			}
		case []NullTime:
			for _, v := range vv {
				if v.Valid {
					args = append(args, v.Time)
				} else {
					args = append(args, nil)
				}
			}
		}
	}
	return args
}

func (a Arguments) add(v interface{}) Arguments {
	if l := len(a); l > 0 {
		// look back if there might be a name.
		if arg := a[l-1]; !arg.isSet {
			// The previous call Name() has set the name and now we set the
			// value, but don't append a new entry.
			arg.isSet = true
			arg.value = v
			a[l-1] = arg
			return a
		}
	}
	return append(a, argument{isSet: true, value: v})
}
func (a Arguments) PlaceHolder() Arguments                  { return a.add(placeHolder(1)) }
func (a Arguments) Null() Arguments                         { return a.add(nil) }
func (a Arguments) Int(i int) Arguments                     { return a.add(i) }
func (a Arguments) Ints(i ...int) Arguments                 { return a.add(i) }
func (a Arguments) Int64(i int64) Arguments                 { return a.add(i) }
func (a Arguments) Int64s(i ...int64) Arguments             { return a.add(i) }
func (a Arguments) Uint(i uint) Arguments                   { return a.add(uint64(i)) }
func (a Arguments) Uints(i ...uint) Arguments               { return a.add(i) }
func (a Arguments) Uint64(i uint64) Arguments               { return a.add(i) }
func (a Arguments) Uint64s(i ...uint64) Arguments           { return a.add(i) }
func (a Arguments) Float64(f float64) Arguments             { return a.add(f) }
func (a Arguments) Float64s(f ...float64) Arguments         { return a.add(f) }
func (a Arguments) Bool(b bool) Arguments                   { return a.add(b) }
func (a Arguments) Bools(b ...bool) Arguments               { return a.add(b) }
func (a Arguments) String(s string) Arguments               { return a.add(s) }
func (a Arguments) Strings(s ...string) Arguments           { return a.add(s) }
func (a Arguments) Time(t time.Time) Arguments              { return a.add(t) }
func (a Arguments) Times(t ...time.Time) Arguments          { return a.add(t) }
func (a Arguments) Bytes(b []byte) Arguments                { return a.add(b) }
func (a Arguments) BytesSlice(b ...[]byte) Arguments        { return a.add(b) }
func (a Arguments) NullString(nv ...NullString) Arguments   { return a.add(nv) }
func (a Arguments) NullFloat64(nv ...NullFloat64) Arguments { return a.add(nv) }
func (a Arguments) NullInt64(nv ...NullInt64) Arguments     { return a.add(nv) }
func (a Arguments) NullBool(nv ...NullBool) Arguments       { return a.add(nv) }
func (a Arguments) NullTime(nv ...NullTime) Arguments       { return a.add(nv) }

// Name sets the name for the following argument. Calling Name two time after
// each other sets the first call to Name to a NULL value. A call to Name should
// always follow a call to a function type like Int, Float64s or NullTime.
func (a Arguments) Name(n string) Arguments { return append(a, argument{name: n}) }

// Reset resets the slice for new usage retaining the already allocated memory.
func (a Arguments) Reset() Arguments { return a[:0] }

// DriverValue adds multiple of the same underlying values to the argument
// slice. When using different values, the last applied value wins and gets
// added to the argument slice.
func (a Arguments) DriverValue(dvs ...driver.Valuer) Arguments {
	// value is a value that drivers must be able to handle.
	// It is either nil or an instance of one of these types:
	//
	//   int64
	//   float64
	//   bool
	//   []byte
	//   string
	//   time.Time
	var arg argument
	var i64s []int64
	var f64s []float64
	var bs []bool
	var bytess [][]byte
	var strs []string
	var times []time.Time
	for _, dv := range dvs {
		// dv cannot be nil
		v, err := dv.Value()
		if err != nil {
			// TODO: Either keep panic or delay the error until another function gets called which also returns an error.
			panic(errors.NewFatal(err, "[dml] Driver.value error for %#v", dv))
		}

		switch t := v.(type) {
		case nil:
		case int64:
			i64s = append(i64s, t)
		case float64:
			f64s = append(f64s, t)
		case bool:
			bs = append(bs, t)
		case []byte:
			bytess = append(bytess, t)
		case string:
			strs = append(strs, t)
		case time.Time:
			times = append(times, t)
		default:
			panic(errors.NewNotSupportedf("[dml] Type %#v not supported in value slice: %#v", t, dvs))
		}
	}

	arg.isSet = true

	switch {
	case len(i64s) > 0:
		arg.value = i64s
	case len(f64s) > 0:
		arg.value = f64s
	case len(bs) > 0:
		arg.value = bs
	case len(bytess) > 0:
		arg.value = bytess
	case len(strs) > 0:
		arg.value = strs
	case len(times) > 0:
		arg.value = times
	}

	a = append(a, arg)
	return a
}

// DriverValues adds each driver.value as its own argument to the argument slice.
func (a Arguments) DriverValues(dvs ...driver.Valuer) Arguments {
	// value is a value that drivers must be able to handle.
	// It is either nil or an instance of one of these types:
	//
	//   int64
	//   float64
	//   bool
	//   []byte
	//   string
	//   time.Time
	for _, dv := range dvs {
		if dv == nil {
			a = a.Null()
			continue
		}
		v, err := dv.Value()
		if err != nil {
			// TODO: Either keep panic or delay the error until another function gets called which also returns an error.
			panic(errors.NewFatal(err, "[dml] Driver.value error for %#v", dv))
		}
		switch t := v.(type) {
		case nil:
			a = a.Null()
		case int64:
			a = a.Int64(t)
		case float64:
			a = a.Float64(t)
		case bool:
			a = a.Bool(t)
		case []byte:
			a = a.Bytes(t)
		case string:
			a = a.String(t)
		case time.Time:
			a = a.Times(t)
		default:
			panic(errors.NewNotSupportedf("[dml] Type %#v not supported in value slice: %#v", t, dvs))
		}
	}
	return a
}

func iFaceToArgs(values ...interface{}) Arguments {
	args := make(Arguments, 0, len(values))
	for _, val := range values {
		switch v := val.(type) {
		case float32:
			args = args.Float64(float64(v))
		case float64:
			args = args.Float64(v)
		case int64:
			args = args.Int64(v)
		case int:
			args = args.Int64(int64(v))
		case int32:
			args = args.Int64(int64(v))
		case int16:
			args = args.Int64(int64(v))
		case int8:
			args = args.Int64(int64(v))
		case uint32:
			args = args.Int64(int64(v))
		case uint16:
			args = args.Int64(int64(v))
		case uint8:
			args = args.Int64(int64(v))
		case bool:
			args = args.Bool(v)
		case string:
			args = args.String(v)
		case []byte:
			args = args.Bytes(v)
		case time.Time:
			args = args.Time(v)
		case *time.Time:
			if v != nil {
				args = args.Time(*v)
			}
		case nil:
			args = args.Null()
		default:
			panic(errors.NewNotSupportedf("[dml] iFaceToArgs type %#v not yet supported", v))
		}
	}
	return args
}
