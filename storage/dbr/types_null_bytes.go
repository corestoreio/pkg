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
	"database/sql/driver"
	"fmt"

	"github.com/corestoreio/errors"
)

// NullBytes is a nullable byte slice. JSON marshals to zero if null. Considered
// null to SQL if zero. NullBytes implements interface Argument.
type NullBytes struct {
	Bytes []byte
	op    Op
	Valid bool
}

func (a NullBytes) toIFace(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Bytes)
	}
	return append(args, nil)
}

func (a NullBytes) writeTo(w queryWriter, _ int) error {
	if a.Valid {
		dialect.EscapeBinary(w, a.Bytes)
		return nil
	}
	_, err := w.WriteString(sqlStrNull)
	return err
}

func (a NullBytes) len() int { return 1 }

// Op sets the SQL operator (IN, =, LIKE, BETWEEN, ...). Please refer to
// the constants Op*.
func (a NullBytes) applyOperator(op Op) Argument {
	a.op = op
	return a
}

func (a NullBytes) operator() Op { return a.op }

// MakeNullBytes creates a new NullBytes. Implements interface Argument.
func MakeNullBytes(b []byte, valid ...bool) NullBytes {
	v := true
	if len(valid) == 1 {
		v = valid[0]
	} else {
		if b == nil {
			v = false
		}
	}
	return NullBytes{
		Bytes: b,
		Valid: v,
	}
}

// GoString prints an optimized Go representation.
func (a NullBytes) GoString() string {
	if !a.Valid {
		return "dbr.NullBytes{}"
	}
	return fmt.Sprintf("dbr.MakeNullBytes(%#v)", a.Bytes)
}

// UnmarshalJSON implements json.Unmarshaler.
// If data is len 0 or nil, it will unmarshal to JSON null.
// If not, it will copy your data slice into NullBytes.
func (a *NullBytes) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		a.Bytes = []byte("null")
	} else {
		a.Bytes = append(a.Bytes[0:0], data...)
	}

	a.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to nil if the text is nil or len 0.
func (a *NullBytes) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		a.Bytes = nil
		a.Valid = false
	} else {
		a.Bytes = append(a.Bytes[0:0], text...)
		a.Valid = true
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if the NullBytes is nil.
func (a NullBytes) MarshalJSON() ([]byte, error) {
	if len(a.Bytes) == 0 || a.Bytes == nil {
		return []byte("null"), nil
	}
	return a.Bytes, nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode nil if the NullBytes is invalid.
func (a NullBytes) MarshalText() ([]byte, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Bytes, nil
}

// SetValid changes this NullBytes's value and also sets it to be non-null.
func (a *NullBytes) SetValid(n []byte) {
	a.Bytes = n
	a.Valid = true
}

// Ptr returns a pointer to this NullBytes's value, or a nil pointer if this NullBytes is null.
func (a NullBytes) Ptr() *[]byte {
	if !a.Valid {
		return nil
	}
	return &a.Bytes
}

// IsZero returns true for null or zero NullBytes's, for future omitempty support (Go 1.4?)
func (a NullBytes) IsZero() bool {
	return !a.Valid
}

// Scan implements the Scanner interface.
func (a *NullBytes) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		a.Bytes = make([]byte, len(v))
		copy(a.Bytes, v)
	case string:
		a.Bytes = []byte(v)
	default:
		return errors.NewNotSupportedf("[dbr] NUllBytes.Scan Type %T not supported", value)
	}

	a.Valid = a.Bytes != nil

	return nil
}

// Value implements the driver Valuer interface.
func (a NullBytes) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Bytes, nil
}
