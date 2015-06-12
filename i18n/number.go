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
	"errors"
	"fmt"
	"io"

	"golang.org/x/text/language"
)

// DefaultNumber default formatter for default locale en-US
var DefaultNumber NumberFormatter

var ErrCannotDetectMinusSign = errors.New("Cannot detect minus sign")

// this is quick implementation and needs some refactorings

type (
	// NumberFormatter knows locale specific format properties about a currency/number
	NumberFormatter interface {
		// FmtNumber formats a number according to the number format of the
		// locale. i and dec represents a floating point number. Only i can be
		// negative. Sign must be either -1 or +1. IF sign is 0 the prefix
		// will be guessed from i. If sign and i are 0 function must
		// return ErrCannotDetectMinusSign.
		FmtNumber(w io.Writer, sign int, i int64, dec int64) error
	}

	Number struct {
		language.Tag
		// @todo
		Symbols Symbols
	}

	Symbols struct {
		Decimal                []byte
		Group                  []byte
		List                   []byte
		PercentSign            []byte
		CurrencySign           []byte
		PlusSign               []byte
		MinusSign              []byte
		Exponential            []byte
		SuperscriptingExponent []byte
		PerMille               []byte
		Infinity               []byte
		Nan                    []byte
	}

	NumberOptFunc func(*Number)
)

func init() {
	DefaultNumber = NewNumber()
}

var _ NumberFormatter = (*Number)(nil)

func NumberTag(locale string) NumberOptFunc {
	return func(n *Number) {
		n.Tag = language.MustParse(locale)
	}
}

func NumberSymbols(s Symbols) NumberOptFunc {
	return func(n *Number) {
		n.Symbols = s
	}
}

func NewNumber(opts ...NumberOptFunc) *Number {
	n := new(Number)
	NumberTag("en-US")(n)
	NumberSymbols(
		Symbols{
			Decimal:                []byte(`.`),
			Group:                  []byte(`,`),
			List:                   []byte(`;`),
			PercentSign:            []byte(`%`),
			CurrencySign:           []byte(`¤`), // ¤ http://en.wikipedia.org/wiki/Currency_sign_(typography)
			PlusSign:               []byte(`+`),
			MinusSign:              []byte(`-`),
			Exponential:            []byte(`E`),
			SuperscriptingExponent: []byte(`×`),
			PerMille:               []byte(`‰`),
			Infinity:               []byte(`∞`),
			Nan:                    []byte(`NaN`),
		},
	)(n)
	for _, o := range opts {
		if o != nil {
			o(n)
		}
	}
	return n
}

// FmtNumber formats a number according to the underlying locale
func (c *Number) FmtNumber(w io.Writer, sign int, i int64, dec int64) error {
	if sign == 0 && i == 0 {
		return ErrCannotDetectMinusSign
	}
	if dec < 0 {
		dec *= -1
	}
	if sign < 0 && i == 0 && dec > 0 {
		w.Write(c.Symbols.MinusSign) // because Dec is always positive ...
	}
	_, err := fmt.Fprintf(w, "%d%s%02d", i, c.Symbols.Decimal, dec) // @todo remove Sprintf
	return err

}
