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
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// https://www.adampalmer.me/iodigitalsec/2013/08/18/mysql_real_escape_string-wont-magically-solve-your-sql-injection-problems/

const (
	sqlStrNullUC             = "NULL"
	sqlStrNullLC             = "null"
	sqlStar                  = "*"
	defaultArgumentsCapacity = 5
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
		panic(errors.NotSupported.Newf("[dml] Unsupported field type: %T", arg.value))
	}
	return buf.String()
}

const (
	argOptionExpandPlaceholder = 1 << iota
	argOptionInterpolate
)

// Arguments a collection of primitive types or slices of primitive types.
// It acts as some kind of prepared statement.
type Arguments struct {
	base builderCommon
	// insertCachedSQL contains the final build SQL string with the correct
	// amount of placeholders.
	insertCachedSQL   []byte
	insertColumnCount uint
	insertRowCount    uint
	Options           uint
	// hasNamedArgs checks if the SQL string in the cachedSQL field contains
	// named arguments. 0 not yet checked, 1=does not contain, 2 = yes
	hasNamedArgs      uint8 // 0 not checked, 1=no, 2=yes
	nextUnnamedArgPos int
	raw               []interface{}
	args              []argument
	recs              []QualifiedRecord
	// rawReturn will be filled with the final primitives which gets returned
	// in the ToSQL function. It gets reset in every call.
	rawReturn    []interface{}
	argsPrepared bool
}

// ToSQL the returned interface slice is owned by the callee.
func (a *Arguments) ToSQL() (string, []interface{}, error) {
	return a.prepareArgs()
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `args` will then be nil.
func (a *Arguments) Interpolate() *Arguments {
	a.Options = a.Options | argOptionInterpolate
	return a
}

// ExpandPlaceHolders repeats the place holders with the provided argument
// count. If the amount of arguments does not match the number of place holders,
// a mismatch error gets returned.
//		ExpandPlaceHolders("SELECT * FROM table WHERE id IN (?) AND status IN (?)", Int(myIntSlice...), String(myStrSlice...))
// Gets converted to:
//		SELECT * FROM table WHERE id IN (?,?) AND status IN (?,?,?)
// The place holders are of course depending on the values in the Arg*
// functions. This function should be generally used when dealing with prepared
// statements or interpolation.
func (a *Arguments) ExpandPlaceHolders() *Arguments {
	a.Options = a.Options | argOptionExpandPlaceholder
	return a
}

// MakeArgs creates a new argument slice with the desired capacity.
func MakeArgs(cap int) *Arguments { return &Arguments{args: make([]argument, 0, cap)} }

func (a *Arguments) isEmpty() bool {
	if a == nil {
		return true
	}
	return len(a.raw) == 0 && len(a.args) == 0 && len(a.recs) == 0
}

func (a *Arguments) argsCount() int {
	if a == nil {
		return 0
	}
	return len(a.args)
}

// multiplyArguments is only applicable when using *Union as a template.
// multiplyArguments repeats the `args` variable n-times to match the number of
// generated SELECT queries in the final UNION statement. It should be called
// after all calls to `StringReplace` have been made.
func (a *Arguments) multiplyArguments() {
	factor := a.base.templateStmtCount
	if a == nil || factor < 2 {
		return
	}
	newA := make([]argument, len(a.args)*factor)
	lArgs := len(a.args)
	for i := 0; i < factor; i++ {
		copy(newA[i*lArgs:], a.args)
	}
	a.args = newA
	return
}

// prepareArgs transforms mainly the Arguments into []interface{} but also
// appends the `args` from the Exec+ or Query+ function. All method receivers
// are not thread safe. The returned interface slce belongs to the callee.
func (a *Arguments) prepareArgs(fncArgs ...interface{}) (string, []interface{}, error) {
	if len(a.base.cachedSQL) == 0 {
		return "", nil, errors.Empty.Newf("[dml] Arguments: The SQL string is empty.")
	}

	if a.base.source == dmlSourceInsert {
		return a.prepareArgsInsert(fncArgs...)
	}

	if a.isEmpty() {
		a.hasNamedArgs = 1
		return string(a.base.cachedSQL), fncArgs, nil
	}

	if a.hasNamedArgs == 0 {
		found := false
		a.base.cachedSQL, a.base.qualifiedColumns, found = extractReplaceNamedArgs(a.base.cachedSQL, a.base.qualifiedColumns)
		a.hasNamedArgs = 1
		if found {
			a.hasNamedArgs = 2
		}
	}
	a.raw = append(a.raw, fncArgs...)
	if err := a.appendConvertedRecordsToArguments(); err != nil {
		return "", nil, errors.WithStack(err)
	}

	// Make a copy of the original SQL statement because it gets modified in the
	// worst case. Best case would be no modification and hence we don't need a
	// bytes.Buffer from the pool! TODO(CYS) optimize this and only acquire a
	// buffer from the pool in the worse case.
	sqlBuf := bufferpool.Get()
	defer bufferpool.Put(sqlBuf)
	if _, err := sqlBuf.Write(a.base.cachedSQL); err != nil {
		return "", nil, errors.WithStack(err)
	}

	// `switch` statement no suitable.
	if a.Options > 0 && len(a.raw) > 0 && len(a.recs) == 0 && len(a.args) == 0 {
		return "", nil, errors.NotAllowed.Newf("[dml] Interpolation/ExpandPlaceholders supports only Records and Arguments and not yet an interface slice.")
	}
	var tmpBuf *bytes.Buffer
	if a.Options > 0 {
		tmpBuf = bufferpool.Get()
		defer bufferpool.Put(tmpBuf)
	}
	if a.Options&argOptionExpandPlaceholder != 0 {
		if phCount := bytes.Count(sqlBuf.Bytes(), placeHolderByte); phCount < a.Len() {
			tmpBuf.Grow(sqlBuf.Len() * 5 / 4)
			tmpBuf.Reset()
			if err := expandPlaceHolders(tmpBuf, sqlBuf.Bytes(), a); err != nil {
				return "", nil, errors.WithStack(err)
			}
			sqlBuf.Reset()
			if _, err := tmpBuf.WriteTo(sqlBuf); err != nil {
				return "", nil, errors.WithStack(err)
			}
			tmpBuf.Reset()
		}
	}
	if a.Options&argOptionInterpolate != 0 {
		if err := writeInterpolateBytes(tmpBuf, sqlBuf.Bytes(), a); err != nil {
			return "", nil, errors.Wrapf(err, "[dml] Interpolation failed: %q", sqlBuf.String())
		}
		//a.Reset() TODO why is this needed?
		return tmpBuf.String(), nil, nil
	}

	a.rawReturn = a.rawReturn[:0]
	a.rawReturn = a.Interfaces(a.raw...) // TODO investigate ret
	a.argsPrepared = true
	return sqlBuf.String(), a.rawReturn, nil
}

func (a *Arguments) prepareArgsInsert(fncArgs ...interface{}) (string, []interface{}, error) {
	if a.argsPrepared {
		if a.Options > 0 {
			if len(a.raw) > 0 && len(a.recs) == 0 && len(a.args) == 0 {
				return "", nil, errors.NotAllowed.Newf("[dml] Interpolation/ExpandPlaceholders supports only Records and Arguments and not yet an interface slice.")
			}

			if a.Options&argOptionInterpolate != 0 {
				sqlBuf := bufferpool.Get()
				defer bufferpool.Put(sqlBuf)

				if err := writeInterpolateBytes(sqlBuf, a.insertCachedSQL, a); err != nil {
					return "", nil, errors.Wrapf(err, "[dml] Interpolation failed: %q", sqlBuf.String())
				}
				return sqlBuf.String(), nil, nil
			}
		}

		a.raw = append(a.raw, fncArgs...)
		a.rawReturn = a.rawReturn[:0]
		a.rawReturn = a.Interfaces(a.raw...)
		return string(a.insertCachedSQL), a.rawReturn, nil
	}

	a.argsPrepared = true

	sqlBuf := bufferpool.GetTwin()
	defer bufferpool.PutTwin(sqlBuf)

	if _, err := sqlBuf.First.Write(a.base.cachedSQL); err != nil {
		return "", nil, errors.WithStack(err)
	}

	{ // Extract arguments from ColumnMapper and append them to `a.args`.
		cm := newColumnMap(MakeArgs(len(a.recs)*5/4), "") // TODO get it from sync.Pool
		// inserting multiple rows retrieved from a collection. There is no qualifier.
		cm.setColumns(a.base.qualifiedColumns)
		for _, qRec := range a.recs {
			if qRec.Qualifier != "" {
				return "", nil, errors.Fatal.Newf("[dml] Qualifier in %T is not supported and not needed.", qRec)
			}
			if err := qRec.Record.MapColumns(cm); err != nil {
				return "", nil, errors.WithStack(err)
			}
		}
		if len(cm.Args.args) > 0 {
			a.args = cm.Args.args // copy fom pool into a.args
		}
	}

	totalArgLen := uint(len(a.args) + len(a.raw))
	{ // Write placeholder list e.g. "VALUES (?,?),(?,?)"
		odkPos := bytes.Index(a.base.cachedSQL, []byte(onDuplicateKeyPart))
		if odkPos > 0 {
			sqlBuf.First.Reset()
			sqlBuf.First.Write(a.base.cachedSQL[:odkPos])
		}

		if a.insertRowCount > 0 {
			columnCount := totalArgLen / a.insertRowCount
			writeInsertPlaceholders(sqlBuf.First, a.insertRowCount, columnCount)

		} else if a.insertColumnCount > 0 {
			rowCount := totalArgLen / a.insertColumnCount
			if rowCount == 0 {
				rowCount = 1
			}
			writeInsertPlaceholders(sqlBuf.First, rowCount, a.insertColumnCount)
		}
		if odkPos > 0 {
			sqlBuf.First.Write(a.base.cachedSQL[odkPos:])
		}
		a.insertCachedSQL = bufTrySizeByResliceOrNew(a.insertCachedSQL, sqlBuf.First.Len())
		copy(a.insertCachedSQL, sqlBuf.First.Bytes())

	}

	if a.Options > 0 {
		if len(a.raw) > 0 && len(a.recs) == 0 && len(a.args) == 0 {
			return "", nil, errors.NotAllowed.Newf("[dml] Interpolation/ExpandPlaceholders supports only Records and Arguments and not yet an interface slice.")
		}

		if a.Options&argOptionInterpolate != 0 {
			if err := writeInterpolateBytes(sqlBuf.Second, sqlBuf.First.Bytes(), a); err != nil {
				return "", nil, errors.Wrapf(err, "[dml] Interpolation failed: %q", sqlBuf.First.String())
			}
			//a.Reset() TODO why is this needed?
			return sqlBuf.Second.String(), nil, nil
		}
	}

	// TODO this interface creation process can be further optimized
	a.raw = append(a.raw, fncArgs...)
	a.rawReturn = a.rawReturn[:0]
	a.rawReturn = a.Interfaces(a.raw...)

	return sqlBuf.First.String(), a.rawReturn, nil
}

func (a *Arguments) appendConvertedRecordsToArguments() error {
	if a.base.templateStmtCount == 0 {
		a.base.templateStmtCount = 1
	}
	if len(a.args) == 0 && len(a.recs) == 0 {
		return nil
	}

	if len(a.args) > 0 && len(a.recs) == 0 && a.base.templateStmtCount > 1 && a.hasNamedArgs < 2 {
		a.multiplyArguments()
		return nil
	}

	if a.argsPrepared { // || a.hasNamedArgs == 1
		return nil
	}

	// TODO refactor prototype and make it performant and beautiful code

	cm := newColumnMap(MakeArgs(len(a.args)+len(a.recs)), "") // can use an arg pool Arguments
	cm.Args = new(Arguments)                                  // TODO for now a new pointer, later should be `a` or sync.Pool or just []argument

	for tsc := 0; tsc < a.base.templateStmtCount; tsc++ { // only in case of UNION statements in combination with a template SELECT, can be optimized later

		// `qualifiedColumns` contains the correct order as the place holders
		// appear in the SQL string.
		for _, identifier := range a.base.qualifiedColumns {
			// identifier can be either: column or qualifier.column or :column
			qualifier, column := splitColumn(identifier)
			// a.base.defaultQualifier is empty in case of INSERT statements

			column, isNamedArg := cutNamedArgStartStr(column) // removes the colon for named arguments
			cm.columns[0] = column                            // length is always one, as created in newColumnMap

			if isNamedArg {
				// if the colon : cannot be found then a simple place holder ? has been detected
				if err := a.MapColumns(cm); err != nil {
					return errors.WithStack(err)
				}
			} else {
				found := false
				for _, qRec := range a.recs {
					if qRec.Qualifier == "" {
						qRec.Qualifier = a.base.defaultQualifier
					}
					if qRec.Qualifier == qualifier {
						if err := qRec.Record.MapColumns(cm); err != nil {
							return errors.WithStack(err)
						}
						found = true
					}
				}
				if !found {
					if pArg, ok := a.nextUnnamedArg(); ok {
						cm.Args.args = append(cm.Args.args, pArg)
					}
				}
			}
		}
	}

	if len(cm.Args.args) > 0 {
		a.args = cm.Args.args
	}

	return nil
}

// nextUnnamedArg returns an unnamed argument by its position.
func (a *Arguments) nextUnnamedArg() (argument, bool) {
	var unnamedCounter int
	lenArg := len(a.args)
	for i := 0; i < lenArg && a.nextUnnamedArgPos >= 0; i++ {
		if arg := a.args[i]; arg.name == "" {
			if unnamedCounter == a.nextUnnamedArgPos {
				a.nextUnnamedArgPos++
				return arg, true
			}
			unnamedCounter++
		}
	}
	a.nextUnnamedArgPos = -1 // nothing found, so no need to further iterate through the []argument slice.
	return argument{}, false
}

// MapColumns allows to merge one argument slice with another depending on the
// matched columns. Each argument in the slice must be a named argument.
// Implements interface ColumnMapper.
func (a *Arguments) MapColumns(cm *ColumnMap) error {
	if cm.Mode() == ColumnMapEntityReadAll {
		cm.Args.args = append(cm.Args.args, a.args...)
		return cm.Err()
	}
	for cm.Next() {
		// now a bit slow ... but will be refactored later with constant time
		// access, but first benchmark it. This for loop can be the 3rd one in the
		// overall chain.
		c := cm.Column()
		for _, arg := range a.args {
			// Case sensitive comparison
			if c != "" && arg.name == c {
				cm.Args.args = append(cm.Args.args, arg)
				break
			}
		}
	}
	return cm.Err()
}

func (a *Arguments) GoString() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "dml.MakeArgs(%d)", len(a.args))
	for _, arg := range a.args {
		buf.WriteString(arg.GoString())
	}
	return buf.String()
}

func (a *Arguments) sliceArgumentLen() int {
	if a == nil {
		return 0
	}
	return len(a.args)
}

// Len returns the total length of all arguments.
func (a *Arguments) Len() int {
	if a == nil {
		return 0
	}
	var l int
	for _, arg := range a.args {
		l += arg.len()
	}
	return l
}

// Write writes all arguments into buf and separates by a comma.
func (a Arguments) Write(buf *bytes.Buffer) error {
	if len(a.args) > 1 {
		buf.WriteByte('(')
	}
	for j, arg := range a.args {
		if j > 0 {
			buf.WriteByte(',')
		}
		if err := arg.writeTo(buf, 0); err != nil {
			return errors.Wrapf(err, "[dml] args write failed at pos %d with argument %#v", j, arg)
		}
	}
	if len(a.args) > 1 {
		buf.WriteByte(')')
	}
	return nil
}

// Interfaces creates an interface slice with flatend values. Each type is one
// of the allowed types in driver.Value. It appends its values to the `args`
// slice.
func (a *Arguments) Interfaces(args ...interface{}) []interface{} {
	const maxInt64 = 1<<63 - 1
	if len(a.args) == 0 {
		return args
	}
	if args == nil {
		args = make([]interface{}, 0, 2*len(a.args))
	}

	for _, arg := range a.args {
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
			panic(errors.NotSupported.Newf("[dml] Unsupported field type: %T", arg.value))
		}
	}
	return args
}

func (a *Arguments) add(v interface{}) *Arguments {
	if a == nil {
		a = MakeArgs(defaultArgumentsCapacity)
	}
	if l := len(a.args); l > 0 {
		// look back if there might be a name.
		if arg := a.args[l-1]; !arg.isSet {
			// The previous call Name() has set the name and now we set the
			// value, but don't append a new entry.
			arg.isSet = true
			arg.value = v
			a.args[l-1] = arg
			return a
		}
	}
	a.args = append(a.args, argument{isSet: true, value: v})
	return a
}

// TODO QualifiedRecord can be removed because we can use Arguments.Name function to qualify a record.

func (a *Arguments) Record(qualifier string, record ColumnMapper) *Arguments {
	a.recs = append(a.recs, Qualify(qualifier, record))
	return a
}

// Arguments sets the internal arguments slice to the provided argument. Those
// are the slices Arguments, records and raw.
func (a *Arguments) Arguments(args *Arguments) *Arguments {
	// maybe deprecated this function.
	a.args = args.args
	a.recs = args.recs
	a.raw = args.raw
	return a
}

func (a *Arguments) Records(records ...QualifiedRecord) *Arguments { a.recs = records; return a }
func (a *Arguments) Raw(raw ...interface{}) *Arguments             { a.raw = raw; return a }

func (a *Arguments) Null() *Arguments                          { return a.add(nil) }
func (a *Arguments) Unsafe(arg interface{}) *Arguments         { return a.add(arg) }
func (a *Arguments) Int(i int) *Arguments                      { return a.add(i) }
func (a *Arguments) Ints(i ...int) *Arguments                  { return a.add(i) }
func (a *Arguments) Int64(i int64) *Arguments                  { return a.add(i) }
func (a *Arguments) Int64s(i ...int64) *Arguments              { return a.add(i) }
func (a *Arguments) Uint(i uint) *Arguments                    { return a.add(uint64(i)) }
func (a *Arguments) Uints(i ...uint) *Arguments                { return a.add(i) }
func (a *Arguments) Uint64(i uint64) *Arguments                { return a.add(i) }
func (a *Arguments) Uint64s(i ...uint64) *Arguments            { return a.add(i) }
func (a *Arguments) Float64(f float64) *Arguments              { return a.add(f) }
func (a *Arguments) Float64s(f ...float64) *Arguments          { return a.add(f) }
func (a *Arguments) Bool(b bool) *Arguments                    { return a.add(b) }
func (a *Arguments) Bools(b ...bool) *Arguments                { return a.add(b) }
func (a *Arguments) String(s string) *Arguments                { return a.add(s) }
func (a *Arguments) Strings(s ...string) *Arguments            { return a.add(s) }
func (a *Arguments) Time(t time.Time) *Arguments               { return a.add(t) }
func (a *Arguments) Times(t ...time.Time) *Arguments           { return a.add(t) }
func (a *Arguments) Bytes(b []byte) *Arguments                 { return a.add(b) }
func (a *Arguments) BytesSlice(b ...[]byte) *Arguments         { return a.add(b) }
func (a *Arguments) NullString(nv NullString) *Arguments       { return a.add(nv) }
func (a *Arguments) NullStrings(nv ...NullString) *Arguments   { return a.add(nv) }
func (a *Arguments) NullFloat64(nv NullFloat64) *Arguments     { return a.add(nv) }
func (a *Arguments) NullFloat64s(nv ...NullFloat64) *Arguments { return a.add(nv) }
func (a *Arguments) NullInt64(nv NullInt64) *Arguments         { return a.add(nv) }
func (a *Arguments) NullInt64s(nv ...NullInt64) *Arguments     { return a.add(nv) }
func (a *Arguments) NullBool(nv NullBool) *Arguments           { return a.add(nv) }
func (a *Arguments) NullBools(nv ...NullBool) *Arguments       { return a.add(nv) }
func (a *Arguments) NullTime(nv NullTime) *Arguments           { return a.add(nv) }
func (a *Arguments) NullTimes(nv ...NullTime) *Arguments       { return a.add(nv) }

// Name sets the name for the following argument. Calling Name two times after
// each other sets the first call to Name to a NULL value. A call to Name should
// always follow a call to a function type like Int, Float64s or NullTime.
// Name may contain the placeholder prefix colon.
func (a *Arguments) Name(n string) *Arguments { a.args = append(a.args, argument{name: n}); return a }

// TODO: maybe use such a function to set the position, but then add a new field: pos int to the argument struct
// func (a *Arguments) Pos(n int) *Arguments { return append(a, argument{name: n}) }

// Reset resets the slice for new usage retaining the already allocated memory.
func (a *Arguments) Reset() *Arguments {
	for i := range a.recs {
		a.recs[i].Qualifier = ""
		a.recs[i].Record = nil
	}
	a.recs = a.recs[:0]
	a.args = a.args[:0]
	a.raw = a.raw[:0]
	a.argsPrepared = false
	a.nextUnnamedArgPos = 0
	return a
}

// DriverValue adds multiple of the same underlying values to the argument
// slice. When using different values, the last applied value wins and gets
// added to the argument slice. For example driver.Values of type `int` will
// result in []int.
func (a *Arguments) DriverValue(dvs ...driver.Valuer) *Arguments {
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
			panic(errors.Fatal.New(err, "[dml] Driver.value error for %#v", dv))
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
			panic(errors.NotSupported.Newf("[dml] Type %#v not supported in value slice: %#v", t, dvs))
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
	if a == nil {
		a = MakeArgs(len(dvs))
	}
	a.args = append(a.args, arg)
	return a
}

// DriverValues adds each driver.Value as its own argument to the argument
// slice. It panics if the underlying type is not one of the allowed of
// interface driver.Valuer.
func (a *Arguments) DriverValues(dvs ...driver.Valuer) *Arguments {
	if a == nil {
		a = MakeArgs(len(dvs))
	}
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
			panic(errors.Fatal.New(err, "[dml] Driver.value error for %#v", dv))
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
			panic(errors.NotSupported.Newf("[dml] Type %#v not supported in value slice: %#v", t, dvs))
		}
	}
	return a
}

func iFaceToArgs(values ...interface{}) *Arguments {
	args := MakeArgs(len(values))
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
			panic(errors.NotSupported.Newf("[dml] iFaceToArgs type %#v not yet supported", v))
		}
	}
	return args
}

// WithDB sets the database query object.
func (a *Arguments) WithDB(db QueryExecPreparer) *Arguments {
	a.base.DB = db
	return a
}

/*********************************************
	LOAD / QUERY and EXEC functions
*********************************************/

var poolColumnMap = sync.Pool{
	New: func() interface{} {
		return new(ColumnMap)
	},
}

func poolColumnMapGet() *ColumnMap {
	return poolColumnMap.Get().(*ColumnMap)
}

func poolColumnMapPut(cm *ColumnMap) {
	cm.reset()
	poolColumnMap.Put(cm)
}

// Exec executes the statement represented by the Insert object. It returns the
// raw database/sql Result or an error if there was one. Regarding
// LastInsertID(): If you insert multiple rows using a single INSERT statement,
// LAST_INSERT_ID() returns the value generated for the first inserted row only.
// The reason for this at to make it possible to reproduce easily the same
// INSERT statement against some other server. If a record resp. and object
// implements the interface LastInsertIDAssigner then the LastInsertID gets
// assigned incrementally to the objects.
func (a *Arguments) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return a.exec(ctx, args...)
}

func (a *Arguments) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	return a.query(ctx, args...)
}

// QueryRow traditional way, allocation heavy.
func (a *Arguments) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	sqlStr, rawArgs, err := a.prepareArgs(args...)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("QueryRowContext", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	return a.base.DB.QueryRowContext(ctx, sqlStr, rawArgs...)
}

// Iterate iterates over the result set by loading only one row each iteration
// and then discarding it. Handles records one by one. Even if there are no rows
// in the query, the callBack function gets executed with a ColumnMap argument
// that indicates ColumnMap.HasRows equals false.
func (a *Arguments) Iterate(ctx context.Context, callBack func(*ColumnMap) error, args ...interface{}) (err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Iterate", log.String("id", a.base.id), log.Err(err))
	}

	r, err := a.query(ctx, args...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] Iterate.Query with query ID %q", a.base.id)
		return
	}
	cmr := poolColumnMapGet()
	defer func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] Iterate.QueryContext.Rows.Close")
		}
		poolColumnMapPut(cmr)
	}()

	for r.Next() {
		if err = cmr.Scan(r); err != nil {
			err = errors.WithStack(err)
			return
		}
		if err = callBack(cmr); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	err = errors.WithStack(r.Err())
	if !cmr.HasRows && err == nil {
		err = errors.WithStack(callBack(cmr))
	}
	return
}

// Load loads data from a query into an object. Load can load a single row or
// muliple-rows. It checks on top if ColumnMapper `s` implements io.Closer, to
// call the custom close function. This is useful for e.g. unlocking a mutex.
func (a *Arguments) Load(ctx context.Context, s ColumnMapper, args ...interface{}) (rowCount uint64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Load", log.String("id", a.base.id), log.Err(err), log.ObjectTypeOf("ColumnMapper", s))
	}

	r, err := a.query(ctx, args...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] Arguments.Load.QueryContext failed with queryID %q and ColumnMapper %T", a.base.id, s)
		return
	}
	cm := poolColumnMapGet()
	defer func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] Arguments.Load.Rows.Close")
		}
		if rc, ok := s.(ioCloser); ok {
			if err2 := rc.Close(); err2 != nil && err == nil {
				err = errors.Wrap(err2, "[dml] Arguments.Load.ColumnMapper.Close")
			}
		}
		poolColumnMapPut(cm)
	}()

	for r.Next() {
		if err = cm.Scan(r); err != nil {
			return 0, errors.WithStack(err)
		}
		if err = s.MapColumns(cm); err != nil {
			return 0, errors.Wrapf(err, "[dml] Arguments.Load failed with queryID %q and ColumnMapper %T", a.base.id, s)
		}
	}
	if err = r.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if cm.HasRows {
		cm.Count++ // because first row is zero but we want the actual row number
	}
	return cm.Count, err
}

// LoadInt64 executes the prepared statement and returns the value as an
// int64. It returns a NotFound error if the query returns nothing.
func (a *Arguments) LoadInt64(ctx context.Context, args ...interface{}) (int64, error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("LoadInt64")
	}
	return loadInt64(a.query(ctx, args...))
}

// LoadInt64s executes the Select and returns the value as a slice of
// int64s.
func (a *Arguments) LoadInt64s(ctx context.Context, args ...interface{}) (ret []int64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadInt64s", log.Int("row_count", len(ret)), log.Err(err))
	}
	ret, err = loadInt64s(a.query(ctx, args...))
	// Do not simplify it because we need ret in the defer. we don't log errors
	// because they get handled.
	return ret, err
}

// LoadUint64 executes the Select and returns the value at an uint64. It returns
// a NotFound error if the query returns nothing. This function comes in handy
// when performing a COUNT(*) query. See function `Select.Count`.
func (a *Arguments) LoadUint64(ctx context.Context, args ...interface{}) (_ uint64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadUint64", log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	var value uint64
	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return 0, errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if !found {
		err = errors.NotFound.Newf("[dml] LoadUint64 value not found")
	}
	return value, err
}

// LoadUint64s executes the Select and returns the value at a slice of uint64s.
func (a *Arguments) LoadUint64s(ctx context.Context, args ...interface{}) (values []uint64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadUint64s", log.Int("row_count", len(values)), log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	values = make([]uint64, 0, 10)
	for rows.Next() {
		var value uint64
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, nil
}

// LoadFloat64 executes the Select and returns the value at an float64. It
// returns a NotFound error if the query returns nothing.
func (a *Arguments) LoadFloat64(ctx context.Context, args ...interface{}) (_ float64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadFloat64", log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	var value float64
	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return 0, errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if !found {
		err = errors.NotFound.Newf("[dml] LoadFloat64 value not found")
	}
	return value, err
}

// LoadFloat64s executes the Select and returns the value at a slice of float64s.
func (a *Arguments) LoadFloat64s(ctx context.Context, args ...interface{}) (_ []float64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadFloat64s", log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	values := make([]float64, 0, 10)
	for rows.Next() {
		var value float64
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, err
}

// LoadString executes the Select and returns the value as a string. It
// returns a NotFound error if the row amount is not equal one.
func (a *Arguments) LoadString(ctx context.Context, args ...interface{}) (_ string, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadString", log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	var value string
	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return "", errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return "", errors.WithStack(err)
	}
	if !found {
		err = errors.NotFound.Newf("[dml] LoadInt64 value not found")
	}
	return value, err
}

// LoadStrings executes the Select and returns a slice of strings.
func (a *Arguments) LoadStrings(ctx context.Context, args ...interface{}) (values []string, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadStrings", log.Int("row_count", len(values)), log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	values = make([]string, 0, 10)
	for rows.Next() {
		var value string
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, err
}

func (a *Arguments) query(ctx context.Context, args ...interface{}) (rows *sql.Rows, err error) {
	var sqlStr string
	var rawArgs []interface{}
	sqlStr, rawArgs, err = a.prepareArgs(args...)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Query", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rows, err = a.base.DB.QueryContext(ctx, sqlStr, rawArgs...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] Query.QueryContext with query %q", sqlStr)
	}
	return
}

func loadInt64(rows *sql.Rows, errIn error) (value int64, err error) {
	if errIn != nil {
		return 0, errors.WithStack(errIn)
	}

	defer func() {
		if cErr := rows.Close(); err == nil && cErr != nil {
			err = errors.WithStack(cErr)
		}
	}()

	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return 0, errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if !found {
		err = errors.NotFound.Newf("[dml] LoadInt64 value not found")
	}
	return value, err
}

func loadInt64s(rows *sql.Rows, errIn error) (_ []int64, err error) {
	if errIn != nil {
		return nil, errors.WithStack(errIn)
	}
	defer func() {
		if cErr := rows.Close(); err == nil && cErr != nil {
			err = errors.WithStack(cErr)
		}
	}()

	values := make([]int64, 0, 16)
	for rows.Next() {
		var value int64
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, nil
}

func (a *Arguments) exec(ctx context.Context, args ...interface{}) (result sql.Result, err error) {
	var sqlStr string
	var rawArgs []interface{}
	sqlStr, rawArgs, err = a.prepareArgs(args...)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Exec", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result, err = a.base.DB.ExecContext(ctx, sqlStr, rawArgs...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] ExecContext with query %q", sqlStr) // err gets catched by the defer
		return
	}

	if a.recs == nil {
		return result, nil
	}
	lID, err := result.LastInsertId()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	for i, rec := range a.recs {
		if a, ok := rec.Record.(LastInsertIDAssigner); ok {
			a.AssignLastInsertID(lID + int64(i))
		}
	}
	return
}
