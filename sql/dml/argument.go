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
	"database/sql"
	"database/sql/driver"
	"math"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
)

// https://www.adampalmer.me/iodigitalsec/2013/08/18/mysql_real_escape_string-wont-magically-solve-your-sql-injection-problems/

const (
	sqlStrNullUC = "NULL"
	sqlStar      = "*"
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

// internalNULLNIL represent an internal indicator that the value NULL should be
// written, if an interface{} is nil, then nothing gets written in function
// writeInterfaceValue.
type internalNULLNIL struct{}

func sliceLen(arg interface{}) (l int, isSlice bool) {
	switch v := arg.(type) {
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

func writeInterfaceValue(arg interface{}, w *bytes.Buffer, pos uint) (err error) {
	var requestPos bool
	if pos > 0 {
		requestPos = true
		pos-- // because we cannot use zero as index 0 when calling writeTo somewhere
	}
	switch v := arg.(type) {
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
	case internalNULLNIL:
		_, err = w.WriteString(sqlStrNullUC)
	case nil:
		// do nothing
		// _, err = w.WriteString("[PLEASE USE type internalNULLNIL]")
	case sql.NamedArg:
		return writeInterfaceValue(v.Value, w, pos)
	default:
		return errors.NotSupported.Newf("[dml] Unsupported field type: %T => %#v", arg, arg)
	}
	return err
}

// multiplyInterfaceValues is only applicable when using *Union as a template.
// multiplyInterfaceValues repeats the `args` variable n-times to match the number of
// generated SELECT queries in the final UNION statement. It should be called
// after all calls to `StringReplace` have been made.
func multiplyInterfaceValues(inArg []interface{}, factor int) []interface{} {
	if factor < 2 {
		return inArg
	}
	lArgs := len(inArg)
	newA := make([]interface{}, lArgs*factor)
	for i := 0; i < factor; i++ {
		copy(newA[i*lArgs:], inArg)
	}
	return newA
}

// Len returns the total length of all arguments.
func totalSliceLenSimple(args []interface{}) (l int) {
	for _, arg := range args {
		al, _ := sliceLen(arg)
		l += al
	}
	return l
}

func totalSliceLen(args []interface{}) (l int, containsAtLeastOneSlice bool) {
	for _, arg := range args {
		al, isSlice := sliceLen(arg)
		if isSlice {
			containsAtLeastOneSlice = true
		}
		l += al
	}
	return
}

// Write writes all arguments into buf and separates by a comma.
func writeInterfaces(buf *bytes.Buffer, args []interface{}) error {
	if len(args) > 1 {
		buf.WriteByte('(')
	}
	for j, arg := range args {
		if j > 0 {
			buf.WriteByte(',')
		}
		if err := writeInterfaceValue(arg, buf, 0); err != nil {
			return errors.Wrapf(err, "[dml] args write failed at pos %d with argument %#v", j, arg)
		}
	}
	if len(args) > 1 {
		buf.WriteByte(')')
	}
	return nil
}

// expandInterfaces creates an interface slice with flatten values. E.g. a
// string slice gets expanded into its strings. Each type is one of the allowed
// types in driver.Value. It appends its values to the `args` slice.
func expandInterfaces(args []interface{}) []interface{} {
	lenArgs := len(args)
	if lenArgs == 0 {
		return nil
	}
	appendTo := make([]interface{}, 0, 3*lenArgs)
	for _, arg := range args {
		appendTo = expandInterface(appendTo, arg)
	}
	return appendTo
}

func expandInterface(appendTo []interface{}, arg interface{}) []interface{} {
	switch vv := arg.(type) {

	case bool, string, []byte, time.Time, float64, int64, nil:
		appendTo = append(appendTo, arg)

	case int:
		appendTo = append(appendTo, int64(vv))
	case []int:
		for _, v := range vv {
			appendTo = append(appendTo, int64(v))
		}
	case int8:
		appendTo = append(appendTo, int64(vv))
	case []int8:
		for _, v := range vv {
			appendTo = append(appendTo, int64(v))
		}
	case int16:
		appendTo = append(appendTo, int64(vv))
	case []int16:
		for _, v := range vv {
			appendTo = append(appendTo, int64(v))
		}
	case int32:
		appendTo = append(appendTo, int64(vv))
	case []int32:
		for _, v := range vv {
			appendTo = append(appendTo, int64(v))
		}

	case []int64:
		for _, v := range vv {
			appendTo = append(appendTo, v)
		}
	case null.Int8:
		appendTo = vv.Append(appendTo)
	case []null.Int8:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}
	case null.Int16:
		appendTo = vv.Append(appendTo)
	case []null.Int16:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}
	case null.Int32:
		appendTo = vv.Append(appendTo)
	case []null.Int32:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}
	case null.Int64:
		appendTo = vv.Append(appendTo)
	case []null.Int64:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}

		// Get send as text in a byte slice. The MySQL/MariaDB Server type
		// casts it into a bigint. If you change this, a test will fail.
	case uint64:
		if vv > math.MaxInt64 {
			appendTo = append(appendTo, strconv.AppendUint([]byte{}, vv, 10))
		} else {
			appendTo = append(appendTo, int64(vv))
		}
	case null.Uint8:
		appendTo = vv.Append(appendTo)
	case []null.Uint8:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}
	case null.Uint16:
		appendTo = vv.Append(appendTo)
	case []null.Uint16:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}
	case null.Uint32:
		appendTo = vv.Append(appendTo) // TODO check all uints for overflow of MaxInt64
	case []null.Uint32:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}
	case null.Uint64:
		appendTo = vv.Append(appendTo)
	case []null.Uint64:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}

	case uint:
		if vv > math.MaxInt64 {
			appendTo = append(appendTo, strconv.AppendUint([]byte{}, uint64(vv), 10))
		} else {
			appendTo = append(appendTo, int64(vv))
		}

	case uint8:
		appendTo = append(appendTo, int64(vv))
	case uint16:
		appendTo = append(appendTo, int64(vv))
	case uint32:
		appendTo = append(appendTo, int64(vv))

	case []uint64:
		for _, v := range vv {
			if v > math.MaxInt64 {
				appendTo = append(appendTo, strconv.AppendUint([]byte{}, v, 10))
			} else {
				appendTo = append(appendTo, int64(v))
			}
		}
	case []uint:
		for _, v := range vv {
			if v > math.MaxInt64 {
				appendTo = append(appendTo, strconv.AppendUint([]byte{}, uint64(v), 10))
			} else {
				appendTo = append(appendTo, int64(v))
			}
		}

	case []float64:
		for _, v := range vv {
			appendTo = append(appendTo, v)
		}
	case null.Float64:
		appendTo = vv.Append(appendTo)
	case []null.Float64:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}

	case []bool:
		for _, v := range vv {
			appendTo = append(appendTo, v)
		}
	case null.Bool:
		appendTo = vv.Append(appendTo)
	case []null.Bool:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}

	case []string:
		for _, v := range vv {
			appendTo = append(appendTo, v)
		}
	case null.String:
		appendTo = vv.Append(appendTo)
	case []null.String:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}

	case [][]byte:
		for _, v := range vv {
			appendTo = append(appendTo, v)
		}

	case []time.Time:
		for _, v := range vv {
			appendTo = append(appendTo, v)
		}
	case null.Time:
		appendTo = vv.Append(appendTo)
	case []null.Time:
		for _, v := range vv {
			appendTo = v.Append(appendTo)
		}
	case sql.NamedArg:
		appendTo = expandInterface(appendTo, vv.Value)
	case []sql.NamedArg:
		for _, v := range vv {
			appendTo = expandInterface(appendTo, v.Value)
		}
	case driver.Valuer:
		dvv, _ := vv.Value()
		appendTo = expandInterface(appendTo, dvv)
	case internalNULLNIL:
		appendTo = expandInterface(appendTo, nil)

	case QualifiedRecord, ColumnMapper:
		// skip and do nothing
	default:
		panic(errors.NotSupported.Newf("[dml] Unsupported field type: %T", arg))
	}
	return appendTo
}

func driverValue(appendTo []interface{}, dvs ...driver.Valuer) ([]interface{}, error) {
	// value is a value that drivers must be able to handle.
	// It is either nil or an instance of one of these types:
	//
	//   int64
	//   float64
	//   bool
	//   []byte
	//   string
	//   time.Time
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

	var arg interface{}
	switch {
	case len(i64s) > 0:
		arg = i64s
	case len(f64s) > 0:
		arg = f64s
	case len(bs) > 0:
		arg = bs
	case len(bytess) > 0:
		arg = bytess
	case len(strs) > 0:
		arg = strs
	case len(times) > 0:
		arg = times
	}
	return append(appendTo, arg), nil
}

func driverValues(appendToArgs []interface{}, dvs ...driver.Valuer) ([]interface{}, error) {
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
			appendToArgs = append(appendToArgs, nil) // TODO consider internal NIL type
			continue
		}
		v, err := dv.Value()
		if err != nil {
			return nil, errors.Fatal.New(err, "[dml] Driver.Values error for %#v", dv)
		}
		var a interface{}
		switch t := v.(type) {
		case nil:
			// nothing to do or TODO consider internal nil type
		case int64, float64, bool, []byte, string, time.Time:
			a = t
		default:
			return nil, errors.NotSupported.Newf("[dml] Type %#v not supported in Driver.Values slice: %#v", t, dvs)
		}
		appendToArgs = append(appendToArgs, a)
	}
	return appendToArgs, nil
}

func iFaceToArgs(args []interface{}, values ...interface{}) ([]interface{}, error) {
	for _, val := range values {
		switch v := val.(type) {
		case float32:
			args = append(args, float64(v))
		case float64:
			args = append(args, v)
		case int64:
			args = append(args, v)
		case int:
			args = append(args, int64(v))
		case int32:
			args = append(args, int64(v))
		case int16:
			args = append(args, int64(v))
		case int8:
			args = append(args, int64(v))
		case uint64:
			args = append(args, int64(v))
		case uint32:
			args = append(args, int64(v))
		case uint16:
			args = append(args, int64(v))
		case uint8:
			args = append(args, int64(v))
		case bool:
			args = append(args, v)
		case string:
			args = append(args, v)
		case []byte:
			args = append(args, v)
		case time.Time:
			args = append(args, v)
		case *time.Time:
			if v != nil {
				args = append(args, *v)
			}
		case nil:
			args = append(args, nil)
		default:
			return nil, errors.NotSupported.Newf("[dml] iFaceToArgs type %#v not yet supported", v)
		}
	}
	return args, nil
}
