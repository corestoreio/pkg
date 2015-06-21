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

	"sync"

	"github.com/corestoreio/csfw/utils"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
)

// DefaultNumber default formatter for default locale en-US
var DefaultNumber NumberFormatter

// DefaultNumberFormat 1,000.000
const DefaultNumberFormat = `#,##0.###`

const (

	// numberBufferSize bytes buffer size. a number including currency sign can
	// be up to 64 bytes. Some runes might need more bytes ...
	numberBufferSize = 64

	// formatBufferSize used for the buffer for prefix and suffix
	formatBufferSize = 16
	// formatSeparator differentiates between the positive and negative format.
	// Negative format is always on the second position.
	formatSeparator = ';'

	floatRoundingCorrection = 0.000000000001
	floatMax64              = math.MaxFloat64 * 0.9999999999
)

var (
	ErrCannotDetectMinusSign = errors.New("Cannot detect minus sign")
	ErrPrecIsTooShort        = errors.New("Argument precision does not match with the amount of digits in dec. Prec is too short.")
)

type (
	// NumberFormatter knows locale specific format properties about a currency/number.
	NumberFormatter interface {
		// FmtNumber formats a number according to the number format of the
		// locale. i and dec represents a floating point number splitted in their
		// integer parts. Only i can be negative. Dec must always be positive. Sign
		// must be either -1 or +1. If sign is 0 the prefix will be guessed
		// from i. If sign and i are 0 function must return ErrCannotDetectMinusSign.
		// If sign is incorrect from i, sign will be adjusted to the prefix of i.
		// Prec specifies the overall precision of dec. E.g. your number is 0.0169
		// and prec is 4 then dec would be 169. Due to the precision the formatter
		// does know to add a leading zero. If prec is shorter than the length of
		// dec then prec will be adjusted to the dec length.
		FmtNumber(w io.Writer, sign int, i int64, prec int, dec int64) (int, error)
		// FmtInt formats an integer according to the format pattern.
		FmtInt64(w io.Writer, i int64) (int, error)
		// FmtFloat64 formats a float value, does internal maybe incorrect rounding.
		FmtFloat64(w io.Writer, f float64) (int, error)
	}

	Number struct {
		//		Tag language.Tag 		@todo

		// Symbols contains all available symbols for formatting any number.
		// Struct not embedded because friendlier in IDE auto completion.
		sym  Symbols
		fo   format
		fneg format // format for negative numbers
		buf  []byte // size numberBufferSize @todo check for a possible race condition
		mu   sync.RWMutex
		// frac will only be set when we're parsing a currency format.
		// So frac will be set by the parent CurrencyFormatter.
		// The Digits in CurrencyFraction will override the precision in the
		// format if different and fracValid is true
		frac      CurrencyFractions
		fracValid bool
	}

	// NumberOptFunc applies options to the Number struct. To read more
	// about the recursion pattern:
	// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
	NumberOptFunc func(*Number) NumberOptFunc
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

// NumberSymbols sets the Symbols tables. The argument will be merged into the
// default Symbols table
func NumberSymbols(s Symbols) NumberOptFunc {
	return func(n *Number) NumberOptFunc {
		previous := n.sym
		n.sym = NewSymbols(s)
		return NumberSymbols(previous)
	}
}

// NumberFormat applies a format to a Number. If you do not have set the second
// argument Symbols (will be merge into) then the default Symbols will be used.
// Only one second argument is supported. If format is empty, fallback to the
// default format.
func NumberFormat(f string, s ...Symbols) NumberOptFunc {
	if f == "" {
		f = DefaultNumberFormat
	}
	var posFmt, negFmt []rune // positive format, negative format
	found := false
	for _, r := range f {
		if r == formatSeparator && !found { // avoid strings.Split
			found = true
			continue // skip semi colon :-)
		}
		if !found {
			posFmt = append(posFmt, r)
		} else {
			negFmt = append(negFmt, r)
		}
	}

	return func(n *Number) NumberOptFunc {
		previousF := string(n.fo.pattern)
		if len(n.fneg.pattern) > 0 {
			previousF = previousF + string(formatSeparator) + string(n.fneg.pattern)
		}
		previousS := n.sym

		if len(s) == 1 {
			n.sym.Merge(s[0])
		}
		n.fo = format{
			parsed:    false,
			pattern:   posFmt,
			precision: 9,
			plusSign:  n.sym.PlusSign, // apply default values
			minusSign: n.sym.MinusSign,
			decimal:   n.sym.Decimal,
			group:     n.sym.Group,
			prefix:    make([]byte, formatBufferSize),
			suffix:    make([]byte, formatBufferSize),
		}
		n.fneg = n.fo // copy default format
		n.fneg.pattern = negFmt
		n.fneg.isNegative = true
		if len(s) == 1 {
			return NumberFormat(previousF, previousS)
		}
		return NumberFormat(previousF)
	}
}

// NewNumber creates a new number type including the default Symbols table
// and default number format. You should only create one type and reuse the
// formatter anywhere else.
func NewNumber(opts ...NumberOptFunc) *Number {
	n := &Number{
		sym: NewSymbols(),
		buf: make([]byte, numberBufferSize),
	}
	NumberFormat(DefaultNumberFormat)(n) // normally that should come from golang.org/x/text package
	//	NumberTag("en-US")(n)
	n.NOptions(opts...)
	return n
}

// NOptions applies number options and returns the last applied previous
// option function. For more details please read here
// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
// This function is thread safe.
func (no *Number) NOptions(opts ...NumberOptFunc) (previous NumberOptFunc) {
	no.mu.Lock()
	for _, o := range opts {
		if o != nil {
			previous = o(no)
		}
	}
	no.mu.Unlock()
	return
}

// GetFormat parses the pattern depended if we have a negative value or not.
// Use this function only for debugging purposes.
// NOT Thread safe.
func (no *Number) GetFormat(isNegative bool) (format, error) {
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

// FmtNumber formats a number according to the number format. Internal rounding
// will be applied. Returns the number bytes written or an error. Thread safe.
// For more details please see the interface documentation.
func (no *Number) FmtNumber(w io.Writer, sign int, intgr int64, prec int, dec int64) (int, error) {
	no.mu.Lock()
	defer no.mu.Unlock()
	no.clearBuf()

	// first check the sign
	switch {
	case sign == 0 && intgr == 0:
		return 0, ErrCannotDetectMinusSign
	case prec < intLen(dec):
		// check for the correct value for prec. prec cannot be shorter than dec. E.g.:
		// dec = 324 and prec = 2 triggers the error because length of dec is 3.
		return 0, ErrPrecIsTooShort
	case intgr < 0:
		sign = -1
	case intgr > 0:
		sign = 1
	}

	usedFmt, err := no.GetFormat(sign < 0)
	if err != nil {
		return 0, log.Error("Number=FmtNumber", "err", err, "format", usedFmt.String())
	}

	var wrote int
	if sign > 0 && usedFmt.plusSign > 0 {
		wrote += utf8.EncodeRune(no.buf, usedFmt.plusSign)
	}

	if no.fracValid {
		// currency has different precision than in the format. e.g.: japanese Yen.
		usedFmt.precision = no.frac.Digits
	}

	dec = usedFmt.adjustDecToPrec(dec, prec)

	if usedFmt.precision == 0 {
		if sign < 0 {
			intgr -= dec // round down
		} else {
			intgr += dec // round up
		}
	}

	// remove minus prefix from intgr if format is neg and minus sign is 0
	if usedFmt.isNegative && usedFmt.minusSign == 0 && intgr < 0 {
		intgr = -intgr
	}

	// generate integer part string
	intStr := strconv.FormatInt(intgr, 10) // maybe convert to byte ...
	if usedFmt.group > 0 {                 // add thousand separator if required
		if intgr < 0 {
			intStr = intStr[1:] // skip the minus sign
		}
		gc := string(usedFmt.group)
		for i := len(intStr); i > 3; {
			i -= 3
			intStr = intStr[:i] + gc + intStr[i:]
		}
		if intgr < 0 {
			intStr = "-" + intStr // add minus sign back
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
	fracStr := strconv.FormatInt(dec, 10)

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
		no.buf = bytes.Replace(no.buf[:wrote], minusSign, mBuf[:mWritten], 1)
	}

	return w.Write(no.buf[:wrote])
}

// FmtInt64 formats an integer according to the format pattern.
// Thread safe
func (no *Number) FmtInt64(w io.Writer, i int64) (int, error) {
	sign := 1
	if i < 0 {
		sign = -sign
	}
	return no.FmtNumber(w, sign, int64(i), 0, 0)
}

// FmtFloat64 formats a float value, does internal maybe incorrect rounding.
// Thread safe
func (no *Number) FmtFloat64(w io.Writer, f float64) (int, error) {
	sign := 1
	if f < 0 {
		sign = -sign
	}

	// Special cases:
	//   NaN = "NaN"
	//   +Inf = "+Infinity"
	//   -Inf = "-Infinity"

	if math.IsNaN(f) {
		return w.Write(no.sym.Nan)
	}

	if f > floatMax64 {
		no.mu.Lock()
		defer no.mu.Unlock()
		no.clearBuf()

		wr := utf8.EncodeRune(no.buf, no.sym.Infinity)
		return w.Write(no.buf[:wr])
	}
	if f < -floatMax64 {
		no.mu.Lock()
		defer no.mu.Unlock()
		no.clearBuf()

		wr := utf8.EncodeRune(no.buf, no.sym.MinusSign)
		wr += utf8.EncodeRune(no.buf[wr:], no.sym.Infinity)
		no.buf = no.buf[:numberBufferSize]
		return w.Write(no.buf[:wr])
	}

	if isInt(f) { // check if float is integer value
		return no.FmtInt64(w, int64(f))
	}

	usedFmt, err := no.GetFormat(sign < 0)
	if err != nil {
		return 0, log.Error("Number=FmtFloat64", "err", err, "format", usedFmt.String())
	}

	// to test the next lines: http://play.golang.org/p/L0ykFv3G4B
	precPow10 := math.Pow10(usedFmt.precision)

	var modf float64
	if f > 0 {
		modf = f + (5 / (precPow10 * 10))
	} else {
		modf = f - (5 / (precPow10 * 10))
	}
	intgr, fracf := math.Modf(modf)

	if fracf < 0 {
		fracf = -fracf
	}

	fracI := int64(utils.Round(fracf*precPow10, 0, usedFmt.precision))

	return no.FmtNumber(w, sign, int64(intgr), intLen(fracI), fracI)
}

func (no *Number) clearBuf() {
	no.buf = no.buf[:numberBufferSize]
	for i := range no.buf {
		no.buf[i] = 0x0 // clear buffer
	}
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

@todo rounding increment in a pattern http://unicode.org/reports/tr35/tr35-numbers.html#Rounding
@todo support more special chars: http://unicode.org/reports/tr35/tr35-numbers.html#Special_Pattern_Characters

*/

// format contains the pattern and acts as a cache
type format struct {
	isNegative bool
	parsed     bool
	pattern    []rune
	precision  int
	plusSign   rune
	minusSign  rune
	decimal    rune
	group      rune
	prefix     []byte
	suffix     []byte
}

// String human friendly printed format for debugging purposes.
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

// decToPrec adapts the fractal value of a float64 number to the format precision
// Rounds the value
func (f *format) adjustDecToPrec(dec int64, prec int) int64 {

	if f.precision > prec {
		// Moving dec values to the correct precision.
		// Edge case when format has a higher precision than prec.
		// E.G.: Format is #,##0.000 and prec=2 and dec=8 (1234.08)
		// the re-calculated dec is then 8*(10^2) = 80 to move
		// 8 to the second place. The new number would be then 1234.080 because
		// the format requires to have 3 decimal digits
		dec *= int64(math.Pow10(f.precision - prec))
	}

	// if the prec is higher than the formatted precision then we have to round
	// the dec value to fit into the precision of the format.
	if prec > 0 && prec > f.precision {
		il10 := math.Pow10(prec)
		ilf := float64(dec) / il10
		prec10 := math.Pow10(f.precision)
		decf := float64((ilf*prec10)+0.55) / prec10 // hmmm that .55 needs to be monitored. everywhere else we have just .5
		decf *= prec10
		decf += floatRoundingCorrection // I'm lovin it 8-)
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
	hasGroup, hasPlus, hasMinus := false, false, false
	precCount := 0
	for _, c := range f.pattern {
		switch c {
		case '+':
			hasPlus = true
		case '-':
			hasMinus = true
		case '#', '0', '.', ',':
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

// isInt checks if float value has no decimals. This function should run
// after checking NaN, minFloat64 and maxFloat64 OR fix it here ;-)
func isInt(f float64) bool {
	return int64(math.Floor(f)) == int64(math.Ceil(f))
}
