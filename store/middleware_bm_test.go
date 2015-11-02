// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/httputils"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	storemock "github.com/corestoreio/csfw/store/mock"
	"golang.org/x/net/context"
)

// Benchmark_WithValidateBaseUrl-4         	  200000	     11339 ns/op	    2716 B/op	      49 allocs/op
func Benchmark_WithValidateBaseUrl(b *testing.B) {

	var configReader = config.NewMockReader(
		config.WithMockValues(config.MockPV{
			config.MockPathScopeDefault(store.PathRedirectToBase):    1,
			config.MockPathScopeStore(1, store.PathSecureInFrontend): true,
			config.MockPathScopeStore(1, store.PathUnsecureBaseURL):  "http://www.corestore.io/",
			config.MockPathScopeStore(1, store.PathSecureBaseURL):    "https://www.corestore.io/",
		}),
	)

	var ctxStoreService = storemock.NewContextService(
		scope.Option{},
		func(ms *storemock.Storage) {
			ms.MockStore = func() (*store.Store, error) {
				return store.NewStore(
					&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true, true)},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
					store.SetStoreConfig(configReader),
				)
			}
		},
	)

	req, err := http.NewRequest(httputils.MethodGet, "https://corestore.io/customer/comments/view?id=1916#tab=ratings", nil)
	if err != nil {
		b.Fatal(err)
	}

	finalHandler := store.WithValidateBaseUrl(configReader)(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return errors.New("This handler should not be called!")
	}))
	want := "https://www.corestore.io/customer/comments/view?id=1916#tab=ratings"
	rec := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := finalHandler.ServeHTTPContext(ctxStoreService, rec, req); err != nil {
			b.Error(err)
		}
		if rec.HeaderMap.Get("Location") != want {
			b.Errorf("Have: %s\nWant: %s", rec.HeaderMap.Get("Location"), want)
		}
		rec.HeaderMap = nil
	}
}
