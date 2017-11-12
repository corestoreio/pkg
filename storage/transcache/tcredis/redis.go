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

package tcredis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/corestoreio/pkg/net/url"
	"github.com/corestoreio/pkg/storage/transcache"
	"github.com/corestoreio/errors"
	"github.com/garyburd/redigo/redis"
)

// WithClient connects to the Redis server. Set ping to true to check if the
// connection works correctly.
//
// For options see: https://godoc.org/gopkg.in/redis.v3#Options
func WithClient(pool *redis.Pool) transcache.Option {
	return func(p *transcache.Processor) error {
		var ping bool
		if p2, ok := p.Cache.(wrapper); ok {
			ping = p2.ping
		}
		w := wrapper{
			Pool: pool,
			ping: ping,
		}
		p.Cache = w
		return doPing(w)
	}
}

// WithURL connects to a Redis server at the given URL using the Redis URI
// scheme. URLs should follow the draft IANA specification for the scheme
// (https://www.iana.org/assignments/uri-schemes/prov/redis). This option
// function sets the connection as cache backend to the Processor.
//
// On error, while parsing the rawURL, this function will leak sensitive data,
// for now.
//
// For example:
// 		redis://localhost:6379/3
// 		redis://localhost:6379/?max_active=50&max_idle=5&idle_timeout=10s
func WithURL(rawURL string) transcache.Option {
	return func(p *transcache.Processor) error {

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
			return errors.NewNotValid(err, "[backend] NewRedis.ParseNoSQLURL. Parameter max_idle not valid in %q", rawURL)
		}
		idleTimeout, err := time.ParseDuration(params.Get("idle_timeout"))
		if err != nil {
			return errors.NewNotValid(err, "[backend] NewRedis.ParseNoSQLURL. Parameter idle_timeout not valid in %q", rawURL)
		}

		pool := &redis.Pool{
			MaxActive:   maxActive,
			MaxIdle:     maxIdle,
			IdleTimeout: idleTimeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", addr)
				if err != nil {
					return nil, errors.NewFatal(err, "[backend] Redis Dial failed")
				}
				if pw != "" {
					if _, err := c.Do("AUTH", pw); err != nil {
						c.Close()
						return nil, errors.NewUnauthorized(err, "[backend] Redis AUTH failed")
					}
				}
				if _, err := c.Do("SELECT", params.Get("db")); err != nil {
					c.Close()
					return nil, errors.NewFatal(err, "[backend] Redis DB select failed")
				}
				return c, nil
			},
		}

		return WithClient(pool)(p)
	}
}

// WithPing pings the redis database and checks if the connection parameters are
// valid.
func WithPing() transcache.Option {
	return func(p *transcache.Processor) error {
		var pool *redis.Pool
		if w, ok := p.Cache.(wrapper); ok && w.Pool != nil {
			pool = w.Pool
		}
		w := wrapper{
			Pool: pool,
			ping: true,
		}
		p.Cache = w
		return doPing(w)
	}
}

func doPing(w wrapper) error {
	if !w.ping || w.Pool == nil {
		return nil
	}
	conn := w.Pool.Get()
	defer conn.Close()

	pong, err := redis.String(conn.Do("PING"))
	if err != nil && err != redis.ErrNil {
		return errors.NewFatalf("[backend] Redis Ping failed: %s", err)
	}
	if pong != "PONG" {
		return errors.NewFatalf("[backend] Redis Ping not Pong: %#v", pong)
	}
	return nil
}

type wrapper struct {
	*redis.Pool
	ping bool
}

func (w wrapper) Set(key []byte, value []byte) error {
	conn := w.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SET", key, value); err != nil {
		return errors.NewFatalf("[tcredis] wrapper.Set.NewStatusCmd: %s", err)
	}
	return nil
}

func (w wrapper) Get(key []byte) ([]byte, error) {
	conn := w.Pool.Get()
	defer conn.Close()

	raw, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err != redis.ErrNil {
			return nil, errors.NewFatalf("[tcredis] wrapper.Get.Cmd: %s", err)
		}
		return nil, keyNotFound{key: key}
	}
	return raw, nil
}

type keyNotFound struct {
	key []byte
}

func (k keyNotFound) Error() string {
	return fmt.Sprintf("[tcredis] The key %q has not been found.", k.key)
}

func (k keyNotFound) NotFound() bool {
	return true
}
