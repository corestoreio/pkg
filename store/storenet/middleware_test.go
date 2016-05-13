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

package storenet_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/store/storenet"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
	"github.com/stretchr/testify/assert"
)

//var middlewareConfigReader *cfgmock.Service
//var middlewareCtxStoreService context.Context
//
//func init() {
//	middlewareConfigReader = cfgmock.NewService(
//		cfgmock.WithPV(cfgmock.PathValue{
//			scope.StrDefault.FQPathInt64(0, backend.Backend.WebURLRedirectToBase.String()):  1,
//			scope.StrStores.FQPathInt64(1, backend.Backend.WebSecureUseInFrontend.String()): true,
//			scope.StrStores.FQPathInt64(1, backend.Backend.WebUnsecureBaseURL.String()):     "http://www.corestore.io/",
//			scope.StrStores.FQPathInt64(1, backend.Backend.WebSecureBaseURL.String()):       "https://www.corestore.io/",
//		}),
//	)
//
//	middlewareCtxStoreService = storemock.WithContextMustService(
//		scope.Option{},
//		func(ms *storemock.Storage) {
//			ms.MockStore = func() (*store.Store, error) {
//				return store.NewStore(
//					&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
//					&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
//					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
//					store.WithStoreConfig(middlewareConfigReader),
//				)
//			}
//		},
//	)
//}
//func finalHandlerWithValidateBaseURL(t *testing.T) ctxhttp.HandlerFunc {
//	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//		assert.NotNil(t, ctx)
//		assert.NotNil(t, w)
//		assert.NotNil(t, r)
//		assert.Empty(t, w.Header().Get("Location"))
//		return nil
//	}
//}
//
//func TestWithValidateBaseUrl_DeactivatedAndShouldNotRedirectWithGETRequest(t *testing.T) {
//
//	mockReader := cfgmock.NewService(
//		cfgmock.WithPV(cfgmock.PathValue{
//			scope.StrDefault.FQPathInt64(0, backend.Backend.WebURLRedirectToBase.String()): 0,
//		}),
//	)
//
//	// no post request but check deactivated
//	w := httptest.NewRecorder()
//	req, err := http.NewRequest(httputil.MethodGet, "http://corestore.io/catalog/product/view", nil)
//	assert.NoError(t, err)
//
//	err = storenet.WithValidateBaseURL(mockReader)(finalHandlerWithValidateBaseURL(t)).ServeHTTPContext(context.Background(), w, req)
//	assert.NoError(t, err)
//}
//
//func TestWithValidateBaseUrl_ActivatedAndShouldNotRedirectWithPOSTRequest(t *testing.T) {
//
//	mockReader := cfgmock.NewService(
//		cfgmock.WithPV(cfgmock.PathValue{
//			scope.StrDefault.FQPathInt64(0, backend.Backend.WebURLRedirectToBase.String()): 301,
//		}),
//	)
//
//	w := httptest.NewRecorder()
//	req, err := http.NewRequest(httputil.MethodGet, "http://corestore.io/catalog/product/view", nil)
//	assert.NoError(t, err)
//
//	mw := storenet.WithValidateBaseURL(mockReader)(finalHandlerWithValidateBaseURL(t))
//
//	err = mw.ServeHTTPContext(context.Background(), w, req)
//	assert.EqualError(t, err, store.errContextProviderNotFound.Error())
//
//	w = httptest.NewRecorder()
//	req, err = http.NewRequest(httputil.MethodPost, "http://corestore.io/catalog/product/view", strings.NewReader(`{ "k1": "v1",  "k2": { "k3": ["va1"]  }}`))
//	assert.NoError(t, err)
//
//	err = mw.ServeHTTPContext(context.Background(), w, req)
//	assert.NoError(t, err)
//
//}

//func TestWithValidateBaseUrl_ActivatedAndShouldRedirectWithGETRequest(t *testing.T) {
//
//	var newReq = func(urlStr string) *http.Request {
//		req, err := http.NewRequest(httputil.MethodGet, urlStr, nil)
//		if err != nil {
//			t.Fatal(err)
//		}
//		return req
//	}
//
//	tests := []struct {
//		rec             *httptest.ResponseRecorder
//		req             *http.Request
//		wantRedirectURL string
//	}{
//		{
//			httptest.NewRecorder(),
//			newReq("http://corestore.io/catalog/product/view/"),
//			"http://www.corestore.io/catalog/product/view/",
//		},
//		{
//			httptest.NewRecorder(),
//			newReq("http://corestore.io/catalog/product/view"),
//			"http://www.corestore.io/catalog/product/view",
//		},
//		{
//			httptest.NewRecorder(),
//			newReq("http://corestore.io"),
//			"http://www.corestore.io",
//		},
//		{
//			httptest.NewRecorder(),
//			newReq("https://corestore.io/catalog/category/view?catid=1916"),
//			"https://www.corestore.io/catalog/category/view?catid=1916",
//		},
//		{
//			httptest.NewRecorder(),
//			newReq("https://corestore.io/customer/comments/view?id=1916#tab=ratings"),
//			"https://www.corestore.io/customer/comments/view?id=1916#tab=ratings",
//		},
//	}
//
//	for i, test := range tests {
//		mw := storenet.WithValidateBaseURL(middlewareConfigReader)(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//			return fmt.Errorf("This handler should not be called! Iindex %d", i)
//		}))
//		assert.NoError(t, mw.ServeHTTPContext(middlewareCtxStoreService, test.rec, test.req), "Index %d", i)
//		assert.Exactly(t, test.wantRedirectURL, test.rec.HeaderMap.Get("Location"), "Index %d", i)
//	}
//}

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

func finalInitStoreHandler(t *testing.T, wantStoreCode string, wantErrBhf errors.BehaviourFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		haveReqStore, err := store.FromContextRequestedStore(r.Context())
		if wantErrBhf != nil {
			assert.True(t, wantErrBhf(err), "\nIndex Error %s", err)
		}
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, wantStoreCode, haveReqStore.StoreCode())
	}
}

var testsMWInitByFormCookie = []struct {
	req           *http.Request
	haveSO        scope.Option
	wantStoreCode string // this is the default store in a scope, lookup in getInitializedStoreService
	wantErrBhf    errors.BehaviourFunc
	wantCookie    string // the newly set cookie
	wantLog       string
}{
	{
		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: storenet.ParamName, Value: "uk"}),
		scope.Option{Store: scope.MockID(1)}, "uk", nil, storenet.ParamName + "=uk;", "fix me5",
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+storenet.HTTPRequestParamStore+"=uk", nil),
		scope.Option{Store: scope.MockID(1)}, "uk", nil, storenet.ParamName + "=uk;", "", // generates a new 1year valid cookie
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+storenet.HTTPRequestParamStore+"=%20uk", nil),
		scope.Option{Store: scope.MockID(1)}, "de", nil, "", "fix me4",
	},
	{
		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: storenet.ParamName, Value: "de"}),
		scope.Option{Group: scope.MockID(1)}, "de", nil, storenet.ParamName + "=de;", "fix me3",
	},
	{
		getMWTestRequest("GET", "http://cs.io", nil),
		scope.Option{Group: scope.MockID(1)}, "at", nil, "", http.ErrNoCookie.Error(),
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+storenet.HTTPRequestParamStore+"=de", nil),
		scope.Option{Group: scope.MockID(1)}, "de", nil, storenet.ParamName + "=de;", "", // generates a new 1y valid cookie
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+storenet.HTTPRequestParamStore+"=at", nil),
		scope.Option{Group: scope.MockID(1)}, "at", nil, storenet.ParamName + "=;", "", // generates a delete cookie
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+storenet.HTTPRequestParamStore+"=cz", nil),
		scope.Option{Group: scope.MockID(1)}, "at", errors.IsNotFound, "", "",
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+storenet.HTTPRequestParamStore+"=uk", nil),
		scope.Option{Group: scope.MockID(1)}, "at", errors.IsUnauthorized, "", "",
	},

	{
		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: storenet.ParamName, Value: "nz"}),
		scope.Option{Website: scope.MockID(2)}, "nz", nil, storenet.ParamName + "=nz;", "fix me2",
	},
	{
		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: storenet.ParamName, Value: "n'z"}),
		scope.Option{Website: scope.MockID(2)}, "au", nil, "", "fix me1",
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+storenet.HTTPRequestParamStore+"=uk", nil),
		scope.Option{Website: scope.MockID(2)}, "au", errors.IsUnauthorized, "", "",
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+storenet.HTTPRequestParamStore+"=nz", nil),
		scope.Option{Website: scope.MockID(2)}, "nz", nil, storenet.ParamName + "=nz;", "",
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+storenet.HTTPRequestParamStore+"=ch", nil),
		scope.Option{Website: scope.MockID(1)}, "at", errors.IsUnauthorized, "", "",
	},
	{
		getMWTestRequest("GET", "http://cs.io/?"+storenet.HTTPRequestParamStore+"=nz", nil),
		scope.Option{Website: scope.MockID(1)}, "at", errors.IsUnauthorized, "", "",
	},
}

func TestWithInitStoreByFormCookie(t *testing.T) {

	debugLogBuf := new(bytes.Buffer)
	lg := log.NewStdLog(log.WithStdWriter(debugLogBuf), log.WithStdLevel(log.StdLevelDebug))

	for i, test := range testsMWInitByFormCookie {

		srv := storemock.NewEurozzyService(test.haveSO, store.WithStorageConfig(cfgmock.NewService()))
		dsv, err := srv.DefaultStoreView()
		ctx := store.WithContextRequestedStore(test.req.Context(), dsv, errors.Wrap(err, "DefaultStoreView"))

		mw := storenet.WithInitStoreByFormCookie(srv, lg)(finalInitStoreHandler(t, test.wantStoreCode, test.wantErrBhf))

		rec := httptest.NewRecorder()
		mw.ServeHTTP(rec, test.req.WithContext(ctx))

		if test.wantLog != "" {
			assert.Contains(t, debugLogBuf.String(), test.wantLog, "\nIndex %d\n", i)
			debugLogBuf.Reset()
			continue
		} else {
			assert.Empty(t, debugLogBuf.String(), "\nIndex %d\n", i)
		}

		newKeks := rec.HeaderMap.Get("Set-Cookie")
		if test.wantCookie != "" {
			assert.Contains(t, newKeks, test.wantCookie, "\nIndex %d\n", i)
		} else {
			assert.Empty(t, newKeks, "%#v", test)
		}
		debugLogBuf.Reset()
	}
}
