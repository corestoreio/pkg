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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestServiceWithBackend_NoBackend(t *testing.T) {
	t.Parallel()

	jwts := MustNewService()
	// a hack for testing to remove the default setting or make it invalid
	jwts.defaultScopeCache = scopedConfig{}

	cr := cfgmock.NewService()
	sc, err := jwts.ConfigByScopedGetter(cr.NewScoped(0, 0))
	assert.EqualError(t, err, "[ctxjwt] Cannot find JWT configuration for Scope(Default) ID(0)")
	assert.Exactly(t, scopedConfig{}, sc)
}

func TestServiceWithBackend_DefaultConfig(t *testing.T) {
	t.Parallel()

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

	assert.NotNil(t, dsc.ErrorHandler)
	assert.NotNil(t, sc.ErrorHandler)
	assert.True(t, jwts.defaultScopeCache.ErrorHandler != nil)
	assert.Exactly(t, DefaultExpire, dsc.Expire)
	assert.False(t, dsc.Key.IsEmpty())
	assert.False(t, sc.Key.IsEmpty())
}

func newStoreServiceWithTokenCtx(initO scope.Option, tokenStoreCode string) context.Context {
	ctx := store.WithContextProvider(context.Background(), storemock.NewEurozzyService(initO))
	tok := csjwt.NewToken(jwtclaim.Map{
		StoreParamName: tokenStoreCode,
	})
	ctx = withContext(ctx, tok)
	return ctx
}

func TestWithInitTokenAndStore_EqualPointers(t *testing.T) {
	t.Parallel()
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
	jwts := MustNewService()
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
