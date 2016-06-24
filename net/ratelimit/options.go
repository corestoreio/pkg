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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

// Option can be used as an argument in NewService to configure it with
// different settings.
type Option func(*Service) error

// OptionFactoryFunc a closure around a scoped configuration to figure out which
// options should be returned depending on the scope brought to you during
// a request.
type OptionFactoryFunc func(config.ScopedGetter) []Option

// WithDefaultConfig applies the default GeoIP configuration settings based for
// a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- Alternative Handler: variable DefaultAlternativeHandler
//		- Logger black hole
//		- Check allow: If allowed countries are empty, all countries are allowed
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache = defaultScopedConfig(h)
			return nil
		}
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		s.scopeCache[h] = defaultScopedConfig(h)
		return nil
	}
}

// WithVaryBy sets a custom Key by http.Request producer. Convenience helper
// function.
func WithVaryBy(vb VaryByer) Option {
	return func(s *Service) error {
		s.VaryByer = vb
		return nil
	}
}

// WithRateLimiter creates a rate limiter for a specific scope with its ID.
// The rate limiter is already warmed up.
func WithRateLimiter(scp scope.Scope, id int64, rl throttled.RateLimiter) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.RateLimiter = rl
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.RateLimiter = rl

		if sc, ok := s.scopeCache[h]; ok {
			sc.RateLimiter = scNew.RateLimiter
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithDeniedHandler sets a custom denied handler for a specific scope. The
// default denied handler returns a simple:
//		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
func WithDeniedHandler(scp scope.Scope, id int64, next http.Handler) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.deniedHandler = next
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.deniedHandler = next

		if sc, ok := s.scopeCache[h]; ok {
			sc.deniedHandler = scNew.deniedHandler
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithLogger applies a logger to the default scope which gets inherited to
// subsequent scopes. Mainly used for debugging. Convenience helper function.
func WithLogger(l log.Logger) Option {
	return func(s *Service) error {
		s.Log = l
		return nil
	}
}

// WithOptionFactory applies a function which lazily loads the option depending
// on the incoming scope within a request. For example applies the backend
// configuration to the service.
//
// WRONG: Once this option function has been set all other manually set option
// functions, which accept a scope and a scope ID as an argument, will be
// overwritten by the new values retrieved from the configuration service.
//
// Example:
//	cfgStruct, err := backendgeoip.NewConfigStructure()
//	if err != nil {
//		panic(err)
//	}
//	pb := backendgeoip.New(cfgStruct)
//
//	geoSrv := geoip.MustNewService(
//		geoip.WithOptionFactory(backendgeoip.PrepareOptions(pb)),
//	)
func WithOptionFactory(f OptionFactoryFunc) Option {
	return func(s *Service) error {
		s.optionFactoryFunc = f
		s.optionInflight = new(singleflight.Group)
		return nil
	}
}

// NewGCRAMemStore creates the default memory based GCRA rate limiter.
// It uses the PkgBackend models to create a ratelimiter for each scope.
func WithGCRAMemStore(scp scope.Scope, id int64, maxKeys int, duration string, requests, burst int) Option {
	return func(s *Service) error {

		rlStore, err := memstore.New(maxKeys)
		if err != nil {
			return err
		}

		rq := throttled.RateQuota{
			MaxRate:  CalculateRate(duration, requests),
			MaxBurst: burst,
		}

		rl, err := throttled.NewGCRARateLimiter(rlStore, rq)
		if err != nil {
			return err
		}

		return WithRateLimiter(scp, id, rl)(s)
	}
}

func CalculateRate(duration string, requests int) (r throttled.Rate) {
	switch duration {
	case "s": // second
		r = throttled.PerSec(requests)
	case "i": // minute
		r = throttled.PerMin(requests)
	case "h": // hour
		r = throttled.PerHour(requests)
	case "d": // day
		r = throttled.PerDay(requests)
	default:
		r = throttled.PerHour(requests)
	}
	return
}
