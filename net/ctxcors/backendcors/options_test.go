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
	"github.com/corestoreio/csfw/net/ctxcors"
	"github.com/corestoreio/csfw/net/ctxcors/backendcors"
	"github.com/corestoreio/csfw/net/ctxcors/internal/corstest"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
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

func newCorsService() *ctxcors.Service {
	return ctxcors.MustNew(
		ctxcors.WithBackend(backendcors.BackendOptions(backend)),
	)
}

func TestAllowedMethodPassthrough(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCtxcorsAllowedOrigins.ToPath, scope.Website, 2):     "http://foobar.com",
		mustToPath(t, backend.NetCtxcorsAllowedMethods.ToPath, scope.Website, 2):     "PUT\nDELETE",
		mustToPath(t, backend.NetCtxcorsOptionsPassthrough.ToPath, scope.Website, 2): true,
	}))
	corstest.TestAllowedMethodPassthrough(t, s, req)
}

func TestDisallowedMethod(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCtxcorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
		mustToPath(t, backend.NetCtxcorsAllowedMethods.ToPath, scope.Website, 2): "PUT\nDELETE",
	}))

	corstest.TestDisallowedMethod(t, s, req)
}

func TestAllowedHeader(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCtxcorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
		mustToPath(t, backend.NetCtxcorsAllowedHeaders.ToPath, scope.Website, 2): "X-Header-1\nx-header-2",
	}))

	corstest.TestAllowedHeader(t, s, req)
}

func TestAllowedWildcardHeader(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, backend.NetCtxcorsAllowedOrigins.ToPath, scope.Website, 2): "http://foobar.com",
		mustToPath(t, backend.NetCtxcorsAllowedHeaders.ToPath, scope.Website, 2): "*",
	}))

	corstest.TestAllowedWildcardHeader(t, s, req)
}

//func TestWithBackendApplied(t *testing.T) {
//
//	be := initBackend(t)
//
//	cfgGet := cfgmock.NewService(
//		cfgmock.WithPV(cfgmock.PathValue{
//			mustToPath(t, be.NetCtxcorsExposedHeaders.FQ, scope.Website, 2):     "X-CoreStore-ID\nContent-Type\n\n",
//			mustToPath(t, be.NetCtxcorsAllowedOrigins.FQ, scope.Website, 2):     "host1.com\nhost2.com\n\n",
//			mustToPath(t, be.NetCtxcorsAllowedMethods.FQ, scope.Default, 0):     "PATCH\nDELETE",
//			mustToPath(t, be.NetCtxcorsAllowedHeaders.FQ, scope.Default, 0):     "Date,X-Header1",
//			mustToPath(t, be.NetCtxcorsAllowCredentials.FQ, scope.Website, 2):   "1",
//			mustToPath(t, be.NetCtxcorsOptionsPassthrough.FQ, scope.Website, 2): "1",
//			mustToPath(t, be.NetCtxcorsMaxAge.FQ, scope.Website, 2):             "2h",
//		}),
//	)
//
//	c := MustNew(WithBackendApplied(be, cfgGet.NewScoped(2, 4)))
//
//	assert.Exactly(t, []string{"X-Corestore-Id", "Content-Type"}, c.exposedHeaders)
//	assert.Exactly(t, []string{"host1.com", "host2.com"}, c.allowedOrigins)
//	assert.Exactly(t, []string{"PATCH", "DELETE"}, c.allowedMethods)
//	assert.Exactly(t, []string{"Date,X-Header1", "Origin"}, c.allowedHeaders)
//	assert.Exactly(t, true, c.AllowCredentials)
//	assert.Exactly(t, true, c.OptionsPassthrough)
//	assert.Exactly(t, "7200", c.maxAge)
//}

//func TestWithBackendAppliedErrors(t *testing.T) {
//
//	be := initBackend(t)
//
//	cfgErr := errors.New("Test Error")
//	cfgGet := cfgmock.NewService(
//		cfgmock.WithBool(func(_ string) (bool, error) {
//			return false, cfgErr
//		}),
//		cfgmock.WithString(func(_ string) (string, error) {
//			return "", cfgErr
//		}),
//	)
//
//	c, err := New(WithBackendApplied(be, cfgGet.NewScoped(223, 43213)))
//	assert.Nil(t, c)
//	assert.EqualError(t, err, "Route net/ctxcors/exposed_headers: Test Error\nRoute net/ctxcors/allowed_origins: Test Error\nRoute net/ctxcors/allowed_methods: Test Error\nRoute net/ctxcors/allowed_headers: Test Error\nRoute net/ctxcors/allow_credentials: Test Error\nRoute net/ctxcors/allow_credentials: Test Error\nRoute net/ctxcors/max_age: Test Error\nMaxAge: Invalid Duration seconds: 0")
//}
