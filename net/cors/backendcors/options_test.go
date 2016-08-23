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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/net/cors"
	"github.com/corestoreio/csfw/net/cors/backendcors"
	corstest "github.com/corestoreio/csfw/net/cors/internal"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

type fataler interface {
	Fatal(args ...interface{})
}

func reqWithStore(method string) *http.Request {
	req := httptest.NewRequest(method, "https://corestore.io/catalog/product/id/33454", nil)
	return req.WithContext(
		scope.WithContext(req.Context(), 2, 5), // website 2 = OZ; store = 5 AU
	)
}

func newCorsService(pv cfgmock.PathValue) *cors.Service {
	return cors.MustNew(
		cors.WithRootConfig(cfgmock.NewService(pv)),
		cors.WithOptionFactory(backendcors.PrepareOptions(backend)),
		cors.WithServiceErrorHandler(mw.ErrorWithPanic),
	)
}

func TestNoConfig(t *testing.T) {
	s := newCorsService(nil)
	req := reqWithStore("GET")
	corstest.TestNoConfig(t, s, req)
}

func TestMatchAllOrigin(t *testing.T) {
	var logBuf = new(log.MutexBuffer)
	s := newCorsService(nil) // STAR is the default value in the element structure
	req := reqWithStore("GET")

	if err := s.Options(cors.WithLogger(logw.NewLog(logw.WithWriter(logBuf), logw.WithLevel(logw.LevelFatal)))); err != nil {
		t.Fatal(err)
	}
	corstest.TestMatchAllOrigin(t, s, req)
	//println("\n", logBuf.String())
}

func TestAllowedOrigin(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
	})
	req := reqWithStore("GET")
	corstest.TestAllowedOrigin(t, s, req)
}

func TestWildcardOrigin(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://*.bar.com",
	})
	req := reqWithStore("GET")
	corstest.TestWildcardOrigin(t, s, req)
}

func TestDisallowedOrigin(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
	})
	req := reqWithStore("GET")
	corstest.TestDisallowedOrigin(t, s, req)
}

func TestDisallowedWildcardOrigin(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://*.bar.com",
	})
	req := reqWithStore("GET")
	corstest.TestDisallowedWildcardOrigin(t, s, req)
}

func TestAllowedOriginFunc(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowOriginRegex.MustFQ(scope.Website, 2): "^http://foo",
	})
	req := reqWithStore("GET")
	corstest.TestAllowedOriginFunc(t, s, req)
}

func TestAllowedMethodNoPassthrough(t *testing.T) {
	var logBuf = new(log.MutexBuffer)

	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsAllowedMethods.MustFQ(scope.Website, 2): "PUT\nDELETE",
		// backend.NetCorsOptionsPassthrough.MustFQ(scope.Website, 2): false, <== this is the default value
	})
	if err := s.Options(cors.WithLogger(logw.NewLog(logw.WithWriter(logBuf), logw.WithLevel(logw.LevelDebug)))); err != nil {
		t.Fatal(err)
	}

	req := reqWithStore("OPTIONS")
	req.Body = ioutil.NopCloser(strings.NewReader("Body of TestAllowedMethod_No_Passthrough"))
	corstest.TestAllowedMethodNoPassthrough(t, s, req)

	if have, want := strings.Count(logBuf.String(), `Service.ConfigByScopedGetter.Inflight.Do`), 1; have != want {
		//println("\n", logBuf.String())
		t.Fatalf("Have: %v Want: %v", have, want)
	}
	if have, want := strings.Count(logBuf.String(), `cors.Service.ConfigByScopedGetter.IsValid`), 88; have <= want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
	//println("\n", logBuf.String())
}

func TestAllowedMethodPassthrough(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2):     "http://foobar.com",
		backend.NetCorsAllowedMethods.MustFQ(scope.Website, 2):     "PUT\nDELETE",
		backend.NetCorsOptionsPassthrough.MustFQ(scope.Website, 2): true,
	})
	req := reqWithStore("OPTIONS")
	req.Body = ioutil.NopCloser(strings.NewReader("Body of TestAllowedMethod_Passthrough"))
	corstest.TestAllowedMethodPassthrough(t, s, req)
}

func TestDisallowedMethod(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsAllowedMethods.MustFQ(scope.Website, 2): "PUT\nDELETE",
	})

	req := reqWithStore("OPTIONS")

	corstest.TestDisallowedMethod(t, s, req)
}

func TestAllowedHeader(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsAllowedHeaders.MustFQ(scope.Website, 2): "X-Header-1\nx-header-2",
	})

	req := reqWithStore("OPTIONS")

	corstest.TestAllowedHeader(t, s, req)
}

func TestAllowedWildcardHeader(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsAllowedHeaders.MustFQ(scope.Website, 2): "*",
	})

	req := reqWithStore("OPTIONS")
	corstest.TestAllowedWildcardHeader(t, s, req)
}

func TestDisallowedHeader(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsAllowedHeaders.MustFQ(scope.Website, 2): "X-Header-1\nx-header-2",
	})

	req := reqWithStore("OPTIONS")
	corstest.TestDisallowedHeader(t, s, req)
}

func TestExposedHeader(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsExposedHeaders.MustFQ(scope.Website, 2): "X-Header-1\nx-header-2",
	})

	req := reqWithStore("GET")
	corstest.TestExposedHeader(t, s, req)
}

func TestAllowedCredentials(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2):   "http://foobar.com",
		backend.NetCorsAllowCredentials.MustFQ(scope.Website, 2): true,
	})

	req := reqWithStore("OPTIONS")
	corstest.TestAllowedCredentials(t, s, req)
}
func TestMaxAge(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.NetCorsAllowedOrigins.MustFQ(scope.Website, 2): "http://foobar.com",
		backend.NetCorsMaxAge.MustFQ(scope.Website, 2):         "30",
	})

	req := reqWithStore("OPTIONS")
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
		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			test.toPath(scope.Website, 2): test.val,
		})
		cfgScp := cfgSrv.NewScoped(2, 0)

		_, err := cors.New(scpFnc(cfgScp)...)
		assert.True(t, test.errBhf(err), "Index %d Error: %+v", i, err)
	}
}
