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
	"math"
	"strconv"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/byteconv"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// Uint16 is a nullable int16. It does not consider zero values to be null.
// It will decode to null, not zero, if null. Uint16 implements interface
// Argument.
type Uint16 struct {
	Uint16 uint16
	Valid  bool // Valid is true if Uint16 is not NULL
}

// MakeUint16 creates a new Uint16. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. Uint16
// implements interface Argument.
func MakeUint16(i uint16) Uint16 {
	return Uint16{
		Uint16: i,
		Valid:  true,
	}
}

// MakeUint16FromByte makes a new Uint16 from a (text) byte slice.
func MakeUint16FromByte(data []byte) (nv Uint16, err error) {
	var i64 uint64
	i64, nv.Valid, err = byteconv.ParseUint(data, 10, 16)
	nv.Uint16 = uint16(i64)
	return
}

// Scan implements the Scanner interface. Approx. >3x times faster than
// database/sql.convertAssign
func (a *Uint16) Scan(value interface{}) (err error) {
	if value == nil {
		a.Uint16, a.Valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		var i64 uint64
		i64, a.Valid, err = byteconv.ParseUint(v, 10, 16)
		a.Valid = err == nil
		a.Uint16 = uint16(i64)
	case int16:
		a.Uint16 = uint16(v)
		a.Valid = true
	case int:
		a.Uint16 = uint16(v)
		a.Valid = true
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not supported in Uint16.Scan", value)
	}
	return
}

// String returns the string representation of the int or null.
func (a Uint16) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatUint(uint64(a.Uint16), 10)
}

// GoString prints an optimized Go representation. Takes are of backticks.
func (a Uint16) GoString() string {
	if !a.Valid {
		return "null.Uint16{}"
	}
	return "null.MakeUint16(" + strconv.FormatUint(uint64(a.Uint16), 10) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null Uint16. It also supports unmarshalling a
// sql.Uint16.
func (a *Uint16) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(bTextNullLC, data) {
		a.Valid = false
		a.Uint16 = 0
		return nil
	}
	if v, ok, err := byteconv.ParseUint(data, 10, 16); ok && err == nil && v >= 0 && v <= math.MaxUint16 {
		a.Valid = true
		a.Uint16 = uint16(v)
		return nil
	}
	var err error
	var v interface{}
	if err = jsonUnMarshalFn(data, &v); err != nil {
		return err
	}

	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int16, to avoid intermediate float64
		err = jsonUnMarshalFn(data, &a.Uint16)
	case map[string]interface{}:
		dto := &struct {
			Uint16 uint16
			Valid  bool
		}{}
		err = jsonUnMarshalFn(data, dto)
		a.Uint16 = dto.Uint16
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[null] json: cannot unmarshal (%T) %#v into Go value of type null.Uint16", v, v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint16 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or sqlStrNullLC.
func (a *Uint16) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 || bytes.Equal(bTextNullLC, text) {
		a.Valid = false
		a.Uint16 = 0
		return nil
	}
	var i64 uint64
	i64, a.Valid, err = byteconv.ParseUint(text, 10, 16)
	a.Uint16 = uint16(i64)
	return
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint16 is null.
func (a Uint16) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return strconv.AppendUint([]byte{}, uint64(a.Uint16), 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint16 is null.
func (a Uint16) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendUint([]byte{}, uint64(a.Uint16), 10), nil
}

// SetValid changes this Uint16's value and also sets it to be non-null.
func (a Uint16) SetValid(n uint16) Uint16 { a.Uint16 = n; a.Valid = true; return a }

// SetNull sets the value to Go's default value and Valid to false.
func (a Uint16) SetNull() Uint16 { return Uint16{} }

// Ptr returns a pointer to this Uint16's value, or a nil pointer if this Uint16 is null.
func (a Uint16) Ptr() *uint16 {
	if !a.Valid {
		return nil
	}
	return &a.Uint16
}

// IsZero returns true for invalid Uint16's, for future omitempty support (Go 1.4?)
// A non-null Uint16 with a 0 value will not be considered zero.
func (a Uint16) IsZero() bool {
	return !a.Valid
}

// Value implements the driver.Valuer interface.
func (a Uint16) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}

	const maxInt16 = 1<<16 - 1
	if a.Uint16 < maxInt16 {
		return int16(a.Uint16), nil
	}
	return strconv.AppendUint([]byte{}, uint64(a.Uint16), 10), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *Uint16) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a Uint16) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Uint16) WriteTo(_ Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		return writeUint64(w, uint64(a.Uint16))
	}
	_, err = w.WriteString(sqlStrNullUC)
	return err
}

// Append appends the value or its nil type to the interface slice.
func (a Uint16) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Uint16)
	}
	return append(args, nil)
}
