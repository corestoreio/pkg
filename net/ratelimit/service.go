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
	"net/http"
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sys/suspend"
	"github.com/corestoreio/csfw/util/errors"
)

// VaryByer is called for each request to generate a key for the
// limiter. If it is nil, all requests use an empty string key.
type VaryByer interface {
	Key(*http.Request) string
}

// HTTPRateLimit faciliates using a Limiter to limit HTTP requests.
type Service struct {
	// Log used for debugging. Defaults to black hole. Panics if nil.
	Log log.Logger

	// optionFactoryFunc optional configuration closure, can be nil. It pulls
	// out the configuration settings during a request and caches the settings
	// in the internal map. ScopedOption requires a config.ScopedGetter. This
	// function gets set via WithOptionFactory()
	optionFactoryFunc OptionFactoryFunc
	// optionFactoryState checks on a per scope.Hash basis if the
	// configuration loading process takes place. Delays the execution of other
	// Goroutines with the same scope.Hash until the configuration has been
	// fully loaded and applied and for that specific scope. This function gets
	// set via WithOptionFactory()
	optionFactoryState suspend.State

	// defaultScopeCache has been extracted from the scopeCache to allow faster
	// access to the standard configuration without accessing a map.
	defaultScopeCache scopedConfig

	// rwmu protects all fields below
	rwmu sync.RWMutex

	// scopeCache internal cache of already created token configurations
	// scoped.Hash relates to the website ID. This can become a bottle neck when
	// multiple website IDs supplied by a request try to access the map. we can
	// use the same pattern like in freecache to create a segment of 256 slice
	// items to evenly distribute the lock.
	scopeCache map[scope.Hash]scopedConfig

	// VaryByer is called for each request to generate a key for the
	// limiter. If it is nil, all requests use an empty string key.
	VaryByer
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
		return nil, errors.Wrap(err, "[shy] Options WithDefaultConfig")
	}
	if err := s.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[shy] Options Any Config")
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
// ScopedGetter. Mainly used within the middleware. A nil argument falls back to
// the default scope configuration. If you have applied the option WithBackend()
// the configuration will be pulled out only one time from the backend service.
func (s *Service) configByScopedGetter(scpGet config.ScopedGetter) scopedConfig {

	h := scope.NewHash(scpGet.Scope())
	// fallback to default scope
	if s.useDefaultConfig(h) {
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.defaultScopeCache", log.Stringer("scope", h), log.Bool("optionFactoryFunc_Nil", s.optionFactoryFunc == nil))
		}
		return s.defaultScopeCache
	}

	sCfg := s.getConfigByScopeID(h, false) // 1. map lookup, but only one lookup during many requests, 99%

	switch {
	case sCfg.isValid() == nil:
		// cached entry found which can contain the configured scoped configuration or
		// the default scope configuration.
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.IsValid", log.Stringer("scope", h))
		}
		return sCfg
	case s.optionFactoryFunc == nil:
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.optionFactoryFunc.nil", log.Stringer("scope", h))
		}
		// When everything has been preconfigured for each scope via functional
		// options then we might have the case where a scope in a request comes
		// in which does not match to a scopedConfiguration entry in the map
		// therefore a fallback to default scope must be provided. This fall
		// back will only be executed once and each scope knows from now on that
		// it has the configuration of the default scope.
		return s.getConfigByScopeID(h, true)
	case s.optionFactoryState.ShouldStart(h.ToUint64()): // 2. map lookup, execute for each scope which needs to initialize the configuration.
		// gets tested by backendshy
		defer s.optionFactoryState.Done(h.ToUint64()) // send Signal and release waiter

		if err := s.Options(s.optionFactoryFunc(scpGet)...); err != nil {
			return scopedConfig{
				lastErr: errors.Wrap(err, "[shy] Options by scpOptionFnc"),
			}
		}
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.ShouldStart", log.Stringer("scope", h), log.Err(sCfg.isValid()))
		}
	case s.optionFactoryState.ShouldWait(h.ToUint64()): // 3. map lookup
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.ShouldWait", log.Stringer("scope", h), log.Err(sCfg.isValid()))
		}
		// gets tested by backendratelimit.
		// Wait here! After optionFactoryState.Done() has been called in
		// optionFactoryState.ShouldStart() we proceed and go to the last return
		// statement to search for the newly set scopedConfig value.
	}

	return s.configByScopedGetter(scpGet)
}

func (s *Service) getConfigByScopeID(hash scope.Hash, useDefault bool) scopedConfig {

	s.rwmu.RLock()
	scpCfg, ok := s.scopeCache[hash]
	s.rwmu.RUnlock()
	if !ok && useDefault {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// all other fields are empty or nil!
		scpCfg.scopeHash = hash
		scpCfg.useDefault = true
		// in very high concurrency this might get executed multiple times but it doesn't
		// matter that much until the entry into the map has been written.
		s.scopeCache[hash] = scpCfg
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.getConfigByScopeID.fallbackToDefault", log.Stringer("scope", hash))
		}
	}
	if !ok && !useDefault {
		return scopedConfig{
			lastErr: errors.Wrap(errConfigNotFound, "[shy] Service.getConfigByScopeID.NotFound"),
		}
	}
	if scpCfg.useDefault {
		return s.defaultScopeCache
	}
	return scpCfg
}
