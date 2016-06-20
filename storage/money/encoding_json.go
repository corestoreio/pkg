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

package money

import (
	"bytes"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/errors"
)

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

var _ Encoder = new(JSONType)
var _ Decoder = new(JSONType)

// NewJSONEncoder creates a new encoder depending on the type. Accepts either
// zero or one argument. Default encoder is JSONLocale.
func NewJSONEncoder(jts ...JSONType) Encoder {
	if len(jts) != 1 {
		return JSONLocale
	}
	return jts[0]
}

// NewJSONDecoder creates a new decoder depending on the type. Accepts either
// zero or one argument. Default decoder is JSONLocale.
func NewJSONDecoder(jts ...JSONType) Decoder {
	if len(jts) != 1 {
		return JSONLocale
	}
	return jts[0]
}

// Encode encodes a money to JSON bytes according to the defined JSONType
func (t JSONType) Encode(c *Money) ([]byte, error) {
	switch t {
	case JSONNumber:
		return jsonNumberMarshal(c)
	case JSONExtended:
		return jsonExtendedMarshal(c)
	default:
		return jsonLocaleMarshal(c)
	}
}

// Decode decodes three different currency representations into a Money struct.
func (t JSONType) Decode(c *Money, b []byte) error {
	if len(b) < 1 { // we must have a valid string
		c.m, c.Valid = 0, false
		return nil
	}
	if false == utf8.Valid(b) { // we must have a valid string
		c.m, c.Valid = 0, false
		return errors.NewNotValidf("[money] Byte slice contains invalid utf8 characters: %q", string(b))
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
			c.m, c.Valid = 0, false
			return nil
		}

		lRunes++
	}

	if isArray { // now it's an error because no colon found
		c.m, c.Valid = 0, false
		return errors.NewNotValidf(`[money] No colon found in JSON array: %q`, string(runes))
	}

	switch {
	case realNumber == lRunes: // real number e.g. -1234.56 without any other stuff
		return c.ParseFloat(string(runes))

	case posSepComma == 0 && posSepDot == 0, // no decimals but included any other stripped of character
		posSepComma == 0 && posSepDot > 0: // currency contains only a dot
		return c.ParseFloat(string(number))

	case posSepComma > 0 && posSepDot == 0: // currency contains only a comma
		for i, r := range number {
			if r == ',' {
				number[i] = '.'
			}
		}
		return c.ParseFloat(string(number))

	case posSepComma > 0 && posSepDot > 0:
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
	return errors.NewNotValidf("[money] Invalid bytes: %q => Number: %q", string(b), string(number))
}

// jsonNumberMarshal generates a number formatted currency string
func jsonNumberMarshal(c *Money) ([]byte, error) {
	if c == nil || c.Valid == false {
		return nullString, nil
	}
	return c.Ftoa(), nil
}

// jsonLocaleMarshal encodes into a locale specific quoted string
func jsonLocaleMarshal(c *Money) ([]byte, error) {
	if c == nil || c.Valid == false {
		return nullString, nil
	}
	var b bytes.Buffer
	b.WriteString(`"`)
	lb, err := c.Localize()
	if err != nil {
		return nil, errors.Wrapf(err, "[money] Currency %#v\nBytes %q", c, lb.String())
	}
	template.JSEscape(&b, lb.Bytes())
	b.WriteString(`"`)
	return b.Bytes(), err
}

// jsonExtendedMarshal encodes a currency into a JSON array: [1234.56, "€", "1.234,56 €"]
func jsonExtendedMarshal(c *Money) ([]byte, error) {
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
		return nil, errors.Wrapf(err, "[money] Currency %#v\nBytes %q", c, lb.String())
	}
	template.JSEscape(&b, lb.Bytes())

	b.WriteRune('"')
	b.WriteRune(']')
	return b.Bytes(), err
}
