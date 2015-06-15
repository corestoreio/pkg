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
	"fmt"
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
		// FmtNumber formats a number according to the number format.
		// Internal rounding will be applied.
		// Returns the number bytes written or an error.
		FmtNumber(w io.Writer, n float64) (int, error)
	}

	Number struct {
		//		Tag language.Tag
		// @todo
		Symbols Symbols
		fo      format
		fneg    format // format for negative numbers
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
	var generalFormat []rune
	var negativeFormat []rune
	found := false
	for _, r := range f {
		if r == ';' && !found {
			found = true
			continue // skip semi colon :-)
		}
		if !found {
			generalFormat = append(generalFormat, r)
		} else {
			negativeFormat = append(negativeFormat, r)
		}
	}

	return func(n *Number) {
		n.fo = format{
			parsed:    false,
			pattern:   generalFormat,
			precision: 9,
			plusSign:  n.Symbols.PlusSign, // apply default values
			minusSign: n.Symbols.MinusSign,
			decimal:   n.Symbols.Decimal,
			group:     0,
		}
		n.fneg = n.fo // copy default format
		n.fneg.pattern = negativeFormat
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
	NumberFormat(`#0.######`)(n) // normally that should come from golang.org/x/text package
	//	NumberTag("en-US")(n)
	return n.Options(opts...)
}

// Options applies options and returns a number pointer
func (no *Number) Options(opts ...NumberOptFunc) *Number {
	for _, o := range opts {
		if o != nil {
			o(no)
		}
	}
	return no
}

func (no *Number) getFormat(isNegative bool) (format, error) {

	if isNegative {
		if false == no.fneg.parsed {
			if err := no.fneg.parse(); err != nil {
				return no.fneg, errgo.Mask(err)
			}
		}
		if true == no.fneg.parsed { // fneg can still be invalid because not available
			no.fneg.minusSign = 0
			return no.fneg, nil
		}
	}

	if false == no.fo.parsed { // parse positive format
		if err := no.fo.parse(); err != nil {
			return no.fo, errgo.Mask(err)
		}
	}
	return no.fo, nil
}

// FmtNumber formats a number according to the number format.
// Internal rounding will be applied.
// Returns the number bytes written or an error.
func (no *Number) FmtNumber(w io.Writer, nFloat float64) (int, error) {
	for i := range no.buf {
		no.buf[i] = 0 // clear buffer
	}

	// Special cases:
	//   NaN = "NaN"
	//   +Inf = "+Infinity"
	//   -Inf = "-Infinity"

	if math.IsNaN(nFloat) {
		return w.Write(no.Symbols.Nan)
	}

	if nFloat > math.MaxFloat64 { // that won't happen
		utf8.EncodeRune(no.buf, no.Symbols.Infinity)
		return w.Write(no.buf)
	}
	if nFloat < -math.MaxFloat64 { // that won't happen
		utf8.EncodeRune(no.buf, no.Symbols.MinusSign)
		utf8.EncodeRune(no.buf, no.Symbols.Infinity)
		return w.Write(no.buf)
	}

	usedFmt, err := no.getFormat(nFloat < -0.000000001)
	//	fmt.Println(usedFmt.String())
	if err != nil {
		return 0, log.Error("Number=FmtNumber", "err", err, "format", usedFmt.String())
	}

	var wrote int
	if nFloat > 0.000000001 && usedFmt.plusSign > 0 {
		wrote += utf8.EncodeRune(no.buf, usedFmt.plusSign)
	}
	if nFloat < -0.000000001 && usedFmt.minusSign > 0 {
		wrote += utf8.EncodeRune(no.buf, usedFmt.minusSign)
	}
	if nFloat < -0.000000001 {
		nFloat = -nFloat // convert to positive value because of the minusSign
	}

	precPow10 := math.Pow10(usedFmt.precision)

	// rounds on 0.5, 0.05, 0.005, 0.0005 ... depending on the amount of fractals in the format
	intf, fracf := math.Modf(nFloat + (5 / (precPow10 * 10)))

	// generate integer part string
	intStr := strconv.FormatInt(int64(intf), 10) // maybe convert to byte ...
	if usedFmt.group > 0 {                       // add thousand separator if required
		for i := len(intStr); i > 3; {
			i -= 3
			intStr = intStr[:i] + string(usedFmt.group) + intStr[i:]
		}
	}

	// no fractional part, we can leave now
	if usedFmt.precision == 0 {

		if usedFmt.prefix > 0 {
			wrote += utf8.EncodeRune(no.buf[wrote:], usedFmt.prefix)
			no.buf = no.buf[:numberBufferSize] // revert back to old size
		}

		no.buf = append(no.buf[:wrote], intStr...)
		no.buf = no.buf[:numberBufferSize] // revert back to old size
		wrote += len(intStr)

		if usedFmt.suffix > 0 {
			wrote += utf8.EncodeRune(no.buf[wrote:], usedFmt.suffix)
			no.buf = no.buf[:numberBufferSize] // revert back to old size
		}
		return w.Write(no.buf[:wrote])
	}

	// generate fractional part
	fracStr := strconv.FormatInt(int64(fracf*precPow10), 10)

	// may need padding
	if len(fracStr) < usedFmt.precision {
		fracStr = "000000000000000"[:usedFmt.precision-len(fracStr)] + fracStr
	}

	if usedFmt.prefix > 0 {

		wrote += utf8.EncodeRune(no.buf[wrote:], usedFmt.prefix)
		no.buf = no.buf[:numberBufferSize] // revert back to old size
	}

	no.buf = append(no.buf[:wrote], intStr...)
	wrote += len(intStr)
	no.buf = no.buf[:numberBufferSize] // revert back to old size

	wrote += utf8.EncodeRune(no.buf[wrote:], usedFmt.decimal)
	no.buf = append(no.buf[:wrote], fracStr...)
	wrote += len(fracStr)
	no.buf = no.buf[:numberBufferSize] // revert back to old size

	if usedFmt.suffix > 0 {
		wrote += utf8.EncodeRune(no.buf[wrote:], usedFmt.suffix)
	}

	return w.Write(no.buf[:wrote])
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
	prefix    rune
	suffix    rune
}

func (f *format) String() string {
	return fmt.Sprintf(
		"Parsed \t%t\nPattern\t%s\nPrec.  \t%d\nPlus\t_%s_\nMinus  \t_%s_\nDecimal\t_%s_\nGroup \t_%s_\nPrefix \t_%s_\nSuffix \t_%s_\n",
		f.parsed,
		string(f.pattern),
		f.precision,
		string(f.plusSign),
		string(f.minusSign),
		string(f.decimal),
		string(f.group),
		string(f.prefix),
		string(f.suffix),
	)
}

// Number patterns affect how numbers are interpreted in a localized context.
// Here are some examples, based on the French locale. The "." shows where the
// decimal point should go. The "," shows where the thousands separator should
// go. A "0" indicates zero-padding: if the number is too short, a zero (in
// the locale's numeric set) will go there. A "#" indicates no padding: if the
// number is too short, nothing goes there. A "¤" shows where the currency sign
// will go. The following illustrates the effects of different patterns for the
// French locale, with the number "1234.567". Notice how the pattern characters
// ',' and '.' are replaced by the characters appropriate for the locale.
// Pattern		Currency	Text
// #,##0.##		n/a			1 234,57
// #,##0.###	n/a			1 234,567
// ###0.#####	n/a			1234,567
// ###0.0000#	n/a			1234,5670
// 00000.0000	n/a			01234,5670
// #,##0.00 ¤	EUR			1 234,57 €
//				JPY			1 235 ¥JP

func (f *format) parse() error {

	if len(f.pattern) == 0 {
		return nil
	}
	f.parsed = true // only IF there is a format

	// collect indices of meaningful formatting directives
	formatDirectiveIndices := make([]int, 0)
	for i, c := range f.pattern {
		if c != '#' && c != '0' { //&& c != '\u00A4' /* ¤ */ {
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
		if formatDirectiveIndices[0] == 0 { // first character is neither # nor 0

			if f.pattern[formatDirectiveIndices[0]] != '+' {
				f.plusSign = 0 // not needed
			}

			f.prefix = f.pattern[formatDirectiveIndices[0]]
			if f.prefix != '\u00A4' { // ¤
				lastIndex := len(formatDirectiveIndices) - 1
				f.suffix = f.pattern[formatDirectiveIndices[lastIndex]]
				formatDirectiveIndices = formatDirectiveIndices[:lastIndex]
			}
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
		// @todo in some rare cases the group separator will be set after 4 digits, not 3.
		// http://unicode.org/reports/tr35/tr35-numbers.html#Number_Symbols
		//		fmt.Printf("\n%#v => %d - %d\n", formatDirectiveIndices,formatDirectiveIndices[1] , formatDirectiveIndices[0])
		if len(formatDirectiveIndices) == 2 {
			diff := (formatDirectiveIndices[1] - formatDirectiveIndices[0])
			if diff != 4 && diff != 3 {
				errF := errgo.Newf("Group separator directive must be followed by 3 digit-specifiers in format: %s", string(f.pattern))
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
			lp := len(f.pattern)
			if f.suffix > 0 {
				lp -= 1
			}
			f.precision = lp - formatDirectiveIndices[0] - 1
		}
	}
	return nil
}
