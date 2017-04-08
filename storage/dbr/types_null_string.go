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
	opt byte
	sql.NullString
}

func (a NullString) toIFace(args *[]interface{}) {
	if a.NullString.Valid {
		*args = append(*args, a.NullString.String)
	} else {
		*args = append(*args, nil)
	}
}

func (a NullString) writeTo(w queryWriter, _ int) error {
	if a.NullString.Valid {
		if !utf8.ValidString(a.NullString.String) {
			return errors.NewNotValidf("[dbr] Argument.WriteTo: StringNull is not UTF-8: %q", a.NullString.String)
		}
		dialect.EscapeString(w, a.NullString.String)
	} else {
		w.WriteString("NULL")
	}
	return nil
}

func (a NullString) len() int { return 1 }
func (a NullString) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a NullString) operator() byte { return a.opt }

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
func (ns NullString) GoString() string {
	if ns.Valid && strings.ContainsRune(ns.String, '`') {
		// `This is my`string`
		ns.String = strings.Join(strings.Split(ns.String, "`"), "`+\"`\"+`")
		// `This is my`+"`"+`string`
	}
	if !ns.Valid {
		return "dbr.NullString{}"
	}
	return "dbr.MakeNullString(`" + ns.String + "`)"
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null NullString.
// It also supports unmarshalling a sql.NullString.
func (s *NullString) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}

	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		s.String = x
	case map[string]interface{}:
		dto := &struct {
			NullString string
			Valid      bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		s.String = dto.NullString
		s.Valid = dto.Valid
	case nil:
		s.Valid = false
		return nil
	default:
		err = errors.NewNotValidf("[null] json: cannot unmarshal %#v into Go value of type dbr.NullString", v)
	}
	s.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullString is dbr.
func (s NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return JSONMarshalFn(s.String)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this NullString is dbr.
func (s NullString) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return []byte(s.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null NullString if the input is a blank string.
func (s *NullString) UnmarshalText(text []byte) error {
	if !utf8.Valid(text) {
		return errors.NewNotValidf("[dbr] Input bytes are not valid UTF-8 encoded.")
	}
	s.String = string(text)
	s.Valid = s.String != ""
	return nil
}

// SetValid changes this NullString's value and also sets it to be non-dbr.
func (s *NullString) SetValid(v string) {
	s.String = v
	s.Valid = true
}

// Ptr returns a pointer to this NullString's value, or a nil pointer if this NullString is dbr.
func (s NullString) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero returns true for null strings, for potential future omitempty support.
func (s NullString) IsZero() bool {
	return !s.Valid
}

type argNullStrings struct {
	opt  byte
	data []NullString
}

func (a argNullStrings) toIFace(args *[]interface{}) {
	for _, s := range a.data {
		if s.Valid {
			*args = append(*args, s.String)
		} else {
			*args = append(*args, nil)
		}
	}
}

func (a argNullStrings) writeTo(w queryWriter, pos int) error {
	if a.operator() != OperatorIn && a.operator() != OperatorNotIn {
		if s := a.data[pos]; s.Valid {
			if !utf8.ValidString(s.String) {
				return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", s.String)
			}
			dialect.EscapeString(w, s.String)
			return nil
		}
		_, err := w.WriteString("NULL")
		return err
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		if v.Valid {
			if !utf8.ValidString(v.String) {
				return errors.NewNotValidf("[dbr] Argument.WriteTo: StringNull is not UTF-8: %q", v.String)
			}
			dialect.EscapeString(w, v.String)
			if i < l {
				w.WriteRune(',')
			}
		} else {
			w.WriteString("NULL")
		}
	}
	_, err := w.WriteRune(')')
	return err
}

func (a argNullStrings) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

func (a argNullStrings) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a argNullStrings) operator() byte { return a.opt }

// ArgNullString adds a nullable string or a slice of nullable strings to the
// argument list. Providing no arguments returns a NULL type. All arguments mut
// be a valid utf-8 string.
func ArgNullString(args ...NullString) Argument {
	if len(args) == 1 {
		return args[0]
	}
	return argNullStrings{data: args}
}
