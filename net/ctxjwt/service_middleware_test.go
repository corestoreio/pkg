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
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestMiddlewareWithParseAndValidateNoStoreProvider(t *testing.T) {
	t.Parallel()

	authHandler, _ := testAuth(t)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	assert.EqualError(t, authHandler.ServeHTTPContext(context.Background(), w, req), store.ErrContextProviderNotFound.Error())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestMiddlewareWithParseAndValidateNoToken(t *testing.T) {
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

func TestMiddlewareWithParseAndValidateHTTPErrorHandler(t *testing.T) {
	t.Parallel()

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.WebsiteID, "euro"),
		store.WithStorageConfig(cr),
	)
	ctx := store.WithContextProvider(context.Background(), srv)

	authHandler, _ := testAuth(t, ctxjwt.WithErrorHandler(scope.DefaultID, 0, func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		tok, err := ctxjwt.FromContext(ctx)
		assert.Nil(t, tok)
		w.WriteHeader(http.StatusTeapot)
		_, err = w.Write([]byte(err.Error()))
		return err
	}))

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	assert.NoError(t, authHandler.ServeHTTPContext(ctx, w, req))
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Equal(t, "no token present in request", w.Body.String())
}

func TestMiddlewareWithParseAndValidateSuccess(t *testing.T) {
	t.Parallel()

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.WebsiteID, "euro"),
		store.WithStorageConfig(cr),
	)
	ctx := store.WithContextProvider(context.Background(), srv)

	jm, err := ctxjwt.NewService()
	assert.NoError(t, err)

	theToken, _, err := jm.GenerateToken(scope.DefaultID, 0, map[string]interface{}{
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
	authHandler := jm.WithParseAndValidate()(finalHandler)

	wRec := httptest.NewRecorder()
	assert.NoError(t, authHandler.ServeHTTPContext(ctx, wRec, req))
	assert.Equal(t, http.StatusTeapot, wRec.Code)
	assert.Equal(t, `I'm more of a coffee pot`, wRec.Body.String())
}

func TestMiddlewareWithParseAndValidateInBlackList(t *testing.T) {
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
	bl.theToken = theToken
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	ctxjwt.SetHeaderAuthorization(req, theToken)

	finalHandler := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		return nil
	})
	authHandler := jm.WithParseAndValidate()(finalHandler)

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
	authHandler := jm.WithParseAndValidate()(final)
	return authHandler, theToken
}
