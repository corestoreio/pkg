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

package jwt

import (
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
)

const (
	claimExpiresAt = "exp"
	claimIssuedAt  = "iat"
	claimKeyID     = "jti"
)

// Service main type for handling JWT authentication, generation, blacklists and
// log outs depending on a scope.
type Service struct {
	// JTI represents the interface to generate a new UUID aka JWT ID
	JTI interface {
		Get() string
	}

	// Blacklist concurrent safe black list service which handles blocked
	// tokens. Default black hole storage. Must be thread safe.
	Blacklist Blacklister

	// Log mostly used for debugging.
	Log log.Logger

	// StoreService used in the middleware to set a new requested store, change
	// store. If nil the requested store extracted from the context won't be
	// changed.
	StoreService store.Requester

	rootConfig config.Getter

	// optionFactoryFunc optional configuration closure, can be nil. It pulls
	// out the configuration settings during a request and caches the settings
	// in the internal map scopeCache. This function gets set via
	// WithOptionFactory()
	optionFactoryFunc OptionFactoryFunc

	// optionInflight checks on a per scope.Hash basis if the configuration
	// loading process takes place. Stops the execution of other Goroutines (aka
	// incoming requests) with the same scope.Hash until the configuration has
	// been fully loaded and applied and for that specific scope. This function
	// gets set via WithOptionFactory()
	optionInflight *singleflight.Group

	// defaultScopeCache has been extracted from the scopeCache to allow faster
	// access to the standard configuration without accessing a map.
	defaultScopeCache ScopedConfig

	// rwmu protects all fields below
	rwmu sync.RWMutex

	// scopeCache internal cache of already created token configurations
	// scoped.Hash relates to the website ID. this can become a bottle neck when
	// multiple website IDs supplied by a request try to access the map.
	scopeCache map[scope.Hash]ScopedConfig
}

// NewService creates a new token service.
// Default values from option function WithDefaultConfig() will be
// applied if passing no functional option arguments.
func NewService(opts ...Option) (*Service, error) {
	s := &Service{
		JTI:        jti{},
		Blacklist:  nullBL{},
		Log:        log.BlackHole{}, // disabled debug and info logging
		scopeCache: make(map[scope.Hash]ScopedConfig),
	}

	if err := s.Options(WithDefaultConfig(scope.Default, 0)); err != nil {
		return nil, errors.Wrap(err, "[jwt] Options WithDefaultConfig")
	}
	if err := s.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[jwt] Options Any Config")
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
		if opt != nil { // might be nil because of package backendjwt
			if err := opt(s); err != nil {
				return errors.Wrap(err, "[jwt] Service.Options")
			}
		}
	}

	s.rwmu.RLock()
	defer s.rwmu.RUnlock()
	for h := range s.scopeCache {
		// This one checks if the configuraion contains only the default or
		// website scope. Store scope is neither allowed nor supported.
		if scp, _ := h.Unpack(); scp > scope.Website {
			return errors.NewNotSupportedf(errServiceUnsupportedScope, h)
		}
	}
	return nil
}

// NewToken creates a new signed JSON web token based on the predefined scoped
// based template token function (WithTemplateToken) and merges the optional 3rd
// argument into the template token claim. The returned token is owned by the
// caller. The tokens Raw field contains the freshly signed byte slice.
// ExpiresAt, IssuedAt and ID are already set and cannot be overwritten, but you
// can access them. It panics if the provided template token has a nil Header or
// Claimer field.
func (s *Service) NewToken(scp scope.Scope, id int64, claim ...csjwt.Claimer) (csjwt.Token, error) {
	var empty csjwt.Token
	now := csjwt.TimeFunc()

	sc := s.ConfigByScopeID(scp, id)
	if err := sc.IsValid(); err != nil {
		return empty, errors.Wrap(err, "[jwt] NewToken.ConfigByScopeID")
	}

	var tk = sc.TemplateToken()

	if len(claim) > 0 && claim[0] != nil {
		if err := csjwt.MergeClaims(tk.Claims, claim...); err != nil {
			return empty, errors.Wrap(err, "[jwt] NewToken.MergeClaims")
		}
	}

	if err := tk.Claims.Set(claimExpiresAt, now.Add(sc.Expire).Unix()); err != nil {
		return empty, errors.Wrap(err, "[jwt] NewToken.Claims.Set EXP")
	}
	if err := tk.Claims.Set(claimIssuedAt, now.Unix()); err != nil {
		return empty, errors.Wrap(err, "[jwt] NewToken.Claims.Set IAT")
	}

	if sc.EnableJTI && s.JTI != nil {
		if err := tk.Claims.Set(claimKeyID, s.JTI.Get()); err != nil {
			return empty, errors.Wrap(err, "[jwt] NewToken.Claims.Set KID")
		}
	}
	var err error
	tk.Raw, err = tk.SignedString(sc.SigningMethod, sc.Key)
	return tk, errors.Wrap(err, "[jwt] NewToken.SignedString")
}

// Logout adds a token securely to a blacklist with the expiration duration.
func (s *Service) Logout(token csjwt.Token) error {
	if len(token.Raw) == 0 || !token.Valid {
		return nil
	}
	return errors.Wrap(s.Blacklist.Set(token.Raw, token.Claims.Expires()), "[jwt] Service.Logout.Blacklist.Set")
}

// Parse parses a token string with the DefaultID scope and returns the
// valid token or an error.
func (s *Service) Parse(rawToken []byte) (csjwt.Token, error) {
	return s.ParseScoped(scope.Default, 0, rawToken)
}

// ParseScoped parses a token based on the applied scope and the scope ID.
// Different configurations are passed to the token parsing function. The black
// list will be checked for containing entries.
func (s *Service) ParseScoped(scp scope.Scope, id int64, rawToken []byte) (csjwt.Token, error) {
	var empty csjwt.Token

	sc := s.ConfigByScopeID(scp, id)
	if err := sc.IsValid(); err != nil {
		return empty, errors.Wrap(err, "[jwt] ParseScoped.ConfigByScopeID")
	}

	token, err := sc.Parse(rawToken)
	if err != nil {
		return empty, errors.Wrap(err, "[jwt] ParseScoped.Parse")
	}

	// todo simplify
	var inBL bool
	isValid := token.Valid && len(token.Raw) > 0
	if isValid {
		inBL = s.Blacklist.Has(token.Raw)
	}
	if isValid && !inBL {
		return token, nil
	}
	if s.Log.IsDebug() {
		s.Log.Debug("jwt.Service.ParseScoped", log.Err(err), log.Bool("inBlackList", inBL), log.String("rawToken", string(rawToken)), log.Marshal("token", token))
	}
	return empty, errors.NewNotValidf(errTokenParseNotValidOrBlackListed)
}

func (s *Service) useDefaultConfig(h scope.Hash) bool {
	return s.optionFactoryFunc == nil && h == scope.DefaultHash && s.defaultScopeCache.IsValid() == nil
}

// ConfigByScopedGetter returns the internal configuration depending on the
// ScopedGetter. Mainly used within the middleware. Exported here to build your
// own middleware. If you have applied the option WithOptionFactory() the
// configuration will be pulled out one time from the backend service.
func (s *Service) ConfigByScopedGetter(sg config.Scoped) ScopedConfig {
	h := scope.NewHash(sg.Scope())
	// fallback to default scope
	if s.useDefaultConfig(h) {
		if s.Log.IsDebug() {
			s.Log.Debug("jwt.Service.ConfigByScopedGetter.defaultScopeCache", log.Stringer("scope", h), log.Bool("optionFactoryFunc_Nil", s.optionFactoryFunc == nil))
		}
		return s.defaultScopeCache
	}

	sCfg := s.getConfigByScopeID(h, false) // 1. map lookup, but only one lookup during many requests, 99%
	switch {
	case sCfg.IsValid() == nil:
		// cached entry found which can contain the configured scoped configuration or
		// the default scope configuration.
		if s.Log.IsDebug() {
			s.Log.Debug("jwt.Service.ConfigByScopedGetter.IsValid", log.Stringer("scope", h))
		}
		return sCfg
	case s.optionFactoryFunc == nil:
		if s.Log.IsDebug() {
			s.Log.Debug("jwt.Service.ConfigByScopedGetter.optionFactoryFunc.nil", log.Stringer("scope", h))
		}
		// When everything has been pre-configured for each scope via functional
		// options then we might have the case where a scope in a request comes
		// in which does not match to a scopedConfiguration entry in the map
		// therefore a fallback to default scope must be provided. This fall
		// back will only be executed once and each scope knows from now on that
		// it has the configuration of the default scope.
		return s.getConfigByScopeID(h, true)
	}

	scpCfgChan := s.optionInflight.DoChan(h.String(), func() (interface{}, error) {
		if err := s.Options(s.optionFactoryFunc(sg)...); err != nil {
			return ScopedConfig{
				lastErr: errors.Wrap(err, "[jwt] Options by scpOptionFnc"),
			}, nil
		}
		if s.Log.IsDebug() {
			s.Log.Debug("jwt.Service.ConfigByScopedGetter.optionInflight.Do", log.Stringer("scope", h), log.Err(sCfg.IsValid()))
		}
		return s.getConfigByScopeID(h, true), nil
	})

	res, ok := <-scpCfgChan
	if !ok {
		return ScopedConfig{lastErr: errors.NewFatalf("[jwt] optionInflight.DoChan returned a closed/unreadable channel")}
	}
	if res.Err != nil {
		return ScopedConfig{lastErr: errors.Wrap(res.Err, "[jwt] optionInflight.DoChan.Error")}
	}
	sCfg, ok = res.Val.(ScopedConfig)
	if !ok {
		sCfg.lastErr = errors.NewFatalf("[jwt] optionInflight.DoChan res.Val cannot be type asserted to scopedConfig")
	}
	return sCfg
}

// ConfigByScopeID returns the internal configuration depending on the scope.
// The following pitfall applies if the option function WithOptionFactory() has
// been used together with the rootConfig (config.Getter): If the scope.Scope is
// equal to scope.Website then the configuration will be pulled out from
// function ConfigByScopedGetter() with the website ID as second argument of
// this function and store ID zero.
//
// If the option function WithOptionFactory() has not been used then the
// configuration will be searched directly in the internal map.
func (s *Service) ConfigByScopeID(scp scope.Scope, id int64) ScopedConfig {
	if s.Log.IsDebug() {
		s.Log.Debug("jwt.Service.ConfigByScopeID", log.Stringer("scope", scope.NewHash(scp, id)), log.Bool("rootConfig_isNil", s.rootConfig == nil))
	}
	if s.useDefaultConfig(scope.NewHash(scp, id)) {
		return s.defaultScopeCache
	}

	if s.rootConfig != nil && scp > scope.Absent && scp < scope.Group {
		// do not forget: there is an automatic fall back to the default scope
		// IF the ScopedGetter cannot find a configuration value for the website
		// scope.
		return s.ConfigByScopedGetter(s.rootConfig.NewScoped(id, 0))
	}
	return s.getConfigByScopeID(scope.NewHash(scp, id), false)
}

func (s *Service) getConfigByScopeID(hash scope.Hash, useDefault bool) ScopedConfig {

	s.rwmu.RLock()
	scpCfg, ok := s.scopeCache[hash]
	s.rwmu.RUnlock()
	if !ok && useDefault {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// all other fields are empty or nil!
		scpCfg.ScopeHash = hash
		scpCfg.UseDefault = true
		// in very high concurrency this might get executed multiple times but
		// it doesn't matter that much until the entry into the map has been
		// written.
		s.scopeCache[hash] = scpCfg
		if s.Log.IsDebug() {
			s.Log.Debug("jwt.Service.getConfigByScopeID.fallbackToDefault", log.Stringer("scope", hash))
		}
	}
	if !ok && !useDefault {
		return ScopedConfig{
			lastErr: errors.Wrap(errConfigNotFound, "[jwt] Service.getConfigByScopeID.NotFound"),
		}
	}
	if scpCfg.UseDefault {
		return s.defaultScopeCache
	}
	return scpCfg
}
