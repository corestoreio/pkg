// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml

import (
	"bytes"
	"database/sql"
	"strings"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// NullString is a nullable string. It supports SQL and JSON serialization.
// It will marshal to null if null. Blank string input will be considered null.
// NullString implements interface Argument.
type NullString struct {
	sql.NullString
}

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
		return "dml.NullString{}"
	}
	return "dml.MakeNullString(`" + a.String + "`)"
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
		err = errors.NewNotValidf("[dml] json: cannot unmarshal %#v into Go value of type dml.NullString", v)
	}
	a.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullString is dml.
func (a NullString) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return sqlBytesNullLC, nil
	}
	return JSONMarshalFn(a.String)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this NullString is dml.
func (a NullString) MarshalText() ([]byte, error) {
	if !a.Valid {
		return nil, nil
	}
	return []byte(a.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null NullString if the input is a blank string.
func (a *NullString) UnmarshalText(text []byte) error {
	if !utf8.Valid(text) {
		return errors.NewNotValidf("[dml] Input bytes are not valid UTF-8 encoded.")
	}
	a.String = string(text)
	a.Valid = a.String != ""
	return nil
}

// SetValid changes this NullString's value and also sets it to be non-dml.
func (a *NullString) SetValid(v string) {
	a.String = v
	a.Valid = true
}

// Ptr returns a pointer to this NullString's value, or a nil pointer if this NullString is dml.
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

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (a NullString) GobEncode() ([]byte, error) {
	return a.Marshal()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (a *NullString) GobDecode(data []byte) error {
	return a.Unmarshal(data)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *NullString) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a NullString) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// Marshal binary encoder for protocol buffers. Implements proto.Marshaler.
func (a NullString) Marshal() ([]byte, error) {
	return a.MarshalText()
}

// MarshalTo binary encoder for protocol buffers which writes into data.
func (a NullString) MarshalTo(data []byte) (n int, err error) {
	if !a.Valid {
		return 0, nil
	}
	n = copy(data, a.String)
	return
}

// Unmarshal binary decoder for protocol buffers. Implements proto.Unmarshaler.
func (a *NullString) Unmarshal(data []byte) error {
	return a.UnmarshalText(data)
}

// Size returns the size of the underlying type. If not valid, the size will be
// 0. Implements proto.Sizer.
func (a NullString) Size() (s int) {
	return len(a.String)
}

func (a NullString) writeTo(w *bytes.Buffer) (err error) {
	if a.Valid {
		if utf8.ValidString(a.String) {
			dialect.EscapeString(w, a.String)
		} else {
			err = errors.NewNotValidf("[dml] NullString.writeTo: String is not UTF-8: %q", a.String)
		}
	} else {
		_, err = w.WriteString(sqlStrNullUC)
	}
	return
}

func (a NullString) append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.String)
	}
	return append(args, nil)
}
