// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/net/geoip"
	"github.com/corestoreio/pkg/net/geoip/backendgeoip"
	"github.com/corestoreio/pkg/net/geoip/maxmindfile"
	"github.com/corestoreio/pkg/net/geoip/maxmindwebservice"
	"github.com/corestoreio/pkg/net/mw"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/util/assert"
)

func init() {
	now := time.Now()
	seed := now.Unix() + int64(now.Nanosecond()) + 12345*int64(os.Getpid())
	rand.Seed(seed)
}

func TestConfiguration_UnregisteredOptionFactoryFunc(t *testing.T) {
	scpCfgSrv := cfgmock.NewService().NewScoped(1, 3)

	srv := geoip.MustNew(
		geoip.WithOptionFactory(backend.PrepareOptionFactory()),
	)
	_, err := srv.ConfigByScopedGetter(scpCfgSrv)
	assert.True(t, errors.IsNotFound(err), "%+v", err)
}

func TestConfiguration_HierarchicalConfig(t *testing.T) {

	scpCfgSrv := cfgmock.NewService(cfgmock.PathValue{
		backend.DataSource.MustFQ():                `file`,
		backend.MaxmindLocalFile.MustFQ():          filePathGeoIP,
		backend.AllowedCountries.MustFQWebsite(1):  `AU,NZ`,
		backend.AlternativeRedirect.MustFQStore(3): `https://signin.corestore.io`,
	}).NewScoped(1, 3)

	srv := geoip.MustNew(
		geoip.WithOptionFactory(backend.PrepareOptionFactory()),
	)
	scpCfg, err := srv.ConfigByScopedGetter(scpCfgSrv)
	assert.NoError(t, err, "%+v", err)

	assert.Exactly(t, []string{`AU`, `NZ`}, scpCfg.AllowedCountries)
}

func TestConfiguration_WithAlternativeRedirect(t *testing.T) {

	t.Run("LocalFile", backend_WithAlternativeRedirect(cfgmock.NewService(cfgmock.PathValue{
		// @see structure.go why scope.Store and scope.Website can be used.
		backend.DataSource.MustFQ():                      `file`,
		backend.MaxmindLocalFile.MustFQ():                filePathGeoIP,
		backend.AlternativeRedirect.MustFQStore(2):       `https://byebye.de.io`,
		backend.AlternativeRedirectCode.MustFQWebsite(1): 307,
		backend.AllowedCountries.MustFQStore(2):          "AT,CH",
	})))

	t.Run("WebService", backend_WithAlternativeRedirect(cfgmock.NewService(cfgmock.PathValue{
		// @see structure.go why scope.Store and scope.Website can be used.
		backend.DataSource.MustFQ():                      `webservice`,
		backend.AlternativeRedirect.MustFQStore(2):       `https://byebye.de.io`,
		backend.AlternativeRedirectCode.MustFQWebsite(1): 307,
		backend.AllowedCountries.MustFQStore(2):          "AT,CH",
		backend.MaxmindWebserviceUserID.MustFQ():         "LiesschenMueller",
		backend.MaxmindWebserviceLicense.MustFQ():        "8x4",
		backend.MaxmindWebserviceTimeout.MustFQ():        "3s",
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
		be.Register(maxmindwebservice.NewOptionFactory(
			&http.Client{
				Transport: cstesting.NewHTTPTrip(200, `{ "continent": { "code": "EU", "geoname_id": 6255148, "names": { "de": "Europa", "en": "Europe", "ru": "Европа", "zh-CN": "欧洲" } }, "country": { "geoname_id": 2921044, "iso_code": "DE", "names": { "de": "Deutschland", "en": "Germany", "es": "Alemania", "fr": "Allemagne", "ja": "ドイツ連邦共和国", "pt-BR": "Alemanha", "ru": "Германия", "zh-CN": "德国" } }, "maxmind": { "queries_remaining": 54321 } }`, nil),
			},
			be.MaxmindWebserviceUserID,
			be.MaxmindWebserviceLicense,
			be.MaxmindWebserviceTimeout,
			be.MaxmindWebserviceRedisURL,
		))
		be.Register(maxmindfile.NewOptionFactory(be.MaxmindLocalFile))

		scpFnc := be.PrepareOptionFactory()
		geoSrv := geoip.MustNew(
			geoip.WithRootConfig(cfgSrv),
			geoip.WithDebugLog(logBuf),
			geoip.WithOptionFactory(scpFnc),
			geoip.WithServiceErrorHandler(mw.ErrorWithPanic),
			geoip.WithErrorHandler(mw.ErrorWithPanic), //  Default Scope ;-)
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

func TestConfiguration_Path_Errors(t *testing.T) {

	tests := []struct {
		toPath string
		val    interface{}
		errBhf errors.BehaviourFunc
	}{
		{backend.AllowedCountries.MustFQ(), struct{}{}, errors.IsNotValid},
		{backend.AlternativeRedirect.MustFQ(), struct{}{}, errors.IsNotValid},
		{backend.AlternativeRedirectCode.MustFQ(), struct{}{}, errors.IsNotValid},
		{backend.MaxmindLocalFile.MustFQ(), "fileNotFound.txt", errors.IsNotFound},
		{backend.DataSource.MustFQ(), struct{}{}, errors.IsNotValid},
	}
	for i, test := range tests {

		cStruct, err := backendgeoip.NewConfigStructure()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		be := backendgeoip.New(cStruct)
		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			test.toPath: test.val,
		})

		gs := geoip.MustNew(
			geoip.WithRootConfig(cfgSrv),
			geoip.WithOptionFactory(be.PrepareOptionFactory()),
		)
		assert.NoError(t, gs.ClearCache())
		_, err = gs.ConfigByScope(0, 0)
		assert.True(t, test.errBhf(err), "Index %d Error: %s", i, err)
	}
}

func TestNewOptionFactoryGeoSourceFile_Invalid_ConfigValue(t *testing.T) {
	name, off := maxmindfile.NewOptionFactory(backend.MaxmindLocalFile)
	assert.Exactly(t, `file`, name)

	cfgSrv := cfgmock.NewService(cfgmock.PathValue{
		backend.MaxmindLocalFile.MustFQ(): struct{}{},
	})

	gs := geoip.MustNew(
		geoip.WithRootConfig(cfgSrv),
		geoip.WithOptionFactory(off),
	)
	assert.NoError(t, gs.ClearCache())
	_, err := gs.ConfigByScope(0, 0)
	assert.True(t, errors.IsNotValid(err), " Error: %+v", err)
}

func TestNewOptionFactoryGeoSourceFile_Empty_ConfigValue(t *testing.T) {
	name, off := maxmindfile.NewOptionFactory(backend.MaxmindLocalFile)
	assert.Exactly(t, `file`, name)

	cfgSrv := cfgmock.NewService(cfgmock.PathValue{
		backend.MaxmindLocalFile.MustFQ(): "",
	})

	gs := geoip.MustNew(
		geoip.WithRootConfig(cfgSrv),
		geoip.WithOptionFactory(off),
	)
	assert.NoError(t, gs.ClearCache())
	_, err := gs.ConfigByScope(0, 0)
	assert.True(t, errors.IsEmpty(err), " Error: %+v", err)
}
