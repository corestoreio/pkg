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
	"database/sql"
	"encoding/json"
	"math"
	"strconv"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/byteconv"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// Int32 is a nullable int32. It does not consider zero values to be null.
// It will decode to null, not zero, if null. Int32 implements interface
// Argument.
type Int32 struct {
	sql.NullInt32
}

// MakeInt32 creates a new Int32. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. Int32
// implements interface Argument.
func MakeInt32(i int32) Int32 {
	return Int32{
		NullInt32: sql.NullInt32{Int32: i, Valid: true},
	}
}

// MakeInt32FromByte makes a new Int32 from a (text) byte slice.
func MakeInt32FromByte(data []byte) (nv Int32, err error) {
	var i64 int64
	i64, nv.Valid, err = byteconv.ParseInt(data)
	nv.Int32 = int32(i64)
	return
}

// Scan implements the Scanner interface. Approx. >3x times faster than
// database/sql.convertAssign.
func (a *Int32) Scan(value interface{}) (err error) {
	// this version BenchmarkSQLScanner/NullInt32_[]byte-4         	20000000	        65.0 ns/op	      32 B/op	       1 allocs/op
	// std lib 		BenchmarkSQLScanner/NullInt32_[]byte-4         	 5000000	       244 ns/op	      56 B/op	       3 allocs/op
	if value == nil {
		a.Int32, a.Valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		var i64 int64
		i64, a.Valid, err = byteconv.ParseInt(v)
		a.Int32 = int32(i64)
	case int32:
		a.Int32 = v
		a.Valid = true
	case int:
		a.Int32 = int32(v)
		a.Valid = true
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not yet supported in Int32.Scan", value)
	}
	return
}

// String returns the string representation of the int or null.
func (a Int32) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatInt(int64(a.Int32), 10)
}

// GoString prints an optimized Go representation. Takes are of backticks.
func (a Int32) GoString() string {
	if !a.Valid {
		return "null.Int32{}"
	}
	return "null.MakeInt32(" + strconv.FormatInt(int64(a.Int32), 10) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null Int32. It also supports unmarshalling a
// sql.Int32.
func (a *Int32) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(bTextNullLC, data) {
		a.Valid = false
		a.Int32 = 0
		return nil
	}
	if v, ok, err := byteconv.ParseInt(data); ok && err == nil && v >= -math.MaxInt32 && v <= math.MaxInt32 {
		a.Int32 = int32(v)
		a.Valid = true
		return nil
	}
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int32, to avoid intermediate float32
		err = json.Unmarshal(data, &a.Int32)
	case map[string]interface{}:
		dto := &struct {
			Int32 int32
			Valid bool
		}{}
		err = json.Unmarshal(data, dto)
		a.Int32 = dto.Int32
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[null] json: cannot unmarshal %#v into Go value of type null.Int32", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int32 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or sqlStrNullLC.
func (a *Int32) UnmarshalText(text []byte) error {
	if len(text) == 0 || bytes.Equal(bTextNullLC, text) {
		a.Valid = false
		a.Int32 = 0
		return nil
	}
	ni, ok, err := byteconv.ParseInt(text)
	a.Int32 = int32(ni)
	a.Valid = ok
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int32 is null.
func (a Int32) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return strconv.AppendInt([]byte{}, int64(a.Int32), 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int32 is null.
func (a Int32) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendInt([]byte{}, int64(a.Int32), 10), nil
}

// SetValid changes this Int32's value and also sets it to be non-null.
func (a *Int32) SetValid(n int32) { a.Int32 = n; a.Valid = true }

// Reset sets the value to Go's default value and Valid to false.
func (a *Int32) Reset() { *a = Int32{} }

// Ptr returns a pointer to this Int32's value, or a nil pointer if this Int32 is null.
func (a Int32) Ptr() *int32 {
	if !a.Valid {
		return nil
	}
	return &a.Int32
}

// SetPtr sets v according to the rules.
func (a *Int32) SetPtr(v *int32) {
	a.Valid = v != nil
	a.Int32 = 0
	if v != nil {
		a.Int32 = *v
	}
}

// IsZero returns true for invalid Int32's, for future omitempty support (Go 1.4?)
// A non-null Int32 with a 0 value will not be considered zero.
func (a Int32) IsZero() bool {
	return !a.Valid
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *Int32) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a Int32) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Int32) WriteTo(_ Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		return writeInt64(w, int64(a.Int32))
	}
	_, err = w.WriteString(sqlStrNullUC)
	return err
}

// Append appends the value or its nil type to the interface slice.
func (a Int32) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Int32)
	}
	return append(args, nil)
}
