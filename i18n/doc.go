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

/*
Package i18n supports string translations with variable substitution, CLDR pluralization,
currency, formats, language, regions and timezones.

Decimals

A decimal number, or just decimal, refers to any number written in decimal
notation, although it is more commonly used to refer to numbers that have
a fractional part separated from the integer part with a decimal separator
(e.g. 11.25) 11 is the integer part, dot the decimal separator and 25 the
fractional part.

Number Format

A format like #,##0.00;(#,##0.00) consists of two parts. The first required
part will be used for negative and positive numbers and if there is a second
part after the semi-colon then this format will be solely used for formatting
of negative numbers. More pattern details can be found:
http://unicode.org/reports/tr35/tr35-numbers.html#Number_Format_Patterns
Formatting with a format like #,##,##0.00 is currently not implemented as
too rarely used.

Number formatting

To instantiate your custom number formatter:

	nf := i18n.NewNumber(
		i18n.NumberFormat("#,##0.00;(#,##0.00)" [, Symbols{Decimal: ',' ... } ] ),
	)
	nf.FmtNumber(w io.Writer, sign int, intgr, dec int64) (int, error)

Sign can be 1 for positive number and -1 for negative.
intgr is the integer part and dec the decimal aka fractal part of your float.
Roundings will be applied if dec does not fit within the decimals specified in the
format.

There are also short hand methods for FmtInt(w io.Writer, i int) (int, error) and
FmtFloat64(w io.Writer, f float64) (int, error).

For more information read the details in the documentation of the functions and types.

Currency formatting

To instantiate your custom currency formatter:

	cf := i18n.NewCurrency(
		CurrencyISO("3-letter ISO 4217 code"),
		CurrencySign(s []byte),
		CurrencyFormat("#,##0.00 ¤" [, Symbols{Decimal: ',' ... } ] ),
		CurrencyFraction(digits, rounding, cashDigits, cashRounding int)
	)
	cf.FmtCurrency(w io.Writer, sign int, i, dec int64) (int, error)

CurrencyFraction: Digits are important when your currency has a different amount
of decimal places as specified in the format. E.g. Japanese Yen has Digits 0
despite the format is #,##0.00 ¤.

@todo: Rounding refers to the Swedish rounding and are a todo in this i18n package. Use the money.Currency type for Swedish rounding.
@todo: CashDigits and CashRounding are currently not implemented.
@todo: something like https://github.com/maximilien/i18n4go

Currency Format

The currency symbol ¤ specifies where the currency sign will be placed.

*/
package i18n
