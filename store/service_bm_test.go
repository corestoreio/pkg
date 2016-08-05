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

package store_test

import (
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
)

// benchmarkStoreService refactor and use a function which generates a huge
// Service containing thousands of websites, groups and stores. Use then build
// tags to create benchmark only tests.
var benchmarkStoreService = storemock.NewEurozzyService(cfgmock.NewService())

func Benchmark_Service_IsAllowedStoreID(b *testing.B) {

	var runner = func(runMode scope.Hash, storeID int64) func(pb *testing.PB) {
		return func(pb *testing.PB) {
			var isA bool
			var stC string
			for pb.Next() {
				var err error
				isA, stC, err = benchmarkStoreService.IsAllowedStoreID(runMode, storeID)
				if err != nil {
					b.Error(err)
				}
				if !isA {
					b.Fatal("StoreID must be allowed")
				}
				if stC == "" {
					b.Fatal("StoreCode cannot be empty")
				}
			}
		}
	}

	b.Run("Store", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.NewHash(scope.Store, 1), 6))
	})
	b.Run("Group", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.NewHash(scope.Group, 1), 2))
	})
	b.Run("Website", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.NewHash(scope.Website, 1), 2))
	})
	b.Run("Default", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.DefaultHash, 2)) // at store
	})
}

func Benchmark_Service_DefaultStoreID(b *testing.B) {

	var runner = func(runMode scope.Hash) func(pb *testing.PB) {
		return func(pb *testing.PB) {
			var bmss int64
			for pb.Next() {
				var err error
				bmss, _, err = benchmarkStoreService.DefaultStoreID(runMode)
				if err != nil {
					b.Fatalf("%+v", err)
				}
				if bmss < 1 {
					b.Fatalf("StoreID must be greater than zero: %d", bmss)
				}
			}
		}
	}

	b.Run("Store", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.NewHash(scope.Store, 1)))
	})
	b.Run("Group", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.NewHash(scope.Group, 2)))
	})
	b.Run("Website", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.NewHash(scope.Website, 1)))
	})
	b.Run("Default", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.DefaultHash))
	})
}

func Benchmark_Service_StoreIDbyCode(b *testing.B) {

	var runner = func(runMode scope.Hash, storeCode string) func(pb *testing.PB) {
		return func(pb *testing.PB) {
			var bmss int64
			for pb.Next() {
				var err error
				bmss, _, err = benchmarkStoreService.StoreIDbyCode(runMode, storeCode)
				if err != nil {
					b.Fatalf("%+v", err)
				}
				if bmss < 1 {
					b.Fatalf("StoreID must be greater than zero: %d", bmss)
				}
			}
		}
	}

	b.Run("Store", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.NewHash(scope.Store, 1), "nz"))
	})
	b.Run("Group", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.NewHash(scope.Group, 2), "uk"))
	})
	b.Run("Website", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(runner(scope.NewHash(scope.Website, 1), "at"))
	})
}

func Benchmark_Service_GetStore(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var bmss store.Store
		for pb.Next() {
			var err error
			bmss, err = benchmarkStoreService.Store(6)
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if err := bmss.Validate(); err != nil {
				b.Fatalf("contains errors: %+v", err)
			}
		}
	})
}

func Benchmark_Service_DefaultStoreView(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var bmss store.Store
		for pb.Next() {
			var err error
			bmss, err = benchmarkStoreService.DefaultStoreView()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if err := bmss.Validate(); err != nil {
				b.Fatalf("contains errors: %+v", err)
			}
		}
	})
}

func Benchmark_Service_GetGroup(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var bmsg store.Group
		for pb.Next() {
			var err error
			bmsg, err = benchmarkStoreService.Group(2)
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if err := bmsg.Validate(); err != nil {
				b.Fatalf("contains errors: %+v", err)
			}
		}
	})
}

func Benchmark_Service_GetWebsite(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var bmsw store.Website
		for pb.Next() {
			var err error
			bmsw, err = benchmarkStoreService.Website(2)
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if err := bmsw.Validate(); err != nil {
				b.Fatalf("contains errors: %+v", err)
			}
		}
	})
}
