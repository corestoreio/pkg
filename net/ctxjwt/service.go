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
	"fmt"
	"sync"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/dgrijalva/jwt-go"
	"github.com/juju/errors"
)

// ErrUnexpectedSigningMethod will be returned if some outside dude tries to trick us
var ErrUnexpectedSigningMethod = errors.New("JWT: Unexpected signing method")

// Blacklister a backend storage to handle blocked tokens.
// Default black hole storage. Must be thread safe.
type Blacklister interface {
	Set(token string, expires time.Duration) error
	Has(token string) bool
}

// Service main object for handling JWT authentication, generation, blacklists and log outs.
type Service struct {
	*cserr.MultiErr

	mu sync.Mutex
	// scopeCache internal cache of already created token configurations
	// scoped.Hash relates to the website ID.
	scopeCache map[scope.Hash]scopedConfig // see freecache to create high concurrent thru put

	// JTI represents the interface to generate a new UUID aka JWT ID
	JTI interface {
		Get() string
	}
	// Blacklist concurrent safe black list service
	Blacklist Blacklister
	backend   *PkgBackend
}

// NewService creates a new token service. If key option will not be
// passed then a HMAC password will be generated.
// Default expire is one hour as in variable DefaultExpire. Default signing
// method is HMAC512. The auto generated password will not be outputted.
func NewService(opts ...Option) (*Service, error) {
	s := new(Service)

	if err := s.Options(opts...); err != nil {
		return nil, s
	}

	if len(s.scopeCache) == 0 {
		if err := s.Options(WithDefaultConfig(scope.DefaultID, 0)); err != nil {
			return nil, s
		}
	}

	if s.Blacklist == nil {
		s.Blacklist = nullBL{}
	}
	if s.JTI == nil {
		s.JTI = jti{}
	}
	return s, nil
}

// MustNewService same as NewService but panics on error.
func MustNewService(opts ...Option) *Service {
	s, err := NewService(opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// Options applies option at creation time or refreshes them.
func (s *Service) Options(opts ...Option) error {
	for _, opt := range opts {
		opt(s)
	}
	if s.HasErrors() {
		return s
	}
	return nil
}

// GenerateToken creates a new JSON web token. The claims argument will be
// assigned after the registered claim names exp and iat have been set.
// If EnableJTI is false the returned argument jti is empty.
// For details of the registered claim names please see
// http://self-issued.info/docs/draft-ietf-oauth-json-web-token.html#rfc.section.4.1
// This function is thread safe.
func (s *Service) GenerateToken(scp scope.Scope, id int64, claims map[string]interface{}) (token, jti string, err error) {

	cfg, err := s.getConfigByScopeID(scp, id)
	if err != nil {
		return "", "", err
	}

	now := time.Now()
	t := jwt.New(cfg.signingMethod)
	t.Claims["exp"] = now.Add(cfg.expire).Unix()
	t.Claims["iat"] = now.Unix()
	for k, v := range claims {
		t.Claims[k] = v
	}
	if cfg.enableJTI && s.JTI != nil {
		jti = s.JTI.Get()
		t.Claims["jti"] = jti
	}

	switch cfg.signingMethod.Alg() {
	case jwt.SigningMethodRS256.Alg(), jwt.SigningMethodRS384.Alg(), jwt.SigningMethodRS512.Alg():
		token, err = t.SignedString(cfg.rsapk)
	case jwt.SigningMethodES256.Alg(), jwt.SigningMethodES384.Alg(), jwt.SigningMethodES512.Alg():
		token, err = t.SignedString(cfg.ecdsapk)
	case jwt.SigningMethodHS256.Alg(), jwt.SigningMethodHS384.Alg(), jwt.SigningMethodHS512.Alg():
		token, err = t.SignedString(cfg.hmacPassword)
	default:
		return "", "", ErrUnexpectedSigningMethod
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

// Parse parses a token string and returns the valid token or an error
func (s *Service) Parse(rawToken string) (*jwt.Token, error) {
	return s.ParseScoped(scope.DefaultID, 0, rawToken)
}

// ParseScoped parses a token based on the applied scope and the scope ID.
// Different configurations are passed to the token parsing function.
// The black list will be checked for containing entries.
func (s *Service) ParseScoped(scp scope.Scope, id int64, rawToken string) (*jwt.Token, error) {

	sc, err := s.getConfigByScopeID(scp, id)

	token, err := jwt.Parse(rawToken, sc.keyFunc)
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

func (s *Service) getConfigByScopedGetter(sg config.ScopedGetter) (scopedConfig, error) {

	sc, err := s.getConfigByScopeID(sg.Scope())
	if err == nil {
		// cached entry found!
		return sc, nil
	}
	if s.backend == nil {
		return scopedConfig{}, errors.Errorf("[ctxjwt.Service] Backend configuration has not been set")
	}

	if err := s.Options(optionsByBackend(s.backend, sg)[:]...); err != nil {
		return scopedConfig{}, errors.Mask(err)
	}

	// after applying the new config try to fetch the new scoped token configuration
	return s.getConfigByScopeID(sg.Scope())
}

func (s *Service) getConfigByScopeID(scp scope.Scope, id int64) (scopedConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// requested scope plus ID
	if scpJWT, ok := s.getScopedConfig(scope.NewHash(scp, id)); ok {
		return scpJWT, nil
	}

	// fallback to default scope
	if scpJWT, ok := s.getScopedConfig(scope.NewHash(scope.DefaultID, 0)); ok {
		return scpJWT, nil
	}

	// give up, nothing found
	return scopedConfig{}, errors.Errorf("[ctxjwt.Service] Cannot find JWT configuration for Scope(%s) and ID %d", scp, id)
}

// getScopedConfig part of lookupScopedConfig and doesn't use a lock because the lock
// has been acquired in lookupScopedConfig()
func (s *Service) getScopedConfig(h scope.Hash) (sc scopedConfig, ok bool) {
	sc, ok = s.scopeCache[h]
	if ok {
		if nil == sc.keyFunc {
			// set the keyFunc and cache it
			s.scopeCache[h] = keyFunc(sc)
		}
	}
	return sc, ok
}

// keyFunc generates the key function for a specific scope and to used in caching
func keyFunc(scpCfg scopedConfig) scopedConfig {
	scpCfg.keyFunc = func(t *jwt.Token) (interface{}, error) {

		if t.Method.Alg() != scpCfg.signingMethod.Alg() {
			if PkgLog.IsDebug() {
				PkgLog.Debug("ctxjwt.keyFunc.SigningMethod", "err", ErrUnexpectedSigningMethod, "token", t, "method", scpCfg.signingMethod.Alg())
			}
			return nil, ErrUnexpectedSigningMethod
		}

		switch t.Method.Alg() {
		case jwt.SigningMethodRS256.Alg(), jwt.SigningMethodRS384.Alg(), jwt.SigningMethodRS512.Alg():
			return &scpCfg.rsapk.PublicKey, nil
		case jwt.SigningMethodES256.Alg(), jwt.SigningMethodES384.Alg(), jwt.SigningMethodES512.Alg():
			return &scpCfg.ecdsapk.PublicKey, nil
		case jwt.SigningMethodHS256.Alg(), jwt.SigningMethodHS384.Alg(), jwt.SigningMethodHS512.Alg():
			return scpCfg.hmacPassword, nil
		default:
			return nil, ErrUnexpectedSigningMethod
		}
	}
	return scpCfg
}
