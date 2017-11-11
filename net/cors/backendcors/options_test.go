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

	"github.com/corestoreio/cspkg/config/cfgmock"
	"github.com/corestoreio/cspkg/net/cors"
	corstest "github.com/corestoreio/cspkg/net/cors/internal"
	"github.com/corestoreio/cspkg/net/mw"
	"github.com/corestoreio/cspkg/store/scope"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/log/logw"
	"github.com/stretchr/testify/assert"
)

// Tests partially loaded configuration and other settings will be applied from
// the backend.
func TestConfiguration_Partially_HierarchicalConfig(t *testing.T) {
	exposedHeaders := []string{"X-Header-1", "X-Header-2"}

	scpCfgSrv := cfgmock.NewService(cfgmock.PathValue{
		// Important backend.ExposedHeaders has not been set and is not
		// available in the backend configuration.
		backend.AllowedOrigins.MustFQWebsite(3): "x.com\ny.com",
		backend.AllowedMethods.MustFQ():         "PUT\nDEL\nCUT",
	}).NewScoped(3, 0)

	srv := cors.MustNew(
		cors.WithSettings(cors.Settings{
			ExposedHeaders: exposedHeaders,
		}, scope.Website.Pack(3)),
		cors.WithMarkPartiallyApplied(true, scope.Website.Pack(3)),
		cors.WithOptionFactory(backend.PrepareOptionFactory()),
	)
	scpCfg, err := srv.ConfigByScopedGetter(scpCfgSrv)
	assert.NoError(t, err, "%+v", err)

	assert.Exactly(t, []string{`x.com`, `y.com`}, scpCfg.AllowedOrigins)
	assert.Exactly(t, []string{"PUT", "DEL", "CUT"}, scpCfg.AllowedMethods)
	assert.Exactly(t, []string{}, scpCfg.ExposedHeaders)
	// TODO: To make the next line possible and remove the above line for checking
	// []string{} there needs some refactoring in cors.WithSettings() to only
	// set values which are available in the backend configuration.
	//assert.Exactly(t, exposedHeaders, scpCfg.ExposedHeaders)
}

func TestConfiguration_HierarchicalConfig(t *testing.T) {

	scpCfgSrv := cfgmock.NewService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(3): "x.com\ny.com",
		backend.AllowedMethods.MustFQ():         "PUT\nDEL\nCUT",
	}).NewScoped(3, 0)

	srv := cors.MustNew(
		cors.WithOptionFactory(backend.PrepareOptionFactory()),
	)
	scpCfg, err := srv.ConfigByScopedGetter(scpCfgSrv)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, []string{`x.com`, `y.com`}, scpCfg.AllowedOrigins)
	assert.Exactly(t, []string{"PUT", "DEL", "CUT"}, scpCfg.AllowedMethods)
}

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
		cors.WithOptionFactory(backend.PrepareOptionFactory()),
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
		backend.AllowedOrigins.MustFQWebsite(2): "http://foobar.com",
	})
	req := reqWithStore("GET")
	corstest.TestAllowedOrigin(t, s, req)
}

func TestWildcardOrigin(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2): "http://*.bar.com",
	})
	req := reqWithStore("GET")
	corstest.TestWildcardOrigin(t, s, req)
}

func TestDisallowedOrigin(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2): "http://foobar.com",
	})
	req := reqWithStore("GET")
	corstest.TestDisallowedOrigin(t, s, req)
}

func TestDisallowedWildcardOrigin(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2): "http://*.bar.com",
	})
	req := reqWithStore("GET")
	corstest.TestDisallowedWildcardOrigin(t, s, req)
}

func TestAllowedOriginFunc(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowOriginRegex.MustFQWebsite(2): "^http://foo",
	})
	req := reqWithStore("GET")
	corstest.TestAllowedOriginFunc(t, s, req)
}

func TestAllowedMethodNoPassthrough(t *testing.T) {
	var logBuf = new(log.MutexBuffer)

	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2): "http://foobar.com",
		backend.AllowedMethods.MustFQWebsite(2): "PUT\nDELETE",
		// backend.NetCorsOptionsPassthrough.MustFQWebsite(2): false, <== this is the default value
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
		backend.AllowedOrigins.MustFQWebsite(2):     "http://foobar.com",
		backend.AllowedMethods.MustFQWebsite(2):     "PUT\nDELETE",
		backend.OptionsPassthrough.MustFQWebsite(2): true,
	})
	req := reqWithStore("OPTIONS")
	req.Body = ioutil.NopCloser(strings.NewReader("Body of TestAllowedMethod_Passthrough"))
	corstest.TestAllowedMethodPassthrough(t, s, req)
}

func TestDisallowedMethod(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2): "http://foobar.com",
		backend.AllowedMethods.MustFQWebsite(2): "PUT\nDELETE",
	})

	req := reqWithStore("OPTIONS")

	corstest.TestDisallowedMethod(t, s, req)
}

func TestAllowedHeader(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2): "http://foobar.com",
		backend.AllowedHeaders.MustFQWebsite(2): "X-Header-1\nx-header-2",
	})

	req := reqWithStore("OPTIONS")

	corstest.TestAllowedHeader(t, s, req)
}

func TestAllowedWildcardHeader(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2): "http://foobar.com",
		backend.AllowedHeaders.MustFQWebsite(2): "*",
	})

	req := reqWithStore("OPTIONS")
	corstest.TestAllowedWildcardHeader(t, s, req)
}

func TestDisallowedHeader(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2): "http://foobar.com",
		backend.AllowedHeaders.MustFQWebsite(2): "X-Header-1\nx-header-2",
	})

	req := reqWithStore("OPTIONS")
	corstest.TestDisallowedHeader(t, s, req)
}

func TestExposedHeader(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2): "http://foobar.com",
		backend.ExposedHeaders.MustFQWebsite(2): "X-Header-1\nx-header-2",
	})

	req := reqWithStore("GET")
	corstest.TestExposedHeader(t, s, req)
}

func TestAllowedCredentials(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2):   "http://foobar.com",
		backend.AllowCredentials.MustFQWebsite(2): true,
	})

	req := reqWithStore("OPTIONS")
	corstest.TestAllowedCredentials(t, s, req)
}
func TestMaxAge(t *testing.T) {
	s := newCorsService(cfgmock.PathValue{
		backend.AllowedOrigins.MustFQWebsite(2): "http://foobar.com",
		backend.MaxAge.MustFQWebsite(2):         "30",
	})

	req := reqWithStore("OPTIONS")
	corstest.TestMaxAge(t, s, req)
}

func TestBackend_Path_Errors(t *testing.T) {

	tests := []struct {
		toPath func(int64) string
		val    interface{}
		errBhf errors.BehaviourFunc
	}{
		{backend.ExposedHeaders.MustFQWebsite, struct{}{}, errors.IsNotValid},
		{backend.AllowedOrigins.MustFQWebsite, struct{}{}, errors.IsNotValid},
		{backend.AllowOriginRegex.MustFQWebsite, struct{}{}, errors.IsNotValid},
		{backend.AllowOriginRegex.MustFQWebsite, "[a-z+", errors.IsFatal},
		{backend.AllowedMethods.MustFQWebsite, struct{}{}, errors.IsNotValid},
		{backend.AllowedHeaders.MustFQWebsite, struct{}{}, errors.IsNotValid},
		{backend.AllowCredentials.MustFQWebsite, struct{}{}, errors.IsNotValid},
		{backend.OptionsPassthrough.MustFQWebsite, struct{}{}, errors.IsNotValid},
		{backend.MaxAge.MustFQWebsite, struct{}{}, errors.IsNotValid},
	}
	for i, test := range tests {

		scpFnc := backend.PrepareOptionFactory()
		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			test.toPath(2): test.val,
		})
		cfgScp := cfgSrv.NewScoped(2, 0)

		_, err := cors.New(scpFnc(cfgScp)...)
		assert.True(t, test.errBhf(err), "Index %d Error: %+v", i, err)
	}
}
