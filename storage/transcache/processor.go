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
	"io"

	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/errors"
)

// Cacher defines a custom cache type to be used as underlying storage of the
// Transcacher. Must be safe for concurrent usage. Caches which implement this
// interface can be found in the subpackages tcbigcache, tcboltdb, tcredis ...
type Cacher interface {
	Set(key, value []byte) (err error)
	Get(key []byte) (value []byte, err error)
	// Close closes the underlying cache service.
	Close() error
}

// Transcacher represents the function for storing and retrieving arbitrary Go
// types.
type Transcacher interface {
	// Set sets the type src with a key
	Set(key []byte, src interface{}) error
	// Get looks up the key and parses the raw data into the destination pointer
	// dst. You have to check yourself if the returned error is of type NotFound
	// or of any other source. Every caching type defines its own NotFound
	// error.
	Get(key []byte, dst interface{}) error
}

// Codecer defines the functions needed to create a new Encoder or Decoder
type Codecer interface {
	NewEncoder(io.Writer) Encoder
	NewDecoder(io.Reader) Decoder
}

// Encoder defines how to encode a type represented by variable src into a byte
// slice. Encoders must write their data into an io.Writer defined in option
// function WithEncoder().
type Encoder interface {
	Encode(src interface{}) error
}

// Decoder defines how to decode a byte slice into variable dst. Please see
// option function WithEncoder() for details how to get the byte slice.
type Decoder interface {
	Decode(dst interface{}) error
}

// Processor handles the encoding, decoding and caching
type Processor struct {
	// Cache exported to allow easy debugging and access to raw values.
	Cache Cacher
	Codec Codecer
}

// NewProcessor creates a new type with no default cache instance and no
// encoder. You must set a caching service or it panics please see the sub
// packages tcbigcache, tcbolddb and tcredis. You must also set an encoder,
// which is not optional ;-)
func NewProcessor(opts ...Option) (*Processor, error) {
	p := new(Processor)
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, errors.Wrap(err, "[transcache] NewProcessor applied options")
		}
	}
	return p, nil
}

// Set sets the type src with a key
func (tr *Processor) Set(key []byte, src interface{}) error {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	enc := tr.Codec.NewEncoder(buf)
	if pc, ok := tr.Codec.(*pooledCodec); ok {
		defer pc.PutEncoder(enc)
	}

	if err := enc.Encode(src); err != nil {
		return errors.NewFatal(err, "[transcache] Set.Encode")
	}

	var copied = make([]byte, buf.Len(), buf.Len())
	copy(copied, buf.Bytes()) // copy the encoded data away because we're reusing the buffer
	return errors.NewFatal(tr.Cache.Set(key, copied), "[transcache] Set.Cache.Set")
}

// Get looks up the key and parses the raw data into the destination pointer
// dst. You have to check yourself if the returned error is of type NotFound or
// of any other source. Every caching type defines its own NotFound error.
func (tr *Processor) Get(key []byte, dst interface{}) error {
	val, err := tr.Cache.Get(key)
	if err != nil {
		return errors.Wrap(err, "[transcache] Get.Cache.Get")
	}
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if _, err := buf.Write(val); err != nil {
		return errors.NewWriteFailed(err, "[transcache] Get.Buffer.Write")
	}
	dec := tr.Codec.NewDecoder(buf)
	if pc, ok := tr.Codec.(*pooledCodec); ok {
		defer pc.PutDecoder(dec)
	}
	if err := dec.Decode(dst); err != nil {
		return errors.NewFatal(err, "[transcache] Get.Decode")
	}
	return nil
}
