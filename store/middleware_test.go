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

	"reflect"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/net/httputils"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	storemock "github.com/corestoreio/csfw/store/mock"
	"github.com/dgrijalva/jwt-go"
	"github.com/juju/errgo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var middlewareConfigReader *config.MockReader
var middlewareCtxStoreService context.Context

func init() {
	middlewareConfigReader = config.NewMockReader(
		config.WithMockValues(config.MockPV{
			config.MockPathScopeDefault(store.PathRedirectToBase):    1,
			config.MockPathScopeStore(1, store.PathSecureInFrontend): true,
			config.MockPathScopeStore(1, store.PathUnsecureBaseURL):  "http://www.corestore.io/",
			config.MockPathScopeStore(1, store.PathSecureBaseURL):    "https://www.corestore.io/",
		}),
	)

	middlewareCtxStoreService = storemock.NewContextService(
		scope.Option{},
		func(ms *storemock.Storage) {
			ms.MockStore = func() (*store.Store, error) {
				return store.NewStore(
					&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
					store.SetStoreConfig(middlewareConfigReader),
				)
			}
		},
	)
}
func finalHandlerWithValidateBaseURL(t *testing.T) ctxhttp.Handler {
	return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		assert.NotNil(t, ctx)
		assert.NotNil(t, w)
		assert.NotNil(t, r)
		assert.Empty(t, w.Header().Get("Location"))
		return nil
	})
}

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

	err = store.WithValidateBaseURL(mockReader)(finalHandlerWithValidateBaseURL(t)).ServeHTTPContext(context.Background(), w, req)
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

	mw := store.WithValidateBaseURL(mockReader)(finalHandlerWithValidateBaseURL(t))

	err = mw.ServeHTTPContext(context.Background(), w, req)
	assert.EqualError(t, err, store.ErrContextServiceNotFound.Error())

	w = httptest.NewRecorder()
	req, err = http.NewRequest(httputils.MethodPost, "http://corestore.io/catalog/product/view", strings.NewReader(`{ "k1": "v1",  "k2": { "k3": ["va1"]  }}`))
	assert.NoError(t, err)

	err = mw.ServeHTTPContext(context.Background(), w, req)
	assert.NoError(t, err)

}

func TestWithValidateBaseUrl_ActivatedAndShouldRedirectWithGETRequest(t *testing.T) {

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
		mw := store.WithValidateBaseURL(middlewareConfigReader)(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return fmt.Errorf("This handler should not be called! Iindex %d", i)
		}))
		assert.NoError(t, mw.ServeHTTPContext(middlewareCtxStoreService, test.rec, test.req), "Index %d", i)
		assert.Exactly(t, test.wantRedirectURL, test.rec.HeaderMap.Get("Location"), "Index %d", i)
	}
}

func getMWTestRequest(m, u string, c *http.Cookie) *http.Request {
	req, err := http.NewRequest(m, u, nil)
	if err != nil {
		panic(err)
	}
	if c != nil {
		req.AddCookie(c)
	}
	return req
}

func finalInitStoreHandler(t *testing.T, wantStoreCode string) ctxhttp.Handler {
	return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		_, haveReqStore, err := store.FromContextReader(ctx)
		if err != nil {
			return err
		}
		assert.Exactly(t, wantStoreCode, haveReqStore.StoreCode())
		return nil
	})
}

var testsMWInitByFormCookie = []struct {
	req           *http.Request
	haveSO        scope.Option
	wantStoreCode string // this is the default store in a scope, lookup in getInitializedStoreService
	wantErr       error
	wantCookie    string // the newly set cookie
	wantLog       string
}{
	{
		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.ParamName, Value: "uk"}),
		scope.Option{Store: scope.MockID(1)}, "uk", nil, store.ParamName + "=uk;", store.ErrStoreCodeInvalid.Error(),
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
		scope.Option{Store: scope.MockID(1)}, "uk", nil, store.ParamName + "=uk;", "", // generates a new 1year valid cookie
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=%20uk", nil),
		scope.Option{Store: scope.MockID(1)}, "de", nil, "", store.ErrStoreCodeInvalid.Error(),
	},
	{
		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.ParamName, Value: "de"}),
		scope.Option{Group: scope.MockID(1)}, "de", nil, store.ParamName + "=de;", store.ErrStoreCodeInvalid.Error(),
	},
	{
		getMWTestRequest("GET", "http://cs.io", nil),
		scope.Option{Group: scope.MockID(1)}, "at", nil, "", http.ErrNoCookie.Error(),
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=de", nil),
		scope.Option{Group: scope.MockID(1)}, "de", nil, store.ParamName + "=de;", "", // generates a new 1y valid cookie
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=at", nil),
		scope.Option{Group: scope.MockID(1)}, "at", nil, store.ParamName + "=;", "", // generates a delete cookie
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=cz", nil),
		scope.Option{Group: scope.MockID(1)}, "at", store.ErrIDNotFoundTableStoreSlice, "", "",
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
		scope.Option{Group: scope.MockID(1)}, "at", store.ErrStoreChangeNotAllowed, "", "",
	},

	{
		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.ParamName, Value: "nz"}),
		scope.Option{Website: scope.MockID(2)}, "nz", nil, store.ParamName + "=nz;", store.ErrStoreCodeInvalid.Error(),
	},
	{
		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: store.ParamName, Value: "n'z"}),
		scope.Option{Website: scope.MockID(2)}, "au", nil, "", store.ErrStoreCodeInvalid.Error(),
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=uk", nil),
		scope.Option{Website: scope.MockID(2)}, "au", store.ErrStoreChangeNotAllowed, "", "",
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=nz", nil),
		scope.Option{Website: scope.MockID(2)}, "nz", nil, store.ParamName + "=nz;", "",
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=ch", nil),
		scope.Option{Website: scope.MockID(1)}, "at", store.ErrStoreNotActive, "", "",
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+store.HTTPRequestParamStore+"=nz", nil),
		scope.Option{Website: scope.MockID(1)}, "at", store.ErrStoreChangeNotAllowed, "", "",
	},
}

func TestWithInitStoreByFormCookie(t *testing.T) {
	errLogBuf.Reset()
	defer errLogBuf.Reset()

	for i, test := range testsMWInitByFormCookie {

		ctx := store.NewContextReader(context.Background(), getInitializedStoreService(test.haveSO))

		mw := store.WithInitStoreByFormCookie()(finalInitStoreHandler(t, test.wantStoreCode))

		rec := httptest.NewRecorder()
		surfErr := mw.ServeHTTPContext(ctx, rec, test.req)
		if test.wantErr != nil {
			var loc string
			if l, ok := surfErr.(errgo.Locationer); ok {
				loc = l.Location().String()
			}

			assert.EqualError(t, surfErr, test.wantErr.Error(), "\nIndex %d\n%s", i, loc)
			errLogBuf.Reset()
			continue
		}

		if test.wantLog != "" {
			assert.Contains(t, errLogBuf.String(), test.wantLog, "\nIndex %d\n", i)
			errLogBuf.Reset()
			continue
		} else {
			assert.Empty(t, errLogBuf.String(), "\nIndex %d\n", i)
		}

		assert.NoError(t, surfErr, "Index %d", i)

		newKeks := rec.HeaderMap.Get("Set-Cookie")
		if test.wantCookie != "" {
			assert.Contains(t, newKeks, test.wantCookie, "\nIndex %d\n", i)
		} else {
			assert.Empty(t, newKeks, "%#v", test)
		}
		errLogBuf.Reset()
	}
}

func TestWithInitStoreByFormCookie_NilCtx(t *testing.T) {
	mw := store.WithInitStoreByFormCookie()(nil)
	surfErr := mw.ServeHTTPContext(context.Background(), nil, nil)
	assert.EqualError(t, surfErr, store.ErrContextServiceNotFound.Error())
}

func newStoreServiceWithTokenCtx(initO scope.Option, tokenStoreCode string) context.Context {
	ctx := store.NewContextReader(context.Background(), getInitializedStoreService(initO))
	tok := jwt.New(jwt.SigningMethodHS256)
	tok.Claims[store.ParamName] = tokenStoreCode
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
		{store.NewContextReader(context.Background(), nil), "de", store.ErrContextServiceNotFound},
		{store.NewContextReader(context.Background(), getInitializedStoreService(scope.Option{Store: scope.MockCode("de")})), "de", ctxjwt.ErrContextJWTNotFound},
		{newStoreServiceWithTokenCtx(scope.Option{Store: scope.MockCode("de")}, "de"), "de", nil},
		{newStoreServiceWithTokenCtx(scope.Option{Store: scope.MockCode("at")}, "ch"), "at", store.ErrStoreNotActive},
		{newStoreServiceWithTokenCtx(scope.Option{Store: scope.MockCode("de")}, "at"), "at", nil},
		{newStoreServiceWithTokenCtx(scope.Option{Store: scope.MockCode("de")}, "a$t"), "de", store.ErrStoreCodeInvalid},
		{newStoreServiceWithTokenCtx(scope.Option{Store: scope.MockCode("at")}, ""), "at", store.ErrStoreCodeInvalid},

		{newStoreServiceWithTokenCtx(scope.Option{Group: scope.MockID(1)}, "de"), "de", nil},
		{newStoreServiceWithTokenCtx(scope.Option{Group: scope.MockID(1)}, "ch"), "at", store.ErrStoreNotActive},
		{newStoreServiceWithTokenCtx(scope.Option{Group: scope.MockID(1)}, " ch"), "at", store.ErrStoreCodeInvalid},
		{newStoreServiceWithTokenCtx(scope.Option{Group: scope.MockID(1)}, "uk"), "at", store.ErrStoreChangeNotAllowed},

		{newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "uk"), "au", store.ErrStoreChangeNotAllowed},
		{newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "nz"), "nz", nil},
		{newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "n z"), "au", store.ErrStoreCodeInvalid},
		{newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "au"), "au", nil},
		{newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, ""), "au", store.ErrStoreCodeInvalid},
	}
	for i, test := range tests {

		mw := store.WithInitStoreByToken()(finalInitStoreHandler(t, test.wantStoreCode))
		rec := httptest.NewRecorder()
		surfErr := mw.ServeHTTPContext(test.ctx, rec, newReq(i))
		if test.wantErr != nil {
			assert.EqualError(t, surfErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, surfErr, "Index %d", i)
	}
}

func TestWithInitStoreByToken_EqualPointers(t *testing.T) {
	// this Test is related to Benchmark_WithInitStoreByToken
	// The returned pointers from store.FromContextReader must be the
	// same for each request with the same request pattern.

	ctx := newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "nz")
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "https://corestore.io/store/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	var equalStorePointer *store.Store
	mw := store.WithInitStoreByToken()(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		_, haveReqStore, err := store.FromContextReader(ctx)
		if err != nil {
			return err
		}

		if equalStorePointer == nil {
			equalStorePointer = haveReqStore
		}

		if "nz" != haveReqStore.StoreCode() {
			t.Errorf("Have: %s\nWant: nz", haveReqStore.StoreCode())
		}

		wantP := reflect.ValueOf(equalStorePointer)
		haveP := reflect.ValueOf(haveReqStore)

		if wantP.Pointer() != haveP.Pointer() {
			t.Errorf("Expecting equal pointers for each request.\nWant: %p\nHave: %p", equalStorePointer, haveReqStore)
		}

		return nil
	}))

	for i := 0; i < 10; i++ {
		if err := mw.ServeHTTPContext(ctx, rec, req); err != nil {
			t.Error(err)
		}
	}
}
