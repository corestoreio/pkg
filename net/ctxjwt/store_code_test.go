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
	"testing"

	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
)

func TestStoreCodeFromClaimFullToken(t *testing.T) {
	t.Parallel()
	s := store.MustNewStore(
		&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
	)

	token := csjwt.NewToken(jwtclaim.Map{
		ctxjwt.StoreParamName: s.StoreCode(),
	})

	so, err := ctxjwt.ScopeOptionFromClaim(token.Claims)
	assert.NoError(t, err)
	assert.EqualValues(t, "de", so.StoreCode())

	so, err = ctxjwt.ScopeOptionFromClaim(nil)
	assert.EqualError(t, store.ErrStoreNotFound, err.Error())
	assert.Nil(t, so.Website)
	assert.Nil(t, so.Group)
	assert.Nil(t, so.Store)

}

func TestStoreCodeFromClaimInvalid(t *testing.T) {
	t.Parallel()

	token2 := csjwt.NewToken(jwtclaim.Map{
		ctxjwt.StoreParamName: "Invalid Cod€",
	})

	so, err := ctxjwt.ScopeOptionFromClaim(token2.Claims)
	assert.EqualError(t, store.ErrStoreCodeInvalid, err.Error())
	assert.Nil(t, so.Website)
	assert.Nil(t, so.Group)
	assert.Nil(t, so.Store)
}

func TestStoreCodeFromClaimNoToken(t *testing.T) {
	t.Parallel()
	tests := []struct {
		token     csjwt.Claimer
		wantErr   error
		wantScope scope.Scope
		wantCode  string
		wantID    int64
	}{
		{
			jwtclaim.Map{},
			store.ErrStoreNotFound,
			scope.Default,
			"",
			0,
		},
		{
			jwtclaim.Map{ctxjwt.StoreParamName: "dede"},
			nil,
			scope.Store,
			"dede",
			scope.UnavailableStoreID,
		},
		{
			jwtclaim.Map{ctxjwt.StoreParamName: "de'de"},
			store.ErrStoreCodeInvalid,
			scope.Default,
			"",
			scope.UnavailableStoreID,
		},
		{
			jwtclaim.Map{ctxjwt.StoreParamName: 1},
			store.ErrStoreNotFound,
			scope.Default,
			"",
			scope.UnavailableStoreID,
		},
	}
	for i, test := range tests {
		so, err := ctxjwt.ScopeOptionFromClaim(test.token)
		testStoreCodeFrom(t, i, err, test.wantErr, so, test.wantScope, test.wantCode, test.wantID)
	}
}

func testStoreCodeFrom(t *testing.T, i int, haveErr, wantErr error, haveScope scope.Option, wantScope scope.Scope, wantCode string, wantID int64) {
	if wantErr != nil {
		assert.EqualError(t, haveErr, wantErr.Error(), "Index: %d", i)

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
