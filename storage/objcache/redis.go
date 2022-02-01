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
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/net/url"
	"github.com/gomodule/redigo/redis"
)

// RedisOption applies several options for the Redis client.
type RedisOption struct {
	KeyPrefix string
}

// NewRedisClient connects to the Redis server and does a ping to check if the
// connection works correctly.
func NewRedisClient(pool *redis.Pool, ro *RedisOption) NewStorageFn {
	return func() (Storager, error) {
		w := makeRedisWrapper(pool, ro)
		w.ping = true
		if err := doPing(w); err != nil {
			return nil, errors.WithStack(err)
		}
		return w, nil
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
		err = errors.NotValid.New(err, "[objcache] NewRedisByURLClient Parameter %q with value %q is invalid", keys[i], params.Get(keys[i]))
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
		err = errors.NotValid.New(err, "[objcache] NewRedisByURLClient Parameter %q with value %q is invalid", keys[i], params.Get(keys[i]))
	}
	return err
}

// NewRedisByURLClient connects to a Redis server at the given URL using the HTTP URI
// scheme. URLs should follow the draft IANA specification for the scheme
// (https://www.iana.org/assignments/uri-schemes/prov/redis). This option
// function sets the connection as cache backend to the Service.
//
// On error, while parsing the rawURL, this function will leak sensitive data,
// for now.
//
// For example:
// 		redis://localhost:6379/?db=3
// 		redis://localhost:6379/?max_active=50&max_idle=5&idle_timeout=10s&max_conn_lifetime=1m&key_prefix=xcache_
func NewRedisByURLClient(rawURL string) NewStorageFn {
	return func() (Storager, error) {
		addr, _, password, params, err := url.ParseConnection(rawURL)
		if err != nil {
			return nil, errors.Wrapf(err, "[objcache] Redis error parsing URL %q", rawURL)
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
			return nil, errors.WithStack(err)
		}

		err = parseInt([]string{
			"max_idle",
			"max_active",
		}, []*int{
			&pool.MaxIdle,
			&pool.MaxActive,
		}, params)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		// if params.Get("tls") == "1" {
		// 	o.TLSConfig = &tls.Config{ServerName: o.Addr} // TODO check if might be wrong the Addr,
		// }

		return NewRedisClient(pool, &RedisOption{
			KeyPrefix: params.Get("key_prefix"),
		})()
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

func makeRedisWrapper(rp *redis.Pool, ro *RedisOption) redisWrapper {
	return redisWrapper{
		Pool: rp,
		// ipf: &sync.Pool{
		// 	New: func() any {
		// 		ifs := make([]any, 0, 10)
		// 		return &ifs
		// 	},
		// },
		keyPrefix: ro.KeyPrefix,
	}
}

type redisWrapper struct {
	*redis.Pool
	ping      bool
	keyPrefix string
	// ifp  *sync.Pool
}

func (w redisWrapper) Set(_ context.Context, keys []string, values [][]byte, expirations []time.Duration) (err error) {
	conn := w.Pool.Get()
	defer func() {
		if err2 := conn.Close(); err == nil && err2 != nil {
			err = err2
		}
	}()

	args := make([]any, 0, len(keys)*3)
	for i, key := range keys {
		e := expirations[i].Seconds() // e = expires in x seconds
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
				err = nil
				val = nil
			} else {
				return nil, errors.Wrapf(err, "[objcache] With keys %v", keys)
			}
		}
		values = append(values, val)
		return
	}

	values, err = redis.ByteSlices(conn.Do("MGET", strSliceToIFaces(nil, keys)...))
	if err != nil {
		err = errors.Wrapf(err, "[objcache] With keys %v", keys)
		return
	}
	if lk, lv := len(keys), len(values); lk != lv {
		err = errors.Mismatch.Newf("[objcache] Length of keys (%d) does not match length of returned bytes (%d). Keys: %v", lk, lv, keys)
	}
	return
}

func strSliceToIFaces(ret []any, sl []string) []any {
	// TODO use a sync.Pool but write before hand appropriate concurrent running benchmarks
	if ret == nil {
		ret = make([]any, 0, len(sl))
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

func (w redisWrapper) Truncate(ctx context.Context) (err error) {
	// TODO flush redis by key prefix
	return nil
}

func (w redisWrapper) Close() error {
	return w.Pool.Close()
}
