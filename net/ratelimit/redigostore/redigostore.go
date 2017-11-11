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

package redigostore

import (
	"time"

	"github.com/corestoreio/cspkg/config"
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/net/ratelimit"
	"github.com/corestoreio/cspkg/net/url"
	"github.com/corestoreio/cspkg/store/scope"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/garyburd/redigo/redis"
	throttledRedis "gopkg.in/throttled/throttled.v2/store/redigostore"
)

// OptionName identifies this package within the register of the
// backendratelimit.Backend type.
const OptionName = `redigostore`

// NewOptionFactory creates a new option factory function for the memstore in the
// backend package to be used for automatic scope based configuration
// initialization. Configuration values are read from argument `be`.
func NewOptionFactory(burst, requests cfgmodel.Int, duration cfgmodel.Str, redisURL cfgmodel.Str) (string, ratelimit.OptionFactoryFunc) {
	return OptionName, func(sg config.Scoped) []ratelimit.Option {

		burst, err := burst.Get(sg)
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[redigostore] RateLimitBurst.Get"))
		}
		req, err := requests.Get(sg)
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[redigostore] RateLimitRequests.Get"))
		}
		durRaw, err := duration.Get(sg)
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[redigostore] RateLimitDuration.Get"))
		}

		if len(durRaw) != 1 {
			return ratelimit.OptionsError(errors.NewFatalf("[redigostore] RateLimitDuration invalid character count: %q. Should be one character long.", durRaw))
		}

		dur := rune(durRaw[0])

		redisURL, err := redisURL.Get(sg)
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[redigostore] RateLimitStorageGcraRedis.Get"))
		}
		if redisURL != "" {
			return []ratelimit.Option{
				WithGCRA(redisURL, dur, req, burst, sg.ScopeIDs()...),
			}
		}
		return ratelimit.OptionsError(errors.NewEmptyf("[redigostore] Redis not active because RateLimitStorageGCRARedis is not set."))
	}
}

// WithGCRA creates a new Redis-based store, using the provided pool to get
// its connections. The keys will have the specified keyPrefix, which
// may be an empty string, and the database index specified by db will
// be selected to store the keys. Any updating operations will reset
// the key TTL to the provided value rounded down to the nearest
// second. Depends on Redis 2.6+ for EVAL support.
//
// URLs should follow the draft IANA specification for the
// scheme (https://www.iana.org/assignments/uri-schemes/prov/redis).
//
//
// For example:
// 		redis://localhost:6379/3
// 		redis://:6380/0 => connects to localhost:6380
// 		redis:// => connects to localhost:6379 with DB 0
// 		redis://empty:myPassword@clusterName.xxxxxx.0001.usw2.cache.amazonaws.com:6379/0
//
// GCRA => https://en.wikipedia.org/wiki/Generic_cell_rate_algorithm
// This function implements a debug log.
func WithGCRA(redisRawURL string, duration rune, requests, burst int, scopeIDs ...scope.TypeID) ratelimit.Option {

	address, password, db, err := url.ParseRedis(redisRawURL)
	if err != nil {
		return func(s *ratelimit.Service) error {
			return errors.Wrap(err, "[ratelimit] url.RedisParseURL")
		}
	}

	pool := &redis.Pool{
		// todo(CS): maybe make this also configurable ...
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

	scID := scope.DefaultTypeID
	if len(scopeIDs) > 0 {
		// the first item in the slice defines the applied scope.
		scID = scopeIDs[0]
	}

	var keyPrefix = "ratelimit_" + scID.String()
	return func(s *ratelimit.Service) error {
		rs, err := throttledRedis.New(pool, keyPrefix, int(db))
		if err != nil {
			return errors.NewFatalf("[redigostore] redigostore.New: %s", err)
		}
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.redigostore.WithGCRA",
				log.Stringer("scope", scope.TypeIDs(scopeIDs)),
				log.String("redis_raw_url", redisRawURL),
				log.String("key_prefix", keyPrefix),
				log.String("duration", string(duration)),
				log.Int("requests", requests),
				log.Int("burst", burst),
			)
		}
		return ratelimit.WithGCRAStore(rs, duration, requests, burst, scopeIDs...)(s)
	}
}
