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
	- http://unicode.org/reports/tr35/tr35-numbers.html#Supplemental_Currency_Data to automatically
	- set the Swedish rounding
	- Currency
*/

import (
	"errors"
	"io"
	"math"
	"strconv"

	"github.com/corestoreio/csfw/i18n"
	"github.com/corestoreio/csfw/utils"
	"github.com/juju/errgo"
	"golang.org/x/text/currency"
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
	Money struct {
		// m money in Guard/DP
		m int64
		// FmtCur to allow language and format specific outputs in a currency format
		FmtCur i18n.CurrencyFormatter
		// FmtNum to allow language and format specific outputs in a number format
		FmtNum i18n.NumberFormatter
		// Valid if false the internal value is NULL
		Valid bool
		// Interval defines how the swedish rounding can be applied.
		Interval Interval

		// Valuta TODO(cs) defines the currency of this money type and allows
		// comparisons and calculations with other currencies.
		Valuta currency.Currency

		Encoder // Encoder default ToJSON
		Decoder // Decoder default FromJSON

		guard  int64
		guardf float64
		prec   int // precision only calculated when changing dp
		dp     int64
		dpf    float64
	}

	// Option used to apply options to the Money struct
	Option func(*Money) Option
)

// WithSwedish sets the Swedish rounding
// http://en.wikipedia.org/wiki/Swedish_rounding
// Errors will be logged
func WithSwedish(i Interval) Option {
	if i >= interval999 {
		PkgLog.Debug("money.Swedish", "err", errors.New("Interval out of scope. Resetting."), "interval", i)
		i = Interval000
	}
	return func(c *Money) Option {
		previous := c.Interval
		c.Interval = i
		return WithSwedish(previous)
	}
}

// WithCashRounding same as Swedish() option function, but:
// Rounding increment, in units of 10-digits. The default is 0, which
// means no rounding is to be done. Therefore, rounding=0 and rounding=1
// have identical behavior. Thus with fraction digits of 2 and rounding
// increment of 5, numeric values are rounded to the nearest 0.05 units
// in formatting. With fraction digits of 0 and rounding increment of
// 50, numeric values are rounded to the nearest 50.
// Possible values: 5, 10, 15, 25, 50, 100.
// todo: refactor to use text/currency package
func WithCashRounding(rounding int) Option {
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

	return func(c *Money) Option {
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
		return WithCashRounding(p)
	}
}

// WithGuard sets the guard
func WithGuard(g int) Option {
	return func(c *Money) Option {
		previous := int(c.guard)
		c.guard, c.guardf = guard(g)
		return WithGuard(previous)
	}
}

// guard generates the guard value. Optimized to reduce allocs
func guard(g int) (int64, float64) {
	if g == 0 {
		g = 1
	}
	return int64(g), float64(g)
}

// WithPrecision sets the precision.
// 2 decimal places => 10^2; 3 decimal places => 10^3; x decimal places => 10^x
// If not a decimal power then falls back to the default value.
func WithPrecision(p int) Option {
	return func(c *Money) Option {
		previous := int(c.dp)
		c.dp, c.dpf, c.prec = precision(p)
		return WithPrecision(previous)
	}
}

// precision internal prec generator. Optimized to reduce allocs
func precision(p int) (int64, float64, int) {
	p64 := int64(p)
	l := int64(math.Log(float64(p64)))
	if p64 != 0 && (l%2) != 0 {
		p64 = int64(gDPi)
	}
	if p64 == 0 { // check for division by zero
		p64 = 1
	}
	return p64, float64(p64), decimals(p64)
}

// New creates a new empty Money struct with package default values.
// Formatter can be overridden after you have created the new type.
// Implements the interfaces: database.Scanner, driver.Valuer,
// json.Marshaller, json.Unmarshaller
func New(opts ...Option) Money {
	c := Money{}
	c.applyDefaults()
	c.Option(opts...)
	return c
}

// applyDefaults used in New() and Scan()
func (m *Money) applyDefaults() {
	if m.guard == 0 {
		m.guard, m.guardf = guard(gGuardi)
	}
	if m.dp == 0 {
		m.dp, m.dpf, m.prec = precision(gDPi)
	}
	if m.Encoder == nil {
		m.Encoder = DefaultJSONEncode
	}
	if m.Decoder == nil {
		m.Decoder = DefaultJSONDecode
	}
	if m.FmtCur == nil {
		m.FmtCur = DefaultFormatterCurrency
	}
	if m.FmtNum == nil {
		m.FmtNum = DefaultFormatterNumber
	}
	m.Interval = gSwedish
}

// Options besides New() also Option() can apply options to the current
// struct. It returns the last set option. More info about the returned function:
// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
func (m *Money) Option(opts ...Option) (previous Option) {
	for _, o := range opts {
		if o != nil {
			previous = o(m)
		}
	}
	return previous
}

// Abs returns the absolute value of Currency
func (m Money) Abs() Money {
	if m.m < 0 {
		return m.Neg()
	}
	return m
}

// Getf gets the float64 value of money (see Raw() for int64)
func (m Money) Getf() float64 {
	return float64(m.m) / m.dpf
}

// Geti gets value of money truncating after decimal precision (see Raw() for no truncation).
// Rounds always down
func (m Money) Geti() int64 {
	return m.m / m.dp
}

// Dec returns the decimals
func (m Money) Dec() int64 {
	return m.Abs().Raw() % m.dp
}

// Raw returns in int64 the value of Currency (also see Geti(), See Getf() for float64)
func (m Money) Raw() int64 {
	return m.m
}

// Set sets the raw Currency field m
func (m Money) Set(i int64) Money {
	m.m = i
	m.Valid = true
	return m
}

// Setf sets a float64 into a Currency type for precision calculations
func (m Money) Setf(f float64) Money {
	fDPf := f * m.dpf
	r := int64(f * m.dpf)
	m.Valid = true
	return m.Set(rnd(r, fDPf-float64(r)))
}

// ParseFloat transforms a string float value into a real float64 value and
// sets it. Current value will be overridden. Returns a logged error.
func (m *Money) ParseFloat(s string) error {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("money.Currency.strconv.ParseFloat", "err", err, "arg", s, "currency", m)
		}
		return errgo.Mask(err)
	}
	m.Valid = true
	*m = m.Setf(f)
	return nil
}

// Sign returns:
//
//	-1 if x <  0
//	+1 if x >=  0
//
func (m Money) Sign() int {
	if m.m < 0 {
		return -1
	}
	return 1
}

// Precision returns the amount of decimal digits
func (m Money) Precision() int {
	return m.prec
}

// Localize for money type representation in a specific locale.
func (m Money) Localize() ([]byte, error) {
	var bufC buf
	_, err := m.LocalizeWriter(&bufC)
	return bufC, err
}

// LocalizeWriter for money type representation in a specific locale.
// Returns the number bytes written or an error.
func (m Money) LocalizeWriter(w io.Writer) (int, error) {
	if false == m.Valid {
		return w.Write(gNaN)
	}
	return m.FmtCur.FmtNumber(w, m.Sign(), m.Geti(), m.Precision(), m.Dec())
}

// String for money type representation in a specific locale.
func (m Money) String() string {
	var bufC buf
	if _, err := m.LocalizeWriter(&bufC); err != nil {
		PkgLog.Debug("money.Currency.String.LocalizeWriter", "err", err, "c", m)
	}
	return string(bufC)
}

// Number prints the currency without any locale specific formatting.
// E.g. useful in JavaScript.
func (m Money) Number() ([]byte, error) {
	var bufC buf
	_, err := m.NumberWriter(&bufC)
	return bufC, err
}

// NumberWriter prints the currency as a locale specific formatted number.
// Returns the number bytes written or an error.
func (m Money) NumberWriter(w io.Writer) (int, error) {
	if false == m.Valid {
		return w.Write(gNaN)
	}
	return m.FmtNum.FmtNumber(w, m.Sign(), m.Geti(), m.Precision(), m.Dec())
}

// Symbol returns the currency symbol: €, $, AU$, CHF depending on the formatter.
func (m Money) Symbol() []byte {
	return m.FmtCur.Sign()
}

// Ftoa converts the internal floating-point number to a byte slice without
// any applied formatting.
func (m Money) Ftoa() []byte {
	return m.FtoaAppend(nil)
}

// FtoaAppend converts the internal floating-point number to a byte slice without
// any applied formatting and appends it to dst and returns the extended buffer.
func (m Money) FtoaAppend(dst []byte) []byte {
	if false == m.Valid {
		return append(dst, gNaN...)
	}
	if dst == nil {
		dst = make([]byte, 0, max(m.Precision()+4, 24))
	}
	return strconv.AppendFloat(dst, m.Getf(), 'f', m.Precision(), 64)
}

// Add adds two Currency types. Returns empty Currency on integer overflow.
// Errors will be logged and a trace is available when the level for tracing has been set.
func (m Money) Add(d Money) Money {
	r := m.m + d.m
	if (r^m.m)&(r^d.m) < 0 {
		PkgLog.Debug("money.Currency.Add.Overflow", "err", ErrOverflow, "m", m, "n", d)
		return New()
	}
	m.m = r
	m.Valid = true
	return m
}

// Sub subtracts one Currency type from another. Returns empty Currency on integer overflow.
// Errors will be logged and a trace is available when the level for tracing has been set.
func (m Money) Sub(d Money) Money {
	r := m.m - d.m
	if (r^m.m)&^(r^d.m) < 0 {
		PkgLog.Debug("money.Currency.Sub.Overflow", "err", ErrOverflow, "m", m, "n", d)
		return New()
	}
	m.m = r
	return m
}

// Mul multiplies two Currency types. Both types must have the same precision.
func (m Money) Mul(d Money) Money {
	// @todo c.m*d.m will overflow int64
	r := utils.Round(float64(m.m*d.m)/m.dpf, .5, 0)
	return m.Set(int64(r))
}

// Div divides one Currency type from another
func (m Money) Div(d Money) Money {
	f := (m.guardf * m.dpf * float64(m.m)) / float64(d.m) / m.guardf
	i := int64(f)
	return m.Set(rnd(i, f-float64(i)))
}

// Mulf multiplies a Currency with a float to return a money-stored type
func (m Money) Mulf(f float64) Money {
	i := m.m * int64(f*m.guardf*m.dpf)
	r := i / m.guard / m.dp
	return m.Set(rnd(r, float64(i)/m.guardf/m.dpf-float64(r)))
}

// Neg returns the negative value of Currency
func (m Money) Neg() Money {
	if m.m != 0 {
		m.m *= -1
	}
	return m
}

// Pow is the power of Currency
func (m Money) Pow(f float64) Money {
	return m.Setf(math.Pow(m.Getf(), f))
}

// Swedish applies the Swedish rounding. You may set the usual options.
// TODO: Consider text/currency package based on Valuta field
func (m Money) Swedish(opts ...Option) Money {
	m.Option(opts...)
	const (
		roundOn float64 = .5
		places  int     = 0
	)
	switch m.Interval {
	case Interval005:
		// NL, SG, SA, CH, TR, CL, IE
		// 5 cent rounding
		return m.Setf(utils.Round(m.Getf()*20, roundOn, places) / 20) // base 5
	case Interval010:
		// New Zealand & Hong Kong
		// 10 cent rounding
		// In Sweden between 1985 and 1992, prices were rounded up for sales
		// ending in 5 öre.
		return m.Setf(utils.Round(m.Getf()*10, roundOn, places) / 10)
	case Interval015:
		// 10 cent rounding, special case
		// Special case: In NZ, it is up to the business to decide if they
		// will round 5¢ intervals up or down. The majority of retailers follow
		// government advice and round it down.
		if m.m%5 == 0 {
			m.m = m.m - 1
		}
		return m.Setf(utils.Round(m.Getf()*10, roundOn, places) / 10)
	case Interval025:
		// round to quarter
		return m.Setf(utils.Round(m.Getf()*4, roundOn, places) / 4)
	case Interval050:
		// 50 cent rounding
		// The system used in Sweden from 1992 to 2010, in Norway from 1993 to 2012,
		// and in Denmark since 1 October 2008 is the following:
		// Sales ending in 1–24 öre round down to 0 öre.
		// Sales ending in 25–49 öre round up to 50 öre.
		// Sales ending in 51–74 öre round down to 50 öre.
		// Sales ending in 75–99 öre round up to the next whole Krone/krona.
		return m.Setf(utils.Round(m.Getf()*2, roundOn, places) / 2)
	case Interval100:
		// The system used in Sweden since 30 September 2010 and used in Norway since 1 May 2012.
		// Sales ending in 1–49 öre/øre round down to 0 öre/øre.
		// Sales ending in 50–99 öre/øre round up to the next whole krona/krone.
		return m.Setf(utils.Round(m.Getf()*1, roundOn, places) / 1) // ;-)
	}
	return m
}

// CompareTo depends on the Valuta field (TODO)
func (m Money) CompareTo(d Money) bool {
	return false
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

// max returns the max number from two numbers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
