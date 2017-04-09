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

	"github.com/corestoreio/errors"
)

// NullBool is a nullable bool. It does not consider false values to be null. It
// will decode to null, not false, if null. NullBool implements interface
// Argument.
type NullBool struct {
	opt byte
	sql.NullBool
}

func (a NullBool) toIFace(args *[]interface{}) {
	if a.NullBool.Valid {
		*args = append(*args, a.NullBool.Bool)
	} else {
		*args = append(*args, nil)
	}
}

func (a NullBool) writeTo(w queryWriter, _ int) error {
	if a.NullBool.Valid {
		dialect.EscapeBool(w, a.Bool)
		return nil
	}
	_, err := w.WriteString("NULL")
	return err
}

func (a NullBool) len() int { return 1 }
func (a NullBool) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a NullBool) operator() byte { return a.opt }

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
func (b *NullBool) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		b.Bool = x
	case map[string]interface{}:
		dto := &struct {
			NullBool bool
			Valid    bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		b.Bool = dto.NullBool
		b.Valid = dto.Valid
	case nil:
		b.Valid = false
		return nil
	default:
		err = errors.NewNotValidf("[dbr] json: cannot unmarshal %#v into Go value of type null.NullBool", v)
	}
	b.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler. It will unmarshal to a
// null NullBool if the input is a blank or not an integer. It will return an
// error if the input is not an integer, blank, or "null".
func (b *NullBool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		b.Valid = false
		return nil
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		b.Valid = false
		return errors.NewNotValidf("[dbr] NullBool invalid input: %q", str)
	}
	b.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullBool is null.
func (b NullBool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this NullBool is null.
func (b NullBool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// SetValid changes this NullBool's value and also sets it to be non-null.
func (b *NullBool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

// Ptr returns a pointer to this NullBool's value, or a nil pointer if this
// NullBool is null.
func (b NullBool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// IsZero returns true for invalid Bools, for future omitempty support (Go 1.4?)
// A non-null NullBool with a 0 value will not be considered zero.
func (b NullBool) IsZero() bool {
	return !b.Valid
}
