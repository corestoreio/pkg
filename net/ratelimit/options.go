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
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/url"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
	"gopkg.in/throttled/throttled.v2/store/redigostore"
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
//		- Denied Handler: http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
//		- VaryByer: returns an empty key
// Example:
//		s := MustNewService(WithDefaultConfig(scope.Store,1), WithVaryBy(scope.Store, 1, myVB))
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

// WithVaryBy allows to set a custom key producer. VaryByer is called for each
// request to generate a key for the limiter. If it is nil, the middleware
// panics. The default VaryByer returns an empty string so that all requests
// uses the same key. VaryByer must be thread safe.
func WithVaryBy(scp scope.Scope, id int64, vb VaryByer) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.VaryByer = vb
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.VaryByer = vb

		if sc, ok := s.scopeCache[h]; ok {
			sc.VaryByer = scNew.VaryByer
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
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

// WithGCRAStore creates a new GCRA rate limiter with a custom storage backend.
// Duration: (s second,i minute,h hour,d day)
func WithGCRAStore(scp scope.Scope, id int64, store throttled.GCRAStore, duration rune, requests, burst int) Option {
	return func(s *Service) error {

		rq := throttled.RateQuota{
			MaxRate:  CalculateRate(duration, requests),
			MaxBurst: burst,
		}

		rl, err := throttled.NewGCRARateLimiter(store, rq)
		if err != nil {
			return errors.NewNotValidf("[ratelimit] throttled.NewGCRARateLimiter: %s", err)
		}

		return WithRateLimiter(scp, id, rl)(s)
	}
}

// WithGCRAMemStore creates the default memory based GCRA rate limiter.
// Duration: (s second,i minute,h hour,d day)
func WithGCRAMemStore(scp scope.Scope, id int64, maxKeys int, duration rune, requests, burst int) Option {
	return func(s *Service) error {
		rlStore, err := memstore.New(maxKeys)
		if err != nil {
			return errors.NewFatalf("[ratelimit] memstore.New MaxKeys(%d): %s", maxKeys, err)
		}
		return WithGCRAStore(scp, id, rlStore, duration, requests, burst)(s)
	}
}

// WithGCRARedis creates a new Redis-based store, using the provided pool to get
// its connections. The keys will have the specified keyPrefix, which
// may be an empty string, and the database index specified by db will
// be selected to store the keys. Any updating operations will reset
// the key TTL to the provided value rounded down to the nearest
// second. Depends on Redis 2.6+ for EVAL support.
func WithGCRARedis(scp scope.Scope, id int64, redisRawUrl string, duration rune, requests, burst int) Option {
	h := scope.NewHash(scp, id)

	address, password, db, err := url.ParseRedis(redisRawUrl)
	if err != nil {
		return func(s *Service) error {
			return errors.Wrap(err, "[ratelimit] url.RedisParseURL")
		}
	}

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address, redis.DialPassword(password))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return func(s *Service) error {
		rs, err := redigostore.New(pool, "ratelimit_"+h.String(), int(db))
		if err != nil {
			return errors.NewFatalf("[ratelimit] redigostore.New: %s", err)
		}
		return WithGCRAStore(scp, id, rs, duration, requests, burst)(s)
	}
}

// CalculateRate calculates the rate depending on the duration (s second,i minute,h hour,d day) and the
// maximum requests. Invalid duration falls back to an hourly calculation.
func CalculateRate(duration rune, requests int) (r throttled.Rate) {
	switch duration {
	case 's': // second
		r = throttled.PerSec(requests)
	case 'i': // minute
		r = throttled.PerMin(requests)
	case 'h': // hour
		r = throttled.PerHour(requests)
	case 'd': // day
		r = throttled.PerDay(requests)
	default:
		r = throttled.PerHour(requests)
	}
	return
}
