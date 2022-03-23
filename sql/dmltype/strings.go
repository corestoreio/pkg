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

package dmltype

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"unicode/utf8"
)

const CSVComma = ","

var supportedCommas = [256]bool{
	'0': false, '1': false, '2': false, '3': false, '4': false, '5': false, '6': false, '7': false,
	'8': false, '9': false,

	'a': false, 'b': false, 'c': false, 'd': false, 'e': false, 'f': false, 'g': false, 'h': false,
	'i': false, 'j': false, 'k': false, 'l': false, 'm': false, 'n': false, 'o': false, 'p': false,
	'q': false, 'r': false, 's': false, 't': false, 'u': false, 'v': false, 'w': false, 'x': false,
	'y': false, 'z': false,

	'A': false, 'B': false, 'C': false, 'D': false, 'E': false, 'F': false, 'G': false, 'H': false,
	'I': false, 'J': false, 'K': false, 'L': false, 'M': false, 'N': false, 'O': false, 'P': false,
	'Q': false, 'R': false, 'S': false, 'T': false, 'U': false, 'V': false, 'W': false, 'X': false,
	'Y': false, 'Z': false,

	'!':  false,
	'$':  false,
	'%':  false,
	'&':  false,
	'(':  false,
	')':  false,
	'*':  false,
	'+':  false,
	',':  true,
	'-':  true,
	'.':  true,
	':':  true,
	';':  true,
	'=':  true,
	'[':  false,
	'\'': false,
	']':  false,
	'_':  true,
	'~':  true,
}

// CSV represents an unmerged slice of strings. You can use package
// slices.String for further modifications of this slice type. It also
// implements Text Marshalers for usage in dml.ColumnMap.Text.
// Strings will be merged and split by comma, hence CSV.
type CSV []string

// Scan satisfies the sql.Scanner interface for CSV. If a string starts with a
// supported split-character, this function will take that character to split
// the string.
func (l *CSV) Scan(src any) error {
	var str string
	switch t := src.(type) {
	case []byte:
		str = string(t)
	case string:
		str = t
	default:
		return fmt.Errorf("[dmltype] 1647806798161 CSV.Scan Unknown type or not yet implemented: %#v", src)
	}

	// bail if only one
	if len(str) == 0 {
		*l = []string{}
		return nil
	}

	// if the first rune contains a supported comma, we take that one
	r, _ := utf8.DecodeRuneInString(str)
	csvComma := CSVComma
	if int(r) < len(supportedCommas) && supportedCommas[r] {
		csvComma = string(r)
	}

	split := strings.Split(str, csvComma)
	split2 := split[:0]
	for _, s := range split {
		s = strings.TrimSpace(s)
		if s != "" {
			split2 = append(split2, s)
		}
	}
	*l = split2
	// append coma to the slice ... might breaks things
	return nil
}

func (l CSV) MarshalText() (text []byte, err error) {
	return l.Bytes()
}

func (l *CSV) UnmarshalText(text []byte) error {
	return l.Scan(text)
}

// Value satisfies the driver.Valuer interface for CSV.
func (l CSV) Value() (driver.Value, error) {
	d, err := l.Bytes()
	if err != nil {
		return nil, err
	}
	return string(d), nil
}

func (l CSV) Bytes() ([]byte, error) {
	return []byte(strings.Join(l, CSVComma)), nil
}
