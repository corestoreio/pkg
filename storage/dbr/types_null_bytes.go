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

	"github.com/corestoreio/csfw/util/null/convert"
)

// NullBytes is a nullable byte slice. JSON marshals to zero if null. Considered
// null to SQL if zero. NullBytes implements interface Argument.
type NullBytes struct {
	opt   byte
	Bytes []byte
	Valid bool
}

func (a NullBytes) toIFace(args *[]interface{}) {
	if a.Valid {
		*args = append(*args, a.Bytes)
	} else {
		*args = append(*args, nil)
	}
}

func (a NullBytes) writeTo(w queryWriter, _ int) error {
	if a.Valid {
		dialect.EscapeBinary(w, a.Bytes)
		return nil
	}
	_, err := w.WriteString("NULL")
	return err
}

func (a NullBytes) len() int { return 1 }
func (a NullBytes) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a NullBytes) operator() byte { return a.opt }

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
func (b NullBytes) GoString() string {
	if !b.Valid {
		return "dbr.NullBytes{}"
	}
	return fmt.Sprintf("dbr.MakeNullBytes(%#v)", b.Bytes)
}

// UnmarshalJSON implements json.Unmarshaler.
// If data is len 0 or nil, it will unmarshal to JSON null.
// If not, it will copy your data slice into NullBytes.
func (b *NullBytes) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		b.Bytes = []byte("null")
	} else {
		b.Bytes = append(b.Bytes[0:0], data...)
	}

	b.Valid = true

	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to nil if the text is nil or len 0.
func (b *NullBytes) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		b.Bytes = nil
		b.Valid = false
	} else {
		b.Bytes = append(b.Bytes[0:0], text...)
		b.Valid = true
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if the NullBytes is nil.
func (b NullBytes) MarshalJSON() ([]byte, error) {
	if len(b.Bytes) == 0 || b.Bytes == nil {
		return []byte("null"), nil
	}
	return b.Bytes, nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode nil if the NullBytes is invalid.
func (b NullBytes) MarshalText() ([]byte, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytes, nil
}

// SetValid changes this NullBytes's value and also sets it to be non-null.
func (b *NullBytes) SetValid(n []byte) {
	b.Bytes = n
	b.Valid = true
}

// Ptr returns a pointer to this NullBytes's value, or a nil pointer if this NullBytes is null.
func (b NullBytes) Ptr() *[]byte {
	if !b.Valid {
		return nil
	}
	return &b.Bytes
}

// IsZero returns true for null or zero NullBytes's, for future omitempty support (Go 1.4?)
func (b NullBytes) IsZero() bool {
	return !b.Valid
}

// Scan implements the Scanner interface.
func (n *NullBytes) Scan(value interface{}) error {
	if value == nil {
		n.Bytes, n.Valid = []byte{}, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Bytes, value)
}

// Value implements the driver Valuer interface.
func (n NullBytes) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bytes, nil
}
