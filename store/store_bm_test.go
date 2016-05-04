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
	"encoding/json"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
)

var benchmarkJSON []byte
var benchmarkJSONStore = store.MustNewStore(
	&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
	&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
	&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
)

// BenchmarkJSONMarshal-4	  300000	      4343 ns/op	    1032 B/op	      12 allocs/op
func BenchmarkJSONMarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkJSON, err = json.Marshal(benchmarkJSONStore)
		if err != nil {
			b.Error(err)
		}
	}
}

// @todo add encoders
// BenchmarkJSONCodec-4  	  500000	      4157 ns/op	    1648 B/op	       4 allocs/op
//func BenchmarkJSONCodec(b *testing.B) {
//	b.ReportAllocs()
//	var buf bytes.Buffer
//	for i := 0; i < b.N; i++ {
//		if err := benchmarkJSONStore.ToJSON(&buf); err != nil {
//			b.Error(err)
//		}
//		benchmarkJSON = buf.Bytes()
//		buf.Reset()
//	}
//}
