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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/net/runmode"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func getReq(m, t string, c *http.Cookie) *http.Request {
	req := httptest.NewRequest(m, t, nil)
	if c != nil {
		req.AddCookie(c)
	}
	return req
}

func finalHandler(t *testing.T, wantRunMode scope.TypeID, wantStoreID, wantWebsiteID int64, wantStoreCtx bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		haveWebsiteID, haveStoreID, haveOK := scope.FromContext(r.Context())
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
		wantRunMode   scope.TypeID
		wantCookie    string // the newly set cookie
		wantRespCode  int
		wantStoreID   int64
		wantWebsiteID int64
		wantCtx       bool
	}{
		// test cases: DefaultRunMode
		{ // request with cookie store UK. store valid; don't set a cookie
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "uk"}),
			storemock.NewDefaultStoreID(999, 111, nil, storemock.NewStoreIDbyCode(888, 222, nil)),
			runmode.Options{ErrorHandler: withRunModeErrH(t, nil, false)},
			scope.DefaultRunMode, "", http.StatusAccepted,
			888, 222, true,
		},
		{ // request with store cookie UK. should delete cookie and trigger error because store not allowed
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "uk"}),
			storemock.NewDefaultStoreID(1, 1, nil, storemock.NewStoreIDbyCode(0, 0, errors.NewNotFoundf("UK not found"))),
			runmode.Options{ErrorHandler: withRunModeErrH(t, errors.IsUnauthorized, true)},
			scope.DefaultRunMode, store.CodeFieldName + `=; Path=/`, http.StatusUnauthorized,
			0, 0, false,
		},
		{ // request with store cookie UK. should delete cookie (because store == DefaultStoreID) and store allowed
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "uk"}),
			storemock.NewDefaultStoreID(135, 136, nil, storemock.NewStoreIDbyCode(135, 136, nil)),
			runmode.Options{ErrorHandler: withRunModeErrH(t, errors.IsUnauthorized, false)},
			scope.DefaultRunMode, store.CodeFieldName + `=; Path=/`, http.StatusAccepted,
			135, 136, true,
		},
		{ // request with store cookie UK; fails because DefaultStoreID returns an error
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "uk"}),
			storemock.NewDefaultStoreID(0, 0, errors.NewNotImplementedf("Upsss!")),
			runmode.Options{ErrorHandler: withRunModeErrH(t, errors.IsNotImplemented, false)},
			scope.DefaultRunMode, ``, http.StatusServiceUnavailable,
			0, 0, false,
		},
		{ // request with store GET param UK; fails because StoreIDbyCode returns an error
			getReq("GET", fmt.Sprintf("http://cs.io?x=%%20y&%s=uk", store.CodeURLFieldName), nil),
			storemock.NewDefaultStoreID(1, 1, nil, storemock.NewStoreIDbyCode(0, 1, errors.NewFatalf("No idea what's fatal ..."))),
			runmode.Options{ErrorHandler: withRunModeErrH(t, errors.IsFatal, false)},
			scope.DefaultRunMode, ``, http.StatusServiceUnavailable,
			0, 0, false,
		},
		{ // request with store GET param U K; ignores invalid store, and sets no cookie
			getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=u%%20k", store.CodeURLFieldName), nil),
			storemock.NewDefaultStoreID(165, 166, nil),
			runmode.Options{ErrorHandler: withRunModeErrH(t, nil, false)},
			scope.DefaultRunMode, "", http.StatusAccepted,
			165, 166, true,
		},
		{ // request with store GET param UK and sets cookie with code uk, because store code was provided via GET
			getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=uk", store.CodeURLFieldName), nil),
			storemock.NewDefaultStoreID(175, 176, nil, storemock.NewStoreIDbyCode(177, 178, nil)),
			runmode.Options{ErrorHandler: withRunModeErrH(t, nil, false)},
			scope.DefaultRunMode, store.CodeFieldName + `=uk; Path=/`, http.StatusAccepted,
			177, 178, true,
		},
		{ // request with store GET param FR; fails because IsAllowedStoreID returns an error
			getReq("GET", fmt.Sprintf("http://cs.io?e=f&%s=fr", store.CodeURLFieldName), nil),
			storemock.NewDefaultStoreID(1, 1, nil, storemock.NewStoreIDbyCode(0, 0, errors.NewAlreadyClosedf("Not in the mood"))),
			runmode.Options{ErrorHandler: withRunModeErrH(t, errors.IsAlreadyClosed, false)},
			scope.DefaultRunMode, ``, http.StatusServiceUnavailable,
			0, 0, false,
		},

		// website runmode
		{ // request with store cookie cn does nothing
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "cn"}),
			storemock.NewDefaultStoreID(0, 0, nil, storemock.NewStoreIDbyCode(44, 33, nil)),
			runmode.Options{ErrorHandler: withRunModeErrH(t, nil, false), RunModeCalculater: scope.Website.Pack(2)},
			scope.Website.Pack(2), ``, http.StatusAccepted,
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

func TestWithRunMode_StoreService(t *testing.T) {
	srv := storemock.NewEurozzyService(cfgmock.NewService())

	tests := []struct {
		req           *http.Request
		options       runmode.Options
		wantStoreID   int64
		wantWebsiteID int64
		wantCtx       bool
		wantRespCode  int
	}{
		{
			getReq("GET", "http://cs.io", nil),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Website, 1)},
			2, 1, true, http.StatusAccepted, // at
		},
		{
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "de"}),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Website, 1)},
			1, 1, true, http.StatusAccepted, // de
		},
		{
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "uk"}),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Website, 1)},
			4, 1, true, http.StatusAccepted, // uk
		},
		{
			getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=ch", store.CodeURLFieldName), nil),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Website, 1)},
			2, 1, true, http.StatusUnauthorized, // ch, but inactive
		},

		{
			getReq("GET", "http://cs.io", nil),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Group, 3)},
			5, 2, true, http.StatusAccepted, // au
		},
		{
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "nz"}),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Group, 3)},
			6, 2, true, http.StatusAccepted, // nz
		},
		{
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "de"}),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Group, 3)},
			5, 2, true, http.StatusUnauthorized, // requesting store DE but not allowed
		},

		{
			getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=", store.CodeURLFieldName), nil),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Store, 1)},
			1, 1, true, http.StatusAccepted, // de
		},
		{
			getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=at", store.CodeURLFieldName), nil),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Store, 1)},
			2, 1, true, http.StatusAccepted, // at
		},
		{
			getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=nz", store.CodeURLFieldName), nil),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Store, 1)},
			6, 2, true, http.StatusAccepted, // nz
		},
		{
			getReq("GET", fmt.Sprintf("http://cs.io?x=y&%s=ch", store.CodeURLFieldName), nil),
			runmode.Options{RunModeCalculater: scope.MakeTypeID(scope.Store, 1)},
			1, 1, true, http.StatusUnauthorized, // ch not active
		},

		{
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "dXe"}),
			runmode.Options{},
			2, 1, true, http.StatusUnauthorized, // dXe non-existent
		},
		{
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "de"}),
			runmode.Options{},
			1, 1, true, http.StatusAccepted, // website euro && store de
		},
		{
			getReq("GET", "http://cs.io", &http.Cookie{Name: store.CodeFieldName, Value: "nz"}),
			runmode.Options{},
			2, 1, true, http.StatusUnauthorized, // at, switch to nz not allowed
		},
	}
	for i, test := range tests {
		test.options.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

		rmmw := runmode.WithRunMode(srv, test.options)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			haveWebsiteID, haveStoreID, haveOK := scope.FromContext(r.Context())
			assert.Exactly(t, test.wantCtx, haveOK, "Context; Index %d", i)
			assert.Exactly(t, test.wantStoreID, haveStoreID, "Store; Index %d", i)
			assert.Exactly(t, test.wantWebsiteID, haveWebsiteID, "Website; Index %d", i)
			assert.NotEmpty(t, scope.FromContextRunMode(r.Context()), "Index %d", i)
			w.WriteHeader(http.StatusAccepted)
		}))
		rec := httptest.NewRecorder()
		rmmw.ServeHTTP(rec, test.req)
		assert.Exactly(t, test.wantRespCode, rec.Code, "Index %d => %s", i, rec.Body)
	}
}
