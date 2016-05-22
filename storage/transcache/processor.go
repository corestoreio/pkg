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

package transcache

import (
	"bytes"
	"sync"

	"github.com/corestoreio/csfw/util/errors"
)

// Cacher defines a custom cache type to be used as underlying storage.
// Must be safe for parallel usage.
type Cacher interface {
	Set(key, value []byte) (err error)
	Get(key []byte) (value []byte, err error)
	Close() error
}

// Encoder defines how to encode a type represented by variable src into
// a byte slice. Encoders must write their data into an io.Writer defined
// in option function WithEncoder().
type Encoder interface {
	Encode(src interface{}) error
}

// Decoder defines how to decode a byte slice into variable dst. Please see
// option function WithEncoder() for details how to get the byte slice.
type Decoder interface {
	Decode(dst interface{}) error
}

// 64 is quite good. there are not yet any benefits from higher values
const encodeShards = 64 // must be power of 2
const encodeShardMask uint64 = encodeShards - 1

type encode struct {
	Encoder
	sync.Mutex
	buf *bytes.Buffer
}

type decode struct {
	Decoder
	sync.Mutex
	buf *bytes.Buffer
}

// Processor handles the encoding, decoding and caching
type Processor struct {
	Hasher
	// Cache exported to allow easy debugging and access to raw values.
	Cache Cacher
	enc   [encodeShards]encode
	dec   [encodeShards]decode
}

// NewProcessor creates a new type with no default cache instance
// and encoding/gob as the underlying encoder. If you use gob please make sure
// to use gob.Register() to register your types.
// You must set a caching service or it panics please see the sub packages
// tcbigcache and tcbolddb.
func NewProcessor(opts ...Option) (*Processor, error) {
	p := &Processor{
		Hasher: newDefaultHasher(),
	}

	for i := 0; i < encodeShards; i++ {
		p.enc[i].buf = new(bytes.Buffer)
		p.dec[i].buf = new(bytes.Buffer)
	}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, errors.Wrap(err, "[transcache] NewProcessor applied options")
		}
	}

	if p.enc[0].Encoder == nil || p.dec[0].Decoder == nil {
		if err := withGob()(p); err != nil {
			return nil, errors.Wrap(err, "[transcache] NewProcessor.Option.WithGob")
		}
	}
	return p, nil
}

func (tr *Processor) shardID(key []byte) uint64 {
	return tr.Hasher.Sum64(key) & encodeShardMask
}

// Set sets the type src with a key
func (tr *Processor) Set(key []byte, src interface{}) error {
	shardID := tr.shardID(key)

	tr.enc[shardID].Lock()
	defer tr.enc[shardID].Unlock()
	if err := tr.enc[shardID].Encode(src); err != nil {
		return errors.NewFatal(err, "[transcache] Set.Encode")
	}

	var buf = make([]byte, tr.enc[shardID].buf.Len(), tr.enc[shardID].buf.Len())
	copy(buf, tr.enc[shardID].buf.Bytes()) // copy the encoded data away because we're reusing the buffer
	tr.enc[shardID].buf.Reset()
	return errors.NewFatal(tr.Cache.Set(key, buf), "[transcache] Set.Cache.Set")
}

// Get looksup the key and parses the raw data into the destination pointer dst.
// You have to check yourself if the returned error is of type NotFound or of
// any other source. Every caching type defines its own NotFound error.
func (tr *Processor) Get(key []byte, dst interface{}) error {
	shardID := tr.shardID(key)
	tr.dec[shardID].Lock()
	defer tr.dec[shardID].Unlock()
	tr.dec[shardID].buf.Reset()

	val, err := tr.Cache.Get(key)
	if err != nil {
		return errors.Wrap(err, "[transcache] Get.Cache.Get")
	}
	if _, err := tr.dec[shardID].buf.Write(val); err != nil {
		return errors.NewWriteFailed(err, "[transcache] Get.Buffer.Write")
	}
	if err := tr.dec[shardID].Decode(dst); err != nil {
		return errors.NewFatal(err, "[transcache] Get.Decode")
	}
	return nil
}
