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
	"github.com/corestoreio/csfw/sys/singleflight"
	"github.com/corestoreio/csfw/util/errors"
)

// Service represents a service manager for GeoIP detection and restriction.
// Please consider the law in your country if you would like to implement geo-blocking.
type Service struct {
	// Log used for debugging. Defaults to black hole. Panics if nil.
	Log log.Logger

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
	defaultScopeCache scopedConfig

	// rwmu protects all fields below
	rwmu sync.RWMutex

	// geoIP searches the country by an IP address. If nil panics during
	// execution in the middleware.
	geoIP CountryRetriever

	// geoIPLoaded checks to only load the GeoIP CountryRetriever once because
	// we may set that within a request. It's defined in the backend
	// configuration but later we need to reset this value to zero to allow
	// reloading.
	geoIPLoaded *uint32

	// scopeCache internal cache of the configurations. scoped.Hash relates to
	// the default,website or store ID.
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

// isGeoIPLoaded checks if the geoip lookup interface has been set by an object.
// this can be adjusted dynamically with the scoped configuration.
func (s *Service) isGeoIPLoaded() bool {
	return atomic.LoadUint32(s.geoIPLoaded) == 1
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
			s.Log.Debug("geoip.Service.ConfigByScopedGetter.defaultScopeCache", log.Stringer("scope", h), log.Bool("optionFactoryFunc_Nil", s.optionFactoryFunc == nil))
		}
		return s.defaultScopeCache
	}

	sCfg := s.getConfigByScopeID(h, false) // 1. map lookup, but only one lookup during many requests, 99%
	isGeoIPLoaded := s.isGeoIPLoaded()

	switch {
	case sCfg.isValid() == nil && isGeoIPLoaded: // not nice isGeoIPLoaded, find something better
		// cached entry found which can contain the configured scoped configuration or
		// the default scope configuration.
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.ConfigByScopedGetter.IsValid", log.Stringer("scope", h))
		}
		return sCfg
	case s.optionFactoryFunc == nil:
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.ConfigByScopedGetter.optionFactoryFunc.nil", log.Stringer("scope", h))
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
		if err := s.Options(s.optionFactoryFunc(scpGet)...); err != nil {
			return scopedConfig{
				lastErr: errors.Wrap(err, "[geoip] Options by scpOptionFnc"),
			}, nil
		}
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.ConfigByScopedGetter.optionInflight.Do", log.Stringer("scope", h), log.Bool("isGeoIPLoaded_before", isGeoIPLoaded), log.Bool("isGeoIPLoaded_after", s.isGeoIPLoaded()), log.Err(sCfg.isValid()))
		}
		return s.getConfigByScopeID(h, true), nil
	})

	select {
	case res, ok := <-scpCfgChan:
		if !ok {
			return scopedConfig{lastErr: errors.NewFatalf("[geoip] optionInflight.DoChan returned a closed/unreadable channel")}
		}
		if res.Err != nil {
			return scopedConfig{lastErr: errors.Wrap(res.Err, "[geoip] optionInflight.DoChan.Error")}
		}
		sCfg, ok = res.Val.(scopedConfig)
		if !ok {
			sCfg.lastErr = errors.NewFatalf("[geoip] optionInflight.DoChan res.Val cannot be type asserted to scopedConfig")
		}
		return sCfg
	}

	return scopedConfig{
		lastErr: errors.NewFatalf("[geoip] Nothing to select from channel for scope: %q", h),
	}
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
			lastErr: errors.Wrap(errConfigNotFound, "[geoip] Service.getConfigByScopeID.NotFound"),
		}
	}
	if scpCfg.useDefault {
		return s.defaultScopeCache
	}
	return scpCfg
}
