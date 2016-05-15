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

package backendcors_test

import (
	"net/http"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/net/mwcors"
	"github.com/corestoreio/csfw/net/mwcors/backendcors"
	"github.com/corestoreio/csfw/net/mwcors/internal/corstest"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func mustToPath(t *testing.T, f func(s scope.Scope, scopeID int64) (cfgpath.Path, error), s scope.Scope, scopeID int64) string {
	p, err := f(s, scopeID)
	if err != nil {
		t.Fatal(errors.PrintLoc(err))
	}
	return p.String()
}

func reqWithStore(method string, cfgOpt ...cfgmock.OptionFunc) *http.Request {
	req, err := http.NewRequest(method, "http://corestore.io/foo", nil)
	if err != nil {
		panic(err)
	}

	return req.WithContext(
		store.WithContextRequestedStore(req.Context(), storemock.MustNewStoreAU(cfgmock.NewService(cfgOpt...))),
	)
}

func newCorsService() *mwcors.Service {
	return mwcors.MustNew(
		mwcors.WithOptionFactory(backendcors.PrepareOptions(backend)),
	)
}

func TestNoConfig(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET")
	corstest.TestNoConfig(t, s, req)
}

func TestMatchAllOrigin(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
	// STAR is the default value in the element structure
	}))
	corstest.TestMatchAllOrigin(t, s, req)
}

func TestAllowedOrigin(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
	}))
	corstest.TestAllowedOrigin(t, s, req)
}

func TestWildcardOrigin(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://*.bar.com",
	}))
	corstest.TestWildcardOrigin(t, s, req)
}

func TestDisallowedOrigin(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
	}))
	corstest.TestDisallowedOrigin(t, s, req)
}

func TestDisallowedWildcardOrigin(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://*.bar.com",
	}))
	corstest.TestDisallowedWildcardOrigin(t, s, req)
}

func TestAllowedOriginFunc(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowOriginRegex.ToPath, scope.Website, 2): "^http://foo",
	}))
	corstest.TestAllowedOriginFunc(t, s, req)
}

func TestAllowedMethod(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
		mustToPath(t, backend.NetCorsAllowedMethods.ToPath, scope.Website, 2): "PUT\nDELETE",
	}))
	corstest.TestAllowedMethod(t, s, req)
}

func TestAllowedMethodPassthrough(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2):     "http://foobar.com",
		mustToPath(t, backend.NetCorsAllowedMethods.ToPath, scope.Website, 2):     "PUT\nDELETE",
		mustToPath(t, backend.NetCorsOptionsPassthrough.ToPath, scope.Website, 2): true,
	}))
	corstest.TestAllowedMethodPassthrough(t, s, req)
}

func TestDisallowedMethod(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
		mustToPath(t, backend.NetCorsAllowedMethods.ToPath, scope.Website, 2): "PUT\nDELETE",
	}))

	corstest.TestDisallowedMethod(t, s, req)
}

func TestAllowedHeader(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
		mustToPath(t, backend.NetCorsAllowedHeaders.ToPath, scope.Website, 2): "X-Header-1\nx-header-2",
	}))

	corstest.TestAllowedHeader(t, s, req)
}

func TestAllowedWildcardHeader(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
		mustToPath(t, backend.NetCorsAllowedHeaders.ToPath, scope.Website, 2): "*",
	}))

	corstest.TestAllowedWildcardHeader(t, s, req)
}

func TestDisallowedHeader(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
		mustToPath(t, backend.NetCorsAllowedHeaders.ToPath, scope.Website, 2): "X-Header-1\nx-header-2",
	}))

	corstest.TestDisallowedHeader(t, s, req)
}

func TestExposedHeader(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
		mustToPath(t, backend.NetCorsExposedHeaders.ToPath, scope.Website, 2): "X-Header-1\nx-header-2",
	}))

	corstest.TestExposedHeader(t, s, req)
}

func TestAllowedCredentials(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2):   "http://foobar.com",
		mustToPath(t, backend.NetCorsAllowCredentials.ToPath, scope.Website, 2): true,
	}))

	corstest.TestAllowedCredentials(t, s, req)
}
func TestMaxAge(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
		mustToPath(t, backend.NetCorsMaxAge.ToPath, scope.Website, 2):         "30s",
	}))

	corstest.TestMaxAge(t, s, req)
}

func TestBackend_Path_Errors(t *testing.T) {

	tests := []struct {
		toPath func(s scope.Scope, scopeID int64) (cfgpath.Path, error)
		val    interface{}
		errBhf errors.BehaviourFunc
	}{
		{backend.NetCorsExposedHeaders.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetCorsAllowedOrigins.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetCorsAllowOriginRegex.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetCorsAllowOriginRegex.ToPath, "[a-z+", errors.IsFatal},
		{backend.NetCorsAllowedMethods.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetCorsAllowedHeaders.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetCorsAllowCredentials.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetCorsOptionsPassthrough.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetCorsMaxAge.ToPath, struct{}{}, errors.IsNotValid},
	}
	for i, test := range tests {

		scpFnc := backendcors.PrepareOptions(backend)
		cfgSrv := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			mustToPath(t, test.toPath, scope.Website, 2): test.val,
		}))
		cfgScp := cfgSrv.NewScoped(2, 0)

		_, err := mwcors.New(scpFnc(cfgScp)...)
		assert.True(t, test.errBhf(err), "Index %d Error: %s", i, err)
	}
}
