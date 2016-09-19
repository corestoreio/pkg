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

package backendgeoip_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/net/geoip/backendgeoip"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func init() {
	now := time.Now()
	seed := now.Unix() + int64(now.Nanosecond()) + 12345*int64(os.Getpid())
	rand.Seed(seed)
}

func mustToPath(t interface {
	Fatalf(string, ...interface{})
}, f func(s scope.Scope, scopeID int64) (cfgpath.Path, error), s scope.Scope, scopeID int64) string {
	p, err := f(s, scopeID)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	return p.String()
}

func TestBackend_WithGeoIP2Webservice_Redis(t *testing.T) {

	t.Run("Error_API", testBackend_WithGeoIP2Webservice_Redis(
		func() *http.Client {
			// http://dev.maxmind.com/geoip/geoip2/web-services/#Errors
			return &http.Client{
				Transport: cstesting.NewHTTPTrip(402, `{"error":"The license key you have provided is out of queries.","code":"OUT_OF_QUERIES"}`, nil),
			}
		},
		func(t *testing.T) http.Handler {
			return http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				panic("Should not get called")
			})
		},
		http.StatusServiceUnavailable,
	))

	t.Run("Error_JSON", testBackend_WithGeoIP2Webservice_Redis(
		func() *http.Client {
			// http://dev.maxmind.com/geoip/geoip2/web-services/#Errors
			return &http.Client{
				Transport: cstesting.NewHTTPTrip(200, `{"error":"The license ... wow this JSON isn't valid.`, nil),
			}
		},
		func(t *testing.T) http.Handler {
			return http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		},
		http.StatusServiceUnavailable,
	))

	var calledSuccessHandler int32
	t.Run("Success", testBackend_WithGeoIP2Webservice_Redis(
		func() *http.Client {
			return &http.Client{
				Transport: cstesting.NewHTTPTrip(200, `{ "continent": { "code": "EU", "geoname_id": 6255148, "names": { "de": "Europa", "en": "Europe", "es": "Europa", "fr": "Europe", "ja": "ヨーロッパ", "pt-BR": "Europa", "ru": "Европа", "zh-CN": "欧洲" } }, "country": { "geoname_id": 2921044, "iso_code": "DE", "names": { "de": "Deutschland", "en": "Germany", "es": "Alemania", "fr": "Allemagne", "ja": "ドイツ連邦共和国", "pt-BR": "Alemanha", "ru": "Германия", "zh-CN": "德国" } }, "registered_country": { "geoname_id": 2921044, "iso_code": "DE", "names": { "de": "Deutschland", "en": "Germany", "es": "Alemania", "fr": "Allemagne", "ja": "ドイツ連邦共和国", "pt-BR": "Alemanha", "ru": "Германия", "zh-CN": "德国" } }, "traits": { "autonomous_system_number": 1239, "autonomous_system_organization": "Linkem IR WiMax Network", "domain": "example.com", "is_anonymous_proxy": true, "is_satellite_provider": true, "isp": "Linkem spa", "ip_address": "1.2.3.4", "organization": "Linkem IR WiMax Network", "user_type": "traveler" }, "maxmind": { "queries_remaining": 54321 } }`, nil),
			}
		},
		func(t *testing.T) http.Handler {
			return http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				cty, ok := geoip.FromContextCountry(r.Context())
				assert.True(t, ok)
				assert.Exactly(t, "DE", cty.Country.IsoCode)
				atomic.AddInt32(&calledSuccessHandler, 1)
			})
		},
		http.StatusOK,
	))
	assert.Exactly(t, int32(80), atomic.LoadInt32(&calledSuccessHandler), "calledSuccessHandler")
}

func testBackend_WithGeoIP2Webservice_Redis(
	hcf func() *http.Client,
	finalHandler func(t *testing.T) http.Handler,
	wantCode int,
) func(*testing.T) {

	return func(t *testing.T) {
		rd := miniredis.NewMiniRedis()
		if err := rd.Start(); err != nil {
			t.Fatal(err)
		}
		defer rd.Close()
		redConURL := fmt.Sprintf("redis://%s/3", rd.Addr())

		// test if we get the correct country and if the country has
		// been successfully stored in redis and can be retrieved.

		cfgStruct, err := backendgeoip.NewConfigStructure()
		if err != nil {
			t.Fatal(err)
		}
		be := backendgeoip.New(cfgStruct)
		be.WebServiceClient = hcf()
		scpFnc := backendgeoip.PrepareOptions(be)

		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			// @see structure.go for the limitation to scope.Default
			mustToPath(t, backend.NetGeoipMaxmindWebserviceUserID.ToPath, scope.Default, 0):   `TestUserID`,
			mustToPath(t, backend.NetGeoipMaxmindWebserviceLicense.ToPath, scope.Default, 0):  `TestLicense`,
			mustToPath(t, backend.NetGeoipMaxmindWebserviceTimeout.ToPath, scope.Default, 0):  `21s`,
			mustToPath(t, backend.NetGeoipMaxmindWebserviceRedisURL.ToPath, scope.Default, 0): redConURL,
		})
		cfgScp := cfgSrv.NewScoped(1, 2) // Website ID 2 == euro / Store ID == 2 Austria ==> here doesn't matter

		geoSrv := geoip.MustNew()

		req := func() *http.Request {
			req, _ := http.NewRequest("GET", "http://corestore.io", nil)
			req.Header.Set("X-Cluster-Client-Ip", "2a02:d180::") // Germany
			return req
		}()

		if err := geoSrv.Options(scpFnc(cfgScp)...); err != nil {
			t.Fatalf("%+v", err)
		}
		// food for the race detector
		hpu := cstesting.NewHTTPParallelUsers(8, 10, 500, time.Millisecond) // 8,10
		hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
			assert.Exactly(t, wantCode, rec.Code)
		}
		hpu.ServeHTTP(req, geoSrv.WithCountryByIP(finalHandler(t)))
	}
}

func TestBackend_WithAlternativeRedirect(t *testing.T) {

	t.Run("LocalFile", backend_WithAlternativeRedirect(cfgmock.NewService(cfgmock.PathValue{
		// @see structure.go why scope.Store and scope.Website can be used.
		mustToPath(t, backend.NetGeoipAlternativeRedirect.ToPath, scope.Store, 2):       `https://byebye.de.io`,
		mustToPath(t, backend.NetGeoipAlternativeRedirectCode.ToPath, scope.Website, 1): 307,
		mustToPath(t, backend.NetGeoipAllowedCountries.ToPath, scope.Store, 2):          "AT,CH",
		mustToPath(t, backend.NetGeoipMaxmindLocalFile.ToPath, scope.Default, 0):        filepath.Join("..", "testdata", "GeoIP2-Country-Test.mmdb"),
	})))

	t.Run("WebService", backend_WithAlternativeRedirect(cfgmock.NewService(cfgmock.PathValue{
		// @see structure.go why scope.Store and scope.Website can be used.
		mustToPath(t, backend.NetGeoipAlternativeRedirect.ToPath, scope.Store, 2):        `https://byebye.de.io`,
		mustToPath(t, backend.NetGeoipAlternativeRedirectCode.ToPath, scope.Website, 1):  307,
		mustToPath(t, backend.NetGeoipAllowedCountries.ToPath, scope.Store, 2):           "AT,CH",
		mustToPath(t, backend.NetGeoipMaxmindWebserviceUserID.ToPath, scope.Default, 0):  "LiesschenMueller",
		mustToPath(t, backend.NetGeoipMaxmindWebserviceLicense.ToPath, scope.Default, 0): "8x4",
		mustToPath(t, backend.NetGeoipMaxmindWebserviceTimeout.ToPath, scope.Default, 0): "3s",
	})))
}

func backend_WithAlternativeRedirect(cfgSrv *cfgmock.Service) func(*testing.T) {
	return func(t *testing.T) {
		logBuf := new(log.MutexBuffer)

		cfgStruct, err := backendgeoip.NewConfigStructure()
		if err != nil {
			t.Fatal(err)
		}
		be := backendgeoip.New(cfgStruct)
		be.WebServiceClient = &http.Client{
			Transport: cstesting.NewHTTPTrip(200, `{ "continent": { "code": "EU", "geoname_id": 6255148, "names": { "de": "Europa", "en": "Europe", "ru": "Европа", "zh-CN": "欧洲" } }, "country": { "geoname_id": 2921044, "iso_code": "DE", "names": { "de": "Deutschland", "en": "Germany", "es": "Alemania", "fr": "Allemagne", "ja": "ドイツ連邦共和国", "pt-BR": "Alemanha", "ru": "Германия", "zh-CN": "德国" } }, "maxmind": { "queries_remaining": 54321 } }`, nil),
		}
		scpFnc := backendgeoip.PrepareOptions(be)
		geoSrv := geoip.MustNew(
			geoip.WithRootConfig(cfgSrv),
			geoip.WithDebugLog(logBuf),
			geoip.WithOptionFactory(scpFnc),
			geoip.WithServiceErrorHandler(mw.ErrorWithPanic),
			geoip.WithErrorHandler(scope.DefaultHash, mw.ErrorWithPanic),
		)

		// if you try to set the allowed countries with this option, they get
		// overwritten by the ScopeConfig service.
		// if err := geoSrv.Options(geoip.WithAllowedCountryCodes(scope.Store, 2, "AT", "CH")); err != nil {
		//	t.Fatalf("%+v", err)
		// }

		// Germany is not allowed and must be redirected to https://byebye.de.io with code 307
		req := func() *http.Request {
			req := httptest.NewRequest("GET", "http://corestore.io", nil)
			req.RemoteAddr = "2a02:d180::"
			return req.WithContext(scope.WithContext(req.Context(), 1, 2)) // Website ID 1 == euro / Store ID == 2 Austria
		}()

		hpu := cstesting.NewHTTPParallelUsers(8, 12, 600, time.Millisecond) // 8, 12
		hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
			assert.Exactly(t, `https://byebye.de.io`, rec.Header().Get("Location"))
			assert.Exactly(t, 307, rec.Code)
		}

		// Food for the race detector
		hpu.ServeHTTP(req,
			geoSrv.WithIsCountryAllowedByIP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not be called")
			})),
		)

		// Min 20 calls IsValid
		// Exactly one call to optionInflight.Do
		if have, want := strings.Count(logBuf.String(), `geoip.WithIsCountryAllowedByIP.checkAllow.false`), 90; have < want {
			t.Errorf("ConfigByScopedGetter.IsValid: Have: %d < Want: %d", have, want)
		}
		if have, want := strings.Count(logBuf.String(), `geoip.Service.ConfigByScopedGetter.Inflight.Do`), 1; have != want {
			t.Errorf("ConfigByScopedGetter.optionInflight.Do: Have: %d Want: %d", have, want)
		}
		//	println("\n", logBuf.String(), "\n")
	}
}

func TestBackend_Path_Errors(t *testing.T) {

	tests := []struct {
		toPath func(s scope.Scope, scopeID int64) (cfgpath.Path, error)
		val    interface{}
		errBhf errors.BehaviourFunc
	}{
		{backend.NetGeoipAllowedCountries.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetGeoipAlternativeRedirect.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetGeoipAlternativeRedirectCode.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetGeoipMaxmindLocalFile.ToPath, "fileNotFound.txt", errors.IsNotFound},
		{backend.NetGeoipMaxmindLocalFile.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetGeoipMaxmindWebserviceUserID.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetGeoipMaxmindWebserviceLicense.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetGeoipMaxmindWebserviceTimeout.ToPath, struct{}{}, errors.IsNotValid},
		{backend.NetGeoipMaxmindWebserviceRedisURL.ToPath, struct{}{}, errors.IsNotValid},
	}
	for i, test := range tests {

		cStruct, err := backendgeoip.NewConfigStructure()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		be := backendgeoip.New(cStruct)
		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			mustToPath(t, test.toPath, scope.Default, 0): test.val,
		})

		gs := geoip.MustNew(
			geoip.WithRootConfig(cfgSrv),
			geoip.WithOptionFactory(backendgeoip.PrepareOptions(be)),
		)
		assert.NoError(t, gs.ClearCache())
		scpdCfg := gs.ConfigByScope(0, 0)
		err = scpdCfg.IsValid()
		assert.True(t, test.errBhf(err), "Index %d Error: %s", i, err)
	}
}
