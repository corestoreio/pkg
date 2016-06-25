// Copyright (c) 2014 Olivier Poitrey <rs@dailymotion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cors

import (
	"net/http"
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
	"github.com/corestoreio/csfw/util/errors"
)

const methodOptions = "OPTIONS"

// Service describes the CrossOriginResourceSharing which is used to create a
// Container Filter that implements CORS. Cross-origin resource sharing (CORS)
// is a mechanism that allows JavaScript on a web page to make XMLHttpRequests
// to another domain, not the domain the JavaScript originated from.
//
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
// http://enable-cors.org/server.html
// http://www.html5rocks.com/en/tutorials/cors/#toc-handling-a-not-so-simple-request
type Service struct {
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

	// scopeCache internal cache of the configurations. scoped.Hash relates to
	// the default,website or store ID.
	scopeCache map[scope.Hash]scopedConfig
}

// New creates a new Cors handler with the provided options.
func New(opts ...Option) (*Service, error) {
	s := &Service{
		scopeCache: make(map[scope.Hash]scopedConfig),
	}
	if err := s.Options(WithDefaultConfig(scope.Default, 0)); err != nil {
		return nil, errors.Wrap(err, "[cors] Options WithDefaultConfig")
	}
	if err := s.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[cors] Options Any Config")
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
				return errors.Wrap(err, "[cors] Service.Options")
			}
		}
	}

	s.rwmu.RLock()
	defer s.rwmu.RUnlock()
	for h := range s.scopeCache {
		if scp, _ := h.Unpack(); scp > scope.Website {
			return errors.NewNotSupportedf(errServiceUnsupportedScope, h)
		}
	}
	return nil
}

// WithCORS to be used as a middleware for net.Handler. The applied
// configuration is used for the all store scopes or if the PkgBackend has been
// provided then on a website specific level. Middleware expects to find in a
// context a store.FromContextProvider().
func (s *Service) WithCORS() mw.Middleware {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			requestedStore, err := store.FromContextRequestedStore(ctx)
			if err != nil {
				if s.defaultScopeCache.log.IsDebug() {
					s.defaultScopeCache.log.Debug("Service.WithCORS.FromContextProvider", log.Err(err), log.HTTPRequest("request", r))
				}
				err = errors.Wrap(err, "[cors] FromContextRequestedStore")
				h.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			// the scpCfg depends on how you have initialized the storeService during app boot.
			// requestedStore.Website.Config is the reason that all options only support
			// website scope and not group or store scope.
			scpCfg := s.configByScopedGetter(requestedStore.Website.Config)
			if err := scpCfg.isValid(); err != nil {
				if s.defaultScopeCache.log.IsDebug() {
					s.defaultScopeCache.log.Debug("Service.WithCORS.configByScopedGetter", log.Err(err), log.Marshal("requestedStore", requestedStore), log.HTTPRequest("request", r))
				}
				err = errors.Wrap(err, "[cors] ConfigByScopedGetter")
				h.ServeHTTP(w, r.WithContext(withContextError(ctx, err)))
				return
			}

			if s.defaultScopeCache.log.IsInfo() {
				s.defaultScopeCache.log.Info("Service.WithCORS.handleActualRequest", log.String("method", r.Method), log.Object("scopedConfig", scpCfg))
			}

			if r.Method == methodOptions {
				if s.defaultScopeCache.log.IsDebug() {
					s.defaultScopeCache.log.Debug("Service.WithCORS.handlePreflight", log.String("method", r.Method), log.Bool("OptionsPassthrough", scpCfg.optionsPassthrough))
				}
				scpCfg.handlePreflight(w, r)
				// Preflight requests are standalone and should stop the chain as some other
				// middleware may not handle OPTIONS requests correctly. One typical example
				// is authentication middleware ; OPTIONS requests won't carry authentication
				// headers (see #1)
				if scpCfg.optionsPassthrough {
					h.ServeHTTP(w, r)
				}
				return
			}
			scpCfg.handleActualRequest(w, r)
			h.ServeHTTP(w, r)
		})
	}
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
		if s.defaultScopeCache.log.IsDebug() {
			s.defaultScopeCache.log.Debug("cors.Service.ConfigByScopedGetter.defaultScopeCache", log.Stringer("scope", h), log.Bool("optionFactoryFunc_Nil", s.optionFactoryFunc == nil))
		}
		return s.defaultScopeCache
	}

	sCfg := s.getConfigByScopeID(h, false) // 1. map lookup, but only one lookup during many requests, 99%

	switch {
	case sCfg.isValid() == nil:
		// cached entry found which can contain the configured scoped configuration or
		// the default scope configuration.
		if s.defaultScopeCache.log.IsDebug() {
			s.defaultScopeCache.log.Debug("cors.Service.ConfigByScopedGetter.IsValid", log.Stringer("scope", h))
		}
		return sCfg
	case s.optionFactoryFunc == nil:
		if s.defaultScopeCache.log.IsDebug() {
			s.defaultScopeCache.log.Debug("cors.Service.ConfigByScopedGetter.optionFactoryFunc.nil", log.Stringer("scope", h))
		}
		// When everything has been pre-configured for each scope via functional
		// options then we might have the case where a scope in a request comes
		// in which does not match to a scopedConfiguration entry in the map
		// therefore a fallback to default scope must be provided. This fall
		// back will only be executed once and each scope knows from now on that
		// it has the configuration of the default scope.
		return s.getConfigByScopeID(h, true)
	}

	// the following code gets tested by backendcors package

	scpCfgChan := s.optionInflight.DoChan(h.String(), func() (interface{}, error) {
		if err := s.Options(s.optionFactoryFunc(scpGet)...); err != nil {
			return scopedConfig{
				lastErr: errors.Wrap(err, "[cors] Options by scpOptionFnc"),
			}, nil
		}
		if s.defaultScopeCache.log.IsDebug() {
			s.defaultScopeCache.log.Debug("cors.Service.ConfigByScopedGetter.optionInflight.DoChan", log.Stringer("scope", h), log.Err(sCfg.isValid()))
		}
		return s.getConfigByScopeID(h, true), nil
	})

	res, ok := <-scpCfgChan
	if !ok {
		return scopedConfig{lastErr: errors.NewFatalf("[cors] optionInflight.DoChan returned a closed/unreadable channel")}
	}
	if res.Err != nil {
		return scopedConfig{lastErr: errors.Wrap(res.Err, "[cors] optionInflight.DoChan.Error")}
	}
	sCfg, ok = res.Val.(scopedConfig)
	if !ok {
		sCfg.lastErr = errors.NewFatalf("[cors] optionInflight.DoChan res.Val cannot be type asserted to scopedConfig")
	}
	return sCfg
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
		if s.defaultScopeCache.log.IsDebug() {
			s.defaultScopeCache.log.Debug("cors.Service.getConfigByScopeID.fallbackToDefault", log.Stringer("scope", hash))
		}
	}
	if !ok && !useDefault {
		return scopedConfig{
			lastErr: errors.Wrap(errConfigNotFound, "[cors] Service.getConfigByScopeID.NotFound"),
		}
	}
	if scpCfg.useDefault {
		return s.defaultScopeCache
	}
	return scpCfg
}
