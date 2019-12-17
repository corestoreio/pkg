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

// Int8 is a nullable int8. It does not consider zero values to be null.
// It will decode to null, not zero, if null. Int8 implements interface
// Argument.
type Int8 struct {
	Int8  int8
	Valid bool // Valid is true if Int8 is not NULL
}

// MakeInt8 creates a new Int8. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. Int8
// implements interface Argument.
func MakeInt8(i int8) Int8 {
	return Int8{
		Int8:  i,
		Valid: true,
	}
}

// MakeInt8FromByte makes a new Int8 from a (text) byte slice.
func MakeInt8FromByte(data []byte) (nv Int8, err error) {
	var i64 int64
	i64, nv.Valid, err = byteconv.ParseInt(data)
	nv.Int8 = int8(i64)
	return
}

// Scan implements the Scanner interface. Approx. >3x times faster than
// database/sql.convertAssign.
func (a *Int8) Scan(value interface{}) (err error) {
	// this version BenchmarkSQLScanner/NullInt8_[]byte-4         	20000000	        65.0 ns/op	      8 B/op	       1 allocs/op
	// std lib 		BenchmarkSQLScanner/NullInt8_[]byte-4         	 5000000	       244 ns/op	      56 B/op	       3 allocs/op
	if value == nil {
		a.Int8, a.Valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		var i64 int64
		i64, a.Valid, err = byteconv.ParseInt(v)
		a.Int8 = int8(i64)
	case int8:
		a.Int8 = v
		a.Valid = true
	case int:
		a.Int8 = int8(v)
		a.Valid = true
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not yet supported in Int8.Scan", value)
	}
	return
}

// String returns the string representation of the int or null.
func (a Int8) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatInt(int64(a.Int8), 10)
}

// GoString prints an optimized Go representation. Takes are of backticks.
func (a Int8) GoString() string {
	if !a.Valid {
		return "null.Int8{}"
	}
	return "null.MakeInt8(" + strconv.FormatInt(int64(a.Int8), 10) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null Int8. It also supports unmarshalling a
// sql.Int8.
func (a *Int8) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(bTextNullLC, data) {
		a.Valid = false
		a.Int8 = 0
		return nil
	}
	if v, ok, err := byteconv.ParseInt(data); ok && err == nil && v >= -math.MaxInt8 && v <= math.MaxInt8 {
		a.Valid = true
		a.Int8 = int8(v)
		return nil
	}

	var err error
	var v interface{}
	if err = jsonUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int8, to avoid intermediate float8
		err = jsonUnMarshalFn(data, &a.Int8)
	case map[string]interface{}:
		dto := &struct {
			Int8  int8
			Valid bool
		}{}
		err = jsonUnMarshalFn(data, dto)
		a.Int8 = dto.Int8
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[null] json: cannot unmarshal %#v into Go value of type null.Int8", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int8 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or sqlStrNullLC.
func (a *Int8) UnmarshalText(text []byte) error {
	if len(text) == 0 || bytes.Equal(bTextNullLC, text) {
		a.Valid = false
		a.Int8 = 0
		return nil
	}
	ni, ok, err := byteconv.ParseInt(text)
	a.Int8 = int8(ni)
	a.Valid = ok
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int8 is null.
func (a Int8) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return strconv.AppendInt([]byte{}, int64(a.Int8), 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int8 is null.
func (a Int8) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendInt([]byte{}, int64(a.Int8), 10), nil
}

// SetValid changes this Int8's value and also sets it to be non-null.
func (a Int8) SetValid(n int8) Int8 { a.Int8 = n; a.Valid = true; return a }

// SetNull sets the value to Go's default value and Valid to false.
func (a Int8) SetNull() Int8 { return Int8{} }

// Ptr returns a pointer to this Int8's value, or a nil pointer if this Int8 is null.
func (a Int8) Ptr() *int8 {
	if !a.Valid {
		return nil
	}
	return &a.Int8
}

// IsZero returns true for invalid Int8's, for future omitempty support (Go 1.4?)
// A non-null Int8 with a 0 value will not be considered zero.
func (a Int8) IsZero() bool {
	return !a.Valid
}

// Value implements the driver.Valuer interface.
func (a Int8) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Int8, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *Int8) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a Int8) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Int8) WriteTo(_ Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		return writeInt64(w, int64(a.Int8))
	}
	_, err = w.WriteString(sqlStrNullUC)
	return err
}

// Append appends the value or its nil type to the interface slice.
func (a Int8) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Int8)
	}
	return append(args, nil)
}
