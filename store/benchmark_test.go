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

	"github.com/corestoreio/pkg/store"
	storemock "github.com/corestoreio/pkg/store/mock"
	jsoniter "github.com/json-iterator/go"
	// jlexer "github.com/mailru/easyjson/jlexer"
	//"github.com/mailru/easyjson/jwriter"
	segjson "github.com/segmentio/encoding/json"
)

var benchmarkServiceJSONBytes []byte

func BenchmarkService_Json_Encoding(b *testing.B) {
	srv := storemock.NewServiceEuroW11G11S19()
	b.ResetTimer()

	// b.Run("easyjsonMEJ_____", func(b *testing.B) {
	//	var jw jwriter.Writer
	//	for i := 0; i < b.N; i++ {
	//		srv.Websites().MarshalEasyJSON(&jw)
	//		benchmarkServiceJSONBytes, _ = jw.BuildBytes()
	//		// b.Fatal(string(benchmarkServiceJSONBytes))
	//	}
	//})
	// b.Run("easyjsonMJ_____", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		benchmarkServiceJSONBytes, _ = srv.Websites().MarshalJSON()
	//		// b.Fatal(string(benchmarkServiceJSONBytes))
	//	}
	//})

	b.Run("stdlibNewEncoder", func(b *testing.B) {
		var buf bytes.Buffer
		je := json.NewEncoder(&buf)
		for i := 0; i < b.N; i++ {
			_ = je.Encode(srv.Websites())
			benchmarkServiceJSONBytes = buf.Bytes()
			buf.Reset()
			// b.Fatal(string(benchmarkServiceJSONBytes))
		}
	})
	b.Run("segmentioNewEncoder", func(b *testing.B) {
		var buf bytes.Buffer
		for i := 0; i < b.N; i++ {
			je := segjson.NewEncoder(&buf)
			_ = je.Encode(srv.Websites())
			benchmarkServiceJSONBytes = buf.Bytes()
			buf.Reset()
			// b.Fatal(string(benchmarkServiceJSONBytes))
		}
	})
	b.Run("stdlibMarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchmarkServiceJSONBytes, _ = json.Marshal(srv.Websites())
			// b.Fatal(string(benchmarkServiceJSONBytes))
		}
	})
	b.Run("segmentioMarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchmarkServiceJSONBytes, _ = segjson.Marshal(srv.Websites())
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

func BenchmarkService_Json_Decoding(b *testing.B) {
	// b.Run("easyjson_______", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		var w store.StoreWebsites
	//		l := jlexer.Lexer{Data: rawJSONData}
	//		w.UnmarshalEasyJSON(&l)
	//		if l.Error() != nil {
	//			b.Fatal(l.Error())
	//		}
	//	}
	//})

	b.Run("stdlibNewDecoder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var w store.StoreWebsites
			je := json.NewDecoder(bytes.NewReader(rawJSONData))
			if err := je.Decode(&w); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("segmentioNewDecoder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var w store.StoreWebsites
			je := segjson.NewDecoder(bytes.NewReader(rawJSONData))
			if err := je.Decode(&w); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("stdlibUnmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var w store.StoreWebsites
			if err := json.Unmarshal(rawJSONData, &w); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("segmentioUnmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var w store.StoreWebsites
			if err := segjson.Unmarshal(rawJSONData, &w); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jsoniterFastestStream", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var w store.StoreWebsites
			je := jsoniter.ConfigFastest.NewDecoder(bytes.NewReader(rawJSONData))
			if err := je.Decode(&w); err != nil {
				b.Fatal(err)
			}
		}
	})
	_ = benchmarkServiceJSONBytes
}

var rawJSONData = []byte(`{
  "data": [
    {
      "code": "admin",
      "name": "Admin",
      "stores": {
        "data": [
          {
            "code": "admin",
            "name": "Admin",
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "name": "Admin",
            "code": "admin"
          }
        ]
      }
    },
    {
      "websiteID": 2,
      "code": "de",
      "name": "Deutschland",
      "sortOrder": 1,
      "defaultGroupID": 2,
      "isDefault": true,
      "stores": {
        "data": [
          {
            "storeID": 2,
            "code": "dede",
            "websiteID": 2,
            "groupID": 2,
            "name": "de",
            "sortOrder": 1,
            "isActive": true
          },
          {
            "storeID": 3,
            "code": "detr",
            "websiteID": 2,
            "groupID": 2,
            "name": "tr",
            "sortOrder": 4,
            "isActive": true
          },
          {
            "storeID": 5,
            "code": "deen",
            "websiteID": 2,
            "groupID": 2,
            "name": "en",
            "sortOrder": 4,
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "groupID": 2,
            "websiteID": 2,
            "name": "b2c",
            "rootCategoryID": 2,
            "defaultStoreID": 2,
            "code": "b2c"
          }
        ]
      }
    },
    {
      "websiteID": 3,
      "code": "ch",
      "name": "Schweiz",
      "sortOrder": 2,
      "defaultGroupID": 3,
      "stores": {
        "data": [
          {
            "storeID": 6,
            "code": "chde",
            "websiteID": 3,
            "groupID": 3,
            "name": "de",
            "sortOrder": 1,
            "isActive": true
          },
          {
            "storeID": 7,
            "code": "chfr",
            "websiteID": 3,
            "groupID": 3,
            "name": "fr",
            "sortOrder": 2,
            "isActive": true
          },
          {
            "storeID": 8,
            "code": "chit",
            "websiteID": 3,
            "groupID": 3,
            "name": "it",
            "sortOrder": 3,
            "isActive": true
          },
          {
            "storeID": 9,
            "code": "chen",
            "websiteID": 3,
            "groupID": 3,
            "name": "en",
            "sortOrder": 4,
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "groupID": 3,
            "websiteID": 3,
            "name": "b2c",
            "rootCategoryID": 2,
            "defaultStoreID": 6,
            "code": "b2c"
          }
        ]
      }
    },
    {
      "websiteID": 4,
      "code": "it",
      "name": "Italien",
      "sortOrder": 3,
      "defaultGroupID": 4,
      "stores": {
        "data": [
          {
            "storeID": 10,
            "code": "itit",
            "websiteID": 4,
            "groupID": 4,
            "name": "it",
            "sortOrder": 1,
            "isActive": true
          },
          {
            "storeID": 11,
            "code": "itde",
            "websiteID": 4,
            "groupID": 4,
            "name": "de",
            "sortOrder": 2,
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "groupID": 4,
            "websiteID": 4,
            "name": "b2c",
            "rootCategoryID": 2,
            "defaultStoreID": 10,
            "code": "b2c"
          }
        ]
      }
    },
    {
      "websiteID": 5,
      "code": "fr",
      "name": "Frankreich",
      "sortOrder": 4,
      "defaultGroupID": 5,
      "stores": {
        "data": [
          {
            "storeID": 12,
            "code": "frfr",
            "websiteID": 5,
            "groupID": 5,
            "name": "fr",
            "sortOrder": 1,
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "groupID": 5,
            "websiteID": 5,
            "name": "b2c",
            "rootCategoryID": 2,
            "defaultStoreID": 12,
            "code": "b2c"
          }
        ]
      }
    },
    {
      "websiteID": 6,
      "code": "be",
      "name": "Belgien",
      "sortOrder": 5,
      "defaultGroupID": 6,
      "stores": {
        "data": [
          {
            "storeID": 13,
            "code": "befr",
            "websiteID": 6,
            "groupID": 6,
            "name": "fr",
            "sortOrder": 1,
            "isActive": true
          },
          {
            "storeID": 14,
            "code": "been",
            "websiteID": 6,
            "groupID": 6,
            "name": "en",
            "sortOrder": 2,
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "groupID": 6,
            "websiteID": 6,
            "name": "b2c",
            "rootCategoryID": 2,
            "defaultStoreID": 13,
            "code": "b2c"
          }
        ]
      }
    },
    {
      "websiteID": 7,
      "code": "lu",
      "name": "Luxemburg",
      "sortOrder": 6,
      "defaultGroupID": 7,
      "stores": {
        "data": [
          {
            "storeID": 15,
            "code": "lufr",
            "websiteID": 7,
            "groupID": 7,
            "name": "fr",
            "sortOrder": 1,
            "isActive": true
          },
          {
            "storeID": 16,
            "code": "lude",
            "websiteID": 7,
            "groupID": 7,
            "name": "de",
            "sortOrder": 2,
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "groupID": 7,
            "websiteID": 7,
            "name": "b2c",
            "rootCategoryID": 2,
            "defaultStoreID": 15,
            "code": "b2c"
          }
        ]
      }
    },
    {
      "websiteID": 8,
      "code": "at",
      "name": "Österreich",
      "sortOrder": 7,
      "defaultGroupID": 8,
      "stores": {
        "data": [
          {
            "storeID": 17,
            "code": "atde",
            "websiteID": 8,
            "groupID": 8,
            "name": "de",
            "sortOrder": 1,
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "groupID": 8,
            "websiteID": 8,
            "name": "b2c",
            "rootCategoryID": 2,
            "defaultStoreID": 17,
            "code": "b2c"
          }
        ]
      }
    },
    {
      "websiteID": 9,
      "code": "int",
      "name": "International",
      "sortOrder": 8,
      "defaultGroupID": 9,
      "stores": {
        "data": [
          {
            "storeID": 18,
            "code": "inten",
            "websiteID": 9,
            "groupID": 9,
            "name": "en",
            "sortOrder": 1,
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "groupID": 9,
            "websiteID": 9,
            "name": "b2c",
            "rootCategoryID": 2,
            "defaultStoreID": 18,
            "code": "b2c"
          }
        ]
      }
    },
    {
      "websiteID": 10,
      "code": "nl",
      "name": "Netherlands",
      "sortOrder": 9,
      "defaultGroupID": 10,
      "stores": {
        "data": [
          {
            "storeID": 19,
            "code": "nlen",
            "websiteID": 10,
            "groupID": 10,
            "name": "en",
            "sortOrder": 1,
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "groupID": 10,
            "websiteID": 10,
            "name": "b2c",
            "rootCategoryID": 2,
            "defaultStoreID": 19,
            "code": "b2c"
          }
        ]
      }
    },
    {
      "websiteID": 11,
      "code": "uk",
      "name": "United Kingdom",
      "sortOrder": 10,
      "defaultGroupID": 11,
      "stores": {
        "data": [
          {
            "storeID": 20,
            "code": "uken",
            "websiteID": 11,
            "groupID": 11,
            "name": "en",
            "sortOrder": 1,
            "isActive": true
          }
        ]
      },
      "storeGroups": {
        "data": [
          {
            "groupID": 11,
            "websiteID": 11,
            "name": "b2c",
            "rootCategoryID": 2,
            "defaultStoreID": 20,
            "code": "b2c"
          }
        ]
      }
    }
  ]
}`)

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
