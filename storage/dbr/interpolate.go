package dbr

import (
	// "fmt"
	"bytes"
	"database/sql/driver"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

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
			return "", ErrArgumentMismatch
		}
		return "", nil
	}

	curVal := 0
	buf := new(bytes.Buffer)

	pos := 0
	for pos < len(sql) {
		r, w := utf8.DecodeRuneInString(sql[pos:])
		pos += w

		switch {
		case r == '?':
			if curVal >= len(vals) {
				return "", ErrArgumentMismatch
			}
			if err := interpolate(buf, vals[curVal]); err != nil {
				return "", err
			}
			curVal++
		case r == '`', r == '\'', r == '"':
			p := strings.IndexRune(sql[pos:], r)
			if p == -1 {
				return "", ErrInvalidSyntax
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
			D.EscapeIdent(buf, col)
			pos += w + 1 // size of ']'
		default:
			buf.WriteRune(r)
		}
	}

	if curVal != len(vals) {
		return "", ErrArgumentMismatch
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
			return ErrNotUTF8
		}
		D.EscapeString(w, str)
	case isFloat(kindOfV):
		var fval = valueOfV.Float()

		w.WriteString(strconv.FormatFloat(fval, 'f', -1, 64))
	case kindOfV == reflect.Bool:
		D.EscapeBool(w, valueOfV.Bool())
	case kindOfV == reflect.Struct:
		if typeOfV := valueOfV.Type(); typeOfV == typeOfTime {
			t := valueOfV.Interface().(time.Time)
			D.EscapeTime(w, t)
		} else {
			return ErrInvalidValue
		}
	case kindOfV == reflect.Slice:
		typeOfV := reflect.TypeOf(v)
		subtype := typeOfV.Elem()
		kindOfSubtype := subtype.Kind()

		sliceLen := valueOfV.Len()
		stringSlice := make([]string, 0, sliceLen)

		switch {
		case sliceLen == 0:
			return ErrInvalidSliceLength
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
					return ErrNotUTF8
				}
				buf := new(bytes.Buffer)
				D.EscapeString(buf, str)
				stringSlice = append(stringSlice, buf.String())
			}
		default:
			return ErrInvalidSliceValue
		}
		w.WriteRune('(')
		w.WriteString(strings.Join(stringSlice, ","))
		w.WriteRune(')')
	default:
		return ErrInvalidValue
	}
	return nil
}
