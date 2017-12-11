// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	sqlBytesNullUC  = []byte(sqlStrNullUC)
	sqlBytesNullLC  = []byte(sqlStrNullLC)
	sqlBytesFalseLC = []byte("false")
	sqlBytesTrueLC  = []byte("true")
)

// QualifiedRecord is a ColumnMapper with a qualifier. A QualifiedRecord gets
// used as arguments to ExecRecord or WithRecords in the SQL statement. If you
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

func (arg *argument) set(v interface{}) {
	arg.isSet = true
	arg.value = v
}

func (arg *argument) len() (l int) {
	switch v := arg.value.(type) {
	case nil, int, int64, uint64, float64, bool, string, []byte, time.Time, NullString, NullInt64, NullFloat64, NullBool, NullTime:
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
		panic(errors.NewNotSupportedf("[dml] Unsupported type: %T => %#v", v, v))
	}
	// default is 0
	return
}

// writeTo mainly used in interpolate function
func (arg argument) writeTo(w *bytes.Buffer, pos uint) (err error) {
	if !arg.isSet {
		return nil
	}
	var requestPos bool
	if pos > 0 {
		requestPos = true
		pos-- // because we cannot use zero as index 0 when calling writeTo somewhere
	}
	switch v := arg.value.(type) {
	case int:
		err = writeInt64(w, int64(v))
	case []int:
		if requestPos {
			err = writeInt64(w, int64(v[pos]))
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = writeInt64(w, int64(v[i]))
			}
			w.WriteByte(')')
		}
	case int64:
		err = writeInt64(w, v)
	case []int64:
		if requestPos {
			err = writeInt64(w, v[pos])
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = writeInt64(w, v[i])
			}
			w.WriteByte(')')
		}
	case NullInt64:
		err = v.writeTo(w)
	case []NullInt64:
		if requestPos {
			err = v[pos].writeTo(w)
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = v[i].writeTo(w)
			}
			w.WriteByte(')')
		}
	case uint64:
		err = writeUint64(w, v)
	case []uint64:
		if requestPos {
			err = writeUint64(w, v[pos])
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = writeUint64(w, v[i])
			}
			w.WriteByte(')')
		}
	case []uint:
		if requestPos {
			err = writeUint64(w, uint64(v[pos]))
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = writeUint64(w, uint64(v[i]))
			}
			w.WriteByte(')')
		}
	case float64:
		err = writeFloat64(w, v)
	case []float64:
		if requestPos {
			err = writeFloat64(w, v[pos])
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = writeFloat64(w, v[i])
			}
			w.WriteByte(')')
		}
	case NullFloat64:
		err = v.writeTo(w)
	case []NullFloat64:
		if requestPos {
			err = v[pos].writeTo(w)
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = v[i].writeTo(w)
			}
			w.WriteByte(')')
		}
	case bool:
		dialect.EscapeBool(w, v)
	case []bool:
		if requestPos {
			dialect.EscapeBool(w, v[pos])
		} else {
			w.WriteByte('(')
			for i, val := range v {
				if i > 0 {
					w.WriteByte(',')
				}
				dialect.EscapeBool(w, val)
			}
			w.WriteByte(')')
		}
	case NullBool:
		v.writeTo(w)
	case []NullBool:
		if requestPos {
			v[pos].writeTo(w)
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = v[i].writeTo(w)
			}
			w.WriteByte(')')
		}
	case string:
		if !utf8.ValidString(v) {
			return errors.NewNotValidf("[dml] Argument.WriteTo: String is not UTF-8: %q", v)
		}
		dialect.EscapeString(w, v)
	case []string:
		if requestPos {
			if nv := v[pos]; utf8.ValidString(nv) {
				dialect.EscapeString(w, nv)
			} else {
				err = errors.NewNotValidf("[dml] Argument.WriteTo: String is not UTF-8: %q", nv)
			}
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				if nv := v[i]; utf8.ValidString(nv) {
					dialect.EscapeString(w, nv)
				} else {
					err = errors.NewNotValidf("[dml] Argument.WriteTo: String is not UTF-8: %q", nv)
				}
			}
			w.WriteByte(')')
		}
	case NullString:
		err = v.writeTo(w)
	case []NullString:
		if requestPos {
			err = v[pos].writeTo(w)
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = v[i].writeTo(w)
			}
			w.WriteByte(')')
		}
	case []byte:
		err = writeBytes(w, v)

	case [][]byte:
		if requestPos {
			err = writeBytes(w, v[pos])
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = writeBytes(w, v[i])
			}
			w.WriteByte(')')
		}
	case time.Time:
		dialect.EscapeTime(w, v)
	case []time.Time:
		if requestPos {
			dialect.EscapeTime(w, v[pos])
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					err = w.WriteByte(',')
				}
				dialect.EscapeTime(w, v[i])
			}
			w.WriteByte(')')
		}
	case NullTime:
		err = v.writeTo(w)
	case []NullTime:
		if requestPos {
			err = v[pos].writeTo(w)
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = v[i].writeTo(w)
			}
			w.WriteByte(')')
		}
	case nil:
		_, err = w.WriteString(sqlStrNullUC)

	default:
		panic(errors.NewNotSupportedf("[dml] Unsupported field type: %T => %#v", arg.value, arg.value))
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
	case NullInt64:
		buf.WriteString(".NullInt64(")
		buf.WriteString(v.GoString())
		buf.WriteByte(')')
	case []NullInt64:
		buf.WriteString(".NullInt64s(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(nv.GoString())
		}
		buf.WriteByte(')')

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
	case NullFloat64:
		buf.WriteString(".NullFloat64(")
		buf.WriteString(v.GoString())
		buf.WriteByte(')')
	case []NullFloat64:
		buf.WriteString(".NullFloat64s(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(nv.GoString())
		}
		buf.WriteByte(')')

	case bool:
		fmt.Fprintf(buf, ".Bool(%v)", v)
	case []bool:
		fmt.Fprintf(buf, ".Bools(%#v...)", v)
	case NullBool:
		buf.WriteString(".NullBool(")
		buf.WriteString(v.GoString())
		buf.WriteByte(')')
	case []NullBool:
		buf.WriteString(".NullBools(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(nv.GoString())
		}
		buf.WriteByte(')')

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
		buf.WriteByte(')')
	case NullString:
		buf.WriteString(".NullString(")
		buf.WriteString(v.GoString())
		buf.WriteByte(')')
	case []NullString:
		buf.WriteString(".NullStrings(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(nv.GoString())
		}
		buf.WriteByte(')')

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
		buf.WriteByte(')')

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
		buf.WriteByte(')')
	case NullTime:
		buf.WriteString(".NullTime(")
		buf.WriteString(v.GoString())
		buf.WriteByte(')')
	case []NullTime:
		buf.WriteString(".NullTimes(")
		for i, nv := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(nv.GoString())
		}
		buf.WriteByte(')')

	case nil:
		fmt.Fprint(buf, ".Null()")
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

// unnamedArgByPos returns an unnamed argument by its position.
func (a Arguments) unnamedArgByPos(pos int) (argument, bool) {
	unnamedCounter := 0
	for _, arg := range a {
		if arg.name == "" {
			if unnamedCounter == pos {
				return arg, true
			}
			unnamedCounter++
		}
	}
	return argument{}, false
}

func (a Arguments) hasNamedArgs() bool {
	for _, arg := range a {
		if arg.name != "" {
			return true
		}
	}
	return false
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
	if len(a) > 1 {
		buf.WriteByte('(')
	}
	for j, arg := range a {
		if j > 0 {
			buf.WriteByte(',')
		}
		if err := arg.writeTo(buf, 0); err != nil {
			return errors.Wrapf(err, "[dml] args write failed at pos %d with argument %#v", j, arg)
		}
	}
	if len(a) > 1 {
		buf.WriteByte(')')
	}
	return nil
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
		case NullInt64:
			args = vv.append(args)
		case []NullInt64:
			for _, v := range vv {
				args = v.append(args)
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
		case NullFloat64:
			args = vv.append(args)
		case []NullFloat64:
			for _, v := range vv {
				args = v.append(args)
			}

		case []bool:
			for _, v := range vv {
				args = append(args, v)
			}
		case NullBool:
			args = vv.append(args)
		case []NullBool:
			for _, v := range vv {
				args = v.append(args)
			}

		case []string:
			for _, v := range vv {
				args = append(args, v)
			}
		case NullString:
			args = vv.append(args)
		case []NullString:
			for _, v := range vv {
				args = v.append(args)
			}

		case [][]byte:
			for _, v := range vv {
				args = append(args, v)
			}

		case []time.Time:
			for _, v := range vv {
				args = append(args, v)
			}
		case NullTime:
			args = vv.append(args)
		case []NullTime:
			for _, v := range vv {
				args = v.append(args)
			}
		default:
			panic(errors.NewNotSupportedf("[dml] Unsupported field type: %T", arg.value))
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

func (a Arguments) Null() Arguments                          { return a.add(nil) }
func (a Arguments) Unsafe(arg interface{}) Arguments         { return a.add(arg) }
func (a Arguments) Int(i int) Arguments                      { return a.add(i) }
func (a Arguments) Ints(i ...int) Arguments                  { return a.add(i) }
func (a Arguments) Int64(i int64) Arguments                  { return a.add(i) }
func (a Arguments) Int64s(i ...int64) Arguments              { return a.add(i) }
func (a Arguments) Uint(i uint) Arguments                    { return a.add(uint64(i)) }
func (a Arguments) Uints(i ...uint) Arguments                { return a.add(i) }
func (a Arguments) Uint64(i uint64) Arguments                { return a.add(i) }
func (a Arguments) Uint64s(i ...uint64) Arguments            { return a.add(i) }
func (a Arguments) Float64(f float64) Arguments              { return a.add(f) }
func (a Arguments) Float64s(f ...float64) Arguments          { return a.add(f) }
func (a Arguments) Bool(b bool) Arguments                    { return a.add(b) }
func (a Arguments) Bools(b ...bool) Arguments                { return a.add(b) }
func (a Arguments) String(s string) Arguments                { return a.add(s) }
func (a Arguments) Strings(s ...string) Arguments            { return a.add(s) }
func (a Arguments) Time(t time.Time) Arguments               { return a.add(t) }
func (a Arguments) Times(t ...time.Time) Arguments           { return a.add(t) }
func (a Arguments) Bytes(b []byte) Arguments                 { return a.add(b) }
func (a Arguments) BytesSlice(b ...[]byte) Arguments         { return a.add(b) }
func (a Arguments) NullString(nv NullString) Arguments       { return a.add(nv) }
func (a Arguments) NullStrings(nv ...NullString) Arguments   { return a.add(nv) }
func (a Arguments) NullFloat64(nv NullFloat64) Arguments     { return a.add(nv) }
func (a Arguments) NullFloat64s(nv ...NullFloat64) Arguments { return a.add(nv) }
func (a Arguments) NullInt64(nv NullInt64) Arguments         { return a.add(nv) }
func (a Arguments) NullInt64s(nv ...NullInt64) Arguments     { return a.add(nv) }
func (a Arguments) NullBool(nv NullBool) Arguments           { return a.add(nv) }
func (a Arguments) NullBools(nv ...NullBool) Arguments       { return a.add(nv) }
func (a Arguments) NullTime(nv NullTime) Arguments           { return a.add(nv) }
func (a Arguments) NullTimes(nv ...NullTime) Arguments       { return a.add(nv) }

// Name sets the name for the following argument. Calling Name two times after
// each other sets the first call to Name to a NULL value. A call to Name should
// always follow a call to a function type like Int, Float64s or NullTime.
// Name may contain the placeholder prefix colon.
func (a Arguments) Name(n string) Arguments { return append(a, argument{name: n}) }

// TODO: maybe use such a function to set the position, but then add a new field: pos int to the argument struct
// func (a Arguments) Pos(n int) Arguments { return append(a, argument{name: n}) }

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

// DriverValues adds each driver.Value as its own argument to the argument
// slice. It panics if the underlying type is not one of the allowed of
// interface driver.Valuer.
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
			a = a.Time(t)
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
