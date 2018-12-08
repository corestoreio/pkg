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
	"encoding/binary"
	"strconv"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/byteconv"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// Uint8 is a nullable int8. It does not consider zero values to be null.
// It will decode to null, not zero, if null. Uint8 implements interface
// Argument.
type Uint8 struct {
	Uint8 uint8
	Valid bool // Valid is true if Uint8 is not NULL
}

// MakeUint8 creates a new Uint8. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. Uint8
// implements interface Argument.
func MakeUint8(i uint8) Uint8 {
	return Uint8{
		Uint8: i,
		Valid: true,
	}
}

// MakeUint8FromByte makes a new Uint8 from a (text) byte slice.
func MakeUint8FromByte(data []byte) (nv Uint8, err error) {
	var i64 uint64
	i64, nv.Valid, err = byteconv.ParseUint(data, 10, 8)
	nv.Uint8 = uint8(i64)
	return
}

// Scan implements the Scanner interface. Approx. >3x times faster than
// database/sql.convertAssign
func (a *Uint8) Scan(value interface{}) (err error) {
	if value == nil {
		a.Uint8, a.Valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		var i64 uint64
		i64, a.Valid, err = byteconv.ParseUint(v, 10, 8)
		a.Valid = err == nil
		a.Uint8 = uint8(i64)
	case int8:
		a.Uint8 = uint8(v)
		a.Valid = true
	case int:
		a.Uint8 = uint8(v)
		a.Valid = true
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not supported in Uint8.Scan", value)
	}
	return
}

// String returns the string representation of the int or null.
func (a Uint8) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatUint(uint64(a.Uint8), 10)
}

// GoString prints an optimized Go representation. Takes are of backticks.
func (a Uint8) GoString() string {
	if !a.Valid {
		return "null.Uint8{}"
	}
	return "null.MakeUint8(" + strconv.FormatUint(uint64(a.Uint8), 10) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null Uint8. It also supports unmarshalling a
// sql.Uint8.
func (a *Uint8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}

	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int8, to avoid intermediate float64
		err = JSONUnMarshalFn(data, &a.Uint8)
	case map[string]interface{}:
		dto := &struct {
			Uint8 uint8
			Valid bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		a.Uint8 = uint8(dto.Uint8)
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[null] json: cannot unmarshal (%T) %#v into Go value of type null.Uint8", v, v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint8 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or sqlStrNullLC.
func (a *Uint8) UnmarshalText(text []byte) (err error) {
	str := string(text)
	if str == "" || str == sqlStrNullLC {
		a.Valid = false
		return nil
	}
	var i64 uint64
	i64, a.Valid, err = byteconv.ParseUint(text, 10, 8)
	a.Uint8 = uint8(i64)
	return
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint8 is null.
func (a Uint8) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return strconv.AppendUint([]byte{}, uint64(a.Uint8), 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint8 is null.
func (a Uint8) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendUint([]byte{}, uint64(a.Uint8), 10), nil
}

// SetValid changes this Uint8's value and also sets it to be non-null.
func (a Uint8) SetValid(n uint8) Uint8 { a.Uint8 = n; a.Valid = true; return a }

// SetNull sets the value to Go's default value and Valid to false.
func (a Uint8) SetNull() Uint8 { return Uint8{} }

// Ptr returns a pointer to this Uint8's value, or a nil pointer if this Uint8 is null.
func (a Uint8) Ptr() *uint8 {
	if !a.Valid {
		return nil
	}
	return &a.Uint8
}

// IsZero returns true for invalid Uint8's, for future omitempty support (Go 1.4?)
// A non-null Uint8 with a 0 value will not be considered zero.
func (a Uint8) IsZero() bool {
	return !a.Valid
}

// Value implements the driver.Valuer interface.
func (a Uint8) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}

	const maxInt8 = 1<<8 - 1
	if a.Uint8 < maxInt8 {
		return int8(a.Uint8), nil
	}
	return strconv.AppendUint([]byte{}, uint64(a.Uint8), 10), nil
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (a Uint8) GobEncode() ([]byte, error) {
	return a.Marshal()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (a *Uint8) GobDecode(data []byte) error {
	return a.Unmarshal(data)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *Uint8) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a Uint8) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// Marshal binary encoder for protocol buffers. Implements proto.Marshaler.
func (a Uint8) Marshal() ([]byte, error) {
	if !a.Valid {
		return nil, nil
	}
	var buf [8]byte
	_, err := a.MarshalTo(buf[:])
	return buf[:], err
}

// MarshalTo binary encoder for protocol buffers which writes into data.
func (a Uint8) MarshalTo(data []byte) (n int, err error) {
	if !a.Valid {
		return 0, nil
	}
	binary.LittleEndian.PutUint16(data, uint16(a.Uint8))
	return 8, nil
}

// Unmarshal binary decoder for protocol buffers. Implements proto.Unmarshaler.
func (a *Uint8) Unmarshal(data []byte) error {
	if len(data) < 2 {
		a.Valid = false
		return nil
	}
	u16 := binary.LittleEndian.Uint16(data)
	a.Uint8 = uint8(u16)
	a.Valid = u16 > 0 && u16 <= (1<<8-1)
	return nil
}

// Size returns the size of the underlying type. If not valid, the size will be
// 0. Implements proto.Sizer.
func (a Uint8) Size() (s int) {
	if a.Valid {
		s = 8
	}
	return
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Uint8) WriteTo(_ Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		return writeUint64(w, uint64(a.Uint8))
	}
	_, err = w.WriteString(sqlStrNullUC)
	return err
}

// Append appends the value or its nil type to the interface slice.
func (a Uint8) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Uint8)
	}
	return append(args, nil)
}
