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

/*
	@todo
	http://unicode.org/reports/tr35/tr35-numbers.html#Supplemental_Currency_Data to automatically
	set the Swedish rounding
*/

import (
	"errors"
	"math"

	"io"

	"github.com/corestoreio/csfw/i18n"
	"github.com/corestoreio/csfw/utils"
	"github.com/corestoreio/csfw/utils/log"
)

var (
	// ErrOverflow occurs on integer overflow
	ErrOverflow = errors.New("Integer Overflow")

	RoundTo = .5
	//	RoundTo  = .5 + (1 / Guardf)
	RoundToN = RoundTo * -1
)

// Interval* constants http://en.wikipedia.org/wiki/Swedish_rounding
const (
	// Interval000 no swedish rounding (default)
	Interval000 Interval = iota
	// Interval005 rounding with 0.05 intervals
	Interval005
	// Interval010 rounding with 0.10 intervals
	Interval010
	// Interval015 same as Interval010 except that 5 will be rounded down.
	// 0.45 => 0.40 or 0.46 => 0.50
	// Special case for New Zealand (a must visit!), it is up to the business
	// to decide if they will round 5¢ intervals up or down. The majority of
	// retailers follow government advice and round it down. Use then Interval015.
	// otherwise use Interval010.
	Interval015
	// Interval025 rounding with 0.25 intervals
	Interval025
	// Interval050 rounding with 0.50 intervals
	Interval050
	// Interval100 rounding with 1.00 intervals
	Interval100
	interval999
)

type (
	// Interval defines the type for the Swedish rounding.
	Interval uint8

	// Currency represents a money aka currency type to avoid rounding errors
	// with floats. Includes options for printing, Swedish rounding,
	// database scanning and JSON en/decoding.
	Currency struct {
		// m money in Guard/DP
		m int64
		// fmtCur to allow language and format specific outputs in a currency format
		fmtCur i18n.CurrencyFormatter
		// fmtNum to allow language and format specific outputs in a number format
		fmtNum i18n.NumberFormatter
		// Valid if false the internal value is NULL
		Valid bool
		// Interval defines how the swedish rounding can be applied.
		Interval Interval

		jm  JSONMarshaller
		jum JSONUnmarshaller

		guard  int64
		guardf float64
		prec   int // precision only calculated when changing dp
		dp     int64
		dpf    float64
	}

	// OptionFunc used to apply options to the Currency struct
	OptionFunc func(*Currency) OptionFunc
)

// FormatCurrency to allow language and format specific outputs in a currency format
func FormatCurrency(cf i18n.CurrencyFormatter) OptionFunc {
	// @todo later idea for those two Format* functions
	// maintain an internal cache of formatters and let the user only pass
	// the option functions from the i18n package. rethink that.
	// Another idea: maintain an internal cache depending on the store ID.
	// Another idea: opts ...i18n.CurrencyOptFunc as 2nd parameter, if first is
	//				 nil and 2nd has been provided, copy DefaultFormatterCurrency
	//				 to a new instance and apply those options.
	//				 create a Clone function for i18n formatter ...
	if cf == nil {
		cf = DefaultFormatterCurrency
	}
	return func(c *Currency) OptionFunc {
		previous := c.fmtCur
		c.fmtCur = cf
		return FormatCurrency(previous)
	}
}

// FormatNumber to allow language and format specific outputs in a number format
func FormatNumber(nf i18n.NumberFormatter) OptionFunc {
	if nf == nil {
		nf = DefaultFormatterNumber
	}
	return func(c *Currency) OptionFunc {
		previous := c.fmtNum
		c.fmtNum = nf
		return FormatNumber(previous)
	}
}

// Swedish sets the Swedish rounding
// http://en.wikipedia.org/wiki/Swedish_rounding
// Errors will be logged
func Swedish(i Interval) OptionFunc {
	if i >= interval999 {
		log.Error("Currency=SetSwedishRounding", "err", errors.New("Interval out of scope. Resetting."), "interval", i)
		i = Interval000
	}
	return func(c *Currency) OptionFunc {
		previous := c.Interval
		c.Interval = i
		return Swedish(previous)
	}
}

// CashRounding same as Swedish() option function, but:
// Rounding increment, in units of 10-digits. The default is 0, which
// means no rounding is to be done. Therefore, rounding=0 and rounding=1
// have identical behavior. Thus with fraction digits of 2 and rounding
// increment of 5, numeric values are rounded to the nearest 0.05 units
// in formatting. With fraction digits of 0 and rounding increment of
// 50, numeric values are rounded to the nearest 50.
// Possible values: 5, 10, 15, 25, 50, 100
func CashRounding(rounding int) OptionFunc {
	// somehow that feels like ... not very nice
	i := Interval000
	switch rounding {
	case 5:
		i = Interval005
	case 10:
		i = Interval010
	case 15:
		i = Interval015
	case 25:
		i = Interval025
	case 50:
		i = Interval050
	case 100:
		i = Interval100
	}

	return func(c *Currency) OptionFunc {
		var p int
		switch c.Interval {
		case Interval005:
			p = 5
		case Interval010:
			p = 10
		case Interval015:
			p = 15
		case Interval025:
			p = 25
		case Interval050:
			p = 50
		case Interval100:
			p = 100
		}
		c.Interval = i
		return CashRounding(p)
	}
}

// SetGuard sets the guard
func Guard(g int) OptionFunc {
	if g == 0 {
		g = 1
	}
	return func(c *Currency) OptionFunc {
		previous := int(c.guard)
		c.guard = int64(g)
		c.guardf = float64(g)
		return Guard(previous)
	}
}

// Precision sets the precision.
// 2 decimal places => 10^2; 3 decimal places => 10^3; x decimal places => 10^x
// If not a decimal power then falls back to the default value.
func Precision(p int) OptionFunc {
	p64 := int64(p)
	l := int64(math.Log(float64(p64)))
	if p64 != 0 && (l%2) != 0 {
		p64 = dp
	}
	prec := decimals(p64)
	if p64 == 0 { // check for division by zero
		p64 = 1
	}
	return func(c *Currency) OptionFunc {
		previous := int(c.dp)
		c.dp = p64
		c.dpf = float64(p64)
		c.prec = prec // amount of decimal digits
		return Precision(previous)
	}
}

// JSONMarshal sets a custom JSON Marshaller
func JSONMarshal(m JSONMarshaller) OptionFunc {
	// @todo not sure if this whole function is needed. jm as JSONMarshaller ... but what if we need mutexes?
	if m == nil {
		m = NewJSONEncoder()
	}
	return func(c *Currency) OptionFunc {
		previous := c.jm
		c.jm = m
		return JSONMarshal(previous)
	}
}

// JSONUnmarshal sets a custom JSON Unmmarshaller
func JSONUnmarshal(um JSONUnmarshaller) OptionFunc {
	// @todo not sure if this whole function is needed. jum as JSONUnmarshaller ... but what if we need mutexes?
	if um == nil {
		um = NewJSONDecoder()
	}
	return func(c *Currency) OptionFunc {
		previous := c.jum
		c.jum = um
		return JSONUnmarshal(previous)
	}
}

// New creates a new empty Currency struct with package default values.
// Formatter can be overridden after you have created the new type.
func New(opts ...OptionFunc) Currency {
	c := Currency{
		guard:  guard,
		guardf: guardf,
		dp:     dp,
		dpf:    dpf,
		prec:   decimals(dp),
		fmtCur: DefaultFormatterCurrency,
		fmtNum: DefaultFormatterNumber,
		jm:     DefaultJSONEncode,
		jum:    DefaultJSONDecode,
	}
	c.Option(opts...)
	return c
}

// Options besides New() also Option() can apply options to the current
// struct. It returns the last set option. More info about the returned function:
// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
func (c *Currency) Option(opts ...OptionFunc) (previous OptionFunc) {
	for _, o := range opts {
		if o != nil {
			previous = o(c)
		}
	}
	return previous
}

// Abs returns the absolute value of Currency
func (c Currency) Abs() Currency {
	if c.m < 0 {
		return c.Neg()
	}
	return c
}

// Getf gets the float64 value of money (see Raw() for int64)
func (c Currency) Getf() float64 {
	return float64(c.m) / c.dpf
}

// Geti gets value of money truncating after decimal precision (see Raw() for no truncation).
// Rounds always down
func (c Currency) Geti() int64 {
	return c.m / c.dp
}

// Dec returns the decimals
func (c Currency) Dec() int64 {
	return c.Abs().Raw() % c.dp
}

// Raw returns in int64 the value of Currency (also see Geti(), See Getf() for float64)
func (c Currency) Raw() int64 {
	return c.m
}

// Set sets the raw Currency field m
func (c Currency) Set(i int64) Currency {
	c.m = i
	c.Valid = true
	return c
}

// Setf sets a float64 into a Currency type for precision calculations
func (c Currency) Setf(f float64) Currency {
	fDPf := f * c.dpf
	r := int64(f * c.dpf)
	c.Valid = true
	return c.Set(rnd(r, fDPf-float64(r)))
}

// Sign returns:
//
//	-1 if x <  0
//	+1 if x >=  0
//
func (c Currency) Sign() int {
	if c.m < 0 {
		return -1
	}
	return 1
}

// Precision returns the amount of decimal digits
func (c Currency) Precision() int {
	return c.prec
}

// Localize for money type representation in a specific locale.
func (c Currency) Localize() ([]byte, error) {
	var bufC buf
	_, err := c.LocalizeWriter(&bufC)
	return bufC, err
}

// LocalizeWriter for money type representation in a specific locale.
// Returns the number bytes written or an error.
func (c Currency) LocalizeWriter(w io.Writer) (int, error) {
	return c.fmtCur.FmtNumber(w, c.Sign(), c.Geti(), c.Precision(), c.Dec())
}

// String for money type representation in a specific locale.
func (c Currency) String() string {
	var bufC buf
	if _, err := c.LocalizeWriter(&bufC); err != nil {
		if log.IsTrace() {
			log.Trace("Currency=String", "err", err, "c", c)
		}
		log.Error("Currency=String", "err", err, "c", c)
	}
	return string(bufC)
}

// Number prints the currency without any locale specific formatting.
// E.g. useful in JavaScript.
func (c Currency) Number() ([]byte, error) {
	var bufC buf
	_, err := c.NumberWriter(&bufC)
	return bufC, err
}

// NumberWriter prints the currency as a locale specific formatted number.
// Returns the number bytes written or an error.
func (c Currency) NumberWriter(w io.Writer) (int, error) {
	return c.fmtNum.FmtNumber(w, c.Sign(), c.Geti(), c.Precision(), c.Dec())
}

// Add adds two Currency types. Returns empty Currency on integer overflow.
// Errors will be logged and a trace is available when the level for tracing has been set.
func (c Currency) Add(d Currency) Currency {
	r := c.m + d.m
	if (r^c.m)&(r^d.m) < 0 {
		if log.IsTrace() {
			log.Trace("Currency=Add", "err", ErrOverflow, "m", c, "n", d)
		}
		log.Error("Currency=Add", "err", ErrOverflow, "m", c, "n", d)
		return New()
	}
	c.m = r
	c.Valid = true
	return c
}

// Sub subtracts one Currency type from another. Returns empty Currency on integer overflow.
// Errors will be logged and a trace is available when the level for tracing has been set.
func (c Currency) Sub(d Currency) Currency {
	r := c.m - d.m
	if (r^c.m)&^(r^d.m) < 0 {
		if log.IsTrace() {
			log.Trace("Currency=Sub", "err", ErrOverflow, "m", c, "n", d)
		}
		log.Error("Currency=Sub", "err", ErrOverflow, "m", c, "n", d)
		return New()
	}
	c.m = r
	return c
}

// Mul multiplies two Currency types. Both types must have the same precision.
func (c Currency) Mul(d Currency) Currency {
	// @todo c.m*d.m will overflow int64
	r := utils.Round(float64(c.m*d.m)/c.dpf, .5, 0)
	return c.Set(int64(r))
}

// Div divides one Currency type from another
func (c Currency) Div(d Currency) Currency {
	f := (c.guardf * c.dpf * float64(c.m)) / float64(d.m) / c.guardf
	i := int64(f)
	return c.Set(rnd(i, f-float64(i)))
}

// Mulf multiplies a Currency with a float to return a money-stored type
func (c Currency) Mulf(f float64) Currency {
	i := c.m * int64(f*c.guardf*c.dpf)
	r := i / c.guard / c.dp
	return c.Set(rnd(r, float64(i)/c.guardf/c.dpf-float64(r)))
}

// Neg returns the negative value of Currency
func (c Currency) Neg() Currency {
	if c.m != 0 {
		c.m *= -1
	}
	return c
}

// Pow is the power of Currency
func (c Currency) Pow(f float64) Currency {
	return c.Setf(math.Pow(c.Getf(), f))
}

// Swedish applies the Swedish rounding. You may set the usual options.
func (c Currency) Swedish(opts ...OptionFunc) Currency {
	c.Option(opts...)
	const (
		roundOn float64 = .5
		places  int     = 0
	)
	switch c.Interval {
	case Interval005:
		// NL, SG, SA, CH, TR, CL, IE
		// 5 cent rounding
		return c.Setf(utils.Round(c.Getf()*20, roundOn, places) / 20) // base 5
	case Interval010:
		// New Zealand & Hong Kong
		// 10 cent rounding
		// In Sweden between 1985 and 1992, prices were rounded up for sales
		// ending in 5 öre.
		return c.Setf(utils.Round(c.Getf()*10, roundOn, places) / 10)
	case Interval015:
		// 10 cent rounding, special case
		// Special case: In NZ, it is up to the business to decide if they
		// will round 5¢ intervals up or down. The majority of retailers follow
		// government advice and round it down.
		if c.m%5 == 0 {
			c.m = c.m - 1
		}
		return c.Setf(utils.Round(c.Getf()*10, roundOn, places) / 10)
	case Interval025:
		// round to quarter
		return c.Setf(utils.Round(c.Getf()*4, roundOn, places) / 4)
	case Interval050:
		// 50 cent rounding
		// The system used in Sweden from 1992 to 2010, in Norway from 1993 to 2012,
		// and in Denmark since 1 October 2008 is the following:
		// Sales ending in 1–24 öre round down to 0 öre.
		// Sales ending in 25–49 öre round up to 50 öre.
		// Sales ending in 51–74 öre round down to 50 öre.
		// Sales ending in 75–99 öre round up to the next whole Krone/krona.
		return c.Setf(utils.Round(c.Getf()*2, roundOn, places) / 2)
	case Interval100:
		// The system used in Sweden since 30 September 2010 and used in Norway since 1 May 2012.
		// Sales ending in 1–49 öre/øre round down to 0 öre/øre.
		// Sales ending in 50–99 öre/øre round up to the next whole krona/krone.
		return c.Setf(utils.Round(c.Getf()*1, roundOn, places) / 1) // ;-)
	}
	return c
}

// rnd rounds int64 remainder rounded half towards plus infinity
// trunc = the remainder of the float64 calc
// r     = the result of the int64 cal
func rnd(r int64, trunc float64) int64 {
	//fmt.Printf("RND 1 r = % v, trunc = %v RoundTo = %v\n", r, trunc, RoundTo)
	if trunc > 0 {
		if trunc >= RoundTo {
			r++
		}
	} else {
		if trunc < RoundToN {
			r--
		}
	}
	//fmt.Printf("RND 2 r = % v, trunc = %v RoundTo = %v\n", r, trunc, RoundTo)
	return r
}

// decimals returns the length of the 10^n calculation, ignoring the leading 1.
// For other int64 you must add + 1.
func decimals(dec int64) int {
	if dec < 1 {
		return 0
	}
	return int(math.Floor(math.Log10(float64(dec))))
}
