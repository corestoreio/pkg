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

// +build csall

package store_test

import (
	"bytes"
	"encoding/json"
	"testing"

	storemock "github.com/corestoreio/pkg/store/mock"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson/jwriter"
)

var benchmarkServiceJSONBytes []byte

func BenchmarkService_Json_Encoding(b *testing.B) {
	srv := storemock.NewServiceEuroW11G11S19()
	b.ResetTimer()

	// name                                           time/op
	// Service_Json_Encoding/easyjson_______-4        9.86µs ± 0%
	// Service_Json_Encoding/stdlibNewEncoder-4      40.60µs ± 0%
	// Service_Json_Encoding/jsoniterFastestStream-4 42.00µs ± 1%
	//
	// name                                           alloc/op
	// Service_Json_Encoding/easyjson_______-4        5.92kB ± 0%
	// Service_Json_Encoding/stdlibNewEncoder-4       6.14kB ± 0%
	// Service_Json_Encoding/jsoniterFastestStream-4 11.00kB ± 0%
	//
	// name                                           allocs/op
	// Service_Json_Encoding/easyjson_______-4          30.0 ± 0%
	// Service_Json_Encoding/stdlibNewEncoder-4         36.0 ± 0%
	// Service_Json_Encoding/jsoniterFastestStream-4    37.0 ± 0%

	// "github.com/mailru/easyjson/jwriter" Version 498e5971837e6d60575592b3e7afdfc9873ece5b on git@github.com:SchumacherFM/easyjson.git
	b.Run("easyjson_______", func(b *testing.B) {
		var jw jwriter.Writer
		for i := 0; i < b.N; i++ {
			srv.Websites().MarshalEasyJSON(&jw)
			benchmarkServiceJSONBytes, _ = jw.BuildBytes()
			// b.Fatal(string(benchmarkServiceJSONBytes))
		}
	})
	// Version 06e34e58150a5cb77afdd3807a93f270a3068456 aka. Go 1.13
	b.Run("stdlibNewEncoder", func(b *testing.B) {
		var buf bytes.Buffer
		for i := 0; i < b.N; i++ {
			je := json.NewEncoder(&buf)
			_ = je.Encode(srv.Websites())
			benchmarkServiceJSONBytes = buf.Bytes()
			buf.Reset()
			// b.Fatal(string(benchmarkServiceJSONBytes))
		}
	})
	// jsoniter "github.com/json-iterator/go" Version 0039f4ac3d5680243e7d0650c581e7ec0885ef5a
	b.Run("jsoniterFastestStream", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stream := jsoniter.ConfigFastest.BorrowStream(nil)
			stream.WriteVal(srv.Websites())
			if stream.Error != nil {
				b.Fatal(stream.Error)
			}
			benchmarkServiceJSONBytes = stream.Buffer()
			jsoniter.ConfigFastest.ReturnStream(stream)
			// b.Fatal(string(benchmarkServiceJSONBytes))
		}
	})
	_ = benchmarkServiceJSONBytes
}

// import (
// 	"testing"
//
// 	"github.com/corestoreio/pkg/config/cfgmock"
// 	"github.com/corestoreio/pkg/store"
// 	"github.com/corestoreio/pkg/store/scope"
// 	"github.com/corestoreio/pkg/store/storemock"
// )
//
// // benchmarkStoreService refactor and use a function which generates a huge
// // Service containing thousands of websites, groups and stores. Use then build
// // tags to create benchmark only tests.
// var benchmarkStoreService = storemock.NewEurozzyService(cfgmock.NewService())
//
// func Benchmark_Service_IsAllowedStoreID(b *testing.B) {
//
// 	var runner = func(runMode scope.TypeID, storeID int64) func(pb *testing.PB) {
// 		return func(pb *testing.PB) {
// 			var isA bool
// 			var stC string
// 			for pb.Next() {
// 				var err error
// 				isA, stC, err = benchmarkStoreService.IsAllowedStoreID(runMode, storeID)
// 				if err != nil {
// 					b.Error(err)
// 				}
// 				if !isA {
// 					b.Fatal("StoreID must be allowed")
// 				}
// 				if stC == "" {
// 					b.Fatal("StoreCode cannot be empty")
// 				}
// 			}
// 		}
// 	}
//
// 	b.Run("Store", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.MakeTypeID(scope.Store, 1), 6))
// 	})
// 	b.Run("Group", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.MakeTypeID(scope.Group, 1), 2))
// 	})
// 	b.Run("Website", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.MakeTypeID(scope.Website, 1), 2))
// 	})
// 	b.Run("Default", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.DefaultTypeID, 2)) // at store
// 	})
// }
//
// func Benchmark_Service_DefaultStoreID(b *testing.B) {
//
// 	var runner = func(runMode scope.TypeID) func(pb *testing.PB) {
// 		return func(pb *testing.PB) {
// 			var bmss int64
// 			for pb.Next() {
// 				var err error
// 				bmss, _, err = benchmarkStoreService.DefaultStoreID(runMode)
// 				if err != nil {
// 					b.Fatalf("%+v", err)
// 				}
// 				if bmss < 1 {
// 					b.Fatalf("StoreID must be greater than zero: %d", bmss)
// 				}
// 			}
// 		}
// 	}
//
// 	b.Run("Store", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.MakeTypeID(scope.Store, 1)))
// 	})
// 	b.Run("Group", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.MakeTypeID(scope.Group, 2)))
// 	})
// 	b.Run("Website", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.MakeTypeID(scope.Website, 1)))
// 	})
// 	b.Run("Default", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.DefaultTypeID))
// 	})
// }
//
// func Benchmark_Service_StoreIDbyCode(b *testing.B) {
//
// 	var runner = func(runMode scope.TypeID, storeCode string) func(pb *testing.PB) {
// 		return func(pb *testing.PB) {
// 			var bmss int64
// 			for pb.Next() {
// 				var err error
// 				bmss, _, err = benchmarkStoreService.StoreIDbyCode(runMode, storeCode)
// 				if err != nil {
// 					b.Fatalf("%+v", err)
// 				}
// 				if bmss < 1 {
// 					b.Fatalf("StoreID must be greater than zero: %d", bmss)
// 				}
// 			}
// 		}
// 	}
//
// 	b.Run("Store", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.MakeTypeID(scope.Store, 1), "nz"))
// 	})
// 	b.Run("Group", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.MakeTypeID(scope.Group, 2), "uk"))
// 	})
// 	b.Run("Website", func(b *testing.B) {
// 		b.ReportAllocs()
// 		b.RunParallel(runner(scope.MakeTypeID(scope.Website, 1), "at"))
// 	})
// }
//
// func Benchmark_Service_GetStore(b *testing.B) {
// 	b.ReportAllocs()
// 	b.RunParallel(func(pb *testing.PB) {
// 		var bmss store.Store
// 		for pb.Next() {
// 			var err error
// 			bmss, err = benchmarkStoreService.Store(6)
// 			if err != nil {
// 				b.Fatalf("%+v", err)
// 			}
// 			if err := bmss.Validate(); err != nil {
// 				b.Fatalf("contains errors: %+v", err)
// 			}
// 		}
// 	})
// }
//
// func Benchmark_Service_DefaultStoreView(b *testing.B) {
// 	b.ReportAllocs()
// 	b.RunParallel(func(pb *testing.PB) {
// 		var bmss store.Store
// 		for pb.Next() {
// 			var err error
// 			bmss, err = benchmarkStoreService.DefaultStoreView()
// 			if err != nil {
// 				b.Fatalf("%+v", err)
// 			}
// 			if err := bmss.Validate(); err != nil {
// 				b.Fatalf("contains errors: %+v", err)
// 			}
// 		}
// 	})
// }
//
// func Benchmark_Service_GetGroup(b *testing.B) {
// 	b.ReportAllocs()
// 	b.RunParallel(func(pb *testing.PB) {
// 		var bmsg store.Group
// 		for pb.Next() {
// 			var err error
// 			bmsg, err = benchmarkStoreService.Group(2)
// 			if err != nil {
// 				b.Fatalf("%+v", err)
// 			}
// 			if err := bmsg.Validate(); err != nil {
// 				b.Fatalf("contains errors: %+v", err)
// 			}
// 		}
// 	})
// }
//
// func Benchmark_Service_GetWebsite(b *testing.B) {
// 	b.ReportAllocs()
// 	b.RunParallel(func(pb *testing.PB) {
// 		var bmsw store.Website
// 		for pb.Next() {
// 			var err error
// 			bmsw, err = benchmarkStoreService.Website(2)
// 			if err != nil {
// 				b.Fatalf("%+v", err)
// 			}
// 			if err := bmsw.Validate(); err != nil {
// 				b.Fatalf("contains errors: %+v", err)
// 			}
// 		}
// 	})
// }
