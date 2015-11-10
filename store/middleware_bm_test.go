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
	"github.com/corestoreio/csfw/utils/log"
	"golang.org/x/net/context"
)

// Benchmark_WithValidateBaseUrl-4         	  200000	      8055 ns/op	    2108 B/op	      29 allocs/op
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

// Benchmark_WithInitStoreByToken-4	  100000	     17297 ns/op	    9112 B/op	     203 allocs/op => old bug
// Benchmark_WithInitStoreByToken-4	 2000000	       810 ns/op	     128 B/op	       5 allocs/op => new
func Benchmark_WithInitStoreByToken(b *testing.B) {
	// see TestWithInitStoreByToken_Alloc_Investigations_TEMP
	b.ReportAllocs()

	wantStoreCode := "nz"
	ctx := newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "nz")

	mw := store.WithInitStoreByToken()(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		_, haveReqStore, err := store.FromContextReader(ctx)
		if err != nil {
			return err
		}

		if wantStoreCode != haveReqStore.StoreCode() {
			b.Errorf("Want: %s\nHave: %s", wantStoreCode, haveReqStore.StoreCode())
		}
		return nil
	}))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "https://corestore.io/store/list/", nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := mw.ServeHTTPContext(ctx, rec, req); err != nil {
			b.Error(err)
		}
	}
}

// Benchmark_WithInitStoreByFormCookie-4	    3000	    481881 ns/op	  189103 B/op	     232 allocs/op => with debug enabled
// Benchmark_WithInitStoreByFormCookie-4	  300000	      4797 ns/op	    1016 B/op	      16 allocs/op => debug disabled
func Benchmark_WithInitStoreByFormCookie(b *testing.B) {
	store.PkgLog.SetLevel(log.StdLevelInfo)
	b.ReportAllocs()

	wantStoreCode := "nz"
	ctx := store.NewContextReader(context.Background(), getInitializedStoreService(scope.Option{Website: scope.MockID(2)}))

	mw := store.WithInitStoreByFormCookie()(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		_, haveReqStore, err := store.FromContextReader(ctx)
		if err != nil {
			return err
		}

		if wantStoreCode != haveReqStore.StoreCode() {
			b.Errorf("Want: %s\nHave: %s", wantStoreCode, haveReqStore.StoreCode())
		}
		return nil
	}))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "https://corestore.io/store/list/", nil)
	if err != nil {
		b.Fatal(err)
	}

	req.AddCookie(&http.Cookie{
		Name:  store.ParamName,
		Value: wantStoreCode,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := mw.ServeHTTPContext(ctx, rec, req); err != nil {
			b.Error(err)
		}
	}
}
