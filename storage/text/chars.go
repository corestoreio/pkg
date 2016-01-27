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
	"unicode/utf8"
)

// Chars avoids storing long string values in Labels, Comments, Hints ...
// in Section, Group, Field elements or etc. A Chars can contain HTML.
// The byte slice should reduce copying long strings because we're only
// copying the slice header. Danger! Other function may modify this
// slice and we cannot yet prevent it. Use it wisely and also copy the data
// away if you want to modify it.
type Chars []byte

func (c Chars) String() string {
	return string(c)
}

// Bytes converts the Chars to a byte slice ;-)
func (c Chars) Bytes() []byte {
	return c
}

// Equal returns true if r2 is equal to current Route. Does not consider
// utf8 EqualFold.
func (c Chars) Equal(b []byte) bool {
	// What is the use case for EqualFold?
	return bytes.Equal(c, b)
}

func (c Chars) IsEmpty() bool {
	return c == nil || len(c) == 0
}

// Clone returns a new allocated slice with copied data.
func (c Chars) Clone() Chars {
	n := make([]byte, len(c), len(c))
	copy(n, c)
	return n
}

// RuneCount counts the number of runes in this byte slice.
// len(Chars) and RuneCount can return different results.
func (c Chars) RuneCount() int {
	return utf8.RuneCount(c)
}

// MarshalText transforms the byte slice into a text slice.
// E.g. used in json.Marshal
func (c Chars) MarshalText() (text []byte, err error) {
	// this is magic in combination with json.Marshal ;-)
	return c, nil
}

// UnmarshalText copies the data from text into the Chars type.
// E.g. used in json.Unmarshal
func (c *Chars) UnmarshalText(text []byte) error {
	buf := make([]byte, len(text), len(text))
	copy(buf, text)
	*c = buf
	return nil
}

// Scan implements the sql.Scanner interface.
func (c *Chars) Scan(value interface{}) error {
	*c = nil
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
	*c = buf

	return nil

}

// Value implements the driver.Valuer interface.
func (c Chars) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return c.Bytes(), nil
}

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

// Hash returns a new 64-bit FNV-1a hash.
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Hash implements FNV-1 and FNV-1a, non-cryptographic hash functions
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See
// http://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function.
func (r Chars) Hash() uint64 {
	var hash uint64 = offset64
	for _, c := range r {
		hash ^= uint64(c)
		hash *= prime64
	}
	return hash
}
