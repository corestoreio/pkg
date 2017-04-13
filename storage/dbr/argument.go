package dbr

import (
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

// https://www.adampalmer.me/iodigitalsec/2013/08/18/mysql_real_escape_string-wont-magically-solve-your-sql-injection-problems/

// queryWriter at used to generate a query.
type queryWriter interface {
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
}

// Operators describe all available comparison operators. The upper case letter
// always negates.
// https://dev.mysql.com/doc/refman/5.7/en/comparison-operators.html
const (
	OperatorNull       byte = 'n' // IS NULL
	OperatorNotNull    byte = 'N' // IS NOT NULL
	OperatorIn         byte = 'i' // IN ?
	OperatorNotIn      byte = 'I' // NOT IN ?
	OperatorBetween    byte = 'b' // BETWEEN ? AND ?
	OperatorNotBetween byte = 'B' // NOT BETWEEN ? AND ?
	OperatorLike       byte = 'l' // LIKE ?
	OperatorNotLike    byte = 'L' // NOT LIKE ?
	OperatorGreatest   byte = 'g' // GREATEST(?,?,?)
	OperatorLeast      byte = 'e' // LEAST(?,?,?)
	OperatorEqual      byte = '=' // = ?
	OperatorNotEqual   byte = '!' // != ?
)

func writeOperator(w queryWriter, operator byte, hasArg bool) (addArg bool) {
	// hasArg argument only used in case we have in the parent caller function a
	// sub-select. sub-selects do not need a place holder.
	switch operator {
	case OperatorNull:
		w.WriteString(" IS NULL")
	case OperatorNotNull:
		w.WriteString(" IS NOT NULL")
	case OperatorIn:
		w.WriteString(" IN ")
		if hasArg {
			w.WriteRune('?')
			addArg = true
		}
	case OperatorNotIn:
		w.WriteString(" NOT IN ")
		if hasArg {
			w.WriteRune('?')
			addArg = true
		}
	case OperatorLike:
		w.WriteString(" LIKE ?")
		addArg = true
	case OperatorNotLike:
		w.WriteString(" NOT LIKE ?")
		addArg = true
	case OperatorBetween:
		w.WriteString(" BETWEEN ? AND ?")
		addArg = true
	case OperatorNotBetween:
		w.WriteString(" NOT BETWEEN ? AND ?")
		addArg = true
	case OperatorGreatest:
		w.WriteString(" GREATEST (?)")
		addArg = true
	case OperatorLeast:
		w.WriteString(" LEAST (?)")
		addArg = true
	case OperatorEqual:
		w.WriteString(" = ")
		if hasArg {
			w.WriteRune('?')
			addArg = true
		}
	case OperatorNotEqual:
		w.WriteString(" != ")
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

// StatementTypes identifies in the ArgumentGenerater interface what kind of
// operation requests the arguments.
const (
	StatementTypeDelete byte = 'd'
	StatementTypeInsert byte = 'i'
	StatementTypeSelect byte = 's'
	StatementTypeUpdate byte = 'u'
)

// ArgumentGenerater knows how to generate a record for one table row.
type ArgumentGenerater interface {
	// GenerateArguments generates a single new database record depending
	// on the requested column names. Each Argument gets mapped to the column
	// name. E.g. first column name at "id" then the first returned Argument in
	// the slice must be an integer.
	// GenerateUpdateArguments generates an argument set to be used in the SET
	// clause columns and in the WHERE statement. The `columns` argument
	// contains the name of the columns which are used in the SET clause. The
	// `where` argument contains a list of column names or even expressions
	// which get used in the WHERE statement. These names allows to filter and
	// generate the needed arguments.
	GenerateArguments(statementType byte, columns, condition []string) (Arguments, error)
}

// Argument transforms your value or values into an interface slice. This
// interface slice gets used in the database query functions at an argument. The
// underlying type in the interface must be one of driver.Value allowed types.
type Argument interface {
	// Operator sets a comparison or logical operator. Please see the constants
	// Operator* for the different flags. An underscore in the argument list of
	// a type indicates that no operator is yet supported.
	Operator(opt byte) Argument
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
	case OperatorIn, OperatorNotIn, OperatorGreatest, OperatorLeast:
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
	opt  byte
	data []time.Time
}

func (a argTimes) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v)
	}
}

func (a argTimes) writeTo(w queryWriter, pos int) error {
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

func (a argTimes) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

func (a argTimes) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a argTimes) operator() byte { return a.opt }

// ArgTime adds a time.Time or a slice of times to the argument list.
// Providing no arguments returns a NULL type.
func ArgTime(args ...time.Time) Argument {
	return argTimes{data: args}
}

type argBytes []byte

func (a argBytes) toIFace(args *[]interface{}) {
	*args = append(*args, []byte(a))
}

func (a argBytes) writeTo(w queryWriter, _ int) error {
	dialect.EscapeBinary(w, a)
	return nil
}

func (a argBytes) len() int                 { return 1 }
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

func (i argNull) len() int                 { return 1 }
func (i argNull) Operator(_ byte) Argument { return i }
func (i argNull) operator() byte {
	switch i {
	case 10:
		return OperatorNull
	case 20:
		return OperatorNotNull
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

// does not allocate when using as a argument but does neither support the Operator function.
//type argString string
//
//func (a argString) toIFace(args *[]interface{}) {
//	*args = append(*args, string(a))
//}
//
//func (a argString) writeTo(w queryWriter, _ int) error {
//	if !utf8.ValidString(string(a)) {
//		return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", a)
//	}
//	dialect.EscapeString(w, string(a))
//	return nil
//}
//
//func (a argString) len() int                 { return 1 }
//func (a argString) Operator(_ byte) Argument { return a }
//func (a argString) operator() byte           { return 0 }

type argStrings struct {
	opt  byte
	data []string
}

func (a argStrings) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v)
	}
}

func (a argStrings) writeTo(w queryWriter, pos int) error {
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

func (a argStrings) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

func (a argStrings) Operator(opt byte) Argument {
	a.opt = opt
	return a
}
func (a argStrings) operator() byte { return a.opt }

// ArgString adds a string or a slice of strings to the argument list.
// Providing no arguments returns a NULL type.
// All arguments mut be a valid utf-8 string.
func ArgString(args ...string) Argument {
	//if len(args) == 1 {
	//	return argString(args[0])
	//}
	return argStrings{data: args}
}

type argBool bool

func (a argBool) toIFace(args *[]interface{}) {
	*args = append(*args, a == true)
}

func (a argBool) writeTo(w queryWriter, _ int) error {
	dialect.EscapeBool(w, a == true)
	return nil
}
func (a argBool) len() int                 { return 1 }
func (a argBool) Operator(_ byte) Argument { return a }
func (a argBool) operator() byte           { return 0 }

type argBools struct {
	opt  byte
	data []bool
}

func (a argBools) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v == true)
	}
}

func (a argBools) writeTo(w queryWriter, pos int) error {
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

func (a argBools) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

func (a argBools) Operator(opt byte) Argument {
	a.opt = opt
	return a
}
func (a argBools) operator() byte { return a.opt }

// ArgBool adds a string or a slice of bools to the argument list.
// Providing no arguments returns a NULL type.
func ArgBool(args ...bool) Argument {
	if len(args) == 1 {
		return argBool(args[0])
	}
	return argBools{data: args}
}

type argInt int

func (a argInt) toIFace(args *[]interface{}) {
	*args = append(*args, int64(a))
}

func (a argInt) writeTo(w queryWriter, _ int) error {
	_, err := w.WriteString(strconv.FormatInt(int64(a), 10))
	return err
}
func (a argInt) len() int                 { return 1 }
func (a argInt) Operator(_ byte) Argument { return a }
func (a argInt) operator() byte           { return 0 }

type argInts struct {
	opt  byte
	data []int
}

func (a argInts) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, int64(v))
	}
}

func (a argInts) writeTo(w queryWriter, pos int) error {
	if isNotIn(a.operator()) {
		_, err := w.WriteString(strconv.FormatInt(int64(a.data[pos]), 10))
		return err
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		w.WriteString(strconv.FormatInt(int64(v), 10))
		if i < l {
			w.WriteRune(',')
		}
	}
	w.WriteRune(')')
	return nil
}

func (a argInts) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

func (a argInts) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a argInts) operator() byte { return a.opt }

// ArgInt adds an integer or a slice of integers to the argument list.
// Providing no arguments returns a NULL type.
func ArgInt(args ...int) Argument {
	if len(args) == 1 {
		return argInt(args[0])
	}
	return argInts{data: args}
}

type argInt64 int64

func (a argInt64) toIFace(args *[]interface{}) {
	*args = append(*args, int64(a))
}

func (a argInt64) writeTo(w queryWriter, _ int) error {
	_, err := w.WriteString(strconv.FormatInt(int64(a), 10))
	return err
}
func (a argInt64) len() int                 { return 1 }
func (a argInt64) Operator(_ byte) Argument { return a }
func (a argInt64) operator() byte           { return 0 }

type argInt64s struct {
	opt  byte
	data []int64
}

func (a argInt64s) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, int64(v))
	}
}

func (a argInt64s) writeTo(w queryWriter, pos int) error {
	if isNotIn(a.operator()) {
		_, err := w.WriteString(strconv.FormatInt(a.data[pos], 10))
		return err
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		w.WriteString(strconv.FormatInt(int64(v), 10))
		if i < l {
			w.WriteRune(',')
		}
	}
	w.WriteRune(')')
	return nil
}

func (a argInt64s) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

func (a argInt64s) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a argInt64s) operator() byte { return a.opt }

// ArgInt64 adds an integer or a slice of integers to the argument list.
// Providing no arguments returns a NULL type.
func ArgInt64(args ...int64) Argument {
	if len(args) == 1 {
		return argInt64(args[0])
	}
	return argInt64s{data: args}
}

type argFloat64 float64

func (a argFloat64) toIFace(args *[]interface{}) {
	*args = append(*args, float64(a))
}

func (a argFloat64) writeTo(w queryWriter, _ int) error {
	_, err := w.WriteString(strconv.FormatFloat(float64(a), 'f', -1, 64))
	return err
}
func (a argFloat64) len() int                 { return 1 }
func (a argFloat64) Operator(_ byte) Argument { return a }
func (a argFloat64) operator() byte           { return 0 }

type argFloat64s struct {
	op   byte
	data []float64
}

func (a argFloat64s) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, float64(v))
	}
}

func (a argFloat64s) writeTo(w queryWriter, pos int) error {
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

func (a argFloat64s) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

func (a argFloat64s) Operator(opt byte) Argument {
	a.op = opt
	return a
}

func (a argFloat64s) operator() byte { return a.op }

// ArgFloat64 adds a float64 or a slice of floats to the argument list.
// Providing no arguments returns a NULL type.
func ArgFloat64(args ...float64) Argument {
	if len(args) == 1 {
		return argFloat64(args[0])
	}
	return argFloat64s{data: args}
}

type expr struct {
	SQL string
	Arguments
}

// Expr at a SQL fragment with placeholders, and a slice of args to replace them
// with. Mostly used in UPDATE statements.
func Expr(sql string, args ...Argument) Argument {
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
func (e *expr) len() int                 { return 1 }
func (e *expr) Operator(_ byte) Argument { return e }
func (e *expr) operator() byte           { return 0 }
