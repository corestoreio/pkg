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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/juju/errors"
	"github.com/pborman/uuid"
)

// jti type to generate a JTI for a token, a unique ID
type jti struct{}

func (j jti) Get() string {
	return uuid.New()
}

const (
	claimExpiresAt = "exp"
	claimIssuedAt  = "iat"
	claimKeyID     = "jti"
)

// Service main type for handling JWT authentication, generation, blacklists
// and log outs depending on a scope.
type Service struct {
	*cserr.MultiErr

	// JTI represents the interface to generate a new UUID aka JWT ID
	JTI interface {
		Get() string
	}
	// Blacklist concurrent safe black list service which handles blocked tokens.
	// Default black hole storage. Must be thread safe.
	Blacklist Blacklister

	// scpOptionFnc optional configuration closure, can be nil. It pulls
	// out the configuration settings during a request and caches the settings in the
	// internal map. ScopedOption requires a config.ScopedGetter
	scpOptionFnc ScopedOption

	defaultScopeCache scopedConfig

	mu sync.RWMutex
	// scopeCache internal cache of already created token configurations
	// scoped.Hash relates to the website ID.
	// this can become a bottle neck when multiple website IDs supplied by a
	// request try to access the map. we can use the same pattern like in freecache
	// to create a segment of 256 slice items to evenly distribute the lock.
	scopeCache map[scope.Hash]scopedConfig // see freecache to create high concurrent thru put
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
	}

	if err := s.Options(WithDefaultConfig(scope.Default, 0)); err != nil {
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
		if scp, _ := h.Unpack(); scp > scope.Website {
			return errors.Errorf("[ctxjwt] Service does not support this: %s. Only default or website are allowed.", h)
		}
	}

	return nil
}

// NewToken creates a new signed JSON web token based on the predefined scoped
// based template token function (WithTemplateToken) and merges the optional
// 3rd argument into the template token claim. Only one claim argument is supported.
// The returned token is owned by the caller. The tokens Raw field contains the
// freshly signed byte slice. ExpiresAt, IssuedAt and ID are already set and cannot
// be overwritten, but you can access them. It panics if the provided template
// token has a nil Header or Claimer field.
func (s *Service) NewToken(scp scope.Scope, id int64, claim ...csjwt.Claimer) (csjwt.Token, error) {
	now := csjwt.TimeFunc()
	var empty csjwt.Token
	cfg, err := s.getConfigByScopeID(true, scope.NewHash(scp, id))
	if err != nil {
		return empty, errors.Mask(err)
	}

	var tk = cfg.TemplateToken()

	if len(claim) > 0 && claim[0] != nil {
		if err := tk.Merge(claim[0]); err != nil {
			return empty, errors.Mask(err)
		}
	}

	if err := tk.Claims.Set(claimExpiresAt, now.Add(cfg.Expire).Unix()); err != nil {
		return empty, errors.Mask(err)
	}
	if err := tk.Claims.Set(claimIssuedAt, now.Unix()); err != nil {
		return empty, errors.Mask(err)
	}

	if cfg.EnableJTI && s.JTI != nil {
		if err := tk.Claims.Set(claimKeyID, s.JTI.Get()); err != nil {
			return empty, errors.Mask(err)
		}
	}

	tk.Raw, err = tk.SignedString(cfg.SigningMethod, cfg.Key)
	return tk, errors.Mask(err)
}

// Logout adds a token securely to a blacklist with the expiration duration.
func (s *Service) Logout(token csjwt.Token) error {
	if len(token.Raw) == 0 || !token.Valid {
		return nil
	}

	return s.Blacklist.Set(token.Raw, token.Claims.Expires())
}

// Parse parses a token string with the DefaultID scope and returns the
// valid token or an error.
func (s *Service) Parse(rawToken []byte) (csjwt.Token, error) {
	return s.ParseScoped(scope.Default, 0, rawToken)
}

// ParseScoped parses a token based on the applied scope and the scope ID.
// Different configurations are passed to the token parsing function.
// The black list will be checked for containing entries.
func (s *Service) ParseScoped(scp scope.Scope, id int64, rawToken []byte) (csjwt.Token, error) {
	var emptyTok csjwt.Token
	sc, err := s.getConfigByScopeID(true, scope.NewHash(scp, id))
	if err != nil {
		return emptyTok, errors.Mask(err)
	}

	token, err := sc.Parse(rawToken)
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

// ConfigByScopedGetter returns the internal configuration depending on the ScopedGetter.
// Mainly used within the middleware. Exported here to build your own middleware.
// A nil argument falls back to the default scope configuration.
// If you have applied the option WithBackend() the configuration will be pulled out
// one time from the backend service.
func (s *Service) ConfigByScopedGetter(sg config.ScopedGetter) (scopedConfig, error) {

	h := scope.DefaultHash
	if sg != nil {
		h = scope.NewHash(sg.Scope())
	}

	if (s.scpOptionFnc == nil || sg == nil) && h == scope.DefaultHash && s.defaultScopeCache.IsValid() {
		return s.defaultScopeCache, nil
	}

	sc, err := s.getConfigByScopeID(false, h)
	if err == nil {
		// cached entry found and ignore the error because we fall back to
		// default scope at the end of this function.
		return sc, nil
	}

	if s.scpOptionFnc != nil {
		if err := s.Options(s.scpOptionFnc(sg)...); err != nil {
			return scopedConfig{}, errors.Mask(err)
		}
	}

	// after applying the new config try to fetch the new scoped token configuration
	return s.getConfigByScopeID(true, h)
}

// ConfigByScopeID returns the internal configuration depending on the scope.
func (s *Service) ConfigByScopeID(scp scope.Scope, id int64) (scopedConfig, error) {
	return s.getConfigByScopeID(true, scope.NewHash(scp, id))
}

func (s *Service) getConfigByScopeID(fallback bool, hash scope.Hash) (scopedConfig, error) {
	var empty scopedConfig
	// requested scope plus ID
	scpCfg, ok := s.getScopedConfig(hash)
	if ok {
		if scpCfg.IsValid() {
			return scpCfg, nil
		}
		return empty, errors.Errorf("[ctxjwt] Incomplete configuration for %s. Missing Signing Method and its Key.", hash)
	}

	const errConfigNotFound = "[ctxjwt] Cannot find JWT configuration for %s"
	if fallback {
		// fallback to default scope
		var err error
		if !s.defaultScopeCache.IsValid() {
			err = errors.Errorf(errConfigNotFound, scope.DefaultHash)
		}
		return s.defaultScopeCache, errors.Mask(err)

	}

	// give up, nothing found
	return empty, errors.Errorf(errConfigNotFound, hash)
}

// getScopedConfig part of lookupScopedConfig and doesn't use a lock because the lock
// has been acquired in lookupScopedConfig()
func (s *Service) getScopedConfig(h scope.Hash) (sc scopedConfig, ok bool) {
	s.mu.RLock()
	sc, ok = s.scopeCache[h]
	s.mu.RUnlock()

	if ok {
		var hasChanges bool
		if nil == sc.KeyFunc {
			sc.initKeyFunc()
			hasChanges = true
		}
		if sc.Expire < 1 {
			sc.Expire = DefaultExpire
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
