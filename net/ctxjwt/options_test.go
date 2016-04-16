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

package ctxjwt_test

import (
	"crypto/x509"
	"path/filepath"
	"testing"

	"time"

	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
)

func TestOptionPartialConfigError(t *testing.T) {
	t.Parallel()
	jwts, err := ctxjwt.NewService(ctxjwt.WithTokenID(scope.Website, 3, true))
	if err != nil {
		t.Fatal(err)
	}

	cl := jwtclaim.Map{}
	theToken, err := jwts.NewToken(scope.Website, 3, cl)
	assert.EqualError(t, err, "[ctxjwt] Incomplete configuration for Scope(Website) ID(3). Missing Signing Method and its Key.")
	assert.Empty(t, theToken)
}

func TestOptionWithTokenID(t *testing.T) {
	t.Parallel()
	jwts, err := ctxjwt.NewService(
		ctxjwt.WithTokenID(scope.Website, 22, true),
		ctxjwt.WithKey(scope.Website, 22, csjwt.WithPasswordRandom()),
	)
	if err != nil {
		t.Fatal(err)
	}

	cl := jwtclaim.Map{}
	theToken, err := jwts.NewToken(scope.Website, 22, cl) // must be a pointer the cl or Get() returns nil
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, theToken)
	id, err := cl.Get(jwtclaim.KeyID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, id, uuidLen)
}

func TestOptionScopedDefaultExpire(t *testing.T) {
	t.Parallel()
	jwts, err := ctxjwt.NewService(
		ctxjwt.WithTokenID(scope.Website, 33, true),
		ctxjwt.WithKey(scope.Website, 33, csjwt.WithPasswordRandom()),
	)
	if err != nil {
		t.Fatal(err)
	}

	cl := jwtclaim.Map{}

	now := time.Now()
	theToken, err := jwts.NewToken(scope.Website, 33, cl) // must be a pointer the cl or Get() returns nil
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, theToken)
	exp, err := cl.Get(jwtclaim.KeyExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	iat, err := cl.Get(jwtclaim.KeyIssuedAt)
	if err != nil {
		t.Fatal(err)
	}

	assert.Exactly(t, now.Unix(), iat.(int64))
	assert.Exactly(t, int(ctxjwt.DefaultExpire.Seconds()), int(time.Unix(exp.(int64), 0).Sub(now).Seconds()+1))
}

func TestOptionWithRSAReaderFail(t *testing.T) {
	t.Parallel()
	jm, err := ctxjwt.NewService(
		ctxjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromPEM([]byte(`invalid pem data`))),
	)
	assert.Nil(t, jm)
	assert.EqualError(t, err, `[csjwt] invalid key: Key must be PEM encoded PKCS1 or PKCS8 private key`)
}

var (
	rsaPrivateKeyFileName        = filepath.Join("..", "..", "util", "csjwt", "test", "test_rsa")
	keyRsaPrivateNoPassword      = ctxjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromFile(rsaPrivateKeyFileName))
	keyRsaPrivateWrongPassword   = ctxjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromFile(rsaPrivateKeyFileName, []byte(`adfasdf`)))
	keyRsaPrivateCorrectPassword = ctxjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromFile(rsaPrivateKeyFileName, []byte("cccamp")))
)

func TestOptionWithRSAFromFileNoOrFailedPassword(t *testing.T) {
	t.Parallel()
	jm, err := ctxjwt.NewService(keyRsaPrivateNoPassword)
	assert.EqualError(t, err, "[csjwt] Missing password to decrypt private key")
	assert.Nil(t, jm)

	jm2, err2 := ctxjwt.NewService(keyRsaPrivateWrongPassword)
	assert.EqualError(t, err2, x509.IncorrectPasswordError.Error())
	assert.Nil(t, jm2)
}

func testRsaOption(t *testing.T, opt ctxjwt.Option) {
	jwts, err := ctxjwt.NewService(opt)
	if err != nil {
		t.Fatal(err)
	}

	theToken, err := jwts.NewToken(scope.Default, 0, jwtclaim.Map{})
	if err != nil {
		t.Fatal(cserr.NewMultiErr(err).VerboseErrors())
	}
	assert.NotEmpty(t, theToken)

	tk, err := jwts.Parse(theToken)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, tk)
	assert.True(t, tk.Valid)
}

func TestOptionWithRSAFromFilePassword(t *testing.T) {
	t.Parallel()
	testRsaOption(t, keyRsaPrivateCorrectPassword)
}

func TestOptionWithRSAFromFileNoPassword(t *testing.T) {
	t.Parallel()
	testRsaOption(t, ctxjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromFile(filepath.Join("..", "..", "util", "csjwt", "test", "test_rsa_np"))))
}
