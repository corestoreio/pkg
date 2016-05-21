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
	"github.com/corestoreio/csfw/storage/typecache"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/garyburd/redigo/redis"
)

var errKeyNotFound = errors.NewNotFoundf(`[tcredis] Key not found`)

// WithDial connects to the Redis server at the given network and
// address using the specified options. Sets the connection as
// cache backend to the Processor.
func WithDial(network, address string, options ...redis.DialOption) typecache.Option {
	return func(p *typecache.Processor) error {
		con, err := redis.Dial(network, address, options...)
		if err != nil {
			return errors.NewFatal(err, "[tcredis] WithDial.redis.Dial")
		}
		return WithCon(con)(p)
	}
}

// WithDialURL connects to a Redis server at the given URL using the Redis
// URI scheme. URLs should follow the draft IANA specification for the
// scheme (https://www.iana.org/assignments/uri-schemes/prov/redis).
// Sets the connection as cache backend to the Processor.
func WithDialURL(rawurl string, options ...redis.DialOption) typecache.Option {
	return func(p *typecache.Processor) error {
		con, err := redis.DialURL(rawurl, options...)
		if err != nil {
			return errors.NewFatal(err, "[tcredis] WithDial.redis.DialURL")
		}
		return WithCon(con)(p)
	}
}

// WithCon sets a connection to a Redis server as cache backend.
func WithCon(con redis.Conn) typecache.Option {
	return func(p *typecache.Processor) error {
		p.Cache = wrapper{con}
		return nil
	}
}

type wrapper struct {
	redis.Conn
}

func (bw wrapper) Set(key []byte, value []byte) (err error) {
	_, err = bw.Do("SET", key, value)
	return errors.Wrap(err, "[tcredis] wrapper.Set.Do")
}

func (bw wrapper) Get(key []byte) ([]byte, error) {
	raw, err := bw.Do("GET", key)
	if raw == nil && err == nil {
		return nil, errKeyNotFound
	}
	if err != nil {
		return nil, errors.NewFatal(err, "[tcredis] wrapper.Get.Do2")
	}
	resp, err := conv.ToByteE(raw)
	return resp, errors.NewFatal(err, "[tcredis] wrapper.Get.conv.ToByte")
}
