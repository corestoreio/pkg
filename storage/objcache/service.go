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
	"encoding"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// Storager defines a custom backend cache type to be used as underlying storage
// of the Service. Must be safe for concurrent usage. Caches which implement
// this interface can be enabled via build tag. The context depends if it is
// supported by a backend cache implementation. All keys and values have the
// same length. `expirations` is a list seconds with the same length as keys &
// values. A second entry defines when a key expires. If the entry is empty, the
// key does not expire.
type Storager[K comparable] interface {
	Set(ctx context.Context, keys []K, values [][]byte, expirations []time.Duration) error
	// Get returns the bytes for given keys. The values slice must have the same
	// length as the keys slice. If one of the keys can't be found, its byte
	// slice must be `nil`.
	Get(ctx context.Context, keys []K) (values [][]byte, err error)
	Delete(ctx context.Context, keys []K) error
	Truncate(ctx context.Context) error
	Close() error
}

type NewStorageFn[K comparable] func() (Storager[K], error)

// Codecer defines the functions needed to create a new Encoder or Decoder
type Codecer interface {
	NewEncoder(io.Writer) Encoder
	NewDecoder(io.Reader) Decoder
}

// Encoder defines how to Encode a type represented by variable src into a byte
// slice. Encoders must write their data into an io.Writer defined in option
// function WithEncoder().
type Encoder interface {
	Encode(src any) error
}

// Decoder defines how to decode a byte slice into variable dst. Please see
// option function WithEncoder() for details how to get the byte slice.
type Decoder interface {
	Decode(dst any) error
}

func writeMarshal(buf *bytes.Buffer, m func() ([]byte, error)) error {
	data, err := m()
	if err != nil {
		return err
	}
	_, err = buf.Write(data)
	return err
}

func encodeOne[K comparable](c Codecer, buf *bytes.Buffer, key K, src any) error {
	switch ot := src.(type) {
	case marshaler:
		if err := writeMarshal(buf, ot.Marshal); err != nil {
			return fmt.Errorf("[objcache] 1643662002029 With key %v and dst type %T: %w", key, src, err)
		}
	case encoding.TextMarshaler:
		if err := writeMarshal(buf, ot.MarshalText); err != nil {
			return fmt.Errorf("[objcache] 1643662168139 With key %v and dst type %T: %w", key, src, err)
		}
	case encoding.BinaryMarshaler:
		if err := writeMarshal(buf, ot.MarshalBinary); err != nil {
			return fmt.Errorf("[objcache] 1643662258793 With key %v and dst type %T: %w", key, src, err)
		}
	default:
		if c == nil {
			return fmt.Errorf("[objcache] Src type %T does not implement Marshal or Codec not set", src)
		}

		enc := c.NewEncoder(buf)
		pc, ok := c.(*pooledCodec)
		err := enc.Encode(src)
		if ok {
			pc.PutEncoder(enc)
		}
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			return fmt.Errorf("[objcache] 1643662314576 With key %v and dst type %T", key, src) // saves an allocation ;-)
		}
	}
	return nil
}

func decodeOne[K comparable](c Codecer, data []byte, key K, dst any) (err error) {
	switch ot := dst.(type) {
	case unmarshaler:
		if err = ot.Unmarshal(data); err != nil {
			return fmt.Errorf("[objcache] 1643662475854 With key %v and dst type %T: %w", key, dst, err)
		}
	case encoding.TextUnmarshaler:
		if err = ot.UnmarshalText(data); err != nil {
			return fmt.Errorf("[objcache] 1643662482931 With key %v and dst type %T: %w", key, dst, err)
		}
	case encoding.BinaryUnmarshaler:
		if err = ot.UnmarshalBinary(data); err != nil {
			return fmt.Errorf("[objcache] 1643662489237 With key %v and dst type %T: %w", key, dst, err)
		}
	default:
		if c == nil {
			return errors.NotImplemented.Newf("[objcache] Dst type %T does not implement Unmarshal or Codec not set.", dst)
		}

		r := bufferpool.GetReader(data)
		defer bufferpool.PutReader(r)
		dec := c.NewDecoder(r)
		pc, ok := c.(*pooledCodec)
		err = dec.Decode(dst)
		if ok {
			pc.PutDecoder(dec)
		}
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			err = fmt.Errorf("[objcache] With key %v and dst type %T: %w", key, dst, err) // saves an allocation ;-)
		}
	}
	return err
}

// Encode encodes all items to their byte slice representation. Returns two
// slices whose indexes match to the other. The data might be appended to the
// optional arguments `keys` and `values`.
func encodeAll[K comparable](codec Codecer, ri *rawItems[K], defaultExpire time.Duration, keys []K, src []any, expires []time.Duration) (_ *rawItems[K], err error) {
	lenExpires := len(expires)
	for i, key := range keys {
		ri.keys = append(ri.keys, key)
		var buf bytes.Buffer // TODO a buffer pool can be used because of the append
		if err := encodeOne(codec, &buf, key, src[i]); err != nil {
			return nil, errors.WithStack(err)
		}
		ri.values = append(ri.values, buf.Bytes())

		e := defaultExpire
		if lenExpires > 0 && expires[i] != 0 {
			e = expires[i]
		}
		ri.expires = append(ri.expires, e)
	}
	return ri, nil
}

func decodeAll[K comparable](codec Codecer, values [][]byte, keys []K, dst []any) error {
	for i, key := range keys {
		if err := decodeOne(codec, values[i], key, dst[i]); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

type rawItems[K comparable] struct {
	keys    []K
	values  [][]byte
	expires []time.Duration
}

// Service handles the encoding, decoding and caching.
type Service[K comparable] struct {
	so                ServiceOptions
	level1            Storager[K]
	level2            Storager[K]
	defaultExpiration time.Duration // in seconds
	rawItemsPool      sync.Pool
}

func (tr *Service[K]) poolGetRawItems() *rawItems[K] {
	return tr.rawItemsPool.Get().(*rawItems[K])
}

func (tr *Service[K]) poolPutRawItems(ri *rawItems[K]) {
	ri.keys = ri.keys[:0]
	ri.expires = ri.expires[:0]
	for i := range ri.values {
		ri.values[i] = ri.values[i][:0]
	}
	ri.values = ri.values[:0]
	tr.rawItemsPool.Put(ri)
}

// NewService creates a new cache service. Arguments level1 and level2 define
// the cache level. For example level1 should be an LRU or another in-memory
// cache while level2 should be accessed via network. Only level2 is requred
// while level1 can be nil.
func NewService[K comparable](level1, level2 NewStorageFn[K], so *ServiceOptions) (_ *Service[K], err error) {
	s := &Service[K]{
		rawItemsPool: sync.Pool{
			// values might lead to bugs, theoretically, but never experienced them.
			New: func() any {
				return &rawItems[K]{
					keys:    make([]K, 0, 3),
					values:  make([][]byte, 0, 3),
					expires: make([]time.Duration, 0, 3),
				}
			},
		},
	}
	if so != nil {
		s.so = *so
	}
	if level1 != nil {
		if s.level1, err = level1(); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if s.level2, err = level2(); err != nil {
		return nil, errors.WithStack(err)
	}

	if so != nil && len(so.PrimeObjects) > 0 {
		so.Codec = newPooledCodec(so.Codec, so.PrimeObjects...)
		so.PrimeObjects = nil
	}

	return s, nil
}

// marshaler is the interface representing objects that can marshal themselves.
type marshaler interface {
	Marshal() ([]byte, error)
}

// Set puts the item in the cache. `src` gets either encoded using the
// previously applied encoder OR `src` gets checked if it implements interface
//		type marshaler interface {
//			Marshal() ([]byte, error)
//		}
// and calls `Marshal`. (also checks for the interfaces in package "encoding").
// Checking for marshaler has precedence. Useful with protobuf.
func (tr *Service[K]) Set(ctx context.Context, key K, src any, expires time.Duration) error {
	ri := tr.poolGetRawItems()
	defer tr.poolPutRawItems(ri)

	if expires == 0 {
		expires = tr.defaultExpiration
	}

	var buf bytes.Buffer
	if err := encodeOne(tr.so.Codec, &buf, key, src); err != nil {
		return errors.WithStack(err)
	}
	ri.keys = append(ri.keys, key)
	ri.values = append(ri.values, buf.Bytes())
	ri.expires = append(ri.expires, expires)

	if tr.level1 != nil {
		if err := tr.level1.Set(ctx, ri.keys, ri.values, ri.expires); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := tr.level2.Set(ctx, ri.keys, ri.values, ri.expires); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// SetMulti allows a cache to write several entities at once. For example using
// Redis MSET. Same logic applies as when using `Set`.
func (tr *Service[K]) SetMulti(ctx context.Context, keys []K, src []any, expires []time.Duration) error {
	if lk, ld := len(keys), len(src); lk != ld {
		return errors.Mismatch.Newf("[objcache] Length of keys (%d) vs length of src (%d) must be equal", lk, ld)
	}

	ri := tr.poolGetRawItems()
	defer tr.poolPutRawItems(ri)

	ri, err := encodeAll(tr.so.Codec, ri, tr.defaultExpiration, keys, src, expires)
	if err != nil {
		return errors.WithStack(err)
	}

	if tr.level1 != nil {
		if err := tr.level1.Set(ctx, ri.keys, ri.values, ri.expires); err != nil {
			return errors.WithStack(err)
		}
	}
	if err := tr.level2.Set(ctx, ri.keys, ri.values, ri.expires); err != nil {
		return errors.WithStack(err)
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
// the Unmarshal gets called (or one of the interface from package `encoding`).
// This type check has precedence before the decoder. You have to check yourself
// if the returned error is of type NotFound or of any other source. Every
// caching type defines its own NotFound error. If dst has no pointer property,
// no error gets returned, instead the passed value stays empty.
func (tr *Service[K]) Get(ctx context.Context, key K, dst any) (err error) {
	// If dst is not pointer ... unlucky you, we don't do checks with reflect.
	// Instead write better tests.
	ri := tr.poolGetRawItems()
	defer tr.poolPutRawItems(ri)

	ri.keys = append(ri.keys, key)

	var vals [][]byte
	if tr.level1 != nil {
		vals, err = tr.level1.Get(ctx, ri.keys)
		if err != nil {
			return errors.Wrapf(err, "[objcache] Level1 with keys %v", ri.keys)
		}
	}
	if lv := len(vals); lv == 0 {
		vals, err = tr.level2.Get(ctx, ri.keys)
		if err != nil {
			return errors.Wrapf(err, "[objcache] Level2 with keys %v", ri.keys)
		}
	}
	if err == nil {
		idst := [1]any{dst}
		if err2 := decodeAll(tr.so.Codec, vals, ri.keys, idst[:]); err2 != nil {
			return errors.WithStack(err2)
		}
	}
	return err
}

// GetMulti allows a cache backend to retrieve several values at once. Same
// decoding logic applies as when calling `Get`.
func (tr *Service[K]) GetMulti(ctx context.Context, keys []K, dst []any) (err error) {
	if lk, ld := len(keys), len(dst); lk != ld {
		return errors.Mismatch.Newf("[objcache] Length of keys (%d) vs length of dst (%d) must be equal", lk, ld)
	}
	var vals [][]byte
	if tr.level1 != nil {
		vals, err = tr.level1.Get(ctx, keys)
		if err != nil && !errors.NotFound.Match(err) {
			return errors.Wrapf(err, "[objcache] Level1 with keys %v", keys)
		}
	}
	if lv := len(vals); lv == 0 {
		vals, err = tr.level2.Get(ctx, keys)
		if err != nil && !errors.NotFound.Match(err) {
			return errors.Wrapf(err, "[objcache] Level2 with keys %v", keys)
		}
	}
	if err != nil && errors.NotFound.Match(err) {
		return errors.WithStack(err)
	}

	if err := decodeAll(tr.so.Codec, vals, keys, dst); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Truncate truncates all caches.
func (tr *Service[K]) Truncate(ctx context.Context) (err error) {
	if tr.level1 != nil {
		if err := tr.level1.Truncate(ctx); err != nil {
			return errors.WithStack(err)
		}
	}
	if err := tr.level2.Truncate(ctx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Delete removes keys from the storage.
func (tr *Service[K]) Delete(ctx context.Context, key ...K) error {
	if tr.level1 != nil {
		if err := tr.level1.Delete(ctx, key); err != nil {
			return errors.WithStack(err)
		}
	}
	if err := tr.level2.Delete(ctx, key); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Close closes the underlying storage engines.
func (tr *Service[K]) Close() error {
	if tr.level1 != nil {
		if err := tr.level1.Close(); err != nil {
			return errors.WithStack(err)
		}
	}
	return tr.level2.Close()
}
