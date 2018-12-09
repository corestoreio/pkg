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
	"math"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
)

// https://www.adampalmer.me/iodigitalsec/2013/08/18/mysql_real_escape_string-wont-magically-solve-your-sql-injection-problems/

const (
	sqlStrNullUC             = "NULL"
	sqlStar                  = "*"
	defaultArgumentsCapacity = 5
)

// QualifiedRecord is a ColumnMapper with a qualifier. A QualifiedRecord gets
// used as arguments to ExecRecord or WithRecords in the SQL statement. If you
// use an alias for the main table/view you must set the alias as the qualifier.
type QualifiedRecord struct {
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

func (arg *argument) sliceLen() (l int, isSlice bool) {
	switch v := arg.value.(type) {
	case nil, int, int64, uint64, float64, bool, string, []byte, time.Time, null.String, null.Int64, null.Float64, null.Bool, null.Time:
		l = 1
	case []int:
		l = len(v)
		isSlice = true
	case []int64:
		l = len(v)
		isSlice = true
	case []uint64:
		l = len(v)
		isSlice = true
	case []uint:
		l = len(v)
		isSlice = true
	case []float64:
		l = len(v)
		isSlice = true
	case []bool:
		l = len(v)
		isSlice = true
	case []string:
		l = len(v)
		isSlice = true
	case [][]byte:
		l = len(v)
		isSlice = true
	case []time.Time:
		l = len(v)
		isSlice = true
	case []null.String:
		l = len(v)
		isSlice = true
	case []null.Int64:
		l = len(v)
		isSlice = true
	case []null.Float64:
		l = len(v)
		isSlice = true
	case []null.Bool:
		l = len(v)
		isSlice = true
	case []null.Time:
		l = len(v)
		isSlice = true
	default:
		panic(errors.NotSupported.Newf("[dml] Unsupported type: %T => %#v", v, v))
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
	case null.Int64:
		err = v.WriteTo(dialect, w)
	case []null.Int64:
		if requestPos {
			err = v[pos].WriteTo(dialect, w)
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = v[i].WriteTo(dialect, w)
			}
			w.WriteByte(')')
		}
	case uint64:
		err = writeUint64(w, v)
	case uint:
		err = writeUint64(w, uint64(v))
	case uint8:
		err = writeUint64(w, uint64(v))
	case uint16:
		err = writeUint64(w, uint64(v))
	case uint32:
		err = writeUint64(w, uint64(v))
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
	case null.Float64:
		err = v.WriteTo(dialect, w)
	case []null.Float64:
		if requestPos {
			err = v[pos].WriteTo(dialect, w)
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = v[i].WriteTo(dialect, w)
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
	case null.Bool:
		v.WriteTo(dialect, w)
	case []null.Bool:
		if requestPos {
			v[pos].WriteTo(dialect, w)
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = v[i].WriteTo(dialect, w)
			}
			w.WriteByte(')')
		}
	case string:
		if !utf8.ValidString(v) {
			return errors.NotValid.Newf("[dml] Argument.WriteTo: String is not UTF-8: %q", v)
		}
		dialect.EscapeString(w, v)
	case []string:
		if requestPos {
			if nv := v[pos]; utf8.ValidString(nv) {
				dialect.EscapeString(w, nv)
			} else {
				err = errors.NotValid.Newf("[dml] Argument.WriteTo: String is not UTF-8: %q", nv)
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
					err = errors.NotValid.Newf("[dml] Argument.WriteTo: String is not UTF-8: %q", nv)
				}
			}
			w.WriteByte(')')
		}
	case null.String:
		err = v.WriteTo(dialect, w)
	case []null.String:
		if requestPos {
			err = v[pos].WriteTo(dialect, w)
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = v[i].WriteTo(dialect, w)
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
	case null.Time:
		err = v.WriteTo(dialect, w)
	case []null.Time:
		if requestPos {
			err = v[pos].WriteTo(dialect, w)
		} else {
			w.WriteByte('(')
			for l, i := len(v), 0; i < l && err == nil; i++ {
				if i > 0 {
					w.WriteByte(',')
				}
				err = v[i].WriteTo(dialect, w)
			}
			w.WriteByte(')')
		}
	case nil:
		_, err = w.WriteString(sqlStrNullUC)

	default:
		panic(errors.NotSupported.Newf("[dml] Unsupported field type: %T => %#v", arg.value, arg.value))
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
	case null.Int64:
		buf.WriteString(".NullInt64(")
		buf.WriteString(v.GoString())
		buf.WriteByte(')')
	case []null.Int64:
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
	case null.Float64:
		buf.WriteString(".NullFloat64(")
		buf.WriteString(v.GoString())
		buf.WriteByte(')')
	case []null.Float64:
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
	case null.Bool:
		buf.WriteString(".NullBool(")
		buf.WriteString(v.GoString())
		buf.WriteByte(')')
	case []null.Bool:
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
	case null.String:
		buf.WriteString(".NullString(")
		buf.WriteString(v.GoString())
		buf.WriteByte(')')
	case []null.String:
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
	case null.Time:
		buf.WriteString(".NullTime(")
		buf.WriteString(v.GoString())
		buf.WriteByte(')')
	case []null.Time:
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
		panic(errors.NotSupported.Newf("[dml] Unsupported field type: %T", arg.value))
	}
	return buf.String()
}

// MakeArgs creates a new argument slice with the desired capacity.
func MakeArgs(cap int) *Artisan { return &Artisan{arguments: make(arguments, 0, cap)} }

type arguments []argument

func (as arguments) Clone() arguments {
	if as == nil {
		return nil
	}
	c := make(arguments, len(as))
	copy(c, as)
	return c
}

// multiplyArguments is only applicable when using *Union as a template.
// multiplyArguments repeats the `args` variable n-times to match the number of
// generated SELECT queries in the final UNION statement. It should be called
// after all calls to `StringReplace` have been made.
func (as arguments) multiplyArguments(factor int) arguments {
	if factor < 2 {
		return as
	}
	lArgs := len(as)
	newA := make(arguments, lArgs*factor)
	for i := 0; i < factor; i++ {
		copy(newA[i*lArgs:], as)
	}
	return newA
}

func (as arguments) GoString() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "dml.MakeArgs(%d)", len(as))
	for _, arg := range as {
		buf.WriteString(arg.GoString())
	}
	return buf.String()
}

// Len returns the total length of all arguments.
func (as arguments) Len() (l int) {
	for _, arg := range as {
		al, _ := arg.sliceLen()
		l += al
	}
	return l
}

func (as arguments) totalSliceLen() (l int, containsAtLeastOneSlice bool) {
	for _, arg := range as {
		al, isSlice := arg.sliceLen()
		if isSlice {
			containsAtLeastOneSlice = true
		}
		l += al
	}
	return
}

// Write writes all arguments into buf and separates by a comma.
func (as arguments) Write(buf *bytes.Buffer) error {
	if len(as) > 1 {
		buf.WriteByte('(')
	}
	for j, arg := range as {
		if j > 0 {
			buf.WriteByte(',')
		}
		if err := arg.writeTo(buf, 0); err != nil {
			return errors.Wrapf(err, "[dml] args write failed at pos %d with argument %#v", j, arg)
		}
	}
	if len(as) > 1 {
		buf.WriteByte(')')
	}
	return nil
}

// Interfaces creates an interface slice with flatend values. Each type is one
// of the allowed types in driver.Value. It appends its values to the `args`
// slice.
func (as arguments) Interfaces(args ...interface{}) []interface{} {
	if len(as) == 0 {
		return args
	}
	if args == nil {
		args = make([]interface{}, 0, 2*len(as))
	}

	for _, arg := range as {
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
		case null.Int64:
			args = vv.Append(args)
		case []null.Int64:
			for _, v := range vv {
				args = v.Append(args)
			}

			// Get send as text in a byte slice. The MySQL/MariaDB Server type
			// casts it into a bigint. If you change this, a test will fail.
		case uint64:
			if vv > math.MaxInt64 {
				args = append(args, strconv.AppendUint([]byte{}, vv, 10))
			} else {
				args = append(args, int64(vv))
			}
		case null.Uint64:
			args = vv.Append(args)
		case []null.Uint64:
			for _, v := range vv {
				args = v.Append(args)
			}

		case uint:
			if vv > math.MaxInt64 {
				args = append(args, strconv.AppendUint([]byte{}, uint64(vv), 10))
			} else {
				args = append(args, int64(vv))
			}

		case uint8:
			args = append(args, int64(vv))
		case uint16:
			args = append(args, int64(vv))
		case uint32:
			args = append(args, int64(vv))

		case []uint64:
			for _, v := range vv {
				if v > math.MaxInt64 {
					args = append(args, strconv.AppendUint([]byte{}, v, 10))
				} else {
					args = append(args, int64(v))
				}
			}
		case []uint:
			for _, v := range vv {
				if v > math.MaxInt64 {
					args = append(args, strconv.AppendUint([]byte{}, uint64(v), 10))
				} else {
					args = append(args, int64(v))
				}
			}

		case []float64:
			for _, v := range vv {
				args = append(args, v)
			}
		case null.Float64:
			args = vv.Append(args)
		case []null.Float64:
			for _, v := range vv {
				args = v.Append(args)
			}

		case []bool:
			for _, v := range vv {
				args = append(args, v)
			}
		case null.Bool:
			args = vv.Append(args)
		case []null.Bool:
			for _, v := range vv {
				args = v.Append(args)
			}

		case []string:
			for _, v := range vv {
				args = append(args, v)
			}
		case null.String:
			args = vv.Append(args)
		case []null.String:
			for _, v := range vv {
				args = v.Append(args)
			}

		case [][]byte:
			for _, v := range vv {
				args = append(args, v)
			}

		case []time.Time:
			for _, v := range vv {
				args = append(args, v)
			}
		case null.Time:
			args = vv.Append(args)
		case []null.Time:
			for _, v := range vv {
				args = v.Append(args)
			}
		default:
			panic(errors.NotSupported.Newf("[dml] Unsupported field type: %T", arg.value))
		}
	}
	return args
}

func (as arguments) add(v interface{}) arguments {
	if l := len(as); l > 0 {
		// look back if there might be a name.
		if arg := as[l-1]; !arg.isSet {
			// The previous call Name() has set the name and now we set the
			// value, but don't append a new entry.
			arg.isSet = true
			arg.value = v
			as[l-1] = arg
			return as
		}
	}
	return append(as, argument{isSet: true, value: v})
}

func driverValue(appendTo arguments, dvs ...driver.Valuer) (arguments, error) {
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
			return nil, errors.Fatal.New(err, "[dml] Driver.value error for %#v", dv)
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
			return nil, errors.NotSupported.Newf("[dml] Type %#v not supported in value slice: %#v", t, dvs)
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
	return append(appendTo, arg), nil
}

func driverValues(appendToArgs arguments, dvs ...driver.Valuer) (arguments, error) {
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
			var a argument
			a.set(nil)
			appendToArgs = append(appendToArgs, a)
			continue
		}
		v, err := dv.Value()
		if err != nil {
			return nil, errors.Fatal.New(err, "[dml] Driver.Values error for %#v", dv)
		}
		var a argument
		switch t := v.(type) {
		case nil:
			a.set(nil)
		case int64, float64, bool, []byte, string, time.Time:
			a.set(t)
		default:
			return nil, errors.NotSupported.Newf("[dml] Type %#v not supported in Driver.Values slice: %#v", t, dvs)
		}
		appendToArgs = append(appendToArgs, a)
	}
	return appendToArgs, nil
}

func iFaceToArgs(args arguments, values ...interface{}) (arguments, error) {
	for _, val := range values {
		switch v := val.(type) {
		case float32:
			args = args.add(float64(v))
		case float64:
			args = args.add(v)
		case int64:
			args = args.add(v)
		case int:
			args = args.add(int64(v))
		case int32:
			args = args.add(int64(v))
		case int16:
			args = args.add(int64(v))
		case int8:
			args = args.add(int64(v))
		case uint32:
			args = args.add(int64(v))
		case uint16:
			args = args.add(int64(v))
		case uint8:
			args = args.add(int64(v))
		case bool:
			args = args.add(v)
		case string:
			args = args.add(v)
		case []byte:
			args = args.add(v)
		case time.Time:
			args = args.add(v)
		case *time.Time:
			if v != nil {
				args = args.add(*v)
			}
		case nil:
			args = args.add(nil)
		default:
			return nil, errors.NotSupported.Newf("[dml] iFaceToArgs type %#v not yet supported", v)
		}
	}
	return args, nil
}
