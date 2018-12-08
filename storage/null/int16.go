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

// Int16 is a nullable int16. It does not consider zero values to be null.
// It will decode to null, not zero, if null. Int16 implements interface
// Argument.
type Int16 struct {
	Int16 int16
	Valid bool // Valid is true if Int16 is not NULL
}

// MakeInt16 creates a new Int16. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. Int16
// implements interface Argument.
func MakeInt16(i int16) Int16 {
	return Int16{
		Int16: i,
		Valid: true,
	}
}

// MakeInt16FromByte makes a new Int16 from a (text) byte slice.
func MakeInt16FromByte(data []byte) (nv Int16, err error) {
	var i64 int64
	i64, nv.Valid, err = byteconv.ParseInt(data)
	nv.Int16 = int16(i64)
	return
}

// Scan implements the Scanner interface. Approx. >3x times faster than
// database/sql.convertAssign.
func (a *Int16) Scan(value interface{}) (err error) {
	// this version BenchmarkSQLScanner/NullInt16_[]byte-4         	20000000	        65.0 ns/op	      16 B/op	       1 allocs/op
	// std lib 		BenchmarkSQLScanner/NullInt16_[]byte-4         	 5000000	       244 ns/op	      56 B/op	       3 allocs/op
	if value == nil {
		a.Int16, a.Valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		var i64 int64
		i64, a.Valid, err = byteconv.ParseInt(v)
		a.Int16 = int16(i64)
	case int16:
		a.Int16 = v
		a.Valid = true
	case int:
		a.Int16 = int16(v)
		a.Valid = true
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not yet supported in Int16.Scan", value)
	}
	return
}

// String returns the string representation of the int or null.
func (a Int16) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatInt(int64(a.Int16), 10)
}

// GoString prints an optimized Go representation. Takes are of backticks.
func (a Int16) GoString() string {
	if !a.Valid {
		return "null.Int16{}"
	}
	return "null.MakeInt16(" + strconv.FormatInt(int64(a.Int16), 10) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null Int16. It also supports unmarshalling a
// sql.Int16.
func (a *Int16) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int16, to avoid intermediate float16
		err = JSONUnMarshalFn(data, &a.Int16)
	case map[string]interface{}:
		dto := &struct {
			Int16 int16
			Valid bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		a.Int16 = dto.Int16
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[null] json: cannot unmarshal %#v into Go value of type null.Int16", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int16 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or sqlStrNullLC.
func (a *Int16) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == sqlStrNullLC {
		a.Valid = false
		return nil
	}
	ni, ok, err := byteconv.ParseInt(text)
	a.Int16 = int16(ni)
	a.Valid = ok
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int16 is null.
func (a Int16) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return strconv.AppendInt([]byte{}, int64(a.Int16), 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int16 is null.
func (a Int16) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendInt([]byte{}, int64(a.Int16), 10), nil
}

// SetValid changes this Int16's value and also sets it to be non-null.
func (a Int16) SetValid(n int16) Int16 { a.Int16 = n; a.Valid = true; return a }

// SetNull sets the value to Go's default value and Valid to false.
func (a Int16) SetNull() Int16 { return Int16{} }

// Ptr returns a pointer to this Int16's value, or a nil pointer if this Int16 is null.
func (a Int16) Ptr() *int16 {
	if !a.Valid {
		return nil
	}
	return &a.Int16
}

// IsZero returns true for invalid Int16's, for future omitempty support (Go 1.4?)
// A non-null Int16 with a 0 value will not be considered zero.
func (a Int16) IsZero() bool {
	return !a.Valid
}

// Value implements the driver.Valuer interface.
func (a Int16) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Int16, nil
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (a Int16) GobEncode() ([]byte, error) {
	return a.Marshal()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (a *Int16) GobDecode(data []byte) error {
	return a.Unmarshal(data)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *Int16) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a Int16) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// Marshal binary encoder for protocol buffers. Implements proto.Marshaler.
func (a Int16) Marshal() ([]byte, error) {
	if !a.Valid {
		return nil, nil
	}
	var buf [8]byte
	_, err := a.MarshalTo(buf[:])
	return buf[:], err
}

// MarshalTo binary encoder for protocol buffers which writes into data.
func (a Int16) MarshalTo(data []byte) (n int, err error) {
	if !a.Valid {
		return 0, nil
	}
	binary.LittleEndian.PutUint16(data, uint16(a.Int16))
	return 2, nil
}

// Unmarshal binary decoder for protocol buffers. Implements proto.Unmarshaler.
func (a *Int16) Unmarshal(data []byte) error {
	if len(data) < 8 {
		a.Valid = false
		return nil
	}
	ui := binary.LittleEndian.Uint16(data)
	a.Int16 = int16(ui)
	a.Valid = true
	return nil
}

// Size returns the size of the underlying type. If not valid, the size will be
// 0. Implements proto.Sizer.
func (a Int16) Size() (s int) {
	if a.Valid {
		s = 8
	}
	return
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Int16) WriteTo(_ Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		return writeInt64(w, int64(a.Int16))
	}
	_, err = w.WriteString(sqlStrNullUC)
	return err
}

// Append appends the value or its nil type to the interface slice.
func (a Int16) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Int16)
	}
	return append(args, nil)
}
