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

func ExampleNewOptionFactoryGeoSourceFile() {

	//logBuf := new(log.MutexBuffer)
	//
	//cfgStruct, err := backendgeoip.NewConfigStructure()
	//if err != nil {
	//	panic(fmt.Sprintf("%+v", err))
	//}
	//be := backendgeoip.New(cfgStruct)
	//be.Register(backendgeoip.NewOptionFactoryFile(be.MaxmindLocalFile))
	//
	//// This configuration says that any incoming request whose IP address does
	//// not belong the countries Austria (AT) or Switzerland (CH) gets redirected
	//// to the URL byebye.de.io with the redirect code 307.
	////
	//// We're using here the cfgmock.NewService which is an in-memory
	//// configuration service. Normally you would use the MySQL database or
	//// consul or etcd.
	//cfgSrv := cfgmock.NewService(cfgmock.PathValue{
	//	// @see structure.go why scope.Store and scope.Website can be used.
	//	backend.DataSource.MustFQ():                      `file`, // file triggers the lookup in
	//	backend.MaxmindLocalFile.MustFQ():                filepath.Join("..", "testdata", "GeoIP2-Country-Test.mmdb"),
	//	backend.AlternativeRedirect.MustFQStore(2):       `https://byebye.de.io`,
	//	backend.AlternativeRedirectCode.MustFQWebsite(1): 307,
	//	backend.AllowedCountries.MustFQStore(2):          "AT,CH",
	//})
	//
	//scpFnc := backendgeoip.PrepareOptions(be)
	//geoSrv := geoip.MustNew(
	//	geoip.WithRootConfig(cfgSrv),
	//	geoip.WithDebugLog(logBuf),
	//	geoip.WithOptionFactory(scpFnc),
	//	// Just for testing and in this example, we let the HTTP Handler panicking on any
	//	// error. You should not do that in production apps.
	//	geoip.WithServiceErrorHandler(mw.ErrorWithPanic),
	//	geoip.WithErrorHandler(scope.DefaultTypeID, mw.ErrorWithPanic),
	//)
	//
	//geoSrv.ServeHTTP()

}
