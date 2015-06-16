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
	"bytes"
	"errors"
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
const formatBufferSize = 16

var ErrCannotDetectMinusSign = errors.New("Cannot detect minus sign")

// this is quick ( 2days 8-) ) implementation and needs some refactorings

type (
	// NumberFormatter knows locale specific format properties about a currency/number
	NumberFormatter interface {
		// FmtNumber formats a number according to the number format of the
		// locale. i and dec represents a floating point
		// number. Only i can be negative. Dec must always be positive. Sign
		// must be either -1 or +1. If sign is 0 the prefix will be guessed
		// from i. If sign and i are 0 function must return ErrCannotDetectMinusSign.
		FmtNumber(w io.Writer, sign int, i, dec int64) (int, error)
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
			group:     n.Symbols.Group,
			prefix:    make([]byte, formatBufferSize),
			suffix:    make([]byte, formatBufferSize),
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
			MinusSign:              '—', // em dash \u2014
			Exponential:            'E',
			SuperscriptingExponent: '×',
			PerMille:               '‰',
			Infinity:               '∞',
			Nan:                    []byte(`NaN`),
		},
		buf: make([]byte, numberBufferSize),
	}
	NumberFormat(`#,##0.###`)(n) // normally that should come from golang.org/x/text package
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
func (no *Number) FmtNumber(w io.Writer, sign int, intgr, dec int64) (int, error) {

	switch {
	case sign == 0 && intgr == 0:
		return 0, ErrCannotDetectMinusSign
	case sign == 0 && intgr < 0:
		sign = -1
	}

	for i := range no.buf {
		no.buf[i] = 0 // clear buffer
	}

	usedFmt, err := no.getFormat(sign < 0)

	if err != nil {
		return 0, log.Error("Number=FmtNumber", "err", err, "format", usedFmt.String())
	}
	precPow10 := math.Pow10(usedFmt.precision)

	// Special cases:
	//   NaN = "NaN"
	//   +Inf = "+Infinity"
	//   -Inf = "-Infinity"
	checkNaN := float64(intgr) + (float64(dec) / precPow10)
	if math.IsNaN(checkNaN) {
		return w.Write(no.Symbols.Nan)
	}
	if checkNaN > math.MaxFloat64 {
		utf8.EncodeRune(no.buf, no.Symbols.Infinity)
		return w.Write(no.buf)
	}
	if checkNaN < -math.MaxFloat64 {
		utf8.EncodeRune(no.buf, no.Symbols.MinusSign)
		utf8.EncodeRune(no.buf, no.Symbols.Infinity)
		return w.Write(no.buf)
	}

	var wrote int
	if sign > 0 && usedFmt.plusSign > 0 {
		wrote += utf8.EncodeRune(no.buf, usedFmt.plusSign)
	}

	// generate integer part string
	intStr := strconv.FormatInt(intgr, 10) // maybe convert to byte ...
	if usedFmt.group > 0 {                 // add thousand separator if required
		for i := len(intStr); i > 3; {
			i -= 3
			intStr = intStr[:i] + string(usedFmt.group) + intStr[i:]
		}
	}

	// no fractional part, we can leave now
	if usedFmt.precision == 0 {

		if lp := len(usedFmt.prefix); lp > 0 {
			no.buf = append(no.buf[:wrote], usedFmt.prefix...)
			wrote += lp
			no.buf = no.buf[:numberBufferSize] // revert back to old size
		}

		no.buf = append(no.buf[:wrote], intStr...)
		no.buf = no.buf[:numberBufferSize] // revert back to old size
		wrote += len(intStr)

		if ls := len(usedFmt.suffix); ls > 0 {
			no.buf = append(no.buf[:wrote], usedFmt.suffix...)
			wrote += ls
			no.buf = no.buf[:numberBufferSize] // revert back to old size
		}

		return w.Write(no.buf[:wrote])
	}

	// generate fractional part, round dec it to large to fit into prec
	fracStr := strconv.FormatInt(usedFmt.decToPrec(dec), 10)

	// may need padding
	if len(fracStr) < usedFmt.precision {
		fracStr = "000000000000000"[:usedFmt.precision-len(fracStr)] + fracStr
	}

	if lp := len(usedFmt.prefix); lp > 0 {
		no.buf = append(no.buf[:wrote], usedFmt.prefix...)
		wrote += lp
		no.buf = no.buf[:numberBufferSize] // revert back to old size
	}

	no.buf = append(no.buf[:wrote], intStr...)
	wrote += len(intStr)
	no.buf = no.buf[:numberBufferSize] // revert back to old size

	// write decimal separator
	wrote += utf8.EncodeRune(no.buf[wrote:], usedFmt.decimal)
	no.buf = append(no.buf[:wrote], fracStr...)
	wrote += len(fracStr)
	no.buf = no.buf[:numberBufferSize] // revert back to old size

	// write suffix
	if ls := len(usedFmt.suffix); ls > 0 {
		no.buf = append(no.buf[:wrote], usedFmt.suffix...)
		wrote += ls
		no.buf = no.buf[:numberBufferSize] // revert back to old size
	}

	// if we have a minus sign replace the minus with the format sign
	if usedFmt.minusSign > 0 {
		var mBuf [4]byte
		mWritten := utf8.EncodeRune(mBuf[:], usedFmt.minusSign)
		wrote += mWritten - 1 // check why we need here a -1 and trim does not work
		no.buf = bytes.Replace(no.buf[:wrote], []byte(`-`), mBuf[:mWritten], 1)
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
	prefix    []byte
	suffix    []byte
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

// decToPrec adapts the fractal value of a float64 number to the precision
// Rounds the value
func (f *format) decToPrec(dec int64) int64 {
	if il := intLen(dec); il > 0 && il > f.precision {
		il10 := math.Pow10(il)
		ilf := float64(dec) / il10
		prec10 := math.Pow10(f.precision)
		decf := float64(int((ilf*prec10)+0.5)) / prec10
		decf *= prec10
		decf += 0.000000000001 // I'm lovin it 8-)
		return int64(decf)
	}
	return dec
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

	pw, sw := 0, 0 // prefixWritten, suffixWritten
	suffixStart, precStart := false, false
	frmt := make([]rune, 0)
	hasGroup, hasPlus, hasMinus := false, false, false
	precCount := 0
	for _, c := range f.pattern {
		switch c {
		case '+':
			hasPlus = true
		case '-':
			hasMinus = true
		case '#', '0', '.', ',':
			frmt = append(frmt, c)
			if false == hasGroup && c == ',' {
				hasGroup = true
			}
			if precStart {
				precCount++
			}
			if false == precStart && c == '.' {
				precStart = true
			}
			suffixStart = true
		default:
			if false == suffixStart { // prefix
				if c > 0 {
					pw += utf8.EncodeRune(f.prefix[pw:], c)
					f.prefix = f.prefix[:formatBufferSize]
				}
			} else if c > 0 { // suffix
				sw += utf8.EncodeRune(f.suffix[sw:], c)
				f.suffix = f.suffix[:formatBufferSize]
			}
		}
	}
	f.prefix = f.prefix[:pw]
	f.suffix = f.suffix[:sw]

	if false == hasGroup {
		f.group = 0
	}
	if false == hasPlus {
		f.plusSign = 0
	}
	if false == hasMinus {
		f.minusSign = 0
	}
	f.precision = precCount

	return nil
}

// intLen returns the length of a positive integer.
// 1 = 1; 10 = 2; 12345 = 5; 0 = 0; -12345 = 0
func intLen(n int64) int {
	if n < 1 {
		return 0
	}
	return int(math.Floor(math.Log10(float64(n)))) + 1
}
