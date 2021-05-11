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
	"strings"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// String is a nullable string. It supports SQL and JSON serialization.
// It will marshal to null if null. Blank string input will be considered null.
// String implements interface Argument.
type String struct {
	// Data should be called String but there is String() method or should be
	// called Value but there is a Value() method. Now it's called Data because
	// there is no Data() method. Calling it Str would be weird...
	Data  string
	Valid bool // Valid is true if String is not NULL
}

// MakeString creates a new String. Setting the second optional argument
// to false, the string will not be valid anymore, hence NULL. String
// implements interface Argument.
func MakeString(s string) String {
	return String{
		Data:  s,
		Valid: true,
	}
}

// Scan implements the Scanner interface. Approx. >2x times faster than
// database/sql.convertAssign.
func (a *String) Scan(value interface{}) (err error) {
	// stdlib		BenchmarkSQLScanner/String-4        	10000000	       117 ns/op	      80 B/op	       3 allocs/op
	// this code	BenchmarkSQLScanner/String-4        	20000000	        78.5 ns/op	      48 B/op	       2 allocs/op
	if value == nil {
		a.Data, a.Valid = "", false
		return nil
	}
	switch v := value.(type) {
	case []byte:
		a.Data = string(v) // must be copied
		a.Valid = err == nil
	case string:
		a.Data = v
		a.Valid = err == nil
	default:
		err = errors.NotSupported.Newf("[dml] Type %T not supported in String.Scan", value)
	}
	return
}

// Value implements the driver Valuer interface.
func (a String) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}
	return a.Data, nil
}

// GoString prints an optimized Go representation.
func (a String) String() string {
	if !a.Valid {
		return "null"
	}
	return a.Data
}

// GoString prints an optimized Go representation. Takes are of backticks.
// Looses the information of the private operator. That might get fixed.
func (a String) GoString() string {
	if a.Valid && strings.ContainsRune(a.Data, '`') {
		// `This is my`string`
		a.Data = strings.Join(strings.Split(a.Data, "`"), "`+\"`\"+`")
		// `This is my`+"`"+`string`
	}
	if !a.Valid {
		return "null.String{}"
	}
	return "null.MakeString(`" + a.Data + "`)"
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null String.
// It also supports unmarshalling a sql.String.
func (a *String) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(bTextNullLC, data) {
		a.Valid = false
		a.Data = ""
		return nil
	}

	var err error
	var v interface{}

	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		a.Data = x
	case map[string]interface{}:
		dto := &struct {
			String string
			Valid  bool
		}{}
		err = json.Unmarshal(data, dto)
		a.Data = dto.String
		a.Valid = dto.Valid
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[dml] json: cannot unmarshal %#v into Go value of type null.String", v)
	}
	a.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this String is null.
func (a String) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return json.Marshal(a.Data)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this String is null.
func (a String) MarshalText() ([]byte, error) {
	if !a.Valid {
		return nil, nil
	}
	return []byte(a.Data), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null String if the input is a blank string.
func (a *String) UnmarshalText(text []byte) error {
	if text != nil && !utf8.Valid(text) {
		return errors.NotValid.Newf("[dml] Input bytes are not valid UTF-8 encoded.")
	}
	a.Data = string(text)
	a.Valid = text != nil
	return nil
}

// SetValid changes this String's value and also sets it to be non-null.
func (a String) SetValid(v string) String { a.Data = v; a.Valid = true; return a }

// SetNull sets the value to Go's default value and Valid to false.
func (a String) SetNull() String { return String{} }

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null.
func (a String) Ptr() *string {
	if !a.Valid {
		return nil
	}
	return &a.Data
}

// IsZero returns true for null strings, for potential future omitempty support.
func (a String) IsZero() bool {
	return !a.Valid
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (a *String) UnmarshalBinary(data []byte) error {
	return a.Unmarshal(data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (a String) MarshalBinary() (data []byte, err error) {
	return a.Marshal()
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a String) WriteTo(d Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		if utf8.ValidString(a.Data) {
			d.EscapeString(w, a.Data)
		} else {
			err = errors.NotValid.Newf("[dml] String.writeTo: String is not UTF-8: %q", a.Data)
		}
	} else {
		_, err = w.WriteString(sqlStrNullUC)
	}
	return
}

// Append appends the value or its nil type to the interface slice.
func (a String) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Data)
	}
	return append(args, nil)
}
