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

package objcache

import (
	"bytes"
	"io"
	"sort"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// Storager defines a custom backend cache type to be used as underlying storage of the
// Transcacher. Must be safe for concurrent usage. Caches which implement this
// interface can be found in the subpackages objcache, tcboltdb, objcache ...
type Storager interface {
	Set(key, value []byte) (err error)
	Get(key []byte) (value []byte, err error)
	// Close closes the underlying cache service.
	Close() error
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

// Manager handles the encoding, decoding and caching
type Manager struct {
	// Cache exported to allow easy debugging and access to raw values.
	cache Storager
	codec Codecer
}

// NewManager creates a new type with no default cache instance and no
// encoder. You must set a caching service or it panics please see the sub
// packages objcache, tcbolddb and objcache. You must also set an encoder,
// which is not optional ;-)
func NewManager(opts ...Option) (*Manager, error) {
	p := new(Manager)
	opts2 := options(opts)
	sort.Stable(opts2)
	for _, opt := range opts2 {
		if err := opt.fn(p); err != nil {
			return nil, errors.Wrap(err, "[objcache] NewManager applied options")
		}
	}
	return p, nil
}

// marshaler is the interface representing objects that can marshal themselves.
type marshaler interface {
	Marshal() ([]byte, error)
}

// Set sets the type src with a key. Src gets either encoded using the
// previously applied encoder OR `src` gets checked if it implements interface
//		type marshaler interface {
//			Marshal() ([]byte, error)
//		}
// and calls `Marshal`. Checking for marshaler has precedence. Useful with
// protobuf.
func (tr *Manager) Set(key []byte, src interface{}) error {
	if om, ok := src.(marshaler); ok {
		data, err := om.Marshal()
		if err != nil {
			return errors.Wrapf(err, "[objcache] With key %q", string(key))
		}
		if err := tr.cache.Set(key, data); err != nil {
			return errors.Wrapf(err, "[objcache] With key %q", string(key))
		}
		return nil
	}

	var buf bytes.Buffer
	enc := tr.codec.NewEncoder(&buf)
	if pc, ok := tr.codec.(*pooledCodec); ok {
		defer pc.PutEncoder(enc)
	}

	if err := enc.Encode(src); err != nil {
		return errors.Wrapf(err, "[objcache] With key %q", string(key))
	}

	if err := tr.cache.Set(key, buf.Bytes()); err != nil {
		return errors.Wrapf(err, "[objcache] With key %q", string(key))
	}
	return nil
}

// unmarshaler is the interface representing objects that can
// unmarshal themselves.  The argument points to data that may be
// overwritten, so implementations should not keep references to the
// buffer.
type unmarshaler interface {
	Unmarshal([]byte) error
}

// Get looks up the key and parses the raw data into the destination pointer
// `dst`. If `dst` implements interface
//		type unmarshaler interface {
// 			Unmarshal([]byte) error
//		}
// the Unmarshal gets called. This type check has precedence before the decoder.
// You have to check yourself if the returned error is of type NotFound or of
// any other source. Every caching type defines its own NotFound error.
func (tr *Manager) Get(key []byte, dst interface{}) error {
	val, err := tr.cache.Get(key)
	if err != nil {
		return errors.Wrapf(err, "[objcache] With key %q and dst type %T", string(key), dst)
	}

	if unm, ok := dst.(unmarshaler); ok {
		if err := unm.Unmarshal(val); err != nil {
			return errors.Wrapf(err, "[objcache] With key %q and dst type %T", string(key), dst)
		}
		return nil
	}

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if _, err := buf.Write(val); err != nil {
		return errors.WriteFailed.New(err, "[objcache] Get.Buffer.Write for key %q", string(key))
	}
	dec := tr.codec.NewDecoder(buf)
	if pc, ok := tr.codec.(*pooledCodec); ok {
		defer pc.PutDecoder(dec)
	}
	if err := dec.Decode(dst); err != nil {
		return errors.Wrapf(err, "[objcache] With key %q and dst type %T", string(key), dst)
	}
	return nil
}

// Close closes the underlying storage engines.
func (tr *Manager) Close() error {
	return tr.cache.Close()
}
