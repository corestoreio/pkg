// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package i18n_test

import (
	"bytes"
	"testing"

	"github.com/corestoreio/pkg/i18n"
)

var bmCurrencySign = []byte("€")
var benchmarkFmtCurrency string

// BenchmarkFmtCurrency_Non_Singleton_Pos	  200000	     11878 ns/op	    3473 B/op	     113 allocs/op
func BenchmarkFmtCurrency_Non_Singleton_Pos(b *testing.B) {
	bmFmtCurrency_Non_Singleton(b, "#,##0.00 ¤", "1,234.57 €", 1, 1234, 3, 567)
}

// BenchmarkFmtCurrency_Non_Singleton_Neg	  100000	     12335 ns/op	    3601 B/op	     117 allocs/op
func BenchmarkFmtCurrency_Non_Singleton_Neg(b *testing.B) {
	bmFmtCurrency_Non_Singleton(b, "¤#,##0.00;(¤#,##0.00)", "(€1,234.57)", -1, -1234, 3, 567)
}

func bmFmtCurrency_Non_Singleton(b *testing.B, format, want string, sign int, intgr int64, prec int, frac int64) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		haveC := i18n.NewCurrency(
			i18n.SetCurrencyFormat(format, testDefaultNumberSymbols),
			i18n.SetCurrencySign(bmCurrencySign),
		)
		var buf bytes.Buffer
		if _, err := haveC.FmtNumber(&buf, sign, intgr, prec, frac); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != want {
			b.Errorf("Missmatch have %s vs want %s", have, want)
		}
		benchmarkFmtCurrency = have
	}
}

// BenchmarkFmtCurrency____Singleton_Pos	 1000000	      1074 ns/op	      48 B/op	       6 allocs/op
func BenchmarkFmtCurrency____Singleton_Pos(b *testing.B) {
	bmFmtCurrency_Singleton(b, "#,##0.00 ¤", "1,234.57 €", 1, 1234, 3, 567)
}

// BenchmarkFmtCurrency____Singleton_Int	 1000000	      1072 ns/op	      48 B/op	       7 allocs/op
func BenchmarkFmtCurrency____Singleton_Int(b *testing.B) {
	bmFmtCurrency_Singleton(b, "#,##0. ¤", "1,234.00 €", 1, 1234, 2, 0) // note: currencyfraction still says 2 digits!!!
}

// BenchmarkFmtCurrency____Singleton_Neg	 1000000	      1116 ns/op	      48 B/op	       6 allocs/op
func BenchmarkFmtCurrency____Singleton_Neg(b *testing.B) {
	bmFmtCurrency_Singleton(b, "¤#,##0.00;(¤#,##0.00)", "(€1,234.57)", -1, -1234, 3, 567)
}

func bmFmtCurrency_Singleton(b *testing.B, format, want string, sign int, intgr int64, prec int, frac int64) {
	b.ReportAllocs()
	haveC := i18n.NewCurrency(
		i18n.SetCurrencyFormat(format, testDefaultNumberSymbols),
		i18n.SetCurrencySign(bmCurrencySign),
	)
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := haveC.FmtNumber(&buf, sign, intgr, prec, frac); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != want {
			b.Errorf("Missmatch have %s vs want %s", have, want)
		}
		benchmarkFmtCurrency = have
		buf.Reset()
	}
}

// BenchmarkFmtCurrencyFloat64_Non_Singleton_Pos	  100000	     12412 ns/op	    3601 B/op	     117 allocs/op
func BenchmarkFmtCurrencyFloat64_Non_Singleton_Pos(b *testing.B) {
	benchmarkFmtCurrencyFloat64_Non_Singleton(b, 123.4567*10, "€1,234.57")
}

// BenchmarkFmtCurrencyFloat64_Non_Singleton_Neg	  100000	     12544 ns/op	    3601 B/op	     117 allocs/op
func BenchmarkFmtCurrencyFloat64_Non_Singleton_Neg(b *testing.B) {
	benchmarkFmtCurrencyFloat64_Non_Singleton(b, -123.4567*10, "(€1,234.57)")
}

func benchmarkFmtCurrencyFloat64_Non_Singleton(b *testing.B, fl float64, want string) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		haveC := i18n.NewCurrency(
			i18n.SetCurrencyFormat("¤#,##0.00;(¤#,##0.00)", testDefaultNumberSymbols),
			i18n.SetCurrencySign(bmCurrencySign),
		)
		var buf bytes.Buffer
		if _, err := haveC.FmtFloat64(&buf, fl); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != want {
			b.Errorf("Missmatch have %s vs want %s", have, want)
		}
		benchmarkFmtCurrency = have
	}
}

// BenchmarkFmtCurrencyFloat64_____Singleton_Pos	 1000000	      1241 ns/op	      48 B/op	       6 allocs/op
func BenchmarkFmtCurrencyFloat64_____Singleton_Pos(b *testing.B) {
	benchmarkFmtCurrencyFloat64_Singleton(b, 123.4567*10, "#,##0.00 €", "1,234.57 €")
}

// BenchmarkFmtCurrencyFloat64_____Singleton_Neg	 1000000	      1324 ns/op	      48 B/op	       6 allocs/op
func BenchmarkFmtCurrencyFloat64_____Singleton_Neg(b *testing.B) {
	benchmarkFmtCurrencyFloat64_Singleton(b, -123.4567*10, "¤#,##0.00;(¤#,##0.00)", "(€1,234.57)")
}

func benchmarkFmtCurrencyFloat64_Singleton(b *testing.B, fl float64, format, want string) {
	b.ReportAllocs()
	haveC := i18n.NewCurrency(
		i18n.SetCurrencyFormat(format, testDefaultNumberSymbols),
		i18n.SetCurrencySign(bmCurrencySign),
	)
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := haveC.FmtFloat64(&buf, fl); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != want {
			b.Errorf("Missmatch have %s vs want %s", have, want)
		}
		benchmarkFmtCurrency = have
		buf.Reset()
	}
}
