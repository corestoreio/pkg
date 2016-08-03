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

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/net/runmode"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func getReq(m, t string, c *http.Cookie) *http.Request {
	req := httptest.NewRequest(m, t, nil)
	if c != nil {
		req.AddCookie(c)
	}
	return req
}

type testStoreService struct {
	isAllowed         bool
	allowedCode       string
	allowedErr        error
	defaultStoreID    int64
	defaultStoreIDErr error

	storeIDbyCode    int64
	storeIDbyCodeErr error
}

func (s testStoreService) IsAllowedStoreID(runMode scope.Hash, storeID int64) (bool, string, error) {
	return s.isAllowed, s.allowedCode, s.allowedErr
}
func (s testStoreService) DefaultStoreID(runMode scope.Hash) (int64, error) {
	return s.defaultStoreID, s.defaultStoreIDErr
}
func (s testStoreService) StoreIDbyCode(runMode scope.Hash, storeCode string) (int64, error) {
	return s.storeIDbyCode, s.storeIDbyCodeErr
}

func finalHandler(t *testing.T, wantRunMode scope.Hash) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		haveRunMode := scope.FromContextRunMode(r.Context())
		assert.Exactly(t, wantRunMode, haveRunMode)
	}
}

func TestWithRunMode(t *testing.T) {

	var withRunModeErrH = func(t assert.TestingT, errBhf errors.BehaviourFunc) mw.ErrorHandler {
		return func(haveErr error) http.Handler {
			if errBhf != nil {
				assert.True(t, errBhf(haveErr), "%+v", haveErr)
			} else {
				assert.NoError(t, haveErr, "%+v", haveErr)
			}
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			})
		}
	}

	var testsWithRunMode = []struct {
		req         *http.Request
		storeCheck  store.StoreChecker
		codeIDMap   store.CodeToIDMapper
		options     runmode.Options
		wantRunMode scope.Hash
		wantCookie  string // the newly set cookie
	}{
		{
			getReq("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "uk"}),
			testStoreService{isAllowed: true, allowedCode: "uk", allowedErr: nil, defaultStoreID: 1, defaultStoreIDErr: nil},
			testStoreService{storeIDbyCode: 999, storeIDbyCodeErr: nil},
			runmode.Options{ErrorHandler: withRunModeErrH(t, nil)},
			scope.DefaultRunMode, runmode.FieldName + "=uk;",
		},
		//{
		//	getReq("GET", "http://cs.io/?"+runmode.URLFieldName+"=uk", nil),
		//	scope.Store.ToHash(1), "uk", nil, runmode.FieldName + "=uk;" , // generates a new 1year valid cookie
		//},
		//{
		//	getReq("GET", "http://cs.io/?"+runmode.URLFieldName+"=%20uk", nil),
		//	scope.Store.ToHash(1), "de", nil, "" ,
		//},
		//{
		//	getReq("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "de"}),
		//	scope.Group.ToHash(1), "de", nil, runmode.FieldName + "=de;" ,
		//},
		//{
		//	getReq("GET", "http://cs.io", nil),
		//	scope.Group.ToHash(1), "at", nil, "",
		//},
		//{
		//	getReq("GET", "http://cs.io/?"+runmode.URLFieldName+"=de", nil),
		//	scope.Group.ToHash(1), "de", nil, runmode.FieldName + "=de;",   // generates a new 1y valid cookie
		//},
		//{
		//	getReq("GET", "http://cs.io/?"+runmode.URLFieldName+"=at", nil),
		//	scope.Group.ToHash(1), "at", nil, runmode.FieldName + "=;",   // generates a delete cookie
		//},
		//{
		//	getReq("GET", "http://cs.io/?"+runmode.URLFieldName+"=cz", nil),
		//	scope.Group.ToHash(1), "at", errors.IsNotFound, "",
		//},
		//{
		//	getReq("GET", "http://cs.io/?"+runmode.URLFieldName+"=uk", nil),
		//	scope.Group.ToHash(1), "at", errors.IsUnauthorized, "",
		//},
		//
		//{
		//	getReq("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "nz"}),
		//	scope.Website.ToHash(2), "nz", nil, runmode.FieldName + "=nz;",
		//},
		//{
		//	getReq("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "n'z"}),
		//	scope.Website.ToHash(2), "au", nil, "",
		//},
		//{
		//	getReq("GET", "http://cs.io/?"+runmode.URLFieldName+"=uk", nil),
		//	scope.Website.ToHash(2), "au", errors.IsUnauthorized, "",
		//},
		//{
		//	getReq("GET", "http://cs.io/?"+runmode.URLFieldName+"=nz", nil),
		//	scope.Website.ToHash(2), "nz", nil, runmode.FieldName + "=nz;",
		//},
		//{
		//	getReq("GET", "http://cs.io/?"+runmode.URLFieldName+"=ch", nil),
		//	scope.Website.ToHash(1), "at", errors.IsUnauthorized, "",
		//},
		//{
		//	getReq("GET", "http://cs.io/?"+runmode.URLFieldName+"=nz", nil),
		//	scope.Website.ToHash(1), "at", errors.IsUnauthorized, "",
		//},
	}

	for _, test := range testsWithRunMode {

		mw := runmode.WithRunMode(test.storeCheck, test.codeIDMap, test.options)(finalHandler(t, test.wantRunMode))
		rec := httptest.NewRecorder()
		mw.ServeHTTP(rec, test.req)
		if test.wantCookie != "" {
			assert.Contains(t, rec.Header().Get("Set-Cookie"), test.wantCookie)
		}
	}

}
