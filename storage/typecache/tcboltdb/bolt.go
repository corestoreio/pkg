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

package tcboltdb

import (
	"os"

	"github.com/boltdb/bolt"
	"github.com/corestoreio/csfw/storage/typecache"
	"github.com/corestoreio/csfw/util/errors"
)

// BucketName global bucket name for all entries
var BucketName = []byte("typecache")

var errKeyNotFound = errors.NewNotFoundf(`[tcboltdb] Key not found`)

// WithFile open creates and opens a bolt database at the given path.
// If the file does not exist then it will be created automatically.
// If the third argument Options doesn't get applied bolt.DefaultOptions
// will be used.
// Creates a new bucket from variable name BucketName if that bucket does
// not exists.
func WithFile(path string, mode os.FileMode, options ...*bolt.Options) typecache.Option {
	return func(p *typecache.Processor) error {
		var opt = bolt.DefaultOptions
		if len(options) == 1 {
			opt = options[0]
		}

		db, err := bolt.Open(path, mode, opt)
		if err != nil {
			return errors.NewFatal(err, "[tcboltdb] bolt.Open")
		}
		return WithDB(db)(p)
	}
}

// WithDB uses an existing DB and creates a new bucket from variable name
// BucketName if that bucket does not exists.
func WithDB(db *bolt.DB) typecache.Option {
	return func(p *typecache.Processor) error {

		err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(BucketName)
			return errors.NewFatal(err, "[tcboltdb] bolt.CreateBucketIfNotExists")
		})

		p.Cache = wrapper{db}

		return err
	}
}

type wrapper struct {
	*bolt.DB
}

func (w wrapper) Set(key []byte, value []byte) (err error) {
	err = w.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketName)
		return errors.NewFatal(b.Put(key, value), "[tcboltdb] boltWrapper.Set.Put")
	})
	return errors.Wrap(err, "[tcboltdb] boltWrapper.Set.Update")
}

func (w wrapper) Get(key []byte) ([]byte, error) {
	var found bool
	var buf []byte
	if err := w.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketName)
		v := b.Get(key)
		found = v != nil
		buf = make([]byte, len(v), len(v))
		copy(buf, v)
		return nil
	}); err != nil {
		return nil, errors.NewFatal(err, "[tcboltdb] boltWrapper.Get.View")
	}

	if !found {
		return nil, errKeyNotFound
	}
	return buf, nil
}
