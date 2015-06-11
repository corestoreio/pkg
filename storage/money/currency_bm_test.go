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

var benchBenchmarkUnformatted string // -123456.79
// BenchmarkUnformatted	 1000000	      2406 ns/op	     161 B/op	       9 allocs/op <= fmt.Sprintf
func BenchmarkUnformatted(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := money.New(money.Precision(100)).Setf(-123456.789)
		benchBenchmarkUnformatted = c.Unformatted()
	}
}
