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

package cshttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/net/cshttp"
	"github.com/stretchr/testify/assert"
)

type h1 struct{}

func (h1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`h1 called`))
}

func TestAdapters(t *testing.T) {
	hndlr := cshttp.Adapt(
		h1{},
		cshttp.SupportXHTTPMethodOverride(),
		cshttp.WithHeader("X-Men", "Y-Women"))
	w := httptest.NewRecorder()
	req, err := http.NewRequest(cshttp.HTTPMethodGet, "http://example.com/foo", nil)
	req.Header.Set(cshttp.HTTPMethodOverrideHeader, cshttp.HTTPMethodPut)
	assert.NoError(t, err)
	hndlr.ServeHTTP(w, req)
	assert.Equal(t, cshttp.HTTPMethodPut, req.Method)
	assert.Equal(t, "h1 called", w.Body.String())
	assert.Equal(t, "Y-Women", w.Header().Get("X-Men"))
}

func TestHttpMethodOverride(t *testing.T) {
	hndlr := cshttp.Adapt(
		h1{},
		cshttp.SupportXHTTPMethodOverride())
	w := httptest.NewRecorder()
	req, err := http.NewRequest(cshttp.HTTPMethodGet, "http://example.com/foo?_method="+cshttp.HTTPMethodPatch, nil)
	assert.NoError(t, err)
	hndlr.ServeHTTP(w, req)
	assert.Equal(t, cshttp.HTTPMethodPatch, req.Method)
	assert.Equal(t, "h1 called", w.Body.String())

	w = httptest.NewRecorder()
	req, err = http.NewRequest(cshttp.HTTPMethodGet, "http://example.com/foo?_method=KARATE", nil)
	assert.NoError(t, err)
	hndlr.ServeHTTP(w, req)
	assert.Equal(t, cshttp.HTTPMethodGet, req.Method)

	w = httptest.NewRecorder()
	req, err = http.NewRequest(cshttp.HTTPMethodGet, "http://example.com/foobar", nil)
	assert.NoError(t, err)
	hndlr.ServeHTTP(w, req)
	assert.Equal(t, cshttp.HTTPMethodGet, req.Method)

}
