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

package storenet_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storenet"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestStoreCodeFromCookie(t *testing.T) {

	var getRootRequest = func(c *http.Cookie) *http.Request {
		rootRequest, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("Root request error: %s", err)
		}
		if c != nil {
			rootRequest.AddCookie(c)
		}
		return rootRequest
	}

	tests := []struct {
		req        *http.Request
		wantErrBhf errors.BehaviourFunc
		wantScope  scope.Scope
		wantCode   string
		wantID     int64
	}{
		{
			getRootRequest(&http.Cookie{Name: storenet.ParamName, Value: "dede"}),
			nil,
			scope.Store,
			"dede",
			scope.UnavailableStoreID,
		},
		{
			getRootRequest(&http.Cookie{Name: storenet.ParamName, Value: "ded'e"}),
			errors.IsNotValid,
			scope.Default,
			"",
			scope.UnavailableStoreID,
		},
		{
			getRootRequest(&http.Cookie{Name: "invalid", Value: "dede"}),
			errors.IsNotFound,
			scope.Default,
			"",
			scope.UnavailableStoreID,
		},
	}
	for i, test := range tests {
		so, err := storenet.CodeFromCookie(test.req)
		testStoreCodeFrom(t, i, err, test.wantErrBhf, so, test.wantScope, test.wantCode, test.wantID)
	}
}

func TestStoreCodeFromRequestGET(t *testing.T) {

	var getRootRequest = func(kv ...string) *http.Request {

		reqURL := "http://corestore.io/"
		var uv url.Values
		if len(kv)%2 == 0 {
			uv = make(url.Values)
			for i := 0; i < len(kv); i = i + 2 {
				uv.Set(kv[i], kv[i+1])
			}
			reqURL = reqURL + "?" + uv.Encode()
		}
		rootRequest, err := http.NewRequest("GET", reqURL, nil)
		if err != nil {
			t.Fatalf("Root request error: %s", err)
		}

		return rootRequest
	}

	tests := []struct {
		req        *http.Request
		wantErrBhf errors.BehaviourFunc
		wantScope  scope.Scope
		wantCode   string
		wantID     int64
	}{
		{
			getRootRequest(storenet.HTTPRequestParamStore, "dede"),
			nil,
			scope.Store,
			"dede",
			scope.UnavailableStoreID,
		},
		{
			getRootRequest(storenet.HTTPRequestParamStore, "ded¢e"),
			errors.IsNotValid,
			scope.Default,
			"",
			scope.UnavailableStoreID,
		},
		{
			getRootRequest("invalid", "dede"),
			errors.IsNotValid,
			scope.Default,
			"",
			scope.UnavailableStoreID,
		},
	}
	for i, test := range tests {
		so, err := storenet.CodeFromRequestGET(test.req)
		testStoreCodeFrom(t, i, err, test.wantErrBhf, so, test.wantScope, test.wantCode, test.wantID)
	}
}

func testStoreCodeFrom(t *testing.T, i int, haveErr error, wantErrBhf errors.BehaviourFunc, haveScope scope.Option, wantScope scope.Scope, wantCode string, wantID int64) {
	if wantErrBhf != nil {
		assert.True(t, wantErrBhf(haveErr), "Index: %d => %s", i, haveErr)
	}
	switch sos := haveScope.Scope(); sos {
	case scope.Store:
		assert.Exactly(t, wantID, haveScope.Store.StoreID(), "Index: %d", i)
	case scope.Group:
		assert.Exactly(t, wantID, haveScope.Group.GroupID(), "Index: %d", i)
	case scope.Website:
		assert.Exactly(t, wantID, haveScope.Website.WebsiteID(), "Index: %d", i)
	case scope.Default:
		assert.Nil(t, haveScope.Store, "Index: %d", i)
		assert.Nil(t, haveScope.Group, "Index: %d", i)
		assert.Nil(t, haveScope.Website, "Index: %d", i)
	default:
		t.Fatalf("Unknown scope: %d", sos)
	}
	assert.Exactly(t, wantScope, haveScope.Scope(), "Index: %d", i)
	assert.Exactly(t, wantCode, haveScope.StoreCode(), "Index: %d", i)
}
