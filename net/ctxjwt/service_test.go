// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"go/build"

	"net/http/httptest"

	"fmt"

	"crypto/x509"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const uuidLen = 36

func TestNewDefault(t *testing.T) {
	jm, err := ctxjwt.NewService()
	assert.NoError(t, err)
	assert.Equal(t, time.Hour, jm.Expire)

	assert.Nil(t, jm.Blacklist.Set("test", time.Hour))
	assert.False(t, jm.Blacklist.Has("test"))

	assert.Len(t, jm.JTI.Get(), uuidLen)

	testClaims := map[string]interface{}{
		"mascot": "gopher",
	}
	theToken, jti, err := jm.GenerateToken(testClaims)
	assert.NoError(t, err)
	assert.Empty(t, jti)
	assert.NotEmpty(t, theToken)

	haveToken, err := jm.Parse(theToken)
	assert.NoError(t, err)
	assert.True(t, haveToken.Valid)
	assert.Equal(t, "gopher", haveToken.Claims["mascot"])

	failedToken, err := jm.Parse(theToken + "c")
	assert.Error(t, err)
	assert.Nil(t, failedToken)

	jmRSA, err := ctxjwt.NewService(ctxjwt.WithRSAFromFile("invalid.key"))
	assert.Nil(t, jmRSA)
	assert.Contains(t, err.Error(), "open invalid.key: no such file or directory")
}

func TestInvalidSigningMethod(t *testing.T) {
	password := []byte(`Rump3lst!lzch3n`)
	jm, err := ctxjwt.NewService(
		ctxjwt.WithPassword(password),
	)
	assert.NoError(t, err)

	tk := jwt.New(jwt.SigningMethodHS256)
	tk.Claims["exp"] = time.Now().Add(time.Hour).Unix()
	tk.Claims["iat"] = time.Now().Unix()
	tk.Header["alg"] = "HS384"
	malformedToken, err := tk.SignedString(password)
	assert.NoError(t, err)

	mt, err := jm.Parse(malformedToken)
	assert.EqualError(t, err, ctxjwt.ErrUnexpectedSigningMethod.Error())
	assert.Nil(t, mt)
}

func TestJTI(t *testing.T) {
	jm, err := ctxjwt.NewService()
	assert.NoError(t, err)
	jm.EnableJTI = true

	theToken, jti, err := jm.GenerateToken(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, jti)
	assert.NotEmpty(t, theToken)
	assert.Len(t, jti, uuidLen)
}

type testBL struct {
	*testing.T
	theToken string
	exp      time.Duration
}

func (b *testBL) Set(theToken string, exp time.Duration) error {
	b.theToken = theToken
	b.exp = exp
	return nil
}
func (b *testBL) Has(_ string) bool { return false }

func TestLogout(t *testing.T) {

	tbl := &testBL{T: t}
	jm, err := ctxjwt.NewService()
	assert.NoError(t, err)
	jm.Blacklist = tbl

	theToken, _, err := jm.GenerateToken(nil)
	assert.NoError(t, err)

	tk, err := jm.Parse(theToken)
	assert.NoError(t, err)
	assert.NotNil(t, tk)

	assert.Nil(t, jm.Logout(nil))
	assert.Nil(t, jm.Logout(tk))
	assert.Equal(t, int(time.Hour.Seconds()), 1+int(tbl.exp.Seconds()))
	assert.Equal(t, theToken, tbl.theToken)
}

var pkFile = filepath.Join(build.Default.GOPATH, "src", "github.com", "corestoreio", "csfw", "net", "ctxjwt", "test_rsa")

func TestRSAEncryptedNoOrFailedPassword(t *testing.T) {
	jm, err := ctxjwt.NewService(ctxjwt.WithRSAFromFile(pkFile))
	assert.EqualError(t, err, ctxjwt.ErrPrivateKeyNoPassword.Error())
	assert.Nil(t, jm)

	jm2, err2 := ctxjwt.NewService(ctxjwt.WithRSAFromFile(pkFile, []byte(`adfasdf`)))
	assert.EqualError(t, err2, x509.IncorrectPasswordError.Error())
	assert.Nil(t, jm2)
}

func testRsaOption(t *testing.T, opt ctxjwt.Option) {
	jm, err := ctxjwt.NewService(opt)
	assert.NoError(t, err)
	assert.NotNil(t, jm)

	theToken, _, err := jm.GenerateToken(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken)

	tk, err := jm.Parse(theToken)
	assert.NoError(t, err)
	assert.NotNil(t, tk)
	assert.True(t, tk.Valid)
}

func TestRSAEncryptedPassword(t *testing.T) {
	pw := []byte("cccamp")
	testRsaOption(t, ctxjwt.WithRSAFromFile(pkFile, pw))
}

func TestRSAWithoutPassword(t *testing.T) {
	pkFileNP := filepath.Join(build.Default.GOPATH, "src", "github.com", "corestoreio", "csfw", "net", "ctxjwt", "test_rsa_np")
	testRsaOption(t, ctxjwt.WithRSAFromFile(pkFileNP))
}

func TestRSAGenerate(t *testing.T) {
	testRsaOption(t, ctxjwt.WithRSAGenerator())
}

func testAuth(t *testing.T, errH ctxhttp.Handler, opts ...ctxjwt.Option) (ctxhttp.Handler, string) {
	jm, err := ctxjwt.NewService(opts...)
	assert.NoError(t, err)
	theToken, _, err := jm.GenerateToken(map[string]interface{}{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	assert.NoError(t, err)

	final := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})
	authHandler := jm.WithParseAndValidate(errH)(final)
	return authHandler, theToken
}

func TestWithParseAndValidateNoToken(t *testing.T) {

	authHandler, _ := testAuth(t, nil)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	authHandler.ServeHTTPContext(context.Background(), w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, w.Body.String(), http.StatusText(http.StatusUnauthorized)+"\n")
}

func TestWithParseAndValidateHTTPErrorHandler(t *testing.T) {

	authHandler, _ := testAuth(t, ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		err, ok := ctxjwt.FromContextWithError(ctx)
		assert.True(t, ok)

		w.WriteHeader(http.StatusTeapot)
		_, err = w.Write([]byte(err.Error()))
		return err
	}))

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	authHandler.ServeHTTPContext(context.Background(), w, req)
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Equal(t, "no token present in request", w.Body.String())
}

func TestWithParseAndValidateSuccess(t *testing.T) {
	jm, err := ctxjwt.NewService()
	assert.NoError(t, err)

	theToken, _, err := jm.GenerateToken(map[string]interface{}{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+theToken)

	finalHandler := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprintf(w, "I'm more of a coffee pot")

		ctxToken, ok := ctxjwt.FromContext(ctx)
		assert.True(t, ok)
		assert.NotNil(t, ctxToken)
		assert.Exactly(t, "bar", ctxToken.Claims["xfoo"].(string))

		return nil
	})
	authHandler := jm.WithParseAndValidate()(finalHandler)

	wRec := httptest.NewRecorder()
	authHandler.ServeHTTPContext(context.Background(), wRec, req)

	assert.Equal(t, http.StatusTeapot, wRec.Code)
	assert.Equal(t, `I'm more of a coffee pot`, wRec.Body.String())
}

type testRealBL struct {
	token string
	exp   time.Duration
}

func (b *testRealBL) Set(t string, exp time.Duration) error {
	b.token = t
	b.exp = exp
	return nil
}
func (b *testRealBL) Has(t string) bool { return b.token == t }

func TestWithParseAndValidateInBlackList(t *testing.T) {
	jm, err := ctxjwt.NewService()
	assert.NoError(t, err)

	bl := &testRealBL{}
	jm.Blacklist = bl
	theToken, _, err := jm.GenerateToken(nil)
	bl.token = theToken
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer   "+theToken)

	finalHandler := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		return nil
	})
	authHandler := jm.WithParseAndValidate()(finalHandler)

	wRec := httptest.NewRecorder()
	authHandler.ServeHTTPContext(context.Background(), wRec, req)

	assert.NotEqual(t, http.StatusTeapot, wRec.Code)
	assert.Equal(t, http.StatusUnauthorized, wRec.Code)
}

// todo add test for form with input field: access_token
