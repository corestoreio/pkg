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

package cfgsource

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/errors"
)

// NotNull* are specifying which type has a non null value
const (
	NotNullString uint8 = iota + 1
	NotNullInt
	NotNullFloat64
	NotNullBool
)

// Pair contains a different typed values and a label for printing in a browser.
// Especially useful for JS API and in total for value validation.
type Pair struct {
	// NotNull defines which type is not null
	NotNull uint8
	String  string  `json:"-"`
	Int     int     `json:"-"`
	Float64 float64 `json:"-"`
	Bool    bool    `json:"-"`
	label   string
	// TODO(cs) add maybe more types and SQL connection ...
}

// Label returns the label and if empty the Value().
func (p Pair) Label() string {
	if p.label == "" {
		return p.Value()
	}
	return p.label
}

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
func (p *Pair) UnmarshalJSON(data []byte) error {

	var rawPair = struct {
		Value interface{}
		Label string
	}{}

	if err := json.Unmarshal(data, &rawPair); err != nil {
		return errors.Wrapf(err, "[source] Unmarshal: %q", string(data))
	}

	p.label = rawPair.Label

	switch vt := rawPair.Value.(type) {
	case string:
		p.NotNull = NotNullString
		p.String = vt
	case float64: // due to the interface{} above int types do not exists
		if math.Abs(vt) < float64(math.MaxInt32) && vt == float64(int64(vt)) { // is int
			p.NotNull = NotNullInt
			p.Int = int(vt)
		} else { // is float
			p.NotNull = NotNullFloat64
			p.Float64 = vt
		}
	case bool:
		p.NotNull = NotNullBool
		p.Bool = vt
	default:
		return errors.Errorf("[source] Cannot detect type for value '%s' in Pair: %#v", rawPair.Value, rawPair)
	}

	return nil
}

// MarshalJSON encodes a pair into JSON
func (p Pair) MarshalJSON() ([]byte, error) {
	var data []byte
	var err error
	switch p.NotNull {
	case NotNullString:
		data, err = json.Marshal(p.String)
		if err != nil {
			return nil, errors.Wrapf(err, "[source] String Marshal: %q", p.String)
		}
	case NotNullInt:
		data, err = json.Marshal(p.Int)
		if err != nil {
			return nil, errors.Wrapf(err, "[source] String Marshal: %v", p.Int)
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
			return nil, errors.Wrapf(err, "[source] Float Marshal: %v", n)
		}
	case NotNullBool:
		data, err = json.Marshal(p.Bool)
		if err != nil {
			return nil, errors.Wrapf(err, "[source] Bool Marshal: %q", p.Bool)
		}
	}

	labelData, err := json.Marshal(p.Label())
	if err != nil {
		return nil, errors.Wrapf(err, "[source] Label Marshal: %q", p.Label())
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)
	buf.WriteString(`{"Value":`)
	if len(data) == 0 {
		_ = buf.WriteByte('"')
		_ = buf.WriteByte('"')
	} else {
		_, _ = buf.Write(data)
	}

	_, _ = buf.WriteString(`,"Label":`)
	_, _ = buf.Write(labelData)
	_ = buf.WriteByte('}')

	return buf.Bytes(), nil
}
