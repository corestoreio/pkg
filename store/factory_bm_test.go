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

package store

import "testing"

var benchmarkFactoryWebsite Website
var benchmarkFactoryWebsiteDefaultGroup Group

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkFactoryWebsiteGetDefaultGroup	  200000	      6081 ns/op	    1712 B/op	      45 allocs/op
// BenchmarkFactoryWebsiteGetDefaultGroup-4	   50000	     26210 ns/op	   10608 B/op	     229 allocs/op
func BenchmarkFactoryWebsiteGetDefaultGroup(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkFactoryWebsite, err = testFactory.Website(1)
		if err != nil {
			b.Error(err)
		}

		benchmarkFactoryWebsiteDefaultGroup, err = benchmarkFactoryWebsite.DefaultGroup()
		if err != nil {
			b.Error(err)
		}
	}
}

var benchmarkFactoryGroup Group
var benchmarkFactoryGroupDefaultStore Store

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkFactoryGroupGetDefaultStore	 1000000	      1916 ns/op	     464 B/op	      14 allocs/op
// BenchmarkFactoryGroupGetDefaultStore-4  	  300000	      5387 ns/op	    2880 B/op	      64 allocs/op
func BenchmarkFactoryGroupGetDefaultStore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkFactoryGroup, err = testFactory.Group(3)
		if err != nil {
			b.Error(err)
		}

		benchmarkFactoryGroupDefaultStore, err = benchmarkFactoryGroup.DefaultStore()
		if err != nil {
			b.Error(err)
		}
	}
}

var benchmarkFactoryStore Store
var benchmarkFactoryStoreID int64
var benchmarkFactoryStoreWebsite Website

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkFactoryStoreGetWebsite	 2000000	       656 ns/op	     176 B/op	       6 allocs/op
// BenchmarkFactoryStoreGetWebsite-4       	   50000	     32968 ns/op	   15280 B/op	     334 allocs/op
func BenchmarkFactoryStoreGetWebsite(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkFactoryStore, err = testFactory.Store(1)
		if err != nil {
			b.Error(err)
		}

		benchmarkFactoryStoreWebsite = benchmarkFactoryStore.Website
		if err := benchmarkFactoryStoreWebsite.Validate(); err != nil {
			b.Errorf("benchmarkFactoryStoreWebsite %+v", err)
		}
	}
}

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkFactoryDefaultStoreView	 2000000	       724 ns/op	     176 B/op	       7 allocs/op
// BenchmarkFactoryDefaultStoreView-4      	   50000	     40856 ns/op	   15296 B/op	     335 allocs/op
func BenchmarkFactoryDefaultStoreView(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkFactoryStoreID, err = testFactory.DefaultStoreID()
		if err != nil {
			b.Error(err)
		}
	}
}
