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
	"net"
	"net/url"
	"regexp"
	"strconv"

	"github.com/corestoreio/csfw/storage/typecache"
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
func WithClient(opt *redis.Options, ping ...bool) typecache.Option {
	return func(p *typecache.Processor) error {
		c := redis.NewClient(opt)
		if len(ping) > 0 && ping[0] {
			if _, err := c.Ping().Result(); err != nil {
				return errors.NewFatal(err, "[tcredis] WithClient Ping")
			}
		}
		p.Cache = wrapper{
			Client: c,
		}
		return nil
	}
}

var pathDBRegexp = regexp.MustCompile(`/(\d*)\z`)

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
func WithURL(rawurl string, opt *redis.Options, ping ...bool) typecache.Option {
	return func(p *typecache.Processor) error {

		u, err := url.Parse(rawurl)
		if err != nil {
			return errors.NewFatal(err, "[tcredis] WithDialURL url.Parse")
		}

		if u.Scheme != "redis" {
			return errors.NewNotValidf("[tcredis] Invalid Redis URL scheme: %q", u.Scheme)
		}

		// As per the IANA draft spec, the host defaults to localhost and
		// the port defaults to 6379.
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			// assume port is missing
			host = u.Host
			port = "6379"
		}
		if host == "" {
			host = "localhost"
		}
		address := net.JoinHostPort(host, port)

		var password string
		if u.User != nil {
			password, _ = u.User.Password()
		}

		var db int64
		match := pathDBRegexp.FindStringSubmatch(u.Path)
		if len(match) == 2 {
			if len(match[1]) > 0 {
				db, err = strconv.ParseInt(match[1], 10, 64)
				if err != nil {
					return errors.NewNotValidf("[tcredis] Invalid database: %q in %q", u.Path[1:], match[1])
				}
			}
		} else if u.Path != "" {
			return errors.NewNotValidf("[tcredis] Invalid database: %q", u.Path[1:])
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
	return errors.NewFatal(cmd.Err(), "[tcredis] wrapper.Set.NewStatusCmd")
}

var errKeyNotFound = errors.NewNotFoundf(`[tcredis] Key not found`)

func (w wrapper) Get(key []byte) ([]byte, error) {

	cmd := redis.NewCmd("GET", key)
	w.Client.Process(cmd)

	if cmd.Err() != nil {
		if cmd.Err().Error() != "redis: nil" { // wow that is ugly, how to do better?
			return nil, errors.NewFatal(cmd.Err(), "[tcredis] wrapper.Get.Cmd")
		}
		return nil, errKeyNotFound
	}

	raw, err := conv.ToByteE(cmd.Val())
	return raw, errors.NewFatal(err, "[tcredis] wrapper.Get.conv.ToByte")
}
