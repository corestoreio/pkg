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
		// FmtNumberDec formats a number according to the number format of the
		// locale. i and dec represents a floating point
		// number. Only i can be negative. Dec must always be positive. Sign
		// must be either -1 or +1. If sign is 0 the prefix will be guessed
		// from i. If sign and i are 0 function must return ErrCannotDetectMinusSign.
		// Prec defines the precision which triggers the amount of leading zeros
		// in the decimals. Prec is a number between 0 and n. If prec is 0
		// no decimal digits will be printed.
		FmtNumber(w io.Writer, sign, prec int, i, dec int64) error
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

// numberFormatMap quick and dirty implementation
var numberFormatMap = []string{"%d", "%d%s%01d", "%d%s%02d", "%d%s%03d", "%d%s%04d", "%d%s%05d", "%d%s%06d", "%d%s%07d", "%d%s%08d", "%d%s%09d", "%d%s%010d", "%d%s%011d", "%d%s%012d", "%d%s%013d", "%d%s%014d"}

// FmtNumber formats a number according to the underlying locale @todo
// and the interface contract.
func (c *Number) FmtNumber(w io.Writer, sign, prec int, i, dec int64) error {
	if sign == 0 && i == 0 {
		return ErrCannotDetectMinusSign
	}
	if dec < 0 {
		dec *= -1
	}
	if sign < 0 && i == 0 && dec > 0 {
		w.Write(c.Symbols.MinusSign) // because Dec is always positive ...
	}

	var err error
	if prec > 0 && prec <= len(numberFormatMap) {
		_, err = fmt.Fprintf(w, numberFormatMap[prec], i, c.Symbols.Decimal, dec) // @todo remove Sprintf
	} else {
		_, err = fmt.Fprintf(w, numberFormatMap[0], i) // @todo remove Sprintf
	}
	return err
}
