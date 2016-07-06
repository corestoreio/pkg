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

type optionFactory struct {
	// OptionFactoryFunc optional configuration closure, can be nil. It pulls
	// out the configuration settings during a request and caches the settings
	// in the internal map. ScopedOption requires a config.ScopedGetter. This
	// function gets set via WithOptionFactory()
	OptionFactoryFunc

	// Group checks on a per scope.Hash basis if the configuration
	// loading process takes place. Stops the execution of other Goroutines (aka
	// incoming requests) with the same scope.Hash until the configuration has
	// been fully loaded and applied and for that specific scope. This function
	// gets set via WithOptionFactory()
	*singleflight.Group

	// rwmu protects all fields below
	sync.RWMutex
	applied map[scope.Hash]struct{}
}

func (of *optionFactory) shouldRun(h scope.Hash) bool {
	if of.OptionFactoryFunc == nil {
		return false
	}
	of.RLock()
	_, ok := of.applied[h]
	of.RUnlock()
	return !ok
}

func (of *optionFactory) done(h scope.Hash) {
	println("optionFactory ==> done!", h.String())
	of.Lock()
	of.applied[h] = struct{}{}
	of.Unlock()
}

type service struct {
	// Log used for debugging. Defaults to black hole. Panics if nil.
	Log log.Logger

	oFactory *optionFactory

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
			Log: log.BlackHole{},
			oFactory: &optionFactory{
				applied: make(map[scope.Hash]struct{}),
			},
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
// from the backend configuration service. The field Inflight handles the
// guaranteed atomic single loading for each scope.
// It returns never nil.
func (s *Service) configByScopedGetter(scpGet config.ScopedGetter) scopedConfig {

	h := scope.NewHash(scpGet.Scope())  // can be store or website or default
	p := scope.NewHash(scpGet.Parent()) // can be website or default

	if s.oFactory.shouldRun(h) {
		//if s.Log.IsDebug() {
		//	s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.OptionFactoryRun",
		//		log.Stringer("requested_scope", h),
		//		log.Stringer("requested_parent_scope", p),
		//		log.Stringer("responded_scope", sCfg.scopeHash),
		//		log.ErrWithKey("responded_scope_valid", sCfg.isValid()),
		//	)
		//}
		//
		defer s.oFactory.Group.Forget(h.String())
		scpCfgChan := s.oFactory.DoChan(h.String(), func() (interface{}, error) {
			defer s.oFactory.done(h)
			if err := s.Options(s.oFactory.OptionFactoryFunc(scpGet)...); err != nil {
				return newScopedConfigError(errors.Wrap(err, prefixError+" Options applied by OptionFactoryFunc")), nil
			}
			sCfg := s.getConfigByHash(h, p)
			if s.Log.IsDebug() {
				s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.Inflight.Do",
					log.Stringer("requested_scope", h),
					log.Stringer("requested_parent_scope", p),
					log.Stringer("responded_scope", sCfg.scopeHash),
					log.ErrWithKey("responded_scope_valid", sCfg.isValid()),
				)
			}
			return sCfg, nil
		})

		res, ok := <-scpCfgChan
		if !ok {
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

	// 99 % of the hits
	if sCfg := s.getConfigByHash(h, p); sCfg.isValid() == nil {
		// cached entry found which can contain the configured scoped configuration or
		// the default scope configuration.

		// println("Hit allowedOriginsAll", sCfg.allowedOriginsAll, sCfg.scopeHash.String(), "Requested: ", h.String(), p.String())

		if s.Log.IsDebug() {
			s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.IsValid",
				log.Stringer("requested_scope", h),
				log.Stringer("requested_parent_scope", p),
				log.Stringer("responded_scope", sCfg.scopeHash),
				log.ErrWithKey("responded_scope_valid", sCfg.isValid()),
			)
		}
		return sCfg
	}

	// finaly fall back to default scope
	p = scope.DefaultHash
	sCfg := s.getConfigByHash(h, p)
	if s.Log.IsDebug() {
		s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.Fallback.Default",
			log.Stringer("requested_scope", h),
			log.Stringer("requested_parent_scope", p),
			log.Stringer("responded_scope", sCfg.scopeHash),
			log.ErrWithKey("responded_scope_valid", sCfg.isValid()),
		)
	}
	return sCfg
}

// runOptionFactory uses the external configuration to load the values and
// executes the functional options. after the options has been applied a call to
// getConfigByHash() will be made to select the correct cached configuration.
//func (s *Service) runOptionFactory(scpGet config.ScopedGetter) scopedConfig {
//	if s.oFactory.OptionFactoryFunc == nil {
//		return newScopedConfigError(errOptionFactoryNotSet)
//	}
//
//	h := scope.NewHash(scpGet.Scope())
//	p := scope.NewHash(scpGet.Parent())
//
//	scpCfgChan := s.oFactory.DoChan(h.String(), func() (interface{}, error) {
//		if err := s.Options(s.oFactory.OptionFactoryFunc(scpGet)...); err != nil {
//			return newScopedConfigError(errors.Wrap(err, prefixError+" Options applied by OptionFactoryFunc")), nil
//		}
//		sCfg := s.getConfigByHash(h, p)
//		if s.Log.IsDebug() {
//			s.Log.Debug(prefixLog+"Service.ConfigByScopedGetter.Inflight.Do",
//				log.Stringer("requested_scope", h),
//				log.Stringer("requested_parent_scope", p),
//				log.Stringer("responded_scope", sCfg.scopeHash),
//				log.ErrWithKey("responded_scope_valid", sCfg.isValid()),
//			)
//		}
//		return sCfg, nil
//	})
//
//	res, ok := <-scpCfgChan
//	if !ok {
//		return newScopedConfigError(errors.NewFatalf(prefixError + " Inflight.DoChan returned a closed/unreadable channel"))
//	}
//
//	if res.Err != nil {
//		return newScopedConfigError(errors.Wrap(res.Err, prefixError+" Inflight.DoChan.Error"))
//	}
//	sCfg, ok := res.Val.(scopedConfig)
//	if !ok {
//		sCfg = newScopedConfigError(errors.NewFatalf(prefixError + " Inflight.DoChan res.Val cannot be type asserted to scopedConfig"))
//	}
//	return sCfg
//}

// getConfigByScopeID returns the correct configuration for a scope and may fall
// back to the next higher scope: store -> website -> default. if an entry for a
// scope cannot be found the next higher get looked up and the pointer of the
// next higher scope gets assigned to the current scope.
func (s *Service) getConfigByHash(hash scope.Hash, parent scope.Hash) (scpCfg scopedConfig) {
	// hash can be store or website scope
	// parent can be website or default scope.

	if parent.Scope() < scope.Default {
		return newScopedConfigError(errors.NewFatalf(prefixError+" Parent scope must be minimum scope.Default: %s", parent))
	}

	// pointer gets dereferenced in a lock to avoid race conditions while
	// reading in middleware the config values because we might execute the
	// functional options for another scope while one scope runs in the
	// middleware.

	// lookup store/website scope. this should hit 99% of the calls of this function.
	s.rwmu.RLock()
	pScpCfg, ok := s.scopeCache[hash]
	if ok && pScpCfg != nil {
		scpCfg = *pScpCfg
	}
	s.rwmu.RUnlock()
	if ok {
		return scpCfg
	}

	// now lock everything until the fall back has been found.
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	// if the store scope cannot be found, fall back to website.
	if !ok && parent.Scope() == scope.Website {
		pScpCfg, ok = s.scopeCache[parent]
		if ok && pScpCfg != nil {
			scpCfg = *pScpCfg
		}
		if ok && pScpCfg != nil {
			s.scopeCache[hash] = pScpCfg
			return scpCfg
		}
	}

	// default config lookup
	if !ok {
		pScpCfg, ok = s.scopeCache[scope.DefaultHash]
		if ok && pScpCfg != nil {
			scpCfg = *pScpCfg
		}
		if ok && pScpCfg != nil {
			s.scopeCache[hash] = pScpCfg
		}
	}

	return scpCfg
}
