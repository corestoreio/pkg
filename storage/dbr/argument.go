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
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/bufferpool"
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
// SQL statement type. In case where stmtType has been isSet to
// SQLStmtInsert|SQLPartValues, the `columns` slice can be empty which means
// that all arguments are requested.
type ArgumentsAppender interface {
	AppendArguments(stmtType int, args ArgUnions, columns []string) (ArgUnions, error)
}

// argUnion is union type for different Go primitives and their slice
// representation. argUnion must be used as a pointer because it slows
// everything down. Check the benchmarks.
type argUnion struct {
	isSet bool
	// name for named place holders sql.NamedArg
	name  string // todo
	value interface{}
}

type placeHolder uint8

func (arg *argUnion) set(v interface{}) {
	arg.isSet = true
	arg.value = v
}

func (arg *argUnion) len() (l int) {
	switch v := arg.value.(type) {
	case nil, int, int64, uint64, float64, bool, string, []byte, time.Time, placeHolder:
		l = 1
	case []int:
		l = len(v)
	case []int64:
		l = len(v)
	case []uint64:
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

func (arg argUnion) writeTo(w *bytes.Buffer, pos int) (err error) {
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
		_, err = w.WriteString(sqlStrNull)

	case uint64:
		err = writeUint64(w, v)
	case []uint64:
		err = writeUint64(w, v[pos])

	case float64:
		err = writeFloat64(w, v)
	case []float64:
		err = writeFloat64(w, v[pos])
	case []NullFloat64:
		if s := v[pos]; s.Valid {
			return writeFloat64(w, s.Float64)
		}
		_, err = w.WriteString(sqlStrNull)

	case bool:
		dialect.EscapeBool(w, v)
	case []bool:
		dialect.EscapeBool(w, v[pos])
	case []NullBool:
		if s := v[pos]; s.Valid {
			dialect.EscapeBool(w, s.Bool)
			return nil
		}
		_, err = w.WriteString(sqlStrNull)

		// TODO(CyS) Cut the printed string in errors if it's longer than XX chars
	case string:
		if !utf8.ValidString(v) {
			return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", v)
		}
		dialect.EscapeString(w, v)
	case []string:
		if !utf8.ValidString(v[pos]) {
			return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", v[pos])
		}
		dialect.EscapeString(w, v[pos])
	case []NullString:
		if s := v[pos]; s.Valid {
			if !utf8.ValidString(s.String) {
				return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", s.String)
			}
			dialect.EscapeString(w, s.String)
		} else {
			_, err = w.WriteString(sqlStrNull)
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
			_, err = w.WriteString(sqlStrNull)
		}

	case nil:
		_, err = w.WriteString(sqlStrNull)
	case placeHolder:
		err = w.WriteByte(placeHolderRune)

	default:
		panic(errors.NewNotSupportedf("[dbr] Unsupported field type: %d", arg.value))
	}
	return err
}

func (arg argUnion) GoString() string {
	buf := new(bytes.Buffer)

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

	case float64:
		fmt.Fprintf(buf, ".Float64(%f)", v)
	case []float64:
		fmt.Fprintf(buf, ".Float64s(%#v...)", v)
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
		fmt.Fprintf(buf, ".Str(%q)", v)
	case []string:
		buf.WriteString(".Strs(")
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
		panic(errors.NewNotSupportedf("[dbr] Unsupported field type: %T", arg.value))
	}
	return buf.String()
}

// args a collection of primitive types or slice of primitive types. Using
// pointers in *argUnion would slow down the program.
type ArgUnions []argUnion

// MakeArgUnions creates a new argument union slice with the desired capacity.
func MakeArgUnions(cap int) ArgUnions {
	return make(ArgUnions, 0, cap)
}

func (a ArgUnions) GoString() string {
	buf := new(bytes.Buffer)
	buf.WriteString("dbr.MakeArgUnions()")
	for _, arg := range a {
		buf.WriteString(arg.GoString())
	}
	return buf.String()
}

// Len returns the total length of all arguments.
func (a ArgUnions) Len() int {
	var l int
	for _, arg := range a {
		l += arg.len()
	}
	return l
}

// String implements fmt.Stringer. Errors will be written in the returned
// string, which might be annoying for now. Can be changed later.
func (a ArgUnions) String() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := a.Write(buf); err != nil {
		return fmt.Sprintf("[dbr] args.String: %+v", err)
	}
	return buf.String()
}

// Write writes all arguments into buf and separated by a colon.
func (a ArgUnions) Write(buf *bytes.Buffer) error {
	buf.WriteByte('(')
	for j, arg := range a {
		l := arg.len()
		for i := 0; i < l; i++ {
			if i > 0 || j > 0 {
				buf.WriteByte(',')
			}
			if err := arg.writeTo(buf, i); err != nil {
				return errors.Wrapf(err, "[dbr] args write failed at pos %d with argument %#v", j, arg)
			}
		}
	}
	return buf.WriteByte(')')
}

// Interfaces creates an interface slice with flat values. Each type is one of
// the allowed in driver.value.
func (a ArgUnions) Interfaces(args ...interface{}) []interface{} {
	const maxInt64 = 1<<63 - 1
	if len(a) == 0 {
		return nil
	}
	if args == nil {
		args = make([]interface{}, 0, 2*len(a))
	}

	for j := 0; j < len(a); j++ { // faster than range
		arg := &a[j]
		switch vv := arg.value.(type) {

		case bool, string, []byte, time.Time, float64, int64, nil:
			args = append(args, vv) // vv is already interface{} !

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

func (a ArgUnions) PlaceHolder() ArgUnions {
	return append(a, argUnion{isSet: true, value: placeHolder(1)})
}
func (a ArgUnions) Null() ArgUnions         { return append(a, argUnion{isSet: true}) }
func (a ArgUnions) Int(i int) ArgUnions     { return append(a, argUnion{isSet: true, value: i}) }
func (a ArgUnions) Ints(i ...int) ArgUnions { return append(a, argUnion{isSet: true, value: i}) }
func (a ArgUnions) Int64(i int64) ArgUnions {
	return append(a, argUnion{isSet: true, value: int64(i)})
}
func (a ArgUnions) Int64s(i ...int64) ArgUnions {
	return append(a, argUnion{isSet: true, value: i})
}
func (a ArgUnions) Uint64(i uint64) ArgUnions {
	return append(a, argUnion{isSet: true, value: i})
}
func (a ArgUnions) Uint64s(i ...uint64) ArgUnions {
	return append(a, argUnion{isSet: true, value: i})
}
func (a ArgUnions) Float64(f float64) ArgUnions {
	return append(a, argUnion{isSet: true, value: f})
}
func (a ArgUnions) Float64s(f ...float64) ArgUnions {
	return append(a, argUnion{isSet: true, value: f})
}
func (a ArgUnions) Bool(f bool) ArgUnions {
	return append(a, argUnion{isSet: true, value: f})
}
func (a ArgUnions) Bools(f ...bool) ArgUnions {
	return append(a, argUnion{isSet: true, value: f})
}
func (a ArgUnions) Str(f string) ArgUnions {
	return append(a, argUnion{isSet: true, value: f})
}
func (a ArgUnions) Strs(f ...string) ArgUnions {
	return append(a, argUnion{isSet: true, value: f})
}
func (a ArgUnions) Bytes(b []byte) ArgUnions {
	return append(a, argUnion{isSet: true, value: b})
}
func (a ArgUnions) BytesSlice(b ...[]byte) ArgUnions {
	return append(a, argUnion{isSet: true, value: b})
}
func (a ArgUnions) Time(t time.Time) ArgUnions {
	return append(a, argUnion{isSet: true, value: t})
}
func (a ArgUnions) Times(t ...time.Time) ArgUnions {
	return append(a, argUnion{isSet: true, value: t})
}
func (a ArgUnions) NullString(nv ...NullString) ArgUnions {
	return append(a, argUnion{isSet: true, value: nv})
}
func (a ArgUnions) NullFloat64(nv ...NullFloat64) ArgUnions {
	return append(a, argUnion{isSet: true, value: nv})
}
func (a ArgUnions) NullInt64(nv ...NullInt64) ArgUnions {
	return append(a, argUnion{isSet: true, value: nv})
}
func (a ArgUnions) NullBool(nv ...NullBool) ArgUnions {
	return append(a, argUnion{isSet: true, value: nv})
}
func (a ArgUnions) NullTime(nv ...NullTime) ArgUnions {
	return append(a, argUnion{isSet: true, value: nv})
}

// DriverValue adds multiple of the same underlying values to the argument
// slice. When using different values, the last applied value wins and gets
// added to the argument slice.
func (a ArgUnions) DriverValue(dvs ...driver.Valuer) ArgUnions {
	// value is a value that drivers must be able to handle.
	// It is either nil or an instance of one of these types:
	//
	//   int64
	//   float64
	//   bool
	//   []byte
	//   string
	//   time.Time
	var arg argUnion
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
			panic(errors.NewFatal(err, "[dbr] Driver.value error for %#v", dv))
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
			panic(errors.NewNotSupportedf("[dbr] Type %#v not supported in value slice: %#v", t, dvs))
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
func (a ArgUnions) DriverValues(dvs ...driver.Valuer) ArgUnions {
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
			panic(errors.NewFatal(err, "[dbr] Driver.value error for %#v", dv))
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
			a = a.Str(t)
		case time.Time:
			a = a.Times(t)
		default:
			panic(errors.NewNotSupportedf("[dbr] Type %#v not supported in value slice: %#v", t, dvs))
		}
	}
	return a
}

func iFaceToArgs(values ...interface{}) ArgUnions {
	args := make(ArgUnions, 0, len(values))
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
			args = args.Str(v)
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
			panic(errors.NewNotSupportedf("[dbr] iFaceToArgs type %#v not yet supported", v))
		}
	}
	return args
}
