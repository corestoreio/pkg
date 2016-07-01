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

package ratelimit

import (
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
	"github.com/corestoreio/csfw/util/errors"
)

// Service creates a middleware that facilitates using a Limiter to limit HTTP
// requests.
type Service struct {
	// Log used for debugging. Defaults to black hole. Panics if nil.
	Log log.Logger

	// optionFactoryFunc optional configuration closure, can be nil. It pulls
	// out the configuration settings during a request and caches the settings
	// in the internal map. ScopedOption requires a config.ScopedGetter. This
	// function gets set via WithOptionFactory()
	optionFactoryFunc OptionFactoryFunc

	// optionInflight checks on a per scope.Hash basis if the configuration
	// loading process takes place. Stops the execution of other Goroutines (aka
	// incoming requests) with the same scope.Hash until the configuration has
	// been fully loaded and applied and for that specific scope. This function
	// gets set via WithOptionFactory()
	optionInflight *singleflight.Group

	// rwmu protects all fields below
	rwmu sync.RWMutex

	// defaultScopeCache has been extracted from the scopeCache to allow faster
	// access to the standard configuration without accessing a map.
	defaultScopeCache scopedConfig

	// scopeCache internal cache of the configurations. scoped.Hash relates to
	// the default,website or store ID.
	scopeCache map[scope.Hash]scopedConfig
}

// New creates a new rate limit middleware.
//
// Default DeniedHandler returns http.StatusTooManyRequests.
//
// Default RateLimiterFactory is the NewGCRAMemStore(). If *PkgBackend has
// been provided the values from the configration will be taken otherwise
// GCRAMemStore() uses the Default* variables.
func New(opts ...Option) (*Service, error) {
	s := &Service{
		Log:        log.BlackHole{},
		scopeCache: make(map[scope.Hash]scopedConfig),
	}
	if err := s.Options(WithDefaultConfig(scope.Default, 0)); err != nil {
		return nil, errors.Wrap(err, "[ratelimit] Options WithDefaultConfig")
	}
	if err := s.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[ratelimit] Options Any Config")
	}
	return s, nil
}

// MustNew same as New() but panics on error. Use only during app start up process.
func MustNew(opts ...Option) *Service {
	c, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return c
}

// Options applies option at creation time or refreshes them.
func (s *Service) Options(opts ...Option) error {
	for _, opt := range opts {
		if opt != nil { // can be nil because of the backend options where we have an array instead of a slice.
			if err := opt(s); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) useDefaultConfig(h scope.Hash) bool {
	return s.optionFactoryFunc == nil && h == scope.DefaultHash && s.defaultScopeCache.isValid() == nil
}

// configByScopedGetter returns the internal configuration depending on the
// ScopedGetter. Mainly used within the middleware.  If you have applied the
// option WithOptionFactory() the configuration will be pulled out only one time
// from the backend configuration service. The field optionInflight handles the
// guaranteed atomic single loading for each scope.
func (s *Service) configByScopedGetter(scpGet config.ScopedGetter) scopedConfig {

	h := scope.NewHash(scpGet.Scope())
	// fallback to default scope
	if s.useDefaultConfig(h) {
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.defaultScopeCache", log.Stringer("scope", h), log.Bool("optionFactoryFunc_Nil", s.optionFactoryFunc == nil))
		}
		return s.defaultScopeCache
	}
	p := scope.NewHash(scpGet.Parent())

	sCfg := s.getConfigByScopeID(h, p) // 1. map lookup, but only one lookup during many requests, 99%

	switch {
	case sCfg.isValid() == nil:
		// cached entry found which can contain the configured scoped configuration or
		// the default scope configuration.
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.IsValid", log.Stringer("scope", h), log.Stringer("parentScope", p))
		}
		return sCfg
	case s.optionFactoryFunc == nil:
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.optionFactoryFunc.nil", log.Stringer("scope", h), log.Stringer("parentScope", p))
		}
		// When everything has been pre-configured for each scope via functional
		// options then we might have the case where a scope in a request comes
		// in which does not match to a scopedConfiguration entry in the map
		// therefore a fallback to default scope must be provided. This fall
		// back will only be executed once and each scope knows from now on that
		// it has the configuration of the default scope.
		return s.getConfigByScopeID(h, scope.DefaultHash)
	}

	scpCfgChan := s.optionInflight.DoChan(h.String(), func() (interface{}, error) {
		if err := s.Options(s.optionFactoryFunc(scpGet)...); err != nil {
			return scopedConfig{
				lastErr: errors.Wrap(err, "[ratelimit] Options by scpOptionFnc"),
			}, nil
		}
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.optionInflight.Do", log.Stringer("scope", h), log.Stringer("parentScope", p), log.Err(sCfg.isValid()))
		}
		return s.getConfigByScopeID(h, p), nil
	})

	res, ok := <-scpCfgChan
	if !ok {
		return scopedConfig{lastErr: errors.NewFatalf("[ratelimit] optionInflight.DoChan returned a closed/unreadable channel")}
	}
	if res.Err != nil {
		return scopedConfig{lastErr: errors.Wrap(res.Err, "[ratelimit] optionInflight.DoChan.Error")}
	}
	sCfg, ok = res.Val.(scopedConfig)
	if !ok {
		sCfg.lastErr = errors.NewFatalf("[ratelimit] optionInflight.DoChan res.Val cannot be type asserted to scopedConfig")
	}
	return sCfg
}

func (s *Service) getScpCfg(hash scope.Hash, parent scope.Hash) scopedConfig {
	s.rwmu.RLock()
	scpCfg, ok := s.scopeCache[hash]
	s.rwmu.RUnlock()

	if !ok && parent.EqualScope(scope.DefaultHash) {
		println("176 DEFAULT CURRENT:", hash.String(), " PARENT:", parent.String(), " VALID:", scpCfg.isValid().Error())

		s.rwmu.Lock()
		scpCfg.fallBackScopeHash = scope.DefaultHash
		s.scopeCache[hash] = scpCfg
		s.rwmu.Unlock()
		return s.defaultScopeCache
	}
	if ok && scpCfg.isValid() == nil {
		return scpCfg
	}
	if ok && scpCfg.fallBackScopeHash.EqualScope(scope.DefaultHash) {
		println("188 DEFAULT CURRENT:", hash.String(), " PARENT:", parent.String(), " VALID:", scpCfg.isValid(), "scpCfg.scope", scpCfg.scopeHash.String())
		return s.defaultScopeCache
	}

	return scpCfg
}

func (s *Service) getConfigByScopeID(hash scope.Hash, parent scope.Hash) scopedConfig {

	scpCfg := s.getScpCfg(hash, parent)

	if scpCfg.fallBackScopeHash.EqualScope(scope.DefaultHash) {
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.getConfigByScopeID.Hash.FallBackDefault",
				log.Stringer("scope", hash),
				log.Stringer("scope_fallback", scpCfg.fallBackScopeHash),
				log.Stringer("scope_parent", parent),
			)
		}
		return scpCfg
	}

	// check website scope
	if scpCfg.fallBackScopeHash == 0 && parent.Scope() == scope.Website {
		scpCfg = s.getScpCfg(parent, 0)

		if err := scpCfg.isValid(); err == nil {
			// we found an entry for a website config
			s.rwmu.Lock()
			dummy := scopedConfig{
				scopeHash:         hash,
				fallBackScopeHash: parent,
			}
			s.scopeCache[hash] = dummy
			s.rwmu.Unlock()
			if s.Log.IsDebug() {
				s.Log.Debug("ratelimit.Service.getConfigByScopeID.Parent.Valid",
					log.Stringer("scope", hash),
					log.Stringer("scope_parent", parent),
				)
			}
		} else if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.getConfigByScopeID.Parent.Invalid",
				log.Stringer("scope", hash),
				log.Stringer("scope_parent", parent),
				log.Err(err),
			)
		}
		// return website config
	}

	if s.Log.IsDebug() {
		s.Log.Debug("ratelimit.Service.getConfigByScopeID.Return",
			log.Stringer("scope", hash),
			log.Stringer("scope_fallback", scpCfg.fallBackScopeHash),
			log.Stringer("scope_parent", parent),
			log.ErrWithKey("scpCfg_is_valid", scpCfg.isValid()),
		)
	}

	return scpCfg
}
