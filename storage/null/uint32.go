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
	"math"
	"strconv"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/byteconv"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// Uint32 is a nullable int32. It does not consider zero values to be null.
// It will decode to null, not zero, if null. Uint32 implements interface
// Argument.
type Uint32 struct {
	Uint32 uint32
	Valid  bool // Valid is true if Uint32 is not NULL
}

// MakeUint32 creates a new Uint32. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. Uint32
// implements interface Argument.
func MakeUint32(i uint32) Uint32 {
	return Uint32{
		Uint32: i,
		Valid:  true,
	}
}

// MakeUint32FromByte makes a new Uint32 from a (text) byte slice.
func MakeUint32FromByte(data []byte) (nv Uint32, err error) {
	var i64 uint64
	i64, nv.Valid, err = byteconv.ParseUint(data, 10, 32)
	nv.Uint32 = uint32(i64)
	return
}

// Scan implements the Scanner interface. Approx. >3x times faster than
// database/sql.convertAssign
func (a *Uint32) Scan(value interface{}) (err error) {
	if value == nil {
		a.Uint32, a.Valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		var i64 uint64
		i64, a.Valid, err = byteconv.ParseUint(v, 10, 32)
		a.Valid = err == nil
		a.Uint32 = uint32(i64)
	case int32:
		a.Uint32 = uint32(v)
		a.Valid = true
	case int:
		a.Uint32 = uint32(v)
		a.Valid = true
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not supported in Uint32.Scan", value)
	}
	return
}

// String returns the string representation of the int or null.
func (a Uint32) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatUint(uint64(a.Uint32), 10)
}

// GoString prints an optimized Go representation. Takes are of backticks.
func (a Uint32) GoString() string {
	if !a.Valid {
		return "null.Uint32{}"
	}
	return "null.MakeUint32(" + strconv.FormatUint(uint64(a.Uint32), 10) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null Uint32. It also supports unmarshalling a
// sql.Uint32.
func (a *Uint32) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(bTextNullLC, data) {
		a.Valid = false
		a.Uint32 = 0
		return nil
	}
	if v, ok, err := byteconv.ParseUint(data, 10, 32); ok && err == nil && v >= 0 && v <= math.MaxUint32 {
		a.Valid = true
		a.Uint32 = uint32(v)
		return nil
	}
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int32, to avoid intermediate float64
		err = json.Unmarshal(data, &a.Uint32)
	case map[string]interface{}:
		dto := &struct {
			Uint32 uint32
			Valid  bool
		}{}
		err = json.Unmarshal(data, dto)
		a.Uint32 = dto.Uint32
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[null] json: cannot unmarshal (%T) %#v into Go value of type null.Uint32", v, v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint32 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or sqlStrNullLC.
func (a *Uint32) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 || bytes.Equal(bTextNullLC, text) {
		a.Valid = false
		a.Uint32 = 0
		return nil
	}
	var i64 uint64
	i64, a.Valid, err = byteconv.ParseUint(text, 10, 32)
	a.Uint32 = uint32(i64)
	return
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint32 is null.
func (a Uint32) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return strconv.AppendUint([]byte{}, uint64(a.Uint32), 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint32 is null.
func (a Uint32) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendUint([]byte{}, uint64(a.Uint32), 10), nil
}

// SetValid changes this Uint32's value and also sets it to be non-null.
func (a Uint32) SetValid(n uint32) Uint32 { a.Uint32 = n; a.Valid = true; return a }

// SetNull sets the value to Go's default value and Valid to false.
func (a Uint32) SetNull() Uint32 { return Uint32{} }

// Ptr returns a pointer to this Uint32's value, or a nil pointer if this Uint32 is null.
func (a Uint32) Ptr() *uint32 {
	if !a.Valid {
		return nil
	}
	return &a.Uint32
}

// IsZero returns true for invalid Uint32's, for future omitempty support (Go 1.4?)
// A non-null Uint32 with a 0 value will not be considered zero.
func (a Uint32) IsZero() bool {
	return !a.Valid
}

// Value implements the driver.Valuer interface.
func (a Uint32) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	const maxInt32 = 1<<31 - 1
	if a.Uint32 < maxInt32 {
		return int32(a.Uint32), nil
	}
	return strconv.AppendUint([]byte{}, uint64(a.Uint32), 10), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *Uint32) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a Uint32) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Uint32) WriteTo(_ Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		return writeUint64(w, uint64(a.Uint32))
	}
	_, err = w.WriteString(sqlStrNullUC)
	return err
}

// Append appends the value or its nil type to the interface slice.
func (a Uint32) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Uint32)
	}
	return append(args, nil)
}
