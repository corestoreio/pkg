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
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/net/geoip/backendgeoip"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cstesting"
)

func BenchmarkWithAlternativeRedirect(b *testing.B) {
	cfgSrv := cfgmock.NewService(cfgmock.PathValue{
		// @see structure.go why scope.Store and scope.Website can be used.
		mustToPath(b, backend.NetGeoipAlternativeRedirect.ToPath, scope.Store, 2):       `https://byebye.de.io`,
		mustToPath(b, backend.NetGeoipAlternativeRedirectCode.ToPath, scope.Website, 1): 307,
		mustToPath(b, backend.NetGeoipAllowedCountries.ToPath, scope.Store, 2):          "AT,CH",
		mustToPath(b, backend.NetGeoipMaxmindLocalFile.ToPath, scope.Default, 0):        filepath.Join("..", "testdata", "GeoIP2-Country-Test.mmdb"),
	})
	b.Run("LocalFile_NoCache", benchmarkWithAlternativeRedirect(cfgSrv))

	cfgSrv = cfgmock.NewService(cfgmock.PathValue{
		// @see structure.go why scope.Store and scope.Website can be used.
		mustToPath(b, backend.NetGeoipAlternativeRedirect.ToPath, scope.Store, 2):        `https://byebye.de.io`,
		mustToPath(b, backend.NetGeoipAlternativeRedirectCode.ToPath, scope.Website, 1):  307,
		mustToPath(b, backend.NetGeoipAllowedCountries.ToPath, scope.Store, 2):           "AT,CH",
		mustToPath(b, backend.NetGeoipMaxmindWebserviceUserID.ToPath, scope.Default, 0):  "LiesschenMueller",
		mustToPath(b, backend.NetGeoipMaxmindWebserviceLicense.ToPath, scope.Default, 0): "8x4",
		mustToPath(b, backend.NetGeoipMaxmindWebserviceTimeout.ToPath, scope.Default, 0): "3s",
	})
	// to fix the speed here ... BigCache_Gob must be optimized
	b.Run("Webservice_BigCache_Gob", benchmarkWithAlternativeRedirect(cfgSrv))
}

func benchmarkWithAlternativeRedirect(cfgSrv *cfgmock.Service) func(b *testing.B) {
	return func(b *testing.B) {
		cfgStruct, err := backendgeoip.NewConfigStructure()
		if err != nil {
			b.Fatal(err)
		}
		be := backendgeoip.New(cfgStruct)
		be.WebServiceClient = &http.Client{
			Transport: cstesting.NewHTTPTrip(200, `{ "continent": { "code": "EU", "geoname_id": 6255148, "names": { "de": "Europa", "en": "Europe", "ru": "Европа", "zh-CN": "欧洲" } }, "country": { "geoname_id": 2921044, "iso_code": "DE", "names": { "de": "Deutschland", "en": "Germany", "es": "Alemania", "fr": "Allemagne", "ja": "ドイツ連邦共和国", "pt-BR": "Alemanha", "ru": "Германия", "zh-CN": "德国" } }, "maxmind": { "queries_remaining": 54321 } }`, nil),
		}
		scpFnc := backendgeoip.PrepareOptions(be)
		geoSrv := geoip.MustNew(geoip.WithOptionFactory(scpFnc))

		// Germany is not allowed and must be redirected to https://byebye.de.io with code 307
		req := func() *http.Request {
			req := httptest.NewRequest("GET", "http://corestore.io", nil)
			req.Header.Set("X-Cluster-Client-Ip", "2a02:d180::")
			// Website ID 1 == euro / Store ID == 2 Austria
			return req.WithContext(scope.WithContext(req.Context(), 1, 2))
		}()

		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				rec := httptest.NewRecorder()
				geoSrv.WithIsCountryAllowedByIP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					c, ok := geoip.FromContextCountry(r.Context())
					if c != nil {
						b.Fatalf("Country must be nil, but is %#v", c)
					}
					if !ok {
						b.Fatal("Failed to find a country pointer in the context")
					}

					panic("Should not be called")

				})).ServeHTTP(rec, req)

				if have, want := rec.Header().Get("Location"), `https://byebye.de.io`; have != want {
					b.Errorf("Have %q Want %q", have, want)
				}
				if have, want := rec.Code, 307; have != want {
					b.Errorf("Have %q Want %q", have, want)
				}
			}
		})
	}
}
