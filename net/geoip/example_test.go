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

package geoip_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/net/geoip/backendgeoip"
	"github.com/corestoreio/csfw/net/geoip/maxmindfile"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
)

func ExampleService_WithIsCountryAllowedByIP() {

	// The scope. Those two values are now hard coded because we cannot access
	// here the database to the website, store_group and store tables.
	const (
		websiteID = 1
		storeID   = 2
	)

	logBuf := new(log.MutexBuffer)

	cfgStruct, err := backendgeoip.NewConfigStructure()
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	backend := backendgeoip.New(cfgStruct)
	backend.Register(maxmindfile.NewOptionFactory(backend.MaxmindLocalFile))

	// This configuration says that any incoming request whose IP address does
	// not belong to the countries Germany (DE), Austria (AT) or Switzerland
	// (CH) gets redirected to the URL byebye.de.io with the redirect code 307.
	//
	// We're using here the cfgmock.NewService which is an in-memory
	// configuration service. Normally you would use the MySQL database or
	// consul or etcd.
	cfgSrv := cfgmock.NewService(cfgmock.PathValue{
		// @see structure.go why scope.Store and scope.Website can be used.
		backend.DataSource.MustFQ():                              `file`, // file triggers the lookup in
		backend.MaxmindLocalFile.MustFQ():                        filepath.Join("testdata", "GeoIP2-Country-Test.mmdb"),
		backend.AlternativeRedirect.MustFQStore(storeID):         `https://byebye.de.io`,
		backend.AlternativeRedirectCode.MustFQWebsite(websiteID): 307,
		backend.AllowedCountries.MustFQStore(storeID):            "DE,AT,CH",
	})

	geoSrv := geoip.MustNew(
		geoip.WithRootConfig(cfgSrv),
		geoip.WithDebugLog(logBuf),
		geoip.WithOptionFactory(backend.PrepareOptionFactory()),
		// Just for testing and in this example, we let the HTTP Handler panicking on any
		// error. You should not do that in production apps.
		geoip.WithServiceErrorHandler(mw.ErrorWithPanic),
	)

	// Set up the incoming request from the outside world. The scope in the
	// context gets set via scope.RunModeCalculater or via JSON web token or via
	// session.
	req := httptest.NewRequest("GET", "https://corestore.io/", nil)
	req = req.WithContext(scope.WithContext(req.Context(), websiteID, storeID))
	req.RemoteAddr = `2a02:d180::` // IP address range of Germany
	rec := httptest.NewRecorder()

	lastHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, ok := geoip.FromContextCountry(r.Context())
		if !ok {
			panic("A programmer made an error")
		}
		fmt.Fprintf(w, "Got allowed country DE:%s EN:%s", c.Country.Names["de"], c.Country.Names["en"])
	})
	geoSrv.WithIsCountryAllowedByIP(lastHandler).ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	// Change the request to an IP address outside the DACH region
	req.RemoteAddr = `2a02:d200::` // Finland
	rec = httptest.NewRecorder()
	// re-run the middleware and access is denied despite the scope in the context
	// is the same.
	geoSrv.WithIsCountryAllowedByIP(lastHandler).ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())
	// Output:
	// Got allowed country DE:Deutschland EN:Germany
	// <a href="https://byebye.de.io">Temporary Redirect</a>.
}
