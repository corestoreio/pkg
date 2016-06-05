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
	"testing"

	"path/filepath"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/net/geoip/backendgeoip"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
)

// TODO: fix bug of too many allocs. they are caused by the 3rd party packages for reading MM files.
// BenchmarkWithAlternativeRedirect-4   	  100000	     13264 ns/op	   17207 B/op	     133 allocs/op
func BenchmarkWithAlternativeRedirect_Database_NoCache(b *testing.B) {
	cfgSrv := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		// @see structure.go why scope.Store and scope.Website can be used.
		mustToPath(b, backend.NetGeoipAlternativeRedirect.ToPath, scope.Store, 2):       `https://byebye.de.io`,
		mustToPath(b, backend.NetGeoipAlternativeRedirectCode.ToPath, scope.Website, 1): 307,
		mustToPath(b, backend.NetGeoipAllowedCountries.ToPath, scope.Store, 2):          "AT,CH",
		mustToPath(b, backend.NetGeoipMaxmindLocalFile.ToPath, scope.Default, 0):        filepath.Join("..", "testdata", "GeoIP2-Country-Test.mmdb"),
	}))
	benchmarkWithAlternativeRedirect(b, cfgSrv)
}

func BenchmarkWithAlternativeRedirect_Webservice_BigCache_Gob(b *testing.B) {

	trip := cstesting.NewHttpTrip(200, `{ "continent": { "code": "EU", "geoname_id": 6255148, "names": { "de": "Europa", "en": "Europe", "es": "Europa", "fr": "Europe", "ja": "ヨーロッパ", "pt-BR": "Europa", "ru": "Европа", "zh-CN": "欧洲" } }, "country": { "geoname_id": 2921044, "iso_code": "DE", "names": { "de": "Deutschland", "en": "Germany", "es": "Alemania", "fr": "Allemagne", "ja": "ドイツ連邦共和国", "pt-BR": "Alemanha", "ru": "Германия", "zh-CN": "德国" } }, "registered_country": { "geoname_id": 2921044, "iso_code": "DE", "names": { "de": "Deutschland", "en": "Germany", "es": "Alemania", "fr": "Allemagne", "ja": "ドイツ連邦共和国", "pt-BR": "Alemanha", "ru": "Германия", "zh-CN": "德国" } }, "traits": { "autonomous_system_number": 1239, "autonomous_system_organization": "Linkem IR WiMax Network", "domain": "example.com", "is_anonymous_proxy": true, "is_satellite_provider": true, "isp": "Linkem spa", "ip_address": "1.2.3.4", "organization": "Linkem IR WiMax Network", "user_type": "traveler" }, "maxmind": { "queries_remaining": 54321 } }`, nil)
	backend.WebServiceClient = &http.Client{
		Transport: trip,
	}

	cfgSrv := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		// @see structure.go why scope.Store and scope.Website can be used.
		mustToPath(b, backend.NetGeoipAlternativeRedirect.ToPath, scope.Store, 2):        `https://byebye.de.io`,
		mustToPath(b, backend.NetGeoipAlternativeRedirectCode.ToPath, scope.Website, 1):  307,
		mustToPath(b, backend.NetGeoipAllowedCountries.ToPath, scope.Store, 2):           "AT,CH",
		mustToPath(b, backend.NetGeoipMaxmindWebserviceUserID.ToPath, scope.Default, 0):  "LiesschenMueller",
		mustToPath(b, backend.NetGeoipMaxmindWebserviceLicense.ToPath, scope.Default, 0): "8x4",
		mustToPath(b, backend.NetGeoipMaxmindWebserviceTimeout.ToPath, scope.Default, 0): "3s",
	}))
	benchmarkWithAlternativeRedirect(b, cfgSrv)
}

func benchmarkWithAlternativeRedirect(b *testing.B, cfgSrv *cfgmock.Service) {

	geoSrv := geoip.MustNew(geoip.WithOptionFactory(backendgeoip.Default()))

	// Germany is not allowed and must be redirected to https://byebye.de.io with code 307
	req := func() *http.Request {
		o, err := scope.SetByCode(scope.Website, "euro")
		if err != nil {
			b.Fatal(err)
		}
		storeSrv := storemock.NewEurozzyService(o)
		req, _ := http.NewRequest("GET", "http://corestore.io", nil)
		req.RemoteAddr = "2a02:d180::"
		atSt, err := storeSrv.Store(scope.MockID(2)) // Austria Store
		if err != nil {
			b.Fatal(errors.PrintLoc(err))
		}
		atSt.Config = cfgSrv.NewScoped(1, 2) // Website ID 1 == euro / Store ID == 2 Austria

		return req.WithContext(store.WithContextRequestedStore(req.Context(), atSt))
	}()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rec := httptest.NewRecorder()
			geoSrv.WithIsCountryAllowedByIP()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				c, err := geoip.FromContextCountry(r.Context())
				if c != nil {
					b.Fatalf("Country must be nil, but is %#v", c)
				}
				if err != nil {
					b.Fatal(errors.PrintLoc(err))
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
