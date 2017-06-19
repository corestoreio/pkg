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
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

// https://www.adampalmer.me/iodigitalsec/2013/08/18/mysql_real_escape_string-wont-magically-solve-your-sql-injection-problems/

// Comparison functions and operators describe all available possibilities.
const (
	Null           Op = 'n'          // IS NULL
	NotNull        Op = 'N'          // IS NOT NULL
	In             Op = '∈'          // IN ?
	NotIn          Op = '∉'          // NOT IN ?
	Between        Op = 'b'          // BETWEEN ? AND ?
	NotBetween     Op = 'B'          // NOT BETWEEN ? AND ?
	Like           Op = 'l'          // LIKE ?
	NotLike        Op = 'L'          // NOT LIKE ?
	Greatest       Op = '≫'          // GREATEST(?,?,?) returns NULL if any argument is NULL.
	Least          Op = '≪'          // LEAST(?,?,?) If any argument is NULL, the result is NULL.
	Equal          Op = '='          // = ?
	NotEqual       Op = '≠'          // != ?
	Exists         Op = '∃'          // EXISTS(subquery)
	NotExists      Op = '∄'          // NOT EXISTS(subquery)
	Less           Op = '<'          // <
	Greater        Op = '>'          // >
	LessOrEqual    Op = '≤'          // <=
	GreaterOrEqual Op = '≥'          // >=
	Regexp         Op = 'r'          // REGEXP ?
	NotRegexp      Op = 'R'          // NOT REGEXP ?
	Xor            Op = '⊻'          // XOR ?
	SpaceShip      Op = '\U0001f680' // a <=> b is equivalent to a = b OR (a IS NULL AND b IS NULL) NULL-safe equal to operator
	Coalesce       Op = 'c'          // Returns the first non-NULL value in the list, or NULL if there are no non-NULL values.
)

// Op the Operator, defines comparison and operator functions used in any
// conditions. The upper case letter always negates.
// https://dev.mysql.com/doc/refman/5.7/en/comparison-operators.html
// https://mariadb.com/kb/en/mariadb/comparison-operators/
type Op rune

// String transforms the rune into a string.
func (o Op) String() string {
	if o < 1 {
		o = Equal
	}
	return string(o)
}

// With allows to use any argument with an operator.
func (o Op) With(arg Argument) Argument {
	return arg.applyOperator(o)
}

// Str uses string values for comparison.
func (o Op) Str(values ...string) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	return &argStrings{data: values, op: o}
}

// NullString uses nullable string values for comparison.
func (o Op) NullString(values ...NullString) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	if len(values) == 1 {
		values[0].op = o
		return values[0]
	}
	return argNullStrings{data: values, op: o}
}

// Float64 uses float64 values for comparison.
func (o Op) Float64(values ...float64) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	return &argFloat64s{data: values, op: o}
}

// NullFloat64 uses nullable float64 values for comparison.
func (o Op) NullFloat64(values ...NullFloat64) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	if len(values) == 1 {
		values[0].op = o
		return values[0]
	}
	return argNullFloat64s{data: values, op: o}
}

// Int64 uses int64 values for comparison.
func (o Op) Int64(values ...int64) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	return &argInt64s{data: values, op: o}
}

// NullInt64 uses nullable int64 values for comparison.
func (o Op) NullInt64(values ...NullInt64) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	if len(values) == 1 {
		values[0].op = o
		return values[0]
	}
	return argNullInt64s{data: values, op: o}
}

// Int uses int values for comparison.
func (o Op) Int(values ...int) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	return &argInts{data: values, op: o}
}

// Bool uses bool values for comparison.
func (o Op) Bool(values ...bool) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	return &argBools{data: values, op: o}
}

// NullBool uses nullable bool values for comparison.
func (o Op) NullBool(value NullBool) Argument {
	value.op = o
	return value
}

// Time uses time.Time values for comparison.
func (o Op) Time(values ...time.Time) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	return &argTimes{data: values, op: o}
}

// NullTime uses nullable time values for comparison.
func (o Op) NullTime(values ...NullTime) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	if len(values) == 1 {
		values[0].op = o
		return values[0]
	}
	return argNullTimes{data: values, op: o}
}

// Null is always a NULL.
func (o Op) Null() Argument {
	return argNull(o)
}

// Bytes uses a byte slice for comparison. Providing a nil argument returns a
// NULL type. Detects between valid UTF-8 strings and binary data. Later gets
// hex encoded.
func (o Op) Bytes(p ...[]byte) Argument {
	if len(p) == 0 {
		return argPlaceHolder(o)
	}
	return argBytes{data: p, op: o}
}

// Value uses driver.Valuers for comparison.
func (o Op) Value(values ...driver.Valuer) Argument {
	if len(values) == 0 {
		return argPlaceHolder(o)
	}
	return &argValue{data: values, op: o}
}

const (
	sqlStrNull = "NULL"
	sqlStar    = "*"
)

// SQL statement types and parts used as bit flag e.g. hint in
// ArgumentAssembler.AssembleArguments.
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

// ArgumentAssembler assembles arguments for CRUD statements. The `stmtType`
// variable contains a bit flag from the constants SQLStmt* and SQLPart* to
// allow the knowledge in which case the function AssembleArguments gets called.
// Any new arguments must be append to variable `args` and then returned.
// Variable `columns` contains the name of the requested columns. E.g. if the
// first requested column names `id` then the first appended argument must be an
// integer. Variable `columns` can additionally contain the names and/or
// expressions used in the WHERE, JOIN or HAVING clauses, if applicable for the
// SQL statement type. In case where stmtType has been set to SQLStmtInsert|SQLPartValues, the
// `columns` slice can be empty which means that all arguments are requested.
type ArgumentAssembler interface {
	AssembleArguments(stmtType int, args Arguments, columns []string) (Arguments, error)
}

// Argument transforms a value or values into an interface slice or encodes them
// into their textual representation to be used directly in a SQL query. This
// interface slice gets used in the database query functions as an argument. The
// underlying primitive type in the interface must be one of driver.Value
// allowed types.
type Argument interface {
	// applyOperator sets a comparison or logical operator. Please see the
	// constants Op for the different flags.
	applyOperator(Op) Argument
	// toIFace appends the value or values to interface slice and returns it.
	toIFace([]interface{}) []interface{}
	// writeTo writes the value correctly escaped to the queryWriter. It must
	// avoid SQL injections.
	writeTo(w queryWriter, position int) error
	// len returns the length of the available values. If the IN clause has been
	// activated then len returns 1. In case of an underlying place holder type
	// the returned length of cahensConstant
	len() int
	operator() Op
}

// Arguments representing multiple arguments.
type Arguments []Argument

func writePlaceHolderList(w queryWriter, arg Argument) {
	w.WriteByte('(')
	for j := 0; j < arg.len(); j++ {
		if j > 0 {
			w.WriteByte(',')
		}
		w.WriteByte('?')
	}
	w.WriteByte(')')
}

// argCount contains the number of primitives within an argument.
func writeOperator(w queryWriter, hasArg bool, arg Argument) (addArg bool) {
	// hasArg argument only used in cases where we have in the parent caller
	// function a sub-select. sub-selects do not need a place holder.
	switch arg.operator() {
	case Null:
		w.WriteString(" IS NULL")
	case NotNull:
		w.WriteString(" IS NOT NULL")
	case In:
		w.WriteString(" IN ")
		if hasArg {
			writePlaceHolderList(w, arg)
			addArg = true
		}
	case NotIn:
		w.WriteString(" NOT IN ")
		if hasArg {
			writePlaceHolderList(w, arg)
			addArg = true
		}
	case Like:
		w.WriteString(" LIKE ?")
		addArg = true
	case NotLike:
		w.WriteString(" NOT LIKE ?")
		addArg = true
	case Regexp:
		w.WriteString(" REGEXP ?")
		addArg = true
	case NotRegexp:
		w.WriteString(" NOT REGEXP ?")
		addArg = true
	case Between:
		w.WriteString(" BETWEEN ? AND ?")
		addArg = true
	case NotBetween:
		w.WriteString(" NOT BETWEEN ? AND ?")
		addArg = true
	case Greatest:
		w.WriteString(" GREATEST ")
		writePlaceHolderList(w, arg)
		addArg = true
	case Least:
		w.WriteString(" LEAST ")
		writePlaceHolderList(w, arg)
		addArg = true
	case Coalesce:
		w.WriteString(" COALESCE ")
		writePlaceHolderList(w, arg)
		addArg = true
	case Xor:
		w.WriteString(" XOR ?")
		addArg = true
	case Exists:
		w.WriteString(" EXISTS ")
		addArg = true
	case NotExists:
		w.WriteString(" NOT EXISTS ")
		addArg = true
	case Equal:
		w.WriteString(" = ")
		if hasArg {
			w.WriteByte('?')
			addArg = true
		}
	case NotEqual:
		w.WriteString(" != ")
		if hasArg {
			w.WriteByte('?')
			addArg = true
		}
	case Less:
		w.WriteString(" < ")
		if hasArg {
			w.WriteByte('?')
			addArg = true
		}
	case Greater:
		w.WriteString(" > ")
		if hasArg {
			w.WriteByte('?')
			addArg = true
		}
	case LessOrEqual:
		w.WriteString(" <= ")
		if hasArg {
			w.WriteByte('?')
			addArg = true
		}
	case GreaterOrEqual:
		w.WriteString(" >= ")
		if hasArg {
			w.WriteByte('?')
			addArg = true
		}
	case SpaceShip:
		w.WriteString(" <=> ")
		if hasArg {
			w.WriteByte('?')
			addArg = true
		}
	default:
		w.WriteString(" = ")
		if hasArg {
			w.WriteByte('?')
			addArg = true
		}
	}
	return
}

// len calculates the total length of all values
func (as Arguments) len() (tl int) {
	for _, a := range as {
		l := a.len()
		if l == cahensConstant {
			l = 1
		}
		tl += l
	}
	return
}

// Interfaces converts the underlying concrete types into an interface slice.
// Each entry in the interface is guaranteed to be one of the following values:
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
			args = append(args, ArgFloat64(float64(v)))
		case float64:
			args = append(args, ArgFloat64(v))
		case int64:
			args = append(args, ArgInt64(v))
		case int:
			args = append(args, ArgInt64(int64(v)))
		case int32:
			args = append(args, ArgInt64(int64(v)))
		case int16:
			args = append(args, ArgInt64(int64(v)))
		case int8:
			args = append(args, ArgInt64(int64(v)))
		case uint32:
			args = append(args, ArgInt64(int64(v)))
		case uint16:
			args = append(args, ArgInt64(int64(v)))
		case uint8:
			args = append(args, ArgInt64(int64(v)))
		case bool:
			args = append(args, ArgBool(v))
		case string:
			args = append(args, ArgString(v))
		case []byte:
			args = append(args, ArgBytes(v))
		case time.Time:
			args = append(args, ArgTime(v))
		case *time.Time:
			if v != nil {
				args = append(args, ArgTime(*v))
			}
		case nil:
			args = append(args, ArgNull())
		default:
			panic(errors.NewNotSupportedf("[dbr] iFaceToArgs type %#v not yet supported", v))
		}
	}
	return args
}

type argValue struct {
	op   Op
	data []driver.Valuer
}

func (a *argValue) toIFace(args []interface{}) []interface{} {
	for _, v := range a.data {
		args = append(args, v)
	}
	return args
}

func writeDriverValuer(w queryWriter, value driver.Valuer) error {
	if value == nil {
		_, err := w.WriteString("NULL")
		return err
	}
	val, err := value.Value()
	if err != nil {
		return errors.Wrapf(err, "[dbr] argValue.WriteTo: %#v", value)
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
		return errors.NewNotSupportedf("[dbr] argValue.WriteTo Type not yet supported: %#v", value)
	}
	return err
}

func (a *argValue) writeTo(w queryWriter, pos int) error {
	return writeDriverValuer(w, a.data[pos])
}

func (a *argValue) len() int {
	return len(a.data)
}

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a *argValue) applyOperator(op Op) Argument {
	a.op = op
	return a
}

func (a *argValue) operator() Op { return a.op }

// ArgValue allows to use any type which implements driver.Valuer interface.
// Implements interface Argument.
func ArgValue(args ...driver.Valuer) Argument {
	return &argValue{
		data: args,
	}
}

type argTimes struct {
	op   Op
	data []time.Time
}

func (a *argTimes) toIFace(args []interface{}) []interface{} {
	for _, v := range a.data {
		args = append(args, v)
	}
	return args
}

func (a *argTimes) writeTo(w queryWriter, pos int) error {
	dialect.EscapeTime(w, a.data[pos])
	return nil
}

func (a *argTimes) len() int {
	return len(a.data)
}

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a *argTimes) applyOperator(op Op) Argument {
	a.op = op
	return a
}

func (a *argTimes) operator() Op { return a.op }

// ArgTime adds a time.Time or a slice of times to the argument list. Providing
// no arguments returns a NULL type. Implements interface Argument.
func ArgTime(args ...time.Time) Argument {
	return &argTimes{data: args}
}

type argBytes struct {
	op   Op
	data [][]byte
}

func (a argBytes) toIFace(args []interface{}) []interface{} {
	return append(args, a.data)
}

func (a argBytes) writeTo(w queryWriter, pos int) (err error) {
	if !utf8.Valid(a.data[pos]) {
		dialect.EscapeBinary(w, a.data[pos])
	} else {
		dialect.EscapeString(w, string(a.data[pos]))
	}
	return nil
}

func (a argBytes) len() int {
	return len(a.data)
}

// Op not supported
func (a argBytes) applyOperator(op Op) Argument { a.op = op; return a }
func (a argBytes) operator() Op                 { return a.op }

// ArgBytes adds a byte slice to the argument list. Providing a nil argument
// returns a NULL type. Detects between valid UTF-8 strings and binary data. Later
// gets hex encoded.
func ArgBytes(p ...[]byte) Argument {
	if p == nil {
		return ArgNull()
	}
	return argBytes{data: p}
}

type argNull rune

func (i argNull) toIFace(args []interface{}) []interface{} {
	return append(args, nil)
}

func (i argNull) writeTo(w queryWriter, _ int) (err error) {
	_, err = w.WriteString("NULL")
	return err
}

func (i argNull) len() int { return 1 }

// Op not supported
func (i argNull) applyOperator(op Op) Argument { return argNull(op) }
func (i argNull) operator() Op {
	if i > 0 {
		return Op(i)
	}
	return Null
}

// ArgNull treats the argument as a SQL `IS NULL` or `NULL`. IN clause not
// supported. Implements interface Argument.
func ArgNull() Argument {
	return argNull(0)
}

// ArgString implements interface Argument.
type ArgString string

func (a ArgString) toIFace(args []interface{}) []interface{} {
	return append(args, string(a))
}

func (a ArgString) writeTo(w queryWriter, _ int) error {
	if !utf8.ValidString(string(a)) {
		return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", a)
	}
	dialect.EscapeString(w, string(a))
	return nil
}

func (a ArgString) len() int { return 1 }

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a ArgString) applyOperator(op Op) Argument {
	return &argStrings{
		data: []string{string(a)},
		op:   op,
	}
}
func (a ArgString) operator() Op { return 0 }

type argStrings struct {
	data []string
	op   Op
}

func (a *argStrings) toIFace(args []interface{}) []interface{} {
	for _, v := range a.data {
		args = append(args, v)
	}
	return args
}

func (a *argStrings) writeTo(w queryWriter, pos int) error {
	if !utf8.ValidString(a.data[pos]) {
		return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", a.data[pos])
	}
	dialect.EscapeString(w, a.data[pos])
	return nil
}

func (a *argStrings) len() int {
	return len(a.data)
}

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a *argStrings) applyOperator(op Op) Argument {
	a.op = op
	return a
}
func (a *argStrings) operator() Op { return a.op }

// ArgBool implements interface Argument.
type ArgBool bool

func (a ArgBool) toIFace(args []interface{}) []interface{} {
	return append(args, a == true)
}

func (a ArgBool) writeTo(w queryWriter, _ int) error {
	dialect.EscapeBool(w, a == true)
	return nil
}
func (a ArgBool) len() int { return 1 }

// Op not supported
func (a ArgBool) applyOperator(_ Op) Argument { return a }
func (a ArgBool) operator() Op                { return 0 }

type argBools struct {
	op   Op
	data []bool
}

func (a *argBools) toIFace(args []interface{}) []interface{} {
	for _, v := range a.data {
		args = append(args, v == true)
	}
	return args
}

func (a *argBools) writeTo(w queryWriter, pos int) error {
	dialect.EscapeBool(w, a.data[pos])
	return nil
}

func (a *argBools) len() int {
	return len(a.data)
}

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a *argBools) applyOperator(op Op) Argument {
	a.op = op
	return a
}
func (a *argBools) operator() Op { return a.op }

// ArgInt implements interface Argument.
type ArgInt int

func (a ArgInt) toIFace(args []interface{}) []interface{} {
	return append(args, int64(a))
}

func (a ArgInt) writeTo(w queryWriter, _ int) error {
	return writeInt64(w, int64(a))
}
func (a ArgInt) len() int { return 1 }

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a ArgInt) applyOperator(op Op) Argument {
	return &argInts{
		op:   op,
		data: []int{int(a)},
	}
}
func (a ArgInt) operator() Op { return 0 }

type argInts struct {
	op   Op
	data []int
}

func (a *argInts) toIFace(args []interface{}) []interface{} {
	for _, v := range a.data {
		args = append(args, int64(v))
	}
	return args
}

func (a *argInts) writeTo(w queryWriter, pos int) error {
	return writeInt64(w, int64(a.data[pos]))
}

func (a *argInts) len() int {
	return len(a.data)
}

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a *argInts) applyOperator(op Op) Argument {
	a.op = op
	return a
}

func (a *argInts) operator() Op { return a.op }

// ArgInt64 implements interface Argument.
type ArgInt64 int64

func (a ArgInt64) toIFace(args []interface{}) []interface{} {
	return append(args, int64(a))
}

func (a ArgInt64) writeTo(w queryWriter, _ int) error {
	return writeInt64(w, int64(a))
}
func (a ArgInt64) len() int { return 1 }

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a ArgInt64) applyOperator(op Op) Argument {
	return &argInt64s{
		op:   op,
		data: []int64{int64(a)},
	}
}
func (a ArgInt64) operator() Op { return 0 }

type argInt64s struct {
	op   Op
	data []int64
}

func (a *argInt64s) toIFace(args []interface{}) []interface{} {
	for _, v := range a.data {
		args = append(args, v)
	}
	return args
}

func (a *argInt64s) writeTo(w queryWriter, pos int) error {
	return writeInt64(w, int64(a.data[pos]))
}

func (a *argInt64s) len() int {
	return len(a.data)
}

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a *argInt64s) applyOperator(op Op) Argument {
	a.op = op
	return a
}

func (a *argInt64s) operator() Op { return a.op }

// ArgFloat64 implements interface Argument.
type ArgFloat64 float64

func (a ArgFloat64) toIFace(args []interface{}) []interface{} {
	return append(args, float64(a))
}

func (a ArgFloat64) writeTo(w queryWriter, _ int) error {
	return writeFloat64(w, float64(a))
}
func (a ArgFloat64) len() int { return 1 }

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a ArgFloat64) applyOperator(op Op) Argument {
	return &argFloat64s{
		op:   op,
		data: []float64{float64(a)},
	}
}
func (a ArgFloat64) operator() Op { return 0 }

type argFloat64s struct {
	op   Op
	data []float64
}

func (a *argFloat64s) toIFace(args []interface{}) []interface{} {
	for _, v := range a.data {
		args = append(args, v)
	}
	return args
}

func (a *argFloat64s) writeTo(w queryWriter, pos int) error {
	return writeFloat64(w, a.data[pos])
}

func (a *argFloat64s) len() int {
	return len(a.data)
}

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a *argFloat64s) applyOperator(op Op) Argument {
	a.op = op
	return a
}

func (a *argFloat64s) operator() Op { return a.op }

type expr struct {
	SQL string
	Arguments
	op Op
}

// ArgExpr at a SQL fragment with placeholders, and a slice of args to replace
// them with. Mostly used in UPDATE statements. Implements interface Argument.
func ArgExpr(sql string, args ...Argument) Argument {
	return &expr{SQL: sql, Arguments: args}
}

func (e *expr) toIFace(args []interface{}) []interface{} {
	for _, a := range e.Arguments {
		args = a.toIFace(args)
	}
	return args
}

func (e *expr) writeTo(w queryWriter, _ int) error {
	w.WriteString(e.SQL)
	return nil
}
func (e *expr) len() int { return 1 }

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (e *expr) applyOperator(op Op) Argument {
	e.op = op
	return e
}

func (e *expr) operator() Op { return e.op }

type argPlaceHolder rune

func (i argPlaceHolder) toIFace(args []interface{}) []interface{} {
	return args //append(args, nil)
}

func (i argPlaceHolder) writeTo(w queryWriter, _ int) (err error) {
	_, err = w.WriteString("? /*PLACEHOLDER*/") // maybe remove /*PLACEHOLDER*/ if it's annoying
	return err
}

func (i argPlaceHolder) len() int {
	return cahensConstant
}

// Op not supported
func (i argPlaceHolder) applyOperator(op Op) Argument { return argPlaceHolder(op) }
func (i argPlaceHolder) operator() Op {
	return Op(i)
}

func (i argPlaceHolder) GoString() string {
	return fmt.Sprintf("argPlaceHolder(%q)", i)
}
