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
	"strconv"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/byteconv"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// Int64 is a nullable int64. It does not consider zero values to be null.
// It will decode to null, not zero, if null. Int64 implements interface
// Argument.
type Int64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// MakeInt64 creates a new Int64. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. Int64
// implements interface Argument.
func MakeInt64(i int64) Int64 {
	return Int64{
		Int64: i,
		Valid: true,
	}
}

// MakeInt64FromByte makes a new Int64 from a (text) byte slice.
func MakeInt64FromByte(data []byte) (nv Int64, err error) {
	nv.Int64, nv.Valid, err = byteconv.ParseInt(data)
	return
}

// Scan implements the Scanner interface. Approx. >3x times faster than
// database/sql.convertAssign.
func (a *Int64) Scan(value interface{}) (err error) {
	// this version BenchmarkSQLScanner/NullInt64_[]byte-4         	20000000	        65.0 ns/op	      32 B/op	       1 allocs/op
	// std lib 		BenchmarkSQLScanner/NullInt64_[]byte-4         	 5000000	       244 ns/op	      56 B/op	       3 allocs/op
	if value == nil {
		a.Int64, a.Valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		a.Int64, a.Valid, err = byteconv.ParseInt(v)
	case int64:
		a.Int64 = v
		a.Valid = true
	case int:
		a.Int64 = int64(v)
		a.Valid = true
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not yet supported in Int64.Scan", value)
	}
	return
}

// String returns the string representation of the int or null.
func (a Int64) String() string {
	if !a.Valid {
		return "null"
	}
	return strconv.FormatInt(a.Int64, 10)
}

// GoString prints an optimized Go representation. Takes are of backticks.
func (a Int64) GoString() string {
	if !a.Valid {
		return "null.Int64{}"
	}
	return "null.MakeInt64(" + strconv.FormatInt(a.Int64, 10) + ")"
}

// UnmarshalJSON implements json.Unmarshaler. It supports number and null input.
// 0 will not be considered a null Int64. It also supports unmarshalling a
// sql.Int64.
func (a *Int64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(bTextNullLC, data) {
		a.Valid = false
		a.Int64 = 0
		return nil
	}
	if v, ok, err := byteconv.ParseInt(data); ok && err == nil {
		a.Int64 = v
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
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &a.Int64)
	case map[string]interface{}:
		dto := &struct {
			Int64 int64
			Valid bool
		}{}
		err = json.Unmarshal(data, dto)
		a.Int64 = dto.Int64
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[null] json: cannot unmarshal %#v into Go value of type null.Int64", v)
	}
	a.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int64 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or sqlStrNullLC.
func (a *Int64) UnmarshalText(text []byte) error {
	if len(text) == 0 || bytes.Equal(bTextNullLC, text) {
		a.Valid = false
		a.Int64 = 0
		return nil
	}
	ni, ok, err := byteconv.ParseInt(text)
	a.Int64 = ni
	a.Valid = ok
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int64 is null.
func (a Int64) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return strconv.AppendInt([]byte{}, a.Int64, 10), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int64 is null.
func (a Int64) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte{}, nil
	}
	return strconv.AppendInt([]byte{}, a.Int64, 10), nil
}

// SetValid changes this Int64's value and also sets it to be non-null.
func (a *Int64) SetValid(n int64) { a.Int64 = n; a.Valid = true }

// Reset sets the value to Go's default value and Valid to false.
func (a *Int64) Reset() { *a = Int64{} }

// Ptr returns a pointer to this value, or a nil pointer if value is null.
func (a Int64) Ptr() *int64 {
	if !a.Valid {
		return nil
	}
	return &a.Int64
}

// SetPtr sets v according to the rules.
func (a *Int64) SetPtr(v *int64) {
	a.Valid = v != nil
	a.Int64 = 0
	if v != nil {
		a.Int64 = *v
	}
}

// IsZero returns true for invalid Int64's, for future omitempty support (Go 1.4?)
// A non-null Int64 with a 0 value will not be considered zero.
func (a Int64) IsZero() bool {
	return !a.Valid
}

// Value implements the driver.Valuer interface.
func (a Int64) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Int64, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *Int64) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a Int64) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Int64) WriteTo(_ Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		return writeInt64(w, a.Int64)
	}
	_, err = w.WriteString(sqlStrNullUC)
	return err
}

// Append appends the value or its nil type to the interface slice.
func (a Int64) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Int64)
	}
	return append(args, nil)
}
