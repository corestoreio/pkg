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

package runmode_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/net/runmode"
	"github.com/corestoreio/pkg/store"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/store/storemock"
	"github.com/corestoreio/log"
)

func BenchmarkWithRunMode(b *testing.B) {
	srv := storemock.NewServiceEuroOZ(cfgmock.NewService())

	var runner = func(req *http.Request, runMode scope.TypeID, wantStoreID, wantWebsiteID int64) func(b *testing.B) {
		return func(b *testing.B) {
			rmmw := runmode.WithRunMode(srv, runmode.Options{
				Log:        log.BlackHole{}, // disabled debug and info logging
				Calculater: runMode,
			})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				haveWebsiteID, haveStoreID, haveOK := scope.FromContext(r.Context())
				if !haveOK {
					b.Fatal("Missing scope context for store and website")
				}
				if have, want := haveStoreID, wantStoreID; have != want {
					b.Fatalf("Have: %v Want AT Store: %v", have, want)
				}
				if have, want := haveWebsiteID, wantWebsiteID; have != want {
					b.Fatalf("Have: %v Want Euro Website: %v", have, want)
				}
				w.WriteHeader(http.StatusTeapot)
			}))

			b.ResetTimer()
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					rec := httptest.NewRecorder() // 4 allocs
					rmmw.ServeHTTP(rec, req)
					if have, want := rec.Code, http.StatusTeapot; have != want {
						b.Fatalf("Have: %v Want: %v", have, want)
					}
				}
			})
			//rec := httptest.NewRecorder() // 4 allocs
			//for i := 0; i < b.N; i++ {
			//	rmmw.ServeHTTP(rec, req)
			//	if have, want := rec.Code, http.StatusTeapot; have != want {
			//		b.Fatalf("Have: %v Want: %v", have, want)
			//	}
			//}
		}
	}

	b.Run("Website Default", runner(
		getReq("GET", "http://cs.io", nil),
		scope.MakeTypeID(scope.Website, 1), 2, 1)) // 2 = at; 1 = euro

	b.Run("Website OZ", runner(
		getReq("GET", "http://cs.io", nil),
		scope.MakeTypeID(scope.Website, 2), 5, 2)) // 5 = au; 2 = oz

	b.Run("Store DE Cookie", runner(
		getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "de"}),
		scope.MakeTypeID(scope.Website, 1), 1, 1)) // 2 = at; 1 = euro
	b.Run("Store DE GET", runner(
		getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=de", store.CodeURLFieldName), nil),
		scope.MakeTypeID(scope.Website, 1), 1, 1)) // 2 = at; 1 = euro

	b.Run("Store UK Cookie", runner(
		getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "uk"}),
		scope.MakeTypeID(scope.Website, 1), 4, 1)) // 4 = uk; 1 = euro
	b.Run("Store UK GET", runner(
		getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=uk", store.CodeURLFieldName), nil),
		scope.MakeTypeID(scope.Website, 1), 4, 1)) // 4 = uk; 1 = euro

}
