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
	"hash/fnv"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/jwt"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/blacklist"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func testAuth_WithRunMode(t *testing.T, finalHandler http.Handler, opts ...jwt.Option) (http.Handler, []byte) {
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(append(opts, jwt.WithConfigGetter(cfg))...)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	theToken, err := jm.NewToken(scope.Default, 0, jwtclaim.Map{
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

	authHandler := jm.WithRunMode(scope.RunMode{}, srv)(finalHandler)
	return authHandler, theToken.Raw
}

func TestService_WithRunMode_NoToken(t *testing.T) {

	//  request calls default unauthorized handler

	authHandler, _ := testAuth_WithRunMode(t, nil,
		jwt.WithErrorHandler(scope.Default, 0, mw.ErrorWithPanic),
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
		jwt.WithErrorHandler(scope.Default, 0, mw.ErrorWithPanic),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		jwt.WithUnauthorizedHandler(scope.Website, 1, func(err error) http.Handler {
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
		}),
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
		jwt.WithErrorHandler(scope.Website, 1, mw.ErrorWithPanic),
		jwt.WithErrorHandler(scope.Default, 0, mw.ErrorWithPanic),
		jwt.WithSigningMethod(scope.Website, 1, nil),
	)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), `[jwt] ScopedConfig Scope(Website) ID(1) is invalid`)
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
		jwt.WithDisable(scope.Website, 1, true), // 1 == euro website
		jwt.WithErrorHandler(scope.Website, 1, mw.ErrorWithPanic),
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
		jwt.WithSingleTokenUsage(scope.Website, 1, true),
		jwt.WithErrorHandler(scope.Website, 1, mw.ErrorWithPanic),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		// default is a null blacklist so we must set one
		jwt.WithBlacklist(blacklist.NewInMemory(fnv.New64a)),
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
		jwt.WithConfigGetter(cfg),
		jwt.WithErrorHandler(scope.Default, 0, mw.ErrorWithPanic),
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

	authHandler := jm.WithRunMode(scope.RunMode{}, storemock.Find{
		StoreIDError: errors.NewNotImplementedf("Sorry Dude"),
	})(final)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusAlreadyReported, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledServiceErrorHandler)
}

func TestService_WithRunMode_StoreIDbyCode_Error(t *testing.T) {
	var calledServiceErrorHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithConfigGetter(cfg),
		jwt.WithErrorHandler(scope.Website, 778, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAlreadyReported)
				assert.True(t, errors.IsNotImplemented(err))
				calledServiceErrorHandler = true
			})
		}),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	authHandler := jm.WithRunMode(scope.RunMode{}, storemock.Find{
		WebsiteIDDefault: 778,
		IDByCodeError:    errors.NewNotImplementedf("Sorry Dude"),
	})(final)

	claimStore := jwtclaim.NewStore()
	claimStore.Store = "'80s FTW"
	theToken, err := jm.NewToken(scope.Website, 778, claimStore)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	jwt.SetHeaderAuthorization(req, theToken.Raw)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusAlreadyReported, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledServiceErrorHandler)
}

func TestService_WithRunMode_IsAllowedStoreID_Error(t *testing.T) {
	var calledServiceErrorHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithConfigGetter(cfg),
		jwt.WithErrorHandler(scope.Website, 778, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAlreadyReported)
				assert.True(t, errors.IsTemporary(err))
				calledServiceErrorHandler = true
			})
		}),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	authHandler := jm.WithRunMode(scope.RunMode{}, storemock.Find{
		WebsiteIDDefault: 778,
		AllowedError:     errors.NewTemporaryf("Sorry Dude"),
	})(final)

	claimStore := jwtclaim.NewStore()
	claimStore.Store = "'80s FTW"
	theToken, err := jm.NewToken(scope.Website, 778, claimStore)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	jwt.SetHeaderAuthorization(req, theToken.Raw)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusAlreadyReported, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledServiceErrorHandler)
}

func TestService_WithRunMode_IsAllowedStoreID_Not(t *testing.T) {
	var calledUnauthorizedHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithConfigGetter(cfg),
		jwt.WithErrorHandler(scope.Website, 889, mw.ErrorWithPanic),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		jwt.WithUnauthorizedHandler(scope.Website, 889, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wID, sID, ok := scope.FromContext(r.Context())
				assert.True(t, ok)
				assert.Exactly(t, int64(901), wID, "scope.FromContext website")
				assert.Exactly(t, int64(902), sID, "scope.FromContext store")

				tk, ok := jwt.FromContext(r.Context())
				assert.True(t, tk.Valid)
				assert.True(t, ok)

				assert.True(t, errors.IsUnauthorized(err), "%+v", err)
				w.WriteHeader(http.StatusTeapot)
				calledUnauthorizedHandler = true
			})
		}),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	authHandler := jm.WithRunMode(scope.RunMode{}, storemock.Find{
		WebsiteIDDefault: 889,
		StoreIDDefault:   890,

		IDByCodeWebsiteID: 901,
		IDByCodeStoreID:   902,

		Allowed:     false, // important
		AllowedCode: "uninteresting",
	})(final)

	claimStore := jwtclaim.NewStore()
	claimStore.Store = "'80s FTW"
	theToken, err := jm.NewToken(scope.Website, 889, claimStore)
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
	var calledUnauthorizedHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithConfigGetter(cfg),
		jwt.WithErrorHandler(scope.Website, 359, mw.ErrorWithPanic),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		jwt.WithUnauthorizedHandler(scope.Website, 359, mw.ErrorWithPanic),
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
		calledUnauthorizedHandler = true
	})

	authHandler := jm.WithRunMode(scope.RunMode{}, storemock.Find{
		WebsiteIDDefault: 359,
		StoreIDDefault:   360,

		IDByCodeWebsiteID: 379,
		IDByCodeStoreID:   380,

		Allowed:     true, // important
		AllowedCode: "uninteresting",
	})(final)

	claimStore := jwtclaim.NewStore()
	claimStore.Store = "'80s FTW"
	theToken, err := jm.NewToken(scope.Website, 359, claimStore)
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
	assert.True(t, calledUnauthorizedHandler)

}

// TestService_WithRunMode_Disabled
// 1. a request with a valid token
// should do nothing and just call the next handler in the chain because the
// token service has been disabled therefore no validation and checking takes
// place.
// 2. valid request with website oz must be passed through with an
// invalid token because JWT disabled
func TestService_WithRunMode_DifferentScopes(t *testing.T) {
	t.Skip("todo")
	var calledUnauthorizedHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithConfigGetter(cfg),
		jwt.WithErrorHandler(scope.Website, 359, mw.ErrorWithPanic),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		jwt.WithUnauthorizedHandler(scope.Website, 359, mw.ErrorWithPanic),
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
		calledUnauthorizedHandler = true
	})

	authHandler := jm.WithRunMode(scope.RunMode{}, storemock.Find{
		WebsiteIDDefault: 359,
		StoreIDDefault:   360,

		IDByCodeWebsiteID: 379,
		IDByCodeStoreID:   380,

		Allowed:     true, // important
		AllowedCode: "uninteresting",
	})(final)

	claimStore := jwtclaim.NewStore()
	claimStore.Store = "'80s FTW"
	theToken, err := jm.NewToken(scope.Website, 359, claimStore)
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
	assert.True(t, calledUnauthorizedHandler)

}

// todo add test for form with input field: access_token
