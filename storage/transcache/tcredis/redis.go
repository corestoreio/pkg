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
	"github.com/corestoreio/csfw/net/url"
	"github.com/corestoreio/csfw/storage/transcache"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/errors"
	"gopkg.in/redis.v3"
)

// I'm happy to replace the redis client with another as long as the other works
// in concurrent situations without race conditions and have the same benchmark perf.

// WithClient connects to the Redis server. Set ping to true to check if the
// connection works correctly.
//
// For options see: https://godoc.org/gopkg.in/redis.v3#Options
func WithClient(opt *redis.Options, ping ...bool) transcache.Option {
	return func(p *transcache.Processor) error {
		c := redis.NewClient(opt)
		if len(ping) > 0 && ping[0] {
			if _, err := c.Ping().Result(); err != nil {
				return errors.NewFatalf("[tcredis] WithClient Ping: %s", err)
			}
		}
		p.Cache = wrapper{
			Client: c,
		}
		return nil
	}
}

// WithURL connects to a Redis server at the given URL using the Redis
// URI scheme. URLs should follow the draft IANA specification for the
// scheme (https://www.iana.org/assignments/uri-schemes/prov/redis).
// This option function sets the connection as cache backend to the Processor.
//
// For redis.Options see: https://godoc.org/gopkg.in/redis.v3#Options
// They can be nil. If not nil, the rawURL will overwrite network,
// address, password and DB.
//
// For example: redis://localhost:6379/3
func WithURL(rawurl string, opt *redis.Options, ping ...bool) transcache.Option {
	return func(p *transcache.Processor) error {

		address, password, db, err := url.ParseRedis(rawurl)
		if err != nil {
			return errors.Wrap(err, "[tcredis] url.RedisParseURL")
		}

		myOpt := &redis.Options{
			Network:  "tcp",
			Addr:     address,
			Password: password,
			DB:       db,
		}
		if opt != nil {
			opt.Network = myOpt.Network
			opt.Addr = myOpt.Addr
			opt.Password = myOpt.Password
			opt.DB = myOpt.DB
		} else {
			opt = myOpt
		}
		return WithClient(opt, ping...)(p)
	}
}

type wrapper struct {
	*redis.Client
}

func (w wrapper) Set(key []byte, value []byte) error {
	cmd := redis.NewStatusCmd("SET", key, value)
	w.Client.Process(cmd)
	if err := cmd.Err(); err != nil {
		return errors.NewFatalf("[tcredis] wrapper.Set.NewStatusCmd: %s", err)
	}
	return nil
}

var errKeyNotFound = errors.NewNotFoundf(`[tcredis] Key not found`)

func (w wrapper) Get(key []byte) ([]byte, error) {

	cmd := redis.NewCmd("GET", key)
	w.Client.Process(cmd)

	if cmd.Err() != nil {
		if cmd.Err().Error() != "redis: nil" { // wow that is ugly, how to do better?
			return nil, errors.NewFatalf("[tcredis] wrapper.Get.Cmd: %s", cmd.Err())
		}
		return nil, errKeyNotFound
	}

	raw, err := conv.ToByteE(cmd.Val())
	if err != nil {
		return nil, errors.NewFatalf("[tcredis] wrapper.Get.conv.ToByte: %s", err)
	}
	return raw, nil
}
