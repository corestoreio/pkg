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
	"sync"
	"time"
	"unicode/utf8"

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
		err = writeUint64(w, uint64(v))
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

type arguments []argument

var pooledArguments = sync.Pool{
	New: func() interface{} {
		var a [32]argument
		return arguments(a[:0])
	},
}

func pooledArgumentsGet() arguments {
	return pooledArguments.Get().(arguments)
}

func pooledArgumentsPut(a arguments, buf *bufferpool.TwinBuffer) {
	a = a[:0]
	pooledArguments.Put(a)
	if buf != nil {
		bufferpool.PutTwin(buf)
	}
}

const (
	argOptionExpandPlaceholder = 1 << iota
	argOptionInterpolate
)

// Arguments a collection of primitive types or slices of primitive types.
// It acts as some kind of prepared statement.
type Arguments struct {
	base builderCommon
	// QualifiedColumnsAliases allows to overwrite the internal qualified
	// columns slice with custom names. Especially in the use case when records
	// are applied. The list of column names in `QualifiedColumnsAliases` gets
	// passed to the ColumnMapper and back to the provided object. The
	// `QualifiedColumnsAliases` slice must have the same length as the
	// qualified columns slice. The order of the alias names must be in the same
	// order as the qualified columns or as the placeholders occur.
	QualifiedColumnsAliases []string
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
	arguments
	recs []QualifiedRecord
}

// WithQualifiedColumnsAliases for documentation please see:
// Arguments.QualifiedColumnsAliases.
func (a *Arguments) WithQualifiedColumnsAliases(aliases ...string) *Arguments {
	a.QualifiedColumnsAliases = aliases
	return a
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
func MakeArgs(cap int) *Arguments { return &Arguments{arguments: make(arguments, 0, cap)} }

func (a *Arguments) isEmpty() bool {
	if a == nil {
		return true
	}
	return len(a.raw) == 0 && len(a.arguments) == 0 && len(a.recs) == 0
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

// prepareArgs transforms mainly the Arguments into []interface{}. It appends
// its arguments to the `extArgs` arguments from the Exec+ or Query+ function.
// This allows for a developer to reuse the interface slice and save
// allocations. All method receivers are not thread safe. The returned interface
// slice is the same as `extArgs`.
func (a *Arguments) prepareArgs(extArgs ...interface{}) (_ string, _ []interface{}, err error) {
	if a.base.ärgErr != nil {
		return "", nil, errors.WithStack(a.base.ärgErr)
	}
	if len(a.base.cachedSQL) == 0 {
		return "", nil, errors.Empty.Newf("[dml] Arguments: The SQL string is empty.")
	}

	if a.base.source == dmlSourceInsert {
		return a.prepareArgsInsert(extArgs...)
	}

	if a.isEmpty() {
		a.hasNamedArgs = 1
		return string(a.base.cachedSQL), extArgs, nil
	}

	if a.hasNamedArgs == 0 {
		found := false
		a.hasNamedArgs = 1
		a.base.cachedSQL, a.base.qualifiedColumns, found = extractReplaceNamedArgs(a.base.cachedSQL, a.base.qualifiedColumns)

		switch la := len(a.arguments); true {
		case found:
			a.hasNamedArgs = 2
		case !found && len(a.recs) == 0 && la > 0:
			for _, arg := range a.arguments {
				if arg.name != "" {
					a.hasNamedArgs = 2
					break
				}
			}
		}
	}

	//var collectedArgs = append(arguments{}, a.arguments...) // TODO sync.Pool or not, investigate via benchmark, Tests are succeeding `-count=99 [-race]`

	sqlBuf := bufferpool.GetTwin()
	collectedArgs := pooledArgumentsGet()
	defer pooledArgumentsPut(collectedArgs, sqlBuf)
	collectedArgs = append(collectedArgs, a.arguments...)

	extArgs = append(extArgs, a.raw...)
	if collectedArgs, err = a.appendConvertedRecordsToArguments(collectedArgs); err != nil {
		return "", nil, errors.WithStack(err)
	}

	// Make a copy of the original SQL statement because it gets modified in the
	// worst case. Best case would be no modification and hence we don't need a
	// bytes.Buffer from the pool! TODO(CYS) optimize this and only acquire a
	// buffer from the pool in the worse case.
	if _, err := sqlBuf.First.Write(a.base.cachedSQL); err != nil {
		return "", nil, errors.WithStack(err)
	}

	// `switch` statement no suitable.
	if a.Options > 0 && len(extArgs) > 0 && len(a.recs) == 0 && len(a.arguments) == 0 {
		return "", nil, errors.NotAllowed.Newf("[dml] Interpolation/ExpandPlaceholders supports only Records and Arguments and not yet an interface slice.")
	}

	if a.Options&argOptionExpandPlaceholder != 0 {
		if phCount := bytes.Count(sqlBuf.First.Bytes(), placeHolderByte); phCount < a.Len() {
			if err := expandPlaceHolders(sqlBuf.Second, sqlBuf.First.Bytes(), collectedArgs); err != nil {
				return "", nil, errors.WithStack(err)
			}
			if _, err := sqlBuf.CopySecondToFirst(); err != nil {
				return "", nil, errors.WithStack(err)
			}
		}
	}
	if a.Options&argOptionInterpolate != 0 {
		if err := writeInterpolateBytes(sqlBuf.Second, sqlBuf.First.Bytes(), collectedArgs); err != nil {
			return "", nil, errors.Wrapf(err, "[dml] Interpolation failed: %q", sqlBuf.String())
		}
		return sqlBuf.Second.String(), nil, nil
	}

	return sqlBuf.First.String(), collectedArgs.Interfaces(extArgs...), nil
}

func (a *Arguments) appendConvertedRecordsToArguments(collectedArgs arguments) (arguments, error) {
	if a.base.templateStmtCount == 0 {
		a.base.templateStmtCount = 1
	}
	if len(a.arguments) == 0 && len(a.recs) == 0 {
		return collectedArgs, nil
	}

	if len(a.arguments) > 0 && len(a.recs) == 0 && a.hasNamedArgs < 2 {
		if a.base.templateStmtCount > 1 {
			collectedArgs = a.multiplyArguments(a.base.templateStmtCount)
		}
		// This is also a case where there are no records and only arguments and
		// those arguments do not contain any name. Then we can skip the column
		// mapper and ignore the qualifiedColumns.
		return collectedArgs, nil
	}

	qualifiedColumns := a.base.qualifiedColumns
	if lqca := len(a.QualifiedColumnsAliases); lqca > 0 {
		if lqca != len(a.base.qualifiedColumns) {
			return nil, errors.Mismatch.Newf("[dml] Argument.Record: QualifiedColumnsAliases slice %v and qualifiedColumns slice %v must have the same length", a.QualifiedColumnsAliases, a.base.qualifiedColumns)
		}
		qualifiedColumns = a.QualifiedColumnsAliases
	}

	// TODO refactor prototype and make it performant and beautiful code
	cm := NewColumnMap(len(a.arguments)+len(a.recs), "") // can use an arg pool Arguments sync.Pool, nope.

	for tsc := 0; tsc < a.base.templateStmtCount; tsc++ { // only in case of UNION statements in combination with a template SELECT, can be optimized later

		// `qualifiedColumns` contains the correct order as the place holders
		// appear in the SQL string.
		for _, identifier := range qualifiedColumns {
			// identifier can be either: column or qualifier.column or :column
			qualifier, column := splitColumn(identifier)
			// a.base.defaultQualifier is empty in case of INSERT statements

			column, isNamedArg := cutNamedArgStartStr(column) // removes the colon for named arguments
			cm.columns[0] = column                            // length is always one, as created in NewColumnMap

			if isNamedArg && len(a.arguments) > 0 {
				// if the colon : cannot be found then a simple place holder ? has been detected
				if err := a.MapColumns(cm); err != nil {
					return collectedArgs, errors.WithStack(err)
				}
			} else {
				found := false
				for _, qRec := range a.recs {
					if qRec.Qualifier == "" && qualifier != "" {
						qRec.Qualifier = a.base.defaultQualifier
					}
					if qRec.Qualifier != "" && qualifier == "" {
						qualifier = a.base.defaultQualifier
					}

					if qRec.Qualifier == qualifier {
						if err := qRec.Record.MapColumns(cm); err != nil {
							return collectedArgs, errors.WithStack(err)
						}
						found = true
					}
				}
				if !found {
					// If the argument cannot be found in the records then we assume the argument
					// has a numerical position and we grab just the next unnamed argument.
					if pArg, ok := a.nextUnnamedArg(); ok {
						cm.arguments = append(cm.arguments, pArg)
					}
				}
			}
		}
		a.nextUnnamedArgPos = 0
	}
	if len(cm.arguments) > 0 {
		collectedArgs = cm.arguments
	}

	return collectedArgs, nil
}

// prepareArgsInsert prepares the special arguments for an INSERT statement. The
// returned interface slice is the same as the `extArgs` slice. extArgs =
// external arguments.
func (a *Arguments) prepareArgsInsert(extArgs ...interface{}) (string, []interface{}, error) {

	//cm := pooledColumnMapGet()
	sqlBuf := bufferpool.GetTwin()
	defer bufferpool.PutTwin(sqlBuf)
	//defer pooledBufferColumnMapPut(cm, sqlBuf, nil)

	cm := NewColumnMap(16)
	cm.setColumns(a.base.qualifiedColumns)
	//defer bufferpool.PutTwin(sqlBuf)
	cm.arguments = append(cm.arguments, a.arguments...)
	{
		if _, err := sqlBuf.First.Write(a.base.cachedSQL); err != nil {
			return "", nil, errors.WithStack(err)
		}

		// Extract arguments from ColumnMapper and append them to `a.args`.
		// inserting multiple rows retrieved from a collection. There is no qualifier.
		//cm := NewColumnMap(MakeArgs(len(a.base.qualifiedColumns)*5/4), a.base.qualifiedColumns...)

		for _, qRec := range a.recs {
			if qRec.Qualifier != "" {
				return "", nil, errors.Fatal.Newf("[dml] Qualifier in %T is not supported and not needed.", qRec)
			}

			if err := qRec.Record.MapColumns(cm); err != nil {
				return "", nil, errors.WithStack(err)
			}
		}
	}

	extArgs = append(extArgs, a.raw...)
	totalArgLen := uint(len(cm.arguments) + len(extArgs))

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
		if len(extArgs) > 0 && len(a.recs) == 0 && len(cm.arguments) == 0 {
			return "", nil, errors.NotAllowed.Newf("[dml] Interpolation/ExpandPlaceholders supports only Records and Arguments and not yet an interface slice.")
		}

		if a.Options&argOptionInterpolate != 0 {
			if err := writeInterpolateBytes(sqlBuf.Second, sqlBuf.First.Bytes(), cm.arguments); err != nil {
				return "", nil, errors.Wrapf(err, "[dml] Interpolation failed: %q", sqlBuf.First.String())
			}
			return sqlBuf.Second.String(), nil, nil
		}
	}

	return sqlBuf.First.String(), cm.arguments.Interfaces(extArgs...), nil
}

// nextUnnamedArg returns an unnamed argument by its position.
func (a *Arguments) nextUnnamedArg() (argument, bool) {
	var unnamedCounter int
	lenArg := len(a.arguments)
	for i := 0; i < lenArg && a.nextUnnamedArgPos >= 0; i++ {
		if arg := a.arguments[i]; arg.name == "" {
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
		cm.arguments = append(cm.arguments, a.arguments...)
		return cm.Err()
	}
	for cm.Next() {
		// now a bit slow ... but will be refactored later with constant time
		// access, but first benchmark it. This for loop can be the 3rd one in the
		// overall chain.
		c := cm.Column()
		for _, arg := range a.arguments {
			// Case sensitive comparison
			if c != "" && arg.name == c {
				cm.arguments = append(cm.arguments, arg)
				break
			}
		}
	}
	return cm.Err()
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
func (as arguments) Len() int {
	var l int
	for _, arg := range as {
		l += arg.len()
	}
	return l
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
	const maxInt64 = 1<<63 - 1
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
		case uint:
			if vv > maxInt64 {
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

func (a *Arguments) add(v interface{}) *Arguments {
	if a == nil {
		a = MakeArgs(defaultArgumentsCapacity)
	}
	a.arguments = a.arguments.add(v)
	return a
}

func (a *Arguments) Record(qualifier string, record ColumnMapper) *Arguments {
	a.recs = append(a.recs, Qualify(qualifier, record))
	return a
}

// Arguments sets the internal arguments slice to the provided argument. Those
// are the slices Arguments, records and raw.
func (a *Arguments) Arguments(args *Arguments) *Arguments {
	// maybe deprecated this function.
	a.arguments = args.arguments
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
func (a *Arguments) Name(n string) *Arguments {
	a.arguments = append(a.arguments, argument{name: n})
	return a
}

// Reset resets the slice for new usage retaining the already allocated memory.
func (a *Arguments) Reset() *Arguments {
	for i := range a.recs {
		a.recs[i].Qualifier = ""
		a.recs[i].Record = nil
	}
	a.recs = a.recs[:0]
	a.arguments = a.arguments[:0]
	a.raw = a.raw[:0]
	a.nextUnnamedArgPos = 0
	return a
}

// DriverValue adds multiple of the same underlying values to the argument
// slice. When using different values, the last applied value wins and gets
// added to the argument slice. For example driver.Values of type `int` will
// result in []int.
func (a *Arguments) DriverValue(dvs ...driver.Valuer) *Arguments {
	if a == nil {
		a = MakeArgs(len(dvs))
	}
	a.arguments, a.base.ärgErr = driverValue(a.arguments, dvs...)
	return a
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

// DriverValues adds each driver.Value as its own argument to the argument
// slice. It panics if the underlying type is not one of the allowed of
// interface driver.Valuer.
func (a *Arguments) DriverValues(dvs ...driver.Valuer) *Arguments {
	if a == nil {
		a = MakeArgs(len(dvs))
	}
	a.arguments, a.base.ärgErr = driverValues(a.arguments, dvs...)
	return a
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
			panic(errors.NotSupported.Newf("[dml] Type %#v not supported in Driver.Values slice: %#v", t, dvs))
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

// WithDB sets the database query object.
func (a *Arguments) WithDB(db QueryExecPreparer) *Arguments {
	a.base.DB = db
	return a
}

/*********************************************
	LOAD / QUERY and EXEC functions
*********************************************/

var pooledColumnMap = sync.Pool{
	New: func() interface{} {
		return NewColumnMap(20, "")
	},
}

func pooledColumnMapGet() *ColumnMap {
	return pooledColumnMap.Get().(*ColumnMap)
}

func pooledBufferColumnMapPut(cm *ColumnMap, buf *bufferpool.TwinBuffer, fn func()) {
	if buf != nil {
		bufferpool.PutTwin(buf)
	}
	if fn != nil {
		fn()
	}
	cm.reset()
	pooledColumnMap.Put(cm)
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
	sqlStr, args, err := a.prepareArgs(args...)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("QueryRowContext", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	return a.base.DB.QueryRowContext(ctx, sqlStr, args...)
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
	cmr := pooledColumnMapGet() // this sync.Pool might not work correctly, write a complex test.
	defer pooledBufferColumnMapPut(cmr, nil, func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] Iterate.QueryContext.Rows.Close")
		}
	})

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
		defer log.WhenDone(a.base.Log).Debug("Load", log.String("id", a.base.id), log.Err(err), log.ObjectTypeOf("ColumnMapper", s), log.Uint64("row_count", rowCount))
	}

	r, err := a.query(ctx, args...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] Arguments.Load.QueryContext failed with queryID %q and ColumnMapper %T", a.base.id, s)
		return
	}
	cm := pooledColumnMapGet()
	defer pooledBufferColumnMapPut(cm, nil, func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] Arguments.Load.Rows.Close")
		}
		if rc, ok := s.(ioCloser); ok {
			if err2 := rc.Close(); err2 != nil && err == nil {
				err = errors.Wrap(err2, "[dml] Arguments.Load.ColumnMapper.Close")
			}
		}
	})

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
	rowCount = cm.Count
	return
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
	return
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
	return
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
	sqlStr, args, err2 := a.prepareArgs(args...)
	err = err2
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Query", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rows, err = a.base.DB.QueryContext(ctx, sqlStr, args...)
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
	sqlStr, args, err2 := a.prepareArgs(args...)
	err = err2
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Exec", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result, err = a.base.DB.ExecContext(ctx, sqlStr, args...)
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
