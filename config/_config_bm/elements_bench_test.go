// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/util/cstesting"
)

const enableGCPauseOutput = false

// BenchmarkSectionSliceValidate	    1000	   1791239 ns/op	  158400 B/op	    4016 allocs/op => Go 1.4.2
// BenchmarkSectionSliceValidate   	    1000	   1636547 ns/op	  158400 B/op	    3213 allocs/op => Go 1.5.0
// BenchmarkSectionSliceValidate   	    1000	   1766386 ns/op	  102783 B/op	    1607 allocs/op => Go 1.5.2
// BenchmarkSectionSliceValidate   	    2000	   1092104 ns/op	  152864 B/op	    2410 allocs/op => cfgpath.Routes
// BenchmarkSectionSliceValidate   	    2000	   1123606 ns/op	  191408 B/op	    2410 allocs/op => cfgpath.Routes with Sum32
func BenchmarkSectionSliceValidate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := packageAllConfiguration.Validate(); err != nil {
			b.Error(err)
		}
	}
	if enableGCPauseOutput {
		b.Log("GC Pause:", cstesting.GCPause())
	}
}

var bsstj string

// BenchmarkSectionSliceToJson	     300	   4336829 ns/op	  973188 B/op	   17254 allocs/op => Go 1.4.2
// BenchmarkSectionSliceToJson 	     500	   3609676 ns/op	  914083 B/op	   14943 allocs/op => Go 1.5.0
// BenchmarkSectionSliceToJson 	     500	   3580314 ns/op	  895303 B/op	   14620 allocs/op => Go 1.5.2
// BenchmarkSectionSliceToJson	     500	   3844865 ns/op	  874724 B/op	   16505 allocs/op => cfgpath.Routes
// BenchmarkSectionSliceToJson 	     500	   3728585 ns/op	  875460 B/op	   16505 allocs/op
func BenchmarkSectionSliceToJson(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if bsstj = packageAllConfiguration.ToJSON(); bsstj == "" {
			b.Error("JSON is empty!")
		}
	}
	if enableGCPauseOutput {
		b.Log("GC Pause:", cstesting.GCPause())
	}
}

var sectionSliceFindFieldByID1 element.Field

// BenchmarkSectionSliceFindFieldByID1		20000000	       92.9 ns/op	       0 B/op	       0 allocs/op => Go 1.4.2 strings
// BenchmarkSectionSliceFindFieldByID1		20000000	       84.1 ns/op	       0 B/op	       0 allocs/op => Go 1.5.0 strings
// BenchmarkSectionSliceFindFieldByID1		20000000	        86.6 ns/op	       0 B/op	       0 allocs/op => Go 1.5.2 strings
// BenchmarkSectionSliceFindFieldByID1	 	 2000000	       890 ns/op	       0 B/op	       0 allocs/op => cfgpath.Routes
// BenchmarkSectionSliceFindFieldByID1-4	 2000000	       751 ns/op	       0 B/op	       0 allocs/op => cfgpath.Routes with Sum32
// BenchmarkSectionSliceFindFieldByID1-4	10000000	       137 ns/op	       0 B/op	       0 allocs/op => cfgpath.Routes with Sum32 + array with pointers
// BenchmarkSectionSliceFindFieldByID1-4	 3000000	       484 ns/op	       0 B/op	       0 allocs/op => removed pointers
func BenchmarkSectionSliceFindFieldByID1(b *testing.B) {
	r := cfgpath.MakeRoute("carriers", "usps", "gateway_url")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		if sectionSliceFindFieldByID1, _, err = packageAllConfiguration.FindField(r); err != nil {
			b.Error(err)
		}
	}
	if sectionSliceFindFieldByID1.ID.String() != "gateway_url" {
		b.Error("Field ID must be gateway_url")
	}
}

// BenchmarkSectionSliceFindFieldByID5	 	 2000000	       587 ns/op	       0 B/op	       0 allocs/op => Go 1.4.2
// BenchmarkSectionSliceFindFieldByID5	 	 3000000	       565 ns/op	       0 B/op	       0 allocs/op => Go 1.5.0
// BenchmarkSectionSliceFindFieldByID5	 	 3000000	       564 ns/op	       0 B/op	       0 allocs/op => Go 1.5.2
// BenchmarkSectionSliceFindFieldByID5	  	  300000	      6077 ns/op	       0 B/op	       0 allocs/op => cfgpath.Routes
// BenchmarkSectionSliceFindFieldByID5-4	  300000	      4580 ns/op	       0 B/op	       0 allocs/op => cfgpath.Routes with Sum32
// BenchmarkSectionSliceFindFieldByID5-4	 2000000	       851 ns/op	       0 B/op	       0 allocs/op => cfgpath.Routes with Sum32 + array with pointers
// BenchmarkSectionSliceFindFieldByID5-4	  500000	      3045 ns/op	       0 B/op	       0 allocs/op => removed pointers
func BenchmarkSectionSliceFindFieldByID5(b *testing.B) {
	routePaths := [...]cfgpath.Route{
		cfgpath.MakeRoute("carriers", "usps", "gateway_url"),
		cfgpath.MakeRoute("wishlist", "email", "number_limit"),
		cfgpath.MakeRoute("tax", "calculation", "apply_tax_on"),
		cfgpath.MakeRoute("sitemap", "generate", "frequency"),
		cfgpath.MakeRoute("sales_email", "creditmemo_comment", "guest_template"),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range routePaths {
			var err error
			if sectionSliceFindFieldByID1, _, err = packageAllConfiguration.FindField(r); err != nil {
				b.Error(err)
			}
		}
	}
}

// BenchmarkSectionSliceFindFieldByID5_Parallel-4	 1000000	      1576 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSectionSliceFindFieldByID5_Parallel(b *testing.B) {
	routePaths := [...]cfgpath.Route{
		cfgpath.MakeRoute("carriers", "usps", "gateway_url"),
		cfgpath.MakeRoute("wishlist", "email", "number_limit"),
		cfgpath.MakeRoute("tax", "calculation", "apply_tax_on"),
		cfgpath.MakeRoute("sitemap", "generate", "frequency"),
		cfgpath.MakeRoute("sales_email", "creditmemo_comment", "guest_template"),
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, r := range routePaths {
				var err error
				if sectionSliceFindFieldByID1, _, err = packageAllConfiguration.FindField(r); err != nil {
					b.Error(err)
				}
			}
		}
	})
	if enableGCPauseOutput {
		b.Log("GC Pause:", cstesting.GCPause())
	}
}
