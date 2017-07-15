// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr

import (
	"bytes"
	"testing"

	"github.com/corestoreio/csfw/util/bufferpool"
)

var preprocessSink string

// BenchmarkInterpolate-4   	  500000	      4013 ns/op	     174 B/op	      11 allocs/op with reflection
// BenchmarkInterpolate-4   	  500000	      3591 ns/op	     174 B/op	      11 allocs/op string
// BenchmarkInterpolate-4   	  500000	      3599 ns/op	     174 B/op	      11 allocs/op []byte
func BenchmarkInterpolate(b *testing.B) {
	ipBuf := bufferpool.Get()
	defer bufferpool.Put(ipBuf)

	const want = `SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND d = 4 AND e = 5 AND f = 6 AND g = 7 AND h = 8 AND i = 9 AND j = 10 AND k = 'Hello' AND l = 1`
	var sqlBytes = []byte("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ? AND g = ? AND h = ? AND i = ? AND j = ? AND k = ? AND l = ?")
	args := Values{
		Int64s{1, -2, 3, 4, 5, 6, 7, 8, 9, 10},
		String("Hello"),
		Bool(true),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := interpolate(ipBuf, sqlBytes, args...); err != nil {
			b.Fatal(err)
		}
		preprocessSink = ipBuf.String()
		ipBuf.Reset()
	}
	if preprocessSink != want {
		b.Fatalf("Have: %v Want: %v", ipBuf.String(), want)
	}
}

var benchmarkIsValidIdentifier int8

// BenchmarkIsValidIdentifier-4   	20000000	        92.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIsValidIdentifier-4   	 5000000	       280 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIsValidIdentifier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkIsValidIdentifier = isValidIdentifier(`store_owner.catalog_product_entity_varchar`)
	}
	if benchmarkIsValidIdentifier != 0 {
		b.Fatalf("Should be zero but got %d", benchmarkIsValidIdentifier)
	}
}

func BenchmarkQuoteAlias(b *testing.B) {
	const want = "(e.price * a.tax * e.weee) AS `final_price`"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have := Quoter.exprAlias("(e.price * a.tax * e.weee)", "final_price"); have != want {
			b.Fatalf("Have %s\nWant %s\n", have, want)
		}
	}
}

// BenchmarkUpdatedColumns_writeOnDuplicateKey-4   	 5000000	       337 ns/op	       0 B/op	       0 allocs/op
func BenchmarkUpdatedColumns_writeOnDuplicateKey(b *testing.B) {
	buf := new(bytes.Buffer)
	args := make(Values, 0, 2)
	uc := UpdatedColumns{
		Columns:   []string{"name", "sku", "stock"},
		Arguments: Values{String("E0S 5D Mark III"), nil, Int64(14)},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := uc.writeOnDuplicateKey(buf); err != nil {
			b.Fatalf("%+v", err)
		}
		var err error
		if args, err = uc.appendArgs(args); err != nil {
			b.Fatalf("%+v", err)
		}
		buf.Reset()
		args = args[:0]
	}
}
