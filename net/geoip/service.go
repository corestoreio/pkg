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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
)

// Service represents a service manager
type Service struct {
	// GeoIP searches the country for an IP address. If nil panics
	// during execution in the middleware
	GeoIP CountryRetriever
	// Log used for debugging. Defaults to black hole. Panics if nil.
	Log log.Logger

	// scopedOptionFunc optional configuration closure, can be nil. It pulls
	// out the configuration settings during a request and caches the settings in the
	// internal map. ScopedOption requires a config.ScopedGetter
	scopedOptionFunc ScopedOptionFunc

	defaultScopeCache scopedConfig

	mu sync.RWMutex
	// scopeCache internal cache of already created token configurations
	// scoped.Hash relates to the website ID.
	// this can become a bottle neck when multiple website IDs supplied by a
	// request try to access the map. we can use the same pattern like in freecache
	// to create a segment of 256 slice items to evenly distribute the lock.
	scopeCache map[scope.Hash]scopedConfig // see freecache to create high concurrent thru put
}

// NewService creates a new GeoIP service to be used as a middleware.
func New(opts ...Option) (*Service, error) {
	s := &Service{
		scopeCache: make(map[scope.Hash]scopedConfig),
		Log:        log.BlackHole{},
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

// configByScopedGetter returns the internal configuration depending on the ScopedGetter.
// Mainly used within the middleware. Exported here to build your own middleware.
// A nil argument falls back to the default scope configuration.
// If you have applied the option WithBackend() the configuration will be pulled out
// one time from the backend service.
func (s *Service) configByScopedGetter(scpcfg config.ScopedGetter) (scopedConfig, error) {

	h := scope.DefaultHash
	if scpcfg != nil {
		h = scope.NewHash(scpcfg.Scope())
	}
	if s.Log.IsDebug() {
		s.Log.Debug("geoip.Service.ConfigByScopedGetter.ScopedGetter", "ScopedGetter_Nil", scpcfg == nil, "scope", h.String())
	}

	// fallback to default scope
	if (s.scopedOptionFunc == nil || scpcfg == nil) && h == scope.DefaultHash && s.defaultScopeCache.isValid() {
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.ConfigByScopedGetter.defaultScopeCache", "ScopedGetter_Nil", scpcfg == nil, "scpOptionFnc_Nil", s.scopedOptionFunc == nil)
		}
		return s.defaultScopeCache, nil
	}

	sc, err := s.getConfigByScopeID(false, h)
	if err == nil {
		// cached entry found and ignore the error because we fall back to
		// default scope at the end of this function.
		return sc, nil
	}

	if s.scopedOptionFunc != nil {
		if s.Log.IsDebug() {
			s.Log.Debug("geoip.Service.ConfigByScopedGetter.scpOptionFnc", "scope", h.String())
		}
		if err := s.Options(s.scopedOptionFunc(scpcfg)...); err != nil {
			return scopedConfig{}, errors.Wrap(err, "[geoip] Options by scpOptionFnc")
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
		return empty, errors.NewNotValidf(errScopedConfigNotValid, hash)
	}

	if fallback {
		// fallback to default scope
		var err error
		if !s.defaultScopeCache.isValid() {
			err = errors.NewNotFoundf(errScopedConfigNotValid, scope.DefaultHash)
		}
		return s.defaultScopeCache, err

	}

	// give up, nothing found
	return empty, errConfigNotFound
}

// getScopedConfig part of lookupScopedConfig and doesn't use a lock because the lock
// has been acquired in lookupScopedConfig()
func (s *Service) getScopedConfig(h scope.Hash) (sc scopedConfig, ok bool) {
	s.mu.RLock()
	sc, ok = s.scopeCache[h]
	s.mu.RUnlock()

	if ok {
		var hasChanges bool
		// do some init stuff ...
		if sc.log == nil {
			sc.log = s.defaultScopeCache.log // copy logger
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
