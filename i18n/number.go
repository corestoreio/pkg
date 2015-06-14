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

package i18n

import (
	"io"
	"math"
	"strconv"

	"unicode/utf8"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
)

// DefaultNumber default formatter for default locale en-US
var DefaultNumber NumberFormatter

// numberBufferSize bytes buffer size. a number including currency sign can
// be up to 64 bytes. Some runes might need more bytes ...
const numberBufferSize = 64

// this is quick implementation and needs some refactorings

type (
	// NumberFormatter knows locale specific format properties about a currency/number
	NumberFormatter interface {
		// FmtNumberDec formats a number according to the number format of the
		// locale. i and dec represents a floating point
		// number. Only i can be negative. Dec must always be positive. Sign
		// must be either -1 or +1. If sign is 0 the prefix will be guessed
		// from i. If sign and i are 0 function must return ErrCannotDetectMinusSign.
		// Prec defines the precision which triggers the amount of leading zeros
		// in the decimals. Prec is a number between 0 and n. If prec is 0
		// no decimal digits will be printed.
		FmtNumber(w io.Writer, n float64) error
	}

	Number struct {
		//		Tag language.Tag
		// @todo
		Symbols Symbols
		f       format
		buf     []byte // size numberBufferSize
	}

	Symbols struct {
		Decimal                rune
		Group                  rune
		List                   rune
		PercentSign            rune
		CurrencySign           rune
		PlusSign               rune
		MinusSign              rune
		Exponential            rune
		SuperscriptingExponent rune
		PerMille               rune
		Infinity               rune
		Nan                    []byte
	}

	// NumberOptFunc applies options to the Number struct
	NumberOptFunc func(*Number)
)

func init() {
	DefaultNumber = NewNumber()
}

var _ NumberFormatter = (*Number)(nil)

//func NumberTag(locale string) NumberOptFunc {
//	return func(n *Number) {
//		n.Tag = language.MustParse(locale)
//	}
//}

// NumberFormat applies a format to a Number. Must be applied after you have
// set the Symbol struct.
func NumberFormat(f string) NumberOptFunc {
	return func(n *Number) {
		n.f = format{
			parsed:    false,
			pattern:   []rune(f),
			precision: 9,
			plusSign:  n.Symbols.PlusSign, // apply default values
			minusSign: n.Symbols.MinusSign,
			decimal:   n.Symbols.Decimal,
			group:     0,
		}
	}
}

func NewNumber(opts ...NumberOptFunc) *Number {
	n := &Number{
		Symbols: Symbols{
			// normally that all should come from golang.org/x/text package
			Decimal:                '.',
			Group:                  ',',
			List:                   ';',
			PercentSign:            '%',
			CurrencySign:           '¤',
			PlusSign:               '+',
			MinusSign:              '-',
			Exponential:            'E',
			SuperscriptingExponent: '×',
			PerMille:               '‰',
			Infinity:               '∞',
			Nan:                    []byte(`NaN`),
		},
		buf: make([]byte, numberBufferSize),
	}
	NumberFormat(`#,###.##`)(n) // normally that should come from golang.org/x/text package
	//	NumberTag("en-US")(n)
	for _, o := range opts {
		if o != nil {
			o(n)
		}
	}
	return n
}

/*
A function to render a number to a string based on
the following user-specified criteria:

* thousands separator
* decimal separator
* decimal precision

The format parameter tells how to render the number n.

Examples of format strings, given n = 12345.6789:

"#,###.##" => "12,345.67"
"#,###." => "12,345"
"#,###" => "12345,678"
"#\u202F###,##" => "12â€¯345,67"
"#.###,###### => 12.345,678900
"" (aka default format) => 12,345.67

The highest precision allowed is 9 digits after the decimal symbol.
There is also a version for integer number, RenderInteger(),
which is convenient for calls within template.
*/

// format some kind of cache
type format struct {
	parsed    bool
	pattern   []rune
	precision int
	plusSign  rune
	minusSign rune
	decimal   rune
	group     rune
}

func (f *format) parse() error {

	// collect indices of meaningful formatting directives
	formatDirectiveIndices := make([]int, 0)
	for i, char := range f.pattern {
		if char != '#' && char != '0' {
			formatDirectiveIndices = append(formatDirectiveIndices, i)
		}
	}

	if len(formatDirectiveIndices) > 0 {
		// Directive at index 0:
		//   Must be a '+'
		//   Raise an error if not the case
		// index: 0123456789
		//        +0.000,000
		//        +000,000.0
		//        +0000.00
		//        +0000
		if formatDirectiveIndices[0] == 0 {
			if f.pattern[formatDirectiveIndices[0]] != '+' {
				errF := errgo.Newf("invalid positive sign directive in format: %s", string(f.pattern))
				if log.IsTrace() {
					log.Trace("Number=FmtNumber", "err", errF)
				}
				return log.Error("Number=FmtNumber", "err", errF)
			}
			// positiveStr = no.Symbols.PlusSign
			formatDirectiveIndices = formatDirectiveIndices[1:]
		} else {
			f.plusSign = 0
		}

		// Two directives:
		//   First is thousands separator
		//   Raise an error if not followed by 3-digit
		// 0123456789
		// 0.000,000
		// 000,000.00
		if len(formatDirectiveIndices) == 2 {
			if (formatDirectiveIndices[1] - formatDirectiveIndices[0]) != 4 {
				errF := errgo.Newf("thousands separator directive must be followed by 3 digit-specifiers in format: %s", string(f.pattern))
				if log.IsTrace() {
					log.Trace("Number=FmtNumber", "err", errF)
				}
				return log.Error("Number=FmtNumber", "err", errF)
			}
			f.group = f.pattern[formatDirectiveIndices[0]]
			formatDirectiveIndices = formatDirectiveIndices[1:]
		}

		// One directive:
		//   Directive is decimal separator
		//   The number of digit-specifier following the separator indicates wanted prec
		// 0123456789
		// 0.00
		// 000,0000
		if len(formatDirectiveIndices) == 1 {
			f.decimal = f.pattern[formatDirectiveIndices[0]]
			f.precision = len(f.pattern) - formatDirectiveIndices[0] - 1
		}
	}
	return nil
}

func (f *format) valid() bool {
	return f.parsed && len(f.pattern) > 0
}

func (no *Number) FmtNumber(w io.Writer, nFloat float64) error {
	for i := range no.buf {
		no.buf[i] = 0 // clear buffer
	}

	// Special cases:
	//   NaN = "NaN"
	//   +Inf = "+Infinity"
	//   -Inf = "-Infinity"

	if math.IsNaN(nFloat) {
		w.Write(no.Symbols.Nan)
		return nil
	}

	if nFloat > math.MaxFloat64 {
		utf8.EncodeRune(no.buf, no.Symbols.Infinity)
		w.Write(no.buf)
		return nil
	}
	if nFloat < -math.MaxFloat64 {
		utf8.EncodeRune(no.buf, no.Symbols.MinusSign)
		utf8.EncodeRune(no.buf, no.Symbols.Infinity)
		w.Write(no.buf)
		return nil
	}

	// default format, all runes

	if false == no.f.valid() {
		if err := no.f.parse(); err != nil {
			return err
		}
	}

	var wrote int
	if nFloat > 0.000000001 && no.f.plusSign > 0 {
		wrote += utf8.EncodeRune(no.buf, no.f.plusSign)
	}
	if nFloat < -0.000000001 {
		wrote += utf8.EncodeRune(no.buf, no.f.minusSign)
		nFloat = -nFloat
	}

	precPow10 := math.Pow10(no.f.precision)

	intf, fracf := math.Modf(nFloat + (5 / (precPow10 * 10)))

	// generate integer part string
	intStr := strconv.FormatInt(int64(intf), 10) // maybe convert to byte ...

	// add thousand separator if required
	if no.f.group > 0 {
		for i := len(intStr); i > 3; {
			i -= 3
			intStr = intStr[:i] + string(no.f.group) + intStr[i:]
		}
	}

	// no fractional part, we can leave now
	if no.f.precision == 0 {
		w.Write(append(no.buf[:wrote], intStr...))
		return nil
	}

	// generate fractional part
	fracStr := strconv.FormatInt(int64(fracf*precPow10), 10)

	// may need padding
	if len(fracStr) < no.f.precision {
		fracStr = "000000000000000"[:no.f.precision-len(fracStr)] + fracStr
	}

	no.buf = append(no.buf[:wrote], intStr...)
	no.buf = no.buf[:numberBufferSize] // revert back to old size

	wPos := wrote + len(intStr)

	wPos += 0
	wPos += utf8.EncodeRune(no.buf[wPos:], no.f.decimal)
	no.buf = append(no.buf[:wPos], fracStr...)

	w.Write(no.buf)

	return nil
}
