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

package jwt

import (
	"hash/fnv"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/corestoreio/csfw/util/blacklist"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/shortid"
	"github.com/stretchr/testify/assert"
)

func TestScopedConfig_ParseFromRequest_Valid(t *testing.T) {
	bl := blacklist.NewInMemory(fnv.New64a)
	sc := newScopedConfig()
	kid := shortid.MustGenerate()
	tk := csjwt.NewToken(jwtclaim.Map{"jti": kid})
	token, err := tk.SignedString(sc.SigningMethod, sc.Key)
	assert.NoError(t, err, "%+v", err)

	req := httptest.NewRequest("GET", "https://token-service.corestore.io", nil)
	SetHeaderAuthorization(req, token)
	reqToken, err := sc.ParseFromRequest(bl, req)
	assert.NoError(t, err, "%+v", err)

	assert.True(t, reqToken.Valid)
	assert.Exactly(t, token, reqToken.Raw)
}

func TestScopedConfig_ParseFromRequest_Invalid_Token(t *testing.T) {
	sc := newScopedConfig()
	tk := csjwt.NewToken(jwtclaim.Map{})
	_, err := tk.SignedString(sc.SigningMethod, sc.Key)
	assert.NoError(t, err, "%+v", err)

	req := httptest.NewRequest("GET", "https://token-service.corestore.io", nil)

	reqToken, err := sc.ParseFromRequest(nil, req)
	assert.True(t, errors.IsNotFound(err), "%+v", err)

	assert.False(t, reqToken.Valid)
	assert.Empty(t, reqToken.Raw)
}

func TestScopedConfig_ParseFromRequest_Invalid_JTI(t *testing.T) {

	sc := newScopedConfig()
	tk := csjwt.NewToken(jwtclaim.Map{})
	token, err := tk.SignedString(sc.SigningMethod, sc.Key)
	assert.NoError(t, err, "%+v", err)

	req := httptest.NewRequest("GET", "https://token-service.corestore.io", nil)
	SetHeaderAuthorization(req, token)
	reqToken, err := sc.ParseFromRequest(nil, req)
	assert.True(t, errors.IsEmpty(err), "%+v", err)

	assert.True(t, reqToken.Valid)
	assert.Exactly(t, token, reqToken.Raw)
}

func TestScopedConfig_ParseFromRequest_In_Blacklist(t *testing.T) {
	bl := blacklist.NewInMemory(fnv.New64a)
	sc := newScopedConfig()
	kid := shortid.MustGenerate()
	assert.NoError(t, bl.Set([]byte(kid), time.Hour))
	tk := csjwt.NewToken(jwtclaim.Map{"jti": kid})
	token, err := tk.SignedString(sc.SigningMethod, sc.Key)
	assert.NoError(t, err, "%+v", err)

	req := httptest.NewRequest("GET", "https://token-service.corestore.io", nil)
	SetHeaderAuthorization(req, token)
	reqToken, err := sc.ParseFromRequest(bl, req)
	assert.True(t, errors.IsNotValid(err), "%+v", err)

	assert.True(t, reqToken.Valid)
	assert.Exactly(t, token, reqToken.Raw)
}

type errBl struct {
	setErr error
	has    bool
}

func (e errBl) Set(id []byte, expires time.Duration) error {
	return e.setErr
}
func (e errBl) Has(id []byte) bool {
	return e.has
}

var _ Blacklister = (*errBl)(nil)

func TestScopedConfig_ParseFromRequest_SingleTokenUsage_BL_Set_Error(t *testing.T) {

	sc := newScopedConfig()
	sc.SingleTokenUsage = true
	kid := shortid.MustGenerate()

	tk := csjwt.NewToken(jwtclaim.Map{"jti": kid})
	token, err := tk.SignedString(sc.SigningMethod, sc.Key)
	assert.NoError(t, err, "%+v", err)

	req := httptest.NewRequest("GET", "https://token-service.corestore.io", nil)
	SetHeaderAuthorization(req, token)
	reqToken, err := sc.ParseFromRequest(errBl{
		setErr: errors.NewAlreadyClosedf("Remote server closed"),
	}, req)
	assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)

	assert.True(t, reqToken.Valid)
	assert.Exactly(t, token, reqToken.Raw)
}

// todo investigate allocs
// 200000	      9072 ns/op	    1529 B/op	      32 allocs/op
func BenchmarkScopedConfig_ParseFromRequest_HS256Fast_FNV64a(b *testing.B) {
	bl := blacklist.NewInMemory(fnv.New64a)

	for i := 0; i < 10000; i++ {
		kid := []byte(shortid.MustGenerate())
		if err := bl.Set(kid, time.Second*time.Duration(i)); err != nil {
			b.Fatalf("%+v", err)
		}
	}

	sc := newScopedConfig()
	kid := shortid.MustGenerate()

	tk := csjwt.NewToken(&jwtclaim.Standard{ID: kid, ExpiresAt: time.Now().Add(time.Hour).Unix()})
	token, err := tk.SignedString(sc.SigningMethod, sc.Key)
	if err != nil {
		b.Fatalf("%+v", err)
	}

	req := httptest.NewRequest("GET", "https://token-service.corestore.io", nil)
	SetHeaderAuthorization(req, token)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reqToken, err := sc.ParseFromRequest(bl, req)
		if err != nil {
			b.Fatalf("%+v", err)
		}
		if !reqToken.Valid {
			b.Fatalf("Token should be valid: %#v", reqToken)
		}
	}
}
