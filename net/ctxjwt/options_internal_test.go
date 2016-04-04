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

	"testing"
)

func TestInternalOptionWithErrorHandler(t *testing.T) {
	t.Parallel()

	jwts := MustNewService()

	defaultErrH := jwts.DefaultErrorHandler

	wsErrH := ctxhttp.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
		http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
		return nil
	})

	if err := jwts.Options(WithErrorHandler(scope.WebsiteID, 22, wsErrH)); err != nil {
		t.Fatal(err)
	}

	cstesting.EqualPointers(t, defaultErrH, jwts.DefaultErrorHandler)
	cstesting.EqualPointers(t, wsErrH, jwts.scopeCache[scope.NewHash(scope.WebsiteID, 22)].errorHandler)
	cstesting.UnEqualPointers(t, defaultErrH, jwts.scopeCache[scope.NewHash(scope.WebsiteID, 22)].errorHandler)

	if err := jwts.Options(WithErrorHandler(scope.DefaultID, 0, wsErrH)); err != nil {
		t.Fatal(err)
	}
	cstesting.UnEqualPointers(t, defaultErrH, jwts.DefaultErrorHandler)
	cstesting.EqualPointers(t, wsErrH, jwts.DefaultErrorHandler)
}
