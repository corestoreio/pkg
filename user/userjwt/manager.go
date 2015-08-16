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

package userjwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/corestoreio/csfw/utils/cast"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/juju/errgo"
)

// ErrUnexpectedSigningMethod will be returned if some outside dude tries to trick us
var ErrUnexpectedSigningMethod = errors.New("JWT: Unexpected signing method")

// PostFormVarPrefix defines the prefix of the imported token claims into the
// request.Form map.
const PostFormVarPrefix = "jwt__"

// AuthManager main object for handling JWT authentication, generation, blacklists and log outs.
type AuthManager struct {
	rsapk    *rsa.PrivateKey
	ecdsapk  *ecdsa.PrivateKey
	password []byte // password for hmac
	hasKey   bool   // must be set to true if one of the three above keys has been set

	lastError error // last error assigned via an OptionFunc

	// Expire defines the duration when the token is about to expire
	Expire time.Duration
	// SigningMethod how to sign the JWT. For default value see the OptionFuncs
	SigningMethod jwt.SigningMethod
	// EnableJTI activates the (JWT ID) Claim, a unique identifier. UUID.
	EnableJTI bool
	// JTI represents the interface to generate a new UUID
	JTI interface {
		Get() string
	}

	// Blacklist a backend storage to handle blocked tokens.
	// Default black hole storage. Must be thread safe.
	Blacklist interface {
		Set(token string, expires time.Duration) error
		Has(token string) bool
	}
	// HTTPErrorHandler defines your specific handler when the token is invalid.
	// Default handler nil and a status StatusUnauthorized will be provided
	HTTPErrorHandler func(error) http.Handler

	// PostFormVarPrefix defines the prefix for the form values when the toke parts will
	// be appended to the *http.Request.Form map. Default in constant PostFormVarPrefix
	PostFormVarPrefix string
}

// New create a new manager. If key option will not be
// passed then a HMAC password will be generated.
// Default expire is one hour. Default signing method is HMAC512. The auto
// generated password will not be output.
func New(opts ...OptionFunc) (*AuthManager, error) {
	a := new(AuthManager)
	for _, opt := range opts {
		opt(a)
	}
	if a.lastError != nil {
		return nil, a.lastError
	}
	if !a.hasKey {
		setAuthManagerDefaults(a)
	}
	if a.lastError != nil {
		return nil, a.lastError
	}
	if a.Expire.Seconds() < 1 {
		a.Expire = time.Hour
	}
	if a.Blacklist == nil {
		a.Blacklist = nullBL{}
	}
	if a.PostFormVarPrefix == "" {
		a.PostFormVarPrefix = PostFormVarPrefix
	}
	if a.JTI == nil {
		a.JTI = jti{}
	}
	if a.HTTPErrorHandler == nil {
		a.HTTPErrorHandler = func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
			})
		}
	}
	return a, nil
}

// GenerateToken creates a new JSON web token. The claims argument will be
// assigned after the registered claim names exp and iat have been set.
// If EnableJTI is false the returned argument jti is empty.
// For details of the registered claim names please see
// http://self-issued.info/docs/draft-ietf-oauth-json-web-token.html#rfc.section.4.1
func (a *AuthManager) GenerateToken(claims map[string]interface{}) (token, jti string, err error) {
	now := time.Now()
	t := jwt.New(a.SigningMethod)
	t.Claims["exp"] = now.Add(a.Expire).Unix()
	t.Claims["iat"] = now.Unix()
	for k, v := range claims {
		t.Claims[k] = v
	}
	if a.EnableJTI && a.JTI != nil {
		jti = a.JTI.Get()
		t.Claims["jti"] = jti
	}

	switch t.Method.Alg() {
	case jwt.SigningMethodRS256.Alg(), jwt.SigningMethodRS384.Alg(), jwt.SigningMethodRS512.Alg():
		token, err = t.SignedString(a.rsapk)
	case jwt.SigningMethodES256.Alg(), jwt.SigningMethodES384.Alg(), jwt.SigningMethodES512.Alg():
		token, err = t.SignedString(a.ecdsapk)
	case jwt.SigningMethodHS256.Alg(), jwt.SigningMethodHS384.Alg(), jwt.SigningMethodHS512.Alg():
		token, err = t.SignedString(a.password)
	default:
		return "", "", errgo.Newf("GenerateToken: Unknown algorithm %s", t.Method.Alg())
	}

	return
}

// Logout adds a token securely to a blacklist with the expiration duration
func (a *AuthManager) Logout(token *jwt.Token) error {
	if token == nil || token.Raw == "" || token.Valid == false {
		return nil
	}

	var exp time.Duration
	if cexp, ok := token.Claims["exp"]; ok {
		if fexp, ok := cexp.(float64); ok {
			tm := time.Unix(int64(fexp), 0)
			if remainer := tm.Sub(time.Now()); remainer > 0 {
				exp = remainer
			}
		}
	}

	return a.Blacklist.Set(token.Raw, exp)
}

// keyFunc runs parallel and concurrent
func (a *AuthManager) keyFunc(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != a.SigningMethod.Alg() {
		return nil, log.Error("userjwt.AuthManager.Authenticate.SigningMethod", "err", ErrUnexpectedSigningMethod, "token", t, "method", a.SigningMethod.Alg())
	} else {
		switch t.Method.Alg() {
		case jwt.SigningMethodRS256.Alg(), jwt.SigningMethodRS384.Alg(), jwt.SigningMethodRS512.Alg():
			return &a.rsapk.PublicKey, nil
		case jwt.SigningMethodES256.Alg(), jwt.SigningMethodES384.Alg(), jwt.SigningMethodES512.Alg():
			return &a.ecdsapk.PublicKey, nil
		case jwt.SigningMethodHS256.Alg(), jwt.SigningMethodHS384.Alg(), jwt.SigningMethodHS512.Alg():
			return a.password, nil
		default:
			return nil, errgo.Newf("Authenticate: Unknown algorithm %s", t.Method.Alg())
		}
	}
}

// Parse parses a token string and returns the valid token or an error
func (a *AuthManager) Parse(rawToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(rawToken, a.keyFunc)
	if token != nil && err == nil && token.Valid && !a.Blacklist.Has(token.Raw) {
		return token, nil
	}
	return nil, errgo.Mask(err)
}

// Authenticate represent a middleware handler for a http router.
// For POST or PUT requests, it also parses the request body as a form and
// put the results into r.Form. The claims of a token will
// be appended to the requests Form map.
func (a *AuthManager) Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := jwt.ParseFromRequest(r, a.keyFunc)

		var inBL bool
		if token != nil {
			inBL = a.Blacklist.Has(token.Raw)
		}
		if token != nil && err == nil && token.Valid && !inBL {
			if err := appendTokenToForm(r, token, a.PostFormVarPrefix); err != nil {
				a.HTTPErrorHandler(
					log.Error("userjwt.AuthManager.Authenticate.appendTokenToForm", "err", err, "r", r, "token", token),
				).ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		} else {
			if log.IsInfo() {
				log.Info("userjwt.AuthManager.Authenticate", "err", err, "token", token, "blacklist", inBL)
			}
			a.HTTPErrorHandler(err).ServeHTTP(w, r) // is that really thread safe or other bug?
		}
	})
}

func appendTokenToForm(r *http.Request, t *jwt.Token, prefix string) error {

	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		if err := r.ParseForm(); err != nil {
			return err
		}
	}

	if r.Form == nil {
		r.Form = make(url.Values)
	}

	for k, v := range t.Claims {
		vs, err := cast.ToStringE(v)
		if err == nil {
			r.Form.Add(prefix+k, vs)
		} else {
			return errgo.Newf("appendTokenToForm: failed to assert to type string: key %s => value %v : %s", k, v, err.Error())
		}
	}

	return nil
}
