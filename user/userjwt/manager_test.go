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

package userjwt_test

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"go/build"

	"net/http/httptest"

	"strings"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/user/userjwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

const uuidLen = 36

func TestNewDefault(t *testing.T) {
	jm, err := userjwt.New()
	assert.NoError(t, err)
	assert.Equal(t, time.Hour, jm.Expire)

	assert.Nil(t, jm.Blacklist.Set("test", time.Hour))
	assert.False(t, jm.Blacklist.Has("test"))
	assert.Equal(t, userjwt.PostFormVarPrefix, jm.PostFormVarPrefix)
	assert.Len(t, jm.JTI.Get(), uuidLen)

	testClaims := map[string]interface{}{
		"mascot": "gopher",
	}
	token, jti, err := jm.GenerateToken(testClaims)
	assert.NoError(t, err)
	assert.Empty(t, jti)
	assert.NotEmpty(t, token)

	haveToken, err := jm.Parse(token)
	assert.NoError(t, err)
	assert.True(t, haveToken.Valid)
	assert.Equal(t, "gopher", haveToken.Claims["mascot"])

	failedToken, err := jm.Parse(token + "c")
	assert.Error(t, err)
	assert.Nil(t, failedToken)

	jmRSA, err := userjwt.New(userjwt.RSAFromFile("invalid.key"))
	assert.Nil(t, jmRSA)
	assert.Contains(t, err.Error(), "open invalid.key: no such file or directory")
}

func TestInvalidSigningMethod(t *testing.T) {
	password := []byte(`Rump3lst!lzch3n`)
	jm, err := userjwt.New(
		userjwt.Password(password),
	)
	assert.NoError(t, err)

	tk := jwt.New(jwt.SigningMethodHS256)
	tk.Claims["exp"] = time.Now().Add(time.Hour).Unix()
	tk.Claims["iat"] = time.Now().Unix()
	tk.Header["alg"] = "HS384"
	malformedToken, err := tk.SignedString(password)
	assert.NoError(t, err)

	mt, err := jm.Parse(malformedToken)
	assert.EqualError(t, err, userjwt.ErrUnexpectedSigningMethod.Error())
	assert.Nil(t, mt)
}

func TestPasswordFromConfig(t *testing.T) {

	cfg := config.NewMockReader(
		config.MockString(func(path string) string {
			if path == userjwt.PathJWTPassword {
				return `Rump3lst!lzch3n`
			}
			return ""
		}),
	)

	jm, err := userjwt.New(
		userjwt.PasswordFromConfig(cfg),
	)
	assert.NoError(t, err)

	token, _, err := jm.GenerateToken(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

}

func TestJTI(t *testing.T) {
	jm, err := userjwt.New()
	assert.NoError(t, err)
	jm.EnableJTI = true

	token, jti, err := jm.GenerateToken(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, jti)
	assert.NotEmpty(t, token)
	assert.Len(t, jti, uuidLen)
}

type testBL struct {
	*testing.T
	token string
	exp   time.Duration
}

func (b *testBL) Set(token string, exp time.Duration) error {
	b.token = token
	b.exp = exp
	return nil
}
func (b *testBL) Has(_ string) bool { return false }

func TestLogout(t *testing.T) {

	tbl := &testBL{T: t}
	jm, err := userjwt.New()
	assert.NoError(t, err)
	jm.Blacklist = tbl

	token, _, err := jm.GenerateToken(nil)
	assert.NoError(t, err)

	tk, err := jm.Parse(token)
	assert.NoError(t, err)
	assert.NotNil(t, tk)

	assert.Nil(t, jm.Logout(nil))
	assert.Nil(t, jm.Logout(tk))
	assert.Equal(t, int(time.Hour.Seconds()), 1+int(tbl.exp.Seconds()))
	assert.Equal(t, token, tbl.token)
}

var pkFile = filepath.Join(build.Default.GOPATH, "src", "github.com", "corestoreio", "csfw", "user", "userjwt", "test_rsa")

func TestRSAEncryptedNoOrFailedPassword(t *testing.T) {
	jm, err := userjwt.New(userjwt.RSAFromFile(pkFile))
	assert.Contains(t, err.Error(), "Private Key is encrypted but password was not set")
	assert.Nil(t, jm)

	jm2, err2 := userjwt.New(userjwt.RSAFromFile(pkFile, []byte(`adfasdf`)))
	assert.Contains(t, err2.Error(), "Private Key decryption failed: x509: decryption password incorrect")
	assert.Nil(t, jm2)
}

func testRsaOption(t *testing.T, opt userjwt.OptionFunc) {
	jm, err := userjwt.New(opt)
	assert.NoError(t, err)
	assert.NotNil(t, jm)

	token, _, err := jm.GenerateToken(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	tk, err := jm.Parse(token)
	assert.NoError(t, err)
	assert.NotNil(t, tk)
	assert.True(t, tk.Valid)
}

func TestRSAEncryptedPassword(t *testing.T) {
	pw := []byte("cccamp")
	testRsaOption(t, userjwt.RSAFromFile(pkFile, pw))
}

func TestRSAWithoutPassword(t *testing.T) {
	pkFileNP := filepath.Join(build.Default.GOPATH, "src", "github.com", "corestoreio", "csfw", "user", "userjwt", "test_rsa_np")
	testRsaOption(t, userjwt.RSAFromFile(pkFileNP))
}

func TestRSAGenerate(t *testing.T) {
	testRsaOption(t, userjwt.RSAGenerate())
}

func testAuth(t *testing.T, opts ...userjwt.OptionFunc) (http.Handler, string) {
	jm, err := userjwt.New(opts...)
	assert.NoError(t, err)
	token, _, err := jm.GenerateToken(map[string]interface{}{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	assert.NoError(t, err)

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// t.Logf("Check Form for map: %v\n", r.Form)
		fmt.Fprint(w, r.Form.Encode())
	})
	authHandler := jm.Authorization(final)
	return authHandler, token
}

func TestAuthorizationNoToken(t *testing.T) {

	authHandler, _ := testAuth(t)

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `no token present in request`)
}

func TestAuthorizationHTTPErrorHandler(t *testing.T) {

	authHandler, _ := testAuth(t, func(a *userjwt.AuthManager) {

		a.HTTPErrorHandler = func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
				w.Write([]byte(err.Error()))
			})
		}
	})

	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Equal(t, "no token present in request", w.Body.String())
}

func TestAuthorizationValidTokenPOST(t *testing.T) {

	authHandler, token := testAuth(t)

	req, err := http.NewRequest("POST", "http://auth.xyz", strings.NewReader("a=b&c=d"))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(7))

	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `a=b&c=d&jwt__exp=`)
	assert.Contains(t, w.Body.String(), `&jwt__xfoo=bar&jwt__zfoo=4711`)

}

func TestAuthorizationValidTokenPOSTNoBody(t *testing.T) {

	authHandler, token := testAuth(t)

	req, err := http.NewRequest("POST", "http://auth.xyz", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `missing form body`)

}
