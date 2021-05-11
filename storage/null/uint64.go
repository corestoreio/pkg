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

// TODO(cys): Remove  MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// Uint64 is a nullable int64. It does not consider zero values to be null.
// It will decode to null, not zero, if null. Uint64 implements interface
// Argument.
type Uint64 struct {
	Uint64 uint64
	Valid  bool // Valid is true if Uint64 is not NULL
}

// MakeUint64 creates a new Uint64. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. Uint64
// implements interface Argument.
func MakeUint64(i uint64) Uint64 {
	return Uint64{
		Uint64: i,
		Valid:  true,
	}
}

// MakeUint64FromByte makes a new Uint64 from a (text) byte slice.
func MakeUint64FromByte(data []byte) (nv Uint64, err error) {
	nv.Uint64, nv.Valid, err = byteconv.ParseUint(data, 10, 64)
	return
}

// Scan implements the Scanner interface. Approx. >3x times faster than
// database/sql.convertAssign
func (a *Uint64) Scan(value interface{}) (err error) {
	if value == nil {
		a.Uint64, a.Valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		a.Uint64, a.Valid, err = byteconv.ParseUint(v, 10, 64)
		a.Valid = err == nil
	case int64:
		a.Uint64 = uint64(v)
		a.Valid = true
	case int:
		a.Uint64 = uint64(v)
		a.Valid = true
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not supported in Uint64.Scan", value)
	}
	return
}

// String returns the string representation of the int or null.
func (a Uint64) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatUint(a.Uint64, 10)
}

// GoString prints an optimized Go representation. Takes are of backticks.
func (a Uint64) GoString() string {
	if !a.Valid {
		return "null.Uint64{}"
	}
	return "null.MakeUint64(" + strconv.FormatUint(a.Uint64, 10) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null Uint64. It also supports unmarshalling a
// sql.Uint64.
func (a *Uint64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(bTextNullLC, data) {
		a.Valid = false
		a.Uint64 = 0
		return nil
	}
	if v, ok, err := byteconv.ParseUint(data, 10, 64); ok && err == nil && v >= 0 && v <= math.MaxUint64 {
		a.Valid = true
		a.Uint64 = v
		return nil
	}
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &a.Uint64)
	case map[string]interface{}:
		dto := &struct {
			Uint64 int64
			Valid  bool
		}{}
		err = json.Unmarshal(data, dto)
		a.Uint64 = uint64(dto.Uint64)
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[null] json: cannot unmarshal %#v into Go value of type null.Uint64", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint64 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or sqlStrNullLC.
func (a *Uint64) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 || bytes.Equal(bTextNullLC, text) {
		a.Valid = false
		a.Uint64 = 0
		return nil
	}
	a.Uint64, a.Valid, err = byteconv.ParseUint(text, 10, 64)
	return
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint64 is null.
func (a Uint64) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return strconv.AppendUint([]byte{}, a.Uint64, 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint64 is null.
func (a Uint64) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendUint([]byte{}, a.Uint64, 10), nil
}

// SetValid changes this Uint64's value and also sets it to be non-null.
func (a Uint64) SetValid(n uint64) Uint64 { a.Uint64 = n; a.Valid = true; return a }

// SetNull sets the value to Go's default value and Valid to false.
func (a Uint64) SetNull() Uint64 { return Uint64{} }

// Ptr returns a pointer to this Uint64's value, or a nil pointer if this Uint64 is null.
func (a Uint64) Ptr() *uint64 {
	if !a.Valid {
		return nil
	}
	return &a.Uint64
}

// IsZero returns true for invalid Uint64's, for future omitempty support (Go 1.4?)
// A non-null Uint64 with a 0 value will not be considered zero.
func (a Uint64) IsZero() bool {
	return !a.Valid
}

// Value implements the driver.Valuer interface.
func (a Uint64) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	const maxInt64 = 1<<63 - 1
	if a.Uint64 < maxInt64 {
		return int64(a.Uint64), nil
	}
	return strconv.AppendUint([]byte{}, a.Uint64, 10), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *Uint64) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a Uint64) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Uint64) WriteTo(_ Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		return writeUint64(w, a.Uint64)
	}
	_, err = w.WriteString(sqlStrNullUC)
	return err
}

// Append appends the value or its nil type to the interface slice.
func (a Uint64) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Uint64)
	}
	return append(args, nil)
}
