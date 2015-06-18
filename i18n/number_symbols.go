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

var (
	// defaultSymbols contains all important default characters for
	// formatting a number and or a currency.
	defaultSymbols = Symbols{
		// normally that all should come from golang.org/x/text package
		Decimal:                '.',
		Group:                  ',',
		List:                   ';',
		PercentSign:            '%',
		CurrencySign:           '¤',
		PlusSign:               '+',
		MinusSign:              '—', // em dash \u2014 ;-)
		Exponential:            'E',
		SuperscriptingExponent: '×',
		PerMille:               '‰',
		Infinity:               '∞',
		Nan:                    []byte(`NaN`),
	}

	minusSign  = []byte(`-`)
	symbolSign = []byte(`¤`)
)

// Symbols general symbols used when formatting. Some are unused because @todo
type Symbols struct {
	Decimal                rune // used
	Group                  rune // used
	List                   rune
	PercentSign            rune
	CurrencySign           rune
	PlusSign               rune // used
	MinusSign              rune // used
	Exponential            rune
	SuperscriptingExponent rune
	PerMille               rune
	Infinity               rune
	Nan                    []byte // used
}

// NewSymbols creates a new non-pointer Symbols type with the
// pre-filled default symbol table. Use arguments to override the default
// symbols.
func NewSymbols(syms ...Symbols) Symbols {
	s := defaultSymbols
	for _, sym := range syms {
		s.Merge(sym)
	}
	return s
}

// String human friendly printing. Shows also the default symbol table ;-)
func (s Symbols) String() string {
	return `Decimal					` + string(s.Decimal) + `
Group					` + string(s.Group) + `
List					` + string(s.List) + `
PercentSign				` + string(s.PercentSign) + `
CurrencySign			` + string(s.CurrencySign) + `
PlusSign				` + string(s.PlusSign) + `
MinusSign				` + string(s.MinusSign) + `
Exponential				` + string(s.Exponential) + `
SuperscriptingExponent	` + string(s.SuperscriptingExponent) + `
PerMille				` + string(s.PerMille) + `
Infinity				` + string(s.Infinity) + `
NaN						` + string(s.Nan) + `
`

}

// Merge merges one Symbols into another ignoring empty values in the argument
// Symbols struct.
func (s *Symbols) Merge(from Symbols) {
	if from.Decimal > 0 {
		s.Decimal = from.Decimal
	}
	if from.Group > 0 {
		s.Group = from.Group
	}
	if from.List > 0 {
		s.List = from.List
	}
	if from.PercentSign > 0 {
		s.PercentSign = from.PercentSign
	}
	if from.CurrencySign > 0 {
		s.CurrencySign = from.CurrencySign
	}
	if from.PlusSign > 0 {
		s.PlusSign = from.PlusSign
	}
	if from.MinusSign > 0 {
		s.MinusSign = from.MinusSign
	}
	if from.Exponential > 0 {
		s.Exponential = from.Exponential
	}
	if from.SuperscriptingExponent > 0 {
		s.SuperscriptingExponent = from.SuperscriptingExponent
	}
	if from.PerMille > 0 {
		s.PerMille = from.PerMille
	}
	if from.Infinity > 0 {
		s.Infinity = from.Infinity
	}
	if from.Nan != nil {
		s.Nan = from.Nan
	}
}
