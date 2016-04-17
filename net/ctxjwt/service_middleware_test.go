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

package ctxjwt_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/store/storenet"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestMiddlewareWithInitTokenNoStoreProvider(t *testing.T) {
	t.Parallel()

	authHandler, _ := testAuth(t, ctxjwt.WithErrorHandler(scope.Default, 0, ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
		tk, err := ctxjwt.FromContext(ctx)
		assert.False(t, tk.Valid)
		assert.EqualError(t, err, store.ErrContextProviderNotFound.Error())
		return nil
	})))

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	assert.NoError(t, authHandler.ServeHTTPContext(context.Background(), w, req))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestMiddlewareWithInitTokenNoToken(t *testing.T) {
	t.Parallel()

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		store.WithStorageConfig(cr),
	)
	ctx := store.WithContextProvider(context.Background(), srv)
	authHandler, _ := testAuth(t)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	assert.NoError(t, authHandler.ServeHTTPContext(ctx, w, req))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, w.Body.String(), http.StatusText(http.StatusUnauthorized)+"\n")
}

func TestMiddlewareWithInitTokenHTTPErrorHandler(t *testing.T) {
	t.Parallel()

	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		store.WithStorageConfig(cfgmock.NewService()), // empty config
	)
	ctx := store.WithContextProvider(context.Background(), srv)

	authHandler, _ := testAuth(t, ctxjwt.WithErrorHandler(scope.Default, 0, ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		tok, err := ctxjwt.FromContext(ctx)
		assert.False(t, tok.Valid)
		w.WriteHeader(http.StatusTeapot)
		_, err = w.Write([]byte(err.Error()))
		return err
	})))

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	assert.NoError(t, authHandler.ServeHTTPContext(ctx, w, req))
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Equal(t, csjwt.ErrTokenNotInRequest.Error(), w.Body.String())
}

func TestMiddlewareWithInitTokenSuccess(t *testing.T) {
	t.Parallel()

	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		store.WithStorageConfig(cfgmock.NewService()),
	)
	ctx := store.WithContextProvider(context.Background(), srv)

	jwts := ctxjwt.MustNewService()

	if err := jwts.Options(ctxjwt.WithErrorHandler(scope.Default, 0,
		ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			token, err := ctxjwt.FromContext(ctx)
			t.Logf("Token: %#v\n", token)
			t.Fatal("Unexpected Error:", cserr.NewMultiErr(err).VerboseErrors())
			return nil
		}),
	)); err != nil {
		t.Fatal(err)
	}

	theToken, err := jwts.NewToken(scope.Default, 0, jwtclaim.Map{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken.Raw)

	req, err := http.NewRequest("GET", "http://corestore.io/customer/account", nil)
	assert.NoError(t, err)
	ctxjwt.SetHeaderAuthorization(req, theToken.Raw)

	finalHandler := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprintf(w, "I'm more of a coffee pot")

		ctxToken, err := ctxjwt.FromContext(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, ctxToken)
		xFoo, err := ctxToken.Claims.Get("xfoo")
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, "bar", xFoo.(string))
		return nil
	})
	authHandler := jwts.WithInitTokenAndStore()(finalHandler)

	wRec := httptest.NewRecorder()
	assert.NoError(t, authHandler.ServeHTTPContext(ctx, wRec, req))
	assert.Equal(t, http.StatusTeapot, wRec.Code)
	assert.Equal(t, `I'm more of a coffee pot`, wRec.Body.String())
}

func TestMiddlewareWithInitTokenInBlackList(t *testing.T) {
	t.Parallel()

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		store.WithStorageConfig(cr),
	)
	ctx := store.WithContextProvider(context.Background(), srv)

	bl := &testRealBL{}
	jm, err := ctxjwt.NewService(
		ctxjwt.WithBlacklist(bl),
	)
	assert.NoError(t, err)

	theToken, err := jm.NewToken(scope.Default, 0, &jwtclaim.Standard{})
	bl.theToken = theToken.Raw
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken.Raw)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	ctxjwt.SetHeaderAuthorization(req, theToken.Raw)

	finalHandler := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		return nil
	})
	authHandler := jm.WithInitTokenAndStore()(finalHandler)

	wRec := httptest.NewRecorder()
	assert.NoError(t, authHandler.ServeHTTPContext(ctx, wRec, req))

	assert.NotEqual(t, http.StatusTeapot, wRec.Code)
	assert.Equal(t, http.StatusUnauthorized, wRec.Code)
}

// todo add test for form with input field: access_token

func testAuth(t *testing.T, opts ...ctxjwt.Option) (ctxhttp.Handler, []byte) {
	jm, err := ctxjwt.NewService(opts...)
	if err != nil {
		t.Fatal(err)
	}

	theToken, err := jm.NewToken(scope.Default, 0, jwtclaim.Map{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	assert.NoError(t, err)

	final := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})
	authHandler := jm.WithInitTokenAndStore()(final)
	return authHandler, theToken.Raw
}

func newStoreServiceWithCtx(initO scope.Option) context.Context {
	ctx := store.WithContextProvider(context.Background(), storemock.NewEurozzyService(initO))

	return ctx
}

func finalInitStoreHandler(t *testing.T, wantStoreCode string) ctxhttp.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		_, haveReqStore, err := store.FromContextProvider(ctx)
		if err != nil {
			t.Fatal(err)
			return err
		}
		assert.Exactly(t, wantStoreCode, haveReqStore.StoreCode())
		return nil
	}
}

func TestWithInitTokenAndStore_Request(t *testing.T) {
	t.Parallel()
	var newReq = func(i int, token []byte) *http.Request {
		req, err := http.NewRequest(httputil.MethodGet, fmt.Sprintf("https://corestore.io/store/list/%d", i), nil)
		if err != nil {
			t.Fatal(err)
		}
		ctxjwt.SetHeaderAuthorization(req, token)
		return req
	}

	tests := []struct {
		ctx            context.Context
		tokenStoreCode string
		wantStoreCode  string
		wantErr        error
	}{
		{store.WithContextProvider(context.Background(), nil), "de", "de", store.ErrContextProviderNotFound},
		{store.WithContextProvider(context.Background(), storemock.NewEurozzyService(scope.Option{Store: scope.MockCode("de")})), "de", "de", csjwt.ErrTokenNotInRequest},
		{newStoreServiceWithCtx(scope.Option{Store: scope.MockCode("de")}), "de", "de", nil},
		{newStoreServiceWithCtx(scope.Option{Store: scope.MockCode("at")}), "ch", "at", store.ErrStoreNotActive},
		{newStoreServiceWithCtx(scope.Option{Store: scope.MockCode("de")}), "at", "at", nil},
		{newStoreServiceWithCtx(scope.Option{Store: scope.MockCode("de")}), "a$t", "de", store.ErrStoreCodeInvalid},
		{newStoreServiceWithCtx(scope.Option{Store: scope.MockCode("at")}), "", "at", store.ErrStoreCodeInvalid},
		//
		{newStoreServiceWithCtx(scope.Option{Group: scope.MockID(1)}), "de", "de", nil},
		{newStoreServiceWithCtx(scope.Option{Group: scope.MockID(1)}), "ch", "at", store.ErrStoreNotActive},
		{newStoreServiceWithCtx(scope.Option{Group: scope.MockID(1)}), " ch", "at", store.ErrStoreCodeInvalid},
		{newStoreServiceWithCtx(scope.Option{Group: scope.MockID(1)}), "uk", "at", store.ErrStoreChangeNotAllowed},
		//
		{newStoreServiceWithCtx(scope.Option{Website: scope.MockID(2)}), "uk", "au", store.ErrStoreChangeNotAllowed},
		{newStoreServiceWithCtx(scope.Option{Website: scope.MockID(2)}), "nz", "nz", nil},
		{newStoreServiceWithCtx(scope.Option{Website: scope.MockID(2)}), "n z", "au", store.ErrStoreCodeInvalid},
		{newStoreServiceWithCtx(scope.Option{Website: scope.MockID(2)}), "au", "au", nil},
		{newStoreServiceWithCtx(scope.Option{Website: scope.MockID(2)}), "", "au", store.ErrStoreCodeInvalid},
	}
	for i, test := range tests {
		jwts := ctxjwt.MustNewService(ctxjwt.WithKey(scope.Default, 0, csjwt.WithPasswordRandom()))

		token, err := jwts.NewToken(scope.Default, 0, jwtclaim.Map{
			storenet.ParamName: test.tokenStoreCode,
		})
		if err != nil {
			t.Fatal(err)
		}

		if test.wantErr != nil {
			if err := jwts.Options(ctxjwt.WithErrorHandler(scope.Default, 0,
				ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
					_, err := ctxjwt.FromContext(ctx)
					assert.EqualError(t, err, test.wantErr.Error(), "Index %d", i)
					return nil
				}),
			)); err != nil {
				t.Fatal(err)
			}
		}
		mw := jwts.WithInitTokenAndStore()(finalInitStoreHandler(t, test.wantStoreCode))
		rec := httptest.NewRecorder()
		if err := mw.ServeHTTPContext(test.ctx, rec, newReq(i, token.Raw)); err != nil {
			t.Fatal(err)
		}
	}
}
