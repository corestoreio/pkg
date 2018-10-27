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
// same length. `expirations` is a list seconds with the same length as keys &
// values. A second entry defines when a key expires. If the entry is empty, the
// key does not expire.
type Storager interface {
	Set(ctx context.Context, keys []string, values [][]byte, expirations []int64) (err error)
	Get(ctx context.Context, keys []string) (values [][]byte, err error)
	Delete(ctx context.Context, keys []string) (err error)
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
	cache             map[int]Storager // read/write randomly from/to a storage
	codec             Codecer
	defaultExpiration time.Duration // TODO implement
	// sync.Pool of *Item
}

type items []*Item

// Item defines a single cache entry with options like cache expiration.
type Item struct {
	Key string
	// Object is a pointer to the current type.
	Object interface{}
	// Expiration in seconds. If 0 no expiration desired.
	Expiration int64

	// ClearObjectAfterSet bool idea to set the object field to nil after Set has been called.
}

// NewItem creates a new item pointer.
func NewItem(key string, object interface{}) *Item {
	return &Item{
		Key:    key,
		Object: object,
	}
}

func (m *Item) encode(c Codecer, w io.Writer) (err error) {
	switch ot := m.Object.(type) {
	case marshaler:
		data, err2 := ot.Marshal()
		if err2 != nil {
			return errors.Wrapf(err2, "[objcache] With key %q and dst type %T", m.Key, m.Object)
		}
		_, err = w.Write(data)
	default:
		enc := c.NewEncoder(w)
		pc, ok := c.(*pooledCodec)
		err = enc.Encode(m.Object)
		if ok {
			pc.PutEncoder(enc)
		}
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			err = errors.Wrapf(err, "[objcache] With key %q and dst type %T", m.Key, m.Object) // saves an allocation ;-)
		}
	}
	return err
}

func (m *Item) decode(c Codecer, data []byte) (err error) {
	switch ot := m.Object.(type) {
	case unmarshaler:
		if err = ot.Unmarshal(data); err != nil {
			return errors.Wrapf(err, "[objcache] With key %q and dst type %T", m.Key, m.Object)
		}
	default:
		r := bytes.NewReader(data)
		dec := c.NewDecoder(r)
		pc, ok := c.(*pooledCodec)
		err = dec.Decode(m.Object)
		if ok {
			pc.PutDecoder(dec)
		}
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			err = errors.Wrapf(err, "[objcache] With key %q and dst type %T", m.Key, m.Object) // saves an allocation ;-)
		}
	}
	return err
}

// Encode encodes all items to their byte slice representation. Returns two
// slices whose indexes match to the other. The data might be appended to the
// optional arguments `keys` and `values`.
func (ms items) encode(codec Codecer) (keys []string, values [][]byte, expires []int64, err error) {
	keys = make([]string, 0, len(ms))
	values = make([][]byte, 0, len(ms))
	expires = make([]int64, 0, len(ms))

	for _, item := range ms {
		keys = append(keys, item.Key)
		var buf bytes.Buffer
		if err := item.encode(codec, &buf); err != nil {
			return nil, nil, nil, errors.WithStack(err)
		}
		values = append(values, buf.Bytes())
		expires = append(expires, item.Expiration)
	}
	return keys, values, expires, nil
}

func (ms items) decode(codec Codecer, values [][]byte) error {
	for i, item := range ms {
		if err := item.decode(codec, values[i]); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (ms items) keys(keys ...string) []string {
	if keys == nil {
		keys = make([]string, 0, len(ms))
	}
	for _, item := range ms {
		keys = append(keys, item.Key)
	}
	return keys
}

// NewService creates a new cache service.
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

// Set sets the item in the cache. `items` gets either encoded using the
// previously applied encoder OR each item of the `items` slice gets checked if
// it implements interface
//		type marshaler interface {
//			Marshal() ([]byte, error)
//		}
// and calls `Marshal`. Checking for marshaler has precedence. Useful with
// protobuf. The argument `items` allows a cache backend to write multiple items
// to its cache with one connection. E.g. Redis MGET/MSET.
func (tr *Service) Set(ctx context.Context, item ...*Item) error {
	itms := items(item)
	keys, values, expires, err := itms.encode(tr.codec)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, c := range tr.cache {
		if err := c.Set(ctx, keys, values, expires); err != nil {
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

// Get looks up the items and parses the raw data into the destination pointer
// `dst`. If `dst` implements interface
//		type unmarshaler interface {
// 			Unmarshal([]byte) error
//		}
// the Unmarshal gets called. This type check has precedence before the decoder.
// You have to check yourself if the returned error is of type NotFound or of
// any other source. Every caching type defines its own NotFound error. The
// argument `items` allows a cache backend to read multiple items at once, e.g.
// Redis MGET/MSET.
func (tr *Service) Get(ctx context.Context, item ...*Item) error {
	itms := items(item)
	keys := itms.keys()
	for _, c := range tr.cache {
		vals, err := c.Get(ctx, keys)
		if err != nil {
			return errors.Wrapf(err, "[objcache] With keys %v", keys)
		}
		if err := itms.decode(tr.codec, vals); err != nil {
			return errors.WithStack(err)
		}
		return nil // TODO implement a better cache selection algorithm instead of random map reads
	}
	return nil
}

// Delete removes keys from the storage.
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
