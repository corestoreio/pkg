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
// of the Service. Must be safe for concurrent usage. Caches which implement
// this interface can be enabled via build tag. The context depends if it is
// supported by a backend cache implementation. All keys and values have the
// same length.
type Storager interface {
	Set(ctx context.Context, items *Items) (err error)
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

// Encoder defines how to Encode a type represented by variable src into a byte
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

// Service handles the encoding, decoding and caching.
type Service struct {
	// Cache exported to allow easy debugging and access to raw values.
	cache map[int]Storager // read/write randomly from/to a storage
	codec Codecer
}

type Items struct {
	items []*Item
	codec Codecer
}

func newItems(i []*Item, c Codecer) *Items {
	// sync pool would be a nice fit here
	return &Items{
		items: i,
		codec: c,
	}
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

func (m *Item) encode(c Codecer, w io.Writer) error {
	switch ot := m.Object.(type) {
	case marshaler:
		data, err := ot.Marshal()
		if err != nil {
			return errors.Wrapf(err, "[objcache] With key %q and dst type %T", m.Key, m.Object)
		}
		_, err = w.Write(data)
		return err
	default:
		enc := c.NewEncoder(w)
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

// Encode encodes all items to they byte slice representation. Returns two slices whose indexes
// match to the other. The data might be appended to the optional arguments `keys` and `values`
func (ms *Items) Encode(keys []string, values [][]byte) (_keys []string, _values [][]byte, err error) {
	if keys == nil {
		keys = make([]string, 0, len(ms.items))
	}
	if values == nil {
		values = make([][]byte, 0, len(ms.items))
	}
	for _, item := range ms.items {
		keys = append(keys, item.Key)
		var buf bytes.Buffer
		if err := item.encode(ms.codec, &buf); err != nil {
			return nil, nil, errors.WithStack(err)
		}
		values = append(values, buf.Bytes())
	}
	return keys, values, nil
}

func (ms *Items) decode(values [][]byte) error {
	for i, item := range ms.items {
		if err := item.decode(ms.codec, values[i]); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Keys returns only all available keys. Might append it to argument `keys`.
func (ms *Items) Keys(keys ...string) []string {
	if keys == nil {
		keys = make([]string, 0, len(ms.items))
	}
	for _, item := range ms.items {
		keys = append(keys, item.Key)
	}
	return keys
}

// NewService creates a new type with no default cache instance and no
// encoder. You must set a caching service or it panics please see the sub
// packages objcache, tcbolddb and objcache. You must also set an encoder,
// which is not optional ;-)
func NewService(opts ...Option) (*Service, error) {
	p := &Service{
		cache: make(map[int]Storager, 6),
	}
	opts2 := options(opts)
	sort.Stable(opts2)
	for _, opt := range opts2 {
		if err := opt.fn(p); err != nil {
			return nil, errors.Wrap(err, "[objcache] NewService applied options")
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
func (tr *Service) Set(ctx context.Context, items ...*Item) error {
	itms := newItems(items, tr.codec)
	for _, c := range tr.cache {

		if err := c.Set(ctx, itms); err != nil {
			return errors.WithStack(err)
		}
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
func (tr *Service) Get(ctx context.Context, items ...*Item) error {
	itms := newItems(items, tr.codec)
	keys := itms.Keys()
	for _, c := range tr.cache {
		vals, err := c.Get(ctx, keys)
		if err != nil {
			return errors.Wrapf(err, "[objcache] With keys %v", keys)
		}
		if err := itms.decode(vals); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Delete removes a key from the storage.
func (tr *Service) Delete(ctx context.Context, key ...string) (err error) {
	for _, c := range tr.cache {
		if err := c.Delete(ctx, key); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Close closes the underlying storage engines.
func (tr *Service) Close() error {
	for _, c := range tr.cache {
		if err := c.Close(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
