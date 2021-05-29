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

package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/byteconv"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// Bool is a nullable bool. It does not consider false values to be null. It
// will decode to null, not false, if null. Bool implements interface
// Argument.
type Bool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// MakeBool creates a new Bool. Implements interface Argument.
func MakeBool(b bool) Bool {
	return Bool{
		Bool:  b,
		Valid: true,
	}
}

// MakeBoolFromByte makes a new Bool from a (text) byte slice.
func MakeBoolFromByte(data []byte) (nv Bool, err error) {
	nv.Bool, nv.Valid, err = byteconv.ParseBool(data)
	return
}

// Scan implements the Scanner interface.
func (a *Bool) Scan(value interface{}) (err error) {
	if value == nil {
		a.Bool, a.Valid = false, false
		return
	}
	switch vt := value.(type) {
	case []byte:
		a.Bool, a.Valid, err = byteconv.ParseBool(vt)
	case bool:
		a.Bool = vt
	case string:
		a.Bool, err = strconv.ParseBool(vt)
	}
	a.Valid = err == nil
	return nil
}

// Value implements the driver Valuer interface.
func (a Bool) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Bool, nil
}

// GoString prints an optimized Go representation.
func (a Bool) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatBool(a.Bool)
}

// GoString prints an optimized Go representation.
func (a Bool) GoString() string {
	if !a.Valid {
		return "null.Bool{}"
	}
	return "null.MakeBool(" + strconv.FormatBool(a.Bool) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null Bool. It also supports unmarshalling a
// sql.Bool.
func (a *Bool) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(bTextNullLC, data) {
		a.Valid = false
		a.Bool = false
		return nil
	}
	if v, ok, err := byteconv.ParseBool(data); ok && err == nil {
		a.Valid = true
		a.Bool = v
		return nil
	}

	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		a.Bool = x
	case map[string]interface{}:
		dto := &struct {
			Bool  bool
			Valid bool
		}{}
		err = json.Unmarshal(data, dto)
		a.Bool = dto.Bool
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[dml] json: cannot unmarshal %#v into Go value of type null.Bool", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler. It will unmarshal to a
// null Bool if the input is a blank or not an integer. It will return an
// error if the input is not an integer, blank, or "null".
func (a *Bool) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 || bytes.Equal(text, bTextNullUC) || bytes.Equal(text, bTextNullLC) {
		a.Valid = false
		a.Bool = false
		return nil
	}
	a.Bool, a.Valid, err = byteconv.ParseBool(text)
	return
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Bool is null.
func (a Bool) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	if !a.Bool {
		return bTextFalseLC, nil
	}
	return bTextTrueLC, nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Bool is null.
func (a Bool) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	if !a.Bool {
		return bTextFalseLC, nil
	}
	return bTextTrueLC, nil
}

// SetValid changes this Bool's value and also sets it to be non-null.
func (a *Bool) SetValid(v bool) { a.Bool = v; a.Valid = true }

// Reset sets the value to Go's default value and Valid to false.
func (a *Bool) Reset() { *a = Bool{} }

// Ptr returns a pointer to this Bool's value, or a nil pointer if this
// Bool is null.
func (a Bool) Ptr() *bool {
	if !a.Valid {
		return nil
	}
	return &a.Bool
}

// SetPtr sets v according to the rules.
func (a *Bool) SetPtr(v *bool) {
	a.Valid = v != nil
	a.Bool = false
	if v != nil {
		a.Bool = *v
	}
}

// IsZero returns true for invalid Bools, for future omitempty support (Go 1.4?)
// A non-null Bool with a 0 value will not be considered zero.
func (a Bool) IsZero() bool {
	return !a.Valid
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *Bool) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a Bool) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Bool) WriteTo(d Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		d.EscapeBool(w, a.Bool)
	} else {
		_, err = w.WriteString(sqlStrNullUC)
	}
	return
}

// Append appends the value or its nil type to the interface slice.
func (a Bool) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Bool)
	}
	return append(args, nil)
}
