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

	"github.com/corestoreio/cspkg/i18n"
)

var benchmarkFmtNumber string

// BenchmarkFmtNumber_Non_Singleton_Pos	  300000	      5413 ns/op	    1760 B/op	      47 allocs/op
func BenchmarkFmtNumber_Non_Singleton_Pos(b *testing.B) {
	bmFmtNumber_Non_Singleton(b, "#,###.##", "1,234.57", 1, 1234, 3, 567)
}

// BenchmarkFmtNumber_Non_Singleton_Neg	  200000	      7320 ns/op	    1888 B/op	      51 allocs/op
func BenchmarkFmtNumber_Non_Singleton_Neg(b *testing.B) {
	bmFmtNumber_Non_Singleton(b, "#,##0.00;(#,##0.00)", "(1,234.57)", -1, -1234, 3, 567)
}

func bmFmtNumber_Non_Singleton(b *testing.B, format, want string, sign int, intgr int64, prec int, frac int64) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		haveNumber := i18n.NewNumber(
			i18n.SetNumberFormat(format, testDefaultNumberSymbols),
		)
		var buf bytes.Buffer
		if _, err := haveNumber.FmtNumber(&buf, sign, intgr, prec, frac); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != want {
			b.Errorf("Missmatch %s vs %s", have, want)
		}
		benchmarkFmtNumber = have
	}
}

// BenchmarkFmtNumber____Singleton_Pos	 2000000	       722 ns/op	      24 B/op	       5 allocs/op
func BenchmarkFmtNumber____Singleton_Pos(b *testing.B) {
	bmFmtNumber_Cached(b, "#,###.##", "1,234.57", 1, 1234, 3, 567)
}

// BenchmarkFmtNumber____Singleton_Int	 3000000	       593 ns/op	      21 B/op	       4 allocs/op
func BenchmarkFmtNumber____Singleton_Int(b *testing.B) {
	bmFmtNumber_Cached(b, "#,###.", "1,234", 1, 1234, 2, 0)
}

// BenchmarkFmtNumber____Singleton_Neg	 2000000	       722 ns/op	      32 B/op	       5 allocs/op
func BenchmarkFmtNumber____Singleton_Neg(b *testing.B) {
	bmFmtNumber_Cached(b, "#,##0.00;(#,##0.00)", "(1,234.57)", -1, -1234, 3, 567)
}

func bmFmtNumber_Cached(b *testing.B, format, want string, sign int, intgr int64, prec int, frac int64) {
	b.ReportAllocs()
	haveNumber := i18n.NewNumber(
		i18n.SetNumberFormat(format, testDefaultNumberSymbols),
	)
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := haveNumber.FmtNumber(&buf, sign, intgr, prec, frac); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != want {
			b.Errorf("Missmatch %s vs %s", have, want)
		}
		benchmarkFmtNumber = have
		buf.Reset()
	}
}

// BenchmarkFmtFloat64_Non_Singleton_Pos	  300000	      5659 ns/op	    1760 B/op	      47 allocs/op
func BenchmarkFmtFloat64_Non_Singleton_Pos(b *testing.B) {
	benchmarkFmtFloat64_Non_Singleton(b, 123.4567*10, "1,234.57")
}

// BenchmarkFmtFloat64_Non_Singleton_Neg	  200000	      6279 ns/op	    1888 B/op	      51 allocs/op
func BenchmarkFmtFloat64_Non_Singleton_Neg(b *testing.B) {
	benchmarkFmtFloat64_Non_Singleton(b, -123.4567*10, "(1,234.57)")
}

func benchmarkFmtFloat64_Non_Singleton(b *testing.B, fl float64, want string) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		haveNumber := i18n.NewNumber(
			i18n.SetNumberFormat("#,##0.00;(#,##0.00)", testDefaultNumberSymbols),
		)
		var buf bytes.Buffer
		if _, err := haveNumber.FmtFloat64(&buf, fl); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != want {
			b.Errorf("Missmatch %s vs %s", have, want)
		}
		benchmarkFmtNumber = have
	}
}

// BenchmarkFmtFloat64_____Singleton_Pos	 2000000	       859 ns/op	      24 B/op	       5 allocs/op
func BenchmarkFmtFloat64_____Singleton_Pos(b *testing.B) {
	benchmarkFmtFloat64_____Singleton(b, 123.4567*10, "#,###.##", "1,234.57")
}

// BenchmarkFmtFloat64_____Singleton_Neg	 2000000	       895 ns/op	      32 B/op	       5 allocs/op
func BenchmarkFmtFloat64_____Singleton_Neg(b *testing.B) {
	benchmarkFmtFloat64_____Singleton(b, -123.4567*10, "#,##0.00;(#,##0.00)", "(1,234.57)")
}

func benchmarkFmtFloat64_____Singleton(b *testing.B, fl float64, format, want string) {
	b.ReportAllocs()
	haveNumber := i18n.NewNumber(
		i18n.SetNumberFormat(format, testDefaultNumberSymbols),
	)
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := haveNumber.FmtFloat64(&buf, fl); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != want {
			b.Errorf("Missmatch %s vs %s", have, want)
		}
		benchmarkFmtNumber = have
		buf.Reset()
	}
}
