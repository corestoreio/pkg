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

package config_test

import (
	"testing"

	"github.com/corestoreio/csfw/util/log"
)

func init() {
	log.SetLevel(log.StdLevelInfo)
}

// BenchmarkSectionSliceValidate	    1000	   1791239 ns/op	  158400 B/op	    4016 allocs/op => Go 1.4.2
// BenchmarkSectionSliceValidate   	    1000	   1636547 ns/op	  158400 B/op	    3213 allocs/op => Go 1.5.0
// BenchmarkSectionSliceValidate   	    1000	   1766386 ns/op	  102783 B/op	    1607 allocs/op => Go 1.5.2
func BenchmarkSectionSliceValidate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := packageAllConfiguration.Validate(); err != nil {
			b.Error(err)
		}
	}
}

var bsstj string

// BenchmarkSectionSliceToJson	     300	   4336829 ns/op	  973188 B/op	   17254 allocs/op => Go 1.4.2
// BenchmarkSectionSliceToJson 	     500	   3609676 ns/op	  914083 B/op	   14943 allocs/op => Go 1.5.0
// BenchmarkSectionSliceToJson 	     500	   3580314 ns/op	  895303 B/op	   14620 allocs/op => Go 1.5.2
func BenchmarkSectionSliceToJson(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if bsstj = packageAllConfiguration.ToJSON(); bsstj == "" {
			b.Error("JSON is empty!")
		}
	}
}

// BenchmarkSectionSliceFindFieldByPath1	20000000	       92.9 ns/op	       0 B/op	       0 allocs/op => Go 1.4.2
// BenchmarkSectionSliceFindFieldByPath1	20000000	       84.1 ns/op	       0 B/op	       0 allocs/op => Go 1.5.0
// BenchmarkSectionSliceFindFieldByPath1	20000000	        86.6 ns/op	       0 B/op	       0 allocs/op => Go 1.5.2
func BenchmarkSectionSliceFindFieldByPath1(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := packageAllConfiguration.FindFieldByPath("carriers", "usps", "gateway_url"); err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkSectionSliceFindFieldByPath5	 2000000	       587 ns/op	       0 B/op	       0 allocs/op => Go 1.4.2
// BenchmarkSectionSliceFindFieldByPath5	 3000000	       565 ns/op	       0 B/op	       0 allocs/op => Go 1.5.0
// BenchmarkSectionSliceFindFieldByPath5	 3000000	       564 ns/op	       0 B/op	       0 allocs/op => Go 1.5.2
func BenchmarkSectionSliceFindFieldByPath5(b *testing.B) {
	var paths = [][]string{
		{"carriers", "usps", "gateway_url"},
		{"wishlist", "email", "number_limit"},
		{"tax", "calculation", "apply_tax_on"},
		{"sitemap", "generate", "frequency"},
		{"sales_email", "creditmemo_comment", "guest_template"},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			if _, err := packageAllConfiguration.FindFieldByPath(path...); err != nil {
				b.Error(err)
			}
		}
	}
}
