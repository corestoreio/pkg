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
	"strings"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/net/cors"
	"github.com/corestoreio/csfw/net/cors/backendcors"
	corstest "github.com/corestoreio/csfw/net/cors/internal"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

type fataler interface {
	Fatal(args ...interface{})
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

func newCorsService() *cors.Service {
	return cors.MustNew(
		cors.WithOptionFactory(backendcors.PrepareOptions(backend)),
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
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
	}))
	corstest.TestAllowedOrigin(t, s, req)
}

func TestWildcardOrigin(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://*.bar.com",
	}))
	corstest.TestWildcardOrigin(t, s, req)
}

func TestDisallowedOrigin(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
	}))
	corstest.TestDisallowedOrigin(t, s, req)
}

func TestDisallowedWildcardOrigin(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://*.bar.com",
	}))
	corstest.TestDisallowedWildcardOrigin(t, s, req)
}

func TestAllowedOriginFunc(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowOriginRegex.MustFQ(scope.Website, 2): "^http://foo",
	}))
	corstest.TestAllowedOriginFunc(t, s, req)
}

func TestAllowedMethod(t *testing.T) {
	var logBuf log.MutexBuffer

	s := newCorsService()
	if err := s.Options(cors.WithLogger(logw.NewLog(logw.WithWriter(&logBuf), logw.WithLevel(logw.LevelDebug)))); err != nil {
		t.Fatal(err)
	}

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsAllowedMethods.MustFQ(scope.Website, 2): "PUT\nDELETE",
	}))
	corstest.TestAllowedMethod(t, s, req)

	if have, want := strings.Count(logBuf.String(), `cors.Service.ConfigByScopedGetter.optionInflight.DoChan`), 1; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
	if have, want := strings.Count(logBuf.String(), `cors.Service.ConfigByScopedGetter.IsValid`), 9; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
}

func TestAllowedMethodPassthrough(t *testing.T) {
	s := newCorsService()
	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2):     "http://foobar.com",
		backend.NetCorsAllowedMethods.MustFQ(scope.Website, 2):     "PUT\nDELETE",
		backend.NetCorsOptionsPassthrough.MustFQ(scope.Website, 2): true,
	}))
	corstest.TestAllowedMethodPassthrough(t, s, req)
}

func TestDisallowedMethod(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsAllowedMethods.MustFQ(scope.Website, 2): "PUT\nDELETE",
	}))

	corstest.TestDisallowedMethod(t, s, req)
}

func TestAllowedHeader(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsAllowedHeaders.MustFQ(scope.Website, 2): "X-Header-1\nx-header-2",
	}))

	corstest.TestAllowedHeader(t, s, req)
}

func TestAllowedWildcardHeader(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsAllowedHeaders.MustFQ(scope.Website, 2): "*",
	}))

	corstest.TestAllowedWildcardHeader(t, s, req)
}

func TestDisallowedHeader(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsAllowedHeaders.MustFQ(scope.Website, 2): "X-Header-1\nx-header-2",
	}))

	corstest.TestDisallowedHeader(t, s, req)
}

func TestExposedHeader(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("GET", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsExposedHeaders.MustFQ(scope.Website, 2): "X-Header-1\nx-header-2",
	}))

	corstest.TestExposedHeader(t, s, req)
}

func TestAllowedCredentials(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2):   "http://foobar.com",
		backend.NetCorsAllowCredentials.MustFQ(scope.Website, 2): true,
	}))

	corstest.TestAllowedCredentials(t, s, req)
}
func TestMaxAge(t *testing.T) {
	s := newCorsService()

	req := reqWithStore("OPTIONS", cfgmock.WithPV(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsMaxAge.MustFQ(scope.Website, 2):         "30s",
	}))

	corstest.TestMaxAge(t, s, req)
}

func TestBackend_Path_Errors(t *testing.T) {

	tests := []struct {
		toPath func(s scope.Scope, scopeID int64) string
		val    interface{}
		errBhf errors.BehaviourFunc
	}{
		{backend.NetCorsExposedHeaders.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.NetCorsAllowedOrigins.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.NetCorsAllowOriginRegex.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.NetCorsAllowOriginRegex.MustFQ, "[a-z+", errors.IsFatal},
		{backend.NetCorsAllowedMethods.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.NetCorsAllowedHeaders.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.NetCorsAllowCredentials.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.NetCorsOptionsPassthrough.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.NetCorsMaxAge.MustFQ, struct{}{}, errors.IsNotValid},
	}
	for i, test := range tests {

		scpFnc := backendcors.PrepareOptions(backend)
		cfgSrv := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			test.toPath(scope.Website, 2): test.val,
		}))
		cfgScp := cfgSrv.NewScoped(2, 0)

		_, err := cors.New(scpFnc(cfgScp)...)
		assert.True(t, test.errBhf(err), "Index %d Error: %+v", i, err)
	}
}
