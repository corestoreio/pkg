// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package text

import (
	"bytes"
	"database/sql/driver"
	"fmt"
)

// LongText avoids storing long string values in Labels, Comments, Hints ...
// in Section, Group or Field elements. A LongText can contain HTML.
// The byte slice should reduce copying long strings because we're only
// copying the slice header.
type Long []byte

func (l Long) String() string {
	return string(l)
}

// Equal returns true if r2 is equal to current Route. Does not consider
// utf8 EqualFold.
func (l Long) Equal(b []byte) bool {
	// What is the use case for EqualFold?
	return bytes.Equal(l, b)
}

func (l Long) IsEmpty() bool {
	return l == nil || len(l) == 0
}

func (l Long) Copy() Long {
	n := make([]byte, len(l), len(l))
	copy(n, l)
	return n
}

// MarshalText transforms the byte slice into a text slice.
// E.g. used in json.Marshal
func (l Long) MarshalText() (text []byte, err error) {
	// this is magic in combination with json.Marshal ;-)
	return l, nil
}

// UnmarshalText copies the data from text into the Long type.
// E.g. used in json.Unmarshal
func (l *Long) UnmarshalText(text []byte) error {
	buf := make([]byte, len(text), len(text))
	copy(buf, text)
	*l = buf
	return nil
}

// Scan implements the Scanner interface.
func (l *Long) Scan(value interface{}) error {
	*l = nil
	if value == nil {
		return nil
	}
	var buf []byte
	switch t := value.(type) {
	case []byte:
		buf = make([]byte, len(t), len(t))
		copy(buf, t)
	case string:
		buf = make([]byte, len(t), len(t))
		copy(buf, t)
	default:
		return fmt.Errorf("Cannot convert value %#v to []byte", value)
	}
	*l = buf

	return nil

}

// Value implements the driver Valuer interface.
func (l Long) Value() (driver.Value, error) {
	if l == nil {
		return nil, nil
	}
	return l, nil
}
