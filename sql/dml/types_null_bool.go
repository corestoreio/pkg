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
	"database/sql"
	"strconv"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/byteconv"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// NullBool is a nullable bool. It does not consider false values to be null. It
// will decode to null, not false, if null. NullBool implements interface
// Argument.
type NullBool struct {
	sql.NullBool
}

// MakeNullBool creates a new NullBool. Implements interface Argument.
func MakeNullBool(b bool, valid ...bool) NullBool {
	v := true
	if len(valid) == 1 {
		v = valid[0]
	}
	return NullBool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: v,
		},
	}
}

// GoString prints an optimized Go representation.
func (a NullBool) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatBool(a.Bool)
}

// GoString prints an optimized Go representation.
func (a NullBool) GoString() string {
	if !a.Valid {
		return "dml.NullBool{}"
	}
	return "dml.MakeNullBool(" + strconv.FormatBool(a.Bool) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null NullBool. It also supports unmarshalling a
// sql.NullBool.
func (a *NullBool) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		a.Bool = x
	case map[string]interface{}:
		dto := &struct {
			NullBool bool
			Valid    bool
		}{}
		err = JSONUnMarshalFn(data, dto)
		a.Bool = dto.NullBool
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NewNotValidf("[dml] json: cannot unmarshal %#v into Go value of type null.NullBool", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler. It will unmarshal to a
// null NullBool if the input is a blank or not an integer. It will return an
// error if the input is not an integer, blank, or "null".
func (a *NullBool) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 || bytes.Equal(text, sqlBytesNullUC) || bytes.Equal(text, sqlBytesNullLC) {
		a.Valid = false
		return nil
	}
	a.NullBool, err = byteconv.ParseNullBool(text)
	return
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullBool is null.
func (a NullBool) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return sqlBytesNullLC, nil
	}
	if !a.Bool {
		return sqlBytesFalseLC, nil
	}
	return sqlBytesTrueLC, nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this NullBool is null.
func (a NullBool) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	if !a.Bool {
		return sqlBytesFalseLC, nil
	}
	return sqlBytesTrueLC, nil
}

// SetValid changes this NullBool's value and also sets it to be non-null.
func (a *NullBool) SetValid(v bool) {
	a.Bool = v
	a.Valid = true
}

// Ptr returns a pointer to this NullBool's value, or a nil pointer if this
// NullBool is null.
func (a NullBool) Ptr() *bool {
	if !a.Valid {
		return nil
	}
	return &a.Bool
}

// IsZero returns true for invalid Bools, for future omitempty support (Go 1.4?)
// A non-null NullBool with a 0 value will not be considered zero.
func (a NullBool) IsZero() bool {
	return !a.Valid
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (a NullBool) GobEncode() ([]byte, error) {
	return a.Marshal()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (a *NullBool) GobDecode(data []byte) error {
	return a.Unmarshal(data)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *NullBool) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a NullBool) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// Marshal binary encoder for protocol buffers. Implements proto.Marshaler.
func (a NullBool) Marshal() ([]byte, error) {
	if !a.Valid {
		return nil, nil
	}
	var buf [1]byte
	_, err := a.MarshalTo(buf[:])
	return buf[:], err
}

// MarshalTo binary encoder for protocol buffers which writes into data.
func (a NullBool) MarshalTo(data []byte) (n int, err error) {
	if !a.Valid {
		return 0, nil
	}
	data[0] = 0
	if a.Bool {
		data[0] = 1
	}
	return 1, nil
}

// Unmarshal binary decoder for protocol buffers. Implements proto.Unmarshaler.
func (a *NullBool) Unmarshal(data []byte) error {
	if len(data) != 1 {
		a.Valid = false
		return nil
	}
	a.Bool = data[0] == 1
	a.Valid = true
	return nil
}

// Size returns the size of the underlying type. If not valid, the size will be
// 0. Implements proto.Sizer.
func (a NullBool) Size() (s int) {
	if a.Valid {
		s = 1
	}
	return
}

func (a NullBool) writeTo(w *bytes.Buffer) (err error) {
	if a.Valid {
		dialect.EscapeBool(w, a.Bool)
	} else {
		_, err = w.WriteString(sqlStrNullUC)
	}
	return
}
func (a NullBool) append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Bool)
	}
	return append(args, nil)
}
