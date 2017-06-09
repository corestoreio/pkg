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
	"strings"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

// NullString is a nullable string. It supports SQL and JSON serialization.
// It will marshal to null if null. Blank string input will be considered null.
// NullString implements interface Argument.
type NullString struct {
	sql.NullString
	op Op
}

func (a NullString) toIFace(args []interface{}) []interface{} {
	if a.NullString.Valid {
		return append(args, a.NullString.String)
	}
	return append(args, nil)
}

func (a NullString) writeTo(w queryWriter, _ int) error {
	if a.NullString.Valid {
		if !utf8.ValidString(a.NullString.String) {
			return errors.NewNotValidf("[dbr] Argument.WriteTo: StringNull is not UTF-8: %q", a.NullString.String)
		}
		dialect.EscapeString(w, a.NullString.String)
	} else {
		w.WriteString(sqlStrNull)
	}
	return nil
}

func (a NullString) len() int { return 1 }

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a NullString) applyOperator(op Op) Argument {
	a.op = op
	return a
}

func (a NullString) operator() Op { return a.op }

// MakeNullString creates a new NullString. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. NullString
// implements interface Argument.
func MakeNullString(s string, valid ...bool) NullString {
	v := true
	if len(valid) == 1 {
		v = valid[0]
	}
	return NullString{
		NullString: sql.NullString{
			String: s,
			Valid:  v,
		},
	}
}

// GoString prints an optimized Go representation. Takes are of backticks.
// Looses the information of the private operator. That might get fixed.
func (a NullString) GoString() string {
	if a.Valid && strings.ContainsRune(a.String, '`') {
		// `This is my`string`
		a.String = strings.Join(strings.Split(a.String, "`"), "`+\"`\"+`")
		// `This is my`+"`"+`string`
	}
	if !a.Valid {
		return "dbr.NullString{}"
	}
	return "dbr.MakeNullString(`" + a.String + "`)"
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null NullString.
// It also supports unmarshalling a sql.NullString.
func (a *NullString) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}

	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		a.String = x
	case map[string]interface{}:
		dto := &struct {
			NullString string
			Valid      bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		a.String = dto.NullString
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NewNotValidf("[dbr] json: cannot unmarshal %#v into Go value of type dbr.NullString", v)
	}
	a.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullString is dbr.
func (a NullString) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return []byte("null"), nil
	}
	return JSONMarshalFn(a.String)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this NullString is dbr.
func (a NullString) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return []byte(a.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null NullString if the input is a blank string.
func (a *NullString) UnmarshalText(text []byte) error {
	if !utf8.Valid(text) {
		return errors.NewNotValidf("[dbr] Input bytes are not valid UTF-8 encoded.")
	}
	a.String = string(text)
	a.Valid = a.String != ""
	return nil
}

// SetValid changes this NullString's value and also sets it to be non-dbr.
func (a *NullString) SetValid(v string) {
	a.String = v
	a.Valid = true
}

// Ptr returns a pointer to this NullString's value, or a nil pointer if this NullString is dbr.
func (a NullString) Ptr() *string {
	if !a.Valid {
		return nil
	}
	return &a.String
}

// IsZero returns true for null strings, for potential future omitempty support.
func (a NullString) IsZero() bool {
	return !a.Valid
}

type argNullStrings struct {
	op   Op
	data []NullString
}

func (a argNullStrings) toIFace(args []interface{}) []interface{} {
	for _, s := range a.data {
		if s.Valid {
			args = append(args, s.String)
		} else {
			args = append(args, nil)
		}
	}
	return args
}

func (a argNullStrings) writeTo(w queryWriter, pos int) error {
	if s := a.data[pos]; s.Valid {
		if !utf8.ValidString(s.String) {
			return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", s.String)
		}
		dialect.EscapeString(w, s.String)
		return nil
	}
	_, err := w.WriteString(sqlStrNull)
	return err
}

func (a argNullStrings) len() int {
	return len(a.data)
}

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a argNullStrings) applyOperator(op Op) Argument {
	a.op = op
	return a
}

func (a argNullStrings) operator() Op { return a.op }

// ArgNullString adds a nullable string or a slice of nullable strings to the
// argument list. Providing no arguments returns a NULL type. All arguments must
// be a valid utf-8 string.
func ArgNullString(args ...NullString) Argument {
	if len(args) == 1 {
		return args[0]
	}
	return argNullStrings{data: args}
}
