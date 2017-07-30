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
	"bytes"
	"database/sql"

	"github.com/corestoreio/errors"
)

// NullBool is a nullable bool. It does not consider false values to be null. It
// will decode to null, not false, if null. NullBool implements interface
// Argument.
type NullBool struct {
	sql.NullBool
}

func (a NullBool) toIFace(args []interface{}) []interface{} {
	if a.NullBool.Valid {
		return append(args, a.NullBool.Bool)
	}
	return append(args, nil)
}

func (a NullBool) writeTo(w *bytes.Buffer, _ int) error {
	if a.NullBool.Valid {
		dialect.EscapeBool(w, a.Bool)
		return nil
	}
	_, err := w.WriteString(sqlStrNull)
	return err
}

func (a NullBool) len() int { return 1 }

// MakeNullBool creates a new NullBool. Implements interface Argument.
func MakeNullBool(b bool, valid ...bool) NullBool {
	v := true
	if len(valid) == 1 {
		v = valid[0]
	}
	return NullBool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: v,
		},
	}
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null NullBool. It also supports unmarshalling a
// sql.NullBool.
func (a *NullBool) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		a.Bool = x
	case map[string]interface{}:
		dto := &struct {
			NullBool bool
			Valid    bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		a.Bool = dto.NullBool
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NewNotValidf("[dbr] json: cannot unmarshal %#v into Go value of type null.NullBool", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler. It will unmarshal to a
// null NullBool if the input is a blank or not an integer. It will return an
// error if the input is not an integer, blank, or "null".
func (a *NullBool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		a.Valid = false
		return nil
	case "true":
		a.Bool = true
	case "false":
		a.Bool = false
	default:
		a.Valid = false
		return errors.NewNotValidf("[dbr] NullBool invalid input: %q", str)
	}
	a.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullBool is null.
func (a NullBool) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return []byte("null"), nil
	}
	if !a.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this NullBool is null.
func (a NullBool) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	if !a.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// SetValid changes this NullBool's value and also sets it to be non-null.
func (a *NullBool) SetValid(v bool) {
	a.Bool = v
	a.Valid = true
}

// Ptr returns a pointer to this NullBool's value, or a nil pointer if this
// NullBool is null.
func (a NullBool) Ptr() *bool {
	if !a.Valid {
		return nil
	}
	return &a.Bool
}

// IsZero returns true for invalid Bools, for future omitempty support (Go 1.4?)
// A non-null NullBool with a 0 value will not be considered zero.
func (a NullBool) IsZero() bool {
	return !a.Valid
}
