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
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"time"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/juju/errors"
	"github.com/pborman/uuid"
)

// ErrUnexpectedSigningMethod will be returned if some outside dude tries to trick us
var ErrUnexpectedSigningMethod = errors.New("JWT: Unexpected signing method")

// DefaultExpire duration when a token expires
var DefaultExpire time.Duration = time.Hour

// Blacklister a backend storage to handle blocked tokens.
// Default black hole storage. Must be thread safe.
type Blacklister interface {
	Set(token string, expires time.Duration) error
	Has(token string) bool
}

// Service main object for handling JWT authentication, generation, blacklists and log outs.
type Service struct {
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
	// Blacklist concurrent safe black list service
	Blacklist Blacklister
}

// NewService creates a new token service. If key option will not be
// passed then a HMAC password will be generated.
// Default expire is one hour as in variable DefaultExpire. Default signing
// method is HMAC512. The auto generated password will not be outputted.
func NewService(opts ...Option) (*Service, error) {
	s := new(Service)
	for _, opt := range opts {
		opt(s)
	}
	if s.lastError != nil {
		return nil, s.lastError
	}
	if !s.hasKey {
		s.hasKey = true
		s.SigningMethod = jwt.SigningMethodHS512
		s.password = []byte(uuid.NewRandom()) // @todo can be better ...
	}
	if s.Expire.Seconds() < 1 {
		s.Expire = DefaultExpire
	}
	if s.Blacklist == nil {
		s.Blacklist = nullBL{}
	}
	if s.JTI == nil {
		s.JTI = jti{}
	}
	return s, nil
}

// GenerateToken creates a new JSON web token. The claims argument will be
// assigned after the registered claim names exp and iat have been set.
// If EnableJTI is false the returned argument jti is empty.
// For details of the registered claim names please see
// http://self-issued.info/docs/draft-ietf-oauth-json-web-token.html#rfc.section.4.1
// This function is thread safe.
func (s *Service) GenerateToken(claims map[string]interface{}) (token, jti string, err error) {
	now := time.Now()
	t := jwt.New(s.SigningMethod)
	t.Claims["exp"] = now.Add(s.Expire).Unix()
	t.Claims["iat"] = now.Unix()
	for k, v := range claims {
		t.Claims[k] = v
	}
	if s.EnableJTI && s.JTI != nil {
		jti = s.JTI.Get()
		t.Claims["jti"] = jti
	}

	switch t.Method.Alg() {
	case jwt.SigningMethodRS256.Alg(), jwt.SigningMethodRS384.Alg(), jwt.SigningMethodRS512.Alg():
		token, err = t.SignedString(s.rsapk)
	case jwt.SigningMethodES256.Alg(), jwt.SigningMethodES384.Alg(), jwt.SigningMethodES512.Alg():
		token, err = t.SignedString(s.ecdsapk)
	case jwt.SigningMethodHS256.Alg(), jwt.SigningMethodHS384.Alg(), jwt.SigningMethodHS512.Alg():
		token, err = t.SignedString(s.password)
	default:
		return "", "", fmt.Errorf("GenerateToken: Unknown algorithm %s", t.Method.Alg())
	}

	return
}

// Logout adds a token securely to a blacklist with the expiration duration
func (s *Service) Logout(token *jwt.Token) error {
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

	return s.Blacklist.Set(token.Raw, exp)
}

// keyFunc runs parallel and concurrent
func (s *Service) keyFunc(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != s.SigningMethod.Alg() {
		if PkgLog.IsDebug() {
			PkgLog.Debug("ctxjwt.AuthManager.Authenticate.SigningMethod", "err", ErrUnexpectedSigningMethod, "token", t, "method", s.SigningMethod.Alg())
		}
		return nil, ErrUnexpectedSigningMethod
	}

	switch t.Method.Alg() {
	case jwt.SigningMethodRS256.Alg(), jwt.SigningMethodRS384.Alg(), jwt.SigningMethodRS512.Alg():
		return &s.rsapk.PublicKey, nil
	case jwt.SigningMethodES256.Alg(), jwt.SigningMethodES384.Alg(), jwt.SigningMethodES512.Alg():
		return &s.ecdsapk.PublicKey, nil
	case jwt.SigningMethodHS256.Alg(), jwt.SigningMethodHS384.Alg(), jwt.SigningMethodHS512.Alg():
		return s.password, nil
	default:
		return nil, fmt.Errorf("ctxjwt.Service.keyFunc: Unknown algorithm %s", t.Method.Alg())
	}
}

// Parse parses a token string and returns the valid token or an error
func (s *Service) Parse(rawToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(rawToken, s.keyFunc)
	var inBL bool
	if token != nil {
		inBL = s.Blacklist.Has(token.Raw)
	}
	if token != nil && err == nil && token.Valid && !inBL {
		return token, nil
	}
	if PkgLog.IsDebug() {
		PkgLog.Debug("ctxjwt.Service.Parse", "err", err, "inBlackList", inBL, "rawToken", rawToken, "token", token)
	}
	return nil, errors.Mask(err)
}
