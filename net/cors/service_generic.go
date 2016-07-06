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

	// optionFactory optional configuration closure, can be nil. It pulls out
	// the configuration settings from a slow backend during a request and
	// caches the settings in the internal map.  This function gets set via
	// WithOptionFactory()
	optionFactory OptionFactoryFunc

	// optionInflight checks on a per scope.Hash basis if the configuration
	// loading process takes place. Stops the execution of other Goroutines (aka
	// incoming requests) with the same scope.Hash until the configuration has
	// been fully loaded and applied for that specific scope. This function gets
	// set via WithOptionFactory()
	optionInflight *singleflight.Group

	// optionAfterApply allows to set a custom function which runs every time
	// after the options has been applied. Gets only executed if not nil.
	optionAfterApply func() error

	// rwmu protects all fields below
	rwmu sync.RWMutex

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
	if s.optionAfterApply != nil {
		return errors.Wrap(s.optionAfterApply(), prefixError+" optionValidation")
	}
	return nil
}

// flushCache cors cache flusher
func (s *Service) flushCache() error {
	s.scopeCache = make(map[scope.Hash]*scopedConfig)
	return nil
}

// configByScopedGetter returns the internal configuration depending on the
// ScopedGetter. Mainly used within the middleware.  If you have applied the
// option WithOptionFactory() the configuration will be pulled out only one time
// from the backend configuration service. The field optionInflight handles the
// guaranteed atomic single loading for each scope.
func (s *Service) configByScopedGetter(scpGet config.ScopedGetter) scopedConfig {

	current := scope.NewHash(scpGet.Scope())   // can be store or website or default
	fallback := scope.NewHash(scpGet.Parent()) // can be website or default

	// 99.9999 % of the hits; 2nd argument must be zero because we must first
	// test if a direct entry can be found; if not we must apply either the
	// optionFactory function or do a fall back to the website scope and/or
	// default scope.
	if sCfg := s.getConfigByHash(current, 0); sCfg.isValid() == nil {
		if s.Log.IsDebug() {
			s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.IsValid",
				log.Stringer("requested_scope", current),
				log.Stringer("requested_fallback_scope", scope.Hash(0)),
				log.Stringer("responded_scope", sCfg.scopeHash),
				log.ErrWithKey("responded_scope_valid", sCfg.isValid()),
			)
		}
		return sCfg
	}

	// load the configuration from the slow backend. optionInflight guarantees
	// that the closure will only be executed once but the returned result gets
	// returned to all waiting goroutines.
	if s.optionFactory != nil {
		res, ok := <-s.optionInflight.DoChan(current.String(), func() (interface{}, error) {
			if err := s.Options(s.optionFactory(scpGet)...); err != nil {
				return newScopedConfigError(errors.Wrap(err, prefixError+" Options applied by OptionFactoryFunc")), nil
			}
			sCfg := s.getConfigByHash(current, fallback)
			if s.Log.IsDebug() {
				s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.Inflight.Do",
					log.Stringer("requested_scope", current),
					log.Stringer("requested_fallback_scope", fallback),
					log.Stringer("responded_scope", sCfg.scopeHash),
					log.ErrWithKey("responded_scope_valid", sCfg.isValid()),
				)
			}
			return sCfg, nil
		})
		if !ok { // unlikely to happen but you'll never know. how to test that?
			return newScopedConfigError(errors.NewFatalf(prefixError + " Inflight.DoChan returned a closed/unreadable channel"))
		}
		if res.Err != nil {
			return newScopedConfigError(errors.Wrap(res.Err, prefixError+" Inflight.DoChan.Error"))
		}
		sCfg, ok := res.Val.(scopedConfig)
		if !ok {
			sCfg = newScopedConfigError(errors.NewFatalf(prefixError + " Inflight.DoChan res.Val cannot be type asserted to scopedConfig"))
		}
		return sCfg
	}

	sCfg := s.getConfigByHash(current, fallback)
	if s.Log.IsDebug() {
		s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.Fallback.Default",
			log.Stringer("requested_scope", current),
			log.Stringer("requested_fallback_scope", fallback),
			log.Stringer("responded_scope", sCfg.scopeHash),
			log.ErrWithKey("responded_scope_valid", sCfg.isValid()),
		)
	}
	return sCfg
}

// getConfigByHash returns the correct configuration for a scope and may fall
// back to the next higher scope: store -> website -> default. if an entry for a
// scope cannot be found the next higher scope gets looked up and the pointer of
// the next higher scope gets assigned to the current scope. this makes sure to
// avoid redundant configurations and enables us to change one scope with an
// impact on all other scopes which depend on the parent scope.
func (s *Service) getConfigByHash(current scope.Hash, fallback scope.Hash) (scpCfg scopedConfig) {
	// current can be store or website scope
	// fallback can be website or default scope. If 0 then no fall back

	// pointer must get dereferenced in a lock to avoid race conditions while
	// reading in middleware the config values because we might execute the
	// functional options for another scope while one scope runs in the
	// middleware.

	// lookup store/website scope. this should hit 99% of the calls of this function.
	s.rwmu.RLock()
	pScpCfg, ok := s.scopeCache[current]
	if ok && pScpCfg != nil {
		scpCfg = *pScpCfg
	}
	s.rwmu.RUnlock()
	if ok {
		return scpCfg
	}
	if fallback == 0 {
		return newScopedConfigError(errConfigNotFound)
	}

	// slow path: now lock everything until the fall back has been found.
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	// if the current scope cannot be found, fall back to fallback scope and
	// apply the maybe found configuration to the current scope configuration.
	if !ok && fallback.Scope() == scope.Website {
		pScpCfg, ok = s.scopeCache[fallback]
		if ok && pScpCfg != nil {
			scpCfg = *pScpCfg
		}
		if ok && pScpCfg != nil {
			s.scopeCache[current] = pScpCfg
			return scpCfg
		}
	}

	// if the current and fallback scope cannot be found, fall back to default
	// scope and apply the maybe found configuration to the current scope
	// configuration.
	if !ok {
		pScpCfg, ok = s.scopeCache[scope.DefaultHash]
		if ok && pScpCfg != nil {
			scpCfg = *pScpCfg
		}
		if ok && pScpCfg != nil {
			s.scopeCache[current] = pScpCfg
		}
	}
	return scpCfg
}
