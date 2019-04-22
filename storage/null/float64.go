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
	"math"
	"strconv"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/byteconv"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// Float64 is a nullable float64. It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Float64 struct {
	Float64 float64
	Valid   bool // Valid is true if Float64 is not NULL
}

// MakeFloat64 creates a new Float64. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. Float64
// implements interface Argument.
func MakeFloat64(f float64) Float64 {
	return Float64{
		Float64: f,
		Valid:   true,
	}
}

// MakeFloat64FromByte makes a new Float64 from a (text) byte slice.
func MakeFloat64FromByte(data []byte) (nv Float64, err error) {
	nv.Float64, nv.Valid, err = byteconv.ParseFloat(data)
	return
}

// Scan implements the Scanner interface. Approx. >3x times faster than
// database/sql.convertAssign.
func (a *Float64) Scan(value interface{}) (err error) {
	// this version BenchmarkSQLScanner/NullFloat64_[]byte-4       	20000000	        79.0 ns/op	      32 B/op	       1 allocs/op
	// std lib 		BenchmarkSQLScanner/NullFloat64_[]byte-4       	 5000000	       266 ns/op	      64 B/op	       3 allocs/op
	if value == nil {
		a.Float64, a.Valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		a.Float64, a.Valid, err = byteconv.ParseFloat(v)
	case float64:
		a.Float64 = v
		a.Valid = true
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not yet supported in Float64.Scan", value)
	}
	return
}

// Value implements the driver Valuer interface.
func (a Float64) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Float64, nil
}

// String returns the string representation of the float or null.
func (a Float64) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatFloat(a.Float64, 'f', -1, 64)
}

// GoString prints an optimized Go representation.
func (a Float64) GoString() string {
	if !a.Valid {
		return "null.Float64{}"
	}
	return "null.MakeFloat64(" + strconv.FormatFloat(a.Float64, 'f', -1, 64) + ")"
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Float64.
// It also supports unmarshalling a sql.Float64.
func (a *Float64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = jsonUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		a.Float64 = x
	case map[string]interface{}:
		dto := &struct {
			Float64 float64
			Valid   bool
		}{}
		err = jsonUnMarshalFn(data, dto)
		a.Float64 = dto.Float64
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[dml] json: cannot unmarshal %#v into Go value of type null.Float64", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Float64 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (a *Float64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == sqlStrNullLC {
		a.Valid = false
		return nil
	}
	var err error
	a.Float64, err = strconv.ParseFloat(string(text), 64)
	a.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Float64 is null.
func (a Float64) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return strconv.AppendFloat([]byte{}, a.Float64, 'f', -1, 64), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Float64 is null.
func (a Float64) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendFloat([]byte{}, a.Float64, 'f', -1, 64), nil
}

// SetValid changes this Float64's value and also sets it to be non-null.
func (a Float64) SetValid(n float64) Float64 { a.Float64 = n; a.Valid = true; return a }

// SetNull sets the value to Go's default value and Valid to false.
func (a Float64) SetNull() Float64 { return Float64{} }

// Ptr returns a pointer to this Float64's value, or a nil pointer if this Float64 is null.
func (a Float64) Ptr() *float64 {
	if !a.Valid {
		return nil
	}
	return &a.Float64
}

// IsZero returns true for invalid Float64s, for future omitempty support (Go 1.4?)
// A non-null Float64 with a 0 value will not be considered zero.
func (a Float64) IsZero() bool {
	return !a.Valid
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (a Float64) GobEncode() ([]byte, error) {
	return a.Marshal()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (a *Float64) GobDecode(data []byte) error {
	return a.Unmarshal(data)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *Float64) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a Float64) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// Marshal binary encoder for protocol buffers. Implements proto.Marshaler.
func (a Float64) Marshal() ([]byte, error) {
	if !a.Valid {
		return nil, nil
	}
	var buf [8]byte
	_, err := a.MarshalTo(buf[:])
	return buf[:], err
}

// MarshalTo binary encoder for protocol buffers which writes into data.
func (a Float64) MarshalTo(data []byte) (n int, err error) {
	if !a.Valid {
		return 0, nil
	}
	binary.LittleEndian.PutUint64(data, math.Float64bits(a.Float64))
	return 8, nil
}

// Unmarshal binary decoder for protocol buffers. Implements proto.Unmarshaler.
func (a *Float64) Unmarshal(data []byte) error {
	if len(data) < 8 {
		a.Valid = false
		return nil
	}

	a.Float64 = math.Float64frombits(binary.LittleEndian.Uint64(data))
	a.Valid = true
	return nil
}

// Size returns the size of the underlying type. If not valid, the size will be
// 0. Implements proto.Sizer.
func (a Float64) Size() (s int) {
	if a.Valid {
		s = 8
	}
	return
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Float64) WriteTo(_ Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		return writeFloat64(w, a.Float64)
	}
	_, err = w.WriteString(sqlStrNullUC)
	return err
}

// Append appends the value or its nil type to the interface slice.
func (a Float64) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Float64)
	}
	return append(args, nil)
}
