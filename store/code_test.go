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

package store_test

import (
	"testing"

	"net/http"
	"net/url"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/store"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestStoreCodeFromClaimFullToken(t *testing.T) {
	s := store.MustNewStore(
		&store.TableStore{StoreID: 1, Code: csdb.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: csdb.NewNullString("admin"), Name: csdb.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: csdb.NewNullBool(false)},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
	)
	token := jwt.New(jwt.SigningMethodHS256)
	s.AddClaim(token.Claims)

	so, err := store.CodeFromClaim(token.Claims)
	assert.NoError(t, err)
	assert.EqualValues(t, "de", so.StoreCode())

	so, err = store.CodeFromClaim(nil)
	assert.EqualError(t, store.ErrStoreNotFound, err.Error())
	assert.Nil(t, so.Website)
	assert.Nil(t, so.Group)
	assert.Nil(t, so.Store)

	token2 := jwt.New(jwt.SigningMethodHS256)
	token2.Claims[store.ParamName] = "Invalid Cod€"
	so, err = store.CodeFromClaim(token2.Claims)
	assert.EqualError(t, store.ErrStoreCodeInvalid, err.Error())
	assert.Nil(t, so.Website)
	assert.Nil(t, so.Group)
	assert.Nil(t, so.Store)
}

func TestStoreCodeFromClaimNoToken(t *testing.T) {
	tests := []struct {
		token     map[string]interface{}
		wantErr   error
		wantScope scope.Scope
		wantCode  string
		wantID    int64
	}{
		{
			map[string]interface{}{},
			store.ErrStoreNotFound,
			scope.DefaultID,
			"",
			0,
		},
		{
			map[string]interface{}{store.ParamName: "dede"},
			nil,
			scope.StoreID,
			"dede",
			scope.UnavailableStoreID,
		},
		{
			map[string]interface{}{store.ParamName: "de'de"},
			store.ErrStoreCodeInvalid,
			scope.DefaultID,
			"",
			scope.UnavailableStoreID,
		},
		{
			map[string]interface{}{store.ParamName: 1},
			store.ErrStoreNotFound,
			scope.DefaultID,
			"",
			scope.UnavailableStoreID,
		},
	}
	for i, test := range tests {
		so, err := store.CodeFromClaim(test.token)
		testStoreCodeFrom(t, i, err, test.wantErr, so, test.wantScope, test.wantCode, test.wantID)
	}
}

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
		req       *http.Request
		wantErr   error
		wantScope scope.Scope
		wantCode  string
		wantID    int64
	}{
		{
			nil,
			store.ErrStoreNotFound,
			scope.DefaultID,
			"",
			0,
		},
		{
			getRootRequest(&http.Cookie{Name: store.ParamName, Value: "dede"}),
			nil,
			scope.StoreID,
			"dede",
			scope.UnavailableStoreID,
		},
		{
			getRootRequest(&http.Cookie{Name: store.ParamName, Value: "ded'e"}),
			store.ErrStoreCodeInvalid,
			scope.DefaultID,
			"",
			scope.UnavailableStoreID,
		},
		{
			getRootRequest(&http.Cookie{Name: "invalid", Value: "dede"}),
			http.ErrNoCookie,
			scope.DefaultID,
			"",
			scope.UnavailableStoreID,
		},
	}
	for i, test := range tests {
		so, err := store.CodeFromCookie(test.req)
		testStoreCodeFrom(t, i, err, test.wantErr, so, test.wantScope, test.wantCode, test.wantID)
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
		req       *http.Request
		wantErr   error
		wantScope scope.Scope
		wantCode  string
		wantID    int64
	}{
		{
			nil,
			store.ErrStoreNotFound,
			scope.DefaultID,
			"",
			0,
		},
		{
			getRootRequest(store.HTTPRequestParamStore, "dede"),
			nil,
			scope.StoreID,
			"dede",
			scope.UnavailableStoreID,
		},
		{
			getRootRequest(store.HTTPRequestParamStore, "ded¢e"),
			store.ErrStoreCodeInvalid,
			scope.DefaultID,
			"",
			scope.UnavailableStoreID,
		},
		{
			getRootRequest("invalid", "dede"),
			store.ErrStoreCodeInvalid,
			scope.DefaultID,
			"",
			scope.UnavailableStoreID,
		},
	}
	for i, test := range tests {
		so, err := store.CodeFromRequestGET(test.req)
		testStoreCodeFrom(t, i, err, test.wantErr, so, test.wantScope, test.wantCode, test.wantID)
	}
}

func testStoreCodeFrom(t *testing.T, i int, haveErr, wantErr error, haveScope scope.Option, wantScope scope.Scope, wantCode string, wantID int64) {
	if wantErr != nil {
		assert.EqualError(t, haveErr, wantErr.Error(), "Index: %d", i)

	}
	switch sos := haveScope.Scope(); sos {
	case scope.StoreID:
		assert.Exactly(t, wantID, haveScope.Store.StoreID(), "Index: %d", i)
	case scope.GroupID:
		assert.Exactly(t, wantID, haveScope.Group.GroupID(), "Index: %d", i)
	case scope.WebsiteID:
		assert.Exactly(t, wantID, haveScope.Website.WebsiteID(), "Index: %d", i)
	case scope.DefaultID:
		assert.Nil(t, haveScope.Store, "Index: %d", i)
		assert.Nil(t, haveScope.Group, "Index: %d", i)
		assert.Nil(t, haveScope.Website, "Index: %d", i)
	default:
		t.Fatalf("Unknown scope: %d", sos)
	}
	assert.Exactly(t, wantScope, haveScope.Scope(), "Index: %d", i)
	assert.Exactly(t, wantCode, haveScope.StoreCode(), "Index: %d", i)
}

func TestValidateStoreCode(t *testing.T) {
	tests := []struct {
		have    string
		wantErr error
	}{
		{"@de", store.ErrStoreCodeInvalid},
		{" de", store.ErrStoreCodeInvalid},
		{"de", nil},
		{"DE", nil},
		{"deCH09_", nil},
		{"_de", store.ErrStoreCodeInvalid},
		{"", store.ErrStoreCodeInvalid},
		{"\U0001f41c", store.ErrStoreCodeInvalid},
		{"au_en", nil},
		{"au-fr", store.ErrStoreCodeInvalid},
		{"Hello GoLang", store.ErrStoreCodeInvalid},
		{"Hello€GoLang", store.ErrStoreCodeInvalid},
		{"HelloGoLdhashdfkjahdjfhaskjdfhuiwehfiawehfuahweldsnjkasfkjkwejqwehqang", store.ErrStoreCodeInvalid},
	}
	for _, test := range tests {
		haveErr := store.CodeIsValid(test.have)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "err codes switched: %#v", test)
		} else {
			assert.NoError(t, haveErr, "%#v", test)
		}
	}
}
