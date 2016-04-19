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
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
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
	_ = ctxjwt.MustNewService(ctxjwt.WithKey(scope.Default, 0, csjwt.WithECPrivateKeyFromFile("non-existent.pem")))
}

func TestServiceNewDefaultBlacklist(t *testing.T) {
	t.Parallel()
	jwts := ctxjwt.MustNewService()

	key := []byte("test")
	assert.Nil(t, jwts.Blacklist.Set(key, time.Hour))
	assert.False(t, jwts.Blacklist.Has(key))
	assert.Len(t, jwts.JTI.Get(), uuidLen)
}

func TestServiceNewDefault(t *testing.T) {
	t.Parallel()
	jwts := ctxjwt.MustNewService()

	testClaims := &jwtclaim.Standard{
		Subject: "gopher",
	}
	theToken, err := jwts.NewToken(scope.Default, 0, testClaims)
	assert.NoError(t, err)
	assert.Empty(t, testClaims.ID)
	assert.NotEmpty(t, theToken.Raw)

	haveToken, err := jwts.Parse(theToken.Raw)
	assert.NoError(t, err)
	assert.True(t, haveToken.Valid)

	mascot, err := haveToken.Claims.Get(jwtclaim.KeySubject)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "gopher", mascot.(string))

	failedToken, err := jwts.Parse(append(text.Chars(theToken.Raw).Clone(), []byte("c")...))
	assert.Error(t, err)
	assert.False(t, failedToken.Valid)
}

func TestServiceNewDefaultRSAError(t *testing.T) {
	t.Parallel()

	jmRSA, err := ctxjwt.NewService(ctxjwt.WithKey(scope.Default, 0, csjwt.WithRSAPrivateKeyFromFile("invalid.key")))
	assert.Nil(t, jmRSA)
	assert.Contains(t, err.Error(), "open invalid.key:") //  no such file or directory OR The system cannot find the file specified.
}

type malformedSigner struct {
	*csjwt.SigningMethodHMAC
}

func (ms malformedSigner) Alg() string {
	return "None"
}

func TestServiceParseInvalidSigningMethod(t *testing.T) {
	t.Parallel()

	ms := &malformedSigner{
		SigningMethodHMAC: csjwt.NewSigningMethodHS256(),
	}

	keyRand := csjwt.WithPasswordRandom()
	jwts := ctxjwt.MustNewService(ctxjwt.WithKey(scope.Default, 0, keyRand))

	tk := csjwt.NewToken(jwtclaim.Map{
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	malformedToken, err := tk.SignedString(ms, keyRand)
	assert.NoError(t, err)

	mt, err := jwts.Parse(malformedToken)
	assert.EqualError(t, err, "[csjwt] token is unverifiable\n[ctxjwt] Unknown signing method - Have: \"None\" Want: \"HS256\"")
	assert.False(t, mt.Valid)
}

type testBL struct {
	*testing.T
	theToken []byte
	exp      time.Duration
}

func (b *testBL) Set(theToken []byte, exp time.Duration) error {
	b.theToken = theToken
	b.exp = exp
	return nil
}
func (b *testBL) Has(_ []byte) bool { return false }

var _ ctxjwt.Blacklister = (*testBL)(nil)

func TestServiceLogout(t *testing.T) {
	t.Parallel()

	tbl := &testBL{T: t}
	jwts := ctxjwt.MustNewService(
		ctxjwt.WithBlacklist(tbl),
	)

	theToken, err := jwts.NewToken(scope.Default, 0, jwtclaim.NewStore())
	assert.NoError(t, err)

	tk, err := jwts.Parse(theToken.Raw)
	assert.NoError(t, err)
	assert.NotNil(t, tk)

	assert.Nil(t, jwts.Logout(csjwt.Token{}))
	assert.Nil(t, jwts.Logout(tk))
	assert.Equal(t, int(time.Hour.Seconds()), 1+int(tbl.exp.Seconds()))
	assert.Equal(t, string(theToken.Raw), string(tbl.theToken))
}

func TestServiceIncorrectConfigurationScope(t *testing.T) {
	t.Parallel()

	jwts, err := ctxjwt.NewService(ctxjwt.WithKey(scope.Store, 33, csjwt.WithPasswordRandom()))
	assert.Nil(t, jwts)
	assert.EqualError(t, err, `[ctxjwt] Service does not support this: Scope(Store) ID(33). Only default or website are allowed.`)
}

func TestService_NewToken_Merge_Maps(t *testing.T) {
	t.Parallel()
	jwts, err := ctxjwt.NewService(
		ctxjwt.WithKey(scope.Website, 3, csjwt.WithPasswordRandom()),
	)
	if err != nil {
		t.Fatal(err)
	}

	// NewToken has an underlying map as a claimer
	theToken, err := jwts.NewToken(scope.Website, 3, jwtclaim.Map{
		"xk1": 2.718281,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, theToken.Raw)
	id, err := theToken.Claims.Get("xk1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, 2.718281, id)
}

func TestService_NewToken_Merge_Structs(t *testing.T) {
	t.Parallel()
	jwts, err := ctxjwt.NewService(
		ctxjwt.WithKey(scope.Website, 4, csjwt.WithPasswordRandom()),
		ctxjwt.WithTemplateToken(scope.Website, 4, func() csjwt.Token {
			s := jwtclaim.NewStore()
			s.Store = "de"
			return csjwt.NewToken(s)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	// NewToken has an underlying jwtclaim.NewStore as a claimer
	theToken, err := jwts.NewToken(scope.Website, 4, jwtclaim.Map{
		jwtclaim.KeyUserID: "0815",
	})
	if err != nil {
		t.Fatal(cserr.NewMultiErr(err).VerboseErrors())
	}
	assert.NotEmpty(t, theToken.Raw)

	storeID, err := theToken.Claims.Get(jwtclaim.KeyStore)
	if err != nil {
		t.Fatal(cserr.NewMultiErr(err).VerboseErrors())
	}
	assert.Exactly(t, "de", storeID)

	userID, err := theToken.Claims.Get(jwtclaim.KeyUserID)
	if err != nil {
		t.Fatal(cserr.NewMultiErr(err).VerboseErrors())
	}
	assert.Exactly(t, "0815", userID)

	expires, err := theToken.Claims.Get(jwtclaim.KeyExpiresAt)
	if err != nil {
		t.Fatal(cserr.NewMultiErr(err).VerboseErrors())
	}
	assert.True(t, conv.ToInt64(expires) > time.Now().Unix())
}

func TestService_NewToken_Merge_Fail(t *testing.T) {
	t.Parallel()
	jwts, err := ctxjwt.NewService(
		ctxjwt.WithKey(scope.Website, 4, csjwt.WithPasswordRandom()),
		ctxjwt.WithTemplateToken(scope.Website, 4, func() csjwt.Token {
			return csjwt.NewToken(&jwtclaim.Standard{})
		}),
	)
	if err != nil {
		t.Fatal(cserr.NewMultiErr(err).VerboseErrors())
	}

	// NewToken has an underlying jwtclaim.NewStore as a claimer
	theToken, err := jwts.NewToken(scope.Website, 4, jwtclaim.Map{
		jwtclaim.KeyUserID: "0815",
	})

	assert.EqualError(t, err, `[csjwt] Cannot set Key "userid" with value 0815. Error: [jwtclaim] Claim "userid" not supported.`)
	assert.Empty(t, theToken.Raw)

}
