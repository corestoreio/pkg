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

package ctxjwt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestServiceWithBackend_NoBackend(t *testing.T) {

	jwts := MustNewService()
	// a hack for testing to remove the default setting or make it invalid
	jwts.defaultScopeCache = scopedConfig{}

	cr := cfgmock.NewService()
	sc, err := jwts.ConfigByScopedGetter(cr.NewScoped(0, 0))
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	assert.Exactly(t, scopedConfig{}, sc)
}

func TestServiceWithBackend_DefaultConfig(t *testing.T) {

	jwts := MustNewService()

	cr := cfgmock.NewService()
	sc, err := jwts.ConfigByScopedGetter(cr.NewScoped(0, 0))
	assert.NoError(t, err)
	dsc, err := defaultScopedConfig()
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, csjwt.HS256, sc.SigningMethod.Alg())
	assert.Exactly(t, dsc.Key.Algorithm(), sc.Key.Algorithm())

	assert.Nil(t, dsc.ErrorHandler)
	assert.Nil(t, sc.ErrorHandler)
	assert.Nil(t, jwts.defaultScopeCache.ErrorHandler)
	assert.Exactly(t, DefaultExpire, dsc.Expire)
	assert.False(t, dsc.Key.IsEmpty())
	assert.False(t, sc.Key.IsEmpty())
}

func TestWithInitTokenAndStore_EqualPointers(t *testing.T) {

	// this Test is related to Benchmark_WithInitTokenAndStore
	// The returned pointers from store.FromContextReader must be the
	// same for each request with the same request pattern.

	var equalStorePointer *store.Store
	jwts := MustNewService()
	mw := jwts.WithInitTokenAndStore()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if _, err := FromContext(ctx); err != nil {
			t.Fatal(err)
		}

		_, haveReqStore, err := store.FromContextProvider(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if equalStorePointer == nil {
			equalStorePointer = haveReqStore
		}

		if have, want := haveReqStore.StoreCode(), "nz"; have != want {
			t.Errorf("Have: %q Want: %q", have, want)
		}
		cstesting.EqualPointers(t, equalStorePointer, haveReqStore)
	}))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(httputil.MethodGet, "https://corestore.io/store/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	sc := jwtclaim.NewStore()
	sc.Store = "nz"
	tok, err := jwts.NewToken(scope.Default, 0, sc)
	if err != nil {
		t.Fatal(err)
	}
	SetHeaderAuthorization(req, tok.Raw)
	// Bind request to a specific Website in this case down under
	ctx := store.WithContextProvider(context.Background(), storemock.NewEurozzyService(scope.Option{Website: scope.MockID(2)}))

	for i := 0; i < 10; i++ {
		mw.ServeHTTP(rec, req.WithContext(ctx))
	}
}
