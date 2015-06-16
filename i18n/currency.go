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

	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
	"golang.org/x/text/language"
)

// DefaultCurrencyName 3-letter ISO 4217 code
const DefaultCurrencyName = "XXX"

// DefaultCurrency represents the package wide default currency locale
// specific formatter.
var DefaultCurrency CurrencyFormatter

var _ CurrencyFormatter = (*Currency)(nil)

type (
	// CurrencyFormatter knows locale specific properties about a currency/number
	CurrencyFormatter interface {
		NumberFormatter
		// FmtCurrency formats a number according to the number format of the
		// locale. i and dec represents a floating point
		// number. Only i can be negative. Dec must always be positive. Sign
		// must be either -1 or +1. If sign is 0 the prefix will be guessed
		// from i. If sign and i are 0 function must return ErrCannotDetectMinusSign.
		FmtCurrency(w io.Writer, sign int, i, dec int64) (int, error)

		// Symbol returns the currency symbol
		Symbol() []byte
	}

	Currency struct {
		// @todo
		*Number
		language.Currency        // maybe one day that will get extended ...
		symbol            []byte // € or USD or ...
		// fraction see description of CurrencyFraction struct
		frac CurrencyFractions
	}

	// CurrencyFractions element contains any number of info elements.
	// No negative values allowed.
	CurrencyFractions struct {
		// Digits the minimum and maximum number of decimal digits normally
		// formatted. The default is 2. For example, in the en_US locale with
		// the default value of 2 digits, the value 1 USD would format as
		// "$1.00", and the value 1.123 USD would format as → "$1.12".
		// Having a format like #,##0.00 ¤ would result in a French locale
		// 1 234,57 € and 1 235 ¥JP. Means for Euro we have 2 digits and
		// for the Yen 0 digits. Default value is 2.
		// ⚠ Warning: Digits will override the decimal/fraction part in the
		// format string ⚠.
		Digits int
		// Rounding increment, in units of 10-digits. The default is 0, which
		// means no rounding is to be done. Therefore, rounding=0 and rounding=1
		// have identical behavior. Thus with fraction digits of 2 and rounding
		// increment of 5, numeric values are rounded to the nearest 0.05 units
		// in formatting. With fraction digits of 0 and rounding increment of
		// 50, numeric values are rounded to the nearest 50.
		// ⚠ Warning: Rounding must be applied in the package money ⚠
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

	// CurrencyOptFunc options function for Currency struct
	CurrencyOptFunc func(*Currency)
)

func init() {
	DefaultCurrency = NewCurrency()
}

// CurrencyISO parses a 3-letter ISO 4217 code and sets it to the Currency
// struct. If parsing fails errors will be logged and falls back to DefaultCurrencyName.
func CurrencyISO(cur string) CurrencyOptFunc {
	return func(c *Currency) {
		lc, err := language.ParseCurrency(cur)
		if err != nil {
			if log.IsTrace() {
				log.Trace("i18n=CurrencyISO", "err", err, "cur", cur)
			}
			log.Error("i18n=CurrencyISO", "err", err, "cur", cur)
			lc = language.MustParseCurrency(DefaultCurrencyName)
		}
		c.Currency = lc
		CurrencySymbol([]byte(lc.String()))(c)
	}
}

// CurrencySymbol sets the currency symbol
func CurrencySymbol(s []byte) CurrencyOptFunc {
	return func(c *Currency) {
		if string(c.symbol) == DefaultCurrencyName {
			s = []byte("\U0001f4b0") // money bag emoji
		}
		c.symbol = s
	}
}

// CurrencyFormat sets the currency format
func CurrencyFormat(f string) CurrencyOptFunc {
	return func(c *Currency) {
		c.Number.NOptions(NumberFormat(f))
	}
}

// CurrencyFraction sets the currency fractions. For details please
// see CurrencyFractions.
func CurrencyFraction(fr CurrencyFractions) CurrencyOptFunc {
	if fr.Digits < 0 {
		fr.Digits = 0
	}
	if fr.Rounding < 0 {
		fr.Rounding = 0
	}
	if fr.CashDigits < 0 {
		fr.CashDigits = 0
	}
	if fr.CashRounding < 0 {
		fr.CashRounding = 0
	}
	return func(c *Currency) {
		c.frac = fr
	}
}

// NewCurrency creates a new Currency pointer with default settings.
func NewCurrency(opts ...CurrencyOptFunc) *Currency {
	c := &Currency{
		Number: NewNumber(),
		frac: CurrencyFractions{
			Digits:     2,
			CashDigits: 2,
		},
	}
	CurrencyISO(DefaultCurrencyName)(c)

	for _, o := range opts {
		if o != nil {
			o(c)
		}
	}
	return c.COptions(opts...)
}

// COptions applies currency options and returns a Currency pointer
func (c *Currency) COptions(opts ...CurrencyOptFunc) *Currency {
	for _, o := range opts {
		if o != nil {
			o(c)
		}
	}
	return c
}

// FmtCurrency formats a number according to the currency format.
// Internal rounding will be applied.
// Returns the number bytes written or an error.
func (c *Currency) FmtCurrency(w io.Writer, sign int, i, dec int64) (int, error) {
	i3 := 0
	var err error

	if i3, err = c.FmtNumber(w, sign, i, dec); err != nil {
		return 0, errgo.Mask(err)
	}

	// now replace ¤ with the real symbol or 3letter ISO code

	return i3, err
}

// Symbol returns the currency symbol
func (c *Currency) Symbol() []byte { return c.symbol }
