// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package request_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/pkg/net/request"
	"github.com/corestoreio/pkg/util/assert"
)

func getReq(accept string) *http.Request {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept", accept)
	return r
}

func TestAcceptsJSON(t *testing.T) {

	t.Run("empty", func(t *testing.T) {
		assert.True(t, request.AcceptsJSON(getReq("")))
	})
	t.Run("wild card / wild card", func(t *testing.T) {
		assert.True(t, request.AcceptsJSON(getReq("*/*")))
	})
	t.Run("application / wild card", func(t *testing.T) {
		assert.True(t, request.AcceptsJSON(getReq("application/*")))
	})
	t.Run("application / json", func(t *testing.T) {
		assert.True(t, request.AcceptsJSON(getReq("application/json")))
	})
	t.Run("application / xml", func(t *testing.T) {
		assert.False(t, request.AcceptsJSON(getReq("application/xml")))
	})
}

func TestAcceptsContentType(t *testing.T) {

	t.Run("empty", func(t *testing.T) {
		assert.True(t, request.AcceptsContentType(getReq(""), "text/plain"))
	})
	t.Run("wild card / wild card", func(t *testing.T) {
		assert.True(t, request.AcceptsContentType(getReq("*/*"), "text/plain"))
	})
	t.Run("application / wild card", func(t *testing.T) {
		assert.True(t, request.AcceptsContentType(getReq("text/*"), "text/plain"))
	})
	t.Run("application / json", func(t *testing.T) {
		assert.True(t, request.AcceptsContentType(getReq("text/plain"), "text/plain"))
	})
	t.Run("application / xml", func(t *testing.T) {
		assert.False(t, request.AcceptsContentType(getReq("text/html"), "text/plain"))
	})
}
