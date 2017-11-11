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

package jwt_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/cspkg/config/cfgmock"
	"github.com/corestoreio/cspkg/net/jwt"
	"github.com/corestoreio/cspkg/net/mw"
	"github.com/corestoreio/cspkg/storage/containable"
	"github.com/corestoreio/cspkg/store/scope"
	"github.com/corestoreio/cspkg/store/storemock"
	"github.com/corestoreio/cspkg/util/csjwt"
	"github.com/corestoreio/cspkg/util/csjwt/jwtclaim"
	"github.com/corestoreio/cspkg/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
)

func testAuth_WithRunMode(t *testing.T, finalHandler http.Handler, opts ...jwt.Option) (http.Handler, []byte) {
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(append(opts, jwt.WithRootConfig(cfg))...)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	theToken, err := jm.NewToken(scope.DefaultTypeID, jwtclaim.Map{
		"xfoo": "baz",
		"zfoo": 4712,
	})
	assert.NoError(t, err)

	if finalHandler == nil {
		finalHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	}
	srv := storemock.NewEurozzyService(cfg)

	authHandler := jm.WithRunMode(scope.DefaultRunMode, srv)(finalHandler)
	return authHandler, theToken.Raw
}

func TestService_WithRunMode_NoToken(t *testing.T) {

	//  request calls default unauthorized handler

	authHandler, _ := testAuth_WithRunMode(t, nil,
		jwt.WithErrorHandler(mw.ErrorWithPanic),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
	)

	req := httptest.NewRequest("GET", "http://auth.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusUnauthorized)+"\n")
}

func TestService_WithRunMode_Custom_UnauthorizedHandler(t *testing.T) {

	// request calls the unauthorized handler of website scope 1 = euro scope

	var calledUnauthorizedHandler bool
	authHandler, _ := testAuth_WithRunMode(t, nil,
		jwt.WithErrorHandler(mw.ErrorWithPanic),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		jwt.WithUnauthorizedHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tk, ok := jwt.FromContext(r.Context())
				assert.False(t, tk.Valid)
				assert.False(t, ok)

				assert.True(t, errors.IsNotFound(err), "%+v", err)
				w.WriteHeader(http.StatusTeapot)

				wID, sID, ok := scope.FromContext(r.Context())
				assert.True(t, ok)
				assert.Exactly(t, int64(1), wID, "scope.FromContext website")
				assert.Exactly(t, int64(2), sID, "scope.FromContext store")

				calledUnauthorizedHandler = true
			})
		}, scope.Website.Pack(1)),
	)

	req := httptest.NewRequest("GET", "http://auth.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledUnauthorizedHandler)
}

func TestService_WithRunMode_Invalid_ScopedConfiguration(t *testing.T) {
	authHandler, _ := testAuth_WithRunMode(t, nil,
		jwt.WithErrorHandler(mw.ErrorWithPanic, scope.Website.Pack(1)),
		jwt.WithErrorHandler(mw.ErrorWithPanic),
		jwt.WithSigningMethod(nil, scope.Website.Pack(1)),
	)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), `[jwt] ScopedConfig Type(Website) ID(1) is invalid`)
}

func TestService_WithRunMode_Disabled(t *testing.T) {
	authHandler, _ := testAuth_WithRunMode(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			wID, sID, ok := scope.FromContext(r.Context())
			assert.True(t, ok)
			assert.Exactly(t, int64(1), wID, "scope.FromContext website")
			assert.Exactly(t, int64(2), sID, "scope.FromContext store")
		}),
		jwt.WithDisable(true, scope.Website.Pack(1)), // 1 == euro website
		jwt.WithErrorHandler(mw.ErrorWithPanic, scope.Website.Pack(1)),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
	)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestService_WithRunMode_SingleUsage(t *testing.T) {
	authHandler, token := testAuth_WithRunMode(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			wID, sID, ok := scope.FromContext(r.Context())
			assert.True(t, ok)
			assert.Exactly(t, int64(1), wID, "scope.FromContext website")
			assert.Exactly(t, int64(2), sID, "scope.FromContext store")
		}),
		jwt.WithSingleTokenUsage(true, scope.Website.Pack(1)),
		jwt.WithErrorHandler(mw.ErrorWithPanic, scope.Website.Pack(1)),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		// default is a null blacklist so we must set one
		jwt.WithBlacklist(containable.NewInMemory()),
	)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	jwt.SetHeaderAuthorization(req, token)

	// 1st request ok
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Empty(t, w.Body.String())

	// 2nd request unauthorized
	w = httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusUnauthorized)+"\n")
}

func TestService_WithRunMode_DefaultStoreID_Error(t *testing.T) {
	var calledServiceErrorHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithRootConfig(cfg),
		jwt.WithErrorHandler(mw.ErrorWithPanic),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAlreadyReported)
				assert.True(t, errors.IsNotImplemented(err))
				calledServiceErrorHandler = true
			})
		}),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	authHandler := jm.WithRunMode(scope.DefaultRunMode,
		storemock.NewDefaultStoreID(0, 0, errors.NewNotImplementedf("Sorry Dude")),
	)(final)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusAlreadyReported, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledServiceErrorHandler, "calledServiceErrorHandler")
}

func TestService_WithRunMode_StoreIDbyCode_Error(t *testing.T) {
	var calledErrorHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithRootConfig(cfg),
		jwt.WithErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAlreadyReported)
				assert.True(t, errors.IsNotImplemented(err))
				calledErrorHandler = true
			})
		}, scope.Website.Pack(778)),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	authHandler := jm.WithRunMode(scope.DefaultRunMode,
		storemock.NewDefaultStoreID(0, 778, nil, storemock.NewStoreIDbyCode(0, 0, errors.NewNotImplementedf("Sorry Dude"))),
	)(final)

	claimStore := jwtclaim.NewStore()
	claimStore.Store = "'80s FTW"
	theToken, err := jm.NewToken(scope.Website.Pack(778), claimStore)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	jwt.SetHeaderAuthorization(req, theToken.Raw)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusAlreadyReported, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledErrorHandler)
}

func TestService_WithRunMode_IsAllowedStoreID_Error(t *testing.T) {
	var calledErrorHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithRootConfig(cfg),
		jwt.WithErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAlreadyReported)
				assert.True(t, errors.IsTemporary(err))
				calledErrorHandler = true
			})
		}, scope.Website.Pack(778)),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	authHandler := jm.WithRunMode(scope.DefaultRunMode,
		storemock.NewDefaultStoreID(0, 778, nil, storemock.NewStoreIDbyCode(0, 0, errors.NewTemporaryf("Sorry Dude"))),
	)(final)

	claimStore := jwtclaim.NewStore()
	claimStore.Store = "'80s FTW"
	theToken, err := jm.NewToken(scope.Website.Pack(778), claimStore)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	jwt.SetHeaderAuthorization(req, theToken.Raw)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusAlreadyReported, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledErrorHandler)
}

func TestService_WithRunMode_IsAllowedStoreID_Not(t *testing.T) {
	var calledUnauthorizedHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithRootConfig(cfg),
		jwt.WithErrorHandler(mw.ErrorWithPanic, scope.Website.Pack(889)),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		jwt.WithUnauthorizedHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wID, sID, ok := scope.FromContext(r.Context())
				assert.True(t, ok)
				assert.Exactly(t, int64(889), wID, "scope.FromContext website")
				assert.Exactly(t, int64(890), sID, "scope.FromContext store")

				tk, ok := jwt.FromContext(r.Context())
				assert.True(t, tk.Valid)
				assert.True(t, ok)

				assert.True(t, errors.IsUnauthorized(err), "%+v", err)
				w.WriteHeader(http.StatusTeapot)
				calledUnauthorizedHandler = true
			})
		}, scope.Website.Pack(889)),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	authHandler := jm.WithRunMode(scope.DefaultRunMode,
		storemock.NewDefaultStoreID(890, 889, nil, storemock.NewStoreIDbyCode(0, 0, errors.NewNotFoundf("Store code not found"))),
	)(final)

	claimStore := jwtclaim.NewStore()
	claimStore.Store = "'80s FTW"
	theToken, err := jm.NewToken(scope.Website.Pack(889), claimStore)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	jwt.SetHeaderAuthorization(req, theToken.Raw)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledUnauthorizedHandler)

}

func TestService_WithRunMode_AllowedToChangeStore(t *testing.T) {
	var calledFinalHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithRootConfig(cfg),
		jwt.WithErrorHandler(mw.ErrorWithPanic, scope.Website.Pack(359)),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		jwt.WithUnauthorizedHandler(mw.ErrorWithPanic, scope.Website.Pack(359)),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk, ok := jwt.FromContext(r.Context())
		assert.True(t, tk.Valid)
		assert.True(t, ok)

		rm := scope.FromContextRunMode(r.Context())
		assert.Exactly(t, scope.DefaultRunMode, rm, "FromContextRunMode")

		wID, sID, ok := scope.FromContext(r.Context())
		assert.True(t, ok)
		assert.Exactly(t, int64(379), wID, "scope.FromContext website")
		assert.Exactly(t, int64(380), sID, "scope.FromContext store")

		w.WriteHeader(http.StatusTeapot)
		calledFinalHandler = true
	})

	authHandler := jm.WithRunMode(scope.DefaultRunMode,
		storemock.NewDefaultStoreID(360, 359, nil, storemock.NewStoreIDbyCode(380, 379, nil)),
	)(final)

	claimStore := jwtclaim.NewStore()
	claimStore.Store = "'80s FTW"
	theToken, err := jm.NewToken(scope.Website.Pack(359), claimStore)
	assert.NoError(t, err)
	if len(theToken.Raw) == 0 {
		t.Fatalf("Token empty: %#v", theToken)
	}

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	jwt.SetHeaderAuthorization(req, theToken.Raw)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledFinalHandler)

}

// TestService_WithRunMode_DifferentScopes
// all requests with a valid token
// 1. runmode default, so website euro and default store to AT. just call the
// next handler in the chain because the token is valid and scope is euro/at.
// 2. runmode website OZ default store AU valid request with store NZ and must
// change scope to NZ
func TestService_WithRunMode_DifferentScopes(t *testing.T) {

	key := csjwt.WithPasswordRandom()
	hs256, err := csjwt.NewSigningMethodHS256Fast(key)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	hs512, err := csjwt.NewSigningMethodHS512Fast(key)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	var calledFinalHandler = new(int32)
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithRootConfig(cfg),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		jwt.WithErrorHandler(mw.ErrorWithPanic, scope.Website.Pack(1)),
		jwt.WithErrorHandler(mw.ErrorWithPanic, scope.Website.Pack(2)),
		jwt.WithUnauthorizedHandler(mw.ErrorWithPanic, scope.Website.Pack(1)),
		jwt.WithUnauthorizedHandler(mw.ErrorWithPanic, scope.Website.Pack(2)),
		jwt.WithStoreCodeFieldName("euro_store", scope.Website.Pack(1)),
		jwt.WithStoreCodeFieldName("oz_store", scope.Website.Pack(2)),
		jwt.WithSigningMethod(hs256, scope.Website.Pack(1)),
		jwt.WithSigningMethod(hs512, scope.Website.Pack(2)),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wID, sID, hasScope := scope.FromContext(r.Context())
		assert.True(t, hasScope)

		tk, ok := jwt.FromContext(r.Context())
		assert.True(t, tk.Valid)
		assert.True(t, ok)

		switch rm := scope.FromContextRunMode(r.Context()); rm {
		case scope.MakeTypeID(scope.Website, 1):
			assert.Exactly(t, int64(1), wID, "scope.FromContext website EURO")
			assert.Exactly(t, int64(2), sID, "scope.FromContext store AT")
			assert.True(t, len(tk.Raw) > 200, "Host %q", r.Host)
		case scope.MakeTypeID(scope.Website, 2):
			assert.Exactly(t, int64(2), wID, "scope.FromContext website OZ")
			assert.Exactly(t, int64(6), sID, "scope.FromContext store NZ")
			assert.True(t, len(tk.Raw) > 240, "Host %q", r.Host)
		default:
			panic(fmt.Sprintf("Invalid runMode: %s", rm))
		}

		w.WriteHeader(http.StatusTeapot)
		atomic.AddInt32(calledFinalHandler, 1)
	})

	srv := storemock.NewEurozzyService(cfg)
	authHandler := jm.WithRunMode(
		scope.RunModeFunc(func(r *http.Request) scope.TypeID {
			switch r.Host {
			case "scope-euro.xyz":
				return scope.MakeTypeID(scope.Website, 1)
			case "scope-oz.co.nz":
				return scope.MakeTypeID(scope.Website, 2)
			}
			panic(fmt.Sprintf("Unkown host: %q", r.Host))
		}), srv)(final)

	{
		euroClaim := jwtclaim.Map{
			"euro_store": "", // we dont want to change the store
		}

		euroToken, err := jm.NewToken(scope.Website.Pack(1), euroClaim)
		assert.NoError(t, err)
		if len(euroToken.Raw) == 0 {
			t.Fatalf("Euro Token empty: %#v", euroToken)
		}

		req := httptest.NewRequest("GET", "http://scope-euro.xyz", nil)
		jwt.SetHeaderAuthorization(req, euroToken.Raw)
		hpu := cstesting.NewHTTPParallelUsers(3, 5, 100, time.Millisecond)
		hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusTeapot, rec.Code)
			assert.Empty(t, rec.Body.String())
		}
		hpu.ServeHTTP(req, authHandler)
	}

	{
		ozClaim := jwtclaim.Map{
			"oz_store": "nz", // switch to store NZ
		}

		ozToken, err := jm.NewToken(scope.Website.Pack(2), ozClaim)
		assert.NoError(t, err)
		if len(ozToken.Raw) == 0 {
			t.Fatalf("OZ Token empty: %#v", ozToken)
		}
		req := httptest.NewRequest("GET", "http://scope-oz.co.nz", nil)
		jwt.SetHeaderAuthorization(req, ozToken.Raw)
		hpu := cstesting.NewHTTPParallelUsers(3, 5, 100, time.Millisecond)
		hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusTeapot, rec.Code)
			assert.Empty(t, rec.Body.String())
		}
		hpu.ServeHTTP(req, authHandler)
	}

	assert.Exactly(t, int32(30), *calledFinalHandler, "calledFinalHandler 2*(3*5)")
}

// todo add test for form with input field: access_token
