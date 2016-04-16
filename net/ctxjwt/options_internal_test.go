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

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cstesting"
	"golang.org/x/net/context"

	"fmt"
	"testing"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/stretchr/testify/assert"
)

func TestInternalOptionWithErrorHandler(t *testing.T) {
	t.Parallel()
	jwts := MustNewService()

	defaultErrH := jwts.defaultScopeCache.errorHandler

	wsErrH := ctxhttp.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
		http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
		return nil
	})

	if err := jwts.Options(WithErrorHandler(scope.Website, 22, wsErrH)); err != nil {
		t.Fatal(err)
	}

	cstesting.EqualPointers(t, defaultErrH, jwts.defaultScopeCache.errorHandler)
	cstesting.EqualPointers(t, wsErrH, jwts.scopeCache[scope.NewHash(scope.Website, 22)].errorHandler)
	cstesting.UnEqualPointers(t, defaultErrH, jwts.scopeCache[scope.NewHash(scope.Website, 22)].errorHandler)

	if err := jwts.Options(WithErrorHandler(scope.Default, 0, wsErrH)); err != nil {
		t.Fatal(err)
	}
	cstesting.UnEqualPointers(t, defaultErrH, jwts.defaultScopeCache.errorHandler)
	cstesting.EqualPointers(t, wsErrH, jwts.defaultScopeCache.errorHandler)
}

func TestInternalOptionNoLeakage(t *testing.T) {
	t.Parallel()
	sc := scopedConfig{
		Key: csjwt.WithPasswordRandom(),
	}
	assert.Exactly(t, `{0 csjwt.Key{/*redacted*/} 0 <nil> <nil> false <nil> <nil> <nil>}`, fmt.Sprintf("%v", sc))
	assert.Exactly(t, `ctxjwt.scopedConfig{scopeHash:0x0, Key:csjwt.Key{/*redacted*/}, expire:0, signingMethod:csjwt.Signer(nil), jwtVerify:(*csjwt.Verification)(nil), enableJTI:false, errorHandler:ctxhttp.Handler(nil), keyFunc:(csjwt.Keyfunc)(nil), newClaims:(func() csjwt.Claimer)(nil)}`, fmt.Sprintf("%#v", sc))
}
