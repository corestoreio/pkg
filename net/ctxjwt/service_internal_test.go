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

package ctxjwt

import (
	"testing"

	"fmt"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/stretchr/testify/assert"
)

func mustToPath(t *testing.T, f func(s scope.Scope, scopeID int64) (cfgpath.Path, error), s scope.Scope, scopeID int64) string {
	p, err := f(s, scopeID)
	if err != nil {
		t.Fatal(err)
	}
	return p.String()
}

func TestServiceWithBackend_NoBackend(t *testing.T) {
	t.Parallel()

	jwts := MustNewService()
	// delete() is a hack for testing to remove the default setting
	delete(jwts.scopeCache, scope.NewHash(scope.DefaultID, 0))

	cr := cfgmock.NewService()
	sc, err := jwts.getConfigByScopedGetter(cr.NewScoped(0, 0))
	assert.EqualError(t, err, "[ctxjwt] Cannot find JWT configuration for Scope(Default) ID(0)")
	assert.Exactly(t, scopedConfig{}, sc)
}

func TestServiceWithBackend_DefaultConfig(t *testing.T) {

	jwts := MustNewService()

	cr := cfgmock.NewService()
	sc, err := jwts.getConfigByScopedGetter(cr.NewScoped(0, 0))
	assert.NoError(t, err)
	dsc, err := defaultScopedConfig()
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, csjwt.HS256, sc.signingMethod.Alg())
	assert.Exactly(t, dsc.Key.Algorithm(), sc.Key.Algorithm())

	assert.True(t, dsc.errorHandler == nil)
	assert.True(t, sc.errorHandler == nil)
	assert.True(t, jwts.DefaultErrorHandler != nil)
	assert.Exactly(t, DefaultExpire, dsc.expire)
	assert.False(t, dsc.Key.IsEmpty())
	assert.False(t, sc.Key.IsEmpty())
}

func TestServiceWithBackend_HMACSHA_Website(t *testing.T) {
	t.Parallel()
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		t.Fatal(err)
	}
	pb := NewBackend(cfgStruct, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	jwts := MustNewService(
		WithBackend(pb),
	)
	pv := cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtSigningMethod.ToPath, scope.DefaultID, 0): "ES384",
		mustToPath(t, pb.NetCtxjwtSigningMethod.ToPath, scope.WebsiteID, 1): "HS512",

		mustToPath(t, pb.NetCtxjwtEnableJTI.ToPath, scope.DefaultID, 0): 0, // disabled
		mustToPath(t, pb.NetCtxjwtEnableJTI.ToPath, scope.WebsiteID, 1): 1, // enabled

		mustToPath(t, pb.NetCtxjwtExpiration.ToPath, scope.DefaultID, 0): "2m",
		mustToPath(t, pb.NetCtxjwtExpiration.ToPath, scope.WebsiteID, 1): "5m1s",

		mustToPath(t, pb.NetCtxjwtHmacPassword.ToPath, scope.DefaultID, 0): "pw1",
		mustToPath(t, pb.NetCtxjwtHmacPassword.ToPath, scope.WebsiteID, 1): "pw2",
	}
	sg := cfgmock.NewService(cfgmock.WithPV(pv)).NewScoped(1, 1)

	scNew, err := jwts.getConfigByScopedGetter(sg)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, scNew.enableJTI)
	assert.Exactly(t, "5m1s", scNew.expire.String())
	assert.Exactly(t, "HS512", scNew.signingMethod.Alg())
	assert.False(t, scNew.Key.IsEmpty())
	assert.Nil(t, scNew.errorHandler)
	assert.NotNil(t, jwts.DefaultErrorHandler)

	// test if cache returns the same scopedConfig
	scCached, err := jwts.getConfigByScopedGetter(sg)
	if err != nil {
		t.Fatal(err)
	}
	// reflect.DeepEqual returns here false because scCached was copied.
	assert.Exactly(t, fmt.Sprintf("%#v", scNew), fmt.Sprintf("%#v", scCached))
}

func TestServiceWithBackend_HMACSHA_Fallback(t *testing.T) {
	t.Parallel()
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		t.Fatal(err)
	}
	pb := NewBackend(cfgStruct, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))

	jwts := MustNewService(
		WithBackend(pb),
	)
	pv := cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtSigningMethod.ToPath, scope.DefaultID, 0): "HS384",

		mustToPath(t, pb.NetCtxjwtEnableJTI.ToPath, scope.DefaultID, 0): 0, // disabled

		mustToPath(t, pb.NetCtxjwtExpiration.ToPath, scope.DefaultID, 0): "2m",

		mustToPath(t, pb.NetCtxjwtHmacPassword.ToPath, scope.DefaultID, 0): "pw1",
	}
	sg := cfgmock.NewService(cfgmock.WithPV(pv)).NewScoped(1, 1)

	scNew, err := jwts.getConfigByScopedGetter(sg)
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, scNew.enableJTI)
	assert.Exactly(t, "2m0s", scNew.expire.String())
	assert.Exactly(t, "HS384", scNew.signingMethod.Alg())
	assert.False(t, scNew.Key.IsEmpty())

	// test if cache returns the same scopedConfig
	scCached, err := jwts.getConfigByScopedGetter(sg)
	if err != nil {
		t.Fatal(err)
	}
	// reflect.DeepEqual returns here false because scCached was copied.
	assert.Exactly(t, fmt.Sprintf("%#v", scNew), fmt.Sprintf("%#v", scCached))

}

func TestServiceWithBackend_UnknownSigningMethod(t *testing.T) {
	t.Parallel()

	pb := NewBackend(nil)
	jwts := MustNewService(WithBackend(pb))
	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtSigningMethod.ToPath, scope.DefaultID, 0): "HS4711",
	}))

	sc, err := jwts.getConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.EqualError(t, err, "[ctxjwt] ConfigSigningMethod: Unknown algorithm HS4711")
	assert.Exactly(t, scopedConfig{}, sc)
}

func TestServiceWithBackend_InvalidExpiration(t *testing.T) {
	t.Parallel()

	pb := NewBackend(nil)
	jwts := MustNewService(WithBackend(pb))
	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtExpiration.ToPath, scope.DefaultID, 0): "Fail",
	}))

	sc, err := jwts.getConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.EqualError(t, err, "time: invalid duration Fail")
	assert.Exactly(t, scopedConfig{}, sc)
}

func TestServiceWithBackend_InvalidJTI(t *testing.T) {
	t.Parallel()

	pb := NewBackend(nil)
	jwts := MustNewService(WithBackend(pb))
	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtEnableJTI.ToPath, scope.DefaultID, 0): []byte(`1`),
	}))

	sc, err := jwts.getConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.EqualError(t, err, "Route net/ctxjwt/enable_jti: Unable to cast []byte{0x31} to bool")
	assert.Exactly(t, scopedConfig{}, sc)
}

func TestServiceWithBackend_RSAFail(t *testing.T) {
	t.Parallel()

	pb := NewBackend(nil, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))
	jwts := MustNewService(WithBackend(pb))
	cr := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		mustToPath(t, pb.NetCtxjwtSigningMethod.ToPath, scope.DefaultID, 0):  "RS256",
		mustToPath(t, pb.NetCtxjwtRSAKey.ToPath, scope.DefaultID, 0):         []byte(`1`),
		mustToPath(t, pb.NetCtxjwtRSAKeyPassword.ToPath, scope.DefaultID, 0): nil,
	}))

	sc, err := jwts.getConfigByScopedGetter(cr.NewScoped(1, 1))
	assert.EqualError(t, err, "[csjwt] invalid key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	assert.Exactly(t, scopedConfig{}, sc)
}

func TestServiceWithBackend_NilScopedGetter(t *testing.T) {
	t.Parallel()

	pb := NewBackend(nil)
	jwts := MustNewService(WithBackend(pb))

	sc, err := jwts.getConfigByScopedGetter(nil)
	assert.NoError(t, err)

	assert.Exactly(t, scope.DefaultHash, sc.scopeHash)
	assert.False(t, sc.Key.IsEmpty())
	assert.Exactly(t, DefaultExpire, sc.expire)
	assert.Exactly(t, csjwt.HS256, sc.signingMethod.Alg())
	assert.False(t, sc.enableJTI)
	assert.True(t, sc.errorHandler == nil)
	assert.True(t, sc.keyFunc != nil)
}
