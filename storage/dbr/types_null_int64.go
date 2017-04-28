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
	"database/sql"
	"strconv"

	"github.com/corestoreio/errors"
)

// NullInt64 is a nullable int64. It does not consider zero values to be null.
// It will decode to null, not zero, if null. NullInt64 implements interface
// Argument.
type NullInt64 struct {
	sql.NullInt64
	op rune
}

func (a NullInt64) toIFace(args []interface{}) []interface{} {
	if a.NullInt64.Valid {
		return append(args, a.NullInt64.Int64)
	}
	return append(args, nil)
}

func (a NullInt64) writeTo(w queryWriter, _ int) error {
	if a.NullInt64.Valid {
		_, err := w.WriteString(strconv.FormatInt(a.NullInt64.Int64, 10))
		return err
	}
	_, err := w.WriteString(sqlStrNull)
	return err
}

func (a NullInt64) len() int { return 1 }

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a NullInt64) Operator(op rune) Argument {
	a.op = op
	return a
}

func (a NullInt64) operator() rune { return a.op }

// MakeNullInt64 creates a new NullInt64. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. NullInt64
// implements interface Argument.
func MakeNullInt64(i int64, valid ...bool) NullInt64 {
	v := true
	if len(valid) == 1 {
		v = valid[0]
	}
	return NullInt64{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: v,
		},
	}
}

// GoString prints an optimized Go representation. Takes are of backticks.
func (a NullInt64) GoString() string {
	if !a.Valid {
		return "dbr.NullInt64{}"
	}
	return "dbr.MakeNullInt64(" + strconv.FormatInt(a.Int64, 10) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null NullInt64. It also supports unmarshalling a
// sql.NullInt64.
func (a *NullInt64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = JSONUnMarshalFn(data, &a.Int64)
	case map[string]interface{}:
		dto := &struct {
			NullInt64 int64
			Valid     bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		a.Int64 = dto.NullInt64
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NewNotValidf("[null] json: cannot unmarshal %#v into Go value of type null.NullInt64", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null NullInt64 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (a *NullInt64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		a.Valid = false
		return nil
	}
	var err error
	a.Int64, err = strconv.ParseInt(string(text), 10, 64)
	a.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullInt64 is null.
func (a NullInt64) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return []byte("null"), nil
	}
	return strconv.AppendInt([]byte{}, a.Int64, 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this NullInt64 is null.
func (a NullInt64) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendInt([]byte{}, a.Int64, 10), nil
}

// SetValid changes this NullInt64's value and also sets it to be non-null.
func (a *NullInt64) SetValid(n int64) {
	a.Int64 = n
	a.Valid = true
}

// Ptr returns a pointer to this NullInt64's value, or a nil pointer if this NullInt64 is null.
func (a NullInt64) Ptr() *int64 {
	if !a.Valid {
		return nil
	}
	return &a.Int64
}

// IsZero returns true for invalid NullInt64's, for future omitempty support (Go 1.4?)
// A non-null NullInt64 with a 0 value will not be considered zero.
func (a NullInt64) IsZero() bool {
	return !a.Valid
}

type argNullInt64s struct {
	op   rune
	data []NullInt64
}

func (a argNullInt64s) toIFace(args []interface{}) []interface{} {
	for _, s := range a.data {
		if s.Valid {
			args = append(args, s.Int64)
		} else {
			args = append(args, nil)
		}
	}
	return args
}

func (a argNullInt64s) writeTo(w queryWriter, pos int) error {
	if a.operator() != In && a.operator() != NotIn {
		if s := a.data[pos]; s.Valid {
			_, err := w.WriteString(strconv.FormatInt(s.Int64, 10))
			return err
		}
		_, err := w.WriteString(sqlStrNull)
		return err
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		if v.Valid {
			w.WriteString(strconv.FormatInt(v.Int64, 10))
		} else {
			w.WriteString(sqlStrNull)
		}
		if i < l {
			w.WriteRune(',')
		}
	}
	_, err := w.WriteRune(')')
	return err
}

func (a argNullInt64s) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

// Operator sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Operator*.
func (a argNullInt64s) Operator(op rune) Argument {
	a.op = op
	return a
}

func (a argNullInt64s) operator() rune { return a.op }

// ArgNullInt64 adds a nullable int64 or a slice of nullable int64s to the
// argument list. Providing no arguments returns a NULL type.
func ArgNullInt64(args ...NullInt64) Argument {
	if len(args) == 1 {
		return args[0]
	}
	return argNullInt64s{data: args}
}
