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

package jwtauth

import (
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
)

const (
	claimExpiresAt = "exp"
	claimIssuedAt  = "iat"
	claimKeyID     = "jti"
)

// Service main type for handling JWT authentication, generation, blacklists
// and log outs depending on a scope.
type Service struct {
	// JTI represents the interface to generate a new UUID aka JWT ID
	JTI interface {
		Get() string
	}
	// Blacklist concurrent safe black list service which handles blocked tokens.
	// Default black hole storage. Must be thread safe.
	Blacklist Blacklister
	// Log mostly used for debugging. todo(CS) add more logging at useful places
	Log log.Logger

	// optionError used by functional option arguments to indicate that one
	// option has triggered an error and hence the other options can
	// skip their process.
	optionError error

	// scpOptionFnc optional configuration closure, can be nil. It pulls
	// out the configuration settings during a request and caches the settings in the
	// internal map. ScopedOption requires a config.ScopedGetter
	scpOptionFnc ScopedOptionFunc

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
		Log:        log.BlackHole{}, // disabled debug and info logging
	}

	if err := s.Options(WithDefaultConfig(scope.Default, 0)); err != nil {
		return nil, errors.Wrap(err, "[jwtauth] Options WithDefaultConfig")
	}
	if err := s.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[jwtauth] Options Any Config")
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
	if s.optionError != nil {
		return s.optionError
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for h := range s.scopeCache {
		if scp, _ := h.Unpack(); scp > scope.Website {
			return errors.NewNotSupportedf(errServiceUnsupportedScope, h)
		}
	}

	return nil
}

// AddError used by functional options to set an error. The error will only be
// then set if there is not yet an error otherwise it gets discarded. You can
// enable debug logging to find out more.
func (s *Service) AddError(err error) {
	if s.optionError != nil {
		if s.Log.IsDebug() {
			s.Log.Debug("jwtauth.Service.AddError", "err", err, "skipped", true, "currentError", s.optionError)
		}
		return
	}
	s.optionError = err
}

// NewToken creates a new signed JSON web token based on the predefined scoped
// based template token function (WithTemplateToken) and merges the optional
// 3rd argument into the template token claim.
// The returned token is owned by the caller. The tokens Raw field contains the
// freshly signed byte slice. ExpiresAt, IssuedAt and ID are already set and cannot
// be overwritten, but you can access them. It panics if the provided template
// token has a nil Header or Claimer field.
func (s *Service) NewToken(scp scope.Scope, id int64, claim ...csjwt.Claimer) (csjwt.Token, error) {
	now := csjwt.TimeFunc()
	var empty csjwt.Token
	cfg, err := s.getConfigByScopeID(true, scope.NewHash(scp, id))
	if err != nil {
		return empty, errors.Wrap(err, "[jwtauth] getConfigByScopeID")
	}

	var tk = cfg.TemplateToken()

	if len(claim) > 0 && claim[0] != nil {
		if err := csjwt.MergeClaims(tk.Claims, claim...); err != nil {
			return empty, errors.Wrap(err, "[jwtauth] MergeClaims")
		}
	}

	if err := tk.Claims.Set(claimExpiresAt, now.Add(cfg.Expire).Unix()); err != nil {
		return empty, errors.Wrap(err, "[jwtauth] Claims.Set EXP")
	}
	if err := tk.Claims.Set(claimIssuedAt, now.Unix()); err != nil {
		return empty, errors.Wrap(err, "[jwtauth] Claims.Set IAT")
	}

	if cfg.EnableJTI && s.JTI != nil {
		if err := tk.Claims.Set(claimKeyID, s.JTI.Get()); err != nil {
			return empty, errors.Wrap(err, "[jwtauth] Claims.Set KID")
		}
	}

	tk.Raw, err = tk.SignedString(cfg.SigningMethod, cfg.Key)
	return tk, errors.Wrap(err, "[jwtauth] SignedString")
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
		return emptyTok, errors.Wrap(err, "[jwtauth] getConfigByScopeID")
	}

	token, err := sc.Parse(rawToken)
	if err != nil {
		return emptyTok, errors.Wrap(err, "[jwtauth] Parse")
	}

	var inBL bool
	isValid := token.Valid && len(token.Raw) > 0
	if isValid {
		inBL = s.Blacklist.Has(token.Raw)
	}
	if isValid && !inBL {
		return token, nil
	}
	if s.Log.IsDebug() {
		s.Log.Debug("jwtauth.Service.Parse", "err", err, "inBlackList", inBL, "rawToken", string(rawToken), "token", token)
	}
	return emptyTok, errors.NewNotValidf(errTokenParseNotValidOrBlackListed)
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
			return scopedConfig{}, errors.Wrap(err, "[jwtauth] Options by scpOptionFnc")
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
		return empty, errors.NewNotValidf(errScopedConfigMissingSigningMethod, hash)
	}

	if fallback {
		// fallback to default scope
		var err error
		if !s.defaultScopeCache.IsValid() {
			err = errors.NewNotFoundf(errConfigNotFound, scope.DefaultHash)
		}
		return s.defaultScopeCache, err

	}

	// give up, nothing found
	return empty, errors.NewNotFoundf(errConfigNotFound, hash)
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
		if hasChanges {
			s.mu.Lock()
			s.scopeCache[h] = sc
			s.mu.Unlock()
		}
	}
	return sc, ok
}
