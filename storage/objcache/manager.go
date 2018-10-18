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
	"context"
	"io"
	"sort"
	"time"

	"github.com/corestoreio/errors"
)

// Storager defines a custom backend cache type to be used as underlying storage
// of the Manager. Must be safe for concurrent usage. Caches which implement
// this interface can be enabled via build tag. The context depends if it is
// supported by a backend cache implementation. All keys and values have the
// same length.
type Storager interface {
	Set(ctx context.Context, key []string, value [][]byte) (err error)
	Get(ctx context.Context, key []string) (value [][]byte, err error)
	Delete(ctx context.Context, key []string) (err error)
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

type Item struct {
	Key string
	// Object is a pointer to the current type.
	Object     interface{}
	Expiration time.Duration // TODO implement
	// More fields will follow
}

func NewItem(key string, object interface{}) *Item {
	return &Item{
		Key:    key,
		Object: object,
	}
}

func (m *Item) encode(c Codecer, buf *bytes.Buffer) error {
	switch ot := m.Object.(type) {
	case marshaler:
		data, err := ot.Marshal()
		if err != nil {
			return errors.Wrapf(err, "[objcache] With key %q and dst type %T", m.Key, m.Object)
		}
		_, err = buf.Write(data)
		return err
	default:
		enc := c.NewEncoder(buf)
		pc, ok := c.(*pooledCodec)
		err := enc.Encode(m.Object)
		if ok {
			pc.PutEncoder(enc)
		}
		if err != nil {
			return errors.Wrapf(err, "[objcache] With key %q and dst type %T", m.Key, m.Object) // saves an allocation ;-)
		}
	}
	return nil
}

func (m *Item) decode(c Codecer, data []byte) error {
	if unm, ok := m.Object.(unmarshaler); ok {
		if err := unm.Unmarshal(data); err != nil {
			return errors.Wrapf(err, "[objcache] With key %q and dst type %T", m.Key, m.Object)
		}
		return nil
	}

	r := bytes.NewReader(data)
	dec := c.NewDecoder(r)
	pc, ok := c.(*pooledCodec)
	err := dec.Decode(m.Object)
	if ok {
		pc.PutDecoder(dec)
	}
	if err != nil && err != io.EOF {
		return errors.Wrapf(err, "[objcache] With key %q and dst type %T", m.Key, m.Object) // saves an allocation ;-)
	}
	return nil
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
func (tr *Manager) Set(ctx context.Context, items ...*Item) error {
	keys := make([]string, 0, len(items))
	vals := make([][]byte, 0, len(items))
	for _, item := range items {
		keys = append(keys, item.Key)
		var buf bytes.Buffer
		if err := item.encode(tr.codec, &buf); err != nil {
			return errors.WithStack(err)
		}
		vals = append(vals, buf.Bytes())
	}
	if err := tr.cache.Set(ctx, keys, vals); err != nil {
		return errors.Wrapf(err, "[objcache] With keys %v", keys)
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
func (tr *Manager) Get(ctx context.Context, items ...*Item) error {
	keys := make([]string, len(items))
	for i, item := range items {
		keys[i] = item.Key
	}

	vals, err := tr.cache.Get(ctx, keys)
	if err != nil {
		return errors.Wrapf(err, "[objcache] With keys %v", keys)
	}

	for i, item := range items {
		val := vals[i]
		if err := item.decode(tr.codec, val); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Delete removes a key from the storage.
func (tr *Manager) Delete(ctx context.Context, key ...string) (err error) {
	if err := tr.cache.Delete(ctx, key); err != nil {
		return errors.Wrapf(err, "[objcache] With keys %v", key)
	}
	return nil
}

// Close closes the underlying storage engines.
func (tr *Manager) Close() error {
	return tr.cache.Close()
}
