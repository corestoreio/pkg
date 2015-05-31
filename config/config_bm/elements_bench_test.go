// Copyright 2015 CoreStore Authors
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

	"github.com/corestoreio/csfw/utils/log"
)

func init() {
	log.Set(log.NewStdLogger())
	log.SetLevel(log.StdLevelDebug)
}

// BenchmarkSectionSliceValidate	    1000	   1457760 ns/op	   43520 B/op	     804 allocs/op
func BenchmarkSectionSliceValidate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := packageAllConfiguration.Validate(); err != nil {
			b.Error(err)
		}
	}
}

var bsstj string

// BenchmarkSectionSliceToJson	     300	   5826655 ns/op	 1378568 B/op	   17275 allocs/op <-- encoding/json NewEncoder io.Pipe with ReadAll
// BenchmarkSectionSliceToJson	     300	   5582783 ns/op	  992781 B/op	   17255 allocs/op <-- encoding/json NewEncoder with buffer
// BenchmarkSectionSliceToJson	     300	   5613962 ns/op	 1134998 B/op	   17258 allocs/op <-- encondig/json.Marshal
func BenchmarkSectionSliceToJson(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if bsstj = packageAllConfiguration.ToJson(); bsstj == "" {
			b.Error("JSON is empty!")
		}
	}
}

// BenchmarkSectionSliceFindFieldByPath1	20000000	       101 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSectionSliceFindFieldByPath1(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := packageAllConfiguration.FindFieldByPath("carriers", "usps", "gateway_url"); err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkSectionSliceFindFieldByPath5	 2000000	       654 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSectionSliceFindFieldByPath5(b *testing.B) {
	var paths = [][]string{
		[]string{"carriers", "usps", "gateway_url"},
		[]string{"wishlist", "email", "number_limit"},
		[]string{"tax", "calculation", "apply_tax_on"},
		[]string{"sitemap", "generate", "frequency"},
		[]string{"sales_email", "creditmemo_comment", "guest_template"},
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
