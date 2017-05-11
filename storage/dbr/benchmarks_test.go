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
)

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
	args := make(Arguments, 0, 2)
	uc := UpdatedColumns{
		Columns:   []string{"name", "sku", "stock"},
		Arguments: Arguments{ArgString("E0S 5D Mark III"), nil, argInt64(14)},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		if args, err = uc.writeOnDuplicateKey(buf, args); err != nil {
			b.Fatalf("%+v", err)
		}
		buf.Reset()
		args = args[:0]
	}
}
