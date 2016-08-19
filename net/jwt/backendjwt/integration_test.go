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

package backendjwt_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/jwt"
	"github.com/corestoreio/csfw/net/jwt/backendjwt"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestServiceWithBackend_HMACSHA_Website(t *testing.T) {
	cfgStruct, err := backendjwt.NewConfigStructure()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pb := backendjwt.New(cfgStruct, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	cfgSrv := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		pb.NetJwtSigningMethod.MustFQ(scope.Default, 0): "ES384",
		pb.NetJwtSigningMethod.MustFQ(scope.Website, 1): "HS512",

		pb.NetJwtSingleTokenUsage.MustFQ(scope.Default, 0): 0, // disabled
		pb.NetJwtSingleTokenUsage.MustFQ(scope.Website, 1): 1, // enabled

		pb.NetJwtDisabled.MustFQ(scope.Default, 0): 0, // disable: disabled 8-)
		pb.NetJwtDisabled.MustFQ(scope.Website, 1): 1, // disable: enabled 8-)

		pb.NetJwtExpiration.MustFQ(scope.Default, 0): "2m",
		pb.NetJwtExpiration.MustFQ(scope.Website, 1): "5m1s",

		pb.NetJwtSkew.MustFQ(scope.Default, 0): "4m",
		pb.NetJwtSkew.MustFQ(scope.Website, 1): "6m1s",

		pb.NetJwtHmacPassword.MustFQ(scope.Default, 0): "pw1",
		pb.NetJwtHmacPassword.MustFQ(scope.Website, 1): "pw2",
	}))

	jwts := jwt.MustNew(
		jwt.WithOptionFactory(backendjwt.PrepareOptions(pb)),
	)

	sg := cfgSrv.NewScoped(1, 0) // only website scope supported

	scNew := jwts.ConfigByScopedGetter(sg)
	if err := scNew.IsValid(); err != nil {
		t.Fatalf("%+v", err)
	}

	assert.True(t, scNew.SingleTokenUsage)
	assert.True(t, scNew.Disabled)
	assert.Exactly(t, "5m1s", scNew.Expire.String())
	assert.Exactly(t, "6m1s", scNew.Skew.String())
	assert.Exactly(t, "HS512", scNew.SigningMethod.Alg())
	assert.False(t, scNew.Key.IsEmpty())
	assert.NotNil(t, scNew.ErrorHandler)

	// test if cache returns the same scopedConfig
	scCached := jwts.ConfigByScopedGetter(sg)
	if err := scCached.IsValid(); err != nil {
		t.Fatalf("%+v", err)
	}
	// reflect.DeepEqual returns here false because scCached was copied.
	assert.Exactly(t, fmt.Sprintf("%#v", scNew), fmt.Sprintf("%#v", scCached))
}

func TestServiceWithBackend_HMACSHA_Fallback(t *testing.T) {

	cfgStruct, err := backendjwt.NewConfigStructure()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pb := backendjwt.New(cfgStruct, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	cfgSrv := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		pb.NetJwtSigningMethod.MustFQ(scope.Default, 0):    "HS384",
		pb.NetJwtSingleTokenUsage.MustFQ(scope.Default, 0): 0, // disabled
		pb.NetJwtDisabled.MustFQ(scope.Default, 0):         1, // disabled active
		pb.NetJwtExpiration.MustFQ(scope.Default, 0):       "2m",
		pb.NetJwtSkew.MustFQ(scope.Default, 0):             "3m",
		pb.NetJwtHmacPassword.MustFQ(scope.Default, 0):     "pw1",
	}))

	jwts := jwt.MustNew(
		jwt.WithOptionFactory(backendjwt.PrepareOptions(pb)),
	)

	sg := cfgSrv.NewScoped(1, 0) // 1 = website euro and 0 no store ID provided like in the middleware

	scNew := jwts.ConfigByScopedGetter(sg)
	if err := scNew.IsValid(); err != nil {
		t.Fatalf("%+v", err)
	}

	assert.False(t, scNew.SingleTokenUsage)
	assert.True(t, scNew.Disabled)
	assert.Exactly(t, "2m0s", scNew.Expire.String())
	assert.Exactly(t, "3m0s", scNew.Skew.String())
	assert.Exactly(t, "HS384", scNew.SigningMethod.Alg())
	assert.False(t, scNew.Key.IsEmpty())

	// test if cache returns the same scopedConfig
	scCached := jwts.ConfigByScopedGetter(sg)
	if err := scCached.IsValid(); err != nil {
		t.Fatalf("%+v", err)
	}
	// reflect.DeepEqual returns here false because scCached was copied.
	assert.Exactly(t, fmt.Sprintf("%#v", scNew), fmt.Sprintf("%#v", scCached))
}

func getJwts(opts ...cfgmodel.Option) (jwts *jwt.Service, pb *backendjwt.Configuration) {
	cfgStruct, err := backendjwt.NewConfigStructure()
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	pb = backendjwt.New(cfgStruct, opts...)
	jwts = jwt.MustNew(jwt.WithOptionFactory(backendjwt.PrepareOptions(pb)))
	return
}

func TestServiceWithBackend_MissingSectionSlice(t *testing.T) {

	pb := backendjwt.New(nil)
	jwts := jwt.MustNew(jwt.WithOptionFactory(backendjwt.PrepareOptions(pb)))

	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		pb.NetJwtSigningMethod.MustFQ(scope.Default, 0): "HS4711",
	}))

	sc := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotFound(sc.IsValid()), "Error: %+v", sc.IsValid())
}

func TestServiceWithBackend_UnknownSigningMethod(t *testing.T) {

	jwts, pb := getJwts()

	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		pb.NetJwtSigningMethod.MustFQ(scope.Default, 0): "HS4711",
	}))

	sc := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotImplemented(sc.IsValid()), "Error: %+v", sc.IsValid())
}

func TestServiceWithBackend_InvalidExpiration(t *testing.T) {

	jwts, pb := getJwts()

	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		pb.NetJwtExpiration.MustFQ(scope.Default, 0): "Fail",
	}))

	sc := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	err := sc.IsValid()
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
}

func TestServiceWithBackend_InvalidSkew(t *testing.T) {

	jwts, pb := getJwts()

	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		pb.NetJwtSkew.MustFQ(scope.Default, 0): "Fail171",
	}))

	sc := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	err := sc.IsValid()
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
}

func TestServiceWithBackend_InvalidJTI(t *testing.T) {

	jwts, pb := getJwts()

	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		pb.NetJwtSingleTokenUsage.MustFQ(scope.Default, 0): []byte(`1`),
	}))

	sc := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	err := sc.IsValid()
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
}

func TestServiceWithBackend_RSAFail(t *testing.T) {

	jwts, pb := getJwts(cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		pb.NetJwtSigningMethod.MustFQ(scope.Default, 0):  "RS256",
		pb.NetJwtRSAKey.MustFQ(scope.Default, 0):         []byte(`1`),
		pb.NetJwtRSAKeyPassword.MustFQ(scope.Default, 0): nil,
	}))

	sc := jwts.ConfigByScopedGetter(cr.NewScoped(1, 0))
	err := sc.IsValid()
	assert.True(t, errors.IsNotSupported(err))
}

// TestServiceWithBackend_WithRunMode_Valid_Request tests that a request
// contains a valid token, loads atomically the backend configuration and
// switches the stores
func TestServiceWithBackend_WithRunMode_Valid_Request(t *testing.T) {

	// setup overall configuration structure
	cfgStruct, err := backendjwt.NewConfigStructure()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// use that configuration structure to apply it to the configuration models.
	pb := backendjwt.New(cfgStruct, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	// create a configuration for websiteID 1. this configuration resides usually in
	// the MySQL core_config_data table.
	cfgSrv := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		pb.NetJwtSigningMethod.MustFQ(scope.Website, 1):    "HS512",
		pb.NetJwtSingleTokenUsage.MustFQ(scope.Website, 1): 1, // enabled
		pb.NetJwtDisabled.MustFQ(scope.Website, 1):         0, // JWT parsing enabled
		pb.NetJwtExpiration.MustFQ(scope.Website, 1):       "5m1s",
		pb.NetJwtSkew.MustFQ(scope.Website, 1):             "6m1s",
		pb.NetJwtHmacPassword.MustFQ(scope.Website, 1):     "pw2",
	}))

	logBuf := new(bytes.Buffer)
	jwts := jwt.MustNew(
		jwt.WithRootConfig(cfgSrv),
		jwt.WithDebugLog(logBuf),
		jwt.WithOptionFactory(backendjwt.PrepareOptions(pb)),
	)

	// our token will be crafted to contain the DE store so the JWT middleware
	// must change the store to Germany, the store code with wich we've started
	// was Austria AT.

	_ = jwts.ConfigByScope(1, 0) // init config, triggers first Log entry: jwt.Service.ConfigByScopedGetter.Inflight.Do

	rawToken := func() []byte {
		stClaim := jwtclaim.NewStore()
		stClaim.Store = "de"
		stClaim.UserID = "hans_wurst"
		newToken, err := jwts.NewToken(scope.Website, 1, stClaim)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		return newToken.Raw
	}()

	if err := jwts.ClearCache(); err != nil {
		// reset the cache to trigger loading the option factory
		t.Fatalf("%+v", err)
	}

	// craft the request which contains the configuration based on the incoming scope
	req := func() *http.Request {
		req := httptest.NewRequest("GET", "http://corestore.io", nil)
		req.Header.Set("X-Cluster-Client-Ip", "2a02:d180::") // Germany
		jwt.SetHeaderAuthorization(req, rawToken)
		return req
	}()

	// food for the race detector
	// the very first request triggers the 2nd log entry: jwt.Service.ConfigByScopedGetter.Inflight.Do
	hpu := cstesting.NewHTTPParallelUsers(4, 10, 100, time.Microsecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		if have, want := rec.Code, 200; have != want {
			t.Errorf("Response Code wrong. Have: %v Want: %v", have, want)
		}
	}
	hpu.ServeHTTP(
		req,
		jwts.WithRunMode(
			scope.RunMode{Mode: scope.NewHash(scope.Website, 1)}, // use euro wesite with default store AT.
			storemock.NewEurozzyService(cfgSrv),
		)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tk, ok := jwt.FromContext(r.Context())
			if !ok {
				panic("Token not found in context")
			}
			if !tk.Valid {
				panic(fmt.Sprintf("Token not valid: %+v", tk))
			}

			// here we must have the new request store and not anymore the Austrian store.
			websiteID, storeID, ok := scope.FromContext(r.Context())
			if !ok {
				panic("Scope not found in context")
			}
			assert.Exactly(t, int64(1), websiteID, "websiteID")
			assert.Exactly(t, int64(1), storeID, "storeID")
		})),
	)

	// println(logBuf.String(), "\n\n")

	containTests := []struct {
		check string
		want  int
	}{
		{`jwt.Service.ConfigByScopedGetter.Inflight.Do`, 2},
		{`Service.WithInitTokenAndStore.Disabled`, 0},
		{`ScopeOptionFromClaim.StoreServiceIsNil`, 0},
		{`jwt.Service.WithRunMode.NextHandler.WithCode`, 40},
		{`jwt.Service.ConfigByScopedGetter.optionFactoryFunc.nil`, 0},
	}
	for _, test := range containTests {
		if have, want := strings.Count(logBuf.String(), test.check), test.want; have != want {
			t.Errorf("%s: Have: %v Want: %v", test.check, have, want)
		}
	}
	if have, want := strings.Count(logBuf.String(), `jwt.Service.ConfigByScopedGetter.IsValid`), 36; have < want {
		t.Errorf("jwt.Service.ConfigByScopedGetter.IsValid: Have: %v Want: %v", have, want)
	}
}
