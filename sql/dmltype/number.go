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
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
)

type numbers interface {
	constraints.Integer | constraints.Float
}

// CSN are Comma Separated Numbers, it represents an unmerged slice of numbers.
// It also implements Text Marshalers for usage in dml.ColumnMap.Text. Strings
// will be merged and split by comma, hence CSN.
type CSN[K numbers] []K

// Scan satisfies the sql.Scanner interface for CSN. If a string starts with a
// supported split-character, this function will take that character to split
// the string.
func (l *CSN[K]) Scan(src any) error {
	var str string
	switch t := src.(type) {
	case []byte:
		str = string(t)
	case string:
		str = t
	default:
		return fmt.Errorf("[dmltype] 1648066847364 CSN.Scan Unknown type or not yet implemented: %#v", src)
	}

	// bail if only one
	if len(str) == 0 {
		*l = []K{}
		return nil
	}

	// if the first rune contains a supported comma, we take that one
	r, _ := utf8.DecodeRuneInString(str)
	csvComma := CSVComma
	if int(r) < len(supportedCommas) && supportedCommas[r] {
		csvComma = string(r)
	}

	split := strings.Split(str, csvComma)

	var values []K
	for _, s := range split {
		s = strings.TrimSpace(s)
		if s != "" {
			val, err := convToK[K](s)
			if err != nil {
				return fmt.Errorf("[dmltype] 1648068389253 CSN.Scan failed to parse %q with: %w", s, err)
			}
			values = append(values, val)
		}
	}
	*l = values

	return nil
}

func convToK[K numbers](val string) (K, error) {
	var empty K
	switch any(empty).(type) {
	case int, int64:
		f, err := strconv.ParseInt(val, 10, 64)
		return K(f), err
	case int32:
		f, err := strconv.ParseInt(val, 10, 32)
		return K(f), err
	case int16:
		f, err := strconv.ParseInt(val, 10, 16)
		return K(f), err
	case int8:
		f, err := strconv.ParseInt(val, 10, 8)
		return K(f), err

	case uint, uint64:
		f, err := strconv.ParseUint(val, 10, 64)
		return K(f), err
	case uint32:
		f, err := strconv.ParseUint(val, 10, 32)
		return K(f), err
	case uint16:
		f, err := strconv.ParseUint(val, 10, 16)
		return K(f), err
	case uint8:
		f, err := strconv.ParseUint(val, 10, 8)
		return K(f), err

	case float32:
		f, err := strconv.ParseFloat(val, 32)
		return K(f), err
	case float64:
		f, err := strconv.ParseFloat(val, 64)
		return K(f), err

	default:
		return empty, fmt.Errorf("[dmltype] 1648068025064 invalid type: %T", empty)
	}
}

func (l CSN[K]) MarshalText() (text []byte, err error) {
	return l.Bytes()
}

func (l *CSN[K]) UnmarshalText(text []byte) error {
	return l.Scan(text)
}

// Value satisfies the driver.Valuer interface for CSN.
func (l CSN[K]) Value() (driver.Value, error) {
	return l.intoBuffer().String(), nil
}

func (l CSN[K]) Bytes() ([]byte, error) {
	return l.intoBuffer().Bytes(), nil
}

func (l CSN[K]) intoBuffer() *bytes.Buffer {
	var buf bytes.Buffer
	for i, v := range l {
		if i > 0 {
			buf.WriteString(CSVComma)
		}
		fmt.Fprintf(&buf, "%v", v)
	}
	return &buf
}
