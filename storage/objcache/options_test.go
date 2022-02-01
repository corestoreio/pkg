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
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/assert"
)

var _ Codecer = (*JSONCodec)(nil)

type JSONCodec struct{}

func (c JSONCodec) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (c JSONCodec) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

func TestWithSimpleSlowCacheMap_Expires(t *testing.T) {
	t.Parallel()

	p, err := NewService[string](NewBlackHoleClient[string](nil), NewCacheSimpleInmemory[string], &ServiceOptions{Codec: JSONCodec{}})
	assert.NoError(t, err)
	defer assert.NoError(t, p.Close())

	t.Run("key not found", func(t *testing.T) {
		err := p.Get(context.TODO(), "upppsss", nil)
		assert.NoError(t, err, "should have kind not found, but got: %+v", err)
	})

	t.Run("key expires", func(t *testing.T) {
		err := p.Set(context.TODO(), "keyEx", 3.14159, time.Second)
		assert.NoError(t, err)
		var f float64
		err = p.Get(context.TODO(), "keyEx", &f)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, 3.14159, f)
		time.Sleep(time.Second * 2)

		var f2 float64
		err = p.Get(context.TODO(), "keyEx", &f2)
		assert.NoError(t, err, "%+v", err)
		assert.Empty(t, f2)
	})
}

var (
	_ marshaler   = (*binary)(nil)
	_ unmarshaler = (*binary)(nil)
)

func TestMakeBinary(t *testing.T) {
	p, err := NewService(NewCacheSimpleInmemory[string], NewCacheSimpleInmemory[string], &ServiceOptions{Codec: JSONCodec{}})
	assert.NoError(t, err)
	defer assert.NoError(t, p.Close())

	t.Run("exists single", func(t *testing.T) {
		b := MakeBinary()
		err := p.Set(context.TODO(), "mb01", b, 0)
		assert.NoError(t, err)

		err = p.Get(context.TODO(), "mb01", &b)
		assert.NoError(t, err)
		assert.True(t, b.IsValid(), "Binary should be valid")
	})

	t.Run("not exists single", func(t *testing.T) {
		b := MakeBinary()
		err = p.Get(context.TODO(), "mb02", &b)
		assert.NoError(t, err)
		assert.False(t, b.IsValid(), "Binary should be valid")
	})

	t.Run("exists multiple", func(t *testing.T) {
		b1 := MakeBinary()
		b2 := MakeBinary()
		b3 := MakeBinary()
		keys := []string{"mb10", "mb20", "mb30"}
		vals := []any{&b1, &b2, &b3}
		err := p.SetMulti(context.TODO(), keys, vals, nil)
		assert.NoError(t, err)

		b1a := MakeBinary()
		b2a := MakeBinary()
		b3a := MakeBinary()
		vals2 := []any{&b1a, &b2a, &b3a}
		err = p.GetMulti(context.TODO(), keys, vals2)
		assert.NoError(t, err)
		assert.True(t, b1a.IsValid(), "Binary b1a should be valid")
		assert.True(t, b2a.IsValid(), "Binary b2a should be valid")
		assert.True(t, b3a.IsValid(), "Binary b3a should be valid")
	})

	t.Run("not exists multiple", func(t *testing.T) {
		b1 := MakeBinary()
		b2 := MakeBinary()
		b3 := MakeBinary()
		keys := []string{"mb10", "mb20a", "mb30a"}
		vals := []any{&b1, &b2, &b3}

		err = p.GetMulti(context.TODO(), keys, vals)
		assert.NoError(t, err)
		assert.True(t, b1.IsValid(), "Binary b1 should be valid")
		assert.False(t, b2.IsValid(), "Binary b2 should not be valid")
		assert.False(t, b3.IsValid(), "Binary b3 should not be valid")
	})
}
