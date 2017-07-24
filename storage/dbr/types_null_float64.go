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

// NullFloat64 is a nullable float64. It does not consider zero values to be null.
// It will decode to null, not zero, if null. NullFloat64 implements interface
// Argument.
type NullFloat64 struct {
	sql.NullFloat64
}

func (a NullFloat64) toIFace(args []interface{}) []interface{} {
	if a.NullFloat64.Valid {
		return append(args, a.NullFloat64.Float64)
	}
	return append(args, nil)
}

func (a NullFloat64) writeTo(w queryWriter, _ int) error {
	if a.NullFloat64.Valid {
		return writeFloat64(w, a.NullFloat64.Float64)
	}
	_, err := w.WriteString(sqlStrNull)
	return err
}

func (a NullFloat64) len() int { return 1 }

// MakeNullFloat64 creates a new NullFloat64. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. NullFloat64
// implements interface Argument.
func MakeNullFloat64(f float64, valid ...bool) NullFloat64 {
	v := true
	if len(valid) == 1 {
		v = valid[0]
	}
	return NullFloat64{
		NullFloat64: sql.NullFloat64{
			Float64: f,
			Valid:   v,
		},
	}
}

// GoString prints an optimized Go representation.
func (a NullFloat64) GoString() string {
	if !a.Valid {
		return "dbr.NullFloat64{}"
	}
	return "dbr.MakeNullFloat64(" + strconv.FormatFloat(a.Float64, 'f', -1, 64) + ")"
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null NullFloat64.
// It also supports unmarshalling a sql.NullFloat64.
func (a *NullFloat64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		a.Float64 = x
	case map[string]interface{}:
		dto := &struct {
			NullFloat64 float64
			Valid       bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		a.Float64 = dto.NullFloat64
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NewNotValidf("[dbr] json: cannot unmarshal %#v into Go value of type null.NullFloat64", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null NullFloat64 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (a *NullFloat64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		a.Valid = false
		return nil
	}
	var err error
	a.Float64, err = strconv.ParseFloat(string(text), 64)
	a.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullFloat64 is null.
func (a NullFloat64) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return []byte("null"), nil
	}
	return strconv.AppendFloat([]byte{}, a.Float64, 'f', -1, 64), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this NullFloat64 is null.
func (a NullFloat64) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendFloat([]byte{}, a.Float64, 'f', -1, 64), nil
}

// SetValid changes this NullFloat64's value and also sets it to be non-null.
func (a *NullFloat64) SetValid(n float64) {
	a.Float64 = n
	a.Valid = true
}

// Ptr returns a pointer to this NullFloat64's value, or a nil pointer if this NullFloat64 is null.
func (a NullFloat64) Ptr() *float64 {
	if !a.Valid {
		return nil
	}
	return &a.Float64
}

// IsZero returns true for invalid Float64s, for future omitempty support (Go 1.4?)
// A non-null NullFloat64 with a 0 value will not be considered zero.
func (a NullFloat64) IsZero() bool {
	return !a.Valid
}

// NullFloat64s adds a nullable float64 or a slice of nullable float64s to the
// argument list. Providing no arguments returns a NULL type.
type NullFloat64s []NullFloat64

func (a NullFloat64s) toIFace(args []interface{}) []interface{} {
	for _, s := range a {
		if s.Valid {
			args = append(args, s.Float64)
		} else {
			args = append(args, nil)
		}
	}
	return args
}

func (a NullFloat64s) writeTo(w queryWriter, pos int) error {
	if s := a[pos]; s.Valid {
		return writeFloat64(w, s.Float64)
	}
	_, err := w.WriteString(sqlStrNull)
	return err
}

func (a NullFloat64s) len() int {
	return len(a)
}
