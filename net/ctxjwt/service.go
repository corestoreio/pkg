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
	"net/http"
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
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

	mu sync.RWMutex
	// scopeCache internal cache of already created token configurations
	// scoped.Hash relates to the website ID.
	// this can become a bottle neck when multiple website IDs supplied by a
	// request try to access the map. we can use the same pattern like in freecache
	// to create a segment of 256 slice items to evenly distribute the lock.
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
		return nil, err
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

	s.mu.RLock()
	defer s.mu.RUnlock()
	for h := range s.scopeCache {
		if scp, _ := h.Unpack(); scp > scope.WebsiteID {
			return errors.Errorf("[ctxjwt] Service does not support this: %s. Only default or website are allowed.", h)
		}
	}

	return nil
}

// ClaimStore creates a new struct with preset ExpiresAt, IssuedAt and ID.
func (s *Service) ClaimStore(scp scope.Scope, id int64) (*jwtclaim.Store, error) {
	cfg, err := s.getConfigByScopeID(true, scope.NewHash(scp, id))
	if err != nil {
		return nil, err
	}

	now := csjwt.TimeFunc()
	c := jwtclaim.NewStore()
	c.ExpiresAt = now.Add(cfg.expire).Unix()
	c.IssuedAt = now.Unix()

	if cfg.enableJTI && s.JTI != nil {
		c.ID = s.JTI.Get()
	}

	return c, nil
}

// NewToken creates a new JSON web token based on the Claimer interface and
// depending on the scope and the scoped based configuration. The returned token slice
// is owned by the caller. ExpiresAt, IssuedAt and ID are already set and cannot
// be overwritten.
func (s *Service) NewToken(scp scope.Scope, id int64, claims csjwt.Claimer) (token text.Chars, err error) {

	cfg, err := s.getConfigByScopeID(true, scope.NewHash(scp, id))
	if err != nil {
		return nil, err
	}

	now := csjwt.TimeFunc()

	if err := claims.Set(jwtclaim.KeyExpiresAt, now.Add(cfg.expire).Unix()); err != nil {
		return nil, errors.Mask(err)
	}
	if err := claims.Set(jwtclaim.KeyIssuedAt, now.Unix()); err != nil {
		return nil, errors.Mask(err)
	}

	if cfg.enableJTI && s.JTI != nil {
		if err := claims.Set(jwtclaim.KeyID, s.JTI.Get()); err != nil {
			return nil, errors.Mask(err)
		}
	}

	return csjwt.NewToken(claims).SignedString(cfg.signingMethod, cfg.Key)
}

// Logout adds a token securely to a blacklist with the expiration duration.
func (s *Service) Logout(token csjwt.Token) error {
	if len(token.Raw) == 0 || token.Valid == false {
		return nil
	}

	return s.Blacklist.Set(token.Raw, token.Claims.Expires())
}

// Parse parses a token string with the DefaultID scope and returns the
// valid token or an error.
func (s *Service) Parse(rawToken []byte) (csjwt.Token, error) {
	return s.ParseScoped(scope.DefaultID, 0, rawToken)
}

// ParseScoped parses a token based on the applied scope and the scope ID.
// Different configurations are passed to the token parsing function.
// The black list will be checked for containing entries.
func (s *Service) ParseScoped(scp scope.Scope, id int64, rawToken []byte) (csjwt.Token, error) {
	var emptyTok csjwt.Token
	sc, err := s.getConfigByScopeID(true, scope.NewHash(scp, id))
	if err != nil {
		return emptyTok, err
	}

	token, err := sc.parseWithClaim(rawToken)
	if err != nil {
		return emptyTok, errors.Mask(err)
	}

	var inBL bool
	isValid := token.Valid && len(token.Raw) > 0
	if isValid {
		inBL = s.Blacklist.Has(token.Raw)
	}
	if isValid && !inBL {
		return token, nil
	}
	if PkgLog.IsDebug() {
		PkgLog.Debug("ctxjwt.Service.Parse", "err", err, "inBlackList", inBL, "rawToken", string(rawToken), "token", token)
	}
	return emptyTok, errors.Mask(err)
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
		// cached entry found and ignore the error because we fall back to
		// default scope at the end of this function.
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
	var empty scopedConfig
	// requested scope plus ID
	scpCfg, ok := s.getScopedConfig(hash)

	if ok {
		if scpCfg.isValid() {
			return scpCfg, nil
		}
		return empty, errors.Errorf("[ctxjwt] Incomplete configuration for %s. Missing Signing Method and its Key.", hash)
	}

	if fallback {
		// fallback to default scope
		scpCfg, ok := s.getScopedConfig(scope.DefaultHash)
		if ok {
			return scpCfg, nil
		}
	}

	// give up, nothing found
	return empty, errors.Errorf("[ctxjwt] Cannot find JWT configuration for %s", hash)
}

// getScopedConfig part of lookupScopedConfig and doesn't use a lock because the lock
// has been acquired in lookupScopedConfig()
func (s *Service) getScopedConfig(h scope.Hash) (sc scopedConfig, ok bool) {
	s.mu.RLock()
	sc, ok = s.scopeCache[h]
	s.mu.RUnlock()

	if ok {
		var hasChanges bool
		if nil == sc.keyFunc {
			sc.keyFunc = getKeyFunc(sc)
			hasChanges = true
		}
		if sc.expire < 1 {
			sc.expire = DefaultExpire
			hasChanges = true
		}
		if hasChanges {
			s.mu.Lock()
			s.scopeCache[h] = sc
			s.mu.Unlock()
		}
	}
	return sc, ok
}
