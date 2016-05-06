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

package jwtauth_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/net/jwtauth"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareWithInitTokenNoStoreProvider(t *testing.T) {

	authHandler, _ := testAuth(t, jwtauth.WithErrorHandler(scope.Default, 0, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk, err := jwtauth.FromContext(r.Context())
		assert.False(t, tk.Valid)
		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	})))

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestMiddlewareWithInitTokenNoToken(t *testing.T) {

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		store.WithStorageConfig(cr),
	)
	ctx := store.WithContextProvider(context.Background(), srv)
	authHandler, _ := testAuth(t, jwtauth.WithErrorHandler(scope.Default, 0, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk, err := jwtauth.FromContext(r.Context())
		assert.False(t, tk.Valid)
		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})))

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized)+"\n", w.Body.String())
}

func TestMiddlewareWithInitTokenHTTPErrorHandler(t *testing.T) {

	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		store.WithStorageConfig(cfgmock.NewService()), // empty config
	)

	authHandler, _ := testAuth(t, jwtauth.WithErrorHandler(scope.Default, 0, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok, err := jwtauth.FromContext(r.Context())
		assert.False(t, tok.Valid)
		w.WriteHeader(http.StatusTeapot)
		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			t.Fatal(err)
		}
	})))

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	ctx := store.WithContextProvider(context.Background(), srv)
	authHandler.ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Contains(t, w.Body.String(), `token not present in request: Not found`)
}

func TestMiddlewareWithInitTokenSuccess(t *testing.T) {

	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		store.WithStorageConfig(cfgmock.NewService()),
	)

	jwts := jwtauth.MustNewService()

	if err := jwts.Options(jwtauth.WithErrorHandler(scope.Default, 0,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := jwtauth.FromContext(r.Context())
			t.Logf("Token: %#v\n", token)
			t.Fatal(errors.PrintLoc(err))
		}),
	)); err != nil {
		t.Fatal(errors.PrintLoc(err))
	}

	theToken, err := jwts.NewToken(scope.Default, 0, jwtclaim.Map{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken.Raw)

	req, err := http.NewRequest("GET", "http://corestore.io/customer/account", nil)
	assert.NoError(t, err)
	jwtauth.SetHeaderAuthorization(req, theToken.Raw)

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprintf(w, "I'm more of a coffee pot")

		ctxToken, err := jwtauth.FromContext(r.Context())
		assert.NoError(t, err)
		assert.NotNil(t, ctxToken)
		xFoo, err := ctxToken.Claims.Get("xfoo")
		if err != nil {
			t.Fatal(errors.PrintLoc(err))
		}
		assert.Exactly(t, "bar", xFoo.(string))

	})
	authHandler := jwts.WithInitTokenAndStore()(finalHandler)

	wRec := httptest.NewRecorder()
	ctx := store.WithContextProvider(context.Background(), srv)

	authHandler.ServeHTTP(wRec, req.WithContext(ctx))
	assert.Equal(t, http.StatusTeapot, wRec.Code)
	assert.Equal(t, `I'm more of a coffee pot`, wRec.Body.String())
}

type testRealBL struct {
	theToken []byte
	exp      time.Duration
}

func (b *testRealBL) Set(t []byte, exp time.Duration) error {
	b.theToken = t
	b.exp = exp
	return nil
}
func (b *testRealBL) Has(t []byte) bool { return bytes.Equal(b.theToken, t) }

var _ jwtauth.Blacklister = (*testRealBL)(nil)

func TestMiddlewareWithInitTokenInBlackList(t *testing.T) {

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		store.WithStorageConfig(cr),
	)

	bl := &testRealBL{}
	jm, err := jwtauth.NewService(
		jwtauth.WithBlacklist(bl),
	)
	assert.NoError(t, err)

	theToken, err := jm.NewToken(scope.Default, 0, &jwtclaim.Standard{})
	bl.theToken = theToken.Raw
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken.Raw)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	jwtauth.SetHeaderAuthorization(req, theToken.Raw)

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := jwtauth.FromContext(r.Context())
		assert.True(t, errors.IsNotValid(err))
		w.WriteHeader(http.StatusUnauthorized)
	})
	authHandler := jm.WithInitTokenAndStore()(finalHandler)

	wRec := httptest.NewRecorder()
	ctx := store.WithContextProvider(context.Background(), srv)
	authHandler.ServeHTTP(wRec, req.WithContext(ctx))

	//assert.NotEqual(t, http.StatusTeapot, wRec.Code)
	assert.Equal(t, http.StatusUnauthorized, wRec.Code)
}

// todo add test for form with input field: access_token

func testAuth(t *testing.T, opts ...jwtauth.Option) (http.Handler, []byte) {
	jm, err := jwtauth.NewService(opts...)
	if err != nil {
		t.Fatal(errors.PrintLoc(err))
	}

	theToken, err := jm.NewToken(scope.Default, 0, jwtclaim.Map{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	assert.NoError(t, err)

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	authHandler := jm.WithInitTokenAndStore()(final)
	return authHandler, theToken.Raw
}

func newStoreServiceWithCtx(initO scope.Option) context.Context {
	return store.WithContextProvider(context.Background(), storemock.NewEurozzyService(initO))
}

func finalInitStoreHandler(t *testing.T, wantStoreCode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, haveReqStore, err := store.FromContextProvider(r.Context())
		if err != nil {
			t.Fatal(errors.PrintLoc(err))
		}
		assert.Exactly(t, wantStoreCode, haveReqStore.StoreCode())
	}
}

func TestWithInitTokenAndStore_Request(t *testing.T) {

	var newReq = func(i int, token []byte) *http.Request {
		req, err := http.NewRequest(httputil.MethodGet, fmt.Sprintf("https://corestore.io/store/list/%d", i), nil)
		if err != nil {
			t.Fatal(errors.PrintLoc(err))
		}
		jwtauth.SetHeaderAuthorization(req, token)
		return req
	}

	tests := []struct {
		ctx            context.Context
		tokenStoreCode string
		wantStoreCode  string
		wantErrBhf     errors.BehaviourFunc
	}{
		{store.WithContextProvider(context.Background(), nil), "de", "de", errors.IsNotFound},
		{store.WithContextProvider(context.Background(), storemock.NewEurozzyService(scope.Option{Store: scope.MockCode("de")})), "de", "de", errors.IsNotFound},
		{newStoreServiceWithCtx(scope.Option{Store: scope.MockCode("de")}), "de", "de", nil},
		{newStoreServiceWithCtx(scope.Option{Store: scope.MockCode("at")}), "ch", "at", errors.IsUnauthorized},
		{newStoreServiceWithCtx(scope.Option{Store: scope.MockCode("de")}), "at", "at", nil},
		{newStoreServiceWithCtx(scope.Option{Store: scope.MockCode("de")}), "a$t", "de", errors.IsNotValid},
		{newStoreServiceWithCtx(scope.Option{Store: scope.MockCode("at")}), "", "at", errors.IsNotValid},
		//
		{newStoreServiceWithCtx(scope.Option{Group: scope.MockID(1)}), "de", "de", nil},
		{newStoreServiceWithCtx(scope.Option{Group: scope.MockID(1)}), "ch", "at", errors.IsUnauthorized},
		{newStoreServiceWithCtx(scope.Option{Group: scope.MockID(1)}), " ch", "at", errors.IsNotValid},
		{newStoreServiceWithCtx(scope.Option{Group: scope.MockID(1)}), "uk", "at", errors.IsUnauthorized},
		//
		{newStoreServiceWithCtx(scope.Option{Website: scope.MockID(2)}), "uk", "au", errors.IsUnauthorized},
		{newStoreServiceWithCtx(scope.Option{Website: scope.MockID(2)}), "nz", "nz", nil},
		{newStoreServiceWithCtx(scope.Option{Website: scope.MockID(2)}), "n z", "au", errors.IsNotValid},
		{newStoreServiceWithCtx(scope.Option{Website: scope.MockID(2)}), "au", "au", nil},
		{newStoreServiceWithCtx(scope.Option{Website: scope.MockID(2)}), "", "au", errors.IsNotValid},
	}
	for i, test := range tests {
		jwts := jwtauth.MustNewService(jwtauth.WithKey(scope.Default, 0, csjwt.WithPasswordRandom()))

		token, err := jwts.NewToken(scope.Default, 0, jwtclaim.Map{
			jwtauth.StoreParamName: test.tokenStoreCode,
		})
		if err != nil {
			t.Fatal(errors.PrintLoc(err))
		}

		if test.wantErrBhf != nil {
			if err := jwts.Options(jwtauth.WithErrorHandler(scope.Default, 0,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, err := jwtauth.FromContext(r.Context())
					assert.True(t, test.wantErrBhf(err), "Index %d => %s", i, err)
				}),
			)); err != nil {
				t.Fatal(errors.PrintLoc(err))
			}
		}
		mw := jwts.WithInitTokenAndStore()(finalInitStoreHandler(t, test.wantStoreCode))
		rec := httptest.NewRecorder()
		mw.ServeHTTP(rec, newReq(i, token.Raw).WithContext(test.ctx))
	}
}
