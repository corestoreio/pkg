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

package dml

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

// NullUint64 is a nullable int64. It does not consider zero values to be null.
// It will decode to null, not zero, if null. NullUint64 implements interface
// Argument.
type NullUint64 struct {
	Uint64 uint64
	Valid  bool // Valid is true if Uint64 is not NULL
}

// MakeNullUint64 creates a new NullUint64. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. NullUint64
// implements interface Argument.
func MakeNullUint64(i uint64, valid ...bool) NullUint64 {
	v := true
	if len(valid) == 1 {
		v = valid[0]
	}
	return NullUint64{
		Uint64: i,
		Valid:  v,
	}
}

// Scan implements the Scanner interface. Approx. >3x times faster than
// database/sql.convertAssign
func (n *NullUint64) Scan(value interface{}) (err error) {
	if value == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		n.Uint64, n.Valid, err = byteconv.ParseUintSQL(v, 10, 64)
		n.Valid = err == nil
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not supported in NullUint64.Scan", value)
	}
	return
}

// String returns the string representation of the int or null.
func (a NullUint64) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatUint(a.Uint64, 10)
}

// GoString prints an optimized Go representation. Takes are of backticks.
func (a NullUint64) GoString() string {
	if !a.Valid {
		return "dml.NullUint64{}"
	}
	return "dml.MakeNullUint64(" + strconv.FormatUint(a.Uint64, 10) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null NullUint64. It also supports unmarshalling a
// sql.NullUint64.
func (a *NullUint64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = JSONUnMarshalFn(data, &a.Uint64)
	case map[string]interface{}:
		dto := &struct {
			NullUint64 int64
			Valid      bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		a.Uint64 = uint64(dto.NullUint64)
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[null] json: cannot unmarshal %#v into Go value of type null.NullUint64", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null NullUint64 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or sqlStrNullLC.
func (a *NullUint64) UnmarshalText(text []byte) (err error) {
	str := string(text)
	if str == "" || str == sqlStrNullLC {
		a.Valid = false
		return nil
	}
	a.Uint64, a.Valid, err = byteconv.ParseUintSQL(text, 10, 64)
	return
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullUint64 is null.
func (a NullUint64) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return sqlBytesNullLC, nil
	}
	return strconv.AppendUint([]byte{}, a.Uint64, 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this NullUint64 is null.
func (a NullUint64) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendUint([]byte{}, a.Uint64, 10), nil
}

// SetValid changes this NullUint64's value and also sets it to be non-null.
func (a *NullUint64) SetValid(n uint64) {
	a.Uint64 = n
	a.Valid = true
}

// Ptr returns a pointer to this NullUint64's value, or a nil pointer if this NullUint64 is null.
func (a NullUint64) Ptr() *uint64 {
	if !a.Valid {
		return nil
	}
	return &a.Uint64
}

// IsZero returns true for invalid NullUint64's, for future omitempty support (Go 1.4?)
// A non-null NullUint64 with a 0 value will not be considered zero.
func (a NullUint64) IsZero() bool {
	return !a.Valid
}

// Value implements the driver.Valuer interface.
func (a NullUint64) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	const maxInt64 = 1<<63 - 1
	if a.Uint64 < maxInt64 {
		return int64(a.Uint64), nil
	}
	return strconv.AppendUint([]byte{}, a.Uint64, 10), nil
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (a NullUint64) GobEncode() ([]byte, error) {
	return a.Marshal()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (a *NullUint64) GobDecode(data []byte) error {
	return a.Unmarshal(data)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *NullUint64) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a NullUint64) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// Marshal binary encoder for protocol buffers. Implements proto.Marshaler.
func (a NullUint64) Marshal() ([]byte, error) {
	if !a.Valid {
		return nil, nil
	}
	var buf [8]byte
	_, err := a.MarshalTo(buf[:])
	return buf[:], err
}

// MarshalTo binary encoder for protocol buffers which writes into data.
func (a NullUint64) MarshalTo(data []byte) (n int, err error) {
	if !a.Valid {
		return 0, nil
	}
	binary.LittleEndian.PutUint64(data, uint64(a.Uint64))
	return 8, nil
}

// Unmarshal binary decoder for protocol buffers. Implements proto.Unmarshaler.
func (a *NullUint64) Unmarshal(data []byte) error {
	if len(data) < 8 {
		a.Valid = false
		return nil
	}
	a.Uint64 = binary.LittleEndian.Uint64(data)
	a.Valid = true
	return nil
}

// Size returns the size of the underlying type. If not valid, the size will be
// 0. Implements proto.Sizer.
func (a NullUint64) Size() (s int) {
	if a.Valid {
		s = 8
	}
	return
}

func (a NullUint64) writeTo(w *bytes.Buffer) error {
	if a.Valid {
		return writeUint64(w, a.Uint64)
	}
	_, err := w.WriteString(sqlStrNullUC)
	return err
}

func (a NullUint64) append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Uint64)
	}
	return append(args, nil)
}
