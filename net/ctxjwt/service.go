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
	"sync"
	"time"

	"net/http"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/dgrijalva/jwt-go"
	"github.com/juju/errors"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
)

// jti type to generate a JTI for a token, a unique ID
type jti struct{}

func (j jti) Get() string {
	return uuid.New()
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
	// Backend optional configuration, can be nil.
	Backend *PkgBackend
	// DefaultErrorHandler global default error handler. Used in the middleware. Fallback to
	// this handler when a scoped based handler isn't available.
	DefaultErrorHandler ctxhttp.Handler
}

// NewService creates a new token service. If key option will not be
// passed then a HMAC password will be generated.
// Default expire is one hour as in variable DefaultExpire. Default signing
// method is HMAC512. The auto generated password will not be outputted.
// The DefaultErrorHandler returns a http.StatusUnauthorized.
func NewService(opts ...Option) (*Service, error) {
	s := &Service{
		scopeCache: make(map[scope.Hash]scopedConfig),
		JTI:        jti{},
		Blacklist:  nullBL{},
		DefaultErrorHandler: ctxhttp.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return nil
		}),
	}

	if err := s.Options(WithDefaultConfig(scope.DefaultID, 0)); err != nil {
		return nil, s
	}
	if err := s.Options(opts...); err != nil {
		return nil, s
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

	cfg, err := s.getConfigByScopeID(true, scope.NewHash(scp, id))
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
		fexp := conv.ToFloat64(cexp)
		if fexp > 0.001 {
			tm := time.Unix(int64(fexp), 0)
			if remainer := tm.Sub(time.Now()); remainer > 0 {
				exp = remainer
			}
		}
	}
	return s.Blacklist.Set([]byte(token.Raw), exp)
}

// Parse parses a token string and returns the valid token or an error
func (s *Service) Parse(rawToken string) (*jwt.Token, error) {
	return s.ParseScoped(scope.DefaultID, 0, rawToken)
}

// ParseScoped parses a token based on the applied scope and the scope ID.
// Different configurations are passed to the token parsing function.
// The black list will be checked for containing entries.
func (s *Service) ParseScoped(scp scope.Scope, id int64, rawToken string) (*jwt.Token, error) {

	sc, err := s.getConfigByScopeID(true, scope.NewHash(scp, id))
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(rawToken, sc.keyFunc)
	var inBL bool
	if token != nil {
		inBL = s.Blacklist.Has([]byte(token.Raw))
	}
	if token != nil && err == nil && token.Valid && !inBL {
		return token, nil
	}
	if PkgLog.IsDebug() {
		PkgLog.Debug("ctxjwt.Service.Parse", "err", err, "inBlackList", inBL, "rawToken", rawToken, "token", token)
	}
	return nil, errors.Mask(err)
}

// getConfigByScopedGetter used in the middleware where sg comes from the store.Website.Config
// A nil argument falls back to the default scope configuration.
func (s *Service) getConfigByScopedGetter(sg config.ScopedGetter) (scopedConfig, error) {

	h := scope.DefaultHash
	if sg != nil {
		h = scope.NewHash(sg.Scope())
	}

	sc, err := s.getConfigByScopeID(false, h)
	if err == nil {
		// cached entry found!
		return sc, nil
	}

	if s.Backend != nil {
		if err := s.Options(optionsByBackend(s.Backend, sg)...); err != nil {
			return scopedConfig{}, err
		}
	}

	// after applying the new config try to fetch the new scoped token configuration
	return s.getConfigByScopeID(true, h)
}

func (s *Service) getConfigByScopeID(fallback bool, hash scope.Hash) (scopedConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// requested scope plus ID
	if scpJWT, ok := s.getScopedConfig(hash); ok {
		return scpJWT, nil
	}

	if fallback {
		// fallback to default scope
		if scpJWT, ok := s.getScopedConfig(scope.DefaultHash); ok {
			return scpJWT, nil
		}
	}

	// give up, nothing found
	scp, id := hash.Unpack()
	return scopedConfig{}, errors.Errorf("[ctxjwt.Service] Cannot find JWT configuration for Scope(%s) and ID %d", scp, id)
}

// getScopedConfig part of lookupScopedConfig and doesn't use a lock because the lock
// has been acquired in lookupScopedConfig()
func (s *Service) getScopedConfig(h scope.Hash) (sc scopedConfig, ok bool) {
	sc, ok = s.scopeCache[h]
	if ok {
		if nil == sc.keyFunc {
			// set the keyFunc and cache it
			sc.keyFunc = getKeyFunc(sc)
			s.scopeCache[h] = sc
		}
	}
	return sc, ok
}
