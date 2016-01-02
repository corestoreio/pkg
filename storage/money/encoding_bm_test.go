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

package money_test

import (
	"strconv"
	"testing"

	"github.com/corestoreio/csfw/storage/money"
)

var benchmarkMoneyScan float64

// Benchmark_MoneyScan	 3000000	       504 ns/op	     136 B/op	       2 allocs/op => Go 1.4.2
// Benchmark_MoneyScan 	 5000000	       386 ns/op	     144 B/op	       2 allocs/op => Go 1.5.0
func Benchmark_MoneyScan(b *testing.B) {
	var d interface{}
	d = []byte{0x37, 0x30, 0x35, 0x2e, 0x39, 0x39, 0x33, 0x33}
	var want float64 = 705.993300
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var c money.Money
		c.Scan(d)
		benchmarkMoneyScan = c.Getf()
		if benchmarkMoneyScan != want {
			b.Errorf("Have: %f\nWant: %f", benchmarkMoneyScan, want)
		}
	}
}

// Benchmark_ParseFloat	30000000	        50.8 ns/op	       0 B/op	       0 allocs/op => Go 1.4.2
// Benchmark_ParseFloat	30000000	        53.1 ns/op	       0 B/op	       0 allocs/op => Go 1.5.0
func Benchmark_ParseFloat(b *testing.B) {
	d := "705.993300"
	var want float64 = 705.993300
	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkMoneyScan, err = strconv.ParseFloat(d, 64)
		if err != nil {
			b.Error(err)
		}
		if benchmarkMoneyScan != want {
			b.Errorf("Have: %f\nWant: %f", benchmarkMoneyScan, want)
		}
	}
}

// Benchmark_JSONUnMarshalSingle__Number	 1000000	      1256 ns/op	     240 B/op	       4 allocs/op => Go 1.4.2
// Benchmark_JSONUnMarshalSingle__Number 	 1000000	      1183 ns/op	     240 B/op	       4 allocs/op => Go 1.5.0
func Benchmark_JSONUnMarshalSingle__Number(b *testing.B) {
	benchmark_JSONUnMarshalSingle(b, []byte(`-1234.56789`), -12345679)
}

// Benchmark_JSONUnMarshalSingle__Locale	 1000000	      1531 ns/op	     272 B/op	       4 allocs/op => Go 1.4.2
// Benchmark_JSONUnMarshalSingle__Locale 	 1000000	      1258 ns/op	     272 B/op	       4 allocs/op => Go 1.5.0
func Benchmark_JSONUnMarshalSingle__Locale(b *testing.B) {
	benchmark_JSONUnMarshalSingle(b, []byte(`-2 345 678,45 €`), -23456784500)
}

// Benchmark_JSONUnMarshalSingle_Extended	 1000000	      1586 ns/op	     368 B/op	       4 allocs/op => Go 1.4.2
// Benchmark_JSONUnMarshalSingle_Extended	 1000000	      1251 ns/op	     368 B/op	       4 allocs/op => Go 1.5.0
func Benchmark_JSONUnMarshalSingle_Extended(b *testing.B) {
	benchmark_JSONUnMarshalSingle(b, []byte(`[-1999.00236, null, null]`), -19990024)
}

var benchmark_JSONUnMarshalSingleValue int64

func benchmark_JSONUnMarshalSingle(b *testing.B, data []byte, want int64) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var c money.Money
		if err := c.UnmarshalJSON(data); err != nil {
			b.Error(err)
		}
		benchmark_JSONUnMarshalSingleValue = c.Raw()
		if benchmark_JSONUnMarshalSingleValue != want {
			b.Errorf("Have: %d\nWant: %d", benchmark_JSONUnMarshalSingleValue, want)
		}
	}
}
