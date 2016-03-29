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
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var _ error = (*ctxjwt.Service)(nil)

const uuidLen = 36

func TestServiceMustNewServicePanic(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), "open non-existent.pem: no such file or directory")
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	_ = ctxjwt.MustNewService(ctxjwt.WithECDSAFromFile(scope.DefaultID, 0, "non-existent.pem"))
}

func TestServiceNewDefault(t *testing.T) {
	t.Parallel()
	jwts := ctxjwt.MustNewService()

	key := []byte("test")
	assert.Nil(t, jwts.Blacklist.Set(key, time.Hour))
	assert.False(t, jwts.Blacklist.Has(key))

	assert.Len(t, jwts.JTI.Get(), uuidLen)

	testClaims := map[string]interface{}{
		"mascot": "gopher",
	}
	theToken, jti, err := jwts.GenerateToken(scope.DefaultID, 0, testClaims)
	assert.NoError(t, err)
	assert.Empty(t, jti)
	assert.NotEmpty(t, theToken)
	haveToken, err := jwts.Parse(theToken)
	assert.NoError(t, err)
	assert.True(t, haveToken.Valid)
	assert.Equal(t, "gopher", haveToken.Claims["mascot"])

	failedToken, err := jwts.Parse(theToken + "c")
	assert.Error(t, err)
	assert.Nil(t, failedToken)

	jmRSA, err := ctxjwt.NewService(ctxjwt.WithRSAFromFile(scope.DefaultID, 0, "invalid.key"))
	assert.Nil(t, jmRSA)
	assert.Contains(t, err.Error(), "open invalid.key:") //  no such file or directory OR The system cannot find the file specified.
}

func TestServiceParseInvalidSigningMethod(t *testing.T) {
	t.Parallel()
	password := []byte(`Rump3lst!lzch3n`)
	jwts := ctxjwt.MustNewService(
		ctxjwt.WithPassword(scope.DefaultID, 0, password),
	)

	tk := jwt.New(jwt.SigningMethodHS256)
	tk.Claims["exp"] = time.Now().Add(time.Hour).Unix()
	tk.Claims["iat"] = time.Now().Unix()
	tk.Header["alg"] = "HS384"
	malformedToken, err := tk.SignedString(password)
	assert.NoError(t, err)

	mt, err := jwts.Parse(malformedToken)
	assert.EqualError(t, err, ctxjwt.ErrUnexpectedSigningMethod.Error())
	assert.Nil(t, mt)
}

func TestServiceLogout(t *testing.T) {
	t.Parallel()

	tbl := &testBL{T: t}
	jwts := ctxjwt.MustNewService(
		ctxjwt.WithBlacklist(tbl),
	)

	theToken, _, err := jwts.GenerateToken(scope.DefaultID, 0, nil)
	assert.NoError(t, err)

	tk, err := jwts.Parse(theToken)
	assert.NoError(t, err)
	assert.NotNil(t, tk)

	assert.Nil(t, jwts.Logout(nil))
	assert.Nil(t, jwts.Logout(tk))
	assert.Equal(t, int(time.Hour.Seconds()), 1+int(tbl.exp.Seconds()))
	assert.Equal(t, theToken, string(tbl.theToken))
}
