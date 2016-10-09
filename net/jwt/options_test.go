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

package jwt_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/jwt"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionWithTemplateToken(t *testing.T) {

	jwts, err := jwt.New(
		// jwt.WithKey(scope.Website.ToHash(3), csjwt.WithPasswordRandom()),
		jwt.WithTemplateToken(func() csjwt.Token {
			sClaim := jwtclaim.NewStore()
			sClaim.Store = "potato"

			h := jwtclaim.NewHeadSegments()
			h.JKU = "https://corestore.io/public.key"

			return csjwt.Token{
				Header: h, // header h has 6 struct fields
				Claims: sClaim,
			}
		}, scope.Website.Pack(3)),
	)
	require.NoError(t, err)

	tkDefault, err := jwts.NewToken(scope.DefaultTypeID, jwtclaim.Map{
		"lang": "ch_DE",
	})
	require.NoError(t, err, "%+v", err)

	tkWebsite, err := jwts.NewToken(scope.Website.Pack(3), &jwtclaim.Standard{
		Audience: "Gophers",
	})
	require.NoError(t, err)

	tkDefaultParsed, err := jwts.ParseScoped(scope.DefaultTypeID, tkDefault.Raw)
	require.NoError(t, err)
	// t.Logf("tkMissing: %#v\n", tkDefaultParsed)
	lng, err := tkDefaultParsed.Claims.Get("lang")
	require.NoError(t, err)
	assert.Exactly(t, "ch_DE", conv.ToString(lng))

	tkWebsiteParsed, err := jwts.ParseScoped(scope.Website.Pack(3), tkWebsite.Raw)
	require.NoError(t, err)
	// t.Logf("tkFull: %#v\n", tkWebsiteParsed)
	claimStore, err := tkWebsiteParsed.Claims.Get(jwtclaim.KeyStore)
	require.NoError(t, err)
	assert.Exactly(t, "potato", conv.ToString(claimStore))

}

func TestOptionWithTokenID(t *testing.T) {

	jwts, err := jwt.New(
		jwt.WithKey(csjwt.WithPasswordRandom(), scope.Website.Pack(22)),
	)
	require.NoError(t, err)

	theToken, err := jwts.NewToken(scope.Website.Pack(22))
	require.NoError(t, err)
	assert.NotEmpty(t, theToken.Raw)

	id, err := theToken.Claims.Get(jwtclaim.KeyID)
	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestOptionScopedDefaultExpire(t *testing.T) {

	jwts, err := jwt.New(
		jwt.WithKey(csjwt.WithPasswordRandom(), scope.Website.Pack(33)),
	)
	require.NoError(t, err)

	now := time.Now()
	theToken, err := jwts.NewToken(scope.Website.Pack(33)) // must be a pointer the cl or Get() returns nil
	require.NoError(t, err)

	assert.NotEmpty(t, theToken.Raw)
	exp, err := theToken.Claims.Get(jwtclaim.KeyExpiresAt)
	require.NoError(t, err)

	iat, err := theToken.Claims.Get(jwtclaim.KeyIssuedAt)
	require.NoError(t, err)

	assert.Exactly(t, now.Unix(), conv.ToInt64(iat))
	assert.Exactly(t, int(jwt.DefaultExpire.Seconds()), int(time.Unix(conv.ToInt64(exp), 0).Sub(now).Seconds()+1))
}

func TestWithMaxSkew_Valid(t *testing.T) {
	jwts, err := jwt.New(
		jwt.WithKey(csjwt.WithPasswordRandom(), scope.Website.Pack(44)),
		jwt.WithSkew(time.Second*5, scope.Website.Pack(44)),
		jwt.WithExpiration(-time.Second*3, scope.Website.Pack(44)),
	)
	require.NoError(t, err)

	newTK, err := jwts.NewToken(scope.Website.Pack(44), jwtclaim.Map{"key1": "value1"})
	assert.NoError(t, err)

	parsedTK, err := jwts.ParseScoped(scope.Website.Pack(44), newTK.Raw)
	assert.NoError(t, err)
	assert.True(t, parsedTK.Valid, "Token must be valid")

	k1, err := parsedTK.Claims.Get("key1")
	require.NoError(t, err)
	assert.Exactly(t, "value1", k1)
}

func TestWithMaxSkew_NotValid(t *testing.T) {
	jwts, err := jwt.New(
		// DefaultScopeID
		jwt.WithKey(csjwt.WithPasswordRandom()),
		jwt.WithSkew(time.Second*1),
		jwt.WithExpiration(-time.Second*3),
	)
	require.NoError(t, err)

	newTK, err := jwts.NewToken(scope.DefaultTypeID, jwtclaim.Map{"key1": "value1"})
	assert.NoError(t, err)

	parsedTK, err := jwts.Parse(newTK.Raw)
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	assert.False(t, parsedTK.Valid, "Token must be NOT valid")

}

func TestOptionWithRSAReaderFail(t *testing.T) {

	jm, err := jwt.New(
		jwt.WithKey(csjwt.WithRSAPrivateKeyFromPEM([]byte(`invalid pem data`))), // scope.DefaultTypeID
	)
	assert.Nil(t, jm)
	assert.True(t, errors.IsNotSupported(err), "Error: %+v", err)
}

var (
	rsaPrivateKeyFileName = filepath.Join("..", "..", "util", "csjwt", "test", "test_rsa")
	// Next three configurations for the DefaultScopeID
	keyRsaPrivateNoPassword      = jwt.WithKey(csjwt.WithRSAPrivateKeyFromFile(rsaPrivateKeyFileName))
	keyRsaPrivateWrongPassword   = jwt.WithKey(csjwt.WithRSAPrivateKeyFromFile(rsaPrivateKeyFileName, []byte(`adfasdf`)))
	keyRsaPrivateCorrectPassword = jwt.WithKey(csjwt.WithRSAPrivateKeyFromFile(rsaPrivateKeyFileName, []byte("cccamp")))
)

func TestOptionWithRSAFromFileNoOrFailedPassword(t *testing.T) {

	jm, err := jwt.New(keyRsaPrivateNoPassword)
	assert.True(t, errors.IsEmpty(err), "Error: %+v", err)
	assert.Nil(t, jm)

	jm2, err := jwt.New(keyRsaPrivateWrongPassword)
	assert.True(t, errors.IsNotValid(err), "Error: %+v\nType %d", err, errors.HasBehaviour(err))
	assert.Nil(t, jm2)
}

func testRsaOption(t *testing.T, opt jwt.Option) {
	jwts, err := jwt.New(opt)
	require.NoError(t, err)

	theToken, err := jwts.NewToken(scope.DefaultTypeID, jwtclaim.Map{})
	require.NoError(t, err)
	assert.NotEmpty(t, theToken.Raw)

	tk, err := jwts.Parse(theToken.Raw)
	require.NoError(t, err)
	assert.NotNil(t, tk)
	assert.True(t, tk.Valid)
}

func TestOptionWithRSAFromFilePassword(t *testing.T) {
	testRsaOption(t, keyRsaPrivateCorrectPassword)
}

func TestOptionWithRSAFromFileNoPassword(t *testing.T) {
	testRsaOption(t, jwt.WithKey(csjwt.WithRSAPrivateKeyFromFile(filepath.Join("..", "..", "util", "csjwt", "test", "test_rsa_np"))))
}
