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

package cors

import (
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
	"github.com/corestoreio/csfw/util/errors"
)

// auto generated: do not edit. See net/gen eric package

const (
	prefixError = `[cors] `
	prefixLog   = `cors.`
)

type service struct {
	// Log used for debugging. Defaults to black hole. Panics if nil.
	Log log.Logger

	// optionValidation allows to set a custom option function which can check
	// if e.g. the scopes of the applied options are correct. Gets only executed
	// if not nil.
	optionValidation func() error

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
	defaultScopeCache *scopedConfig

	// scopeCache internal cache of the configurations. scoped.Hash relates to
	// the default,website or store ID.
	scopeCache map[scope.Hash]*scopedConfig
}

func newService(opts ...Option) (*Service, error) {
	s := &Service{
		service: service{
			Log:        log.BlackHole{},
			scopeCache: make(map[scope.Hash]*scopedConfig),
		},
	}
	if err := s.Options(WithDefaultConfig(scope.Default, 0)); err != nil {
		return nil, errors.Wrap(err, prefixError+" Options WithDefaultConfig")
	}
	if err := s.Options(opts...); err != nil {
		return nil, errors.Wrap(err, prefixError+" Options Any Config")
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
				return errors.Wrap(err, prefixError+" Service.Options")
			}
		}
	}
	if s.optionValidation != nil {
		return errors.Wrap(s.optionValidation(), prefixError+" optionValidation")
	}
	return nil
}

// flushCache cors cache flusher
func (s *Service) flushCache() error {
	s.scopeCache = make(map[scope.Hash]*scopedConfig)
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
// It returns never nil.
func (s *Service) configByScopedGetter(scpGet config.ScopedGetter) *scopedConfig {

	h := scope.NewHash(scpGet.Scope())
	// fallback to default scope
	if s.useDefaultConfig(h) {
		if s.Log.IsDebug() {
			s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.defaultScopeCache", log.Stringer("scope", h), log.Bool("optionFactoryFunc_Nil", s.optionFactoryFunc == nil))
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
			s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.IsValid", log.Stringer("scope", h), log.Stringer("parentScope", p))
		}
		return sCfg
	case s.optionFactoryFunc == nil:
		if s.Log.IsDebug() {
			s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.optionFactoryFunc.nil", log.Stringer("scope", h), log.Stringer("parentScope", p))
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
			return newScopedConfigError(errors.Wrap(err, prefixError+" Options by scpOptionFnc")), nil
		}
		if s.Log.IsDebug() {
			s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.optionInflight.Do", log.Stringer("scope", h), log.Stringer("parentScope", p), log.Err(sCfg.isValid()))
		}
		return s.getConfigByScopeID(h, p), nil
	})

	res, ok := <-scpCfgChan
	if !ok {
		return newScopedConfigError(errors.NewFatalf(prefixError + " optionInflight.DoChan returned a closed/unreadable channel"))
	}

	if res.Err != nil {
		return newScopedConfigError(errors.Wrap(res.Err, prefixError+" optionInflight.DoChan.Error"))
	}
	sCfg, ok = res.Val.(*scopedConfig)
	if !ok {
		sCfg.lastErr = errors.NewFatalf(prefixError + " optionInflight.DoChan res.Val cannot be type asserted to scopedConfig")
	}
	return sCfg
}

// getScpCfg returns the config for argument "hash" and uses argument "parent"
// only to check if it must fall back to the default scope. if so the "parent"
// field gets the defaultScopeCache assigned.
func (s *Service) getScpCfg(hash scope.Hash, parent scope.Hash) *scopedConfig {
	s.rwmu.RLock()
	scpCfg, ok := s.scopeCache[hash]
	s.rwmu.RUnlock()

	if ok {
		return scpCfg
	}
	if !ok && parent.EqualScope(scope.DefaultHash) {
		s.rwmu.Lock()
		scpCfg = s.defaultScopeCache
		s.scopeCache[hash] = s.defaultScopeCache
		s.rwmu.Unlock()
		if s.Log.IsDebug() {
			s.Log.Debug(prefixLog+"Service.getScpCfg.DefaultScopeCache",
				log.Stringer("arg_scope", hash),
				log.Stringer("arg_scope_parent", parent),
				log.String("scope_applied", scpCfg.printScope()),
			)
		}
		return s.defaultScopeCache
	}
	return scpCfg
}

// getConfigByScopeID returns the correct configuration for a scope and may fall back
// to the next higher scope: store -> website -> default.
func (s *Service) getConfigByScopeID(hash scope.Hash, parent scope.Hash) *scopedConfig {

	scpCfg := s.getScpCfg(hash, parent)

	if scpCfg.isValid() == nil {
		if s.Log.IsDebug() {
			s.Log.Debug(prefixLog+"Service.getConfigByScopeID.Hash.Valid.Cached",
				log.Stringer("arg_scope", hash),
				log.Stringer("arg_scope_parent", parent),
				log.String("scope_applied", scpCfg.printScope()),
			)
		}
		return scpCfg
	}

	// lookup parent configuration scope
	if parent.Scope() == scope.Website {
		// overwrite store scope with website scope pointer
		scpCfg = s.getScpCfg(parent, 0)

		if err := scpCfg.isValid(); err == nil {
			// we found an entry for a website config
			s.rwmu.Lock()
			s.scopeCache[hash] = scpCfg // set the hash to the parent website configuration
			s.rwmu.Unlock()
			if s.Log.IsDebug() {
				s.Log.Debug(prefixLog+"Service.getConfigByScopeID.Parent.Valid.New",
					log.Stringer("arg_scope", hash),
					log.Stringer("arg_scope_parent", parent),
					log.String("scope_applied", scpCfg.printScope()),
				)
			}
		} else if s.Log.IsDebug() {
			s.Log.Debug(prefixLog+"Service.getConfigByScopeID.Parent.Invalid",
				log.Stringer("arg_scope", hash),
				log.Stringer("arg_scope_parent", parent),
				log.Err(err),
			)
		}
		// return website config
	}

	if s.Log.IsDebug() {
		s.Log.Debug(prefixLog+"Service.getConfigByScopeID.Return",
			log.Stringer("arg_scope", hash),
			log.Stringer("arg_scope_parent", parent),
			log.ErrWithKey("scp_cfg_is_valid", scpCfg.isValid()),
		)
	}

	return scpCfg
}
