// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package valuelabel

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/corestoreio/csfw/util/bufferpool"
)

var _ json.Marshaler = (*Pair)(nil)
var _ json.Unmarshaler = (*Pair)(nil)

// NotNull* are specifying which type has a non null value
const (
	NotNullString uint = iota + 1
	NotNullInt
	NotNullFloat64
	NotNullBool
)

// Pair contains a stringyfied value and a label for printing in a browser / JS api.
type Pair struct {
	// NotNull defines which type is not null
	NotNull uint
	String  string  `json:"-"`
	Int     int     `json:"-"`
	Float64 float64 `json:"-"`
	Bool    bool    `json:"-"`
	label   string
	// TODO(cs) add maybe more types and SQL connection ...
}

// Label returns the the label. A function for consistency to Value()
func (p Pair) Label() string { return p.label }

// Value returns the underlying value as a string
func (p Pair) Value() string {
	var s string
	switch p.NotNull {
	case NotNullString:
		s = p.String
	case NotNullInt:
		s = strconv.Itoa(p.Int)
	case NotNullFloat64:
		s = strconv.FormatFloat(p.Float64, 'f', 4, 64)
	case NotNullBool:
		s = strconv.FormatBool(p.Bool)
	}
	return s
}

// UnmarshalJSON decodes a pair from JSON
func (p Pair) UnmarshalJSON(data []byte) error {
	// todo Pair.UnmarshalJSON

	println(string(data))

	return json.Unmarshal(data, &p.String)
	//return nil
}

// MarshalJSON encodes a pair into JSON
func (p Pair) MarshalJSON() ([]byte, error) {
	var data []byte
	var err error
	switch p.NotNull {
	case NotNullString:
		data, err = json.Marshal(p.String)
		if err != nil {
			return nil, err
		}
	case NotNullInt:
		data, err = json.Marshal(p.Int)
		if err != nil {
			return nil, err
		}
	case NotNullFloat64:
		var n = p.Float64
		switch {
		case math.IsInf(p.Float64, 1):
			n = math.MaxFloat64
		case math.IsInf(p.Float64, -1):
			n = -math.MaxFloat64
		case math.IsNaN(p.Float64):
			n = 0.0
		}
		data, err = json.Marshal(n)
		if err != nil {
			return nil, err
		}
	case NotNullBool:
		data, err = json.Marshal(p.Bool)
		if err != nil {
			return nil, err
		}
	}

	labelData, err := json.Marshal(p.label)
	if err != nil {
		return nil, err
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)
	buf.WriteString(`{"Value":`)
	if len(data) == 0 {
		buf.WriteByte('"')
		buf.WriteByte('"')
	} else {
		buf.Write(data)
	}

	buf.WriteString(`,"Label":`)
	buf.Write(labelData)
	buf.WriteByte('}')

	return buf.Bytes(), nil
}
