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

package geoip

import (
	"sync"
	"sync/atomic"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// Service represents a service manager for GeoIP detection.
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
	optionFactoryState scope.HashState

	// defaultScopeCache has been extracted from the scopeCache to allow faster
	// access to the standard configuration without accessing a map.
	defaultScopeCache scopedConfig

	// rwmu protects all fields below
	rwmu sync.RWMutex

	// geoIP searches the country for an IP address. If nil panics
	// during execution in the middleware
	geoIP CountryRetriever
	// geoIPLoaded checks to only load the GeoIP CountryRetriever once
	// because we may set that within a request because it is defined
	// in the backend configuration but later we need to reset
	// this value to zero to allow reloading.
	geoIPLoaded *uint32

	// scopeCache internal cache of already created token configurations
	// scoped.Hash relates to the website ID. This can become a bottle neck when
	// multiple website IDs supplied by a request try to access the map. we can
	// use the same pattern like in freecache to create a segment of 256 slice
	// items to evenly distribute the lock.
	scopeCache map[scope.Hash]scopedConfig
}

// NewService creates a new GeoIP service to be used as a middleware.
func New(opts ...Option) (*Service, error) {
	s := &Service{
		geoIPLoaded: new(uint32),
		scopeCache:  make(map[scope.Hash]scopedConfig),
		Log:         log.BlackHole{},
	}
	if err := s.Options(WithDefaultConfig(scope.Default, 0)); err != nil {
		return nil, errors.Wrap(err, "[geoip] Options WithDefaultConfig")
	}
	if err := s.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[geoip] Options Any Config")
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

// Closes the underlying GeoIP CountryRetriever service and resets the internal loading state.
func (s *Service) Close() error {
	atomic.StoreUint32(s.geoIPLoaded, 0)
	return s.geoIP.Close()
}

func (s *Service) isGeoIPLoaded() bool {
	return atomic.LoadUint32(s.geoIPLoaded) == 1
}

func (s *Service) useDefaultConfig(h scope.Hash) bool {
	return s.optionFactoryFunc == nil && h == scope.DefaultHash && s.defaultScopeCache.isValid() == nil
}

// configByScopedGetter returns the internal configuration depending on the ScopedGetter.
// Mainly used within the middleware. Exported here to build your own middleware.
// A nil argument falls back to the default scope configuration.
// If you have applied the option WithBackend() the configuration will be pulled out
// one time from the backend service.
func (s *Service) configByScopedGetter(scpGet config.ScopedGetter) scopedConfig {

	h := scope.NewHash(scpGet.Scope())
	// fallback to default scope
	if s.useDefaultConfig(h) {
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.ConfigByScopedGetter.defaultScopeCache", log.Stringer("scope", h), log.Bool("optionFactoryFunc_Nil", s.optionFactoryFunc == nil))
		}
		return s.defaultScopeCache
	}

	sCfg := s.getConfigByScopeID(h, false) // 1. map lookup, but only one lookup during many requests
	switch {
	case sCfg.isValid() == nil && s.isGeoIPLoaded(): // not nice isGeoIPLoaded, find something better
		// cached entry found which can contain the configured scoped configuration or
		// the default scope configuration.
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.ConfigByScopedGetter.IsValid", log.Stringer("scope", h))
		}
		return sCfg
	case s.optionFactoryState.ShouldStart(h): // 2. map lookup, execute for each scope which needs to initialize the configuration.
		// gets tested by backendgeoip
		defer s.optionFactoryState.Done(h) // send Signal and release waiter

		if err := s.Options(s.optionFactoryFunc(scpGet)...); err != nil {
			return scopedConfig{
				lastErr: errors.Wrap(err, "[geoip] Options by scpOptionFnc"),
			}
		}
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.ConfigByScopedGetter.ShouldStart", log.Stringer("scope", h))
		}
	case s.optionFactoryState.ShouldWait(h): // 3. map lookup
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.ConfigByScopedGetter.ShouldWait", log.Stringer("scope", h))
		}
		// gets tested by backendgeoip
		// Wait! After optionFactoryState.Done() has been called in optionFactoryState.Done() we
		// proceed here and go to the last return statement to search for the
		// newly set scopedConfig value.
	}

	// after applying the new config try to fetch the new scoped configuration
	// or if not found fall back to the default scope.
	return s.getConfigByScopeID(h, true) // 4. map lookup
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
			s.Log.Debug("geoip.Service.getConfigByScopeID.fallbackToDefault", log.Stringer("scope", hash))
		}
	}
	if !ok && !useDefault {
		return scopedConfig{
			lastErr: errors.Wrap(errConfigNotFound, "[geoip] Service.getConfigByScopeID"),
		}
	}
	if scpCfg.useDefault {
		return s.defaultScopeCache
	}
	return scpCfg
}
