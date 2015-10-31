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
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/httputils"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	storemock "github.com/corestoreio/csfw/store/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestWithValidateBaseUrlNoRedirectGET(t *testing.T) {

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

func TestWithValidateBaseUrlNoRedirectPOST(t *testing.T) {

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
	assert.Contains(t, err.Error(), "Cannot extract config.Reader from context")

	w = httptest.NewRecorder()
	req, err = http.NewRequest(httputils.MethodPost, "http://corestore.io/catalog/product/view", strings.NewReader(`{ "k1": "v1",  "k2": { "k3": ["va1"]  }}`))
	assert.NoError(t, err)

	err = mw.ServeHTTPContext(context.Background(), w, req)
	assert.NoError(t, err)

}

func TestWithValidateBaseUrlNoRedirectValidBaseURL(t *testing.T) {

	var configReader = config.NewMockReader(
		config.WithMockValues(config.MockPV{
			config.MockPathScopeDefault(store.PathRedirectToBase):   1,
			config.MockPathScopeStore(1, store.PathUnsecureBaseURL): "http://www.corestore.io/",
			config.MockPathScopeStore(1, store.PathSecureBaseURL):   "https://www.corestore.io/",
		}),
	)

	var ctxStoreService = storemock.NewContextService(
		scope.Option{},
		func(ms *storemock.Storage) {
			ms.MockStore = func() (*store.Store, error) {
				return store.NewStore(
					&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
					store.SetStoreConfig(configReader),
				)
			}
		},
	)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "http://corestore.io/catalog/product/view", nil)
	assert.NoError(t, err)

	mw := store.WithValidateBaseUrl(configReader)(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		assert.NotNil(t, ctx)
		assert.NotNil(t, w)
		assert.NotNil(t, r)

		assert.Empty(t, w.Header().Get("Location"))

		return nil
	}))

	err = mw.ServeHTTPContext(ctxStoreService, w, req)
	assert.NoError(t, err)

	//	w = httptest.NewRecorder()
	//	req, err = http.NewRequest(httputils.MethodPost, "http://corestore.io/catalog/product/view", strings.NewReader(`{ "k1": "v1",  "k2": { "k3": ["va1"]  }}`))
	//	assert.NoError(t, err)
	//
	//	err = mw.ServeHTTPContext(context.Background(), w, req)
	//	assert.NoError(t, err)

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
//func TestInitByRequestGeneral(t *testing.T) {
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
//}
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

//func TestWithInitStoreByToken(t *testing.T) {
//
//	getToken := func(code string) *jwt.Token {
//		t := jwt.New(jwt.SigningMethodHS256)
//		t.Claims[store.CookieName] = code
//		return t
//	}
//
//	tests := []struct {
//		haveSO             scope.Option
//		haveCodeToken      string
//		haveScopeType      scope.Scope
//		wantStoreCode      string           // this is the default store in a scope, lookup in getInitializedStoreService
//		wantTokenStoreCode scope.StoreCoder // can be nil
//		wantErr            error
//	}{
//		{scope.Option{Store: scope.MockCode("de")}, "de", scope.StoreID, "de", scope.MockCode("de"), nil},
//		{scope.Option{Store: scope.MockCode("de")}, "at", scope.StoreID, "de", scope.MockCode("at"), nil},
//		{scope.Option{Store: scope.MockCode("de")}, "a$t", scope.StoreID, "de", nil, nil},
//		{scope.Option{Store: scope.MockCode("at")}, "ch", scope.StoreID, "at", nil, store.ErrStoreNotActive},
//		{scope.Option{Store: scope.MockCode("at")}, "", scope.StoreID, "at", nil, nil},
//
//		{scope.Option{Group: scope.MockID(1)}, "de", scope.GroupID, "at", scope.MockCode("de"), nil},
//		{scope.Option{Group: scope.MockID(1)}, "ch", scope.GroupID, "at", nil, store.ErrStoreNotActive},
//		{scope.Option{Group: scope.MockID(1)}, " ch", scope.GroupID, "at", nil, nil},
//		{scope.Option{Group: scope.MockID(1)}, "uk", scope.GroupID, "at", nil, store.ErrStoreChangeNotAllowed},
//
//		{scope.Option{Website: scope.MockID(2)}, "uk", scope.WebsiteID, "au", nil, store.ErrStoreChangeNotAllowed},
//		{scope.Option{Website: scope.MockID(2)}, "nz", scope.WebsiteID, "au", scope.MockCode("nz"), nil},
//		{scope.Option{Website: scope.MockID(2)}, "n z", scope.WebsiteID, "au", nil, nil},
//		{scope.Option{Website: scope.MockID(2)}, "", scope.WebsiteID, "au", nil, nil},
//	}
//	for _, test := range tests {
//
//		haveStore, haveErr := getInitializedStoreService.InitByToken(nil, test.haveScopeType)
//		assert.Nil(t, haveStore)
//		assert.EqualError(t, store.ErrAppStoreNotSet, haveErr.Error())
//
//		if err := getInitializedStoreService.Init(test.haveSO); err != nil {
//			t.Fatal(err)
//		}
//
//		if s, err := getInitializedStoreService.Store(); err == nil {
//			assert.EqualValues(t, test.wantStoreCode, s.Data.Code.String)
//		} else {
//			assert.EqualError(t, err, store.ErrStoreNotFound.Error())
//			t.Fail()
//		}
//
//		haveStore, haveErr = getInitializedStoreService.InitByToken(getToken(test.haveCodeToken).Claims, test.haveScopeType)
//		if test.wantErr != nil {
//			assert.Nil(t, haveStore, "%#v", test)
//			assert.Error(t, haveErr, "%#v", test)
//			assert.EqualError(t, test.wantErr, haveErr.Error())
//		} else {
//			if test.wantTokenStoreCode != nil {
//				assert.NotNil(t, haveStore, "%#v", test)
//				assert.NoError(t, haveErr)
//				assert.Equal(t, test.wantTokenStoreCode.StoreCode(), haveStore.Data.Code.String)
//			} else {
//				assert.Nil(t, haveStore, "%#v", test)
//				assert.NoError(t, haveErr, "%#v", test)
//			}
//
//		}
//		getInitializedStoreService.ClearCache(true)
//	}
//}
