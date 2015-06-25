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
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
)

var (
	_          json.Unmarshaler = (*Currency)(nil)
	_          json.Marshaler   = (*Currency)(nil)
	_          sql.Scanner      = (*Currency)(nil)
	nullString                  = []byte(`null`)
)

// ErrDecodeMissingColon can be returned on malformed JSON value when decoding a currency.
var ErrDecodeMissingColon = errors.New("No colon found in JSON array")

const (
	// JSONNumber encodes/decodes a currency as a number string to directly use
	// in e.g. JavaScript
	JSONNumber JSONType = 1 << iota
	// JSONLocale encodes/decodes a currency according to its locale format.
	// Decoding: Considers the locale if the currency symbol is valid.
	JSONLocale
	// JSONExtended encodes/decodes a currency into a JSON array:
	// [1234.56, "€", "1.234,56 €"].
	// Decoding: Considers the locale if the currency symbol is valid.
	JSONExtended
)

// JSONType defines the type of the marshaller/unmarshaller
type JSONType uint8

type (
	// JSONMarshaller interface for JSON encoding
	JSONMarshaller interface {
		// MarshalJSON encodes the currency
		MarshalJSON(*Currency) ([]byte, error)
	}
	// JSONUnmarshaller interface for JSON decoding
	JSONUnmarshaller interface {
		// UnmarshalJSON reads the bytes and decodes them into the currency
		UnmarshalJSON(*Currency, []byte) error
	}
)

// MarshalJSON generates JSON output depending on the Marshaller.
func (c Currency) MarshalJSON() ([]byte, error) {
	return c.jm.MarshalJSON(&c)
}

// UnmarshalJSON reads JSON and fills the currency struct depending on the Unmarshaller.
func (c *Currency) UnmarshalJSON(src []byte) error {
	c.applyDefaults()
	if src == nil {
		c.m, c.Valid = 0, false
		return nil
	}
	return c.jum.UnmarshalJSON(c, src)
}

// Scan scans a value into the Currency struct. Returns an error on data loss.
// Errors will be logged. Initial default settings are the guard and precision value.
func (c *Currency) Scan(src interface{}) error {
	c.applyDefaults()

	if src == nil {
		c.m, c.Valid = 0, false
		if log.IsDebug() {
			log.Debug("Currency=Scan", "case", 89, "c", c, "src", src)
		}
		return nil
	}

	if b, ok := src.([]byte); ok {
		return c.ParseFloat(string(b))
	}
	return log.Error("Currency=Scan", "err", errgo.Newf("Unsupported Type for value. Supported: []byte"), "src", src)
}

// NewJSONEncoder creates a new encoder depending on the type.
// Accepts either zero or one argument.
// Default encoder is JSONLocale
func NewJSONEncoder(jts ...JSONType) JSONMarshaller {
	if len(jts) != 1 {
		return JSONLocale
	}
	return jts[0]
}

// NewJSONDecoder creates a new decoder depending on the type.
// Accepts either zero or one argument.
// Default decoder is JSONLocale
func NewJSONDecoder(jts ...JSONType) JSONUnmarshaller {
	if len(jts) != 1 {
		return JSONLocale
	}
	return jts[0]
}

var _ JSONMarshaller = new(JSONType)
var _ JSONUnmarshaller = new(JSONType)

// MarshalJSON encodes a currency to JSON bytes according to the defined JSONType
func (t JSONType) MarshalJSON(c *Currency) ([]byte, error) {
	switch t {
	case JSONNumber:
		return jsonNumberMarshal(c)
	case JSONExtended:
		return jsonExtendedMarshal(c)
	default:
		return jsonLocaleMarshal(c)
	}
}

// UnmarshalJSON decodes three different currency representations into a currency
// struct.
func (t JSONType) UnmarshalJSON(c *Currency, b []byte) error {
	if len(b) < 1 || false == utf8.Valid(b) { // we must have a valid string
		if log.IsDebug() {
			log.Debug("JSONType=UnmarshalJSON", "case", "invalid_bytes", "c", c, "bytes", string(b))
		}
		c.m, c.Valid = 0, false
		return nil
	}

	runes := bytes.Runes(b)
	lenRunes := len(runes)
	var realNumber, isNull, lRunes, posSepComma, posSepDot int
	var isArray bool
	number := make([]rune, 0, lenRunes)
	// atm not needed because currency symbol depends on the formatter
	//symbol := make([]rune, 0, lenRunes)

	// strip quotes
	if lenRunes > 1 && runes[0] == '"' && runes[lenRunes-1] == '"' {
		runes = runes[1 : lenRunes-1]
	}
	lenRunes = len(runes)

	if 0 == lenRunes {
		if log.IsDebug() {
			log.Debug("JSONType=UnmarshalJSON", "case", "lenRunes=0", "c", c, "bytes", string(b))
		}
		c.m, c.Valid = 0, false
		return nil
	}

OuterLoop:
	for i, r := range runes {

		switch {
		case unicode.IsSpace(r):
			continue
		case r == '[':
			isArray = true // [999.0000,"$","$ 999.00"] only until the first comma will be considered.
		case unicode.IsNumber(r): // 1234.56
			number = append(number, r)
			realNumber++
		case r == '.', r == '-': // -1234.56
			number = append(number, r)
			realNumber++
		case r == ',': // -1,234.56 or -1.234,56 or -1 234,56
			if isArray { // we stop after the first colon, because then the 2nd entry starts in the array
				isArray = false
				break OuterLoop
			}
			number = append(number, r)
			//case unicode.IsLetter(r), unicode.IsSymbol(r):
			//	symbol = append(symbol, r)
		}

		if posSepComma == 0 && r == ',' { // check for first occurrence of the comma
			posSepComma = i
		}
		if posSepDot == 0 && r == '.' {
			posSepDot = i
		}

		switch unicode.ToLower(r) {
		case 'n', 'u', 'l':
			isNull++
		}

		if isNull == 4 {
			if log.IsDebug() {
				log.Debug("JSONType=UnmarshalJSON", "case", "isNull", "c", c, "bytes", string(b), "runes", string(runes))
			}
			c.m, c.Valid = 0, false
			return nil
		}

		lRunes++
	}

	if isArray { // now it's an error because no colon found
		c.m, c.Valid = 0, false
		return log.Error("JSONType=UnmarshalJSON", "err", ErrDecodeMissingColon, "bytes", string(b), "number", string(number))
	}

	// I'll keep this redundant IF cases because if the unit tests and coverage checks.
	// Once it's working we can merge them.

	if realNumber == lRunes { // real number e.g. -1234.56 without any other stuff
		return c.ParseFloat(string(runes))
	}
	if posSepComma == 0 && posSepDot == 0 { // no decimals but included any other stripped of character
		return c.ParseFloat(string(number))
	}
	if posSepComma == 0 && posSepDot > 0 { // currency contains only a dot
		return c.ParseFloat(string(number))
	}
	if posSepComma > 0 && posSepDot == 0 { // currency contains only a comma
		for i, r := range number {
			if r == ',' {
				number[i] = '.'
			}
		}
		return c.ParseFloat(string(number))
	}
	if posSepComma > 0 && posSepDot > 0 {

		replaceChar := ','           // number is 12,211,232.45 or 1,234.56
		if posSepDot < posSepComma { // number is 12.211.232,45 or 1.234,56
			replaceChar = '.'
		}

		var i int
		for i < len(number) {
			switch {
			case replaceChar == '.' && number[i] == ',':
				number[i] = '.' // replace decimal comma with a dot to create fractals
			case number[i] == replaceChar:
				number = append(number[:i], number[i+1:]...) // cut comma
				i = 0                                        // restart loop
			}
			i++
		}
		return c.ParseFloat(string(number))
	}

	c.m, c.Valid = 0, false
	return log.Error("JSONType=UnmarshalJSON", "err", errors.New("Invalid bytes"), "bytes", string(b), "number", string(number))
}

// jsonNumberMarshal generates a number formatted currency string
func jsonNumberMarshal(c *Currency) ([]byte, error) {
	if c == nil || c.Valid == false {
		return nullString, nil
	}
	return c.Ftoa(), nil
}

// jsonLocaleMarshal encodes into a locale specific quoted string
func jsonLocaleMarshal(c *Currency) ([]byte, error) {
	if c == nil || c.Valid == false {
		return nullString, nil
	}
	var b bytes.Buffer
	b.WriteString(`"`)
	lb, err := c.Localize()
	if err != nil {
		return nil, log.Error("JSONLocale=MarshalJSON", "err", err, "currency", c, "bytes", lb)
	}
	template.JSEscape(&b, lb)
	b.WriteString(`"`)
	return b.Bytes(), err
}

// jsonExtendedMarshal encodes a currency into a JSON array: [1234.56, "€", "1.234,56 €"]
func jsonExtendedMarshal(c *Currency) ([]byte, error) {
	if c == nil || c.Valid == false {
		return nullString, nil
	}
	var b bytes.Buffer
	b.WriteRune('[')
	b.Write(c.Ftoa())
	b.WriteString(`, "`)
	b.WriteString(template.JSEscapeString(string(c.Symbol())))
	b.WriteString(`", "`)
	lb, err := c.Localize()
	if err != nil {
		return nil, log.Error("JSONLocale=MarshalJSON", "err", err, "currency", c, "bytes", lb)
	}
	template.JSEscape(&b, lb)

	b.WriteRune('"')
	b.WriteRune(']')
	return b.Bytes(), err
}
