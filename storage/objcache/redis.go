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
	gourl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/net/url"
	"github.com/gomodule/redigo/redis"
)

// WithRedisClient connects to the Redis server and does a ping to check if the
// connection works correctly.
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

func parseDuration(keys []string, ds []*time.Duration, params gourl.Values) (err error) {
	i := 0
	for ; i < len(keys) && err == nil; i++ {
		if p := params.Get(keys[i]); p != "" {
			*(ds[i]), err = time.ParseDuration(p)
		}
	}
	if err != nil {
		err = errors.NotValid.New(err, "[objcache] WithRedisURL Parameter %q with value %q is invalid", keys[i], params.Get(keys[i]))
	}
	return err
}

func parseInt(keys []string, ds []*int, params gourl.Values) (err error) {
	i := 0
	for ; i < len(keys) && err == nil; i++ {
		if p := params.Get(keys[i]); p != "" {
			*(ds[i]), err = strconv.Atoi(p)
		}
	}
	if err != nil {
		err = errors.NotValid.New(err, "[objcache] WithRedisURL Parameter %q with value %q is invalid", keys[i], params.Get(keys[i]))
	}
	return err
}

// WithRedisURL connects to a Redis server at the given URL using the HTTP URI
// scheme. URLs should follow the draft IANA specification for the scheme
// (https://www.iana.org/assignments/uri-schemes/prov/redis). This option
// function sets the connection as cache backend to the Service.
//
// On error, while parsing the rawURL, this function will leak sensitive data,
// for now.
//
// For example:
// 		redis://localhost:6379/?db=3
// 		redis://localhost:6379/?max_active=50&max_idle=5&idle_timeout=10s&max_conn_lifetime=1m
func WithRedisURL(rawURL string) Option {
	return Option{
		sortOrder: 2,
		fn: func(p *Service) (err error) {
			addr, _, password, params, err := url.ParseConnection(rawURL)
			if err != nil {
				return errors.Wrapf(err, "[objcache] Redis error parsing URL %q", rawURL)
			}
			pool := &redis.Pool{
				Dial: func() (redis.Conn, error) {
					c, err := redis.Dial("tcp", addr)
					if err != nil {
						return nil, errors.Fatal.New(err, "[objcache] Redis Dial failed")
					}
					if password != "" {
						if _, err := c.Do("AUTH", password); err != nil {
							c.Close()
							return nil, errors.Unauthorized.New(err, "[objcache] Redis AUTH failed")
						}
					}
					if _, err := c.Do("SELECT", params.Get("db")); err != nil {
						c.Close()
						return nil, errors.Fatal.New(err, "[objcache] Redis DB select failed")
					}
					return c, nil
				},
			}

			// add some more params if missing
			err = parseDuration([]string{
				"max_conn_lifetime",
				"idle_timeout",
			}, []*time.Duration{
				&pool.MaxConnLifetime,
				&pool.IdleTimeout,
			}, params)
			if err != nil {
				return errors.WithStack(err)
			}

			err = parseInt([]string{
				"max_idle",
				"max_active",
			}, []*int{
				&pool.MaxIdle,
				&pool.MaxActive,
			}, params)
			if err != nil {
				return errors.WithStack(err)
			}

			// if params.Get("tls") == "1" {
			// 	o.TLSConfig = &tls.Config{ServerName: o.Addr} // TODO check if might be wrong the Addr,
			// }

			return WithRedisClient(pool).fn(p)
		},
	}
}

func doPing(w redisWrapper) error {
	if !w.ping {
		return nil
	}
	conn := w.Pool.Get()
	defer conn.Close()

	pong, err := redis.String(conn.Do("PING"))
	if err != nil && err != redis.ErrNil {
		return errors.Fatal.Newf("[objcache] Redis Ping failed: %s", err)
	}
	if pong != "PONG" {
		return errors.ConnectionFailed.Newf("[objcache] Redis Ping not Pong: %#v", pong)
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

func (w redisWrapper) Set(ctx context.Context, keys []string, values [][]byte, expirations []int64) (err error) {
	conn := w.Pool.Get()
	defer func() {
		if err2 := conn.Close(); err == nil && err2 != nil {
			err = err2
		}
	}()

	args := make([]interface{}, 0, len(keys)*3)
	for i, key := range keys {
		e := expirations[i] // e = expires in x seconds
		if e < 1 {
			args = append(args, key, values[i])
		} else {
			if _, err2 := conn.Do("SETEX", key, e, values[i]); err2 != nil {
				err = errors.Wrapf(err2, "[objcache] With key %q", key)
				return
			}
		}
	}

	if la := len(args); la > 0 && la%2 == 0 {
		if _, err2 := conn.Do("MSET", args...); err2 != nil {
			err = errors.Wrapf(err2, "[objcache] With keys %v", keys)
			return
		}
	}

	return err
}

func (w redisWrapper) Get(_ context.Context, keys []string) (values [][]byte, err error) {
	conn := w.Pool.Get()
	defer func() {
		if err2 := conn.Close(); err == nil && err2 != nil {
			err = err2
		}
	}()

	if len(keys) == 1 {
		var val []byte
		val, err = redis.Bytes(conn.Do("GET", keys[0]))
		if err != nil {
			if err == redis.ErrNil {
				err = ErrKeyNotFound(keys[0])
			} else {
				err = errors.Wrapf(err, "[objcache] With keys %v", keys)
			}
		} else {
			values = append(values, val)
		}
		return
	}

	// TODO optimize for length==1
	values, err = redis.ByteSlices(conn.Do("MGET", strSliceToIFaces(nil, keys)...))
	if err != nil {
		err = errors.Wrapf(err, "[objcache] With keys %v", keys)
		return
	}
	if lk, lv := len(keys), len(values); lk != lv {
		err = ErrKeyNotFound(strings.Join(keys, ", "))
		return
	}
	for i, key := range keys {
		if values[i] == nil {
			err = ErrKeyNotFound(key)
			return
		}
	}
	return
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

func (w redisWrapper) Delete(_ context.Context, keys []string) (err error) {
	conn := w.Pool.Get()
	defer func() {
		if err2 := conn.Close(); err == nil && err2 != nil {
			err = err2
		}
	}()
	if _, err = conn.Do("DEL", strSliceToIFaces(nil, keys)...); err != nil {
		err = errors.Wrapf(err, "[objcache] With keys %v", keys)
	}
	return
}
