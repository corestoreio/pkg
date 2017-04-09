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
	opt byte
	sql.NullInt64
}

func (a NullInt64) toIFace(args *[]interface{}) {
	if a.NullInt64.Valid {
		*args = append(*args, a.NullInt64.Int64)
	} else {
		*args = append(*args, nil)
	}
}

func (a NullInt64) writeTo(w queryWriter, _ int) error {
	if a.NullInt64.Valid {
		_, err := w.WriteString(strconv.FormatInt(a.NullInt64.Int64, 10))
		return err
	}
	_, err := w.WriteString("NULL")
	return err
}

func (a NullInt64) len() int { return 1 }
func (a NullInt64) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a NullInt64) operator() byte { return a.opt }

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
func (i *NullInt64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = JSONUnMarshalFn(data, &i.Int64)
	case map[string]interface{}:
		dto := &struct {
			NullInt64 int64
			Valid     bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		i.Int64 = dto.NullInt64
		i.Valid = dto.Valid
	case nil:
		i.Valid = false
		return nil
	default:
		err = errors.NewNotValidf("[null] json: cannot unmarshal %#v into Go value of type null.NullInt64", v)
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null NullInt64 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *NullInt64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullInt64 is null.
func (i NullInt64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return strconv.AppendInt([]byte{}, i.Int64, 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this NullInt64 is null.
func (i NullInt64) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return strconv.AppendInt([]byte{}, i.Int64, 10), nil
}

// SetValid changes this NullInt64's value and also sets it to be non-null.
func (i *NullInt64) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

// Ptr returns a pointer to this NullInt64's value, or a nil pointer if this NullInt64 is null.
func (i NullInt64) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsZero returns true for invalid NullInt64's, for future omitempty support (Go 1.4?)
// A non-null NullInt64 with a 0 value will not be considered zero.
func (i NullInt64) IsZero() bool {
	return !i.Valid
}

type argNullInt64s struct {
	opt  byte
	data []NullInt64
}

func (a argNullInt64s) toIFace(args *[]interface{}) {
	for _, s := range a.data {
		if s.Valid {
			*args = append(*args, s.Int64)
		} else {
			*args = append(*args, nil)
		}
	}
}

func (a argNullInt64s) writeTo(w queryWriter, pos int) error {
	if a.operator() != OperatorIn && a.operator() != OperatorNotIn {
		if s := a.data[pos]; s.Valid {
			_, err := w.WriteString(strconv.FormatInt(s.Int64, 10))
			return err
		}
		_, err := w.WriteString("NULL")
		return err
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		if v.Valid {
			w.WriteString(strconv.FormatInt(v.Int64, 10))
		} else {
			w.WriteString("NULL")
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

func (a argNullInt64s) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a argNullInt64s) operator() byte { return a.opt }

// ArgNullInt64 adds a nullable int64 or a slice of nullable int64s to the
// argument list. Providing no arguments returns a NULL type.
func ArgNullInt64(args ...NullInt64) Argument {
	if len(args) == 1 {
		return args[0]
	}
	return argNullInt64s{data: args}
}
