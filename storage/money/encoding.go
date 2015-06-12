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

package money

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/corestoreio/csfw/utils/log"
)

var (
	_          json.Unmarshaler = (*Currency)(nil)
	_          json.Marshaler   = (*Currency)(nil)
	_          sql.Scanner      = (*Currency)(nil)
	nullString                  = []byte(`null`)
	quotes                      = []byte(`"`)
	colon                       = []byte(`,`)
)

type (
	JSONMarshaller interface {
		MarshalJSON(*Currency) ([]byte, error)
	}
	JSONUnmarshaller interface {
		UnmarshalJSON(*Currency, []byte) error
	}
)

// MarshalJSON generates JSON output depending on the Marshaller.
func (c Currency) MarshalJSON() ([]byte, error) {
	return c.jm.MarshalJSON(&c)
}

// UnmarshalJSON reads JSON and fills the currency struct depending on the Unmarshaller.
func (c *Currency) UnmarshalJSON(b []byte) error {
	return c.jum.UnmarshalJSON(c, b)
}

// Scan scans a value into the Currency struct. Returns an error on data loss.
// Errors will be logged.
func (c *Currency) Scan(value interface{}) error {
	// @todo quick write down without tests so add tests 8-)
	if value == nil {
		c.m, c.Valid = 0, false
		return nil
	}
	if c.guard == 0 {
		c.Option(Guard(guardi))
	}
	if c.dp == 0 {
		c.Option(Precision(dpi))
	}

	if rb, ok := value.(*sql.RawBytes); ok {
		f, err := atof64([]byte(*rb))
		if err != nil {
			return log.Error("Currency=Scan", "err", err)
		}
		c.Valid = true
		c.Setf(f)
	}
	return nil
}

func atof64(bVal []byte) (f float64, err error) {
	bVal = bytes.Replace(bVal, colon, nil, -1)
	//	s := string(bVal)
	//	s1 := strings.Replace(s, ",", "", -1)
	f, err = strconv.ParseFloat(string(bVal), 64)
	return f, err
}

// JSONNumber encodes/decodes a currency as a number string to directly use in e.g. JavaScript
type JSONNumber struct{}

// JSONLocale encodes/decodes a currency according to its locale format.
// Considers the locale if a the currency symbol is valid.
type JSONLocale struct{}

// JSONExtended encodes/decodes a currency into a JSON array: [1234.56, "€", "1.234,56 €"].
// Considers the locale if a the currency symbol is valid.
type JSONExtended struct{}

var _ JSONMarshaller = new(JSONNumber)
var _ JSONUnmarshaller = new(JSONNumber)
var _ JSONMarshaller = new(JSONLocale)
var _ JSONUnmarshaller = new(JSONLocale)
var _ JSONMarshaller = new(JSONExtended)
var _ JSONUnmarshaller = new(JSONExtended)

// MarshalJSON generates a number formatted currency string
func (je JSONNumber) MarshalJSON(c *Currency) ([]byte, error) {
	if c == nil {
		return nullString, nil
	}
	if c.Valid == false {
		return nullString, nil
	}
	return c.NumberByte(), nil
}

// UnmarshalJSON decodes a string number into the Currency.
func (je JSONNumber) UnmarshalJSON(c *Currency, b []byte) error {
	f, err := atof64(b)
	if err != nil {
		return log.Error("JSONNumber=UnmarshalJSON", "err", err, "currency", c, "bytes", b)
	}
	c.Setf(f)
	return nil
}

// MarshalJSON encodes into a locale specific quoted string
func (jl JSONLocale) MarshalJSON(c *Currency) ([]byte, error) {
	if c == nil {
		return nullString, nil
	}
	if c.Valid == false {
		return nullString, nil
	}
	var b buf
	b.Write(quotes)
	b.Write(c.Localize())
	b.Write(quotes)
	return b, nil
}

// UnmarshalJSON decodes a fully localized string into a currency struct @todo
// Considers the locale if a the currency symbol is valid.
func (jl JSONLocale) UnmarshalJSON(c *Currency, b []byte) error {
	// @todo trim currency symbol, replace thousands separator, etc ...
	return errors.New("@todo unmarshal of localized bytes")
}

// MarshalJSON encodes a currency into a JSON array: [1234.56, "€", "1.234,56 €"]
func (je JSONExtended) MarshalJSON(c *Currency) ([]byte, error) {
	if c == nil {
		return nullString, nil
	}
	if c.Valid == false {
		return nullString, nil
	}

	return nil, errors.New(`@todo encodes a currency into a JSON array: [1234.56, "€", "1.234,56 €"]`)
}

// UnmarshalJSON decodes a JSON array: [1234.56, "€", "1.234,56 €"] int a currency struct.
// Considers the locale if a the currency symbol is valid.
func (je JSONExtended) UnmarshalJSON(c *Currency, b []byte) error {
	// @todo trim currency symbol, replace thousands separator, etc ...
	return errors.New("@todo unmarshal of [3]array")
}
