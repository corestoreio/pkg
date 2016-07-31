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

package runmode_test

//
//func getMWTestRequest(m, u string, c *http.Cookie) *http.Request {
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
//func finalInitStoreHandler(t *testing.T, wantStoreCode string, wantErrBhf errors.BehaviourFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		haveReqStore, err := store.FromContextRequestedStore(r.Context())
//		if wantErrBhf != nil {
//			assert.True(t, wantErrBhf(err), "\nIndex Error %s", err)
//		}
//		if err != nil {
//			t.Fatal(err)
//		}
//		assert.Exactly(t, wantStoreCode, haveReqStore.StoreCode())
//	}
//}
//
//var testsMWInitByFormCookie = []struct {
//	req           *http.Request
//	haveSO        scope.Option
//	wantStoreCode string // this is the default store in a scope, lookup in getInitializedStoreService
//	wantErrBhf    errors.BehaviourFunc
//	wantCookie    string // the newly set cookie
//	wantLog       string
//}{
//	{
//		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "uk"}),
//		scope.Option{Store: scope.MockID(1)}, "uk", nil, runmode.FieldName + "=uk;", "fix me5",
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io/?"+runmode.URLFieldName+"=uk", nil),
//		scope.Option{Store: scope.MockID(1)}, "uk", nil, runmode.FieldName + "=uk;", "", // generates a new 1year valid cookie
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io/?"+runmode.URLFieldName+"=%20uk", nil),
//		scope.Option{Store: scope.MockID(1)}, "de", nil, "", "fix me4",
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "de"}),
//		scope.Option{Group: scope.MockID(1)}, "de", nil, runmode.FieldName + "=de;", "fix me3",
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io", nil),
//		scope.Option{Group: scope.MockID(1)}, "at", nil, "", http.ErrNoCookie.Error(),
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io/?"+runmode.URLFieldName+"=de", nil),
//		scope.Option{Group: scope.MockID(1)}, "de", nil, runmode.FieldName + "=de;", "", // generates a new 1y valid cookie
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io/?"+runmode.URLFieldName+"=at", nil),
//		scope.Option{Group: scope.MockID(1)}, "at", nil, runmode.FieldName + "=;", "", // generates a delete cookie
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io/?"+runmode.URLFieldName+"=cz", nil),
//		scope.Option{Group: scope.MockID(1)}, "at", errors.IsNotFound, "", "",
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io/?"+runmode.URLFieldName+"=uk", nil),
//		scope.Option{Group: scope.MockID(1)}, "at", errors.IsUnauthorized, "", "",
//	},
//
//	{
//		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "nz"}),
//		scope.Option{Website: scope.MockID(2)}, "nz", nil, runmode.FieldName + "=nz;", "fix me2",
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "n'z"}),
//		scope.Option{Website: scope.MockID(2)}, "au", nil, "", "fix me1",
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io/?"+runmode.URLFieldName+"=uk", nil),
//		scope.Option{Website: scope.MockID(2)}, "au", errors.IsUnauthorized, "", "",
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io/?"+runmode.URLFieldName+"=nz", nil),
//		scope.Option{Website: scope.MockID(2)}, "nz", nil, runmode.FieldName + "=nz;", "",
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io/?"+runmode.URLFieldName+"=ch", nil),
//		scope.Option{Website: scope.MockID(1)}, "at", errors.IsUnauthorized, "", "",
//	},
//	{
//		getMWTestRequest("GET", "http://cs.io/?"+runmode.URLFieldName+"=nz", nil),
//		scope.Option{Website: scope.MockID(1)}, "at", errors.IsUnauthorized, "", "",
//	},
//}
//
//func TestWithInitStoreByFormCookie(t *testing.T) {
//
//	debugLogBuf := new(bytes.Buffer)
//	lg := logw.NewLog(logw.WithWriter(debugLogBuf), logw.WithLevel(logw.LevelDebug))
//
//	for i, test := range testsMWInitByFormCookie {
//
//		srv := storemock.NewEurozzyService(test.haveSO, store.WithStorageConfig(cfgmock.NewService()))
//		dsv, err := srv.DefaultStoreView()
//		ctx := store.WithContextRequestedStore(test.req.Context(), dsv, errors.Wrap(err, "DefaultStoreView"))
//
//		mw := runmode.WithInitStoreByFormCookie(srv, lg)(finalInitStoreHandler(t, test.wantStoreCode, test.wantErrBhf))
//
//		rec := httptest.NewRecorder()
//		mw.ServeHTTP(rec, test.req.WithContext(ctx))
//
//		if test.wantLog != "" {
//			assert.Contains(t, debugLogBuf.String(), test.wantLog, "\nIndex %d\n", i)
//			debugLogBuf.Reset()
//			continue
//		} else {
//			assert.Empty(t, debugLogBuf.String(), "\nIndex %d\n", i)
//		}
//
//		newKeks := rec.HeaderMap.Get("Set-Cookie")
//		if test.wantCookie != "" {
//			assert.Contains(t, newKeks, test.wantCookie, "\nIndex %d\n", i)
//		} else {
//			assert.Empty(t, newKeks, "%#v", test)
//		}
//		debugLogBuf.Reset()
//	}
//}
