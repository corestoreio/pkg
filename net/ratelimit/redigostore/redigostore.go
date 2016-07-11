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

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/net/url"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/garyburd/redigo/redis"
	throttledRedis "gopkg.in/throttled/throttled.v2/store/redigostore"
)

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
func WithGCRA(scp scope.Scope, id int64, redisRawURL string, duration rune, requests, burst int) ratelimit.Option {
	h := scope.NewHash(scp, id)

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

	var keyPrefix = "ratelimit_" + h.String()
	return func(s *ratelimit.Service) error {
		rs, err := throttledRedis.New(pool, keyPrefix, int(db))
		if err != nil {
			return errors.NewFatalf("[redigostore] redigostore.New: %s", err)
		}
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.redigostore.WithGCRA",
				log.Stringer("scope", scp),
				log.Int64("scope_id", id),
				log.String("redis_raw_url", redisRawURL),
				log.String("key_prefix", keyPrefix),
				log.String("duration", string(duration)),
				log.Int("requests", requests),
				log.Int("burst", burst),
			)
		}
		return ratelimit.WithGCRAStore(scp, id, rs, duration, requests, burst)(s)
	}
}
