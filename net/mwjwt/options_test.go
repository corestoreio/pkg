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

package mwjwt_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/mwjwt"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestOptionWithTemplateToken(t *testing.T) {

	jwts, err := mwjwt.NewService(
		// mwjwt.WithKey(scope.Website, 3, csjwt.WithPasswordRandom()),
		mwjwt.WithTemplateToken(scope.Website, 3, func() csjwt.Token {
			sClaim := jwtclaim.NewStore()
			sClaim.Store = "potato"

			h := jwtclaim.NewHeadSegments()
			h.JKU = "https://corestore.io/public.key"

			return csjwt.Token{
				Header: h, // header h has 6 struct fields
				Claims: sClaim,
			}
		}),
	)
	cstesting.FatalIfError(t, err)

	tkDefault, err := jwts.NewToken(scope.Default, 0, jwtclaim.Map{
		"lang": "ch_DE",
	})
	cstesting.FatalIfError(t, err)

	tkWebsite, err := jwts.NewToken(scope.Website, 3, &jwtclaim.Standard{
		Audience: "Gophers",
	})
	cstesting.FatalIfError(t, err)

	tkDefaultParsed, err := jwts.ParseScoped(scope.Default, 0, tkDefault.Raw)
	cstesting.FatalIfError(t, err)
	// t.Logf("tkMissing: %#v\n", tkDefaultParsed)
	lng, err := tkDefaultParsed.Claims.Get("lang")
	cstesting.FatalIfError(t, err)
	assert.Exactly(t, "ch_DE", conv.ToString(lng))

	tkWebsiteParsed, err := jwts.ParseScoped(scope.Website, 3, tkWebsite.Raw)
	cstesting.FatalIfError(t, err)
	// t.Logf("tkFull: %#v\n", tkWebsiteParsed)
	claimStore, err := tkWebsiteParsed.Claims.Get(jwtclaim.KeyStore)
	cstesting.FatalIfError(t, err)
	assert.Exactly(t, "potato", conv.ToString(claimStore))

}

func TestOptionWithTokenID(t *testing.T) {

	jwts, err := mwjwt.NewService(
		mwjwt.WithTokenID(scope.Website, 22, true),
		mwjwt.WithKey(scope.Website, 22, csjwt.WithPasswordRandom()),
	)
	cstesting.FatalIfError(t, err)

	theToken, err := jwts.NewToken(scope.Website, 22)
	cstesting.FatalIfError(t, err)
	assert.NotEmpty(t, theToken.Raw)

	id, err := theToken.Claims.Get(jwtclaim.KeyID)
	cstesting.FatalIfError(t, err)
	assert.Len(t, id, uuidLen)
}

func TestOptionScopedDefaultExpire(t *testing.T) {

	jwts, err := mwjwt.NewService(
		mwjwt.WithTokenID(scope.Website, 33, true),
		mwjwt.WithKey(scope.Website, 33, csjwt.WithPasswordRandom()),
	)
	cstesting.FatalIfError(t, err)

	now := time.Now()
	theToken, err := jwts.NewToken(scope.Website, 33) // must be a pointer the cl or Get() returns nil
	cstesting.FatalIfError(t, err)

	assert.NotEmpty(t, theToken.Raw)
	exp, err := theToken.Claims.Get(jwtclaim.KeyExpiresAt)
	cstesting.FatalIfError(t, err)

	iat, err := theToken.Claims.Get(jwtclaim.KeyIssuedAt)
	cstesting.FatalIfError(t, err)

	assert.Exactly(t, now.Unix(), conv.ToInt64(iat))
	assert.Exactly(t, int(mwjwt.DefaultExpire.Seconds()), int(time.Unix(conv.ToInt64(exp), 0).Sub(now).Seconds()+1))
}

func TestWithMaxSkew_Valid(t *testing.T) {
	jwts, err := mwjwt.NewService(
		mwjwt.WithKey(scope.Website, 44, csjwt.WithPasswordRandom()),
		mwjwt.WithSkew(scope.Website, 44, time.Second*5),
		mwjwt.WithExpiration(scope.Website, 44, -time.Second*3),
	)
	cstesting.FatalIfError(t, err)

	newTK, err := jwts.NewToken(scope.Website, 44, jwtclaim.Map{"key1": "value1"})
	assert.NoError(t, err)

	parsedTK, err := jwts.ParseScoped(scope.Website, 44, newTK.Raw)
	assert.NoError(t, err)
	assert.True(t, parsedTK.Valid, "Token must be valid")

	k1, err1 := parsedTK.Claims.Get("key1")
	cstesting.FatalIfError(t, err1)
	assert.Exactly(t, "value1", k1)
}

func TestWithMaxSkew_NotValid(t *testing.T) {
	jwts, err := mwjwt.NewService(
		mwjwt.WithKey(scope.Default, 0, csjwt.WithPasswordRandom()),
		mwjwt.WithSkew(scope.Default, 0, time.Second*1),
		mwjwt.WithExpiration(scope.Default, 0, -time.Second*3),
	)
	cstesting.FatalIfError(t, err)

	newTK, err := jwts.NewToken(scope.Default, 0, jwtclaim.Map{"key1": "value1"})
	assert.NoError(t, err)

	parsedTK, err := jwts.Parse(newTK.Raw)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
	assert.False(t, parsedTK.Valid, "Token must be NOT valid")

}

func TestOptionWithRSAReaderFail(t *testing.T) {

	jm, err := mwjwt.NewService(
		mwjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromPEM([]byte(`invalid pem data`))),
	)
	assert.Nil(t, jm)
	assert.True(t, errors.IsNotSupported(err), "Error: %s", err)
}

var (
	rsaPrivateKeyFileName        = filepath.Join("..", "..", "util", "csjwt", "test", "test_rsa")
	keyRsaPrivateNoPassword      = mwjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromFile(rsaPrivateKeyFileName))
	keyRsaPrivateWrongPassword   = mwjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromFile(rsaPrivateKeyFileName, []byte(`adfasdf`)))
	keyRsaPrivateCorrectPassword = mwjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromFile(rsaPrivateKeyFileName, []byte("cccamp")))
)

func TestOptionWithRSAFromFileNoOrFailedPassword(t *testing.T) {

	jm, err := mwjwt.NewService(keyRsaPrivateNoPassword)
	assert.True(t, errors.IsEmpty(err), "Error: %s", err)
	assert.Nil(t, jm)

	jm2, err2 := mwjwt.NewService(keyRsaPrivateWrongPassword)
	assert.True(t, errors.IsNotValid(err2), "Error: %s", err2)
	assert.Nil(t, jm2)
}

func testRsaOption(t *testing.T, opt mwjwt.Option) {
	jwts, err := mwjwt.NewService(opt)
	cstesting.FatalIfError(t, err)

	theToken, err := jwts.NewToken(scope.Default, 0, jwtclaim.Map{})
	cstesting.FatalIfError(t, err)
	assert.NotEmpty(t, theToken.Raw)

	tk, err := jwts.Parse(theToken.Raw)
	cstesting.FatalIfError(t, err)
	assert.NotNil(t, tk)
	assert.True(t, tk.Valid)
}

func TestOptionWithRSAFromFilePassword(t *testing.T) {

	testRsaOption(t, keyRsaPrivateCorrectPassword)
}

func TestOptionWithRSAFromFileNoPassword(t *testing.T) {

	testRsaOption(t, mwjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromFile(filepath.Join("..", "..", "util", "csjwt", "test", "test_rsa_np"))))
}
