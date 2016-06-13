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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/storage/suspend"
	"github.com/corestoreio/csfw/store/scope"
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

// WithVaryBy ...
func WithVaryBy(vb VaryByer) Option {
	return func(s *Service) {
		s.VaryByer = vb
	}
}

// WithScopedRateLimiter creates a rate limiter for a specific scope with its ID.
// The rate limiter is already warmed up.
func WithScopedRateLimiter(scp scope.Scope, id int64, rl throttled.RateLimiter) Option {
	return func(s *Service) {
		s.mu.Lock()
		s.scopedRLs[scope.NewHash(scp, id)] = rl
		s.mu.Unlock()
	}
}

// WithRateLimiterFactory ...
func WithRateLimiterFactory(rlf RateLimiterFactory) Option {
	return func(s *Service) {
		s.RateLimiterFactory = rlf
	}
}

// WithLogger applies a logger to the default scope which gets inherited to
// subsequent scopes. Mainly used for debugging.
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
		s.optionFactoryState = suspend.NewState()
		return nil
	}
}

// DefaultRequests number of requests allowed per time period.
// Used when *PkgBackend has not been provided.
var DefaultRequests = 100

// DefaultBurst defines the number of requests that
// will be allowed to exceed the rate in a single burst and must be
// greater than or equal to zero.
// Used when *PkgBackend has not been provided.
var DefaultBurst = 20

// DefaultDuration per second (s), minute (i), hour (h), day (d)
// Used when *PkgBackend has not been provided.
var DefaultDuration = "h"

const MemStoreMaxKeys = 65536

// NewGCRAMemStore creates the default memory based GCRA rate limiter.
// It uses the PkgBackend models to create a ratelimiter for each scope.
func NewGCRAMemStore(maxKeys int) RateLimiterFactory {
	return func(be *PkgBackend, sg config.ScopedGetter) (throttled.RateLimiter, error) {

		rlStore, err := memstore.New(maxKeys)
		if err != nil {
			return nil, err
		}

		rq, err := rateQuota(be, sg)
		if err != nil {
			return nil, err
		}

		rl, err := throttled.NewGCRARateLimiter(rlStore, rq)
		if err != nil {
			return nil, err
		}

		return rl, nil
	}
}

// rateQuota creates a new quota for the GCRARateLimiter
func rateQuota(be *PkgBackend, sg config.ScopedGetter) (rq throttled.RateQuota, err error) {

	if be == nil {
		return throttled.RateQuota{
			MaxRate:  calculateRate(DefaultDuration, DefaultRequests),
			MaxBurst: DefaultBurst,
		}, nil
	}

	burst, err := be.RateLimitBurst.Get(sg)
	if err != nil {
		err = errors.Mask(err)
		return
	}
	request, err := be.RateLimitRequests.Get(sg)
	if err != nil {
		err = errors.Mask(err)
		return
	}
	if request == 0 {
		request = DefaultRequests
	}

	rate, err := be.RateLimitDuration.Get(sg, request)
	err = errors.Mask(err)

	rq.MaxRate = rate
	rq.MaxBurst = burst
	return
}
