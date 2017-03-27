package dbr

import (
	"database/sql/driver"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

type fieldType uint8

// Type* constants define all available types which a field can contain.
const (
	typeBool fieldType = iota + 1
	typeBools
	typeInt64
	typeInt64s
	typeFloat64
	typeFloat64s
	typeString
	typeStrings
	typeByte
	typeTime
	typeInterfaces
	typeFields
)

// Argument transforms your value or values into an interface slice. This
// interface slice gets used in the database query functions as an argument. The
// underlying type in the interface must be one of driver.Value allowed types.
type Argument interface {
	appendTo(*[]interface{})
}

type arg struct {
	// fieldType specifies the used type. If 0 this struct is empty
	fieldType

	int64
	int64s []int64
	float64
	float64s []float64
	bool
	bools []bool
	string
	strings []string
	bytes   []byte
	time    time.Time
	// already converted
	ifaces []interface{}
}

func (a arg) appendTo(args *[]interface{}) {

	switch a.fieldType {
	case typeBool:
		*args = append(*args, a.bool)
	case typeBools:
		for _, v := range a.bools {
			*args = append(*args, v)
		}
	case typeInt64:
		*args = append(*args, a.int64)
	case typeInt64s:
		for _, v := range a.int64s {
			*args = append(*args, v)
		}
	case typeFloat64:
		*args = append(*args, a.float64)
	case typeFloat64s:
		for _, v := range a.float64s {
			*args = append(*args, v)
		}
	case typeString:
		*args = append(*args, a.string)
	case typeStrings:
		for _, v := range a.strings {
			*args = append(*args, v)
		}
	case typeByte:
		*args = append(*args, a.bytes)
	case typeTime:
		*args = append(*args, a.time)
	case typeInterfaces:
		*args = append(*args, a.ifaces...)
	}
}

// ArgString adds a string or a slice of strings to the argument list.
func ArgString(args ...string) Argument {
	if len(args) == 1 {
		return arg{
			fieldType: typeString,
			string:    args[0],
		}
	}
	return arg{
		fieldType: typeStrings,
		strings:   args,
	}
}

// ArgInt adds an integer or a slice of integers to the argument list.
func ArgInt(args ...int) Argument {
	if len(args) == 1 {
		return arg{
			fieldType: typeInt64,
			int64:     int64(args[0]),
		}
	}
	args2 := make([]interface{}, len(args))
	for i, a := range args {
		args2[i] = int64(a)
	}
	return arg{
		fieldType: typeInterfaces,
		ifaces:    args2,
	}
}

// Repeat takes a SQL string and repeats the question marks with the provided
// arguments. If the amount of arguments does not match the number of questions
// marks, a Mismatch error gets returned. The arguments are getting converted to
// an interface slice to easy passing into the db.Query/db.Exec/etc functions as
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
			args[i].appendTo(&retArgs)
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

func isUint(k reflect.Kind) bool {
	return (k == reflect.Uint) ||
		(k == reflect.Uint8) ||
		(k == reflect.Uint16) ||
		(k == reflect.Uint32) ||
		(k == reflect.Uint64)
}

func isInt(k reflect.Kind) bool {
	return (k == reflect.Int) ||
		(k == reflect.Int8) ||
		(k == reflect.Int16) ||
		(k == reflect.Int32) ||
		(k == reflect.Int64)
}

func isFloat(k reflect.Kind) bool {
	return (k == reflect.Float32) ||
		(k == reflect.Float64)
}

// sql is like "id = ? OR username = ?"
// vals is like []interface{}{4, "bob"}
// NOTE that vals can only have values of certain types:
//   - Integers (signed and unsigned)
//   - floats
//   - strings (that are valid utf-8)
//   - booleans
//   - times
var typeOfTime = reflect.TypeOf(time.Time{})

// Preprocess takes an SQL string with placeholders and a list of arguments to
// replace them with. It returns a blank string and error if the number of placeholders
// does not match the number of arguments.
func Preprocess(sql string, vals []interface{}) (string, error) {
	// Get the number of arguments to add to this query
	if sql == "" {
		if len(vals) != 0 {
			return "", errors.NewNotValidf(errArgMismatch)
		}
		return "", nil
	}

	curVal := 0
	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	pos := 0
	for pos < len(sql) {
		r, w := utf8.DecodeRuneInString(sql[pos:])
		pos += w

		switch {
		case r == '?':
			if curVal >= len(vals) {
				return "", errors.NewNotValidf(errArgMismatch)
			}
			if err := interpolate(buf, vals[curVal]); err != nil {
				return "", err
			}
			curVal++
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

	if curVal != len(vals) {
		return "", errors.NewNotValidf(errArgMismatch)
	}
	return buf.String(), nil
}

func interpolate(w QueryWriter, v interface{}) error {
	valuer, ok := v.(driver.Valuer)
	if ok {
		val, err := valuer.Value()
		if err != nil {
			return err
		}
		v = val
	}

	valueOfV := reflect.ValueOf(v)
	kindOfV := valueOfV.Kind()

	switch {
	case v == nil:
		w.WriteString("NULL")
	case isInt(kindOfV):
		var ival = valueOfV.Int()

		w.WriteString(strconv.FormatInt(ival, 10))
	case isUint(kindOfV):
		var uival = valueOfV.Uint()

		w.WriteString(strconv.FormatUint(uival, 10))
	case kindOfV == reflect.String:
		var str = valueOfV.String()

		if !utf8.ValidString(str) {
			return errors.NewNotValidf(errNotUTF8)
		}
		dialect.EscapeString(w, str)
	case isFloat(kindOfV):
		var fval = valueOfV.Float()

		w.WriteString(strconv.FormatFloat(fval, 'f', -1, 64))
	case kindOfV == reflect.Bool:
		dialect.EscapeBool(w, valueOfV.Bool())
	case kindOfV == reflect.Struct:
		if typeOfV := valueOfV.Type(); typeOfV == typeOfTime {
			t := valueOfV.Interface().(time.Time)
			dialect.EscapeTime(w, t)
		} else {
			return errors.NewNotValidf("[dbr] Interpolate: Invalid value for time")
		}
	case kindOfV == reflect.Slice:
		typeOfV := reflect.TypeOf(v)
		subtype := typeOfV.Elem()
		kindOfSubtype := subtype.Kind()

		sliceLen := valueOfV.Len()
		stringSlice := make([]string, 0, sliceLen)

		switch {
		case sliceLen == 0:
			return errors.NewNotValidf("[dbr] Interpolate: Invalid slice length")
		case isInt(kindOfSubtype):
			for i := 0; i < sliceLen; i++ {
				var ival = valueOfV.Index(i).Int()
				stringSlice = append(stringSlice, strconv.FormatInt(ival, 10))
			}
		case isUint(kindOfSubtype):
			for i := 0; i < sliceLen; i++ {
				var uival = valueOfV.Index(i).Uint()
				stringSlice = append(stringSlice, strconv.FormatUint(uival, 10))
			}
		case kindOfSubtype == reflect.String:
			for i := 0; i < sliceLen; i++ {
				var str = valueOfV.Index(i).String()
				if !utf8.ValidString(str) {
					return errors.NewNotValidf(errNotUTF8)
				}
				var buf = bufferpool.Get()
				dialect.EscapeString(buf, str)
				stringSlice = append(stringSlice, buf.String())
				bufferpool.Put(buf)
			}
		default:
			return errors.NewNotValidf("[dbr] Interpolate: Invalid slice value")
		}
		w.WriteRune('(')
		w.WriteString(strings.Join(stringSlice, ","))
		w.WriteRune(')')
	default:
		return errors.NewNotValidf("[dbr] Interpolate: Invalid value")
	}
	return nil
}
