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
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestConfiguration_HierarchicalConfig(t *testing.T) {

	// allow all scopes for testing
	//backend.Skew.Field.Scopes = scope.PermStoreReverse

	scpCfgSrv := cfgmock.NewService(cfgmock.PathValue{
		backend.SingleTokenUsage.MustFQ():     `1`,
		backend.Expiration.MustFQWebsite(3):   `66s`,
		backend.Skew.MustFQ():                 `33s`,
		backend.HmacPassword.MustFQWebsite(3): `This is a secure encrypted password.`,
	}).NewScoped(3, 0)

	srv := jwt.MustNew(
		jwt.WithOptionFactory(backendjwt.PrepareOptions(backend)),
	)
	scpCfg, err := srv.ConfigByScopedGetter(scpCfgSrv)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.True(t, scpCfg.SingleTokenUsage)
	assert.Exactly(t, time.Second*33, scpCfg.Skew)
	assert.Exactly(t, time.Second*66, scpCfg.Expire)
}

func TestServiceWithBackend_HMACSHA_Website(t *testing.T) {
	cfgStruct, err := backendjwt.NewConfigStructure()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pb := backendjwt.New(cfgStruct, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	cfgSrv := cfgmock.NewService(cfgmock.PathValue{
		pb.SigningMethod.MustFQ():         "ES384",
		pb.SigningMethod.MustFQWebsite(1): "HS512",

		pb.SingleTokenUsage.MustFQ():         0, // disabled
		pb.SingleTokenUsage.MustFQWebsite(1): 1, // enabled

		//pb.Disabled.MustFQ():         0, // disable: disabled 8-)
		//pb.Disabled.MustFQWebsite(1): 1, // disable: enabled 8-)

		pb.Expiration.MustFQ():         "2m",
		pb.Expiration.MustFQWebsite(1): "5m1s",

		pb.Skew.MustFQ():         "4m",
		pb.Skew.MustFQWebsite(1): "6m1s",

		pb.HmacPassword.MustFQ():         "pw1",
		pb.HmacPassword.MustFQWebsite(1): "pw2",
	})

	jwts := jwt.MustNew(jwt.WithOptionFactory(backendjwt.PrepareOptions(pb)))

	sg := cfgSrv.NewScoped(1, 0) // only website scope supported

	scNew, err := jwts.ConfigByScopedGetter(sg)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	assert.True(t, scNew.SingleTokenUsage, "SingleTokenUsage")
	assert.False(t, scNew.Disabled, "Disabled")
	assert.Exactly(t, "5m1s", scNew.Expire.String(), "Expire")
	assert.Exactly(t, "6m1s", scNew.Skew.String(), "Skew")
	assert.Exactly(t, "HS512", scNew.SigningMethod.Alg(), "SigningMethod")
	assert.False(t, scNew.Key.IsEmpty())
	assert.NotNil(t, scNew.ErrorHandler)

	// test if cache returns the same scopedConfig
	scCached, err := jwts.ConfigByScopedGetter(sg)
	if err != nil {
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

	cfgSrv := cfgmock.NewService(cfgmock.PathValue{
		pb.SigningMethod.MustFQ():    "HS384",
		pb.SingleTokenUsage.MustFQ(): 0, // disabled
		pb.Expiration.MustFQ():       "2m",
		pb.Skew.MustFQ():             "3m",
		pb.HmacPassword.MustFQ():     "pw1",
	})

	jwts := jwt.MustNew(
		jwt.WithOptionFactory(backendjwt.PrepareOptions(pb)),
	)

	sg := cfgSrv.NewScoped(1, 0) // 1 = website euro and 0 no store ID provided like in the middleware

	scNew, err := jwts.ConfigByScopedGetter(sg)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	assert.False(t, scNew.SingleTokenUsage)
	assert.False(t, scNew.Disabled)
	assert.Exactly(t, "2m0s", scNew.Expire.String())
	assert.Exactly(t, "3m0s", scNew.Skew.String())
	assert.Exactly(t, "HS384", scNew.SigningMethod.Alg())
	assert.False(t, scNew.Key.IsEmpty())

	// test if cache returns the same scopedConfig
	scCached, err := jwts.ConfigByScopedGetter(sg)
	if err != nil {
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

	cr := cfgmock.NewService(cfgmock.PathValue{
		pb.SigningMethod.MustFQ(): "HS4711",
	})

	_, err := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotFound(err), "Error: %+v", err)
}

func TestServiceWithBackend_UnknownSigningMethod(t *testing.T) {

	jwts, pb := getJwts()

	cr := cfgmock.NewService(cfgmock.PathValue{
		pb.SigningMethod.MustFQ(): "HS4711",
	})

	_, err := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotImplemented(err), "Error: %+v", err)
}

func TestServiceWithBackend_InvalidExpiration(t *testing.T) {

	jwts, pb := getJwts()

	cr := cfgmock.NewService(cfgmock.PathValue{
		pb.Expiration.MustFQ(): "Fail",
	})

	_, err := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
}

func TestServiceWithBackend_InvalidSkew(t *testing.T) {

	jwts, pb := getJwts()

	cr := cfgmock.NewService(cfgmock.PathValue{
		pb.Skew.MustFQ(): "Fail171",
	})

	_, err := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
}

func TestServiceWithBackend_InvalidJTI(t *testing.T) {

	jwts, pb := getJwts()

	cr := cfgmock.NewService(cfgmock.PathValue{
		pb.SingleTokenUsage.MustFQ(): []byte(`1`),
	})

	_, err := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
}

func TestServiceWithBackend_RSAFail(t *testing.T) {

	jwts, pb := getJwts(cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	cr := cfgmock.NewService(cfgmock.PathValue{
		pb.SigningMethod.MustFQ():  "RS256",
		pb.RSAKey.MustFQ():         []byte(`1`),
		pb.RSAKeyPassword.MustFQ(): nil,
	})

	_, err := jwts.ConfigByScopedGetter(cr.NewScoped(1, 0))
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
	cfgSrv := cfgmock.NewService(cfgmock.PathValue{
		pb.SigningMethod.MustFQWebsite(1):    "HS512",
		pb.SingleTokenUsage.MustFQWebsite(1): 1, // enabled
		pb.Disabled.MustFQWebsite(1):         0, // JWT parsing enabled
		pb.Expiration.MustFQWebsite(1):       "5m1s",
		pb.Skew.MustFQWebsite(1):             "6m1s",
		pb.HmacPassword.MustFQWebsite(1):     "pw2",
	})

	logBuf := new(bytes.Buffer)
	jwts := jwt.MustNew(
		jwt.WithRootConfig(cfgSrv),
		jwt.WithDebugLog(logBuf),
		jwt.WithOptionFactory(backendjwt.PrepareOptions(pb)),
		jwt.WithServiceErrorHandler(mw.ErrorWithPanic),
		jwt.WithErrorHandler(scope.Website.Pack(1), mw.ErrorWithPanic),
		jwt.WithTriggerOptionFactories(scope.Website.Pack(1), true), // because we load more from the OptionFactories
	)

	// our token will be crafted to contain the DE store so the JWT middleware
	// must change the store to Germany, the store code with wich we've started
	// was Austria AT.

	xcfg, err := jwts.ConfigByScope(1, 0) // init config, triggers first Log entry: jwt.Service.ConfigByScopedGetter.Inflight.Do
	if err != nil && !errors.IsTemporary(err) {
		t.Fatalf("Expecting temporary error because config marked as partially loaded: %+v\n", err)
	}
	if xcfg.SigningMethod.Alg() != csjwt.HS512 {
		t.Fatalf("Expecting SigningMethod to be Have: %q Want %q Scope %s", xcfg.SigningMethod.Alg(), csjwt.HS512, xcfg.ScopeID)
	}

	rawToken := func() []byte {
		stClaim := jwtclaim.NewStore()
		stClaim.Store = "de"
		stClaim.UserID = "hans_wurst"
		newToken, err := jwts.NewToken(scope.Website.Pack(1), stClaim)
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
	hpu := cstesting.NewHTTPParallelUsers(4, 10, 100, time.Millisecond)
	// setting it to time.Microsecond above causes loading of the wrong config under high load for the initial
	// cache filling.
	// test with $ go test -race -run=TestServiceWithBackend_WithRunMode_Valid_Request -count=10 .
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		if have, want := rec.Code, 200; have != want {
			t.Errorf("Response Code wrong. Have: %v Want: %v\n\n%s", have, want, rec.Body)
		}
	}
	hpu.ServeHTTP(
		req,
		jwts.WithRunMode(
			scope.MakeTypeID(scope.Website, 1), // use euro website with default store AT.
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
