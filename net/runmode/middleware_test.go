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

	"fmt"
	"github.com/corestoreio/csfw/log"
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

var _ store.Finder = (*testStoreService)(nil)

type testStoreService struct {
	isAllowed   bool
	allowedCode string
	allowedErr  error

	defaultStoreID    int64
	defaultWebsiteID  int64
	defaultStoreIDErr error

	idByCodeStore   int64
	idByCodeWebsite int64
	idByCodeErr     error
}

func (s testStoreService) IsAllowedStoreID(runMode scope.Hash, storeID int64) (bool, string, error) {
	return s.isAllowed, s.allowedCode, s.allowedErr
}
func (s testStoreService) DefaultStoreID(runMode scope.Hash) (int64, int64, error) {
	return s.defaultStoreID, s.defaultWebsiteID, s.defaultStoreIDErr
}
func (s testStoreService) StoreIDbyCode(runMode scope.Hash, storeCode string) (int64, int64, error) {
	return s.idByCodeStore, s.idByCodeWebsite, s.idByCodeErr
}

func finalHandler(t *testing.T, wantRunMode scope.Hash, wantStoreID, wantWebsiteID int64, wantStoreCtx bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		haveStoreID, haveWebsiteID, haveOK := scope.FromContext(r.Context())
		assert.Exactly(t, wantStoreCtx, haveOK)
		assert.Exactly(t, wantStoreID, haveStoreID)
		assert.Exactly(t, wantWebsiteID, haveWebsiteID)

		haveRunMode := scope.FromContextRunMode(r.Context())
		assert.Exactly(t, wantRunMode, haveRunMode)
		w.WriteHeader(http.StatusAccepted)
	}
}

func TestWithRunMode(t *testing.T) {

	var withRunModeErrH = func(t assert.TestingT, errBhf errors.BehaviourFunc, wantStoreIDCtx bool) mw.ErrorHandler {
		return func(haveErr error) http.Handler {
			code := http.StatusNoContent // just the default OK
			if errBhf != nil {
				assert.True(t, errBhf(haveErr), "%+v", haveErr)
				code = http.StatusServiceUnavailable
			} else {
				assert.NoError(t, haveErr, "%+v", haveErr)
			}
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _, haveOK := scope.FromContext(r.Context())
				assert.Exactly(t, wantStoreIDCtx, haveOK, "Scope Context in Error Handler")
				w.WriteHeader(code)
			})
		}
	}

	var testsWithRunMode = []struct {
		req           *http.Request
		storeFinder   store.Finder
		options       runmode.Options
		wantRunMode   scope.Hash
		wantCookie    string // the newly set cookie
		wantRespCode  int
		wantStoreID   int64
		wantWebsiteID int64
		wantCtx       bool
	}{
		// test cases: DefaultRunMode
		{ // request with cookie store UK. store valid; don't set a cookie
			getReq("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "uk"}),
			testStoreService{
				isAllowed: true, allowedCode: "uk", allowedErr: nil,
				defaultStoreID: 999, defaultWebsiteID: 111, defaultStoreIDErr: nil,
				idByCodeStore: 888, idByCodeWebsite: 222, idByCodeErr: nil},
			runmode.Options{ErrorHandler: withRunModeErrH(t, nil, false)},
			scope.DefaultRunMode, "", http.StatusAccepted,
			888, 222, true,
		},
		{ // request with store cookie UK. should delete cookie and trigger error because store not allowed
			getReq("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "uk"}),
			testStoreService{
				isAllowed: false, allowedCode: "", allowedErr: nil,
				defaultStoreID: 1, defaultWebsiteID: 1, defaultStoreIDErr: nil,
				idByCodeStore: 999, idByCodeWebsite: 888, idByCodeErr: nil},
			runmode.Options{ErrorHandler: withRunModeErrH(t, errors.IsUnauthorized, true)},
			scope.DefaultRunMode, runmode.FieldName + `=; Path=/`, http.StatusServiceUnavailable,
			0, 0, false,
		},
		{ // request with store cookie UK. should delete cookie (because store == DefaultStoreID) and store allowed
			getReq("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "uk"}),
			testStoreService{
				isAllowed: true, allowedCode: "uk", allowedErr: nil,
				defaultStoreID: 135, defaultWebsiteID: 136, defaultStoreIDErr: nil,
				idByCodeStore: 135, idByCodeWebsite: 136, idByCodeErr: nil},
			runmode.Options{ErrorHandler: withRunModeErrH(t, errors.IsUnauthorized, false)},
			scope.DefaultRunMode, runmode.FieldName + `=; Path=/`, http.StatusAccepted,
			135, 136, true,
		},
		{ // request with store cookie UK; fails because DefaultStoreID returns an error
			getReq("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "uk"}),
			testStoreService{
				isAllowed: false, allowedCode: "", allowedErr: nil,
				defaultStoreID: 0, defaultWebsiteID: 0, defaultStoreIDErr: errors.NewNotImplementedf("Upsss!"),
				idByCodeStore: 0, idByCodeWebsite: 0, idByCodeErr: nil},
			runmode.Options{ErrorHandler: withRunModeErrH(t, errors.IsNotImplemented, false)},
			scope.DefaultRunMode, ``, http.StatusServiceUnavailable,
			0, 0, false,
		},
		{ // request with store GET param UK; fails because StoreIDbyCode returns an error
			getReq("GET", fmt.Sprintf("http://cs.io?x=%%20y&%s=uk", runmode.URLFieldName), nil),
			testStoreService{
				isAllowed: true, allowedCode: "uk", allowedErr: nil,
				defaultStoreID: 1, defaultWebsiteID: 1, defaultStoreIDErr: nil,
				idByCodeStore: 0, idByCodeWebsite: 1, idByCodeErr: errors.NewFatalf("No idea what's fatal ...")},
			runmode.Options{ErrorHandler: withRunModeErrH(t, errors.IsFatal, false)},
			scope.DefaultRunMode, ``, http.StatusServiceUnavailable,
			0, 0, false,
		},
		{ // request with store GET param U K; ignores invalid store, and sets no cookie
			getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=u%%20k", runmode.URLFieldName), nil),
			testStoreService{
				isAllowed: true, allowedCode: "gb", allowedErr: nil,
				defaultStoreID: 165, defaultWebsiteID: 166, defaultStoreIDErr: nil,
				idByCodeStore: 0, idByCodeWebsite: 0, idByCodeErr: nil},
			runmode.Options{ErrorHandler: withRunModeErrH(t, nil, false)},
			scope.DefaultRunMode, "", http.StatusAccepted,
			165, 166, true,
		},
		{ // request with store GET param UK and sets cookie with new code gb
			getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=uk", runmode.URLFieldName), nil),
			testStoreService{
				isAllowed: true, allowedCode: "gb", allowedErr: nil,
				defaultStoreID: 175, defaultWebsiteID: 176, defaultStoreIDErr: nil,
				idByCodeStore: 177, idByCodeWebsite: 178, idByCodeErr: nil},
			runmode.Options{ErrorHandler: withRunModeErrH(t, nil, false)},
			scope.DefaultRunMode, runmode.FieldName + `=gb; Path=/`, http.StatusAccepted,
			177, 178, true,
		},
		{ // request; fails because IsAllowedStoreID returns an error
			getReq("GET", "http://cs.io", nil),
			testStoreService{
				isAllowed: false, allowedCode: "", allowedErr: errors.NewAlreadyClosedf("Not in the mood"),
				defaultStoreID: 1, defaultWebsiteID: 1, defaultStoreIDErr: nil,
				idByCodeStore: 0, idByCodeWebsite: 0, idByCodeErr: nil},
			runmode.Options{ErrorHandler: withRunModeErrH(t, errors.IsAlreadyClosed, true)},
			scope.DefaultRunMode, ``, http.StatusServiceUnavailable,
			0, 0, false,
		},

		// website runmode
		{ // request with store cookie cn ...
			getReq("GET", "http://cs.io", &http.Cookie{Name: runmode.FieldName, Value: "cn"}),
			testStoreService{
				isAllowed: true, allowedCode: "cn", allowedErr: nil,
				defaultStoreID: 0, defaultWebsiteID: 0, defaultStoreIDErr: nil,
				idByCodeStore: 44, idByCodeWebsite: 33, idByCodeErr: nil},
			runmode.Options{ErrorHandler: withRunModeErrH(t, nil, false), RunMode: scope.RunMode{Mode: scope.Website.ToHash(2)}},
			scope.Website.ToHash(2), ``, http.StatusAccepted,
			44, 33, true,
		},
	}

	for i, test := range testsWithRunMode {
		test.options.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

		rmmw := runmode.WithRunMode(test.storeFinder, test.options)(finalHandler(t, test.wantRunMode, test.wantStoreID, test.wantWebsiteID, test.wantCtx))
		rec := httptest.NewRecorder()
		rmmw.ServeHTTP(rec, test.req)
		if c := rec.Header().Get("Set-Cookie"); test.wantCookie != "" {
			assert.Contains(t, c, test.wantCookie, "Index %d", i)
		} else {
			assert.Empty(t, c, "Index %d", i)
		}
		assert.Exactly(t, http.StatusText(test.wantRespCode), http.StatusText(rec.Code), "Index %d", i)
	}

}
