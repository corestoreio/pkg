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
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestMiddlewareWithInitTokenNoStoreProvider(t *testing.T) {
	t.Parallel()

	authHandler, _ := testAuth(t, ctxjwt.WithErrorHandler(scope.DefaultID, 0, ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
		tk, err := ctxjwt.FromContext(ctx)
		assert.Nil(t, tk)
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
		scope.MustSetByCode(scope.WebsiteID, "euro"),
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

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.WebsiteID, "euro"),
		store.WithStorageConfig(cr),
	)
	ctx := store.WithContextProvider(context.Background(), srv)

	authHandler, _ := testAuth(t, ctxjwt.WithErrorHandler(scope.DefaultID, 0, ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		tok, err := ctxjwt.FromContext(ctx)
		assert.Nil(t, tok)
		w.WriteHeader(http.StatusTeapot)
		_, err = w.Write([]byte(err.Error()))
		return err
	})))

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	assert.NoError(t, authHandler.ServeHTTPContext(ctx, w, req))
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Equal(t, "no token present in request", w.Body.String())
}

func TestMiddlewareWithInitTokenSuccess(t *testing.T) {
	t.Parallel()

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.WebsiteID, "euro"),
		store.WithStorageConfig(cr),
	)
	ctx := store.WithContextProvider(context.Background(), srv)

	jwts := ctxjwt.MustNewService()

	theToken, _, err := jwts.GenerateToken(scope.DefaultID, 0, map[string]interface{}{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	ctxjwt.SetHeaderAuthorization(req, theToken)

	finalHandler := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprintf(w, "I'm more of a coffee pot")

		ctxToken, err := ctxjwt.FromContext(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, ctxToken)
		assert.Exactly(t, "bar", ctxToken.Claims["xfoo"].(string))

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
		scope.MustSetByCode(scope.WebsiteID, "euro"),
		store.WithStorageConfig(cr),
	)
	ctx := store.WithContextProvider(context.Background(), srv)

	bl := &testRealBL{}
	jm, err := ctxjwt.NewService(
		ctxjwt.WithBlacklist(bl),
	)
	assert.NoError(t, err)

	theToken, _, err := jm.GenerateToken(scope.DefaultID, 0, nil)
	bl.theToken = []byte(theToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	ctxjwt.SetHeaderAuthorization(req, theToken)

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

func testAuth(t *testing.T, opts ...ctxjwt.Option) (ctxhttp.Handler, string) {
	jm, err := ctxjwt.NewService(opts...)
	if err != nil {
		t.Fatal(err)
	}

	theToken, _, err := jm.GenerateToken(scope.DefaultID, 0, map[string]interface{}{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	assert.NoError(t, err)

	final := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})
	authHandler := jm.WithInitTokenAndStore()(final)
	return authHandler, theToken
}

func newStoreServiceWithTokenCtx(initO scope.Option, tokenStoreCode string) context.Context {
	ctx := store.WithContextProvider(context.Background(), storemock.NewEurozzyService(initO))
	tok := csjwt.New(csjwt.SigningMethodHS256)
	tok.Claims[storenet.ParamName] = tokenStoreCode
	ctx = ctxjwt.WithContext(ctx, tok)
	return ctx
}

func finalInitStoreHandler(t *testing.T, wantStoreCode string) ctxhttp.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		_, haveReqStore, err := store.FromContextProvider(ctx)
		if err != nil {
			return err
		}
		assert.Exactly(t, wantStoreCode, haveReqStore.StoreCode())
		return nil
	}
}

func TestWithInitTokenAndStore(t *testing.T) {

	var newReq = func(i int) *http.Request {
		req, err := http.NewRequest(httputil.MethodGet, fmt.Sprintf("https://corestore.io/store/list/%d", i), nil)
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
		{store.WithContextProvider(context.Background(), nil), "de", store.ErrContextProviderNotFound},
		{store.WithContextProvider(context.Background(), storemock.NewEurozzyService(scope.Option{Store: scope.MockCode("de")})), "de", ctxjwt.ErrContextJWTNotFound},
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
		jwts := ctxjwt.MustNewService()
		mw := jwts.WithInitTokenAndStore()(finalInitStoreHandler(t, test.wantStoreCode))
		rec := httptest.NewRecorder()
		surfErr := mw.ServeHTTPContext(test.ctx, rec, newReq(i))
		if test.wantErr != nil {
			assert.EqualError(t, surfErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, surfErr, "Index %d", i)
	}
}

func TestWithInitTokenAndStore_EqualPointers(t *testing.T) {
	// this Test is related to Benchmark_WithInitTokenAndStore
	// The returned pointers from store.FromContextReader must be the
	// same for each request with the same request pattern.

	ctx := newStoreServiceWithTokenCtx(scope.Option{Website: scope.MockID(2)}, "nz")
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(httputil.MethodGet, "https://corestore.io/store/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	var equalStorePointer *store.Store
	jwts := ctxjwt.MustNewService()
	mw := jwts.WithInitTokenAndStore()(ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		_, haveReqStore, err := store.FromContextProvider(ctx)
		if err != nil {
			return err
		}

		if equalStorePointer == nil {
			equalStorePointer = haveReqStore
		}

		if "nz" != haveReqStore.StoreCode() {
			t.Errorf("Have: %s\nWant: nz", haveReqStore.StoreCode())
		}
		cstesting.EqualPointers(t, equalStorePointer, haveReqStore)

		return nil
	}))

	for i := 0; i < 10; i++ {
		if err := mw.ServeHTTPContext(ctx, rec, req); err != nil {
			t.Error(err)
		}
	}
}
