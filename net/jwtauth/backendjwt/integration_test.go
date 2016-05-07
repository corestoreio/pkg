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
	"fmt"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/net/jwtauth"
	"github.com/corestoreio/csfw/net/jwtauth/backendjwt"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func mustToPath(t *testing.T, f func(s scope.Scope, scopeID int64) (cfgpath.Path, error), s scope.Scope, scopeID int64) string {
	p, err := f(s, scopeID)
	if err != nil {
		t.Fatal(errors.PrintLoc(err))
	}
	return p.String()
}

func TestServiceWithBackend_HMACSHA_Website(t *testing.T) {

	cfgStruct, err := backendjwt.NewConfigStructure()
	if err != nil {
		t.Fatal(errors.PrintLoc(err))
	}
	pb := backendjwt.New(cfgStruct, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	jwts := jwtauth.MustNewService(
		jwtauth.WithBackend(backendjwt.BackendOptions(pb)),
	)

	pv := cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtSigningMethod.ToPath, scope.Default, 0): "ES384",
		mustToPath(t, pb.NetCtxjwtSigningMethod.ToPath, scope.Website, 1): "HS512",

		mustToPath(t, pb.NetCtxjwtEnableJTI.ToPath, scope.Default, 0): 0, // disabled
		mustToPath(t, pb.NetCtxjwtEnableJTI.ToPath, scope.Website, 1): 1, // enabled

		mustToPath(t, pb.NetCtxjwtExpiration.ToPath, scope.Default, 0): "2m",
		mustToPath(t, pb.NetCtxjwtExpiration.ToPath, scope.Website, 1): "5m1s",

		mustToPath(t, pb.NetCtxjwtHmacPassword.ToPath, scope.Default, 0): "pw1",
		mustToPath(t, pb.NetCtxjwtHmacPassword.ToPath, scope.Website, 1): "pw2",
	}
	sg := cfgmock.NewService(cfgmock.WithPV(pv)).NewScoped(1, 0) // only website scope supported

	scNew, err := jwts.ConfigByScopedGetter(sg)
	if err != nil {
		t.Fatal(errors.PrintLoc(err))
	}

	assert.True(t, scNew.EnableJTI)
	assert.Exactly(t, "5m1s", scNew.Expire.String())
	assert.Exactly(t, "HS512", scNew.SigningMethod.Alg())
	assert.False(t, scNew.Key.IsEmpty())
	assert.Nil(t, scNew.ErrorHandler)

	// test if cache returns the same scopedConfig
	scCached, err := jwts.ConfigByScopedGetter(sg)
	if err != nil {
		t.Fatal(errors.PrintLoc(err))
	}
	// reflect.DeepEqual returns here false because scCached was copied.
	assert.Exactly(t, fmt.Sprintf("%#v", scNew), fmt.Sprintf("%#v", scCached))
}

func TestServiceWithBackend_HMACSHA_Fallback(t *testing.T) {

	cfgStruct, err := backendjwt.NewConfigStructure()
	if err != nil {
		t.Fatal(errors.PrintLoc(err))
	}
	pb := backendjwt.New(cfgStruct, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	jwts := jwtauth.MustNewService(
		jwtauth.WithBackend(backendjwt.BackendOptions(pb)),
	)

	pv := cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtSigningMethod.ToPath, scope.Default, 0): "HS384",

		mustToPath(t, pb.NetCtxjwtEnableJTI.ToPath, scope.Default, 0): 0, // disabled

		mustToPath(t, pb.NetCtxjwtExpiration.ToPath, scope.Default, 0): "2m",

		mustToPath(t, pb.NetCtxjwtHmacPassword.ToPath, scope.Default, 0): "pw1",
	}

	sg := cfgmock.NewService(cfgmock.WithPV(pv)).NewScoped(1, 0) // 1 = website euro and 0 no store ID provided like in the middleware

	scNew, err := jwts.ConfigByScopedGetter(sg)
	if err != nil {
		t.Fatal(errors.PrintLoc(err))
	}

	assert.False(t, scNew.EnableJTI)
	assert.Exactly(t, "2m0s", scNew.Expire.String())
	assert.Exactly(t, "HS384", scNew.SigningMethod.Alg())
	assert.False(t, scNew.Key.IsEmpty())

	// test if cache returns the same scopedConfig
	scCached, err := jwts.ConfigByScopedGetter(sg)
	if err != nil {
		t.Fatal(errors.PrintLoc(err))
	}
	// reflect.DeepEqual returns here false because scCached was copied.
	assert.Exactly(t, fmt.Sprintf("%#v", scNew), fmt.Sprintf("%#v", scCached))
}

func getJwts(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) (jwts *jwtauth.Service, pb *backendjwt.Backend) {
	pb = backendjwt.New(cfgStruct, opts...)
	jwts = jwtauth.MustNewService(jwtauth.WithBackend(backendjwt.BackendOptions(pb)))
	return
}

func TestServiceWithBackend_UnknownSigningMethod(t *testing.T) {

	jwts, pb := getJwts(nil)

	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtSigningMethod.ToPath, scope.Default, 0): "HS4711",
	}))

	sc, err := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotImplemented(err), "Error: %s", err)
	assert.False(t, sc.IsValid())
}

func TestServiceWithBackend_InvalidExpiration(t *testing.T) {

	jwts, pb := getJwts(nil)

	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtExpiration.ToPath, scope.Default, 0): "Fail",
	}))

	sc, err := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
	assert.False(t, sc.IsValid())
}

func TestServiceWithBackend_InvalidJTI(t *testing.T) {

	jwts, pb := getJwts(nil)

	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtEnableJTI.ToPath, scope.Default, 0): []byte(`1`),
	}))

	sc, err := jwts.ConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
	assert.False(t, sc.IsValid())
}

func TestServiceWithBackend_RSAFail(t *testing.T) {

	jwts, pb := getJwts(nil, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtSigningMethod.ToPath, scope.Default, 0):  "RS256",
		mustToPath(t, pb.NetCtxjwtRSAKey.ToPath, scope.Default, 0):         []byte(`1`),
		mustToPath(t, pb.NetCtxjwtRSAKeyPassword.ToPath, scope.Default, 0): nil,
	}))

	sc, err := jwts.ConfigByScopedGetter(cr.NewScoped(1, 0))
	assert.True(t, errors.IsNotSupported(err))
	assert.False(t, sc.IsValid())
}

func TestServiceWithBackend_NilScopedGetter(t *testing.T) {

	jwts, _ := getJwts(nil)

	sc, err := jwts.ConfigByScopedGetter(nil) // don't do that in production !!!
	assert.NoError(t, err)

	assert.Exactly(t, scope.DefaultHash, sc.ScopeHash)
	assert.False(t, sc.Key.IsEmpty())
	assert.Exactly(t, jwtauth.DefaultExpire, sc.Expire)
	assert.Exactly(t, csjwt.HS256, sc.SigningMethod.Alg())
	assert.False(t, sc.EnableJTI)
	assert.Nil(t, sc.ErrorHandler)
	assert.NotNil(t, sc.KeyFunc)
}
