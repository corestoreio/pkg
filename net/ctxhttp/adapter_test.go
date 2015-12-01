// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ctxhttp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxmw"
	"github.com/corestoreio/csfw/net/httputils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type h1 struct{}

func (h1) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte(`h1 called`))
	return err
}

func TestAdapters(t *testing.T) {

	hndlr := ctxhttp.Chain(
		h1{},
		ctxmw.WithXHTTPMethodOverride(),
		ctxmw.WithHeader("X-Men", "Y-Women"),
	)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "http://example.com/foo", nil)
	req.Header.Set(httputils.MethodOverrideHeader, httputils.MethodPut)
	assert.NoError(t, err)

	a := ctxhttp.NewAdapter(context.Background(), hndlr)
	a.ServeHTTP(w, req)

	assert.Equal(t, httputils.MethodPut, req.Method)
	assert.Equal(t, "h1 called", w.Body.String())
	assert.Equal(t, "Y-Women", w.Header().Get("X-Men"))
}

func TestDefaultAdapterErrFunc(t *testing.T) {
	anErr := errors.New("This error should be returned")

	w := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "http://example.com/foo?_method=KARATE", nil)
	assert.NoError(t, err)
	ctxhttp.DefaultAdapterErrFunc(w, req, anErr)
	assert.Exactly(t, http.StatusBadRequest, w.Code)
	assert.Exactly(t, "Bad Request\nApp Error: This error should be returned\n", w.Body.String())
}

func TestHttpMethodOverride(t *testing.T) {
	hndlr := ctxhttp.Chain(
		h1{},
		ctxmw.WithXHTTPMethodOverride())
	w := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "http://example.com/foo?_method="+httputils.MethodPatch, nil)
	assert.NoError(t, err)
	assert.NoError(t, hndlr.ServeHTTPContext(context.Background(), w, req))
	assert.Equal(t, httputils.MethodPatch, req.Method)
	assert.Equal(t, "h1 called", w.Body.String())

	w = httptest.NewRecorder()
	req, err = http.NewRequest(httputils.MethodGet, "http://example.com/foo?_method=KARATE", nil)
	assert.NoError(t, err)
	assert.NoError(t, hndlr.ServeHTTPContext(context.Background(), w, req))
	assert.Equal(t, httputils.MethodGet, req.Method)

	w = httptest.NewRecorder()
	req, err = http.NewRequest(httputils.MethodGet, "http://example.com/foobar", nil)
	assert.NoError(t, err)
	assert.NoError(t, hndlr.ServeHTTPContext(context.Background(), w, req))
	assert.Equal(t, httputils.MethodGet, req.Method)

}
