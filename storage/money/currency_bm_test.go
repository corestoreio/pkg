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

package money_test

import (
	"testing"

	"github.com/corestoreio/csfw/storage/money"
)

type buf []byte

// Write writes len(p) bytes from p to the Buffer.
func (b *buf) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

var benchBenchmarkNumberWriter string
var bufferNumberWriter buf

// Benchmark_NumberWriter	 1000000	      1692 ns/op	     264 B/op	      13 allocs/op
func Benchmark_NumberWriter(b *testing.B) {
	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := money.New(money.Precision(100)).Setf(-123456.789)
		_, err = c.NumberWriter(&bufferNumberWriter)
		if err != nil {
			b.Error(err)
		}
		benchBenchmarkNumberWriter = string(bufferNumberWriter)
		bufferNumberWriter = bufferNumberWriter[:0]
	}
}

var benchmarkMoneyNewGetf float64

// Benchmark_MoneyNewGetf	 2000000	       771 ns/op	     208 B/op	       7 allocs/op
func Benchmark_MoneyNewGetf(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := money.New(money.Precision(100)).Setf(-123456.789)
		benchmarkMoneyNewGetf = c.Getf()
	}
}
