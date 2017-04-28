package dbr

import (
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

// https://www.adampalmer.me/iodigitalsec/2013/08/18/mysql_real_escape_string-wont-magically-solve-your-sql-injection-problems/

// Comparions functions and operators describe all available possibilities. The upper case letter
// always negates.
// https://dev.mysql.com/doc/refman/5.7/en/comparison-operators.html
const (
	Null           byte = 'n' // IS NULL
	NotNull        byte = 'N' // IS NOT NULL
	In             byte = 'i' // IN ?
	NotIn          byte = 'I' // NOT IN ?
	Between        byte = 'b' // BETWEEN ? AND ?
	NotBetween     byte = 'B' // NOT BETWEEN ? AND ?
	Like           byte = 'l' // LIKE ?
	NotLike        byte = 'L' // NOT LIKE ?
	Greatest       byte = 'g' // GREATEST(?,?,?)
	Least          byte = 'a' // LEAST(?,?,?)
	Equal          byte = '=' // = ?
	NotEqual       byte = '!' // != ?
	Exists         byte = 'e' // EXISTS(subquery)
	NotExists      byte = 'E' // NOT EXISTS(subquery)
	Less           byte = '<' // <
	Greater        byte = '>' // >
	LessOrEqual    byte = '{' // <=
	GreaterOrEqual byte = '}' // >=
	Regexp         byte = 'r' // REGEXP ?
	NotRegexp      byte = 'R' // NOT REGEXP ?
	Xor            byte = 'o' // XOR ?
)

const (
	sqlStrNull = "NULL"
)

func writeOperator(w queryWriter, operator byte, hasArg bool) (addArg bool) {
	// hasArg argument only used in case we have in the parent caller function a
	// sub-select. sub-selects do not need a place holder.
	switch operator {
	case Null:
		w.WriteString(" IS NULL")
	case NotNull:
		w.WriteString(" IS NOT NULL")
	case In:
		w.WriteString(" IN ")
		if hasArg {
			w.WriteRune('?')
			addArg = true
		}
	case NotIn:
		w.WriteString(" NOT IN ")
		if hasArg {
			w.WriteRune('?')
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
		w.WriteString(" GREATEST (?)")
		addArg = true
	case Least:
		w.WriteString(" LEAST (?)")
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
			w.WriteRune('?')
			addArg = true
		}
	case NotEqual:
		w.WriteString(" != ")
		if hasArg {
			w.WriteRune('?')
			addArg = true
		}
	case Less:
		w.WriteString(" < ")
		if hasArg {
			w.WriteRune('?')
			addArg = true
		}
	case Greater:
		w.WriteString(" > ")
		if hasArg {
			w.WriteRune('?')
			addArg = true
		}
	case LessOrEqual:
		w.WriteString(" <= ")
		if hasArg {
			w.WriteRune('?')
			addArg = true
		}
	case GreaterOrEqual:
		w.WriteString(" >= ")
		if hasArg {
			w.WriteRune('?')
			addArg = true
		}
	default:
		w.WriteString(" = ")
		if hasArg {
			w.WriteRune('?')
			addArg = true
		}
	}
	return
}

// InsertArgProducer produces arguments for a SQL INSERT statement. Any new
// arguments must be append to variable `args` and then returned. Variable
// `columns` contains the name of the requested columns. E.g. if the first
// requested column names `id` then the first appended argument must be an
// integer. An empty or nil `columns` variable must append all requested columns
// to the `args` variable.
type InsertArgProducer interface {
	ProduceInsertArgs(args Arguments, columns []string) (Arguments, error)
}

// UpdateArgProducer produces arguments for a SQL UPDATE statement. Any new
// arguments must be append to variable `args` and then returned. Variable
// `columns` contains the name of the requested columns. E.g. if the first
// requested column names `id` then the first appended argument must be an
// integer. Vairable `condition` contains the names and/or expressions used in
// the WHERE or ON clause.
type UpdateArgProducer interface {
	ProduceUpdateArgs(args Arguments, columns, condition []string) (Arguments, error)
}

// Argument transforms your value or values into an interface slice or encodes
// them into textual representation to be used directly in a SQL query. This
// interface slice gets used in the database query functions at an argument. The
// underlying type in the interface must be one of driver.Value allowed types.
type Argument interface {
	// Operator sets a comparison or logical operator. Please see the constants
	// Operator* for the different flags. An underscore in the argument list of
	// a type indicates that no operator is yet supported.
	Operator(byte) Argument
	toIFace(*[]interface{})
	// writeTo writes the value correctly escaped to the queryWriter. It must
	// avoid SQL injections.
	writeTo(w queryWriter, position int) error
	// len returns the length of the available values. If the IN clause has been
	// activated then len returns 1.
	len() int
	operator() byte
}

// Arguments representing multiple arguments.
type Arguments []Argument

// len calculates the total length of all values
func (as Arguments) len() (l int) {
	for _, a := range as {
		l += a.len()
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
		a.toIFace(&ret)
	}
	return ret
}

func isNotIn(o byte) bool {
	switch o {
	case In, NotIn, Greatest, Least:
		return false
	}
	return true
}

//func ArgValuer(args ...driver.Valuer) (Argument, error) {
//	if len(args) == 1 {
//		dv, err := args[0].Value()
//		if err != nil {
//			return nil, errors.Wrap(err, "[dbr] args[0].Value")
//		}
//		switch v := dv.(type) {
//		case int64:
//			return ArgInt64(v), nil
//		case []int64:
//			return ArgInt64(v...), nil
//		case float64:
//			return ArgFloat64(v), nil
//		case []float64:
//			return ArgFloat64(v...), nil
//		case bool:
//			return ArgBool(v), nil
//		case []bool:
//			return ArgBool(v...), nil
//		case []byte:
//			return ArgBytes(v), nil
//		case string:
//			return ArgString(v), nil
//		case []string:
//			return ArgString(v...), nil
//		case time.Time:
//			return ArgTime(v), nil
//		case []time.Time:
//			return ArgTime(v...), nil
//		case nil:
//			return ArgNull(), nil
//		default:
//			return nil, errors.NewNotSupportedf("[dbr] Argument %#v not supported", dv)
//		}
//	}
//	return nil, nil
//}

type argTimes struct {
	op   byte
	data []time.Time
}

func (a *argTimes) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v)
	}
}

func (a *argTimes) writeTo(w queryWriter, pos int) error {
	if isNotIn(a.operator()) {
		dialect.EscapeTime(w, a.data[pos])
		return nil
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		dialect.EscapeTime(w, v)
		if i < l {
			w.WriteRune(',')
		}
	}
	w.WriteRune(')')
	return nil
}

func (a *argTimes) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a *argTimes) Operator(op byte) Argument {
	a.op = op
	return a
}

func (a *argTimes) operator() byte { return a.op }

// ArgTime adds a time.Time or a slice of times to the argument list.
// Providing no arguments returns a NULL type.
func ArgTime(args ...time.Time) Argument {
	return &argTimes{data: args}
}

type argBytes []byte

func (a argBytes) toIFace(args *[]interface{}) {
	*args = append(*args, []byte(a))
}

func (a argBytes) writeTo(w queryWriter, _ int) error {
	dialect.EscapeBinary(w, a)
	return nil
}

func (a argBytes) len() int { return 1 }

// Operator not supported
func (a argBytes) Operator(_ byte) Argument { return a }
func (a argBytes) operator() byte           { return 0 }

// ArgBytes adds a byte slice to the argument list.
// Providing a nil argument returns a NULL type.
// IN clause not supported.
func ArgBytes(p []byte) Argument {
	if p == nil {
		return ArgNull()
	}
	return argBytes(p)
}

type argNull uint8

func (i argNull) toIFace(args *[]interface{}) {
	*args = append(*args, nil)
}

func (i argNull) writeTo(w queryWriter, _ int) error {
	_, err := w.WriteString(`NULL`)
	return err
}

func (i argNull) len() int { return 1 }

// Operator not supported
func (i argNull) Operator(_ byte) Argument { return i }
func (i argNull) operator() byte {
	switch i {
	case 10:
		return Null
	case 20:
		return NotNull
	}
	return 0
}

// ArgNull treats the argument as a SQL `IS NULL` or `NULL`.
// IN clause not supported.
func ArgNull() Argument {
	return argNull(10)
}

// ArgNotNull treats the argument as a SQL `IS NOT NULL`.
// IN clause not supported.
func ArgNotNull() Argument {
	return argNull(20)
}

// argString implements interface Argument but does not allocate.
type argString string

func (a argString) toIFace(args *[]interface{}) {
	*args = append(*args, string(a))
}

func (a argString) writeTo(w queryWriter, _ int) error {
	if !utf8.ValidString(string(a)) {
		return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", a)
	}
	dialect.EscapeString(w, string(a))
	return nil
}

func (a argString) len() int { return 1 }

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a argString) Operator(op byte) Argument {
	return &argStrings{
		data: []string{string(a)},
		op:   op,
	}
}
func (a argString) operator() byte { return 0 }

type argStrings struct {
	data []string
	op   byte
}

func (a *argStrings) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v)
	}
}

func (a *argStrings) writeTo(w queryWriter, pos int) error {
	if isNotIn(a.operator()) {
		if !utf8.ValidString(a.data[pos]) {
			return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", a.data[pos])
		}
		dialect.EscapeString(w, a.data[pos])
		return nil
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		if !utf8.ValidString(v) {
			return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", v)
		}
		dialect.EscapeString(w, v)
		if i < l {
			w.WriteRune(',')
		}
	}
	w.WriteRune(')')
	return nil
}

func (a *argStrings) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a *argStrings) Operator(op byte) Argument {
	a.op = op
	return a
}
func (a *argStrings) operator() byte { return a.op }

// ArgString adds a string or a slice of strings to the argument list.
// Providing no arguments returns a NULL type.
// All arguments mut be a valid utf-8 string.
func ArgString(args ...string) Argument {
	if len(args) == 1 {
		return argString(args[0])
	}
	return &argStrings{data: args}
}

type argBool bool

func (a argBool) toIFace(args *[]interface{}) {
	*args = append(*args, a == true)
}

func (a argBool) writeTo(w queryWriter, _ int) error {
	dialect.EscapeBool(w, a == true)
	return nil
}
func (a argBool) len() int { return 1 }

// Operator not supported
func (a argBool) Operator(_ byte) Argument { return a }
func (a argBool) operator() byte           { return 0 }

type argBools struct {
	op   byte
	data []bool
}

func (a *argBools) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v == true)
	}
}

func (a *argBools) writeTo(w queryWriter, pos int) error {
	if isNotIn(a.operator()) {
		dialect.EscapeBool(w, a.data[pos])
		return nil
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		dialect.EscapeBool(w, v == true)
		if i < l {
			w.WriteRune(',')
		}
	}
	w.WriteRune(')')
	return nil
}

func (a *argBools) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a *argBools) Operator(op byte) Argument {
	a.op = op
	return a
}
func (a *argBools) operator() byte { return a.op }

// ArgBool adds a string or a slice of bools to the argument list.
// Providing no arguments returns a NULL type.
func ArgBool(args ...bool) Argument {
	if len(args) == 1 {
		return argBool(args[0])
	}
	return &argBools{data: args}
}

// argInt implements interface Argument but does not allocate.
type argInt int

func (a argInt) toIFace(args *[]interface{}) {
	*args = append(*args, int64(a))
}

func (a argInt) writeTo(w queryWriter, _ int) error {
	_, err := w.WriteString(strconv.FormatInt(int64(a), 10))
	return err
}
func (a argInt) len() int { return 1 }

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a argInt) Operator(op byte) Argument {
	return &argInts{
		op:   op,
		data: []int{int(a)},
	}
}
func (a argInt) operator() byte { return 0 }

type argInts struct {
	op   byte
	data []int
}

func (a *argInts) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, int64(v))
	}
}

func (a *argInts) writeTo(w queryWriter, pos int) error {
	if isNotIn(a.operator()) {
		_, err := w.WriteString(strconv.Itoa(a.data[pos]))
		return err
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		w.WriteString(strconv.Itoa(v))
		if i < l {
			w.WriteRune(',')
		}
	}
	w.WriteRune(')')
	return nil
}

func (a *argInts) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a *argInts) Operator(op byte) Argument {
	a.op = op
	return a
}

func (a *argInts) operator() byte { return a.op }

// ArgInt adds an integer or a slice of integers to the argument list.
// Providing no arguments returns a NULL type.
func ArgInt(args ...int) Argument {
	if len(args) == 1 {
		return argInt(args[0])
	}
	return &argInts{data: args}
}

// argInt64 implements interface Argument but does not allocate.
type argInt64 int64

func (a argInt64) toIFace(args *[]interface{}) {
	*args = append(*args, int64(a))
}

func (a argInt64) writeTo(w queryWriter, _ int) error {
	_, err := w.WriteString(strconv.FormatInt(int64(a), 10))
	return err
}
func (a argInt64) len() int { return 1 }

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a argInt64) Operator(op byte) Argument {
	return &argInt64s{
		op:   op,
		data: []int64{int64(a)},
	}
}
func (a argInt64) operator() byte { return 0 }

type argInt64s struct {
	op   byte
	data []int64
}

func (a *argInt64s) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v)
	}
}

func (a *argInt64s) writeTo(w queryWriter, pos int) error {
	if isNotIn(a.operator()) {
		_, err := w.WriteString(strconv.FormatInt(a.data[pos], 10))
		return err
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		w.WriteString(strconv.FormatInt(v, 10))
		if i < l {
			w.WriteRune(',')
		}
	}
	w.WriteRune(')')
	return nil
}

func (a *argInt64s) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a *argInt64s) Operator(op byte) Argument {
	a.op = op
	return a
}

func (a *argInt64s) operator() byte { return a.op }

// ArgInt64 adds an integer or a slice of integers to the argument list.
// Providing no arguments returns a NULL type.
func ArgInt64(args ...int64) Argument {
	if len(args) == 1 {
		return argInt64(args[0])
	}
	return &argInt64s{data: args}
}

type argFloat64 float64

func (a argFloat64) toIFace(args *[]interface{}) {
	*args = append(*args, float64(a))
}

func (a argFloat64) writeTo(w queryWriter, _ int) error {
	_, err := w.WriteString(strconv.FormatFloat(float64(a), 'f', -1, 64))
	return err
}
func (a argFloat64) len() int { return 1 }

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a argFloat64) Operator(op byte) Argument {
	return &argFloat64s{
		op:   op,
		data: []float64{float64(a)},
	}
}
func (a argFloat64) operator() byte { return 0 }

type argFloat64s struct {
	op   byte
	data []float64
}

func (a *argFloat64s) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v)
	}
}

func (a *argFloat64s) writeTo(w queryWriter, pos int) error {
	if isNotIn(a.operator()) {
		_, err := w.WriteString(strconv.FormatFloat(a.data[pos], 'f', -1, 64))
		return err
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		w.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
		if i < l {
			w.WriteRune(',')
		}
	}
	w.WriteRune(')')
	return nil
}

func (a *argFloat64s) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a *argFloat64s) Operator(op byte) Argument {
	a.op = op
	return a
}

func (a *argFloat64s) operator() byte { return a.op }

// ArgFloat64 adds a float64 or a slice of floats to the argument list.
// Providing no arguments returns a NULL type.
func ArgFloat64(args ...float64) Argument {
	if len(args) == 1 {
		return argFloat64(args[0])
	}
	return &argFloat64s{data: args}
}

type expr struct {
	SQL string
	Arguments
	op byte
}

// ArgExpr at a SQL fragment with placeholders, and a slice of args to replace them
// with. Mostly used in UPDATE statements.
func ArgExpr(sql string, args ...Argument) Argument {
	return &expr{SQL: sql, Arguments: args}
}

func (e *expr) toIFace(args *[]interface{}) {
	for _, a := range e.Arguments {
		a.toIFace(args)
	}
}

func (e *expr) writeTo(w queryWriter, _ int) error {
	w.WriteString(e.SQL)
	return nil
}
func (e *expr) len() int { return 1 }

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (e *expr) Operator(op byte) Argument {
	e.op = op
	return e
}

func (e *expr) operator() byte { return e.op }

// for type subQuery see function SubSelect

//type argSubSelect struct {
//	// buf contains the cached SQL string
//	buf *bytes.Buffer
//	// args contains the arguments after calling ToSQL
//	args Arguments
//	s    *Select
//	op   byte
//}

// I don't know anymore where I would have needed this ... but once the idea
// and a real world use case pops up, I'm gonna implement it. Until then use the function
// SubSelect(rawStatementOrColumnName string, operator byte, s *Select) ConditionArg
//// ArgSubSelect
//// The written sub-select gets wrapped in parenthesis: (SELECT ...)
//func ArgSubSelect(s *Select) Argument {
//	return &argSubSelect{s: s}
//}
//
//func (e *argSubSelect) toIFace(args *[]interface{}) {
//
//	if e.buf == nil {
//		e.buf = new(bytes.Buffer)
//		var err error
//		e.args, err = e.s.toSQL(e.buf) // can be optimized later
//		if err != nil {
//			*args = append(*args, err) // not that optimal :-(
//		} else {
//			for _, a := range e.args {
//				a.toIFace(args)
//			}
//		}
//		return
//	}
//	for _, a := range e.args {
//		a.toIFace(args)
//	}
//}
//
//func (e *argSubSelect) writeTo(w queryWriter, _ int) (err error) {
//	if e.buf == nil {
//		e.buf = new(bytes.Buffer)
//		e.buf.WriteRune('(')
//		e.args, err = e.s.toSQL(e.buf)
//		if err != nil {
//			return errors.Wrap(err, "[dbr] argSubSelect.writeTo")
//		}
//		e.buf.WriteRune(')')
//	}
//	_, err = w.WriteString(e.buf.String())
//	return err
//}
//
//func (e *argSubSelect) len() int { return 1 }
//
//// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
//// the constants Operator*.
//func (e *argSubSelect) Operator(op byte) Argument {
//	e.op = op
//	return e
//}
//func (e *argSubSelect) operator() byte { return e.op }
