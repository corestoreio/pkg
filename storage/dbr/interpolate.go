package dbr

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

// https://www.adampalmer.me/iodigitalsec/2013/08/18/mysql_real_escape_string-wont-magically-solve-your-sql-injection-problems/

// queryWriter at used to generate a query.
type queryWriter interface {
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
}

// If the Argument interface returns in the null() function we can have three
// states: 0 means ignoring; argOptionNull writes IS NULL and argOptionNotNull
// writes IS NOT NULL. not all types do support NULL ... but can be implemented
// later.
const (
	argOptionNull uint = 1 << iota
	argOptionNotNull
	argOptionIN
	argOptionBetween // TODO implement between, extend pubic interface
)

// RecordGenerater knows how to generate a record for one database row.
type RecordGenerater interface {
	// Record creates a single new database record depending on the requested
	// column names. Each Argument gets mapped to the column name. E.g. first
	// column name at "id" then the first returned Argument in the slice must be
	// an integer.
	Record(columns ...string) (Arguments, error)
}

// Argument transforms your value or values into an interface slice. This
// interface slice gets used in the database query functions at an argument. The
// underlying type in the interface must be one of driver.Value allowed types.
type Argument interface {
	// IN sets an internal flag that the slice value will be used for
	// an IN clause query. Default: all values will be treated as single
	// arguments.
	IN() Argument
	// Between
	// BETWEEN() Argument
	// we must use interface at an argument because of the nested `where` functions.
	toIFace(*[]interface{})
	// writeTo writes the value correctly escaped to the queryWriter. It must avoid
	// SQL injections.
	writeTo(w queryWriter, position int) error
	// len returns the length of the available values. If the IN clause has been activated
	// then len returns 1.
	len() int
	options() uint
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

// Interfaces conerts the underlying concrete types into an interface slice.
// Each entry in the interface is guaranteed to be one of the following values:
// []byte, bool, float64, int64, string or time.Time.
// Use driver.IsValue() for a check.
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

//func ArgValuer(args ...driver.Valuer) Argument {
//
//}

type argTime struct {
	time.Time
}

func (at argTime) toIFace(args *[]interface{}) {
	*args = append(*args, at.Time)
}

func (at argTime) writeTo(w queryWriter, _ int) error {
	dialect.EscapeTime(w, at.Time)
	return nil
}
func (at argTime) len() int      { return 1 }
func (at argTime) IN() Argument  { return at }
func (at argTime) options() uint { return 0 }

type argTimes struct {
	opts uint
	data []time.Time
}

func (a argTimes) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v)
	}
}

func (a argTimes) writeTo(w queryWriter, pos int) error {
	if a.options() == 0 {
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
	if a.options() == 0 {
		return len(a.data)
	}
	return 1
}

func (a argTimes) IN() Argument {
	a.opts = argOptionIN
	return a
}

func (a argTimes) options() uint { return a.opts }

// ArgTime adds a time.Time or a slice of times to the argument list.
// Providing no arguments returns a NULL type.
func ArgTime(args ...time.Time) Argument {
	if len(args) == 1 {
		return argTime{Time: args[0]}
	}
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

func (a argBytes) len() int      { return 1 }
func (a argBytes) IN() Argument  { return a }
func (a argBytes) options() uint { return 0 }

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

func (i argNull) len() int     { return 1 }
func (i argNull) IN() Argument { return i }
func (i argNull) options() uint {
	switch i {
	case 10:
		return argOptionNull
	case 20:
		return argOptionNotNull
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

func (s argString) len() int      { return 1 }
func (s argString) IN() Argument  { return s }
func (s argString) options() uint { return 0 }

type argStrings struct {
	opts uint
	data []string
}

func (a argStrings) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v)
	}
}

func (a argStrings) writeTo(w queryWriter, pos int) error {
	if a.options() == 0 {
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
	if a.options() == 0 {
		return len(a.data)
	}
	return 1
}

func (a argStrings) IN() Argument {
	a.opts = argOptionIN
	return a
}
func (a argStrings) options() uint { return a.opts }

// ArgString adds a string or a slice of strings to the argument list.
// Providing no arguments returns a NULL type.
// All arguments mut be a valid utf-8 string.
func ArgString(args ...string) Argument {
	if len(args) == 1 {
		return argString(args[0])
	}
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
func (a argBool) len() int      { return 1 }
func (a argBool) IN() Argument  { return a }
func (a argBool) options() uint { return 0 }

type argBools struct {
	opts uint
	data []bool
}

func (a argBools) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, v == true)
	}
}

func (a argBools) writeTo(w queryWriter, pos int) error {
	if a.options() == 0 {
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
	if a.options() == 0 {
		return len(a.data)
	}
	return 1
}

func (a argBools) IN() Argument {
	a.opts = argOptionIN
	return a
}
func (a argBools) options() uint { return a.opts }

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
func (a argInt) len() int      { return 1 }
func (a argInt) IN() Argument  { return a }
func (a argInt) options() uint { return 0 }

type argInts struct {
	opts uint
	data []int
}

func (a argInts) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, int64(v))
	}
}

func (a argInts) writeTo(w queryWriter, pos int) error {
	if a.options() == 0 {
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
	if a.options() == 0 {
		return len(a.data)
	}
	return 1
}

func (a argInts) IN() Argument {
	a.opts = argOptionIN
	return a
}

func (a argInts) options() uint { return a.opts }

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
func (a argInt64) len() int      { return 1 }
func (a argInt64) IN() Argument  { return a }
func (a argInt64) options() uint { return 0 }

type argInt64s struct {
	opts uint
	data []int64
}

func (a argInt64s) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, int64(v))
	}
}

func (a argInt64s) writeTo(w queryWriter, pos int) error {
	if a.options() == 0 {
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
	if a.options() == 0 {
		return len(a.data)
	}
	return 1
}

func (a argInt64s) IN() Argument {
	a.opts = argOptionIN
	return a
}

func (a argInt64s) options() uint { return a.opts }

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
func (a argFloat64) len() int      { return 1 }
func (a argFloat64) IN() Argument  { return a }
func (a argFloat64) options() uint { return 0 }

type argFloat64s struct {
	opts uint
	data []float64
}

func (a argFloat64s) toIFace(args *[]interface{}) {
	for _, v := range a.data {
		*args = append(*args, float64(v))
	}
}

func (a argFloat64s) writeTo(w queryWriter, pos int) error {
	if a.options() == 0 {
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
	if a.options() == 0 {
		return len(a.data)
	}
	return 1
}

func (a argFloat64s) IN() Argument {
	a.opts = argOptionIN
	return a
}

func (a argFloat64s) options() uint { return a.opts }

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

func (e *expr) writeTo(w queryWriter, _ int) error { return nil }
func (e *expr) len() int                           { return 1 }
func (e *expr) IN() Argument                       { return e }
func (e *expr) options() uint                      { return 0 }

// Repeat takes a SQL string and repeats the question marks with the provided
// arguments. If the amount of arguments does not match the number of questions
// marks, a Mismatch error gets returned. The arguments are getting converted to
// an interface slice to easy passing into the db.Query/db.Exec/etc functions at
// an argument.
//		Repeat("SELECT * FROM table WHERE id IN (?) AND status IN (?)", ArgInt(myIntSlice...), ArgString(myStrSlice...))
// Gets converted to:
//		SELECT * FROM table WHERE id IN (?,?) AND status IN (?,?,?)
// The questions marks are of course depending on the values in the Arg*
// functions. This function should be generally used when dealing with prepared
// statements.
func Repeat(sql string, args ...Argument) (string, []interface{}, error) {
	const qMarkStr = `?`
	const qMarkRne = '?'

	markCount := strings.Count(sql, qMarkStr)
	if want := len(args); markCount != want || want == 0 {
		return "", nil, errors.NewMismatchf("[dbr] Repeat: Number of %s:%d do not match the number of repetitions: %d", qMarkStr, markCount, want)
	}

	retArgs := make([]interface{}, 0, len(args)*2)

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	n := markCount
	i := 0
	for i < n {
		m := strings.Index(sql, qMarkStr)
		if m < 0 {
			break
		}
		buf.WriteString(sql[:m])

		if i < len(args) {
			prevLen := len(retArgs)
			args[i].toIFace(&retArgs)
			reps := len(retArgs) - prevLen
			for r := 0; r < reps; r++ {
				buf.WriteByte(qMarkRne)
				if r < reps-1 {
					buf.WriteByte(',')
				}
			}
		}
		sql = sql[m+len(qMarkStr):]
		i++
	}
	buf.WriteString(sql)
	return buf.String(), retArgs, nil
}

// Preprocess takes an SQL string with placeholders and a list of arguments to
// replace them with. It returns a blank string and error if the number of placeholders
// does not match the number of arguments.
func Preprocess(sql string, args ...Argument) (string, error) {
	// Get the number of arguments to add to this query
	if sql == "" {
		if len(args) != 0 {
			return "", errors.NewNotValidf("[dbr] Arguments are imbalanced")
		}
		return "", nil
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	qCountTotal := 0
	qCount := -1
	argIndex := 0
	argLength := 0
	if len(args) > 0 {
		argLength = args[0].len()
	}
	pos := 0
	for pos < len(sql) {
		r, w := utf8.DecodeRuneInString(sql[pos:])
		pos += w

		switch {
		case r == '?':
			if qCount >= argLength {
				return "", errors.NewNotValidf("[dbr] Arguments are imbalanced. Placeholder count: %d Current argument count: %d", qCount, args[argIndex].len())
			}

			if qCount < argLength-1 {
				qCount++
			} else {
				qCount = 0 // next argument set starts
				argIndex++
				if argIndex >= len(args) {
					return "", errors.NewNotValidf("[dbr] Arguments are imbalanced. Argument Index %d but argument count was %d", argIndex, len(args)-1)
				}
				argLength = args[argIndex].len()
			}
			if argLength == 0 {
				return "", errors.NewEmptyf("[dbr] Empty Argument for position %d", qCountTotal+1)
			}

			if err := args[argIndex].writeTo(buf, qCount); err != nil {
				return "", errors.Wrap(err, "[dbr] Preprocess writeTo arguments")
			}

			qCountTotal++
		case r == '`', r == '\'', r == '"':
			p := strings.IndexRune(sql[pos:], r)
			if p == -1 {
				return "", errors.NewNotValidf("[dbr] Preprocess: Invalid syntax")
			}
			if r == '"' {
				r = '\''
			}
			buf.WriteRune(r)
			buf.WriteString(sql[pos : pos+p])
			buf.WriteRune(r)
			pos += p + 1
		case r == '[':
			w := strings.IndexRune(sql[pos:], ']')
			col := sql[pos : pos+w]
			dialect.EscapeIdent(buf, col)
			pos += w + 1 // size of ']'
		default:
			buf.WriteRune(r)
		}
	}

	if al := Arguments(args).len(); qCountTotal != al {
		return "", errors.NewNotValidf("[dbr] Arguments are imbalanced. Placeholders: %d Current argument count: %d or %d", qCountTotal, al, len(args))
	}
	return buf.String(), nil
}
