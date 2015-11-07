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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/net/httputils"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	storemock "github.com/corestoreio/csfw/store/mock"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestWithValidateBaseUrl_DeactivatedAndShouldNotRedirectWithGETRequest(t *testing.T) {

	mockReader := config.NewMockReader(
		config.WithMockValues(config.MockPV{
			config.MockPathScopeDefault(store.PathRedirectToBase): 0,
		}),
	)

	// no post request but check deactivated
	w := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "http://corestore.io/catalog/product/view", nil)
	assert.NoError(t, err)

	err = store.WithValidateBaseUrl(mockReader)(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		assert.NotNil(t, ctx)
		assert.NotNil(t, w)
		assert.NotNil(t, r)

		assert.Empty(t, w.Header().Get("Location"))

		return nil
	})).ServeHTTPContext(context.Background(), w, req)
	assert.NoError(t, err)
}

func TestWithValidateBaseUrl_ActivatedAndShouldNotRedirectWithPOSTRequest(t *testing.T) {

	mockReader := config.NewMockReader(
		config.WithMockValues(config.MockPV{
			config.MockPathScopeDefault(store.PathRedirectToBase): 301,
		}),
	)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "http://corestore.io/catalog/product/view", nil)
	assert.NoError(t, err)

	mw := store.WithValidateBaseUrl(mockReader)(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		assert.NotNil(t, ctx)
		assert.NotNil(t, w)
		assert.NotNil(t, r)

		assert.Empty(t, w.Header().Get("Location"))

		return nil
	}))

	err = mw.ServeHTTPContext(context.Background(), w, req)
	assert.EqualError(t, err, store.ErrContextServiceNotFound.Error())

	w = httptest.NewRecorder()
	req, err = http.NewRequest(httputils.MethodPost, "http://corestore.io/catalog/product/view", strings.NewReader(`{ "k1": "v1",  "k2": { "k3": ["va1"]  }}`))
	assert.NoError(t, err)

	err = mw.ServeHTTPContext(context.Background(), w, req)
	assert.NoError(t, err)

}

func TestWithValidateBaseUrl_ActivatedAndShouldRedirectWithGETRequest(t *testing.T) {

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

	var newReq = func(urlStr string) *http.Request {
		req, err := http.NewRequest(httputils.MethodGet, urlStr, nil)
		if err != nil {
			t.Fatal(err)
		}
		return req
	}

	tests := []struct {
		rec             *httptest.ResponseRecorder
		req             *http.Request
		wantRedirectURL string
	}{
		{
			httptest.NewRecorder(),
			newReq("http://corestore.io/catalog/product/view/"),
			"http://www.corestore.io/catalog/product/view/",
		},
		{
			httptest.NewRecorder(),
			newReq("http://corestore.io/catalog/product/view"),
			"http://www.corestore.io/catalog/product/view",
		},
		{
			httptest.NewRecorder(),
			newReq("http://corestore.io"),
			"http://www.corestore.io",
		},
		{
			httptest.NewRecorder(),
			newReq("https://corestore.io/catalog/category/view?catid=1916"),
			"https://www.corestore.io/catalog/category/view?catid=1916",
		},
		{
			httptest.NewRecorder(),
			newReq("https://corestore.io/customer/comments/view?id=1916#tab=ratings"),
			"https://www.corestore.io/customer/comments/view?id=1916#tab=ratings",
		},
	}

	for i, test := range tests {
		i = i
		mw := store.WithValidateBaseUrl(configReader)(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return fmt.Errorf("This handler should not be called! Iindex %d", i)
		}))
		assert.NoError(t, mw.ServeHTTPContext(ctxStoreService, test.rec, test.req), "Index %d", i)
		assert.Exactly(t, test.wantRedirectURL, test.rec.HeaderMap.Get("Location"), "Index %d", i)
	}
}

//func getTestRequest(m, u string, c *http.Cookie) *http.Request {
//	req, err := http.NewRequest(m, u, nil)
//	if err != nil {
//		panic(err)
//	}
//	if c != nil {
//		req.AddCookie(c)
//	}
//	return req
//}
//
//var testsInitByRequest = []struct {
//	req                  *http.Request
//	haveSO               scope.Option
//	haveScopeType        scope.Scope
//	wantStoreCode        string           // this is the default store in a scope, lookup in getInitializedStoreService
//	wantRequestStoreCode scope.StoreCoder // can be nil in tests
//	wantErr              error
//	wantCookie           string
//}{
//	{
//		getTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "uk"}),
//		scope.Option{Store: scope.MockID(1)}, scope.StoreID, "de", scope.MockCode("uk"), nil, store.CookieName + "=uk;",
//	},
//	{
//		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
//		scope.Option{Store: scope.MockID(1)}, scope.StoreID, "de", scope.MockCode("uk"), nil, store.CookieName + "=uk;", // generates a new 1y valid cookie
//	},
//	{
//		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=%20uk", nil),
//		scope.Option{Store: scope.MockID(1)}, scope.StoreID, "de", scope.MockCode("uk"), store.ErrStoreCodeInvalid, "",
//	},
//	{
//		getTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "de"}),
//		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", scope.MockCode("de"), nil, store.CookieName + "=de;",
//	},
//	{
//		getTestRequest("GET", "http://cs.io", nil),
//		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", nil, store.ErrUnsupportedScope, "",
//	},
//	{
//		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=de", nil),
//		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", scope.MockCode("de"), nil, store.CookieName + "=de;", // generates a new 1y valid cookie
//	},
//	{
//		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=at", nil),
//		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", scope.MockCode("at"), nil, store.CookieName + "=;", // generates a delete cookie
//	},
//	{
//		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=cz", nil),
//		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", nil, store.ErrIDNotFoundTableStoreSlice, "",
//	},
//	{
//		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
//		scope.Option{Group: scope.MockID(1)}, scope.GroupID, "at", nil, store.ErrStoreChangeNotAllowed, "",
//	},
//
//	{
//		getTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "nz"}),
//		scope.Option{Website: scope.MockID(2)}, scope.WebsiteID, "au", scope.MockCode("nz"), nil, store.CookieName + "=nz;",
//	},
//	{
//		getTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.CookieName, Value: "n'z"}),
//		scope.Option{Website: scope.MockID(2)}, scope.WebsiteID, "au", nil, store.ErrStoreCodeInvalid, "",
//	},
//	{
//		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
//		scope.Option{Website: scope.MockID(2)}, scope.WebsiteID, "au", nil, store.ErrStoreChangeNotAllowed, "",
//	},
//	{
//		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=nz", nil),
//		scope.Option{Website: scope.MockID(2)}, scope.WebsiteID, "au", scope.MockCode("nz"), nil, store.CookieName + "=nz;",
//	},
//	{
//		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=ch", nil),
//		scope.Option{Website: scope.MockID(1)}, scope.WebsiteID, "at", nil, store.ErrStoreNotActive, "",
//	},
//	{
//		getTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=nz", nil),
//		scope.Option{Website: scope.MockID(1)}, scope.DefaultID, "at", scope.MockCode("nz"), store.ErrStoreChangeNotAllowed, "",
//	},
//}
//
func TestWithInitStoreByRequest(t *testing.T) {
	//	errLogBuf.Reset()
	//	defer errLogBuf.Reset()
	//
	//	for _, test := range testsInitByRequest {
	//		if _, haveErr := getInitializedStoreService.InitByRequest(nil, nil, test.haveScopeType); haveErr != nil {
	//			assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
	//		} else {
	//			t.Fatal("InitByRequest should return an error if used without running Init() first.")
	//		}
	//
	//		if err := getInitializedStoreService.Init(test.haveSO); err != nil {
	//			assert.EqualError(t, store.ErrUnsupportedScope, err.Error())
	//			t.Log("continuing for loop because of expected store.ErrUnsupportedScopeGroup")
	//			getInitializedStoreService.ClearCache(true)
	//			continue
	//		}
	//
	//		if s, err := getInitializedStoreService.Store(); err == nil {
	//			assert.EqualValues(t, test.wantStoreCode, s.Data.Code.String)
	//		} else {
	//			assert.EqualError(t, err, store.ErrStoreNotFound.Error())
	//			t.Log("continuing for loop because of expected store.ErrStoreNotFound")
	//			getInitializedStoreService.ClearCache(true)
	//			continue
	//		}
	//		getInitializedStoreService.ClearCache(true)
	//	}
}

//
//func TestInitByRequestInDepth(t *testing.T) {
//	errLogBuf.Reset()
//	defer errLogBuf.Reset()
//
//	for i, test := range testsInitByRequest {
//		if err := getInitializedStoreService.Init(test.haveSO); err != nil {
//			assert.EqualError(t, store.ErrUnsupportedScope, err.Error())
//			t.Log("continuing for loop because of expected store.ErrUnsupportedScopeGroup")
//			getInitializedStoreService.ClearCache(true)
//			continue
//		}
//
//		resRec := httptest.NewRecorder()
//
//		haveStore, haveErr := getInitializedStoreService.InitByRequest(resRec, test.req, test.haveScopeType)
//		if test.wantErr != nil {
//			assert.Nil(t, haveStore)
//			assert.Error(t, haveErr, "Index %d", i)
//			assert.EqualError(t, haveErr, test.wantErr.Error(), "\nIndex: %d\nError: %s", i, errLogBuf.String())
//		} else {
//
//			assert.NoError(t, haveErr, "Test: %#v\n\n%s\n\n", test, errLogBuf.String())
//
//			if test.wantRequestStoreCode != nil {
//				assert.NotNil(t, haveStore, "URL Query: %#v\nCookies %#v", test.req.URL.Query(), test.req.Cookies())
//				assert.EqualValues(t, test.wantRequestStoreCode.StoreCode(), haveStore.Data.Code.String)
//
//				newKeks := resRec.HeaderMap.Get("Set-Cookie")
//				if test.wantCookie != "" {
//					assert.Contains(t, newKeks, test.wantCookie, "%#v", test)
//					//					t.Logf(
//					//						"\nwantRequestStoreCode: %s\nCookie Str: %#v\n",
//					//						test.wantRequestStoreCode.Code(),
//					//						newKeks,
//					//					)
//				} else {
//					assert.Empty(t, newKeks, "%#v", test)
//				}
//
//			} else {
//				assert.Nil(t, haveStore, "%#v", haveStore)
//			}
//		}
//		getInitializedStoreService.ClearCache(true)
//	}
//}

func newStoreServiceWithTokenCtx(initO scope.Option, tokenStoreCode string) context.Context {
	ctx := store.NewContextReader(context.Background(), getInitializedStoreService(initO), nil)
	tok := jwt.New(jwt.SigningMethodHS256)
	tok.Claims[store.CookieName] = tokenStoreCode
	ctx = ctxjwt.NewContext(ctx, tok)
	return ctx
}

func TestWithInitStoreByToken(t *testing.T) {

	var newReq = func(i int) *http.Request {
		req, err := http.NewRequest(httputils.MethodGet, fmt.Sprintf("https://corestore.io/store/list/%d", i), nil)
		if err != nil {
			t.Fatal(err)
		}
		return req
	}

	tests := []struct {
		ctx           context.Context
		wantStoreCode string
		wantErr       error
	}{
		{store.NewContextReader(context.Background(), nil, nil), "de", store.ErrContextServiceNotFound},
		{store.NewContextReader(context.Background(), getInitializedStoreService(scope.Option{Store: scope.MockCode("de")}), nil), "de", ctxjwt.ErrContextJWTNotFound},
		{newStoreServiceWithTokenCtx(scope.Option{Store: scope.MockCode("de")}, "de"), "de", nil},
		{newStoreServiceWithTokenCtx(scope.Option{Store: scope.MockCode("at")}, "ch"), "at", store.ErrStoreNotActive},
		{newStoreServiceWithTokenCtx(scope.Option{Store: scope.MockCode("de")}, "at"), "at", nil},
		{newStoreServiceWithTokenCtx(scope.Option{Store: scope.MockCode("de")}, "a$t"), "de", store.ErrStoreCodeInvalid},
		{newStoreServiceWithTokenCtx(scope.Option{Store: scope.MockCode("at")}, ""), "at", store.ErrStoreCodeEmpty},

		{newStoreServiceWithTokenCtx(scope.Option{Group: scope.MockID(1)}, "de"), "de", nil},
		{newStoreServiceWithTokenCtx(scope.Option{Group: scope.MockID(1)}, "ch"), "at", store.ErrStoreNotActive},
		{newStoreServiceWithTokenCtx(scope.Option{Group: scope.MockID(1)}, " ch"), "at", store.ErrStoreCodeInvalid},
		{newStoreServiceWithTokenCtx(scope.Option{Group: scope.MockID(1)}, "uk"), "at", store.ErrStoreChangeNotAllowed},

		{newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "uk"), "au", store.ErrStoreChangeNotAllowed},
		{newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "nz"), "nz", nil},
		{newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "n z"), "au", store.ErrStoreCodeInvalid},
		{newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "au"), "au", nil},
		{newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, ""), "au", store.ErrStoreCodeEmpty},
	}
	for i, test := range tests {

		mw := store.WithInitStoreByToken()(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			_, haveReqStore, err := store.FromContextReader(ctx)
			if err != nil {
				return err
			}

			assert.Exactly(t, test.wantStoreCode, haveReqStore.StoreCode(), "Index %d", i)
			return nil
		}))
		rec := httptest.NewRecorder()
		surfErr := mw.ServeHTTPContext(test.ctx, rec, newReq(i))
		if test.wantErr != nil {
			assert.EqualError(t, surfErr, test.wantErr.Error(), "Index %d", i)
			continue
		}

		assert.NoError(t, surfErr, "Index %d", i)
	}
}

func TestWithInitStoreByToken_Alloc_Investigations_TEMP(t *testing.T) {
	// this Test is related to Benchmark_WithInitStoreByToken

	ctx := newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "nz")
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "https://corestore.io/store/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	mw := store.WithInitStoreByToken()(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		_, haveReqStore, err := store.FromContextReader(ctx)
		if err != nil {
			return err
		}
		if "nz" != haveReqStore.StoreCode() {
			t.Errorf("Have: %s\nWant: nz", haveReqStore.StoreCode())
		}
		return nil
	}))

	for i := 0; i < 2; i++ {
		if err := mw.ServeHTTPContext(ctx, rec, req); err != nil {
			t.Error(err)
		}
	}
}
