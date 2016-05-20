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
	"github.com/corestoreio/csfw/util/errors"
	"github.com/garyburd/redigo/redis"
)

//BucketName global bucket name for all entries
//var BucketName = []byte("typecache")

var errKeyNotFound = errors.NewNotFoundf(`[tcboltdb] Key not found`)

func WithDial(network, address string, options ...redis.DialOption) typecache.Option {
	return func(p *typecache.Processor) error {

		_, _ = redis.Dial(network, address, options...)

		return nil
	}
}

type wrapper struct {
	redis.Conn
}

func (bw wrapper) Set(key []byte, value []byte) (err error) {

	return errors.Wrap(err, "[tcboltdb] boltWrapper.Set.Update")
}

func (bw wrapper) Get(key []byte) ([]byte, error) {

	return nil, nil
}
