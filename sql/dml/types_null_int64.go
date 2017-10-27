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

package dml

import (
	"database/sql"
	"database/sql/driver"
	"strconv"

	"github.com/corestoreio/errors"
)

// NullInt64 is a nullable int64. It does not consider zero values to be null.
// It will decode to null, not zero, if null. NullInt64 implements interface
// Argument.
type NullInt64 struct {
	sql.NullInt64
}

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
		return "dml.NullInt64{}"
	}
	return "dml.MakeNullInt64(" + strconv.FormatInt(a.Int64, 10) + ")"
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
// It will return an error if the input is not an integer, blank, or sqlStrNullLC.
func (a *NullInt64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == sqlStrNullLC {
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
		return []byte(sqlStrNullLC), nil
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

// Value implements the driver.Valuer interface.
func (a NullInt64) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Int64, nil
}
