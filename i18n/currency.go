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
	"io"
	"sync"

	"github.com/corestoreio/csfw/utils/log"
	"golang.org/x/text/currency"
)

// DefaultCurrencyName 3-letter ISO 4217 code
const DefaultCurrencyName = "XXX"

// DefaultCurrencyFormat Symbol no-breaking-space 1,000.00
const DefaultCurrencyFormat = "¤\u00a0#,##0.00"

// DefaultCurrencySign is the classical Dollar $
var DefaultCurrencySign = []byte("$")

// DefaultCurrency represents the package wide default currency locale
// specific formatter.
var DefaultCurrency CurrencyFormatter

const currencyBufferSize = numberBufferSize + (2 * formatBufferSize)

var _ CurrencyFormatter = (*Currency)(nil)

type (
	// CurrencyFormatter knows locale specific properties about a currency/number
	CurrencyFormatter interface {
		// NumberFormatter functions for formatting. Please see interface description
		// about the contract.
		NumberFormatter
		// Sign returns the currency sign. Can be one character or a 3-letter ISO 4217 code.
		Sign() []byte
	}

	// Currency represents a locale specific formatter for a currency.
	// You must use NewCurrency() to create a new type.
	Currency struct {
		*Number
		// ISO contains the 3-letter ISO 4217 currency code.
		// Maybe one day that will get extended in text/currency package ...
		ISO currency.Currency
		sgn []byte // € or USD or ...
		buf buf
		mu  sync.RWMutex
	}

	// CurrencyFractions element contains any number of info elements.
	// No negative values allowed ⚠.
	CurrencyFractions struct {
		// Digits the minimum and maximum number of decimal digits normally
		// formatted. The default is 2. For example, in the en_US locale with
		// the default value of 2 digits, the value 1 USD would format as
		// "$1.00", and the value 1.123 USD would format as → "$1.12".
		// Having a format like #,##0.00 ¤ would result in a French locale
		// 1 234,57 € and 1 235 ¥JP. Means for Euro we have 2 digits and
		// for the Yen 0 digits. Default value is 2.
		// ⚠ Warning: Digits will override the fraction part in the
		// format string, if Digits is > 0 ⚠.
		Digits int
		// Rounding increment, in units of 10-digits. The default is 0, which
		// means no rounding is to be done. Therefore, rounding=0 and rounding=1
		// have identical behavior. Thus with fraction digits of 2 and rounding
		// increment of 5, numeric values are rounded to the nearest 0.05 units
		// in formatting. With fraction digits of 0 and rounding increment of
		// 50, numeric values are rounded to the nearest 50.
		// ⚠ Warning: Rounding must be applied in the package money ⚠ @todo
		Rounding int
		// CashDigits the number of decimal digits to be used when formatting
		// quantities used in cash transactions (as opposed to a quantity that
		// would appear in a more formal setting, such as on a bank statement).
		// If absent, the value of "digits" should be used as a default.
		// Default value is 2. @todo
		CashDigits int
		// CashRounding increment, in units of 10-cashDigits. The default is 0,
		// which means no rounding is to be done; and as with rounding, this
		// has the same effect as cashRounding="1". This is the rounding increment
		// to be used when formatting quantities used in cash transactions (as
		// opposed to a quantity that would appear in a more formal setting,
		// such as on a bank statement). If absent, the value of "rounding"
		// should be used as a default. @todo
		CashRounding int
	}

	// CurrencyOptions applies options to the Currency struct. To read more
	// about the recursion pattern:
	// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
	CurrencyOptions func(*Currency) CurrencyOptions
)

func init() {
	DefaultCurrency = NewCurrency()
}

// SetCurrencyISO parses a 3-letter ISO 4217 code and sets it to the Currency
// struct. If parsing fails errors will be logged and falls back to DefaultCurrencyName.
// Calling this function sets also the CurrencySign() to the at the moment
// 3-letter ISO code. (Missing feature in text/currency package)
// This function is called in NewCurrency().
func SetCurrencyISO(cur string) CurrencyOptions {
	return func(c *Currency) CurrencyOptions {
		previous := c.ISO.String()
		lc, err := currency.ParseISO(cur)
		if err != nil {
			if log.IsTrace() {
				log.Trace("i18n.CurrencyISO.ParseCurrency.error", "err", err, "cur", cur)
			}
			log.Error("i18n.CurrencyISO.ParseCurrency", "err", err, "cur", cur)
			lc = currency.MustParseISO(DefaultCurrencyName)
		}
		c.ISO = lc
		SetCurrencySign([]byte(lc.String()))(c)
		return SetCurrencyISO(previous)
	}
}

// SetCurrencySign sets the currency sign.
func SetCurrencySign(s []byte) CurrencyOptions {
	if string(s) == DefaultCurrencyName || len(s) == 0 {
		s = DefaultCurrencySign
	}
	return func(c *Currency) CurrencyOptions {
		previous := c.sgn
		c.sgn = s
		return SetCurrencySign(previous)
	}
}

// SetCurrencySymbols sets the Symbols tables. The argument will be merged into the
// default Symbols table
func SetCurrencySymbols(s Symbols) CurrencyOptions {
	return func(c *Currency) CurrencyOptions {
		previous := c.sym
		c.sym = NewSymbols(s)
		return SetCurrencySymbols(previous)
	}
}

// SetCurrencyFormat applies a format (e.g.: #,##0.00 ¤) to a Number.
// If you do not have set the second argument Symbols (will be merge into) then the
// default Symbols will be used. Only one second argument is supported. If format is
// empty, fallback to the default format.
// A "¤" shows where the currency sign will go.
func SetCurrencyFormat(f string, s ...Symbols) CurrencyOptions {
	if f == "" {
		f = DefaultCurrencyFormat
	}
	return func(c *Currency) CurrencyOptions {
		previousF := string(c.fo.pattern)
		if len(c.fneg.pattern) > 0 {
			previousF = previousF + string(formatSeparator) + string(c.fneg.pattern)
		}
		previousS := c.sym

		c.NSetOptions(SetNumberFormat(f, s...))
		if len(s) == 1 {
			return SetCurrencyFormat(previousF, previousS)
		}
		return SetCurrencyFormat(previousF)
	}
}

// SetCurrencyFraction sets the currency fractions. For details please
// see CurrencyFractions. Values below 0 will be reset to zero.
func SetCurrencyFraction(digits, rounding, cashDigits, cashRounding int) CurrencyOptions {
	if digits < 0 {
		digits = 0
	}
	if rounding < 0 {
		rounding = 0
	}
	if cashDigits < 0 {
		cashDigits = 0
	}
	if cashRounding < 0 {
		cashRounding = 0
	}
	return func(c *Currency) CurrencyOptions {
		prevD := c.frac.Digits
		prevR := c.frac.Rounding
		prevCD := c.frac.CashDigits
		prevCR := c.frac.CashRounding
		c.frac = CurrencyFractions{
			Digits:       digits,
			Rounding:     rounding,
			CashDigits:   cashDigits,
			CashRounding: cashRounding,
		}
		c.fracValid = true
		return SetCurrencyFraction(prevD, prevR, prevCD, prevCR)
	}
}

// NewCurrency creates a new Currency pointer with default settings.
// To change the symbols depending on the locale: see documentation.
func NewCurrency(opts ...CurrencyOptions) *Currency {
	c := &Currency{
		Number: NewNumber(),
		buf:    make(buf, currencyBufferSize),
	}
	SetCurrencyISO(DefaultCurrencyName)(c)
	SetCurrencyFormat(DefaultCurrencyFormat)(c)
	SetCurrencyFraction(2, 0, 2, 0)(c)
	c.CSetOptions(opts...)
	return c
}

// CSetOptions applies currency options and returns the last applied previous
// option function. For more details please read here
// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
// This function is thread safe.
func (c *Currency) CSetOptions(opts ...CurrencyOptions) (previous CurrencyOptions) {
	c.mu.Lock()
	for _, o := range opts {
		if o != nil {
			previous = o(c)
		}
	}
	c.mu.Unlock()
	return
}

// FmtNumber formats a number according to the currency format. Internal rounding
// will be applied. Returns the number bytes written or an error. Thread safe.
// For more details please see the interface documentation.
func (c *Currency) FmtNumber(w io.Writer, sign int, intgr int64, prec int, frac int64) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clearBuf()

	if _, err := c.Number.FmtNumber(&c.buf, sign, intgr, prec, frac); err != nil {
		return 0, log.Error("i18n.Currency.FmtNumber.FmtNumber", "err", err, "buffer", string(c.buf), "sign", sign, "i", intgr, "prec", prec, "frac", frac)
	}
	return c.flushBuf(w)
}

// FmtInt64 formats an integer according to the currency format pattern.
// Thread safe
func (c *Currency) FmtInt64(w io.Writer, i int64) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clearBuf()

	if _, err := c.Number.FmtInt64(&c.buf, i); err != nil {
		return 0, log.Error("i18n.Currency.FmtInt64.FmtInt64", "err", err, "buffer", string(c.buf), "int", i)
	}
	return c.flushBuf(w)
}

// FmtCurrencyFloat64 formats a float value, does internal maybe incorrect rounding.
// Returns the number bytes written or an error. Thread safe.
func (c *Currency) FmtFloat64(w io.Writer, f float64) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clearBuf()

	if _, err := c.Number.FmtFloat64(&c.buf, f); err != nil {
		return 0, log.Error("i18n.Currency.FmtFloat64.FmtFloat64", "err", err, "buffer", string(c.buf), "float64", f)
	}
	return c.flushBuf(w)
}

// flushBuf replaces the typographical symbol sign with the real sign.
func (c *Currency) flushBuf(w io.Writer) (int, error) {
	// now replace ¤ with the real symbol or what ever
	c.buf = bytes.Replace(c.buf, symbolSign, c.sgn, 1)
	return w.Write(c.buf)
}

// Sign returns the currency sign. Can be one character or a 3-letter ISO 4217 code.
func (c *Currency) Sign() []byte { return c.sgn }

// clearBuf iterates over the buffer and sets each element to 0 and resizes
// the buffer to length 0. Must be protected by a mutex.
func (c *Currency) clearBuf() {
	for i := range c.buf {
		c.buf[i] = 0 // clear buffer
	}
	c.buf = c.buf[:0]
}
