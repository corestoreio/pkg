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
	}

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
		c.Number.Options(NumberFormat(f))
	}
}

func NewCurrency(opts ...CurrencyOptFunc) *Currency {
	c := new(Currency)
	CurrencyISO(DefaultCurrencyName)(c)
	c.Number = NewNumber()
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
