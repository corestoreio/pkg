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

package jwt

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/stretchr/testify/assert"
)

func TestInternalOptionWithErrorHandler(t *testing.T) {

	jwts := MustNew()

	wsErrH := mw.ErrorWithStatusCode(http.StatusAccepted)

	if err := jwts.Options(WithErrorHandler(scope.Website, 22, wsErrH)); err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, jwts.defaultScopeCache.ErrorHandler)
	cstesting.EqualPointers(t, wsErrH, jwts.scopeCache[scope.NewHash(scope.Website, 22)].ErrorHandler)

	if err := jwts.Options(WithErrorHandler(scope.Default, 0, wsErrH)); err != nil {
		t.Fatal(err)
	}
	cstesting.EqualPointers(t, wsErrH, jwts.defaultScopeCache.ErrorHandler)
}

func TestInternalOptionNoLeakage(t *testing.T) {

	sc := ScopedConfig{
		Key: csjwt.WithPasswordRandom(),
	}
	assert.Contains(t, fmt.Sprintf("%v", sc), `csjwt.Key{/*redacted*/}`)
	assert.Contains(t, fmt.Sprintf("%#v", sc), `csjwt.Key{/*redacted*/}`)
}
