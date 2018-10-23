// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// +build redis csall

package objcache

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/net/url"
	"github.com/garyburd/redigo/redis"
)

// WithRedisClient connects to the Redis server. Set ping to true to check if the
// connection works correctly.
//
// For options see: https://godoc.org/gopkg.in/redis.v3#Options
func WithRedisClient(pool *redis.Pool) Option {
	return Option{
		sortOrder: 1,
		fn: func(p *Service) error {
			w := makeRedisWrapper(pool)
			w.ping = true
			if err := doPing(w); err != nil {
				return errors.WithStack(err)
			}
			p.cache[len(p.cache)+1] = w
			return nil
		},
	}
}

// WithRedisURL connects to a Redis server at the given URL using the Redis URI
// scheme. URLs should follow the draft IANA specification for the scheme
// (https://www.iana.org/assignments/uri-schemes/prov/redis). This option
// function sets the connection as cache backend to the Service.
//
// On error, while parsing the rawURL, this function will leak sensitive data,
// for now.
//
// For example:
// 		redis://localhost:6379/3
// 		redis://localhost:6379/?max_active=50&max_idle=5&idle_timeout=10s
func WithRedisURL(rawURL string) Option {
	return Option{
		sortOrder: 2,
		fn: func(p *Service) error {

			addr, _, pw, params, err := url.ParseConnection(rawURL)
			if err != nil {
				return errors.Wrapf(err, "[backend] Redis error parsing URL %q", rawURL)
			}
			maxActive, err := strconv.Atoi(params.Get("max_active"))
			if err != nil {
				return errors.Wrapf(err, "[backend] NewRedis.ParseNoSQLURL. Parameter max_active not valid in %q", rawURL)
			}
			maxIdle, err := strconv.Atoi(params.Get("max_idle"))
			if err != nil {
				return errors.NotValid.New(err, "[backend] NewRedis.ParseNoSQLURL. Parameter max_idle not valid in %q", rawURL)
			}
			idleTimeout, err := time.ParseDuration(params.Get("idle_timeout"))
			if err != nil {
				return errors.NotValid.New(err, "[backend] NewRedis.ParseNoSQLURL. Parameter idle_timeout not valid in %q", rawURL)
			}

			pool := &redis.Pool{
				MaxActive:   maxActive,
				MaxIdle:     maxIdle,
				IdleTimeout: idleTimeout,
				Dial: func() (redis.Conn, error) {
					c, err := redis.Dial("tcp", addr)
					if err != nil {
						return nil, errors.Fatal.New(err, "[backend] Redis Dial failed")
					}
					if pw != "" {
						if _, err := c.Do("AUTH", pw); err != nil {
							c.Close()
							return nil, errors.Unauthorized.New(err, "[backend] Redis AUTH failed")
						}
					}
					if _, err := c.Do("SELECT", params.Get("db")); err != nil {
						c.Close()
						return nil, errors.Fatal.New(err, "[backend] Redis DB select failed")
					}
					return c, nil
				},
			}

			return WithRedisClient(pool).fn(p)
		},
	}
}

func doPing(w redisWrapper) error {
	if !w.ping || w.Pool == nil {
		return nil
	}
	conn := w.Pool.Get()
	defer conn.Close()

	pong, err := redis.String(conn.Do("PING"))
	if err != nil && err != redis.ErrNil {
		return errors.Fatal.Newf("[backend] Redis Ping failed: %s", err)
	}
	if pong != "PONG" {
		return errors.ConnectionFailed.Newf("[backend] Redis Ping not Pong: %#v", pong)
	}
	return nil
}

func makeRedisWrapper(rp *redis.Pool) redisWrapper {
	return redisWrapper{
		Pool: rp,
		// ipf: &sync.Pool{
		// 	New: func() interface{} {
		// 		ifs := make([]interface{}, 0, 10)
		// 		return &ifs
		// 	},
		// },
	}
}

type redisWrapper struct {
	*redis.Pool
	ping bool
	// ifp  *sync.Pool
}

func (w redisWrapper) Set(_ context.Context, items *Items) error {
	conn := w.Pool.Get()
	defer conn.Close()

	keys, values, err := items.Encode(nil, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	if _, err := conn.Do("MSET", byteSliceToIFaces(strSliceToIFaces(nil, keys), values)...); err != nil {
		return errors.Wrapf(err, "[objcache] With keys %v", keys)
	}

	return nil
}

func (w redisWrapper) Get(_ context.Context, keys []string) (values [][]byte, err error) {
	conn := w.Pool.Get()
	defer conn.Close()
	// TODO optimize for length==1
	values, err = redis.ByteSlices(conn.Do("MGET", strSliceToIFaces(nil, keys)...))
	if err != nil {
		if err != redis.ErrNil {
			return nil, errors.Wrapf(err, "[objcache] With keys %v", keys)
		}
		return nil, ErrKeyNotFound(strings.Join(keys, ", "))
	}
	if lk, lv := len(keys), len(values); lk != lv {
		return nil, ErrKeyNotFound(strings.Join(keys, ", "))
	}
	return values, nil
}

func strSliceToIFaces(ret []interface{}, sl []string) []interface{} {
	// TODO use a sync.Pool but write before hand appropriate concurrent running benchmarks
	if ret == nil {
		ret = make([]interface{}, 0, len(sl))
	}
	for _, s := range sl {
		ret = append(ret, s)
	}
	return ret
}

func byteSliceToIFaces(ret []interface{}, sl [][]byte) []interface{} {
	if ret == nil {
		ret = make([]interface{}, 0, len(sl))
	}
	for _, s := range sl {
		ret = append(ret, s)
	}
	return ret
}

func (w redisWrapper) Delete(_ context.Context, keys []string) error {
	conn := w.Pool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", strSliceToIFaces(nil, keys)...); err != nil {
		return errors.Wrapf(err, "[objcache] With keys %v", keys)
	}
	return nil
}
